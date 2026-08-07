# Enable the Cloud KMS API — the control plane that owns key rings and keys.
# disable_on_destroy is false: tearing down one key ring must never disable
# the API for everything else in the project (especially here, where other
# rings' keys may be actively encrypting production data).
resource "google_project_service" "cloudkms_api" {
  project = local.project_id
  service = "cloudkms.googleapis.com"

  disable_dependent_services = true
  disable_on_destroy         = false
}

# The KMS key ring — the permanent organizational container for crypto keys.
# Every field is ForceNew, and GCP has no delete API for key rings: on
# destroy the provider removes the ring from state only, leaving the (free,
# inert) ring in the project forever. IAM granted on the ring flows down to
# every key inside it, which makes the ring the blast-radius boundary —
# group keys by environment or data domain, not one ring per key.
resource "google_kms_key_ring" "this" {
  name     = local.key_ring_name
  location = local.location
  project  = local.project_id

  depends_on = [
    google_project_service.cloudkms_api,
  ]
}
