---
title: "Presets"
description: "Ready-to-deploy configuration presets for Service Entry"
type: "preset-list"
componentSlug: "service-entry"
componentTitle: "Service Entry"
provider: "kubernetes"
icon: "package"
order: 200
presets:
  - slug: "01-external-https-api"
    rank: "01"
    title: "Reach an External HTTPS API"
    excerpt: "The canonical ServiceEntry: register an external service (a SaaS API, a partner endpoint) so mesh workloads can call it as a first-class destination, with TLS routed by SNI and the host resolved via..."
  - slug: "02-static-mesh-internal-endpoints"
    rank: "02"
    title: "Bring Static Endpoints Into the Mesh"
    excerpt: "Register a service that has a fixed set of backing IPs (a VM-hosted database, a legacy service, an appliance) as a MESH_INTERNAL destination with STATIC resolution. Mesh workloads then reach it by..."
  - slug: "03-dynamic-dns-wildcard-egress"
    rank: "03"
    title: "Dynamic-DNS Wildcard Egress"
    excerpt: "Registers a whole wildcard domain (`*.example-saas.com`) in the mesh's service registry with `DYNAMIC_DNS` resolution: the proxy resolves the ACTUAL hostname each request asks for, at request time —..."
---

# Service Entry Presets

Ready-to-deploy configuration presets for Service Entry. Each preset is a complete manifest you can copy, customize, and deploy.
