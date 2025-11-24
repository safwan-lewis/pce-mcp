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
)

type GetClusterHardwareByIdArg struct {
	ClusterId string
}

type GetClusterHardwareByIdResponse struct {
	Partial   bool `json:"partial"`
	Vcpus     int  `json:"vcpus"`
	MemoryGB  int  `json:"memory_gb"`
	StorageGB int  `json:"storage_gb"`
}

func GetClusterHardwareById(ctx context.Context, c *Client, arg *GetClusterHardwareByIdArg) (*GetClusterHardwareByIdResponse, *APIError) {
	if arg == nil || arg.ClusterId == "" {
		return nil, NewAPIError(400, "cluster_id is required")
	}

	path := c.ExpandPath("/v1/clusters/{cluster_id}/hardware", map[string]string{"cluster_id": arg.ClusterId})

	var resp GetClusterHardwareByIdResponse
	if apiErr := c.Get(ctx, path, nil, &resp); apiErr != nil {
		return nil, apiErr
	}
	return &resp, nil
}

type GetClusterLicensingByIdArg struct {
	ClusterId string
}
type GetClusterLicensingByIdResponse struct {
	Cluster struct {
		Status     string `json:"status"`
		NextExpiry string `json:"next_expiry"`
	}
	// `node_id` -> { ok: bool, expiry: string }
	Nodes map[string]struct {
		Ok     bool   `json:"ok"`
		Expiry string `json:"expiry"`
	}
}

func GetClusterLicensingById(ctx context.Context, c *Client, arg *GetClusterLicensingByIdArg) (*GetClusterLicensingByIdResponse, *APIError) {
	if arg == nil || arg.ClusterId == "" {
		return nil, NewAPIError(400, "cluster_id is required")
	}

	path := c.ExpandPath("/v1/clusters/{cluster_id}/licensing", map[string]string{"cluster_id": arg.ClusterId})

	var resp GetClusterLicensingByIdResponse
	if apiErr := c.Get(ctx, path, nil, &resp); apiErr != nil {
		return nil, apiErr
	}
	return &resp, nil
}

type ListClustersArg struct {
	DatacenterId string
}
type ListClustersResponse = []ClusterList

func ListClusters(ctx context.Context, c *Client, arg *ListClustersArg) (*ListClustersResponse, *APIError) {
	if arg == nil || arg.DatacenterId == "" {
		return nil, NewAPIError(400, "datacenter_id is required")
	}

	path := "/v1/clusters"
	query := make(url.Values)
	query.Set("datacenter_id", arg.DatacenterId)

	var resp ListClustersResponse
	if apiErr := c.Get(ctx, path, query, &resp); apiErr != nil {
		return nil, apiErr
	}
	return &resp, nil
}

type GetClusterByIdArg struct {
	ClusterId string
}
type GetClusterByIdResponse = ClusterDetail

func GetClusterById(ctx context.Context, c *Client, arg *GetClusterByIdArg) (*GetClusterByIdResponse, *APIError) {
	if arg == nil || arg.ClusterId == "" {
		return nil, NewAPIError(400, "cluster_id is required")
	}

	path := c.ExpandPath("/v1/clusters/{cluster_id}", map[string]string{"cluster_id": arg.ClusterId})

	var resp GetClusterByIdResponse
	if apiErr := c.Get(ctx, path, nil, &resp); apiErr != nil {
		return nil, apiErr
	}
	return &resp, nil
}

type InitializeClusterArg struct {
	NodeId string `json:"node_id"`
}
type InitializeClusterResponse struct {
	Id      string `json:"id"`
	JoinKey string `json:"join_key"`
}

func InitializeCluster(ctx context.Context, c *Client, arg *InitializeClusterArg) (*InitializeClusterResponse, *APIError) {
	if arg == nil || arg.NodeId == "" {
		return nil, NewAPIError(400, "node_id is required")
	}

	path := "/v1/clusters"

	payload, err := json.Marshal(arg)
	if err != nil {
		return nil, NewAPIError(500, "failed to encode request payload")
	}

	var resp InitializeClusterResponse
	if apiErr := c.Post(ctx, path, nil, bytes.NewReader(payload), &resp); apiErr != nil {
		return nil, apiErr
	}
	return &resp, nil
}

type UpdateClusterArg struct {
	ClusterId   string
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
}
type UpdateClusterResponse = struct{}

func UpdateCluster(ctx context.Context, c *Client, arg *UpdateClusterArg) (*UpdateClusterResponse, *APIError) {
	if arg == nil || arg.ClusterId == "" {
		return nil, NewAPIError(400, "cluster_id is required")
	}

	path := c.ExpandPath("/v1/clusters/{cluster_id}", map[string]string{"cluster_id": arg.ClusterId})

	payloadMap := make(map[string]interface{})
	if arg.Name != "" {
		payloadMap["name"] = arg.Name
	}
	if arg.Description != "" {
		payloadMap["description"] = arg.Description
	}

	payload, err := json.Marshal(payloadMap)
	if err != nil {
		return nil, NewAPIError(500, "failed to encode request payload")
	}

	var resp UpdateClusterResponse
	if apiErr := c.Patch(ctx, path, nil, bytes.NewReader(payload), &resp); apiErr != nil {
		return nil, apiErr
	}
	return &resp, nil
}

type GetClusterJoinKeyArg struct {
	ClusterId string
}
type GetClusterJoinKeyResponse struct {
	JoinKey string `json:"join_key"`
}

func GetClusterJoinKey(ctx context.Context, c *Client, arg *GetClusterJoinKeyArg) (*GetClusterJoinKeyResponse, *APIError) {
	if arg == nil || arg.ClusterId == "" {
		return nil, NewAPIError(400, "cluster_id is required")
	}

	path := c.ExpandPath("/v1/clusters/{cluster_id}/joinkey", map[string]string{"cluster_id": arg.ClusterId})

	var resp GetClusterJoinKeyResponse
	if apiErr := c.Get(ctx, path, nil, &resp); apiErr != nil {
		return nil, apiErr
	}
	return &resp, nil
}

type GetClusterMetricsArg struct {
	ClusterId    string
	ForceRefresh bool
}
type ClusterNodeMetrics struct {
	Id      string  `json:"id"`
	Cpu     float64 `json:"cpu"`
	Memory  float64 `json:"memory"`
	Storage float64 `json:"storage"`
}
type GetClusterMetricsResponse struct {
	Partial bool                 `json:"partial"`
	Nodes   []ClusterNodeMetrics `json:"nodes"`
	Total   struct {
		Cpu     float64 `json:"cpu"`
		Memory  float64 `json:"memory"`
		Storage float64 `json:"storage"`
	} `json:"total"`
}

func GetClusterMetrics(ctx context.Context, c *Client, arg *GetClusterMetricsArg) (*GetClusterMetricsResponse, *APIError) {
	if arg == nil || arg.ClusterId == "" {
		return nil, NewAPIError(400, "cluster_id is required")
	}

	path := c.ExpandPath("/v1/clusters/{cluster_id}/metrics", map[string]string{"cluster_id": arg.ClusterId})

	query := make(url.Values)
	if arg.ForceRefresh {
		query.Set("force_refresh", "true")
	}

	var resp GetClusterMetricsResponse
	if apiErr := c.Get(ctx, path, query, &resp); apiErr != nil {
		return nil, apiErr
	}
	return &resp, nil
}

type ListClusterStoragePoolsArg struct {
	ClusterId string
}
type ListClusterStoragePoolsResponse = []StoragePoolList

type StoragePoolList struct {
	Id            string `json:"id"`
	ClusterId     string `json:"cluster_id"`
	Name          string `json:"name"`
	Type          int    `json:"type"`
	Description   string `json:"description"`
	NodeCount     int    `json:"node_count"`
	CanHoldImages bool   `json:"can_hold_images"`
	IsDefault     bool   `json:"is_default"`
	Creation      string `json:"creation"`
}

func ListClusterStoragePools(ctx context.Context, c *Client, arg *ListClusterStoragePoolsArg) (*ListClusterStoragePoolsResponse, *APIError) {
	if arg == nil || arg.ClusterId == "" {
		return nil, NewAPIError(400, "cluster_id is required")
	}

	path := c.ExpandPath("/v1/clusters/{cluster_id}/storage_pools", map[string]string{"cluster_id": arg.ClusterId})

	var resp ListClusterStoragePoolsResponse
	if apiErr := c.Get(ctx, path, nil, &resp); apiErr != nil {
		return nil, apiErr
	}
	return &resp, nil
}
