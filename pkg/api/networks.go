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
	"context"
	"net/url"
)

// VSwitch represents a standalone vSwitch on a node
type VSwitch struct {
	ID          string  `json:"id"`
	Type        int     `json:"type"`
	NodeID      string  `json:"node_id"`
	Name        string  `json:"name"`
	Description *string `json:"description"`
}

type ListVSwitchesArg struct {
	NodeID string
}

type ListVSwitchesResponse = []VSwitch

func ListVSwitches(ctx context.Context, c *Client, arg *ListVSwitchesArg) (*ListVSwitchesResponse, *APIError) {
	if arg == nil || arg.NodeID == "" {
		return nil, NewAPIError(400, "node_id is required")
	}

	path := "/v1/networks/vswitches/standalone/"
	query := make(url.Values)
	query.Set("node_id", arg.NodeID)

	var resp ListVSwitchesResponse
	if apiErr := c.Get(ctx, path, query, &resp); apiErr != nil {
		return nil, apiErr
	}
	return &resp, nil
}


