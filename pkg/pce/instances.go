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
	"encoding/json"
	"fmt"

	"github.com/PextraCloud/pce-mcp/internal/session"
	"github.com/PextraCloud/pce-mcp/pkg/api"
	"github.com/PextraCloud/pce-mcp/pkg/api/enum"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

const instancesHelpText = `\n\nInstances are virtual machines (QEMU/KVM) or containers (LXC) that run on nodes within clusters.
They utilize the compute resources of the nodes to perform various tasks and services.` + hierarchyHelpText

type getInstancesInNodeOrClusterResult struct {
	Instances *api.GetInstancesByIdResponse `json:"instances"`
}

func GetInstancesInCluster() (mcp.Tool, server.ToolHandlerFunc) {
	return mcp.NewTool("get_instances_in_cluster",
		mcp.WithDescription(fmt.Sprintf("Retrieve instances deployed on all nodes within a specific cluster%s", instancesHelpText)),
		mcp.WithToolAnnotation(mcp.ToolAnnotation{
			Title:        "Get Instances in Cluster",
			ReadOnlyHint: mcp.ToBoolPtr(true),
		}),
		mcp.WithString("cluster_id",
			mcp.Required(),
			mcp.Description("Unique cluster id (format: cls-<xxx>)"),
		),
	), handleGetInstancesInCluster
}

func handleGetInstancesInCluster(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	clusterId, err := requiredParam[string](req, "cluster_id")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	client, err := session.GetSession("sessionId")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	instances, getErr := api.GetInstancesById(ctx, client, &api.GetInstancesByIdArg{
		ClusterId: clusterId,
	})
	if getErr != nil {
		return mcp.NewToolResultError(getErr.Error()), nil
	}

	return mcp.NewToolResultJSON(&getInstancesInNodeOrClusterResult{
		Instances: instances,
	})
}

func GetInstancesInNode() (mcp.Tool, server.ToolHandlerFunc) {
	return mcp.NewTool("get_instances_in_node",
		mcp.WithDescription(fmt.Sprintf("Retrieve instances deployed on a specific node%s", instancesHelpText)),
		mcp.WithToolAnnotation(mcp.ToolAnnotation{
			Title:        "Get Instances in Node",
			ReadOnlyHint: mcp.ToBoolPtr(true),
		}),
		mcp.WithString("node_id",
			mcp.Required(),
			mcp.Description("Unique node id (format: node-<xxx>)"),
		),
	), handleGetInstancesInNode
}

func handleGetInstancesInNode(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	nodeId, err := requiredParam[string](req, "node_id")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	client, err := session.GetSession("sessionId")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	instances, getErr := api.GetInstancesById(ctx, client, &api.GetInstancesByIdArg{
		NodeId: nodeId,
	})
	if getErr != nil {
		return mcp.NewToolResultError(getErr.Error()), nil
	}

	return mcp.NewToolResultJSON(&getInstancesInNodeOrClusterResult{
		Instances: instances,
	})
}

func PowerInstance() (mcp.Tool, server.ToolHandlerFunc) {
	return mcp.NewTool("power_instance",
		mcp.WithDescription("Perform a power action on a specific instance (start, stop, restart, kill)"),
		mcp.WithToolAnnotation(mcp.ToolAnnotation{
			Title:        "Power Instance",
			ReadOnlyHint: mcp.ToBoolPtr(false),
		}),
		mcp.WithString("node_id",
			mcp.Required(),
			mcp.Description("Unique node id (format: node-<xxx>)"),
		),
		mcp.WithString("instance_id",
			mcp.Required(),
			mcp.Description("Unique instance id (format: inst-<xxx>)"),
		),
		mcp.WithString("action",
			mcp.Enum(
				string(enum.InstancePowerActionStart),
				string(enum.InstancePowerActionStop),
				string(enum.InstancePowerActionRestart),
				string(enum.InstancePowerActionKill),
			),
			mcp.Required(),
			mcp.Description("Power action to perform. start = power on, stop = graceful shutdown (may do nothing if the guest OS does not support it), restart = graceful reboot, kill = immediate power off, like pulling the power plug"),
		),
	), handlePowerInstance
}

func handlePowerInstance(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	nodeId, err := requiredParam[string](req, "node_id")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	instanceId, err := requiredParam[string](req, "instance_id")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	// Using `string` instead of  `enum.InstancePowerAction` since `requiredParam`` does not support enum types
	action, err := requiredParam[string](req, "action")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	client, err := session.GetSession("sessionId")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	res, powerErr := api.PowerInstance(ctx, client, &api.PowerInstanceArg{
		NodeId:     nodeId,
		InstanceId: instanceId,
		Action:     enum.InstancePowerAction(action),
	})
	if powerErr != nil {
		return mcp.NewToolResultError(powerErr.Error()), nil
	}

	return mcp.NewToolResultJSON(struct {
		Message string `json:"message"`
		TaskId  string `json:"task_id"`
	}{
		Message: "Power action initiated successfully",
		TaskId:  res.TaskId,
	})
}

func GetInstanceById() (mcp.Tool, server.ToolHandlerFunc) {
	return mcp.NewTool("get_instance_by_id",
		mcp.WithDescription(fmt.Sprintf("Retrieve detailed information about a specific instance%s", instancesHelpText)),
		mcp.WithToolAnnotation(mcp.ToolAnnotation{
			Title:        "Get Instance By ID",
			ReadOnlyHint: mcp.ToBoolPtr(true),
		}),
		mcp.WithString("instance_id",
			mcp.Required(),
			mcp.Description("Unique instance id (format: inst-<xxx>)"),
		),
	), handleGetInstanceById
}

func handleGetInstanceById(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	instanceId, err := requiredParam[string](req, "instance_id")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	client, err := session.GetSession("sessionId")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	instance, getErr := api.GetInstanceById(ctx, client, &api.GetInstanceByIdArg{
		InstanceId: instanceId,
	})
	if getErr != nil {
		return mcp.NewToolResultError(getErr.Error()), nil
	}

	return mcp.NewToolResultJSON(instance)
}

func DeleteInstance() (mcp.Tool, server.ToolHandlerFunc) {
	return mcp.NewTool("delete_instance",
		mcp.WithDescription("Delete an instance. This action is irreversible and will destroy the instance and all its data."),
		mcp.WithToolAnnotation(mcp.ToolAnnotation{
			Title:           "Delete Instance",
			DestructiveHint: mcp.ToBoolPtr(true),
		}),
		mcp.WithString("instance_id",
			mcp.Required(),
			mcp.Description("Unique instance id (format: inst-<xxx>)"),
		),
		mcp.WithString("node_id",
			mcp.Required(),
			mcp.Description("Unique node id (format: node-<xxx>)"),
		),
		mcp.WithBoolean("destroy_volumes",
			mcp.Description("Whether to destroy all volumes associated with the instance. This action is irreversible."),
			mcp.DefaultBool(false),
		),
		mcp.WithBoolean("are_you_sure",
			mcp.Required(),
			mcp.Description("A safety check to prevent accidental deletions. Must be set to true to proceed with deletion."),
		),
	), handleDeleteInstance
}

func handleDeleteInstance(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	instanceId, err := requiredParam[string](req, "instance_id")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	nodeId, err := requiredParam[string](req, "node_id")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	areYouSure, err := requiredParam[bool](req, "are_you_sure")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	if !areYouSure {
		return mcp.NewToolResultError("Deletion not confirmed. Set 'are_you_sure' to true to proceed."), nil
	}

	destroyVolumes, err := optionalParam[bool](req, "destroy_volumes")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	client, err := session.GetSession("sessionId")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	_, deleteErr := api.DeleteInstance(ctx, client, &api.DeleteInstanceArg{
		InstanceId:     instanceId,
		NodeId:         nodeId,
		DestroyVolumes: destroyVolumes,
	})
	if deleteErr != nil {
		return mcp.NewToolResultError(deleteErr.Error()), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("Instance %s deleted successfully.", instanceId)), nil
}

func GetInstanceMetrics() (mcp.Tool, server.ToolHandlerFunc) {
	return mcp.NewTool("get_instance_metrics",
		mcp.WithDescription("Retrieve metrics (CPU, memory, network usage) for a specific instance"),
		mcp.WithToolAnnotation(mcp.ToolAnnotation{
			Title:        "Get Instance Metrics",
			ReadOnlyHint: mcp.ToBoolPtr(true),
		}),
		mcp.WithString("instance_id",
			mcp.Required(),
			mcp.Description("Unique instance id (format: inst-<xxx>)"),
		),
		mcp.WithString("node_id",
			mcp.Required(),
			mcp.Description("Unique node id (format: node-<xxx>)"),
		),
	), handleGetInstanceMetrics
}

func handleGetInstanceMetrics(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	instanceId, err := requiredParam[string](req, "instance_id")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	nodeId, err := requiredParam[string](req, "node_id")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	client, err := session.GetSession("sessionId")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	metrics, getErr := api.GetInstanceMetrics(ctx, client, &api.GetInstanceMetricsArg{
		InstanceId: instanceId,
		NodeId:     nodeId,
	})
	if getErr != nil {
		return mcp.NewToolResultError(getErr.Error()), nil
	}

	return mcp.NewToolResultJSON(metrics)
}

func GetInstanceConsole() (mcp.Tool, server.ToolHandlerFunc) {
	return mcp.NewTool("get_instance_console",
		mcp.WithDescription("Get a console session token for connecting to an instance's console. For QEMU instances, use with noVNC. For other instance types, use with Xterm.js or similar."),
		mcp.WithToolAnnotation(mcp.ToolAnnotation{
			Title:        "Get Instance Console",
			ReadOnlyHint: mcp.ToBoolPtr(false),
		}),
		mcp.WithString("instance_id",
			mcp.Required(),
			mcp.Description("Unique instance id (format: inst-<xxx>)"),
		),
		mcp.WithString("node_id",
			mcp.Required(),
			mcp.Description("Unique node id (format: node-<xxx>)"),
		),
	), handleGetInstanceConsole
}

func handleGetInstanceConsole(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	instanceId, err := requiredParam[string](req, "instance_id")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	nodeId, err := requiredParam[string](req, "node_id")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	client, err := session.GetSession("sessionId")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	result, getErr := api.GetInstanceConsole(ctx, client, &api.GetInstanceConsoleArg{
		InstanceId: instanceId,
		NodeId:     nodeId,
	})
	if getErr != nil {
		return mcp.NewToolResultError(getErr.Error()), nil
	}

	return mcp.NewToolResultJSON(result)
}

func BackupInstance() (mcp.Tool, server.ToolHandlerFunc) {
	return mcp.NewTool("backup_instance",
		mcp.WithDescription("Create a backup of an instance. This will create a PXI image file that can be used to restore or deploy the instance later."),
		mcp.WithToolAnnotation(mcp.ToolAnnotation{
			Title: "Backup Instance",
		}),
		mcp.WithString("instance_id",
			mcp.Required(),
			mcp.Description("Unique instance id (format: inst-<xxx>)"),
		),
		mcp.WithString("node_id",
			mcp.Required(),
			mcp.Description("Unique node id (format: node-<xxx>)"),
		),
	), handleBackupInstance
}

func handleBackupInstance(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	instanceId, err := requiredParam[string](req, "instance_id")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	nodeId, err := requiredParam[string](req, "node_id")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	client, err := session.GetSession("sessionId")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	res, backupErr := api.BackupInstance(ctx, client, &api.BackupInstanceArg{
		InstanceId: instanceId,
		NodeId:     nodeId,
	})
	if backupErr != nil {
		return mcp.NewToolResultError(backupErr.Error()), nil
	}

	return mcp.NewToolResultJSON(struct {
		Message string `json:"message"`
		TaskId  string `json:"task_id"`
	}{
		Message: "Backup initiated successfully",
		TaskId:  res.TaskId,
	})
}

func DeployInstance() (mcp.Tool, server.ToolHandlerFunc) {
	return mcp.NewTool("deploy_instance",
		mcp.WithDescription("Deploy a new instance using the v2 deployment API. Specify either node_id or cluster_id for placement, and provide instance configuration including CPU, memory, networks, and metadata."),
		mcp.WithToolAnnotation(mcp.ToolAnnotation{
			Title: "Deploy Instance",
		}),
		// Destination: either node_id or cluster_id (mutually exclusive)
		mcp.WithString("node_id",
			mcp.Description("Unique node id (format: node-<xxx>) where to deploy the instance. Either node_id or cluster_id must be provided, but not both."),
		),
		mcp.WithString("cluster_id",
			mcp.Description("Unique cluster id (format: clus-<xxx>) for smart instance placement. Either node_id or cluster_id must be provided, but not both."),
		),
		// Basic instance info
		mcp.WithString("name",
			mcp.Required(),
			mcp.MinLength(3),
			mcp.MaxLength(255),
			mcp.Pattern("^[a-zA-Z0-9_-]+$"),
			mcp.Description("The name of the instance (3-255 characters, alphanumeric, underscores, hyphens)."),
		),
		mcp.WithString("description",
			mcp.MaxLength(512),
			mcp.Description("Description of the instance (max 512 characters)."),
		),
		mcp.WithString("architecture",
			mcp.Required(),
			mcp.Description("The CPU architecture of the instance (e.g., 'x86_64', 'aarch64'). Accepted values can be retrieved from the node capabilities endpoint."),
		),
		mcp.WithString("image",
			mcp.Description("The name of the image to deploy the instance from. Required for LXC instances."),
		),
		mcp.WithNumber("type",
			mcp.Required(),
			mcp.Min(0),
			mcp.Max(3),
			mcp.Description("Instance type: 0=Docker, 1=LXC, 2=QEMU, 3=Podman."),
		),
		// CPU configuration
		mcp.WithNumber("cpu_sockets",
			mcp.Required(),
			mcp.Min(1),
			mcp.Max(4),
			mcp.Description("The number of CPU sockets (1-4)."),
		),
		mcp.WithNumber("cpu_cores",
			mcp.Required(),
			mcp.Min(1),
			mcp.Description("The number of CPU cores per socket (minimum 1)."),
		),
		mcp.WithNumber("cpu_threads",
			mcp.Required(),
			mcp.Min(1),
			mcp.Description("The number of CPU threads per core (minimum 1)."),
		),
		mcp.WithString("cpu_affinity",
			mcp.Description("The CPU affinity string. Comma-separated list, where each item is one of: non-negative integer (e.g., '0', '1', '2'), range of non-negative integers (e.g., '0-3', '4-7'), caret (^) followed by a non-negative integer to exclude (e.g., '^2' excludes CPU 2). Examples: '0-3', '0,1,2,3', '0-7,^2,^5'."),
		),
		// Memory
		mcp.WithNumber("memory_mb",
			mcp.Required(),
			mcp.Min(64),
			mcp.Description("The amount of physical memory in megabytes (MB) for the instance (minimum 64 MB)."),
		),
		// Optional flags
		mcp.WithBoolean("imds_enabled",
			mcp.Description("Whether to enable access to the Instance Metadata Service (IMDS) for this instance. Defaults to false. This will add an internal network interface to the instance."),
		),
		mcp.WithBoolean("autostart",
			mcp.Description("Whether the instance should start automatically when the node boots. Defaults to false."),
		),
		// Complex nested structures as JSON
		mcp.WithString("networks_json",
			mcp.Required(),
			mcp.Description("JSON array of network interface configurations. Each network must have: vswitch_id (string), port_group_name (string), and optionally: mac (string), bandwidth (object), model (string), ipv4 (object), ipv6 (object). Example: [{\"vswitch_id\":\"vsw-xxx\",\"port_group_name\":\"spg0\"}]"),
		),
		mcp.WithString("metadata_json",
			mcp.Required(),
			mcp.Description("JSON object with instance-type-specific metadata. For LXC: {\"_type\":\"lxc\",\"init\":\"/sbin/init\"}. For QEMU: {\"_type\":\"qemu\",\"cpu_model\":\"...\",\"machine_type\":\"...\",\"firmware\":\"bios|efi\",\"secure_boot\":{...},\"new_volumes\":[...],\"existing_volumes\":[...]}. Refer to API spec for full structure."),
		),
	), handleDeployInstance
}

func handleDeployInstance(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Extract destination (node_id or cluster_id)
	nodeId, err := optionalParam[string](req, "node_id")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	clusterId, err := optionalParam[string](req, "cluster_id")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	if nodeId == "" && clusterId == "" {
		return mcp.NewToolResultError("either node_id or cluster_id is required"), nil
	}
	if nodeId != "" && clusterId != "" {
		return mcp.NewToolResultError("only one of node_id or cluster_id should be provided"), nil
	}

	// Build destination
	destination := api.DeployInstanceV2Destination{}
	if nodeId != "" {
		destination.NodeId = &nodeId
	} else {
		destination.ClusterId = &clusterId
	}

	// Extract basic fields
	name, err := requiredParam[string](req, "name")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	description, err := optionalParam[string](req, "description")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	architecture, err := requiredParam[string](req, "architecture")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	image, err := optionalParam[string](req, "image")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	typeVal, err := requiredParam[float64](req, "type")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	instanceType := enum.InstanceTypeEnum(int(typeVal))
	if instanceType < 0 || instanceType > 3 {
		return mcp.NewToolResultError("type must be between 0 and 3 (0=Docker, 1=LXC, 2=QEMU, 3=Podman)"), nil
	}

	// Extract CPU configuration
	cpuSockets, err := requiredParam[float64](req, "cpu_sockets")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	cpuCores, err := requiredParam[float64](req, "cpu_cores")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	cpuThreads, err := requiredParam[float64](req, "cpu_threads")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	cpuAffinity, err := optionalParam[string](req, "cpu_affinity")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	cpu := api.DeployInstanceV2CPU{
		Sockets:  int(cpuSockets),
		Cores:    int(cpuCores),
		Threads:  int(cpuThreads),
		Affinity: cpuAffinity,
	}

	// Extract memory
	memoryMB, err := requiredParam[float64](req, "memory_mb")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	if memoryMB < 64 {
		return mcp.NewToolResultError("memory_mb must be at least 64"), nil
	}

	// Extract optional flags
	imdsEnabled, err := optionalParam[bool](req, "imds_enabled")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	autostart, err := optionalParam[bool](req, "autostart")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	// Extract networks JSON
	networksJSON, err := requiredParam[string](req, "networks_json")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	networksBytes := []byte(networksJSON)
	if !json.Valid(networksBytes) {
		return mcp.NewToolResultError("networks_json must be valid JSON"), nil
	}
	var networks []api.DeployInstanceV2Network
	if err := json.Unmarshal(networksBytes, &networks); err != nil {
		return mcp.NewToolResultError("failed to parse networks_json: " + err.Error()), nil
	}
	if len(networks) == 0 {
		return mcp.NewToolResultError("at least one network is required"), nil
	}

	// Extract metadata JSON
	metadataJSON, err := requiredParam[string](req, "metadata_json")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	metadataBytes := []byte(metadataJSON)
	if !json.Valid(metadataBytes) {
		return mcp.NewToolResultError("metadata_json must be valid JSON"), nil
	}

	// Build the deployment argument
	deployArg := api.DeployInstanceV2Arg{
		Destination: destination,
		Name:        name,
		Architecture: architecture,
		Type:        instanceType,
		CPU:         cpu,
		Memory:      int(memoryMB),
		Networks:    networks,
		Metadata:    json.RawMessage(metadataBytes),
	}

	if description != "" {
		deployArg.Description = description
	}
	if image != "" {
		deployArg.Image = image
	}
	if imdsEnabled {
		deployArg.ImdsEnabled = true
	}
	if autostart {
		deployArg.Autostart = true
	}

	client, err := session.GetSession("sessionId")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	res, deployErr := api.DeployInstanceV2(ctx, client, &deployArg)
	if deployErr != nil {
		return mcp.NewToolResultError(deployErr.Error()), nil
	}

	return mcp.NewToolResultJSON(struct {
		Message string `json:"message"`
		TaskId  string `json:"task_id"`
	}{
		Message: "Instance deployment initiated successfully",
		TaskId:  res.TaskId,
	})
}

func UpdateInstance() (mcp.Tool, server.ToolHandlerFunc) {
	return mcp.NewTool("update_instance",
		mcp.WithDescription("Update basic instance information such as name and description."),
		mcp.WithToolAnnotation(mcp.ToolAnnotation{
			Title: "Update Instance",
		}),
		mcp.WithString("instance_id",
			mcp.Required(),
			mcp.Description("Unique instance id (format: inst-<xxx>)"),
		),
		mcp.WithString("name",
			mcp.Description("New name for the instance."),
		),
		mcp.WithString("description",
			mcp.Description("New description for the instance."),
		),
	), handleUpdateInstance
}

func handleUpdateInstance(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	instanceId, err := requiredParam[string](req, "instance_id")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	name, err := optionalParam[string](req, "name")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	description, err := optionalParam[string](req, "description")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	if name == "" && description == "" {
		return mcp.NewToolResultError("at least one of name or description must be provided"), nil
	}

	client, err := session.GetSession("sessionId")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	_, updateErr := api.UpdateInstance(ctx, client, &api.UpdateInstanceArg{
		InstanceId:  instanceId,
		Name:        name,
		Description: description,
	})
	if updateErr != nil {
		return mcp.NewToolResultError(updateErr.Error()), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("Instance %s updated successfully.", instanceId)), nil
}

func RestoreInstance() (mcp.Tool, server.ToolHandlerFunc) {
	return mcp.NewTool("restore_instance",
		mcp.WithDescription("Restore an instance from a backup PXI image. This overwrites the existing instance data."),
		mcp.WithToolAnnotation(mcp.ToolAnnotation{
			Title:           "Restore Instance",
			DestructiveHint: mcp.ToBoolPtr(true),
		}),
		mcp.WithString("instance_id",
			mcp.Required(),
			mcp.Description("Unique instance id (format: inst-<xxx>)"),
		),
		mcp.WithString("node_id",
			mcp.Required(),
			mcp.Description("Unique node id (format: node-<xxx>)"),
		),
		mcp.WithString("image_name",
			mcp.Required(),
			mcp.Description("Name of the PXI image to restore."),
		),
		mcp.WithString("volume_mappings_json",
			mcp.Required(),
			mcp.Description("JSON object describing volume mappings. Refer to the API documentation for structure."),
		),
	), handleRestoreInstance
}

func handleRestoreInstance(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	instanceId, err := requiredParam[string](req, "instance_id")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	nodeId, err := requiredParam[string](req, "node_id")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	imageName, err := requiredParam[string](req, "image_name")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	volumeMappingsStr, err := requiredParam[string](req, "volume_mappings_json")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	volumeMappingsBytes := []byte(volumeMappingsStr)
	if !json.Valid(volumeMappingsBytes) {
		return mcp.NewToolResultError("volume_mappings_json must be valid JSON"), nil
	}

	client, err := session.GetSession("sessionId")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	res, restoreErr := api.RestoreInstance(ctx, client, &api.RestoreInstanceArg{
		InstanceId:     instanceId,
		NodeId:         nodeId,
		ImageName:      imageName,
		VolumeMappings: json.RawMessage(volumeMappingsBytes),
	})
	if restoreErr != nil {
		return mcp.NewToolResultError(restoreErr.Error()), nil
	}

	return mcp.NewToolResultJSON(struct {
		Message string `json:"message"`
		TaskId  string `json:"task_id"`
	}{
		Message: "Instance restore initiated successfully",
		TaskId:  res.TaskId,
	})
}

func GetInstanceDevices() (mcp.Tool, server.ToolHandlerFunc) {
	return mcp.NewTool("get_instance_devices",
		mcp.WithDescription("Get all devices attached to an instance. This includes storage volumes, network interfaces, and other devices."),
		mcp.WithToolAnnotation(mcp.ToolAnnotation{
			Title:        "Get Instance Devices",
			ReadOnlyHint: mcp.ToBoolPtr(true),
		}),
		mcp.WithString("instance_id",
			mcp.Required(),
			mcp.Description("Unique instance id (format: inst-<xxx>)"),
		),
		mcp.WithString("node_id",
			mcp.Required(),
			mcp.Description("Unique node id (format: node-<xxx>)"),
		),
	), handleGetInstanceDevices
}

func handleGetInstanceDevices(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	instanceId, err := requiredParam[string](req, "instance_id")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	nodeId, err := requiredParam[string](req, "node_id")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	client, err := session.GetSession("sessionId")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	devices, getErr := api.GetInstanceDevices(ctx, client, &api.GetInstanceDevicesArg{
		InstanceId: instanceId,
		NodeId:     nodeId,
	})
	if getErr != nil {
		return mcp.NewToolResultError(getErr.Error()), nil
	}

	return mcp.NewToolResultJSON(struct {
		Devices *api.GetInstanceDevicesResponse `json:"devices"`
	}{
		Devices: devices,
	})
}

func AttachDeviceToInstance() (mcp.Tool, server.ToolHandlerFunc) {
	return mcp.NewTool("attach_device_to_instance",
		mcp.WithDescription("Attach a device to an instance (storage volume, network interface, etc.). The instance may need to be rebooted for changes to take effect."),
		mcp.WithToolAnnotation(mcp.ToolAnnotation{
			Title: "Attach Device to Instance",
		}),
		mcp.WithString("instance_id",
			mcp.Required(),
			mcp.Description("Unique instance id (format: inst-<xxx>)"),
		),
		mcp.WithString("node_id",
			mcp.Required(),
			mcp.Description("Unique node id (format: node-<xxx>)"),
		),
		mcp.WithString("payload_json",
			mcp.Required(),
			mcp.Description("JSON payload for the device to attach. Must include: name, type, and metadata. Refer to the API spec for device-specific metadata structure."),
		),
	), handleAttachDeviceToInstance
}

func handleAttachDeviceToInstance(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	instanceId, err := requiredParam[string](req, "instance_id")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	nodeId, err := requiredParam[string](req, "node_id")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	payloadStr, err := requiredParam[string](req, "payload_json")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	payloadBytes := []byte(payloadStr)
	if !json.Valid(payloadBytes) {
		return mcp.NewToolResultError("payload_json must be valid JSON"), nil
	}

	client, err := session.GetSession("sessionId")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	res, attachErr := api.AttachDeviceToInstance(ctx, client, &api.AttachDeviceToInstanceArg{
		InstanceId: instanceId,
		NodeId:     nodeId,
		Payload:    json.RawMessage(payloadBytes),
	})
	if attachErr != nil {
		return mcp.NewToolResultError(attachErr.Error()), nil
	}

	return mcp.NewToolResultJSON(struct {
		Message     string `json:"message"`
		DeviceName  string `json:"device_name"`
		NeedsReboot bool   `json:"needs_reboot"`
	}{
		Message:     "Device attached successfully",
		DeviceName:  res.Name,
		NeedsReboot: res.NeedsReboot,
	})
}