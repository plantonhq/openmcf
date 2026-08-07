# KubernetesIstio Terraform module.
#
# Installs the Istio control plane from the official Helm charts as real Helm
# releases, in upstream's own order:
#
#   CRDs (module-owned, server-side apply)
#   -> base (validation plumbing) -> istiod (the control plane)
#   -> [ambient or cni.enabled] istio-cni (node agent)
#   -> [ambient] ztunnel (per-node L4 proxy)
#
# The CRDs deliberately apply OUTSIDE the base release (the chart installs
# them with base.excludedCRDs covering the whole bundle): Helm refuses to
# adopt CRDs that already exist without ITS ownership metadata, so a cluster
# running the CRDs-only KubernetesIstioBaseCrds kind could never upgrade to
# the full mesh if the chart owned them — server-side-applied CRDs are
# co-ownable by both kinds, making that migration a plain redeploy.
#
# The typed spec renders into per-chart values (locals.*_typed_values); each
# release's helm_values escape hatch is passed as a SECOND values document,
# which the provider merges over the first with Helm -f semantics — the exact
# semantic twin of the Pulumi module's build*Values + mergeMaps.
#
# Deliberately NO gateway release: istiod implements the Kubernetes Gateway
# API, so north-south gateways are composed from KubernetesGateway resources
# (gateway_class_name: istio) and istiod provisions their deployments itself.

# The optional installation namespace. Created before the releases; deleted
# with the resource.
resource "kubernetes_namespace_v1" "istio" {
  count = var.spec.create_namespace ? 1 : 0

  metadata {
    name   = local.namespace
    labels = local.labels
  }
}

# The Istio CRDs, module-owned (never Helm-owned — see the header). The
# bundle version is pinned to spec.version so the installed CRD schema
# matches the control plane and the typed Istio kinds' generated SDK.
data "http" "istio_crds" {
  url = local.crd_bundle_url

  request_headers = {
    Accept = "application/yaml"
  }
}

resource "kubectl_manifest" "istio_crds" {
  # Keyed by each CRD's OWN NAME (never the split index): the name is the
  # document's identity, so state addresses stay stable across bundle
  # reorderings AND the address key feeds the composed import ID blind
  # (from_address_key in the import map).
  for_each = {
    for doc in split("---", data.http.istio_crds.response_body) :
    yamldecode(doc).metadata.name => doc
    if trimspace(doc) != "" && can(yamldecode(doc).metadata.name)
  }

  yaml_body = each.value

  server_side_apply = true
  force_conflicts   = true
}

# base: the default-revision validation-webhook plumbing (CRDs excluded —
# module-owned above).
resource "helm_release" "base" {
  name       = "istio-base"
  repository = local.helm_chart_repo
  chart      = "base"
  version    = local.version
  namespace  = local.namespace

  # The module owns namespace creation (create_namespace flag).
  create_namespace = false

  wait            = true
  atomic          = true
  cleanup_on_fail = true
  timeout         = 300

  values = concat(
    [yamlencode(local.base_typed_values)],
    try(var.spec.helm_values.base, "") != "" ? [var.spec.helm_values.base] : []
  )

  depends_on = [kubernetes_namespace_v1.istio, kubectl_manifest.istio_crds]
}

# istiod: the control plane. Waiting for readiness is the whole promise — a
# control plane whose webhooks and discovery service are not serving rejects
# every mesh-config apply and every injection.
resource "helm_release" "istiod" {
  name       = local.istiod_release_name
  repository = local.helm_chart_repo
  chart      = "istiod"
  version    = local.version
  namespace  = local.namespace

  create_namespace = false

  wait            = true
  atomic          = true
  cleanup_on_fail = true
  timeout         = 600

  values = concat(
    [yamlencode(local.istiod_typed_values)],
    try(var.spec.helm_values.istiod, "") != "" ? [var.spec.helm_values.istiod] : []
  )

  depends_on = [helm_release.base]
}

# istio-cni: the node agent DaemonSet. Always installed in ambient mode (it is
# how traffic reaches ztunnel); opt-in in sidecar mode (replaces the injected
# privileged init-container).
resource "helm_release" "cni" {
  count = local.install_cni ? 1 : 0

  name       = "istio-cni"
  repository = local.helm_chart_repo
  chart      = "cni"
  version    = local.version
  namespace  = local.namespace

  create_namespace = false

  wait            = true
  atomic          = true
  cleanup_on_fail = true
  timeout         = 600

  values = concat(
    [yamlencode(local.cni_typed_values)],
    try(var.spec.helm_values.cni, "") != "" ? [var.spec.helm_values.cni] : []
  )

  depends_on = [helm_release.istiod]
}

# ztunnel: the ambient per-node L4 proxy DaemonSet.
resource "helm_release" "ztunnel" {
  count = local.ambient ? 1 : 0

  name       = "ztunnel"
  repository = local.helm_chart_repo
  chart      = "ztunnel"
  version    = local.version
  namespace  = local.namespace

  create_namespace = false

  wait            = true
  atomic          = true
  cleanup_on_fail = true
  timeout         = 600

  values = concat(
    [yamlencode(local.ztunnel_typed_values)],
    try(var.spec.helm_values.ztunnel, "") != "" ? [var.spec.helm_values.ztunnel] : []
  )

  depends_on = [helm_release.istiod]
}
