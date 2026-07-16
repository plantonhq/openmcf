# AzureContainerAppEnvironmentManagedCertificate -- Design Research

## The Resource

An Azure-managed certificate on a Container App Environment
(`Microsoft.App/managedEnvironments/managedCertificates`). The component
maps onto `azurerm_container_app_environment_managed_certificate`
(azurerm v4.x,
`internal/services/containerapps/container_app_environment_managed_certificate_resource.go`,
Microsoft.App API 2025-07-01), parity-verified against pulumi-azure v6
(`containerapp.EnvironmentManagedCertificate`) -- zero bridge lag.

## Field Mapping (azurerm -> spec)

| azurerm | spec | Notes |
|---|---|---|
| `name` | `certificate_name` | The provider's certificate-name validator as a CEL. ForceNew |
| `container_app_environment_id` | `container_app_environment_id` | FK to `AzureContainerAppEnvironment.environment_id`. ForceNew |
| `subject_name` | `subject_name` | FQDN CEL that also rejects wildcards (managed certificates cannot cover them -- documented Azure behavior the provider leaves to apply time). ForceNew |
| `domain_control_validation` | `domain_control_validation` | Closed CNAME/HTTP enum; unset deploys HTTP (the provider default, sent explicitly by both engines). ForceNew |
| `tags` | `tags` | The only in-place update |

Computed: `validation_token`, exported as an output.

## Decomposition Decision

Separate from `AzureContainerAppEnvironmentCertificate` -- two different
ARM types with disjoint inputs (a domain to validate vs uploaded
material) and opposite lifecycles (a long-running operation that blocks
on public-DNS domain validation vs a synchronous PUT). See that kind's
research doc for the full reasoning.

## Front-Loaded Contracts

- The subject-name CEL rejects wildcards up front -- Azure would reject
  them at issuance, ~30 minutes into the operation.
- The DNS prerequisites (asuid TXT + CNAME/HTTP routing record) are
  public-DNS state no spec rule can see; taught prominently on the spec
  header, both modules, and the presets instead.

## Recorded Skips (with reasons)

Nothing skipped: the azurerm surface is exactly the five fields above.

## Operational Behavior Worth Knowing

- **Create blocks until Azure's domain validation succeeds** (a
  long-running operation with a ~30-minute provider timeout). The
  required DNS records must resolve on the PUBLIC internet -- an Azure
  DNS zone that is not registrar-delegated cannot satisfy validation.
- Azure attaches the issued certificate to the matching custom-domain
  binding asynchronously; the binding deploys certificate-less first.
- Renewals are Azure's job; no rotation surface exists on the resource.
