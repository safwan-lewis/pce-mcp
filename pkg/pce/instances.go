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
		mcp.WithDescription("Deploy a new instance using the v2 deployment API. Provide the payload in JSON format as described in the API specification."),
		mcp.WithToolAnnotation(mcp.ToolAnnotation{
			Title: "Deploy Instance",
		}),
		mcp.WithString("payload_json",
			mcp.Required(),
			mcp.Description("JSON payload for the deployment request. Refer to the API spec for the supported fields."),
		),
	), handleDeployInstance
}

func handleDeployInstance(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
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

	res, deployErr := api.DeployInstanceV2(ctx, client, &api.DeployInstanceV2Arg{
		Payload: json.RawMessage(payloadBytes),
	})
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