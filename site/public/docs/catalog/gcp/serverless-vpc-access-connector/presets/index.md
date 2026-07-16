---
title: "Presets"
description: "Ready-to-deploy configuration presets for Serverless VPC Access Connector"
type: "preset-list"
componentSlug: "serverless-vpc-access-connector"
componentTitle: "Serverless VPC Access Connector"
provider: "gcp"
icon: "package"
order: 200
presets:
  - slug: "01-private-egress-basic"
    rank: "01"
    title: "Private egress — basic"
    excerpt: "The standard connector for giving Cloud Functions, Cloud Run, and App Engine access to private VPC resources (Cloud SQL private IP, Memorystore, internal load balancers). Network placement carves a..."
  - slug: "02-high-throughput"
    rank: "02"
    title: "High throughput — production"
    excerpt: "A production-scale connector for serverless fleets that push serious traffic into the VPC. `e2-standard-4` instances carry roughly a 1 Gbps class of throughput each (versus ~200 Mbps for `e2-micro`),..."
  - slug: "03-shared-vpc-subnet"
    rank: "03"
    title: "Shared VPC — subnet placement"
    excerpt: "The required shape on Shared VPC: the connector cannot carve its own range out of a host-project network, so a network admin creates a dedicated `/28` subnetwork in the host project..."
---

# Serverless VPC Access Connector Presets

Ready-to-deploy configuration presets for Serverless VPC Access Connector. Each preset is a complete manifest you can copy, customize, and deploy.
