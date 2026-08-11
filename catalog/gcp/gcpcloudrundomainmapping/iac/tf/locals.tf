locals {
  # Honor the spec contract: an empty project_id falls back to the provider's
  # default project (ambient credentials decide).
  project_id = var.spec.project_id != "" ? var.spec.project_id : null

  # Platform attribution labels stored on the mapping object. User labels
  # merge in first so the platform attribution labels can never be
  # clobbered by a spec label with the same key. metadata.name keys the
  # name label (the GcpDnsZone basis for domain-named kinds): the
  # mapping's cloud-side name is the domain itself, whose dots and
  # 253-char budget don't fit the Knative/K8s 63-char label-value
  # contract.
  base_labels = {
    "planton-ai_resource" = "true"
    "planton-ai_name"     = var.metadata.name
    "planton-ai_kind"     = "gcpcloudrundomainmapping"
  }

  id_label = (
    var.metadata.id != null && var.metadata.id != ""
  ) ? { "planton-ai_id" = var.metadata.id } : {}

  org_label = (
    var.metadata.org != null && var.metadata.org != ""
  ) ? { "planton-ai_organization" = var.metadata.org } : {}

  env_label = (
    var.metadata.env != null && var.metadata.env != ""
  ) ? { "planton-ai_environment" = var.metadata.env } : {}

  platform_labels = merge(local.base_labels, local.id_label, local.org_label, local.env_label)
  final_labels    = merge(var.spec.labels, local.platform_labels)

  # The Cloud Run v1 API requires the mapping's namespace to equal the
  # project ID or project number. Fallback chain: spec.namespace → the
  # spec's project → the provider's resolved default project (the
  # count-gated data source below supplies that last case only).
  namespace_needs_client_config = var.spec.namespace == "" && var.spec.project_id == ""

  namespace = (
    var.spec.namespace != "" ? var.spec.namespace :
    var.spec.project_id != "" ? var.spec.project_id :
    data.google_client_config.default[0].project
  )
}
