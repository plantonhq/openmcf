# AzureContainerApp: Research & Design Documentation

## Executive Summary

The Container App is the family's continuously running workload kind: containers, KEDA scaling, ingress, secrets, registries, Dapr, and identity, modeled to the full azurerm v4 surface of `azurerm_container_app`. The provider's cross-field contracts (registry auth modes, Key Vault secret/identity pairing, traffic-weight targeting, per-probe-type thresholds) are front-loaded as validation rules, and closed enums replace every string vocabulary.

## Azure Deployment Landscape

### The revision model

Every change to the template (containers, scale, volumes) creates a revision. SINGLE mode keeps one revision active (new replaces old); MULTIPLE keeps several active with `traffic_weight` splitting -- the blue-green/canary mechanism. `latest_revision_fqdn` bypasses splitting for direct verification.

### Scaling

Scaling is KEDA under the hood: HTTP/TCP concurrency rules, the Azure Queue rule, and custom rules covering the KEDA scaler catalog (the provider validates the type against its pinned allowlist, mirrored in the spec). Custom rules optionally execute under a managed identity (`identity_id` -- the literal "System" or a referenced user-assigned identity) instead of connection-string secrets. Replica bounds (0-300), the cooldown period, and the polling interval complete the scaler surface.

### Volumes

Four storage types: `EmptyDir` (ephemeral), `AzureFile` (SMB) and `NfsAzureFile` (NFS) via environment storage registrations, and `Secret` (mounts the app's secrets as files). The share-backed types pair with `storage_name` -- a reference to `AzureContainerAppEnvironmentStorage`.

### Ingress

External or environment-internal exposure; auto/http/http2/tcp transports (`exposed_port` only with tcp); required traffic weights; client certificate modes for mTLS; CORS; IP allow/deny lists (Azure evaluates the list as a single allowlist or denylist -- all rules share one action).

## Design Decisions

### 1. Location is inherited, never configured

`azurerm_container_app` has no location argument -- the app runs where its environment runs. The spec has no region field.

### 2. Probe contracts are per-type, on a shared message

Azure's probe types diverge: liveness defaults `initial_delay` to 1 (0 for the others), `success_count_threshold` exists only on readiness, and the `failure_count_threshold` ceiling is 30/48/240 for liveness/readiness/startup. One shared probe message carries the union; container-level validation rules narrow it per type, and both engines send the per-type initial-delay default when unset.

### 3. The ingress mTLS vocabulary is lowercase on the wire

ARM's `IngressClientCertificateMode` string constants are lowercase (`accept`/`require`/`ignore`) even though the SDK's Go identifiers are capitalized -- both modules map the enum to the lowercase wire values explicitly.

### 4. The KEDA scaler vocabulary is closed

The provider validates `custom_rule_type` against a pinned allowlist (~60 scalers); the spec mirrors it so a typo fails at validation time instead of plan time.

### 5. Every provider cross-field contract is a validation rule

Registry auth (identity XOR username+password, username/password travel together), secret shape (value XOR Key Vault reference; identity iff Key Vault), traffic weights (latest XOR suffix, exactly one), volume/storage-name pairing, exposed-port-requires-TCP, and queue/custom auth trigger parameters -- all front-loaded, mirroring `ValidateContainerAppRegistry` and the resource's CustomizeDiff.

### 6. `custom_domain_verification_id` is provider-Sensitive

The Terraform output carries `sensitive = true` (OpenTofu rejects the configuration at plan time without it); the environment's same-named output is not marked Sensitive by the provider and stays plain.

## Terraform Provider Analysis

### Source Files

- `internal/services/containerapps/container_app_resource.go` (schema, CustomizeDiff)
- `internal/services/containerapps/helpers/container_apps.go` (template/ingress/secret/registry/scale schemas and expanders)

### Key Behaviors

- `revision_suffix` is Optional+Computed and only sent on change (the service generates one otherwise)
- Probe removal is explicitly zeroed on update (`ContainerAppProbesRemoved`)
- Ingress custom domains are read-only on the app resource -- per-app domains are the separate custom-domain association's job

### API Version

`Microsoft.App` `2025-07-01` (via the provider's typed SDK).

## Pulumi Provider Analysis

### Package

`github.com/pulumi/pulumi-azure/sdk/v6/go/azure/containerapp` -- `App` and its nested arg types.

### Field Mapping

The bridge carries the complete surface: `AppIngressArgs` includes both `Cors` and `ClientCertificateMode`; `AppTemplateArgs` includes the cooldown/polling/termination-grace dials; custom scale rules include `IdentityId`. Every spec field maps 1:1; no parity exceptions.

## Downstream Dependencies

### References this kind consumes

- `container_app_environment_id` -> `AzureContainerAppEnvironment.environment_id`
- `volumes[].storage_name` -> `AzureContainerAppEnvironmentStorage.storage_name`
- `identity.user_assigned_identity_ids[]` -> `AzureUserAssignedIdentity.identity_id`

### Deliberately skipped surface

- Per-app custom domains (`azurerm_container_app_custom_domain`) and the environment certificate family are the Container Apps TLS/domain surface -- lifecycle-bearing resources that warrant their own kinds rather than inline fields.

## References

- Azure Container Apps: https://learn.microsoft.com/azure/container-apps/overview
- Revisions: https://learn.microsoft.com/azure/container-apps/revisions
- KEDA scalers: https://keda.sh/docs/scalers/
