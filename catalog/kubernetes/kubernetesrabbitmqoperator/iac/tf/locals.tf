# Computed values for the KubernetesRabbitMqOperator module. Every
# resolution here has an exact twin in the Pulumi module — keep them in
# lockstep: same patched Deployment, same outputs.
#
# HCL DISCIPLINE: conditional keys are contributed as merge() of
# `cond ? { key = value } : {}` singleton maps, or written as
# `key = cond ? value : null` inside ONE object literal pruned with
# `{ for k, v in {...} : k => v if v != null }` — a ternary whose branches
# are differently-shaped objects fails plan-time type unification.
# Optional nested blocks are read with try(): HCL's && does NOT
# short-circuit.

locals {
  # Pinned rabbitmq/cluster-operator release tag.
  #
  # MUST stay in sync with `rabbitmq_cluster_operator_release` in
  # pkg/kubernetes/kubernetestypes/Makefile and the Pulumi module's
  # OperatorRelease constant, so the installed CRD schema matches the
  # crd2pulumi-generated typed SDK that KubernetesRabbitMq is built
  # against. Always an exact release TAG, never a branch — tag pinning
  # keeps installs reproducible.
  operator_release = "v2.22.3"

  # The released single-file manifest — the operator's OFFICIAL
  # distribution (it has no Helm chart). The release pipeline pins the
  # operator image tag inside this asset.
  manifest_url = "https://github.com/rabbitmq/cluster-operator/releases/download/${local.operator_release}/cluster-operator.yml"

  # Fixed handles from the release manifest (namespace and names are
  # baked into the manifest's own cross-references — not configurable).
  namespace            = "rabbitmq-system"
  deployment_name      = "rabbitmq-cluster-operator"
  metrics_service_name = "rabbitmq-cluster-operator-metrics-service"
  crd_name             = "rabbitmqclusters.rabbitmq.com"
  metrics_port         = 8080
  metrics_endpoint     = "http://${local.metrics_service_name}.${local.namespace}.svc.cluster.local:${local.metrics_port}/metrics"

  # ---- manifest documents ------------------------------------------------------
  # Each document keyed by its COMPOSED IDENTITY
  # `apiVersion//kind//name[//namespace]` — the exact ID form the kubectl
  # importer takes, so state addresses stay stable across manifest
  # reorderings AND the address key feeds the composed import ID blind
  # (from_address_key in the import map). Cluster-scoped documents render
  # 3-part keys (the importer rejects a trailing `//`).
  manifest_documents = [
    for doc in split("\n---\n", data.http.cluster_operator_manifest.response_body) :
    yamldecode(doc)
    if trimspace(doc) != "" && can(yamldecode(doc).kind)
  ]

  documents_by_id = {
    for doc in local.manifest_documents :
    join("//", concat(
      [doc.apiVersion, doc.kind, doc.metadata.name],
      try(doc.metadata.namespace, "") != "" ? [doc.metadata.namespace] : []
    )) => doc
  }

  # ---- typed overrides for the operator Deployment ------------------------------
  # Env entries appended to the manifest's own (which carry
  # OPERATOR_NAMESPACE via fieldRef). Absent OPERATOR_SCOPE_NAMESPACE =
  # the operator watches ALL namespaces (the upstream default).
  extra_env = concat(
    length(try(var.spec.watch_namespaces, [])) > 0 ? [{
      name  = "OPERATOR_SCOPE_NAMESPACE"
      value = join(",", var.spec.watch_namespaces)
    }] : [],
    try(var.spec.default_rabbitmq_image, "") != "" ? [{
      name  = "DEFAULT_RABBITMQ_IMAGE"
      value = var.spec.default_rabbitmq_image
    }] : [],
    try(var.spec.default_user_updater_image, "") != "" ? [{
      name  = "DEFAULT_USER_UPDATER_IMAGE"
      value = var.spec.default_user_updater_image
    }] : []
  )

  # The operator image override folds repo:tag; empty keeps the release
  # manifest's pinned image.
  operator_image_repo = try(var.spec.operator_image.repo, "")
  operator_image_tag  = try(var.spec.operator_image.tag, "")
  operator_image = (
    local.operator_image_repo != "" && local.operator_image_tag != "" ?
    "${local.operator_image_repo}:${local.operator_image_tag}" :
    "${local.operator_image_repo}${local.operator_image_tag}"
  )

  operator_resources = try(var.spec.resources, null) == null ? null : {
    for k, v in {
      requests = try(var.spec.resources.requests, null) == null ? null : {
        for rk, rv in {
          cpu    = var.spec.resources.requests.cpu
          memory = var.spec.resources.requests.memory
        } : rk => rv if rv != null && rv != ""
      }
      limits = try(var.spec.resources.limits, null) == null ? null : {
        for rk, rv in {
          cpu    = var.spec.resources.limits.cpu
          memory = var.spec.resources.limits.memory
        } : rk => rv if rv != null && rv != ""
      }
    } : k => v if v != null
  }

  operator_tolerations = [
    for t in try(var.spec.tolerations, []) : { for k, v in {
      key               = t.key != "" ? t.key : null
      operator          = t.operator != "" ? t.operator : null
      value             = t.value != "" ? t.value : null
      effect            = t.effect != "" ? t.effect : null
      tolerationSeconds = t.toleration_seconds
    } : k => v if v != null }
  ]

  image_pull_secret_names = distinct(concat(
    try(var.spec.image_pull_secrets, []),
    try(var.spec.operator_image.pull_secret_name, "") != "" ? [var.spec.operator_image.pull_secret_name] : []
  ))

  # ---- the patched Deployment ----------------------------------------------------
  # The release manifest ships exactly one container ("operator") — the
  # patch rebuilds the document with the overrides merged in and leaves
  # every other field untouched (faithful distribution). Pulumi twin: the
  # deploymentTransformation function.
  deployment_id = "apps/v1//Deployment//${local.deployment_name}//${local.namespace}"

  original_deployment = local.documents_by_id[local.deployment_id]
  original_pod_spec   = local.original_deployment.spec.template.spec
  original_container  = local.original_deployment.spec.template.spec.containers[0]

  patched_container = merge(
    local.original_container,
    length(local.extra_env) > 0 ? {
      env = concat(try(local.original_container.env, []), local.extra_env)
    } : {},
    local.operator_image != "" ? { image = local.operator_image } : {},
    local.operator_resources != null ? { resources = local.operator_resources } : {}
  )

  patched_pod_spec = merge(
    local.original_pod_spec,
    { containers = [local.patched_container] },
    length(try(var.spec.node_selector, {})) > 0 ? { nodeSelector = var.spec.node_selector } : {},
    length(local.operator_tolerations) > 0 ? { tolerations = local.operator_tolerations } : {},
    length(local.image_pull_secret_names) > 0 ? {
      imagePullSecrets = [for n in local.image_pull_secret_names : { name = n }]
    } : {}
  )

  patched_deployment = merge(
    local.original_deployment,
    {
      spec = merge(
        local.original_deployment.spec,
        {
          template = merge(
            local.original_deployment.spec.template,
            { spec = local.patched_pod_spec }
          )
        }
      )
    }
  )

  # The applied set: every document verbatim, with the Deployment
  # swapped for its patched form. merge() — never a per-element ternary:
  # a ternary whose branches are differently-shaped objects (the patched
  # Deployment vs a ServiceAccount) fails HCL plan-time type unification.
  applied_documents = merge(
    local.documents_by_id,
    { (local.deployment_id) = local.patched_deployment }
  )
}
