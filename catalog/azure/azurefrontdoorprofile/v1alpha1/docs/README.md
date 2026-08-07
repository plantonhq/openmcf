# AzureFrontDoorProfile -- Design Research

## Scope

The profile is deliberately just the container for a Front Door
(Standard/Premium) deployment: SKU tier, origin response timeout,
managed identity, access-log scrubbing, and tags. The delivery surface
-- endpoints, origin groups, origins, routes -- decomposes into
first-class kinds referencing the profile, mirroring Azure's own ARM
child-resource model (each child has its own resource type, lifecycle,
and API under `Microsoft.Cdn/profiles/*`).

Source of truth: `azurerm_cdn_frontdoor_profile` (terraform-provider-azurerm
v4.80, `internal/services/cdn/cdn_frontdoor_profile_resource.go`),
parity-verified against pulumi-azure v6 (`cdn.FrontdoorProfile`).

## Field mapping

| Spec field | azurerm attribute | Notes |
| --- | --- | --- |
| `resource_group` | `resource_group_name` | FK to AzureResourceGroup |
| `profile_name` | `name` | ForceNew; 2-90 chars, letters/digits/hyphens (the provider's own regex) |
| `sku` | `sku_name` | Closed enum; unspecified deploys Standard_AzureFrontDoor; ForceNew, and the provider's CustomizeDiff rejects Premium -> Standard outright |
| `response_timeout_seconds` | `response_timeout_seconds` | 16-240, default 120 |
| `identity` | `identity` | SystemAssigned / UserAssigned / both; UAI ids are FKs to AzureUserAssignedIdentity |
| `log_scrubbing_variables` | `log_scrubbing_rule` (set of blocks) | Modeled as a repeated enum because the block's only configurable field is `match_variable` -- the provider forces `operator = EqualsAny` and an empty selector (the only shape the service accepts on profiles). Presence enables scrubbing; an empty list disables it |
| `tags` | `tags` | Merged over Planton-derived identity tags; user tags win |

## Outputs

| Output | Source | Consumers |
| --- | --- | --- |
| `profile_id` | resource id | AzureFrontDoorEndpoint / AzureFrontDoorOriginGroup parent FK |
| `profile_name` | `name` | human/portal orientation |
| `resource_guid` | `resource_guid` | apex-domain afdverify traffic validation |
| `identity_principal_id` | `identity[0].principal_id` | Key Vault grants for bring-your-own certificates (empty without a system-assigned identity) |

## Behavior worth knowing

- **Front Door is global**: the provider forces `location = "global"`;
  there is no region field anywhere in the family.
- **SKU is one-way**: ForceNew AND the Premium -> Standard downgrade is
  rejected even as a replace. Standard -> Premium replaces the profile
  and everything nested under it.
- **Deletion is slow**: profile deletion runs several minutes (the
  provider allows up to 6 h); endpoints/routes under it delete first in
  composed teardowns.

## Recorded skips (with reasons)

- **Endpoints/origin groups/routes as inline lists** -- deliberately
  NOT modeled: each is a first-class kind (independent lifecycle,
  many-per-parent, FK-referenced), which is what keeps per-region
  stamps and per-app endpoints composable.
- **Custom domains, secrets (BYO certificates), rule sets, Front Door
  WAF policy + security policy** -- separate ARM resources with their
  own lifecycles; they are their own kinds (rule sets/custom
  domains/secrets and the WAF pair land as the family completes; the
  route's `rule_set_ids`/`custom_domain_ids` references arrive with
  those kinds).
