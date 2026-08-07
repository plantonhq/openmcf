---
title: "Kafka on Confluent Cloud"
description: "Kafka on Confluent Cloud deployment documentation"
icon: "package"
order: 100
componentName: "confluentkafka"
---

# Kafka on Confluent Cloud

Deploys a Confluent Cloud Kafka cluster with configurable cluster type, multi-zone availability, and optional private networking across AWS, Azure, and GCP regions. Supports Basic, Standard, Enterprise, and Dedicated cluster tiers with automatic type-specific configuration. Integrates with Planton's Confluent Provider Connection for API key management.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **Confluent Cloud Kafka Cluster** -- a managed Kafka cluster of the specified type (Basic, Standard, Enterprise, or Dedicated), placed in the configured cloud provider region and associated with a Confluent Cloud environment with the specified availability zone setting
- **Network Association** -- created only when `networkConfig` is provided, associates the cluster with a pre-existing Confluent Cloud network resource for private connectivity (PrivateLink on AWS, Private Link on Azure, Private Service Connect on GCP)

## Before You Deploy

### Planton Setup

- **Confluent Provider Connection** -- an active connection in the Connect module with a Confluent Cloud API key and secret. Map it as the default for your environment, or specify it explicitly when creating the Cloud Resource.
- **Planton Runner** -- required when using Runner-based credential delivery. Not needed for inline API key authentication.

### Confluent Cloud Account

- **A Confluent Cloud environment** -- the `environmentId` of an existing Confluent Cloud environment where the cluster will be created. Environments act as the parent container for clusters and other Confluent resources.
- **A Confluent Cloud network** (optional) -- required for private networking with Enterprise or Dedicated clusters. Must be pre-created in the same environment. Provide the network ID via `networkConfig.networkId`.
- **Sufficient CKU quota** (optional) -- required only for Dedicated clusters. Confluent Kafka Units are provisioned capacity units that determine throughput. Minimum: 1 CKU.
- **Cloud provider and region** -- verify that your chosen `cloud` and `region` combination supports the desired `clusterType`. Not all cluster types are available in every region.
- **Cluster naming** (optional) -- set `displayName` to customize the cluster name shown in the Confluent Cloud console. Defaults to `metadata.name` when not specified.

## Deploy

### Console

Open the deployment store, find **Kafka on Confluent Cloud**, and click **Deploy**. The creation wizard walks you through environment and connection configuration, and spec fields including cluster type, cloud provider, region, and availability.

### CLI

Create a manifest and apply it:

```yaml
apiVersion: confluent.planton.dev/v1
kind: ConfluentKafka
metadata:
  name: event-bus
  org: acme-corp
  env: prod
spec:
  cloud: AWS
  region: us-east-2
  availability: MULTI_ZONE
  environmentId: env-abc123
```

```shell
planton apply -f confluent-kafka.yaml
```

This creates a Standard Kafka cluster with multi-zone availability on AWS us-east-2, associated with the specified Confluent Cloud environment. The default cluster type is Standard when `clusterType` is not specified. No private networking or dedicated capacity is configured. A Stack Job tracks the provisioning in real time.

For a Dedicated cluster with private networking, add `clusterType: DEDICATED`, `dedicatedConfig`, and `networkConfig` to the spec.

## Key Configuration

These are the most important decisions when configuring a Confluent Cloud Kafka cluster. Explore the full field reference in the [API Explorer](#api-explorer) tab.

**Cluster type** -- The `clusterType` field determines the deployment model and capabilities. BASIC is multi-tenant and single-zone only, suitable for development. STANDARD is multi-tenant with elastic scaling for production. ENTERPRISE adds private networking support on shared infrastructure. DEDICATED provides single-tenant isolation with provisioned CKU capacity and private networking.

**Availability** -- Set `availability` to `MULTI_ZONE` for production workloads requiring a 99.99% SLA. Use `SINGLE_ZONE` for development and testing where cost is prioritized over resilience. Multi-zone is required for Standard and Dedicated clusters in production.

**Cloud and region** -- The `cloud` field accepts `AWS`, `AZURE`, or `GCP`. The `region` field accepts the cloud-specific region code (e.g., `us-east-2` for AWS, `us-central1` for GCP, `eastus` for Azure). Choose the region closest to your producers and consumers.

**Private networking** -- Set `networkConfig.networkId` to associate the cluster with a Confluent Cloud network for private connectivity. Available only for Enterprise and Dedicated cluster types. Without this configuration, the cluster uses public internet access.

**Dedicated capacity** -- When `clusterType` is `DEDICATED`, set `dedicatedConfig.cku` to the number of Confluent Kafka Units. CKU determines throughput capacity and can be scaled up or down after creation but not to zero.

## Outputs and Dependencies

### What This Component Consumes

This component has no foreign key dependencies.

### What This Component Provides

After provisioning, `status.outputs` contains values that downstream Cloud Resources can consume via ValueFromRef:

| Output | Description | Common Downstream Use |
|--------|-------------|----------------------|
| `id` | Provider-assigned unique ID for the Kafka cluster | API operations, resource tracking |
| `bootstrap_endpoint` | Bootstrap endpoint for Kafka client connections (SASL_SSL protocol) | Producer and consumer client configuration |
| `crn` | Confluent Resource Name for RBAC and API references | Access control policies, Confluent CLI operations |
| `rest_endpoint` | REST endpoint for the Kafka cluster | REST Proxy access, admin API operations |

## Common Patterns

No presets are available yet. Configure directly using the fields documented in the [API Explorer](#api-explorer) tab.

## Works With

This component operates independently and does not reference other deployment components.