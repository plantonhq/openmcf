# Kubernetes Ingress Terraform module.
#
# Creation deliberately does NOT block on a controller claiming the Ingress
# (wait_for_load_balancer = false; the Pulumi module's skipAwait annotation is
# the exact same choice). An Ingress object is valid without a controller —
# infra charts routinely deploy the workload and its exposure before the
# ingress controller wave — and blocking every deploy until a controller
# populates the load-balancer status would couple this kind to cluster addon
# ordering. The load-balancer address handles surface through the outputs as
# soon as a controller reconciles the object.

resource "kubernetes_ingress_v1" "ingress" {
  wait_for_load_balancer = false

  metadata {
    name        = var.spec.name
    namespace   = local.namespace
    labels      = local.labels
    annotations = local.annotations
  }

  spec {
    # "" defers to the cluster's default IngressClass.
    ingress_class_name = var.spec.ingress_class_name != "" ? var.spec.ingress_class_name : null

    # Backend for requests no rule matches. Exactly one of number/name is set
    # per the spec's CEL contract — the provider requires the same.
    dynamic "default_backend" {
      for_each = var.spec.default_backend != null ? [var.spec.default_backend] : []
      content {
        service {
          # service_name is a flattened StringValueOrRef (default kind
          # KubernetesService) — this is where a workload's exported `service`
          # output lands when charts wire exposure to a workload.
          name = default_backend.value.service_name
          port {
            number = default_backend.value.port_number > 0 ? default_backend.value.port_number : null
            name   = default_backend.value.port_name != "" ? default_backend.value.port_name : null
          }
        }
      }
    }

    dynamic "tls" {
      for_each = var.spec.tls
      content {
        hosts = length(tls.value.hosts) > 0 ? tls.value.hosts : null
        # A flattened StringValueOrRef to a KubernetesSecret. With cert-manager
        # the Secret need not exist yet — the issuer annotation instructs
        # cert-manager to create it under exactly this name.
        secret_name = tls.value.secret_name != "" ? tls.value.secret_name : null
      }
    }

    dynamic "rule" {
      for_each = var.spec.rules
      content {
        host = rule.value.host != "" ? rule.value.host : null
        http {
          dynamic "path" {
            for_each = rule.value.paths
            content {
              path      = path.value.path != "" ? path.value.path : null
              path_type = lookup(local.path_type_map, coalesce(path.value.path_type, "prefix"), "Prefix")
              backend {
                service {
                  name = path.value.backend.service_name
                  port {
                    number = path.value.backend.port_number > 0 ? path.value.backend.port_number : null
                    name   = path.value.backend.port_name != "" ? path.value.backend.port_name : null
                  }
                }
              }
            }
          }
        }
      }
    }
  }
}
