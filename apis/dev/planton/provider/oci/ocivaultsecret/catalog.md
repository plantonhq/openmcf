# Vault Secret on OCI

Deploys an Oracle Cloud Infrastructure Vault Secret -- a named piece of sensitive data (credential, certificate, API key) stored in an OCI KMS Vault and encrypted by a master encryption key. Supports two mutually exclusive content modes: explicit content (user-provided base64 data) and auto-generation (OCI generates a passphrase, SSH key, or random bytes using a template). Lifecycle management includes secret expiry rules, content reuse policies, and scheduled rotation against Autonomous Database or OCI Functions targets. The component integrates with Planton's Provider Connections for OCI credential management and supports ValueFromRef wiring to compartments, vaults, and encryption keys.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **Vault Secret** -- a `vault.Secret` in the specified compartment, within the target vault, encrypted by the specified KMS key. The secret holds either user-provided content or auto-generated content depending on the configuration mode.
- **Freeform Tags** -- resource metadata tags (organization, environment, resource kind, resource ID) applied to the secret

## Before You Deploy

### Planton Setup

- **OCI Provider Connection** -- an active connection in the Connect module with credentials for the target OCI tenancy. Map it as the default for your environment, or specify it explicitly when creating the Cloud Resource.
- **Planton Runner** -- required when using Runner-based credential delivery. Not needed for inline credentials.

### OCI Tenancy

- A compartment to place the secret in. Provide the compartment OCID directly or reference an OciCompartment Cloud Resource via ValueFromRef.
- A KMS vault to contain the secret. Provide the vault OCID directly or reference an OciKmsVault Cloud Resource via ValueFromRef. Immutable after creation.
- A symmetric KMS encryption key within the vault. Provide the key OCID directly or reference an OciKmsKey Cloud Resource via ValueFromRef. Immutable after creation.
- For explicit content: base64-encoded secret data.
- For auto-generation: a generation template name and, for passphrases, a desired length.
- For rotation: an Autonomous Database OCID or Functions function OCID as the rotation target.

## Deploy

### Console

Open the deployment store, find **Vault Secret on OCI**, and click **Deploy**. The creation wizard walks you through preset selection, environment and connection configuration, and spec fields. Start from the **Explicit Credential** preset in the [Presets](#presets) tab to pre-populate a secret with user-provided content and a 90-day expiry rule.

### CLI

```yaml
apiVersion: oci.planton.dev/v1
kind: OciVaultSecret
metadata:
  name: app-credential
  org: acme-corp
  env: prod
spec:
  compartmentId:
    value: "ocid1.compartment.oc1..example"
  secretName: app-db-password
  vaultId:
    value: "ocid1.vault.oc1..example"
  keyId:
    value: "ocid1.key.oc1..example"
  secretContent:
    content: "cGFzc3dvcmQxMjM="
    stage: CURRENT
```

```shell
planton apply -f vault-secret.yaml
```

This creates a vault secret with explicit base64-encoded content. No expiry rules, reuse policies, or rotation are configured. The `secretName` is immutable after creation.

### InfraChart

When deploying as part of a multi-resource environment, use ValueFromRef to wire the secret to a compartment, vault, and encryption key deployed in the same InfraPipeline:

```yaml
spec:
  compartmentId:
    valueFrom:
      kind: OciCompartment
      name: security-compartment
      fieldPath: status.outputs.compartmentId
  vaultId:
    valueFrom:
      kind: OciKmsVault
      name: platform-vault
      fieldPath: status.outputs.vaultId
  keyId:
    valueFrom:
      kind: OciKmsKey
      name: secrets-key
      fieldPath: status.outputs.keyId
```

The InfraPipeline resolves the dependency graph, deploys the compartment, vault, and key first, then provisions the secret with the resolved values.

## Key Configuration

These are the most important decisions when configuring a vault secret. Explore the full field reference in the [API Explorer](#api-explorer) tab.

**Content mode** -- Two mutually exclusive modes. Set `secretContent` with base64-encoded data for explicit content. Alternatively, set `enableAutoGeneration: true` with a `secretGenerationContext` for OCI-generated content (passphrase, SSH key, or random bytes). CEL validation rules enforce mutual exclusivity -- setting both produces a validation error.

**Auto-generation** -- When `enableAutoGeneration` is `true`, set `secretGenerationContext.generationType` to `passphrase`, `ssh_key`, or `bytes`. For passphrases, `passphraseLength` is required. The `generationTemplate` names a provider-defined template (e.g., `SECRET_TPL_DBAAS_DEFAULT` for database credentials). An optional `secretTemplate` provides a structural template for the generated value.

**Secret rules** -- Add entries to `secretRules` for lifecycle policies. `secret_expiry_rule` controls version expiry via `secretVersionExpiryInterval` (ISO 8601, e.g., `P90D`) and optionally blocks content retrieval after expiry. `secret_reuse_rule` prevents reuse of previous content values, optionally enforced on deleted versions.

**Rotation** -- Set `rotationConfig` for scheduled rotation against a target system. `targetSystemDetails.targetSystemType` is `adb` (Autonomous Database) or `function` (OCI Functions). When `isScheduledRotationEnabled` is `true`, `rotationInterval` (ISO 8601, e.g., `P30D`) controls the rotation frequency. The target system's credentials are updated automatically during rotation.

**Immutable fields** -- `secretName`, `vaultId`, and `keyId` are all ForceNew -- changing any of them destroys and recreates the secret. Choose these values carefully at creation time.

## Outputs and Dependencies

### What This Component Consumes

| Dependency | Field | ValueFromRef Path |
|------------|-------|-------------------|
| **OciCompartment** | `compartmentId` | `status.outputs.compartmentId` |
| **OciKmsVault** | `vaultId` | `status.outputs.vaultId` |
| **OciKmsKey** | `keyId` | `status.outputs.keyId` |
| **OciAutonomousDatabase** (optional, rotation) | `rotationConfig.targetSystemDetails.adbId` | `status.outputs.autonomousDatabaseId` |

### What This Component Provides

After provisioning, `status.outputs` contains values that downstream Cloud Resources can consume via ValueFromRef:

| Output | Description | Common Downstream Use |
|--------|-------------|----------------------|
| `secret_id` | OCID of the vault secret | Application configuration, IAM policy scoping, rotation target references |
| `current_version_number` | Version number of the currently active secret version | Version tracking, rotation monitoring |

## Common Patterns

Browse the [Presets](#presets) tab for ready-to-deploy configurations.

**Explicit credential** -- A secret with user-provided base64 content, a 90-day version expiry rule that blocks retrieval after expiry. The standard pattern for manually managed credentials (API keys, third-party tokens). Start from the **Explicit Credential** preset.

**Auto-generated passphrase** -- A secret with OCI-generated passphrase content, scheduled 30-day rotation against an Autonomous Database, and a reuse rule preventing content recycling. Designed for database credentials with fully automated lifecycle management. Start from the **Auto-Generated Passphrase** preset.

## Works With

- [**Compartment on OCI**](/cloud-catalog/oci-compartment) -- provides the compartment that scopes this secret
- [**KMS Vault on OCI**](/cloud-catalog/oci-kms-vault) -- provides the vault that contains this secret
- [**KMS Key on OCI**](/cloud-catalog/oci-kms-key) -- provides the master encryption key that protects this secret
- [**Autonomous Database on OCI**](/cloud-catalog/oci-autonomous-database) -- provides the rotation target for database credential secrets