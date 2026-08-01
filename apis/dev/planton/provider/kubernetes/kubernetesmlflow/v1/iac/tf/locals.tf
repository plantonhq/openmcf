# KubernetesMlflow locals — every resolution here has an exact twin in the
# Pulumi module's locals.go; keep them in lockstep. This is a
# MODULE-OWNED-MANIFESTS kind (MLflow publishes no Helm chart; the
# official image is the distribution), so these locals ARE the deployment
# contract.

locals {
  name      = var.metadata.name
  namespace = var.spec.namespace

  # Resource-identity labels stamped on every module-created object.
  labels = merge(
    {
      "planton.ai/resource"      = "true"
      "planton.ai/resource-name" = var.metadata.name
      "planton.ai/resource-kind" = "KubernetesMlflow"
    },
    try(var.metadata.id, "") != "" ? { "planton.ai/resource-id" = var.metadata.id } : {},
    try(var.metadata.org, "") != "" ? { "planton.ai/organization" = var.metadata.org } : {},
    try(var.metadata.env, "") != "" ? { "planton.ai/environment" = var.metadata.env } : {}
  )

  # Selector labels for the Deployment/Service pairing — immutable once
  # deployed.
  selector_labels = {
    "app.kubernetes.io/name"     = "mlflow"
    "app.kubernetes.io/instance" = var.metadata.name
  }

  # ------------------------------ server --------------------------------
  # The OFFICIAL image (ghcr.io/mlflow/mlflow) at this kind's pinned
  # release — identical constants in the Pulumi module's vars.
  image_repository = try(var.spec.server.image.repository, "") != "" && try(var.spec.server.image.repository, "") != null ? var.spec.server.image.repository : "ghcr.io/mlflow/mlflow"
  image_tag        = try(var.spec.server.image.tag, "") != "" && try(var.spec.server.image.tag, "") != null ? var.spec.server.image.tag : "v3.15.0"
  image            = "${local.image_repository}:${local.image_tag}"

  server_port = 5000
  replicas    = try(var.spec.server.replicas, null) != null ? var.spec.server.replicas : 1
  workers     = try(var.spec.server.workers, null) != null ? var.spec.server.workers : 4

  # --------------------------- backend store ----------------------------
  backend_postgres = try(var.spec.backend_store.postgres, null)
  backend_mysql    = try(var.spec.backend_store.mysql, null)
  backend_sqlite   = try(var.spec.backend_store.sqlite_pvc, null)

  backend_type = local.backend_postgres != null ? "postgres" : (local.backend_mysql != null ? "mysql" : "sqlite")

  db_protocol = local.backend_type == "postgres" ? "postgresql" : "mysql+pymysql"
  db_host     = local.backend_type == "postgres" ? local.backend_postgres.host : (local.backend_type == "mysql" ? local.backend_mysql.host : "")
  db_port = local.backend_type == "postgres" ? (
    try(local.backend_postgres.port, null) != null ? local.backend_postgres.port : 5432
    ) : (local.backend_type == "mysql" ? (
      try(local.backend_mysql.port, null) != null ? local.backend_mysql.port : 3306
  ) : 0)
  db_name = local.backend_type == "postgres" ? (
    try(local.backend_postgres.database_name, "") != "" && try(local.backend_postgres.database_name, "") != null ? local.backend_postgres.database_name : "mlflow"
    ) : (local.backend_type == "mysql" ? (
      try(local.backend_mysql.database_name, "") != "" && try(local.backend_mysql.database_name, "") != null ? local.backend_mysql.database_name : "mlflow"
  ) : "")
  db_user = local.backend_type == "postgres" ? (
    try(local.backend_postgres.username, "") != "" && try(local.backend_postgres.username, "") != null ? local.backend_postgres.username : "mlflow"
    ) : (local.backend_type == "mysql" ? (
      try(local.backend_mysql.username, "") != "" && try(local.backend_mysql.username, "") != null ? local.backend_mysql.username : "mlflow"
  ) : "")
  db_password_secret = local.backend_type == "postgres" ? local.backend_postgres.password_secret.secret_name : (
    local.backend_type == "mysql" ? local.backend_mysql.password_secret.secret_name : ""
  )
  db_password_secret_key_raw = local.backend_type == "postgres" ? try(local.backend_postgres.password_secret.secret_key, "") : (
    local.backend_type == "mysql" ? try(local.backend_mysql.password_secret.secret_key, "") : ""
  )
  db_password_secret_key = local.db_password_secret_key_raw != "" && local.db_password_secret_key_raw != null ? local.db_password_secret_key_raw : "password"

  # Filesystem contract on the data PVC (sqlite arm): tracking db + auth
  # db side by side. Four slashes = sqlite absolute-path URI.
  data_mount_path        = "/mlflow/data"
  artifacts_mount_path   = "/mlflow/artifacts"
  sqlite_backend_uri     = "sqlite:///${local.data_mount_path}/mlflow.db"
  sqlite_auth_db_uri     = "sqlite:///${local.data_mount_path}/basic_auth.db"
  backend_uri_secret     = "${local.name}-backend-uri"
  data_pvc_enabled       = local.backend_type == "sqlite"
  data_pvc_name          = "${local.name}-data"
  data_pvc_size          = try(local.backend_sqlite.storage_size, "") != "" && try(local.backend_sqlite.storage_size, "") != null ? local.backend_sqlite.storage_size : "5Gi"
  data_pvc_storage_class = try(local.backend_sqlite.storage_class, "")

  # --------------------------- artifact store ---------------------------
  # The server PROXIES artifact traffic (--artifacts-destination) —
  # clients never carry store credentials.
  artifact_s3_compatible = try(var.spec.artifact_store.s3_compatible, null)
  artifact_aws_s3        = try(var.spec.artifact_store.aws_s3, null)
  artifact_gcs           = try(var.spec.artifact_store.gcs, null)
  artifact_azure         = try(var.spec.artifact_store.azure_blob, null)
  artifact_pvc           = try(var.spec.artifact_store.pvc, null)

  artifact_type = local.artifact_s3_compatible != null ? "s3_compatible" : (
    local.artifact_aws_s3 != null ? "aws_s3" : (
      local.artifact_gcs != null ? "gcs" : (
        local.artifact_azure != null ? "azure_blob" : "pvc"
      )
    )
  )

  artifact_destination_by_type = {
    s3_compatible = local.artifact_s3_compatible == null ? "" : (
      try(local.artifact_s3_compatible.prefix, "") != "" ? "s3://${local.artifact_s3_compatible.bucket}/${local.artifact_s3_compatible.prefix}" : "s3://${local.artifact_s3_compatible.bucket}"
    )
    aws_s3 = local.artifact_aws_s3 == null ? "" : (
      try(local.artifact_aws_s3.prefix, "") != "" ? "s3://${local.artifact_aws_s3.bucket}/${local.artifact_aws_s3.prefix}" : "s3://${local.artifact_aws_s3.bucket}"
    )
    gcs = local.artifact_gcs == null ? "" : (
      try(local.artifact_gcs.prefix, "") != "" ? "gs://${local.artifact_gcs.bucket}/${local.artifact_gcs.prefix}" : "gs://${local.artifact_gcs.bucket}"
    )
    azure_blob = local.artifact_azure == null ? "" : (
      try(local.artifact_azure.prefix, "") != "" ? "wasbs://${local.artifact_azure.container}@${local.artifact_azure.storage_account}.blob.core.windows.net/${local.artifact_azure.prefix}" : "wasbs://${local.artifact_azure.container}@${local.artifact_azure.storage_account}.blob.core.windows.net"
    )
    pvc = local.artifacts_mount_path
  }
  artifact_destination = local.artifact_destination_by_type[local.artifact_type]

  artifacts_pvc_enabled       = local.artifact_type == "pvc"
  artifacts_pvc_name          = "${local.name}-artifacts"
  artifacts_pvc_size          = try(local.artifact_pvc.storage_size, "") != "" && try(local.artifact_pvc.storage_size, "") != null ? local.artifact_pvc.storage_size : "10Gi"
  artifacts_pvc_storage_class = try(local.artifact_pvc.storage_class, "")

  s3_access_key_id_key     = try(local.artifact_s3_compatible.credentials_secret.access_key_id_key, "") != "" && try(local.artifact_s3_compatible.credentials_secret.access_key_id_key, "") != null ? local.artifact_s3_compatible.credentials_secret.access_key_id_key : "admin_access_key_id"
  s3_secret_access_key_key = try(local.artifact_s3_compatible.credentials_secret.secret_access_key_key, "") != "" && try(local.artifact_s3_compatible.credentials_secret.secret_access_key_key, "") != null ? local.artifact_s3_compatible.credentials_secret.secret_access_key_key : "admin_secret_access_key"

  aws_access_key_id_key     = try(local.artifact_aws_s3.credentials_secret.access_key_id_key, "") != "" && try(local.artifact_aws_s3.credentials_secret.access_key_id_key, "") != null ? local.artifact_aws_s3.credentials_secret.access_key_id_key : "access_key_id"
  aws_secret_access_key_key = try(local.artifact_aws_s3.credentials_secret.secret_access_key_key, "") != "" && try(local.artifact_aws_s3.credentials_secret.secret_access_key_key, "") != null ? local.artifact_aws_s3.credentials_secret.secret_access_key_key : "secret_access_key"

  gcs_credentials_secret_key = try(local.artifact_gcs.credentials_secret.secret_key, "") != "" && try(local.artifact_gcs.credentials_secret.secret_key, "") != null ? local.artifact_gcs.credentials_secret.secret_key : "credentials.json"
  gcs_credentials_mount_path = "/etc/mlflow/gcs"

  azure_credentials_secret_key = try(local.artifact_azure.credentials_secret.secret_key, "") != "" && try(local.artifact_azure.credentials_secret.secret_key, "") != null ? local.artifact_azure.credentials_secret.secret_key : "access_key"

  # ------------------------------- auth ----------------------------------
  # SECURED BY DEFAULT: basic auth is ON unless the spec disables it —
  # upstream's open server and its admin/password1234 example never
  # ship.
  auth_enabled   = try(var.spec.auth.enabled, null) != null ? var.spec.auth.enabled : true
  admin_username = try(var.spec.auth.admin_username, "") != "" && try(var.spec.auth.admin_username, "") != null ? var.spec.auth.admin_username : "admin"

  admin_secret_byo          = try(var.spec.auth.admin_password_secret, null) != null
  admin_secret_module_owned = local.auth_enabled && !local.admin_secret_byo
  admin_secret_name         = local.auth_enabled ? (local.admin_secret_byo ? var.spec.auth.admin_password_secret.secret_name : "${local.name}-admin-auth") : ""
  admin_secret_key_raw      = local.admin_secret_byo ? try(var.spec.auth.admin_password_secret.secret_key, "") : "password"
  admin_secret_key          = local.admin_secret_key_raw != "" && local.admin_secret_key_raw != null ? local.admin_secret_key_raw : "password"

  auth_config_secret_name = "${local.name}-auth-config"
  auth_config_mount_path  = "/etc/mlflow/auth"
  auth_config_file_name   = "basic_auth.ini"
  default_permission      = try(var.spec.auth.default_permission, "") != "" && try(var.spec.auth.default_permission, "") != null ? var.spec.auth.default_permission : "READ"

  # -------------------------------- gc -----------------------------------
  gc_enabled    = try(var.spec.gc.enabled, false)
  gc_schedule   = try(var.spec.gc.schedule, "") != "" && try(var.spec.gc.schedule, "") != null ? var.spec.gc.schedule : "0 3 * * *"
  gc_older_than = try(var.spec.gc.older_than, "") != "" && try(var.spec.gc.older_than, "") != null ? var.spec.gc.older_than : "30d"

  # ------------------------------ metrics --------------------------------
  metrics_enabled         = try(var.spec.metrics.enabled, false)
  service_monitor_enabled = try(var.spec.metrics.service_monitor_enabled, false)

  # ------------------------------ service --------------------------------
  service_type        = try(var.spec.service.type, "") != "" && try(var.spec.service.type, "") != null ? var.spec.service.type : "ClusterIP"
  service_annotations = try(var.spec.service.annotations, {})

  # ---------------------------- server args ------------------------------
  # `mlflow server` command line — shape follows upstream's own
  # deployment reference. The backend URI deliberately rides env (from
  # the module's Secret), never a pod argument.
  server_args = concat(
    [
      "mlflow", "server",
      "--host", "0.0.0.0",
      "--port", tostring(local.server_port),
      "--workers", tostring(local.workers),
      "--artifacts-destination", local.artifact_destination,
      "--serve-artifacts",
    ],
    local.auth_enabled ? ["--app-name", "basic-auth"] : [],
    local.metrics_enabled ? ["--expose-prometheus", "/tmp/metrics"] : [],
    try(var.spec.extra_args, [])
  )

  # ------------------------------ env vars -------------------------------
  # Rendered as dynamic env blocks in main.tf; credential values are all
  # Secret-sourced (secretKeyRef) — nothing credential-bearing lands in
  # the pod spec.
  plain_env = merge(
    local.backend_type == "sqlite" ? { MLFLOW_BACKEND_STORE_URI = local.sqlite_backend_uri } : {},
    local.auth_enabled ? { MLFLOW_AUTH_CONFIG_PATH = "${local.auth_config_mount_path}/${local.auth_config_file_name}" } : {},
    local.artifact_type == "s3_compatible" ? { MLFLOW_S3_ENDPOINT_URL = local.artifact_s3_compatible.endpoint } : {},
    local.artifact_type == "aws_s3" ? { AWS_DEFAULT_REGION = local.artifact_aws_s3.region } : {},
    local.artifact_type == "gcs" && try(local.artifact_gcs.credentials_secret, null) != null ? {
      GOOGLE_APPLICATION_CREDENTIALS = "${local.gcs_credentials_mount_path}/${local.gcs_credentials_secret_key}"
    } : {},
    try(var.spec.extra_env, {})
  )

  secret_env = merge(
    local.backend_type != "sqlite" ? {
      MLFLOW_BACKEND_STORE_URI = { secret_name = local.backend_uri_secret, secret_key = "uri" }
    } : {},
    local.artifact_type == "s3_compatible" ? {
      AWS_ACCESS_KEY_ID     = { secret_name = local.artifact_s3_compatible.credentials_secret.secret_name, secret_key = local.s3_access_key_id_key }
      AWS_SECRET_ACCESS_KEY = { secret_name = local.artifact_s3_compatible.credentials_secret.secret_name, secret_key = local.s3_secret_access_key_key }
    } : {},
    local.artifact_type == "aws_s3" && try(local.artifact_aws_s3.credentials_secret, null) != null ? {
      AWS_ACCESS_KEY_ID     = { secret_name = local.artifact_aws_s3.credentials_secret.secret_name, secret_key = local.aws_access_key_id_key }
      AWS_SECRET_ACCESS_KEY = { secret_name = local.artifact_aws_s3.credentials_secret.secret_name, secret_key = local.aws_secret_access_key_key }
    } : {},
    local.artifact_type == "azure_blob" ? {
      AZURE_STORAGE_ACCESS_KEY = { secret_name = local.artifact_azure.credentials_secret.secret_name, secret_key = local.azure_credentials_secret_key }
    } : {},
    {
      for name, ref in try(var.spec.extra_env_from_secret, {}) : name => {
        secret_name = ref.secret_name
        secret_key  = try(ref.secret_key, "") != "" && try(ref.secret_key, "") != null ? ref.secret_key : "password"
      }
    }
  )

  # Strategy follows the volume truth: any RWO PVC binds one pod, so
  # updates must Recreate; stateless shapes roll.
  deployment_strategy = (local.data_pvc_enabled || local.artifacts_pvc_enabled) ? "Recreate" : "RollingUpdate"

  # ------------------------------- outputs -------------------------------
  tracking_endpoint    = "http://${local.name}.${local.namespace}.svc.cluster.local:${local.server_port}"
  port_forward_command = "kubectl port-forward svc/${local.name} -n ${local.namespace} ${local.server_port}:${local.server_port}"

  admin_password_secret_output_name = local.auth_enabled ? local.admin_secret_name : ""
  admin_password_secret_output_key  = local.auth_enabled ? local.admin_secret_key : ""
  backend_uri_secret_output_name    = local.backend_type != "sqlite" ? local.backend_uri_secret : ""
}
