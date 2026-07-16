# AzureContainerRegistry Terraform Module

## Overview

This Terraform module provisions an Azure Container Registry using the
`azurerm` provider (`~> 4.0`). It creates an `azurerm_container_registry`
-- the managed, private OCI registry the platform's workloads pull their
images from -- covering the full surface: SKU, admin account, network
posture (public access, IP rule set, bypass option, data endpoints),
policies (quarantine, retention, content trust, export), geo-replications,
managed identity, and customer-managed-key encryption.

The SKU is the registry's feature gate; the Premium-only fields
(geo-replication, zone redundancy, network rules, policies, CMK) are
enforced at spec validation, mirroring ARM's own gates, so the module
receives only deployable shapes. An unset `sku` deploys the STANDARD
baseline (azurerm requires an explicit value; the default is materialized
in locals).

Lifecycle notes worth knowing before operating this resource: name and
region are the registry's identity -- changing either replaces it and its
CONTENTS DO NOT MIGRATE; every image would need re-pushing. Zone
redundancy (home replica) and CMK encryption are likewise fixed at
creation. The SKU changes in place, but downgrading requires every
Premium-only feature to be unset first. Geo-replications are managed
inline and add/remove in place; azurerm expects them in alphabetical
location order, and the module produces exactly that by iterating a
location-keyed map, keeping manifests order-insensitive.

Grants are composed, never bundled: AcrPull/AcrPush role assignments are
standalone AzureRoleAssignment resources scoped to the registry's ARM ID,
and the identity that unwraps a CMK key is a referenced first-class
AzureUserAssignedIdentity that must hold get/wrapKey/unwrapKey on the
key's vault before the registry is created.

## Resources Created

- `azurerm_container_registry.main` -- the registry, including inline
  geo-replications, network rule set, identity, and encryption blocks

## Variables

| Variable | Type | Description |
|----------|------|-------------|
| `metadata` | object | Planton metadata (name, org, env) |
| `spec` | object | Container registry specification |

Key spec fields:

| Field | Required | Description |
|-------|----------|-------------|
| `region` | yes | The home replica's region; additional regions are geo-replications |
| `resource_group` | yes | Resource group name |
| `registry_name` | yes | 5-50 lowercase alphanumerics, globally unique (becomes `{name}.azurecr.io`) |
| `sku` | no | `BASIC` / `STANDARD` / `PREMIUM` enum name; unset applies the STANDARD baseline |
| `admin_user_enabled` | no | Built-in admin account (default false; Entra auth is the production path) |
| `public_network_access_enabled` | no | Default true; false = private-endpoints-only (Premium) |
| `zone_redundancy_enabled` | no | Home replica across availability zones (Premium; fixed at creation) |
| `anonymous_pull_enabled` | no | Unauthenticated pulls (Standard/Premium; makes every repo public) |
| `data_endpoint_enabled` | no | Dedicated `{name}.{region}.data.azurecr.io` endpoints for exact firewall allowlisting (Premium) |
| `quarantine_policy_enabled` | no | Hold pushed images until scanning tooling passes them (Premium) |
| `retention_policy_in_days` | no | Purge untagged manifests after N days, 0-365; unset keeps them forever (Premium) |
| `trust_policy_enabled` | no | Docker Content Trust / image signing (Premium) |
| `export_policy_enabled` | no | Default true; disabling requires Premium + public access off |
| `network_rule_bypass_option` | no | `AZURE_SERVICES` (Azure's default) or `NONE` |
| `network_rule_set` | no | Default action (`ALLOW`/`DENY`) plus IPv4 CIDR allowlist (Premium) |
| `georeplications` | no | Additional regions with per-replica zone redundancy, regional endpoint, tags (Premium; must not contain the home region) |
| `identity` | no | `SYSTEM_ASSIGNED` / `USER_ASSIGNED` / `SYSTEM_AND_USER_ASSIGNED` plus resolved identity ARM IDs |
| `encryption` | no | CMK: unwrapping identity's client id + Key Vault key ID (Premium; fixed at creation) |
| `tags` | no | User tags, merged over metadata-derived tags (user wins) |

## Outputs

| Output | Description |
|--------|-------------|
| `container_registry_id` | Full ARM ID of the registry -- the join key AKS clusters and AcrPull/AcrPush role assignments reference |
| `container_registry_name` | The registry's name as deployed |
| `login_server` | The hostname images are tagged with and pulled from, e.g. `myregistry.azurecr.io` |
| `admin_username` | Admin account username (empty unless `admin_user_enabled`) |
| `admin_password` | One of the admin account's two rotatable passwords (sensitive; empty unless enabled) |
| `system_assigned_identity_principal_id` | Principal ID of the system-assigned identity (empty unless the identity type includes it) |
| `data_endpoint_host_names` | Dedicated regional data-endpoint hostnames (empty unless `data_endpoint_enabled`) |

## Usage

```hcl
module "container_registry" {
  source = "./iac/tf"

  metadata = { name = "prod-registry", org = "mycompany", env = "production" }

  spec = {
    region        = "eastus"
    resource_group = "prod-rg"
    registry_name = "prodregistry01"
    sku           = "PREMIUM"
    zone_redundancy_enabled  = true
    retention_policy_in_days = 30
    georeplications = [
      { location = "westeurope", zone_redundancy_enabled = true }
    ]
  }
}
```

## Required Permissions

The deploying credential needs
`Microsoft.ContainerRegistry/registries/write` on the resource group --
held via Contributor or Owner. For CMK encryption, the referenced
user-assigned identity (not the deployer) must hold get/wrapKey/unwrapKey
on the Key Vault key's vault before the registry is created.
