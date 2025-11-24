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
	"github.com/PextraCloud/pce-mcp/pkg/pce"
	"github.com/mark3labs/mcp-go/server"
)

func addOrganizationTools(s *server.MCPServer) {
	s.AddTool(pce.ListOrganizations())
	s.AddTool(pce.GetOrganizationById())
	s.AddTool(pce.GetCurrentOrganization())
	// s.AddTool(pce.ListOrganizationAuditLogsById())
	// s.AddTool(pce.ListOrganizationUserLockoutsById())
	s.AddTool(pce.CreateOrganization())
	s.AddTool(pce.DeleteOrganizationById())
}

func addUserTools(s *server.MCPServer) {
	s.AddTool(pce.ListUsersInOrganizationById())
	s.AddTool(pce.InvalidateUserSessionsById())
	s.AddTool(pce.DeleteUserById())
}

func addDatacenterTools(s *server.MCPServer) {
	s.AddTool(pce.ListDatacenters())
	s.AddTool(pce.GetDatacenterById())
	s.AddTool(pce.CreateDatacenter())
	s.AddTool(pce.UpdateDatacenter())
	s.AddTool(pce.DeleteDatacenterById())
}

func addClusterTools(s *server.MCPServer) {
	s.AddTool(pce.ListClusters())
	s.AddTool(pce.GetClusterById())
	s.AddTool(pce.InitializeCluster())
	s.AddTool(pce.UpdateCluster())
	s.AddTool(pce.GetClusterJoinKey())
	s.AddTool(pce.GetClusterMetrics())
	s.AddTool(pce.ListClusterStoragePools())
	s.AddTool(pce.GetClusterHardwareById())
	s.AddTool(pce.GetClusterLicensingById())
}

func addNodeTools(s *server.MCPServer) {
	s.AddTool(pce.GetNodeById())
	s.AddTool(pce.GetCurrentNode())
	s.AddTool(pce.GetNodeHardwareById())
	s.AddTool(pce.GetNodeLicenseById())
	s.AddTool(pce.GetNodeStoragePoolsById())
	s.AddTool(pce.GetImages())
	s.AddTool(pce.GetNodePciDevicesById())
	s.AddTool(pce.GetNodeMetrics())
	s.AddTool(pce.GetNodeLogs())
	s.AddTool(pce.GetNodeCapabilities())
	s.AddTool(pce.WakeNode())
	s.AddTool(pce.GetNodeConsole())
	s.AddTool(pce.GetNodeNetworkInterfaces())
	s.AddTool(pce.GetNodeTasks())
	s.AddTool(pce.GetNodeJobs())
	s.AddTool(pce.GetNodeSshKeys())
}

func addInstanceTools(s *server.MCPServer) {
	s.AddTool(pce.GetInstancesInNode())
	s.AddTool(pce.GetInstancesInCluster())
	s.AddTool(pce.DeployInstance())
	s.AddTool(pce.GetInstanceById())
	s.AddTool(pce.UpdateInstance())
	s.AddTool(pce.DeleteInstance())
	s.AddTool(pce.RestoreInstance())
	s.AddTool(pce.GetInstanceMetrics())
	s.AddTool(pce.GetInstanceConsole())
	s.AddTool(pce.BackupInstance())
	s.AddTool(pce.PowerInstance())
}

func AddTools(s *server.MCPServer) {
	addOrganizationTools(s)
	addUserTools(s)
	addDatacenterTools(s)
	addClusterTools(s)
	addNodeTools(s)
	addInstanceTools(s)
}
