# AWS Network Load Balancer: Design and Research Notes

This document is the architecture-focused reference for the AwsNlb component:
Layer-4 load balancing fundamentals, NLB-specific behavior and constraints,
and how Planton models the NLB as one of a set of composable ELBv2 kinds.

---

## Layer 4 vs Layer 7

| Layer | Name | What It Inspects | Examples |
|-------|------|------------------|----------|
| **Layer 4** | Transport | IP, port, protocol (TCP/UDP) | NLB, HAProxy (TCP mode) |
| **Layer 7** | Application | HTTP headers, path, host, body | ALB, Nginx, HAProxy (HTTP mode) |

**NLB is Layer 4**: it forwards connections based solely on destination port
and protocol. It does not inspect HTTP headers, paths, or request content. A
client connecting on port 443 gets forwarded to a target — the NLB does not
care whether the bytes are HTTPS, TLS-wrapped gRPC, or a proprietary
protocol.

The implications shape the whole component model:

- **No content-based routing**: there is nothing to write a rule against, so
  `AwsLbListenerRule` does not apply to NLBs. Routing is purely by
  port/protocol — one listener, one forward.
- **No HTTP-level features**: no redirects, no fixed responses, no
  authentication actions, no WAF at the NLB level, no Lambda targets. NLB
  listeners forward, and AWS rejects every other action type.
- **Ultra-low latency and flow-based scaling**: no HTTP parsing means minimal
  overhead; the NLB scales on connection flows (new/active connections), not
  requests, and handles millions of them.

---

## How Planton Models It: Composable Kinds

The load balancer carries no routing configuration — that is deliberate, and
it mirrors how AWS models ELBv2:

- **`AwsNlb`** (this component) owns what is truly load-balancer-wide: node
  placement with optional static IPs, optional security groups, and traffic
  distribution behavior.
- **`AwsLbListener`** attaches by ARN
  (`status.outputs.load_balancer_arn`) and owns a port, a protocol
  (`TCP`/`UDP`/`TCP_UDP`/`TLS`), TLS material (certificate, security policy,
  ALPN) for TLS listeners, and a TCP idle timeout. The action set is
  forward-only — exactly the split AWS enforces at Layer 4.
- **`AwsLbTargetGroup`** owns the destination: target type (`instance`,
  `ip`, `alb`), health checks, deregistration/draining, client-IP
  preservation, Proxy Protocol v2, and connection termination.

The split matters because the pieces change at different speeds and are
owned by different teams. The NLB — and the Elastic IPs partners have
allowlisted against it — is the most static artifact in the request path;
listeners change with certificates; target groups change with the services
behind them. TLS termination versus passthrough, health-check tuning, and
client-IP preservation are all decisions that live on the listener and
target group, not here.

An NLB with no listeners is a valid intermediate state: provision the load
balancer, hand its ARN to whoever adds ports.

---

## Static IP Architecture (Subnet Mappings)

### Why Static IPs Matter

ALBs expose only a dynamic DNS name; the underlying addresses change with
scaling and maintenance. For partner allowlisting, firewall rules, legacy
integrations, and DNS pinning, **static IPs are a hard requirement** — and
they are the primary reason architectures choose NLB.

### The Model

`subnetMappings` is the placement primitive: each mapping pins one NLB node
to a subnet (one per Availability Zone), and optionally assigns it a static
address:

- **Internet-facing**: set `allocationId` per mapping — an Elastic IP,
  typically referenced from an `AwsElasticIp` resource via `valueFrom`
  (`status.outputs.allocation_id`). The node's public IP then survives every
  scaling event, instance replacement, and AWS maintenance window. At most
  one Elastic IP per mapping.
- **Internal**: optionally set `privateIpv4Address` per mapping to pin the
  node to a specific address inside the subnet's CIDR — useful when
  downstream systems reference the NLB by IP. When omitted, AWS assigns one.

At least one mapping is required; AWS recommends two or more for zonal
redundancy.

### Add-Only Semantics

AWS allows **adding** subnet mappings to an existing NLB but not removing
them. Start with the minimum set you need and add AZs deliberately —
retreating from one means replacing the load balancer.

---

## Traffic Distribution

### Cross-Zone Load Balancing

Unlike ALB (always on), NLB **defaults cross-zone distribution off**, and the
spec keeps it an explicit opt-in (`crossZoneLoadBalancingEnabled`) because it
is a real cost decision: with cross-zone off, NLB-to-target traffic stays
inside each AZ and costs nothing extra; with it on, an NLB node in AZ-1 may
forward to a target in AZ-2 and the inter-AZ transfer is billed. Enable it
when target distribution across AZs is uneven and per-target balance matters
more than the transfer cost.

### DNS Client Routing Policy

`dnsRecordClientRoutingPolicy` controls how the NLB's DNS name resolves for
clients across AZs:

- **`any_availability_zone`** (default): clients may be routed to any AZ —
  maximum spillover capacity.
- **`availability_zone_affinity`**: clients resolve to the node in their
  resolver's AZ — lowest latency and cross-zone traffic, best when targets
  are evenly distributed.
- **`partial_availability_zone_affinity`**: 85% of queries stay in the
  resolver's AZ, 15% spill elsewhere — a hedge between the two.

### Zonal Shift

`zonalShiftEnabled` allows Amazon Application Recovery Controller to shift
this NLB's traffic away from an impaired Availability Zone without a deploy —
cheap insurance for multi-AZ deployments.

---

## Security Groups: Optional, and a One-Way Door

Unlike ALB (where security groups are effectively required), an NLB can run
**without** security groups, accepting all traffic on its listener ports.
That is a reasonable posture for internal NLBs in trusted network segments
where NACLs and routing provide isolation.

The critical constraint: **once security groups are attached, they can never
be fully removed** — at least one must remain for the life of the load
balancer. You can swap groups, never retreat to "none". Attach groups only
when committed to operating them long-term; the decision is effectively
create-time.

### PrivateLink Enforcement

`enforceSecurityGroupInboundRulesOnPrivateLinkTraffic` (`on`/`off`) decides
whether inbound security-group rules are evaluated for traffic arriving
through PrivateLink VPC endpoints. AWS defaults it to `on` for NLBs created
with security groups. Set `off` when the endpoint service should admit any
consumer the endpoint policy allows, regardless of the NLB's own group
rules; keep `on` for defense in depth when exposing services via
PrivateLink. It is only meaningful when security groups are attached.

---

## Observability: TLS-Only Access Logs

NLB access logs (`accessLogs`, delivered to S3) capture **TLS-listener
traffic only** — an AWS limitation, not a component choice. Plain TCP and
UDP flows are never logged; if you need flow-level visibility for those,
VPC Flow Logs are the tool. As with the ALB, presence of the block implies
enabled, and the bucket must carry the regional ELB log-delivery bucket
policy — delivery fails silently otherwise.

---

## DNS and Alias Records

AWS assigns a DNS name like `edge-nlb-abc123.elb.us-west-2.amazonaws.com`,
resolving to one address per AZ; with Elastic IPs those addresses are stable.
When `dns.enabled` is true, the component creates Route53 **alias** A records
pointing each hostname at the NLB. Alias records rather than CNAMEs because
they work at the zone apex, cost nothing per query, and inherit the NLB's
health. Both engines create the records with
`evaluate_target_health = false`: target-health evaluation only changes
behavior under failover/weighted routing policies, and a simple alias should
not pay for health evaluation.

---

## Naming and Lifecycle

The NLB name comes from `metadata.name`; AWS caps load balancer names at 32
characters and both IaC modules truncate deterministically, reporting the
final name in the `load_balancer_name` output. The scheme (`internal`) is
immutable — changing it replaces the load balancer. Most other attributes
update in place. `deleteProtectionEnabled` is recommended for production for
two reasons: deleting an NLB silently orphans every listener attached to it,
and any Elastic IPs pinned to it start billing as unattached addresses.

---

## The NLB-in-Front-of-ALB Pattern

```
Internet → NLB (static IPs, Layer 4) → ALB (Layer 7 routing) → services
```

When an architecture needs both static IPs and content-based routing,
compose the two: an `AwsNlb` with Elastic IPs, an `AwsLbListener` on the NLB
forwarding to an `AwsLbTargetGroup` with `targetType: alb` whose registered
target is the ALB, and the usual listener/rule/target-group graph on the
`AwsAlb` behind it. Every hop in that chain is an explicit reference between
components.

---

## Deliberate Scoping Omissions

The spec deliberately does not model two subnet-mapping surfaces. They are
niche enough that carrying them would cost every reader comprehension for
capabilities almost no architecture pulls; they are deferred until real
architectures pull them:

- **IPv6 source-NAT prefixes**: per-mapping `/80` prefixes used when
  dual-stack NLBs front IPv4-only targets through UDP — a corner of the
  dual-stack story with its own addressing constraints.
- **Secondary private IPs per subnet**: additional per-node private
  addresses beyond the primary (used by some appliance-style deployments).

---

## Summary: Key Design Decisions

| Topic | NLB Behavior | Recommendation |
|-------|--------------|----------------|
| **Static IPs** | Elastic IP per subnet mapping | Reference `AwsElasticIp` resources for allowlisting, firewall rules, DNS pinning |
| **Routing** | Port/protocol only; no rules | Attach `AwsLbListener` (forward-only) and `AwsLbTargetGroup`; per-request routing belongs on an ALB |
| **Cross-zone** | Off by AWS default; billed inter-AZ when on | Enable only when target distribution is uneven |
| **Security groups** | Optional; last one can never be removed | Attach only when committed; treat as create-time |
| **PrivateLink** | SG rules enforced on endpoint traffic by default | Keep `on` for defense in depth |
| **Access logs** | TLS listeners only | Use VPC Flow Logs for TCP/UDP visibility |
| **Subnet mappings** | Add-only | Start minimal; add AZs deliberately |
| **Scheme** | Immutable | Changing `internal` replaces the load balancer |

---

## Dual-Engine Implementation

`AwsNlb` ships both a Terraform/OpenTofu module and a Pulumi (Go) module at
behavioral parity. Both create the load balancer with identity tags, apply
the same 32-character name truncation, send only explicitly set attributes
(so AWS defaults survive), treat access-log presence as enabled, create the
same alias records with the same `evaluate_target_health = false` decision,
and export the same four outputs (`load_balancer_arn`,
`load_balancer_name`, `load_balancer_dns_name`,
`load_balancer_hosted_zone_id`). Whichever engine a team standardizes on,
the NLB behaves identically.
