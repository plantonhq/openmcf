# AzureLoadBalancer

Azure Load Balancer is a Layer 4 (TCP/UDP) network load balancer that distributes incoming traffic across healthy backend instances. It operates at the transport layer, making routing decisions based on IP address and port combinations without inspecting application-layer content.

## When to Use

Use AzureLoadBalancer when you need:

- **Layer 4 load balancing** -- TCP/UDP traffic distribution without HTTP inspection
- **High availability** -- automatic health checks remove unhealthy backends from rotation
- **Internal service routing** -- private load balancing within a VNet for multi-tier architectures
- **Public internet ingress** -- internet-facing traffic distribution to backend pools
- **Explicit outbound SNAT** -- deliberately sized egress port budgets through public frontends
- **Port forwarding** -- inbound NAT rules reaching individual instances (single-target or per-pool-member)
- **HA ports** -- forward all ports and protocols (protocol `ALL`) for NVA or SQL AlwaysOn scenarios

For Layer 7 (HTTP/HTTPS) load balancing with path-based routing, SSL termination, or WAF, use **AzureApplicationGateway** instead.

## Key Configuration

### Frontends: public vs internal, per frontend

Each entry in `frontend_ip_configurations` is either public (references an `AzurePublicIp` or an `AzurePublicIpPrefix`) or internal (references an `AzureSubnet`, with an optional pinned private address and availability zones). One load balancer can mix both -- for example a public frontend for ingress and an internal one for east-west traffic. Rules, NAT rules, and outbound rules target a frontend by name; the name may be omitted when exactly one frontend exists.

### SKU and tier

`STANDARD` (the default) is the production SKU: zone redundancy, SLA, outbound rules, HA ports. `GATEWAY` is the niche SKU for chaining network virtual appliances -- its backend pools carry tunnel interfaces, and the subscription needs the `Microsoft.Network/AllowGatewayLoadBalancer` feature registered. Basic is not modeled: Azure retired it in September 2025.

`sku_tier: GLOBAL` creates a cross-region load balancer whose backend members are the frontends of regional load balancers (declared as pool `addresses` referencing regional frontend IDs).

### Backend pools and membership

Pools are named containers. NIC-based membership is expressed **from the member side** -- a network interface's ip_configuration or a virtual machine scale set's network profile references `status.outputs.backend_pool_ids.<pool-name>` -- which is Azure's own attachment model. Vnet-scoped IP-based members (appliances addressed by IP) are declared inline via each pool's `addresses`.

### Health probes

`PROBE_TCP` (default) checks that the port opens; `PROBE_HTTP`/`PROBE_HTTPS` GET a `request_path` and require HTTP 200 -- prefer them when the workload exposes a health endpoint. `probe_threshold` sets how many consecutive successes re-admit a recovered instance.

### Rules, NAT rules, and outbound rules

- **Load-balancing rules** map a frontend port/protocol to a backend pool and port, with `load_distribution` (session persistence), `tcp_reset_enabled`, `floating_ip_enabled` (Direct Server Return), and `disable_outbound_snat`.
- **Inbound NAT rules** forward frontend ports to individual instances: single-target rules are completed by a NIC-side association referencing `status.outputs.nat_rule_ids.<rule-name>`; pool-style rules give every pool member its own frontend port from a range.
- **Outbound rules** configure explicit SNAT through public frontends with a deliberately sized per-instance port budget -- combine with `disable_outbound_snat` on the load-balancing rules that share the pool.

## Outputs

| Output | Description |
|--------|-------------|
| `load_balancer_id` | Azure Resource Manager ID |
| `load_balancer_name` | Load balancer name |
| `private_ip_address` | First internal frontend's private address (empty when all frontends are public) |
| `private_ip_addresses` | All internal frontends' private addresses |
| `frontend_ip_configuration_ids` | Frontend IDs keyed by frontend name (gateway chaining, GLOBAL-tier pools) |
| `backend_pool_ids` | Pool IDs keyed by pool name -- the member-side association seam |
| `probe_ids` | Probe IDs keyed by probe name (scale-set rolling-upgrade health probe) |
| `nat_rule_ids` | Inbound NAT rule IDs keyed by rule name -- the NIC NAT-rule association seam |

## Related Resources

- **AzurePublicIp** / **AzurePublicIpPrefix** -- public frontend addresses
- **AzureSubnet** -- internal frontend placement
- **AzureNetworkInterface** -- joins pools and completes single-target NAT rules from the member side
- **AzureVirtualMachineScaleSet** -- scale-set instances join pools through their network profile
- **AzureApplicationGateway** -- Layer 7 alternative with HTTP routing and WAF
- **AzureDnsRecord** -- DNS records pointing at frontend addresses

---

© Planton. Licensed under [Apache-2.0](https://github.com/plantonhq/planton/blob/main/LICENSE).
