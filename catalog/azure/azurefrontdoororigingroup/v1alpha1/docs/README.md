# AzureFrontDoorOriginGroup -- Design Research

## Scope

The origin group is the pool-level behavior of a Front Door backend:
health probing, latency-aware selection, session affinity, and the
recovery ramp. Backends decompose into AzureFrontDoorOrigin (a separate
ARM child resource with an independent lifecycle) so per-region stamps
compose without touching the group.

Source of truth: `azurerm_cdn_frontdoor_origin_group`
(terraform-provider-azurerm v4.80,
`internal/services/cdn/cdn_frontdoor_origin_group_resource.go`),
parity-verified against pulumi-azure v6 (`cdn.FrontdoorOriginGroup`).

## Field mapping

| Spec field | azurerm attribute | Notes |
| --- | --- | --- |
| `profile_id` | `cdn_frontdoor_profile_id` | FK to AzureFrontDoorProfile; ForceNew |
| `origin_group_name` | `name` | ForceNew; 2-90 chars (the provider's own regex) |
| `load_balancing.sample_size` | `load_balancing.sample_size` | 0-255, default 4 |
| `load_balancing.successful_samples_required` | `load_balancing.successful_samples_required` | 0-255, default 3 |
| `load_balancing.additional_latency_in_milliseconds` | `load_balancing.additional_latency_in_milliseconds` | 0-1000, default 50 |
| `health_probe.protocol` | `health_probe.protocol` | Http/Https enum; block absence disables probing |
| `health_probe.interval_in_seconds` | `health_probe.interval_in_seconds` | 1-255, required in the block |
| `health_probe.request_type` | `health_probe.request_type` | HEAD (default) / GET |
| `health_probe.path` | `health_probe.path` | default "/" |
| `session_affinity_enabled` | `session_affinity_enabled` | default true |
| `restore_traffic_time_to_healed_or_new_endpoint_in_minutes` | same | 0-50, default 10 |

## Outputs

| Output | Source | Consumers |
| --- | --- | --- |
| `origin_group_id` | resource id | AzureFrontDoorOrigin parent FK; AzureFrontDoorRoute destination FK |
| `origin_group_name` | `name` | human/portal orientation |

## Behavior worth knowing

- **`load_balancing` is required by Azure on every group** -- both
  modules always send the block, materializing Azure's defaults when
  the spec omits it; an unset spec block and Azure's defaults are the
  same thing by design.
- **Health-probe absence is a real behavior** (probing disabled), not a
  defaults shortcut -- the Front Door API needs an explicit null to
  disable probes, and azurerm ships a PATCH workaround client so
  unrelated updates don't silently null probe settings. Both engines
  inherit that from the provider layer.
- **No tags**: ARM does not support tags on origin groups; identity
  tags live on the profile.

## Recorded skips (with reasons)

- None -- the azurerm surface is fully modeled.
