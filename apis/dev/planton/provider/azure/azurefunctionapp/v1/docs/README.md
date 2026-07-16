# AzureFunctionApp: Research & Design Documentation

## Executive Summary

Azure Linux Function App (`Microsoft.Web/sites`, kind `functionapp,linux`) is App Service's serverless compute platform: functions triggered by HTTP, queues, timers, blobs, and other Azure events. The hosting plan decides the cost model -- Consumption (pay-per-execution, scale to zero), Elastic Premium (pre-warmed instances, no cold starts), or Dedicated tiers.

This document captures the research and design rationale behind the `AzureFunctionApp` Planton component (enum 443, id_prefix `azfn`). The component models the complete `azurerm_linux_function_app` surface: all three storage bindings, the full runtime matrix, serverless scaling dials, built-in authentication (Easy Auth v2), scheduled backups, sticky settings, deployment hardening, storage mounts, VNet integration with image-pull and backup routing, managed identity, and user tags.

## Azure Deployment Landscape

### The storage requirement

Every Function App needs an Azure Storage account for runtime state (trigger coordination, execution logs, the content share). Three binding forms exist, and the provider enforces their pairing rules at create time -- the spec front-loads them:

1. **Account name + access key** -- the classic form; the key is static credential material.
2. **Account name + managed identity** (`storage_uses_managed_identity`) -- credential-free; the app's identity needs Storage Blob Data Owner and Storage Queue Data Contributor on the account. Conflicts with the access key.
3. **Key Vault secret ID** (`storage_key_vault_secret_id`) -- the connection string lives in a vault secret; exactly-one-of with the account name.

### Plan-tier behavior

| Tier | Scaling | Cold start | Cost controls |
|------|---------|-----------|----------------|
| Consumption (Y1) | 0-200 automatic | Yes | `app_scale_limit`, `daily_memory_time_quota` |
| Elastic Premium (EP*) | elastic to 100 | Eliminated via `elastic_instance_minimum` + `pre_warmed_instance_count` | `app_scale_limit`, the plan's `maximum_elastic_worker_count` |
| Dedicated (B*/S*/P*) | plan-controlled | Use `always_on` | plan `worker_count` |

Flex Consumption (FC1) is a SEPARATE Azure resource type (`azurerm_function_app_flex_consumption`) with a container-endpoint storage model, top-level runtime selection, and per-instance memory sizing -- recorded as a standalone-kind candidate rather than folded into this spec.

### Built-in authentication (Easy Auth v2)

Identical to the Web App surface: the platform authenticates requests before function code runs, against Entra ID, social providers, or any OpenID Connect provider, with client secrets referenced by app setting name.

## Design Decisions

### 1. All three storage bindings, with the contracts front-loaded

The storage account was previously a hard-required field; the full provider surface makes it one of three binding forms. The spec enforces exactly-one-of (account name XOR Key Vault secret ID) and the key-XOR-managed-identity conflict at manifest time, exactly as the provider enforces them at plan time.

### 2. Model auth v2 only; v1 is a recorded skip

Same rationale as the Web App kind: `auth_settings_v2` covers everything the legacy v1 block does with the current vocabulary; modeling both would create conflicting ways to express the same intent.

### 3. Closed vocabularies are proto enums

Client certificate mode, TLS versions, FTPS state, load balancing mode, pipeline mode, IP restriction actions, connection-string types, storage mount types, identity types, backup frequency units, and the auth flow vocabularies -- all closed sets, modeled as enums and mapped row-by-row to Azure's wire spellings in both IaC modules. Runtime versions and cipher suites stay CEL-validated strings.

### 4. Provider contracts front-loaded as spec rules

- Storage: exactly-one binding; access key XOR managed identity.
- The VNet-riding toggles require `virtual_network_subnet_id`.
- `health_check_eviction_time_in_min` requires `health_check_path`.
- CORS credentials cannot combine with a wildcard origin.
- Sticky settings must name at least one setting; HTTP-time formats validated; the auth pairings (Entra credential XOR, token store XOR, custom-proxy headers) mirror the provider's expand-time checks.
- The identity block pairs `identity_ids` with USER_ASSIGNED types.

Apply-time contracts that reference the PLAN's tier (backup needs Standard+, `always_on` rejected on Free/Shared, VNet features rejected on Consumption) stay documented on the fields -- the plan's SKU is a different resource's property, invisible to manifest validation.

### 5. Deployment hardening is first-class

Basic-auth publishing toggles (Web Deploy + FTP), `ftps_state` defaulting to DISABLED, and the secret-bearing `site_credential_password` output documenting exactly what those toggles revoke.

### 6. Deliberate scope boundaries

- **Flex Consumption** -- its own resource type, its own future kind.
- **Deployment slots** -- standalone-kind candidate; `sticky_settings` lives here because it belongs to the production app.
- **Windows Function Apps** -- the legacy path; Linux-first platform.

## Terraform Provider Analysis

### Source Files

- `internal/services/appservice/linux_function_app_resource.go` -- resource implementation
- `internal/services/appservice/helpers/function_app_schema.go` -- site_config + function application_stack
- `internal/services/appservice/helpers/shared_schema.go` / `common_web_app_schema.go` -- IP restrictions, CORS, sticky settings, backup, storage mounts
- `internal/services/appservice/helpers/auth_v2_schema.go` -- the Easy Auth v2 model

### Key Behaviors

1. **Storage ExactlyOneOf**: `storage_account_name` XOR `storage_key_vault_secret_id`; the access key conflicts with managed identity.
2. **Function docker is a nested block** (unlike the Web App's flattened fields); both modules pass the nested shape through.
3. **CustomizeDiff**: VNet image pull cannot be disabled on ASE-hosted apps and cannot be enabled on Consumption; backup is rejected on Dynamic/Basic tiers; Dynamic <-> non-Dynamic plan moves force replacement.
4. **Outbound IP lists**: the provider exposes list-typed attributes; both modules export real lists.
5. **site_credential** is a computed block carrying the publishing password -- exported as secret-bearing outputs.

### API Version

- Azure API: `Microsoft.Web` version `2023-12-01`
- Resource ID: `/subscriptions/{sub}/resourceGroups/{rg}/providers/Microsoft.Web/sites/{name}`

## Pulumi Provider Analysis

- `github.com/pulumi/pulumi-azure/sdk/v6/go/azure/appservice`, resource `appservice.LinuxFunctionApp`
- The provider is built through the shared `pulumiazureprovider` builder (static client secret, keyless web identity, or ambient chain)
- All spec blocks map 1:1 onto `LinuxFunctionAppArgs`; enum wire maps mirror the Terraform module's lookup tables row for row

## Downstream Dependencies

AzureFunctionApp is a leaf: nothing references its outputs. It references, via `StringValueOrRef`:

| Field | Kind | Output consumed |
|-------|------|-----------------|
| `resource_group` | AzureResourceGroup | `resource_group_name` |
| `service_plan_id` | AzureServicePlan | `service_plan_id` |
| `storage_account_name` / `storage_account_access_key` / `storage_mounts[].access_key` | AzureStorageAccount | `storage_account_name` / `primary_access_key` |
| `application_insights_connection_string` | AzureApplicationInsights | `connection_string` |
| `virtual_network_subnet_id` (+ per-rule IP restriction subnets) | AzureSubnet | `subnet_id` |
| `identity.identity_ids`, `key_vault_reference_identity_id` | AzureUserAssignedIdentity | `identity_id` |

### Infra Charts

| Chart | Role |
|-------|------|
| `function-app-environment` | Leaf workload (Service Plan + Storage -> Function App) |

## References

- [Azure Functions documentation](https://learn.microsoft.com/en-us/azure/azure-functions/)
- [Terraform azurerm_linux_function_app](https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs/resources/linux_function_app)
- [Azure Functions hosting options](https://learn.microsoft.com/en-us/azure/azure-functions/functions-scale)
- [App Service built-in authentication](https://learn.microsoft.com/en-us/azure/app-service/overview-authentication-authorization)
