# ClusterIP Service fronting the Deployment's pods — one Service port per app
# container port, service_port defaulting to the container port. No Service is
# created when the app container exposes no ports (workers, consumers).
#
# External exposure attaches through first-class ingress kinds referencing this
# Service by name; the module itself creates no exposure infrastructure.
resource "kubernetes_service_v1" "this" {
  count = local.create_service ? 1 : 0

  metadata {
    name      = local.kube_service_name
    namespace = local.namespace
    labels    = local.final_labels
  }

  spec {
    type     = "ClusterIP"
    selector = local.selector_labels

    dynamic "port" {
      for_each = local.service_ports
      content {
        name         = port.value.name
        protocol     = port.value.protocol
        port         = port.value.port
        target_port  = port.value.target_port
        app_protocol = port.value.app_protocol
      }
    }
  }

  depends_on = [kubernetes_deployment_v1.this]
}
