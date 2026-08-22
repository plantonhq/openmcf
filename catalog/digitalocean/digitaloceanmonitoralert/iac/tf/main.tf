# DigitalOcean Monitor Alert
#
# Provisions an alert policy on DigitalOcean's built-in metrics for
# Droplets, load balancers, and managed database clusters, modeling the
# complete digitalocean_monitor_alert resource surface. Every field updates
# in place. The policy's tags select tagged Droplets as targets (they are
# alert TARGETING, not resource labels, so no Planton labels are merged
# in).

resource "digitalocean_monitor_alert" "alert" {
  description = var.spec.description
  type        = var.spec.metric_type
  compare     = var.spec.compare

  # DigitalOcean stores the threshold as a 32-bit float; more than 7
  # significant digits are truncated server-side.
  value  = var.spec.value
  window = var.spec.window

  # Unset (null) defers to the provider's default, enabled; the provider
  # sends the value as a pointer, so an explicit false is transmitted.
  enabled = var.spec.enabled

  # An omitted entities set is simply not sent; DigitalOcean then resolves
  # targets from tags.
  entities = length(local.entities) > 0 ? local.entities : null
  tags     = length(var.spec.tags) > 0 ? var.spec.tags : null

  # The provider caps this block at exactly one.
  alerts {
    email = var.spec.alerts.emails
    dynamic "slack" {
      for_each = var.spec.alerts.slack
      content {
        channel = slack.value.channel
        url     = slack.value.url
      }
    }
  }
}
