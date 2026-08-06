---
title: "Container App Custom Domain"
description: "Container App Custom Domain deployment documentation"
icon: "package"
order: 100
componentName: "azurecontainerappcustomdomain"
---

# Azure Container App Custom Domain

Binds a custom domain to a Container App -- your own hostname serving the app, with TLS from an Azure-managed or bring-your-own certificate.

## What Gets Created

When you deploy an AzureContainerAppCustomDomain resource, Planton provisions:

- **Custom-domain binding** -- an `azurerm_container_app_custom_domain` entry on the referenced app's ingress, certificate-less (managed flow) or bound to an environment certificate over SNI (bring-your-own flow)

## Prerequisites

- **Azure credentials** configured via environment variables or Planton provider config
- **An AzureContainerApp with ingress enabled** (the binding attaches to the app's ingress)
- **Public DNS records that already resolve** (Azure validates during creation):
  - TXT at `asuid.{domain}` carrying the app's `custom_domain_verification_id` output
  - a CNAME from the domain to the app's `ingress_fqdn` output (or an A record to the environment's static IP for apex domains)
- For bring-your-own TLS: an **AzureContainerAppEnvironmentCertificate** on the app's environment

## Quick Start

Create a file `custom-domain.yaml`:

```yaml
apiVersion: azure.planton.dev/v1alpha1
kind: AzureContainerAppCustomDomain
metadata:
  name: app-custom-domain
  annotations:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: dev.AzureContainerAppCustomDomain.app-custom-domain
spec:
  domainName: app.example.com
  containerAppId:
    valueFrom:
      kind: AzureContainerApp
      name: my-app
      fieldPath: status.outputs.container_app_id
```

Deploy:

```shell
planton apply -f custom-domain.yaml
```

With the domain's zone on Azure DNS, declare the verification and routing records as [Azure DNS Record](/docs/catalog/azure/dns-record) resources in the same composition. For the managed-TLS flow, follow the binding with an [Azure Container App Environment Managed Certificate](/docs/catalog/azure/container-app-environment-managed-certificate) for the same hostname -- Azure attaches it automatically once issued. For bring-your-own TLS, add the certificate reference and `certificateBindingType: SNI_ENABLED`.

## Key Outputs

| Output | Purpose |
|--------|---------|
| `custom_domain_id` | The binding's identifier |
| `managed_certificate_id` | The attached managed certificate, once Azure fills it in (empty for BYO bindings) |

## Related Resources

- [Azure Container App](/docs/catalog/azure/container-app) -- the app the domain binds to
- [Azure DNS Record](/docs/catalog/azure/dns-record) -- the verification and routing records
- [Azure Container App Environment Managed Certificate](/docs/catalog/azure/container-app-environment-managed-certificate) / [Environment Certificate](/docs/catalog/azure/container-app-environment-certificate) -- the two TLS flows
