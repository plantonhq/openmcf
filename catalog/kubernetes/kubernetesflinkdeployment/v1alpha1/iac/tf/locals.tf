# Computed values for the KubernetesFlinkDeployment module. Every
# resolution here has an exact twin in the Pulumi module (locals.go /
# flinkdeployment_cr.go) — keep them in lockstep: same rendered CR body,
# same outputs.
#
# HCL DISCIPLINE: conditional keys are contributed via merge() of
# `cond ? { key = value } : {}` singleton maps, or one object literal
# pruned with `{ for k, v in {...} : k => v if v != null }` — a ternary
# whose branches are differently-shaped objects fails plan-time type
# unification. Optional nested blocks are read with try(): HCL's && does
# NOT short-circuit. Optional scalars in string templates resolve with
# try(coalesce(...)).

locals {
  # metadata.name is the CR name — the operator's naming root: the
  # JobManager REST Service is `<name>-rest` and TaskManager pods take
  # `<name>-taskmanager-N-M`.
  flinkdeployment_name = var.metadata.name
  namespace            = var.spec.namespace
  api_version          = "flink.apache.org/v1beta1"

  # Planton governance labels on the module-created objects.
  labels = merge(
    {
      "planton.ai/resource"      = "true"
      "planton.ai/resource-name" = var.metadata.name
      "planton.ai/resource-kind" = "KubernetesFlinkDeployment"
    },
    try(var.metadata.id, "") != "" ? { "planton.ai/resource-id" = var.metadata.id } : {},
    try(var.metadata.org, "") != "" ? { "planton.ai/organization" = var.metadata.org } : {},
    try(var.metadata.env, "") != "" ? { "planton.ai/environment" = var.metadata.env } : {}
  )

  # ---- image / serviceAccount / mode -------------------------------------------------
  # VERSION/IMAGE LOCKSTEP: the default image derives from flink_version
  # (`flink:<major>.<minor>` — v2_1 → flink:2.1) by stripping the "v" and
  # replacing "_" with ".". A custom image must carry exactly that Flink
  # version — the operator shapes its submission protocol from the
  # declared version and a mismatch fails at runtime, not at apply.
  default_image = "flink:${replace(trimprefix(var.spec.flink_version, "v"), "_", ".")}"
  image         = try(var.spec.image, "") != "" ? var.spec.image : local.default_image

  # Empty = "flink" — the account the KubernetesFlinkOperator's chart
  # creates with reconcile RBAC.
  service_account = try(coalesce(var.spec.service_account, "flink"), "flink")

  # "native" is the CR default — mode renders only on divergence.
  standalone = try(var.spec.mode, "") == "standalone"

  # ---- flinkConfiguration -------------------------------------------------------------
  # Module-owned keys rendered from spec.state. The typed fields are the
  # truth: they merge LAST, so a colliding user flink_configuration entry
  # loses, deliberately.
  #
  # "high-availability.type: kubernetes" is the CURRENT key form at the
  # pinned operator (its own test/e2e manifests use exactly this key —
  # test-deployment-key-value-configuration.yaml, multi-sessionjob.yaml —
  # and examples/basic-checkpoint-ha.yaml carries the same key in nested
  # YAML form); the legacy `high-availability: <factory-class>` form is
  # not used.
  #
  # s3.path.style.access renders with its effective value whenever s3 is
  # set: the spec default (true) is correct for in-cluster object stores;
  # explicit false is the AWS-S3-itself posture. Credentials NEVER go
  # here — flinkConfiguration renders into a ConfigMap in clear text;
  # they ride pod env from Secret refs (see the pod template).
  s3_set = try(var.spec.state.s3, null) != null

  state_flink_configuration = merge(
    try(var.spec.state.checkpoints_dir, "") != "" ? { "state.checkpoints.dir" = var.spec.state.checkpoints_dir } : {},
    try(var.spec.state.savepoints_dir, "") != "" ? { "state.savepoints.dir" = var.spec.state.savepoints_dir } : {},
    try(var.spec.state.high_availability.enabled, false) ? {
      "high-availability.type"       = "kubernetes"
      "high-availability.storageDir" = var.spec.state.high_availability.storage_dir
    } : {},
    local.s3_set ? {
      "s3.endpoint"          = var.spec.state.s3.endpoint
      "s3.path.style.access" = tostring(try(coalesce(var.spec.state.s3.path_style_access, true), true))
    } : {}
  )

  flink_configuration = merge(
    try(var.spec.flink_configuration, {}),
    local.state_flink_configuration
  )

  # ---- spec.job ------------------------------------------------------------------------
  # Set = APPLICATION cluster (the cluster runs exactly this job); absent
  # = SESSION cluster. state renders only on divergence from the CR
  # default "running" (i.e. "suspended"), upgradeMode only on divergence
  # from "stateless" (i.e. "last-state"/"savepoint"). Null-pruned, never
  # merge() over conditional `{...}:{}` lists — that idiom silently
  # UNIFIES primitive-only sibling objects into map(string) and fails
  # plan-time type-checking on mixed shapes (the discipline header's
  # second class).
  job_block = try(var.spec.job, null) == null ? null : {
    for k, v in {
      jarURI                = var.spec.job.jar_uri
      entryClass            = try(var.spec.job.entry_class, "") != "" ? var.spec.job.entry_class : null
      args                  = length(try(var.spec.job.args, [])) > 0 ? var.spec.job.args : null
      parallelism           = try(var.spec.job.parallelism, null)
      state                 = try(var.spec.job.state, "") == "suspended" ? "suspended" : null
      upgradeMode           = contains(["last-state", "savepoint"], try(coalesce(var.spec.job.upgrade_mode, "stateless"), "stateless")) ? var.spec.job.upgrade_mode : null
      initialSavepointPath  = try(var.spec.job.initial_savepoint_path, "") != "" ? var.spec.job.initial_savepoint_path : null
      allowNonRestoredState = try(var.spec.job.allow_non_restored_state, false) ? true : null
      savepointTriggerNonce = try(var.spec.job.savepoint_trigger_nonce, 0) != 0 ? var.spec.job.savepoint_trigger_nonce : null
    } : k => v if v != null
  }

  # ---- jobManager / taskManager --------------------------------------------------------
  # The CR's Resource takes cpu as a NUMBER and memory as a string, and
  # the operator's validator REQUIRES resource memory on BOTH tiers
  # (DefaultValidator: "JobManager resource memory must be defined…" /
  # "TaskManager resource memory must be defined…") — so both blocks
  # render ALWAYS, defaulting to cpu 1.0 / memory "2Gi" (the spec's
  # default_container_resources values) when the tier or its resources
  # are unset.
  #
  # quantityToCores: the spec carries Kubernetes CPU quantities
  # (requests.cpu preferred, else limits.cpu); the CR wants cores as a
  # number — "1000m" → 1.0, "500m" → 0.5, "2" → 2.0. Millicores are
  # detected by regex and divided by 1000; anything else parses as a
  # plain number. Identical semantics in the Pulumi twin
  # (quantityToCores in flinkdeployment_cr.go).
  jm_cpu_quantity = try(var.spec.job_manager.resources.requests.cpu, "") != "" ? var.spec.job_manager.resources.requests.cpu : try(var.spec.job_manager.resources.limits.cpu, "")
  jm_cpu          = local.jm_cpu_quantity == "" ? 1.0 : (can(regex("^([0-9.]+)m$", local.jm_cpu_quantity)) ? tonumber(replace(local.jm_cpu_quantity, "m", "")) / 1000 : tonumber(local.jm_cpu_quantity))

  # memory: limits.memory preferred (the ceiling Flink sizes its JVM
  # from), else requests.memory, else the "2Gi" default.
  jm_memory = try(var.spec.job_manager.resources.limits.memory, "") != "" ? var.spec.job_manager.resources.limits.memory : (try(var.spec.job_manager.resources.requests.memory, "") != "" ? var.spec.job_manager.resources.requests.memory : "2Gi")

  tm_cpu_quantity = try(var.spec.task_manager.resources.requests.cpu, "") != "" ? var.spec.task_manager.resources.requests.cpu : try(var.spec.task_manager.resources.limits.cpu, "")
  tm_cpu          = local.tm_cpu_quantity == "" ? 1.0 : (can(regex("^([0-9.]+)m$", local.tm_cpu_quantity)) ? tonumber(replace(local.tm_cpu_quantity, "m", "")) / 1000 : tonumber(local.tm_cpu_quantity))

  tm_memory = try(var.spec.task_manager.resources.limits.memory, "") != "" ? var.spec.task_manager.resources.limits.memory : (try(var.spec.task_manager.resources.requests.memory, "") != "" ? var.spec.task_manager.resources.requests.memory : "2Gi")

  # JobManager replicas render only past 1 (standbys — the spec-level
  # rule requires state.high_availability for them).
  job_manager_block = merge(
    { resource = { cpu = local.jm_cpu, memory = local.jm_memory } },
    try(var.spec.job_manager.replicas, null) != null && try(var.spec.job_manager.replicas, 1) > 1 ? { replicas = var.spec.job_manager.replicas } : {}
  )

  # TaskManager replicas render whenever declared — meaningful in
  # standalone mode only (native mode derives worker count from the
  # job's parallelism).
  task_manager_block = merge(
    { resource = { cpu = local.tm_cpu, memory = local.tm_memory } },
    try(var.spec.task_manager.replicas, null) != null ? { replicas = var.spec.task_manager.replicas } : {}
  )

  # ---- podTemplate ----------------------------------------------------------------------
  # Rendered only when it would carry something: scheduling, image-pull
  # Secrets, or the S3 credential env. "flink-main-container" is the
  # operator's merge contract for the main container (its own
  # examples/pod-template.yaml: "Do not change the main container
  # name").
  #
  # S3 credentials ride pod env from Secret refs — NEVER into
  # flinkConfiguration, which renders into a ConfigMap in clear text.
  # ENABLE_BUILT_IN_PLUGINS activates the named S3 filesystem plugin
  # from the image's bundled (disabled-by-default) plugin set.
  # Credential env entries carry valueFrom; the plugin entry carries
  # value. Terraform cannot concat those object shapes (or ternary
  # them against null/[]) — "Inconsistent conditional result types"
  # (verified live once builtin_plugin_jar reached the variable
  # schema). Encode each entry as JSON (uniform string list), then
  # decode so the rendered podTemplate gets the right keys.
  s3_env_json = local.s3_set ? concat(
    [
      jsonencode({
        name = "AWS_ACCESS_KEY_ID"
        valueFrom = {
          secretKeyRef = {
            name = var.spec.state.s3.access_key_secret.name
            key  = var.spec.state.s3.access_key_secret.key
          }
        }
      }),
      jsonencode({
        name = "AWS_SECRET_ACCESS_KEY"
        valueFrom = {
          secretKeyRef = {
            name = var.spec.state.s3.secret_key_secret.name
            key  = var.spec.state.s3.secret_key_secret.key
          }
        }
      }),
    ],
    try(var.spec.state.s3.builtin_plugin_jar, "") != "" ? [
      jsonencode({
        name  = "ENABLE_BUILT_IN_PLUGINS"
        value = var.spec.state.s3.builtin_plugin_jar
      })
    ] : []
  ) : []

  s3_env = [for e in local.s3_env_json : jsondecode(e)]

  scheduling_tolerations = [
    for t in try(var.spec.scheduling.tolerations, []) : { for k, v in {
      key               = t.key != "" ? t.key : null
      operator          = t.operator != "" ? t.operator : null
      value             = t.value != "" ? t.value : null
      effect            = t.effect != "" ? t.effect : null
      tolerationSeconds = t.toleration_seconds
    } : k => v if v != null }
  ]

  pod_template_spec = merge(
    length(try(var.spec.image_pull_secrets, [])) > 0 ? {
      imagePullSecrets = [for secret_name in var.spec.image_pull_secrets : { name = secret_name }]
    } : {},
    length(try(var.spec.scheduling.node_selector, {})) > 0 ? { nodeSelector = var.spec.scheduling.node_selector } : {},
    length(local.scheduling_tolerations) > 0 ? { tolerations = local.scheduling_tolerations } : {},
    try(var.spec.scheduling.priority_class_name, "") != "" ? { priorityClassName = var.spec.scheduling.priority_class_name } : {},
    length(local.s3_env) > 0 ? {
      containers = [{ name = "flink-main-container", env = local.s3_env }]
    } : {}
  )

  render_pod_template = try(var.spec.scheduling, null) != null || length(try(var.spec.image_pull_secrets, [])) > 0 || local.s3_set

  pod_template = {
    apiVersion = "v1"
    kind       = "Pod"
    metadata   = { name = "pod-template" }
    spec       = local.pod_template_spec
  }

  # ---- the FlinkDeployment CR spec body --------------------------------------------------
  # Field names are the CR's own JSON keys (verified against the pinned
  # operator's API classes: FlinkDeploymentSpec / JobSpec /
  # JobManagerSpec / TaskManagerSpec / Resource). Values render ONLY
  # when declared so the operator's defaulting stays authoritative —
  # except flinkVersion/image/serviceAccount (ALWAYS: the deployment's
  # identity) and jobManager/taskManager (ALWAYS: the validator requires
  # resource on both tiers). Pulumi twin: flinkDeploymentSpecBody.
  flinkdeployment_spec = {
    for k, v in {
      flinkVersion       = var.spec.flink_version
      image              = local.image
      serviceAccount     = local.service_account
      mode               = local.standalone ? "standalone" : null
      flinkConfiguration = length(local.flink_configuration) > 0 ? local.flink_configuration : null
      job                = local.job_block
      restartNonce       = try(var.spec.restart_nonce, 0) != 0 ? var.spec.restart_nonce : null
      jobManager         = local.job_manager_block
      taskManager        = local.task_manager_block
      logConfiguration   = length(try(var.spec.log_configuration, {})) > 0 ? var.spec.log_configuration : null
      podTemplate        = local.render_pod_template ? local.pod_template : null
    } : k => v if v != null
  }

  # ---- outputs ----------------------------------------------------------------------------
  # All derived blind from the operator's naming contract: the JobManager
  # REST Service is `<name>-rest` (IngressUtils' REST_SVC_NAME_SUFFIX),
  # serving the Flink REST API and web UI on 8081.
  rest_service  = "${local.flinkdeployment_name}-rest"
  rest_endpoint = "${local.rest_service}.${local.namespace}.svc.cluster.local:8081"

  port_forward_command = "kubectl port-forward -n ${local.namespace} service/${local.rest_service} 8081:8081"
}
