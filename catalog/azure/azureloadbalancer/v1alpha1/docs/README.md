# Azure Load Balancer Deployment: One Resource, Many Moving Parts

## Introduction: The Load Balancer Is a Composite

Azure Load Balancer looks like one resource in the portal and *is* one resource in ARM — but it is a composite. Inside the single `Microsoft.Network/loadBalancers` resource live five kinds of sub-resources, each with its own name, its own properties, and its own ARM ID:

- **Frontend IP configurations** — the addresses that receive traffic
- **Backend address pools** — named containers for the instances that serve it
- **Health probes** — the checks that decide which instances are serving
- **Load-balancing rules** — the frontend-port-to-backend-pool mappings
- **Inbound NAT rules** and **outbound rules** — port forwarding in, explicit SNAT out

None of these sub-resources has a life without the load balancer: a backend pool, probe, or rule exists only inside its parent, and is deployed, versioned, and destroyed with it. That is why they are configured as one unit, in one spec. What *is* independent — **membership** of a pool — is deliberately not here: a network interface or a virtual machine scale set joins a pool by referencing the pool's exported ID from its own side. That is Azure's own attachment model, and this API mirrors it. The name-keyed ID maps in the outputs (`backend_pool_ids`, `nat_rule_ids`, `probe_ids`, `frontend_ip_configuration_ids`) are the seams everything else composes through.

This document explains the load balancer's anatomy, the SKU and tier decisions, the shape of each sub-resource, and the design choices behind the API — including what is deliberately not modeled, and why.

## Layer 4, and Where It Sits

Azure Load Balancer is a Layer 4 (transport layer) traffic distributor: it inspects IP addresses and ports, never HTTP headers or URLs. That makes it the right tool for high-throughput, protocol-agnostic distribution — TCP and UDP services, databases, east-west traffic inside a virtual network, HA ports for appliances — and the wrong tool for anything that needs URL routing, TLS termination, or a WAF (that is Application Gateway's job), global anycast (Front Door), or DNS-level failover (Traffic Manager).

Two properties follow from operating at Layer 4 and matter operationally: the load balancer adds essentially no latency (no parsing, no proxying — it rewrites flows in the fabric), and it is protocol-blind — a "healthy" backend is whatever the health probe says, nothing more.

## SKU: Standard Is the Product, Gateway Is the Niche

The `sku` field carries two values, and an unspecified spec applies STANDARD:

**STANDARD** is the production SKU: zone-redundant frontends, a 99.99% SLA, backend pools up to 1000 instances, HTTPS probes, HA ports, outbound rules, and secure-by-default semantics (traffic is dropped unless a network security group allows it). Virtually every deployment is Standard.

**GATEWAY** is the service-chaining SKU for network virtual appliances. A gateway load balancer sits invisibly in front of *other* load balancers' frontends: traffic addressed to a chained frontend is first steered through the gateway's backend pool of appliances (firewalls, packet inspectors) via VXLAN tunnels, then delivered. Its backend pools must declare **tunnel interfaces** (conventionally an internal/external pair at identifiers 800/801, ports 10800/10801), and its rules may target two pools at once — the dual-tunnel pattern. Consuming a gateway is a one-field affair on the other side: a frontend (or a NIC configuration) sets `gatewayLoadBalancerFrontendIpConfigurationId` to the gateway's frontend ID.

**Basic is not modeled.** Azure retired the Basic SKU in September 2025; there is no production case for it and no migration path worth encoding.

One operational note on Gateway: the SKU requires the `Microsoft.Network/AllowGatewayLoadBalancer` feature registered on the subscription, which Azure grants through a support ticket rather than self-service registration. The API models the SKU fully — spec, validations, and both IaC engines — but live end-to-end verification of Gateway deployments is excluded from the standard test matrix for exactly this reason: the feature gate is per-subscription and ticket-bound.

## Tier: Regional Is the Shape, Global Is the Exception

`skuTier` selects between REGIONAL (the default — a load balancer serving one region, which is virtually every deployment) and GLOBAL — a cross-region load balancer whose backend pool members are not instances but *the frontends of regional load balancers*. A GLOBAL-tier pool declares its members inline: each address sets `loadBalancerFrontendIpConfigurationId` to a regional frontend's ARM ID (available from that load balancer's `frontend_ip_configuration_ids` output). GLOBAL requires the STANDARD SKU — the spec enforces it.

## Frontends: Public vs Internal Is Per-Frontend, Not Per-LB

There is no "internal load balancer" flag. Each frontend IP configuration is independently public or internal, and one load balancer can mix both — a public frontend for ingress and an internal one for east-west traffic on the same resource. The mode is determined by which address source the frontend references, and exactly one may be set:

- **`subnetId`** → an internal frontend: a private address in the referenced `AzureSubnet`, dynamically allocated or pinned via `privateIpAddress` (pin only when DNS, firewall rules, or service discovery hold the address elsewhere — a dynamic address is already stable for the frontend's lifetime). Internal frontends also carry the zone posture (`zones: ["1","2","3"]` for zone redundancy, a single zone to pin) and optionally an address family (`privateIpAddressVersion: IPV6` for dual-stack subnets).
- **`publicIpAddressId`** → a public frontend fronted by a first-class `AzurePublicIp` (Standard SKU), so the address is visible in the resource graph, allowlistable, and reusable. A public frontend's zone posture comes from the referenced public IP, not from the frontend.
- **`publicIpPrefixId`** → a public frontend drawing from an `AzurePublicIpPrefix` — the SNAT-scaling shape, used with outbound rules so downstream partners can allowlist one contiguous CIDR block.

Every rule, NAT rule, and outbound rule targets a frontend **by its `name`**. When the load balancer has exactly one frontend, the name may be omitted and defaults to it — the single-frontend case (most deployments) stays terse, while multi-frontend specs are forced to be explicit. Spec-level validation enforces exactly that boundary.

Two lifecycle facts worth planning around: Azure does not allow removing *all* frontends from an existing load balancer (going from some to none replaces the resource), and changing a frontend's zones replaces that frontend.

## Backend Pools: Declared Here, Joined from Elsewhere

A backend pool in this spec is deliberately thin: a **name**, and optionally the IP-based membership machinery. That thinness is the design.

**NIC-based membership — the common case — lives on the member.** In Azure's model, a network interface joins a pool through an association declared NIC-side; a scale set joins through its network profile. The load balancer never lists its NIC members. This API preserves that: the pool's ARM ID is exported as `status.outputs.backend_pool_ids.<pool-name>`, and an `AzureNetworkInterface`'s ip configuration references it in `loadBalancerBackendAddressPoolIds`. Members join and leave without the load balancer's spec changing — no lifecycle coupling, no shared-ownership fights over one resource.

**IP-based membership is the exception, and it is inline.** Pools whose members are addressed by IP rather than by NIC (appliances, cross-VNet members, GLOBAL-tier pools) declare `addresses` inline, which requires the pool's `virtualNetworkId` (a reference to an `AzureVirtualNetwork`). `synchronousMode` (AUTOMATIC or MANUAL) governs how IP-based members synchronize, exists only for vnet-scoped pools, and is fixed at creation. Each address is either an `ipAddress` (REGIONAL tier) or a `loadBalancerFrontendIpConfigurationId` (GLOBAL tier) — exactly one, spec-enforced.

A frontend-only load balancer with no pools at all is legal — a resource carrying only inbound NAT rules.

## Health Probes: The Only Honest Signal

Probes decide which pool members receive traffic. A probe fires every `intervalInSeconds` (default 15, minimum 5); after `numberOfProbes` consecutive failures (default 2) the instance leaves rotation; after `probeThreshold` consecutive successes (default 1) it returns. `probeThreshold` is the flap dampener — raise it when instances oscillate near the health boundary.

The `protocol` enum carries three values: `PROBE_TCP` (port open means healthy — the default), `PROBE_HTTP`, and `PROBE_HTTPS` (a GET to `requestPath` must return 200 — certificate validity is not checked for HTTPS). Prefer the HTTP variants whenever the workload has a health endpoint: a process can hold a port open while entirely unhealthy. `requestPath` is required for the HTTP variants and forbidden for TCP — the spec enforces the pairing, so a TCP probe with a path (or an HTTP probe without one) fails validation, not deployment.

**Why the `PROBE_` prefix?** Protobuf enum value names are scoped to the *package*, not the enum — two enums in the same package cannot both declare a value named `TCP`. The transport protocol enum (used by rules, NAT rules, and outbound rules) owns the bare `TCP`/`UDP` names, so the probe protocol enum carries the `PROBE_` prefix. Both IaC modules map the full value names to ARM's values internally.

A probe's ARM ID is exported as `status.outputs.probe_ids.<probe-name>` — the seam a virtual machine scale set's rolling-upgrade policy references as its `health_probe_id`.

## Rules: The Routing Table

A load-balancing rule maps a frontend port/protocol to a backend pool and port. Each rule names its frontend (omittable with exactly one frontend), its pool(s) via `backendPoolNames` (one pool on STANDARD; two only on GATEWAY — the dual-tunnel pattern), and optionally its gating probe via `probeName`. A rule without a probe load-balances blindly — legal, but production rules should probe.

The tuning fields:

- **`protocol`** — `TCP`, `UDP`, or `ALL`. `ALL` creates an **HA-ports rule**: every port, both protocols, the appliance/NVA pattern — valid only on an internal STANDARD frontend, with both ports set to 0. The spec enforces the pairing in both directions: ALL requires zero ports, non-ALL requires non-zero ports.
- **`loadDistribution`** — session persistence. `DEFAULT` (5-tuple hash, effectively no persistence, best distribution) is right for stateless backends; `SOURCE_IP` (2-tuple) and `SOURCE_IP_PROTOCOL` (3-tuple) pin a client to one backend for legacy stateful workloads.
- **`idleTimeoutInMinutes`** (4–100, default 4) and **`tcpResetEnabled`** — raise the timeout for long-lived connections (WebSockets, database pools), and turn on TCP reset so dropped connections fail fast instead of dying silently on the next write. Azure defaults reset to off; production TCP rules generally want it on.
- **`floatingIpEnabled`** — Direct Server Return: the backend sees the *frontend's* IP as the destination. Required for SQL Server AlwaysOn listeners and some clustering schemes; the backend OS must carry a matching loopback.
- **`disableOutboundSnat`** — turns off the implicit outbound SNAT this rule would otherwise grant its pool. Set it whenever the pool's egress is handled properly — by an explicit outbound rule (required to combine the two on one pool) or a NAT gateway on the subnet — because implicit SNAT has a small, exhaustion-prone port budget.

## Inbound NAT Rules: Port Forwarding, Two Modes

An inbound NAT rule forwards a frontend port to an individual instance — port forwarding, not load balancing. One message models both of Azure's modes, and each rule is exactly one of them (spec-enforced XOR):

**Single-target mode** sets `frontendPort` (and `backendPort`). The rule itself names no instance — the attachment is completed from the member side, by a NIC's NAT-rule association referencing `status.outputs.nat_rule_ids.<rule-name>`. The same member-side seam as pool membership, for the same reason: which instance receives the forward is the member's declaration, not the load balancer's.

**Pool-style mode** sets `backendPoolName` plus `frontendPortStart`/`frontendPortEnd`: every member of the pool automatically gets its own dedicated frontend port from the range — per-instance SSH across a scale set is the canonical use. The range must be at least as large as the pool's member count.

Pool-style NAT rules are the modern replacement for the legacy **NAT pool** mechanism (`azurerm_lb_nat_pool` / `inboundNatPools`), which is deliberately not modeled: it predates pool-style NAT rules, exists only for scale-set port ranges, and offers nothing the pool-style rule does not do better on the current API surface.

## Outbound Rules: Explicit SNAT

By default, a backend instance's internet egress rides on implicit SNAT with a small shared port budget — the classic source of mysterious connection failures at scale. An outbound rule makes egress explicit: which public frontends the pool's traffic is SNATed through (`frontendIpConfigurationNames` — more frontends multiply the port budget), and how many SNAT ports each instance gets (`allocatedOutboundPorts`, default 1024).

Sizing note: each frontend IP carries a budget of 64,000 ports. Setting `allocatedOutboundPorts: 0` lets Azure divide the budget evenly across the pool's current size — convenient, but reallocation churns connections whenever the pool scales. Production pools should size explicitly (budget ÷ maximum instances, in multiples of 8).

Outbound rules are legal only with public frontends on the STANDARD SKU, and combining an outbound rule with load-balancing rules on the same pool requires `disableOutboundSnat: true` on those rules. Pools that egress heavily may prefer a NAT gateway on the subnet instead — the outbound rule is the load-balancer-native alternative, strongest when the egress addresses must be the load balancer's own frontends (or a prefix's known CIDR).

## Validation That Mirrors ARM

The spec encodes the contracts Azure Resource Manager (and the provider) would otherwise enforce at deploy time, so a manifest that cannot deploy fails at validation with a message naming the rule:

- **Frontend shape** — at most one address source per frontend (subnet / public IP / prefix); a pinned private address requires a subnet; zones only on internal frontends.
- **SKU pairings** — GLOBAL tier requires STANDARD SKU; GATEWAY requires tunnel interfaces on *every* pool, and tunnel interfaces are forbidden off GATEWAY; a two-pool rule is GATEWAY-only.
- **HA ports** — protocol ALL pairs with zero frontend/backend ports, and only with them.
- **NAT rule modes** — single-target (`frontendPort`) XOR pool-style (`backendPoolName`), and pool-style requires both range bounds while single-target forbids them.
- **Probe pairing** — `requestPath` required for PROBE_HTTP/PROBE_HTTPS, forbidden for PROBE_TCP.
- **Pool addresses** — each inline address is an IP member XOR a regional-frontend member; IP members require the pool's virtual network.
- **Referential integrity by name** — every rule's `backendPoolNames`, `probeName`, and frontend name; every NAT rule's frontend and pool; every outbound rule's frontends and pool — all must match sub-resources actually declared in the spec, and a frontend name may be omitted only when exactly one frontend exists. A typo'd cross-reference is a validation error, never a failed deploy.

## The Planton Approach

Planton provides a declarative, protobuf-based API for the Azure Load Balancer. Two design decisions define its shape.

### Bundled Sub-Resources, Member-Side Membership

The composite is one spec because its parts share one lifecycle — but the spec stops exactly where Azure's ownership model stops. Pool membership, NAT-rule attachment, and probe consumption belong to the *members*, and the outputs are built for it: four **name-keyed ID maps** let any consumer reference a sub-resource by the name the spec gave it, without parsing ARM IDs or knowing the load balancer's internals.

| Output | Keyed by | Referenced by |
|--------|----------|---------------|
| `backend_pool_ids` | pool name | a NIC's `loadBalancerBackendAddressPoolIds`, a scale set's network profile |
| `nat_rule_ids` | NAT rule name | a NIC's `loadBalancerInboundNatRuleIds` (single-target rules) |
| `probe_ids` | probe name | a scale set's rolling-upgrade `health_probe_id` |
| `frontend_ip_configuration_ids` | frontend name | gateway chaining, GLOBAL-tier pool addresses |

Alongside the maps: `load_balancer_id`, `load_balancer_name`, `private_ip_address` (the first internal frontend's address — what internal DNS records point at), and `private_ip_addresses` (all internal frontends, in declaration order). A public frontend's address lives on its referenced `AzurePublicIp`, where it always was.

### References Where Kinds Exist, ARM IDs Where They Don't

The resource group, subnets, public IPs, public IP prefixes, and pool virtual networks are first-class references to their Planton kinds. Gateway-frontend chaining IDs and GLOBAL-tier regional-frontend IDs are plain ARM ID strings — they address *sub-resources* of another load balancer, which the reference machinery reaches through that load balancer's `frontend_ip_configuration_ids` map output.

### Example: Public Web Tier

One public frontend, one pool, an HTTP probe, and a rule — with egress made explicit:

```yaml
apiVersion: azure.planton.dev/v1alpha1
kind: AzureLoadBalancer
metadata:
  name: web-lb
spec:
  region: eastus
  resourceGroup:
    valueFrom:
      name: prod-rg
  name: web-lb
  frontendIpConfigurations:
    - name: public
      publicIpAddressId:
        valueFrom:
          name: web-lb-ip
  backendPools:
    - name: web
  healthProbes:
    - name: http-health
      protocol: PROBE_HTTP
      port: 8080
      requestPath: /healthz
  rules:
    - name: https
      protocol: TCP
      frontendPort: 443
      backendPort: 8443
      backendPoolNames: [web]
      probeName: http-health
      tcpResetEnabled: true
      disableOutboundSnat: true
  outboundRules:
    - name: web-egress
      frontendIpConfigurationNames: [public]
      backendPoolName: web
      protocol: ALL
      allocatedOutboundPorts: 4000
```

The rule omits `frontendIpConfigurationName` — with exactly one frontend, it defaults. Members join by referencing `status.outputs.backend_pool_ids.web` from their NICs.

### Example: Internal LB with Zone Redundancy and Admin NAT

```yaml
apiVersion: azure.planton.dev/v1alpha1
kind: AzureLoadBalancer
metadata:
  name: app-ilb
spec:
  region: eastus
  resourceGroup:
    valueFrom:
      name: prod-rg
  name: app-ilb
  frontendIpConfigurations:
    - name: internal
      subnetId:
        valueFrom:
          name: app-subnet
      privateIpAddress: 10.0.2.10
      zones: ["1", "2", "3"]
  backendPools:
    - name: app
  healthProbes:
    - name: tcp-health
      port: 5432
  rules:
    - name: postgres
      protocol: TCP
      frontendPort: 5432
      backendPort: 5432
      backendPoolNames: [app]
      probeName: tcp-health
      idleTimeoutInMinutes: 30
      tcpResetEnabled: true
  natRules:
    - name: ssh-admin
      protocol: TCP
      frontendPort: 2222
      backendPort: 22
```

The pinned frontend address is what DNS records key on; the zone list makes it survive a zone outage. The single-target NAT rule declares the forward — the admin host's NIC completes it by referencing `status.outputs.nat_rule_ids.ssh-admin`.

## Common Anti-Patterns to Avoid

**❌ Anti-Pattern 1: Rules Without Probes**

A rule with no `probeName` distributes traffic to every pool member, dead or alive.

**✅ Solution:** Every production rule references a probe — and an HTTP probe over a TCP probe wherever a health endpoint exists, because an open port is not health.

---

**❌ Anti-Pattern 2: Relying on Implicit SNAT at Scale**

Pools that egress through a load-balancing rule's implicit SNAT share a small port budget; exhaustion surfaces as intermittent, hard-to-attribute connection failures.

**✅ Solution:** An explicit outbound rule with sized `allocatedOutboundPorts` (plus `disableOutboundSnat` on the pool's rules), or a NAT gateway on the subnet.

---

**❌ Anti-Pattern 3: Pinning Private Addresses by Reflex**

Dynamic frontend addresses are already stable for the frontend's lifetime. A pinned address is one more manually-managed value that can collide.

**✅ Solution:** Pin only when something *outside* the platform holds the address — external DNS, partner firewall rules, appliance peering. Otherwise let Azure allocate and read the address from `private_ip_address`.

---

**❌ Anti-Pattern 4: Declaring Membership on the Wrong Side**

Trying to enumerate a pool's members in the load balancer spec couples every scale event to this resource and fights Azure's model.

**✅ Solution:** NIC-based members reference `status.outputs.backend_pool_ids.<name>` from their own specs; only IP-based members (appliances, GLOBAL-tier regional frontends) are declared inline, via `addresses`.

---

**❌ Anti-Pattern 5: Skipping the Zone Posture on Internal Frontends**

An internal frontend with no `zones` has no zone guarantee — one zone outage can take the address down while the backends survive.

**✅ Solution:** `zones: ["1", "2", "3"]` on internal frontends is the production default posture. (Public frontends inherit theirs from the referenced public IP.)

## Conclusion: Declare the Topology, Reference the Seams

The load balancer spec is a routing topology: addresses in, pools and probes in the middle, rules wiring them together, NAT and outbound rules for the traffic that doesn't fit the load-balanced shape. Everything with the load balancer's lifecycle lives in the one spec; everything with its own lifecycle — membership — references the name-keyed output maps from its own side, exactly as Azure intends.

Declare the topology once, validate it before it deploys, and let the members come and go without touching it.
