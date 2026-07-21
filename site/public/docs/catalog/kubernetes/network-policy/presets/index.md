---
title: "Presets"
description: "Ready-to-deploy configuration presets for Network Policy"
type: "preset-list"
componentSlug: "network-policy"
componentTitle: "Network Policy"
provider: "kubernetes"
icon: "package"
order: 200
presets:
  - slug: "01-default-deny-all"
    rank: "01"
    title: "Default Deny All"
    excerpt: "This preset creates the namespace-wide security baseline: it denies ALL inbound and outbound traffic for every pod in the namespace. NetworkPolicies are additive allows with no deny rule, so \"deny..."
  - slug: "02-allow-same-namespace"
    rank: "02"
    title: "Allow Same Namespace"
    excerpt: "This preset isolates a namespace from the rest of the cluster while keeping it open internally: all pods in the namespace accept inbound traffic from all other pods in the SAME namespace, and nothing..."
  - slug: "03-allow-from-namespace"
    rank: "03"
    title: "Allow From Namespace"
    excerpt: "This preset allows inbound traffic to one workload's pods from ANY pod in a specific other namespace, on one port. Selecting the source namespace uses the automatic `kubernetes.io/metadata.name:..."
  - slug: "04-allow-dns-egress"
    rank: "04"
    title: "Allow DNS Egress"
    excerpt: "This preset allows all pods in the namespace to resolve DNS: egress on UDP and TCP port 53 to the cluster DNS pods in `kube-system`. It is the mandatory companion to any deny-all-egress posture — a..."
---

# Network Policy Presets

Ready-to-deploy configuration presets for Network Policy. Each preset is a complete manifest you can copy, customize, and deploy.
