---
title: "Premium with Private Origins"
description: "This preset creates a Premium-tier profile with a system-assigned managed identity -- the shape for locked-down architectures where backends disable public access entirely and Front Door reaches them..."
type: "preset"
rank: "02"
presetSlug: "02-premium-private-origins"
componentSlug: "front-door-profile"
componentTitle: "Front Door Profile"
provider: "azure"
icon: "package"
order: 2
---

# Premium with Private Origins

This preset creates a Premium-tier profile with a system-assigned
managed identity -- the shape for locked-down architectures where
backends disable public access entirely and Front Door reaches them
over Private Link.

## When to Use

- Origins (App Service, Storage, internal load balancers) that must not
  be reachable from the public internet
- Deployments that will attach the managed WAF rule sets
  (Microsoft_DefaultRuleSet, Bot Manager)
- Custom domains that will carry bring-your-own Key Vault certificates

## Key Configuration Choices

- **`sku: PREMIUM`** -- required for Private Link origins and managed
  WAF rules; the upgrade is one-way (Azure refuses a downgrade), so
  this is a deliberate commitment
- **System-assigned identity** -- provisioned now so certificate wiring
  later is a grant, not a redeploy; the `identity_principal_id` output
  is the principal to grant Key Vault access to
- **Private Link itself lives on the ORIGIN** -- configure it per
  backend on AzureFrontDoorOrigin's `privateLink` block

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<resource-group-name>` | The AzureResourceGroup's Planton resource name | Your Azure composition |
| `profileName` (example value) | 2-90 chars, letters/digits/hyphens -- rename to your convention | Your naming convention |

## Downstream Wiring

Grant the identity read access to your certificate vault:

```yaml
# On an AzureRoleAssignment
principalId:
  valueFrom:
    kind: AzureFrontDoorProfile
    name: my-secure-front-door
    fieldPath: status.outputs.identity_principal_id
```
