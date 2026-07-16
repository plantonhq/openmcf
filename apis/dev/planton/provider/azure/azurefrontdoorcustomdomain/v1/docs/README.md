# AzureFrontDoorCustomDomain -- Design Research

## Scope

The custom domain is Front Door's ownership-proven hostname node: it
carries the hostname, its TLS posture, and the DNS-validation
lifecycle. It is a first-class kind because domains are many-per-profile
with independent lifecycles, are referenced by routes
(`custom_domain_ids`) and by Front Door security policies (the WAF
association scope), and their validation workflow is an operational
surface of its own.

Source of truth: `azurerm_cdn_frontdoor_custom_domain`
(terraform-provider-azurerm v4.80,
`internal/services/cdn/cdn_frontdoor_custom_domain_resource.go`),
parity-verified against pulumi-azure v6 (`cdn.FrontdoorCustomDomain`).

## Field mapping

| Spec field | azurerm attribute | Notes |
| --- | --- | --- |
| `profile_id` | `cdn_frontdoor_profile_id` | FK to AzureFrontDoorProfile; ForceNew |
| `domain_name` | `name` | ForceNew; 2-260, letter/digit edges, hyphens (the provider's regex) -- the ARM name, NOT the hostname |
| `host_name` | `host_name` | ForceNew; FQDN validation mirrors the provider's (<=253 chars, >=2 labels, each <=63, wildcard only as the whole first label) |
| `dns_zone_id` | `dns_zone_id` | optional FK to AzureDnsZone (`zone_id` output); updatable |
| `tls.certificate_type` | `tls.certificate_type` | Managed (default) / Customer; updatable |
| `tls.secret_id` | `tls.cdn_frontdoor_secret_id` | FK to AzureFrontDoorSecret; required iff Customer, forbidden iff Managed (the provider's create-time contract, front-loaded) |
| `tls.cipher_suite.*` | `tls.cipher_suite` block | type (TLS12_2022 / TLS12_2023 / Customized) + custom tls12/tls13 suite lists |

## Validation contracts (front-loaded as CELs)

Mirrored from the provider's CustomizeDiff + create-time validators:

- managed certificates: host_name <= 64 characters and never a wildcard
- secret_id required with CUSTOMER_CERTIFICATE, forbidden with
  MANAGED_CERTIFICATE
- CUSTOMIZED requires custom_ciphers (and the predefined sets forbid
  them); at least one TLS 1.2 suite; tls13, when pinned, must list BOTH
  mandatory suites; suite values from the provider's allowlists (the
  four ECDHE-RSA TLS 1.2 suites -- the provider deliberately rejects
  the DHE suites the SDK enum carries)

## Recorded skips (with reasons)

- **`tls.minimum_version` / `tls.minimum_tls_version`** -- NOT modeled.
  Azure retired TLS 1.0/1.1 (March 2025); the provider accepts exactly
  one value (`TLS12`) on the current field and keeps the deprecated
  alias only for pre-retirement state. A field with a single legal
  value is a constant, not a configuration choice -- both modules send
  `TLS12` unconditionally. The field lands as a real enum if Azure ever
  raises the floor to TLS 1.3.

## Lifecycle notes

- **Creation never blocks on validation**: the domain lands in
  pending-validation and exports `validation_token` +
  `expiration_date`. The provider's own documented pattern has the DNS
  TXT record depend on the domain resource, which is only possible
  because create returns pre-validation.
- **Long provider timeouts are intentional** (create/delete 12 h,
  update 24 h): updates wait for the domain to leave transitional
  states, and BYO Private Link approvals can be manual. Both engines
  inherit them from the provider layer.
- ForceNew: `domain_name`, `profile_id`, `host_name`. The TLS block and
  DNS zone update in place.

## Outputs

| Output | Source | Consumers |
| --- | --- | --- |
| `custom_domain_id` | resource id | `AzureFrontDoorRoute.custom_domain_ids`; Front Door security policies (WAF scoping) |
| `host_name` | `host_name` | the CNAME source at the DNS provider |
| `validation_token` | validation properties | the `_dnsauth.<host_name>` TXT record value (public by design -- it proves DNS control, not secret possession) |
| `expiration_date` | validation properties | token-freshness monitoring |

## Composition

RG -> profile -> (secret ->) custom domain -> route
(`custom_domain_ids`); with Azure DNS, the zone reference lets Azure
watch for the validation record, and the TXT/CNAME records themselves
are AzureDnsRecord resources in the zone.
