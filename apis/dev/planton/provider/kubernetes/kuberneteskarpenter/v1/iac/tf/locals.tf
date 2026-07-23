# Computed values for the KubernetesKarpenter module. Every resolution here
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
# short-circuit, so chained null checks still dereference the null — and
# var.spec is typed 'any', so an absent attribute is an error, not a null.

locals {
  # Chart identity — must stay byte-identical with the Pulumi module's vars:
  # cross-engine chart drift deploys two different products from one
  # manifest. OCI ENGINE ASYMMETRY: this provider takes the registry path
  # in `repository` and the bare chart name in `chart`, joining them
  # internally; Pulumi's helm.v3.Release needs the JOINED
  # "oci://public.ecr.aws/karpenter/<chart>" reference with no repository
  # opts. Same chart bytes, different wiring.
  helm_oci_repo   = "oci://public.ecr.aws/karpenter"
  helm_chart_name = "karpenter"
  crd_chart_name  = "karpenter-crd"

  # Release names FIXED: Karpenter owns the cluster-wide karpenter.sh label
  # domain, its CRDs, and node lifecycle — one installation per cluster is
  # an upstream constraint, so the names never derive from metadata.name.
  release_name     = "karpenter"
  crd_release_name = "karpenter-crd"

  # Chart version resolved to the pinned default when unset, so both engines
  # install the same charts whether or not the platform's defaulting
  # middleware ran — mirror of the Pulumi module's DefaultChartVersion. The
  # karpenter and karpenter-crd charts version together with the controller
  # (both 1.14.0 = Karpenter 1.14.0), so ONE version pins BOTH releases.
  chart_version = coalesce(try(var.spec.chart_version, null), "1.14.0")

  namespace = var.spec.namespace

  # Resource-identity labels stamped on the namespace this module creates
  # (never injected into the charts' own resources — Helm owns those).
  labels = merge(
    {
      "planton.ai/resource"      = "true"
      "planton.ai/resource-name" = var.metadata.name
      "planton.ai/resource-kind" = "KubernetesKarpenter"
    },
    var.metadata.id != null && var.metadata.id != "" ? { "planton.ai/resource-id" = var.metadata.id } : {},
    var.metadata.org != null && var.metadata.org != "" ? { "planton.ai/organization" = var.metadata.org } : {},
    var.metadata.env != null && var.metadata.env != "" ? { "planton.ai/environment" = var.metadata.env } : {}
  )

  # ---- CRD lifecycle -------------------------------------------------------
  # The CRDs install as a dedicated karpenter-crd release — upstream's
  # supported mechanism for keeping them upgradable (Helm installs the
  # copies bundled inside the main chart once and NEVER upgrades them).
  # keep_on_uninstall stamps the standard Helm resource-policy annotation
  # onto every CRD through the CRD chart's additionalAnnotations — without
  # it a plain uninstall cascade-deletes every NodePool/EC2NodeClass/
  # NodeClaim in the cluster.
  crds_install = try(var.spec.crds.install, null) != null ? var.spec.crds.install : true
  crds_keep    = try(var.spec.crds.keep_on_uninstall, null) != null ? var.spec.crds.keep_on_uninstall : true

  # Values document for the CRD release (its whole values surface is this
  # one knob) — twin of the Pulumi module's buildCrdHelmValues.
  crd_values = local.crds_keep ? {
    additionalAnnotations = { "helm.sh/resource-policy" = "keep" }
  } : {}

  # ---- AWS provider arm ------------------------------------------------------
  aws_enabled = try(var.spec.aws, null) != null

  # IRSA rides the service-account annotation the EKS webhook watches — not
  # a settings entry. Empty means EKS Pod Identity (association configured
  # cloud-side, no annotation needed).
  service_account_values = try(var.spec.aws.irsa_role_arn, "") != "" ? {
    annotations = { "eks.amazonaws.com/role-arn" = var.spec.aws.irsa_role_arn }
  } : null

  # ---- chart-default resolution ----------------------------------------------
  # Fields whose spec default mirrors the chart default render with the
  # default APPLIED — explicit and byte-identical across engines regardless
  # of whether the platform's defaulting middleware ran. Every fallback
  # literal mirrors the served chart's values.yaml — twin of the Pulumi
  # module's vars defaults table. (coalesce also skips empty strings, which
  # matches the Go side treating "" as unset.)
  replicas            = coalesce(try(var.spec.controller.replicas, null), 2)
  log_level           = coalesce(try(var.spec.controller.log_level, null), "info")
  batch_max_duration  = coalesce(try(var.spec.batching.max_duration, null), "10s")
  batch_idle_duration = coalesce(try(var.spec.batching.idle_duration, null), "1s")
  preference_policy   = coalesce(try(var.spec.scheduling.preference_policy, null), "Respect")
  min_values_policy   = coalesce(try(var.spec.scheduling.min_values_policy, null), "Strict")
  priority_class_name = coalesce(try(var.spec.controller_scheduling.priority_class_name, null), "system-cluster-critical")

  # ---- feature gates -----------------------------------------------------------
  # The WHOLE map renders with defaults applied — the deployment template
  # composes the FEATURE_GATES env var from all six keys unconditionally,
  # so explicit is safer than sparse (and reservedCapacity's default is
  # TRUE, unlike the other five).
  feature_gates_values = {
    nodeRepair              = try(var.spec.feature_gates.node_repair, false)
    nodeOverlay             = try(var.spec.feature_gates.node_overlay, false)
    reservedCapacity        = try(var.spec.feature_gates.reserved_capacity, null) != null ? var.spec.feature_gates.reserved_capacity : true
    spotToSpotConsolidation = try(var.spec.feature_gates.spot_to_spot_consolidation, false)
    staticCapacity          = try(var.spec.feature_gates.static_capacity, false)
    capacityBuffer          = try(var.spec.feature_gates.capacity_buffer, false)
  }

  # ---- settings block ------------------------------------------------------------
  # clusterName is the one value the chart REFUSES to render without
  # (deployment.yaml wraps it in `required`) — always rendered, never
  # conditional.
  settings_values = {
    for k, v in {
      clusterName     = var.spec.cluster.name
      clusterEndpoint = try(var.spec.cluster.endpoint, "") != "" ? var.spec.cluster.endpoint : null
      eksControlPlane = try(var.spec.cluster.eks_control_plane, false) ? true : null
      clusterCABundle = try(var.spec.cluster.ca_bundle, "") != "" ? var.spec.cluster.ca_bundle : null

      interruptionQueue = try(var.spec.aws.interruption_queue, "") != "" ? var.spec.aws.interruption_queue : null
      isolatedVPC       = try(var.spec.aws.isolated_vpc, false) ? true : null
      enableZonalShift  = try(var.spec.aws.enable_zonal_shift, false) ? true : null

      # TYPE FIDELITY: the chart's reservedENIs default is the STRING "0" —
      # tostring() keeps the served chart's type (the Pulumi module renders
      # strconv.Itoa for the same reason). Rendered whenever the AWS arm is
      # present, default applied.
      reservedENIs = local.aws_enabled ? tostring(coalesce(try(var.spec.aws.reserved_enis, null), 0)) : null

      # TYPE FIDELITY (inverse case): the chart's vmMemoryOverheadPercent
      # default is the NUMBER 0.075 — the spec carries it as a string, so
      # tonumber() restores the chart's type (the Pulumi module parses a
      # float for the same reason).
      vmMemoryOverheadPercent = local.aws_enabled ? tonumber(coalesce(try(var.spec.aws.vm_memory_overhead_percent, null), "0.075")) : null

      batchMaxDuration  = local.batch_max_duration
      batchIdleDuration = local.batch_idle_duration

      preferencePolicy = local.preference_policy
      minValuesPolicy  = local.min_values_policy

      featureGates = local.feature_gates_values
    } : k => v if v != null
  }

  # ---- controller container resources ---------------------------------------------
  # Resources live under the controller CONTAINER block
  # (controller.resources), unlike replicas/logLevel which are top-level —
  # the chart's layout, not ours. KEDA-pattern null-prune of the shared
  # ContainerResources shape.
  controller_resources = try(var.spec.controller.resources, null) == null ? null : {
    for k, v in {
      limits = try(var.spec.controller.resources.limits, null) == null ? null : {
        for lk, lv in {
          cpu    = try(var.spec.controller.resources.limits.cpu, "") != "" ? var.spec.controller.resources.limits.cpu : null
          memory = try(var.spec.controller.resources.limits.memory, "") != "" ? var.spec.controller.resources.limits.memory : null
        } : lk => lv if lv != null
      }
      requests = try(var.spec.controller.resources.requests, null) == null ? null : {
        for rk, rv in {
          cpu    = try(var.spec.controller.resources.requests.cpu, "") != "" ? var.spec.controller.resources.requests.cpu : null
          memory = try(var.spec.controller.resources.requests.memory, "") != "" ? var.spec.controller.resources.requests.memory : null
        } : rk => rv if rv != null
      }
    } : k => v if v != null && length(v) > 0
  }

  # ---- typed chart values (twin of the Pulumi module's buildHelmValues) ------
  # No fullnameOverride: the release name ("karpenter") contains the chart
  # name, so the chart's fullname template resolves to the release name —
  # there is nothing for an override to pin.
  typed_values = {
    for k, v in {
      settings       = local.settings_values
      serviceAccount = local.service_account_values

      replicas = local.replicas
      logLevel = local.log_level

      controller = local.controller_resources != null && length(coalesce(local.controller_resources, {})) > 0 ? {
        resources = local.controller_resources
      } : null

      # Where Karpenter itself runs — NOT what it provisions. The chart
      # already pins controller pods away from Karpenter-provisioned nodes
      # (nodeAffinity on karpenter.sh/nodepool DoesNotExist).
      priorityClassName = local.priority_class_name

      # nodeSelector MERGES onto the chart's kubernetes.io/os=linux default
      # (Helm deep-merges maps) — entries here narrow, never replace.
      nodeSelector = length(try(var.spec.controller_scheduling.node_selector, {})) > 0 ? var.spec.controller_scheduling.node_selector : null

      # tolerations REPLACE the chart's default list (CriticalAddonsOnly) —
      # Helm replaces lists wholesale — so render only when the spec
      # provides them; an empty list would silently DROP the default.
      tolerations = length(try(var.spec.controller_scheduling.tolerations, [])) > 0 ? [
        for t in var.spec.controller_scheduling.tolerations : {
          for tk, tv in {
            key               = try(t.key, "") != "" ? t.key : null
            operator          = try(t.operator, "") != "" ? t.operator : null
            value             = try(t.value, "") != "" ? t.value : null
            effect            = try(t.effect, "") != "" ? t.effect : null
            tolerationSeconds = try(t.toleration_seconds, null)
          } : tk => tv if tv != null
        }
      ] : null

      hostNetwork = try(var.spec.controller_scheduling.host_network, false) ? true : null

      serviceMonitor = try(var.spec.prometheus.service_monitor, false) ? { enabled = true } : null
    } : k => v if v != null
  }
}
