locals {
  # Honor the spec contract: an empty project_id falls back to the provider's
  # default project. Passing null (instead of "") lets the google provider
  # resolve its own project from configuration or the GOOGLE_PROJECT /
  # GOOGLE_CLOUD_PROJECT environment chain.
  project_id = var.spec.project_id != "" ? var.spec.project_id : null

  # Connector name defaults to metadata.name. GCP caps connector names at
  # 25 characters — shorter than most resource names — so a metadata.name
  # over the cap fails at the API, not silently truncated here.
  connector_name = var.spec.connector_name != "" ? var.spec.connector_name : var.metadata.name

  network       = var.spec.network != "" ? var.spec.network : null
  ip_cidr_range = var.spec.ip_cidr_range != "" ? var.spec.ip_cidr_range : null
  machine_type  = var.spec.machine_type != "" ? var.spec.machine_type : null

  # NOTE ON LABELS: google_vpc_access_connector has no labels surface at all
  # (verified against the released provider schema), so this module — unlike
  # nearly every other GCP kind — attaches no platform attribution labels.
  # Both engines skip labels identically; attribution rides on the connector
  # name and the Planton control plane's own records.
}
