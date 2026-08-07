# Computed values for the KubernetesIstio module. Every resolution here has an
# exact twin in the Pulumi module's locals.go / values.go — keep them in
# lockstep.
#
# HCL DISCIPLINE (applies to every conditional object in this file):
# conditional entries are written as `key = cond ? value : null` inside ONE
# object literal, pruned with `{ for k, v in {...} : k => v if v != null }`.
# The two tempting alternatives are both broken: `cond ? {...} : {}` ternaries
# fail plan-time type unification when branches carry different attributes,
# and `merge(concat(cond ? [{...}] : [], ...)...)` silently UNIFIES
# primitive-only sibling objects into map(string) — numbers and booleans
# arrive in the chart values as strings. The null-prune form preserves every
# value's type.
#
# Optional nested blocks are read with try(): HCL's && does NOT short-circuit,
# so chained null checks still dereference the null.

locals {
  # Chart identity — must stay byte-identical with the Pulumi module's vars:
  # cross-engine chart-name drift deploys two different products from one
  # manifest. Istio versions its charts in lockstep with the product, so one
  # version pin drives all releases.
  helm_chart_repo = "https://istio-release.storage.googleapis.com/charts"

  # Version resolved to the pinned default when unset, so both engines
  # install the same charts whether or not the platform's defaulting
  # middleware ran — mirror of the Pulumi module's DefaultVersion.
  version = coalesce(var.spec.version, "1.30.3")

  namespace = var.spec.namespace

  # The spec's revision normalized to "" when null/absent — a typed optional
  # attribute with no default is NULL when unset, and `try(x, "")` does NOT
  # normalize that (try only catches errors, and reading an existing null
  # attribute is not an error). Not coalesce(): it rejects an all-null/empty
  # argument list, and "" is exactly the value we want for "unset".
  spec_revision = var.spec.revision != null ? var.spec.revision : ""

  # Control-plane revision: spec value or "default" (the chart's own
  # vocabulary for the unnamed revision).
  revision = local.spec_revision != "" ? local.spec_revision : "default"

  # istiod release name gains a "-<revision>" suffix for a named revision,
  # matching how the chart names the istiod Deployment/Service themselves.
  istiod_release_name = local.spec_revision != "" ? "istiod-${local.spec_revision}" : "istiod"

  # The istiod Service name equals the release-derived resource name —
  # exported as the discovery-address handle.
  istiod_service_name = local.istiod_release_name

  # Data plane mode: the cni and ztunnel releases install only in ambient
  # mode (plus cni when spec.cni.enabled in sidecar mode).
  ambient     = try(var.spec.dataplane_mode, "") == "ambient"
  install_cni = local.ambient || try(var.spec.cni.enabled, false)

  # Resolved trust domain (spec value or the upstream default) — exported for
  # AuthorizationPolicy principal authoring.
  trust_domain = coalesce(try(var.spec.mesh_config.trust_domain, null), "cluster.local")

  dataplane_mode = coalesce(var.spec.dataplane_mode, "sidecar")

  # Resource-identity labels stamped on the namespace this module creates
  # (never injected into the charts' own resources — Helm owns those).
  labels = merge(
    {
      "planton.ai/resource"      = "true"
      "planton.ai/resource-name" = var.metadata.name
      "planton.ai/resource-kind" = "KubernetesIstio"
    },
    var.metadata.id != null && var.metadata.id != "" ? { "planton.ai/resource-id" = var.metadata.id } : {},
    var.metadata.org != null && var.metadata.org != "" ? { "planton.ai/organization" = var.metadata.org } : {},
    var.metadata.env != null && var.metadata.env != "" ? { "planton.ai/environment" = var.metadata.env } : {}
  )

  # ---- shared renderings --------------------------------------------------

  # ContainerResources -> chart resources shape (null when nothing set).
  istiod_resources = try(var.spec.istiod.resources, null) == null ? null : {
    for k, v in {
      limits = try(var.spec.istiod.resources.limits, null) == null ? null : {
        for lk, lv in {
          cpu    = try(var.spec.istiod.resources.limits.cpu, "") != "" ? var.spec.istiod.resources.limits.cpu : null
          memory = try(var.spec.istiod.resources.limits.memory, "") != "" ? var.spec.istiod.resources.limits.memory : null
        } : lk => lv if lv != null
      }
      requests = try(var.spec.istiod.resources.requests, null) == null ? null : {
        for rk, rv in {
          cpu    = try(var.spec.istiod.resources.requests.cpu, "") != "" ? var.spec.istiod.resources.requests.cpu : null
          memory = try(var.spec.istiod.resources.requests.memory, "") != "" ? var.spec.istiod.resources.requests.memory : null
        } : rk => rv if rv != null
      }
    } : k => v if v != null && length(v) > 0
  }

  proxy_resources = try(var.spec.proxy.resources, null) == null ? null : {
    for k, v in {
      limits = try(var.spec.proxy.resources.limits, null) == null ? null : {
        for lk, lv in {
          cpu    = try(var.spec.proxy.resources.limits.cpu, "") != "" ? var.spec.proxy.resources.limits.cpu : null
          memory = try(var.spec.proxy.resources.limits.memory, "") != "" ? var.spec.proxy.resources.limits.memory : null
        } : lk => lv if lv != null
      }
      requests = try(var.spec.proxy.resources.requests, null) == null ? null : {
        for rk, rv in {
          cpu    = try(var.spec.proxy.resources.requests.cpu, "") != "" ? var.spec.proxy.resources.requests.cpu : null
          memory = try(var.spec.proxy.resources.requests.memory, "") != "" ? var.spec.proxy.resources.requests.memory : null
        } : rk => rv if rv != null
      }
    } : k => v if v != null && length(v) > 0
  }

  ztunnel_resources = try(var.spec.ztunnel.resources, null) == null ? null : {
    for k, v in {
      limits = try(var.spec.ztunnel.resources.limits, null) == null ? null : {
        for lk, lv in {
          cpu    = try(var.spec.ztunnel.resources.limits.cpu, "") != "" ? var.spec.ztunnel.resources.limits.cpu : null
          memory = try(var.spec.ztunnel.resources.limits.memory, "") != "" ? var.spec.ztunnel.resources.limits.memory : null
        } : lk => lv if lv != null
      }
      requests = try(var.spec.ztunnel.resources.requests, null) == null ? null : {
        for rk, rv in {
          cpu    = try(var.spec.ztunnel.resources.requests.cpu, "") != "" ? var.spec.ztunnel.resources.requests.cpu : null
          memory = try(var.spec.ztunnel.resources.requests.memory, "") != "" ? var.spec.ztunnel.resources.requests.memory : null
        } : rk => rv if rv != null
      }
    } : k => v if v != null && length(v) > 0
  }

  # Tolerations pass through with null-pruned optional keys.
  istiod_tolerations = length(try(var.spec.istiod.tolerations, [])) == 0 ? null : [
    for t in var.spec.istiod.tolerations : {
      for k, v in {
        key               = try(t.key, "") != "" ? t.key : null
        operator          = try(t.operator, "") != "" ? t.operator : null
        value             = try(t.value, "") != "" ? t.value : null
        effect            = try(t.effect, "") != "" ? t.effect : null
        tolerationSeconds = try(t.toleration_seconds, null)
      } : k => v if v != null
    }
  ]

  # Shared image-source knobs (base/istiod/cni read them under global.*;
  # ztunnel reads them top-level).
  image_global_values = {
    for k, v in {
      hub              = try(var.spec.images.hub, "") != "" ? var.spec.images.hub : null
      variant          = try(var.spec.images.variant, null)
      imagePullSecrets = length(try(var.spec.images.image_pull_secrets, [])) > 0 ? var.spec.images.image_pull_secrets : null
    } : k => v if v != null
  }

  # The istio/base CRDs-only bundle at the pinned version (module-owned CRDs
  # — see main.tf).
  crd_bundle_url = "https://raw.githubusercontent.com/istio/istio/${local.version}/manifests/charts/base/files/crd-all.gen.yaml"

  # Every CRD in the pinned bundle, handed to the base chart as
  # base.excludedCRDs so the release templates NO CRDs (module-owned above).
  # Pinned knowledge of the 1.30.3 bundle — reconcile together with the
  # version pin (twin of the Pulumi module's CrdNames).
  crd_names = [
    "authorizationpolicies.security.istio.io",
    "destinationrules.networking.istio.io",
    "envoyfilters.networking.istio.io",
    "gateways.networking.istio.io",
    "peerauthentications.security.istio.io",
    "proxyconfigs.networking.istio.io",
    "requestauthentications.security.istio.io",
    "serviceentries.networking.istio.io",
    "sidecars.networking.istio.io",
    "telemetries.telemetry.istio.io",
    "trafficextensions.extensions.istio.io",
    "virtualservices.networking.istio.io",
    "wasmplugins.extensions.istio.io",
    "workloadentries.networking.istio.io",
    "workloadgroups.networking.istio.io",
  ]

  # ---- base values ---------------------------------------------------------
  base_typed_values = {
    for k, v in {
      # The default-revision validating webhook must point at THIS control
      # plane's revision ("default" for the unnamed revision).
      defaultRevision = local.revision
      # Exclude the ENTIRE CRD bundle from the release (module-owned CRDs).
      base = {
        excludedCRDs = local.crd_names
      }
      global = local.namespace != "istio-system" ? {
        istioNamespace = local.namespace
      } : null
    } : k => v if v != null
  }

  # ---- istiod values --------------------------------------------------------
  istiod_proxy_values = {
    for pk, pv in {
      resources     = local.proxy_resources
      logLevel      = try(var.spec.proxy.log_level, null)
      autoInject    = try(var.spec.proxy.auto_inject, null)
      clusterDomain = try(var.spec.proxy.cluster_domain, null) != null && try(var.spec.proxy.cluster_domain, "") != "cluster.local" ? var.spec.proxy.cluster_domain : null
    } : pk => pv if pv != null
  }

  istiod_global_values = {
    for k, v in {
      istioNamespace = local.namespace != "istio-system" ? local.namespace : null
      logging = try(var.spec.istiod.log_level, null) != null && try(var.spec.istiod.log_level, "") != "default:info" ? {
        level = var.spec.istiod.log_level
      } : null
      defaultPodDisruptionBudget = try(var.spec.istiod.pod_disruption_budget, null) != null ? {
        enabled = var.spec.istiod.pod_disruption_budget
      } : null
      priorityClassName = try(var.spec.istiod.priority_class_name, "") != "" ? var.spec.istiod.priority_class_name : null
      multiCluster = try(var.spec.mesh_config.cluster_name, "") != "" ? {
        clusterName = var.spec.mesh_config.cluster_name
      } : null
      network          = try(var.spec.mesh_config.network, "") != "" ? var.spec.mesh_config.network : null
      meshID           = try(var.spec.mesh_config.mesh_id, "") != "" ? var.spec.mesh_config.mesh_id : null
      proxy            = length(local.istiod_proxy_values) > 0 ? local.istiod_proxy_values : null
      hub              = try(local.image_global_values.hub, null)
      variant          = try(local.image_global_values.variant, null)
      imagePullSecrets = try(local.image_global_values.imagePullSecrets, null)
    } : k => v if v != null
  }

  istiod_sidecar_injector_values = {
    for sk, sv in {
      enableNamespacesByDefault = try(var.spec.sidecar_injector.enable_namespaces_by_default, false) ? true : null
      rewriteAppHTTPProbe       = try(var.spec.sidecar_injector.rewrite_app_http_probe, null)
    } : sk => sv if sv != null
  }

  istiod_mesh_config_values = {
    for k, v in {
      trustDomain = try(var.spec.mesh_config.trust_domain, null) != null && try(var.spec.mesh_config.trust_domain, "") != "cluster.local" ? var.spec.mesh_config.trust_domain : null
      outboundTrafficPolicy = try(var.spec.mesh_config.outbound_traffic_policy_mode, null) != null ? {
        mode = var.spec.mesh_config.outbound_traffic_policy_mode
      } : null
      accessLogFile         = try(var.spec.mesh_config.access_log_file, "") != "" ? var.spec.mesh_config.access_log_file : null
      enablePrometheusMerge = try(var.spec.mesh_config.enable_prometheus_merge, null)
    } : k => v if v != null
  }

  istiod_typed_values = {
    for k, v in {
      # Ambient mode rides the chart's own profile overlay (the upstream
      # install path: --set profile=ambient).
      profile      = local.ambient ? "ambient" : null
      revision     = local.spec_revision != "" ? local.spec_revision : null
      replicaCount = try(var.spec.istiod.replicas, null)
      # A fixed replica count only holds if the chart's HPA is off (the spec
      # forbids replicas+autoscale, so this is safe); otherwise honor the
      # typed autoscale block.
      autoscaleEnabled = try(var.spec.istiod.replicas, null) != null ? false : try(var.spec.istiod.autoscale.enabled, null)
      autoscaleMin     = try(var.spec.istiod.autoscale.min_replicas, null)
      autoscaleMax     = try(var.spec.istiod.autoscale.max_replicas, null)
      cpu = try(var.spec.istiod.autoscale.target_cpu_utilization_percent, null) != null ? {
        targetAverageUtilization = var.spec.istiod.autoscale.target_cpu_utilization_percent
      } : null
      resources              = local.istiod_resources
      nodeSelector           = length(try(var.spec.istiod.node_selector, {})) > 0 ? var.spec.istiod.node_selector : null
      tolerations            = local.istiod_tolerations
      meshConfig             = length(local.istiod_mesh_config_values) > 0 ? local.istiod_mesh_config_values : null
      global                 = length(local.istiod_global_values) > 0 ? local.istiod_global_values : null
      sidecarInjectorWebhook = length(local.istiod_sidecar_injector_values) > 0 ? local.istiod_sidecar_injector_values : null
      # In sidecar mode with the node-level CNI agent installed, istiod's
      # injector must emit pods that rely on it instead of the privileged
      # istio-init container. (Ambient's profile overlay handles this itself.)
      cni = !local.ambient && local.install_cni ? { enabled = true } : null
      # The per-GatewayClass deployment overlay for istiod's Gateway API
      # auto-provisioning.
      gatewayClasses = try(var.spec.gateway_defaults.service_type, null) != null ? {
        istio = { service = { spec = { type = var.spec.gateway_defaults.service_type } } }
      } : null
    } : k => v if v != null
  }

  # ---- cni values -------------------------------------------------------------
  cni_typed_values = {
    for k, v in {
      profile           = local.ambient ? "ambient" : null
      revision          = local.spec_revision != "" ? local.spec_revision : null
      excludeNamespaces = length(try(var.spec.cni.exclude_namespaces, [])) > 0 ? var.spec.cni.exclude_namespaces : null
      cniBinDir         = try(var.spec.cni.cni_bin_dir, null)
      cniConfDir        = try(var.spec.cni.cni_conf_dir, null)
      chained           = try(var.spec.cni.chained, null)
      global            = length(local.image_global_values) > 0 ? local.image_global_values : null
    } : k => v if v != null
  }

  # ---- ztunnel values -----------------------------------------------------------
  # The ztunnel chart reads hub/variant/imagePullSecrets at the top level,
  # unlike the global.* convention of the other three charts.
  ztunnel_typed_values = {
    for k, v in {
      revision         = local.spec_revision != "" ? local.spec_revision : null
      istioNamespace   = local.namespace != "istio-system" ? local.namespace : null
      resources        = local.ztunnel_resources
      logLevel         = try(var.spec.ztunnel.log_level, null)
      hub              = try(local.image_global_values.hub, null)
      variant          = try(local.image_global_values.variant, null)
      imagePullSecrets = try(local.image_global_values.imagePullSecrets, null)
    } : k => v if v != null
  }
}
