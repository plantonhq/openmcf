# AzureContainerAppEnvironmentCertificate

Store a bring-your-own TLS certificate on a Container App Environment, shared by every app in it and bound to app custom domains through `AzureContainerAppCustomDomain`.

## Overview

Certificates live on the ENVIRONMENT: upload once, bind to any number of app domains. The certificate arrives one of two ways -- set exactly one:

1. **Inline PFX** -- a base64-encoded PKCS#12 bundle (certificate chain + private key) with its password. Rotation is manual.
2. **Key Vault reference** -- the environment pulls the certificate from Azure Key Vault and follows renewals automatically when the reference is versionless.

For certificates Azure should provision and renew end to end (free, domain-validated), use `AzureContainerAppEnvironmentManagedCertificate` instead -- this kind is for certificates you bring: EV/OV chains, org-mandated CAs, or wildcard certificates.

## Key Features

- **Two sources**: inline PFX (with sensitive handling for the blob and password) or Key Vault (versionless secret reference; renewals propagate)
- **Certificate facts as outputs**: subject, issuer, issue/expiry dates, thumbprint -- read back from the uploaded material for expiry monitoring
- **Composable**: the Key Vault path defaults to an `AzureKeyVaultCertificate`'s versionless secret face; the identity to a UAI's ARM id

## When to Use

- Custom domains with EV/OV or org-mandated CA certificates
- Wildcard certificates covering many app hostnames (managed certificates cannot do wildcards)
- Certificates already managed in Key Vault that Container Apps should track

## Spec Highlights

| Field | Notes |
| --- | --- |
| `certificate_name` | How bindings and the portal refer to it (commonly the domain). ForceNew |
| `container_app_environment_id` | The environment, by ARM id reference. ForceNew |
| `certificate_blob_base64` + `certificate_password` | The inline PFX pair (both sensitive). ForceNew |
| `certificate_key_vault` | Versionless secret reference + the identity that reads it ("System" or a UAI). ForceNew |
| `tags` | The only field Azure updates in place |

## Outputs

| Output | Purpose |
| --- | --- |
| `certificate_id` | The binding seam `AzureContainerAppCustomDomain` consumes |
| `subject_name`, `issuer`, `issue_date`, `expiration_date`, `thumbprint` | Certificate facts for monitoring |

## Quick Example

```yaml
apiVersion: azure.planton.dev/v1
kind: AzureContainerAppEnvironmentCertificate
metadata:
  name: app-tls-cert
spec:
  certificateName: app.example.com
  containerAppEnvironmentId:
    valueFrom:
      kind: AzureContainerAppEnvironment
      name: my-env
      fieldPath: status.outputs.environment_id
  certificateKeyVault:
    keyVaultSecretId:
      valueFrom:
        kind: AzureKeyVaultCertificate
        name: my-cert
        fieldPath: status.outputs.versionless_secret_id
```

The Key Vault path requires the environment's managed identity to hold read access to the vault's secrets (Key Vault Secrets User under RBAC) -- compose an `AzureRoleAssignment` for it.

---

© Planton. Licensed under [Apache-2.0](https://github.com/plantonhq/planton/blob/main/LICENSE).
