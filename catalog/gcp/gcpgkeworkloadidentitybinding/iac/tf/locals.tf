locals {
  # Honor the spec contract: an empty project_id falls back to the
  # provider's default project. The member string needs a concrete project,
  # so the fallback is made concrete by reading the provider's resolved
  # project from google_client_config.
  pool_project = (
    var.spec.project_id != ""
    ? var.spec.project_id
    : data.google_client_config.current.project
  )

  # The workload-identity principal, constructed from its parts so a typo'd
  # member string is impossible by construction:
  #   serviceAccount:<pool-project>.svc.id.goog[<namespace>/<ksa>]
  workload_identity_member = "serviceAccount:${local.pool_project}.svc.id.goog[${var.spec.ksa_namespace}/${var.spec.ksa_name}]"

  # The provider requires the fully-qualified service-account resource name,
  # not the bare email. The "-" project wildcard lets the IAM API infer the
  # SA's own project from the email — correct even when the GSA lives in a
  # different project than the workload-identity pool. Identical construction
  # in the Pulumi module.
  service_account_id = "projects/-/serviceAccounts/${var.spec.service_account_email}"
}

# The provider's own resolved configuration — the source of the default
# pool project when spec.project_id is omitted.
data "google_client_config" "current" {
}
