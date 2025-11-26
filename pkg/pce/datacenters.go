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

const datacentersHelpText = `\n\nDatacenters are logical groups of clusters, typically located in various geographical locations.
They are used to organize infrastructure and can help optimize latency when Pextra CloudEnvironment is deployed across multiple datacenters.` + hierarchyHelpText

func ListDatacenters() (mcp.Tool, server.ToolHandlerFunc) {
	return mcp.NewTool("list_datacenters",
		mcp.WithDescription(fmt.Sprintf("Retrieve a list of all datacenters in an organization%s", datacentersHelpText)),
		mcp.WithToolAnnotation(mcp.ToolAnnotation{
			Title:        "List Datacenters",
			ReadOnlyHint: mcp.ToBoolPtr(true),
		}),
		mcp.WithString("organization_id",
			mcp.Required(),
			mcp.Description("Unique organization id (format: org-<xxx>)"),
		),
	), handleListDatacenters
}

func handleListDatacenters(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	orgId, err := requiredParam[string](req, "organization_id")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	client, err := session.GetSession("sessionId")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	datacenters, listErr := api.ListDatacenters(ctx, client, &api.ListDatacentersArg{
		OrganizationId: orgId,
	})
	if listErr != nil {
		return mcp.NewToolResultError(listErr.Error()), nil
	}

	// Wrap the array in an object to match MCP library expectations
	result := map[string]interface{}{
		"datacenters": datacenters,
	}
	return mcp.NewToolResultJSON(result)
}

func GetDatacenterById() (mcp.Tool, server.ToolHandlerFunc) {
	return mcp.NewTool("get_datacenter_by_id",
		mcp.WithDescription(fmt.Sprintf("Retrieve detailed information about a specific datacenter, including all clusters within it%s", datacentersHelpText)),
		mcp.WithToolAnnotation(mcp.ToolAnnotation{
			Title:        "Get Datacenter By ID",
			ReadOnlyHint: mcp.ToBoolPtr(true),
		}),
		mcp.WithString("datacenter_id",
			mcp.Required(),
			mcp.Description("Unique datacenter id (format: dc-<xxx>)"),
		),
	), handleGetDatacenterById
}

func handleGetDatacenterById(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	datacenterId, err := requiredParam[string](req, "datacenter_id")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	client, err := session.GetSession("sessionId")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	datacenter, getErr := api.GetDatacenterById(ctx, client, &api.GetDatacenterByIdArg{
		DatacenterId: datacenterId,
	})
	if getErr != nil {
		return mcp.NewToolResultError(getErr.Error()), nil
	}

	return mcp.NewToolResultJSON(datacenter)
}

func CreateDatacenter() (mcp.Tool, server.ToolHandlerFunc) {
	return mcp.NewTool("create_datacenter",
		mcp.WithDescription(fmt.Sprintf("Create a new datacenter in an organization%s", datacentersHelpText)),
		mcp.WithToolAnnotation(mcp.ToolAnnotation{
			Title: "Create Datacenter",
		}),
		mcp.WithString("name",
			mcp.Required(),
			mcp.MinLength(nameDefaultMinLength),
			mcp.MaxLength(nameDefaultMaxLength),
			mcp.Pattern(nameRegex(nameDefaultMinLength, nameDefaultMaxLength)),
			mcp.Description("The name of the new datacenter."),
		),
		mcp.WithString("organization_id",
			mcp.Required(),
			mcp.Description("Unique organization id (format: org-<xxx>)"),
		),
		mcp.WithString("description",
			mcp.MaxLength(descriptionDefaultMaxLength),
			mcp.Description("A brief description of the datacenter."),
		),
	), handleCreateDatacenter
}

func handleCreateDatacenter(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	name, err := requiredParam[string](req, "name")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	orgId, err := requiredParam[string](req, "organization_id")
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

	datacenter, createErr := api.CreateDatacenter(ctx, client, &api.CreateDatacenterArg{
		Name:           name,
		OrganizationId: orgId,
		Description:    description,
	})
	if createErr != nil {
		return mcp.NewToolResultError(createErr.Error()), nil
	}

	return mcp.NewToolResultJSON(datacenter)
}

func UpdateDatacenter() (mcp.Tool, server.ToolHandlerFunc) {
	return mcp.NewTool("update_datacenter",
		mcp.WithDescription("Update details of an existing datacenter (name, description, location)"),
		mcp.WithToolAnnotation(mcp.ToolAnnotation{
			Title: "Update Datacenter",
		}),
		mcp.WithString("datacenter_id",
			mcp.Required(),
			mcp.Description("Unique datacenter id (format: dc-<xxx>)"),
		),
		mcp.WithString("name",
			mcp.MinLength(nameDefaultMinLength),
			mcp.MaxLength(nameDefaultMaxLength),
			mcp.Pattern(nameRegex(nameDefaultMinLength, nameDefaultMaxLength)),
			mcp.Description("The new name of the datacenter."),
		),
		mcp.WithString("description",
			mcp.MaxLength(descriptionDefaultMaxLength),
			mcp.Description("A brief description of the datacenter."),
		),
		mcp.WithNumber("latitude",
			mcp.Min(-90),
			mcp.Max(90),
			mcp.Description("The latitude of the datacenter location (-90 to 90)."),
		),
		mcp.WithNumber("longitude",
			mcp.Min(-180),
			mcp.Max(180),
			mcp.Description("The longitude of the datacenter location (-180 to 180)."),
		),
	), handleUpdateDatacenter
}

func handleUpdateDatacenter(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	datacenterId, err := requiredParam[string](req, "datacenter_id")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	client, err := session.GetSession("sessionId")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	updateArg := &api.UpdateDatacenterArg{
		DatacenterId: datacenterId,
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

	// Check if latitude or longitude are provided
	_, hasLat := req.GetArguments()["latitude"]
	_, hasLon := req.GetArguments()["longitude"]
	
	if hasLat || hasLon {
		if !hasLat || !hasLon {
			return mcp.NewToolResultError("both latitude and longitude must be provided if location is being updated"), nil
		}
		latitude, err := requiredParam[float64](req, "latitude")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		longitude, err := requiredParam[float64](req, "longitude")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		updateArg.Location = &struct {
			Latitude  float64 `json:"latitude"`
			Longitude float64 `json:"longitude"`
		}{
			Latitude:  latitude,
			Longitude: longitude,
		}
	}

	datacenter, updateErr := api.UpdateDatacenter(ctx, client, updateArg)
	if updateErr != nil {
		return mcp.NewToolResultError(updateErr.Error()), nil
	}

	return mcp.NewToolResultJSON(datacenter)
}

func DeleteDatacenterById() (mcp.Tool, server.ToolHandlerFunc) {
	return mcp.NewTool("delete_datacenter_by_id",
		mcp.WithDescription("Delete an existing datacenter. The datacenter must be empty of any clusters before it can be deleted. This action is irreversible."),
		mcp.WithToolAnnotation(mcp.ToolAnnotation{
			Title:           "Delete Datacenter By ID",
			DestructiveHint: mcp.ToBoolPtr(true),
		}),
		mcp.WithString("datacenter_id",
			mcp.Required(),
			mcp.Description("Unique datacenter id (format: dc-<xxx>)"),
		),
		mcp.WithBoolean("are_you_sure",
			mcp.Required(),
			mcp.Description("A safety check to prevent accidental deletions. Must be set to true to proceed with deletion."),
		),
	), handleDeleteDatacenterById
}

func handleDeleteDatacenterById(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	datacenterId, err := requiredParam[string](req, "datacenter_id")
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

	_, deleteErr := api.DeleteDatacenterById(ctx, client, &api.DeleteDatacenterByIdArg{
		DatacenterId: datacenterId,
	})
	if deleteErr != nil {
		return mcp.NewToolResultError(deleteErr.Error()), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("Datacenter %s deleted successfully.", datacenterId)), nil
}

