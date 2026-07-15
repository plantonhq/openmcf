---
title: "Front Door Custom Domain"
description: "Front Door Custom Domain deployment documentation"
icon: "package"
order: 100
componentName: "azurefrontdoorcustomdomain"
---

# Azure Front Door Custom Domain

Creates a custom domain inside an AzureFrontDoorProfile -- your own hostname served by Front Door with managed or bring-your-own TLS. The domain deploys pending DNS validation and exports the TXT challenge (`validation_token`); AzureFrontDoorRoute resources serve it through their `customDomainIds`.

## What Gets Created

When you deploy an AzureFrontDoorCustomDomain resource, Planton provisions:

- **Front Door Custom Domain** -- an `azurerm_cdn_frontdoor_custom_domain` on the referenced profile, in the pending-validation state until the DNS TXT challenge is published

## Prerequisites

- **Azure credentials** configured via environment variables or Planton provider config
- **An AzureFrontDoorProfile** to create the domain in (referenced through `profileId`)
- **Control of the hostname's DNS** -- validation publishes a TXT record; serving traffic needs a CNAME to the endpoint
- **For BYO certificates**: an AzureFrontDoorSecret, and a one-time grant of Key Vault read access to Front Door's service principal (`Microsoft.AzureFrontDoor-Cdn`)

## Quick Start

Create a file `custom-domain.yaml`:

```yaml
apiVersion: azure.planton.dev/v1
kind: AzureFrontDoorCustomDomain
metadata:
  name: www-example-com
  annotations:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: dev.AzureFrontDoorCustomDomain.www-example-com
spec:
  profileId:
    valueFrom:
      kind: AzureFrontDoorProfile
      name: my-front-door
      fieldPath: status.outputs.profile_id
  domainName: www-example-com
  hostName: www.example.com
  tls: {}
```

Deploy:

```shell
planton apply -f custom-domain.yaml
```

Then read the `validation_token` output, publish it as a TXT record at `_dnsauth.www.example.com`, wait for Azure to approve the domain, CNAME the hostname to the endpoint's `host_name`, and attach the domain to a route. An empty `tls` block means Azure's managed certificate; wildcards and EV/OV certificates use `certificateType: CUSTOMER_CERTIFICATE` with a `secretId` reference.

## Key Outputs

| Output | Purpose |
|--------|---------|
| `custom_domain_id` | The ARM id -- what routes reference in `customDomainIds` and security policies scope WAFs to |
| `host_name` | The hostname to CNAME to the endpoint once validated |
| `validation_token` | The DNS TXT challenge for `_dnsauth.<host_name>` |
| `expiration_date` | When the current token expires |

## Related Resources

- [Azure Front Door Profile](/docs/catalog/azure/front-door-profile) -- the parent profile
- [Azure Front Door Secret](/docs/catalog/azure/front-door-secret) -- the BYO certificate node
- [Azure Front Door Route](/docs/catalog/azure/front-door-route) -- serves this domain
- [Azure DNS Zone](/docs/catalog/azure/dns-zone) -- hosts the validation and CNAME records
