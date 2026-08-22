# A notification policy: one alert type fanned out to email, PagerDuty, and
# webhook destinations, optionally narrowed by filters. A plain CRUD
# resource (real create/update/delete; only the account forces replacement).
#
# The spec flattens each mechanism to its identity (an address or a UUID);
# this module rebuilds the API's object rows. `enabled` is sent only when
# set -- Cloudflare's default is true.
resource "cloudflare_notification_policy" "main" {
  account_id = var.spec.account_id
  name       = var.spec.name
  alert_type = var.spec.alert_type

  description    = try(var.spec.description, "") != "" ? var.spec.description : null
  enabled        = try(var.spec.enabled, null)
  alert_interval = try(var.spec.alert_interval, "") != "" ? var.spec.alert_interval : null

  mechanisms = {
    email = length(try(var.spec.mechanisms.emails, [])) > 0 ? [
      for address in var.spec.mechanisms.emails : { id = address }
    ] : null
    pagerduty = length(try(var.spec.mechanisms.pagerduty_ids, [])) > 0 ? [
      for id in var.spec.mechanisms.pagerduty_ids : { id = id }
    ] : null
    webhooks = length(try(var.spec.mechanisms.webhook_ids, [])) > 0 ? [
      for id in var.spec.mechanisms.webhook_ids : { id = id }
    ] : null
  }

  # Every filter is a list of strings; only declared filters are sent, so
  # the payload holds exactly the fields the alert type reads.
  filters = try(var.spec.filters, null) != null ? {
    actions                         = length(try(var.spec.filters.actions, [])) > 0 ? var.spec.filters.actions : null
    affected_asns                   = length(try(var.spec.filters.affected_asns, [])) > 0 ? var.spec.filters.affected_asns : null
    affected_components             = length(try(var.spec.filters.affected_components, [])) > 0 ? var.spec.filters.affected_components : null
    affected_locations              = length(try(var.spec.filters.affected_locations, [])) > 0 ? var.spec.filters.affected_locations : null
    airport_code                    = length(try(var.spec.filters.airport_code, [])) > 0 ? var.spec.filters.airport_code : null
    alert_trigger_preferences       = length(try(var.spec.filters.alert_trigger_preferences, [])) > 0 ? var.spec.filters.alert_trigger_preferences : null
    alert_trigger_preferences_value = length(try(var.spec.filters.alert_trigger_preferences_value, [])) > 0 ? var.spec.filters.alert_trigger_preferences_value : null
    enabled                         = length(try(var.spec.filters.enabled, [])) > 0 ? var.spec.filters.enabled : null
    environment                     = length(try(var.spec.filters.environment, [])) > 0 ? var.spec.filters.environment : null
    event                           = length(try(var.spec.filters.event, [])) > 0 ? var.spec.filters.event : null
    event_source                    = length(try(var.spec.filters.event_source, [])) > 0 ? var.spec.filters.event_source : null
    event_type                      = length(try(var.spec.filters.event_type, [])) > 0 ? var.spec.filters.event_type : null
    group_by                        = length(try(var.spec.filters.group_by, [])) > 0 ? var.spec.filters.group_by : null
    health_check_id                 = length(try(var.spec.filters.health_check_id, [])) > 0 ? var.spec.filters.health_check_id : null
    incident_impact                 = length(try(var.spec.filters.incident_impact, [])) > 0 ? var.spec.filters.incident_impact : null
    input_id                        = length(try(var.spec.filters.input_id, [])) > 0 ? var.spec.filters.input_id : null
    insight_class                   = length(try(var.spec.filters.insight_class, [])) > 0 ? var.spec.filters.insight_class : null
    limit                           = length(try(var.spec.filters.limit, [])) > 0 ? var.spec.filters.limit : null
    logo_tag                        = length(try(var.spec.filters.logo_tag, [])) > 0 ? var.spec.filters.logo_tag : null
    megabits_per_second             = length(try(var.spec.filters.megabits_per_second, [])) > 0 ? var.spec.filters.megabits_per_second : null
    new_health                      = length(try(var.spec.filters.new_health, [])) > 0 ? var.spec.filters.new_health : null
    new_status                      = length(try(var.spec.filters.new_status, [])) > 0 ? var.spec.filters.new_status : null
    packets_per_second              = length(try(var.spec.filters.packets_per_second, [])) > 0 ? var.spec.filters.packets_per_second : null
    pool_id                         = length(try(var.spec.filters.pool_id, [])) > 0 ? var.spec.filters.pool_id : null
    pop_names                       = length(try(var.spec.filters.pop_names, [])) > 0 ? var.spec.filters.pop_names : null
    product                         = length(try(var.spec.filters.product, [])) > 0 ? var.spec.filters.product : null
    project_id                      = length(try(var.spec.filters.project_id, [])) > 0 ? var.spec.filters.project_id : null
    protocol                        = length(try(var.spec.filters.protocol, [])) > 0 ? var.spec.filters.protocol : null
    query_tag                       = length(try(var.spec.filters.query_tag, [])) > 0 ? var.spec.filters.query_tag : null
    requests_per_second             = length(try(var.spec.filters.requests_per_second, [])) > 0 ? var.spec.filters.requests_per_second : null
    selectors                       = length(try(var.spec.filters.selectors, [])) > 0 ? var.spec.filters.selectors : null
    services                        = length(try(var.spec.filters.services, [])) > 0 ? var.spec.filters.services : null
    slo                             = length(try(var.spec.filters.slo, [])) > 0 ? var.spec.filters.slo : null
    status                          = length(try(var.spec.filters.status, [])) > 0 ? var.spec.filters.status : null
    target_hostname                 = length(try(var.spec.filters.target_hostname, [])) > 0 ? var.spec.filters.target_hostname : null
    target_ip                       = length(try(var.spec.filters.target_ip, [])) > 0 ? var.spec.filters.target_ip : null
    target_zone_name                = length(try(var.spec.filters.target_zone_name, [])) > 0 ? var.spec.filters.target_zone_name : null
    traffic_exclusions              = length(try(var.spec.filters.traffic_exclusions, [])) > 0 ? var.spec.filters.traffic_exclusions : null
    tunnel_id                       = length(try(var.spec.filters.tunnel_id, [])) > 0 ? var.spec.filters.tunnel_id : null
    tunnel_name                     = length(try(var.spec.filters.tunnel_name, [])) > 0 ? var.spec.filters.tunnel_name : null
    type                            = length(try(var.spec.filters.type, [])) > 0 ? var.spec.filters.type : null
    where                           = length(try(var.spec.filters.where, [])) > 0 ? var.spec.filters.where : null
    zones                           = length(try(var.spec.filters.zones, [])) > 0 ? var.spec.filters.zones : null
  } : null
}
