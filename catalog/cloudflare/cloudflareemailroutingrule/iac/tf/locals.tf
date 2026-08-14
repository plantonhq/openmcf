locals {
  resource_name = coalesce(try(var.metadata.name, null), "cloudflare-email-routing-rule")

  # Map each typed action onto the provider's generic {type, value[]}:
  # forward -> the destination addresses; worker -> the single script name;
  # drop -> no values.
  actions = [
    for a in var.spec.actions : {
      type = a.type
      value = a.type == "forward" ? a.forward_to : (
        a.type == "worker" ? [a.worker] : []
      )
    }
  ]
}
