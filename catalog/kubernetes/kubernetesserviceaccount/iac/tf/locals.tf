# Local values and computed configuration

locals {
  # Build combined labels
  standard_labels = {
    "managed-by"    = "planton"
    "resource"      = var.metadata.name
    "resource-kind" = "KubernetesServiceAccount"
  }

  labels = merge(local.standard_labels, var.spec.labels)

  # Workload-identity annotation translation: the selected arm becomes the exact
  # ServiceAccount annotation the cloud's webhook/agent watches. This is the entire
  # Kubernetes-side half of the federation; the cloud-side half (IAM binding, trust
  # policy, federated credential) is owned by the referenced cloud identity resource.
  #
  # try() is required here: HCL's && does NOT short-circuit, so an expression like
  # `var.spec.workload_identity != null && var.spec.workload_identity.gke != null`
  # would still evaluate (and error on) the attribute access when workload_identity
  # is null. try() swallows the traversal error and yields the empty-map fallback.
  workload_identity_annotations = merge(
    # GKE Workload Identity: pod tokens impersonate this GCP service account.
    try({ "iam.gke.io/gcp-service-account" = var.spec.workload_identity.gke.service_account_email }, {}),
    # EKS IRSA: the pod-identity webhook injects credentials for this IAM role.
    try({ "eks.amazonaws.com/role-arn" = var.spec.workload_identity.eks.role_arn }, {}),
    # Azure AD Workload Identity: the client ID of the managed identity / Entra app.
    try({ "azure.workload.identity/client-id" = var.spec.workload_identity.aks.client_id }, {}),
    # tenant-id is annotated only when explicitly set (cross-tenant scenarios);
    # otherwise the Azure webhook falls back to its default tenant.
    try(
      var.spec.workload_identity.aks.tenant_id != null
      ? { "azure.workload.identity/tenant-id" = var.spec.workload_identity.aks.tenant_id }
      : {},
      {}
    )
  )

  # Workload-identity annotations are merged last so they win on collision: the
  # typed workload_identity field is the authoritative expression of the cloud
  # binding, and a stray user annotation must not silently override it.
  annotations = merge(var.spec.annotations, local.workload_identity_annotations)

  # The bound cloud identity handle (email / role ARN / client ID) for outputs.
  # try() falls through each unset arm's traversal error to the next candidate,
  # ending at "" when no workload identity is configured.
  workload_identity_handle = try(
    var.spec.workload_identity.gke.service_account_email,
    var.spec.workload_identity.eks.role_arn,
    var.spec.workload_identity.aks.client_id,
    ""
  )

  # The exact string cloud trust configuration (IAM trust policies, federated
  # credentials) matches on — exported so downstream never re-assembles it.
  rbac_subject = "system:serviceaccount:${var.spec.namespace}:${var.spec.name}"
}
