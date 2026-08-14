# Create the Azure Event Grid domain -- ONE publishing endpoint and one
# pair of access keys serving many event streams (domain topics), the
# multi-tenant pattern. The domain's name becomes a PUBLIC DNS hostname
# ({name}.{region}.eventgrid.azure.net), unique across all Azure
# customers in the region. Free at rest; billing is per operation.
#
# Domain topics come from two postures: auto-managed (the defaults --
# Azure creates a topic when its first subscription arrives and deletes
# it with the last one) or pinned (both flags false -- topics exist
# only as declared AzureEventgridDomainTopic resources).
resource "azurerm_eventgrid_domain" "main" {
  name                = var.spec.name
  resource_group_name = var.spec.resource_group
  location            = var.spec.region

  # The provider defaults to EventGridSchema; the platform sends its
  # own default explicitly so the rendered plan states the schema --
  # one schema for every topic in the domain. Create-only: changing it
  # replaces the domain.
  input_schema = var.spec.input_schema

  # Custom-schema envelope mappings -- sent only when they carry at
  # least one field (the built-in schemas need no mapping). Empty
  # strings inside a sent block are dropped by the provider's expand
  # (it maps only non-empty source fields).
  dynamic "input_mapping_fields" {
    for_each = local.input_mapping_fields_set ? [var.spec.input_mapping_fields] : []
    content {
      id           = input_mapping_fields.value.id != "" ? input_mapping_fields.value.id : null
      topic        = input_mapping_fields.value.topic != "" ? input_mapping_fields.value.topic : null
      event_time   = input_mapping_fields.value.event_time != "" ? input_mapping_fields.value.event_time : null
      event_type   = input_mapping_fields.value.event_type != "" ? input_mapping_fields.value.event_type : null
      subject      = input_mapping_fields.value.subject != "" ? input_mapping_fields.value.subject : null
      data_version = input_mapping_fields.value.data_version != "" ? input_mapping_fields.value.data_version : null
    }
  }

  dynamic "input_mapping_default_values" {
    for_each = local.input_mapping_default_values_set ? [var.spec.input_mapping_default_values] : []
    content {
      event_type   = input_mapping_default_values.value.event_type != "" ? input_mapping_default_values.value.event_type : null
      subject      = input_mapping_default_values.value.subject != "" ? input_mapping_default_values.value.subject : null
      data_version = input_mapping_default_values.value.data_version != "" ? input_mapping_default_values.value.data_version : null
    }
  }

  # Domain-topic lifecycle flags -- always sent (platform defaults
  # mirror Azure's auto-managed posture).
  auto_create_topic_with_first_subscription = var.spec.auto_create_topic_with_first_subscription
  auto_delete_topic_with_last_subscription  = var.spec.auto_delete_topic_with_last_subscription

  # Always sent (platform defaults mirror Azure's). NOTE the provider
  # seam: on create it renders false as ARM's "Disabled", and local
  # auth inverts to ARM's disableLocalAuth.
  public_network_access_enabled = var.spec.public_network_access_enabled
  local_auth_enabled            = var.spec.local_auth_enabled

  # The provider clears rules on update by sending an EMPTY list -- an
  # attribute-mode list mirrors that exactly (rule deletions propagate).
  # "Allow" is Azure's only legal action on this resource at v5, sent
  # explicitly.
  inbound_ip_rule = [
    for ip_mask in var.spec.inbound_ip_rules : {
      ip_mask = ip_mask
      action  = "Allow"
    }
  ]

  dynamic "identity" {
    for_each = var.spec.identity != null ? [var.spec.identity] : []
    content {
      type         = local.identity_type_map[identity.value.type]
      identity_ids = length(identity.value.identity_ids) > 0 ? identity.value.identity_ids : null
    }
  }

  tags = local.final_tags
}
