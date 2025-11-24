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
	"fmt"

	"github.com/PextraCloud/pce-mcp/internal/session"
	"github.com/PextraCloud/pce-mcp/pkg/api"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

const clustersHelpText = `\n\nClusters are groups of nodes that work together to provide compute resources.
They are organized within datacenters and are used to manage the distribution of workloads across nodes.` + hierarchyHelpText

func ListClusters() (mcp.Tool, server.ToolHandlerFunc) {
	return mcp.NewTool("list_clusters",
		mcp.WithDescription(fmt.Sprintf("Retrieve a list of all clusters in a datacenter%s", clustersHelpText)),
		mcp.WithToolAnnotation(mcp.ToolAnnotation{
			Title:        "List Clusters",
			ReadOnlyHint: mcp.ToBoolPtr(true),
		}),
		mcp.WithString("datacenter_id",
			mcp.Required(),
			mcp.Description("Unique datacenter id (format: dc-<xxx>)"),
		),
	), handleListClusters
}

func handleListClusters(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	datacenterId, err := requiredParam[string](req, "datacenter_id")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	client, err := session.GetSession("sessionId")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	clusters, listErr := api.ListClusters(ctx, client, &api.ListClustersArg{
		DatacenterId: datacenterId,
	})
	if listErr != nil {
		return mcp.NewToolResultError(listErr.Error()), nil
	}

	return mcp.NewToolResultJSON(clusters)
}

func GetClusterById() (mcp.Tool, server.ToolHandlerFunc) {
	return mcp.NewTool("get_cluster_by_id",
		mcp.WithDescription(fmt.Sprintf("Retrieve detailed information about a specific cluster, including all nodes within it%s", clustersHelpText)),
		mcp.WithToolAnnotation(mcp.ToolAnnotation{
			Title:        "Get Cluster By ID",
			ReadOnlyHint: mcp.ToBoolPtr(true),
		}),
		mcp.WithString("cluster_id",
			mcp.Required(),
			mcp.Description("Unique cluster id (format: cls-<xxx>)"),
		),
	), handleGetClusterById
}

func handleGetClusterById(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	clusterId, err := requiredParam[string](req, "cluster_id")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	client, err := session.GetSession("sessionId")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	cluster, getErr := api.GetClusterById(ctx, client, &api.GetClusterByIdArg{
		ClusterId: clusterId,
	})
	if getErr != nil {
		return mcp.NewToolResultError(getErr.Error()), nil
	}

	return mcp.NewToolResultJSON(cluster)
}

func InitializeCluster() (mcp.Tool, server.ToolHandlerFunc) {
	return mcp.NewTool("initialize_cluster",
		mcp.WithDescription("Initialize a standalone node into a new cluster. This converts the node from standalone mode to cluster mode."),
		mcp.WithToolAnnotation(mcp.ToolAnnotation{
			Title: "Initialize Cluster",
		}),
		mcp.WithString("node_id",
			mcp.Required(),
			mcp.Description("Unique node id (format: node-<xxx>)"),
		),
	), handleInitializeCluster
}

func handleInitializeCluster(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	nodeId, err := requiredParam[string](req, "node_id")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	client, err := session.GetSession("sessionId")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	result, initErr := api.InitializeCluster(ctx, client, &api.InitializeClusterArg{
		NodeId: nodeId,
	})
	if initErr != nil {
		return mcp.NewToolResultError(initErr.Error()), nil
	}

	return mcp.NewToolResultJSON(result)
}

func UpdateCluster() (mcp.Tool, server.ToolHandlerFunc) {
	return mcp.NewTool("update_cluster",
		mcp.WithDescription("Update details of an existing cluster (name, description)"),
		mcp.WithToolAnnotation(mcp.ToolAnnotation{
			Title: "Update Cluster",
		}),
		mcp.WithString("cluster_id",
			mcp.Required(),
			mcp.Description("Unique cluster id (format: cls-<xxx>)"),
		),
		mcp.WithString("name",
			mcp.MinLength(nameDefaultMinLength),
			mcp.MaxLength(nameDefaultMaxLength),
			mcp.Pattern(nameRegex(nameDefaultMinLength, nameDefaultMaxLength)),
			mcp.Description("The new name of the cluster."),
		),
		mcp.WithString("description",
			mcp.MaxLength(descriptionDefaultMaxLength),
			mcp.Description("A brief description of the cluster."),
		),
	), handleUpdateCluster
}

func handleUpdateCluster(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	clusterId, err := requiredParam[string](req, "cluster_id")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	client, err := session.GetSession("sessionId")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	updateArg := &api.UpdateClusterArg{
		ClusterId: clusterId,
	}

	name, err := optionalParam[string](req, "name")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	if name != "" {
		updateArg.Name = name
	}

	description, err := optionalParam[string](req, "description")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	if description != "" {
		updateArg.Description = description
	}

	_, updateErr := api.UpdateCluster(ctx, client, updateArg)
	if updateErr != nil {
		return mcp.NewToolResultError(updateErr.Error()), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("Cluster %s updated successfully.", clusterId)), nil
}

func GetClusterJoinKey() (mcp.Tool, server.ToolHandlerFunc) {
	return mcp.NewTool("get_cluster_join_key",
		mcp.WithDescription("Retrieve the join key for a cluster. This key is used to join additional nodes to the cluster."),
		mcp.WithToolAnnotation(mcp.ToolAnnotation{
			Title:        "Get Cluster Join Key",
			ReadOnlyHint: mcp.ToBoolPtr(true),
		}),
		mcp.WithString("cluster_id",
			mcp.Required(),
			mcp.Description("Unique cluster id (format: cls-<xxx>)"),
		),
	), handleGetClusterJoinKey
}

func handleGetClusterJoinKey(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	clusterId, err := requiredParam[string](req, "cluster_id")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	client, err := session.GetSession("sessionId")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	result, getErr := api.GetClusterJoinKey(ctx, client, &api.GetClusterJoinKeyArg{
		ClusterId: clusterId,
	})
	if getErr != nil {
		return mcp.NewToolResultError(getErr.Error()), nil
	}

	return mcp.NewToolResultJSON(result)
}

func GetClusterMetrics() (mcp.Tool, server.ToolHandlerFunc) {
	return mcp.NewTool("get_cluster_metrics",
		mcp.WithDescription("Retrieve metrics (CPU, memory, storage usage) for all nodes in a cluster"),
		mcp.WithToolAnnotation(mcp.ToolAnnotation{
			Title:        "Get Cluster Metrics",
			ReadOnlyHint: mcp.ToBoolPtr(true),
		}),
		mcp.WithString("cluster_id",
			mcp.Required(),
			mcp.Description("Unique cluster id (format: cls-<xxx>)"),
		),
		mcp.WithBoolean("force_refresh",
			mcp.Description("Whether to force a refresh of the metrics cache"),
			mcp.DefaultBool(false),
		),
	), handleGetClusterMetrics
}

func handleGetClusterMetrics(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	clusterId, err := requiredParam[string](req, "cluster_id")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	forceRefresh, err := optionalParam[bool](req, "force_refresh")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	client, err := session.GetSession("sessionId")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	metrics, getErr := api.GetClusterMetrics(ctx, client, &api.GetClusterMetricsArg{
		ClusterId:    clusterId,
		ForceRefresh: forceRefresh,
	})
	if getErr != nil {
		return mcp.NewToolResultError(getErr.Error()), nil
	}

	return mcp.NewToolResultJSON(metrics)
}

func ListClusterStoragePools() (mcp.Tool, server.ToolHandlerFunc) {
	return mcp.NewTool("list_cluster_storage_pools",
		mcp.WithDescription("Retrieve a list of all storage pools in a cluster"),
		mcp.WithToolAnnotation(mcp.ToolAnnotation{
			Title:        "List Cluster Storage Pools",
			ReadOnlyHint: mcp.ToBoolPtr(true),
		}),
		mcp.WithString("cluster_id",
			mcp.Required(),
			mcp.Description("Unique cluster id (format: cls-<xxx>)"),
		),
	), handleListClusterStoragePools
}

func handleListClusterStoragePools(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	clusterId, err := requiredParam[string](req, "cluster_id")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	client, err := session.GetSession("sessionId")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	pools, listErr := api.ListClusterStoragePools(ctx, client, &api.ListClusterStoragePoolsArg{
		ClusterId: clusterId,
	})
	if listErr != nil {
		return mcp.NewToolResultError(listErr.Error()), nil
	}

	return mcp.NewToolResultJSON(&struct {
		Pools *api.ListClusterStoragePoolsResponse `json:"pools"`
	}{
		Pools: pools,
	})
}

func GetClusterHardwareById() (mcp.Tool, server.ToolHandlerFunc) {
	return mcp.NewTool("get_cluster_hardware_by_id",
		mcp.WithDescription(fmt.Sprintf("Retrieve aggregated hardware information about all nodes in a specific cluster%s", clustersHelpText)),
		mcp.WithToolAnnotation(mcp.ToolAnnotation{
			Title:        "Get Cluster Hardware By ID",
			ReadOnlyHint: mcp.ToBoolPtr(true),
		}),
		mcp.WithString("cluster_id",
			mcp.Required(),
			mcp.Description("Unique cluster id (format: cls-<xxx>)"),
		),
	), handleGetClusterHardwareById
}

func handleGetClusterHardwareById(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	clusterId, err := requiredParam[string](req, "cluster_id")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	client, err := session.GetSession("sessionId")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	hardware, getErr := api.GetClusterHardwareById(ctx, client, &api.GetClusterHardwareByIdArg{
		ClusterId: clusterId,
	})
	if getErr != nil {
		return mcp.NewToolResultError(getErr.Error()), nil
	}

	return mcp.NewToolResultJSON(hardware)
}

func GetClusterLicensingById() (mcp.Tool, server.ToolHandlerFunc) {
	return mcp.NewTool("get_cluster_licensing_by_id",
		mcp.WithDescription(fmt.Sprintf("Retrieve aggregated licensing information about all nodes in a specific cluster. Only nodes that have licensing information will be included in the response. There is no separate 'cluster license', licenses only exist for individual nodes.%s", clustersHelpText)),
		mcp.WithToolAnnotation(mcp.ToolAnnotation{
			Title:        "Get Cluster License By ID",
			ReadOnlyHint: mcp.ToBoolPtr(true),
		}),
		mcp.WithString("cluster_id",
			mcp.Required(),
			mcp.Description("Unique cluster id (format: cls-<xxx>)"),
		),
	), handleGetClusterLicensingById
}

func handleGetClusterLicensingById(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	clusterId, err := requiredParam[string](req, "cluster_id")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	client, err := session.GetSession("sessionId")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	license, getErr := api.GetClusterLicensingById(ctx, client, &api.GetClusterLicensingByIdArg{
		ClusterId: clusterId,
	})
	if getErr != nil {
		return mcp.NewToolResultError(getErr.Error()), nil
	}
	return mcp.NewToolResultJSON(license)
}
