# Computed values for the KubernetesCertManager module. Every resolution here
# has an exact twin in the Pulumi module's locals.go / values.go — keep them
# in lockstep.
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
  helm_chart_name = "cert-manager"
  helm_chart_repo = "https://charts.jetstack.io"

  # Release name fixed to the chart name: one cert-manager per cluster is an
  # upstream architectural constraint (cluster-scoped CRDs and webhooks).
  release_name = local.helm_chart_name

  # Chart version resolved to the pinned default when unset, so both engines
  # install the same chart whether or not the platform's defaulting
  # middleware ran — mirror of the Pulumi module's DefaultChartVersion.
  chart_version = coalesce(var.spec.chart_version, "v1.20.3")

  namespace = var.spec.namespace

  # The chart derives the controller ServiceAccount name from the release
  # name — with the release named "cert-manager" the ServiceAccount is
  # "cert-manager". Exported for cloud-side identity bindings.
  service_account_name = local.helm_chart_name

  # Resolved cluster-resource namespace: explicit spec value, else the
  # installation namespace (cert-manager's own default). Exported —
  # KubernetesClusterIssuer materializes credential Secrets here.
  cluster_resource_namespace = coalesce(var.spec.cluster_resource_namespace, local.namespace)

  # Resource-identity labels stamped on the namespace this module creates
  # (never injected into the chart's own resources — Helm owns those).
  labels = merge(
    {
      "planton.ai/resource"      = "true"
      "planton.ai/resource-name" = var.metadata.name
      "planton.ai/resource-kind" = "KubernetesCertManager"
    },
    var.metadata.id != null && var.metadata.id != "" ? { "planton.ai/resource-id" = var.metadata.id } : {},
    var.metadata.org != null && var.metadata.org != "" ? { "planton.ai/organization" = var.metadata.org } : {},
    var.metadata.env != null && var.metadata.env != "" ? { "planton.ai/environment" = var.metadata.env } : {}
  )

  # ---- workload identity annotations -----------------------------------
  # The chart creates the controller ServiceAccount; the identity annotation
  # rides serviceAccount.annotations. AKS additionally needs the
  # azure.workload.identity/use pod label.
  workload_identity_annotations = merge(
    try(var.spec.workload_identity.gke, null) != null ? {
      "iam.gke.io/gcp-service-account" = var.spec.workload_identity.gke.service_account_email
    } : {},
    try(var.spec.workload_identity.eks, null) != null ? {
      "eks.amazonaws.com/role-arn" = var.spec.workload_identity.eks.role_arn
    } : {},
    try(var.spec.workload_identity.aks, null) != null ? merge(
      {
        "azure.workload.identity/client-id" = var.spec.workload_identity.aks.client_id
      },
      try(var.spec.workload_identity.aks.tenant_id, null) != null ? {
        "azure.workload.identity/tenant-id" = var.spec.workload_identity.aks.tenant_id
      } : {}
    ) : {}
  )

  # ---- per-component resources (shared ContainerResources shape) --------
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

  webhook_resources = try(var.spec.webhook.resources, null) == null ? null : {
    for k, v in {
      limits = try(var.spec.webhook.resources.limits, null) == null ? null : {
        for lk, lv in {
          cpu    = try(var.spec.webhook.resources.limits.cpu, "") != "" ? var.spec.webhook.resources.limits.cpu : null
          memory = try(var.spec.webhook.resources.limits.memory, "") != "" ? var.spec.webhook.resources.limits.memory : null
        } : lk => lv if lv != null
      }
      requests = try(var.spec.webhook.resources.requests, null) == null ? null : {
        for rk, rv in {
          cpu    = try(var.spec.webhook.resources.requests.cpu, "") != "" ? var.spec.webhook.resources.requests.cpu : null
          memory = try(var.spec.webhook.resources.requests.memory, "") != "" ? var.spec.webhook.resources.requests.memory : null
        } : rk => rv if rv != null
      }
    } : k => v if v != null && length(v) > 0
  }

  cainjector_resources = try(var.spec.cainjector.resources, null) == null ? null : {
    for k, v in {
      limits = try(var.spec.cainjector.resources.limits, null) == null ? null : {
        for lk, lv in {
          cpu    = try(var.spec.cainjector.resources.limits.cpu, "") != "" ? var.spec.cainjector.resources.limits.cpu : null
          memory = try(var.spec.cainjector.resources.limits.memory, "") != "" ? var.spec.cainjector.resources.limits.memory : null
        } : lk => lv if lv != null
      }
      requests = try(var.spec.cainjector.resources.requests, null) == null ? null : {
        for rk, rv in {
          cpu    = try(var.spec.cainjector.resources.requests.cpu, "") != "" ? var.spec.cainjector.resources.requests.cpu : null
          memory = try(var.spec.cainjector.resources.requests.memory, "") != "" ? var.spec.cainjector.resources.requests.memory : null
        } : rk => rv if rv != null
      }
    } : k => v if v != null && length(v) > 0
  }

  # ---- typed chart values (twin of the Pulumi module's buildHelmValues) --
  typed_values = {
    for k, v in {
      # Planton default: install CRDs with the release; keep them (and all
      # certificate data cluster-wide) on uninstall unless explicitly told
      # otherwise.
      crds = {
        enabled = try(var.spec.crds.install, null) != null ? var.spec.crds.install : true
        keep    = try(var.spec.crds.keep_on_uninstall, null) != null ? var.spec.crds.keep_on_uninstall : true
      }

      replicaCount = try(var.spec.replicas, null)
      resources    = local.controller_resources

      global = (
        try(var.spec.log_level, null) != null ||
        (try(var.spec.leader_election_namespace, null) != null && try(var.spec.leader_election_namespace, "kube-system") != "kube-system")
        ) ? {
        for gk, gv in {
          logLevel = try(var.spec.log_level, null)
          leaderElection = (
            try(var.spec.leader_election_namespace, null) != null && try(var.spec.leader_election_namespace, "kube-system") != "kube-system"
          ) ? { namespace = var.spec.leader_election_namespace } : null
        } : gk => gv if gv != null
      } : null

      clusterResourceNamespace  = var.spec.cluster_resource_namespace != "" ? var.spec.cluster_resource_namespace : null
      enableCertificateOwnerRef = var.spec.enable_certificate_owner_ref ? true : null

      featureGates = length(var.spec.feature_gates) > 0 ? join(",", [
        for k in sort(keys(var.spec.feature_gates)) : "${k}=${var.spec.feature_gates[k]}"
      ]) : null

      dns01RecursiveNameservers = (
        try(var.spec.dns01_self_check, null) != null && length(try(var.spec.dns01_self_check.recursive_nameservers, [])) > 0
      ) ? join(",", var.spec.dns01_self_check.recursive_nameservers) : null
      dns01RecursiveNameserversOnly = (
        try(var.spec.dns01_self_check.recursive_nameservers_only, false) && length(try(var.spec.dns01_self_check.recursive_nameservers, [])) > 0
      ) ? true : null

      maxConcurrentChallenges = try(var.spec.max_concurrent_challenges, null)

      serviceAccount = length(local.workload_identity_annotations) > 0 ? {
        annotations = local.workload_identity_annotations
      } : null
      podLabels = try(var.spec.workload_identity.aks, null) != null ? {
        "azure.workload.identity/use" = "true"
      } : null

      imageRegistry = var.spec.image_registry != "" ? var.spec.image_registry : null

      prometheus = try(var.spec.prometheus, null) == null ? null : {
        for pk, pv in {
          enabled = try(var.spec.prometheus.enabled, null)
          servicemonitor = try(var.spec.prometheus.service_monitor, false) ? {
            for sk, sv in {
              enabled  = true
              interval = try(var.spec.prometheus.service_monitor_interval, null)
              labels   = length(try(var.spec.prometheus.service_monitor_labels, {})) > 0 ? var.spec.prometheus.service_monitor_labels : null
            } : sk => sv if sv != null
          } : null
        } : pk => pv if pv != null
      }

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

      podDisruptionBudget = var.spec.pod_disruption_budget ? { enabled = true } : null

      webhook = try(var.spec.webhook, null) == null ? null : {
        for wk, wv in {
          replicaCount   = try(var.spec.webhook.replicas, null)
          timeoutSeconds = try(var.spec.webhook.timeout_seconds, null)
          hostNetwork    = try(var.spec.webhook.host_network, false) ? true : null
          securePort     = try(var.spec.webhook.secure_port, null)
          resources      = local.webhook_resources
        } : wk => wv if wv != null
      }

      cainjector = try(var.spec.cainjector, null) == null ? null : {
        for ck, cv in {
          enabled      = try(var.spec.cainjector.enabled, null)
          replicaCount = try(var.spec.cainjector.replicas, null)
          resources    = local.cainjector_resources
        } : ck => cv if cv != null
      }

      startupapicheck = try(var.spec.startupapicheck, null) == null ? null : {
        for sk, sv in {
          enabled = try(var.spec.startupapicheck.enabled, null)
          timeout = try(var.spec.startupapicheck.timeout, null)
        } : sk => sv if sv != null
      }
    } : k => v if v != null && v != {}
  }
}
