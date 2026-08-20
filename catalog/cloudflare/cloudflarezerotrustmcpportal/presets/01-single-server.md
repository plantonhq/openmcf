# Single-server portal

The starter shape: one registered MCP server published on one Access-protected hostname. The `server_id` here is a literal; in a full IaC estate, reference the `CloudflareZeroTrustMcpServer` resource's output instead so the dependency graph carries the edge.
