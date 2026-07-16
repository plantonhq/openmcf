variable "metadata" {
  description = "Metadata for the resource, including name and labels"
  type = object({
    name    = string,
    id      = optional(string),
    org     = optional(string),
    env     = optional(string),
    labels  = optional(map(string)),
    tags    = optional(list(string)),
    version = optional(object({ id = string, message = string }))
  })
}

variable "spec" {
  description = "Azure Load Balancer specification"
  type = object({
    # The Azure region where the load balancer will be created.
    region = string

    # The resource group name. References are resolved to a literal name
    # by the platform before the module runs.
    resource_group = string

    # The name of the load balancer, unique within the resource group.
    name = string

    # The SKU, as the spec enum's name string (STANDARD/GATEWAY). Unset
    # applies STANDARD -- the production SKU. Basic is not modeled (Azure
    # retired it September 2025).
    sku = optional(string)

    # The SKU tier, as the spec enum's name string (REGIONAL/GLOBAL).
    # Unset applies REGIONAL. GLOBAL (cross-region) requires STANDARD.
    sku_tier = optional(string)

    # Edge Zone pinning for edge-computing workloads (fixed at creation).
    edge_zone = optional(string)

    # The frontends that receive traffic (at least one). Each is public
    # (public IP or prefix) or internal (subnet); rules target a frontend
    # by name.
    frontend_ip_configurations = list(object({
      # A label for this frontend, unique within the load balancer.
      name = string

      # INTERNAL frontend: the subnet the private address lives in, as a
      # resolved ARM ID.
      subnet_id = optional(string)

      # PUBLIC frontend: the fronting public IP, as a resolved ARM ID.
      public_ip_address_id = optional(string)

      # PUBLIC frontend: a public IP prefix to draw from, as a resolved
      # ARM ID (the SNAT-scaling shape).
      public_ip_prefix_id = optional(string)

      # For internal frontends: pin a specific private address (unset
      # lets Azure allocate dynamically).
      private_ip_address = optional(string)

      # The address family, as the spec enum's name string (IPV4/IPV6).
      # Unset applies Azure's default (IPv4).
      private_ip_address_version = optional(string)

      # Availability zones for an internal frontend's private address
      # (["1","2","3"] for zone redundancy).
      zones = optional(list(string), [])

      # Gateway-SKU load balancer frontend to chain this frontend behind
      # (NVA service chaining), as an ARM ID.
      gateway_load_balancer_frontend_ip_configuration_id = optional(string)
    }))

    # Backend address pools. Membership is expressed member-side (NIC and
    # VMSS references to the exported pool IDs) except for the inline
    # IP-based addresses.
    backend_pools = optional(list(object({
      # The pool name, unique within the load balancer.
      name = string

      # The virtual network of IP-based members, as a resolved ARM ID.
      # Required when addresses or synchronous_mode are set.
      virtual_network_id = optional(string)

      # IP-based membership sync mode, as the spec enum's name string
      # (AUTOMATIC/MANUAL). Unset sends nothing.
      synchronous_mode = optional(string)

      # GATEWAY SKU only: the VXLAN tunnel interfaces through which
      # chained traffic reaches the NVAs in this pool.
      tunnel_interfaces = optional(list(object({
        # The VXLAN network identifier (conventionally 800/801).
        identifier = number

        # The tunnel port (conventionally 10800/10801).
        port = number

        # Encapsulation protocol, as the spec enum's name string
        # (TUNNEL_PROTOCOL_NONE/NATIVE/VXLAN).
        protocol = string

        # Direction, as the spec enum's name string
        # (TUNNEL_TYPE_NONE/INTERNAL/EXTERNAL).
        type = string
      })), [])

      # Inline IP-based members (requires virtual_network_id), or -- for
      # GLOBAL-tier pools -- regional load balancer frontends.
      addresses = optional(list(object({
        # A label for this member, unique within the pool.
        name = string

        # REGIONAL tier: the member's private IP inside the pool's vnet.
        ip_address = optional(string)

        # GLOBAL tier: the regional load balancer frontend this member
        # represents, as an ARM ID.
        load_balancer_frontend_ip_configuration_id = optional(string)
      })), [])
    })), [])

    # Health probes gating rule backends.
    health_probes = optional(list(object({
      # The probe name, unique within the load balancer.
      name = string

      # The probe protocol, as the spec enum's name string
      # (PROBE_TCP/PROBE_HTTP/PROBE_HTTPS). Unset applies Azure's
      # default (Tcp).
      protocol = optional(string)

      # The port to probe (1-65535).
      port = number

      # The GET path for HTTP/HTTPS probes (required for them, forbidden
      # for TCP -- spec-level validation enforces both).
      request_path = optional(string)

      # Interval between probes in seconds (min 5, default 15).
      interval_in_seconds = optional(number, 15)

      # Consecutive failures before unhealthy (default 2).
      number_of_probes = optional(number, 2)

      # Consecutive successes required before healthy again (default 1).
      probe_threshold = optional(number, 1)
    })), [])

    # Load-balancing rules mapping a frontend port/protocol to a backend
    # pool and port.
    rules = optional(list(object({
      # The rule name, unique within the load balancer.
      name = string

      # The frontend this rule listens on, by frontend name. Optional
      # when exactly one frontend exists.
      frontend_ip_configuration_name = optional(string)

      # Transport protocol, as the spec enum's name string (TCP/UDP/ALL).
      # ALL creates an HA-ports rule (internal STANDARD frontends, ports 0).
      protocol = string

      # Frontend port (0-65534; 0 only for HA ports).
      frontend_port = number

      # Backend port (0-65535; 0 only for HA ports).
      backend_port = number

      # The target pool(s) by name (two only on GATEWAY SKU).
      backend_pool_names = list(string)

      # The gating probe by name (optional -- but production rules
      # should probe).
      probe_name = optional(string)

      # Session persistence, as the spec enum's name string
      # (DEFAULT/SOURCE_IP/SOURCE_IP_PROTOCOL). Unset applies Azure's
      # default (5-tuple, no persistence).
      load_distribution = optional(string)

      # TCP idle timeout in minutes (4-100, default 4).
      idle_timeout_in_minutes = optional(number, 4)

      # Floating IP / Direct Server Return (SQL AlwaysOn et al).
      floating_ip_enabled = optional(bool, false)

      # Send TCP reset on idle-timeout drop so clients fail fast.
      tcp_reset_enabled = optional(bool, false)

      # Disable this rule's implicit outbound SNAT (set when the pool
      # egresses via an explicit outbound rule or a NAT gateway).
      disable_outbound_snat = optional(bool, false)
    })), [])

    # Inbound NAT rules: port forwarding to individual instances.
    # Single-target mode sets frontend_port (the NIC association
    # completes the attachment member-side); pool-style mode sets
    # backend_pool_name + frontend_port_start/end.
    nat_rules = optional(list(object({
      # The NAT rule name, unique within the load balancer.
      name = string

      # The frontend this rule listens on, by frontend name. Optional
      # when exactly one frontend exists.
      frontend_ip_configuration_name = optional(string)

      # Transport protocol, as the spec enum's name string (TCP/UDP/ALL).
      protocol = string

      # SINGLE-TARGET mode: the one frontend port to forward.
      frontend_port = optional(number, 0)

      # The backend instance port receiving the forwarded traffic.
      backend_port = number

      # POOL-STYLE mode: the pool whose members each get a dedicated
      # frontend port from the range below.
      backend_pool_name = optional(string)

      # First frontend port of the pool-style range (inclusive).
      frontend_port_start = optional(number, 0)

      # Last frontend port of the pool-style range (inclusive).
      frontend_port_end = optional(number, 0)

      # Floating IP / Direct Server Return for the forwarded traffic.
      floating_ip_enabled = optional(bool, false)

      # Send TCP reset on idle-timeout drop.
      tcp_reset_enabled = optional(bool, false)

      # TCP idle timeout in minutes (4-30, default 4).
      idle_timeout_in_minutes = optional(number, 4)
    })), [])

    # Outbound rules: explicit SNAT for a pool's egress via public
    # frontends (STANDARD SKU, public frontends only).
    outbound_rules = optional(list(object({
      # The outbound rule name, unique within the load balancer.
      name = string

      # The PUBLIC frontends whose addresses egress uses, by frontend
      # name (at least one; more frontends = more SNAT ports).
      frontend_ip_configuration_names = list(string)

      # The pool whose members' egress this rule governs, by pool name.
      backend_pool_name = string

      # Transport protocol, as the spec enum's name string (TCP/UDP/ALL).
      protocol = string

      # SNAT ports per backend instance (default 1024; 0 lets Azure
      # divide the budget, which churns on scale events).
      allocated_outbound_ports = optional(number, 1024)

      # Send TCP reset on idle-timeout drop for outbound connections.
      tcp_reset_enabled = optional(bool, false)

      # Idle timeout in minutes for outbound connections (default 4).
      idle_timeout_in_minutes = optional(number, 4)
    })), [])

    # Free-form user tags, merged over the metadata-derived tags (user
    # tags win on key collision).
    tags = optional(map(string), {})
  })
}
