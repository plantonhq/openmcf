---
title: "Dynamic Group"
description: "Dynamic Group deployment documentation"
icon: "package"
order: 100
componentName: "ocidynamicgroup"
---

# Dynamic Group on OCI

Deploys an Oracle Cloud Infrastructure dynamic group for enabling workload identity -- the mechanism that lets compute instances, OKE pods, and Functions authenticate to OCI services without stored credentials. A dynamic group uses a matching rule to select which OCI resources are members, and combined with an OciIdentityPolicy, enables the credential-less authentication pattern: the dynamic group defines *who* (matching rule), and the policy defines *what they can do* (statements). The component integrates with Planton's Provider Connections for OCI credential management and supports ValueFromRef wiring to the tenancy compartment.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **Identity Dynamic Group** -- an `oci_identity_dynamic_group` in the tenancy with the provided name, description, and matching rule. The dynamic group name defaults to `metadata.name` if not explicitly set in the spec.
- **Freeform Tags** -- resource metadata tags (organization, environment, resource kind, resource ID) applied to the dynamic group

## Before You Deploy

### Planton Setup

- **OCI Provider Connection** -- an active connection in the Connect module with credentials for the target OCI tenancy. Map it as the default for your environment, or specify it explicitly when creating the Cloud Resource.
- **Planton Runner** -- required when using Runner-based credential delivery. Not needed for inline credentials.

### OCI Tenancy

- The tenancy OCID (root compartment). Dynamic groups are tenancy-level IAM resources and must be created in the tenancy root compartment, not in a child compartment. Provide the tenancy OCID directly or reference an OciCompartment Cloud Resource that represents the tenancy root via ValueFromRef.
- Knowledge of OCI's matching rule syntax: `Any {condition}` matches resources satisfying any single condition, `All {condition, condition}` requires all conditions. Common conditions include `instance.compartment.id`, `resource.type`, and `tag.<namespace>.<key>.value`.
- A companion OciIdentityPolicy is needed to grant the dynamic group's members actual permissions -- without a policy, group membership alone grants no access.

## Deploy

### Console

Open the deployment store, find **Dynamic Group on OCI**, and click **Deploy**. The creation wizard walks you through preset selection, environment and connection configuration, and spec fields. Start from the **Compute Instance Principal** preset in the [Presets](#presets) tab to pre-populate a dynamic group matching all compute instances in a compartment.

### CLI

```yaml
apiVersion: oci.planton.dev/v1
kind: OciDynamicGroup
metadata:
  name: compute-workers
  org: acme-corp
  env: prod
spec:
  compartmentId:
    value: "ocid1.tenancy.oc1..example"
  description: "All compute instances in the production compartment"
  matchingRule: "Any {instance.compartment.id = 'ocid1.compartment.oc1..production'}"
```

```shell
planton apply -f dynamic-group.yaml
```

This creates a dynamic group named `compute-workers` in the tenancy. All compute instances in the specified compartment automatically become members. To grant these members permissions, create a companion OciIdentityPolicy with statements referencing this dynamic group.

### InfraChart

When deploying as part of a multi-resource environment, use ValueFromRef to wire the dynamic group to the tenancy compartment:

```yaml
spec:
  compartmentId:
    valueFrom:
      kind: OciCompartment
      name: tenancy-root
      fieldPath: status.outputs.compartmentId
```

The InfraPipeline resolves the dependency graph, deploys the compartment first, then provisions the dynamic group with the resolved tenancy OCID.

## Key Configuration

These are the most important decisions when configuring a dynamic group. Explore the full field reference in the [API Explorer](#api-explorer) tab.

**Matching rule** -- The `matchingRule` defines which OCI resources are members of this group. Use `Any {condition}` when a resource needs to satisfy only one condition (e.g., all instances in a compartment). Use `All {condition, condition}` when multiple conditions must all be met (e.g., Functions of a specific type in a specific compartment). Common conditions: `instance.compartment.id` for compute instances, `resource.type = 'fnfunc'` for Functions, `tag.<namespace>.<key>.value` for tag-based matching.

**Tenancy-level placement** -- The `compartmentId` field must be the tenancy OCID, not a child compartment OCID. Dynamic groups are tenancy-level IAM resources in OCI. The matching rule itself can reference any compartment to scope membership, but the group resource lives at the tenancy root.

**Name immutability** -- The `name` field (or `metadata.name` fallback) cannot be changed after creation and must be unique across all groups (including user groups) in the tenancy. Choose a descriptive, stable name. Dynamic group names appear in policy statements (`Allow dynamic-group <name> to ...`), so renaming requires updating all referencing policies.

## Outputs and Dependencies

### What This Component Consumes

| Dependency | Field | ValueFromRef Path |
|------------|-------|-------------------|
| **OciCompartment** | `compartmentId` | `status.outputs.compartmentId` |

### What This Component Provides

After provisioning, `status.outputs` contains values that downstream Cloud Resources can consume via ValueFromRef:

| Output | Description | Common Downstream Use |
|--------|-------------|----------------------|
| `dynamic_group_id` | OCID of the created dynamic group | Policy auditing, resource management |

## Common Patterns

Browse the [Presets](#presets) tab for ready-to-deploy configurations.

**Compute instance principal** -- A dynamic group matching all compute instances in a compartment for instance principal authentication. Pairs with an OciIdentityPolicy granting the group access to services like Vault, KMS, and Object Storage. Start from the **Compute Instance Principal** preset.

**Functions workload identity** -- A dynamic group matching all OCI Functions in a compartment for serverless workload identity. Uses `All {resource.type = 'fnfunc', resource.compartment.id = '...'}` to restrict membership to Functions only. Start from the **Functions Workload Identity** preset.

## Works With

- [**Compartment on OCI**](/cloud-catalog/oci-compartment) -- provides the tenancy OCID referenced by `compartmentId`; matching rules typically reference compartment OCIDs to scope group membership