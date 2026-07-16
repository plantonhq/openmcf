---
title: "Mail MX Records"
description: "The domain's mail-exchange record set: a primary and a backup mail server, each with its own delivery preference in one record set (Azure stores all MX values for the apex together)."
type: "preset"
rank: "03"
presetSlug: "03-mail-mx-records"
componentSlug: "dns-record"
componentTitle: "DNS Record"
provider: "azure"
icon: "package"
order: 3
---

# Mail MX Records

The domain's mail-exchange record set: a primary and a backup mail server, each with its own delivery preference in one record set (Azure stores all MX values for the apex together).

## When to Use

- Routing the domain's email to your mail provider (use the provider's published exchanges and preferences -- e.g. Google Workspace and Microsoft 365 each document theirs)
- Adding a backup exchange at a higher preference number

## Key Configuration Choices

- `preference` -- lower is tried first; equal preferences load-balance
- One record set carries ALL the domain's mail servers -- never split them across resources (a second apex MX resource conflicts rather than merging)
- Pair with TXT records (SPF, DKIM, DMARC) declared as separate AzureDnsRecord resources on their own names

## Placeholders to Replace

| Placeholder | Description | Where to Find |
|---|---|---|
| `<your-resource-group>` | The zone's resource group | `AzureResourceGroup.status.outputs.resource_group_name` |
| `example.com` | Replace with the zone name | `AzureDnsZone.status.outputs.zone_name` |
| `mail1.example.com` / `mail2.example.com` | Replace with the mail exchanges | Your mail provider's setup documentation |

## Related Presets

- `04-domain-verification-txt` -- the TXT records mail authentication and domain-validation flows need
