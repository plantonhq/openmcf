---
title: "Self-Managed Primary CNI with Kube-Proxy Replacement"
description: "This preset makes Cilium the cluster's primary CNI on a self-managed (kubeadm-style) cluster AND replaces kube-proxy entirely with Cilium's eBPF service load-balancing. The operator hands out pod IPs..."
type: "preset"
rank: "03"
presetSlug: "03-self-managed-primary-kpr"
componentSlug: "cilium"
componentTitle: "Cilium"
provider: "kubernetes"
icon: "package"
order: 3
---

# Self-Managed Primary CNI with Kube-Proxy Replacement

This preset makes Cilium the cluster's primary CNI on a self-managed
(kubeadm-style) cluster AND replaces kube-proxy entirely with Cilium's
eBPF service load-balancing. The operator hands out pod IPs from an
explicit cluster pool, and inter-node traffic is VXLAN-encapsulated so the
underlying network needs no awareness of pod CIDRs.

## When to Use

- Self-managed clusters (kubeadm, Cluster API, bare metal) created without
  a CNI addon and with kube-proxy skipped
  (`kubeadm init --skip-phases=addon/kube-proxy`)
- You want the kube-proxy-free posture: eBPF service resolution instead of
  iptables chains, and the datapath features that require it (Cilium's
  Gateway API support, for one, is gated on kube-proxy replacement)
- Networks where you control the pod CIDR plan and want it explicit rather
  than the chart's blanket `10.0.0.0/8` default

## Key Configuration Choices

- **`kubeProxyReplacement: true`** — Cilium serves ClusterIP/NodePort/
  LoadBalancer traffic in eBPF; the cluster must run WITHOUT kube-proxy
- **`k8sServiceHost` + `k8sServicePort: 6443`** — required with kube-proxy
  replacement (the spec's CEL rule enforces it): before Cilium's own
  load-balancing is up there is no kube-proxy to resolve the in-cluster
  `kubernetes.default` Service, so the agent needs the API server's real
  address
- **`ipam.mode: cluster-pool` with `10.42.0.0/16`, mask size 24** — the
  operator carves a /24 (254 pod IPs) per node from a deliberately chosen
  pool instead of the chart's `10.0.0.0/8`; pick CIDRs that overlap neither
  the node network nor the Service CIDR
- **`routing: tunnel` / `vxlan`** — encapsulation works on any fabric;
  switch to `native` mode (with `ipv4NativeRoutingCidr`) only when the
  underlying network can route pod CIDRs

## Placeholders to Replace

| Placeholder | Description | Where to Find |
|---|---|---|
| `<api-server-endpoint>` | API server address the agent dials before service load-balancing exists — a hostname or IP reachable from every node (the control-plane endpoint / load-balancer in front of the API servers) | Your kubeadm `controlPlaneEndpoint` / cluster provisioning config |

Also review `clusterPoolIpv4PodCidrs` (`10.42.0.0/16`) against your
network plan before deploying.

## Related Presets

- **01-kind-dev-cluster** — the same primary-CNI idea at laptop scale
- **04-production-observability** — combine with this preset's fields for
  a fully observed, encrypted production dataplane
