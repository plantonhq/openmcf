---
title: "Custom Mode VPC with Regional Routing"
description: "This preset creates a VPC in custom subnet mode with regional routing. Custom mode gives full control over subnet CIDR ranges and regions."
type: "preset"
rank: "01"
presetSlug: "01-custom-mode-regional"
componentSlug: "vpc"
componentTitle: "VPC"
provider: "gcp"
icon: "package"
order: 1
---

# Custom Mode VPC with Regional Routing

This preset creates a VPC in custom subnet mode with regional routing. Custom mode gives full control over subnet CIDR ranges and regions.

Private connectivity to Google managed services (Cloud SQL, AlloyDB, Memorystore) is composed separately: reserve a `GcpGlobalAddress` with `purpose: VPC_PEERING` and connect it through a `GcpServiceNetworkingConnection` on this network.

## When to Use

- Production workloads within a single GCP region
- Any project following the GCP best practice of custom-mode VPCs over auto-mode

## What It Configures

- **Custom mode** (`autoCreateSubnetworks: false`) -- no auto-created subnets; you define subnets explicitly with `GcpSubnetwork`
- **Regional routing** (`routingMode: REGIONAL`) -- Cloud Router advertises routes within one region only (simpler, lower cost)
- **Description** -- documents intent for operators browsing the console

## After Deployment

Create subnets with `GcpSubnetwork` referencing this VPC's `network_self_link` output. For managed-service private IP, add the `GcpGlobalAddress` (VPC_PEERING range) + `GcpServiceNetworkingConnection` pair.
