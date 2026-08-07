---
title: "Identity Policy"
description: "Identity Policy deployment documentation"
icon: "package"
order: 100
componentName: "ociidentitypolicy"
---

# Identity Policy on OCI

Deploys an Oracle Cloud Infrastructure IAM policy for granting access to compartment resources. Each policy contains one or more human-readable statements written in OCI's policy language and is attached to a compartment, granting permissions within that compartment and all of its children. Tenancy-level policies can be created by attaching the policy to the tenancy root compartment. The component integrates with Planton's Provider Connections for OCI credential management and supports ValueFromRef wiring to compartments.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **Identity Policy** -- an `oci_identity_policy` in the specified compartment with name, description, and one or more policy statements. The policy name defaults to `metadata.name` if not explicitly set in the spec.
- **Freeform Tags** -- resource metadata tags (organization, environment, resource kind, resource ID) applied to the policy

## Before You Deploy

### Planton Setup

- **OCI Provider Connection** -- an active connection in the Connect module with credentials for the target OCI tenancy. Map it as the default for your environment, or specify it explicitly when creating the Cloud Resource.
- **Planton Runner** -- required when using Runner-based credential delivery. Not needed for inline credentials.

### OCI Tenancy

- A compartment OCID or tenancy OCID where the policy will be created. Tenancy-level policies grant access across all compartments. Compartment-scoped policies grant access within that compartment and its children. Provide the OCID directly or reference an OciCompartment Cloud Resource via ValueFromRef.
- Knowledge of OCI's policy language syntax: `Allow <subject> to <verb> <resource-type> in <location> [where <conditions>]`. Subjects can be groups (`group GroupName`) or dynamic groups (`dynamic-group DynGroupName`).

## Deploy

### Console

Open the deployment store, find **Identity Policy on OCI**, and click **Deploy**. The creation wizard walks you through preset selection, environment and connection configuration, and spec fields. Start from the **Compartment Admin** preset in the [Presets](#presets) tab to pre-populate a policy granting a group full administrative access.

### CLI

```yaml
apiVersion: oci.planton.dev/v1
kind: OciIdentityPolicy
metadata:
  name: platform-admin-policy
  org: acme-corp
  env: prod
spec:
  compartmentId:
    value: "ocid1.compartment.oc1..example"
  description: "Grants administrative access to the platform compartment"
  statements:
    - "Allow group PlatformAdmins to manage all-resources in compartment platform"
```

```shell
planton apply -f policy.yaml
```

This creates an IAM policy named `platform-admin-policy` attached to the specified compartment with a single statement granting full administrative access. The `versionDate` field is not set, so the policy evaluates using current service behavior.

### InfraChart

When deploying as part of a multi-resource environment, use ValueFromRef to wire the policy to a compartment deployed in the same InfraPipeline:

```yaml
spec:
  compartmentId:
    valueFrom:
      kind: OciCompartment
      name: platform-compartment
      fieldPath: status.outputs.compartmentId
```

The InfraPipeline resolves the dependency graph, deploys the compartment first, then provisions the policy with the resolved compartment OCID.

## Key Configuration

These are the most important decisions when configuring an identity policy. Explore the full field reference in the [API Explorer](#api-explorer) tab.

**Policy statements** -- Provide at least one statement in `statements`. Each statement follows OCI's policy language: `Allow <subject> to <verb> <resource-type> in <location>`. Use `manage all-resources` for administrative access, `inspect all-resources` for read-only auditing, or target specific resource families (`virtual-network-family`, `secret-family`, `object-family`) for least-privilege grants.

**Policy scope** -- Attach the policy to the compartment it governs. Policies grant access within their compartment and all children. For tenancy-wide policies (e.g., auditor access across all compartments), set `compartmentId` to the tenancy OCID and use `in tenancy` in statements instead of `in compartment <name>`.

**Policy name immutability** -- The `name` field (or `metadata.name` fallback) cannot be changed after creation. Choose a descriptive, stable name. Renaming requires destroying and recreating the policy.

**Version date** -- Set `versionDate` (YYYY-MM-DD format) to pin policy evaluation to OCI service behavior on that date. This prevents future OCI service changes from altering how statements are interpreted. Omit for policies that should follow current behavior.

## Outputs and Dependencies

### What This Component Consumes

| Dependency | Field | ValueFromRef Path |
|------------|-------|-------------------|
| **OciCompartment** | `compartmentId` | `status.outputs.compartmentId` |

### What This Component Provides

After provisioning, `status.outputs` contains values that downstream Cloud Resources can consume via ValueFromRef:

| Output | Description | Common Downstream Use |
|--------|-------------|----------------------|
| `policy_id` | OCID of the created policy | Policy auditing and management |

## Common Patterns

Browse the [Presets](#presets) tab for ready-to-deploy configurations.

**Compartment admin** -- A policy granting a group full administrative access to all resources within a compartment. The most common starting point for team-level access. Start from the **Compartment Admin** preset.

**Dynamic group service access** -- A policy granting a dynamic group access to specific OCI services for workload identity. Used with OciDynamicGroup to enable credential-less authentication for compute instances and OKE pods. Start from the **Service Access** preset.

**Read-only auditor** -- A tenancy-level policy granting a group inspect-level visibility across all compartments for compliance and security auditing. Start from the **Read-Only Auditor** preset.

## Works With

- [**Compartment on OCI**](/cloud-catalog/oci-compartment) -- provides the compartment that scopes this policy and determines which resources it governs