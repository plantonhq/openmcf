# KubernetesIngressNginx Guide

The judgment this guide carries: the ingress controller is shared-cluster
infrastructure — one installation serves every application environment on
the cluster. Proposing a per-app controller, or an Ingress with no
controller layer at all, are the two ways agent-composed exposure goes
wrong before a single routing rule matters.

## Where it sits in the public-endpoint chain

An endpoint that must serve public HTTPS at a real domain needs four
layers, and only one of them is the app's own manifest:

1. **This controller** — once per cluster, in the shared-cluster chart.
2. **KubernetesIngress** objects — per application, selecting this
   controller through its ingress class (the controller's
   `status.outputs.ingress_class_name`).
3. **The certificate chain** — cert-manager plus a cluster issuer; the
   checklist lives in the
   [KubernetesClusterIssuer guide](../../kubernetesclusterissuer/v1alpha1/GUIDE.md).
4. **The DNS record** — published by
   [KubernetesExternalDns](../../kubernetesexternaldns/v1alpha1/GUIDE.md) from
   the Ingress's hostname.

A second controller instance is for a real traffic tier (public vs
internal), never for a second application — the multi-instance rules are
on [reference.md](reference.md).

## Namespace ownership — the infra exception

Unlike shared application namespaces, this controller conventionally owns
a dedicated single-tenant namespace ("ingress-nginx"), so
`createNamespace: true` is the normal shape here — the
[namespace-ownership pattern](../../../patterns/namespace-ownership.md)'s
sole-tenant case, not a violation of it.

## On the diagram

The controller renders as the cluster's entry node in the shared-cluster layer;
each application's Ingress draws the edge from that entry into its
workload. An architecture showing workloads with no entry node is
promising reachability it cannot deliver.

## Pairs well with

- KubernetesIngress — the per-application routing rules that select this
  controller.
- KubernetesCertManager + KubernetesClusterIssuer — the HTTPS half of the
  chain; the controller's `defaultTlsCertificate` can reference a
  KubernetesCertificate's secret output for a cluster-wide default.
- KubernetesExternalDns — turns Ingress hostnames into real DNS records.
- KubernetesMetricsServer — required before enabling the controller's
  autoscaling block.
