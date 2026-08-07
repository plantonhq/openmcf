---
title: "KMS Vault"
description: "KMS Vault deployment documentation"
icon: "package"
order: 100
componentName: "ocikmsvault"
---

# KMS Vault on OCI

Deploys an Oracle Cloud Infrastructure Key Management Service vault -- an HSM-backed container for encryption keys used by Compute, Block Volume, Object Storage, Database, and other OCI services. The vault exposes crypto and management endpoints that downstream OciKmsKey resources consume for key creation and cryptographic operations. Supports shared (Default), dedicated (Virtual Private), and external (BYOK/EKMS) vault types. The component integrates with Planton's Provider Connections for OCI credential management and supports ValueFromRef wiring to compartments.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **KMS Vault** -- a `kms.Vault` in the specified compartment with configurable vault type (shared HSM, dedicated HSM, or external key manager). The vault exposes crypto and management endpoints as outputs.
- **Freeform Tags** -- resource metadata tags (organization, environment, resource kind, resource ID) applied to the vault

## Before You Deploy

### Planton Setup

- **OCI Provider Connection** -- an active connection in the Connect module with credentials for the target OCI tenancy. Map it as the default for your environment, or specify it explicitly when creating the Cloud Resource.
- **Planton Runner** -- required when using Runner-based credential delivery. Not needed for inline credentials.

### OCI Tenancy

- A compartment to place the vault in. Provide the compartment OCID directly or reference an OciCompartment Cloud Resource via ValueFromRef.
- For external vault type only: IDCS OAuth credentials (client app ID, secret, and account URL) for connecting to the third-party key manager, and a pre-existing KMS private endpoint OCID for network connectivity.

## Deploy

### Console

Open the deployment store, find **KMS Vault on OCI**, and click **Deploy**. The creation wizard walks you through preset selection, environment and connection configuration, and spec fields. Start from the **Shared Vault** preset in the [Presets](#presets) tab to pre-populate a cost-effective vault with a shared HSM partition.

### CLI

```yaml
apiVersion: oci.planton.dev/v1
kind: OciKmsVault
metadata:
  name: platform-vault
  org: acme-corp
  env: prod
spec:
  compartmentId:
    value: "ocid1.compartment.oc1..example"
  vaultType: default_vault
```

```shell
planton apply -f vault.yaml
```

This creates a shared-HSM vault suitable for most workloads. The vault OCID, crypto endpoint, and management endpoint are exported as stack outputs. No display name is configured -- it defaults to `metadata.name`.

### InfraChart

When deploying as part of a multi-resource environment, use ValueFromRef to wire the vault to a compartment deployed in the same InfraPipeline:

```yaml
spec:
  compartmentId:
    valueFrom:
      kind: OciCompartment
      name: security-compartment
      fieldPath: status.outputs.compartmentId
```

The InfraPipeline resolves the dependency graph, deploys the compartment first, then provisions the vault with the resolved compartment OCID.

## Key Configuration

These are the most important decisions when configuring a KMS vault. Explore the full field reference in the [API Explorer](#api-explorer) tab.

**Vault type** -- The `vaultType` field is immutable after creation. `default_vault` uses a shared HSM partition (lower cost, suitable for most workloads). `virtual_private` allocates a dedicated HSM partition (higher throughput limits, required for high-volume cryptographic operations or compliance mandates). `external` connects to a third-party key manager via IDCS OAuth for BYOK/EKMS scenarios.

**External key manager** -- When `vaultType` is `external`, the `externalKeyManagerMetadata` message is required. All sub-fields (endpoint URL, OAuth credentials, private endpoint ID) are immutable after creation. This is enforced by CEL validation rules on the spec. The external vault pattern is for organizations that must retain key material outside OCI.

**Display name** -- The `displayName` field is optional. When omitted, the vault uses `metadata.name`. In practice, name the vault after the team, project, or environment it serves (e.g., `platform-prod`, `data-team-staging`).

## Outputs and Dependencies

### What This Component Consumes

| Dependency | Field | ValueFromRef Path |
|------------|-------|-------------------|
| **OciCompartment** | `compartmentId` | `status.outputs.compartmentId` |

### What This Component Provides

After provisioning, `status.outputs` contains values that downstream Cloud Resources can consume via ValueFromRef:

| Output | Description | Common Downstream Use |
|--------|-------------|----------------------|
| `vault_id` | OCID of the KMS vault | OciKmsKey, OciVaultSecret |
| `crypto_endpoint` | Service endpoint for cryptographic operations (encrypt, decrypt, sign, verify) | Applications performing client-side encryption via the OCI KMS API |
| `management_endpoint` | Service endpoint for key management operations (create, import, rotate keys) | OciKmsKey `managementEndpoint` field -- required to create keys in this vault |

## Common Patterns

Browse the [Presets](#presets) tab for ready-to-deploy configurations.

**Shared vault** -- A default vault with shared HSM partition providing FIPS 140-2 Level 3 certified key storage at lower cost. The recommended starting point for most workloads. Start from the **Shared Vault** preset.

**Dedicated vault** -- A virtual private vault with a dedicated HSM partition for high-throughput cryptographic operations or compliance regimes requiring dedicated hardware (PCI-DSS Level 1, HIPAA). Start from the **Dedicated Vault** preset.

## Works With

- [**Compartment on OCI**](/cloud-catalog/oci-compartment) -- provides the compartment that scopes this vault