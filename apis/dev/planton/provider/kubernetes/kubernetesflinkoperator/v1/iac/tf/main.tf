# KubernetesFlinkOperator Terraform module.
#
# Installs the Apache Flink Kubernetes Operator from the official ASF
# chart — served per version at
# https://downloads.apache.org/flink/flink-kubernetes-operator-<version>/
# — as a single Helm release named after metadata.name. The operator
# reconciles FlinkDeployment (declared with KubernetesFlinkDeployment)
# and the other flink.apache.org custom resources into running Flink
# clusters.
#
# CRD LIFECYCLE: the chart ships its four flink.apache.org CRDs from its
# crds/ DIRECTORY — Helm installs them once, never upgrades them on chart
# upgrades, and LEAVES them on uninstall (no release ownership metadata).
# That upstream posture is exactly the keep-on-uninstall this catalog
# wants for workload-bearing CRDs, so the module neither re-owns nor
# templates them — chart-version bumps that change CRDs are applied
# manually per the upstream release notes.
#
# THE WEBHOOK LIFECYCLE (chart truth at 1.15.0): with the webhook enabled
# (the upstream default this module keeps), the chart renders
# cert-manager Issuer/Certificate resources UNCONDITIONALLY — cert-manager
# is this kind's registry prerequisite, there is no self-signed fallback —
# and both webhook configurations are failurePolicy Fail. webhook_enabled=
# false removes the webhook, the certificate machinery, and the
# cert-manager dependency; the operator still validates in its reconcile
# loop.
#
# THE KEYSTORE PASSWORD (why this module generates a credential): the
# chart's default webhook keystore Secret is a HARDCODED PUBLIC PASSWORD
# ("password1234", base64 in templates/webhook/secret.yaml behind
# keystore.useDefaultPassword=true). It must never ship — the module
# generates a random password, materializes it as a module-owned Secret,
# and wires webhook.keystore.passwordSecretRef at it. useDefaultPassword=
# false is additionally RE-PINNED after the escape-hatch merge so the
# public default cannot resurface through helm_values.
#
# The typed spec renders into chart values (locals.typed_values); the
# helm_values escape hatch is passed as a SECOND values document, which
# the provider merges over the first with Helm -f semantics; the keystore
# re-pin rides a THIRD document merged after both — the exact semantic
# twin of the Pulumi module's buildHelmValues + mergeMaps + post-merge
# re-pin.

# The optional installation namespace. Created before the release; deleted
# with the resource (pre-existing-namespace installs leave create_namespace
# false).
resource "kubernetes_namespace_v1" "flink_operator" {
  count = try(var.spec.create_namespace, false) ? 1 : 0

  metadata {
    name   = local.namespace
    labels = local.labels
  }
}

# The webhook keystore password (webhook-enabled installs only — the
# disabled arm creates no random/Secret resources at all). Letters+digits
# only: letters-only-safe alphabets avoid whole config-parser bug classes
# (the credential lands in a JVM keystore env value); the 32-char length
# compensates the smaller alphabet. The generation-shape arguments are
# ignored after creation so an IMPORTED credential never silently
# regenerates: rotation stays an explicit verb, never plan fallout.
# Twin: the Pulumi module's RandomPassword in keystore_secret.go.
resource "random_password" "webhook_keystore" {
  count = local.webhook_enabled ? 1 : 0

  length  = 32
  special = false

  lifecycle {
    ignore_changes = [
      length, special, upper, lower, numeric,
      min_lower, min_numeric, min_special, min_upper, override_special,
    ]
  }
}

# The module-owned keystore-password Secret the chart's
# webhook.keystore.passwordSecretRef points at — the replacement for the
# chart's hardcoded-public-password default Secret.
resource "kubernetes_secret_v1" "webhook_keystore" {
  count = local.webhook_enabled ? 1 : 0

  metadata {
    name      = local.webhook_keystore_secret_name
    namespace = local.namespace
    labels    = local.labels
  }

  data = {
    password = random_password.webhook_keystore[0].result
  }

  depends_on = [
    kubernetes_namespace_v1.flink_operator,
  ]
}

# The operator release.
resource "helm_release" "flink_operator" {
  name = local.release_name
  # The repository URL carries the chart version — the chart is served
  # from a versioned Apache downloads directory.
  repository = "https://downloads.apache.org/flink/flink-kubernetes-operator-${local.chart_version}/"
  chart      = local.helm_chart_name
  version    = local.chart_version
  namespace  = local.namespace

  # The module owns namespace creation (create_namespace flag).
  create_namespace = false

  # Wait for the operator Deployment to become Available — a JVM with a
  # 30s-initial-delay startup probe, plus (webhook arm) a cert-manager
  # certificate the webhook container mounts; an unpullable image, an
  # absent cert-manager, or a broken config should fail THIS apply with
  # a readiness timeout, not surface later as FlinkDeployments that
  # mysteriously never reconcile.
  wait            = true
  atomic          = true
  cleanup_on_fail = true
  timeout         = 600

  # Documents merged in order by the provider (helm -f semantics): the
  # typed rendering first, the user's escape hatch second, and — webhook
  # arm only — the keystore re-pin LAST, so the chart's hardcoded-
  # password default can never resurface through helm_values.
  values = concat(
    [yamlencode(local.typed_values)],
    try(var.spec.helm_values, "") != "" ? [var.spec.helm_values] : [],
    local.webhook_enabled ? [yamlencode(local.keystore_repin_values)] : []
  )

  depends_on = [
    kubernetes_namespace_v1.flink_operator,
    kubernetes_secret_v1.webhook_keystore,
  ]

  lifecycle {
    precondition {
      # 63-char Kubernetes name limit minus the module's longest
      # derived child name, "<name>-webhook-keystore" (17-char suffix
      # — the published 45-char budget keeps one character of
      # headroom). The chart's webhook
      # Service/certificate/issuer names are CHART-FIXED
      # ("flink-operator-webhook-service", "flink-operator-serving-
      # cert", "flink-operator-selfsigned-issuer") — not fullname-
      # derived, so they are excluded from the budget. Twin of the
      # Pulumi module's fail-loud guard.
      condition     = length(var.metadata.name) <= 45
      error_message = "metadata.name must be 45 characters or fewer: the module derives \"<name>-webhook-keystore\" and Kubernetes caps names at 63."
    }
  }
}
