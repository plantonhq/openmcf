---
title: "SAE Application"
description: "SAE Application deployment documentation"
icon: "package"
order: 100
componentName: "alicloudsaeapplication"
---

# AliCloud SAE Application

Deploys an Alibaba Cloud Serverless App Engine (SAE) application with configurable package type, discrete CPU/memory tiers, optional VPC networking, health checks, update strategies, and SLS log collection. SAE is a container-based serverless compute platform that handles provisioning, scaling, load balancing, and log collection. The component integrates with Planton's Provider Connections for AliCloud credential management and supports ValueFromRef wiring to VPCs, VSwitches, and security groups.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **SAE Application** -- an `alicloud_sae_application` with the configured package type, instance count, compute sizing, and optional VPC networking, health checks, update strategy, and SLS log collection
- **AliCloud Tags** -- resource metadata tags (organization, environment, resource kind, resource ID) merged with user-provided `tags`

## Before You Deploy

### Planton Setup

- **AliCloud Provider Connection** -- an active connection in the Connect module with credentials for the target Alibaba Cloud account. Map it as the default for your environment, or specify it explicitly when creating the Cloud Resource.
- **Planton Runner** -- required when using Runner-based credential delivery. Not needed for inline credentials.

### Alibaba Cloud Account

- A container image in ACR (for `packageType: Image`) or a JAR/WAR/ZIP package accessible via URL (for `FatJar`, `War`, `PythonZip`, `PhpZip`).
- A VPC, VSwitch, and security group if deploying in VPC mode. SAE provides managed networking by default -- VPC configuration is optional. Provide IDs directly or reference AliCloudVpc, AliCloudVswitch, and AliCloudSecurityGroup Cloud Resources via ValueFromRef. VPC ID is immutable after creation.
- An ACR Enterprise Edition instance ID if pulling images from a private ACR EE registry (not needed for ACR Personal Edition).

## Deploy

### Console

Open the deployment store, find **AliCloud SAE Application**, and click **Deploy**. The creation wizard walks you through preset selection, environment and connection configuration, and spec fields. Start from the **Container Image Production** preset in the [Presets](#presets) tab to pre-populate a VPC-connected container application with health checks and a rolling update strategy.

### CLI

```yaml
apiVersion: alicloud.planton.dev/v1
kind: AliCloudSaeApplication
metadata:
  name: api-service
  org: acme-corp
  env: prod
spec:
  region: cn-hangzhou
  appName: api-service
  packageType: Image
  replicas: 2
  cpu: 2000
  memory: 4096
  imageUrl: registry.cn-hangzhou.aliyuncs.com/my-ns/api-service:v1.0.0
```

```shell
planton apply -f sae-application.yaml
```

This creates an SAE application running 2 instances with 2 vCPUs and 4 GB memory each, using SAE's managed networking. No VPC attachment, health checks, or update strategy are configured.

### InfraChart

When deploying as part of a multi-resource environment, use ValueFromRef to wire the application to a VPC, VSwitch, and security group deployed in the same InfraPipeline:

```yaml
spec:
  vpcId:
    valueFrom:
      kind: AliCloudVpc
      name: app-vpc
      fieldPath: status.outputs.vpc_id
  vswitchId:
    valueFrom:
      kind: AliCloudVswitch
      name: app-vswitch
      fieldPath: status.outputs.vswitch_id
  securityGroupId:
    valueFrom:
      kind: AliCloudSecurityGroup
      name: app-sg
      fieldPath: status.outputs.security_group_id
```

The InfraPipeline resolves the dependency graph, deploys the VPC, VSwitch, and security group first, then provisions the SAE application with the resolved values.

## Key Configuration

These are the most important decisions when configuring an SAE application. Explore the full field reference in the [API Explorer](#api-explorer) tab.

**Package type** -- Set `packageType` to `Image` (container), `FatJar` (Spring Boot), `War` (servlet), `PythonZip`, or `PhpZip`. This is immutable after creation. Image deployments use `imageUrl`; package deployments use `packageUrl` and `packageVersion`. Java applications additionally configure `jdk` and `jarStartOptions`.

**Compute sizing** -- CPU and memory are specified as discrete tiers, not arbitrary values. CPU options: 500, 1000, 2000, 4000, 8000, 16000, 32000 millicores. Memory options: 1024, 2048, 4096, 8192, 12288, 16384, 24576, 32768, 65536, 131072 MB. Set `replicas` to control horizontal scale and `minReadyInstances` to maintain availability during deployments.

**Health checks** -- Configure `liveness` and `readiness` probes using HTTP GET, TCP socket, or exec checks. SAE restarts instances that fail the liveness check and withholds traffic from instances that fail the readiness check. Java applications with heavy startup (model loading, JIT warmup) should set a longer `initialDelaySeconds` on the liveness probe.

**Update strategy** -- Set `updateStrategy.type` to `BatchUpdate` (sequential batches) or `GrayBatchUpdate` (canary-style phased release). Control the number of batches, wait time between batches, and whether each batch requires manual approval (`releaseType: manual`) or proceeds automatically (`releaseType: auto`).

## Outputs and Dependencies

### What This Component Consumes

| Dependency | Field | ValueFromRef Path |
|------------|-------|-------------------|
| **AliCloudVpc** (optional) | `vpcId` | `status.outputs.vpc_id` |
| **AliCloudVswitch** (optional) | `vswitchId` | `status.outputs.vswitch_id` |
| **AliCloudSecurityGroup** (optional) | `securityGroupId` | `status.outputs.security_group_id` |

### What This Component Provides

After provisioning, `status.outputs` contains values that downstream Cloud Resources can consume via ValueFromRef:

| Output | Description | Common Downstream Use |
|--------|-------------|----------------------|
| `app_id` | The SAE application ID assigned by Alibaba Cloud | Monitoring dashboards, SLB backend bindings |
| `app_name` | The application name (mirrors the spec input) | Service discovery, operational references |

## Common Patterns

Browse the [Presets](#presets) tab for ready-to-deploy configurations.

**Container image production** -- A VPC-connected container application with 3 replicas, 4 vCPUs, 8 GB memory, HTTP liveness and readiness probes, and a 2-batch rolling update strategy. Start from the **Container Image Production** preset.

**Java FatJar production** -- A VPC-connected Spring Boot application with JDK 17, G1GC tuning, Spring Actuator health endpoints, and batch update strategy. Start from the **Java FatJar Production** preset.

**Container image development** -- A minimal single-replica container application with 0.5 vCPU and 1 GB memory for development and testing. Uses SAE managed networking (no VPC). Start from the **Container Image Development** preset.

## Works With

- [**AliCloud VPC**](/cloud-catalog/ali-cloud-vpc) -- provides the VPC for private network deployment
- [**AliCloud VSwitch**](/cloud-catalog/ali-cloud-vswitch) -- provides the subnet and availability zone placement
- [**AliCloud Security Group**](/cloud-catalog/ali-cloud-security-group) -- provides network access control for the application instances