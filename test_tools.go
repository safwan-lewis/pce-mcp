package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// MCP Request structure
type MCPRequest struct {
	JSONRPC string      `json:"jsonrpc"`
	ID      int         `json:"id"`
	Method  string      `json:"method"`
	Params  interface{} `json:"params,omitempty"`
}

// MCP Response structure
type MCPResponse struct {
	JSONRPC string      `json:"jsonrpc"`
	ID      int         `json:"id"`
	Result  interface{} `json:"result,omitempty"`
	Error   *MCPError   `json:"error,omitempty"`
}

type MCPError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type ToolsListResult struct {
	Tools []Tool `json:"tools"`
}

type Tool struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description,omitempty"`
	InputSchema map[string]interface{} `json:"inputSchema,omitempty"`
}

// Expected tools by category
var expectedTools = map[string][]string{
	"organizations": {
		"list_organizations",
		"get_organization_by_id",
		"get_current_organization",
		"list_organization_audit_logs",
		"create_organization",
		"delete_organization_by_id",
		"list_organization_auth_providers",
		"delete_organization_auth_provider",
		"list_organization_roles",
		"create_organization_role",
		"update_organization_role",
		"list_organization_ai_providers",
		"list_supported_ai_providers",
	},
	"users": {
		"list_users_in_organization_by_id",
		"invalidate_user_sessions_by_id",
		"delete_user_by_id",
	},
	"datacenters": {
		"list_datacenters",
		"get_datacenter_by_id",
		"create_datacenter",
		"update_datacenter",
		"delete_datacenter_by_id",
	},
	"clusters": {
		"list_clusters",
		"get_cluster_by_id",
		"initialize_cluster",
		"update_cluster",
		"get_cluster_join_key",
		"get_cluster_metrics",
		"list_cluster_storage_pools",
		"get_cluster_hardware_by_id",
		"get_cluster_licensing_by_id",
	},
	"nodes": {
		"get_node_by_id",
		"get_current_node",
		"get_node_hardware_by_id",
		"get_node_license_by_id",
		"get_node_storagepools_by_id",
		"get_images",
		"get_node_pcidevices_by_id",
		"get_node_metrics",
		"get_node_logs",
		"get_node_capabilities",
		"wake_node",
		"get_node_console",
		"get_node_network_interfaces",
		"get_node_tasks",
		"get_node_jobs",
		"get_node_ssh_keys",
	},
	"instances": {
		"get_instances_in_node",
		"get_instances_in_cluster",
		"deploy_instance",
		"get_instance_by_id",
		"update_instance",
		"delete_instance",
		"restore_instance",
		"get_instance_metrics",
		"get_instance_console",
		"backup_instance",
		"power_instance",
	},
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run test_tools.go <path-to-pce-mcp-binary>")
		fmt.Println("Example: go run test_tools.go ./pce-mcp.exe")
		os.Exit(1)
	}

	binaryPath := os.Args[1]
	baseURL := "https://demo.pextra.cloud"
	
	// Start the server process
	cmd := exec.Command(binaryPath, "serve", "--base-url", baseURL, "--sse-addr=", "--http-addr=")
	
	// Setup stdio pipes
	stdin, err := cmd.StdinPipe()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating stdin pipe: %v\n", err)
		os.Exit(1)
	}
	defer stdin.Close()

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating stdout pipe: %v\n", err)
		os.Exit(1)
	}
	defer stdout.Close()

	cmd.Stderr = os.Stderr

	// Start the server
	fmt.Printf("Starting MCP server: %s\n", binaryPath)
	if err := cmd.Start(); err != nil {
		fmt.Fprintf(os.Stderr, "Error starting server: %v\n", err)
		os.Exit(1)
	}
	defer cmd.Process.Kill()

	// Wait a bit for server to initialize
	// In a real implementation, you'd wait for initialize response

	// Send initialize request
	initReq := MCPRequest{
		JSONRPC: "2.0",
		ID:      1,
		Method:  "initialize",
		Params: map[string]interface{}{
			"protocolVersion": "2024-11-05",
			"capabilities":    map[string]interface{}{},
			"clientInfo": map[string]interface{}{
				"name":    "test-tool-registration",
				"version": "1.0.0",
			},
		},
	}

	// Send tools/list request
	toolsReq := MCPRequest{
		JSONRPC: "2.0",
		ID:      2,
		Method:  "tools/list",
	}

	// Send requests
	requests := []MCPRequest{initReq, toolsReq}
	
	scanner := bufio.NewScanner(stdout)
	responseChan := make(chan MCPResponse, 10)
	
	// Start reading responses
	go func() {
		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			if line == "" {
				continue
			}
			
			var resp MCPResponse
			if err := json.Unmarshal([]byte(line), &resp); err != nil {
				// Skip invalid JSON
				continue
			}
			responseChan <- resp
		}
	}()

	// Send requests
	encoder := json.NewEncoder(stdin)
	for _, req := range requests {
		if err := encoder.Encode(req); err != nil {
			fmt.Fprintf(os.Stderr, "Error sending request: %v\n", err)
			os.Exit(1)
		}
	}

	// Wait for and process responses
	var toolsList *ToolsListResult
	responseCount := 0
	
	for responseCount < 2 {
		select {
		case resp := <-responseChan:
			responseCount++
			
			if resp.Error != nil {
				fmt.Printf("Received error response (ID %d): %s\n", resp.ID, resp.Error.Message)
				continue
			}
			
			if resp.ID == 2 && resp.Result != nil {
				// This should be tools/list response
				resultBytes, err := json.Marshal(resp.Result)
				if err != nil {
					fmt.Fprintf(os.Stderr, "Error marshaling result: %v\n", err)
					continue
				}
				
				var tools ToolsListResult
				if err := json.Unmarshal(resultBytes, &tools); err != nil {
					fmt.Fprintf(os.Stderr, "Error unmarshaling tools list: %v\n", err)
					continue
				}
				
				toolsList = &tools
			}
		}
	}

	if toolsList == nil {
		fmt.Fprintf(os.Stderr, "Error: Did not receive tools list response\n")
		os.Exit(1)
	}

	// Verify tools
	fmt.Printf("\n=== Tool Registration Test ===\n\n")
	fmt.Printf("Total tools registered: %d\n", len(toolsList.Tools))
	
	// Build map of registered tools
	registeredTools := make(map[string]bool)
	for _, tool := range toolsList.Tools {
		registeredTools[tool.Name] = true
	}

	// Check all expected tools
	allExpected := []string{}
	for _, tools := range expectedTools {
		allExpected = append(allExpected, tools...)
	}

	fmt.Printf("Expected tools: %d\n\n", len(allExpected))

	missing := []string{}
	extra := []string{}

	for _, expected := range allExpected {
		if !registeredTools[expected] {
			missing = append(missing, expected)
		}
	}

	for _, tool := range toolsList.Tools {
		found := false
		for _, expected := range allExpected {
			if tool.Name == expected {
				found = true
				break
			}
		}
		if !found {
			extra = append(extra, tool.Name)
		}
	}

	// Report results
	if len(missing) == 0 && len(extra) == 0 {
		fmt.Printf("✅ SUCCESS: All %d expected tools are registered!\n\n", len(allExpected))
	} else {
		if len(missing) > 0 {
			fmt.Printf("❌ Missing tools (%d):\n", len(missing))
			for _, tool := range missing {
				fmt.Printf("  - %s\n", tool)
			}
			fmt.Println()
		}
		
		if len(extra) > 0 {
			fmt.Printf("⚠️  Extra tools found (%d):\n", len(extra))
			for _, tool := range extra {
				fmt.Printf("  - %s\n", tool)
			}
			fmt.Println()
		}
	}

	// Show tools by category
	fmt.Println("Tools by category:")
	totalRegistered := 0
	for category, tools := range expectedTools {
		count := 0
		for _, tool := range tools {
			if registeredTools[tool] {
				count++
			}
		}
		totalRegistered += count
		fmt.Printf("  %s: %d/%d tools\n", category, count, len(tools))
	}

	if len(missing) > 0 {
		os.Exit(1)
	}
}

