---
title: "Container App Environment Certificate"
description: "Container App Environment Certificate deployment documentation"
icon: "package"
order: 100
componentName: "azurecontainerappenvironmentcertificate"
---

# Azure Container App Environment Certificate

Stores a bring-your-own TLS certificate on a Container App Environment -- uploaded as a PFX or pulled from Azure Key Vault -- for binding to app custom domains.

## What Gets Created

When you deploy an AzureContainerAppEnvironmentCertificate resource, Planton provisions:

- **Environment certificate** -- an `azurerm_container_app_environment_certificate` on the referenced environment, sourced from an inline PFX bundle or a Key Vault secret reference

## Prerequisites

- **Azure credentials** configured via environment variables or Planton provider config
- **An AzureContainerAppEnvironment** to store the certificate on
- For the Key Vault path: an **AzureKeyVaultCertificate** and read access to the vault's secrets for the environment's managed identity (Key Vault Secrets User under RBAC)

## Quick Start

Create a file `certificate.yaml`:

```yaml
apiVersion: azure.planton.dev/v1
kind: AzureContainerAppEnvironmentCertificate
metadata:
  name: app-tls-cert
  annotations:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: dev.AzureContainerAppEnvironmentCertificate.app-tls-cert
spec:
  certificateName: app.example.com
  containerAppEnvironmentId:
    valueFrom:
      kind: AzureContainerAppEnvironment
      name: my-env
      fieldPath: status.outputs.environment_id
  certificateBlobBase64: $secret/app-cert-pfx
  certificatePassword: $secret/app-cert-password
```

Deploy:

```shell
planton apply -f certificate.yaml
```

The blob and password are secret-bearing fields -- reference managed secrets rather than pasting values. Bind the certificate to an app's hostname with [Azure Container App Custom Domain](/docs/catalog/azure/container-app-custom-domain), referencing this resource's `certificate_id` output with `certificateBindingType: SNI_ENABLED`.

Everything except `tags` is ForceNew -- rotating an inline certificate replaces the resource and briefly re-binds the domains using it. The Key Vault source avoids that treadmill: a versionless secret reference follows renewals automatically.

## Key Outputs

| Output | Purpose |
|--------|---------|
| `certificate_id` | What custom-domain bindings reference |
| `subject_name` / `issuer` / `thumbprint` | Identity of the installed certificate |
| `issue_date` / `expiration_date` | Expiry monitoring (alarm on this for inline PFX certificates) |

## Related Resources

- [Azure Container App Custom Domain](/docs/catalog/azure/container-app-custom-domain) -- binds the certificate to an app hostname
- [Azure Container App Environment Managed Certificate](/docs/catalog/azure/container-app-environment-managed-certificate) -- the free, Azure-renewed alternative
- [Azure Key Vault Certificate](/docs/catalog/azure/key-vault-certificate) -- the Key Vault source this kind can track
