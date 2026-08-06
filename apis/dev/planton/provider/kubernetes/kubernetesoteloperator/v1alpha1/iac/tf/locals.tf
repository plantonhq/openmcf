# Computed values for the KubernetesOtelOperator module. Every resolution
# here has an exact twin in the Pulumi module's locals.go / values.go —
# keep them in lockstep.
#
# HCL DISCIPLINE (applies to every conditional object in this file):
# conditional entries are written as `key = cond ? value : null` inside ONE
# object literal, pruned with `{ for k, v in {...} : k => v if v != null }`.
# The two tempting alternatives are both broken: `cond ? {...} : {}`
# ternaries fail plan-time type unification when branches carry different
# attributes, and merge() over conditional lists silently UNIFIES
# primitive-only sibling objects into map(string). The null-prune form
# preserves every value's type.
#
# Optional nested blocks are read with try(): HCL's && does NOT
# short-circuit, so chained null checks still dereference the null.

locals {
  # Chart identity — must stay byte-identical with the Pulumi module's
  # vars: cross-engine chart drift deploys two different products from one
  # manifest.
  helm_chart_repo = "https://open-telemetry.github.io/opentelemetry-helm-charts"
  helm_chart_name = "opentelemetry-operator"

  # Release name = metadata.name.
  release_name = var.metadata.name

  # Chart version resolved to the pinned default when unset, so both
  # engines install the same chart whether or not the platform's
  # defaulting middleware ran — mirror of the Pulumi module's
  # DefaultChartVersion. 0.120.0 is the newest SERVED stable chart
  # (= operator appVersion 0.156.0, verified against the repository
  # index). Bumping it requires re-staging ../crds from the new pin —
  # the staged files ARE the chart's CRDs at this version.
  chart_version = coalesce(try(var.spec.chart_version, null), "0.120.0")

  namespace = var.spec.namespace

  # Resource-identity labels stamped on the namespace this module creates
  # (never injected into the chart's own resources — Helm owns those).
  labels = merge(
    {
      "planton.ai/resource"      = "true"
      "planton.ai/resource-name" = var.metadata.name
      "planton.ai/resource-kind" = "KubernetesOtelOperator"
    },
    var.metadata.id != null && var.metadata.id != "" ? { "planton.ai/resource-id" = var.metadata.id } : {},
    var.metadata.org != null && var.metadata.org != "" ? { "planton.ai/organization" = var.metadata.org } : {},
    var.metadata.env != null && var.metadata.env != "" ? { "planton.ai/environment" = var.metadata.env } : {}
  )

  # ---- module-owned CRDs (twin of the Pulumi module's crds.go) --------------
  # One entry per staged CRD file, keyed by the CRD's OWN metadata.name
  # (parsed from the file — never a positional index), so state addresses
  # stay stable across file renames and reorderings.
  #
  # THE STAGED FILES ARE RENDERED, TOKENIZED CHART TEMPLATES: unlike
  # plain-YAML CRD bundles, this chart TEMPLATES its CRDs — the collector
  # CRD carries the cert-manager.io/inject-ca-from annotation and a
  # version-conversion webhook clientConfig, both derived from the
  # RELEASE's identity. The files under ../crds were rendered from the
  # pinned chart with those release-derived values replaced by
  # __PLANTON_RELEASE_NAME__ / __PLANTON_NAMESPACE__ tokens, substituted
  # here (and identically in the Pulumi module) — so the kept CRDs always
  # point at THIS release's webhook Service and cert-manager Certificate.
  #
  # skip_crds is the bring-your-own-CRDs arm: the CRDs are owned
  # elsewhere (a GitOps-managed bundle) and this module must not touch
  # them. The filter clause keeps the for_each type stable (never a
  # `cond ? {} : {...}` ternary — the type-unification class).
  crd_manifests = {
    for f in fileset("${path.module}/../crds", "*.yaml") :
    yamldecode(file("${path.module}/../crds/${f}")).metadata.name =>
    replace(
      replace(file("${path.module}/../crds/${f}"), "__PLANTON_RELEASE_NAME__", local.release_name),
      "__PLANTON_NAMESPACE__", local.namespace
    )
    if !try(var.spec.skip_crds, false)
  }

  # ---- webhook artifact names (twins of the Pulumi module's locals) ---------
  # The module pins the chart's fullnameOverride to the resource name
  # (see typed_values), so every chart-derived name below hangs off
  # metadata.name. The 30-character budget in main.tf comes from the
  # longest suffix: "-controller-manager-service-cert" (33 chars) against
  # the Kubernetes 63-character name limit.
  webhook_service          = "${local.release_name}-webhook"
  webhook_cert_secret_name = "${local.release_name}-controller-manager-service-cert"

  # ---- default collector image (chart: manager.collectorImage) --------------
  # The spec carries ONE image string; the chart takes repository and tag
  # separately (and renders --collector-image only when BOTH are present —
  # a repository-only override deep-merges with the chart's default tag,
  # so the flag still renders). A tag exists when the LAST ":" comes after
  # the last "/" (registry ports carry ":" too — "reg.example.com:5000/x"
  # has no tag). Twin of the Pulumi module's splitImageRef.
  collector_image_raw   = try(var.spec.default_collector_image, "")
  collector_image_parts = split(":", local.collector_image_raw)
  collector_image_has_tag = length(local.collector_image_parts) > 1 && !strcontains(
    element(local.collector_image_parts, length(local.collector_image_parts) - 1), "/"
  )
  collector_image = {
    for k, v in {
      repository = local.collector_image_raw != "" ? (
        local.collector_image_has_tag ?
        join(":", slice(local.collector_image_parts, 0, length(local.collector_image_parts) - 1)) :
        local.collector_image_raw
      ) : null
      tag = local.collector_image_has_tag ? element(local.collector_image_parts, length(local.collector_image_parts) - 1) : null
    } : k => v if v != null
  }

  # ---- operator manager container resources (shared ContainerResources) -----
  # Twin of the Pulumi module's resourcesMap. The chart ships NO default
  # requests/limits for the manager — the resources key renders only when
  # the spec sets them. Helm deep-merges per key, so a partial block
  # overrides only the halves it carries.
  manager_resources = try(var.spec.resources, null) == null ? null : {
    for k, v in {
      limits = try(var.spec.resources.limits, null) == null ? null : {
        for lk, lv in {
          cpu    = try(var.spec.resources.limits.cpu, "") != "" ? var.spec.resources.limits.cpu : null
          memory = try(var.spec.resources.limits.memory, "") != "" ? var.spec.resources.limits.memory : null
        } : lk => lv if lv != null
      }
      requests = try(var.spec.resources.requests, null) == null ? null : {
        for rk, rv in {
          cpu    = try(var.spec.resources.requests.cpu, "") != "" ? var.spec.resources.requests.cpu : null
          memory = try(var.spec.resources.requests.memory, "") != "" ? var.spec.resources.requests.memory : null
        } : rk => rv if rv != null
      }
    } : k => v if v != null && length(v) > 0
  }

  # ---- manager image (air-gap/private-mirror registry replacement) ----------
  # image_registry replaces ONLY the registry part of the manager image
  # (the one image this component's pods pull); the path stays the
  # upstream one. The default collector image the operator INJECTS into
  # CRs is mirrored via default_collector_image instead — collector pods
  # pull that one, not this component. Twin of the Pulumi module.
  manager_image = {
    for k, v in {
      repository = try(var.spec.image_registry, "") != "" ? "${var.spec.image_registry}/open-telemetry/opentelemetry-operator/opentelemetry-operator" : null
    } : k => v if v != null
  }

  # ---- the manager block (every key renders only when the spec sets it) -----
  manager_values = {
    for k, v in {
      resources      = local.manager_resources != null && length(local.manager_resources) > 0 ? local.manager_resources : null
      image          = length(local.manager_image) > 0 ? local.manager_image : null
      collectorImage = length(local.collector_image) > 0 ? local.collector_image : null

      # Plain bool (no presence): false IS the chart default, so only
      # true renders. Requires the monitoring.coreos.com CRDs on the
      # cluster (KubernetesKubePrometheusStack).
      serviceMonitor = try(var.spec.service_monitor_enabled, false) ? { enabled = true } : null
    } : k => v if v != null
  }

  # ---- cert-manager issuer reference ----------------------------------------
  # Rendered only when the spec names an issuer; empty means the chart
  # creates its own self-signed Issuer (the default posture). cert-manager
  # itself is NOT optional — see the crds.create/certManager re-pin note
  # in main.tf.
  issuer_ref = try(var.spec.webhook.issuer_ref, null) == null ? null : {
    for k, v in {
      kind = try(var.spec.webhook.issuer_ref.kind, "") != "" ? var.spec.webhook.issuer_ref.kind : null
      name = try(var.spec.webhook.issuer_ref.name, "") != "" ? var.spec.webhook.issuer_ref.name : null
    } : k => v if v != null
  }

  # ---- typed chart values (twin of the Pulumi module's buildHelmValues) ------
  # Chart-default-matching values render only on divergence, so the
  # rendered values stay minimal on both engines — with TWO deliberate
  # exceptions pinned in main.tf's values list: crds.create=false and
  # admissionWebhooks.certManager.enabled=true (see main.tf).
  typed_values = {
    for k, v in {
      # crds.create: false ALWAYS — never conditional, never a spec knob.
      # The chart templates its CRDs release-owned (a Helm uninstall
      # would cascade-delete every collector in the cluster); the module
      # owns the CRDs instead (kubectl_manifest.crds in main.tf).
      crds = { create = false }

      # fullnameOverride pins the chart's fullname to the resource name
      # (the catalog's Helm-kind identity convention). Load-bearing: the
      # staged CRDs' conversion webhook and inject-ca-from annotation
      # point at names derived from it (see crd_manifests above).
      fullnameOverride = local.release_name

      # Rendered on presence — an explicit 1 re-states the chart default
      # harmlessly; 2 gives a warm standby behind leader election.
      replicaCount = try(var.spec.replicas, null)

      manager = length(local.manager_values) > 0 ? local.manager_values : null

      admissionWebhooks = local.issuer_ref != null && length(local.issuer_ref) > 0 ? {
        certManager = { issuerRef = local.issuer_ref }
      } : null

      # Pull secrets ride the chart's TOP-LEVEL imagePullSecrets list
      # (raw Kubernetes object list, piped into the pod spec — verified
      # in the deployment template; this chart does NOT nest them under
      # manager).
      imagePullSecrets = length(try(var.spec.image_pull_secrets, [])) > 0 ? [for s in var.spec.image_pull_secrets : { name = s }] : null

      # Scheduling keys are TOP-LEVEL in this chart. nodeSelector
      # deep-merges over the chart's default {kubernetes.io/os: linux},
      # so the OS pin survives a spec selector.
      nodeSelector = length(try(var.spec.scheduling.node_selector, {})) > 0 ? var.spec.scheduling.node_selector : null
      tolerations = length(try(var.spec.scheduling.tolerations, [])) > 0 ? [
        for t in var.spec.scheduling.tolerations : {
          for tk, tv in {
            key               = try(t.key, "") != "" ? t.key : null
            operator          = try(t.operator, "") != "" ? t.operator : null
            value             = try(t.value, "") != "" ? t.value : null
            effect            = try(t.effect, "") != "" ? t.effect : null
            tolerationSeconds = try(t.toleration_seconds, null)
          } : tk => tv if tv != null
        }
      ] : null
      priorityClassName = try(var.spec.scheduling.priority_class_name, "") != "" ? var.spec.scheduling.priority_class_name : null
    } : k => v if v != null
  }
}
