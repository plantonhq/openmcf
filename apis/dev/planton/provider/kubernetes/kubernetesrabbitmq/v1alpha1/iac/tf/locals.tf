# Computed values for the KubernetesRabbitMq module. Every resolution here
# has an exact twin in the Pulumi module — keep them in lockstep: same
# rendered CR body, same outputs.
#
# HCL DISCIPLINE (applies to every conditional object in this file):
# conditional keys are contributed as merge() of `cond ? { key = value } : {}`
# singleton maps, or written as `key = cond ? value : null` inside ONE object
# literal pruned with `{ for k, v in {...} : k => v if v != null }`. The
# tempting alternative — a ternary whose branches are differently-shaped
# objects — fails plan-time type unification. The singleton/null-prune forms
# preserve every value's type: replicas and seconds render as YAML numbers,
# booleans as booleans.
#
# Optional nested blocks are read with try(): HCL's && does NOT
# short-circuit, so chained null checks still dereference the null.

locals {
  # The cluster name is metadata.name — the naming root the operator
  # derives every object from: the client Service `<name>`, the headless
  # Service `<name>-nodes`, the credentials Secret `<name>-default-user`,
  # the StatefulSet `<name>-server` and each pod's PVC
  # `persistence-<name>-server-<i>`.
  cluster_name = var.metadata.name
  namespace    = var.spec.namespace

  # Resource-identity labels stamped on the module-created objects
  # (namespace, RabbitmqCluster). The operator derives ITS objects'
  # identity from the cluster name; these labels tie the family back to
  # the Planton resource.
  labels = merge(
    {
      "planton.ai/resource"      = "true"
      "planton.ai/resource-name" = var.metadata.name
      "planton.ai/resource-kind" = "KubernetesRabbitMq"
    },
    var.metadata.id != null && var.metadata.id != "" ? { "planton.ai/resource-id" = var.metadata.id } : {},
    var.metadata.org != null && var.metadata.org != "" ? { "planton.ai/organization" = var.metadata.org } : {},
    var.metadata.env != null && var.metadata.env != "" ? { "planton.ai/environment" = var.metadata.env } : {}
  )

  replicas = try(coalesce(var.spec.replicas), 1)

  # The CR takes ONE image string — the shared ContainerImage folds into
  # `repo:tag`. Empty = omitted, so the operator's compiled-in default
  # (or the operator kind's fleet-wide default_rabbitmq_image) applies.
  image_repo = try(var.spec.image.repo, "")
  image_tag  = try(var.spec.image.tag, "")
  image = (
    local.image_repo != "" && local.image_tag != "" ?
    "${local.image_repo}:${local.image_tag}" :
    "${local.image_repo}${local.image_tag}"
  )

  # image.pull_secret_name joins image_pull_secrets, deduplicated — a
  # private image override naturally travels with its own pull secret.
  image_pull_secret_names = distinct(concat(
    try(var.spec.image_pull_secrets, []),
    try(var.spec.image.pull_secret_name, "") != "" ? [var.spec.image.pull_secret_name] : []
  ))

  # ---- TLS posture -----------------------------------------------------------
  # TLS is on when a certificate Secret is referenced; the plain listeners
  # close only when disable_non_tls_listeners additionally asks. The
  # effective client ports below follow (5671/15671 when closed).
  tls_secret_name        = try(var.spec.tls.secret_name, "")
  tls_enabled            = local.tls_secret_name != ""
  plain_listeners_closed = local.tls_enabled && try(var.spec.tls.disable_non_tls_listeners, false)

  amqp_port         = local.plain_listeners_closed ? 5671 : 5672
  management_port   = local.plain_listeners_closed ? 15671 : 15672
  management_scheme = local.plain_listeners_closed ? "https" : "http"

  # ---- service block ---------------------------------------------------------
  # Rendered only when the spec diverges from the operator's own defaults
  # (ClusterIP, no annotations) — presence discipline keeps the CR minimal.
  service_type = {
    ""                 = ""
    "cluster_ip"       = ""
    "load_balancer"    = "LoadBalancer"
    "node_port"        = "NodePort"
  }[try(coalesce(var.spec.service.type), "")]

  ip_family_policy = {
    ""                             = ""
    "ip_family_policy_unspecified" = ""
    "single_stack"                 = "SingleStack"
    "prefer_dual_stack"            = "PreferDualStack"
    "require_dual_stack"           = "RequireDualStack"
  }[try(coalesce(var.spec.service.ip_family_policy), "")]

  service_body = { for k, v in {
    type           = local.service_type != "" ? local.service_type : null
    annotations    = length(try(var.spec.service.annotations, {})) > 0 ? var.spec.service.annotations : null
    labels         = length(try(var.spec.service.labels, {})) > 0 ? var.spec.service.labels : null
    ipFamilyPolicy = local.ip_family_policy != "" ? local.ip_family_policy : null
  } : k => v if v != null }

  # ---- persistence -----------------------------------------------------------
  # The ephemeral posture is the CR's own storage-0 + emptyDir mechanism;
  # otherwise the resolved disk size (spec default 10Gi) and the optional
  # storage class render as the persistent volume claim template. ONE
  # null-pruned object literal — never a ternary whose branches are
  # differently-shaped objects (plan-time type unification).
  is_ephemeral = try(var.spec.ephemeral, false)

  persistence_body = { for k, v in {
    storage          = local.is_ephemeral ? "0" : try(coalesce(var.spec.disk_size), "10Gi")
    emptyDir         = local.is_ephemeral ? {} : null
    storageClassName = (!local.is_ephemeral && try(var.spec.storage_class, "") != "") ? var.spec.storage_class : null
  } : k => v if v != null }

  # ---- resources (requests/limits) ------------------------------------------
  resources_body = try(var.spec.resources, null) == null ? null : {
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

  # ---- rabbitmq configuration block ------------------------------------------
  rabbitmq_body = { for k, v in {
    additionalPlugins = length(try(var.spec.configuration.additional_plugins, [])) > 0 ? var.spec.configuration.additional_plugins : null
    additionalConfig  = try(var.spec.configuration.additional_config, "") != "" ? var.spec.configuration.additional_config : null
    advancedConfig    = try(var.spec.configuration.advanced_config, "") != "" ? var.spec.configuration.advanced_config : null
    envConfig         = try(var.spec.configuration.env_config, "") != "" ? var.spec.configuration.env_config : null
    erlangInetConfig  = try(var.spec.configuration.erlang_inet_config, "") != "" ? var.spec.configuration.erlang_inet_config : null
  } : k => v if v != null }

  # ---- tls block --------------------------------------------------------------
  tls_body = { for k, v in {
    secretName             = local.tls_enabled ? local.tls_secret_name : null
    caSecretName           = try(var.spec.tls.ca_secret_name, "") != "" ? var.spec.tls.ca_secret_name : null
    disableNonTLSListeners = local.plain_listeners_closed ? true : null
  } : k => v if v != null }

  # ---- tolerations -------------------------------------------------------------
  tolerations = [
    for t in try(var.spec.tolerations, []) : { for k, v in {
      key               = t.key != "" ? t.key : null
      operator          = t.operator != "" ? t.operator : null
      value             = t.value != "" ? t.value : null
      effect            = t.effect != "" ? t.effect : null
      tolerationSeconds = t.toleration_seconds
    } : k => v if v != null }
  ]

  # ---- affinity ----------------------------------------------------------------
  # node_selector renders as REQUIRED node affinity with one In-match per
  # label (the CR has no nodeSelector field; behaviorally identical for
  # exact matches — Pulumi twin renders the same shape).
  # spread_across_nodes renders as REQUIRED pod anti-affinity on the
  # operator's own `app.kubernetes.io/name: <cluster>` pod label over the
  # hostname topology.
  node_affinity_body = length(try(var.spec.node_selector, {})) == 0 ? null : {
    requiredDuringSchedulingIgnoredDuringExecution = {
      nodeSelectorTerms = [{
        matchExpressions = [
          for k, v in var.spec.node_selector : {
            key      = k
            operator = "In"
            values   = [v]
          }
        ]
      }]
    }
  }

  pod_anti_affinity_body = try(var.spec.spread_across_nodes, false) != true ? null : {
    requiredDuringSchedulingIgnoredDuringExecution = [{
      topologyKey = "kubernetes.io/hostname"
      labelSelector = {
        matchLabels = {
          "app.kubernetes.io/name" = local.cluster_name
        }
      }
    }]
  }

  affinity_body = { for k, v in {
    nodeAffinity    = local.node_affinity_body
    podAntiAffinity = local.pod_anti_affinity_body
  } : k => v if v != null }

  # ---- secret backend -----------------------------------------------------------
  vault_backend = try(var.spec.secret_backend.vault, null) == null ? null : { for k, v in {
    role            = var.spec.secret_backend.vault.role
    defaultUserPath = var.spec.secret_backend.vault.default_user_path
    annotations     = length(try(var.spec.secret_backend.vault.annotations, {})) > 0 ? var.spec.secret_backend.vault.annotations : null
    tls = try(var.spec.secret_backend.vault.pki_issuer_path, "") == "" ? null : {
      pkiIssuerPath = var.spec.secret_backend.vault.pki_issuer_path
    }
  } : k => v if v != null }

  external_secret_name = try(var.spec.secret_backend.external_secret_name, "")

  # ONE null-pruned object literal — the chained per-arm ternary selecting
  # differently-shaped objects is the same plan-time type-unification
  # failure as the two-branch form.
  secret_backend_body = { for k, v in {
    vault          = local.vault_backend
    externalSecret = local.external_secret_name != "" ? { name = local.external_secret_name } : null
  } : k => v if v != null }

  # ---- the RabbitmqCluster manifest ----------------------------------------------
  # Optionals with operator-default values render only when the incoming
  # spec DECLARES them (presence carried by the converter's snake_case
  # object) — but tfvars objects flatten presence, so the two
  # operator-default-valued knobs (termination grace 604800, delay start
  # 30) render whenever they differ from the operator defaults; equal
  # values are omitted to keep the CR minimal and byte-identical with the
  # Pulumi rendering of an unset field.
  rabbitmq_cluster_manifest = {
    apiVersion = "rabbitmq.com/v1beta1"
    kind       = "RabbitmqCluster"
    metadata = {
      name      = local.cluster_name
      namespace = local.namespace
      labels    = local.labels
    }
    spec = { for k, v in {
      replicas         = local.replicas
      image            = local.image != "" ? local.image : null
      imagePullSecrets = length(local.image_pull_secret_names) > 0 ? [for n in local.image_pull_secret_names : { name = n }] : null
      service          = length(local.service_body) > 0 ? local.service_body : null
      persistence      = local.persistence_body
      resources        = local.resources_body
      rabbitmq         = length(local.rabbitmq_body) > 0 ? local.rabbitmq_body : null
      tls              = length(local.tls_body) > 0 ? local.tls_body : null
      tolerations      = length(local.tolerations) > 0 ? local.tolerations : null
      affinity         = length(local.affinity_body) > 0 ? local.affinity_body : null

      terminationGracePeriodSeconds = try(coalesce(var.spec.termination_grace_period_seconds), 604800) != 604800 ? var.spec.termination_grace_period_seconds : null
      delayStartSeconds             = try(coalesce(var.spec.delay_start_seconds), 30) != 30 ? var.spec.delay_start_seconds : null
      skipPostDeploySteps           = try(var.spec.skip_post_deploy_steps, false) ? true : null
      autoEnableAllFeatureFlags     = try(var.spec.auto_enable_all_feature_flags, false) ? true : null
      secretBackend                 = length(local.secret_backend_body) > 0 ? local.secret_backend_body : null
    } : k => v if v != null }
  }

  # ---- outputs-facing handles -------------------------------------------------
  # Operator naming contract: client Service `<name>`, headless Service
  # `<name>-nodes`, credentials Secret `<name>-default-user`. The Vault
  # backend replaces the generated Secret entirely (empty handle).
  service_name             = local.cluster_name
  headless_service_name    = "${local.cluster_name}-nodes"
  amqp_endpoint            = "${local.service_name}.${local.namespace}.svc.cluster.local:${local.amqp_port}"
  management_endpoint      = "${local.management_scheme}://${local.service_name}.${local.namespace}.svc.cluster.local:${local.management_port}"
  default_user_secret_name = local.vault_backend != null ? "" : "${local.cluster_name}-default-user"
  port_forward_command     = "kubectl port-forward svc/${local.service_name} -n ${local.namespace} ${local.management_port}:${local.management_port}"
}
