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
	"context"
	"net/url"
)

// VolumeMetadata represents the metadata of a volume
type VolumeMetadata struct {
	Driver string `json:"driver"`
}

// Volume represents a storage volume
type Volume struct {
	ID            string          `json:"id"`
	StoragePoolID string          `json:"storage_pool_id"`
	NodeID        string          `json:"node_id"`
	Name          string          `json:"name"`
	FQName        string          `json:"fq_name"`
	Description   *string         `json:"description"`
	Status        int             `json:"status"`
	Size          float64         `json:"size"`
	Metadata      VolumeMetadata  `json:"metadata"`
	Attached      bool            `json:"attached"`
	AttachedTo    *string         `json:"attached_to"`
	Creation      string          `json:"creation"`
}

type ListVolumesArg struct {
	NodeID        string
	StoragePoolID *string
}

type ListVolumesResponse = []Volume

func ListVolumes(ctx context.Context, c *Client, arg *ListVolumesArg) (*ListVolumesResponse, *APIError) {
	if arg == nil || arg.NodeID == "" {
		return nil, NewAPIError(400, "node_id is required")
	}

	path := "/v1/volumes/"
	query := make(url.Values)
	query.Set("node_id", arg.NodeID)
	if arg.StoragePoolID != nil && *arg.StoragePoolID != "" {
		query.Set("storage_pool_id", *arg.StoragePoolID)
	}

	var resp ListVolumesResponse
	if apiErr := c.Get(ctx, path, query, &resp); apiErr != nil {
		return nil, apiErr
	}
	return &resp, nil
}


