# Zero Trust MCP server registration: the upstream tool server whose prompts
# and tools Cloudflare proxies, audits, and gates behind Access. Identity is
# user-supplied: server_id (the provider's id), hostname, and auth_type all
# force replacement.
#
# The two credential arguments (auth_credentials, client_secret) are
# WRITE-ONLY at Cloudflare -- stored encrypted, never returned by any read --
# so out-of-band rotation is invisible to IaC; rotate by changing the value
# here. Both are declared sensitive so they never print in plans.
resource "cloudflare_zero_trust_access_ai_controls_mcp_server" "main" {
  account_id = var.spec.account_id
  id         = var.spec.server_id
  name       = var.spec.name
  hostname   = var.spec.hostname
  auth_type  = var.spec.auth_type

  auth_credentials = try(var.spec.auth_credentials, "") != "" ? sensitive(var.spec.auth_credentials) : null
  client_secret    = try(var.spec.client_secret, "") != "" ? sensitive(var.spec.client_secret) : null

  description = try(var.spec.description, "") != "" ? var.spec.description : null

  is_shared_oauth_callback_enabled = try(var.spec.is_shared_oauth_callback_enabled, null)
  secure_web_gateway               = try(var.spec.secure_web_gateway, null)

  # Override rows: an omitted enabled is not sent -- Cloudflare's default
  # keeps the prompt/tool available.
  updated_prompts = length(try(var.spec.updated_prompts, [])) > 0 ? [
    for override in var.spec.updated_prompts : {
      name        = override.name
      alias       = try(override.alias, "") != "" ? override.alias : null
      description = try(override.description, "") != "" ? override.description : null
      enabled     = try(override.enabled, null)
    }
  ] : null

  updated_tools = length(try(var.spec.updated_tools, [])) > 0 ? [
    for override in var.spec.updated_tools : {
      name        = override.name
      alias       = try(override.alias, "") != "" ? override.alias : null
      description = try(override.description, "") != "" ? override.description : null
      enabled     = try(override.enabled, null)
    }
  ] : null
}
