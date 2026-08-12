# Create the Azure Event Grid namespace -- the capacity-scaled hub of
# the newer Event Grid: it hosts CloudEvents namespace topics and an
# optional MQTT broker behind one set of regional endpoints. "Standard"
# is the SKU's only legal value at v5 (deliberately not part of the
# spec); capacity is throughput units, updatable in place.
resource "azurerm_eventgrid_namespace" "main" {
  name                = var.spec.name
  resource_group_name = var.spec.resource_group
  location            = var.spec.region
  sku                 = "Standard"

  # Platform default 1 TU -- always sent, so the rendered plan states
  # the provisioned throughput.
  capacity = coalesce(var.spec.capacity, 1)

  # Platform default true, mapped to the provider's Enabled/Disabled
  # tokens -- always sent (mirrors Azure's own default).
  public_network_access = coalesce(var.spec.public_network_access_enabled, true) ? "Enabled" : "Disabled"

  # "Allow" is Azure's only legal rule action on this resource at v5,
  # sent explicitly. Rules only take effect while public network access
  # is enabled.
  dynamic "inbound_ip_rule" {
    for_each = var.spec.inbound_ip_rules
    content {
      ip_mask = inbound_ip_rule.value
      action  = "Allow"
    }
  }

  dynamic "identity" {
    for_each = var.spec.identity != null ? [var.spec.identity] : []
    content {
      type         = local.identity_type_map[identity.value.type]
      identity_ids = length(identity.value.identity_ids) > 0 ? identity.value.identity_ids : null
    }
  }

  # The MQTT broker block. Presence is the enable switch (the provider
  # hardcodes state Enabled when the block is sent) and the WHOLE block
  # is create-only -- changing it replaces the namespace. The session
  # dials carry platform defaults (1 session, 1 hour), mirroring the
  # provider's own schema defaults; the route topic is sent only when
  # set.
  dynamic "topic_spaces_configuration" {
    for_each = var.spec.topic_spaces_configuration != null ? [var.spec.topic_spaces_configuration] : []
    content {
      alternative_authentication_name_source          = topic_spaces_configuration.value.alternative_authentication_name_sources
      maximum_client_sessions_per_authentication_name = coalesce(topic_spaces_configuration.value.maximum_client_sessions_per_authentication_name, 1)
      maximum_session_expiry_in_hours                 = coalesce(topic_spaces_configuration.value.maximum_session_expiry_in_hours, 1)
      route_topic_id                                  = topic_spaces_configuration.value.route_topic_id != "" ? topic_spaces_configuration.value.route_topic_id : null

      dynamic "dynamic_routing_enrichment" {
        for_each = topic_spaces_configuration.value.dynamic_routing_enrichments
        content {
          key   = dynamic_routing_enrichment.value.key
          value = dynamic_routing_enrichment.value.value
        }
      }

      # The provider pins every static enrichment's value type to
      # String -- nothing to send beyond the pair.
      dynamic "static_routing_enrichment" {
        for_each = topic_spaces_configuration.value.static_routing_enrichments
        content {
          key   = static_routing_enrichment.value.key
          value = static_routing_enrichment.value.value
        }
      }
    }
  }

  tags = local.final_tags
}
