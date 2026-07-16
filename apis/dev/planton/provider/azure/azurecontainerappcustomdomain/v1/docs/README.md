# AzureContainerAppCustomDomain -- Design Research

## The Resource

A custom-domain binding on a Container App. Azure models it as an entry
in the app's `configuration.ingress.customDomains` array -- NOT a
standalone ARM resource -- and the provider realizes it as a locked
read-modify-write of the parent app
(`internal/services/containerapps/container_app_custom_domain_resource.go`,
Microsoft.App API 2025-07-01, azurerm v4.x), parity-verified against
pulumi-azure v6 (`containerapp.CustomDomain`) -- zero bridge lag.

## Field Mapping (azurerm -> spec)

| azurerm | spec | Notes |
|---|---|---|
| `name` | `domain_name` | FQDN CEL (no wildcard label -- per-app bindings cover one concrete hostname). ForceNew |
| `container_app_id` | `container_app_id` | FK to `AzureContainerApp.container_app_id`. ForceNew |
| `container_app_environment_certificate_id` | `container_app_environment_certificate_id` | FK to `AzureContainerAppEnvironmentCertificate.certificate_id`. ForceNew |
| `certificate_binding_type` | `certificate_binding_type` | Closed SNI_ENABLED/DISABLED/AUTO enum mapping to the provider's SniEnabled/Disabled/Auto. Provider `RequiredWith` the certificate id, front-loaded as a bidirectional pairing CEL. ForceNew |

Computed: `container_app_environment_managed_certificate_id`, exported
as the `managed_certificate_id` output.

## The Two-Flow Drift Design

Azure attaches managed certificates to the binding ASYNCHRONOUSLY, out
of band -- the provider documents an ignore-changes lifecycle on the two
certificate fields for exactly that flow. Applying the ignore
unconditionally would swallow a BYO user's legitimate certificate change
(a silent-data-loss class), so the modules dispatch:

- **Terraform**: two count-gated variants of the same resource -- `byo`
  (certificate declared; full drift detection) and `managed` (no
  certificate; `ignore_changes` on the certificate fields) -- because a
  lifecycle block is static per resource. Outputs coalesce per
  attribute.
- **Pulumi**: one resource with `IgnoreChanges` applied conditionally
  (its resource options support that directly).

Behavior is identical; only the mechanics differ per engine's lifecycle
model.

## Front-Loaded Contracts

- Certificate id and binding type together-or-neither (message CEL) --
  the provider's RequiredWith plus the managed flow's both-absent shape.
- The ingress-must-exist and asuid-TXT-must-resolve requirements are
  cross-resource / public-DNS state no spec rule can see; taught on the
  spec header, both modules, and the presets.

## Recorded Skips (with reasons)

Nothing skipped: the azurerm surface is exactly the four arguments
above.

## Operational Behavior Worth Knowing

- **Create blocks on Azure validating the asuid TXT record against
  public DNS** (the provider performs no pre-check; Azure enforces it
  during the polled operation). Publish the records first.
- All fields ForceNew -- the binding has no update surface; changing
  anything replaces it (a brief re-bind).
- The provider locks the parent app during create/delete (binding
  changes serialize per app), and re-packs the app's secrets on every
  write -- an internal mechanic both engines inherit from the provider.
