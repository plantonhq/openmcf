# The subscription, addressed by the parent topic's ARM ID (azurerm's
# v4 child-addressing grain). Consumer semantics -- locks, delivery
# counts, sessions, dead-lettering -- live here, not on the topic.
resource "azurerm_servicebus_subscription" "main" {
  # ForceNew: renaming replaces the subscription and resets its read
  # position -- undelivered messages in the old subscription are lost.
  name     = var.spec.subscription_name
  topic_id = var.spec.topic_id

  # Required: Azure has no server default for a subscription's delivery
  # tolerance.
  max_delivery_count = var.spec.max_delivery_count

  # Lifecycle dials -- unset leaves Azure's defaults in place (lock
  # PT1M, TTL inherited from the topic, never auto-delete).
  lock_duration       = var.spec.lock_duration
  default_message_ttl = var.spec.default_message_ttl
  auto_delete_on_idle = var.spec.auto_delete_on_idle

  dead_lettering_on_message_expiration      = var.spec.dead_lettering_on_message_expiration
  dead_lettering_on_filter_evaluation_error = var.spec.dead_lettering_on_filter_evaluation_error

  # ForceNew: the session model is fixed at creation.
  requires_session = var.spec.requires_session

  batched_operations_enabled = var.spec.batched_operations_enabled

  # Routing chains: targets are entity NAMES in the same namespace (not
  # ARM ids). The classic fan-out-then-collect pattern: subscriptions
  # filter, forwarding funnels matches into a work queue.
  forward_to                        = var.spec.forward_to
  forward_dead_lettered_messages_to = var.spec.forward_dead_lettered_messages_to

  status = local.status

  # The JMS 2.0 client-affine binding. Azure stores the entity as
  # {name}${client_id}$D internally; the provider round-trips the
  # user-facing name.
  client_scoped_subscription_enabled = var.spec.client_scoped_subscription != null

  dynamic "client_scoped_subscription" {
    for_each = var.spec.client_scoped_subscription != null ? [var.spec.client_scoped_subscription] : []
    content {
      client_id                               = client_scoped_subscription.value.client_id
      is_client_scoped_subscription_shareable = client_scoped_subscription.value.shareable
    }
  }
}

# Folded filter rules -- a rule has no life outside its subscription and
# nothing references one, so they ship as part of the subscription
# document. OR semantics: a message is delivered once if ANY rule
# matches -- ALONGSIDE Azure's auto-created "$Default" catch-all, which
# cannot be declared here (the provider's import check refuses to adopt
# the service-created rule; the spec reserves the name). Restrictive
# delivery = remove the catch-all once, out-of-band.
resource "azurerm_servicebus_subscription_rule" "rules" {
  for_each = { for r in var.spec.rules : r.rule_name => r }

  name            = each.value.rule_name
  subscription_id = azurerm_servicebus_subscription.main.id

  filter_type = local.filter_type_map[each.value.filter_type]

  # The SQL path: expression required with SqlFilter (spec-enforced
  # XOR). The optional action annotates matched messages before
  # delivery.
  sql_filter = each.value.sql_filter
  action     = each.value.action

  # The correlation path: equality matching on correlation properties --
  # cheaper than SQL at high throughput.
  dynamic "correlation_filter" {
    for_each = each.value.correlation_filter != null ? [each.value.correlation_filter] : []
    content {
      correlation_id      = correlation_filter.value.correlation_id
      message_id          = correlation_filter.value.message_id
      to                  = correlation_filter.value.to
      reply_to            = correlation_filter.value.reply_to
      label               = correlation_filter.value.label
      session_id          = correlation_filter.value.session_id
      reply_to_session_id = correlation_filter.value.reply_to_session_id
      content_type        = correlation_filter.value.content_type
      properties          = length(correlation_filter.value.properties) > 0 ? correlation_filter.value.properties : null
    }
  }
}
