# KubernetesDeployment Guide

The judgment this guide carries: a Deployment is almost never the whole
architecture. It deliberately exports handles instead of embedding
exposure, autoscaling beyond the basics, or namespace ownership — the
completeness of a proposal depends on composing the right neighbors, not
on filling more of this spec.

## Exposure is composed — finish the chain

This kind creates the Service but nothing that makes it reachable from
outside the cluster. When the user's app needs a public endpoint, the
proposal must also include an exposure kind (KubernetesIngress, or
KubernetesHttpRoute behind a Gateway) pointed at the Service this kind
exports — `status.outputs.kube_endpoint` is the in-cluster handle, and
`status.outputs.selector_labels` feeds NetworkPolicy pod selectors. Public
HTTPS additionally needs the certificate chain; the composition checklist
lives in the [KubernetesClusterIssuer guide](../../kubernetesclusterissuer/v1/GUIDE.md).

## Two ways to autoscale — pick by metric surface

- `spec.availability.horizontalPodAutoscaling` (inline) covers CPU and
  memory targets — right for the common case, zero extra nodes.
- The standalone KubernetesHorizontalPodAutoscaler kind carries the full
  autoscaling/v2 surface (custom metrics, external metrics, scaling
  behavior policies) and references this kind's
  `status.outputs.deployment_name`. Reach for it when scaling on anything
  beyond CPU/memory — and never declare both for the same workload: two
  controllers fighting over one replica count.

## Namespace ownership

`spec.namespace` is a required foreign key targeting KubernetesNamespace.
`createNamespace: true` makes this workload the namespace's owner in IaC
state — safe only when it is the sole tenant. Any shared namespace wants a
dedicated KubernetesNamespace component instead; the failure story and the
`valueFrom` wiring: [namespace-ownership pattern](../../../patterns/namespace-ownership.md).

## Pipeline-attached services

When this resource is a Service Hub deploy target, the pipeline owns
`spec.version` and `spec.container.app.image` (the deploy-target contract
on [reference.md](reference.md)). Compose the initial manifest with a real
image, but do not build automation that edits those two paths out-of-band
— the pipeline writes them on every build.

## On the diagram

Each composed piece — namespace, ingress/route, certificate chain,
standalone autoscaler — renders as its own node with reference edges into
this one. An architecture that buries exposure in annotations or namespace
creation in a flag shows a lone deployment node and hides exactly the
infrastructure a reviewer needs to see.

## Pairs well with

- KubernetesNamespace — the namespace owner (pattern above).
- KubernetesIngress / KubernetesHttpRoute — external exposure, wired to
  the exported Service.
- KubernetesHorizontalPodAutoscaler — advanced autoscaling (judgment
  above).
- KubernetesConfigMap / KubernetesSecret / KubernetesExternalSecret —
  configuration and credentials mounted by the containers.
