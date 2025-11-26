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

const nodesHelpText = `\n\nNodes are servers that provide the compute resources within Pextra CloudEnvironment (PCE). They are organized in clusters, and host instances (virtual machines or containers) that run workloads.` + hierarchyHelpText

func GetNodeById() (mcp.Tool, server.ToolHandlerFunc) {
	return mcp.NewTool("get_node_by_id",
		mcp.WithDescription(fmt.Sprintf("Retrieve detailed information about a specific node%s", nodesHelpText)),
		mcp.WithToolAnnotation(mcp.ToolAnnotation{
			Title:        "Get Node By ID",
			ReadOnlyHint: mcp.ToBoolPtr(true),
		}),
		mcp.WithString("node_id",
			mcp.Required(),
			mcp.Description("Unique node id (format: node-<xxx>)"),
		),
	), handleGetNodeById
}

func handleGetNodeById(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	nodeId, err := requiredParam[string](req, "node_id")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	client, err := session.GetSession("sessionId")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	node, getErr := api.GetNodeById(ctx, client, &api.GetNodeByIdArg{
		NodeId: nodeId,
	})
	if getErr != nil {
		return mcp.NewToolResultError(getErr.Error()), nil
	}

	return mcp.NewToolResultJSON(node)
}

func GetCurrentNode() (mcp.Tool, server.ToolHandlerFunc) {
	return mcp.NewTool("get_current_node",
		mcp.WithDescription("Retrieve detailed information about the current node"),
		mcp.WithToolAnnotation(mcp.ToolAnnotation{
			Title:        "Get Current Node",
			ReadOnlyHint: mcp.ToBoolPtr(true),
		}),
	), handleGetCurrentNode
}

func handleGetCurrentNode(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	client, err := session.GetSession("sessionId")
	if err != nil {
		fmt.Println("Error retrieving session:", err)
		return mcp.NewToolResultError(err.Error()), nil
	}

	// Get node ID through healthcheck endpoint
	health, healthErr := api.RunHealthcheck(ctx, client, &api.RunHealthcheckArg{})
	if healthErr != nil {
		return mcp.NewToolResultError(healthErr.Error()), nil
	}
	// This should never happen
	if !health.Healthy {
		return mcp.NewToolResultError("Node is not healthy"), nil
	}

	nodeId := health.Id
	node, getErr := api.GetNodeById(ctx, client, &api.GetNodeByIdArg{
		NodeId: nodeId,
	})
	if getErr != nil {
		return mcp.NewToolResultError(getErr.Error()), nil
	}

	return mcp.NewToolResultJSON(node)
}

func GetNodeHardwareById() (mcp.Tool, server.ToolHandlerFunc) {
	return mcp.NewTool("get_node_hardware_by_id",
		mcp.WithDescription("Retrieve detailed hardware information (CPU, memory, disks, USB devices) about a specific node."),
		mcp.WithToolAnnotation(mcp.ToolAnnotation{
			Title:        "Get Node Hardware By ID",
			ReadOnlyHint: mcp.ToBoolPtr(true),
		}),
		mcp.WithString("node_id",
			mcp.Required(),
			mcp.Description("Unique node id (format: node-<xxx>)"),
		),
	), handleGetNodeHardwareById
}

func handleGetNodeHardwareById(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	nodeId, err := requiredParam[string](req, "node_id")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	client, err := session.GetSession("sessionId")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	hardware, getErr := api.GetNodeHardwareById(ctx, client, &api.GetNodeHardwareByIdArg{
		NodeId: nodeId,
	})
	if getErr != nil {
		return mcp.NewToolResultError(getErr.Error()), nil
	}

	return mcp.NewToolResultJSON(hardware)
}

func GetNodeLicenseById() (mcp.Tool, server.ToolHandlerFunc) {
	return mcp.NewTool("get_node_license_by_id",
		mcp.WithDescription("Retrieve license information about a specific node."),
		mcp.WithToolAnnotation(mcp.ToolAnnotation{
			Title:        "Get Node License By ID",
			ReadOnlyHint: mcp.ToBoolPtr(true),
		}),
		mcp.WithString("node_id",
			mcp.Required(),
			mcp.Description("Unique node id (format: node-<xxx>)"),
		),
		mcp.WithBoolean("include_key",
			mcp.Description("Whether to include the license key in the response. This may expose sensitive information."),
			mcp.DefaultBool(false),
		),
	), handleGetNodeLicenseById
}

func handleGetNodeLicenseById(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	nodeId, err := requiredParam[string](req, "node_id")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	includeKey, err := optionalParam[bool](req, "include_key")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	client, err := session.GetSession("sessionId")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	license, getErr := api.GetNodeLicenseById(ctx, client, &api.GetNodeByIdArg{
		NodeId: nodeId,
	})
	if getErr != nil {
		return mcp.NewToolResultError(getErr.Error()), nil
	}

	// Redact license key if not explicitly requested
	if !includeKey {
		license.Key = "REDACTED"
	}
	return mcp.NewToolResultJSON(license)
}

func GetNodeStoragePoolsById() (mcp.Tool, server.ToolHandlerFunc) {
	return mcp.NewTool("get_node_storagepools_by_id",
		mcp.WithDescription("Retrieve storage pool information for a specific node."),
		mcp.WithToolAnnotation(mcp.ToolAnnotation{
			Title:        "Get Node Storage Pools By ID",
			ReadOnlyHint: mcp.ToBoolPtr(true),
		}),
		mcp.WithString("node_id",
			mcp.Required(),
			mcp.Description("Unique node id (format: node-<xxx>)"),
		),
	), handleGetNodeStoragePoolsById
}

type getNodeStoragePoolsByIdResult struct {
	Pools *api.GetNodeStoragePoolsByIdResponse `json:"pools"`
}

func handleGetNodeStoragePoolsById(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	nodeId, err := requiredParam[string](req, "node_id")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	client, err := session.GetSession("sessionId")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	pools, getErr := api.GetNodeStoragePoolsById(ctx, client, &api.GetNodeStoragePoolsByIdArg{
		NodeId: nodeId,
	})
	if getErr != nil {
		return mcp.NewToolResultError(getErr.Error()), nil
	}

	return mcp.NewToolResultJSON(&getNodeStoragePoolsByIdResult{
		Pools: pools,
	})
}

func GetNodePciDevicesById() (mcp.Tool, server.ToolHandlerFunc) {
	return mcp.NewTool("get_node_pcidevices_by_id",
		mcp.WithDescription("Retrieve PCI device information for a specific node. Use this tool if the user asks for GPUs, which are PCI devices."),
		mcp.WithToolAnnotation(mcp.ToolAnnotation{
			Title:        "Get Node PCI Devices By ID",
			ReadOnlyHint: mcp.ToBoolPtr(true),
		}),
		mcp.WithString("node_id",
			mcp.Required(),
			mcp.Description("Unique node id (format: node-<xxx>)"),
		),
	), handleGetNodePciDevicesById
}

type getNodePciDevicesByIdResult struct {
	Devices *api.GetNodePciDevicesByIdResponse `json:"devices"`
}

func handleGetNodePciDevicesById(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	nodeId, err := requiredParam[string](req, "node_id")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	client, err := session.GetSession("sessionId")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	devices, getErr := api.GetNodePciDevicesById(ctx, client, &api.GetNodePciDevicesByIdArg{
		NodeId: nodeId,
	})
	if getErr != nil {
		return mcp.NewToolResultError(getErr.Error()), nil
	}

	return mcp.NewToolResultJSON(&getNodePciDevicesByIdResult{
		Devices: devices,
	})
}

func GetNodeMetrics() (mcp.Tool, server.ToolHandlerFunc) {
	return mcp.NewTool("get_node_metrics",
		mcp.WithDescription("Retrieve metrics (CPU, memory, network, etc.) for a specific node. Metrics are returned as time-series data points."),
		mcp.WithToolAnnotation(mcp.ToolAnnotation{
			Title:        "Get Node Metrics",
			ReadOnlyHint: mcp.ToBoolPtr(true),
		}),
		mcp.WithString("node_id",
			mcp.Required(),
			mcp.Description("Unique node id (format: node-<xxx>)"),
		),
		mcp.WithString("interval",
			mcp.Required(),
			mcp.Description("Time interval for metrics: '1h', '4h', '1d', '7d', '30d', or '1y'"),
		),
		mcp.WithString("include",
			mcp.Description("Comma-separated list of metrics to include (e.g., 'cpu:percent-idle,memory:percent-used'). Default: 'cpu:percent-idle,memory:percent-used'"),
		),
	), handleGetNodeMetrics
}

func handleGetNodeMetrics(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	nodeId, err := requiredParam[string](req, "node_id")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	interval, err := requiredParam[string](req, "interval")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	include, err := optionalParam[string](req, "include")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	client, err := session.GetSession("sessionId")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	metrics, getErr := api.GetNodeMetrics(ctx, client, &api.GetNodeMetricsArg{
		NodeId:   nodeId,
		Interval: interval,
		Include:  include,
	})
	if getErr != nil {
		return mcp.NewToolResultError(getErr.Error()), nil
	}

	// Wrap the array in an object to match MCP library expectations
	result := map[string]interface{}{
		"metrics": metrics,
	}
	return mcp.NewToolResultJSON(result)
}

func GetNodeLogs() (mcp.Tool, server.ToolHandlerFunc) {
	return mcp.NewTool("get_node_logs",
		mcp.WithDescription("Retrieve system logs for a specific node. Logs are returned as syslog entries with priority, facility, and message information."),
		mcp.WithToolAnnotation(mcp.ToolAnnotation{
			Title:        "Get Node Logs",
			ReadOnlyHint: mcp.ToBoolPtr(true),
		}),
		mcp.WithString("node_id",
			mcp.Required(),
			mcp.Description("Unique node id (format: node-<xxx>)"),
		),
		mcp.WithNumber("since",
			mcp.Description("Query logs from this time (hours ago). Default is 1 hour ago."),
		),
		mcp.WithNumber("until",
			mcp.Description("Query logs to this time (hours ago). Default is now."),
		),
		mcp.WithNumber("max_points",
			mcp.Description("Maximum number of log entries to return. Default is 1000, maximum is 10000."),
		),
	), handleGetNodeLogs
}

func handleGetNodeLogs(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	nodeId, err := requiredParam[string](req, "node_id")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	arg := &api.GetNodeLogsArg{NodeId: nodeId}
	
	// Check if parameters were provided
	if _, ok := req.GetArguments()["since"]; ok {
		since, err := optionalParam[float64](req, "since")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		sinceInt := int(since)
		arg.Since = &sinceInt
	}
	if _, ok := req.GetArguments()["until"]; ok {
		until, err := optionalParam[float64](req, "until")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		untilInt := int(until)
		arg.Until = &untilInt
	}
	if _, ok := req.GetArguments()["max_points"]; ok {
		maxPoints, err := optionalParam[float64](req, "max_points")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		maxPointsInt := int(maxPoints)
		arg.MaxPoints = &maxPointsInt
	}

	client, err := session.GetSession("sessionId")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	logs, getErr := api.GetNodeLogs(ctx, client, arg)
	if getErr != nil {
		return mcp.NewToolResultError(getErr.Error()), nil
	}

	// Wrap the array in an object to match MCP library expectations
	result := map[string]interface{}{
		"logs": logs,
	}
	return mcp.NewToolResultJSON(result)
}

func GetNodeCapabilities() (mcp.Tool, server.ToolHandlerFunc) {
	return mcp.NewTool("get_node_capabilities",
		mcp.WithDescription("Retrieve virtualization capabilities for a specific node, including supported architectures, CPU models, machine types, and features."),
		mcp.WithToolAnnotation(mcp.ToolAnnotation{
			Title:        "Get Node Capabilities",
			ReadOnlyHint: mcp.ToBoolPtr(true),
		}),
		mcp.WithString("node_id",
			mcp.Required(),
			mcp.Description("Unique node id (format: node-<xxx>)"),
		),
	), handleGetNodeCapabilities
}

func handleGetNodeCapabilities(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	nodeId, err := requiredParam[string](req, "node_id")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	client, err := session.GetSession("sessionId")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	capabilities, getErr := api.GetNodeCapabilities(ctx, client, &api.GetNodeCapabilitiesArg{
		NodeId: nodeId,
	})
	if getErr != nil {
		return mcp.NewToolResultError(getErr.Error()), nil
	}

	return mcp.NewToolResultJSON(capabilities)
}

func WakeNode() (mcp.Tool, server.ToolHandlerFunc) {
	return mcp.NewTool("wake_node",
		mcp.WithDescription("Send a Wake-on-LAN (WOL) packet to a powered-off node to power it on. The node must have WOL enabled on its network interface and be connected to power."),
		mcp.WithToolAnnotation(mcp.ToolAnnotation{
			Title: "Wake Node",
		}),
		mcp.WithString("node_id",
			mcp.Required(),
			mcp.Description("Unique node id (format: node-<xxx>)"),
		),
	), handleWakeNode
}

func handleWakeNode(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	nodeId, err := requiredParam[string](req, "node_id")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	client, err := session.GetSession("sessionId")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	_, wakeErr := api.WakeNode(ctx, client, &api.WakeNodeArg{
		NodeId: nodeId,
	})
	if wakeErr != nil {
		return mcp.NewToolResultError(wakeErr.Error()), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("Wake-on-LAN packet sent to node %s", nodeId)), nil
}

func GetNodeConsole() (mcp.Tool, server.ToolHandlerFunc) {
	return mcp.NewTool("get_node_console",
		mcp.WithDescription("Get a console session token for connecting to a node's console. Use this token with Xterm.js via the consoleproxy websocket."),
		mcp.WithToolAnnotation(mcp.ToolAnnotation{
			Title:        "Get Node Console",
			ReadOnlyHint: mcp.ToBoolPtr(false),
		}),
		mcp.WithString("node_id",
			mcp.Required(),
			mcp.Description("Unique node id (format: node-<xxx>)"),
		),
	), handleGetNodeConsole
}

func handleGetNodeConsole(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	nodeId, err := requiredParam[string](req, "node_id")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	client, err := session.GetSession("sessionId")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	result, getErr := api.GetNodeConsole(ctx, client, &api.GetNodeConsoleArg{
		NodeId: nodeId,
	})
	if getErr != nil {
		return mcp.NewToolResultError(getErr.Error()), nil
	}

	return mcp.NewToolResultJSON(result)
}

func GetNodeNetworkInterfaces() (mcp.Tool, server.ToolHandlerFunc) {
	return mcp.NewTool("get_node_network_interfaces",
		mcp.WithDescription("Retrieve physical network interfaces (NICs) for a specific node, including their MAC addresses, MTU, and vSwitch assignments."),
		mcp.WithToolAnnotation(mcp.ToolAnnotation{
			Title:        "Get Node Network Interfaces",
			ReadOnlyHint: mcp.ToBoolPtr(true),
		}),
		mcp.WithString("node_id",
			mcp.Required(),
			mcp.Description("Unique node id (format: node-<xxx>)"),
		),
	), handleGetNodeNetworkInterfaces
}

func handleGetNodeNetworkInterfaces(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	nodeId, err := requiredParam[string](req, "node_id")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	client, err := session.GetSession("sessionId")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	interfaces, getErr := api.GetNodeNetworkInterfaces(ctx, client, &api.GetNodeNetworkInterfacesArg{
		NodeId: nodeId,
	})
	if getErr != nil {
		return mcp.NewToolResultError(getErr.Error()), nil
	}

	return mcp.NewToolResultJSON(interfaces)
}

func GetNodeTasks() (mcp.Tool, server.ToolHandlerFunc) {
	return mcp.NewTool("get_node_tasks",
		mcp.WithDescription("Retrieve tasks (long-lived operations) running on a specific node. Tasks may be triggered by API calls or by the node itself."),
		mcp.WithToolAnnotation(mcp.ToolAnnotation{
			Title:        "Get Node Tasks",
			ReadOnlyHint: mcp.ToBoolPtr(true),
		}),
		mcp.WithString("node_id",
			mcp.Required(),
			mcp.Description("Unique node id (format: node-<xxx>)"),
		),
	), handleGetNodeTasks
}

func handleGetNodeTasks(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	nodeId, err := requiredParam[string](req, "node_id")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	client, err := session.GetSession("sessionId")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	tasks, getErr := api.GetNodeTasks(ctx, client, &api.GetNodeTasksArg{
		NodeId: nodeId,
	})
	if getErr != nil {
		return mcp.NewToolResultError(getErr.Error()), nil
	}

	// Wrap the array in an object to match MCP library expectations
	result := map[string]interface{}{
		"tasks": tasks,
	}
	return mcp.NewToolResultJSON(result)
}

func GetNodeJobs() (mcp.Tool, server.ToolHandlerFunc) {
	return mcp.NewTool("get_node_jobs",
		mcp.WithDescription("Retrieve the status of all system jobs for a specific node. Jobs are periodic tasks that run automatically on nodes."),
		mcp.WithToolAnnotation(mcp.ToolAnnotation{
			Title:        "Get Node Jobs",
			ReadOnlyHint: mcp.ToBoolPtr(true),
		}),
		mcp.WithString("node_id",
			mcp.Required(),
			mcp.Description("Unique node id (format: node-<xxx>)"),
		),
	), handleGetNodeJobs
}

func handleGetNodeJobs(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	nodeId, err := requiredParam[string](req, "node_id")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	client, err := session.GetSession("sessionId")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	jobs, getErr := api.GetNodeJobs(ctx, client, &api.GetNodeJobsArg{
		NodeId: nodeId,
	})
	if getErr != nil {
		return mcp.NewToolResultError(getErr.Error()), nil
	}

	// Wrap the array in an object to match MCP library expectations
	result := map[string]interface{}{
		"jobs": jobs,
	}
	return mcp.NewToolResultJSON(result)
}

func GetNodeSshKeys() (mcp.Tool, server.ToolHandlerFunc) {
	return mcp.NewTool("get_node_ssh_keys",
		mcp.WithDescription("Retrieve internal SSH keys for a specific node. This includes keys managed by PCE (system keys and authorized keys), but not manually added keys. Private keys are not included for security reasons."),
		mcp.WithToolAnnotation(mcp.ToolAnnotation{
			Title:        "Get Node SSH Keys",
			ReadOnlyHint: mcp.ToBoolPtr(true),
		}),
		mcp.WithString("node_id",
			mcp.Required(),
			mcp.Description("Unique node id (format: node-<xxx>)"),
		),
	), handleGetNodeSshKeys
}

func handleGetNodeSshKeys(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	nodeId, err := requiredParam[string](req, "node_id")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	client, err := session.GetSession("sessionId")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	keys, getErr := api.GetNodeSshKeys(ctx, client, &api.GetNodeSshKeysArg{
		NodeId: nodeId,
	})
	if getErr != nil {
		return mcp.NewToolResultError(getErr.Error()), nil
	}

	// Wrap the array in an object to match MCP library expectations
	result := map[string]interface{}{
		"ssh_keys": keys,
	}
	return mcp.NewToolResultJSON(result)
}