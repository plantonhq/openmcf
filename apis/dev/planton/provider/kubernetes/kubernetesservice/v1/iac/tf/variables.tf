# Input variables for KubernetesService Terraform module.
# These mirror the KubernetesServiceSpec protobuf schema; StringValueOrRef
# fields (namespace) arrive flattened to plain strings, and enum fields arrive
# as the proto enum value names (e.g. "cluster_ip", "prefer_same_zone").

variable "metadata" {
  description = "Metadata for the service resource"
  type = object({
    name = string
    id   = optional(string)
    org  = optional(string)
    env  = optional(string)
  })
}

variable "spec" {
  description = "Specification for the Kubernetes Service"
  type = object({
    namespace   = optional(string, "default")
    name        = string
    labels      = optional(map(string), {})
    annotations = optional(map(string), {})

    type     = optional(string, "cluster_ip")
    selector = optional(map(string), {})

    ports = optional(list(object({
      name = optional(string, "")
      # TCP / UDP / SCTP
      protocol = optional(string, "TCP")
      # L7 hint: an IANA name ("http") or a kubernetes.io/-prefixed name.
      app_protocol = optional(string, "")
      port         = number
      # Number ("8080") or named container port ("http"); "" means same as port.
      target_port = optional(string, "")
      node_port   = optional(number, 0)
    })), [])

    headless = optional(bool, false)
    # Static virtual IP; "" lets the cluster allocate. Immutable.
    cluster_ip_address = optional(string, "")
    # CNAME target for ExternalName services.
    external_dns_name = optional(string, "")
    external_ips      = optional(list(string), [])

    external_traffic_policy = optional(string, "cluster")
    internal_traffic_policy = optional(string)
    traffic_distribution    = optional(string)

    session_affinity                 = optional(string, "none")
    session_affinity_timeout_seconds = optional(number)

    load_balancer_source_ranges = optional(list(string), [])
    load_balancer_class         = optional(string, "")
    # Tri-state: null defers to the Kubernetes default (true).
    allocate_load_balancer_node_ports = optional(bool)
    health_check_node_port            = optional(number, 0)

    publish_not_ready_addresses = optional(bool, false)

    # Dual-stack: entries are "ipv4" / "ipv6"; policy is "single_stack" /
    # "prefer_dual_stack" / "require_dual_stack".
    ip_families      = optional(list(string), [])
    ip_family_policy = optional(string)
  })
}
