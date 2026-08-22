# Cloudflare waiting room: a virtual queue on a host+path, plus the room's
# bypass rules. Cloudflare models the rules as a separate per-room list with
# full-replacement updates -- the rules resource below owns that WHOLE list, so
# the manifest's bypass_rules are the room's entire rule set (rules added
# outside the manifest are overwritten on apply and cleared on destroy).
#
# Advanced-subscription fields (additional_routes, custom_page_html,
# disable_session_renewal, json_response_enabled, non-default queueing_method,
# infinite_queue, non-off turnstile modes) fail at the API on plans without the
# add-on -- the entitlement wall is Cloudflare's.
resource "cloudflare_waiting_room" "main" {
  zone_id = var.spec.zone_id

  name                 = var.spec.name
  host                 = var.spec.host
  path                 = var.spec.path
  new_users_per_minute = var.spec.new_users_per_minute
  total_active_users   = var.spec.total_active_users

  session_duration     = var.spec.session_duration
  suspended            = var.spec.suspended
  queue_all            = var.spec.queue_all
  queueing_method      = var.spec.queueing_method
  queueing_status_code = var.spec.queueing_status_code

  cookie_attributes = var.spec.cookie_attributes != null ? {
    samesite = var.spec.cookie_attributes.samesite
    secure   = var.spec.cookie_attributes.secure
  } : null

  cookie_suffix             = var.spec.cookie_suffix != "" ? var.spec.cookie_suffix : null
  custom_page_html          = var.spec.custom_page_html != "" ? var.spec.custom_page_html : null
  default_template_language = var.spec.default_template_language
  description               = var.spec.description != "" ? var.spec.description : null
  disable_session_renewal   = var.spec.disable_session_renewal
  json_response_enabled     = var.spec.json_response_enabled

  additional_routes = length(var.spec.additional_routes) > 0 ? [
    for route in var.spec.additional_routes : {
      host = route.host
      path = route.path
    }
  ] : null

  enabled_origin_commands = length(var.spec.enabled_origin_commands) > 0 ? var.spec.enabled_origin_commands : null

  turnstile_action = var.spec.turnstile_action
  turnstile_mode   = var.spec.turnstile_mode
}

# The room's bypass rules. The action is fixed to bypass_waiting_room (the only
# action Cloudflare supports on waiting-room rules) -- the module supplies it so
# manifests never repeat a constant. enabled defaults to TRUE upstream and here.
resource "cloudflare_waiting_room_rules" "main" {
  count = length(var.spec.bypass_rules) > 0 ? 1 : 0

  zone_id         = var.spec.zone_id
  waiting_room_id = cloudflare_waiting_room.main.id

  rules = [
    for rule in var.spec.bypass_rules : {
      action      = "bypass_waiting_room"
      expression  = rule.expression
      description = rule.description != "" ? rule.description : null
      enabled     = rule.enabled != null ? rule.enabled : true
    }
  ]
}
