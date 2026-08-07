# Computed values for the KubernetesRayCluster module. Every resolution
# here has an exact twin in the Pulumi module (locals.go /
# raycluster_cr.go) — keep them in lockstep: same rendered CR body, same
# outputs.
#
# HCL DISCIPLINE: conditional keys are contributed via merge() of
# `cond ? { key = value } : {}` singleton maps, or one object literal
# pruned with `{ for k, v in {...} : k => v if v != null }` — a ternary
# whose branches are differently-shaped objects fails plan-time type
# unification. Optional nested blocks are read with try(): HCL's && does
# NOT short-circuit. Optional scalars in string templates resolve with
# try(coalesce(...)).

locals {
  # metadata.name is the CR name — the operator's naming root: the head
  # Service is `<name>-head-svc` (GenerateHeadServiceName, kuberay
  # ray-operator controllers/ray/utils/util.go), head/worker pod names
  # prefix it, and in token auth mode the operator's generated
  # bearer-token Secret is named EXACTLY the cluster name
  # (reconcileAuthSecret → utils.CheckName(instance.Name), kuberay
  # raycluster_controller.go).
  cluster_name = var.metadata.name
  namespace    = var.spec.namespace
  api_version  = "ray.io/v1"

  # Planton governance labels on the module-created objects.
  labels = merge(
    {
      "planton.ai/resource"      = "true"
      "planton.ai/resource-name" = var.metadata.name
      "planton.ai/resource-kind" = "KubernetesRayCluster"
    },
    try(var.metadata.id, "") != "" ? { "planton.ai/resource-id" = var.metadata.id } : {},
    try(var.metadata.org, "") != "" ? { "planton.ai/organization" = var.metadata.org } : {},
    try(var.metadata.env, "") != "" ? { "planton.ai/environment" = var.metadata.env } : {}
  )

  # ---- image -----------------------------------------------------------------------
  # VERSION/IMAGE LOCKSTEP: ray_version must match the Ray inside the
  # image — the operator reads rayVersion to shape its commands but runs
  # the image as given; a mismatch fails at runtime, not at apply.
  # Deriving the default image from ray_version keeps the lockstep true
  # by construction; a custom image overrides it deliberately.
  image = try(var.spec.image, "") != "" ? var.spec.image : "rayproject/ray:${var.spec.ray_version}"

  # ---- spec.auth → authOptions -------------------------------------------------------
  # This catalog's default is TOKEN auth (secure-by-default) while the
  # operator's own nil-authOptions default is DISABLED — so authOptions
  # renders ALWAYS, never left to the CR default: an absent block would
  # silently deploy the legacy open cluster (anyone reaching the
  # dashboard port runs arbitrary code). Empty auth or empty mode means
  # token; only an explicit "disabled" opts out.
  token_auth_enabled = try(coalesce(var.spec.auth.mode, "token"), "token") != "disabled"

  existing_token_secret_name = try(var.spec.auth.existing_token_secret_name, "")

  auth_options = merge(
    { mode = local.token_auth_enabled ? "token" : "disabled" },
    # secretName tells the operator to skip generating a token Secret
    # and read the bring-your-own one (data key `auth_token`) instead.
    local.token_auth_enabled && local.existing_token_secret_name != "" ? { secretName = local.existing_token_secret_name } : {}
  )

  # ---- spec.autoscaling → enableInTreeAutoscaling / autoscalerOptions ----------------
  autoscaling_enabled = try(var.spec.autoscaling.enabled, false)

  autoscaler_resources_requests = {
    for k, v in {
      cpu    = try(var.spec.autoscaling.resources.requests.cpu, "") != "" ? var.spec.autoscaling.resources.requests.cpu : null
      memory = try(var.spec.autoscaling.resources.requests.memory, "") != "" ? var.spec.autoscaling.resources.requests.memory : null
    } : k => v if v != null
  }

  autoscaler_resources_limits = {
    for k, v in {
      cpu    = try(var.spec.autoscaling.resources.limits.cpu, "") != "" ? var.spec.autoscaling.resources.limits.cpu : null
      memory = try(var.spec.autoscaling.resources.limits.memory, "") != "" ? var.spec.autoscaling.resources.limits.memory : null
    } : k => v if v != null
  }

  autoscaler_resources = merge(
    length(local.autoscaler_resources_limits) > 0 ? { limits = local.autoscaler_resources_limits } : {},
    length(local.autoscaler_resources_requests) > 0 ? { requests = local.autoscaler_resources_requests } : {}
  )

  # Keys render only when declared — the operator's defaults (60s idle
  # timeout, Default upscaling) stay authoritative otherwise.
  autoscaler_options = merge(
    {
      for k, v in {
        idleTimeoutSeconds = try(var.spec.autoscaling.idle_timeout_seconds, null)
        upscalingMode      = try(var.spec.autoscaling.upscaling_mode, "") != "" ? var.spec.autoscaling.upscaling_mode : null
      } : k => v if v != null
    },
    length(local.autoscaler_resources) > 0 ? { resources = local.autoscaler_resources } : {}
  )

  # ---- spec.gcs_fault_tolerance → gcsFaultToleranceOptions ---------------------------
  # STATE TRUTH: without this block the head pod's GCS holds the
  # cluster's control state in memory — losing the head loses every
  # job, actor, and worker registration. With it, state lives in the
  # external Redis-protocol store and a replaced head RECOVERS.
  gcs_ft_enabled = try(var.spec.gcs_fault_tolerance.enabled, false)

  # redisAddress arrives pre-resolved (literal or KubernetesValkey
  # reference) and is ALWAYS rendered inside the block — the spec's CEL
  # guarantees it is set whenever enabled is true.
  gcs_ft_options = merge(
    { redisAddress = try(var.spec.gcs_fault_tolerance.redis_address, "") },
    try(var.spec.gcs_fault_tolerance.redis_password_secret, null) != null ? {
      redisPassword = {
        valueFrom = {
          secretKeyRef = {
            name = var.spec.gcs_fault_tolerance.redis_password_secret.name
            key  = var.spec.gcs_fault_tolerance.redis_password_secret.key
          }
        }
      }
    } : {},
    # Empty means the operator derives one from the cluster's UID (safe
    # default); set explicitly only when state must survive
    # delete-and-recreate of the CR itself.
    try(var.spec.gcs_fault_tolerance.external_storage_namespace, "") != "" ? {
      externalStorageNamespace = var.spec.gcs_fault_tolerance.external_storage_namespace
    } : {}
  )

  # ---- spec.head → headGroupSpec ------------------------------------------------------
  # schedule_tasks_on_head empty = false: production heads stay
  # unloaded (a task-loaded head starves the GCS).
  schedule_tasks_on_head = try(coalesce(var.spec.head.schedule_tasks_on_head, false), false)

  # The module OWNS two rayStartParams and merges them LAST so user
  # entries cannot override them:
  #   dashboard-host "0.0.0.0" — the dashboard binds localhost
  #     otherwise and the head Service cannot answer (sample-verified:
  #     upstream ray-cluster.complete.yaml sets exactly this);
  #   num-cpus "0" (unless schedule_tasks_on_head) — keeps the Ray
  #     scheduler from placing application work on the head.
  head_ray_start_params = merge(
    try(var.spec.head.ray_start_params, {}),
    { "dashboard-host" = "0.0.0.0" },
    local.schedule_tasks_on_head ? {} : { "num-cpus" = "0" }
  )

  head_resources_requests = {
    for k, v in {
      cpu    = try(var.spec.head.resources.requests.cpu, "") != "" ? var.spec.head.resources.requests.cpu : null
      memory = try(var.spec.head.resources.requests.memory, "") != "" ? var.spec.head.resources.requests.memory : null
    } : k => v if v != null
  }

  head_resources_limits = {
    for k, v in {
      cpu    = try(var.spec.head.resources.limits.cpu, "") != "" ? var.spec.head.resources.limits.cpu : null
      memory = try(var.spec.head.resources.limits.memory, "") != "" ? var.spec.head.resources.limits.memory : null
    } : k => v if v != null
  }

  head_resources = merge(
    length(local.head_resources_limits) > 0 ? { limits = local.head_resources_limits } : {},
    length(local.head_resources_requests) > 0 ? { requests = local.head_resources_requests } : {}
  )

  head_tolerations = [
    for t in try(var.spec.head.scheduling.tolerations, []) : { for k, v in {
      key               = t.key != "" ? t.key : null
      operator          = t.operator != "" ? t.operator : null
      value             = t.value != "" ? t.value : null
      effect            = t.effect != "" ? t.effect : null
      tolerationSeconds = t.toleration_seconds
    } : k => v if v != null }
  ]

  # Pod-template scheduling: unlike CRs that model only affinity, the
  # RayCluster embeds a full corev1 pod template — nodeSelector renders
  # verbatim.
  head_scheduling = merge(
    length(try(var.spec.head.scheduling.node_selector, {})) > 0 ? { nodeSelector = var.spec.head.scheduling.node_selector } : {},
    length(local.head_tolerations) > 0 ? { tolerations = local.head_tolerations } : {},
    try(var.spec.head.scheduling.priority_class_name, "") != "" ? { priorityClassName = var.spec.head.scheduling.priority_class_name } : {}
  )

  head_group_spec = {
    rayStartParams = local.head_ray_start_params
    template = {
      spec = merge(
        {
          containers = [
            merge(
              {
                name  = "ray-head"
                image = local.image
              },
              length(local.head_resources) > 0 ? { resources = local.head_resources } : {}
            )
          ]
        },
        local.head_scheduling
      )
    }
  }

  # ---- spec.worker_groups → workerGroupSpecs ------------------------------------------
  # Parallel per-group lists (HCL for-expressions have no let-bindings);
  # worker_group_specs zips them by index.
  #
  # extra_resource_limits land in LIMITS ONLY: extended resources
  # (nvidia.com/gpu and friends) must not be requested without limits —
  # Kubernetes rejects requests-without-limits for extended resources,
  # and Ray discovers accelerators from the container LIMITS.
  worker_container_limits = [
    for g in try(var.spec.worker_groups, []) : merge(
      {
        for k, v in {
          cpu    = try(g.resources.limits.cpu, "") != "" ? g.resources.limits.cpu : null
          memory = try(g.resources.limits.memory, "") != "" ? g.resources.limits.memory : null
        } : k => v if v != null
      },
      try(g.extra_resource_limits, {})
    )
  ]

  worker_container_requests = [
    for g in try(var.spec.worker_groups, []) : {
      for k, v in {
        cpu    = try(g.resources.requests.cpu, "") != "" ? g.resources.requests.cpu : null
        memory = try(g.resources.requests.memory, "") != "" ? g.resources.requests.memory : null
      } : k => v if v != null
    }
  ]

  worker_scheduling = [
    for g in try(var.spec.worker_groups, []) : merge(
      length(try(g.scheduling.node_selector, {})) > 0 ? { nodeSelector = g.scheduling.node_selector } : {},
      length(try(g.scheduling.tolerations, [])) > 0 ? {
        tolerations = [
          for t in g.scheduling.tolerations : { for k, v in {
            key               = t.key != "" ? t.key : null
            operator          = t.operator != "" ? t.operator : null
            value             = t.value != "" ? t.value : null
            effect            = t.effect != "" ? t.effect : null
            tolerationSeconds = t.toleration_seconds
          } : k => v if v != null }
        ]
      } : {},
      try(g.scheduling.priority_class_name, "") != "" ? { priorityClassName = g.scheduling.priority_class_name } : {}
    )
  ]

  # replicas/minReplicas/maxReplicas render only when declared — the v1
  # CRD defaults all three (replicas 0, minReplicas 0, maxReplicas
  # maxInt32; verified in the pinned ray.io_rayclusters.yaml schema).
  #
  # rayStartParams renders ALWAYS, {} when empty: NOT required by the v1
  # CRD (only the retired v1alpha1 schema listed it under `required`;
  # verified in the pinned CRD), but the operator's own Go type
  # serializes it unconditionally and every upstream sample renders
  # `rayStartParams: {}` explicitly — matching that keeps SSA diffs
  # quiet.
  worker_group_specs = [
    for i, g in try(var.spec.worker_groups, []) : merge(
      { groupName = g.name },
      try(g.replicas, null) != null ? { replicas = g.replicas } : {},
      try(g.min_replicas, null) != null ? { minReplicas = g.min_replicas } : {},
      try(g.max_replicas, null) != null ? { maxReplicas = g.max_replicas } : {},
      {
        rayStartParams = try(g.ray_start_params, {})
        template = {
          spec = merge(
            {
              containers = [
                merge(
                  {
                    name  = "ray-worker"
                    image = local.image
                  },
                  length(local.worker_container_limits[i]) > 0 || length(local.worker_container_requests[i]) > 0 ? {
                    resources = merge(
                      length(local.worker_container_limits[i]) > 0 ? { limits = local.worker_container_limits[i] } : {},
                      length(local.worker_container_requests[i]) > 0 ? { requests = local.worker_container_requests[i] } : {}
                    )
                  } : {}
                )
              ]
            },
            local.worker_scheduling[i]
          )
        }
      }
    )
  ]

  # ---- the RayCluster CR spec body ----------------------------------------------------
  # Field names are the CRD's own JSON keys (verified against the pinned
  # ray.io/v1 schema and raycluster_types.go). Values render ONLY when
  # declared so the operator's defaulting stays authoritative — except
  # authOptions (see above) and headGroupSpec.rayStartParams (module-
  # owned entries). No upgradeStrategy, no managedBy, no
  # headServiceAnnotations — unmodeled. Pulumi twin: rayClusterSpecBody.
  raycluster_spec = merge(
    { rayVersion = var.spec.ray_version },
    # suspend deletes head and worker PODS but keeps the declaration
    # (and, with GCS fault tolerance, the external state).
    try(var.spec.suspend, false) ? { suspend = true } : {},
    local.autoscaling_enabled ? { enableInTreeAutoscaling = true } : {},
    local.autoscaling_enabled && length(local.autoscaler_options) > 0 ? { autoscalerOptions = local.autoscaler_options } : {},
    local.gcs_ft_enabled ? { gcsFaultToleranceOptions = local.gcs_ft_options } : {},
    { authOptions = local.auth_options },
    { headGroupSpec = local.head_group_spec },
    length(local.worker_group_specs) > 0 ? { workerGroupSpecs = local.worker_group_specs } : {}
  )

  # ---- outputs -----------------------------------------------------------------------
  # All derived blind from the operator's naming contract: the head
  # Service is `<name>-head-svc` (GenerateHeadServiceName,
  # controllers/ray/utils/util.go — "%s-%s-%s" of name, "head", "svc").
  head_service = "${local.cluster_name}-head-svc"

  client_endpoint    = "${local.head_service}.${local.namespace}.svc.cluster.local:10001"
  dashboard_endpoint = "${local.head_service}.${local.namespace}.svc.cluster.local:8265"
  gcs_endpoint       = "${local.head_service}.${local.namespace}.svc.cluster.local:6379"

  # The bearer-token credential handle (token mode only): the
  # bring-your-own Secret when named, else the operator-generated Secret
  # named EXACTLY the cluster name — reconcileAuthSecret uses
  # utils.CheckName(instance.Name) verbatim (raycluster_controller.go;
  # CheckName only rewrites names past 50 characters, which the 40-char
  # name budget rules out). The data key is the operator's
  # RAY_AUTH_TOKEN_SECRET_KEY constant, "auth_token". Empty when auth is
  # disabled — no Secret exists then.
  auth_token_secret_name = local.token_auth_enabled ? (
    local.existing_token_secret_name != "" ? local.existing_token_secret_name : local.cluster_name
  ) : ""

  port_forward_command = "kubectl port-forward -n ${local.namespace} service/${local.head_service} 8265:8265"
}
