# KubernetesMlflow Terraform module.
#
# Deploys one MLflow tracking server + model registry as MODULE-OWNED
# manifests (MLflow publishes no Helm chart; the official
# ghcr.io/mlflow/mlflow image is the distribution). SECURED BY DEFAULT:
# basic authentication is ON unless the spec disables it — upstream's
# open server and its admin/password1234 example never ship. Nothing
# credential-bearing lands in any rendered pod spec: the backend-store
# URI reaches the server as an env var from a module-composed Secret, the
# admin password lives only inside Secrets, and the basic-auth ini is
# Secret-mounted. The exact same resource set renders from the Pulumi
# module — keep them in lockstep.

# The optional installation namespace. Created before everything; deleted
# with the resource.
resource "kubernetes_namespace_v1" "mlflow" {
  count = var.spec.create_namespace ? 1 : 0

  metadata {
    name   = local.namespace
    labels = local.labels
  }
}

# ---------------------------------------------------------------------------
# The module-generated admin password. Generation-shape arguments are
# ignored after creation so an IMPORTED credential never silently
# regenerates (Pulumi twin: IgnoreChanges on the same arguments).
# Letters+digits only: the password lands inside an ini file and in
# users' MLFLOW_TRACKING_PASSWORD env vars — symbol-free avoids quoting
# bugs; the larger length compensates the smaller alphabet.
# ---------------------------------------------------------------------------

resource "random_password" "admin_password" {
  count       = local.admin_secret_module_owned ? 1 : 0
  length      = 24
  special     = false
  min_upper   = 2
  min_lower   = 2
  min_numeric = 2

  lifecycle {
    ignore_changes = [length, special, upper, lower, numeric, min_lower, min_numeric, min_special, min_upper, override_special]
  }
}

resource "kubernetes_secret_v1" "admin_auth" {
  count = local.admin_secret_module_owned ? 1 : 0

  metadata {
    name      = local.admin_secret_name
    namespace = local.namespace
    labels    = local.labels
  }
  type = "Opaque"
  data = {
    "password" = random_password.admin_password[0].result
  }

  depends_on = [kubernetes_namespace_v1.mlflow]
}

# ---------------------------------------------------------------------------
# The referenced credential reads. The referenced Secrets are created by
# OTHER components (e.g. the KubernetesPostgres app Secret), so they
# exist before this module applies. depends_on a module-created resource
# DEFERS the read to apply time on a fresh plan — an offline plan (no
# cluster) stays green with "(known after apply)" while a real apply
# reads the live value (Pulumi twin: the DryRun-gated GetSecret). The
# dependency anchor is a marker ConfigMap because every module-created
# Secret here is conditional.
# ---------------------------------------------------------------------------

resource "kubernetes_config_map_v1" "read_anchor" {
  count = local.backend_type != "sqlite" || (local.auth_enabled && !local.admin_secret_module_owned) ? 1 : 0

  metadata {
    name      = "${local.name}-read-anchor"
    namespace = local.namespace
    labels    = local.labels
  }
  data = {
    # Marker only — this ConfigMap exists to defer the credential
    # data-source reads below to apply time; it carries no configuration.
    purpose = "defers-referenced-secret-reads-to-apply"
  }

  depends_on = [kubernetes_namespace_v1.mlflow]
}

data "kubernetes_secret_v1" "db_password" {
  count = local.backend_type != "sqlite" ? 1 : 0

  metadata {
    name      = local.db_password_secret
    namespace = local.namespace
  }

  depends_on = [kubernetes_config_map_v1.read_anchor]
}

data "kubernetes_secret_v1" "byo_admin_password" {
  count = local.auth_enabled && !local.admin_secret_module_owned ? 1 : 0

  metadata {
    name      = local.admin_secret_name
    namespace = local.namespace
  }

  depends_on = [kubernetes_config_map_v1.read_anchor]
}

# ---------------------------------------------------------------------------
# Composed Secrets.
# ---------------------------------------------------------------------------

# The backend-store URI (`postgresql://user:pass@host:port/db`) — the
# exported handle other tools can mount to reach the same tracking
# database. Userinfo urlencoded (Pulumi parity).
resource "kubernetes_secret_v1" "backend_uri" {
  count = local.backend_type != "sqlite" ? 1 : 0

  metadata {
    name      = local.backend_uri_secret
    namespace = local.namespace
    labels    = local.labels
  }
  type = "Opaque"
  data = {
    "uri" = "${local.db_protocol}://${urlencode(local.db_user)}:${urlencode(data.kubernetes_secret_v1.db_password[0].data[local.db_password_secret_key])}@${local.db_host}:${local.db_port}/${local.db_name}"
  }

  depends_on = [kubernetes_namespace_v1.mlflow]
}

# The Flask CSRF signing key: the auth app raises at create_app without
# MLFLOW_FLASK_SERVER_SECRET_KEY (server source at the pin).
# Letters+digits only (env-consumed); module-generated — there is no
# user-facing reason to model a knob for a purely internal signing key.
resource "random_password" "flask_secret_key" {
  count       = local.auth_enabled ? 1 : 0
  length      = 40
  special     = false
  min_upper   = 2
  min_lower   = 2
  min_numeric = 2

  # An IMPORTED key never silently regenerates: rotation stays an
  # explicit verb, never plan fallout (the Pulumi twin ignores the same
  # generation-shape arguments).
  lifecycle {
    ignore_changes = [length, special, upper, lower, numeric, min_lower, min_numeric, min_special, min_upper, override_special]
  }
}

# The basic-auth ini — the server's own contract
# (mlflow/server/auth/basic_auth.ini at the pin). The auth database
# follows the backend store: the same database on the postgres/mysql
# arms (upstream-supported; keeps auth state as durable as tracking
# state), a sqlite file beside the tracking data on the sqlite arm.
# The Secret's second key is the Flask CSRF signing key the auth app
# REFUSES to start without — env-wired from this shared Secret so every
# replica carries one consistent value (upstream's own multi-server
# requirement).
resource "kubernetes_secret_v1" "auth_config" {
  count = local.auth_enabled ? 1 : 0

  metadata {
    name      = local.auth_config_secret_name
    namespace = local.namespace
    labels    = local.labels
  }
  type = "Opaque"
  data = {
    (local.auth_config_file_name) = <<-EOT
      [mlflow]
      default_permission = ${local.default_permission}
      database_uri = ${local.backend_type == "sqlite" ? local.sqlite_auth_db_uri : "${local.db_protocol}://${urlencode(local.db_user)}:${urlencode(data.kubernetes_secret_v1.db_password[0].data[local.db_password_secret_key])}@${local.db_host}:${local.db_port}/${local.db_name}"}
      admin_username = ${local.admin_username}
      admin_password = ${local.admin_secret_module_owned ? random_password.admin_password[0].result : data.kubernetes_secret_v1.byo_admin_password[0].data[local.admin_secret_key]}
      authorization_function = mlflow.server.auth:authenticate_request_basic_auth
    EOT
    flask_secret_key              = random_password.flask_secret_key[0].result
  }

  depends_on = [kubernetes_namespace_v1.mlflow]
}

# ---------------------------------------------------------------------------
# PVCs (sqlite backend / local artifacts). Both ReadWriteOnce — the CEL
# contract caps replicas at 1 whenever either exists.
# ---------------------------------------------------------------------------

resource "kubernetes_persistent_volume_claim_v1" "data" {
  count = local.data_pvc_enabled ? 1 : 0

  metadata {
    name      = local.data_pvc_name
    namespace = local.namespace
    labels    = local.labels
  }
  spec {
    access_modes = ["ReadWriteOnce"]
    resources {
      requests = {
        storage = local.data_pvc_size
      }
    }
    storage_class_name = local.data_pvc_storage_class != "" ? local.data_pvc_storage_class : null
  }

  # Bind-timing doctrine: WaitForFirstConsumer claims are correctly
  # Pending until the server pod schedules.
  wait_until_bound = false

  depends_on = [kubernetes_namespace_v1.mlflow]
}

resource "kubernetes_persistent_volume_claim_v1" "artifacts" {
  count = local.artifacts_pvc_enabled ? 1 : 0

  metadata {
    name      = local.artifacts_pvc_name
    namespace = local.namespace
    labels    = local.labels
  }
  spec {
    access_modes = ["ReadWriteOnce"]
    resources {
      requests = {
        storage = local.artifacts_pvc_size
      }
    }
    storage_class_name = local.artifacts_pvc_storage_class != "" ? local.artifacts_pvc_storage_class : null
  }

  wait_until_bound = false

  depends_on = [kubernetes_namespace_v1.mlflow]
}

# ---------------------------------------------------------------------------
# The tracking-server Deployment.
# ---------------------------------------------------------------------------

resource "kubernetes_deployment_v1" "mlflow" {
  metadata {
    name      = local.name
    namespace = local.namespace
    labels    = local.labels
  }

  spec {
    replicas = local.replicas

    # Strategy follows the volume truth: any RWO PVC binds one pod, so
    # updates must Recreate; stateless shapes roll.
    strategy {
      type = local.deployment_strategy
    }

    selector {
      match_labels = local.selector_labels
    }

    template {
      metadata {
        labels = merge(local.labels, local.selector_labels)
      }

      spec {
        container {
          name    = "mlflow"
          image   = local.image
          command = local.server_args

          port {
            name           = "http"
            container_port = local.server_port
            protocol       = "TCP"
          }

          dynamic "env" {
            for_each = local.plain_env
            content {
              name  = env.key
              value = env.value
            }
          }

          dynamic "env" {
            for_each = local.secret_env
            content {
              name = env.key
              value_from {
                secret_key_ref {
                  name = env.value.secret_name
                  key  = env.value.secret_key
                }
              }
            }
          }

          # /health is the server's own unauthenticated health contract
          # (upstream's deployment reference probes it) — it answers
          # even with basic-auth on.
          liveness_probe {
            http_get {
              path = "/health"
              port = "http"
            }
            initial_delay_seconds = 15
            period_seconds        = 20
          }

          readiness_probe {
            http_get {
              path = "/health"
              port = "http"
            }
            initial_delay_seconds = 5
            period_seconds        = 10
          }

          dynamic "resources" {
            for_each = try(var.spec.server.resources, null) != null ? [var.spec.server.resources] : []
            content {
              requests = {
                for k, v in {
                  cpu    = try(resources.value.requests.cpu, "")
                  memory = try(resources.value.requests.memory, "")
                } : k => v if v != ""
              }
              limits = {
                for k, v in {
                  cpu    = try(resources.value.limits.cpu, "")
                  memory = try(resources.value.limits.memory, "")
                } : k => v if v != ""
              }
            }
          }

          dynamic "volume_mount" {
            for_each = local.data_pvc_enabled ? [1] : []
            content {
              name       = "data"
              mount_path = local.data_mount_path
            }
          }

          dynamic "volume_mount" {
            for_each = local.artifacts_pvc_enabled ? [1] : []
            content {
              name       = "artifacts"
              mount_path = local.artifacts_mount_path
            }
          }

          dynamic "volume_mount" {
            for_each = local.auth_enabled ? [1] : []
            content {
              name       = "auth-config"
              mount_path = local.auth_config_mount_path
              read_only  = true
            }
          }

          dynamic "volume_mount" {
            for_each = local.artifact_type == "gcs" && try(local.artifact_gcs.credentials_secret, null) != null ? [1] : []
            content {
              name       = "gcs-credentials"
              mount_path = local.gcs_credentials_mount_path
              read_only  = true
            }
          }
        }

        dynamic "volume" {
          for_each = local.data_pvc_enabled ? [1] : []
          content {
            name = "data"
            persistent_volume_claim {
              claim_name = local.data_pvc_name
            }
          }
        }

        dynamic "volume" {
          for_each = local.artifacts_pvc_enabled ? [1] : []
          content {
            name = "artifacts"
            persistent_volume_claim {
              claim_name = local.artifacts_pvc_name
            }
          }
        }

        dynamic "volume" {
          for_each = local.auth_enabled ? [1] : []
          content {
            name = "auth-config"
            secret {
              secret_name = local.auth_config_secret_name
              # Only the ini materializes as a file — the Secret's
              # flask_secret_key rides env, never the filesystem.
              items {
                key  = local.auth_config_file_name
                path = local.auth_config_file_name
              }
            }
          }
        }

        dynamic "volume" {
          for_each = local.artifact_type == "gcs" && try(local.artifact_gcs.credentials_secret, null) != null ? [1] : []
          content {
            name = "gcs-credentials"
            secret {
              secret_name = local.artifact_gcs.credentials_secret.secret_name
            }
          }
        }

        dynamic "toleration" {
          for_each = try(var.spec.scheduling.tolerations, [])
          content {
            key                = toleration.value.key != "" ? toleration.value.key : null
            operator           = toleration.value.operator != "" ? toleration.value.operator : null
            value              = toleration.value.value != "" ? toleration.value.value : null
            effect             = toleration.value.effect != "" ? toleration.value.effect : null
            toleration_seconds = try(toleration.value.toleration_seconds, null)
          }
        }

        node_selector = length(try(var.spec.scheduling.node_selector, {})) > 0 ? var.spec.scheduling.node_selector : null
      }
    }
  }

  # Rollout waiting stays ON (the provider default): a server that
  # cannot reach its backend store or artifact credentials should fail
  # THIS apply, not the first experiment run.

  depends_on = [
    kubernetes_namespace_v1.mlflow,
    kubernetes_secret_v1.admin_auth,
    kubernetes_secret_v1.backend_uri,
    kubernetes_secret_v1.auth_config,
    kubernetes_persistent_volume_claim_v1.data,
    kubernetes_persistent_volume_claim_v1.artifacts,
  ]
}

# ---------------------------------------------------------------------------
# The Service (the front door — ClusterIP by default; exposure composes
# from first-class kinds over the exported handle).
# ---------------------------------------------------------------------------

resource "kubernetes_service_v1" "mlflow" {
  metadata {
    name        = local.name
    namespace   = local.namespace
    labels      = local.labels
    annotations = length(local.service_annotations) > 0 ? local.service_annotations : null
  }

  spec {
    type     = local.service_type
    selector = local.selector_labels

    port {
      name        = "http"
      port        = local.server_port
      target_port = "http"
      protocol    = "TCP"
    }
  }

  depends_on = [kubernetes_namespace_v1.mlflow]
}

# ---------------------------------------------------------------------------
# The `mlflow gc` CronJob — permanent removal of soft-deleted
# runs/experiments older than the retention window. Same image, same
# backend env; artifact-store credentials ride along so gc can delete the
# artifacts too. The backend URI rides env expansion ($(VAR) — resolved
# by the kubelet from the Secret-sourced env var) so the URI itself never
# appears in the rendered spec.
# ---------------------------------------------------------------------------

resource "kubernetes_cron_job_v1" "gc" {
  count = local.gc_enabled ? 1 : 0

  metadata {
    name      = "${local.name}-gc"
    namespace = local.namespace
    labels    = local.labels
  }

  spec {
    schedule           = local.gc_schedule
    concurrency_policy = "Forbid"

    job_template {
      metadata {
        labels = local.labels
      }
      spec {
        backoff_limit = 2
        template {
          metadata {
            labels = local.labels
          }
          spec {
            restart_policy = "OnFailure"

            container {
              name  = "mlflow-gc"
              image = local.image
              command = [
                "mlflow", "gc",
                "--older-than", local.gc_older_than,
                "--backend-store-uri", "$(MLFLOW_BACKEND_STORE_URI)",
              ]

              dynamic "env" {
                for_each = local.plain_env
                content {
                  name  = env.key
                  value = env.value
                }
              }

              dynamic "env" {
                for_each = local.secret_env
                content {
                  name = env.key
                  value_from {
                    secret_key_ref {
                      name = env.value.secret_name
                      key  = env.value.secret_key
                    }
                  }
                }
              }

              dynamic "volume_mount" {
                for_each = local.data_pvc_enabled ? [1] : []
                content {
                  name       = "data"
                  mount_path = local.data_mount_path
                }
              }
            }

            dynamic "volume" {
              for_each = local.data_pvc_enabled ? [1] : []
              content {
                name = "data"
                persistent_volume_claim {
                  claim_name = local.data_pvc_name
                }
              }
            }

            node_selector = length(try(var.spec.scheduling.node_selector, {})) > 0 ? var.spec.scheduling.node_selector : null
          }
        }
      }
    }
  }

  depends_on = [kubernetes_deployment_v1.mlflow]
}

# ---------------------------------------------------------------------------
# The optional ServiceMonitor (monitoring.coreos.com/v1) — requires the
# Prometheus operator CRDs on the cluster (a KubernetesKubePrometheusStack
# composes naturally); deploying without them fails loudly by design.
# kubectl_manifest plans without the CRD existing.
# ---------------------------------------------------------------------------

resource "kubectl_manifest" "service_monitor" {
  count = local.service_monitor_enabled ? 1 : 0

  yaml_body = yamlencode({
    apiVersion = "monitoring.coreos.com/v1"
    kind       = "ServiceMonitor"
    metadata = {
      name      = "${local.name}-metrics"
      namespace = local.namespace
      labels    = local.labels
    }
    spec = {
      selector = {
        matchLabels = local.selector_labels
      }
      endpoints = [
        {
          port     = "http"
          path     = "/metrics"
          interval = "30s"
        }
      ]
    }
  })

  depends_on = [kubernetes_deployment_v1.mlflow]
}
