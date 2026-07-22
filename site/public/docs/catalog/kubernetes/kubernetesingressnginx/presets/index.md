---
title: "Presets"
description: "Ready-to-deploy configuration presets for KubernetesIngressNginx"
type: "preset-list"
componentSlug: "kubernetesingressnginx"
componentTitle: "KubernetesIngressNginx"
provider: "kubernetes"
icon: "package"
order: 200
presets:
  - slug: "01-aws-nlb-public"
    rank: "01"
    title: "AWS NLB Public Entry (EKS)"
    excerpt: "This preset installs the cluster's public entry point on EKS: an ingress-nginx controller behind an AWS Network Load Balancer, running two replicas, owning the standard `nginx` class as the cluster..."
  - slug: "02-internal-only"
    rank: "02"
    title: "Internal-Only Controller (Second Instance)"
    excerpt: "This preset installs a SECOND ingress-nginx controller dedicated to private traffic — the other half of the public + internal split. It runs in its own namespace with its own IngressClass..."
  - slug: "03-bare-metal-daemonset"
    rank: "03"
    title: "Bare Metal / Edge (DaemonSet + Host Ports)"
    excerpt: "This preset installs the controller in the standard on-prem/edge shape: a DaemonSet running one controller pod per node, answering directly on node ports 80/443 through hostPorts, with a NodePort..."
  - slug: "04-tls-default-cert"
    rank: "04"
    title: "Cluster-Wide Default TLS Certificate (cert-manager Composition)"
    excerpt: "This preset installs the public entry controller with a cluster-wide default TLS certificate wired from cert-manager: a KubernetesCertificate's secret output flows into `defaultTlsCertificate` as a..."
---

# KubernetesIngressNginx Presets

Ready-to-deploy configuration presets for KubernetesIngressNginx. Each preset is a complete manifest you can copy, customize, and deploy.
