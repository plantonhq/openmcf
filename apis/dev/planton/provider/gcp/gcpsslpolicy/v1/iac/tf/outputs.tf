# Exactly one of the two resources exists (see the count guards in main.tf),
# so each output selects whichever branch was created with one(concat(...)).

# The self-link — the value a target HTTPS (or SSL) proxy references in its
# ssl_policy field; the composition handle that hardens TLS at the load
# balancer.
output "self_link" {
  description = "Self-link URI of the SSL policy"
  value       = one(concat(google_compute_ssl_policy.this[*].self_link, google_compute_region_ssl_policy.this[*].self_link))
}

# The name as it exists in GCP.
output "ssl_policy_name" {
  description = "Name of the SSL policy"
  value       = one(concat(google_compute_ssl_policy.this[*].name, google_compute_region_ssl_policy.this[*].name))
}

# The cipher suites the policy actually enables, computed by GCP from the
# profile (or copied from custom_features on CUSTOM) — the list a compliance
# auditor asks for.
output "enabled_features" {
  description = "Cipher suites enabled by the policy as computed by GCP"
  value       = one(concat(google_compute_ssl_policy.this[*].enabled_features, google_compute_region_ssl_policy.this[*].enabled_features))
}

# Region of a regional SSL policy; empty for a global one, so downstream
# composition can confirm scope compatibility.
output "region" {
  description = "Region of the SSL policy (empty for global)"
  value       = local.is_regional ? var.spec.region : ""
}
