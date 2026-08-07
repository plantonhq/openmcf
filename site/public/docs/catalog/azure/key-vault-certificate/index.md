---
title: "Key Vault Certificate"
description: "Key Vault Certificate deployment documentation"
icon: "package"
order: 100
componentName: "azurekeyvaultcertificate"
---

# Azure Key Vault Certificate

Deploys an X.509 certificate inside an Azure Key Vault -- the TLS building block the vault manages end to end: private key custody, enrollment or import, renewal, and expiry notification. A vault certificate is really three linked objects sharing one name: the certificate (public part + policy), its private KEY, and a SECRET whose value is the full bundle. That third face is the composition seam -- Azure services that terminate TLS consume the certificate THROUGH its secret identifier, which is why the `secret_id` outputs are what downstream kinds reference. The component integrates with Planton's Provider Connections for Azure credential management and ValueFromRef for dependency wiring.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **Key Vault Certificate** -- in the referenced vault, either vault-GENERATED (self-signed or CA-issued from a policy) or IMPORTED from an existing PFX/PEM bundle
- **Issuance policy** -- when configured: the issuer, the private-key shape (RSA/EC, optionally HSM-backed, exportable or vault-bound), near-expiry actions (auto-renew or contact notification), and the secret encoding (PKCS#12 or PEM)
- **X.509 content** -- for generated certificates: the subject, SANs, key usages, EKU OIDs, and validity window renewals re-issue for

## The Certificate in the Vault Family

- **AzureKeyVault** -- the parent security boundary, referenced by `keyVaultId`; its authorization mode, network posture, and deletion safety govern this certificate wholesale
- **TLS terminators** -- Application Gateway listeners and trusted roots, Front Door secrets, firewall-policy TLS inspection, and Container Apps environment certificates all consume this certificate's secret identifiers

## Before You Deploy

### Planton Setup

- **Azure Provider Connection** -- an active connection in the Connect module with credentials for the target Azure subscription. Map it as the default for your environment, or specify it explicitly when creating the Cloud Resource.
- **Planton Runner** -- required when using Runner-based credential delivery. Not needed for inline credentials or browser OAuth authentication modes.

### Azure Subscription

- **An Azure Key Vault** for the certificate to live in. Reference an AzureKeyVault Cloud Resource via ValueFromRef. HSM-backed private keys require the vault on the PREMIUM SKU.
- **Data-plane certificate permissions** for the deploying credential -- the "Key Vault Administrator" or "Key Vault Certificates Officer" RBAC role, or certificate permissions in a legacy access policy.
- **For imports** -- the base64-encoded PFX/PEM bundle stored as a Planton org secret (it carries the private key; the plaintext never enters the manifest).
- **For CA issuance** -- a CA issuer (DigiCert, GlobalSign) configured on the vault out-of-band.

## Deploy

### Console

Open the deployment store, find **Azure Key Vault Certificate**, and click **Deploy**. The creation wizard walks you through placement, the generate-vs-import source choice (imported bundles ride org-secret references), the issuance policy, and -- for generated certificates -- the X.509 content with live SAN guidance. Start from the **Self-Signed Auto-Renew** preset in the [Presets](#presets) tab.

### CLI

Create a manifest and apply it:

```yaml
apiVersion: azure.planton.dev/v1
kind: AzureKeyVaultCertificate
metadata:
  name: web-tls
  org: acme-corp
  env: prod
spec:
  name: web-tls
  keyVaultId:
    valueFrom:
      kind: AzureKeyVault
      name: platform-vault
      fieldPath: status.outputs.key_vault_id
  certificatePolicy:
    issuerName: Self
    keyProperties:
      exportable: true
      keyType: RSA
      keySize: 2048
    lifetimeActions:
      - actionType: AUTO_RENEW
        trigger:
          lifetimePercentage: 80
    secretProperties:
      contentType: PKCS12
    x509CertificateProperties:
      subject: CN=internal.example.com
      subjectAlternativeNames:
        dnsNames:
          - internal.example.com
      keyUsage:
        - DIGITAL_SIGNATURE
        - KEY_ENCIPHERMENT
      validityInMonths: 12
```

```shell
planton apply -f key-vault-certificate.yaml
```

This creates the fully-hands-off internal-TLS shape: self-signed, renewing itself at 80% lifetime, with consumers following each renewal through the versionless secret ID.

An import brings an existing bundle through an org-secret reference:

```yaml
spec:
  name: purchased-tls
  keyVaultId:
    valueFrom:
      kind: AzureKeyVault
      name: platform-vault
      fieldPath: status.outputs.key_vault_id
  certificate:
    contents: $secret/web-cert-bundle
    password: $secret/web-cert-password
```

### InfraChart

TLS terminators wire to the certificate through its secret outputs in the same InfraPipeline:

```yaml
spec:
  keyVaultSecretId:
    valueFrom:
      kind: AzureKeyVaultCertificate
      name: web-tls
      fieldPath: status.outputs.versionless_secret_id
```

The InfraPipeline resolves the dependency graph -- vault, then certificate, then the gateway that serves it.

## Key Configuration

These are the most important decisions when configuring a Key Vault certificate. Explore the full field reference in the [API Explorer](#api-explorer) tab.

**Generate vs import** -- generate (policy only) keeps the private key inside the vault from birth and enables automatic renewal; import (the `certificate` block) brings an existing PFX/PEM. An import may ALSO carry a policy to govern the imported material; without one, Azure derives it from the bundle.

**The issuer drives renewal** -- "Self" self-signs (internal TLS; auto-renew works); a configured CA issuer automates public issuance end to end; "Unknown" mints a CSR for an out-of-band CA -- renewal is inherently manual there, and only the contact-notification action applies.

**Exportability** -- whether the private key rides inside the secret face consumers read. TLS terminators that load the certificate themselves need `exportable: true`; keep it false when only vault-integrated services consume it.

**Near-expiry actions** -- one action per trigger, each firing on exactly ONE of lifetime-percentage or days-before-expiry. Auto-renew at 80% lifetime is the convention; actions update without minting a new certificate version.

**SANs over subject** -- modern TLS clients validate SANs, not the subject CN. Every hostname the certificate serves must appear in `dnsNames`.

## Outputs and Dependencies

### What This Component Consumes

| Dependency | Field | ValueFromRef Path |
|------------|-------|-------------------|
| **AzureKeyVault** | `keyVaultId` | `status.outputs.key_vault_id` |

### What This Component Provides

After provisioning, `status.outputs` contains values that downstream Cloud Resources can consume via ValueFromRef:

| Output | Description | Common Downstream Use |
|--------|-------------|----------------------|
| `versionless_secret_id` | Versionless SECRET-face ID (`https://{vault}.vault.azure.net/secrets/{name}`) | THE TLS-terminator reference -- Application Gateway and Container Apps environment certificates follow renewals through it |
| `secret_id` | Versioned secret-face ID of the current version | Consumers pinning one issued version (firewall-policy TLS inspection) |
| `certificate_id` | Versioned data-plane certificate ID | Front Door secrets and certificate-object consumers |
| `versionless_id` | Versionless data-plane certificate ID | Certificate-object references that follow renewal |
| `certificate_name` | The certificate's name in the vault | Operator tooling |
| `version` | The current version identifier | Audit trails |
| `thumbprint` | The SHA-1 thumbprint, hex-encoded | Pinning and inventory reconciliation |
| `resource_manager_id` | Versioned ARM resource ID | ARM-level references |
| `resource_manager_versionless_id` | Versionless ARM resource ID | ARM-level references that follow renewal |

## Common Patterns

Browse the [Presets](#presets) tab for ready-to-deploy configurations.

**Self-signed auto-renew** -- the fully-hands-off internal-TLS shape: the vault creates, signs, and renews forever. Start from the **Self-Signed Auto-Renew** preset.

**Imported bundle** -- an externally-purchased certificate brought in through org-secret references. Start from the **Imported Bundle** preset.

**CA pending CSR** -- the "Unknown"-issuer flow: the vault holds a CSR for an out-of-band CA, with contact notification before expiry. Start from the **CA Pending CSR** preset.

## Works With

- [**Azure Key Vault**](/cloud-catalog/azure-key-vault) -- the parent vault whose governance this certificate inherits
- [**Azure Application Gateway**](/cloud-catalog/azure-application-gateway) -- terminates TLS with this certificate through its secret ID
- [**Azure Key Vault Key**](/cloud-catalog/azure-key-vault-key) -- the sibling kind for raw cryptographic keys (CMK)
- [**Azure Role Assignment**](/cloud-catalog/azure-role-assignment) -- grants TLS consumers secret-read access on the vault
