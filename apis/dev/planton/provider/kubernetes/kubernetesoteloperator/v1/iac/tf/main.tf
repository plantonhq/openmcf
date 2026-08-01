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
# CRDs: crds.create pins false unconditionally in the rendered values, and
# kubectl_manifest.crds below applies the staged CRD files with
# apply_only = true — destroy leaves them on the cluster. The release
# depends on the CRDs so the operator never starts against an
# unregistered API group.
#
# THE CONVERSION-TRUST COUPLING (why cert-manager is required): the
# collector CRD carries a version-conversion webhook and the
# cert-manager.io/inject-ca-from annotation. Because the CRDs are
# retained past the release's lifetime, their conversion trust must be
# kept current by a RUNNING reconciler — cert-manager's CA injector —
# not by a certificate embedded once at install time. The staged CRD
# files are tokenized renders of the pinned chart (see locals.tf) so the
# kept CRDs always point at THIS release's webhook Service and
# Certificate.
#
# The typed spec renders into chart values (locals.typed_values); the
# helm_values escape hatch is passed as a SECOND values document, which
# the provider merges over the first with Helm -f semantics — the exact
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

# The module-owned opentelemetry.io CRDs, one resource per staged file,
# keyed by each CRD's OWN metadata.name (never a positional index) so
# state addresses stay stable across file renames and reorderings.
#
# KEEP-ON-UNINSTALL (the load-bearing option): apply_only = true makes the
# provider's Delete a NO-OP (verified in the provider source: "When true,
# Delete is a no-op") — destroying this module removes the CRDs from state
# WITHOUT deleting them from the cluster, so an operator uninstall never
# cascade-deletes OpenTelemetryCollector resources. The exact semantic
# twin of the Pulumi module's retainOnDelete on each CRD.
#
# Server-side apply is REQUIRED, not just preferred: the collector CRD is
# ~418 KB — far past the 262144-byte cap on the client-side
# last-applied-configuration annotation.
resource "kubectl_manifest" "crds" {
  for_each = local.crd_manifests

  yaml_body = each.value

  server_side_apply = true
  force_conflicts   = true
  apply_only        = true
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

  # Documents merged in order by the provider (helm -f semantics): the
  # typed rendering first, the user's escape hatch second — and the TWO
  # design-load-bearing keys re-pinned LAST, the deliberate exceptions to
  # the escape hatch's last-word contract (twin of the Pulumi module):
  #   - crds.create=false: the module owns the CRD lifecycle; handing
  #     them to Helm would arm the uninstall cascade-delete this design
  #     exists to prevent.
  #   - admissionWebhooks.certManager.enabled=true: the kept CRDs'
  #     conversion trust rides cert-manager's CA injector; disabling it
  #     would leave module-owned CRDs pointing at a Certificate that no
  #     longer exists and silently break collector-CR conversion.
  values = concat(
    [yamlencode(local.typed_values)],
    try(var.spec.helm_values, "") != "" ? [var.spec.helm_values] : [],
    [yamlencode({
      crds              = { create = false }
      admissionWebhooks = { certManager = { enabled = true } }
    })]
  )

  depends_on = [
    kubernetes_namespace_v1.otel_operator,
    kubectl_manifest.crds,
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
    precondition {
      # FAIL LOUDLY when the staged CRD files did not travel with the
      # module: fileset() over a missing ../crds directory returns
      # EMPTY and for_each would silently plan ZERO CRDs — the operator
      # then runs against whatever CRDs happen to exist (caught live: a
      # lane "passed" riding a previous install's retained CRDs). Four
      # is the staged count at chart 0.120.0 — restage ../crds and
      # update this count together with chart_version.
      condition     = try(var.spec.skip_crds, false) || length(local.crd_manifests) == 4
      error_message = "The staged CRD files under ../crds did not travel with the module (expected 4, found a different count) — the module owns the CRD lifecycle and cannot install without them. Deploy the module with its full iac/ directory tree."
    }
  }
}
