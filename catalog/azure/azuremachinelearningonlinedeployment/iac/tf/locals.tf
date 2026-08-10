locals {
  resource_id = (
    var.metadata.id != null && var.metadata.id != ""
    ? var.metadata.id
    : var.metadata.name
  )

  # PARITY-EXCEPTION: resource_kind here is the family-wide snake-case
  # literal and resource_id falls back to metadata.name, while the Pulumi
  # module emits the lowered CloudResourceKind enum string and omits
  # resource_id when metadata.id is empty. Output-neutral (tags never feed
  # stack outputs); aligning the two shapes is a family-wide convention
  # change, not a per-kind fix.
  base_tags = {
    "resource"      = "true"
    "resource_id"   = local.resource_id
    "resource_kind" = "azure_machine_learning_online_deployment"
    "resource_name" = var.metadata.name
  }

  org_tag = (
    var.metadata.org != null && var.metadata.org != ""
  ) ? { "organization" = var.metadata.org } : {}

  env_tag = (
    var.metadata.env != null && var.metadata.env != ""
  ) ? { "environment" = var.metadata.env } : {}

  # Metadata-derived tags first, then the user's spec tags merged over them:
  # user tags deliberately win so an org's governance conventions (cost
  # center, owner) can override the derived values where they collide.
  final_tags = merge(local.base_tags, local.org_tag, local.env_tag, var.spec.tags)

  # The three health probes share one ARM shape; each is assembled
  # key-by-key so unset optionals are omitted and the service's own
  # defaults apply. Keys here are the ARM property names the probes merge
  # into the deployment's properties under.
  probe_specs = {
    livenessProbe  = var.spec.liveness_probe
    readinessProbe = var.spec.readiness_probe
    startupProbe   = var.spec.startup_probe
  }

  probe_bodies = {
    for arm_name, probe in local.probe_specs : arm_name => merge(
      probe.failure_threshold != null ? { failureThreshold = probe.failure_threshold } : {},
      probe.success_threshold != null ? { successThreshold = probe.success_threshold } : {},
      probe.initial_delay != "" ? { initialDelay = probe.initial_delay } : {},
      probe.period != "" ? { period = probe.period } : {},
      probe.timeout != "" ? { timeout = probe.timeout } : {},
    ) if probe != null
  }

  code_configuration = var.spec.code_configuration != null ? {
    codeConfiguration = merge(
      var.spec.code_configuration.code_id != "" ? { codeId = var.spec.code_configuration.code_id } : {},
      { scoringScript = var.spec.code_configuration.scoring_script },
    )
  } : {}

  request_settings = var.spec.request_settings != null ? {
    requestSettings = merge(
      var.spec.request_settings.max_concurrent_requests_per_instance != null ? {
        maxConcurrentRequestsPerInstance = var.spec.request_settings.max_concurrent_requests_per_instance
      } : {},
      var.spec.request_settings.request_timeout != "" ? {
        requestTimeout = var.spec.request_settings.request_timeout
      } : {},
    )
  } : {}

  # Model data collection. Each collection's two-value
  # Enabled/Disabled surface rides the spec's bool.
  data_collector = var.spec.data_collector != null ? {
    dataCollector = merge(
      {
        collections = {
          for name, collection in var.spec.data_collector.collections : name => merge(
            { dataCollectionMode = collection.enabled ? "Enabled" : "Disabled" },
            collection.data_id != "" ? { dataId = collection.data_id } : {},
            collection.client_id != "" ? { clientId = collection.client_id } : {},
            collection.sampling_rate != null ? { samplingRate = collection.sampling_rate } : {},
          )
        }
      },
      var.spec.data_collector.rolling_rate != "" ? { rollingRate = var.spec.data_collector.rolling_rate } : {},
      var.spec.data_collector.request_logging != null && length(try(var.spec.data_collector.request_logging.capture_headers, [])) > 0 ? {
        requestLogging = { captureHeaders = var.spec.data_collector.request_logging.capture_headers }
      } : {},
    )
  } : {}

  # The ARM properties object. endpointComputeType is PINNED to Managed:
  # this kind models the managed compute type only (the Kubernetes and
  # AzureMLCompute variants are a recorded deferral -- they require
  # attached compute whose supported story does not exist yet). Scale
  # settings are deliberately absent: the managed variant's only legal
  # mode is Default (fixed instance count via the SKU capacity below);
  # TargetUtilization is Kubernetes-only.
  deployment_properties = merge(
    {
      endpointComputeType = "Managed"
      appInsightsEnabled  = var.spec.app_insights_enabled
    },
    var.spec.instance_type != "" ? { instanceType = var.spec.instance_type } : {},
    var.spec.model != "" ? { model = var.spec.model } : {},
    var.spec.model_mount_path != "" ? { modelMountPath = var.spec.model_mount_path } : {},
    var.spec.environment_id != "" ? { environmentId = var.spec.environment_id } : {},
    length(var.spec.environment_variables) > 0 ? { environmentVariables = var.spec.environment_variables } : {},
    var.spec.egress_public_network_access_enabled != null ? {
      egressPublicNetworkAccess = var.spec.egress_public_network_access_enabled ? "Enabled" : "Disabled"
    } : {},
    var.spec.description != "" ? { description = var.spec.description } : {},
    length(var.spec.properties) > 0 ? { properties = var.spec.properties } : {},
    local.probe_bodies,
    local.code_configuration,
    local.request_settings,
    local.data_collector,
  )
}
