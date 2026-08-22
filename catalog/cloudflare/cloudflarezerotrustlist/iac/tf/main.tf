# Cloudflare Zero Trust list: a reusable, account-scoped set of values
# (domains, IPs, URLs, emails, serials) referenced by Gateway policies and
# device-posture rules. The list TYPE is immutable -- Terraform replaces the
# list (new ID) if it changes. URL-type values are normalized by the API, a
# known upstream drift source at provider v5.23.0 (its own acceptance test
# expects a non-empty plan for URL lists).
resource "cloudflare_zero_trust_list" "main" {
  account_id = var.spec.account_id

  name        = var.spec.name
  type        = var.spec.type
  description = var.spec.description != "" ? var.spec.description : null

  items = local.items
}
