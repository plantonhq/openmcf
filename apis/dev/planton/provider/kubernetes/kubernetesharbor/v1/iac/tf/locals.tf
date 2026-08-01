# Computed values. Every resolution here has an exact twin in the Pulumi
# module (locals.go / values.go) — keep them in lockstep.

locals {
  namespace    = var.spec.namespace
  release_name = var.metadata.name

  # Chart identity — the SERVED index truth (https://helm.goharbor.io);
  # chart 1.19.1 = Harbor 2.15.1.
  helm_chart_name       = "harbor"
  helm_chart_repo       = "https://helm.goharbor.io"
  default_chart_version = "1.19.1"
  chart_version         = try(coalesce(var.spec.chart_version), "") != "" ? var.spec.chart_version : local.default_chart_version

  # Planton governance labels for module-created satellites (namespace,
  # credential Secrets) — never injected into the chart's own
  # resources; Helm owns those.
  labels = merge(
    {
      "planton.ai/resource"      = "true"
      "planton.ai/resource-name" = var.metadata.name
      "planton.ai/resource-kind" = "KubernetesHarbor"
    },
    try(var.metadata.id, "") != "" ? { "planton.ai/resource-id" = var.metadata.id } : {},
    try(var.metadata.org, "") != "" ? { "planton.ai/organization" = var.metadata.org } : {},
    try(var.metadata.env, "") != "" ? { "planton.ai/environment" = var.metadata.env } : {},
  )

  # ------------------------------ expose ---------------------------------
  # The spec's Kubernetes-conventional Service type maps onto the
  # chart's camel-case expose.type values. Only clusterIP / nodePort /
  # loadBalancer ever render (never ingress or route — exposure
  # composes from the catalog's exposure kinds).
  expose_type = (
    try(var.spec.expose.service_type, "") == "NodePort" ? "nodePort" :
    try(var.spec.expose.service_type, "") == "LoadBalancer" ? "loadBalancer" : "clusterIP"
  )

  tls_enabled       = try(var.spec.expose.tls.enabled, false)
  tls_cert_secret   = local.tls_enabled ? try(coalesce(var.spec.expose.tls.cert_secret_name), "") : ""
  front_door_scheme = local.tls_enabled ? "https" : "http"
  front_door_port   = local.tls_enabled ? 443 : 80

  expose_tls_block = { for k, v in {
    enabled = local.tls_enabled
    # Chart-generated self-signed cert (auto) is regenerated on EVERY
    # apply — documented as labs-only on the spec field; the secret arm
    # is the cert-manager seam.
    certSource = local.tls_enabled ? (local.tls_cert_secret != "" ? "secret" : "auto") : null
    secret     = local.tls_enabled && local.tls_cert_secret != "" ? { secretName = local.tls_cert_secret } : null
    auto       = local.tls_enabled && local.tls_cert_secret == "" ? { commonName = local.release_name } : null
  } : k => v if v != null }

  # The active arm's Service NAME is pinned to metadata.name: the chart
  # default is the literal "harbor", which would collide between two
  # installs in one namespace.
  expose_node_ports = { for k, v in {
    http  = local.expose_type == "nodePort" && try(var.spec.expose.node_ports.http, null) != null ? { nodePort = var.spec.expose.node_ports.http } : null
    https = local.expose_type == "nodePort" && try(var.spec.expose.node_ports.https, null) != null ? { nodePort = var.spec.expose.node_ports.https } : null
  } : k => v if v != null }

  expose_arm = { for k, v in {
    name         = local.release_name
    annotations  = length(try(var.spec.expose.service_annotations, {})) > 0 ? var.spec.expose.service_annotations : null
    ports        = local.expose_type == "nodePort" && length(local.expose_node_ports) > 0 ? local.expose_node_ports : null
    IP           = local.expose_type == "loadBalancer" && try(coalesce(var.spec.expose.load_balancer_ip), "") != "" ? var.spec.expose.load_balancer_ip : null
    sourceRanges = local.expose_type == "loadBalancer" && length(try(var.spec.expose.source_ranges, [])) > 0 ? var.spec.expose.source_ranges : null
  } : k => v if v != null }

  expose_block = { for k, v in {
    type         = local.expose_type
    tls          = local.expose_tls_block
    clusterIP    = local.expose_type == "clusterIP" ? local.expose_arm : null
    nodePort     = local.expose_type == "nodePort" ? local.expose_arm : null
    loadBalancer = local.expose_type == "loadBalancer" ? local.expose_arm : null
  } : k => v if v != null }

  # --------------------------- credential names --------------------------
  admin_generated   = try(coalesce(var.spec.admin_auth.existing_secret_name), "") == ""
  admin_secret_name = local.admin_generated ? "${var.metadata.name}-admin-auth" : var.spec.admin_auth.existing_secret_name
  admin_secret_key  = local.admin_generated ? "HARBOR_ADMIN_PASSWORD" : (try(coalesce(var.spec.admin_auth.existing_secret_key), "") != "" ? var.spec.admin_auth.existing_secret_key : "HARBOR_ADMIN_PASSWORD")

  internal_auth_secret_name = "${var.metadata.name}-internal-auth"

  # The internal registry basic-auth username (not secret material —
  # the module-generated PASSWORD pair replaces the chart's public
  # default).
  registry_credential_username = "harbor_registry_user"

  # ------------------------------ database -------------------------------
  internal_database = try(var.spec.database.internal, null) != null
  database_external = try(var.spec.database.external, null)

  database_internal_block = local.internal_database ? { for k, v in {
    # The chart's ONLY intake for this credential is a value — the
    # random_password reference marks the whole rendered values
    # document sensitive (redacted in plans); the Pulumi twin wraps it
    # as a secret Output.
    password     = try(random_password.internal_database[0].result, null)
    resources    = local.component_resources_by_name["database_internal"]
    shmSizeLimit = try(coalesce(var.spec.database.internal.shm_size_limit), "") != "" ? var.spec.database.internal.shm_size_limit : null
    nodeSelector = local.sched_node_selector
    tolerations  = local.sched_tolerations
  } : k => v if v != null } : null

  database_external_block = local.database_external != null ? {
    host     = local.database_external.host
    port     = tostring(coalesce(try(local.database_external.port, null), 5432)) # the chart takes the port as a string
    username = local.database_external.username
    # The chart contract pins the key inside this Secret to `password`
    # — exactly the KubernetesPostgres application Secret's key, so its
    # credential composes as-is.
    existingSecret = local.database_external.password_secret_name
    coreDatabase   = try(coalesce(local.database_external.core_database), "") != "" ? local.database_external.core_database : "registry"
    sslmode        = try(coalesce(local.database_external.ssl_mode), "") != "" ? local.database_external.ssl_mode : "disable"
  } : null

  database_block = { for k, v in {
    type     = local.internal_database ? "internal" : "external"
    internal = local.database_internal_block
    external = local.database_external_block
  } : k => v if v != null }

  # -------------------------------- redis --------------------------------
  internal_redis = try(var.spec.cache.internal, null) != null
  redis_external = try(var.spec.cache.external, null)

  # Declared password → materialized by the module into
  # `<name>-redis-auth` under the chart's contract keys; only the NAME
  # renders in values.
  redis_declared_password = local.redis_external != null ? try(coalesce(local.redis_external.password), "") : ""
  redis_auth_secret_name  = local.redis_declared_password != "" ? "${var.metadata.name}-redis-auth" : ""

  redis_internal_block_raw = { for k, v in {
    resources    = local.component_resources_by_name["redis_internal"]
    nodeSelector = local.sched_node_selector
    tolerations  = local.sched_tolerations
  } : k => v if v != null }

  redis_external_block = local.redis_external != null ? { for k, v in {
    addr              = local.redis_external.addr
    sentinelMasterSet = try(coalesce(local.redis_external.sentinel_master_set), "") != "" ? local.redis_external.sentinel_master_set : null
    username          = try(coalesce(local.redis_external.username), "") != "" ? local.redis_external.username : null
    existingSecret = (
      try(coalesce(local.redis_external.existing_secret_name), "") != "" ? local.redis_external.existing_secret_name :
      local.redis_declared_password != "" ? local.redis_auth_secret_name : null
    )
    tlsOptions = try(local.redis_external.tls_enabled, false) ? { enable = true } : null
  } : k => v if v != null } : null

  redis_block = { for k, v in {
    type     = local.internal_redis ? "internal" : "external"
    internal = local.internal_redis && length(local.redis_internal_block_raw) > 0 ? local.redis_internal_block_raw : null
    external = local.redis_external_block
  } : k => v if v != null }

  # ------------------------------ storage --------------------------------
  storage_fs    = try(var.spec.storage.filesystem, null)
  storage_s3    = try(var.spec.storage.s3, null)
  storage_gcs   = try(var.spec.storage.gcs, null)
  storage_azure = try(var.spec.storage.azure, null)

  # Declared storage credentials → materialized into
  # `<name>-storage-auth` under each driver's contract keys.
  storage_declared_credentials = (
    (local.storage_s3 != null && try(coalesce(local.storage_s3.credentials.access_key), "") != "") ||
    (local.storage_gcs != null && try(coalesce(local.storage_gcs.key_data), "") != "") ||
    (local.storage_azure != null && try(coalesce(local.storage_azure.account_key), "") != "")
  )
  storage_auth_secret_name = local.storage_declared_credentials ? "${var.metadata.name}-storage-auth" : ""

  storage_s3_block = local.storage_s3 != null ? { for k, v in {
    bucket         = local.storage_s3.bucket
    region         = local.storage_s3.region
    regionendpoint = try(coalesce(local.storage_s3.endpoint), "") != "" ? local.storage_s3.endpoint : null
    existingSecret = (
      try(coalesce(local.storage_s3.credentials.existing_secret_name), "") != "" ? local.storage_s3.credentials.existing_secret_name :
      try(coalesce(local.storage_s3.credentials.access_key), "") != "" ? local.storage_auth_secret_name : null
    )
    encrypt       = try(local.storage_s3.encrypt, false) ? true : null
    secure        = coalesce(try(local.storage_s3.secure, null), true)
    skipverify    = try(local.storage_s3.skip_verify, false) ? true : null
    rootdirectory = try(coalesce(local.storage_s3.root_directory), "") != "" ? local.storage_s3.root_directory : null
    storageclass  = try(coalesce(local.storage_s3.storage_class), "") != "" ? local.storage_s3.storage_class : null
  } : k => v if v != null } : null

  storage_gcs_block = local.storage_gcs != null ? { for k, v in {
    bucket              = local.storage_gcs.bucket
    useWorkloadIdentity = try(local.storage_gcs.use_workload_identity, false) ? true : null
    existingSecret = (
      try(local.storage_gcs.use_workload_identity, false) ? null :
      try(coalesce(local.storage_gcs.existing_secret_name), "") != "" ? local.storage_gcs.existing_secret_name :
      try(coalesce(local.storage_gcs.key_data), "") != "" ? local.storage_auth_secret_name : null
    )
    rootdirectory = try(coalesce(local.storage_gcs.root_directory), "") != "" ? local.storage_gcs.root_directory : null
    chunksize     = try(local.storage_gcs.chunk_size, null) != null ? tostring(local.storage_gcs.chunk_size) : null # the chart takes it as a string
  } : k => v if v != null } : null

  storage_azure_block = local.storage_azure != null ? { for k, v in {
    accountname = local.storage_azure.account_name
    container   = local.storage_azure.container
    existingSecret = (
      try(coalesce(local.storage_azure.existing_secret_name), "") != "" ? local.storage_azure.existing_secret_name :
      try(coalesce(local.storage_azure.account_key), "") != "" ? local.storage_auth_secret_name : null
    )
    realm = try(coalesce(local.storage_azure.realm), "") != "" && try(local.storage_azure.realm, "") != "core.windows.net" ? local.storage_azure.realm : null
  } : k => v if v != null } : null

  image_chart_storage = { for k, v in {
    type = (
      local.storage_s3 != null ? "s3" :
      local.storage_gcs != null ? "gcs" :
      local.storage_azure != null ? "azure" : "filesystem"
    )
    filesystem = local.storage_fs != null ? { rootdirectory = "/storage" } : null
    s3         = local.storage_s3_block
    gcs        = local.storage_gcs_block
    azure      = local.storage_azure_block
    # Required for in-cluster S3 stores whose endpoint clients cannot
    # reach (SeaweedFS, MinIO) — a redirect would hand clients an
    # unreachable URL.
    disableredirect = local.storage_s3 != null && try(local.storage_s3.disable_redirect, false) ? true : null
  } : k => v if v != null }

  # ---------------------------- persistence ------------------------------
  keep_volumes    = coalesce(try(var.spec.keep_volumes_on_uninstall, null), true)
  resource_policy = local.keep_volumes ? "keep" : ""

  pvc_block = { for k, v in {
    registry = local.storage_fs != null ? { for k2, v2 in {
      size         = try(coalesce(local.storage_fs.disk_size), "") != "" ? local.storage_fs.disk_size : null
      storageClass = try(coalesce(local.storage_fs.storage_class), "") != "" ? local.storage_fs.storage_class : null
      accessMode   = try(coalesce(local.storage_fs.access_mode), "") != "" ? local.storage_fs.access_mode : null
    } : k2 => v2 if v2 != null } : null
    jobservice = try(coalesce(var.spec.jobservice.log_disk_size), "") != "" ? {
      jobLog = { size = var.spec.jobservice.log_disk_size }
    } : null
    database = local.internal_database ? { for k2, v2 in {
      size         = try(coalesce(var.spec.database.internal.disk_size), "") != "" ? var.spec.database.internal.disk_size : null
      storageClass = try(coalesce(var.spec.database.internal.storage_class), "") != "" ? var.spec.database.internal.storage_class : null
    } : k2 => v2 if v2 != null } : null
    redis = local.internal_redis ? { for k2, v2 in {
      size         = try(coalesce(var.spec.cache.internal.disk_size), "") != "" ? var.spec.cache.internal.disk_size : null
      storageClass = try(coalesce(var.spec.cache.internal.storage_class), "") != "" ? var.spec.cache.internal.storage_class : null
    } : k2 => v2 if v2 != null } : null
    trivy = local.trivy_enabled && try(coalesce(var.spec.trivy.disk_size), "") != "" ? {
      size = var.spec.trivy.disk_size
    } : null
  } : k => v if v != null && v != {} }

  persistence_block = { for k, v in {
    enabled               = true
    resourcePolicy        = local.resource_policy
    imageChartStorage     = local.image_chart_storage
    persistentVolumeClaim = length(local.pvc_block) > 0 ? local.pvc_block : null
  } : k => v if v != null }

  # ---------------------------- scheduling -------------------------------
  # Fanned onto EVERY component (split placement would separate
  # components that share volumes and fail over together);
  # per-component placement rides helm_values.
  sched_node_selector = length(try(var.spec.scheduling.node_selector, {})) > 0 ? var.spec.scheduling.node_selector : null
  sched_tolerations = length(try(var.spec.scheduling.tolerations, [])) > 0 ? [
    for t in var.spec.scheduling.tolerations : {
      for k, v in {
        key               = try(coalesce(t.key), "")
        operator          = try(coalesce(t.operator), "")
        value             = try(coalesce(t.value), "")
        effect            = try(coalesce(t.effect), "")
        tolerationSeconds = try(t.toleration_seconds, null)
      } : k => v if v != "" && v != null
    }
  ] : null

  # --------------------------- component sizing --------------------------
  component_resources_specs = {
    core              = try(var.spec.core.resources, null)
    portal            = try(var.spec.portal.resources, null)
    registry          = try(var.spec.registry.resources, null)
    jobservice        = try(var.spec.jobservice.resources, null)
    nginx             = try(var.spec.nginx.resources, null)
    trivy             = try(var.spec.trivy.resources, null)
    database_internal = try(var.spec.database.internal.resources, null)
    redis_internal    = try(var.spec.cache.internal.resources, null)
  }

  component_resources_by_name = { for name, r in local.component_resources_specs : name => (
    r != null ? {
      for k, v in {
        requests = try(r.requests, null) != null ? {
          for k2, v2 in {
            cpu    = try(coalesce(r.requests.cpu), "")
            memory = try(coalesce(r.requests.memory), "")
          } : k2 => v2 if v2 != ""
        } : null
        limits = try(r.limits, null) != null ? {
          for k2, v2 in {
            cpu    = try(coalesce(r.limits.cpu), "")
            memory = try(coalesce(r.limits.memory), "")
          } : k2 => v2 if v2 != ""
        } : null
      } : k => v if v != null && v != {}
    } : null
  ) }

  core_block = { for k, v in {
    replicas           = try(var.spec.core.replicas, null)
    resources          = local.component_resources_by_name["core"]
    existingSecret     = local.internal_auth_secret_name
    existingXsrfSecret = local.internal_auth_secret_name
    nodeSelector       = local.sched_node_selector
    tolerations        = local.sched_tolerations
  } : k => v if v != null }

  portal_block = { for k, v in {
    replicas     = try(var.spec.portal.replicas, null)
    resources    = local.component_resources_by_name["portal"]
    nodeSelector = local.sched_node_selector
    tolerations  = local.sched_tolerations
  } : k => v if v != null }

  registry_block = { for k, v in {
    replicas       = try(var.spec.registry.replicas, null)
    existingSecret = local.internal_auth_secret_name
    credentials = {
      username       = local.registry_credential_username
      existingSecret = local.internal_auth_secret_name
    }
    # Applied to the registry (distribution) container; the small
    # registryctl sidecar keeps chart defaults — size it via
    # helm_values if ever needed.
    registry     = local.component_resources_by_name["registry"] != null ? { resources = local.component_resources_by_name["registry"] } : null
    nodeSelector = local.sched_node_selector
    tolerations  = local.sched_tolerations
  } : k => v if v != null }

  jobservice_block = { for k, v in {
    replicas       = try(var.spec.jobservice.replicas, null)
    resources      = local.component_resources_by_name["jobservice"]
    existingSecret = local.internal_auth_secret_name
    maxJobWorkers  = try(var.spec.jobservice.max_job_workers, null)
    nodeSelector   = local.sched_node_selector
    tolerations    = local.sched_tolerations
  } : k => v if v != null }

  nginx_block = { for k, v in {
    replicas     = try(var.spec.nginx.replicas, null)
    resources    = local.component_resources_by_name["nginx"]
    nodeSelector = local.sched_node_selector
    tolerations  = local.sched_tolerations
  } : k => v if v != null }

  # -------------------------------- trivy --------------------------------
  # Unset block = chart truth: enabled.
  trivy_enabled = coalesce(try(var.spec.trivy.enabled, null), true)

  # SINGLE null-pruned object, never `enabled ? {rich} : {enabled=false}`
  # — the HCL type-unification class again: `enabled` is the boolean
  # itself; every other key is additionally gated on it and pruned.
  trivy_block = { for k, v in {
    enabled          = local.trivy_enabled
    replicas         = local.trivy_enabled ? try(var.spec.trivy.replicas, null) : null
    resources        = local.trivy_enabled ? local.component_resources_by_name["trivy"] : null
    skipUpdate       = local.trivy_enabled && try(var.spec.trivy.skip_update, false) ? true : null
    skipJavaDBUpdate = local.trivy_enabled && try(var.spec.trivy.skip_java_db_update, false) ? true : null
    offlineScan      = local.trivy_enabled && try(var.spec.trivy.offline_scan, false) ? true : null
    dbRepository     = local.trivy_enabled && length(try(var.spec.trivy.db_repositories, [])) > 0 ? var.spec.trivy.db_repositories : null
    javaDBRepository = local.trivy_enabled && length(try(var.spec.trivy.java_db_repositories, [])) > 0 ? var.spec.trivy.java_db_repositories : null
    # The chart accepts this only as a value (it renders it into its
    # own trivy Secret) — sensitive-wrapped so it is redacted in plans.
    gitHubToken   = local.trivy_enabled && try(coalesce(var.spec.trivy.github_token), "") != "" ? sensitive(var.spec.trivy.github_token) : null
    severity      = local.trivy_enabled && try(coalesce(var.spec.trivy.severity), "") != "" ? var.spec.trivy.severity : null
    ignoreUnfixed = local.trivy_enabled && try(var.spec.trivy.ignore_unfixed, false) ? true : null
    timeout       = local.trivy_enabled && try(coalesce(var.spec.trivy.timeout), "") != "" ? var.spec.trivy.timeout : null
    nodeSelector  = local.trivy_enabled ? local.sched_node_selector : null
    tolerations   = local.trivy_enabled ? local.sched_tolerations : null
  } : k => v if v != null }

  # ----------------------------- internal TLS ----------------------------
  internal_tls_enabled = try(var.spec.internal_tls.enabled, false)
  internal_tls_secrets = try(var.spec.internal_tls.cert_secrets, null)

  internal_tls_block = local.internal_tls_enabled ? { for k, v in {
    enabled            = true
    strong_ssl_ciphers = try(var.spec.internal_tls.strong_ssl_ciphers, false)
    # Chart-generated internal CA (auto) is regenerated on EVERY apply
    # (rolls every component) — documented as labs-only on the spec
    # field; the secret arm is the cert-manager seam.
    certSource = local.internal_tls_secrets != null ? "secret" : "auto"
    core       = local.internal_tls_secrets != null ? { secretName = local.internal_tls_secrets.core_secret_name } : null
    jobservice = local.internal_tls_secrets != null ? { secretName = local.internal_tls_secrets.jobservice_secret_name } : null
    registry   = local.internal_tls_secrets != null ? { secretName = local.internal_tls_secrets.registry_secret_name } : null
    portal     = local.internal_tls_secrets != null ? { secretName = local.internal_tls_secrets.portal_secret_name } : null
    trivy      = local.internal_tls_secrets != null && try(coalesce(local.internal_tls_secrets.trivy_secret_name), "") != "" ? { secretName = local.internal_tls_secrets.trivy_secret_name } : null
  } : k => v if v != null } : null

  # ------------------------------- metrics -------------------------------
  metrics_enabled = try(var.spec.metrics.enabled, false)

  metrics_block = local.metrics_enabled ? { for k, v in {
    enabled = true
    serviceMonitor = try(var.spec.metrics.service_monitor_enabled, false) ? { for k2, v2 in {
      enabled          = true
      interval         = try(coalesce(var.spec.metrics.service_monitor_interval), "") != "" ? var.spec.metrics.service_monitor_interval : null
      additionalLabels = length(try(var.spec.metrics.service_monitor_labels, {})) > 0 ? var.spec.metrics.service_monitor_labels : null
    } : k2 => v2 if v2 != null } : null
  } : k => v if v != null } : null

  exporter_block = local.metrics_enabled ? { for k, v in {
    nodeSelector = local.sched_node_selector
    tolerations  = local.sched_tolerations
  } : k => v if v != null } : null

  # ------------------------------- values --------------------------------
  typed_helm_values_raw = {
    fullnameOverride = local.release_name
    externalURL      = var.spec.external_url
    logLevel         = try(coalesce(var.spec.log_level), "") != "" ? var.spec.log_level : "info"

    expose = local.expose_block

    existingSecretAdminPassword    = local.admin_secret_name
    existingSecretAdminPasswordKey = local.admin_secret_key
    existingSecretSecretKey        = local.internal_auth_secret_name

    core       = local.core_block
    portal     = local.portal_block
    registry   = local.registry_block
    jobservice = local.jobservice_block
    nginx      = local.nginx_block

    database    = local.database_block
    redis       = local.redis_block
    persistence = local.persistence_block
    trivy       = local.trivy_block

    internalTLS = local.internal_tls_block
    metrics     = local.metrics_block
    # try(length(...)) — never `!= {}`: comparing a computed map against
    # an empty object literal is the type-unification edge, and HCL's
    # `&&` must not be trusted to short-circuit a null (WA-003).
    exporter = try(length(local.exporter_block), 0) > 0 ? local.exporter_block : null

    cache = try(var.spec.cache_layer.enabled, false) ? { for k, v in {
      enabled     = true
      expireHours = try(var.spec.cache_layer.expire_hours, null)
    } : k => v if v != null } : null

    proxy = try(coalesce(var.spec.outbound_proxy.http_proxy), "") != "" || try(coalesce(var.spec.outbound_proxy.https_proxy), "") != "" ? { for k, v in {
      httpProxy  = try(coalesce(var.spec.outbound_proxy.http_proxy), "") != "" ? var.spec.outbound_proxy.http_proxy : null
      httpsProxy = try(coalesce(var.spec.outbound_proxy.https_proxy), "") != "" ? var.spec.outbound_proxy.https_proxy : null
      noProxy    = try(coalesce(var.spec.outbound_proxy.no_proxy), "") != "" ? var.spec.outbound_proxy.no_proxy : null
    } : k => v if v != null } : null

    imagePullSecrets = length(try(var.spec.image_pull_secrets, [])) > 0 ? [
      for s in var.spec.image_pull_secrets : { name = s }
    ] : null

    updateStrategy = try(coalesce(var.spec.update_strategy), "") != "" && try(var.spec.update_strategy, "") != "RollingUpdate" ? {
      type = var.spec.update_strategy
    } : null

    caBundleSecretName = try(coalesce(var.spec.ca_bundle_secret_name), "") != "" ? var.spec.ca_bundle_secret_name : null
  }

  typed_helm_values = { for k, v in local.typed_helm_values_raw : k => v if v != null }

  # --------------------------- air-gap mirror ----------------------------
  # The chart pins every component image to
  # docker.io/goharbor/<component>:v<appVersion>; the mirror replaces
  # the docker.io host, keeping the goharbor/* paths and chart-default
  # tags (repository-only overrides — chart-default tags stay in
  # force). Rendered as a SEPARATE values document merged after the
  # typed rendering (twin of values.go applyImageMirror).
  # SINGLE null-pruned object, never `cond ? {rich} : {}`: two-arm
  # conditionals with different object shapes are the HCL
  # type-unification class this program keeps catching — every key is
  # individually gated on the mirror and pruned.
  image_mirror = try(coalesce(var.spec.image_registry), "")
  image_mirror_values = { for k, v in {
    nginx      = local.image_mirror != "" ? { image = { repository = "${local.image_mirror}/goharbor/nginx-photon" } } : null
    portal     = local.image_mirror != "" ? { image = { repository = "${local.image_mirror}/goharbor/harbor-portal" } } : null
    core       = local.image_mirror != "" ? { image = { repository = "${local.image_mirror}/goharbor/harbor-core" } } : null
    jobservice = local.image_mirror != "" ? { image = { repository = "${local.image_mirror}/goharbor/harbor-jobservice" } } : null
    registry = local.image_mirror != "" ? {
      registry   = { image = { repository = "${local.image_mirror}/goharbor/registry-photon" } }
      controller = { image = { repository = "${local.image_mirror}/goharbor/harbor-registryctl" } }
    } : null
    trivy    = local.image_mirror != "" && local.trivy_enabled ? { image = { repository = "${local.image_mirror}/goharbor/trivy-adapter-photon" } } : null
    exporter = local.image_mirror != "" && local.metrics_enabled ? { image = { repository = "${local.image_mirror}/goharbor/harbor-exporter" } } : null
    database = local.image_mirror != "" && local.internal_database ? { internal = { image = { repository = "${local.image_mirror}/goharbor/harbor-db" } } } : null
    redis    = local.image_mirror != "" && local.internal_redis ? { internal = { image = { repository = "${local.image_mirror}/goharbor/redis-photon" } } } : null
  } : k => v if v != null }

  # ------------------------------ outputs --------------------------------
  kube_endpoint = "${local.front_door_scheme}://${local.release_name}.${local.namespace}.svc.cluster.local:${local.front_door_port}"

  port_forward_local_port = local.tls_enabled ? 8443 : 8080
  port_forward_command    = "kubectl port-forward -n ${local.namespace} svc/${local.release_name} ${local.port_forward_local_port}:${local.front_door_port}"

  # NAME BUDGET (chart truth at 1.19.1): the chart truncates its
  # fullname at 63 and then APPENDS component suffixes — the longest,
  # `-jobservice-internal-tls` (24 chars), renders whenever internalTLS
  # runs in auto mode. Enforced by the release precondition in main.tf;
  # the Pulumi twin enforces the same budget.
  max_name_length = 39
}
