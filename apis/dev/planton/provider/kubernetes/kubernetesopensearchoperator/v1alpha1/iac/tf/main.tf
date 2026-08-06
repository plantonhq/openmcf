# KubernetesOpenSearchOperator Terraform module.
#
# Installs the OpenSearch Kubernetes Operator from the official
# opensearch-operator Helm chart as a single Helm release named after
# metadata.name. The operator reconciles OpenSearchCluster custom
# resources (declared through KubernetesOpenSearch) into running search
# clusters with managed TLS, security bootstrap, rolling upgrades and
# Dashboards.
#
# CRD LIFECYCLE: the chart templates its ten CRDs as release-owned
# resources with NO keep-on-uninstall knob upstream, so a Helm-owned
# install would cascade-delete every OpenSearchCluster (and its data) on
# uninstall. The module therefore OWNS the CRDs: installCRDs pins false
# unconditionally in the rendered values, and kubectl_manifest.crds below
# applies the staged CRD files with apply_only = true — destroy leaves
# them on the cluster. The release depends on the CRDs so the operator
# never starts against an unregistered API group.
#
# The typed spec renders into chart values (locals.typed_values); the
# helm_values escape hatch is passed as a SECOND values document, which
# the provider merges over the first with Helm -f semantics — the exact
# semantic twin of the Pulumi module's buildHelmValues + mergeMaps.

# The optional installation namespace. Created before the release; deleted
# with the resource (pre-existing-namespace installs leave create_namespace
# false).
resource "kubernetes_namespace_v1" "opensearch_operator" {
  count = try(var.spec.create_namespace, false) ? 1 : 0

  metadata {
    name   = local.namespace
    labels = local.labels
  }
}

# The module-owned OpenSearch CRDs, one resource per staged file, keyed by
# each CRD's OWN metadata.name (never a positional index) so state
# addresses stay stable across file renames and reorderings.
#
# KEEP-ON-UNINSTALL (the load-bearing option): apply_only = true makes the
# provider's Delete a NO-OP (verified in the provider source: "When true,
# Delete is a no-op") — destroying this module removes the CRDs from state
# WITHOUT deleting them from the cluster, so an operator uninstall never
# cascade-deletes OpenSearchCluster resources and their data. The exact
# semantic twin of the Pulumi module's retainOnDelete on each CRD.
#
# Server-side apply keeps the CRDs co-ownable (no client-side
# last-applied-configuration annotation bloat on the megabyte-scale
# schemas, and a re-adopting apply never conflicts with itself).
resource "kubectl_manifest" "crds" {
  for_each = local.crd_manifests

  yaml_body = each.value

  server_side_apply = true
  force_conflicts   = true
  apply_only        = true
}

# The operator release.
resource "helm_release" "opensearch_operator" {
  name       = local.release_name
  repository = local.helm_chart_repo
  chart      = local.helm_chart_name
  version    = local.chart_version
  namespace  = local.namespace

  # The module owns namespace creation (create_namespace flag).
  create_namespace = false

  # Wait for the operator to become Available — an operator that never
  # becomes ready (an unpullable image from a private mirror is the
  # classic case) should fail THIS apply with a readiness timeout, not
  # surface later as OpenSearch clusters that mysteriously never
  # reconcile.
  wait            = true
  atomic          = true
  cleanup_on_fail = true
  timeout         = 600

  # Documents merged in order by the provider (helm -f semantics): the
  # typed rendering first, the user's escape hatch second — and
  # installCRDs re-pinned LAST, the one deliberate exception to the
  # escape hatch's last-word contract (twin of the Pulumi module). The
  # module owns the CRD lifecycle; letting an override hand them to Helm
  # would arm the cascade-delete this design exists to prevent.
  values = concat(
    [yamlencode(local.typed_values)],
    try(var.spec.helm_values, "") != "" ? [var.spec.helm_values] : [],
    [yamlencode({ installCRDs = false })]
  )

  depends_on = [
    kubernetes_namespace_v1.opensearch_operator,
    kubectl_manifest.crds,
  ]

  lifecycle {
    precondition {
      # FAIL LOUDLY when the staged CRD files did not travel with the
      # module: fileset() over a missing ../crds directory returns
      # EMPTY and for_each would silently plan ZERO CRDs — the operator
      # then runs against whatever CRDs happen to exist (the class was
      # caught live elsewhere: a lane "passed" riding a previous
      # install's retained CRDs). Ten is the staged count at chart
      # 2.8.0 — restage ../crds and update this count together with
      # chart_version. Twin of the Pulumi module's fail-loud guard.
      # The guard sits on the release (not kubectl_manifest.crds): a
      # for_each resource with an empty map evaluates ZERO
      # preconditions, which is exactly the failure being guarded.
      condition     = length(local.crd_manifests) == 10
      error_message = "The staged CRD files under ../crds did not travel with the module (expected 10, found a different count) — the module owns the CRD lifecycle and cannot install without them. Deploy the module with its full iac/ directory tree."
    }
  }
}
