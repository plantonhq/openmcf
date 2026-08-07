# Kubernetes Service Terraform module.
#
# Optional spec fields are sent to the API only when the user set them —
# sending an empty value is NOT the same as omitting the field for several of
# these (clusterIP, healthCheckNodePort, loadBalancerClass are immutable or
# type-gated, and the API rejects them on the wrong service type), so each one
# is guarded by the same condition the Pulumi module uses.

resource "kubernetes_service_v1" "service" {
  metadata {
    name        = var.spec.name
    namespace   = local.namespace
    labels      = local.labels
    annotations = local.annotations
  }

  spec {
    type = local.service_type

    # Selector does not apply to ExternalName (a pure DNS alias) or to
    # selectorless services fronting manually-managed endpoints.
    selector = (!local.is_external_name && length(try(var.spec.selector, {})) > 0) ? var.spec.selector : null

    # Headless IS clusterIP "None"; otherwise honor an explicitly requested
    # static cluster IP. Both are create-time-only (the field is immutable).
    cluster_ip = (
      var.spec.headless
      ? "None"
      : (try(var.spec.cluster_ip_address, "") != "" ? var.spec.cluster_ip_address : null)
    )

    # CNAME target for ExternalName services.
    external_name = local.is_external_name ? var.spec.external_dns_name : null

    # External IPs are a proxying concept the API rejects on ExternalName —
    # gated like ports/selector even though the CEL rule already blocks the
    # combination at validation time.
    external_ips = (!local.is_external_name && length(try(var.spec.external_ips, [])) > 0) ? var.spec.external_ips : null

    # Traffic policies are type-gated by the API: external only for
    # externally-reachable types, internal never for ExternalName.
    external_traffic_policy = local.is_externally_reachable ? local.external_traffic_policy : null
    internal_traffic_policy = local.is_external_name ? null : local.internal_traffic_policy

    session_affinity = local.session_affinity
    dynamic "session_affinity_config" {
      for_each = (local.session_affinity == "ClientIP" && try(var.spec.session_affinity_timeout_seconds, null) != null) ? [1] : []
      content {
        client_ip {
          timeout_seconds = var.spec.session_affinity_timeout_seconds
        }
      }
    }

    # LoadBalancer-only knobs — the spec's CEL rules guarantee they are unset
    # for other types, and the API would reject them there anyway.
    load_balancer_source_ranges       = (local.is_load_balancer && length(try(var.spec.load_balancer_source_ranges, [])) > 0) ? var.spec.load_balancer_source_ranges : null
    load_balancer_class               = (local.is_load_balancer && try(var.spec.load_balancer_class, "") != "") ? var.spec.load_balancer_class : null
    allocate_load_balancer_node_ports = local.is_load_balancer ? try(var.spec.allocate_load_balancer_node_ports, null) : null
    health_check_node_port            = (local.is_load_balancer && try(var.spec.health_check_node_port, 0) > 0) ? var.spec.health_check_node_port : null

    publish_not_ready_addresses = var.spec.publish_not_ready_addresses ? true : null

    # Dual-stack: families and policy are only sent when requested; the
    # cluster otherwise assigns from its own configuration.
    ip_families      = length(local.ip_families) > 0 ? local.ip_families : null
    ip_family_policy = local.ip_family_policy

    # Ports do not apply to ExternalName services.
    dynamic "port" {
      for_each = local.is_external_name ? [] : var.spec.ports
      content {
        name         = port.value.name != "" ? port.value.name : null
        protocol     = port.value.protocol
        app_protocol = port.value.app_protocol != "" ? port.value.app_protocol : null
        port         = port.value.port
        # IntOrString upstream: "" means "same as port" (the API's identity
        # mapping) — encoded explicitly because the provider requires a value.
        target_port = port.value.target_port != "" ? port.value.target_port : port.value.port
        node_port   = port.value.node_port > 0 ? port.value.node_port : null
      }
    }
  }

  lifecycle {
    # PARITY-EXCEPTION: the Terraform kubernetes provider (v3.2.x) does not
    # expose spec.trafficDistribution, so only the Pulumi engine can apply
    # this field. Fail the plan loudly instead of silently dropping a set
    # field — deploy through the Pulumi module when traffic_distribution
    # matters.
    precondition {
      condition     = try(var.spec.traffic_distribution, null) == null
      error_message = "traffic_distribution is not supported by the Terraform kubernetes provider; deploy this manifest with the Pulumi engine or unset the field."
    }
  }
}
