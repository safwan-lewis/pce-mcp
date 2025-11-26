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
package pce

import (
	"context"

	"github.com/PextraCloud/pce-mcp/internal/session"
	"github.com/PextraCloud/pce-mcp/pkg/api"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func ListVolumes() (mcp.Tool, server.ToolHandlerFunc) {
	return mcp.NewTool("list_volumes",
		mcp.WithDescription("List all volumes on a specific node or storage pool. Volumes are virtual storage devices that can be attached to instances."),
		mcp.WithToolAnnotation(mcp.ToolAnnotation{
			Title:        "List Volumes",
			ReadOnlyHint: mcp.ToBoolPtr(true),
		}),
		mcp.WithString("node_id",
			mcp.Required(),
			mcp.Description("Unique node id (format: node-<xxx>)"),
		),
		mcp.WithString("storage_pool_id",
			mcp.Description("Optional storage pool id to filter volumes (format: pool-<xxx>)"),
		),
	), handleListVolumes
}

type listVolumesResult struct {
	Volumes *api.ListVolumesResponse `json:"volumes"`
}

func handleListVolumes(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	nodeId, err := requiredParam[string](req, "node_id")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	storagePoolId, err := optionalParam[string](req, "storage_pool_id")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	var client *api.Client
	if client, err = session.GetSession("sessionId"); err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	var spIdPtr *string
	if storagePoolId != "" {
		spIdPtr = &storagePoolId
	}

	volumes, listErr := api.ListVolumes(ctx, client, &api.ListVolumesArg{
		NodeID:        nodeId,
		StoragePoolID: spIdPtr,
	})
	if listErr != nil {
		return mcp.NewToolResultError(listErr.Error()), nil
	}

	return mcp.NewToolResultJSON(&listVolumesResult{
		Volumes: volumes,
	})
}


