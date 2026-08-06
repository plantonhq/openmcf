# One kind, two ARM types: SAS authorization rules exist at namespace and
# event hub scope with byte-identical shapes (name + the
# listen/send/manage trio + six key/connection-string outputs). The
# spec's exactly-one-scope XOR picks the resource; everything else is
# shared. Both azurerm resources still use legacy name addressing, so the
# locals parse the parent names from the resolved ARM id.

# Namespace-wide rights: every hub in the namespace.
resource "azurerm_eventhub_namespace_authorization_rule" "namespace_scoped" {
  count = local.is_namespace_scoped ? 1 : 0

  # ForceNew: renaming replaces the rule and regenerates its keys.
  name                = var.spec.rule_name
  namespace_name      = local.namespace_scope.ns
  resource_group_name = local.namespace_scope.rg

  listen = local.listen
  send   = local.send
  manage = local.manage
}

# Single-hub rights -- the least-privilege choice for per-stream
# producers and consumers.
resource "azurerm_eventhub_authorization_rule" "hub_scoped" {
  count = local.is_hub_scoped ? 1 : 0

  name                = var.spec.rule_name
  namespace_name      = local.hub_scope.ns
  eventhub_name       = local.hub_scope.hub
  resource_group_name = local.hub_scope.rg

  listen = local.listen
  send   = local.send
  manage = local.manage
}
