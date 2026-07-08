---
title: "Front Door Secret"
description: "Front Door Secret deployment documentation"
icon: "package"
order: 100
componentName: "azurefrontdoorsecret"
---

# Azure Front Door Secret

Creates a secret inside an AzureFrontDoorProfile -- the bring-your-own TLS certificate node wrapping an AzureKeyVaultCertificate. AzureFrontDoorCustomDomain resources reference the secret to terminate TLS with the wrapped certificate; a versionless certificate reference makes Key Vault rotation propagate automatically.

## What Gets Created

When you deploy an AzureFrontDoorSecret resource, Planton provisions:

- **Front Door Secret** -- an `azurerm_cdn_frontdoor_secret` on the referenced profile, wrapping the referenced Key Vault certificate

## Prerequisites

- **Azure credentials** configured via environment variables or Planton provider config
- **An AzureFrontDoorProfile** to create the secret in (referenced through `profileId`)
- **An AzureKeyVaultCertificate** holding the certificate to serve
- **A one-time tenant grant**: Front Door's service principal (`Microsoft.AzureFrontDoor-Cdn`) needs Key Vault read access (e.g. the "Key Vault Secrets User" role) before the secret can deploy

## Quick Start

Create a file `front-door-secret.yaml`:

```yaml
apiVersion: azure.planton.dev/v1
kind: AzureFrontDoorSecret
metadata:
  name: wildcard-cert
  labels:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: dev.AzureFrontDoorSecret.wildcard-cert
spec:
  profileId:
    valueFrom:
      kind: AzureFrontDoorProfile
      name: my-front-door
      fieldPath: status.outputs.profile_id
  secretName: wildcard-example-com
  keyVaultCertificateId:
    valueFrom:
      kind: AzureKeyVaultCertificate
      name: my-wildcard-cert
      fieldPath: status.outputs.versionless_id
```

Deploy:

```shell
planton apply -f front-door-secret.yaml
```

The versionless reference (the default) makes Front Door follow the certificate's latest Key Vault version -- renewals propagate with zero redeploys. Reference the versioned `certificate_id` output instead to pin one exact version for change-controlled rollouts.

## Key Outputs

| Output | Purpose |
|--------|---------|
| `secret_id` | The ARM id -- what a custom domain's `tls.secretId` references |
| `secret_name` | The secret's name inside the profile |
| `subject_alternative_names` | The DNS names the wrapped certificate covers -- confirm a domain's hostname before attaching |

## Related Resources

- [Azure Front Door Profile](/docs/catalog/azure/front-door-profile) -- the parent profile
- [Azure Front Door Custom Domain](/docs/catalog/azure/front-door-custom-domain) -- terminates TLS with this secret
- [Azure Key Vault Certificate](/docs/catalog/azure/key-vault-certificate) -- the wrapped certificate
