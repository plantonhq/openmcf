locals {
  resource_name = coalesce(try(var.metadata.name, null), "cloudflare-email-routing-zone")

  catch_all = try(var.spec.catch_all, null)

  # Map each typed catch-all action onto the provider's generic {type, value[]}:
  # forward -> the destination addresses; worker -> the single script name;
  # drop -> no values.
  catch_all_actions = local.catch_all == null ? [] : [
    for a in local.catch_all.actions : {
      type = a.type
      value = a.type == "forward" ? a.forward_to : (
        a.type == "worker" ? [a.worker] : []
      )
    }
  ]
}
