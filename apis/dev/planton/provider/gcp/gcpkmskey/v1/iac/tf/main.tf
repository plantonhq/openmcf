# Enable the Cloud KMS API on the key ring's project. The key inherits its
# project from the ring path, but google_project_service needs the project
# ID itself — extract it from the ring path's second segment
# (projects/{project}/locations/...). disable_on_destroy is false: tearing
# down one key must never disable the API for everything else in the
# project (other keys may be actively encrypting production data).
resource "google_project_service" "cloudkms_api" {
  project = split("/", local.key_ring_id)[1]
  service = "cloudkms.googleapis.com"

  disable_dependent_services = true
  disable_on_destroy         = false
}

# The KMS crypto key — the resource downstream services reference for
# customer-managed encryption (CMEK). Lifecycle sharp edges, all taught by
# the API rather than invented here:
#
#   - The key can never be deleted from GCP. On destroy the provider
#     destroys all key versions (data encrypted under them becomes
#     unrecoverable once the destroy-scheduled window elapses), disables
#     automatic rotation, and removes the key from state — the key object
#     itself remains, permanently and at no cost, in the ring.
#
#   - Only rotation_period, version_template.algorithm, and labels update
#     in place; every other field is ForceNew, which for an undeletable
#     resource means "abandon and create under a new name".
resource "google_kms_crypto_key" "this" {
  name     = local.key_name
  key_ring = local.key_ring_id
  labels   = local.final_labels

  purpose = local.purpose

  # Rotation mints a new primary version on the cadence; old versions stay
  # decryptable until destroyed. Only valid for ENCRYPT_DECRYPT keys
  # (enforced pre-deploy by the spec).
  rotation_period = local.rotation_period

  # The recovery window for destroyed versions (default 30 days).
  destroy_scheduled_duration = local.destroy_scheduled_duration

  # Create-time-only flag; required for import_only keys, where GCP must
  # never generate material.
  skip_initial_version_creation = var.spec.skip_initial_version_creation

  # BYOK: the key may only ever hold imported versions.
  import_only = var.spec.import_only ? true : null

  # EKM connection for EXTERNAL_VPC keys (the spec enforces the pairing
  # pre-deploy). Null for the SOFTWARE/HSM/EXTERNAL protection levels.
  crypto_key_backend = local.crypto_key_backend

  dynamic "version_template" {
    for_each = local.version_template != null ? [local.version_template] : []
    content {
      algorithm        = version_template.value.algorithm
      protection_level = version_template.value.protection_level != "" ? version_template.value.protection_level : null
    }
  }

  depends_on = [
    google_project_service.cloudkms_api,
  ]
}
