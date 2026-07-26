# Computed values for the KubernetesTektonOperator module. Every
# resolution here has an exact twin in the Pulumi module — keep them in
# lockstep: same ConfigMap patch, same patched Deployments, same outputs.
#
# HCL DISCIPLINE: conditional keys are contributed as merge() of
# `cond ? { key = value } : {}` singleton maps — a ternary whose branches
# are differently-shaped objects fails plan-time type unification.
# Optional nested blocks are read with try(): HCL's && does NOT
# short-circuit.

locals {
  # Pinned tektoncd/operator release tag.
  #
  # MUST stay in sync with the Pulumi module's OperatorRelease constant:
  # the TektonConfig surface the KubernetesTekton kind renders is
  # designed against this release's operator API. Always an exact
  # release TAG, never a branch — tag pinning keeps installs
  # reproducible.
  operator_release = "v0.80.0"

  # The released single-file manifest — the operator's OFFICIAL
  # distribution (the in-repo Helm chart is unpublished, version
  # "devel"). The GitHub release asset is immutable per tag; the old
  # storage.googleapis.com release host is dead.
  manifest_url = "https://github.com/tektoncd/operator/releases/download/${local.operator_release}/release.yaml"

  # Fixed handles from the release manifest (namespace and names are
  # baked into the manifest's own cross-references — not configurable).
  namespace                       = "tekton-operator"
  operator_deployment_name        = "tekton-operator"
  webhook_deployment_name         = "tekton-operator-webhook"
  config_defaults_config_map_name = "tekton-config-defaults"

  # ---- manifest documents ------------------------------------------------------
  # Each document keyed by its COMPOSED IDENTITY
  # `apiVersion//kind//name[//namespace]` — the exact ID form the kubectl
  # importer takes, so state addresses stay stable across manifest
  # reorderings AND the address key feeds the composed import ID blind
  # (from_address_key in the import map). Cluster-scoped documents render
  # 3-part keys (the importer rejects a trailing `//`).
  manifest_documents = [
    for doc in split("\n---\n", data.http.tekton_operator_manifest.response_body) :
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

  # ---- the AUTOINSTALL_COMPONENTS patch ------------------------------------------
  # ALWAYS patched to "false" (the release ships "true") — a design
  # invariant, not a knob: with auto-install on, the operator creates its
  # own TektonConfig named `config` at startup — the exact object the
  # KubernetesTekton declaration kind renders — and the two managers then
  # fight over the same fields through server-side apply. Disabling it
  # makes the declaration kind the single owner. Pulumi twin:
  # autoInstallTransformation.
  config_defaults_id = "v1//ConfigMap//${local.config_defaults_config_map_name}//${local.namespace}"

  original_config_defaults = local.documents_by_id[local.config_defaults_id]
  patched_config_defaults = merge(
    local.original_config_defaults,
    {
      data = merge(
        try(local.original_config_defaults.data, {}),
        { AUTOINSTALL_COMPONENTS = "false" }
      )
    }
  )

  # ---- typed overrides for the two Deployments -----------------------------------
  # Image overrides fold repo:tag; empty keeps the release manifest's
  # digest-pinned images.
  operator_image_repo = try(var.spec.operator_image.repo, "")
  operator_image_tag  = try(var.spec.operator_image.tag, "")
  operator_image = (
    local.operator_image_repo != "" && local.operator_image_tag != "" ?
    "${local.operator_image_repo}:${local.operator_image_tag}" :
    "${local.operator_image_repo}${local.operator_image_tag}"
  )

  webhook_image_repo = try(var.spec.webhook_image.repo, "")
  webhook_image_tag  = try(var.spec.webhook_image.tag, "")
  webhook_image = (
    local.webhook_image_repo != "" && local.webhook_image_tag != "" ?
    "${local.webhook_image_repo}:${local.webhook_image_tag}" :
    "${local.webhook_image_repo}${local.webhook_image_tag}"
  )

  operator_resources = try(var.spec.operator_resources, null) == null ? null : {
    for k, v in {
      requests = try(var.spec.operator_resources.requests, null) == null ? null : {
        for rk, rv in {
          cpu    = var.spec.operator_resources.requests.cpu
          memory = var.spec.operator_resources.requests.memory
        } : rk => rv if rv != null && rv != ""
      }
      limits = try(var.spec.operator_resources.limits, null) == null ? null : {
        for rk, rv in {
          cpu    = var.spec.operator_resources.limits.cpu
          memory = var.spec.operator_resources.limits.memory
        } : rk => rv if rv != null && rv != ""
      }
    } : k => v if v != null
  }

  webhook_resources = try(var.spec.webhook_resources, null) == null ? null : {
    for k, v in {
      requests = try(var.spec.webhook_resources.requests, null) == null ? null : {
        for rk, rv in {
          cpu    = var.spec.webhook_resources.requests.cpu
          memory = var.spec.webhook_resources.requests.memory
        } : rk => rv if rv != null && rv != ""
      }
      limits = try(var.spec.webhook_resources.limits, null) == null ? null : {
        for rk, rv in {
          cpu    = var.spec.webhook_resources.limits.cpu
          memory = var.spec.webhook_resources.limits.memory
        } : rk => rv if rv != null && rv != ""
      }
    } : k => v if v != null
  }

  pod_tolerations = [
    for t in try(var.spec.tolerations, []) : { for k, v in {
      key               = t.key != "" ? t.key : null
      operator          = t.operator != "" ? t.operator : null
      value             = t.value != "" ? t.value : null
      effect            = t.effect != "" ? t.effect : null
      tolerationSeconds = t.toleration_seconds
    } : k => v if v != null }
  ]

  # image_pull_secrets joined with both image overrides' own
  # pull_secret_name entries, deduplicated — a private image override
  # naturally travels with its own pull secret (Pulumi twin:
  # pullSecretNames).
  image_pull_secret_names = distinct(concat(
    try(var.spec.image_pull_secrets, []),
    try(var.spec.operator_image.pull_secret_name, "") != "" ? [var.spec.operator_image.pull_secret_name] : [],
    try(var.spec.webhook_image.pull_secret_name, "") != "" ? [var.spec.webhook_image.pull_secret_name] : []
  ))

  # ---- the patched Deployments ----------------------------------------------------
  # Both Deployments get the same pod-level overrides; per-Deployment
  # image/resources apply to EVERY container in that Deployment (the
  # operator Deployment's two containers share one image upstream).
  # Pulumi twin: deploymentTransformation.
  operator_deployment_id = "apps/v1//Deployment//${local.operator_deployment_name}//${local.namespace}"
  webhook_deployment_id  = "apps/v1//Deployment//${local.webhook_deployment_name}//${local.namespace}"

  pod_spec_overrides = merge(
    length(try(var.spec.node_selector, {})) > 0 ? { nodeSelector = var.spec.node_selector } : {},
    length(local.pod_tolerations) > 0 ? { tolerations = local.pod_tolerations } : {},
    length(local.image_pull_secret_names) > 0 ? {
      imagePullSecrets = [for n in local.image_pull_secret_names : { name = n }]
    } : {}
  )

  original_operator_deployment = local.documents_by_id[local.operator_deployment_id]
  patched_operator_deployment = merge(
    local.original_operator_deployment,
    {
      spec = merge(
        local.original_operator_deployment.spec,
        {
          template = merge(
            local.original_operator_deployment.spec.template,
            {
              spec = merge(
                local.original_operator_deployment.spec.template.spec,
                {
                  containers = [
                    for c in local.original_operator_deployment.spec.template.spec.containers :
                    merge(
                      c,
                      local.operator_image != "" ? { image = local.operator_image } : {},
                      local.operator_resources != null ? { resources = local.operator_resources } : {}
                    )
                  ]
                },
                local.pod_spec_overrides
              )
            }
          )
        }
      )
    }
  )

  original_webhook_deployment = local.documents_by_id[local.webhook_deployment_id]
  patched_webhook_deployment = merge(
    local.original_webhook_deployment,
    {
      spec = merge(
        local.original_webhook_deployment.spec,
        {
          template = merge(
            local.original_webhook_deployment.spec.template,
            {
              spec = merge(
                local.original_webhook_deployment.spec.template.spec,
                {
                  containers = [
                    for c in local.original_webhook_deployment.spec.template.spec.containers :
                    merge(
                      c,
                      local.webhook_image != "" ? { image = local.webhook_image } : {},
                      local.webhook_resources != null ? { resources = local.webhook_resources } : {}
                    )
                  ]
                },
                local.pod_spec_overrides
              )
            }
          )
        }
      )
    }
  )

  # The applied set: every document verbatim, with the ConfigMap and the
  # two Deployments swapped for their patched forms. merge() — never a
  # per-element ternary: a ternary whose branches are differently-shaped
  # objects fails HCL plan-time type unification.
  applied_documents = merge(
    local.documents_by_id,
    {
      (local.config_defaults_id)     = local.patched_config_defaults
      (local.operator_deployment_id) = local.patched_operator_deployment
      (local.webhook_deployment_id)  = local.patched_webhook_deployment
    }
  )
}
