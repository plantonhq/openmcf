# KubernetesKarapace Terraform module.
#
# Deploys one Karapace schema registry — a MODULE-OWNED-MANIFESTS kind
# (Karapace ships no Helm chart or operator), so the module renders core
# Kubernetes objects directly:
#
#   1. the namespace (optional, create_namespace),
#   2. the SASL password Secret (only when spec.kafka.sasl.password is a
#      literal — never-plaintext-env contract, see the secret resource),
#   3. the registry Deployment `<metadata.name>` and its ClusterIP
#      Service,
#   4. the REST-proxy Deployment `<metadata.name>-rest` and its Service
#      (optional, rest_proxy.enabled) — the same engine image with the
#      role flags flipped, wired to the registry Service.
#
# Both roles are configured purely through KARAPACE_* environment
# variables (pydantic env mechanism: config key X → KARAPACE_<upper X>),
# matching upstream's own compose reference. Schema storage is
# Kafka-native — no PVCs, no databases; the schemas live in a compacted
# topic on the connected Kafka cluster.

resource "kubernetes_namespace_v1" "this" {
  count = var.spec.create_namespace ? 1 : 0

  metadata {
    name   = local.namespace
    labels = local.base_labels
  }
}

# Literal spec.kafka.sasl.password materialized into a module-owned
# Secret. WHY A SECRET AND NOT A PLAIN ENV VALUE: the pod spec is
# readable by anyone with get-pod RBAC (and lands in audit logs, kubectl
# describe output, and controller caches); Secret VALUES have their own,
# stricter ACL. The KARAPACE_SASL_PLAIN_PASSWORD env var references this
# Secret (or the user-supplied password_secret) via secretKeyRef — the
# password never appears in any pod spec. Not created when the spec
# references an existing Secret.
resource "kubernetes_secret_v1" "sasl_password" {
  count = local.create_sasl_secret ? 1 : 0

  metadata {
    name      = local.sasl_secret_name
    namespace = local.namespace
    labels    = local.base_labels
  }

  type = "Opaque"
  # The provider's `data` argument takes plaintext and base64-encodes it
  # on the wire.
  data = {
    password = var.spec.kafka.sasl.password
  }

  depends_on = [kubernetes_namespace_v1.this]
}

# The schema-registry Deployment: one `karapace` container running
# upstream's registry entrypoint. The production image declares NO
# ENTRYPOINT/CMD of its own — upstream's compose reference starts the
# registry role with `python3 -m karapace` (and the REST role with
# `python3 -m karapace.kafka_rest_apis`); the KARAPACE_KARAPACE_REGISTRY
# / KARAPACE_KARAPACE_REST flags select which API surface the started
# process serves.
resource "kubernetes_deployment_v1" "registry" {
  metadata {
    name      = local.registry_name
    namespace = local.namespace
    labels    = local.registry_labels
  }

  wait_for_rollout = true

  spec {
    replicas = local.registry_replicas

    # The selector is the role's immutable pod-selection identity — the
    # role-specific "app" value keeps the registry Service from routing
    # to REST-proxy pods (same image, same namespace) and vice versa.
    selector {
      match_labels = local.registry_selector_labels
    }

    template {
      metadata {
        labels = local.registry_labels
      }

      spec {
        node_selector = length(try(var.spec.node_selector, {})) > 0 ? var.spec.node_selector : null

        # Scheduling knobs are registry-scoped per the spec contract
        # ("Node selector for the registry pods"); the REST-proxy role
        # carries only replicas/port/resources.
        dynamic "toleration" {
          for_each = try(var.spec.tolerations, [])
          content {
            key                = try(toleration.value.key, "") != "" ? toleration.value.key : null
            operator           = try(toleration.value.operator, "") != "" ? toleration.value.operator : null
            value              = try(toleration.value.value, "") != "" ? toleration.value.value : null
            effect             = try(toleration.value.effect, "") != "" ? toleration.value.effect : null
            toleration_seconds = try(toleration.value.toleration_seconds, null)
          }
        }

        container {
          name    = "karapace"
          image   = local.image
          command = ["python3", "-m", "karapace"]

          port {
            name           = "http"
            container_port = local.registry_port
          }

          dynamic "env" {
            for_each = local.registry_env_head
            content {
              name  = env.value.name
              value = env.value.value
            }
          }

          # PER-POD advertised hostname via the downward API. Upstream's
          # compose reference gives every instance its OWN identity
          # (KARAPACE_ADVERTISED_HOSTNAME = that container's hostname,
          # one per replica) — never a shared name: the leader publishes
          # `advertised_protocol://advertised_hostname:port` through the
          # consumer group and followers forward writes to it, so a
          # shared (Service) name would make followers forward to
          # themselves. The Kubernetes twin of compose's resolvable
          # container hostname is the POD IP (a Deployment pod's bare
          # name does not resolve in cluster DNS — a follower forwarding
          # to it would fail; the IP is directly reachable pod-to-pod).
          # config.py falls back to `host` when unset, which is 0.0.0.0
          # here — so it must be set explicitly.
          env {
            name = "KARAPACE_ADVERTISED_HOSTNAME"
            value_from {
              field_ref {
                field_path = "status.podIP"
              }
            }
          }

          dynamic "env" {
            for_each = local.kafka_connection_env
            content {
              name  = env.value.name
              value = env.value.value
            }
          }

          # The SASL password ALWAYS arrives via secretKeyRef — either
          # the referenced existing Secret or the module-materialized
          # one — never as a plaintext env value.
          dynamic "env" {
            for_each = local.kafka_sasl != null ? [1] : []
            content {
              name = "KARAPACE_SASL_PLAIN_PASSWORD"
              value_from {
                secret_key_ref {
                  name = local.sasl_password_secret_name
                  key  = local.sasl_password_secret_key
                }
              }
            }
          }

          dynamic "env" {
            for_each = local.registry_env_tail
            content {
              name  = env.value.name
              value = env.value.value
            }
          }

          dynamic "resources" {
            for_each = try(var.spec.resources, null) != null ? [1] : []
            content {
              limits   = local.registry_resource_limits
              requests = local.registry_resource_requests
            }
          }

          # /_health is unauthenticated by contract (config.py skip-auth
          # paths), so the probes keep working with HTTP authentication
          # enabled. Scheme follows server_tls — kubelet probes hit the
          # pod directly and must speak whatever the server serves.
          readiness_probe {
            http_get {
              path   = local.health_check_path
              port   = tostring(local.registry_port)
              scheme = local.registry_probe_scheme
            }
            initial_delay_seconds = 10
            period_seconds        = 10
            timeout_seconds       = 5
            failure_threshold     = 6
          }

          # Give the engine time to replay the schemas topic before
          # liveness can restart it.
          liveness_probe {
            http_get {
              path   = local.health_check_path
              port   = tostring(local.registry_port)
              scheme = local.registry_probe_scheme
            }
            initial_delay_seconds = 10
            period_seconds        = 10
            timeout_seconds       = 5
            failure_threshold     = 6
          }

          dynamic "volume_mount" {
            for_each = local.registry_volumes
            content {
              name       = volume_mount.value.name
              mount_path = volume_mount.value.mount_path
              read_only  = true
            }
          }
        }

        dynamic "volume" {
          for_each = local.registry_volumes
          content {
            name = volume.value.name
            secret {
              secret_name = volume.value.secret_name
            }
          }
        }
      }
    }
  }

  depends_on = [
    kubernetes_namespace_v1.this,
    kubernetes_secret_v1.sasl_password,
  ]
}

resource "kubernetes_service_v1" "registry" {
  metadata {
    name      = local.registry_name
    namespace = local.namespace
    labels    = local.registry_labels
  }

  spec {
    type     = "ClusterIP"
    selector = local.registry_selector_labels

    port {
      name        = "http"
      port        = local.registry_port
      target_port = local.registry_port
    }
  }

  depends_on = [kubernetes_deployment_v1.registry]
}

# The optional REST-proxy Deployment: the SAME image with the role flags
# flipped and upstream's REST entrypoint, wired to the registry Service
# for schema lookups and to the same Kafka cluster (identical connection
# env and TLS mounts) for produce/consume. Always serves plain HTTP —
# spec.server_tls covers the registry API only.
resource "kubernetes_deployment_v1" "rest_proxy" {
  count = local.rest_enabled ? 1 : 0

  metadata {
    name      = local.rest_name
    namespace = local.namespace
    labels    = local.rest_labels
  }

  wait_for_rollout = true

  spec {
    replicas = local.rest_replicas

    selector {
      match_labels = local.rest_selector_labels
    }

    template {
      metadata {
        labels = local.rest_labels
      }

      spec {
        container {
          name    = "karapace"
          image   = local.image
          command = ["python3", "-m", "karapace.kafka_rest_apis"]

          port {
            name           = "http"
            container_port = local.rest_port
          }

          dynamic "env" {
            for_each = local.rest_env_head
            content {
              name  = env.value.name
              value = env.value.value
            }
          }

          # Each proxy replica advertises ITSELF (rest_base_uri derives
          # from it) — same downward-API identity as the registry role.
          env {
            name = "KARAPACE_ADVERTISED_HOSTNAME"
            value_from {
              field_ref {
                field_path = "status.podIP"
              }
            }
          }

          dynamic "env" {
            for_each = local.kafka_connection_env
            content {
              name  = env.value.name
              value = env.value.value
            }
          }

          dynamic "env" {
            for_each = local.kafka_sasl != null ? [1] : []
            content {
              name = "KARAPACE_SASL_PLAIN_PASSWORD"
              value_from {
                secret_key_ref {
                  name = local.sasl_password_secret_name
                  key  = local.sasl_password_secret_key
                }
              }
            }
          }

          dynamic "env" {
            for_each = local.rest_env_tail
            content {
              name  = env.value.name
              value = env.value.value
            }
          }

          dynamic "resources" {
            for_each = try(var.spec.rest_proxy.resources, null) != null ? [1] : []
            content {
              limits   = local.rest_resource_limits
              requests = local.rest_resource_requests
            }
          }

          readiness_probe {
            http_get {
              path   = local.health_check_path
              port   = tostring(local.rest_port)
              scheme = "HTTP"
            }
            initial_delay_seconds = 10
            period_seconds        = 10
            timeout_seconds       = 5
            failure_threshold     = 6
          }

          liveness_probe {
            http_get {
              path   = local.health_check_path
              port   = tostring(local.rest_port)
              scheme = "HTTP"
            }
            initial_delay_seconds = 10
            period_seconds        = 10
            timeout_seconds       = 5
            failure_threshold     = 6
          }

          dynamic "volume_mount" {
            for_each = local.rest_volumes
            content {
              name       = volume_mount.value.name
              mount_path = volume_mount.value.mount_path
              read_only  = true
            }
          }
        }

        dynamic "volume" {
          for_each = local.rest_volumes
          content {
            name = volume.value.name
            secret {
              secret_name = volume.value.secret_name
            }
          }
        }
      }
    }
  }

  depends_on = [
    kubernetes_namespace_v1.this,
    kubernetes_secret_v1.sasl_password,
  ]
}

resource "kubernetes_service_v1" "rest_proxy" {
  count = local.rest_enabled ? 1 : 0

  metadata {
    name      = local.rest_name
    namespace = local.namespace
    labels    = local.rest_labels
  }

  spec {
    type     = "ClusterIP"
    selector = local.rest_selector_labels

    port {
      name        = "http"
      port        = local.rest_port
      target_port = local.rest_port
    }
  }

  depends_on = [kubernetes_deployment_v1.rest_proxy]
}
