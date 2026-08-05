# Computed values for the KubernetesArgoWorkflows module. Every resolution
# here has an exact twin in the Pulumi module's locals.go / values.go —
# keep them in lockstep.
#
# SECRET DISCIPLINE (load-bearing): nothing in this module transports
# credential material. Artifact-store keys ride the chart's own secret
# SELECTORS (name/key pairs resolved by the workload at runtime); the
# archive database credentials ride the same selector contract; SSO client
# credentials (when opened via helm_values) ride the chart's secret
# selectors too. Keyless arms lean on ambient pod identity
# (IRSA/workload identity) on the RUNNER service account.
#
# HCL DISCIPLINE (applies to every conditional object in this file):
# conditional entries are written as `key = cond ? value : null` inside ONE
# object literal, pruned with `{ for k, v in {...} : k => v if v != null }`.
# `cond ? {...} : {}` ternaries fail plan-time type unification when
# branches carry different attributes, and merge() over primitive-only
# sibling objects silently UNIFIES them into map(string). Optional nested
# blocks are read with try() (HCL's && does NOT short-circuit); optional
# scalars inside optional blocks with try(coalesce(x), null).

locals {
  # Chart identity — must stay byte-identical with the Pulumi module's
  # vars. The argo-workflows chart is served from the argoproj
  # GitHub-pages index; chart 1.0.23 ships Argo Workflows v4.0.8.
  helm_chart_name = "argo-workflows"
  helm_chart_repo = "https://argoproj.github.io/argo-helm"

  # Release name — metadata.name, NOT a fixed chart name: several engines
  # can coexist (pair with controller.instance_id so each reconciles only
  # its own Workflows). fullnameOverride below pins every chart child
  # name to this.
  release_name = var.metadata.name

  # Chart version resolved to the pinned default when unset, so both
  # engines install the same chart whether or not the platform's
  # defaulting middleware ran.
  chart_version = coalesce(var.spec.chart_version, "1.0.23")

  namespace = var.spec.namespace

  # Resource-identity labels stamped on the module-created satellites
  # (the namespace — never injected into the chart's own resources; Helm
  # owns those).
  labels = merge(
    {
      "planton.ai/resource"      = "true"
      "planton.ai/resource-name" = var.metadata.name
      "planton.ai/resource-kind" = "KubernetesArgoWorkflows"
    },
    var.metadata.id != null && var.metadata.id != "" ? { "planton.ai/resource-id" = var.metadata.id } : {},
    var.metadata.org != null && var.metadata.org != "" ? { "planton.ai/organization" = var.metadata.org } : {},
    var.metadata.env != null && var.metadata.env != "" ? { "planton.ai/environment" = var.metadata.env } : {}
  )

  # The Argo server Service is `<fullname>-server` (the chart appends
  # each component's name to the fullname, which fullnameOverride pins to
  # the resource name). Port 2746. Feeds the exported handles.
  server_enabled      = try(coalesce(var.spec.server.enabled), null) != null ? var.spec.server.enabled : true
  server_service_name = local.server_enabled ? "${local.release_name}-server" : ""

  # The server speaks plain HTTP unless server.secure turns on its
  # self-signed TLS listener; the exported endpoint follows the scheme.
  server_secure = try(var.spec.server.secure, false)
  server_scheme = local.server_secure ? "https" : "http"

  # Name of the ServiceAccount workflow pods run as — the identity to
  # annotate for IRSA/workload identity when workflows touch cloud APIs.
  workflow_service_account = try(coalesce(var.spec.workflow_service_account), null) != null ? var.spec.workflow_service_account : "argo-workflow"

  # ---- shared resources rendering -------------------------------------------
  resources_or_null = {
    controller = try(var.spec.controller.resources, null)
    server     = try(var.spec.server.resources, null)
  }

  rendered_resources = {
    for name, r in local.resources_or_null : name => r == null ? null : {
      for rk, rv in {
        requests = try(r.requests, null) == null ? null : {
          for qk, qv in {
            cpu    = r.requests.cpu != "" ? r.requests.cpu : null
            memory = r.requests.memory != "" ? r.requests.memory : null
          } : qk => qv if qv != null
        }
        limits = try(r.limits, null) == null ? null : {
          for lk, lv in {
            cpu    = r.limits.cpu != "" ? r.limits.cpu : null
            memory = r.limits.memory != "" ? r.limits.memory : null
          } : lk => lv if lv != null
        }
      } : rk => rv if rv != null && rv != {}
    }
  }

  # ---- scheduling (rendered per component — the chart has no global block) --
  # A single null-pruned object, NEVER an outer `declared ? {...} : {}`
  # ternary (differently-shaped branches fail plan-time type
  # unification); with nothing declared the comprehension yields {}.
  scheduling_block = {
    for k, v in {
      nodeSelector = length(try(var.spec.scheduling.node_selector, {})) > 0 ? var.spec.scheduling.node_selector : null
      tolerations = length(try(var.spec.scheduling.tolerations, [])) > 0 ? [
        for t in var.spec.scheduling.tolerations : {
          for tk, tv in {
            key               = t.key != "" ? t.key : null
            operator          = t.operator != "" ? t.operator : null
            value             = t.value != "" ? t.value : null
            effect            = t.effect != "" ? t.effect : null
            tolerationSeconds = try(t.toleration_seconds, null)
          } : tk => tv if tv != null
        }
      ] : null
      priorityClassName = try(var.spec.scheduling.priority_class_name, "") != "" ? var.spec.scheduling.priority_class_name : null
    } : k => v if v != null
  }

  # ---- workflow archive (controller.persistence) ------------------------------
  # The chart takes the engine section keyed by name (postgresql | mysql)
  # with secret SELECTORS for the credentials — resolved by the
  # controller at runtime, never rendered as values.
  archive_declared = try(var.spec.archive, null) != null

  archive_engine_key = local.archive_declared ? (var.spec.archive.engine == "mysql" ? "mysql" : "postgresql") : ""

  archive_username_key = local.archive_declared ? (
    try(coalesce(var.spec.archive.credentials_secret.username_key), null) != null ? var.spec.archive.credentials_secret.username_key : "username"
  ) : ""
  archive_password_key = local.archive_declared ? (
    try(coalesce(var.spec.archive.credentials_secret.password_key), null) != null ? var.spec.archive.credentials_secret.password_key : "password"
  ) : ""

  archive_engine_block = local.archive_declared ? {
    for k, v in {
      host      = var.spec.archive.host
      port      = try(coalesce(var.spec.archive.port), null) != null ? var.spec.archive.port : (var.spec.archive.engine == "mysql" ? 3306 : 5432)
      database  = var.spec.archive.database
      tableName = "argo_workflows"
      userNameSecret = {
        name = var.spec.archive.credentials_secret.name
        key  = local.archive_username_key
      }
      passwordSecret = {
        name = var.spec.archive.credentials_secret.name
        key  = local.archive_password_key
      }
      ssl     = var.spec.archive.ssl_mode != "" && var.spec.archive.ssl_mode != "disable" ? true : null
      sslMode = var.spec.archive.ssl_mode != "" ? var.spec.archive.ssl_mode : null
    } : k => v if v != null
  } : null

  persistence_block = local.archive_declared ? {
    archive                    = true
    (local.archive_engine_key) = local.archive_engine_block
  } : null

  # ---- controller block ------------------------------------------------------------
  # workflowNamespaces is ALWAYS rendered: the chart's own default is
  # ["default"], and its workflow-role template creates the runner
  # SA/Role/RoleBinding in every listed namespace PLUS the release
  # namespace — leaving the default in place makes every install leak
  # runner RBAC into the cluster's `default` namespace. An explicit []
  # keeps the runner identity release-namespace-only (lists REPLACE under
  # Helm -f merge semantics).
  #
  # instanceID is a STRUCTURED chart block ({enabled, useReleaseName,
  # explicitID} — templates read .enabled directly); the spec's plain
  # string maps to the enabled+explicitID shape, never a bare string.
  controller_block = merge(
    {
      for k, v in {
        replicas           = try(coalesce(var.spec.controller.replicas), null)
        resources          = local.rendered_resources.controller
        workflowNamespaces = length(try(var.spec.controller.workflow_namespaces, [])) > 0 ? var.spec.controller.workflow_namespaces : []
        instanceID = try(var.spec.controller.instance_id, "") != "" ? {
          enabled    = true
          explicitID = var.spec.controller.instance_id
        } : null
        parallelism = try(coalesce(var.spec.controller.parallelism), null)
        namespaceParallelism = try(coalesce(var.spec.controller.namespace_parallelism), null)
        persistence          = local.persistence_block
        retentionPolicy = try(var.spec.retention_policy, null) != null ? {
          for rk, rv in {
            completed = try(coalesce(var.spec.retention_policy.completed), null)
            failed    = try(coalesce(var.spec.retention_policy.failed), null)
            errored   = try(coalesce(var.spec.retention_policy.errored), null)
          } : rk => rv if rv != null
        } : null
        serviceMonitor = var.spec.service_monitor_enabled ? { enabled = true } : null
      } : k => v if v != null && v != {}
    },
    local.scheduling_block
  )

  # ---- server block -----------------------------------------------------------------
  # Scheduling keys are folded into the SAME comprehension gated per key
  # (never `enabled ? scheduling_block : {}` — the type-unification
  # class); a disabled server renders only enabled=false.
  server_block = {
    for k, v in {
      enabled           = local.server_enabled ? null : false
      replicas          = local.server_enabled ? try(coalesce(var.spec.server.replicas), null) : null
      resources         = local.server_enabled ? local.rendered_resources.server : null
      authModes         = local.server_enabled && length(try(var.spec.server.auth_modes, [])) > 0 ? var.spec.server.auth_modes : null
      secure            = local.server_enabled && local.server_secure ? true : null
      baseHref          = local.server_enabled && try(var.spec.server.base_href, "") != "" ? var.spec.server.base_href : null
      nodeSelector      = local.server_enabled ? try(local.scheduling_block.nodeSelector, null) : null
      tolerations       = local.server_enabled ? try(local.scheduling_block.tolerations, null) : null
      priorityClassName = local.server_enabled ? try(local.scheduling_block.priorityClassName, null) : null
    } : k => v if v != null && v != {}
  }

  # ---- artifact repository ---------------------------------------------------------
  # Exactly one backend renders (proto oneof). Secret material rides the
  # chart's secret selectors; keyless arms lean on ambient pod identity
  # (useSDKCreds) on the runner service account.
  s3_declared    = try(var.spec.artifact_repository.s3, null) != null
  gcs_declared   = try(var.spec.artifact_repository.gcs, null) != null
  azure_declared = try(var.spec.artifact_repository.azure, null) != null
  artifact_repository_declared = (
    local.s3_declared || local.gcs_declared || local.azure_declared
  )

  # The s3 credentials Secret: its name resolves a KubernetesSeaweedFs
  # credentials Secret by reference (the generator flattens the
  # StringValueOrRef to a plain string); the key names default to that
  # kind's generated `-s3-secret` convention (admin pair — mirror of the
  # proto field options), so that Secret composes with zero key
  # configuration. The key-name fields are proto `optional` scalars and
  # arrive as null in tfvars — hence the coalesce-guarded reads.
  s3_credentials_declared = try(var.spec.artifact_repository.s3.credentials_secret, null) != null
  s3_access_key_id_key = local.s3_credentials_declared ? (
    try(coalesce(var.spec.artifact_repository.s3.credentials_secret.access_key_id_key), "") != "" ? var.spec.artifact_repository.s3.credentials_secret.access_key_id_key : "admin_access_key_id"
  ) : null
  s3_secret_access_key_key = local.s3_credentials_declared ? (
    try(coalesce(var.spec.artifact_repository.s3.credentials_secret.secret_access_key_key), "") != "" ? var.spec.artifact_repository.s3.credentials_secret.secret_access_key_key : "admin_secret_access_key"
  ) : null

  s3_block = local.s3_declared ? {
    for k, v in {
      bucket   = var.spec.artifact_repository.s3.bucket
      endpoint = try(var.spec.artifact_repository.s3.endpoint, "") != "" ? var.spec.artifact_repository.s3.endpoint : null
      region   = try(var.spec.artifact_repository.s3.region, "") != "" ? var.spec.artifact_repository.s3.region : null
      insecure = try(var.spec.artifact_repository.s3.insecure, false) ? true : null
      # Keyless: sign with the pod's ambient identity (IRSA / workload
      # identity on the runner service account).
      useSDKCreds = try(var.spec.artifact_repository.s3.use_ambient_credentials, false) ? true : null
      accessKeySecret = local.s3_credentials_declared ? {
        name = var.spec.artifact_repository.s3.credentials_secret.secret_name
        key  = local.s3_access_key_id_key
      } : null
      secretKeySecret = local.s3_credentials_declared ? {
        name = var.spec.artifact_repository.s3.credentials_secret.secret_name
        key  = local.s3_secret_access_key_key
      } : null
    } : k => v if v != null
  } : null

  gcs_block = local.gcs_declared ? {
    for k, v in {
      bucket = var.spec.artifact_repository.gcs.bucket
      serviceAccountKeySecret = try(var.spec.artifact_repository.gcs.credentials_secret_name, "") != "" ? {
        name = var.spec.artifact_repository.gcs.credentials_secret_name
        key  = "serviceAccountKey"
      } : null
    } : k => v if v != null
  } : null

  azure_block = local.azure_declared ? {
    for k, v in {
      endpoint  = var.spec.artifact_repository.azure.endpoint
      container = var.spec.artifact_repository.azure.container
      # Keyless: managed identity / workload identity when no account key
      # is declared.
      useSDKCreds = try(var.spec.artifact_repository.azure.credentials_secret_name, "") == "" ? true : null
      accountKeySecret = try(var.spec.artifact_repository.azure.credentials_secret_name, "") != "" ? {
        name = var.spec.artifact_repository.azure.credentials_secret_name
        key  = "account-access-key"
      } : null
    } : k => v if v != null
  } : null

  artifact_repository_block = local.artifact_repository_declared ? {
    for k, v in {
      archiveLogs = try(var.spec.artifact_repository.archive_logs, false) ? true : null
      s3          = local.s3_block
      gcs         = local.gcs_block
      azure       = local.azure_block
    } : k => v if v != null
  } : null

  # ---- typed chart values (twin of the Pulumi module's buildHelmValues) --
  helm_values = {
    for k, v in {
      # fullnameOverride pins every child name (`<name>-server`,
      # `<name>-workflow-controller`, ...) to the resource name; the
      # exported outputs are built from that contract.
      fullnameOverride = local.release_name

      # CRD lifecycle: full-schema CRDs arrive via the chart's hook Job,
      # which DOWNLOADS them from the chart's GitHub release at install
      # time (internet-at-install); crds.full=false falls back to
      # templated minified CRDs for air-gapped clusters.
      crds = try(var.spec.crds, null) != null ? {
        for ck, cv in {
          install = try(coalesce(var.spec.crds.install), null) != null ? var.spec.crds.install : true
          keep    = try(coalesce(var.spec.crds.keep), null) != null ? var.spec.crds.keep : true
          full    = try(coalesce(var.spec.crds.full_schema), null) != null ? var.spec.crds.full_schema : true
          upgradeJob = try(var.spec.crds.base_url, "") != "" ? {
            crdBaseURL = var.spec.crds.base_url
          } : null
        } : ck => cv if cv != null
      } : null

      # The registry/tag override maps onto the chart's SPLIT image
      # values: images.tag is the shared tag; each component's
      # image.registry moves to the mirror while the upstream repository
      # paths stay (the split-image discipline — a combined mapping
      # would break every mirror override identically in both engines).
      images = try(var.spec.image.tag, "") != "" || try(var.spec.image.pull_secret_name, "") != "" ? {
        for ik, iv in {
          tag         = try(var.spec.image.tag, "") != "" ? var.spec.image.tag : null
          pullSecrets = try(var.spec.image.pull_secret_name, "") != "" ? [{ name = var.spec.image.pull_secret_name }] : null
        } : ik => iv if iv != null
      } : null

      controller = merge(
        local.controller_block,
        try(var.spec.image.registry, "") != "" ? { image = { registry = var.spec.image.registry } } : {}
      )

      server = merge(
        local.server_block,
        local.server_enabled && try(var.spec.image.registry, "") != "" ? { image = { registry = var.spec.image.registry } } : {}
      )

      executor = try(var.spec.image.registry, "") != "" ? {
        image = { registry = var.spec.image.registry }
      } : null

      # The chart DEFAULTS workflow.serviceAccount.create to FALSE — an
      # engine whose runner ServiceAccount does not exist rejects every
      # submission. The module always creates it (with the configured
      # name); the chart's workflow.rbac.create default (true) then binds
      # the runner Role to it in every watched namespace.
      workflow = {
        serviceAccount = {
          create = true
          name   = local.workflow_service_account
        }
      }

      artifactRepository = local.artifact_repository_block
    } : k => v if v != null && v != {}
  }
}
