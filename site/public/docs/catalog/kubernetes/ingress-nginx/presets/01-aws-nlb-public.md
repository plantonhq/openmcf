---
title: "AWS NLB Public Entry (EKS)"
description: "This preset installs the cluster's public entry point on EKS: an ingress-nginx controller behind an AWS Network Load Balancer, running two replicas, owning the standard `nginx` class as the cluster..."
type: "preset"
rank: "01"
presetSlug: "01-aws-nlb-public"
componentSlug: "ingress-nginx"
componentTitle: "Ingress NGINX"
provider: "kubernetes"
icon: "package"
order: 1
---

# AWS NLB Public Entry (EKS)

This preset installs the cluster's public entry point on EKS: an
ingress-nginx controller behind an AWS Network Load Balancer, running two
replicas, owning the standard `nginx` class as the cluster default. This is
the standard production posture on AWS.

## When to Use

- EKS clusters serving public HTTP(S) traffic through Ingress resources
- The cluster's primary (or only) ingress controller — the class every
  Ingress without an explicit `ingressClassName` should land on
- Production deployments where client source IPs must survive to NGINX

## Key Configuration Choices

- **NLB annotations** (`aws-load-balancer-type: external` +
  `nlb-target-type: ip`) — provisions a Network Load Balancer targeting pod
  IPs directly through the AWS Load Balancer Controller, instead of the
  legacy in-tree classic ELB
- **`externalTrafficPolicy: local`** — traffic is delivered only to nodes
  running a controller pod, preserving the client source IP and skipping an
  extra hop; LB health checks handle the node selection
- **`replicas: 2`** — survives a node failure and rolls without dropping
  traffic; the chart automatically adds a PodDisruptionBudget
  (minAvailable 1) when replicas exceed one
- **`isDefaultClass: true`** — Ingress resources created without an
  `ingressClassName` are assigned to this controller; at most one class per
  cluster may claim default
- **Class `nginx`** — the conventional name; its controller identifier
  derives to the chart default `k8s.io/ingress-nginx`

After deploy, the NLB's DNS name lands in the `load_balancer_hostname`
output (AWS populates a hostname, not an IP) — point DNS records at it, or
let KubernetesExternalDns publish them.

## Placeholders to Replace

None — this preset is deployable as-is on an EKS cluster with the AWS Load
Balancer Controller installed.

## Related Presets

- **02-internal-only** — a second controller instance for private traffic
- **04-tls-default-cert** — add a cluster-wide default TLS certificate from
  cert-manager
