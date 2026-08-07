# KubernetesIstio Guide

The judgment this guide carries: this installs the mesh's BRAIN
(istiod) and nothing user-facing — north-south exposure and mesh policy
are both composed from other kinds. The one decision that reshapes every
workload is sidecar vs ambient, and it is chosen here.

## What this is, and what it is NOT

istiod validates mesh config, issues workload certificates, and programs
the data plane. It deploys NO gateways and NO ingress. Two families
compose around it:

- **North-south exposure** is Gateway API, natively: create a
  KubernetesGateway with `gatewayClassName: istio` and istiod provisions
  the gateway; routes attach to it. The
  [Gateway anchor](../../kubernetesgateway/v1alpha1/GUIDE.md) is that story
  (Istio is one of its two implementing controllers).
- **Mesh traffic policy** composes from the typed Istio kinds
  (KubernetesDestinationRule, KubernetesPeerAuthentication,
  KubernetesAuthorizationPolicy, KubernetesTelemetry, ...). Those are
  faithful upstream-API mirrors — their reference pages are the
  authority, and they require only the Istio CRDs (installed here, or
  standalone via KubernetesIstioBaseCrds on clusters that use the policy
  APIs without a full mesh).

## Sidecar vs ambient — the data-plane decision

The reshaping choice (the field doc on [reference.md](reference.md)):

- **`sidecar`** (default) injects a proxy container into every enrolled
  workload pod — mature, full L7 per pod, at the cost of a container in
  every pod and a restart to enroll.
- **`ambient`** runs no sidecars — a per-node ztunnel DaemonSet carries
  mTLS/L4 and waypoint proxies add L7 only where needed; enrollment is a
  namespace/pod label, no pod restart. Lower per-pod overhead, newer.

State the choice in any mesh proposal — it changes what every workload
in the mesh looks like and how enrollment works.

## Once per cluster

One control plane per cluster; its namespace (istiod, and in ambient
mode the CNI agent and ztunnel) is the
[namespace-ownership pattern](../../../_patterns/namespace-ownership.md)'s
sole-tenant case. The policy kinds depend on its CRDs exactly as the
[operator-prerequisite pattern](../../../_patterns/operator-prerequisite.md)
describes.

## On the diagram

istiod renders as a shared-cluster node. Gateways (Istio-classed) and the
mesh policy kinds render as their own nodes; the mesh membership of
workloads is a label, not a drawn edge — a reviewer confirms the control
plane exists whenever Istio-classed gateways or Istio policy kinds appear.

## Pairs well with

- KubernetesGateway (class `istio`) — north-south exposure istiod
  provisions.
- KubernetesDestinationRule / PeerAuthentication / AuthorizationPolicy /
  Telemetry — mesh policy (upstream-mirror kinds; reference pages are the
  authority).
- KubernetesIstioBaseCrds — the CRDs alone, for policy-API use without a
  full mesh.
