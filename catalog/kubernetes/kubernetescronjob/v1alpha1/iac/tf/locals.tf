# Computed values shared across the module's resources.
#
# The load-bearing decisions here:
#  - selector_labels are OUR pod-location labels (metadata.name based). The Job
#    controller adds its own controller-uid/job-name labels to each stamped-out
#    run's pods; these labels are stamped on the pod template so
#    `kubectl get pods -l` and `kubectl logs -l` find every run's pods without
#    knowing controller internals. There is deliberately NO version label.
#  - all_containers normalizes the app container (default name "app") and the
#    sidecars into ONE uniform list, so cron_job.tf renders every container
#    through a single dynamic block — identical semantics for app and sidecars.
#  - pod_volumes derives the pod-level volume list from the union of every
#    container's mounts, de-duplicated by name (two containers sharing an
#    EmptyDir declare the same mount name and source).

locals {
  resource_id = (
    var.metadata.id != null && var.metadata.id != ""
    ? var.metadata.id
    : var.metadata.name
  )

  selector_labels = {
    "app"           = var.metadata.name
    "resource_name" = var.metadata.name
  }

  base_labels = merge(local.selector_labels, {
    "resource"      = "true"
    "resource_id"   = local.resource_id
    "resource_kind" = "KubernetesCronJob"
  })

  org_label = (
    var.metadata.org != null && var.metadata.org != ""
  ) ? { "organization" = var.metadata.org } : {}

  env_label = (
    var.metadata.env != null && try(var.metadata.env, "") != ""
  ) ? { "environment" = var.metadata.env } : {}

  final_labels = merge(local.base_labels, local.org_label, local.env_label)

  namespace = var.spec.namespace

  # Satellite resource names, prefixed with metadata.name so multiple
  # instances sharing a namespace never collide.
  env_secret_name        = "${var.metadata.name}-env-secrets"
  image_pull_secret_name = "${var.metadata.name}-image-pull"

  # ---------------------------------------------------------------------------
  # Container normalization: app first (named "app" unless the spec names it),
  # then sidecars in declared order. Kubernetes shows containers in declaration
  # order and tooling conventionally treats the first as primary. All work
  # definition lives under spec.job_template.
  # ---------------------------------------------------------------------------
  app_container = merge(var.spec.job_template.container.app, {
    name = try(var.spec.job_template.container.app.name, "") != "" ? var.spec.job_template.container.app.name : "app"
  })

  all_containers = concat([local.app_container], try(var.spec.job_template.container.sidecars, []))

  init_containers = try(var.spec.job_template.pod.init_containers, [])

  # Literal secret env values across ALL containers (app, sidecars, init),
  # materialized into one workload-scoped Secret. secretRef entries are wired
  # directly as env references and never pass through this map.
  env_secret_data = merge([
    for c in concat(local.all_containers, local.init_containers) : {
      for s in try(c.env.secrets, []) : s.name => s.value
      if try(s.value, "") != "" && try(s.secret_ref, null) == null
    }
  ]...)

  # Pod-level volumes: union of every container's mounts, first declaration of
  # a name wins. A mount with NO source references a volume defined elsewhere
  # and contributes no pod volume.
  volume_mounts_flat = flatten([
    for c in concat(local.all_containers, local.init_containers) : try(c.volume_mounts, [])
  ])

  pod_volumes = {
    for vm in local.volume_mounts_flat : vm.name => vm...
  }

  # Pod-level image pull secrets: spec-listed names plus the module-created
  # docker-config secret when configured. ServiceAccount-attached pull secrets
  # need no entry here.
  image_pull_secret_names = concat(
    try(var.spec.job_template.pod.image_pull_secrets, []),
    local.create_image_pull_secret ? [local.image_pull_secret_name] : []
  )

  create_image_pull_secret = try(var.docker_config_json, "") != ""

  # Pod template labels: controller labels win over user pod labels so a user
  # label can never break pod location.
  pod_template_labels = merge(try(var.spec.job_template.pod.labels, {}), local.final_labels)

  # Restart policy defaults to "Never": one pod per attempt, so failed pods
  # survive for post-mortem inspection — and pod_failure_policy requires it.
  # coalesce treats both null (unset optional) and "" as absent.
  restart_policy = try(coalesce(var.spec.job_template.restart_policy, "Never"), "Never")

  # Concurrency policy defaults to "Forbid" — deliberately safer than
  # upstream's "Allow". Overlapping cron runs are the classic
  # scheduled-workload incident (two backups writing the same target, two
  # migrations racing), so the spec documents Forbid as OUR default and the
  # module applies it explicitly rather than letting Kubernetes fall back to
  # Allow when the field is unset.
  concurrency_policy = try(var.spec.concurrency_policy, "") != "" && try(var.spec.concurrency_policy, null) != null ? var.spec.concurrency_policy : "Forbid"
}

variable "docker_config_json" {
  description = "Docker registry credential (dockerconfigjson) injected by the platform at deploy time; empty when pulling public images or when pull secrets are attached to the ServiceAccount."
  type        = string
  default     = ""
}
