# Enable the Compute Engine API so a fresh project can host the NEG.
# disable_on_destroy is false: tearing down one NEG must never disable the API
# for everything else in the project.
resource "google_project_service" "compute_api" {
  project = local.project_id
  service = "compute.googleapis.com"

  disable_dependent_services = true
  disable_on_destroy         = false
}

# A regional network endpoint group — the bridge that lets a backend service
# target a Cloud Run/Functions/App Engine workload, a Private Service Connect
# endpoint, or an external origin instead of a group of VMs.
#
# The whole resource is immutable (every field is ForceNew): any change
# destroys and recreates the NEG. Because an in-use NEG cannot be deleted, a
# NEG referenced by a backend service should be recreated create-before-destroy
# to avoid a resourceInUseByAnotherResource error.
#
# The endpoint type gates which nested block is set; the spec's CEL rules
# enforce the "exactly one serverless block for SERVERLESS, none otherwise" and
# PSC/internet coherence before deploy, so this module stays declarative and
# sends whatever the spec set.
resource "google_compute_region_network_endpoint_group" "this" {
  name        = local.network_endpoint_group_name
  project     = local.project_id
  region      = var.spec.region
  description = var.spec.description != "" ? var.spec.description : null

  network_endpoint_type = local.network_endpoint_type

  # PSC / INTERNET / GCE_VM_IP_PORTMAP fields (unset for serverless).
  network            = local.network
  subnetwork         = local.subnetwork
  psc_target_service = local.psc_target_service

  dynamic "psc_data" {
    for_each = local.psc_data != null ? [local.psc_data] : []
    content {
      producer_port = psc_data.value.producer_port
    }
  }

  # Serverless targets — exactly one is set for a SERVERLESS NEG (enforced by
  # the spec's CEL). service/function arrive as resolved strings.
  dynamic "cloud_run" {
    for_each = local.cloud_run != null ? [local.cloud_run] : []
    content {
      service  = cloud_run.value.service
      tag      = cloud_run.value.tag
      url_mask = cloud_run.value.url_mask
    }
  }

  dynamic "cloud_function" {
    for_each = local.cloud_function != null ? [local.cloud_function] : []
    content {
      function = cloud_function.value.function
      url_mask = cloud_function.value.url_mask
    }
  }

  # The App Engine block may be empty (routes to the default app), so it is
  # emitted whenever present even with all leaves null.
  dynamic "app_engine" {
    for_each = local.app_engine != null ? [local.app_engine] : []
    content {
      service  = app_engine.value.service
      version  = app_engine.value.version
      url_mask = app_engine.value.url_mask
    }
  }

  deletion_policy = var.spec.deletion_policy != "" ? var.spec.deletion_policy : null

  depends_on = [google_project_service.compute_api]
}
