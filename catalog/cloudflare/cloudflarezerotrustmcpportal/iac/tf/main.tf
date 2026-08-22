# Zero Trust MCP portal: the Access-protected endpoint that publishes a
# curated collection of registered MCP servers. The portal's id is
# user-supplied and forces replacement; hostname and name update in place.
#
# The servers rows are a SET at the provider (the backend ignores
# declaration order and returns its own canonical order), so reordering
# spec rows never plans a change.
resource "cloudflare_zero_trust_access_ai_controls_mcp_portal" "main" {
  account_id = var.spec.account_id
  id         = var.spec.portal_id
  hostname   = var.spec.hostname
  name       = var.spec.name

  description = try(var.spec.description, "") != "" ? var.spec.description : null

  allow_code_mode    = try(var.spec.allow_code_mode, null)
  secure_web_gateway = try(var.spec.secure_web_gateway, null)

  # Published-server rows: omitted booleans are not sent -- Cloudflare's
  # defaults keep the server enabled with on-behalf authentication.
  servers = length(try(var.spec.servers, [])) > 0 ? [
    for server in var.spec.servers : {
      server_id        = server.server_id
      default_disabled = try(server.default_disabled, null)
      on_behalf        = try(server.on_behalf, null)
      updated_prompts = length(try(server.updated_prompts, [])) > 0 ? [
        for override in server.updated_prompts : {
          name        = override.name
          alias       = try(override.alias, "") != "" ? override.alias : null
          description = try(override.description, "") != "" ? override.description : null
          enabled     = try(override.enabled, null)
        }
      ] : null
      updated_tools = length(try(server.updated_tools, [])) > 0 ? [
        for override in server.updated_tools : {
          name        = override.name
          alias       = try(override.alias, "") != "" ? override.alias : null
          description = try(override.description, "") != "" ? override.description : null
          enabled     = try(override.enabled, null)
        }
      ] : null
    }
  ] : null
}
