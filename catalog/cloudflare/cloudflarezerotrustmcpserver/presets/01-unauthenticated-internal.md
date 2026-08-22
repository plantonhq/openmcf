# Unauthenticated internal server

The simplest shape: an internal MCP server that trusts its network path (Cloudflare authenticates the USER at the portal; the upstream takes Cloudflare's word for it). The sharp tool (`delete_page`) is disabled at the server level -- nobody can invoke it through Cloudflare, regardless of portal.
