# KubernetesArgocd Terraform module.
#
# Installs Argo CD from the official Helm chart as a real Helm release. The
# typed spec renders into chart values (locals.helm_values); the initial
# admin password stays APPLICATION-owned (Argo CD generates it at first
# start into the fixed-name `argocd-initial-admin-secret`); SSO client
# secrets ride Argo CD's `$<secret>:<key>` runtime indirection against
# labeled Secrets; the helm_values escape hatch is passed as a SECOND
# values document, which the provider merges over the first with Helm -f
# semantics — the exact semantic twin of the Pulumi module's
# buildHelmValues + mergeMaps.

# The optional installation namespace. Created before the release; deleted
# with the resource.
resource "kubernetes_namespace_v1" "argocd" {
  count = var.spec.create_namespace ? 1 : 0

  metadata {
    name   = local.namespace
    labels = local.labels
  }
}

resource "helm_release" "argocd" {
  name       = local.release_name
  repository = local.helm_chart_repo
  chart      = local.helm_chart_name
  version    = local.chart_version
  namespace  = local.namespace

  # The module owns namespace creation (create_namespace flag).
  create_namespace = false

  # Wait for the control plane to become Ready — a server that never
  # starts (bad OIDC issuer, unschedulable redis-ha, missing external
  # Redis Secret) should fail THIS apply, not the first login. The budget
  # covers all seven components plus the redis-secret-init hook Job.
  wait            = true
  atomic          = true
  cleanup_on_fail = true
  timeout         = 600

  # Documents merged in order by the provider (helm -f semantics): the
  # typed rendering first, the user's escape hatch second — and
  # fullnameOverride re-pinned LAST, the one deliberate exception to the
  # escape hatch's last-word contract (twin of the Pulumi module). Every
  # child name (`<name>-server`, `<name>-application-controller`, ...) —
  # and the exported outputs built from them — derive from the fullname;
  # letting an override move it would break every output.
  values = concat(
    [yamlencode(local.helm_values)],
    try(var.spec.helm_values, "") != "" ? [var.spec.helm_values] : [],
    [yamlencode({ fullnameOverride = local.release_name })]
  )

  lifecycle {
    # FAIL LOUDLY on names past the chart's fullname budget: every child
    # name is `<fullname>-<component>` truncated at 63 characters — the
    # longest component suffix ("-applicationset-controller", 26 chars)
    # would truncate SILENTLY past 37, breaking the naming contract the
    # exported outputs are built on. Twin: the Pulumi module's
    # Resources() guard.
    precondition {
      condition     = length(var.metadata.name) <= 37
      error_message = "The argo-cd chart derives child names from the resource name and silently truncates past 37 characters (63 minus its longest component suffix), which would break the naming contract — use a name of at most 37 characters."
    }
  }

  depends_on = [
    kubernetes_namespace_v1.argocd,
  ]
}
