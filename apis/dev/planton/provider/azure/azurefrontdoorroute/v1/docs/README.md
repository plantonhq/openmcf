# AzureFrontDoorRoute -- Design Research

## Scope

The route matches client requests on an endpoint by URL pattern and
forwards them to an origin group, with protocol policy and edge
caching. It is a first-class kind (not a list on the endpoint) because
routes are many-per-endpoint with independent lifecycles and reference
resources across the family.

Source of truth: `azurerm_cdn_frontdoor_route`
(terraform-provider-azurerm v4.80,
`internal/services/cdn/cdn_frontdoor_route_resource.go`),
parity-verified against pulumi-azure v6 (`cdn.FrontdoorRoute`).

## Field mapping

| Spec field | azurerm attribute | Notes |
| --- | --- | --- |
| `endpoint_id` | `cdn_frontdoor_endpoint_id` | FK to AzureFrontDoorEndpoint; ForceNew (the route's ARM parent) |
| `route_name` | `name` | ForceNew; 2-90 chars (the provider's own regex) |
| `origin_group_id` | `cdn_frontdoor_origin_group_id` | FK to AzureFrontDoorOriginGroup; updatable in place |
| `origin_ids` | `cdn_frontdoor_origin_ids` | Repeated FKs to AzureFrontDoorOrigin; NEVER sent to the API -- pure provisioning-order references (ARM rejects a route whose group has no origins) |
| `patterns_to_match` | `patterns_to_match` | min 1; each starts with "/" |
| `supported_protocols` | `supported_protocols` | closed enum, 1-2 unique values |
| `forwarding_protocol` | `forwarding_protocol` | closed enum; unspecified deploys MatchRequest |
| `https_redirect_enabled` | `https_redirect_enabled` | default true |
| `link_to_default_domain` | `link_to_default_domain` | default true; false requires associated custom domains (Azure apply-time error until the custom-domain kind exists) |
| `enabled` | `enabled` | default true |
| `origin_path` | `cdn_frontdoor_origin_path` | prepended on the origin side |
| `cache.*` | `cache` block | absence disables caching (the provider transmits an explicit null) |

## Validation contracts

- **HTTPS redirect requires both protocols** -- the provider validates
  this at apply (`validate.SupportsBothHttpAndHttps`); the spec
  front-loads it as a message CEL that treats the ABSENT
  `https_redirect_enabled` as true (its documented default).
- **Cache content types come from Azure's fixed allowlist** -- the
  provider's `frontDoorContentTypes()` list, encoded as an items CEL so
  a typo fails at validation, not at apply.
- **Query-string names must not contain commas** -- Azure transports
  the list as a CSV string; provider-validated, spec-front-loaded.
- **`link_to_default_domain: false` without custom domains** is an
  apply-time Azure error, deliberately NOT a spec CEL: the field's
  contract changes when custom-domain references land with that kind,
  and inventing a stricter interim rule would be a false constraint.

## Deferred references (recorded)

- `rule_set_ids` and `custom_domain_ids` -- azurerm carries both; they
  reference kinds that do not exist in the catalog yet (Front Door rule
  sets, custom domains). They land WITH those kinds, wired as typed
  references, in the sessions that forge them. Until then a route
  serves the endpoint's default domain, which is also azurerm's default
  posture.

## Outputs

| Output | Source | Consumers |
| --- | --- | --- |
| `route_id` | resource id | management/diagnostics |
| `route_name` | `name` | human/portal orientation |

No hostname output on purpose: the client-facing hostname lives on the
ENDPOINT's outputs; the route is policy attached to that hostname.

## Recorded skips (with reasons)

- `cdn_frontdoor_custom_domain_ids` / `cdn_frontdoor_rule_set_ids` --
  see "Deferred references" above.
- The `cdn_frontdoor_custom_domain_association` resource is a
  Terraform-side ordering construct with no server-side object; when
  custom domains land, the association is realized inside module logic,
  never as a catalog kind.
