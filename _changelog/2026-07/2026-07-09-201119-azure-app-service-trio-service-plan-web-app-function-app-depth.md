# Azure App Service Trio: Service Plan, Linux Web App, and Function App Depth Rework

**Date**: July 9, 2026
**Type**: Enhancement (breaking spec rework)
**Components**: Azure Provider, API Definitions, IAC Modules (Terraform + Pulumi), E2E Framework, Secret Coverage

## Summary

Reworked the three Azure App Service kinds — `AzureServicePlan` (442), `AzureLinuxWebApp` (444), and `AzureFunctionApp` (443) — from their earlier 80/20 scoping to the full azurerm v4.80 surface, with both IaC engines at 100% behavioral parity on the shared Azure Pulumi provider builder, and proved all three live on both engines end to end (create → verify → destroy → verify-gone) for the first time. The app-hosting family now matches the depth bar set across the rest of the reworked Azure catalog.

## Problem Statement / Motivation

The App Service trio predated the catalog's depth standard. The Service Plan validated its SKU as a free-form string ("deferred to the Azure API to avoid maintaining a 50+ entry whitelist" — an implementation-convenience call), and both app kinds were missing entire azurerm blocks an advanced organization reaches: platform authentication (Easy Auth v2), scheduled backups, slot-swap sticky settings, auto-heal, hardened publishing toggles, VNet image pull, TLS cipher floors, and user tags. Their Pulumi modules also built the Azure provider inline, which silently breaks keyless (OIDC web-identity) authentication, and none of the three had ever been proven against real Azure.

### Pain Points

- Free-form SKU string on the plan; no SKU-tier gating (zone balancing, premium auto-scale, elastic worker ceiling, ASE pairing) validated before apply
- No `auth_settings_v2`, `backup`, `sticky_settings`, `auto_heal_setting` (web), client-certificate controls on the apps
- Function App storage contract (`storage_account_name` XOR `storage_key_vault_secret_id`; access key XOR managed identity) unmodeled
- Inline `azure.NewProvider` in all three Pulumi modules (keyless-auth breaker)
- String+CEL vocabularies instead of the catalog's closed proto enums
- No registry `prerequisites`, no E2E scenarios/verifiers, design-decision breadcrumbs in public artifacts

## Solution / What's New

### AzureServicePlan (442)

- Closed 52-value `AzureServicePlanSku` enum mirroring azurerm's `AllKnownServicePlanSkus()` row for row (Free/Shared, Basic, Standard, Premium v2–v4 incl. memory-optimized, Consumption Y1, Elastic Premium, Flex Consumption FC1, Isolated v1/v2 incl. memory-optimized, Workflow), family-contiguous numbering so SKU-gated rules are range checks; wire maps mechanically diffed identical across both engines
- All four provider CustomizeDiff contracts front-loaded as message CELs: premium-only auto-scale, elastic-worker-count SKU gate, zone-balancing SKU gate, ASE-requires-Isolated
- New fields: `app_service_environment_id` (ARM-ID CEL), `premium_plan_auto_scale_enabled`, user `tags`; `name` → `service_plan_name`; outputs renamed `service_plan_id`/`service_plan_name` + computed `kind`/`reserved` (both app FKs repointed)
- `AzureServicePlanOsType` enum (LINUX default / WINDOWS / WINDOWS_CONTAINER)

### AzureLinuxWebApp (444) and AzureFunctionApp (443)

- **Easy Auth v2 modeled fully** on both kinds: platform switches, unauthenticated action, forward-proxy conventions, the required `login` block (token store, cookie/nonce lifetimes with format CELs), and all nine identity-provider blocks (Entra ID with both credential forms mutually exclusive, Apple, Static Web App, custom OIDC list, Facebook, GitHub, Google, Microsoft, Twitter) — provider secrets referenced by app-setting NAME, never inline
- **Backup** (SAS URL `sensitive`, required schedule with cadence/retention CELs), **sticky settings** (not-both-empty CEL), **auto-heal** (web app: requests/status-code/slow-request triggers with the provider's format contracts as CELs)
- Function App storage contracts as CELs: exactly-one of `storage_account_name` | `storage_key_vault_secret_id`; access key conflicts with managed identity; plus `daily_memory_time_quota`, `enabled`, `zip_deploy_file`
- Hardened-publishing toggles (`webdeploy`/`ftp_publish_basic_authentication_enabled`), `vnet_image_pull_enabled` + `virtual_network_backup_restore_enabled` (both require the subnet, CEL-enforced), `minimum_tls_cipher_suite` (Azure's 17-suite vocabulary), site_config gaps closed (API Management/API definition, default documents, managed pipeline mode, remote debugging, load-balancing mode, health-check eviction), user `tags`
- Closed-vocabulary convergence: TLS versions, FTPS state, client-cert mode, load balancing, pipeline mode, IP-restriction action, connection-string types, storage-mount types, log levels, identity types — all proto enums with per-engine wire maps
- Outputs: `site_credential_name`/`site_credential_password` (both marked sensitive — azurerm marks the whole block sensitive), `possible_outbound_ip_addresses`, `hosting_environment_id`
- FK hygiene: Function App `StorageMount.access_key` gained its missing `default_kind`; registry `prerequisites` added (`[AzureServicePlan]` on both apps, `[AzureResourceGroup]` on the plan)
- Breadcrumbs stripped: all "DD0x"/"80/20" references removed from protos, READMEs, and docs; rationale restated timelessly (Linux-only and Flex-Consumption-is-its-own-resource recorded as deliberate boundaries)

### Cross-cutting

- All three Pulumi modules migrated to the shared `pulumiazureprovider.Get` builder (keyless-ready); provider.tf carries the canonical empty block + keyless comment
- First-ever E2E for the family: verifiers (`service_plan.go`, `linux_web_app.go`, `function_app.go` via ARM GetByID), scenarios, profiles, runner entrypoints

## Validation

**Offline gate (all green):** spec tests 45/84/60 Ginkgo cases (every CEL happy + error path; six Function App error-path cases added in review); targeted + release-equivalent builds; `make build-go`; `planton secret-coverage --check`; `planton validate-refs --check`; `pkg/outputs` conformance ×3; `planton validate-outputs` ×3; full `tofu plan` on the plan + web app hack manifests; `tofu validate` on the function app (its full offline plan stops at azurerm's own plan-time Service Plan tier read — a provider-inherent boundary, recorded in the audit; the live E2E carries the full plan/apply proof); presets ×9 validate; audits ×3 at 100% Fully Complete, PARITY ✅ COVERAGE ✅.

**Live dual-engine E2E green 6/6, zero orphans:**

| Scenario | Pulumi | Terraform |
|---|---|---|
| ServicePlan minimal (B1 + tags) | 142s | 177s |
| LinuxWebApp minimal (RG → plan → app chain) | 269s | 298s |
| FunctionApp minimal (RG → plan + scenario-local storage → app) | 906s | 418s |

Post-run sweep: subscription fully clean (no resource groups, plans, sites, or storage accounts).

**Live findings fixed in-session:**
- Regional App Service capacity/quota classes: new PAYG subscriptions carry zero Basic-tier VM quota in some regions (a 401 wearing a quota message) and `eastus`/`westus3` rejected creates on capacity — scenarios re-regioned to `centralus`; class recorded in `e2e/README.md`
- Terraform outputs surfacing provider-`Sensitive` attributes (`site_credential[0].name`, the function app's `custom_domain_verification_id`) failed OpenTofu's plan-time sensitivity check — fixed with `sensitive = true` on both kinds

## Workflow Uplift (durable learnings)

- `forge/flow/013-terraform-module.mdc`: outputs surfacing provider-`Sensitive` attributes must be marked sensitive — name-like fields can be sensitive too; `terraform validate` cannot catch it
- `forge/flow/008-hack-manifest.mdc`: the stale-CLI class (rebuild the CLI after proto changes before `planton tofu plan`; "unknown field" means a stale binary, not a bad manifest) and the offline-plan boundary class (some resources read the cloud in CustomizeDiff)
- `e2e/README.md`: the App Service regional quota/capacity class and the plan-pins-the-region pairing

## Related Hygiene

A stale `planton` CLI binary falsely reported two `pkg/secretcoverage/baseline.yaml` entries (`GcpCloudFunction`/`GcpCloudRun` env-secret maps) as stale, and they were briefly removed. The fields are genuine unannotated gaps — the repo's authoritative gate (`go test ./pkg/secretcoverage/...`) failed on the removal — so the entries were restored in a follow-up commit. The lesson is the same stale-CLI class recorded in the forge rules: rebuild the CLI before trusting its analyzers, and treat the in-repo test as the source of truth for the secret-coverage gate.

## Impact

The Azure app-hosting family (compute tier + both app kinds) now carries the same completeness, validation depth, parity discipline, and live proof as the rest of the reworked catalog. Anyone deploying Linux web apps or functions through Planton gets platform authentication, backups, hardened publishing, and VNet-integrated postures as first-class, validated spec surface instead of undocumented gaps.

---

**Status**: ✅ Production Ready
