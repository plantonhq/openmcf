# Computed values for the KubernetesClusterAutoscaler module. Every
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
# preserves every value's type. Cross-arm string picks use
# try(coalesce(...), null) instead of chained per-arm ternaries.
#
# Optional nested blocks are read with try(): HCL's && does NOT
# short-circuit, so chained null checks still dereference the null — and
# an absent optional attribute must never be dereferenced bare.

locals {
  # Chart identity — must stay byte-identical with the Pulumi module's vars:
  # cross-engine chart-name drift deploys two different products from one
  # manifest.
  helm_chart_name = "cluster-autoscaler"
  helm_chart_repo = "https://kubernetes.github.io/autoscaler"

  # Release name FIXED to the chart name: the autoscaler leader-elects and
  # owns the cluster-wide scaling decision — a second installation would
  # fight the first over every scale-up, so one installation per cluster is
  # the operating model.
  release_name = local.helm_chart_name

  # Chart version resolved to the pinned default when unset, so both engines
  # install the same chart whether or not the platform's defaulting
  # middleware ran — mirror of the Pulumi module's DefaultChartVersion
  # (chart 9.59.0 ships autoscaler 1.35.0 — chart and app versions move
  # SEPARATELY; the chart pin governs).
  chart_version = coalesce(try(var.spec.chart_version, null), "9.59.0")

  namespace = var.spec.namespace

  # Resource-identity labels stamped on the namespace this module creates
  # (never injected into the chart's own resources — Helm owns those).
  labels = merge(
    {
      "planton.ai/resource"      = "true"
      "planton.ai/resource-name" = var.metadata.name
      "planton.ai/resource-kind" = "KubernetesClusterAutoscaler"
    },
    var.metadata.id != null && var.metadata.id != "" ? { "planton.ai/resource-id" = var.metadata.id } : {},
    var.metadata.org != null && var.metadata.org != "" ? { "planton.ai/organization" = var.metadata.org } : {},
    var.metadata.env != null && var.metadata.env != "" ? { "planton.ai/environment" = var.metadata.env } : {}
  )

  # ---- provider arm (exactly one; proto oneof) -------------------------------
  aws         = try(var.spec.aws, null)
  azure       = try(var.spec.azure, null)
  gce         = try(var.spec.gce, null)
  cluster_api = try(var.spec.cluster_api, null)
  civo        = try(var.spec.civo, null)
  kwok        = try(var.spec.kwok, null)

  # The chart's cloudProvider value — selects the per-provider command/env
  # blocks in the chart's deployment.yaml AND which keys the chart's
  # credential Secret carries (templates/secret.yaml). The key is the
  # autoscaler binary's provider name ("clusterapi", not "cluster_api").
  # one() over the pruned key list instead of a chained per-arm ternary.
  cloud_provider = one([
    for k, v in {
      aws        = local.aws
      azure      = local.azure
      gce        = local.gce
      clusterapi = local.cluster_api
      civo       = local.civo
      kwok       = local.kwok
    } : k if v != null
  ])

  # ServiceAccountName DERIVED from the chart's naming, verified in
  # templates/_helpers.tpl: with rbac.serviceAccount.name unset the service
  # account takes the fullname template, whose default name is
  # "<cloudProvider>-<chartName>" (NOT the bare chart name). That never
  # equals the release name, so fullname renders
  # "<release>-<cloudProvider>-<chartName>" — e.g. for the aws arm:
  # "cluster-autoscaler-aws-cluster-autoscaler" (well under the 63-char
  # truncation for every supported arm).
  service_account_name = "${local.release_name}-${local.cloud_provider}-${local.helm_chart_name}"

  # ---- autoDiscovery ---------------------------------------------------------
  # CHART GATE (verified in templates/deployment.yaml line 1): the
  # Deployment only renders when autoDiscovery.clusterName/namespace/labels
  # or autoscalingGroups is set. autoscalingGroupsnamePrefix (the GCE
  # contract) does NOT satisfy the gate — the chart README explicitly
  # requires autoDiscovery.clusterName ("any-name"; unused by the gce
  # provider blocks) for GCE, and kwok has no typed gate key at all. Both
  # arms therefore render autoDiscovery.clusterName = metadata.name — a
  # benign, deterministic gate value. Without it the release would
  # "succeed" while installing NO autoscaler pod.
  #
  # At most one source is non-null (proto oneof); coalesce also skips the
  # azure empty-string default, and try() turns the all-null case into null.
  auto_discovery_cluster_name = try(coalesce(
    try(local.aws.auto_discovery.cluster_name, null),
    try(local.azure.cluster_name, null),
    (local.gce != null || local.kwok != null) ? var.metadata.name : null
  ), null)

  auto_discovery_values = {
    for k, v in {
      clusterName = local.auto_discovery_cluster_name
      # Cluster API: renders --node-group-auto-discovery=clusterapi:
      # namespace=<ns> via the chart's capiAutodiscoveryConfig helper — and
      # satisfies the Deployment render gate. With namespace empty the gate
      # needs autoDiscovery.clusterName/labels via helm_values (upstream
      # requires at least one discovery dimension for clusterapi).
      namespace = try(local.cluster_api.namespace, "") != "" ? local.cluster_api.namespace : null
      # AWS tag-based discovery: the chart's default tags are the standard
      # k8s.io/cluster-autoscaler/enabled + .../<clusterName> pair — tags
      # render only on explicit override.
      tags = length(try(local.aws.auto_discovery.tags, [])) > 0 ? local.aws.auto_discovery.tags : null
    } : k => v if v != null
  }

  # ---- static node groups (aws XOR azure — proto-validated) -------------------
  # Each entry renders one --nodes=<min>:<max>:<name> flag. concat of the
  # two arms' (identically typed) lists avoids a per-arm ternary; at most
  # one list is non-empty.
  static_node_groups = concat(try(local.aws.node_groups, []), try(local.azure.node_groups, []))
  autoscaling_groups_values = length(local.static_node_groups) > 0 ? [
    for g in local.static_node_groups : {
      name    = g.name
      minSize = g.min_size
      maxSize = g.max_size
    }
  ] : null

  # GCE managed instance groups by name prefix: each entry renders one
  # --node-group-auto-discovery=mig:namePrefix=<name>,min=..,max=.. flag.
  # NOTE the chart's key really is "autoscalingGroupsnamePrefix"
  # (lowercase n — values.yaml).
  gce_prefix_values = local.gce == null ? null : [
    for g in local.gce.instance_group_prefixes : {
      name    = g.name
      minSize = g.min_size
      maxSize = g.max_size
    }
  ]

  # ---- rbac / service-account annotations -------------------------------------
  # The chart forwards service-account annotations verbatim
  # (templates/serviceaccount.yaml) — the EKS/GKE webhooks pick the keyless
  # identity up from these well-known keys. At most one entry is non-null
  # (arms are mutually exclusive).
  service_account_annotations = {
    for k, v in {
      "eks.amazonaws.com/role-arn"     = try(local.aws.irsa_role_arn, "") != "" ? local.aws.irsa_role_arn : null
      "iam.gke.io/gcp-service-account" = try(local.gce.workload_identity_service_account_email, "") != "" ? local.gce.workload_identity_service_account_email : null
    } : k => v if v != null
  }

  rbac_values = {
    for k, v in {
      serviceAccount = length(local.service_account_annotations) > 0 ? {
        annotations = local.service_account_annotations
      } : null
      # rbac.clusterScoped=false switches the chart from ClusterRole to a
      # namespaced Role — the least-privilege posture the chart documents
      # "most useful for Cluster-API".
      clusterScoped = try(local.cluster_api.namespace_scoped_rbac, false) ? false : null
    } : k => v if v != null
  }

  # ---- scaling flags → extraArgs ------------------------------------------------
  # The chart renders every extraArgs entry as --<key>=<value>
  # (deployment.yaml) — flag names carry no leading dashes. Every value is
  # a STRING deliberately (they are CLI flag text; tostring keeps the
  # booleans in lockstep with the Pulumi module's strconv rendering).
  #
  # Presence-aware optional bools: upstream defaults are true for the
  # skip-nodes flags and scale-down-enabled, so an explicit false MUST
  # render — only absence stays silent.
  scaling    = try(var.spec.scaling, null)
  scale_down = try(var.spec.scaling.scale_down, null)

  scaling_extra_args = {
    for k, v in {
      expander                           = try(local.scaling.expander, "") != "" ? local.scaling.expander : null
      "balance-similar-node-groups"      = try(local.scaling.balance_similar_node_groups, false) ? "true" : null
      "scan-interval"                    = try(local.scaling.scan_interval, "") != "" ? local.scaling.scan_interval : null
      "max-node-provision-time"          = try(local.scaling.max_node_provision_time, "") != "" ? local.scaling.max_node_provision_time : null
      "skip-nodes-with-local-storage"    = try(local.scaling.skip_nodes_with_local_storage, null) != null ? tostring(local.scaling.skip_nodes_with_local_storage) : null
      "skip-nodes-with-system-pods"      = try(local.scaling.skip_nodes_with_system_pods, null) != null ? tostring(local.scaling.skip_nodes_with_system_pods) : null
      "scale-down-enabled"               = try(local.scale_down.enabled, null) != null ? tostring(local.scale_down.enabled) : null
      "scale-down-utilization-threshold" = try(local.scale_down.utilization_threshold, "") != "" ? local.scale_down.utilization_threshold : null
      "scale-down-unneeded-time"         = try(local.scale_down.unneeded_time, "") != "" ? local.scale_down.unneeded_time : null
      "scale-down-delay-after-add"       = try(local.scale_down.delay_after_add, "") != "" ? local.scale_down.delay_after_add : null
      "scale-down-delay-after-delete"    = try(local.scale_down.delay_after_delete, "") != "" ? local.scale_down.delay_after_delete : null
      "scale-down-delay-after-failure"   = try(local.scale_down.delay_after_failure, "") != "" ? local.scale_down.delay_after_failure : null
    } : k => v if v != null
  }

  # PRECEDENCE (comment mirrored in the Pulumi module): the typed scaling
  # block renders first, then spec.extra_args merges OVER it — user entries
  # win on key collision. The chart's own extraArgs defaults
  # (logtostderr/stderrthreshold/v) stay chart-side: Helm coalesces our
  # extraArgs map over the chart default per key, so unspecified defaults
  # survive on both engines identically.
  extra_args_values = merge(local.scaling_extra_args, try(var.spec.extra_args, {}))

  # ---- deployment sizing and scheduling --------------------------------------
  deployment_resources = try(var.spec.deployment.resources, null) == null ? null : {
    for k, v in {
      limits = try(var.spec.deployment.resources.limits, null) == null ? null : {
        for lk, lv in {
          cpu    = try(var.spec.deployment.resources.limits.cpu, "") != "" ? var.spec.deployment.resources.limits.cpu : null
          memory = try(var.spec.deployment.resources.limits.memory, "") != "" ? var.spec.deployment.resources.limits.memory : null
        } : lk => lv if lv != null
      }
      requests = try(var.spec.deployment.resources.requests, null) == null ? null : {
        for rk, rv in {
          cpu    = try(var.spec.deployment.resources.requests.cpu, "") != "" ? var.spec.deployment.resources.requests.cpu : null
          memory = try(var.spec.deployment.resources.requests.memory, "") != "" ? var.spec.deployment.resources.requests.memory : null
        } : rk => rv if rv != null
      }
    } : k => v if v != null && length(v) > 0
  }

  tolerations_values = length(try(var.spec.deployment.tolerations, [])) > 0 ? [
    for t in var.spec.deployment.tolerations : {
      for tk, tv in {
        key               = try(t.key, "") != "" ? t.key : null
        operator          = try(t.operator, "") != "" ? t.operator : null
        value             = try(t.value, "") != "" ? t.value : null
        effect            = try(t.effect, "") != "" ? t.effect : null
        tolerationSeconds = try(t.toleration_seconds, null)
      } : tk => tv if tv != null
    }
  ] : null

  # ---- own telemetry -----------------------------------------------------------
  # serviceMonitor.selector is a plain label map rendered onto the
  # ServiceMonitor's metadata.labels; Helm's per-key coalesce replaces the
  # chart's default {release: prometheus-operator} entry.
  service_monitor_values = try(var.spec.prometheus.service_monitor, false) ? {
    for k, v in {
      enabled  = true
      selector = try(var.spec.prometheus.service_monitor_selector_release, "") != "" ? { release = var.spec.prometheus.service_monitor_selector_release } : null
    } : k => v if v != null
  } : null

  # ---- typed chart values (twin of the Pulumi module's buildHelmValues) --------
  # One flat null-pruned literal, matching the chart's mostly-flat key
  # layout. Per-arm keys reference only their own arm, so an unselected
  # arm's entries resolve to null and prune away. Credentials
  # (awsSecretAccessKey, azureClientSecret, civoApiKey) flow only into
  # these chart values, which the chart materializes into its own Secret
  # (templates/secret.yaml) and wires to the pod via secretKeyRef — never
  # log or output them.
  typed_values = {
    for k, v in {
      cloudProvider = local.cloud_provider

      # -- aws: awsRegion also gates the chart's AWS env block
      # (deployment.yaml renders AWS_REGION and the key envs only when
      # awsRegion != ""); the chart's Secret renders only when BOTH keys
      # are set.
      awsRegion          = try(local.aws.region, null)
      awsAccessKeyID     = try(local.aws.access_keys.access_key_id, null)
      awsSecretAccessKey = try(local.aws.access_keys.secret_access_key, null)

      # -- azure: both land in the chart's credential Secret and reach the
      # pod as ARM_SUBSCRIPTION_ID / ARM_RESOURCE_GROUP via secretKeyRef.
      # Exactly one identity posture per proto validation. VERIFIED: no
      # template adds the azure.workload.identity/use pod label — clusters
      # relying on the azure-workload-identity webhook add podLabels via
      # helm_values.
      azureSubscriptionID               = try(local.azure.subscription_id, null)
      azureResourceGroup                = try(local.azure.resource_group, null)
      azureUseWorkloadIdentityExtension = try(local.azure.identity.use_workload_identity, false) ? true : null
      azureUseManagedIdentityExtension  = try(local.azure.identity.use_managed_identity, false) ? true : null
      azureUserAssignedIdentityID       = (try(local.azure.identity.use_managed_identity, false) && try(local.azure.identity.user_assigned_identity_id, "") != "") ? local.azure.identity.user_assigned_identity_id : null
      azureTenantID                     = try(local.azure.identity.service_principal.tenant_id, null)
      azureClientID                     = try(local.azure.identity.service_principal.client_id, null)
      azureClientSecret                 = try(local.azure.identity.service_principal.client_secret, null)

      # -- gce
      autoscalingGroupsnamePrefix = local.gce_prefix_values

      # -- cluster_api: mode/kubeconfig secret; both match the chart's own
      # defaults when unset — rendered only when set.
      clusterAPIMode             = (try(local.cluster_api.mode, null) != null && try(local.cluster_api.mode, "") != "") ? local.cluster_api.mode : null
      clusterAPIKubeconfigSecret = try(local.cluster_api.kubeconfig_secret, "") != "" ? local.cluster_api.kubeconfig_secret : null

      # -- civo: all four land in the chart's credential Secret
      # (api-url/api-key/cluster-id/region) and reach the pod as CIVO_*
      # env vars via secretKeyRef. NO gate key is injected for civo: the
      # chart requires autoscalingGroups (the Civo node pools), which the
      # typed spec does not model — they ride helm_values, and that same
      # document satisfies the Deployment render gate.
      civoClusterID = try(local.civo.cluster_id, null)
      civoRegion    = try(local.civo.region, null)
      civoApiKey    = try(local.civo.api_key, null)
      civoApiUrl    = (try(local.civo.api_url, null) != null && try(local.civo.api_url, "") != "") ? local.civo.api_url : null

      # -- kwok: reaches the pod as KWOK_PROVIDER_CONFIGMAP and names the
      # ConfigMap the chart itself creates for kwok.
      kwokConfigMapName = (try(local.kwok.config_map_name, null) != null && try(local.kwok.config_map_name, "") != "") ? local.kwok.config_map_name : null

      # -- cross-arm assemblies
      autoDiscovery     = length(local.auto_discovery_values) > 0 ? local.auto_discovery_values : null
      autoscalingGroups = local.autoscaling_groups_values
      rbac              = length(local.rbac_values) > 0 ? local.rbac_values : null
      extraArgs         = length(local.extra_args_values) > 0 ? local.extra_args_values : null

      # -- deployment sizing and scheduling. priorityClassName matches the
      # chart's own default ("system-cluster-critical") when unset —
      # rendered only when set. Replicas leader-elect; extras are warm
      # standbys.
      replicaCount      = try(var.spec.deployment.replicas, null)
      resources         = local.deployment_resources
      priorityClassName = (try(var.spec.deployment.priority_class_name, null) != null && try(var.spec.deployment.priority_class_name, "") != "") ? var.spec.deployment.priority_class_name : null
      nodeSelector      = length(try(var.spec.deployment.node_selector, {})) > 0 ? var.spec.deployment.node_selector : null
      tolerations       = local.tolerations_values

      serviceMonitor = local.service_monitor_values
    } : k => v if v != null
  }
}
