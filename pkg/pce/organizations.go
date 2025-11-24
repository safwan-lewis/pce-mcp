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
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

const organizationsHelpText = `\n\nOrganizations are the top-level entities within the Pextra CloudEnvironment (PCE) hierarchy.
They represent distinct tenants within the cloud, each with its own users, storage, network configurations, and compute resources.` + hierarchyHelpText

func ListOrganizations() (mcp.Tool, server.ToolHandlerFunc) {
	return mcp.NewTool("list_organizations",
		mcp.WithDescription(fmt.Sprintf("Retrieve a list of all organizations accessible to the user%s", organizationsHelpText)),
		mcp.WithToolAnnotation(mcp.ToolAnnotation{
			Title:        "List Organizations",
			ReadOnlyHint: mcp.ToBoolPtr(true),
		}),
	), handleListOrganizations
}

func handleListOrganizations(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	client, err := session.GetSession("sessionId")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	orgs, listErr := api.ListOrganizations(ctx, client, &api.ListOrganizationsArg{})
	if listErr != nil {
		return mcp.NewToolResultError(listErr.Error()), nil
	}

	return mcp.NewToolResultJSON(orgs)
}

func GetOrganizationById() (mcp.Tool, server.ToolHandlerFunc) {
	return mcp.NewTool("get_organization_by_id",
		mcp.WithDescription("Retrieve detailed information about a specific organization. This includes the datacenters, clusters, and nodes within the organization. Use this tool when asked for a tree/hierarchical view of infrastructure. Use this tool if a list of nodes, clusters, or datacenters is needed."),
		mcp.WithToolAnnotation(mcp.ToolAnnotation{
			Title:        "Get Organization By ID",
			ReadOnlyHint: mcp.ToBoolPtr(true),
		}),
		mcp.WithString("organization_id",
			mcp.Required(),
			mcp.Description("Unique organization id (format: org-<xxx>)"),
		),
	), handleGetOrganizationById
}

func handleGetOrganizationById(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	orgId, err := requiredParam[string](req, "organization_id")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	var client *api.Client
	if client, err = session.GetSession("sessionId"); err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	org, getErr := api.GetOrganizationById(ctx, client, &api.GetOrganizationByIdArg{
		OrganizationId: orgId,
	})
	if getErr != nil {
		return mcp.NewToolResultError(getErr.Error()), nil
	}

	return mcp.NewToolResultJSON(org)
}

func GetCurrentOrganization() (mcp.Tool, server.ToolHandlerFunc) {
	return mcp.NewTool("get_current_organization",
		mcp.WithDescription("Retrieve detailed information about the organization associated with the current node. Use this tool when asked for a tree/hierarchical view of infrastructure for the current organization."),
		mcp.WithToolAnnotation(mcp.ToolAnnotation{
			Title:        "Get Current Organization",
			ReadOnlyHint: mcp.ToBoolPtr(true),
		}),
	), handleGetCurrentOrganization
}

func handleGetCurrentOrganization(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
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

	// Retrieve organization using organization ID from node
	organizationId := node.Node.OrganizationId
	org, orgErr := api.GetOrganizationById(ctx, client, &api.GetOrganizationByIdArg{
		OrganizationId: organizationId,
	})
	if orgErr != nil {
		return mcp.NewToolResultError(orgErr.Error()), nil
	}
	return mcp.NewToolResultJSON(org)
}

/*func ListOrganizationAuditLogsById() (mcp.Tool, server.ToolHandlerFunc) {
	return mcp.NewTool("list_organization_audit_logs_by_id",
		mcp.WithDescription("Retrieve audit logs for a specific organization. Audit logs provide a record of actions and events that have occurred within the organization, useful for tracking changes and ensuring compliance."),
		mcp.WithToolAnnotation(mcp.ToolAnnotation{
			Title:        "List Organization Audit Logs By ID",
			ReadOnlyHint: mcp.ToBoolPtr(true),
		}),
		mcp.WithString("organization_id",
			mcp.Required(),
			mcp.Description("Unique organization id (format: org-<xxx>)"),
		),
		mcp.WithNumber("entries",
			mcp.Min(0),
			mcp.Max(5000),
			mcp.DefaultNumber(50),
			mcp.Description("The number of audit log entries to retrieve. Default is 50."),
		),
		pageNum,
	), handleListOrganizationAuditLogsById
}

func ListOrganizationUserLockoutsById() (mcp.Tool, server.ToolHandlerFunc) {
	return mcp.NewTool("list_organization_user_lockouts_by_id",
		mcp.WithDescription("Retrieve a list of user lockouts for a specific organization. User lockouts occur when users are temporarily prevented from accessing their accounts due to multiple failed login attempts or security policies."),
		mcp.WithToolAnnotation(mcp.ToolAnnotation{
			Title:        "List Organization User Lockouts By ID",
			ReadOnlyHint: mcp.ToBoolPtr(true),
		}),
		mcp.WithString("organization_id",
			mcp.Required(),
			mcp.Description("Unique organization id (format: org-<xxx>)"),
		),
		mcp.WithNumber("entries",
			mcp.Min(0),
			mcp.Max(5000),
			mcp.DefaultNumber(50),
			mcp.Description("The number of user lockout entries to retrieve. Default is 50."),
		),
		pageNum,
	), handleListOrganizationUserLockoutsById
}*/

func CreateOrganization() (mcp.Tool, server.ToolHandlerFunc) {
	return mcp.NewTool("create_organization",
		mcp.WithDescription("Create a new organization within the Pextra CloudEnvironment (PCE). Organizations are top-level entities that represent distinct tenants within the cloud."),
		mcp.WithToolAnnotation(mcp.ToolAnnotation{
			Title: "Create Organization",
		}),
		mcp.WithString("name",
			mcp.Required(),
			mcp.MinLength(nameDefaultMinLength),
			mcp.MaxLength(nameDefaultMaxLength),
			mcp.Pattern(nameRegex(nameDefaultMinLength, nameDefaultMaxLength)), // redundant, but being explicit
			mcp.Description("The name of the new organization."),
		),
		mcp.WithString("description",
			mcp.MaxLength(descriptionDefaultMaxLength),
			mcp.Description("A brief description of the organization."),
		),
	), handleCreateOrganization
}

func handleCreateOrganization(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	name, err := requiredParam[string](req, "name")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	description, err := optionalParam[string](req, "description")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	client, err := session.GetSession("sessionId")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	org, createErr := api.CreateOrganization(ctx, client, &api.CreateOrganizationArg{
		Name:        name,
		Description: description,
	})
	if createErr != nil {
		return mcp.NewToolResultError(createErr.Error()), nil
	}

	return mcp.NewToolResultJSON(org)
}

func DeleteOrganizationById() (mcp.Tool, server.ToolHandlerFunc) {
	return mcp.NewTool("delete_organization_by_id",
		mcp.WithDescription("Delete an existing organization. The organization must be empty of any datacenters (and consequently clusters and nodes) before it can be deleted."),
		mcp.WithToolAnnotation(mcp.ToolAnnotation{
			Title:           "Delete Organization By ID",
			DestructiveHint: mcp.ToBoolPtr(true),
		}),
		mcp.WithString("organization_id",
			mcp.Required(),
			mcp.Description("Unique organization id (format: org-<xxx>)"),
		),
		mcp.WithBoolean("are_you_sure",
			mcp.Required(),
			mcp.Description("A safety check to prevent accidental deletions. Must be set to true to proceed with deletion."),
		),
	), handleDeleteOrganizationById
}

func handleDeleteOrganizationById(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	orgId, err := requiredParam[string](req, "organization_id")
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

	client, err := session.GetSession("sessionId")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	_, deleteErr := api.DeleteOrganizationById(ctx, client, &api.DeleteOrganizationByIdArg{
		OrganizationId: orgId,
	})
	if deleteErr != nil {
		return mcp.NewToolResultError(deleteErr.Error()), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("Organization %s deleted successfully.", orgId)), nil
}

func ListOrganizationAuditLogs() (mcp.Tool, server.ToolHandlerFunc) {
	return mcp.NewTool("list_organization_audit_logs",
		mcp.WithDescription("Retrieve audit logs for a specific organization. Audit logs provide a record of actions and events that have occurred within the organization, useful for tracking changes and ensuring compliance."),
		mcp.WithToolAnnotation(mcp.ToolAnnotation{
			Title:        "List Organization Audit Logs",
			ReadOnlyHint: mcp.ToBoolPtr(true),
		}),
		mcp.WithString("organization_id",
			mcp.Required(),
			mcp.Description("Unique organization id (format: org-<xxx>)"),
		),
		mcp.WithNumber("entries",
			mcp.Min(1),
			mcp.Max(5000),
			mcp.DefaultNumber(50),
			mcp.Description("The number of audit log entries to retrieve. Default is 50."),
		),
		pageNum,
	), handleListOrganizationAuditLogs
}

func handleListOrganizationAuditLogs(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	orgId, err := requiredParam[string](req, "organization_id")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	arg := &api.ListOrganizationAuditLogsArg{OrganizationId: orgId}
	
	// Check if parameters were provided
	if _, ok := req.GetArguments()["entries"]; ok {
		entries, err := optionalParam[float64](req, "entries")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		entriesInt := int(entries)
		arg.Entries = &entriesInt
	}
	if _, ok := req.GetArguments()["page"]; ok {
		page, err := optionalParam[float64](req, "page")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		pageInt := int(page)
		arg.Page = &pageInt
	}

	client, err := session.GetSession("sessionId")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	logs, listErr := api.ListOrganizationAuditLogs(ctx, client, arg)
	if listErr != nil {
		return mcp.NewToolResultError(listErr.Error()), nil
	}

	return mcp.NewToolResultJSON(logs)
}

func ListOrganizationAuthProviders() (mcp.Tool, server.ToolHandlerFunc) {
	return mcp.NewTool("list_organization_auth_providers",
		mcp.WithDescription("Retrieve a list of authentication providers configured for a specific organization. Auth providers handle user authentication (e.g., local, LDAP)."),
		mcp.WithToolAnnotation(mcp.ToolAnnotation{
			Title:        "List Organization Auth Providers",
			ReadOnlyHint: mcp.ToBoolPtr(true),
		}),
		mcp.WithString("organization_id",
			mcp.Required(),
			mcp.Description("Unique organization id (format: org-<xxx>)"),
		),
	), handleListOrganizationAuthProviders
}

func handleListOrganizationAuthProviders(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	orgId, err := requiredParam[string](req, "organization_id")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	client, err := session.GetSession("sessionId")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	providers, listErr := api.ListOrganizationAuthProviders(ctx, client, &api.ListOrganizationAuthProvidersArg{
		OrganizationId: orgId,
	})
	if listErr != nil {
		return mcp.NewToolResultError(listErr.Error()), nil
	}

	return mcp.NewToolResultJSON(providers)
}

func DeleteOrganizationAuthProvider() (mcp.Tool, server.ToolHandlerFunc) {
	return mcp.NewTool("delete_organization_auth_provider",
		mcp.WithDescription("Delete an authentication provider from an organization. This action cannot be undone and will affect all users associated with this provider."),
		mcp.WithToolAnnotation(mcp.ToolAnnotation{
			Title:           "Delete Organization Auth Provider",
			DestructiveHint: mcp.ToBoolPtr(true),
		}),
		mcp.WithString("organization_id",
			mcp.Required(),
			mcp.Description("Unique organization id (format: org-<xxx>)"),
		),
		mcp.WithString("auth_provider_id",
			mcp.Required(),
			mcp.Description("Unique auth provider id"),
		),
		mcp.WithBoolean("are_you_sure",
			mcp.Required(),
			mcp.Description("A safety check to prevent accidental deletions. Must be set to true to proceed with deletion."),
		),
	), handleDeleteOrganizationAuthProvider
}

func handleDeleteOrganizationAuthProvider(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	orgId, err := requiredParam[string](req, "organization_id")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	authProviderId, err := requiredParam[string](req, "auth_provider_id")
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

	client, err := session.GetSession("sessionId")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	_, deleteErr := api.DeleteOrganizationAuthProvider(ctx, client, &api.DeleteOrganizationAuthProviderArg{
		OrganizationId: orgId,
		AuthProviderId: authProviderId,
	})
	if deleteErr != nil {
		return mcp.NewToolResultError(deleteErr.Error()), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("Auth provider %s deleted successfully.", authProviderId)), nil
}

func ListOrganizationRoles() (mcp.Tool, server.ToolHandlerFunc) {
	return mcp.NewTool("list_organization_roles",
		mcp.WithDescription("Retrieve a list of roles defined for a specific organization. Roles define sets of permissions that can be assigned to users."),
		mcp.WithToolAnnotation(mcp.ToolAnnotation{
			Title:        "List Organization Roles",
			ReadOnlyHint: mcp.ToBoolPtr(true),
		}),
		mcp.WithString("organization_id",
			mcp.Required(),
			mcp.Description("Unique organization id (format: org-<xxx>)"),
		),
	), handleListOrganizationRoles
}

func handleListOrganizationRoles(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	orgId, err := requiredParam[string](req, "organization_id")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	client, err := session.GetSession("sessionId")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	roles, listErr := api.ListOrganizationRoles(ctx, client, &api.ListOrganizationRolesArg{
		OrganizationId: orgId,
	})
	if listErr != nil {
		return mcp.NewToolResultError(listErr.Error()), nil
	}

	return mcp.NewToolResultJSON(roles)
}

func CreateOrganizationRole() (mcp.Tool, server.ToolHandlerFunc) {
	return mcp.NewTool("create_organization_role",
		mcp.WithDescription("Create a new role for an organization. Roles define sets of permissions that can be assigned to users. Permissions must be specified as an array of objects with 'resource' and 'actions' fields."),
		mcp.WithToolAnnotation(mcp.ToolAnnotation{
			Title: "Create Organization Role",
		}),
		mcp.WithString("organization_id",
			mcp.Required(),
			mcp.Description("Unique organization id (format: org-<xxx>)"),
		),
		mcp.WithString("name",
			mcp.Required(),
			mcp.Description("The name of the role (3-32 characters)"),
		),
		mcp.WithString("description",
			mcp.Description("Optional description of the role (max 512 characters)"),
		),
		mcp.WithString("permissions_json",
			mcp.Required(),
			mcp.Description("JSON string representing an array of permission objects. Each permission must have 'resource' (string) and 'actions' (array of strings) fields. Example: '[{\"resource\":\"instances\",\"actions\":[\"read\",\"write\"]}]'"),
		),
	), handleCreateOrganizationRole
}

func handleCreateOrganizationRole(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	orgId, err := requiredParam[string](req, "organization_id")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	name, err := requiredParam[string](req, "name")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	description, err := optionalParam[string](req, "description")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	permissionsJson, err := requiredParam[string](req, "permissions_json")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	var permissions []api.Permission
	if err := json.Unmarshal([]byte(permissionsJson), &permissions); err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Invalid permissions JSON: %v", err)), nil
	}

	client, err := session.GetSession("sessionId")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	role, createErr := api.CreateOrganizationRole(ctx, client, &api.CreateOrganizationRoleArg{
		OrganizationId: orgId,
		Name:           name,
		Description:    description,
		Permissions:    permissions,
	})
	if createErr != nil {
		return mcp.NewToolResultError(createErr.Error()), nil
	}

	return mcp.NewToolResultJSON(role)
}

func UpdateOrganizationRole() (mcp.Tool, server.ToolHandlerFunc) {
	return mcp.NewTool("update_organization_role",
		mcp.WithDescription("Update an existing role for an organization. You can update the name, description, and/or permissions. Omitted fields will remain unchanged."),
		mcp.WithToolAnnotation(mcp.ToolAnnotation{
			Title: "Update Organization Role",
		}),
		mcp.WithString("organization_id",
			mcp.Required(),
			mcp.Description("Unique organization id (format: org-<xxx>)"),
		),
		mcp.WithString("role_id",
			mcp.Required(),
			mcp.Description("Unique role id"),
		),
		mcp.WithString("name",
			mcp.Description("The new name of the role (3-32 characters). Omit to leave unchanged."),
		),
		mcp.WithString("description",
			mcp.Description("The new description of the role (max 512 characters). Omit to leave unchanged."),
		),
		mcp.WithString("permissions_json",
			mcp.Description("JSON string representing an array of permission objects. Each permission must have 'resource' (string) and 'actions' (array of strings) fields. Omit to leave unchanged."),
		),
	), handleUpdateOrganizationRole
}

func handleUpdateOrganizationRole(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	orgId, err := requiredParam[string](req, "organization_id")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	roleId, err := requiredParam[string](req, "role_id")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	arg := &api.UpdateOrganizationRoleArg{
		OrganizationId: orgId,
		RoleId:         roleId,
	}
	
	// Check if parameters were provided and set them
	if _, ok := req.GetArguments()["name"]; ok {
		nameVal, err := optionalParam[string](req, "name")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		arg.Name = &nameVal
	}
	if _, ok := req.GetArguments()["description"]; ok {
		descriptionVal, err := optionalParam[string](req, "description")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		arg.Description = &descriptionVal
	}
	if _, ok := req.GetArguments()["permissions_json"]; ok {
		permissionsJsonVal, err := optionalParam[string](req, "permissions_json")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		var permissions []api.Permission
		if err := json.Unmarshal([]byte(permissionsJsonVal), &permissions); err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Invalid permissions JSON: %v", err)), nil
		}
		arg.Permissions = permissions
	}

	client, err := session.GetSession("sessionId")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	role, updateErr := api.UpdateOrganizationRole(ctx, client, arg)
	if updateErr != nil {
		return mcp.NewToolResultError(updateErr.Error()), nil
	}

	return mcp.NewToolResultJSON(role)
}

func ListOrganizationAiProviders() (mcp.Tool, server.ToolHandlerFunc) {
	return mcp.NewTool("list_organization_ai_providers",
		mcp.WithDescription("Retrieve a list of AI providers configured for a specific organization. AI providers enable AI features in the web UI."),
		mcp.WithToolAnnotation(mcp.ToolAnnotation{
			Title:        "List Organization AI Providers",
			ReadOnlyHint: mcp.ToBoolPtr(true),
		}),
		mcp.WithString("organization_id",
			mcp.Required(),
			mcp.Description("Unique organization id (format: org-<xxx>)"),
		),
	), handleListOrganizationAiProviders
}

func handleListOrganizationAiProviders(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	orgId, err := requiredParam[string](req, "organization_id")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	client, err := session.GetSession("sessionId")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	providers, listErr := api.ListOrganizationAiProviders(ctx, client, &api.ListOrganizationAiProvidersArg{
		OrganizationId: orgId,
	})
	if listErr != nil {
		return mcp.NewToolResultError(listErr.Error()), nil
	}

	return mcp.NewToolResultJSON(providers)
}

func ListSupportedAiProviders() (mcp.Tool, server.ToolHandlerFunc) {
	return mcp.NewTool("list_supported_ai_providers",
		mcp.WithDescription("Retrieve a list of supported AI provider types that can be configured for organizations. This includes information about whether each provider is self-hosted."),
		mcp.WithToolAnnotation(mcp.ToolAnnotation{
			Title:        "List Supported AI Providers",
			ReadOnlyHint: mcp.ToBoolPtr(true),
		}),
	), handleListSupportedAiProviders
}

func handleListSupportedAiProviders(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	client, err := session.GetSession("sessionId")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	providers, listErr := api.ListSupportedAiProviders(ctx, client, &api.ListSupportedAiProvidersArg{})
	if listErr != nil {
		return mcp.NewToolResultError(listErr.Error()), nil
	}

	return mcp.NewToolResultJSON(providers)
}