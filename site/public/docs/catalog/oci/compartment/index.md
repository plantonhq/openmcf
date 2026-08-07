---
title: "Compartment"
description: "Compartment deployment documentation"
icon: "package"
order: 100
componentName: "ocicompartment"
---

# Compartment on OCI

Deploys an Oracle Cloud Infrastructure compartment for hierarchical resource isolation and access control. Compartments are OCI's fundamental organizational primitive -- every resource exists within exactly one compartment. This component creates a single compartment within a parent compartment or tenancy, and nested hierarchies are built by chaining OciCompartment resources where each child references its parent via `compartmentId`. The component integrates with Planton's Provider Connections for OCI credential management and supports ValueFromRef wiring for parent compartment resolution.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **Identity Compartment** -- an `oci_identity_compartment` within the specified parent compartment or tenancy, with name, description, and configurable delete protection
- **Freeform Tags** -- resource metadata tags (organization, environment, resource kind, resource ID) applied to the compartment

## Before You Deploy

### Planton Setup

- **OCI Provider Connection** -- an active connection in the Connect module with credentials for the target OCI tenancy. Map it as the default for your environment, or specify it explicitly when creating the Cloud Resource.
- **Planton Runner** -- required when using Runner-based credential delivery. Not needed for inline credentials.

### OCI Tenancy

- A parent compartment OCID or the tenancy OCID where this compartment will be created. For top-level compartments, use the tenancy OCID. For nested compartments, provide the parent compartment OCID directly or reference another OciCompartment Cloud Resource via ValueFromRef.

## Deploy

### Console

Open the deployment store, find **Compartment on OCI**, and click **Deploy**. The creation wizard walks you through preset selection, environment and connection configuration, and spec fields. Start from the **Project** preset in the [Presets](#presets) tab to pre-populate a long-lived compartment with delete protection enabled.

### CLI

```yaml
apiVersion: oci.planton.dev/v1
kind: OciCompartment
metadata:
  name: platform-compartment
  org: acme-corp
  env: prod
spec:
  compartmentId:
    value: "ocid1.tenancy.oc1..example"
  description: "Platform team infrastructure and shared services"
```

```shell
planton apply -f compartment.yaml
```

This creates a compartment named `platform-compartment` under the tenancy root. Delete protection is enabled by default -- destroying the IaC resource does not delete the compartment from OCI.

### InfraChart

When deploying nested compartments as part of a multi-resource environment, use ValueFromRef to wire a child compartment to its parent:

```yaml
spec:
  compartmentId:
    valueFrom:
      kind: OciCompartment
      name: platform-compartment
      fieldPath: status.outputs.compartmentId
```

The InfraPipeline resolves the dependency graph, deploys the parent compartment first, then provisions the child compartment with the resolved OCID.

## Key Configuration

These are the most important decisions when configuring a compartment. Explore the full field reference in the [API Explorer](#api-explorer) tab.

**Delete protection** -- The `enableDelete` flag defaults to `false`, meaning the compartment is retained in OCI even when the IaC resource is destroyed. This is OCI's safety mechanism to prevent accidental deletion of compartments containing active resources. Set to `true` only for ephemeral or development compartments. To delete a protected compartment, first set `enableDelete: true`, apply, then destroy -- this two-step process is intentional friction.

**Compartment name** -- The `name` field falls back to `metadata.name` if not provided. The name must be unique among siblings within the parent compartment. It appears in the OCI Console and is referenced in IAM policy statements.

**Description** -- Required by the OCI API. Use it to document what resources or teams the compartment is intended for. This is the first thing operators see when navigating the compartment hierarchy in the OCI Console.

## Outputs and Dependencies

### What This Component Consumes

| Dependency | Field | ValueFromRef Path |
|------------|-------|-------------------|
| **OciCompartment** | `compartmentId` | `status.outputs.compartmentId` |

### What This Component Provides

After provisioning, `status.outputs` contains values that downstream Cloud Resources can consume via ValueFromRef:

| Output | Description | Common Downstream Use |
|--------|-------------|----------------------|
| `compartment_id` | OCID of the created compartment | OciVcn, OciSubnet, OciSecurityGroup, OciIdentityPolicy, OciDynamicGroup, OciKmsVault, OciKmsKey, and virtually every other OCI component |

## Common Patterns

Browse the [Presets](#presets) tab for ready-to-deploy configurations.

**Project compartment** -- A long-lived compartment for a project, team, or workload with delete protection enabled. The standard choice for production, staging, and shared-services compartments. Start from the **Project** preset.

**Sandbox compartment** -- An ephemeral compartment for development, CI/CD pipelines, or proof-of-concept work with `enableDelete: true` for automated teardown. Start from the **Sandbox** preset.

## Works With

- [**Compartment on OCI**](/cloud-catalog/oci-compartment) -- provides the parent compartment for nested hierarchies