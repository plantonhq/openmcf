---
title: "Front Door Secret"
description: "Front Door Secret deployment documentation"
icon: "package"
order: 100
componentName: "azurefrontdoorsecret"
---

# Azure Front Door Secret

Deploys a Front Door secret -- the bring-your-own TLS certificate node inside an Azure Front Door (Standard/Premium) profile. A secret wraps a Key Vault certificate so custom domains can terminate TLS with it: the domain's `tls.secretId` references this secret, and this secret references the AzureKeyVaultCertificate that actually holds the key material. One certificate -- typically a wildcard or multi-SAN cert -- serves many domains, and rotating it is a single operation, never a per-domain edit. The component integrates with Planton's Provider Connections for Azure credential management and ValueFromRef for dependency wiring to the profile and the certificate.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **Front Door Secret** -- a named child of the profile wrapping the referenced Key Vault certificate
- **Certificate binding** -- versionless (Front Door follows the certificate's latest version; Key Vault rotation propagates automatically) or version-pinned (one exact certificate ships until the secret is replaced)

## The Secret in the Front Door Family

The secret is the TLS material bridge between Key Vault and the edge:

- **AzureFrontDoorProfile** -- the parent container, referenced by `profileId`
- **AzureKeyVaultCertificate** -- the certificate the secret wraps, referenced by `keyVaultCertificateId` (versionless by default)
- **AzureFrontDoorCustomDomain** -- terminates TLS with this secret via `tls.secretId` when its certificate type is CUSTOMER_CERTIFICATE

## Before You Deploy

### Planton Setup

- **Azure Provider Connection** -- an active connection in the Connect module with credentials for the target Azure subscription.
- **Planton Runner** -- required when using Runner-based credential delivery.

### Azure Subscription

- **An Azure Front Door Profile** the secret nests under.
- **A CA-issued Key Vault certificate with its complete chain** -- Azure rejects self-signed certificates ("the certificate chain includes an invalid number of certificates"). Use a certificate enrolled through a CA integration or an imported PKCS#12 carrying leaf plus issuer.
- **The one-time vault grant** -- Front Door reads Key Vault with Microsoft's own service principal (the `Microsoft.AzureFrontDoor-Cdn` enterprise application). Grant it read access on the vault -- e.g. the "Key Vault Secrets User" role on an RBAC-mode vault -- before the first secret deploys.

## Deploy

### Console

Open the deployment store, find **Azure Front Door Secret**, and click **Deploy**. The wizard walks you through the parent profile, the secret name, and the Key Vault certificate -- whose reference form (versionless vs versioned) decides the rotation story. Start from the **Rotating BYO Certificate** preset in the [Presets](#presets) tab.

### CLI

```yaml
apiVersion: azure.planton.dev/v1
kind: AzureFrontDoorSecret
metadata:
  name: wildcard-cert-secret
  org: acme-corp
  env: prod
spec:
  profileId:
    valueFrom:
      kind: AzureFrontDoorProfile
      name: cdn-profile
      fieldPath: status.outputs.profile_id
  secretName: wildcard-example-com
  keyVaultCertificateId:
    valueFrom:
      kind: AzureKeyVaultCertificate
      name: wildcard-example-com-cert
      fieldPath: status.outputs.versionless_id
```

```shell
planton apply -f front-door-secret.yaml
```

This creates a secret following the certificate's latest version. Custom domains reference the `secret_id` output next.

### InfraChart

```yaml
spec:
  profileId:
    valueFrom:
      kind: AzureFrontDoorProfile
      name: cdn-profile
      fieldPath: status.outputs.profile_id
  secretName: wildcard-example-com
  keyVaultCertificateId:
    valueFrom:
      kind: AzureKeyVaultCertificate
      name: wildcard-example-com-cert
      fieldPath: status.outputs.versionless_id
```

## Key Configuration

**Secret name** -- 2–260 characters; letters, digits, and hyphens; must start and end with a letter or digit. Unique within the profile. The name is a segment of the secret's ARM ID (the exact string domains reference), so renaming replaces the secret under a new ID.

**Key Vault certificate** -- referenced by its Key Vault certificate identifier (a vault data-plane URL, not an ARM ID). The versionless form (no trailing version segment, the default reference) tells Front Door to follow the certificate's LATEST version; a versioned identifier pins one exact certificate.

**Immutability** -- Azure exposes no update on Front Door secrets: changing any field replaces the secret. That is safe in practice, because certificate ROTATION happens inside Key Vault (new versions), not by editing the secret.

## Outputs and Dependencies

### What This Component Consumes

| Dependency | Field | ValueFromRef Path |
|------------|-------|-------------------|
| **AzureFrontDoorProfile** | `profileId` | `status.outputs.profile_id` |
| **AzureKeyVaultCertificate** | `keyVaultCertificateId` | `status.outputs.versionless_id` |

### What This Component Provides

| Output | Description | Common Downstream Use |
|--------|-------------|----------------------|
| `secret_id` | ARM resource ID of the secret | AzureFrontDoorCustomDomain.`tls.secretId` |
| `secret_name` | The secret's name within its profile | Operator tooling |
| `subject_alternative_names` | The DNS names the wrapped certificate covers | Confirming a domain's host name is covered before attaching |

## Presets

| Preset | Rank | Description |
|--------|------|-------------|
| Rotating BYO Certificate | 1 | Versionless reference -- Key Vault rotation propagates automatically |
| Pinned Certificate Version | 2 | Versioned reference -- one exact certificate until the secret is replaced |

## Related Components

- **AzureFrontDoorProfile** -- the parent container
- **AzureKeyVaultCertificate** -- holds the key material this secret wraps
- **AzureFrontDoorCustomDomain** -- terminates TLS with this secret
