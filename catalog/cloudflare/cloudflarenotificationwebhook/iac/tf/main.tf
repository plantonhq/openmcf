# A webhook destination for notification policies. Cloudflare infers the
# destination type (slack, datadog, generic, ...) from the URL and reports
# it as an output. A plain CRUD resource (real create/update/delete; only
# the account forces replacement).
#
# `secret` is WRITE-ONLY at the API: sent on create/update, never returned
# by any read -- Cloudflare cannot drift-detect it and an imported webhook
# has no secret in state.
resource "cloudflare_notification_policy_webhooks" "main" {
  account_id = var.spec.account_id
  name       = var.spec.name
  url        = var.spec.url

  secret = try(var.spec.secret, "") != "" ? var.spec.secret : null
}
