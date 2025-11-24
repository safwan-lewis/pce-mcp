/*
Copyright 2025 Pextra Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	https://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/url"
)

type GetNodeByIdArg struct {
	NodeId string
}
type GetNodeByIdResponse struct {
	Node      NodeDetail     `json:"node"`
	Instances []InstanceList `json:"instances"`
	Hostname  string         `json:"hostname"`
	Os        struct {
		Kernel       string `json:"kernel"`
		Architecture string `json:"architecture"`
		Uefi         bool   `json:"uefi"`
	} `json:"os"`
	Time struct {
		Time     int64  `json:"time"`
		Timezone string `json:"timezone"`
		Uptime   int64  `json:"uptime"`
	} `json:"time"`
	PceVersion string `json:"pce_version"`
}

func GetNodeById(ctx context.Context, c *Client, arg *GetNodeByIdArg) (*GetNodeByIdResponse, *APIError) {
	if arg == nil || arg.NodeId == "" {
		return nil, NewAPIError(400, "node_id is required")
	}

	path := c.ExpandPath("/v1/nodes/{node_id}", map[string]string{"node_id": arg.NodeId})

	var resp GetNodeByIdResponse
	if apiErr := c.Get(ctx, path, nil, &resp); apiErr != nil {
		return nil, apiErr
	}
	return &resp, nil
}

type GetNodeHardwareByIdArg struct {
	NodeId string
}

type GetNodeHardwareByIdResponse struct {
	Vcpus  int                  `json:"vcpus"`
	CPU    NodeHardwareCpu      `json:"cpu"`
	Memory []NodeHardwareMemory `json:"memory"`
	Disks  []NodeHardwareDisk   `json:"disks"`
	Usb    []NodeHardwareUsb    `json:"usb"`
}

func GetNodeHardwareById(ctx context.Context, c *Client, arg *GetNodeHardwareByIdArg) (*GetNodeHardwareByIdResponse, *APIError) {
	if arg == nil || arg.NodeId == "" {
		return nil, NewAPIError(400, "node_id is required")
	}

	path := c.ExpandPath("/v1/nodes/{node_id}/hardware", map[string]string{"node_id": arg.NodeId})

	var resp GetNodeHardwareByIdResponse
	if apiErr := c.Get(ctx, path, nil, &resp); apiErr != nil {
		return nil, apiErr
	}
	return &resp, nil
}

type GetNodeLicenseByIdArg struct {
	NodeId string
}
type GetNodeLicenseByIdResponse struct {
	Key    string `json:"key"`
	Expiry string `json:"expiry"`
	Valid  bool   `json:"valid"`
}

func GetNodeLicenseById(ctx context.Context, c *Client, arg *GetNodeByIdArg) (*GetNodeLicenseByIdResponse, *APIError) {
	if arg == nil || arg.NodeId == "" {
		return nil, NewAPIError(400, "node_id is required")
	}

	path := c.ExpandPath("/v1/nodes/{node_id}/license", map[string]string{"node_id": arg.NodeId})

	var resp GetNodeLicenseByIdResponse
	if apiErr := c.Get(ctx, path, nil, &resp); apiErr != nil {
		return nil, apiErr
	}
	return &resp, nil
}

type GetNodeStoragePoolsByIdArg struct {
	NodeId string
}
type GetNodeStoragePoolsByIdResponse = []StoragePoolDetail

func GetNodeStoragePoolsById(ctx context.Context, c *Client, arg *GetNodeStoragePoolsByIdArg) (*GetNodeStoragePoolsByIdResponse, *APIError) {
	if arg == nil || arg.NodeId == "" {
		return nil, NewAPIError(400, "node_id is required")
	}

	path := c.ExpandPath("/v1/nodes/{node_id}/storage/pools", map[string]string{"node_id": arg.NodeId})

	var resp GetNodeStoragePoolsByIdResponse
	if apiErr := c.Get(ctx, path, nil, &resp); apiErr != nil {
		return nil, apiErr
	}
	return &resp, nil
}

type GetNodePciDevicesByIdArg struct {
	NodeId string
}
type GetNodePciDevicesByIdResponse = []NodePciDevice

func GetNodePciDevicesById(ctx context.Context, c *Client, arg *GetNodePciDevicesByIdArg) (*GetNodePciDevicesByIdResponse, *APIError) {
	if arg == nil || arg.NodeId == "" {
		return nil, NewAPIError(400, "node_id is required")
	}

	path := c.ExpandPath("/v1/nodes/{node_id}/hardware/pci", map[string]string{"node_id": arg.NodeId})

	var resp GetNodePciDevicesByIdResponse
	if apiErr := c.Get(ctx, path, nil, &resp); apiErr != nil {
		return nil, apiErr
	}
	return &resp, nil
}

type GetNodeMetricsArg struct {
	NodeId   string
	Interval string // "1h", "4h", "1d", "7d", "30d", "1y"
	Include  string // Comma-separated list of metrics
}

type NodeMetricData struct {
	Name   string      `json:"name"`
	Step   int         `json:"step"`
	Header []string    `json:"header"`
	Start  int64       `json:"start"`
	End    int64       `json:"end"`
	Data   [][]float64 `json:"data"`
}

type GetNodeMetricsResponse = []NodeMetricData

func GetNodeMetrics(ctx context.Context, c *Client, arg *GetNodeMetricsArg) (*GetNodeMetricsResponse, *APIError) {
	if arg == nil || arg.NodeId == "" || arg.Interval == "" {
		return nil, NewAPIError(400, "node_id and interval are required")
	}

	path := c.ExpandPath("/v1/nodes/{node_id}/metrics", map[string]string{"node_id": arg.NodeId})

	query := make(url.Values)
	query.Set("interval", arg.Interval)
	if arg.Include != "" {
		query.Set("include", arg.Include)
	}

	var resp GetNodeMetricsResponse
	if apiErr := c.Get(ctx, path, query, &resp); apiErr != nil {
		return nil, apiErr
	}
	return &resp, nil
}

type GetNodeLogsArg struct {
	NodeId   string
	Since    *int // Hours ago
	Until    *int // Hours ago
	MaxPoints *int // Max number of log points
}

type NodeLogEntry struct {
	Priority   string `json:"PRIORITY"`
	Facility   string `json:"FACILITY,omitempty"`
	Identifier string `json:"IDENTIFIER,omitempty"`
	Message    string `json:"MESSAGE"`
	Matches    int    `json:"matches"`
}

type GetNodeLogsResponse = []NodeLogEntry

func GetNodeLogs(ctx context.Context, c *Client, arg *GetNodeLogsArg) (*GetNodeLogsResponse, *APIError) {
	if arg == nil || arg.NodeId == "" {
		return nil, NewAPIError(400, "node_id is required")
	}

	path := c.ExpandPath("/v1/nodes/{node_id}/logs", map[string]string{"node_id": arg.NodeId})

	query := make(url.Values)
	if arg.Since != nil {
		query.Set("since", fmt.Sprintf("%d", *arg.Since))
	}
	if arg.Until != nil {
		query.Set("until", fmt.Sprintf("%d", *arg.Until))
	}
	if arg.MaxPoints != nil {
		query.Set("maxPoints", fmt.Sprintf("%d", *arg.MaxPoints))
	}

	var resp GetNodeLogsResponse
	if apiErr := c.Get(ctx, path, query, &resp); apiErr != nil {
		return nil, apiErr
	}
	return &resp, nil
}

type GetNodeCapabilitiesArg struct {
	NodeId string
}

// NodeCapability represents capabilities for a specific architecture
type NodeCapability struct {
	Arch          string                      `json:"arch"`
	WordSize      int                         `json:"wordsize"`
	Emulator      string                      `json:"emulator"`
	KvmSupported  bool                        `json:"kvm_supported"`
	Firmware      map[string]interface{}      `json:"firmware"`
	CpuModels     []NodeCpuModel              `json:"cpu_models"`
	DefaultMachine string                     `json:"default_machine,omitempty"`
	Machines      []NodeMachine               `json:"machines"`
	Features      map[string]bool             `json:"features"`
}

type NodeCpuModel struct {
	Name         string `json:"name"`
	Vendor       string `json:"vendor"`
	Experimental bool   `json:"experimental"`
}

type NodeMachine struct {
	Name    string `json:"name"`
	MaxVcpus int    `json:"max_vcpus"`
}

type GetNodeCapabilitiesResponse = map[string]NodeCapability

func GetNodeCapabilities(ctx context.Context, c *Client, arg *GetNodeCapabilitiesArg) (*GetNodeCapabilitiesResponse, *APIError) {
	if arg == nil || arg.NodeId == "" {
		return nil, NewAPIError(400, "node_id is required")
	}

	path := c.ExpandPath("/v1/nodes/{node_id}/capabilities", map[string]string{"node_id": arg.NodeId})

	var resp GetNodeCapabilitiesResponse
	if apiErr := c.Get(ctx, path, nil, &resp); apiErr != nil {
		return nil, apiErr
	}
	return &resp, nil
}

type WakeNodeArg struct {
	NodeId string
}

type WakeNodeResponse = struct{}

func WakeNode(ctx context.Context, c *Client, arg *WakeNodeArg) (*WakeNodeResponse, *APIError) {
	if arg == nil || arg.NodeId == "" {
		return nil, NewAPIError(400, "node_id is required")
	}

	path := c.ExpandPath("/v1/nodes/{node_id}/wol", map[string]string{"node_id": arg.NodeId})

	payload, err := json.Marshal(map[string]interface{}{})
	if err != nil {
		return nil, NewAPIError(500, "failed to encode request payload")
	}

	var resp WakeNodeResponse
	if apiErr := c.Post(ctx, path, nil, bytes.NewReader(payload), &resp); apiErr != nil {
		return nil, apiErr
	}
	return &resp, nil
}

type GetNodeConsoleArg struct {
	NodeId string
}

type GetNodeConsoleResponse struct {
	SessionToken string `json:"session_token"`
}

func GetNodeConsole(ctx context.Context, c *Client, arg *GetNodeConsoleArg) (*GetNodeConsoleResponse, *APIError) {
	if arg == nil || arg.NodeId == "" {
		return nil, NewAPIError(400, "node_id is required")
	}

	path := c.ExpandPath("/v1/nodes/{node_id}/console", map[string]string{"node_id": arg.NodeId})

	payload, err := json.Marshal(map[string]interface{}{})
	if err != nil {
		return nil, NewAPIError(500, "failed to encode request payload")
	}

	var resp GetNodeConsoleResponse
	if apiErr := c.Post(ctx, path, nil, bytes.NewReader(payload), &resp); apiErr != nil {
		return nil, apiErr
	}
	return &resp, nil
}

type GetNodeNetworkInterfacesArg struct {
	NodeId string
}

type NodeNetworkInterface struct {
	Name          string `json:"name"`
	Mac           string `json:"mac,omitempty"`
	Unused        bool   `json:"unused"`
	Misconfigured bool   `json:"misconfigured"`
	Mtu           int    `json:"mtu"`
	Vswitch       string `json:"vswitch,omitempty"`
}

type GetNodeNetworkInterfacesResponse struct {
	Networks []NodeNetworkInterface `json:"networks"`
}

func GetNodeNetworkInterfaces(ctx context.Context, c *Client, arg *GetNodeNetworkInterfacesArg) (*GetNodeNetworkInterfacesResponse, *APIError) {
	if arg == nil || arg.NodeId == "" {
		return nil, NewAPIError(400, "node_id is required")
	}

	path := c.ExpandPath("/v1/nodes/{node_id}/network/interfaces", map[string]string{"node_id": arg.NodeId})

	var resp GetNodeNetworkInterfacesResponse
	if apiErr := c.Get(ctx, path, nil, &resp); apiErr != nil {
		return nil, apiErr
	}
	return &resp, nil
}

type GetNodeTasksArg struct {
	NodeId string
}

type NodeTask struct {
	Id        string `json:"id"`
	InstanceId string `json:"instance_id,omitempty"`
	Type      string `json:"type"`
	Status    string `json:"status"`
	StartTime string `json:"start_time,omitempty"`
	EndTime   string `json:"end_time,omitempty"`
	Progress  int    `json:"progress,omitempty"`
}

type GetNodeTasksResponse = []NodeTask

func GetNodeTasks(ctx context.Context, c *Client, arg *GetNodeTasksArg) (*GetNodeTasksResponse, *APIError) {
	if arg == nil || arg.NodeId == "" {
		return nil, NewAPIError(400, "node_id is required")
	}

	path := c.ExpandPath("/v1/nodes/{node_id}/tasks", map[string]string{"node_id": arg.NodeId})

	var resp GetNodeTasksResponse
	if apiErr := c.Get(ctx, path, nil, &resp); apiErr != nil {
		return nil, apiErr
	}
	return &resp, nil
}

type GetNodeJobsArg struct {
	NodeId string
}

type NodeJobStatus struct {
	Id     string `json:"id"`
	Name   string `json:"name"`
	Status string `json:"status"`
	NextRun string `json:"next_run,omitempty"`
	LastRun string `json:"last_run,omitempty"`
}

type GetNodeJobsResponse = []NodeJobStatus

func GetNodeJobs(ctx context.Context, c *Client, arg *GetNodeJobsArg) (*GetNodeJobsResponse, *APIError) {
	if arg == nil || arg.NodeId == "" {
		return nil, NewAPIError(400, "node_id is required")
	}

	path := c.ExpandPath("/v1/nodes/{node_id}/jobs", map[string]string{"node_id": arg.NodeId})

	var resp GetNodeJobsResponse
	if apiErr := c.Get(ctx, path, nil, &resp); apiErr != nil {
		return nil, apiErr
	}
	return &resp, nil
}

type GetNodeSshKeysArg struct {
	NodeId string
}

type SshKey struct {
	Id      string `json:"id"`
	Type    string `json:"type"`
	Key     string `json:"key"`
	Comment string `json:"comment,omitempty"`
}

type GetNodeSshKeysResponse struct {
	AuthorizedKeys []SshKey `json:"authorized_keys"`
	SystemKeys     []SshKey `json:"system_keys"`
}

func GetNodeSshKeys(ctx context.Context, c *Client, arg *GetNodeSshKeysArg) (*GetNodeSshKeysResponse, *APIError) {
	if arg == nil || arg.NodeId == "" {
		return nil, NewAPIError(400, "node_id is required")
	}

	path := c.ExpandPath("/v1/nodes/{node_id}/ssh_keys", map[string]string{"node_id": arg.NodeId})

	var resp GetNodeSshKeysResponse
	if apiErr := c.Get(ctx, path, nil, &resp); apiErr != nil {
		return nil, apiErr
	}
	return &resp, nil
}
