# pce-mcp

A production-ready [Model Context Protocol (MCP)](https://modelcontextprotocol.io/) server for [Pextra CloudEnvironment®](https://pextra.cloud/cloudenvironment).

Manage your Pextra CloudEnvironment® resources using MCP-compatible clients and tools.

## Install

1. Download the latest release binary from the [Releases](https://github.com/PextraCloud/pce-mcp/releases) page for your platform.

2. Follow the [Configuration](#configuration) and [Usage](#usage) sections below to run the server.

3. Connect your MCP-compatible clients and tools to the server and begin prompting!

This produces the `pce-mcp` binary in the repository root when using `go build`.

## Configuration

You can configure the server using flags or environment variables. Flags override environment variables.

Flags:

-   `--base-url` (required): Pextra CloudEnvironment® base URL (e.g., `https://192.168.1.27:5007`).
-   `--sse-addr` (default `:2222`): SSE server listen address; set to empty string to disable.
-   `--http-addr` (default `:2223`): HTTP server listen address; set to empty string to disable.
-   `--disable-stdio` (default `false`): Disable the stdio server.
-   `--tls-ca-cert` (default `""`): Path to a custom CA certificate file for the PCE API client. If set, TLS verification will use this CA instead of the system CAs. Mutually exclusive with `--tls-skip-verify`.
-   `--tls-skip-verify` (default `false`): Disable TLS verification for the PCE API client. This exposes you to man-in-the-middle attacks and is not recommended for production use. Mutually exclusive with `--tls-ca-cert`.
-   `--timeout` (default `10`): PCE API request timeout in seconds.
-   `--headers` (default `""`): Custom HTTP headers to include in the PCE API client requests, formatted as a key=value pairs, can be specified multiple times. Example: `--headers "Authorization=Basic xxx" --headers "X-Custom-Header=Value"`.
-   `--safety-level` (default `read-only`): Safety level for tool operations to prevent dangerous operations when using with AI agents. Options:
    - `read-only`: Only enables read operations (Get, List tools). Safe for monitoring and exploration.
    - `update`: Enables read operations plus non-destructive write operations (Create, Update, Deploy, Backup, Power). Allows configuration changes but prevents data loss.
    - `delete`: Enables all operations including destructive ones (Delete, Restore). Use with caution, especially with AI agents.

Environment variables (fallbacks if corresponding flag is not set):

-   `BASE_URL`
-   `SSE_ADDR`
-   `HTTP_ADDR`
-   `DISABLE_STDIO` (e.g., `true`/`false`)
-   `TLS_CA_CERT` (file path)
-   `TLS_SKIP_VERIFY` (e.g., `true`/`false`)
-   `TIMEOUT` (integer seconds)
-   `SAFETY_LEVEL` (one of: `read-only`, `update`, `delete`)

## Usage

Run the server:

```bash
./pce-mcp serve --base-url "https://localhost:5007"
```

Disable transports as needed:

```bash
# Disable SSE
./pce-mcp serve --base-url "https://localhost:5007" --sse-addr=

# Disable HTTP
./pce-mcp serve --base-url "https://localhost:5007" --http-addr=

# Disable stdio server
./pce-mcp serve --base-url "https://localhost:5007" --disable-stdio

# Disable two transports (only stdio server enabled)
./pce-mcp serve --base-url "https://localhost:5007" --sse-addr= --http-addr=
```

Use environment variables as fallbacks:

```bash
export BASE_URL="https://localhost:5007"
export SSE_ADDR=":3000"
export HTTP_ADDR=":3001"
export DISABLE_STDIO=true
export TLS_SKIP_VERIFY=true
export TIMEOUT=15
export SAFETY_LEVEL=read-only
./pce-mcp serve
```

### Safety Levels

When using the server with AI agents, it's recommended to use a restricted safety level to prevent accidental destructive operations:

```bash
# Read-only mode - safe for exploration and monitoring
./pce-mcp serve --base-url "https://localhost:5007" --safety-level read-only

# Update mode - allows configuration changes but prevents deletions
./pce-mcp serve --base-url "https://localhost:5007" --safety-level update

# Delete mode - full access (use with caution)
./pce-mcp serve --base-url "https://localhost:5007" --safety-level delete
```

**Tool Categories:**
- **Read-only tools** (available in all modes): List and Get operations for querying resources
- **Update tools** (available in `update` and `delete` modes): Create, Update, Deploy, Backup, Power operations
- **Destructive tools** (only in `delete` mode): Delete and Restore operations that can cause data loss

## Development

Build:

```bash
go build ./...
```

Run:

```bash
go run ./... serve --base-url "https://localhost:5007"
```

### Tests

Integration tests (Linux only, optional tools required; some tests may need root):

```bash
go test -v -tags=integration -cover ./...
```

Coverage report:

```bash
go test -v -tags=integration -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

## License

This repository is licensed under the [Apache License 2.0](./LICENSE).
