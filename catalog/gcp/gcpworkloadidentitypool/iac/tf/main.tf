# The Workload Identity Pool — the trust boundary external identities
# federate into. The pool holds no issuer configuration; providers
# (google_iam_workload_identity_pool_provider) attach one per external issuer,
# and IAM principals are built from the pool's resource name.
#
# workload_identity_pool_id, project, and mode are immutable (the API rejects
# mode updates even though a plan may show one): changing any of them destroys
# and recreates the pool, invalidating every principal built from the old pool
# name. display_name, description, disabled, and the inline certificate/trust
# configs all update in place.
#
# GCP soft-deletes pools: after destroy, the pool remains DELETED for ~30 days
# and its ID cannot be reused until permanent deletion. Unlike custom roles
# there is NO undelete-on-create — recreating a pool with a soft-deleted ID
# fails outright — so treat pool IDs as long-lived and prefer `disabled` for
# temporary shutoffs.
resource "google_iam_workload_identity_pool" "this" {
  workload_identity_pool_id = var.spec.workload_identity_pool_id
  project                   = local.project_id
  display_name              = var.spec.display_name
  description               = var.spec.description
  # Always sent: `disabled` is the documented kill switch, and an explicit
  # false keeps re-enable flows (true -> false) working.
  disabled = var.spec.disabled
  mode     = local.mode

  # The mTLS trust-domain surface. Both optional; FEDERATION_ONLY pools
  # normally set neither.
  dynamic "inline_certificate_issuance_config" {
    for_each = local.certificate_issuance != null ? [local.certificate_issuance] : []
    content {
      ca_pools                   = inline_certificate_issuance_config.value.ca_pools
      key_algorithm              = inline_certificate_issuance_config.value.key_algorithm
      lifetime                   = inline_certificate_issuance_config.value.lifetime
      rotation_window_percentage = inline_certificate_issuance_config.value.rotation_window_percentage
    }
  }

  dynamic "inline_trust_config" {
    for_each = var.spec.inline_trust_config != null ? [var.spec.inline_trust_config] : []
    content {
      dynamic "additional_trust_bundles" {
        for_each = inline_trust_config.value.additional_trust_bundles
        content {
          trust_domain = additional_trust_bundles.value.trust_domain

          dynamic "trust_anchors" {
            for_each = additional_trust_bundles.value.trust_anchors
            content {
              pem_certificate = trust_anchors.value.pem_certificate
            }
          }
        }
      }
    }
  }
}
