# EKS Chaining (on top of the AWS VPC CNI)

This preset runs Cilium ON TOP of the AWS VPC CNI on an EKS cluster —
CNI chaining, the no-rip-and-replace path. The AWS VPC CNI keeps doing
what EKS depends on (IPAM from the VPC, ENI wiring, VPC-routable pod IPs)
while Cilium attaches eBPF programs to every pod for NetworkPolicy
enforcement (including Cilium's L7-aware policies), eBPF load-balancing,
and Hubble flow observability.

## When to Use

- EKS clusters whose VPC networking must stay untouched, but which need
  real NetworkPolicy enforcement or L7-aware Cilium policies
- Adding flow-level observability (Hubble) to an existing EKS cluster
  without a networking migration
- Any cluster where replacing the incumbent CNI is off the table

Prefer **03-self-managed-primary-kpr** (or `cloud.aws_eni` with
`ipam.mode: eni`) when you CAN make Cilium the primary CNI — chaining
trades some datapath features for the zero-migration path.

## Key Configuration Choices

- **`cni.chainingMode: aws-cni`** — chain onto the AWS VPC CNI; IPAM and
  basic routing stay with the incumbent. aws-cni chaining implies the
  chaining target, so `chainingTarget` stays unset.
- **`cni.exclusive: false`** — required with chaining (the spec's CEL rule
  enforces it): exclusive mode renames non-Cilium CNI configurations in
  `/etc/cni/net.d`, which would break the very CNI the chain depends on.
- **`cloud` stays empty** — the AWS VPC CNI *is* the cloud integration
  here; `cloud.aws_eni` is the mutually-exclusive primary-CNI ENI datapath.
- **Chart defaults otherwise** — Hubble enabled in the agent, policy
  enforcement `default`, 2 operator replicas (fine on multi-node EKS).

## Placeholders to Replace

No placeholders — this preset is directly deployable.

## Related Presets

- **01-kind-dev-cluster** — the local/dev primary-CNI posture
- **03-self-managed-primary-kpr** — Cilium as primary CNI with kube-proxy
  replacement
- **04-production-observability** — layer Hubble relay/UI/metrics,
  Prometheus, and encryption on top of a production install
