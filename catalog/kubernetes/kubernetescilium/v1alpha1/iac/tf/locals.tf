# Computed values for the KubernetesCilium module. Every resolution here has
# an exact twin in the Pulumi module's locals.go / values.go — keep them in
# lockstep.
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
# short-circuit, so chained null checks still dereference the null.

locals {
  # Chart identity — must stay byte-identical with the Pulumi module's vars:
  # cross-engine chart-name drift deploys two different products from one
  # manifest.
  helm_chart_name = "cilium"
  helm_chart_repo = "https://helm.cilium.io"

  # Release name FIXED to the chart name: Cilium is the node dataplane —
  # the agent DaemonSet, operator, and generated CNI configuration are
  # cluster singletons, so one dataplane per cluster is an upstream
  # constraint.
  release_name = local.helm_chart_name

  # Chart version resolved to the pinned default when unset, so both engines
  # install the same chart whether or not the platform's defaulting
  # middleware ran — mirror of the Pulumi module's DefaultChartVersion.
  chart_version = coalesce(var.spec.chart_version, "1.19.6")

  namespace = var.spec.namespace

  # Cluster identity resolved to the chart's own default — the name this
  # cluster carries in Hubble flows and any future Cluster Mesh.
  cluster_name = coalesce(var.spec.cluster_name, "default")

  # Fixed chart-template names for the outputs (verified against the chart
  # templates — hubble-relay/hubble-ui Services and the GatewayClass carry
  # no release-derived prefix): set when the component exists, empty
  # otherwise.
  hubble_relay_service_name = try(var.spec.hubble.relay, false) ? "hubble-relay" : ""
  hubble_ui_service_name    = try(var.spec.hubble.ui, false) ? "hubble-ui" : ""
  gateway_class_name        = var.spec.gateway_api ? "cilium" : ""

  # Resource-identity labels stamped on the namespace this module creates
  # (never injected into the chart's own resources — Helm owns those).
  labels = merge(
    {
      "planton.ai/resource"      = "true"
      "planton.ai/resource-name" = var.metadata.name
      "planton.ai/resource-kind" = "KubernetesCilium"
    },
    var.metadata.id != null && var.metadata.id != "" ? { "planton.ai/resource-id" = var.metadata.id } : {},
    var.metadata.org != null && var.metadata.org != "" ? { "planton.ai/organization" = var.metadata.org } : {},
    var.metadata.env != null && var.metadata.env != "" ? { "planton.ai/environment" = var.metadata.env } : {}
  )

  # ---- IPAM -----------------------------------------------------------------
  ipam_operator = {
    for k, v in {
      clusterPoolIPv4PodCIDRList = length(try(var.spec.ipam.cluster_pool_ipv4_pod_cidrs, [])) > 0 ? var.spec.ipam.cluster_pool_ipv4_pod_cidrs : null
      clusterPoolIPv4MaskSize    = try(var.spec.ipam.cluster_pool_ipv4_mask_size, null)
    } : k => v if v != null
  }
  ipam_values_raw = {
    for k, v in {
      mode     = try(var.spec.ipam.mode, null) != null && try(var.spec.ipam.mode, "") != "" ? var.spec.ipam.mode : null
      operator = length(local.ipam_operator) > 0 ? local.ipam_operator : null
    } : k => v if v != null
  }
  ipam_values = length(local.ipam_values_raw) > 0 ? local.ipam_values_raw : null

  # ---- CNI installation / chaining --------------------------------------------
  cni_values_raw = {
    for k, v in {
      chainingMode   = try(var.spec.cni.chaining_mode, null) != null && try(var.spec.cni.chaining_mode, "") != "" ? var.spec.cni.chaining_mode : null
      chainingTarget = try(var.spec.cni.chaining_target, "") != "" ? var.spec.cni.chaining_target : null
      # Optional bool: presence (not truth) decides rendering — an explicit
      # false is exactly the value chaining setups must send (the CEL rule
      # enforces it), while unset keeps the chart default (true).
      exclusive = try(var.spec.cni.exclusive, null)
    } : k => v if v != null
  }
  cni_values = length(local.cni_values_raw) > 0 ? local.cni_values_raw : null

  # ---- Hubble observability ------------------------------------------------------
  # hubble.metrics.enabled is upstream's LIST of metric families (null
  # disables) — not a boolean despite the name.
  hubble_metrics = length(try(var.spec.hubble.metrics, [])) > 0 ? {
    for k, v in {
      enabled        = var.spec.hubble.metrics
      serviceMonitor = try(var.spec.hubble.metrics_service_monitor, false) ? { enabled = true } : null
    } : k => v if v != null
  } : null
  hubble_values_raw = {
    for k, v in {
      # Chart default is enabled=true, so only an EXPLICIT false is
      # rendered (an explicit true is the default — nothing to say).
      enabled = try(var.spec.hubble.enabled, null) == false ? false : null
      relay   = try(var.spec.hubble.relay, false) ? { enabled = true } : null
      ui      = try(var.spec.hubble.ui, false) ? { enabled = true } : null
      metrics = local.hubble_metrics
    } : k => v if v != null
  }
  hubble_values = length(local.hubble_values_raw) > 0 ? local.hubble_values_raw : null

  # ---- transparent encryption --------------------------------------------------------
  encryption_values = try(var.spec.encryption.enabled, false) ? {
    for k, v in {
      enabled        = true
      type           = try(var.spec.encryption.type, null) != null && try(var.spec.encryption.type, "") != "" ? var.spec.encryption.type : null
      nodeEncryption = try(var.spec.encryption.node_encryption, false) ? true : null
    } : k => v if v != null
  } : null

  # ---- shared ContainerResources shape --------------------------------------
  # Top-level resources = the agent container (the cilium DaemonSet).
  agent_resources = try(var.spec.agent_resources, null) == null ? null : {
    for k, v in {
      limits = try(var.spec.agent_resources.limits, null) == null ? null : {
        for lk, lv in {
          cpu    = try(var.spec.agent_resources.limits.cpu, "") != "" ? var.spec.agent_resources.limits.cpu : null
          memory = try(var.spec.agent_resources.limits.memory, "") != "" ? var.spec.agent_resources.limits.memory : null
        } : lk => lv if lv != null
      }
      requests = try(var.spec.agent_resources.requests, null) == null ? null : {
        for rk, rv in {
          cpu    = try(var.spec.agent_resources.requests.cpu, "") != "" ? var.spec.agent_resources.requests.cpu : null
          memory = try(var.spec.agent_resources.requests.memory, "") != "" ? var.spec.agent_resources.requests.memory : null
        } : rk => rv if rv != null
      }
    } : k => v if v != null && length(v) > 0
  }

  operator_resources = try(var.spec.operator.resources, null) == null ? null : {
    for k, v in {
      limits = try(var.spec.operator.resources.limits, null) == null ? null : {
        for lk, lv in {
          cpu    = try(var.spec.operator.resources.limits.cpu, "") != "" ? var.spec.operator.resources.limits.cpu : null
          memory = try(var.spec.operator.resources.limits.memory, "") != "" ? var.spec.operator.resources.limits.memory : null
        } : lk => lv if lv != null
      }
      requests = try(var.spec.operator.resources.requests, null) == null ? null : {
        for rk, rv in {
          cpu    = try(var.spec.operator.resources.requests.cpu, "") != "" ? var.spec.operator.resources.requests.cpu : null
          memory = try(var.spec.operator.resources.requests.memory, "") != "" ? var.spec.operator.resources.requests.memory : null
        } : rk => rv if rv != null
      }
    } : k => v if v != null && length(v) > 0
  }

  # ---- operator (sizing + telemetry merged into ONE map) --------------------------------
  # The operator map collects keys from TWO spec arms (operator sizing and
  # prometheus) — built once so the arms merge into one map instead of the
  # later overwriting the earlier. Twin of the Pulumi module's
  # operatorValues.
  operator_values_raw = {
    for k, v in {
      replicas  = try(var.spec.operator.replicas, null)
      resources = local.operator_resources
      prometheus = try(var.spec.prometheus.enabled, false) ? {
        for pk, pv in {
          enabled        = true
          serviceMonitor = try(var.spec.prometheus.service_monitor, false) ? { enabled = true } : null
        } : pk => pv if pv != null
      } : null
    } : k => v if v != null
  }
  operator_values = length(local.operator_values_raw) > 0 ? local.operator_values_raw : null

  # ---- typed chart values (twin of the Pulumi module's buildHelmValues) ------
  # No fullnameOverride here (unlike sibling modules): the cilium chart
  # names its workloads with FIXED names — DaemonSet "cilium", Deployment
  # "cilium-operator" — regardless of the release name, so there is nothing
  # to pin.
  typed_values = {
    for k, v in {
      cluster = try(var.spec.cluster_name, null) != null && try(var.spec.cluster_name, "") != "" ? { name = var.spec.cluster_name } : null

      ipam = local.ipam_values

      routingMode           = try(var.spec.routing.mode, null) != null && try(var.spec.routing.mode, "") != "" ? var.spec.routing.mode : null
      tunnelProtocol        = try(var.spec.routing.tunnel_protocol, null) != null && try(var.spec.routing.tunnel_protocol, "") != "" ? var.spec.routing.tunnel_protocol : null
      ipv4NativeRoutingCIDR = try(var.spec.routing.ipv4_native_routing_cidr, "") != "" ? var.spec.routing.ipv4_native_routing_cidr : null
      autoDirectNodeRoutes  = try(var.spec.routing.auto_direct_node_routes, false) ? true : null

      # TRAP: kubeProxyReplacement is a STRING in the chart's values.yaml
      # (historically it took "strict"/"partial"; today "true"/"false") —
      # the string keeps the rendered document byte-identical with what the
      # chart declares and with the Pulumi module. Only rendered when true
      # (chart default is "false").
      kubeProxyReplacement = var.spec.kube_proxy_replacement ? "true" : null

      k8sServiceHost = var.spec.k8s_service_host != "" ? var.spec.k8s_service_host : null
      # k8sServicePort is also a string in values.yaml (default ""), so the
      # number renders as its decimal string — the Pulumi twin uses
      # strconv.Itoa for the same reason.
      k8sServicePort = try(var.spec.k8s_service_port, null) != null ? tostring(var.spec.k8s_service_port) : null

      cni = local.cni_values

      # Cloud-provider datapath integrations. AWS ENI pairs with ipam mode
      # "eni" (pods draw VPC-routable IPs from ENIs Cilium manages).
      eni       = try(var.spec.cloud.aws_eni, false) ? { enabled = true } : null
      aksbyocni = try(var.spec.cloud.aks_byocni, false) ? { enabled = true } : null
      gke       = try(var.spec.cloud.gke, false) ? { enabled = true } : null

      hubble = local.hubble_values

      encryption = local.encryption_values

      policyEnforcementMode = try(var.spec.policy_enforcement_mode, null) != null && try(var.spec.policy_enforcement_mode, "") != "" ? var.spec.policy_enforcement_mode : null

      gatewayAPI = var.spec.gateway_api ? { enabled = true } : null

      bandwidthManager = try(var.spec.bandwidth_manager.enabled, false) ? {
        for bk, bv in {
          enabled = true
          bbr     = try(var.spec.bandwidth_manager.bbr, false) ? true : null
        } : bk => bv if bv != null
      } : null

      operator = local.operator_values

      resources = local.agent_resources

      # Cilium's own telemetry: one spec toggle drives BOTH components —
      # agent metrics here, operator metrics via operator.prometheus above
      # (exposing only one of the two would be a confusing half-telemetry
      # posture).
      prometheus = try(var.spec.prometheus.enabled, false) ? {
        for pk, pv in {
          enabled        = true
          serviceMonitor = try(var.spec.prometheus.service_monitor, false) ? { enabled = true } : null
        } : pk => pv if pv != null
      } : null
    } : k => v if v != null
  }
}
