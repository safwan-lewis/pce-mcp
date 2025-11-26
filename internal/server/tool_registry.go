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
	"github.com/PextraCloud/pce-mcp/internal/config"
)

// ToolSafetyLevel maps tool names to their required safety level.
// Tools are categorized based on their annotations:
// - ReadOnlyHint=true: SafetyLevelReadOnly
// - ReadOnlyHint=false, DestructiveHint=false: SafetyLevelUpdate
// - DestructiveHint=true: SafetyLevelDelete
var ToolSafetyLevel = map[string]config.SafetyLevel{
	// Read-only tools (Get, List operations)
	"list_organizations":                      config.SafetyLevelReadOnly,
	"get_organization_by_id":                  config.SafetyLevelReadOnly,
	"get_current_organization":                config.SafetyLevelReadOnly,
	"list_organization_audit_logs":            config.SafetyLevelReadOnly,
	"list_organization_auth_providers":        config.SafetyLevelReadOnly,
	"list_organization_roles":                 config.SafetyLevelReadOnly,
	"list_organization_ai_providers":          config.SafetyLevelReadOnly,
	"list_supported_ai_providers":             config.SafetyLevelReadOnly,
	"list_users_in_organization_by_id":        config.SafetyLevelReadOnly,
	"list_datacenters":                        config.SafetyLevelReadOnly,
	"get_datacenter_by_id":                    config.SafetyLevelReadOnly,
	"list_clusters":                           config.SafetyLevelReadOnly,
	"get_cluster_by_id":                       config.SafetyLevelReadOnly,
	"get_cluster_join_key":                    config.SafetyLevelReadOnly,
	"get_cluster_metrics":                     config.SafetyLevelReadOnly,
	"list_cluster_storage_pools":              config.SafetyLevelReadOnly,
	"get_cluster_hardware_by_id":              config.SafetyLevelReadOnly,
	"get_cluster_licensing_by_id":             config.SafetyLevelReadOnly,
	"get_node_by_id":                          config.SafetyLevelReadOnly,
	"get_current_node":                        config.SafetyLevelReadOnly,
	"get_node_hardware_by_id":                 config.SafetyLevelReadOnly,
	"get_node_license_by_id":                  config.SafetyLevelReadOnly,
	"get_node_storagepools_by_id":             config.SafetyLevelReadOnly,
	"get_images":                              config.SafetyLevelReadOnly,
	"get_node_pcidevices_by_id":               config.SafetyLevelReadOnly,
	"get_node_metrics":                        config.SafetyLevelReadOnly,
	"get_node_logs":                           config.SafetyLevelReadOnly,
	"get_node_capabilities":                   config.SafetyLevelReadOnly,
	"get_node_console":                        config.SafetyLevelReadOnly,
	"get_node_network_interfaces":             config.SafetyLevelReadOnly,
	"get_node_tasks":                          config.SafetyLevelReadOnly,
	"get_node_jobs":                           config.SafetyLevelReadOnly,
	"get_node_ssh_keys":                       config.SafetyLevelReadOnly,
	"list_vswitches":                          config.SafetyLevelReadOnly,
	"list_volumes":                            config.SafetyLevelReadOnly,
	"get_instances_in_node":                   config.SafetyLevelReadOnly,
	"get_instances_in_cluster":                config.SafetyLevelReadOnly,
	"get_instance_by_id":                      config.SafetyLevelReadOnly,
	"get_instance_metrics":                    config.SafetyLevelReadOnly,
	"get_instance_console":                    config.SafetyLevelReadOnly,
	"get_instance_devices":                    config.SafetyLevelReadOnly,

	// Update/Write tools (Create, Update, Deploy, Backup, Power operations)
	"create_organization":                     config.SafetyLevelUpdate,
	"create_organization_role":                config.SafetyLevelUpdate,
	"update_organization_role":                config.SafetyLevelUpdate,
	"create_datacenter":                       config.SafetyLevelUpdate,
	"update_datacenter":                       config.SafetyLevelUpdate,
	"initialize_cluster":                      config.SafetyLevelUpdate,
	"update_cluster":                          config.SafetyLevelUpdate,
	"wake_node":                               config.SafetyLevelUpdate,
	"deploy_instance":                         config.SafetyLevelUpdate,
	"update_instance":                         config.SafetyLevelUpdate,
	"power_instance":                          config.SafetyLevelUpdate,
	"backup_instance":                         config.SafetyLevelUpdate,
	"attach_device_to_instance":               config.SafetyLevelUpdate,

	// Destructive tools (Delete, Restore operations)
	"delete_organization_by_id":               config.SafetyLevelDelete,
	"delete_organization_auth_provider":       config.SafetyLevelDelete,
	"invalidate_user_sessions_by_id":          config.SafetyLevelDelete,
	"delete_user_by_id":                       config.SafetyLevelDelete,
	"delete_datacenter_by_id":                 config.SafetyLevelDelete,
	"delete_instance":                         config.SafetyLevelDelete,
	"restore_instance":                        config.SafetyLevelDelete,
}

// ShouldEnableTool returns true if the tool should be enabled based on the configured safety level.
func ShouldEnableTool(toolName string, safetyLevel config.SafetyLevel) bool {
	requiredLevel, exists := ToolSafetyLevel[toolName]
	if !exists {
		// If tool is not in registry, default to requiring delete level (most restrictive)
		// This is a safety measure for unknown tools
		requiredLevel = config.SafetyLevelDelete
	}
	return safetyLevel >= requiredLevel
}

