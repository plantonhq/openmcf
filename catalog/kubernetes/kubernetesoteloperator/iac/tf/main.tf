# KubernetesOtelOperator Terraform module.
#
# Installs the OpenTelemetry Operator from the official
# opentelemetry-operator Helm chart as a single Helm release named after
# metadata.name. The operator reconciles OpenTelemetryCollector custom
# resources (declared through KubernetesOtelCollector) into running
# collector fleets, defaulting and validating them through its admission
# webhooks.
#
# CRD LIFECYCLE: the chart templates its opentelemetry.io CRDs as
# release-owned resources, so a Helm-owned install would cascade-delete
# every collector declaration on uninstall. The module therefore OWNS the
# CRDs through the catalog's derive-branch primitive (the generated
# helm_crds.tf): the pinned chart is rendered at plan time with the
# release's own values plus the CRD switch turned on, each
# CustomResourceDefinition is applied keyed by its own name as a kept
# resource (retained on destroy unless crds.keep_on_uninstall is false;
# re-adopted on reinstall; refused when the manifest lowers chart_version
# below what the cluster carries), and the release below installs with
# skip_crds = true and crds.create pinned false so Helm never touches
# them. The release depends on the CRDs so the operator never starts
# against an unregistered API group.
#
# THE CONVERSION-TRUST COUPLING (why cert-manager is required): the
# collector CRD carries a version-conversion webhook and the
# cert-manager.io/inject-ca-from annotation, both rendered from THIS
# release's identity because the CRD render runs with the release's own
# values. Because the CRDs are retained past the release's lifetime,
# their conversion trust must be kept current by a RUNNING reconciler —
# cert-manager's CA injector — not by a certificate embedded once at
# install time.
#
# The typed spec renders into chart values (locals.typed_values); the
# helm_values escape hatch is the SECOND values document and the two
# load-bearing re-pins the THIRD (locals.helm_release_values) — the exact
# semantic twin of the Pulumi module's buildHelmValues + mergeMaps.

# The optional installation namespace. Created before the release; deleted
# with the resource (pre-existing-namespace installs leave create_namespace
# false).
resource "kubernetes_namespace_v1" "otel_operator" {
  count = try(var.spec.create_namespace, false) ? 1 : 0

  metadata {
    name   = local.namespace
    labels = local.labels
  }
}

# The operator release.
resource "helm_release" "otel_operator" {
  name       = local.release_name
  repository = local.helm_chart_repo
  chart      = local.helm_chart_name
  version    = local.chart_version
  namespace  = local.namespace

  # The module owns namespace creation (create_namespace flag).
  create_namespace = false

  # Wait for the operator to become Available — the manager pod mounts
  # the cert-manager-issued webhook Secret, so an install without a
  # working cert-manager (or with an unpullable image) should fail THIS
  # apply with a readiness timeout, not surface later as collectors that
  # mysteriously never reconcile.
  wait            = true
  atomic          = true
  cleanup_on_fail = true
  timeout         = 600

  # The same values list the CRD render consumed (see locals.tf for the
  # merge order and the two re-pins).
  values = local.helm_release_values

  # The module owns the CRDs (helm_crds.tf); Helm must never install its
  # own copy of them, whichever way crds.create is set.
  skip_crds = true

  depends_on = [
    kubernetes_namespace_v1.otel_operator,
    kubectl_manifest.helm_crds,
  ]

  lifecycle {
    precondition {
      # 63-char Kubernetes name limit minus the chart's longest derived
      # suffix, "-controller-manager-service-cert" (33 chars). The module
      # pins fullnameOverride to this name, so the budget is exact. Twin
      # of the Pulumi module's fail-loud guard.
      condition     = length(var.metadata.name) <= 30
      error_message = "metadata.name must be 30 characters or fewer: the chart derives \"<name>-controller-manager-service-cert\" (33-char suffix) and Kubernetes caps names at 63."
    }
  }
}
