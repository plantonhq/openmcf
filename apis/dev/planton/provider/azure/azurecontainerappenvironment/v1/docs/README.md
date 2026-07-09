# AzureContainerAppEnvironment: Research & Design Documentation

## Executive Summary

The Container App Environment is the container kind of the Azure Container Apps family: the secure boundary providing networking, logging, Dapr infrastructure, and compute capacity for every Container App and Job inside it. The spec models the full azurerm v4 surface of `azurerm_container_app_environment` -- networking modes, logging destinations, workload profiles (dedicated, GPU, serverless GPU), mTLS, managed identity, custom DNS suffix, and tags -- with the provider's cross-field contracts front-loaded as validation rules.

## Azure Deployment Landscape

### The environment/workload split

Azure models Container Apps as a two-level hierarchy: the managed environment (`Microsoft.App/managedEnvironments`) owns the shared infrastructure, and workloads (`containerApps`, `jobs`) plus platform registrations (`storages`, `daprComponents`) live inside it. The catalog mirrors that grain exactly: this kind is the parent; the family kinds reference its `environment_id` output.

### Networking modes

| Mode | Configuration | Behavior |
| --- | --- | --- |
| External | no subnet | Azure-managed networking, public static IP, public default domain |
| VNet-injected | `infrastructure_subnet_id` | Apps run inside the customer VNet (`/21`+ subnet); private connectivity to VNet resources |
| Internal | subnet + `internal_load_balancer_enabled` | Private static IP; apps reachable only inside the VNet |

Zone redundancy also requires the subnet -- both pairings are spec-enforced, mirroring the provider's `RequiredWith` contracts.

### Logging destinations

`logs_destination` is a three-way choice: `log-analytics` (paired with a workspace), `azure-monitor` (routed through diagnostic settings configured separately), or absent (streaming-only). In azurerm v4 the attribute is Optional+Computed with legacy inference from the workspace's presence; the spec models it as a closed enum whose unspecified value reproduces that inference (workspace present -> log-analytics), and the CustomizeDiff pairing rules are front-loaded as validations.

### Workload profiles

The profile vocabulary is Azure's own 14-value SKU list: `Consumption`, serverless GPU (`Consumption-GPU-NC8as-T4`, `Consumption-GPU-NC24-A100`), general-purpose `D4-D32`, memory-optimized `E4-E32`, and dedicated GPU `NC24/48/96-A100`. Instance counts only apply to dedicated families (spec-enforced; the provider special-cases Consumption at expand time). The 0-profiles-to-some transition (either direction) forces replacement -- the provider's CustomizeDiff -- while changes within a non-empty set apply in place.

## Design Decisions

### 1. The custom DNS suffix folds into the environment

ARM models the environment-wide custom domain as a PATCH on the environment itself -- the association resource's ID IS the environment ID, exactly one per environment. That fails every split test (no independent lifecycle, not referenced, one-per-parent), so it lives in the spec as an optional message; both engines realize it through the standalone association resource.

### 2. Workload profile types are a closed enum

The provider validates against a hardcoded SKU list, so the spec mirrors it row for row and both modules map enum values to wire spellings explicitly -- a vocabulary drift fails loudly at plan/preview time instead of deploying a wrong profile.

### 3. `dapr_application_insights_connection_string` is sensitive and ForceNew

The connection string embeds the telemetry instrumentation key (a write credential) and ARM never returns it on read -- which is also why the provider marks it ForceNew.

### 4. Identity follows the family identity shape

A closed type enum (SYSTEM_ASSIGNED / USER_ASSIGNED / SYSTEM_AND_USER_ASSIGNED) plus `AzureUserAssignedIdentity` references, with the ids-match-type contract as a validation rule.

### 5. `public_network_access` stays optional-unspecified

Azure derives the value from the network configuration when unset (Enabled externally, Disabled behind an ILB), so the enum's unspecified value sends nothing; the cannot-be-Enabled-with-ILB conflict is front-loaded.

## Terraform Provider Analysis

### Source Files

- `internal/services/containerapps/container_app_environment_resource.go` (schema, CustomizeDiff)
- `internal/services/containerapps/container_app_environment_custom_domain_resource.go` (the folded custom-domain association)
- `internal/services/containerapps/helpers/container_app_environment.go` (workload profile schema)

### Key Behaviors

- Create fetches the Log Analytics shared key itself when a workspace is linked; Read resolves the workspace back from its customer ID
- Update runs through a PATCH workaround client (the SDK's update payload drops fields otherwise)
- The workload-profile set carries a diff suppressor for the extra Consumption profile the API returns
- mTLS writes both `PeerAuthentication.Mtls` and `PeerTrafficConfiguration.Encryption`

### API Version

`Microsoft.App` `2025-07-01` (via the provider's typed SDK).

## Pulumi Provider Analysis

### Package

`github.com/pulumi/pulumi-azure/sdk/v6/go/azure/containerapp` -- `Environment`, `EnvironmentCustomDomain`.

### Field Mapping

Every spec field maps 1:1 onto `EnvironmentArgs` (logs destination, public network access, mTLS, infra RG, identity, tags, workload profiles); the custom domain maps onto the standalone `EnvironmentCustomDomain` resource, parented to the environment. No parity exceptions.

## Downstream Dependencies

### Resources that reference AzureContainerAppEnvironment

- `AzureContainerApp.container_app_environment_id`
- `AzureContainerAppJob.container_app_environment_id`
- `AzureContainerAppEnvironmentStorage.container_app_environment_id`
- `AzureContainerAppEnvironmentDaprComponent.container_app_environment_id`

### Deliberately skipped surface

- `azurerm_container_app_environment_certificate` / `_managed_certificate` and per-app custom domains are the Container Apps TLS/domain family -- separate lifecycle-bearing resources that warrant their own kinds; they are not part of this component.

## References

- Azure Container Apps environments: https://learn.microsoft.com/azure/container-apps/environment
- Workload profiles: https://learn.microsoft.com/azure/container-apps/workload-profiles-overview
- Custom DNS suffix: https://learn.microsoft.com/azure/container-apps/custom-domains-managed-certificates
