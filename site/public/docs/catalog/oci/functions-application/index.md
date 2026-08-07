---
title: "Functions Application"
description: "Functions Application deployment documentation"
icon: "package"
order: 100
componentName: "ocifunctionsapplication"
---

# Functions Application on OCI

Deploys an Oracle Cloud Infrastructure Functions Application -- the organizational container and shared execution environment for serverless functions. The application defines the network placement (subnets and security groups), processor architecture, shared configuration environment variables, optional container image signature verification via KMS, and APM tracing. Individual functions are deployed as code artifacts via `fn deploy` or CI/CD pipelines, not managed by this component. The application integrates with Planton's Provider Connections for OCI credential management and supports ValueFromRef wiring to compartments, subnets, security groups, and KMS keys.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **Functions Application** -- a serverless function container in the specified compartment with subnet placement, processor architecture, optional NSG restrictions, shared config map, optional image signature verification policy, and optional APM tracing
- **Freeform Tags** -- resource metadata tags (organization, environment, resource kind, resource ID) applied to the application

## Before You Deploy

### Planton Setup

- **OCI Provider Connection** -- an active connection in the Connect module with credentials for the target OCI tenancy. Map it as the default for your environment, or specify it explicitly when creating the Cloud Resource.
- **Planton Runner** -- required when using Runner-based credential delivery. Not needed for inline credentials.

### OCI Tenancy

- A compartment to place the application in. Provide the compartment OCID directly or reference an OciCompartment Cloud Resource via ValueFromRef.
- At least one private subnet for function execution. Functions run inside these subnets and can reach any resources accessible from them. Subnets are immutable after creation.
- Optionally, one or more network security groups to restrict inbound and outbound traffic for functions.
- For image signature verification: one or more OCI KMS keys used to sign container images.
- For distributed tracing: an OCI APM domain OCID.

## Deploy

### Console

Open the deployment store, find **Functions Application on OCI**, and click **Deploy**. The creation wizard walks you through preset selection, environment and connection configuration, and spec fields. Start from the **Standard x86** preset in the [Presets](#presets) tab to pre-populate an application with x86 architecture, a single subnet, and APM tracing enabled.

### CLI

```yaml
apiVersion: oci.planton.dev/v1
kind: OciFunctionsApplication
metadata:
  name: order-processing
  org: acme-corp
  env: prod
spec:
  compartmentId:
    value: "ocid1.compartment.oc1..example"
  subnetIds:
    - value: "ocid1.subnet.oc1..example"
  shape: generic_x86
```

```shell
planton apply -f functions-app.yaml
```

This creates a Functions Application with x86 architecture placed in a single subnet. No NSGs, image policy, APM tracing, or shared config are configured.

### InfraChart

When deploying as part of a multi-resource environment, use ValueFromRef to wire the application to a compartment, subnet, and security group deployed in the same InfraPipeline:

```yaml
spec:
  compartmentId:
    valueFrom:
      kind: OciCompartment
      name: app-compartment
      fieldPath: status.outputs.compartmentId
  subnetIds:
    - valueFrom:
        kind: OciSubnet
        name: functions-subnet
        fieldPath: status.outputs.subnetId
  networkSecurityGroupIds:
    - valueFrom:
        kind: OciSecurityGroup
        name: functions-nsg
        fieldPath: status.outputs.networkSecurityGroupId
```

The InfraPipeline resolves the dependency graph, deploys the compartment, subnet, and security group first, then provisions the application with the resolved values.

## Key Configuration

These are the most important decisions when configuring a Functions Application. Explore the full field reference in the [API Explorer](#api-explorer) tab.

**Processor architecture** -- Set `shape` to `generic_x86` for standard Intel/AMD execution, `generic_arm` for Ampere A1 (lower cost per invocation), or `generic_x86_arm` for multi-architecture support. The shape is immutable after creation -- changing it forces application recreation.

**Subnet placement** -- Provide at least one subnet OCID in `subnetIds`. Functions execute within these subnets and inherit their routing and connectivity. Use private subnets with a NAT gateway for outbound internet access. Subnets are immutable after creation.

**Image signature verification** -- Set `imagePolicyConfig.isPolicyEnabled` to `true` and provide KMS key OCIDs in `keyDetails` to enforce that only signed container images can be deployed as functions. At least one key is required when the policy is enabled.

**Shared configuration** -- Use `config` to pass environment variables to all functions in the application. Keys must be ASCII letters, digits, and underscores (no leading digit). Total key-value size is capped at 4 KB.

## Outputs and Dependencies

### What This Component Consumes

| Dependency | Field | ValueFromRef Path |
|------------|-------|-------------------|
| **OciCompartment** | `compartmentId` | `status.outputs.compartmentId` |
| **OciSubnet** | `subnetIds` | `status.outputs.subnetId` |
| **OciSecurityGroup** (optional) | `networkSecurityGroupIds` | `status.outputs.networkSecurityGroupId` |
| **OciKmsKey** (optional) | `imagePolicyConfig.keyDetails[].kmsKeyId` | `status.outputs.keyId` |

### What This Component Provides

After provisioning, `status.outputs` contains values that downstream Cloud Resources can consume via ValueFromRef:

| Output | Description | Common Downstream Use |
|--------|-------------|----------------------|
| `application_id` | OCID of the functions application | Function deployment targets (`fn deploy`), IAM policy scoping, monitoring configuration |

## Common Patterns

Browse the [Presets](#presets) tab for ready-to-deploy configurations.

**Standard x86** -- An application with Intel/AMD architecture, a single subnet, NSG restriction, and APM tracing enabled. The default starting point for most serverless workloads. Start from the **Standard x86** preset.

**Secure production** -- An application with image signature verification enforced via KMS, NSG-restricted networking, and APM tracing. Designed for regulated environments where only signed container images may run. Start from the **Secure Production** preset.

## Works With

- [**Compartment on OCI**](/cloud-catalog/oci-compartment) -- provides the compartment that scopes this application
- [**Subnet on OCI**](/cloud-catalog/oci-subnet) -- provides the network subnets where functions execute
- [**Network Security Group on OCI**](/cloud-catalog/oci-security-group) -- provides network security rules for function traffic
- [**KMS Key on OCI**](/cloud-catalog/oci-kms-key) -- provides signing keys for container image verification