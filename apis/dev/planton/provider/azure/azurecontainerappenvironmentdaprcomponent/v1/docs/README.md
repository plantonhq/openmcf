# AzureContainerAppEnvironmentDaprComponent -- Design Research

## The Resource

A Dapr component registration
(`Microsoft.App/managedEnvironments/daprComponents`) is a pluggable
backend behind Dapr's building blocks -- state stores, pub/sub brokers,
secret stores, bindings -- registered once on the environment and
consumed by Dapr-enabled apps through the Dapr runtime. The component
maps onto `azurerm_container_app_environment_dapr_component` (azurerm
v4.x, `internal/services/containerapps/container_app_environment_dapr_component_resource.go`),
parity-verified against pulumi-azure v6
(`containerapp.EnvironmentDaprComponent`).

## Field Mapping (azurerm -> spec)

| azurerm | spec | Notes |
|---|---|---|
| `container_app_environment_id` | `container_app_environment_id` | Single parent ARM-id FK to `AzureContainerAppEnvironment.environment_id`. ForceNew |
| `name` | `component_name` | What app code passes to the Dapr API; the provider's name contract mirrored as a CEL (lowercase alphanumerics/hyphens, <=60). ForceNew |
| `component_type` | `component_type` | Dapr's dotted type notation (e.g. `state.azure.blobstorage`); an open string -- Dapr's component catalog evolves independently of Azure and the provider does not validate it. ForceNew |
| `version` | `version` | Required; "v1" for virtually all stable components. Updatable |
| `init_timeout` | `init_timeout` | Optional, default "5s"; the provider's interval format (`[0-9]+[smh]`) mirrored as a CEL |
| `ignore_errors` | `ignore_errors` | Optional bool, default false -- left false so a broken backend fails loudly at sidecar startup |
| `secret` | `secrets[]` | Dapr component secrets are name/value pairs (NOT the app-level secret schema -- no Key Vault reference or identity fields on this resource); values are sensitive |
| `metadata` | `metadata[]` | name + (value XOR secret_name), the XOR as a message CEL |
| `scopes` | `scopes[]` | The consuming apps' `dapr.app_id`s; empty exposes the component to every Dapr-enabled app |

No tags: ARM does not support tags on Dapr component registrations.

## Decomposition Decision

First-class kind, not a fold into the environment: components are
many-per-environment with independent lifecycles (a new state store or
broker registration should never mean an environment update), matching
ARM's own child-resource grain.

## Front-Loaded Contracts

- **Metadata value XOR secret_name** (message CEL) -- connection strings
  and keys must travel as component secrets, never inline metadata.
- **init_timeout format** (`^[0-9]+[smh]$`) mirrors the provider's
  validator exactly.

## Recorded Skips (with reasons)

Nothing skipped: the azurerm surface is exactly the fields above and the
spec models all of them. Note the provider expands this resource's
`secret` blocks as Dapr secrets (name/value only) -- the Key Vault
fields on the shared app-level secret schema do not exist on this ARM
resource, so the spec deliberately models a narrower secret message.

## Operational Behavior Worth Knowing

- Component configuration errors typically surface when the SIDECAR
  initialises, not when the registration deploys -- a healthy apply does
  not prove the metadata is right for the backend. Keep `ignore_errors`
  false so misconfiguration fails visibly.
- Scoping is by the app's `dapr.app_id` (its Dapr identity), not the
  Container App resource name.
