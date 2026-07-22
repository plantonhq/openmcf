---
title: "Bare Metal / Edge (DaemonSet + Host Ports)"
description: "This preset installs the controller in the standard on-prem/edge shape: a DaemonSet running one controller pod per node, answering directly on node ports 80/443 through hostPorts, with a NodePort..."
type: "preset"
rank: "03"
presetSlug: "03-bare-metal-daemonset"
componentSlug: "kubernetesingressnginx"
componentTitle: "KubernetesIngressNginx"
provider: "kubernetes"
icon: "package"
order: 3
---

# Bare Metal / Edge (DaemonSet + Host Ports)

This preset installs the controller in the standard on-prem/edge shape: a
DaemonSet running one controller pod per node, answering directly on node
ports 80/443 through hostPorts, with a NodePort Service instead of a cloud
load balancer. Point your external load balancer, VIP (keepalived/MetalLB
in L2 mode), or DNS round-robin at the node addresses.

## When to Use

- Bare-metal, on-prem, and edge clusters with no cloud LB controller
- kind and other local clusters (also see the note below on why
  `load_balancer` must not be used here)
- Setups where an external L4 balancer you manage fronts the nodes

## Key Configuration Choices

- **`controllerKind: daemon_set`** — one controller per (selected) node:
  every node is an entry point, capacity scales with the node count, and no
  second hop is needed to reach a controller pod
- **`hostPorts: true`** — the controller's HTTP/HTTPS listeners bind node
  ports 80/443 via hostPort, staying inside the CNI. The alternative is
  `hostNetwork: true` (binds the node's interfaces directly; the module
  then sets `dnsPolicy: ClusterFirstWithHostNet`) — the two are mutually
  exclusive, choose one
- **`service.type: node_port`** — on clusters WITHOUT a cloud LB
  controller, a `load_balancer` service type FAILS LOUDLY at install: the
  module's Helm readiness wait includes the LB address, and an address that
  can never arrive times out. This is deliberate — the failure names the
  real problem instead of leaving a silently Pending entry point. Use
  `httpNodePort`/`httpsNodePort` to pin the node ports when your external
  balancer needs fixed targets

Scope the DaemonSet to dedicated edge nodes with `nodeSelector` (plus
`tolerations` if those nodes are tainted).

## Placeholders to Replace

None — this preset is deployable as-is on any cluster.

## Related Presets

- **01-aws-nlb-public** — the managed-cloud equivalent of this posture
- **04-tls-default-cert** — add a cluster-wide default TLS certificate
