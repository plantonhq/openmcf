# Create namespace for Zalando Postgres Operator (conditionally)
resource "kubernetes_namespace_v1" "postgres_operator" {
  count = var.spec.create_namespace ? 1 : 0

  metadata {
    name   = local.namespace
    labels = local.final_labels
  }
}

# Create Secret for R2 backup credentials (only when backup is configured)
resource "kubernetes_secret_v1" "backup_credentials" {
  count = local.has_backup_config ? 1 : 0

  metadata {
    name      = local.backup_secret_name
    namespace = local.namespace
    labels    = local.final_labels
  }

  data = {
    AWS_ACCESS_KEY_ID     = local.r2_access_key_id
    AWS_SECRET_ACCESS_KEY = local.r2_secret_key
  }

  type = "Opaque"

  depends_on = [
    kubernetes_namespace_v1.postgres_operator
  ]
}

# Create ConfigMap for backup configuration (only when backup is configured)
resource "kubernetes_config_map_v1" "backup_config" {
  count = local.has_backup_config ? 1 : 0

  metadata {
    name      = local.backup_configmap_name
    namespace = local.namespace
    labels    = local.final_labels
  }

  data = {
    # WAL-G flags (defaults to true if not explicitly disabled)
    USE_WALG_BACKUP        = tostring(local.enable_wal_g_backup)
    USE_WALG_RESTORE       = tostring(local.enable_wal_g_restore)
    CLONE_USE_WALG_RESTORE = tostring(local.enable_clone_wal_g_restore)

    # S3/R2 configuration
    WALG_S3_PREFIX       = local.walg_s3_prefix
    AWS_ENDPOINT         = local.r2_endpoint
    AWS_REGION           = "auto" # R2 uses "auto" region
    AWS_FORCE_PATH_STYLE = "true" # Required for R2

    # Backup schedule
    BACKUP_SCHEDULE = local.backup_schedule

    # Credentials (reference from Secret)
    AWS_ACCESS_KEY_ID     = local.r2_access_key_id
    AWS_SECRET_ACCESS_KEY = local.r2_secret_key
  }

  depends_on = [
    kubernetes_namespace_v1.postgres_operator,
    kubernetes_secret_v1.backup_credentials
  ]
}

# Deploy Zalando Postgres Operator using Helm
resource "helm_release" "postgres_operator" {
  name       = local.helm_chart_name
  repository = local.helm_chart_repository
  chart      = local.helm_chart_name
  version    = local.helm_chart_version
  namespace  = local.namespace

  # Wait for deployment to complete
  wait          = true
  wait_for_jobs = true
  timeout       = 180

  # Configure inherited labels that will be propagated to all PostgreSQL databases
  set {
    name  = "configKubernetes.inherited_labels[0]"
    value = "resource"
  }

  set {
    name  = "configKubernetes.inherited_labels[1]"
    value = "organization"
  }

  set {
    name  = "configKubernetes.inherited_labels[2]"
    value = "environment"
  }

  set {
    name  = "configKubernetes.inherited_labels[3]"
    value = "resource_kind"
  }

  set {
    name  = "configKubernetes.inherited_labels[4]"
    value = "resource_id"
  }

  # Configure pod environment ConfigMap (only when backup is configured)
  dynamic "set" {
    for_each = local.has_backup_config ? [1] : []
    content {
      name  = "configKubernetes.pod_environment_configmap"
      value = "${local.namespace}/${local.backup_configmap_name}"
    }
  }

  # Configure operator container resources
  set {
    name  = "resources.requests.cpu"
    value = var.spec.container.resources.requests.cpu
  }

  set {
    name  = "resources.requests.memory"
    value = var.spec.container.resources.requests.memory
  }

  set {
    name  = "resources.limits.cpu"
    value = var.spec.container.resources.limits.cpu
  }

  set {
    name  = "resources.limits.memory"
    value = var.spec.container.resources.limits.memory
  }

  depends_on = [
    kubernetes_namespace_v1.postgres_operator,
    kubernetes_secret_v1.backup_credentials,
    kubernetes_config_map_v1.backup_config
  ]
}

