# AWS LB Target Group: The Composition Point of ELBv2

## What a Target Group Is

An ELBv2 target group is the set of backends a load balancer sends traffic
to, together with the policy for deciding whether each backend may receive
it. It is not a detail of the load balancer: a target group has its own ARN,
its own lifecycle, and a many-to-many relationship with the routing layer.
One group can receive traffic from several listeners; one listener can spread
traffic across several weighted groups; an ECS service or auto-scaling group
registers targets into a group without knowing anything about the listeners
in front of it.

That independence is what makes the target group the composition point of AWS
load balancing. The load balancer and its listeners are shared,
slow-changing infrastructure; target groups come and go with the services
behind them. `AwsLbTargetGroup` models the group as its own component so that
deploying a service means creating a group and a rule -- never editing the
shared load balancer.

## One Kind for Both Families

The same component serves Application and Network Load Balancers, because
that is exactly how AWS models ELBv2: a target group's `protocol` decides the
family it can attach to (`HTTP`/`HTTPS` for ALB; `TCP`/`UDP`/`TCP_UDP`/`TLS`
for NLB), and the family decides which tuning fields apply. Splitting the
kind by family would duplicate the shared 80% of the surface (health checks,
stickiness, deregistration, target registration) and invent a distinction AWS
itself does not draw. Field documentation and validations mark every
family-specific field instead:

- **ALB only**: `protocolVersion` (HTTP1/HTTP2/GRPC), `slowStartSeconds`,
  `loadBalancingAlgorithmType`, `loadBalancingAnomalyMitigation`, cookie
  stickiness.
- **NLB only**: `preserveClientIp`, `proxyProtocolV2`,
  `connectionTermination`, `targetHealthState`, source-IP stickiness.
- **Lambda only**: `lambdaMultiValueHeadersEnabled`; port, protocol, and VPC
  do not apply at all.

Gateway Load Balancer (GENEVE) target groups are deliberately not modeled:
there is no gateway-load-balancer kind to compose them with, so the GENEVE
surface would be dead weight. If a gateway kind ever exists, GENEVE support
belongs in this component's scope. The target-control port (agent-based
target optimization) and QUIC server IDs on registrations are likewise not
modeled -- both serve niche protocol surfaces with no composition story here
yet, and both are additive when real architectures pull for them.

## Why Target Registrations Fold Into the Group

AWS exposes target registration as a separate attachment resource, and most
IaC tools mirror that. This component folds static registrations into the
group's `targets` list instead, because a registration is pure glue: it has
no configuration beyond "this ID, this port", no referenceable identity, and
no independent lifecycle worth tracking. A separate registration kind would
flood the resource graph with nodes nothing can reference.

The fold only covers *static* registration -- standalone EC2 instances, fixed
IPs, a Lambda function, an inner ALB. Dynamic registrars (ECS services,
auto-scaling groups, Kubernetes load-balancer controllers) register and
deregister their own targets at runtime and should never be listed in the
spec; most architectures leave `targets` empty.

## Health Checks Are Independent of Traffic

The health-check protocol is deliberately decoupled from the traffic
protocol, and using that decoupling is usually the right call for NLB
groups: a TCP check only proves the port is open, while an HTTP check against
a readiness endpoint proves the application can serve. A TCP target group for
a database or gRPC service typically runs `healthCheck.protocol: HTTP`
against an admin port.

Constraints worth knowing:

- TCP health checks take no `path` and no `matcher` (nothing to match).
- HTTP/HTTPS traffic protocols do not accept TCP health checks.
- gRPC groups (`protocolVersion: GRPC`) match gRPC status codes (default
  `0`), not HTTP codes, and health-check paths are fully-qualified method
  names.
- Health checks can only be disabled for `lambda` groups; every other type
  requires them.

## Immutability and the 32-Character Name

Name, `port`, `protocol`, `protocolVersion`, `vpcId`, `targetType`, and
`ipAddressType` are create-only in AWS: changing any of them replaces the
target group, and the IaC engine re-creates dependent listener attachments in
the same operation. Treat those fields as identity, not configuration --
plan protocol or target-type changes as a second group plus a weighted
traffic shift rather than an in-place edit.

The group name comes from `metadata.name` and AWS caps it at 32 characters;
both IaC modules truncate longer names deterministically. The
`target_group_name` output reports the name actually used.

`vpcId` is required for `instance`, `ip`, and `alb` target types and ignored
for `lambda`. That requiredness is enforced by the IaC modules rather than a
proto rule (message-level CEL cannot inspect reference-typed fields), so a
missing VPC fails fast at deploy time with a clear error rather than at the
AWS API.

## Draining, Stickiness, and Rollout Behavior

Several fields exist to make deployments boring:

- **`deregistrationDelaySeconds`** (default 300) keeps a draining target
  receiving in-flight responses; short-lived HTTP services usually lower it
  to 30-60 to speed up rollouts.
- **`slowStartSeconds`** (ALB) ramps traffic to newly registered targets so
  caches warm before full load; it is incompatible with stickiness and with
  the `least_outstanding_requests` algorithm -- AWS rejects the combinations.
- **`connectionTermination`** (NLB) closes established connections when the
  drain expires instead of waiting for clients; long-lived protocols
  (WebSocket, gRPC streams, databases) pin draining targets without it.
- **Stickiness** skews load distribution by design; prefer stateless targets
  and reserve `app_cookie`/`lb_cookie`/`source_ip` for workloads that truly
  keep per-client state.

## Group-Level Health Policy

Beyond per-target checks, the component exposes AWS's group-wide policy:

- **DNS failover** (`targetGroupHealth.dnsFailover`): when healthy targets
  drop below a count or percentage, the load balancer's DNS stops resolving
  to the affected zone, shifting clients elsewhere. Thresholds are strings
  because AWS accepts `"off"` to disable a criterion.
- **Unhealthy-state routing** (`targetGroupHealth.unhealthyStateRouting`):
  below the threshold, the load balancer routes to all targets, unhealthy
  included -- a partially working target beats a rejected request during a
  mass failure.
- **Unhealthy connection handling** (`targetHealthState`, NLB TCP/TLS):
  choose whether established connections to a target that turns unhealthy are
  kept (ride out a transient blip) or terminated (strict fail-fast).

## Dual-Engine Implementation

`AwsLbTargetGroup` ships both a Terraform/OpenTofu module and a Pulumi (Go)
module at behavioral parity. Both derive identity tags from metadata, apply
the same name truncation, enforce the same VPC requiredness, and export the
same outputs (`target_group_arn`, `target_group_name`, `arn_suffix`).
Whichever engine a team standardizes on, the group behaves identically.
