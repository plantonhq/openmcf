---
title: "Parallel VPC cleanup"
description: "Ten cleanup tasks with five running concurrently — the standard map-reduce shape on Cloud Run jobs. Tasks reach private APIs through direct VPC egress (`networkInterfaces` + `PRIVATE_RANGES_ONLY`)..."
type: "preset"
rank: "02"
presetSlug: "02-parallel-vpc-cleanup"
componentSlug: "cloud-run-job"
componentTitle: "Cloud Run Job"
provider: "gcp"
icon: "package"
order: 2
---

# Parallel VPC cleanup

Ten cleanup tasks with five running concurrently — the standard map-reduce shape on Cloud Run jobs. Tasks reach private APIs through direct VPC egress (`networkInterfaces` + `PRIVATE_RANGES_ONLY`) and pull credentials from Secret Manager at start.

The runtime identity needs `roles/secretmanager.secretAccessor` on the secret and network reachability to the cleanup targets via [GcpSubnetwork](/docs/catalog/gcp/gcpsubnetwork) firewall rules.
