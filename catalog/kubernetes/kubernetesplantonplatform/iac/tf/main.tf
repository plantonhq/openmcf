# KubernetesPlantonPlatform Terraform module.
#
# Declares one PlantonPlatform custom resource that the Planton operator
# (KubernetesPlantonOperator) reconciles into a running self-hosted
# platform — control plane, console, identity server, databases, secrets
# manager, and in-cluster runner, in one namespace — plus the optional
# owning namespace.
#
# The CR is rendered from null-pruned locals (locals.platform_spec): keys
# render ONLY when the manifest declared them, so the operator's own
# defaulting stays authoritative for everything unset. The exact twin of
# the Pulumi module's apiextensions.NewCustomResource + platformSpecBody.
#
# DESTROY: platform teardown is Kubernetes garbage collection — every
# operator-created object is owner-referenced to the CR, so deletion
# completes even when the operator itself is already gone. The delete
# timeout is headroom, not an expected wait.

# The optional owning namespace. Created before the CR; deleted with the
# resource (pre-existing-namespace installs leave create_namespace false).
resource "kubernetes_namespace_v1" "planton_platform" {
  count = try(var.spec.create_namespace, false) ? 1 : 0

  metadata {
    name   = local.namespace
    labels = local.labels
  }
}

# The PlantonPlatform declaration.
resource "kubectl_manifest" "planton_platform" {
  yaml_body = yamlencode({
    apiVersion = local.api_version
    kind       = local.cr_kind
    metadata = {
      # Namespaced, named from THIS resource's metadata.name — the prefix
      # of every object the operator creates for the platform.
      name      = local.platform_name
      namespace = local.namespace
      labels    = local.labels
    }
    spec = local.platform_spec
  })

  server_side_apply = true
  force_conflicts   = true

  timeouts {
    delete = "15m"
  }

  depends_on = [kubernetes_namespace_v1.planton_platform]
}
