# KubernetesIngress Guide

The judgment this guide carries: Ingress or Gateway API is a decision the
CLUSTER already made for you more often than not — read the cluster's
shape before picking — and an Ingress is only one link of a four-part
chain that a complete proposal names in full.

## Ingress vs Gateway API (the comparison home)

Both express north-south HTTP(S) exposure; choose by what serves the
cluster, not by novelty:

- **Ingress + KubernetesIngressNginx** — the classic lane: one shared
  controller, host/path routing, controller-specific behavior through
  annotations. The right default on clusters whose entry point is
  ingress-nginx.
- **Gateway API** (KubernetesGateway + route kinds) — the right lane when
  the cluster's networking layer ALREADY implements it: Istio implements
  Gateway API natively (an Istio cluster's north-south path IS the
  Gateway family — its own spec says it deploys no ingress), and Cilium
  provides it behind its `gatewayApi` flag. It is also where
  protocol-typed routing (gRPC, TCP, TLS passthrough) lives. The family
  anchor: [KubernetesGateway guide](../../kubernetesgateway/v1/GUIDE.md).

What breaks when chosen wrong: an Ingress on a mesh-native cluster
bypasses the mesh's own entry (two entry stacks to operate); Gateway
kinds on a cluster with only ingress-nginx sit inert — no implementing
controller ever claims them.

## The chain this kind completes (and never completes alone)

1. A controller claiming the `ingressClassName`
   ([KubernetesIngressNginx guide](../../kubernetesingressnginx/v1/GUIDE.md)) —
   an Ingress without one is valid and inert (Layer 1 states it).
2. This kind's host rules → the workload's exported Service, wired by
   `valueFrom` (the backend `serviceName` is a foreign key).
3. TLS: the `tls` block names a Secret that
   [cert-manager issues](../../kubernetesclusterissuer/v1/GUIDE.md) via
   the issuer annotation — leave the Secret non-existent, cert-manager
   creates it under exactly that name.
4. DNS: the hostname becomes real through
   [KubernetesExternalDns](../../kubernetesexternaldns/v1/GUIDE.md).

## On the diagram

Ingress renders between the controller's entry node and the workload's
Service, with the `valueFrom` backend edge drawn. TLS-by-annotation adds
no node — which is why the issuer chain must be checked deliberately in
the shared-cluster layer.

## Pairs well with

- KubernetesIngressNginx — the controller lane this kind rides.
- KubernetesCertificate / the cluster issuer chain — the TLS Secret's
  origin.
- KubernetesExternalDns — the hostname's publisher.
- KubernetesService — the backend handle, wired by reference.
