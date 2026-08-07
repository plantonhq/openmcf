---
title: "Front Door Custom Domain"
description: "Front Door Custom Domain deployment documentation"
icon: "package"
order: 100
componentName: "azurefrontdoorcustomdomain"
---

# Azure Front Door Custom Domain

Deploys a custom domain inside an Azure Front Door (Standard/Premium) profile -- your own hostname (www.example.com, or the wildcard *.example.com) served at Microsoft's edge instead of the generated `*.azurefd.net` endpoint hostname. TLS is always on: Azure manages the certificate end to end by default, or the domain serves a bring-your-own certificate wrapped by an AzureFrontDoorSecret. The component integrates with Planton's Provider Connections for Azure credential management and ValueFromRef for dependency wiring to the profile, the DNS zone, and the secret.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **Front Door Custom Domain** -- a named child of the profile in a pending-validation state, exporting a `validation_token`
- **TLS configuration** -- an Azure-managed DV certificate (free, auto-rotated) or a customer certificate served from a Front Door secret, optionally with a pinned cipher-suite policy
- **DNS zone binding** -- when `dnsZoneId` references an AzureDnsZone, Front Door watches that zone for the validation records

## The Two-Step Lifecycle

A custom domain goes live in two moves:

1. **Create** -- the domain deploys immediately as pending and exports a `validation_token`.
2. **Validate & point** -- publish the token as a TXT record at `_dnsauth.<host_name>` (Azure flips the domain to approved), then CNAME the hostname to the endpoint's `host_name` output so traffic arrives.

Both records are AzureDnsRecord resources when the zone lives in Azure DNS; with external DNS, publish the same two records at your provider.

## The Domain in the Front Door Family

- **AzureFrontDoorProfile** -- the parent container, referenced by `profileId`
- **AzureFrontDoorSecret** -- the BYO certificate, referenced by `tls.secretId` with CUSTOMER_CERTIFICATE
- **AzureDnsZone / AzureDnsRecord** -- where the validation TXT and traffic CNAME live
- **AzureFrontDoorRoute** -- serves this hostname by listing the domain in its `customDomainIds` (the route side owns the attachment)
- **AzureFrontDoorSecurityPolicy** -- scopes a WAF to this domain through the same ID

## Before You Deploy

### Planton Setup

- **Azure Provider Connection** -- an active connection in the Connect module with credentials for the target Azure subscription.
- **Planton Runner** -- required when using Runner-based credential delivery.

### Azure Subscription

- **An Azure Front Door Profile** the domain nests under.
- **Control of the hostname's DNS** -- validation publishes a TXT record under it.
- **For wildcards or hostnames over 64 characters**: an AzureFrontDoorSecret wrapping a customer certificate (Azure's managed certificates cannot cover them).

## Deploy

### Console

Open the deployment store, find **Azure Front Door Custom Domain**, and click **Deploy**. The wizard walks you through the hostname (wildcards flagged live when they force a customer certificate), the dot-free ARM resource name, the TLS decision, and the DNS zone binding with a live preview of the exact validation records. Start from the **Managed Certificate** preset in the [Presets](#presets) tab.

### CLI

```yaml
apiVersion: azure.planton.dev/v1
kind: AzureFrontDoorCustomDomain
metadata:
  name: www-example-com
  org: acme-corp
  env: prod
spec:
  profileId:
    valueFrom:
      kind: AzureFrontDoorProfile
      name: cdn-profile
      fieldPath: status.outputs.profile_id
  domainName: www-example-com
  hostName: www.example.com
  dnsZoneId:
    valueFrom:
      kind: AzureDnsZone
      name: example-com-zone
      fieldPath: status.outputs.zone_id
  tls: {}
```

```shell
planton apply -f front-door-custom-domain.yaml
```

The empty `tls` block deploys Azure's managed certificate. After deploy, publish the `validation_token` output as a TXT record at `_dnsauth.www.example.com`, then CNAME `www.example.com` to the endpoint's hostname.

### InfraChart

```yaml
spec:
  profileId:
    valueFrom:
      kind: AzureFrontDoorProfile
      name: cdn-profile
      fieldPath: status.outputs.profile_id
  domainName: www-example-com
  hostName: www.example.com
  tls:
    certificateType: CUSTOMER_CERTIFICATE
    secretId:
      valueFrom:
        kind: AzureFrontDoorSecret
        name: wildcard-cert-secret
        fieldPath: status.outputs.secret_id
```

## Key Configuration

**Hostname** -- the FQDN the domain serves; only the FIRST label may be the wildcard `*`. ForceNew: the hostname IS the domain's identity. With an Azure-managed certificate the hostname caps at 64 characters and cannot be a wildcard.

**Domain resource name** -- the ARM label, NOT the hostname (dots forbidden); convention is the hostname with dots as hyphens. ForceNew: renaming replaces the domain and restarts validation.

**Certificate type** -- unspecified deploys MANAGED_CERTIFICATE (Azure's default). CUSTOMER_CERTIFICATE requires `tls.secretId` and unlocks wildcards, EV/OV, and org-mandated CAs.

**Cipher suite** -- unset serves Azure's default suite set. Pin TLS12_2022/TLS12_2023 or hand-pick suites (CUSTOMIZED: at least one TLS 1.2 suite; TLS 1.3 suites both-or-none). The minimum TLS version is not configurable -- every domain floors at TLS 1.2.

## Outputs and Dependencies

### What This Component Consumes

| Dependency | Field | ValueFromRef Path |
|------------|-------|-------------------|
| **AzureFrontDoorProfile** | `profileId` | `status.outputs.profile_id` |
| **AzureDnsZone** (optional) | `dnsZoneId` | `status.outputs.zone_id` |
| **AzureFrontDoorSecret** (customer certificate) | `tls.secretId` | `status.outputs.secret_id` |

### What This Component Provides

| Output | Description | Common Downstream Use |
|--------|-------------|----------------------|
| `custom_domain_id` | ARM resource ID of the domain | AzureFrontDoorRoute.`customDomainIds`; AzureFrontDoorSecurityPolicy.`domainIds` |
| `host_name` | The hostname the domain serves | The CNAME source at your DNS provider |
| `validation_token` | The DNS challenge (public by design) | The `_dnsauth.<host>` TXT record's value |
| `expiration_date` | When the current token expires (RFC-3339) | Re-validation scheduling |

## Presets

| Preset | Rank | Description |
|--------|------|-------------|
| Managed Certificate | 1 | Single-host domain with Azure's free auto-rotated certificate |
| Wildcard BYO Certificate | 2 | Wildcard hostname served from a Front Door secret |
| Hardened Ciphers | 3 | Customer certificate with a pinned cipher-suite policy |

## Related Components

- **AzureFrontDoorProfile** -- the parent container
- **AzureFrontDoorSecret** -- wraps the BYO certificate
- **AzureDnsZone / AzureDnsRecord** -- host the validation TXT and traffic CNAME
- **AzureFrontDoorRoute** -- attaches this domain to serving paths
- **AzureFrontDoorSecurityPolicy** -- scopes a WAF to this domain
