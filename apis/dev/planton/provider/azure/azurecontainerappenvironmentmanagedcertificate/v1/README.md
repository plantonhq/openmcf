# AzureContainerAppEnvironmentManagedCertificate

Provision a TLS certificate Azure issues and renews end to end -- free, domain-validated, for one custom domain on a Container App Environment.

## Overview

The managed certificate lives on the environment and attaches to the matching `AzureContainerAppCustomDomain` binding asynchronously once issued. **Deployment blocks on domain validation**: Azure only issues the certificate after proving you control `subject_name` against public DNS, so the required records must exist BEFORE this resource deploys:

1. A TXT record at `asuid.{subject_name}` carrying the app's `custom_domain_verification_id` output, and
2. the routing record for the hostname itself -- a CNAME to the app's `ingress_fqdn` (CNAME validation) or an HTTP-reachable path (HTTP validation).

With the domain's zone on Azure DNS, both records are `AzureDnsRecord` resources composed in the same deployment. The typical sequence: bind the domain (certificate-less) with `AzureContainerAppCustomDomain`, then deploy this certificate for the same hostname.

## Key Features

- **Zero-cost, zero-touch TLS**: Azure issues and rotates the certificate before expiry
- **Two validation methods**: CNAME (subdomains pointing at the app) or HTTP (apex/A-record domains)
- **Composable**: the environment by reference; the validation records compose through `AzureDnsRecord`

## When to Use

- Standard TLS for app custom domains -- the default choice unless you need EV/OV, wildcards, or an org-mandated CA (bring your own with `AzureContainerAppEnvironmentCertificate` for those)

## Spec Highlights

| Field | Notes |
| --- | --- |
| `certificate_name` | The certificate's name on the environment. ForceNew |
| `container_app_environment_id` | The environment, by ARM id reference. ForceNew |
| `subject_name` | The one domain covered -- no wildcards, no extra SANs. ForceNew |
| `domain_control_validation` | HTTP (default) or CNAME. ForceNew |
| `tags` | The only field Azure updates in place |

## Outputs

| Output | Purpose |
| --- | --- |
| `certificate_id` | The managed certificate's ARM ID |
| `validation_token` | Azure's domain-validation token (informational once issued) |

## Quick Example

```yaml
apiVersion: azure.planton.dev/v1
kind: AzureContainerAppEnvironmentManagedCertificate
metadata:
  name: app-managed-cert
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
