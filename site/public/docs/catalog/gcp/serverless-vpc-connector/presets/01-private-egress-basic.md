---
title: "Private egress — basic"
description: "The standard connector for giving Cloud Functions, Cloud Run, and App Engine access to private VPC resources (Cloud SQL private IP, Memorystore, internal load balancers). Network placement carves a..."
type: "preset"
rank: "01"
presetSlug: "01-private-egress-basic"
componentSlug: "serverless-vpc-connector"
componentTitle: "Serverless VPC Connector"
provider: "gcp"
icon: "package"
order: 1
---

# Private egress — basic

The standard connector for giving Cloud Functions, Cloud Run, and App Engine access to private VPC resources (Cloud SQL private IP, Memorystore, internal load balancers). Network placement carves a dedicated `/28` out of the VPC — pick any range that overlaps no existing subnet or route.

One connector serves every serverless workload in its region: attach it by reference from [GcpCloudFunction](/docs/catalog/gcp/gcpcloudfunction), [GcpCloudRun](/docs/catalog/gcp/gcpcloudrun), or [GcpCloudRunJob](/docs/catalog/gcp/gcpcloudrunjob) rather than creating one per service. `maxInstances: 4` caps idle spend while leaving room to scale; note the fleet never scales in on its own after a burst.
