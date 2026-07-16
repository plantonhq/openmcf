# AzureFrontDoorEndpoint -- Design Research

## Scope

The endpoint is the public entry point of a Front Door deployment: a
generated, globally unique hostname that routes attach to. It is a
first-class kind (not a list on the profile) because endpoints are
many-per-profile with independent lifecycles -- one profile commonly
fronts several applications.

Source of truth: `azurerm_cdn_frontdoor_endpoint`
(terraform-provider-azurerm v4.80,
`internal/services/cdn/cdn_frontdoor_endpoint_resource.go`),
parity-verified against pulumi-azure v6 (`cdn.FrontdoorEndpoint`).

## Field mapping

| Spec field | azurerm attribute | Notes |
| --- | --- | --- |
| `profile_id` | `cdn_frontdoor_profile_id` | FK to AzureFrontDoorProfile; ForceNew |
| `endpoint_name` | `name` | ForceNew; 2-46 chars (the provider's own regex); the hostname prefix |
| `enabled` | `enabled` | default true; false stops traffic at the edge without deletion |
| `tags` | `tags` | Merged over Planton-derived identity tags; user tags win |

## Outputs

| Output | Source | Consumers |
| --- | --- | --- |
| `endpoint_id` | resource id | AzureFrontDoorRoute parent FK; (later) security-policy domain associations |
| `endpoint_name` | `name` | human/portal orientation |
| `host_name` | `host_name` (computed) | AzureDnsRecord CNAME targets; the client-facing address |

## Behavior worth knowing

- **The hostname carries a generated hash** (`{name}-{hash}.z01.azurefd.net`),
  so endpoint names need only per-profile uniqueness -- but the hash
  changes on recreation, which is why `endpoint_name` renames (ForceNew)
  break DNS.
- **Disabled endpoints still expose their hostname** -- DNS can be
  prepared before launch; Front Door answers with errors until enabled.

## Recorded skips (with reasons)

- None -- the azurerm surface is four attributes and all are modeled.
