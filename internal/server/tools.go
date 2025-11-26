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
package server

import (
	"log"

	"github.com/PextraCloud/pce-mcp/internal/config"
	"github.com/PextraCloud/pce-mcp/pkg/pce"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// addToolConditionally adds a tool only if it's allowed by the safety level.
func addToolConditionally(s *server.MCPServer, toolName string, toolFunc func() (mcp.Tool, server.ToolHandlerFunc), safetyLevel config.SafetyLevel) {
	if ShouldEnableTool(toolName, safetyLevel) {
		tool, handler := toolFunc()
		s.AddTool(tool, handler)
	} else {
		log.Printf("Tool '%s' disabled by safety level '%s'", toolName, safetyLevel.String())
	}
}

func addOrganizationTools(s *server.MCPServer, safetyLevel config.SafetyLevel) {
	addToolConditionally(s, "list_organizations", pce.ListOrganizations, safetyLevel)
	addToolConditionally(s, "get_organization_by_id", pce.GetOrganizationById, safetyLevel)
	addToolConditionally(s, "get_current_organization", pce.GetCurrentOrganization, safetyLevel)
	addToolConditionally(s, "list_organization_audit_logs", pce.ListOrganizationAuditLogs, safetyLevel)
	addToolConditionally(s, "create_organization", pce.CreateOrganization, safetyLevel)
	addToolConditionally(s, "delete_organization_by_id", pce.DeleteOrganizationById, safetyLevel)
	addToolConditionally(s, "list_organization_auth_providers", pce.ListOrganizationAuthProviders, safetyLevel)
	addToolConditionally(s, "delete_organization_auth_provider", pce.DeleteOrganizationAuthProvider, safetyLevel)
	addToolConditionally(s, "list_organization_roles", pce.ListOrganizationRoles, safetyLevel)
	addToolConditionally(s, "create_organization_role", pce.CreateOrganizationRole, safetyLevel)
	addToolConditionally(s, "update_organization_role", pce.UpdateOrganizationRole, safetyLevel)
	addToolConditionally(s, "list_organization_ai_providers", pce.ListOrganizationAiProviders, safetyLevel)
	addToolConditionally(s, "list_supported_ai_providers", pce.ListSupportedAiProviders, safetyLevel)
}

func addUserTools(s *server.MCPServer, safetyLevel config.SafetyLevel) {
	addToolConditionally(s, "list_users_in_organization_by_id", pce.ListUsersInOrganizationById, safetyLevel)
	addToolConditionally(s, "invalidate_user_sessions_by_id", pce.InvalidateUserSessionsById, safetyLevel)
	addToolConditionally(s, "delete_user_by_id", pce.DeleteUserById, safetyLevel)
}

func addDatacenterTools(s *server.MCPServer, safetyLevel config.SafetyLevel) {
	addToolConditionally(s, "list_datacenters", pce.ListDatacenters, safetyLevel)
	addToolConditionally(s, "get_datacenter_by_id", pce.GetDatacenterById, safetyLevel)
	addToolConditionally(s, "create_datacenter", pce.CreateDatacenter, safetyLevel)
	addToolConditionally(s, "update_datacenter", pce.UpdateDatacenter, safetyLevel)
	addToolConditionally(s, "delete_datacenter_by_id", pce.DeleteDatacenterById, safetyLevel)
}

func addClusterTools(s *server.MCPServer, safetyLevel config.SafetyLevel) {
	addToolConditionally(s, "list_clusters", pce.ListClusters, safetyLevel)
	addToolConditionally(s, "get_cluster_by_id", pce.GetClusterById, safetyLevel)
	addToolConditionally(s, "initialize_cluster", pce.InitializeCluster, safetyLevel)
	addToolConditionally(s, "update_cluster", pce.UpdateCluster, safetyLevel)
	addToolConditionally(s, "get_cluster_join_key", pce.GetClusterJoinKey, safetyLevel)
	addToolConditionally(s, "get_cluster_metrics", pce.GetClusterMetrics, safetyLevel)
	addToolConditionally(s, "list_cluster_storage_pools", pce.ListClusterStoragePools, safetyLevel)
	addToolConditionally(s, "get_cluster_hardware_by_id", pce.GetClusterHardwareById, safetyLevel)
	addToolConditionally(s, "get_cluster_licensing_by_id", pce.GetClusterLicensingById, safetyLevel)
}

func addNodeTools(s *server.MCPServer, safetyLevel config.SafetyLevel) {
	addToolConditionally(s, "get_node_by_id", pce.GetNodeById, safetyLevel)
	addToolConditionally(s, "get_current_node", pce.GetCurrentNode, safetyLevel)
	addToolConditionally(s, "get_node_hardware_by_id", pce.GetNodeHardwareById, safetyLevel)
	addToolConditionally(s, "get_node_license_by_id", pce.GetNodeLicenseById, safetyLevel)
	addToolConditionally(s, "get_node_storagepools_by_id", pce.GetNodeStoragePoolsById, safetyLevel)
	addToolConditionally(s, "get_images", pce.GetImages, safetyLevel)
	addToolConditionally(s, "get_node_pcidevices_by_id", pce.GetNodePciDevicesById, safetyLevel)
	addToolConditionally(s, "get_node_metrics", pce.GetNodeMetrics, safetyLevel)
	addToolConditionally(s, "get_node_logs", pce.GetNodeLogs, safetyLevel)
	addToolConditionally(s, "get_node_capabilities", pce.GetNodeCapabilities, safetyLevel)
	addToolConditionally(s, "wake_node", pce.WakeNode, safetyLevel)
	addToolConditionally(s, "get_node_console", pce.GetNodeConsole, safetyLevel)
	addToolConditionally(s, "get_node_network_interfaces", pce.GetNodeNetworkInterfaces, safetyLevel)
	addToolConditionally(s, "get_node_tasks", pce.GetNodeTasks, safetyLevel)
	addToolConditionally(s, "get_node_jobs", pce.GetNodeJobs, safetyLevel)
	addToolConditionally(s, "get_node_ssh_keys", pce.GetNodeSshKeys, safetyLevel)
	addToolConditionally(s, "list_vswitches", pce.ListVSwitches, safetyLevel)
	addToolConditionally(s, "list_volumes", pce.ListVolumes, safetyLevel)
}

func addInstanceTools(s *server.MCPServer, safetyLevel config.SafetyLevel) {
	addToolConditionally(s, "get_instances_in_node", pce.GetInstancesInNode, safetyLevel)
	addToolConditionally(s, "get_instances_in_cluster", pce.GetInstancesInCluster, safetyLevel)
	addToolConditionally(s, "deploy_instance", pce.DeployInstance, safetyLevel)
	addToolConditionally(s, "get_instance_by_id", pce.GetInstanceById, safetyLevel)
	addToolConditionally(s, "update_instance", pce.UpdateInstance, safetyLevel)
	addToolConditionally(s, "delete_instance", pce.DeleteInstance, safetyLevel)
	addToolConditionally(s, "restore_instance", pce.RestoreInstance, safetyLevel)
	addToolConditionally(s, "get_instance_metrics", pce.GetInstanceMetrics, safetyLevel)
	addToolConditionally(s, "get_instance_console", pce.GetInstanceConsole, safetyLevel)
	addToolConditionally(s, "backup_instance", pce.BackupInstance, safetyLevel)
	addToolConditionally(s, "power_instance", pce.PowerInstance, safetyLevel)
	addToolConditionally(s, "get_instance_devices", pce.GetInstanceDevices, safetyLevel)
	addToolConditionally(s, "attach_device_to_instance", pce.AttachDeviceToInstance, safetyLevel)
}

// AddTools registers all tools with the server, filtered by the configured safety level.
func AddTools(s *server.MCPServer, safetyLevel config.SafetyLevel) {
	addOrganizationTools(s, safetyLevel)
	addUserTools(s, safetyLevel)
	addDatacenterTools(s, safetyLevel)
	addClusterTools(s, safetyLevel)
	addNodeTools(s, safetyLevel)
	addInstanceTools(s, safetyLevel)
}
