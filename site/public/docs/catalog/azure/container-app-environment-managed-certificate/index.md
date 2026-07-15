---
title: "Container App Environment Managed Certificate"
description: "Container App Environment Managed Certificate deployment documentation"
icon: "package"
order: 100
componentName: "azurecontainerappenvironmentmanagedcertificate"
---

# Azure Container App Environment Managed Certificate

Provisions a free, Azure-issued and Azure-renewed TLS certificate for one custom domain on a Container App Environment.

## What Gets Created

When you deploy an AzureContainerAppEnvironmentManagedCertificate resource, Planton provisions:

- **Managed certificate** -- an `azurerm_container_app_environment_managed_certificate` on the referenced environment; Azure validates domain ownership, issues the certificate, and renews it before expiry

## Prerequisites

- **Azure credentials** configured via environment variables or Planton provider config
- **An AzureContainerAppEnvironment** to provision the certificate on
- **Public DNS records that already resolve** (deployment blocks on validation):
  - TXT at `asuid.{subject_name}` carrying the app's `custom_domain_verification_id` output
  - a CNAME from the hostname to the app's `ingress_fqdn` (for CNAME validation) or HTTP reachability (for HTTP validation)
- Typically an existing **AzureContainerAppCustomDomain** binding for the same hostname (deployed certificate-less first; Azure attaches this certificate to it once issued)

## Quick Start

Create a file `managed-cert.yaml`:

```yaml
apiVersion: azure.planton.dev/v1
kind: AzureContainerAppEnvironmentManagedCertificate
metadata:
  name: app-managed-cert
  annotations:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: dev.AzureContainerAppEnvironmentManagedCertificate.app-managed-cert
spec:
  certificateName: app-example-com
  containerAppEnvironmentId:
    valueFrom:
      kind: AzureContainerAppEnvironment
      name: my-env
      fieldPath: status.outputs.environment_id
  subjectName: app.example.com
  domainControlValidation: CNAME
```

Deploy:

```shell
planton apply -f managed-cert.yaml
```

The deployment polls until Azure's domain validation succeeds -- publish the DNS records first or the operation fails around the 30-minute mark. With the domain's zone on Azure DNS, declare both records as [Azure DNS Record](/docs/catalog/azure/dns-record) resources in the same composition.

Managed certificates cover exactly one name: no wildcards, no additional SANs -- bring your own certificate for those.

## Key Outputs

| Output | Purpose |
|--------|---------|
| `certificate_id` | The managed certificate's ARM ID |
| `validation_token` | Azure's domain-validation token (informational once issued) |

## Related Resources

- [Azure Container App Custom Domain](/docs/catalog/azure/container-app-custom-domain) -- the binding Azure attaches this certificate to
- [Azure DNS Record](/docs/catalog/azure/dns-record) -- the validation and routing records
- [Azure Container App Environment Certificate](/docs/catalog/azure/container-app-environment-certificate) -- bring your own certificate (EV/OV, wildcards, org CAs)
