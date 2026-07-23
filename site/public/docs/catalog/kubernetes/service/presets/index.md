---
title: "Presets"
description: "Ready-to-deploy configuration presets for Service"
type: "preset-list"
componentSlug: "service"
componentTitle: "Service"
provider: "kubernetes"
icon: "package"
order: 200
presets:
  - slug: "01-cluster-ip-app"
    rank: "01"
    title: "ClusterIP for an Application"
    excerpt: "This preset creates the default kind of Service: a cluster-internal virtual IP and DNS name in front of an application's pods. Anything inside the cluster reaches the app at..."
  - slug: "02-public-load-balancer"
    rank: "02"
    title: "Public Load Balancer"
    excerpt: "This preset asks the cloud provider to provision an external load balancer in front of the selected pods. Once deployed, the provider's address lands in the stack outputs — `load_balancer_ip` on..."
  - slug: "03-headless-statefulset"
    rank: "03"
    title: "Headless Service for StatefulSet Peers"
    excerpt: "This preset creates a headless service (`clusterIP: None`): no virtual IP is allocated and DNS returns the pod IPs directly, with each pod also getting its own stable name..."
  - slug: "04-external-name"
    rank: "04"
    title: "ExternalName Alias"
    excerpt: "This preset creates a pure DNS alias: cluster DNS answers lookups for `prod-db.<namespace>.svc.cluster.local` with a CNAME to `db.prod.example.com`. No proxying, no selectors, no ports — traffic..."
---

# Service Presets

Ready-to-deploy configuration presets for Service. Each preset is a complete manifest you can copy, customize, and deploy.
