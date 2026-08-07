# The headless governing Service of the StatefulSet. Headless
# (cluster_ip = "None") is what gives each replica its stable per-pod DNS name
# (<pod>.<service>.<namespace>.svc.cluster.local) — a regular ClusterIP would
# only load-balance and could never address individual members. It is ALWAYS
# created: the StatefulSet API requires spec.service_name to reference an
# existing Service regardless of whether the app exposes ports, and it must
# exist before the StatefulSet so pods can resolve peer DNS at startup.
#
# External exposure attaches through first-class ingress kinds referencing this
# Service by name; the module itself creates no exposure infrastructure.
resource "kubernetes_service_v1" "this" {
  metadata {
    name      = local.kube_service_name
    namespace = local.namespace
    labels    = local.final_labels
  }

  spec {
    # "None" is the headless marker — no virtual IP, DNS resolves straight to
    # pod IPs, and each pod gets its own stable DNS record.
    cluster_ip = "None"
    selector   = local.selector_labels

    # Peers must discover each other BEFORE they can pass readiness — a
    # bootstrapping member needs DNS for pods that are themselves still
    # bootstrapping. Publishing not-ready addresses breaks that deadlock.
    publish_not_ready_addresses = true

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

  depends_on = [kubernetes_namespace.this]
}
