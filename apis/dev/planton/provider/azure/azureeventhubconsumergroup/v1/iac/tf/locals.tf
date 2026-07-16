locals {
  # Consumer groups carry no Azure tags: ARM does not support tags on
  # Event Hubs entities (hubs/consumer groups/schema groups), so the
  # platform's identity tags live on the parent namespace.

  # azurerm still addresses consumer groups by discrete names (resource
  # group, namespace, event hub) rather than the parent's ARM ID. The
  # module derives those names from the spec's single parent reference,
  # so the spec stays on the ARM-id grain with no redundant fields that
  # could contradict each other.
  #
  # Expected shape:
  #   /subscriptions/{sub}/resourceGroups/{rg}/providers/Microsoft.EventHub/namespaces/{ns}/eventhubs/{hub}
  # The anchored regex fails the plan loudly on a malformed ID instead
  # of computing wrong names. Segment literals are matched in ARM's
  # canonical camelCase.
  event_hub_id_parts = regex(
    "/resourceGroups/(?P<rg>[^/]+)/providers/Microsoft.EventHub/namespaces/(?P<ns>[^/]+)/eventhubs/(?P<hub>[^/]+)$",
    var.spec.event_hub_id
  )

  resource_group_name = local.event_hub_id_parts["rg"]
  namespace_name      = local.event_hub_id_parts["ns"]
  event_hub_name      = local.event_hub_id_parts["hub"]
}
