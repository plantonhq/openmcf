# AzureFrontDoorOrigin -- Design Research

## Scope

The origin is one backend in a Front Door origin group. It decomposes
from the group (rather than folding in as a list) because origins have
genuinely independent lifecycles -- per-region stamps, blue/green
swaps, per-origin Private Link approval workflows -- and because there
is no member-side kind that could express pool membership instead.

Source of truth: `azurerm_cdn_frontdoor_origin`
(terraform-provider-azurerm v4.80,
`internal/services/cdn/cdn_frontdoor_origin_resource.go`),
parity-verified against pulumi-azure v6 (`cdn.FrontdoorOrigin`).

## Field mapping

| Spec field | azurerm attribute | Notes |
| --- | --- | --- |
| `origin_group_id` | `cdn_frontdoor_origin_group_id` | FK to AzureFrontDoorOriginGroup; ForceNew |
| `origin_name` | `name` | ForceNew; 2-90 chars (the provider's own regex) |
| `host_name` | `host_name` | hostname, IPv4, or IPv6 |
| `certificate_name_check_enabled` | `certificate_name_check_enabled` | The provider requires it explicitly; modeled optional-with-default-true and always sent |
| `origin_host_header` | `origin_host_header` | unset sends the origin's hostname (Azure behavior) |
| `http_port` / `https_port` | same | 1-65535, defaults 80/443 |
| `priority` / `weight` | same | 1-5 / 1-1000, defaults 1/500 |
| `enabled` | `enabled` | default true; false drains the origin |
| `private_link.*` | `private_link` block | location, target id, target type, request message |

## Validation contracts (the provider enforces these at APPLY; the spec
front-loads what is statically checkable)

- **Private Link requires certificate-name checking** -- spec message
  CEL (statically checkable).
- **`target_type` XOR Private Link Service target** -- when the target
  is not a PLS, `target_type` is required; when it IS a PLS, it must be
  absent. Spec CEL keyed on the `/privateLinkServices/` id segment.
- **Private Link is PREMIUM-profile only** -- NOT statically checkable
  (the SKU lives on the referenced profile); Azure rejects the apply
  with a clear error. Documented on the field.

## Engine dialect note

The private-link `target_type` secondary values differ per engine:
azurerm spells `blob_secondary`/`web_secondary` where the pulumi bridge
expects `blobSecondary`/`webSecondary`. Both land on the same ARM group
id; each module maps the spec enum to its provider's vocabulary
(`GATEWAY` -> `"Gateway"` -- the capital G is Azure's own).

## Outputs

| Output | Source | Consumers |
| --- | --- | --- |
| `origin_id` | resource id | AzureFrontDoorRoute `origin_ids` (deploy ordering) |
| `origin_name` | `name` | human/portal orientation |

## Recorded skips (with reasons)

- None -- the azurerm surface is fully modeled. The provider's
  `azuresdkhacks` update client (always re-sending `origin_host_header`
  so PATCH does not clear it) is provider-internal behavior both
  engines inherit.
