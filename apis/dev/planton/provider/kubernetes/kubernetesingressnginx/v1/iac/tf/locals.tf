# Computed values for the KubernetesIngressNginx module. Every resolution
# here has an exact twin in the Pulumi module's locals.go / values.go — keep
# them in lockstep.
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
  # manifest. The repository serves this chart as "ingress-nginx" — the
  # repo-prefixed spelling "kubernetes-ingress-nginx" does not exist in the
  # index and fails at install time.
  helm_chart_name = "ingress-nginx"
  helm_chart_repo = "https://kubernetes.github.io/ingress-nginx"

  # Release name = metadata.name, NOT a fixed chart name: multiple
  # controller instances per cluster (public + internal split) are a
  # first-class upstream pattern, so each manifest gets its own release.
  # The chart fullname is pinned to this too (typed_values), which also
  # isolates leader election per instance (electionID defaults to
  # "<fullname>-leader" in the chart).
  release_name = var.metadata.name

  # Chart version resolved to the pinned default when unset, so both engines
  # install the same chart whether or not the platform's defaulting
  # middleware ran — mirror of the Pulumi module's DefaultChartVersion.
  chart_version = coalesce(var.spec.chart_version, "4.15.1")

  namespace = var.spec.namespace

  # IngressClass this instance owns (spec default: "nginx").
  ingress_class_name = coalesce(try(var.spec.ingress_class.name, null), "nginx")

  # spec.controller of the IngressClass. Empty spec value derives: the chart
  # default "k8s.io/ingress-nginx" for class "nginx", otherwise
  # "k8s.io/<class-name>" so additional controllers isolate automatically
  # without the user inventing a vocabulary.
  ingress_class_controller_value = (
    try(var.spec.ingress_class.controller_value, "") != ""
    ? var.spec.ingress_class.controller_value
    : (local.ingress_class_name == "nginx" ? "k8s.io/ingress-nginx" : "k8s.io/${local.ingress_class_name}")
  )

  # Deterministic chart-derived object names ("<fullname>-controller" and
  # its "-internal" sibling) — what verification and downstream composition
  # key off.
  controller_service_name = "${var.metadata.name}-controller"
  internal_enabled        = try(var.spec.service.internal.enabled, false)
  internal_service_name   = local.internal_enabled ? "${var.metadata.name}-controller-internal" : ""

  service_type_map = {
    "load_balancer" = "LoadBalancer"
    "node_port"     = "NodePort"
    "cluster_ip"    = "ClusterIP"
  }
  # LoadBalancer is both the spec default and the chart default.
  service_type     = lookup(local.service_type_map, try(var.spec.service.type, "load_balancer"), "LoadBalancer")
  is_load_balancer = local.service_type == "LoadBalancer"

  autoscaling_enabled = try(var.spec.autoscaling.enabled, false)

  # Admission-webhook values that differ from the chart defaults (twin of
  # the Pulumi module's `admission` map — see the usage site for why this
  # is length-gated).
  admission_webhooks_pruned = try(var.spec.admission_webhooks, null) == null ? {} : {
    for ak, av in {
      enabled        = try(var.spec.admission_webhooks.enabled, null) == false ? false : null
      failurePolicy  = try(var.spec.admission_webhooks.failure_policy, "fail") == "ignore" ? "Ignore" : null
      timeoutSeconds = try(var.spec.admission_webhooks.timeout_seconds, null)
    } : ak => av if av != null
  }

  # Resource-identity labels stamped on the namespace this module creates
  # (never injected into the chart's own resources — Helm owns those).
  labels = merge(
    {
      "planton.ai/resource"      = "true"
      "planton.ai/resource-name" = var.metadata.name
      "planton.ai/resource-kind" = "KubernetesIngressNginx"
    },
    var.metadata.id != null && var.metadata.id != "" ? { "planton.ai/resource-id" = var.metadata.id } : {},
    var.metadata.org != null && var.metadata.org != "" ? { "planton.ai/organization" = var.metadata.org } : {},
    var.metadata.env != null && var.metadata.env != "" ? { "planton.ai/environment" = var.metadata.env } : {}
  )

  # ---- shared ContainerResources shape ----------------------------------
  controller_resources = try(var.spec.resources, null) == null ? null : {
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

  default_backend_resources = try(var.spec.default_backend.resources, null) == null ? null : {
    for k, v in {
      limits = try(var.spec.default_backend.resources.limits, null) == null ? null : {
        for lk, lv in {
          cpu    = try(var.spec.default_backend.resources.limits.cpu, "") != "" ? var.spec.default_backend.resources.limits.cpu : null
          memory = try(var.spec.default_backend.resources.limits.memory, "") != "" ? var.spec.default_backend.resources.limits.memory : null
        } : lk => lv if lv != null
      }
      requests = try(var.spec.default_backend.resources.requests, null) == null ? null : {
        for rk, rv in {
          cpu    = try(var.spec.default_backend.resources.requests.cpu, "") != "" ? var.spec.default_backend.resources.requests.cpu : null
          memory = try(var.spec.default_backend.resources.requests.memory, "") != "" ? var.spec.default_backend.resources.requests.memory : null
        } : rk => rv if rv != null
      }
    } : k => v if v != null && length(v) > 0
  }

  # ---- default TLS certificate flag --------------------------------------
  # The chart exposes no first-class key for the default certificate;
  # upstream's own documented mechanism is the controller flag. Namespace
  # defaults to the installation namespace.
  default_tls_secret = try(var.spec.default_tls_certificate.secret_name, "")
  default_tls_ref = local.default_tls_secret != "" ? format(
    "%s/%s",
    try(var.spec.default_tls_certificate.namespace, "") != "" ? var.spec.default_tls_certificate.namespace : local.namespace,
    local.default_tls_secret
  ) : ""

  # ---- controller service values ------------------------------------------
  service_values = {
    for k, v in {
      # Set the type explicitly even at the default so rendered values are
      # self-describing.
      type        = local.service_type
      annotations = length(try(var.spec.service.annotations, {})) > 0 ? var.spec.service.annotations : null
      # "local" preserves client source IPs; "cluster" is the Kubernetes
      # default, so only Local needs rendering.
      externalTrafficPolicy    = try(var.spec.service.external_traffic_policy, "cluster") == "local" ? "Local" : null
      loadBalancerSourceRanges = length(try(var.spec.service.load_balancer_source_ranges, [])) > 0 ? var.spec.service.load_balancer_source_ranges : null
      loadBalancerClass        = try(var.spec.service.load_balancer_class, "") != "" ? var.spec.service.load_balancer_class : null
      enableHttp               = try(var.spec.service.enable_http, null) == false ? false : null
      enableHttps              = try(var.spec.service.enable_https, null) == false ? false : null
      nodePorts = (
        try(var.spec.service.http_node_port, null) != null || try(var.spec.service.https_node_port, null) != null
        ) ? {
        for nk, nv in {
          http  = try(var.spec.service.http_node_port, null)
          https = try(var.spec.service.https_node_port, null)
        } : nk => nv if nv != null
      } : null
      internal = local.internal_enabled ? {
        enabled     = true
        annotations = try(var.spec.service.internal.annotations, {})
      } : null
    } : k => v if v != null
  }

  # ---- controller block (twin of the Pulumi module's buildHelmValues) ------
  controller_values = {
    for k, v in {
      # The IngressClass is the instance's identity — always rendered. The
      # legacy ingress.class-annotation vocabulary must track the class name
      # too, or annotation-based Ingresses mis-route on non-default-named
      # instances.
      ingressClassResource = {
        for ik, iv in {
          name            = local.ingress_class_name
          controllerValue = local.ingress_class_controller_value
          default         = try(var.spec.ingress_class.is_default_class, false) ? true : null
        } : ik => iv if iv != null
      }
      ingressClass             = local.ingress_class_name
      watchIngressWithoutClass = try(var.spec.ingress_class.watch_ingress_without_class, false) ? true : null

      # Replicas vs autoscaling: when the HPA owns the count, do not also
      # pin replicaCount — the chart's Deployment template omits replicas
      # under autoscaling, avoiding a rollout tug-of-war.
      replicaCount = local.autoscaling_enabled ? null : try(var.spec.replicas, null)
      autoscaling = local.autoscaling_enabled ? {
        for ak, av in {
          enabled                           = true
          minReplicas                       = try(var.spec.autoscaling.min_replicas, null)
          maxReplicas                       = try(var.spec.autoscaling.max_replicas, null)
          targetCPUUtilizationPercentage    = try(var.spec.autoscaling.target_cpu_utilization_percent, null)
          targetMemoryUtilizationPercentage = try(var.spec.autoscaling.target_memory_utilization_percent, null)
        } : ak => av if av != null
      } : null

      resources = local.controller_resources
      service   = local.service_values

      kind = try(var.spec.controller_kind, "deployment") == "daemon_set" ? "DaemonSet" : null
      # Keep in-cluster name resolution working from the host network —
      # without this the controller resolves through the node's DNS and
      # cannot see cluster Services.
      hostNetwork = var.spec.host_network ? true : null
      dnsPolicy   = var.spec.host_network ? "ClusterFirstWithHostNet" : null
      hostPort    = var.spec.host_ports ? { enabled = true } : null

      config                  = length(var.spec.nginx_config) > 0 ? var.spec.nginx_config : null
      allowSnippetAnnotations = var.spec.allow_snippet_annotations ? true : null
      extraArgs               = local.default_tls_ref != "" ? { "default-ssl-certificate" = local.default_tls_ref } : null

      # Gated on the PRUNED map being non-empty, not on block presence: a
      # block carrying only chart defaults prunes to {}, and emitting an
      # empty map would diverge from the Pulumi twin (which omits the key)
      # — harmless to the chart, but rendered values must stay
      # byte-identical across engines.
      admissionWebhooks = length(local.admission_webhooks_pruned) > 0 ? local.admission_webhooks_pruned : null

      metrics = try(var.spec.metrics.enabled, false) ? {
        for mk, mv in {
          enabled = true
          serviceMonitor = try(var.spec.metrics.service_monitor, false) ? {
            for sk, sv in {
              enabled          = true
              scrapeInterval   = try(var.spec.metrics.service_monitor_interval, null)
              additionalLabels = length(try(var.spec.metrics.service_monitor_labels, {})) > 0 ? var.spec.metrics.service_monitor_labels : null
            } : sk => sv if sv != null
          } : null
        } : mk => mv if mv != null
      } : null

      nodeSelector = length(var.spec.node_selector) > 0 ? var.spec.node_selector : null
      tolerations = length(var.spec.tolerations) > 0 ? [
        for t in var.spec.tolerations : {
          for tk, tv in {
            key               = t.key != "" ? t.key : null
            operator          = t.operator != "" ? t.operator : null
            value             = t.value != "" ? t.value : null
            effect            = t.effect != "" ? t.effect : null
            tolerationSeconds = try(t.toleration_seconds, null)
          } : tk => tv if tv != null
        }
      ] : null
      priorityClassName = var.spec.priority_class_name != "" ? var.spec.priority_class_name : null
    } : k => v if v != null
  }

  # ---- default backend --------------------------------------------------------
  default_backend_image = try(var.spec.default_backend.image, "")
  default_backend_values = try(var.spec.default_backend.enabled, false) ? {
    for k, v in {
      enabled      = true
      replicaCount = try(var.spec.default_backend.replicas, null)
      # Spec carries "repository:tag"; the chart wants them split.
      image = local.default_backend_image != "" ? {
        for ik, iv in {
          repository = (
            length(regexall(":", local.default_backend_image)) > 0
            ? join(":", slice(split(":", local.default_backend_image), 0, length(split(":", local.default_backend_image)) - 1))
            : local.default_backend_image
          )
          tag = (
            length(regexall(":", local.default_backend_image)) > 0
            ? element(split(":", local.default_backend_image), length(split(":", local.default_backend_image)) - 1)
            : null
          )
        } : ik => iv if iv != null
      } : null
      resources = local.default_backend_resources
    } : k => v if v != null
  } : null

  # ---- typed chart values (twin of the Pulumi module's buildHelmValues) -----
  typed_values = {
    for k, v in {
      # Pin the chart's fullname to the release name (= metadata.name):
      # every chart object then carries a deterministic, manifest-derived
      # name — what verification, imports, and multi-instance coexistence
      # (including per-instance leader election) all key off.
      fullnameOverride = local.release_name
      controller       = local.controller_values
      global           = var.spec.image_registry != "" ? { image = { registry = var.spec.image_registry } } : null
      defaultBackend   = local.default_backend_values
      tcp              = length(var.spec.tcp_services) > 0 ? var.spec.tcp_services : null
      udp              = length(var.spec.udp_services) > 0 ? var.spec.udp_services : null
    } : k => v if v != null
  }
}
