# AWS Application Load Balancer: Design and Research Notes

## Introduction

The AWS Application Load Balancer (ALB) sits at a critical junction in modern
cloud architectures: it is the entry point for user traffic, the routing
substrate for microservices, and the TLS termination point that protects
application backends. Despite its ubiquity, ALB deployment remains
surprisingly error-prone when managed manually, with common misconfigurations
that undermine high availability, security, and operational stability.

This document covers what makes ALBs architecturally interesting (the split
between the load balancer and its routing graph), how the deployment
landscape evolved from console clicks to control planes, and how Planton
models the ALB as one of four composable components. The core position:
production ALBs should never be created through the AWS Console. They should
be declared as code, version-controlled, and continuously reconciled by a
control plane.

## The ALB Deployment Landscape

ALB management spans a spectrum from fully manual to continuously reconciled.

### Level 0: Manual Provisioning via AWS Console (The Anti-Pattern)

The console wizard is approachable for learning and **highly susceptible to
critical configuration errors** in production:

1. **Single-AZ deployment**: the wizard allows selecting multiple subnets in
   the *same* Availability Zone, creating an ALB with zero zonal redundancy —
   defeating the primary purpose of load balancing.
2. **Misconfigured security groups**: overly permissive rules, missing the
   security-group-chaining pattern between the ALB and its targets, or
   allowing direct internet access to targets that bypasses the load balancer
   entirely.
3. **Invisible drift**: every hand-edit to an attribute (idle timeout, desync
   mode, logging) exists only in the console's memory. Six months later,
   nobody knows why staging and production behave differently.

**Verdict**: acceptable for learning. Unacceptable for anything that must be
reproducible, auditable, or secure.

### Level 1: Scripted Provisioning with AWS CLI and SDKs

Imperative scripting (`aws elbv2 create-load-balancer`, Boto3) makes the
sequence explicit — and makes the script author responsible for ordering,
idempotency, and error recovery across `create-load-balancer`,
`modify-load-balancer-attributes`, and the DNS calls. The API documentation
does codify one requirement the console lets you miss: **subnets from at
least two Availability Zones are mandatory**. Suitable for one-off
automation; not for declarative, state-managed infrastructure.

### Level 2: Infrastructure as Code (Terraform, CloudFormation, CDK, Pulumi)

The dominant paradigm. Declarative tools model the ALB estate as granular
resources — `aws_lb`, `aws_lb_listener`, `aws_lb_listener_rule`,
`aws_lb_target_group` — track state, and compute minimal diffs. That
granularity is the important lesson: **AWS itself models the load balancer
and its routing graph as separate resources with separate lifecycles**, and
every serious IaC tool preserves that split rather than flattening it.

### Level 3: Control Planes

A CLI-based tool runs, applies, and exits. A control plane continuously
observes infrastructure and reconciles it against declared intent — the model
behind the AWS Load Balancer Controller for Kubernetes and Crossplane, and
the model Planton's protobuf-defined APIs are built for.

## How Planton Models It: Four Composable Kinds

The essential design decision is that **the load balancer carries no routing
configuration**. Planton splits the ELBv2 surface exactly where AWS does:

- **`AwsAlb`** (this component): placement (subnets, scheme), security
  groups, and the HTTP behavior knobs AWS models as load balancer
  attributes.
- **`AwsLbListener`**: attaches to the ALB by ARN
  (`status.outputs.load_balancer_arn`); owns one port, its TLS certificates
  and policy, and the default action taken when no rule matches.
- **`AwsLbListenerRule`**: attaches to a listener; owns one service's routing
  conditions (host, path, headers) and its action, at an explicit priority.
- **`AwsLbTargetGroup`**: the destination — targets, health checks, draining,
  stickiness.

Why this split matters in practice:

- **Change frequency**: the ALB changes rarely (placement is nearly
  immutable, attributes occasionally). Listeners change with certificate
  rotations. Rules and target groups change with every service deployment.
  Bundling them means every service deploy touches the shared load balancer's
  manifest — the highest-blast-radius resource in the request path.
- **Ownership**: a platform team typically owns the ALB and its listeners;
  service teams own their rules and target groups. Separate kinds give each
  artifact a natural owner and review path.
- **Reference topology**: rules need a listener ARN, listeners need the load
  balancer ARN, forward actions need target group ARNs. With separate kinds,
  each edge is an explicit `valueFrom` reference and the architecture graph
  tells the truth.

An ALB with no listeners is a valid, useful intermediate state: provision the
load balancer, hand its ARN to the teams that add ports and routes.

## What the ALB Itself Owns

### Placement

- **`subnets` (2+ required)**: the AWS minimum for an ALB, and the mechanism
  for zonal redundancy. Public subnets for internet-facing ALBs, private for
  internal ones. The spec enforces `min_items = 2` at the API layer so a
  single-AZ ALB cannot be declared.
- **`internal`**: the scheme. Internet-facing ALBs get public DNS; internal
  ALBs are reachable only inside the VPC. The scheme is immutable in AWS —
  changing it replaces the load balancer, which is why the field is
  documented as a replacement trigger rather than an in-place edit.
- **`ipAddressType`**: `ipv4`, `dualstack`, or
  `dualstack-without-public-ipv4`. The last is the interesting one: public
  IPv6 with private IPv4 avoids AWS's public-IPv4 charges for clients that
  can reach IPv6, while CloudFront or dual-stack clients still work.

### Security Groups

Security groups on an ALB are effectively mandatory for production even
though AWS will silently attach the VPC's default group when none is given.
The production pattern is two groups: an ALB group admitting client traffic
on listener ports, and a target group admitting the application port *only
from the ALB's group* — the chaining that guarantees no traffic bypasses the
load balancer. Both are modeled by `AwsSecurityGroup` and referenced here.

### HTTP Behavior Attributes

These are the knobs AWS models as `modify-load-balancer-attributes`, and they
update in place:

- **`idleTimeoutSeconds`** (1–4000, AWS default 60): raise above the
  application's slowest response to avoid 504s on long-running requests, and
  keep it above upstream keep-alive intervals so the ALB is never the one
  closing first.
- **`clientKeepAliveSeconds`** (60–604800, AWS default 3600): how long a
  client connection may live across requests.
- **`http2Enabled`**: modeled as an *optional* bool so that an explicit
  `false` ("downgrade clients to HTTP/1.1") is distinguishable from unset
  ("keep the AWS default of true"). Plain proto3 bools cannot express that
  tri-state.
- **`desyncMitigationMode`** (`monitor`/`defensive`/`strictest`): protection
  against HTTP request-smuggling. `defensive` (the AWS default) blocks
  ambiguous requests likely to poison caches; `strictest` blocks everything
  not RFC 7230 compliant — the right choice when clients are known-good
  (internal tools, API-only traffic).
- **`dropInvalidHeaderFields`**: drop rather than forward header names that
  are not valid HTTP tokens; pairs with desync mitigation as
  header-smuggling hardening.
- **`xffHeaderProcessingMode`** (`append`/`preserve`/`remove`) and
  **`xffClientPortEnabled`**: X-Forwarded-For control. `preserve` and
  `remove` matter when the ALB sits behind another proxy layer (CloudFront,
  an edge WAF) whose XFF chain must win.
- **`preserveHostHeader`**: forward the client's original Host header
  untouched — required by applications that generate absolute URLs or do
  their own virtual-host routing.
- **`tlsVersionAndCipherSuiteHeadersEnabled`**: inject `x-amzn-tls-version`
  and `x-amzn-tls-cipher-suite` toward targets, for applications that audit
  their clients' TLS posture.
- **`wafFailOpenEnabled`**: what happens when an attached WAF is unreachable.
  Fail-open favors availability; fail-closed (the AWS default) favors
  security. This is a deliberate business decision, so the spec exposes it
  rather than picking.
- **`zonalShiftEnabled`**: allows Amazon Application Recovery Controller to
  drain an impaired Availability Zone without a deploy.

The IaC modules send only explicitly set attributes to AWS, so everything
left unset keeps its AWS default rather than a module opinion — and AWS
default changes flow through without a code change.

### Log Delivery

The ALB has three distinct S3 log streams, and the spec models them with one
shared shape (`bucket` + `prefix`):

- **Access logs**: one entry per request — the raw material for Athena
  queries during an incident (`WHERE elb_status_code >= 500`).
- **Connection logs**: one entry per client connection, with TLS handshake
  details. This is the only place TLS negotiation failures appear, because a
  failed handshake never becomes a request.
- **Health-check logs**: one entry per probe result, for debugging flapping
  targets without packet captures.

Presence of the block implies enabled — a bucket with logging off is
meaningless, so there is no separate `enabled` flag. The bucket must carry
the regional ELB log-delivery bucket policy; without it, delivery fails
silently (an AWS behavior worth knowing before an incident, not during).

### DNS

`dns` creates Route53 **alias** A records for each hostname. Alias records
rather than CNAMEs because they work at the zone apex, cost nothing per
query, and inherit the ALB's health. The records are created with
`evaluate_target_health = false` in both engines: target-health evaluation
only changes behavior under failover/weighted routing policies, and a simple
alias should not pay for health evaluation. Architectures that need
DNS-level failover own their records through `AwsRoute53DnsRecord` with the
policy they intend.

### Naming and Lifecycle

The ALB name comes from `metadata.name`; AWS caps load balancer names at 32
characters and both IaC modules truncate deterministically, reporting the
final name in the `load_balancer_name` output. Scheme and (practically)
subnet topology are identity; everything else — attributes, security groups,
logs, DNS — updates in place. `deleteProtectionEnabled` is recommended for
production because deleting an ALB silently orphans every listener and rule
attached to it.

## Deliberate Scoping Omissions

The spec deliberately does not model several ALB surfaces. These are niche
enough that carrying them would cost every reader comprehension for
capabilities almost no architecture pulls; they are deferred until real
architectures pull them:

- **`subnet_mapping`-style address pinning**: ALBs accept per-subnet mappings
  (notably private IPv4 pinning for internal load balancers), but static
  addressing is the NLB's story — architectures that need pinned addresses
  put an NLB in front. The plain `subnets` list covers the ALB case.
- **IPAM pools**: drawing public IPv4 addresses from a BYOIP IPAM pool.
- **Customer-owned IPv4 pools (Outposts)**: on-premises ALB placement.
- **Minimum load-balancer capacity reservations**: pre-provisioned LCU
  capacity for known traffic spikes, a paid reservation surface with its own
  lifecycle.

## Cost Notes: LCUs and Rule Evaluations

ALB pricing is a fixed hourly charge plus Load Balancer Capacity Units,
billed on the highest of four dimensions per hour: new connections (25/s),
active connections (3,000/min), processed bytes (1 GB/h), and rule
evaluations (1,000/s, first 10 rules free). The non-obvious dimension is rule
evaluations: a gateway ALB with 50 path-based rules processes all of them for
many requests, so a complex routing graph can cost more than the same traffic
on two simpler ALBs. This is a reason the composable design matters — rules
are visible, countable artifacts rather than lines buried in a monolithic
manifest.

## ALB vs. NLB

| Aspect | ALB (Layer 7) | NLB (Layer 4) |
|--------|---------------|----------------|
| Protocols | HTTP, HTTPS, WebSockets, gRPC | TCP, UDP, TLS, TCP_UDP |
| Routing | Host, path, headers, query — via listener rules | Port/protocol only — no rules |
| Static IPs | No (dynamic DNS) | Yes — Elastic IP per subnet mapping |
| Security groups | Effectively required | Optional (one-way door once attached) |
| Typical use | Web apps, APIs, microservices | Non-HTTP traffic, static-IP allowlisting, PrivateLink |

Both are fronted by the same `AwsLbListener` and `AwsLbTargetGroup`
components; the NLB counterpart to this component is `AwsNlb`. The classic
combination — an NLB with static IPs forwarding to an ALB for Layer-7 routing
— is expressed as an `AwsLbTargetGroup` with `targetType: alb`.

## Dual-Engine Implementation

`AwsAlb` ships both a Terraform/OpenTofu module and a Pulumi (Go) module at
behavioral parity. Both create the load balancer with identity tags, apply
the same 32-character name truncation, send only explicitly set attributes,
treat log-block presence as enabled, create the same alias records with the
same `evaluate_target_health = false` decision, and export the same four
outputs (`load_balancer_arn`, `load_balancer_name`,
`load_balancer_dns_name`, `load_balancer_hosted_zone_id`). Whichever engine
a team standardizes on, the ALB behaves identically.
