# KubernetesTekton Terraform module.
#
# Renders the cluster's TektonConfig — the single declaration the Tekton
# Operator (a KubernetesTektonOperator install, the registry
# prerequisite) reconciles into running components via
# TektonInstallerSet resources.
#
# THE CR NAME IS FIXED: the operator's admission webhook allows exactly
# one TektonConfig per cluster and requires the name `config` —
# metadata.name of the Planton resource keys the state identity only.
#
# DESTROY SEMANTICS (the reason this kind exists separately from the
# operator): deleting the TektonConfig triggers the operator to tear
# down every component it installed — the TektonInstallerSet finalizers
# are processed by the RUNNING operator, and this deletion blocks until
# that teardown completes (a 15-minute delete timeout covers the
# full-profile teardown). Destroying this resource BEFORE the operator
# is exactly what makes the teardown clean.

resource "kubectl_manifest" "tekton_config" {
  yaml_body = yamlencode({
    apiVersion = local.api_version
    kind       = "TektonConfig"
    metadata = {
      # Cluster-scoped, operator-required fixed name.
      name   = local.tekton_config_name
      labels = local.labels
    }
    spec = local.tekton_config_spec
  })

  server_side_apply = true
  force_conflicts   = true

  timeouts {
    delete = "15m"
  }
}
