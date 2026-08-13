# Enable the Compute Engine API so a fresh project can host the SSL policy.
# disable_on_destroy is false: tearing down one policy must never disable the
# API for everything else in the project.
resource "google_project_service" "compute_api" {
  project = local.project_id
  service = "compute.googleapis.com"

  disable_dependent_services = true
  disable_on_destroy         = false
}

# A Compute Engine SSL policy — the control for which TLS versions and cipher
# suites a load balancer accepts from clients. Target HTTPS (and SSL) proxies
# reference the policy's self_link; without one, GCP's permissive default
# applies (min TLS 1.0, COMPATIBLE ciphers).
#
# One kind, two provider resources: GCP models global and regional SSL
# policies as separate API collections with an identical surface, so this
# module creates google_compute_ssl_policy when spec.region is empty and
# google_compute_region_ssl_policy when it is set. Exactly one of the two
# count guards is 1 — never both.
#
# name, project, and description are immutable (ForceNew): changing any of
# them destroys and recreates the policy, briefly breaking every proxy that
# references the old self_link. profile, min_tls_version, and custom_features
# all update IN PLACE — a policy is shared configuration, so tightening the
# TLS floor for an entire proxy fleet is a single-resource change that
# applies on the next client handshake.
#
# profile and min_tls_version deliberately fall through to the API defaults
# (COMPATIBLE / TLS_1_0) when unset — hardcoding them here would silently pin
# behavior the provider may evolve.
resource "google_compute_ssl_policy" "this" {
  count = local.is_regional ? 0 : 1

  name        = local.ssl_policy_name
  project     = local.project_id
  description = var.spec.description != "" ? var.spec.description : null

  profile         = local.profile
  min_tls_version = local.min_tls_version

  # Only sent with the CUSTOM profile (the variable validation enforces the
  # pairing); GCP rejects the field on predefined profiles.
  custom_features = length(var.spec.custom_features) > 0 ? var.spec.custom_features : null

  # Post-quantum rollout stance; empty falls through to the API default
  # (DEFAULT — GCP's own timeline).
  post_quantum_key_exchange = var.spec.post_quantum_key_exchange != "" ? var.spec.post_quantum_key_exchange : null

  # Empty defers to the provider default (DELETE).
  deletion_policy = var.spec.deletion_policy != "" ? var.spec.deletion_policy : null

  depends_on = [google_project_service.compute_api]
}

# The regional variant — identical configuration surface, addressed under
# regions/<region>/sslPolicies. Regional proxies can only reference SSL
# policies in their own region.
resource "google_compute_region_ssl_policy" "this" {
  count = local.is_regional ? 1 : 0

  name        = local.ssl_policy_name
  project     = local.project_id
  region      = var.spec.region
  description = var.spec.description != "" ? var.spec.description : null

  profile         = local.profile
  min_tls_version = local.min_tls_version

  custom_features = length(var.spec.custom_features) > 0 ? var.spec.custom_features : null

  post_quantum_key_exchange = var.spec.post_quantum_key_exchange != "" ? var.spec.post_quantum_key_exchange : null

  # Empty defers to the provider default (DELETE).
  deletion_policy = var.spec.deletion_policy != "" ? var.spec.deletion_policy : null

  depends_on = [google_project_service.compute_api]
}
