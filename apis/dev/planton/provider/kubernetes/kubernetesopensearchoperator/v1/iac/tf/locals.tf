# Computed values for the KubernetesOpenSearchOperator module. Every
# resolution here has an exact twin in the Pulumi module's locals.go /
# values.go — keep them in lockstep.
#
# HCL DISCIPLINE (applies to every conditional object in this file):
# conditional entries are written as `key = cond ? value : null` inside ONE
# object literal, pruned with `{ for k, v in {...} : k => v if v != null }`.
# The two tempting alternatives are both broken: `cond ? {...} : {}`
# ternaries fail plan-time type unification when branches carry different
# attributes, and `merge(concat(cond ? [{...}] : [], ...)...)` silently
# UNIFIES primitive-only sibling objects into map(string) — numbers and
# booleans arrive in the chart values as strings. The null-prune form
# preserves every value's type.
#
# Optional nested blocks are read with try(): HCL's && does NOT
# short-circuit, so chained null checks still dereference the null — and
# var.spec is typed 'any', so an absent attribute is an error, not a null.

locals {
  # Chart identity — must stay byte-identical with the Pulumi module's
  # vars: cross-engine chart drift deploys two different products from one
  # manifest.
  helm_chart_repo = "https://opensearch-project.github.io/opensearch-k8s-operator/"
  helm_chart_name = "opensearch-operator"

  # Release name = metadata.name.
  release_name = var.metadata.name

  # Chart version resolved to the pinned default when unset, so both
  # engines install the same chart whether or not the platform's
  # defaulting middleware ran — mirror of the Pulumi module's
  # DefaultChartVersion. 2.8.0 is the newest SERVED chart whose default
  # manager image is a STABLE operator release (the 2.8.3+/3.0.x served
  # charts default to a prerelease image and the 3.x line migrates the
  # CRDs to the opensearch.org API group).
  chart_version = coalesce(try(var.spec.chart_version, null), "2.8.0")

  namespace = var.spec.namespace

  # Resource-identity labels stamped on the namespace this module creates
  # (never injected into the chart's own resources — Helm owns those).
  labels = merge(
    {
      "planton.ai/resource"      = "true"
      "planton.ai/resource-name" = var.metadata.name
      "planton.ai/resource-kind" = "KubernetesOpenSearchOperator"
    },
    var.metadata.id != null && var.metadata.id != "" ? { "planton.ai/resource-id" = var.metadata.id } : {},
    var.metadata.org != null && var.metadata.org != "" ? { "planton.ai/organization" = var.metadata.org } : {},
    var.metadata.env != null && var.metadata.env != "" ? { "planton.ai/environment" = var.metadata.env } : {}
  )

  # ---- module-owned CRDs (twin of the Pulumi module's crds.go) --------------
  # One entry per staged CRD file, keyed by the CRD's OWN metadata.name
  # (parsed from the file — never a positional index), so state addresses
  # stay stable across file renames and reorderings.
  crd_manifests = {
    for f in fileset("${path.module}/../crds", "*.yaml") :
    yamldecode(file("${path.module}/../crds/${f}")).metadata.name => file("${path.module}/../crds/${f}")
  }

  # ---- deployment name (twin of the Pulumi module's deploymentName) ---------
  # The module pins the chart's fullnameOverride to the resource name (see
  # typed_values), so the fullname IS the release name — verified live:
  # without the pin the chart's default fullname
  # ("<release>-opensearch-operator") pushes its metrics Service name
  # ("<fullname>-controller-manager-metrics-service") past Kubernetes'
  # 63-character limit for ordinary release names (the chart truncates the
  # fullname but not the names built from it, so the install fails at the
  # API server).
  deployment_name = "${local.release_name}-controller-manager"

  # ---- operator manager container resources (shared ContainerResources) -----
  # Twin of the Pulumi module's resourcesMap. The chart SHIPS default
  # requests/limits (requests 100m/350Mi, limits 200m/500Mi) — the
  # resources key renders only when the spec sets them, so the chart
  # defaults survive an empty spec. Helm deep-merges per key, so a
  # partial spec block overrides only the halves it carries.
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

  # ---- manager.image (rendered per half — deep-merges over the chart
  # defaults, leaving pullPolicy and the appVersion-derived tag intact) ------
  manager_image = {
    for k, v in {
      repository = try(var.spec.image.repository, "") != "" ? var.spec.image.repository : null
      tag        = try(var.spec.image.tag, "") != "" ? var.spec.image.tag : null
    } : k => v if v != null
  }

  # ---- the manager block (every key renders only when the spec sets it) -----
  manager_values = {
    for k, v in {
      watchNamespace = try(var.spec.watch_namespace, "") != "" ? var.spec.watch_namespace : null

      # loglevel (the chart's lowercase key) renders whenever the spec
      # carries a value; the chart default is "info".
      loglevel = try(var.spec.log_level, "") != "" ? var.spec.log_level : null
      dnsBase  = try(var.spec.dns_base, "") != "" ? var.spec.dns_base : null

      # Rendered on presence — an explicit true re-states the chart
      # default harmlessly, an explicit false is the actual opt-out.
      parallelRecoveryEnabled = try(var.spec.parallel_recovery_enabled, null)

      # Plain bool (no presence): false IS the chart default, so only
      # true renders.
      pprofEndpointsEnabled = try(var.spec.pprof_endpoints_enabled, false) ? true : null

      resources = local.manager_resources != null && length(local.manager_resources) > 0 ? local.manager_resources : null

      # Pull secrets ride the chart's manager.imagePullSecrets list (raw
      # Kubernetes object list, piped into the pod spec).
      imagePullSecrets = length(try(var.spec.image_pull_secrets, [])) > 0 ? [for s in var.spec.image_pull_secrets : { name = s }] : null

      image = length(local.manager_image) > 0 ? local.manager_image : null
    } : k => v if v != null
  }

  # ---- typed chart values (twin of the Pulumi module's buildHelmValues) ------
  # Chart-default-matching values render only on divergence, so the
  # rendered values stay minimal on both engines — with ONE deliberate
  # exception: installCRDs is pinned false unconditionally (see main.tf).
  typed_values = {
    for k, v in {
      # installCRDs: false ALWAYS — never conditional, never a spec knob.
      # The chart templates its ten CRDs release-owned with NO
      # keep-on-uninstall knob, so a Helm-owned install would
      # cascade-delete every OpenSearchCluster (and its data) on
      # uninstall. The module owns the CRDs instead
      # (kubectl_manifest.crds in main.tf).
      installCRDs = false

      # fullnameOverride pins the chart's fullname to the resource name
      # (the catalog's Helm-kind identity convention). Load-bearing:
      # without it the chart's default fullname pushes its metrics
      # Service name past Kubernetes' 63-character limit for ordinary
      # release names — verified live. Twin of the Pulumi module.
      fullnameOverride = local.release_name

      manager = length(local.manager_values) > 0 ? local.manager_values : null

      # The chart nests a single enable flag. Rendered on presence — an
      # explicit true re-states the chart default harmlessly, an
      # explicit false is the actual opt-out.
      # The sidecar's image repository is ALWAYS re-pointed at the
      # maintainer's own quay.io repository (same tag as the chart
      # pins): the chart's default, gcr.io/kubebuilder/kube-rbac-proxy,
      # was DELETED upstream (verified live), so a default-posture
      # install can never pull it. Twin of the Pulumi module;
      # overridable via helm_values.
      kubeRbacProxy = merge(
        { image = { repository = "quay.io/brancz/kube-rbac-proxy" } },
        try(var.spec.kube_rbac_proxy_enabled, null) != null ? { enable = var.spec.kube_rbac_proxy_enabled } : {}
      )

      # Plain bool (no presence): false IS the chart default
      # (ClusterRoleBindings), so only true renders. Spec CEL requires
      # watch_namespace alongside it.
      useRoleBindings = try(var.spec.use_role_bindings, false) ? true : null

      nodeSelector = length(try(var.spec.node_selector, {})) > 0 ? var.spec.node_selector : null
      tolerations = length(try(var.spec.tolerations, [])) > 0 ? [
        for t in var.spec.tolerations : {
          for tk, tv in {
            key               = try(t.key, "") != "" ? t.key : null
            operator          = try(t.operator, "") != "" ? t.operator : null
            value             = try(t.value, "") != "" ? t.value : null
            effect            = try(t.effect, "") != "" ? t.effect : null
            tolerationSeconds = try(t.toleration_seconds, null)
          } : tk => tv if tv != null
        }
      ] : null
    } : k => v if v != null
  }
}
