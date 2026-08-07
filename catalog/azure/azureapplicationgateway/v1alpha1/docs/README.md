# AzureApplicationGateway -- Design Research

## The Resource

An Azure Application Gateway (`Microsoft.Network/applicationGateways`) is
the Layer 7 load balancer and reverse proxy: host/path routing, TLS
termination (including mutual TLS), request/response rewriting, TCP/TLS
layer-4 proxying, and WAF enforcement. The component maps onto
`azurerm_application_gateway` (azurerm v4.x,
`internal/services/network/application_gateway_resource.go`),
parity-verified against pulumi-azure v6 (`network.ApplicationGateway`).

## Bundling (locked) and the Referenceable-Insides Seam

The gateway's sub-objects -- frontends, ports, listeners, pools, settings,
rules, path maps, probes, certificates, profiles, redirects, rewrites --
are ONE atomic ARM resource in Azure (azurerm models them fully inline),
wire to each other BY NAME, and have no life outside their gateway: the
fold test passes for all of them. What other resources need to reach is
exported as name-keyed map outputs (`backend_address_pool_ids`,
`frontend_ip_configuration_ids`), so pool membership composes member-side
(NIC/VMSS association resources) without splitting the gateway apart.

The WAF policy is the opposite verdict: its own kind
(`AzureWebApplicationFirewallPolicy`) referenced by ARM ID at three levels
(gateway, listener, path rule) -- independent lifecycle, FK-referenced,
shared across gateways.

## Field Mapping (azurerm → spec)

| azurerm | spec | Notes |
|---|---|---|
| `name`/`location`/`resource_group_name` | `name`/`region`/`resource_group` | ForceNew trio (everything else updates in place) |
| `gateway_ip_configuration` | `subnet_id` | The dedicated-subnet anchor, modeled as ONE subnet FK; the config name is derived plumbing (see skips for the max-2 nuance) |
| `sku.name`/`sku.tier` | `sku` enum | BASIC/STANDARD_V2/WAF_V2 -- name and tier always carry the same value on the modeled SKUs, so one enum covers both; v1 SKUs are RETIRED (azurerm blocks create/update into them at plan time) and not modeled |
| `sku.capacity` / `autoscale_configuration` | `capacity` XOR `autoscale` | The XOR is spec CEL; Basic's 1-2 cap and no-autoscale gates mirrored |
| `zones` | `zones` | "1"/"2"/"3", unique, ForceNew |
| `identity` | `identity` | UAI required for Key Vault certificate references |
| `frontend_ip_configuration` | `frontend_ip_configurations[]` | public (PIP FK) XOR private (subnet + allocation + optional static address) -- CEL-enforced |
| `frontend_port` | `frontend_ports[]` | Declared once, referenced by name |
| `backend_address_pool` | `backend_address_pools[]` | FQDNs/IPs; per-pool ARM ID exported in the map output |
| `backend_http_settings` | `backend_http_settings[]` | Full surface incl. draining, trusted roots, path prefix, dedicated connections; `cookie_based_affinity` (Enabled/Disabled) modeled as a boolean |
| `http_listener` | `http_listeners[]` | `host_names` only (the singular `host_name` is subsumed -- see skips); per-listener WAF FK, custom errors, SNI, profiles |
| `request_routing_rule` | `request_routing_rules[]` | `priority` REQUIRED (Azure mandates it on the v2 SKUs this catalog models); backend XOR redirect XOR path-map CEL |
| `url_path_map` + `path_rule` | `url_path_maps[]` | Full surface incl. per-path-rule WAF FK; default backend XOR redirect CEL |
| `probe` | `probes[]` | All four protocols; the HTTP-vs-TCP field pairings (path/host/match vs proxy-protocol) as one CEL |
| `ssl_certificate` | `ssl_certificates[]` | Key Vault secret FK (→ AzureKeyVaultCertificate versionless_secret_id) XOR inline PFX data+password (both sensitive) |
| `trusted_root_certificate` | `trusted_root_certificates[]` | Key Vault XOR inline |
| `trusted_client_certificate` | `trusted_client_certificates[]` | mTLS client-CA bundles (sensitive) |
| `ssl_profile` | `ssl_profiles[]` | mTLS verification + per-profile policy (shared SslPolicy message) |
| `ssl_policy` | `ssl_policy` | predefined XOR custom XOR disabled-protocols CEL |
| `redirect_configuration` | `redirect_configurations[]` | listener XOR URL target CEL |
| `rewrite_rule_set` | `rewrite_rule_sets[]` | Conditions, header edits, URL rewrite (+ reroute) |
| `listener`/`backend`/`routing_rule` | `listeners`/`backends`/`routing_rules` | The layer-4 (TCP/TLS) trio; protocol pairings CEL-enforced; azurerm's `backend` is named `backends`→`Layer4BackendSettings` for clarity against pools |
| `firewall_policy_id` / `force_firewall_policy_association` | same | FK → AzureWebApplicationFirewallPolicy; WAF_V2-gated by CEL |
| `custom_error_configuration` | `custom_error_configurations[]` | Top-level + per-listener |
| `private_link_configuration` | `private_link_configurations[]` | NAT ip-configurations with subnet FKs |
| `http2_enabled` / `fips_enabled` / `global` | same (+ `global_configuration`) | Buffering block requires both fields (azurerm's contract, mirrored as CEL) |
| `tags` | `tags` | User tags merged over Planton-derived tags |

## Recorded Skips (with reasons)

- **`authentication_certificate`** (top-level + per-backend-settings) --
  the v1-SKU backend-authentication mechanism; v1 SKUs are retired in
  azurerm 4.x (creation is blocked at plan time), so this is dead surface
  for every deployable gateway. `trusted_root_certificates` is the v2
  grain.
- **`http_listener.host_name` (singular)** -- subsumed by `host_names`
  (multi + wildcards); azurerm keeps both only for backwards
  compatibility and forbids using them together.
- **`backend_http_settings.sni_name` / `sni_validation_enabled` /
  `certificate_chain_validation_enabled`** -- brand-new azurerm 4.80
  fields ABSENT from pulumi-azure v6.38 (the latest release). A
  one-engine-only input would ship a silent-drop divergence, so they are
  skipped until the pulumi SDK carries them (same precedent as the Key
  Vault key's release_policy).
- **A second `gateway_ip_configuration`** -- azurerm allows Max2, but a
  gateway operationally lives in one dedicated subnet; the second config
  exists for an in-place subnet-migration maneuver, not a steady state a
  user declares. One subnet FK keeps the spec honest.
- **`private_endpoint_connection` (computed)** -- reverse-link
  convenience; private endpoints are first-class resources on the
  consuming side.
- **`ssl_certificate.public_cert_data` (computed)** -- a derived
  convenience attribute, not composition-relevant.

## Design Decisions

- **`priority` required on L7 routing rules** -- optional in azurerm's
  schema only because v1 SKUs predate it; Azure rejects v2 gateways
  without it, and this catalog models v2 (+Basic, which also requires
  it). Fail at validation, not after a 20-minute apply.
- **One shared protocol enum** (HTTP/HTTPS/TCP/TLS) with per-field CEL
  restrictions -- one vocabulary to learn; the restrictions document
  which layer each block belongs to.
- **`cookie_based_affinity` as a boolean** -- azurerm's required
  Enabled/Disabled string is a bool wearing a costume; the module maps it
  back. `cookie_based_affinity_enabled` defaults to Azure's real default
  (off).
- **Frontends made explicit** (repeated, public XOR private) -- replacing
  the previous auto-derived single public frontend; private-only and
  dual-frontend gateways are first-class, and the map output makes each
  frontend chainable.
- **Certificates keep both arms** (Key Vault XOR inline) -- Key Vault is
  the production grain (renewals propagate through the versionless
  secret ID); inline PFX remains for certificates not yet migrated, with
  data and password sensitive-annotated.
- **BASIC SKU feature gates as spec CEL** (no autoscale, capacity <= 2,
  no mTLS, no URL rewriting) -- azurerm enforces these at plan/create
  time; the spec surfaces them at validation time.

## Operational Behavior Worth Knowing

- **Applies run 15-25 minutes** -- create AND most updates; plan
  pipelines and E2E timeouts accordingly. Destroys run ~5-10 minutes.
- **The subnet must be dedicated** -- no other resource types; NSGs on it
  need Azure's documented gateway exceptions (GatewayManager inbound,
  65200-65535 for v2).
- **Key Vault ordering** -- the user-assigned identity needs GET on the
  vault's secrets BEFORE the gateway create starts, or the deploy fails
  ~15 minutes in.
- **A WAF policy must be same-region** and cannot be deleted while
  attached.
- **Sub-object names are the wiring** -- renaming a pool or listener
  rewires every reference to it inside the spec, and changes the map
  output keys downstream members reference. Rename with the same care as
  moving a resource.

## Composition

- `resource_group` → `AzureResourceGroup.status.outputs.resource_group_name`
- `subnet_id` (+ private frontends, private-link NAT configs) →
  `AzureSubnet.status.outputs.subnet_id`
- `frontend_ip_configurations[].public_ip_address_id` →
  `AzurePublicIp.status.outputs.public_ip_id` (Standard, static)
- `identity.identity_ids[]` →
  `AzureUserAssignedIdentity.status.outputs.identity_id`
- `ssl_certificates[].key_vault_secret_id` /
  `trusted_root_certificates[].key_vault_secret_id` →
  `AzureKeyVaultCertificate.status.outputs.versionless_secret_id`
- `firewall_policy_id` (gateway / listener / path rule) →
  `AzureWebApplicationFirewallPolicy.status.outputs.policy_id`
- `backend_address_pool_ids` map output is consumed by
  `AzureNetworkInterface.ip_configurations[].application_gateway_backend_address_pool_ids`
  and `AzureVirtualMachineScaleSet` NIC templates
- `application_gateway_id` output is consumed by
  `AzureAksCluster.addons.ingress_application_gateway.gateway_id` (AGIC)
- A public frontend's ADDRESS lives on the referenced `AzurePublicIp`
  (`ip_address` output) -- DNS records point there
