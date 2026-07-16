---
title: "Domain Verification TXT"
description: "The ownership-proof TXT record that custom-domain flows require before they bind a hostname: Container Apps checks `asuid.{host}`, Front Door checks `_dnsauth.{host}`, and most SaaS domain..."
type: "preset"
rank: "04"
presetSlug: "04-domain-verification-txt"
componentSlug: "dns-record"
componentTitle: "DNS Record"
provider: "azure"
icon: "package"
order: 4
---

# Domain Verification TXT

The ownership-proof TXT record that custom-domain flows require before they bind a hostname: Container Apps checks `asuid.{host}`, Front Door checks `_dnsauth.{host}`, and most SaaS domain verifications follow the same pattern.

## When to Use

- Binding a custom domain to a Container App: publish the app's `custom_domain_verification_id` output at `asuid.{host}` BEFORE deploying the AzureContainerAppCustomDomain binding
- Validating a Front Door custom domain: publish the domain's `validation_token` output at `_dnsauth.{host}`
- Any provider-issued domain-ownership token (Google, Microsoft 365, certificate authorities)

## Key Configuration Choices

- `name` -- the service defines it; underscore-led and dotted names are fully supported
- `ttlSeconds: 60` -- validation services poll public DNS; a short TTL makes the proof visible fast
- The token can be referenced instead of pasted -- e.g. `valueFrom` an `AzureContainerApp`'s `custom_domain_verification_id` output

## Placeholders to Replace

| Placeholder | Description | Where to Find |
|---|---|---|
| `<your-resource-group>` | The zone's resource group | `AzureResourceGroup.status.outputs.resource_group_name` |
| `example.com` | Replace with the zone name | `AzureDnsZone.status.outputs.zone_name` |
| `asuid.app` | Replace `app` with the hostname being verified, relative to the zone | The custom-domain binding you are creating |
| `<verification-token-from-app-output>` | The service-issued token | `AzureContainerApp.status.outputs.custom_domain_verification_id`, `AzureFrontDoorCustomDomain.status.outputs.validation_token`, or the service's setup screen |

## Related Presets

- `03-mail-mx-records` -- pair with SPF/DKIM/DMARC TXT records for mail
