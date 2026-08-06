# AzureContainerAppCustomDomain

Bind a custom domain to a Container App: your own hostname serving the app instead of the generated `*.azurecontainerapps.io` address.

## Overview

DNS must prove ownership before the binding deploys -- Azure validates during creation, and the deployment fails without these records resolving publicly:

1. A TXT record at `asuid.{domain_name}` carrying the app's `custom_domain_verification_id` output, and
2. a CNAME from the hostname to the app's `ingress_fqdn` output (or an A record to the environment's static IP for apex domains).

TLS comes one of two ways: leave the certificate fields unset for the **Azure-managed flow** (bind first, then an `AzureContainerAppEnvironmentManagedCertificate` for the same hostname issues and attaches asynchronously), or set both `container_app_environment_certificate_id` and `certificate_binding_type` for a **bring-your-own certificate**.

## Key Features

- **Both TLS flows**: Azure-managed (free, auto-renewed) and bring-your-own (EV/OV, wildcards, org CAs)
- **Drift-safe managed flow**: both engines ignore Azure's asynchronous certificate attachment so it never reads as spurious drift -- while BYO bindings keep full drift detection
- **Composable**: the app by reference; the certificate by the environment certificate's `certificate_id` output; the DNS records through `AzureDnsRecord`

## When to Use

- Serving any Container App from your own hostname -- the last mile of the web-domain story

## Spec Highlights

| Field | Notes |
| --- | --- |
| `domain_name` | The hostname; one concrete name per binding (the environment's custom DNS suffix is the wildcard mechanism). ForceNew |
| `container_app_id` | The app (ingress must be enabled). ForceNew |
| `container_app_environment_certificate_id` + `certificate_binding_type` | Set together for BYO TLS (SNI_ENABLED); leave both unset for the managed flow. ForceNew |

## Outputs

| Output | Purpose |
| --- | --- |
| `custom_domain_id` | The binding's synthetic identifier |
| `managed_certificate_id` | The attached managed certificate, once Azure fills it in (empty for BYO) |

## Quick Example

```yaml
apiVersion: azure.planton.dev/v1alpha1
kind: AzureContainerAppCustomDomain
metadata:
  name: app-custom-domain
spec:
  domainName: app.example.com
  containerAppId:
    valueFrom:
      kind: AzureContainerApp
      name: my-app
      fieldPath: status.outputs.container_app_id
```

The managed-certificate flow: this binding deploys certificate-less, then a managed certificate for `app.example.com` attaches automatically once issued.

---

© Planton. Licensed under [Apache-2.0](https://github.com/plantonhq/planton/blob/main/LICENSE).
