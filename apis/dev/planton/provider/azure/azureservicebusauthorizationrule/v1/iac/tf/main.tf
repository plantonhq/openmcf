# One kind, three ARM types: SAS authorization rules exist at namespace,
# queue, and topic scope with byte-identical shapes (name + the
# listen/send/manage trio + six key/connection-string outputs). The spec's
# exactly-one-scope XOR picks the resource; everything else is shared.
#
# After create/delete on a geo-DR-paired Premium namespace, azurerm waits
# for pairing replication to settle -- no module-side ordering is needed.

# Namespace-wide rights: every queue and topic in the namespace.
resource "azurerm_servicebus_namespace_authorization_rule" "namespace_scoped" {
  count = local.is_namespace_scoped ? 1 : 0

  # ForceNew: renaming replaces the rule and regenerates its keys.
  name         = var.spec.rule_name
  namespace_id = var.spec.namespace_id

  listen = local.listen
  send   = local.send
  manage = local.manage
}

# Single-queue rights -- the least-privilege choice for point-to-point
# workloads.
resource "azurerm_servicebus_queue_authorization_rule" "queue_scoped" {
  count = local.is_queue_scoped ? 1 : 0

  name     = var.spec.rule_name
  queue_id = var.spec.queue_id

  listen = local.listen
  send   = local.send
  manage = local.manage
}

# Single-topic rights (sending to the topic, receiving through its
# subscriptions).
resource "azurerm_servicebus_topic_authorization_rule" "topic_scoped" {
  count = local.is_topic_scoped ? 1 : 0

  name     = var.spec.rule_name
  topic_id = var.spec.topic_id

  listen = local.listen
  send   = local.send
  manage = local.manage
}
