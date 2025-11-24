# Server Testing Guide

## Prerequisites

1. ✅ Go installed and working (Go 1.25.4)
2. ✅ Build successful (all packages compile)
3. ✅ Binary created (`pce-mcp.exe`)
4. ✅ **PCE Instance URL** - Configured: `https://demo.pextra.cloud`

## Server Configuration

The server requires a `--base-url` parameter pointing to your PCE instance.

**Required:**
- `--base-url`: PCE base URL (e.g., `https://192.168.1.27:5007`)

**Optional:**
- `--tls-skip-verify`: Skip TLS verification (useful for self-signed certs)
- `--timeout`: Request timeout in seconds (default: 10)
- `--sse-addr`: SSE server address (default: `:2222`)
- `--http-addr`: HTTP server address (default: `:2223`)
- `--disable-stdio`: Disable stdio transport

## Testing Steps

### Step 1: Verify Server Starts ✅ COMPLETED

Server successfully started and registered all tools:

```powershell
.\pce-mcp.exe serve --base-url "https://demo.pextra.cloud" --sse-addr= --http-addr=
```

**Result:**
- ✅ Server started successfully
- ✅ Base URL validated: `https://demo.pextra.cloud`
- ✅ All 57 tools registered
- ✅ Stdio server ready for connections

**Actual output:**
```
Serving stdio
SSE server disabled (empty address)
HTTP server disabled (empty address)
```

### Step 2: Test Tool Discovery ✅ COMPLETED

Tool registration verified using automated test script (`test_tools.go`):

```powershell
go run test_tools.go .\pce-mcp.exe
```

**Test Results:**
- ✅ All 57 tools successfully registered
- ✅ No missing tools
- ✅ No extra/unexpected tools
- ✅ All categories verified (Organizations, Users, Datacenters, Clusters, Nodes, Instances)

### Step 3: Test Tool Execution

Test individual tools by calling them with appropriate parameters using an MCP client.

## Common Test Scenarios

### Test with Self-Signed Certificate

```powershell
go run ./... serve --base-url "https://localhost:5007" --tls-skip-verify
```

### Test with Custom Timeout

```powershell
go run ./... serve --base-url "https://your-pce-instance:5007" --timeout 30
```

### Test with All Transports Disabled (Validation Only)

```powershell
go run ./... serve --base-url "https://your-pce-instance:5007" --sse-addr= --http-addr= --disable-stdio
```

This will validate configuration and exit immediately.

## Using the Binary

Instead of `go run`, you can use the compiled binary:

```powershell
.\pce-mcp.exe serve --base-url "https://demo.pextra.cloud" --sse-addr= --http-addr=
```

**Verified working:** ✅ Server starts successfully with demo instance

## MCP Client Testing

To test with an MCP client (like Claude Desktop):

1. Configure the client to connect to the server
2. Verify all 57 tools are discoverable
3. Test tool parameters are correctly defined
4. Execute tools with valid parameters

## Troubleshooting

### "invalid base-url" error
- Verify the URL format is correct
- Check if the URL is accessible from your machine

### TLS/certificate errors
- Use `--tls-skip-verify` for self-signed certificates
- Or provide `--tls-ca-cert` with path to CA certificate

### Server doesn't start
- Check if ports are already in use (2222, 2223)
- Verify Go is in PATH
- Check for compilation errors

### No tools discovered
- Verify all tools are registered in `internal/server/tools.go`
- Check server logs for registration errors
- Ensure server started successfully

## Expected Tool Count ✅ VERIFIED

The server successfully registers **57 tools total**:
- Organizations: 13/13 tools ✅
- Users: 3/3 tools ✅
- Datacenters: 5/5 tools ✅
- Clusters: 9/9 tools ✅
- Nodes: 16/16 tools ✅
- Instances: 11/11 tools ✅

**Verification:** All 57 tools confirmed registered via automated test script.

## Test Script

A test script (`test_tools.go`) is available to verify tool registration:

```powershell
go run test_tools.go .\pce-mcp.exe
```

The script:
- Starts the MCP server
- Connects via stdio protocol
- Retrieves tools list
- Verifies all expected tools are present
- Reports any missing or extra tools

