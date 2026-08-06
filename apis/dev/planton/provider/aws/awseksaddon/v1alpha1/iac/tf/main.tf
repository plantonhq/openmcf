# The managed add-on composes onto its neighbors instead of embedding
# them: the cluster attaches by reference, and the add-on's IAM identity
# -- when it needs one beyond the node role -- is a referenced AwsIamRole
# wired through IRSA or EKS Pod Identity (this module never modifies a
# role it merely references).
#
# AWS keys the add-on on (cluster, addon_name); the name and the
# namespace config are create-only, while version, configuration, the
# IAM wiring, and conflict handling update in place (the add-on rolls
# its own pods on version changes).
resource "aws_eks_addon" "this" {
  cluster_name = var.spec.cluster_name
  addon_name   = var.spec.addon_name

  # Empty means the AWS default version for the cluster's Kubernetes
  # version -- the never-goes-stale choice. AWS reports the resolved
  # version back through the addon_version output either way.
  addon_version = var.spec.addon_version != "" ? var.spec.addon_version : null

  # AWS's conflict handling is asymmetric (create: NONE/OVERWRITE;
  # update: +PRESERVE) -- the spec's CEL rules enforce the split before
  # anything reaches the API. Unset falls back to AWS's NONE, which
  # fails loudly on conflicts instead of adopting silently.
  resolve_conflicts_on_create = var.spec.resolve_conflicts_on_create != "" ? var.spec.resolve_conflicts_on_create : null
  resolve_conflicts_on_update = var.spec.resolve_conflicts_on_update != "" ? var.spec.resolve_conflicts_on_update : null

  configuration_values = var.spec.configuration_values != "" ? var.spec.configuration_values : null

  # IRSA: the referenced role must already exist and trust the cluster's
  # OIDC provider; empty means the add-on's pods use node-role permissions.
  service_account_role_arn = var.spec.service_account_role_arn != "" ? var.spec.service_account_role_arn : null

  # EKS Pod Identity: the modern no-OIDC-provider alternative. Each
  # association binds one service account to one referenced role.
  dynamic "pod_identity_association" {
    for_each = var.spec.pod_identity_associations
    content {
      role_arn        = pod_identity_association.value.role_arn
      service_account = pod_identity_association.value.service_account
    }
  }

  # Deleting with preserve leaves the add-on's Kubernetes resources
  # running as self-managed software -- the no-outage way to hand an
  # add-on's lifecycle back to cluster operators.
  preserve = var.spec.preserve

  # Create-only in AWS: changing the namespace replaces the add-on.
  dynamic "namespace_config" {
    for_each = var.spec.namespace_config != null ? [var.spec.namespace_config] : []
    content {
      namespace = namespace_config.value.namespace
    }
  }

  tags = local.aws_tags
}
