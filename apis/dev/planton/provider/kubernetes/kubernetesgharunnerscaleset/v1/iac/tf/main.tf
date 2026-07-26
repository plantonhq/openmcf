# KubernetesGhaRunnerScaleSet Terraform module.
#
# Deploys a GitHub Actions runner scale set from the official OCI chart
# as a real Helm release: the AutoscalingRunnerSet the controller (a
# KubernetesGhaRunnerScaleSetController install, the registry
# prerequisite) reconciles into a listener and ephemeral runner pods.
#
# SECRET DISCIPLINE: on the declared PAT / GitHub App arms the module
# materializes the credential as the `<name>-github-auth` Secret BEFORE
# the release and passes only its NAME into chart values (the chart's
# pre-defined-secret form) — credential material never rides rendered
# values. The existing-Secret arm references the user's own Secret.
#
# NO HELM WAIT: the chart's workload is the AutoscalingRunnerSet custom
# resource — the controller creates the listener AFTER the release
# returns, and the listener needs valid GitHub credentials to come up.
# The E2E verifier owns the listener-registered proof instead (the
# CR-kind precedent).

# The optional installation namespace. Created before the release; deleted
# with the resource.
resource "kubernetes_namespace_v1" "gha_runner_scale_set" {
  count = var.spec.create_namespace ? 1 : 0

  metadata {
    name   = local.namespace
    labels = local.labels
  }
}

# The materialized GitHub credential (declared arms only) — the chart's
# pre-defined-secret key contract. Pulumi twin: githubAuthSecret.
resource "kubernetes_secret_v1" "github_auth" {
  count = local.materialize_auth_secret ? 1 : 0

  metadata {
    name      = local.github_auth_secret_name
    namespace = local.namespace
    labels    = local.labels
  }

  data = local.github_auth_secret_data

  depends_on = [
    kubernetes_namespace_v1.gha_runner_scale_set,
  ]
}

resource "helm_release" "gha_runner_scale_set" {
  name       = local.release_name
  repository = local.helm_oci_repo
  chart      = local.helm_chart_name
  version    = local.chart_version
  namespace  = local.namespace

  # The module owns namespace creation (create_namespace flag).
  create_namespace = false

  # No Helm wait — see the module comment (the workload is a CR the
  # controller reconciles after the release returns).
  wait    = false
  timeout = 300

  # Documents merged in order by the provider (helm -f semantics): the
  # typed rendering first, the user's escape hatch second — and
  # githubConfigSecret re-pinned LAST, the one deliberate exception to
  # the escape hatch's last-word contract (twin of the Pulumi module).
  # The credential contract (a Secret NAME, never inline material) is
  # load-bearing; letting an override move it would break the secret
  # discipline.
  values = concat(
    [yamlencode(local.typed_helm_values)],
    try(var.spec.helm_values, "") != "" ? [var.spec.helm_values] : [],
    [yamlencode({ githubConfigSecret = local.github_auth_secret_name })]
  )

  lifecycle {
    # FAIL LOUDLY past the chart's own scale-set-name budget: the chart
    # template fails installs at >45 characters (a GitHub registration
    # limit) — catch it at plan instead of mid-apply. The spec CEL caps
    # the explicit field; this guard covers the metadata.name fallback.
    # Twin: the Pulumi module's Resources() guard.
    precondition {
      condition     = length(local.runner_scale_set_name) <= 45
      error_message = "GitHub caps runner scale set registrations at 45 characters — set spec.runner_scale_set_name (or a shorter metadata.name)."
    }
  }

  depends_on = [
    kubernetes_namespace_v1.gha_runner_scale_set,
    kubernetes_secret_v1.github_auth,
  ]
}
