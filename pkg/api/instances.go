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
	"net/url"

	"github.com/PextraCloud/pce-mcp/pkg/api/enum"
)

type GetInstancesByIdArg struct {
	// Either `NodeId` or `ClusterId` must be provided
	NodeId    string
	ClusterId string
}
type GetInstancesByIdResponse = []InstanceList

func GetInstancesById(ctx context.Context, c *Client, arg *GetInstancesByIdArg) (*GetInstancesByIdResponse, *APIError) {
	if arg == nil {
		return nil, NewAPIError(400, "either node_id or cluster_id is required")
	}
	if arg.NodeId != "" && arg.ClusterId != "" {
		return nil, NewAPIError(400, "only one of node_id or cluster_id should be provided")
	}

	query := make(url.Values)
	if arg.NodeId != "" {
		query.Set("node_id", arg.NodeId)
	} else if arg.ClusterId != "" {
		query.Set("cluster_id", arg.ClusterId)
	} else {
		return nil, NewAPIError(400, "either node_id or cluster_id is required")
	}

	path := "/v1/instances"

	var resp GetInstancesByIdResponse
	if apiErr := c.Get(ctx, path, query, &resp); apiErr != nil {
		return nil, apiErr
	}
	return &resp, nil
}

type PowerInstanceArg struct {
	InstanceId string
	NodeId     string
	Action     enum.InstancePowerAction
}
type PowerInstanceResponse struct {
	TaskId string `json:"task_id"`
}

func PowerInstance(ctx context.Context, c *Client, arg *PowerInstanceArg) (*PowerInstanceResponse, *APIError) {
	if arg == nil || arg.NodeId == "" || arg.InstanceId == "" {
		return nil, NewAPIError(400, "node_id and instance_id are required")
	}
	if !enum.InstancePowerAction.IsValid(arg.Action) {
		return nil, NewAPIError(400, "invalid action")
	}

	path := c.ExpandPath("/v1/instances/{instance_id}/power", map[string]string{"instance_id": arg.InstanceId})

	query := make(url.Values)
	query.Set("node_id", arg.NodeId)

	payload, err := json.Marshal(map[string]string{"action": string(arg.Action)})
	if err != nil {
		return nil, NewAPIError(500, "failed to encode request payload")
	}

	var resp PowerInstanceResponse
	if apiErr := c.Post(ctx, path, query, bytes.NewReader(payload), &resp); apiErr != nil {
		return nil, apiErr
	}
	return &resp, nil
}

type GetInstanceByIdArg struct {
	InstanceId string
}
type InstanceDetail struct {
	Instance InstanceFull `json:"instance"`
}
// InstanceFull contains detailed instance information
// For now using InstanceList as base, can be extended with additional fields from instance.full schema
type InstanceFull = InstanceList

type GetInstanceByIdResponse = InstanceDetail

func GetInstanceById(ctx context.Context, c *Client, arg *GetInstanceByIdArg) (*GetInstanceByIdResponse, *APIError) {
	if arg == nil || arg.InstanceId == "" {
		return nil, NewAPIError(400, "instance_id is required")
	}

	path := c.ExpandPath("/v1/instances/{instance_id}", map[string]string{"instance_id": arg.InstanceId})

	var resp GetInstanceByIdResponse
	if apiErr := c.Get(ctx, path, nil, &resp); apiErr != nil {
		return nil, apiErr
	}
	return &resp, nil
}

type DeleteInstanceArg struct {
	InstanceId     string
	NodeId         string
	DestroyVolumes bool
}
type DeleteInstanceResponse = struct{}

func DeleteInstance(ctx context.Context, c *Client, arg *DeleteInstanceArg) (*DeleteInstanceResponse, *APIError) {
	if arg == nil || arg.InstanceId == "" || arg.NodeId == "" {
		return nil, NewAPIError(400, "instance_id and node_id are required")
	}

	path := c.ExpandPath("/v1/instances/{instance_id}", map[string]string{"instance_id": arg.InstanceId})

	query := make(url.Values)
	query.Set("node_id", arg.NodeId)
	if arg.DestroyVolumes {
		query.Set("destroy_volumes", "true")
	}

	var resp DeleteInstanceResponse
	if apiErr := c.Delete(ctx, path, query, &resp); apiErr != nil {
		return nil, apiErr
	}
	return &resp, nil
}

type GetInstanceMetricsArg struct {
	InstanceId string
	NodeId     string
}
type InstanceNetworkMetrics struct {
	Rx float64 `json:"rx"`
	Tx float64 `json:"tx"`
}
type GetInstanceMetricsResponse struct {
	Timestamp   string                              `json:"timestamp"`
	CpuPercent  string                              `json:"cpu_percent"`
	Memory      struct {
		Percent string `json:"percent"`
		Usage   string `json:"usage"`
	} `json:"memory"`
	Network map[string]InstanceNetworkMetrics `json:"network"`
}

func GetInstanceMetrics(ctx context.Context, c *Client, arg *GetInstanceMetricsArg) (*GetInstanceMetricsResponse, *APIError) {
	if arg == nil || arg.InstanceId == "" || arg.NodeId == "" {
		return nil, NewAPIError(400, "instance_id and node_id are required")
	}

	path := c.ExpandPath("/v1/instances/{instance_id}/metrics", map[string]string{"instance_id": arg.InstanceId})

	query := make(url.Values)
	query.Set("node_id", arg.NodeId)

	var resp GetInstanceMetricsResponse
	if apiErr := c.Get(ctx, path, query, &resp); apiErr != nil {
		return nil, apiErr
	}
	return &resp, nil
}

type GetInstanceConsoleArg struct {
	InstanceId string
	NodeId     string
}
type GetInstanceConsoleResponse struct {
	SessionToken string `json:"session_token"`
}

func GetInstanceConsole(ctx context.Context, c *Client, arg *GetInstanceConsoleArg) (*GetInstanceConsoleResponse, *APIError) {
	if arg == nil || arg.InstanceId == "" || arg.NodeId == "" {
		return nil, NewAPIError(400, "instance_id and node_id are required")
	}

	path := c.ExpandPath("/v1/instances/{instance_id}/console", map[string]string{"instance_id": arg.InstanceId})

	query := make(url.Values)
	query.Set("node_id", arg.NodeId)

	payload, err := json.Marshal(map[string]interface{}{})
	if err != nil {
		return nil, NewAPIError(500, "failed to encode request payload")
	}

	var resp GetInstanceConsoleResponse
	if apiErr := c.Post(ctx, path, query, bytes.NewReader(payload), &resp); apiErr != nil {
		return nil, apiErr
	}
	return &resp, nil
}

type BackupInstanceArg struct {
	InstanceId string
	NodeId     string
}
type BackupInstanceResponse struct {
	TaskId string `json:"task_id"`
}

func BackupInstance(ctx context.Context, c *Client, arg *BackupInstanceArg) (*BackupInstanceResponse, *APIError) {
	if arg == nil || arg.InstanceId == "" || arg.NodeId == "" {
		return nil, NewAPIError(400, "instance_id and node_id are required")
	}

	path := c.ExpandPath("/v1/instances/{instance_id}/backup", map[string]string{"instance_id": arg.InstanceId})

	query := make(url.Values)
	query.Set("node_id", arg.NodeId)

	payload, err := json.Marshal(map[string]interface{}{})
	if err != nil {
		return nil, NewAPIError(500, "failed to encode request payload")
	}

	var resp BackupInstanceResponse
	if apiErr := c.Post(ctx, path, query, bytes.NewReader(payload), &resp); apiErr != nil {
		return nil, apiErr
	}
	return &resp, nil
}

type DeployInstanceV2Arg struct {
	Payload json.RawMessage
}
type DeployInstanceV2Response struct {
	TaskId string `json:"task_id"`
}

func DeployInstanceV2(ctx context.Context, c *Client, arg *DeployInstanceV2Arg) (*DeployInstanceV2Response, *APIError) {
	if arg == nil || len(arg.Payload) == 0 {
		return nil, NewAPIError(400, "payload is required")
	}

	path := "/v2/instances"

	var resp DeployInstanceV2Response
	if apiErr := c.Post(ctx, path, nil, bytes.NewReader(arg.Payload), &resp); apiErr != nil {
		return nil, apiErr
	}
	return &resp, nil
}

type UpdateInstanceArg struct {
	InstanceId  string
	Name        string
	Description string
}
type UpdateInstanceResponse = struct{}

func UpdateInstance(ctx context.Context, c *Client, arg *UpdateInstanceArg) (*UpdateInstanceResponse, *APIError) {
	if arg == nil || arg.InstanceId == "" {
		return nil, NewAPIError(400, "instance_id is required")
	}
	if arg.Name == "" && arg.Description == "" {
		return nil, NewAPIError(400, "at least one of name or description must be provided")
	}

	path := c.ExpandPath("/v1/instances/{instance_id}", map[string]string{"instance_id": arg.InstanceId})

	payloadMap := map[string]any{
		"data": map[string]any{
			"action": "update_info_fields",
		},
	}
	data := payloadMap["data"].(map[string]any)
	if arg.Name != "" {
		data["name"] = arg.Name
	}
	if arg.Description != "" {
		data["description"] = arg.Description
	}

	payload, err := json.Marshal(payloadMap)
	if err != nil {
		return nil, NewAPIError(500, "failed to encode request payload")
	}

	var resp UpdateInstanceResponse
	if apiErr := c.Patch(ctx, path, nil, bytes.NewReader(payload), &resp); apiErr != nil {
		return nil, apiErr
	}
	return &resp, nil
}

type RestoreInstanceArg struct {
	InstanceId     string
	NodeId         string
	ImageName      string
	VolumeMappings json.RawMessage
}
type RestoreInstanceResponse struct {
	TaskId string `json:"task_id"`
}

func RestoreInstance(ctx context.Context, c *Client, arg *RestoreInstanceArg) (*RestoreInstanceResponse, *APIError) {
	if arg == nil || arg.InstanceId == "" || arg.NodeId == "" || arg.ImageName == "" {
		return nil, NewAPIError(400, "instance_id, node_id, and image_name are required")
	}
	if len(arg.VolumeMappings) == 0 {
		return nil, NewAPIError(400, "volume_mappings is required")
	}

	path := c.ExpandPath("/v1/instances/{instance_id}/restore", map[string]string{"instance_id": arg.InstanceId})

	query := make(url.Values)
	query.Set("node_id", arg.NodeId)

	payloadMap := map[string]any{
		"image_name":      arg.ImageName,
		"volume_mappings": json.RawMessage(arg.VolumeMappings),
	}
	payload, err := json.Marshal(payloadMap)
	if err != nil {
		return nil, NewAPIError(500, "failed to encode request payload")
	}

	var resp RestoreInstanceResponse
	if apiErr := c.Post(ctx, path, query, bytes.NewReader(payload), &resp); apiErr != nil {
		return nil, apiErr
	}
	return &resp, nil
}

type GetInstanceDevicesArg struct {
	InstanceId string
	NodeId     string
}
type InstanceDevice struct {
	Name        string          `json:"name"`
	Type        int             `json:"type"`
	Description *string         `json:"description"`
	Metadata    json.RawMessage `json:"metadata"`
}
type GetInstanceDevicesResponse []InstanceDevice

func GetInstanceDevices(ctx context.Context, c *Client, arg *GetInstanceDevicesArg) (*GetInstanceDevicesResponse, *APIError) {
	if arg == nil || arg.InstanceId == "" || arg.NodeId == "" {
		return nil, NewAPIError(400, "instance_id and node_id are required")
	}

	path := c.ExpandPath("/v1/instances/{instance_id}/devices/", map[string]string{"instance_id": arg.InstanceId})

	query := make(url.Values)
	query.Set("node_id", arg.NodeId)

	var resp GetInstanceDevicesResponse
	if apiErr := c.Get(ctx, path, query, &resp); apiErr != nil {
		return nil, apiErr
	}
	return &resp, nil
}

type AttachDeviceToInstanceArg struct {
	InstanceId string
	NodeId     string
	Payload    json.RawMessage
}
type AttachDeviceToInstanceResponse struct {
	Name        string `json:"name"`
	NeedsReboot bool   `json:"needs_reboot"`
}

func AttachDeviceToInstance(ctx context.Context, c *Client, arg *AttachDeviceToInstanceArg) (*AttachDeviceToInstanceResponse, *APIError) {
	if arg == nil || arg.InstanceId == "" || arg.NodeId == "" {
		return nil, NewAPIError(400, "instance_id and node_id are required")
	}
	if len(arg.Payload) == 0 {
		return nil, NewAPIError(400, "payload is required")
	}

	path := c.ExpandPath("/v1/instances/{instance_id}/devices/", map[string]string{"instance_id": arg.InstanceId})

	query := make(url.Values)
	query.Set("node_id", arg.NodeId)

	var resp AttachDeviceToInstanceResponse
	if apiErr := c.Post(ctx, path, query, bytes.NewReader(arg.Payload), &resp); apiErr != nil {
		return nil, apiErr
	}
	return &resp, nil
}