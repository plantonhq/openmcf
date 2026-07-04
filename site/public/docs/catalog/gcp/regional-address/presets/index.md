---
title: "Presets"
description: "Ready-to-deploy configuration presets for Regional Address"
type: "preset-list"
componentSlug: "regional-address"
componentTitle: "Regional Address"
provider: "gcp"
icon: "package"
order: 200
presets:
  - slug: "01-external-nat-ip"
    rank: "01"
    title: "External NAT IP"
    excerpt: "This preset reserves a regional external IPv4 address for use with Cloud NAT, regional load balancers, or VM instances. It is the simplest GcpAddress configuration — a project, a name, a region, and..."
  - slug: "02-internal-lb-vip"
    rank: "02"
    title: "Internal Load Balancer VIP"
    excerpt: "This preset reserves a regional internal IP address with the `SHARED_LOADBALANCER_VIP` purpose — the address type used by internal Application Load Balancers when multiple backends share a single VIP."
  - slug: "03-internal-gce-endpoint"
    rank: "03"
    title: "Internal GCE Endpoint"
    excerpt: "This preset reserves a regional internal IP address with the `GCE_ENDPOINT` purpose within a specific subnetwork — the address type used for VM instances, alias IP ranges, and similar compute..."
---

# Regional Address Presets

Ready-to-deploy configuration presets for Regional Address. Each preset is a complete manifest you can copy, customize, and deploy.
