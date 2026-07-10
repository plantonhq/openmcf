locals {
  # The scope discriminator: exactly one parent ID is set (spec-enforced
  # XOR), and it picks which of the three azurerm resources materializes.
  # Azure models the three scopes as separate ARM types with identical
  # shapes -- one Planton kind dispatching beats three near-duplicate
  # kinds.
  is_namespace_scoped = var.spec.namespace_id != null && var.spec.namespace_id != ""
  is_queue_scoped     = var.spec.queue_id != null && var.spec.queue_id != ""
  is_topic_scoped     = var.spec.topic_id != null && var.spec.topic_id != ""

  # The rights trio -- azurerm defaults all three to false and rejects
  # an all-false rule at plan time (front-loaded as spec CELs here).
  listen = var.spec.listen != null ? var.spec.listen : false
  send   = var.spec.send != null ? var.spec.send : false
  manage = var.spec.manage != null ? var.spec.manage : false

  # Exactly one of the three resources exists; coalescing PER ATTRIBUTE
  # (never the whole resource object) yields the rule's face for the
  # outputs regardless of scope. Attribute-level splicing matters:
  # coalescing whole resources would taint non-sensitive outputs (id,
  # name) with the key attributes' sensitivity and fail the plan.
  rule_id = coalesce(
    try(azurerm_servicebus_namespace_authorization_rule.namespace_scoped[0].id, null),
    try(azurerm_servicebus_queue_authorization_rule.queue_scoped[0].id, null),
    try(azurerm_servicebus_topic_authorization_rule.topic_scoped[0].id, null),
  )
  rule_name = coalesce(
    try(azurerm_servicebus_namespace_authorization_rule.namespace_scoped[0].name, null),
    try(azurerm_servicebus_queue_authorization_rule.queue_scoped[0].name, null),
    try(azurerm_servicebus_topic_authorization_rule.topic_scoped[0].name, null),
  )
  primary_key = coalesce(
    try(azurerm_servicebus_namespace_authorization_rule.namespace_scoped[0].primary_key, null),
    try(azurerm_servicebus_queue_authorization_rule.queue_scoped[0].primary_key, null),
    try(azurerm_servicebus_topic_authorization_rule.topic_scoped[0].primary_key, null),
  )
  secondary_key = coalesce(
    try(azurerm_servicebus_namespace_authorization_rule.namespace_scoped[0].secondary_key, null),
    try(azurerm_servicebus_queue_authorization_rule.queue_scoped[0].secondary_key, null),
    try(azurerm_servicebus_topic_authorization_rule.topic_scoped[0].secondary_key, null),
  )
  primary_connection_string = coalesce(
    try(azurerm_servicebus_namespace_authorization_rule.namespace_scoped[0].primary_connection_string, null),
    try(azurerm_servicebus_queue_authorization_rule.queue_scoped[0].primary_connection_string, null),
    try(azurerm_servicebus_topic_authorization_rule.topic_scoped[0].primary_connection_string, null),
  )
  secondary_connection_string = coalesce(
    try(azurerm_servicebus_namespace_authorization_rule.namespace_scoped[0].secondary_connection_string, null),
    try(azurerm_servicebus_queue_authorization_rule.queue_scoped[0].secondary_connection_string, null),
    try(azurerm_servicebus_topic_authorization_rule.topic_scoped[0].secondary_connection_string, null),
  )
  # Alias faces are only populated when the namespace carries a geo-DR
  # pairing -- they may legitimately be empty, so no coalesce (which
  # rejects all-null) but a try-chain that falls back to "".
  primary_connection_string_alias = try(
    azurerm_servicebus_namespace_authorization_rule.namespace_scoped[0].primary_connection_string_alias,
    azurerm_servicebus_queue_authorization_rule.queue_scoped[0].primary_connection_string_alias,
    azurerm_servicebus_topic_authorization_rule.topic_scoped[0].primary_connection_string_alias,
    "",
  )
  secondary_connection_string_alias = try(
    azurerm_servicebus_namespace_authorization_rule.namespace_scoped[0].secondary_connection_string_alias,
    azurerm_servicebus_queue_authorization_rule.queue_scoped[0].secondary_connection_string_alias,
    azurerm_servicebus_topic_authorization_rule.topic_scoped[0].secondary_connection_string_alias,
    "",
  )

  # Authorization rules carry no Azure tags: ARM does not support tags
  # on Service Bus entities, so the platform's identity tags live on the
  # parent namespace.
}
