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
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/url"
)

type ListOrganizationsArg struct{}
type ListOrganizationsResponse = []OrganizationDetail

func ListOrganizations(ctx context.Context, c *Client, arg *ListOrganizationsArg) (*ListOrganizationsResponse, *APIError) {
	path := c.ExpandPath("/v1/organizations", nil)

	var resp ListOrganizationsResponse
	if apiErr := c.Get(ctx, path, nil, &resp); apiErr != nil {
		return nil, apiErr
	}
	return &resp, nil
}

type GetOrganizationByIdArg struct {
	OrganizationId string
}
type GetOrganizationByIdResponse = OrganizationDetail

func GetOrganizationById(ctx context.Context, c *Client, arg *GetOrganizationByIdArg) (*GetOrganizationByIdResponse, *APIError) {
	if arg == nil {
		return nil, NewAPIError(400, "organization_id is required")
	}

	path := c.ExpandPath("/v1/organizations/{organization_id}", map[string]string{"organization_id": arg.OrganizationId})

	var resp GetOrganizationByIdResponse
	if apiErr := c.Get(ctx, path, nil, &resp); apiErr != nil {
		return nil, apiErr
	}
	return &resp, nil
}

type CreateOrganizationArg struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
}
type CreateOrganizationResponse = struct {
	Id string `json:"id"`
}

func CreateOrganization(ctx context.Context, c *Client, arg *CreateOrganizationArg) (*CreateOrganizationResponse, *APIError) {
	if arg == nil || arg.Name == "" {
		return nil, NewAPIError(400, "name is required")
	}

	path := c.ExpandPath("/v1/organizations", nil)

	payload, err := json.Marshal(map[string]string{"name": arg.Name, "description": arg.Description})
	if err != nil {
		return nil, NewAPIError(500, "failed to encode request payload")
	}

	var resp CreateOrganizationResponse
	if apiErr := c.Post(ctx, path, nil, bytes.NewReader(payload), &resp); apiErr != nil {
		return nil, apiErr
	}
	return &resp, nil
}

type DeleteOrganizationByIdArg struct {
	OrganizationId string
}
type DeleteOrganizationByIdResponse = struct{}

func DeleteOrganizationById(ctx context.Context, c *Client, arg *DeleteOrganizationByIdArg) (*DeleteOrganizationByIdResponse, *APIError) {
	if arg == nil {
		return nil, NewAPIError(400, "organization_id is required")
	}

	path := c.ExpandPath("/v1/organizations/{organization_id}", map[string]string{"organization_id": arg.OrganizationId})

	var resp DeleteOrganizationByIdResponse
	if apiErr := c.Delete(ctx, path, nil, &resp); apiErr != nil {
		return nil, apiErr
	}
	return &resp, nil
}

type ListOrganizationAuditLogsArg struct {
	OrganizationId string
	Entries        *int
	Page           *int
}

type AuditLogEntry struct {
	Id          string `json:"id"`
	Timestamp   string `json:"timestamp"`
	UserId      string `json:"user_id,omitempty"`
	Username    string `json:"username,omitempty"`
	Action      string `json:"action"`
	Resource    string `json:"resource,omitempty"`
	ResourceId  string `json:"resource_id,omitempty"`
	IpAddress   string `json:"ip_address,omitempty"`
	UserAgent   string `json:"user_agent,omitempty"`
	Success     bool   `json:"success"`
	Details     string `json:"details,omitempty"`
}

type ListOrganizationAuditLogsResponse = []AuditLogEntry

func ListOrganizationAuditLogs(ctx context.Context, c *Client, arg *ListOrganizationAuditLogsArg) (*ListOrganizationAuditLogsResponse, *APIError) {
	if arg == nil || arg.OrganizationId == "" {
		return nil, NewAPIError(400, "organization_id is required")
	}

	path := c.ExpandPath("/v1/organizations/{organization_id}/audit", map[string]string{"organization_id": arg.OrganizationId})

	query := make(url.Values)
	if arg.Entries != nil {
		query.Set("entries", fmt.Sprintf("%d", *arg.Entries))
	}
	if arg.Page != nil {
		query.Set("page", fmt.Sprintf("%d", *arg.Page))
	}

	var resp ListOrganizationAuditLogsResponse
	if apiErr := c.Get(ctx, path, query, &resp); apiErr != nil {
		return nil, apiErr
	}
	return &resp, nil
}

type ListOrganizationAuthProvidersArg struct {
	OrganizationId string
}

type AuthProvider struct {
	Id             string `json:"id"`
	OrganizationId string `json:"organization_id"`
	Name           string `json:"name"`
	Enabled        bool   `json:"enabled"`
	Type           int    `json:"type"` // 1 = local, 2 = LDAP
	Creation       string `json:"creation"`
	UserCount      int    `json:"user_count"`
}

type ListOrganizationAuthProvidersResponse = []AuthProvider

func ListOrganizationAuthProviders(ctx context.Context, c *Client, arg *ListOrganizationAuthProvidersArg) (*ListOrganizationAuthProvidersResponse, *APIError) {
	if arg == nil || arg.OrganizationId == "" {
		return nil, NewAPIError(400, "organization_id is required")
	}

	path := c.ExpandPath("/v1/organizations/{organization_id}/auth-providers", map[string]string{"organization_id": arg.OrganizationId})

	var resp ListOrganizationAuthProvidersResponse
	if apiErr := c.Get(ctx, path, nil, &resp); apiErr != nil {
		return nil, apiErr
	}
	return &resp, nil
}

type DeleteOrganizationAuthProviderArg struct {
	OrganizationId string
	AuthProviderId string
}

type DeleteOrganizationAuthProviderResponse = struct{}

func DeleteOrganizationAuthProvider(ctx context.Context, c *Client, arg *DeleteOrganizationAuthProviderArg) (*DeleteOrganizationAuthProviderResponse, *APIError) {
	if arg == nil || arg.OrganizationId == "" || arg.AuthProviderId == "" {
		return nil, NewAPIError(400, "organization_id and auth_provider_id are required")
	}

	path := c.ExpandPath("/v1/organizations/{organization_id}/auth-providers/{auth_provider_id}", map[string]string{
		"organization_id": arg.OrganizationId,
		"auth_provider_id": arg.AuthProviderId,
	})

	var resp DeleteOrganizationAuthProviderResponse
	if apiErr := c.Delete(ctx, path, nil, &resp); apiErr != nil {
		return nil, apiErr
	}
	return &resp, nil
}

type ListOrganizationRolesArg struct {
	OrganizationId string
}

type Permission struct {
	Resource string   `json:"resource"`
	Actions  []string `json:"actions"`
}

type Role struct {
	Id          string       `json:"id"`
	Name        string       `json:"name"`
	Description string       `json:"description,omitempty"`
	Permissions []Permission `json:"permissions"`
	UserCount   int          `json:"user_count,omitempty"`
}

type ListOrganizationRolesResponse = []Role

func ListOrganizationRoles(ctx context.Context, c *Client, arg *ListOrganizationRolesArg) (*ListOrganizationRolesResponse, *APIError) {
	if arg == nil || arg.OrganizationId == "" {
		return nil, NewAPIError(400, "organization_id is required")
	}

	path := c.ExpandPath("/v1/organizations/{organization_id}/roles", map[string]string{"organization_id": arg.OrganizationId})

	var resp ListOrganizationRolesResponse
	if apiErr := c.Get(ctx, path, nil, &resp); apiErr != nil {
		return nil, apiErr
	}
	return &resp, nil
}

type CreateOrganizationRoleArg struct {
	OrganizationId string       `json:"-"`
	Name           string       `json:"name"`
	Description    string       `json:"description,omitempty"`
	Permissions    []Permission `json:"permissions"`
}

type CreateOrganizationRoleResponse struct {
	Id string `json:"id"`
}

func CreateOrganizationRole(ctx context.Context, c *Client, arg *CreateOrganizationRoleArg) (*CreateOrganizationRoleResponse, *APIError) {
	if arg == nil || arg.OrganizationId == "" || arg.Name == "" {
		return nil, NewAPIError(400, "organization_id, name, and permissions are required")
	}
	if len(arg.Permissions) == 0 {
		return nil, NewAPIError(400, "permissions are required")
	}

	path := c.ExpandPath("/v1/organizations/{organization_id}/roles", map[string]string{"organization_id": arg.OrganizationId})

	payload, err := json.Marshal(map[string]interface{}{
		"name":        arg.Name,
		"description": arg.Description,
		"permissions": arg.Permissions,
	})
	if err != nil {
		return nil, NewAPIError(500, "failed to encode request payload")
	}

	var resp CreateOrganizationRoleResponse
	if apiErr := c.Post(ctx, path, nil, bytes.NewReader(payload), &resp); apiErr != nil {
		return nil, apiErr
	}
	return &resp, nil
}

type UpdateOrganizationRoleArg struct {
	OrganizationId string        `json:"-"`
	RoleId         string        `json:"-"`
	Name           *string       `json:"name,omitempty"`
	Description    *string       `json:"description,omitempty"`
	Permissions    []Permission  `json:"permissions,omitempty"`
}

type UpdateOrganizationRoleResponse = Role

func UpdateOrganizationRole(ctx context.Context, c *Client, arg *UpdateOrganizationRoleArg) (*UpdateOrganizationRoleResponse, *APIError) {
	if arg == nil || arg.OrganizationId == "" || arg.RoleId == "" {
		return nil, NewAPIError(400, "organization_id and role_id are required")
	}

	path := c.ExpandPath("/v1/organizations/{organization_id}/roles/{role_id}", map[string]string{
		"organization_id": arg.OrganizationId,
		"role_id":         arg.RoleId,
	})

	payloadMap := make(map[string]interface{})
	if arg.Name != nil {
		payloadMap["name"] = *arg.Name
	}
	if arg.Description != nil {
		payloadMap["description"] = *arg.Description
	}
	if arg.Permissions != nil {
		payloadMap["permissions"] = arg.Permissions
	}

	payload, err := json.Marshal(payloadMap)
	if err != nil {
		return nil, NewAPIError(500, "failed to encode request payload")
	}

	var resp UpdateOrganizationRoleResponse
	if apiErr := c.Patch(ctx, path, nil, bytes.NewReader(payload), &resp); apiErr != nil {
		return nil, apiErr
	}
	return &resp, nil
}

type ListOrganizationAiProvidersArg struct {
	OrganizationId string
}

type OrganizationAiProvider struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Type        string `json:"type"`
	DefaultModel string `json:"default_model,omitempty"`
	BaseUrl     string `json:"base_url,omitempty"`
	Enabled     bool   `json:"enabled"`
}

type ListOrganizationAiProvidersResponse = []OrganizationAiProvider

func ListOrganizationAiProviders(ctx context.Context, c *Client, arg *ListOrganizationAiProvidersArg) (*ListOrganizationAiProvidersResponse, *APIError) {
	if arg == nil || arg.OrganizationId == "" {
		return nil, NewAPIError(400, "organization_id is required")
	}

	path := c.ExpandPath("/v1/organizations/{organization_id}/ai/providers", map[string]string{"organization_id": arg.OrganizationId})

	var resp ListOrganizationAiProvidersResponse
	if apiErr := c.Get(ctx, path, nil, &resp); apiErr != nil {
		return nil, apiErr
	}
	return &resp, nil
}

type ListSupportedAiProvidersArg struct{}

type SupportedAiProvider struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	CompanyUrl  string `json:"company_url"`
	SelfHosted  bool   `json:"self_hosted"`
}

type ListSupportedAiProvidersResponse = []SupportedAiProvider

func ListSupportedAiProviders(ctx context.Context, c *Client, arg *ListSupportedAiProvidersArg) (*ListSupportedAiProvidersResponse, *APIError) {
	path := c.ExpandPath("/v1/organizations/ai/providers/supported", nil)

	var resp ListSupportedAiProvidersResponse
	if apiErr := c.Get(ctx, path, nil, &resp); apiErr != nil {
		return nil, apiErr
	}
	return &resp, nil
}
