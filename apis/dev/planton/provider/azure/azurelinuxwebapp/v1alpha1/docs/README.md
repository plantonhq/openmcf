# AzureLinuxWebApp: Research & Design Documentation

## Executive Summary

Azure Linux Web App (`Microsoft.Web/sites`, kind `app,linux`) is App Service's managed hosting platform for long-running HTTP workloads: web frontends, APIs, and containerized services. Unlike serverless Function Apps that are event-driven and scale to zero, Web Apps are always-on.

This document captures the research and design rationale behind the `AzureLinuxWebApp` Planton component (enum 444, id_prefix `azweb`). The component models the complete `azurerm_linux_web_app` surface: the full runtime matrix, site configuration with auto-heal, built-in authentication (Easy Auth v2), scheduled backups, sticky settings, deployment hardening, storage mounts, dual-destination logging, VNet integration with image-pull and backup routing, managed identity, and user tags.

## Azure Deployment Landscape

### How a Web App relates to its plan

Every Web App runs on an App Service Plan (`AzureServicePlan`), which fixes the region, OS, VM size, and instance count. Multiple apps share a plan's instances; per-site scaling on the plan lets one app scale independently. The plan tier gates features the app can use:

| Capability | Minimum tier |
|------------|--------------|
| Always On | Basic |
| Custom domains + SSL | Basic |
| Staging slots, scheduled `backup` | Standard |
| VNet integration | Basic (regional) |
| Zone redundancy, premium auto-scale | Premium |
| Single-tenant compute | Isolated (ASE) |

### Runtime model

Exactly one runtime is selected in `site_config.application_stack`: .NET, Node.js, Python, PHP, Ruby (legacy), Go (legacy), Java (with an application server: embedded SE, Tomcat, or JBoss EAP), or a Docker container. Containers pull from any registry -- anonymously, with a username/password, or with the app's managed identity against ACR (`container_registry_use_managed_identity`), optionally over the VNet (`vnet_image_pull_enabled`).

### Built-in authentication (Easy Auth v2)

App Service can authenticate requests at the platform layer, before application code runs. The v2 model is the current surface: platform settings (enabled, runtime version), global gating (`require_authentication`, `unauthenticated_action`, `excluded_paths`), login/session behavior (token store, cookie expiration, nonce), and per-provider blocks -- Microsoft Entra ID, Apple, Azure Static Web Apps, Facebook, GitHub, Google, Microsoft account, Twitter, and any number of custom OpenID Connect providers. Every provider references its client secret by APP SETTING NAME; the secret value lives in `app_settings` (commonly as a Key Vault reference), never in the auth configuration itself.

### Auto-heal

App Service can recycle an unhealthy app automatically when a trigger fires: total request volume, status-code patterns (single codes or ranges, optionally narrowed by sub-status/win32 status/path), or slow requests (globally or per path). Linux supports exactly one heal action -- recycling the process -- with a minimum-uptime guard against recycle loops.

## Design Decisions

### 1. Model auth v2 only; v1 is a recorded skip

`auth_settings` (v1) is the legacy Easy Auth surface kept by Azure for backward compatibility; `auth_settings_v2` covers everything v1 does with the current provider vocabulary. Modeling both would create two ways to express the same intent with silently conflicting semantics. New configuration always uses v2.

### 2. Secrets in the auth block are referenced, never carried

The auth model's `*_setting_name` fields are pointers into `app_settings` -- the platform's own design. The spec preserves this: no secret values flow through `auth_settings_v2`, so the block is never secret-bearing; the referenced app settings pair naturally with Key Vault references.

### 3. The recycle action is implicit

Linux auto-heal supports exactly one action. A one-value knob is a constant, not a field -- both IaC modules send `Recycle`, and the spec configures only the trigger and the minimum-uptime guard.

### 4. Closed vocabularies are proto enums

Client certificate mode, TLS versions, FTPS state, load balancing mode, pipeline mode, IP restriction actions, connection-string types, storage mount types, log levels, identity types, backup frequency units, and the auth flow vocabularies (unauthenticated action, forward proxy convention, cookie expiration convention) are all closed sets in the provider -- modeled as enums, validated at manifest time, and mapped row-by-row to Azure's wire spellings in both IaC modules. Open-ended values (runtime versions, cipher suite identifiers) stay CEL-validated strings.

### 5. Provider contracts are front-loaded as spec rules

- The three VNet-riding toggles (`vnet_route_all_enabled`, `vnet_image_pull_enabled`, `virtual_network_backup_restore_enabled`) require `virtual_network_subnet_id`.
- `health_check_eviction_time_in_min` requires `health_check_path`.
- CORS credentials cannot combine with a wildcard origin.
- Sticky settings must name at least one setting.
- HTTP logs take exactly one destination (file system XOR blob storage).
- The auto-heal trigger must carry at least one condition.
- Java version/server/server-version travel together.
- Entra ID auth takes a secret setting name XOR a certificate thumbprint; the login token store takes a path XOR a SAS setting name; custom forward-proxy headers require the CUSTOM convention.
- The identity block pairs `identity_ids` with USER_ASSIGNED types.

One contract stays apply-time: `always_on` is rejected on Free/Shared plans (the plan's SKU is a different resource's property, invisible to manifest validation).

### 6. Deployment hardening is first-class

The basic-auth publishing toggles (`webdeploy_publish_basic_authentication_enabled`, `ftp_publish_basic_authentication_enabled`) close the classic credential-based deployment paths, and the `site_credential_password` output documents exactly what those toggles revoke. `ftps_state` (the FTP endpoint's TLS posture) defaults to DISABLED.

### 7. Opinionated defaults, each documented on its field

- `https_only` defaults true (Azure's default is false).
- `use_32_bit_worker` defaults false -- 64-bit workers for production (the provider's own default is true).
- `client_affinity_enabled` defaults false -- stateless apps balance better.

### 8. Deliberate scope boundaries

- **Deployment slots** mirror the entire app surface with an independent lifecycle -- a standalone-kind candidate, not a field here. `sticky_settings` lives on this spec because it belongs to the production app.
- **Custom hostname bindings and certificates** are separate Azure resources that compose with the app by reference.
- **Windows Web Apps** are the legacy .NET Framework path; the platform targets Linux-first runtimes and containers.

## Terraform Provider Analysis

### Source Files

- `internal/services/appservice/linux_web_app_resource.go` -- resource implementation
- `internal/services/appservice/helpers/linux_web_app_schema.go` -- site_config + application_stack + Linux auto-heal
- `internal/services/appservice/helpers/shared_schema.go` -- IP restrictions, CORS, sticky settings, identity
- `internal/services/appservice/helpers/common_web_app_schema.go` -- backup, logs, storage mounts, connection strings
- `internal/services/appservice/helpers/auth_v2_schema.go` -- the Easy Auth v2 model

### Key Behaviors

1. **Name**: 2-60 characters, alphanumeric + hyphens, globally unique (forms the `azurewebsites.net` hostname). ForceNew with region and resource group; the plan reference is NOT ForceNew.
2. **Docker fields are flattened** on the provider (`docker_image_name` carries `image:tag`); the spec keeps the honest nested shape and both modules flatten identically.
3. **Outbound IP attributes** arrive as comma-joined strings; both modules split them into real lists so the repeated outputs populate identically.
4. **site_credential** is a computed block carrying the publishing password -- exported as secret-bearing outputs.
5. **auth v2 create-time checks**: the provider validates provider-block pairings (e.g. the Entra credential XOR) at expand time; the spec front-loads them.

### API Version

- Azure API: `Microsoft.Web` version `2023-12-01`
- Resource ID: `/subscriptions/{sub}/resourceGroups/{rg}/providers/Microsoft.Web/sites/{name}`

## Pulumi Provider Analysis

- `github.com/pulumi/pulumi-azure/sdk/v6/go/azure/appservice`, resource `appservice.LinuxWebApp`
- The provider is built through the shared `pulumiazureprovider` builder (static client secret, keyless web identity, or ambient chain)
- All spec blocks map 1:1 onto `LinuxWebAppArgs` (auth v2 = `AuthSettingsV2`, backup = `Backup`, sticky = `StickySettings`, auto-heal = `SiteConfig.AutoHealSetting`); enum wire maps mirror the Terraform module's lookup tables row for row

## Downstream Dependencies

AzureLinuxWebApp is a leaf: nothing references its outputs. It references, via `StringValueOrRef`:

| Field | Kind | Output consumed |
|-------|------|-----------------|
| `resource_group` | AzureResourceGroup | `resource_group_name` |
| `service_plan_id` | AzureServicePlan | `service_plan_id` |
| `application_insights_connection_string` | AzureApplicationInsights | `connection_string` |
| `virtual_network_subnet_id` (+ per-rule IP restriction subnets) | AzureSubnet | `subnet_id` |
| `identity.identity_ids`, `key_vault_reference_identity_id` | AzureUserAssignedIdentity | `identity_id` |
| `storage_mounts[].access_key` | AzureStorageAccount | `primary_access_key` |

### Infra Charts

| Chart | Role |
|-------|------|
| `web-app-environment` | Leaf workload (Service Plan -> Linux Web App) |

## References

- [App Service Linux documentation](https://learn.microsoft.com/en-us/azure/app-service/overview)
- [Terraform azurerm_linux_web_app](https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs/resources/linux_web_app)
- [App Service built-in authentication](https://learn.microsoft.com/en-us/azure/app-service/overview-authentication-authorization)
- [App Service auto-heal](https://learn.microsoft.com/en-us/azure/app-service/overview-diagnostics#auto-healing)
