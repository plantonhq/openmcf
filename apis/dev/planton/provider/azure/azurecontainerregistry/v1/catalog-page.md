# Azure Container Registry

Deploys an Azure Container Registry (ACR) — the managed, private OCI registry that stores the container images and artifacts a platform's workloads run. The SKU is the registry's feature gate, and the spec mirrors Azure's own tiering rather than hiding it: every Premium-only field (geo-replication, zone redundancy, network isolation, policies, CMK encryption) is validated against the chosen tier, so a misconfigured manifest fails at validation, not at apply. What the registry composes with is referenced, never created here: the user-assigned identity that unwraps a customer-managed encryption key is a first-class `AzureUserAssignedIdentity`, AKS clusters pull by referencing the `container_registry_id` output, and `AcrPull`/`AcrPush` grants are standalone `AzureRoleAssignment` resources scoped to the registry.

## What Gets Created

When you deploy an AzureContainerRegistry resource, Planton provisions:

- **Container Registry** — a `containerservice.Registry` in the specified region and resource group, with the configured SKU, admin-account setting, network posture (public access, IP rule set, bypass option, data endpoints), policies (quarantine, retention, content trust, export), managed identity, and CMK encryption
- **Geo-Replications** — for Premium registries, one replication per entry in `georeplications`, each its own tracked ARM resource with its own zone redundancy, regional endpoint, and tags
- **Azure Tags** — Planton-derived metadata tags merged with user tags (user wins on key collision), applied to the registry for tracking and governance

Nothing else is created here. Image-pull and image-push permissions are composed with the standalone `AzureRoleAssignment` resource against the registry's ARM ID, private endpoints for a fully private registry are composed with `AzurePrivateEndpoint`, and the identity that unwraps a CMK key is a referenced `AzureUserAssignedIdentity`.

## Prerequisites

- **Azure credentials** configured via environment variables or Planton provider config
- **An Azure Resource Group** where the registry will be created (can reference an AzureResourceGroup resource)
- **A globally unique registry name** — 5-50 lowercase alphanumerics, unique across all of Azure because it becomes the registry's DNS name (`{name}.azurecr.io`)
- **For CMK encryption**: an `AzureUserAssignedIdentity` holding get/wrapKey/unwrapKey permission on the Key Vault key's vault

## Quick Start

Create a file `acr.yaml`:

```yaml
apiVersion: azure.planton.dev/v1
kind: AzureContainerRegistry
metadata:
  name: my-registry
  annotations:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: dev.AzureContainerRegistry.my-registry
spec:
  region: eastus
  resourceGroup:
    value: my-rg
  registryName: myregistry01
```

Deploy:

```shell
planton apply -f acr.yaml
```

This creates a Standard-tier registry (the baseline applied when the SKU is unspecified) with the admin account disabled and public network access enabled — the production single-region shape. Images are pushed to and pulled from `myregistry01.azurecr.io` (the `login_server` output).

## Configuration Reference

### Required Fields

| Field | Type | Description | Validation |
|-------|------|-------------|------------|
| `region` | `string` | Azure region the registry's home replica lives in (e.g., `eastus`, `westeurope`). Additional regions are geo-replications, not a region change. Changing the region replaces the registry. | Required, minimum length 1 |
| `resourceGroup` | `StringValueOrRef` | Azure Resource Group name. Can reference an AzureResourceGroup resource via `valueFrom`. | Required |
| `registryName` | `string` | Globally unique registry name; becomes the `{name}.azurecr.io` login server. Changing it replaces the registry and its contents do not migrate — every image would need re-pushing. | Required, `^[a-z0-9]{5,50}$` |

### Optional Fields

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `sku` | `enum` | `STANDARD` | The registry's feature gate. `BASIC` (cost-optimized dev/test), `STANDARD` (production baseline: 100 GiB storage, webhooks, anonymous pull available), `PREMIUM` (geo-replication, zone redundancy, network isolation, policies, CMK, 500 GiB, highest throughput). Changes in place, but downgrading requires Premium-only features to be unset first. |
| `adminUserEnabled` | `bool` | `false` | Enables the built-in admin account (one username, two rotatable passwords, surfaced in outputs). Production authentication should be Entra ID; enable only for consumers that can use nothing but a static username/password. |
| `publicNetworkAccessEnabled` | `bool` | `true` | `false` takes the registry fully private (private endpoints only) and requires `PREMIUM`. For a public registry restricted to known addresses, keep `true` and use `networkRuleSet` instead. |
| `zoneRedundancyEnabled` | `bool` | `false` | Spreads the home replica across availability zones. `PREMIUM` only, fixed at creation. Each geo-replication declares its own zone redundancy separately. |
| `anonymousPullEnabled` | `bool` | `false` | Allows unauthenticated pulls — makes every repository publicly readable. `STANDARD` or `PREMIUM` only. The right shape for public artifact distribution, and only that. |
| `dataEndpointEnabled` | `bool` | `false` | Gives the registry dedicated regional data endpoints (`{name}.{region}.data.azurecr.io`) instead of shared storage endpoints, so egress firewalls can allowlist exact hostnames. `PREMIUM` only. |
| `quarantinePolicyEnabled` | `bool` | `false` | Quarantines newly pushed images until scanning tooling marks them passed; unquarantined clients cannot pull them. `PREMIUM` only. |
| `retentionPolicyInDays` | `int32` | unset | Days after which untagged manifests are purged (0 = immediately; unset keeps them forever, Azure's default). `PREMIUM` only. Range: 0-365. The hygiene lever that keeps CI push churn from growing storage without bound. |
| `trustPolicyEnabled` | `bool` | `false` | Enables Docker Content Trust: clients with content trust enabled can push signed images and verify signatures at pull. `PREMIUM` only. |
| `exportPolicyEnabled` | `bool` | `true` | `false` blocks exporting artifacts out of the registry — a data-exfiltration control. Requires `PREMIUM` and `publicNetworkAccessEnabled` explicitly `false` (validation enforces the pairing, as ARM does). |
| `networkRuleBypassOption` | `enum` | `AZURE_SERVICES` | Whether trusted Azure services (ACR Tasks, Microsoft Defender) may reach a network-restricted registry. `NONE` closes even that door. |
| `networkRuleSet` | `object` | unset | Access rules for a public registry: `defaultAction` (`ALLOW`/`DENY`) plus `ipRules` (IPv4 CIDR allowlist). `PREMIUM` only. A real allowlist sets `DENY`; ARM only supports allow rules, so entries carry no per-rule action. |
| `georeplications` | `object[]` | `[]` | Additional regions to replicate to, each with `location`, `zoneRedundancyEnabled`, `regionalEndpointEnabled`, and `tags`. `PREMIUM` only; must not contain the registry's own region (the home replica is implicit). Replicas add and remove in place. |
| `identity` | `object` | unset | The registry's managed identity: `type` (`SYSTEM_ASSIGNED`, `USER_ASSIGNED`, `SYSTEM_AND_USER_ASSIGNED`) and, for user-assigned flavors, `identityIds` referencing `AzureUserAssignedIdentity` resources. Required for CMK encryption. |
| `encryption` | `object` | unset | Customer-managed-key encryption: `identityClientId` (defaults to referencing an `AzureUserAssignedIdentity`'s `client_id` output) and `keyVaultKeyId` (the key's full Key Vault ID). `PREMIUM` only, fixed at creation, and requires a `USER_ASSIGNED` identity. |
| `tags` | `map<string, string>` | `{}` | Additional tags applied to the registry, merged over Planton-derived tags (user wins on collision). |

## Examples

### Standard Production Registry

The single-region production baseline — Standard tier, admin account off, Entra-based authentication:

```yaml
apiVersion: azure.planton.dev/v1
kind: AzureContainerRegistry
metadata:
  name: prod-registry
  annotations:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: prod.AzureContainerRegistry.prod-registry
spec:
  region: eastus
  resourceGroup:
    value: prod-rg
  registryName: prodregistry01
  sku: STANDARD
```

### Premium Geo-Replicated Registry

A Premium registry with a zone-redundant home replica, replicas in two more regions, and automatic purging of untagged manifests. Pushes propagate to all replicas automatically; each replica serves its region's pulls locally:

```yaml
apiVersion: azure.planton.dev/v1
kind: AzureContainerRegistry
metadata:
  name: global-registry
  annotations:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: prod.AzureContainerRegistry.global-registry
spec:
  region: eastus
  resourceGroup:
    value: prod-rg
  registryName: globalregistry01
  sku: PREMIUM
  zoneRedundancyEnabled: true
  retentionPolicyInDays: 30
  georeplications:
    - location: westeurope
      zoneRedundancyEnabled: true
    - location: southeastasia
      zoneRedundancyEnabled: true
```

### Network-Restricted Public Registry

A Premium registry that stays publicly addressable but denies every connection not on the CIDR allowlist, with dedicated data endpoints so egress firewalls can allowlist exact hostnames (surfaced in the `data_endpoint_host_names` output):

```yaml
apiVersion: azure.planton.dev/v1
kind: AzureContainerRegistry
metadata:
  name: restricted-registry
  annotations:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: prod.AzureContainerRegistry.restricted-registry
spec:
  region: eastus
  resourceGroup:
    value: prod-rg
  registryName: restrictedreg01
  sku: PREMIUM
  dataEndpointEnabled: true
  networkRuleSet:
    defaultAction: DENY
    ipRules:
      - ipRange: 203.0.113.0/24
      - ipRange: 198.51.100.32/28
```

### Fully Private Registry with CMK Encryption

A locked-down registry: public access off (private endpoints only), export disabled as a data-exfiltration control, and data encrypted with a Key Vault key unwrapped by a referenced user-assigned identity. The identity must hold get/wrapKey/unwrapKey on the vault before the registry is created:

```yaml
apiVersion: azure.planton.dev/v1
kind: AzureContainerRegistry
metadata:
  name: locked-registry
  annotations:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: prod.AzureContainerRegistry.locked-registry
spec:
  region: eastus
  resourceGroup:
    value: prod-rg
  registryName: lockedregistry01
  sku: PREMIUM
  publicNetworkAccessEnabled: false
  exportPolicyEnabled: false
  identity:
    type: USER_ASSIGNED
    identityIds:
      - valueFrom:
          name: acr-cmk-identity
  encryption:
    identityClientId:
      valueFrom:
        name: acr-cmk-identity
    keyVaultKeyId:
      value: https://prod-vault.vault.azure.net/keys/acr-cmk
```

### Composing Pull Access for a CI Identity

Image-pull and image-push permissions are never bundled into the registry — they are standalone `AzureRoleAssignment` resources scoped to the registry's ARM ID, so grants stay visible, auditable, and independently lifecycled:

```yaml
apiVersion: azure.planton.dev/v1
kind: AzureRoleAssignment
metadata:
  name: ci-acr-pull
spec:
  scope:
    valueFrom:
      kind: AzureContainerRegistry
      name: prod-registry
      fieldPath: status.outputs.container_registry_id
  roleDefinitionName: AcrPull
  principalId:
    valueFrom:
      name: ci-runner-identity
```

## Stack Outputs

After deployment, the following outputs are available in `status.outputs`:

| Output | Type | Description |
|--------|------|-------------|
| `container_registry_id` | `string` | Azure Resource Manager ID of the registry — the primary output; AKS clusters reference it to wire image pulls, and `AcrPull`/`AcrPush` role assignments scope to it |
| `container_registry_name` | `string` | The registry's name as deployed |
| `login_server` | `string` | The hostname images are tagged with and pulled from, e.g. `myregistry.azurecr.io` |
| `admin_username` | `string` | The admin account's username (the registry name); populated only when `adminUserEnabled` is true |
| `admin_password` | `string` | One of the admin account's two rotatable passwords; populated only when `adminUserEnabled` is true |
| `system_assigned_identity_principal_id` | `string` | Principal (object) ID of the registry's system-assigned identity; populated only when the identity type includes `SYSTEM_ASSIGNED` |
| `data_endpoint_host_names` | `string[]` | The dedicated regional data-endpoint hostnames (home region plus each geo-replication); populated only when `dataEndpointEnabled` is true — the exact hostnames an egress firewall must allowlist |

## Related Components

- [AzureResourceGroup](/docs/catalog/azure/azureresourcegroup) -- provides the resource group for registry placement
- [AzureUserAssignedIdentity](/docs/catalog/azure/azureuserassignedidentity) -- the identity that unwraps a customer-managed encryption key, referenced via `identity.identityIds` and `encryption.identityClientId`
- [AzureRoleAssignment](/docs/catalog/azure/azureroleassignment) -- composes `AcrPull`/`AcrPush` grants against the registry's `container_registry_id`
- [AzureAksCluster](/docs/catalog/azure/azureakscluster) -- AKS clusters reference the registry via its `container_registry_id` output to wire image pulls
- [AzurePrivateEndpoint](/docs/catalog/azure/azureprivateendpoint) -- private connectivity for registries with `publicNetworkAccessEnabled: false`
- [AzureKeyVault](/docs/catalog/azure/azurekeyvault) -- holds the customer-managed encryption key referenced by `encryption.keyVaultKeyId`
