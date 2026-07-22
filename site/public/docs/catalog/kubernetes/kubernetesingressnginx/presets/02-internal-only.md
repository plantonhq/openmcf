---
title: "Internal-Only Controller (Second Instance)"
description: "This preset installs a SECOND ingress-nginx controller dedicated to private traffic — the other half of the public + internal split. It runs in its own namespace with its own IngressClass..."
type: "preset"
rank: "02"
presetSlug: "02-internal-only"
componentSlug: "kubernetesingressnginx"
componentTitle: "KubernetesIngressNginx"
provider: "kubernetes"
icon: "package"
order: 2
---

# Internal-Only Controller (Second Instance)

This preset installs a SECOND ingress-nginx controller dedicated to private
traffic — the other half of the public + internal split. It runs in its own
namespace with its own IngressClass (`nginx-internal`) behind an
internal-scheme load balancer that is only reachable inside the VPC/VNet.
Internal services publish Ingresses with `ingressClassName:
nginx-internal`; nothing about the public instance changes.

## When to Use

- Clusters serving both public and private traffic that need real
  separation — different LBs, different classes, independent lifecycles
- Admin panels, internal APIs, and dashboards that must never be reachable
  from the internet
- As the second instance next to **01-aws-nlb-public** (or any public
  controller)

For a lighter-weight variant — the SAME controller pods answering both a
public and a private address — use `service.internal.enabled` on a single
instance instead. Two instances is the recommended posture when the traffic
sets should not share fate.

## Key Configuration Choices

- **Distinct class** (`nginx-internal`) — every controller instance on a
  cluster must own its own IngressClass; the controller identifier derives
  automatically to `k8s.io/nginx-internal`, so the two instances never
  fight over Ingresses
- **Own namespace** (`ingress-nginx-internal`) — instances are fully
  independent: separate release, separate resources, separate
  leader-election identity (all derived from `metadata.name`)
- **The internal-LB annotation** — what makes the load balancer private.
  Per cloud:
  - AWS: `service.beta.kubernetes.io/aws-load-balancer-scheme: "internal"`
    (kept alongside the NLB annotations, as in the YAML)
  - GCP: `networking.gke.io/load-balancer-type: "Internal"`
  - Azure: `service.beta.kubernetes.io/azure-load-balancer-internal: "true"`

The YAML shows the AWS variant; swap the annotations map for your cloud's
recipe.

## Placeholders to Replace

None — this preset is deployable as-is on EKS; substitute the annotation
recipe on GKE/AKS.

## Related Presets

- **01-aws-nlb-public** — the public half of the split
- **04-tls-default-cert** — a default certificate for the internal domain
