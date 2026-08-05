# Computed values for the KubernetesVelero module. Every resolution here
# has an exact twin in the Pulumi module's locals.go / values.go — keep
# them in lockstep.
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
#
# ARM SELECTION DISCIPLINE: the backend oneof is spec-enforced to exactly
# one arm, so arm-dependent scalars are picked with
# `one([for x in [...] : x if x != null])` / per-provider lookup maps —
# never chained per-arm ternaries.

locals {
  # Chart identity — must stay byte-identical with the Pulumi module's
  # vars: cross-engine chart-name drift deploys two different products
  # from one manifest.
  helm_chart_name = "velero"
  helm_chart_repo = "https://vmware-tanzu.github.io/helm-charts"

  # Release name fixed to the chart name: Velero's CRDs and node-agent
  # are cluster-scoped and one server owns the backup records in the
  # store — one installation per cluster is an upstream constraint. The
  # fixed name also collapses the chart's velero.fullname to "velero"
  # (the release name contains the chart name), making every derived
  # name deterministic.
  release_name = local.helm_chart_name

  # Chart version resolved to the pinned default when unset, so both
  # engines install the same chart whether or not the platform's
  # defaulting middleware ran — mirror of the Pulumi module's
  # DefaultChartVersion (chart 12.1.0 ships Velero 1.18.1).
  chart_version = coalesce(var.spec.chart_version, "12.1.0")

  namespace = var.spec.namespace

  # Chart-derived ServiceAccount name (templates/_helpers.tpl
  # "velero.serverServiceAccount"): serviceAccount.server.create defaults
  # true and the module never sets a name, so the name is
  # "<velero.fullname>-server" — with the release fixed to "velero" that
  # is "velero-server". Exported for cloud-side keyless identity
  # bindings (IRSA trust policies, GCP WI bindings, Azure federated
  # credentials).
  service_account_name = "velero-server"

  # Name the module gives the default BackupStorageLocation — what
  # Backup/Schedule resources reference through storageLocation.
  backup_storage_location_name = "default"

  # Resource-identity labels stamped on the namespace this module creates
  # (never injected into the chart's own resources — Helm owns those).
  labels = merge(
    {
      "planton.ai/resource"      = "true"
      "planton.ai/resource-name" = var.metadata.name
      "planton.ai/resource-kind" = "KubernetesVelero"
    },
    var.metadata.id != null && var.metadata.id != "" ? { "planton.ai/resource-id" = var.metadata.id } : {},
    var.metadata.org != null && var.metadata.org != "" ? { "planton.ai/organization" = var.metadata.org } : {},
    var.metadata.env != null && var.metadata.env != "" ? { "planton.ai/environment" = var.metadata.env } : {}
  )

  # ---- active backend arm -----------------------------------------------
  # Exactly one arm is set (spec-enforced oneof). Everything
  # arm-dependent below derives from these three handles plus the
  # provider-keyed lookup maps.
  s3    = try(var.spec.backup_storage.s3, null)
  gcs   = try(var.spec.backup_storage.gcs, null)
  azure = try(var.spec.backup_storage.azure_blob, null)

  # Velero provider name of the active arm (also the BSL/VSL provider).
  provider_name = one(concat(
    local.s3 != null ? ["aws"] : [],
    local.gcs != null ? ["gcp"] : [],
    local.azure != null ? ["azure"] : []
  ))

  # Official plugin images at the versions paired with the chart's
  # Velero release (chart 12.1.0 = Velero 1.18); overridable via
  # backup_storage.plugin_image for private-registry mirrors or pins.
  plugin_default_images = {
    aws   = "velero/velero-plugin-for-aws:v1.14.2"
    gcp   = "velero/velero-plugin-for-gcp:v1.14.2"
    azure = "velero/velero-plugin-for-microsoft-azure:v1.14.2"
  }
  # Init-container names (upstream convention from the chart README's
  # install example: velero-plugin-for-<PROVIDER>).
  plugin_names = {
    aws   = "velero-plugin-for-aws"
    gcp   = "velero-plugin-for-gcp"
    azure = "velero-plugin-for-azure"
  }

  # The provider plugin initContainer the chart expects verbatim under
  # `initContainers` (values.yaml documents the shape: image + a mount of
  # the shared `plugins` dir at /target, where the server discovers
  # plugin binaries).
  plugin_init_container = {
    name  = local.plugin_names[local.provider_name]
    image = var.spec.backup_storage.plugin_image != "" ? var.spec.backup_storage.plugin_image : local.plugin_default_images[local.provider_name]
    volumeMounts = [
      {
        mountPath = "/target"
        name      = "plugins"
      }
    ]
  }

  # ---- default BackupStorageLocation --------------------------------------
  # Velero's generic "bucket" is the blob CONTAINER on Azure.
  bsl_bucket = one(concat(
    local.s3 != null ? [local.s3.bucket] : [],
    local.gcs != null ? [local.gcs.bucket] : [],
    local.azure != null ? [local.azure.container] : []
  ))

  # Provider-specific BSL config. The chart template QUOTES every config
  # value ({{ $value | quote }}), so config entries are strings — the
  # path-style flag renders as the string "true".
  s3_bsl_config = local.s3 == null ? null : {
    for k, v in {
      region           = local.s3.region
      s3Url            = local.s3.s3_url != "" ? local.s3.s3_url : null
      s3ForcePathStyle = local.s3.force_path_style ? "true" : null
      kmsKeyId         = local.s3.kms_key_id != "" ? local.s3.kms_key_id : null
    } : k => v if v != null
  }
  # GCS: config.serviceAccount is the values.yaml-documented key for
  # workload identity ("Specify the service account here if you want to
  # use workload identity instead of providing the key file"); no config
  # at all otherwise.
  gcs_bsl_config = try(local.gcs.workload_identity_service_account_email, "") != "" ? {
    serviceAccount = local.gcs.workload_identity_service_account_email
  } : null
  azure_bsl_config = local.azure == null ? null : {
    resourceGroup  = local.azure.resource_group
    storageAccount = local.azure.storage_account
    subscriptionId = local.azure.subscription_id
  }
  bsl_config = try(
    one([for c in [local.s3_bsl_config, local.gcs_bsl_config, local.azure_bsl_config] : c if c != null]),
    null
  )

  # caCert is an ITEM-LEVEL BSL key (the chart template renders it under
  # the BSL's spec.objectStorage.caCert), NOT a config entry.
  backup_storage_location = {
    for k, v in {
      name     = local.backup_storage_location_name
      provider = local.provider_name
      bucket   = local.bsl_bucket
      default  = true
      prefix   = var.spec.backup_storage.prefix != "" ? var.spec.backup_storage.prefix : null
      caCert   = try(local.s3.ca_cert, "") != "" ? local.s3.ca_cert : null
      config   = local.bsl_config
    } : k => v if v != null
  }

  # ---- default VolumeSnapshotLocation ---------------------------------------
  # snapshotsEnabled matches the chart's own default (true) — rendered
  # only on explicit opt-out. Provider-required VSL config per the chart
  # values.yaml comments: region for aws; resourceGroup + subscriptionId
  # for azure; nothing for gcp.
  snapshots_enabled = try(var.spec.volume_snapshots.enabled, null) != null ? var.spec.volume_snapshots.enabled : true

  vsl_configs = {
    aws   = local.s3 == null ? null : { region = local.s3.region }
    gcp   = null
    azure = local.azure == null ? null : { resourceGroup = local.azure.resource_group, subscriptionId = local.azure.subscription_id }
  }

  volume_snapshot_location = {
    for k, v in {
      name     = local.backup_storage_location_name
      provider = local.provider_name
      config   = local.vsl_configs[local.provider_name]
    } : k => v if v != null
  }

  # ---- credential posture ----------------------------------------------------
  # The `cloud` credentials-file content per arm (chart values.yaml:
  # credentials.secretContents.cloud is "the entire content of your IAM
  # credentials file"; the AWS format is documented inline, GCP/Azure per
  # the plugin READMEs it links).
  #
  # s3 access keys → AWS shared-credentials file.
  s3_cloud = try(local.s3.access_keys, null) == null ? null : "[default]\naws_access_key_id=${local.s3.access_keys.access_key_id}\naws_secret_access_key=${local.s3.access_keys.secret_access_key}\n"

  # gcs key → the service-account JSON verbatim.
  gcs_cloud = try(local.gcs.service_account_key_json, "") != "" ? local.gcs.service_account_key_json : null

  # azure service principal → the full AZURE_* environment-file lines.
  azure_sp_cloud = try(local.azure.service_principal, null) == null ? null : "AZURE_SUBSCRIPTION_ID=${local.azure.subscription_id}\nAZURE_TENANT_ID=${local.azure.service_principal.tenant_id}\nAZURE_CLIENT_ID=${local.azure.service_principal.client_id}\nAZURE_CLIENT_SECRET=${local.azure.service_principal.client_secret}\nAZURE_RESOURCE_GROUP=${local.azure.resource_group}\nAZURE_CLOUD_NAME=AzurePublicCloud\n"

  # azure workload identity → unlike AWS/GCP the Azure plugin STILL reads
  # a `cloud` file for the non-credential parameters; the federated token
  # itself rides the client-id annotation + the azure.workload.identity/use
  # pod label below. A present client id selects the posture (the same
  # presence semantics as the S3 IRSA and GCS WI arms — the tfvars
  # converter flattens the spec's reference to the resolved string).
  azure_wi       = try(local.azure.workload_identity_client_id, "") != ""
  azure_wi_cloud = local.azure_wi ? "AZURE_SUBSCRIPTION_ID=${local.azure.subscription_id}\nAZURE_RESOURCE_GROUP=${local.azure.resource_group}\nAZURE_CLOUD_NAME=AzurePublicCloud\n" : null

  # At most one of these is non-null (spec-enforced credential XORs).
  credentials_cloud = try(
    one([for c in [local.s3_cloud, local.gcs_cloud, local.azure_sp_cloud, local.azure_wi_cloud] : c if c != null]),
    null
  )

  # The secret-bearing values document, isolated so main.tf can wrap it
  # with sensitive() — typed_values below stays fully visible in plans.
  credentials_values_doc = local.credentials_cloud == null ? null : yamlencode({
    credentials = {
      secretContents = {
        cloud = local.credentials_cloud
      }
    }
  })

  # Keyless-posture identity annotations on the chart-created server
  # ServiceAccount.
  server_service_account_annotations = merge(
    try(local.s3.irsa_role_arn, "") != "" ? {
      "eks.amazonaws.com/role-arn" = local.s3.irsa_role_arn
    } : {},
    try(local.gcs.workload_identity_service_account_email, "") != "" ? {
      "iam.gke.io/gcp-service-account" = local.gcs.workload_identity_service_account_email
    } : {},
    local.azure_wi ? {
      "azure.workload.identity/client-id" = local.azure.workload_identity_client_id
    } : {}
  )

  # ---- file-system backup ------------------------------------------------------
  node_agent_resources = try(var.spec.fs_backup.node_agent_resources, null) == null ? null : {
    for k, v in {
      limits = try(var.spec.fs_backup.node_agent_resources.limits, null) == null ? null : {
        for lk, lv in {
          cpu    = try(var.spec.fs_backup.node_agent_resources.limits.cpu, "") != "" ? var.spec.fs_backup.node_agent_resources.limits.cpu : null
          memory = try(var.spec.fs_backup.node_agent_resources.limits.memory, "") != "" ? var.spec.fs_backup.node_agent_resources.limits.memory : null
        } : lk => lv if lv != null
      }
      requests = try(var.spec.fs_backup.node_agent_resources.requests, null) == null ? null : {
        for rk, rv in {
          cpu    = try(var.spec.fs_backup.node_agent_resources.requests.cpu, "") != "" ? var.spec.fs_backup.node_agent_resources.requests.cpu : null
          memory = try(var.spec.fs_backup.node_agent_resources.requests.memory, "") != "" ? var.spec.fs_backup.node_agent_resources.requests.memory : null
        } : rk => rv if rv != null
      }
    } : k => v if v != null && length(v) > 0
  }

  node_agent_tolerations = length(try(var.spec.fs_backup.node_agent_tolerations, [])) > 0 ? [
    for t in var.spec.fs_backup.node_agent_tolerations : {
      for tk, tv in {
        key               = t.key != "" ? t.key : null
        operator          = t.operator != "" ? t.operator : null
        value             = t.value != "" ? t.value : null
        effect            = t.effect != "" ? t.effect : null
        tolerationSeconds = try(t.toleration_seconds, null)
      } : tk => tv if tv != null
    }
  ] : null

  node_agent_map = {
    for k, v in {
      resources   = local.node_agent_resources
      tolerations = local.node_agent_tolerations
    } : k => v if v != null
  }

  # ---- scheduled backups ----------------------------------------------------------
  # The chart's `schedules` value is a MAP keyed by schedule name (the
  # rendered Schedule object is named "velero-<key>" — velero.fullname
  # plus the key). Name uniqueness is spec-enforced. The three optional
  # template booleans render presence-aware: unset means "Velero
  # decides", which is different from false.
  schedule_templates = {
    for s in var.spec.schedules : s.name => {
      for k, v in {
        ttl                      = s.ttl != "" ? s.ttl : null
        includedNamespaces       = length(s.included_namespaces) > 0 ? s.included_namespaces : null
        excludedNamespaces       = length(s.excluded_namespaces) > 0 ? s.excluded_namespaces : null
        includedResources        = length(s.included_resources) > 0 ? s.included_resources : null
        excludedResources        = length(s.excluded_resources) > 0 ? s.excluded_resources : null
        labelSelector            = length(s.label_selector) > 0 ? { matchLabels = s.label_selector } : null
        includeClusterResources  = try(s.include_cluster_resources, null)
        snapshotVolumes          = try(s.snapshot_volumes, null)
        defaultVolumesToFsBackup = try(s.default_volumes_to_fs_backup, null)
        storageLocation          = s.storage_location != "" ? s.storage_location : null
      } : k => v if v != null
    }
  }

  schedules = length(var.spec.schedules) > 0 ? {
    for s in var.spec.schedules : s.name => {
      for k, v in {
        schedule = s.schedule
        paused   = s.paused ? true : null
        template = length(local.schedule_templates[s.name]) > 0 ? local.schedule_templates[s.name] : null
      } : k => v if v != null
    }
  } : null

  # ---- deployment sizing / scheduling -----------------------------------------------
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

  deployment_tolerations = length(try(var.spec.deployment.tolerations, [])) > 0 ? [
    for t in var.spec.deployment.tolerations : {
      for tk, tv in {
        key               = t.key != "" ? t.key : null
        operator          = t.operator != "" ? t.operator : null
        value             = t.value != "" ? t.value : null
        effect            = t.effect != "" ? t.effect : null
        tolerationSeconds = try(t.toleration_seconds, null)
      } : tk => tv if tv != null
    }
  ] : null

  # ---- own telemetry -----------------------------------------------------------------
  # metrics.enabled matches the chart's own default (true) — rendered
  # only on explicit opt-out. The ServiceMonitor is opt-in (it needs the
  # Prometheus operator CRDs on the cluster or the release fails).
  metrics_map = {
    for k, v in {
      enabled        = try(var.spec.prometheus.enabled, null) == false ? false : null
      serviceMonitor = try(var.spec.prometheus.service_monitor, false) ? { enabled = true } : null
    } : k => v if v != null
  }

  # ---- velero-server configuration block ----------------------------------------------
  # The chart's `configuration` collects the BSL/VSL lists AND the
  # velero-server CLI flags — several spec sections contribute.
  # defaultVolumesToFsBackup is a SERVER flag, so it lives here even
  # though the spec groups it with fs_backup.
  configuration = {
    for k, v in {
      backupStorageLocation       = [local.backup_storage_location]
      volumeSnapshotLocation      = local.snapshots_enabled ? [local.volume_snapshot_location] : null
      features                    = try(var.spec.volume_snapshots.enable_csi, false) ? "EnableCSI" : null
      defaultSnapshotMoveData     = try(var.spec.volume_snapshots.default_snapshot_move_data, false) ? true : null
      defaultVolumesToFsBackup    = try(var.spec.fs_backup.default_volumes_to_fs_backup, false) ? true : null
      defaultBackupTTL            = try(var.spec.server.default_backup_ttl, "") != "" ? var.spec.server.default_backup_ttl : null
      defaultItemOperationTimeout = try(var.spec.server.default_item_operation_timeout, "") != "" ? var.spec.server.default_item_operation_timeout : null
      garbageCollectionFrequency  = try(var.spec.server.garbage_collection_frequency, "") != "" ? var.spec.server.garbage_collection_frequency : null
      restoreOnlyMode             = try(var.spec.server.restore_only_mode, false) ? true : null
      logLevel                    = try(var.spec.server.log_level, "") != "" ? var.spec.server.log_level : null
      logFormat                   = try(var.spec.server.log_format, "") != "" ? var.spec.server.log_format : null
    } : k => v if v != null
  }

  # ---- typed chart values (twin of the Pulumi module's buildHelmValues) --
  typed_values = {
    for k, v in {
      # The chart's crds/-directory CRDs are Helm-native
      # keep-on-uninstall; upgradeCRDs (chart default true) is the job
      # that keeps them current across upgrades, and cleanUpCRDs (chart
      # default false) is the DESTRUCTIVE CI-oriented delete. Both render
      # only when they differ from the chart default.
      upgradeCRDs = try(var.spec.crds.upgrade_on_install, true) == false ? false : null
      cleanUpCRDs = try(var.spec.crds.cleanup_on_uninstall, false) ? true : null

      initContainers = [local.plugin_init_container]
      configuration  = local.configuration

      # useSecret here; the `cloud` content itself rides the separate
      # sensitive() document (main.tf) so plans stay readable.
      credentials = { useSecret = local.credentials_cloud != null }

      serviceAccount = length(local.server_service_account_annotations) > 0 ? {
        server = { annotations = local.server_service_account_annotations }
      } : null
      # The AKS webhook only injects the federated token into pods
      # carrying the azure.workload.identity/use label.
      podLabels = local.azure_wi ? { "azure.workload.identity/use" = "true" } : null

      snapshotsEnabled = local.snapshots_enabled ? null : false

      # deployNodeAgent matches the chart's own default (false) —
      # rendered only when the DaemonSet is wanted.
      deployNodeAgent = try(var.spec.fs_backup.deploy_node_agent, false) ? true : null
      nodeAgent       = length(local.node_agent_map) > 0 ? local.node_agent_map : null

      schedules = local.schedules

      # Top-level chart keys applying to the Velero server Deployment
      # (the node-agent has its own block above).
      resources         = local.deployment_resources
      priorityClassName = try(var.spec.deployment.priority_class_name, "") != "" ? var.spec.deployment.priority_class_name : null
      nodeSelector      = length(try(var.spec.deployment.node_selector, {})) > 0 ? var.spec.deployment.node_selector : null
      tolerations       = local.deployment_tolerations

      metrics = length(local.metrics_map) > 0 ? local.metrics_map : null
    } : k => v if v != null
  }
}
