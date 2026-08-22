# DigitalOcean Uptime Check
#
# Provisions an availability/latency probe on an EXTERNAL endpoint plus its
# alert rules, modeling the digitalocean_uptime_check resource surface and
# composing one digitalocean_uptime_alert per spec alert row. The rules are
# composed here because they cannot exist without the check, and because
# the standalone alert resource leaves the parent check id mutable upstream
# -- re-pointing it orphans the alert on the old check, a corruption class
# this composition makes unrepresentable.

resource "digitalocean_uptime_check" "check" {
  name   = var.spec.check_name
  target = var.spec.target

  # Unset defers to the provider's default, https.
  type = var.spec.type != "" ? var.spec.type : null

  # Always declared (spec-required): the provider never reconciles a
  # DigitalOcean-defaulted region set, so an omitted value would leave
  # every subsequent plan trying to remove what the API chose.
  regions = var.spec.regions

  # Unset (null) defers to the provider's default, enabled.
  enabled = var.spec.enabled
}

resource "digitalocean_uptime_alert" "alerts" {
  for_each = local.alerts

  check_id = digitalocean_uptime_check.check.id
  name     = each.value.alert_name
  type     = each.value.type

  # Milliseconds for latency, days before expiry for ssl_expiry; down and
  # down_global carry no threshold (an unset value is sent as the API's
  # accepted zero).
  threshold = each.value.threshold

  comparison = each.value.comparison != "" ? each.value.comparison : null
  period     = each.value.period != "" ? each.value.period : null

  notifications {
    email = each.value.notifications.emails
    dynamic "slack" {
      for_each = each.value.notifications.slack
      content {
        channel = slack.value.channel
        url     = slack.value.url
      }
    }
  }
}
