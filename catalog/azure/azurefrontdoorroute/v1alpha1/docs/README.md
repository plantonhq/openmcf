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
| `rule_set_ids` | `cdn_frontdoor_rule_set_ids` | Repeated FKs to AzureFrontDoorRuleSet -- the attached delivery policies; same-profile membership is Azure's apply-time contract |
| `patterns_to_match` | `patterns_to_match` | min 1; each starts with "/" |
| `supported_protocols` | `supported_protocols` | closed enum, 1-2 unique values |
| `forwarding_protocol` | `forwarding_protocol` | closed enum; unspecified deploys MatchRequest |
| `https_redirect_enabled` | `https_redirect_enabled` | default true |
| `custom_domain_ids` | `cdn_frontdoor_custom_domain_ids` | Repeated FKs to AzureFrontDoorCustomDomain -- the route side owns the attachment; domains must be DNS-validated before traffic flows |
| `link_to_default_domain` | `link_to_default_domain` | default true; false requires at least one custom domain (spec CEL) |
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
- **`link_to_default_domain: false` requires at least one custom
  domain** -- azurerm's exact create/update contract (a route must
  answer on SOME hostname), front-loaded as a message CEL.

## Empty-vs-absent collections (both modules)

Front Door treats an EMPTY `rule_set_ids`/`custom_domain_ids`
collection differently from an ABSENT one (empty means "disassociate",
which only matters on update). Both modules normalize empty lists to
null/omitted so absence and emptiness agree.

## Outputs

| Output | Source | Consumers |
| --- | --- | --- |
| `route_id` | resource id | management/diagnostics |
| `route_name` | `name` | human/portal orientation |

No hostname output on purpose: the client-facing hostname lives on the
ENDPOINT's outputs; the route is policy attached to that hostname.

## Recorded skips (with reasons)

- The `cdn_frontdoor_custom_domain_association` resource is a
  Terraform-side ordering construct with no server-side object -- it
  exists to break dependency cycles when domains and routes live in
  separate Terraform configurations. Planton's decomposed kinds have no
  such cycle (a domain never references a route), so the route's
  `custom_domain_ids` field IS the attachment; no association resource
  exists in either module and no catalog kind models it.
