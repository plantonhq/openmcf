# AzureKeyVault

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `azure.planton.dev/v1alpha1`

**AzureKeyVaultSpec** defines the configuration for creating an Azure Key
Vault: the tenant-scoped container where an organization's encryption keys,
TLS certificates, and application secrets live behind one security
boundary.

The vault is where governance is set -- authorization mode (Azure RBAC vs
legacy access policies), network isolation (public access, IP allowlists,
VNet service endpoints), deletion safety (soft delete and purge
protection), and the pricing tier that gates HSM-backed keys. Azure orgs
deliberately run FEW vaults with MANY objects inside, because every one of
those controls applies vault-wide.

What lives inside the vault is composed, never bundled here: encryption
keys are first-class AzureKeyVaultKey resources and TLS certificates are
first-class AzureKeyVaultCertificate resources, each referencing this
vault's key_vault_id output. Secret VALUES are deliberately out of scope
for infrastructure-as-code -- provision the vault here and manage secret
content through a secrets-management workflow, so plaintext never enters
deployment manifests or state.

The vault is created in the deploying credential's Azure AD tenant (a
vault cannot be managed cross-tenant, so modeling the tenant would only
invite a contradiction); the tenant_id output reports it.

## Example

```yaml
apiVersion: azure.planton.dev/v1alpha1
kind: AzureKeyVault
metadata:
  name: test-kv
  org: test-org
  env: dev
spec:
  region: eastus
  resource_group:
    value: test-rg
  vault_name: planton-hack-kv
  sku: PREMIUM
  rbac_authorization_enabled: false
  access_policies:
    - object_id:
        value: 00000000-0000-0000-0000-000000000001
      key_permissions:
        - KEY_GET
        - KEY_UNWRAP_KEY
        - KEY_WRAP_KEY
      secret_permissions:
        - SECRET_GET
      certificate_permissions:
        - CERTIFICATE_GET
      storage_permissions:
        - STORAGE_GET
  enabled_for_deployment: true
  enabled_for_disk_encryption: true
  purge_protection_enabled: true
  soft_delete_retention_days: 7
  network_acls:
    default_action: DENY
    bypass: NONE
    ip_rules:
      - "203.0.113.0/24"
    virtual_network_subnet_ids:
      - value: /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/test-rg/providers/Microsoft.Network/virtualNetworks/test-vnet/subnets/test-subnet
  tags:
    team: security
```

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.region` | `string` | yes |  |  |
| `spec.resourceGroup` | `string \| valueFrom` | yes |  | AzureResourceGroup (`status.outputs.resource_group_name`) |
| `spec.vaultName` | `string` | yes |  |  |
| `spec.sku` | `enum` |  |  |  |
| `spec.rbacAuthorizationEnabled` | `bool` |  | `true` |  |
| `spec.accessPolicies` | `[]AzureKeyVaultAccessPolicy` |  |  |  |
| `spec.accessPolicies[].objectId` | `string \| valueFrom` | yes |  | AzureUserAssignedIdentity (`status.outputs.principal_id`) |
| `spec.accessPolicies[].tenantId` | `string` |  |  |  |
| `spec.accessPolicies[].applicationId` | `string` |  |  |  |
| `spec.accessPolicies[].keyPermissions` | `[]enum` |  |  |  |
| `spec.accessPolicies[].secretPermissions` | `[]enum` |  |  |  |
| `spec.accessPolicies[].certificatePermissions` | `[]enum` |  |  |  |
| `spec.accessPolicies[].storagePermissions` | `[]enum` |  |  |  |
| `spec.enabledForDeployment` | `bool` |  |  |  |
| `spec.enabledForDiskEncryption` | `bool` |  |  |  |
| `spec.enabledForTemplateDeployment` | `bool` |  |  |  |
| `spec.publicNetworkAccessEnabled` | `bool` |  | `true` |  |
| `spec.purgeProtectionEnabled` | `bool` |  |  |  |
| `spec.softDeleteRetentionDays` | `int32` |  | `90` |  |
| `spec.networkAcls` | `AzureKeyVaultNetworkAcls` |  |  |  |
| `spec.networkAcls.defaultAction` | `enum` | yes |  |  |
| `spec.networkAcls.bypass` | `enum` |  |  |  |
| `spec.networkAcls.ipRules` | `[]string` |  |  |  |
| `spec.networkAcls.virtualNetworkSubnetIds` | `[]string \| valueFrom` |  |  | AzureSubnet (`status.outputs.subnet_id`) |
| `spec.tags` | `map<string, string>` |  |  |  |

## Field Details

### spec.region

`string` · required

The Azure region where the Key Vault will be deployed (e.g. "eastus",
"westeurope"). Key Vault is a regional service; changing the region
replaces the vault.

- rule: {"required":true,"string":{"minLen":"1"}}

### spec.resourceGroup

`string | valueFrom` · required

The Azure resource group the vault will be created in. Can be a literal
resource-group name or a reference to an AzureResourceGroup's name
output.

- references: AzureResourceGroup (`status.outputs.resource_group_name`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: AzureResourceGroup, name: <that resource's name>, fieldPath: status.outputs.resource_group_name}} -- a bare string does not parse

### spec.vaultName

`string` · required

The name of the vault: 3-24 characters of letters, digits, and hyphens,
starting with a letter, ending with a letter or digit, no consecutive
hyphens -- and GLOBALLY unique across all of Azure, because it becomes
the vault's DNS name ({name}.vault.azure.net, the vault_uri output).
Changing the name replaces the vault. A deleted vault's name stays
reserved for the soft-delete retention window unless the vault is
purged.

- rule: vault_name cannot contain consecutive hyphens ("--")
- rule: {"required":true,"string":{"pattern":"^[a-zA-Z][a-zA-Z0-9-]{1,22}[a-zA-Z0-9]$"}}

### spec.sku

`enum`

The pricing tier. Unspecified applies STANDARD: software-protected
keys, the right tier for most workloads. PREMIUM adds HSM-backed keys
(FIPS 140-2 Level 3) -- required before any AzureKeyVaultKey in this
vault can use the EC_HSM/RSA_HSM key types, and for compliance regimes
that mandate hardware protection (PCI-DSS, FedRAMP High). The SKU can
be changed in place.

Allowed values (use exactly as shown):

- `azure_key_vault_sku_unspecified` -- Not specified: STANDARD.
- `STANDARD` -- Software-protected keys and secrets -- the right tier for most workloads.
- `PREMIUM` -- Adds HSM-backed keys (FIPS 140-2 Level 3) -- required for the EC_HSM/RSA_HSM key types and hardware-protection compliance regimes.

### spec.rbacAuthorizationEnabled

`bool` · optional (explicit presence)

Whether data-plane authorization uses Azure RBAC (true) or legacy
vault access policies (false). Azure's own guidance is RBAC: grants
become ordinary role assignments (compose AzureRoleAssignment with
roles like "Key Vault Administrator" or "Key Vault Secrets User"),
participate in PIM and access reviews, and support fine-grained scopes.
ARM's default for a new vault is the legacy access-policy mode, but
this spec defaults to RBAC as the recommended posture -- set false
explicitly (and populate access_policies) to run the legacy mode.
Switching modes on a live vault requires Microsoft.Authorization write
permission (Owner / User Access Administrator) on the vault.

- default: `true`

### spec.accessPolicies

`[]AzureKeyVaultAccessPolicy`

Legacy access-policy grants (up to 1024). Only honored when
rbac_authorization_enabled is false -- ARM stores but IGNORES access
policies on an RBAC-mode vault, so populating both is almost always a
mistake. Each entry grants one Azure AD principal (user, group,
service principal, or managed identity) explicit permission lists over
keys, secrets, certificates, and managed-storage objects in this
vault. Prefer RBAC for new vaults; model policies only when an
existing workload or org standard requires the legacy mode.

- rule: {"repeated":{"maxItems":"1024"}}

### spec.accessPolicies[].objectId

`string | valueFrom` · required

The object ID of the Azure AD principal being granted access (a user,
group, service principal, or managed identity). Defaults to
referencing an AzureUserAssignedIdentity's principal_id output in
composed environments -- note this is the identity's PRINCIPAL id
(the directory object), not its client id.

- references: AzureUserAssignedIdentity (`status.outputs.principal_id`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: AzureUserAssignedIdentity, name: <that resource's name>, fieldPath: status.outputs.principal_id}} -- a bare string does not parse

### spec.accessPolicies[].tenantId

`string` · optional (explicit presence)

The Azure AD tenant the principal lives in. Leave unset to use the
vault's own tenant (the deploying credential's tenant) -- the correct
value for virtually every grant, since access policies cannot span
tenants in practice.

- rule: {"string":{"uuid":true}}

### spec.accessPolicies[].applicationId

`string` · optional (explicit presence)

For grants to a specific application acting on behalf of the
principal (the rarely-used compound-identity flow): that
application's client ID. Leave unset for ordinary grants.

- rule: {"string":{"uuid":true}}

### spec.accessPolicies[].keyPermissions

`[]enum`

Key permissions granted to the principal. KEY_GET + cryptographic
operations (KEY_ENCRYPT/KEY_DECRYPT, KEY_WRAP_KEY/KEY_UNWRAP_KEY,
KEY_SIGN/KEY_VERIFY) cover consumers; management permissions
(KEY_CREATE, KEY_DELETE, rotation-policy permissions) belong to
operators.

- rule: {"repeated":{"items":{"enum":{"definedOnly":true,"notIn":[0]}}}}

Allowed values (use exactly as shown):

- `azure_key_vault_key_permission_unspecified` -- Not specified -- invalid; grant an explicit permission.
- `KEY_GET` -- Read key metadata and the public part.
- `KEY_LIST` -- List keys in the vault.
- `KEY_UPDATE` -- Update a key's attributes.
- `KEY_CREATE` -- Create new keys.
- `KEY_IMPORT` -- Import externally generated key material.
- `KEY_DELETE` -- Delete keys (soft delete).
- `KEY_RECOVER` -- Recover soft-deleted keys.
- `KEY_BACKUP` -- Download a protected key backup.
- `KEY_RESTORE` -- Restore a key from backup.
- `KEY_DECRYPT` -- Decrypt with the key.
- `KEY_ENCRYPT` -- Encrypt with the key.
- `KEY_UNWRAP_KEY` -- Unwrap (decrypt) a data-encryption key -- the permission customer-managed-key consumers need.
- `KEY_WRAP_KEY` -- Wrap (encrypt) a data-encryption key.
- `KEY_VERIFY` -- Verify a signature with the public key.
- `KEY_SIGN` -- Sign with the private key.
- `KEY_PURGE` -- Permanently delete (purge) a soft-deleted key.
- `KEY_RELEASE` -- Release the key to a trusted execution environment (secure key release).
- `KEY_ROTATE` -- Rotate the key to a new version on demand.
- `KEY_GET_ROTATION_POLICY` -- Read the key's rotation policy.
- `KEY_SET_ROTATION_POLICY` -- Set the key's rotation policy.

### spec.accessPolicies[].secretPermissions

`[]enum`

Secret permissions granted to the principal. SECRET_GET is the
consumer permission; SECRET_SET/SECRET_DELETE are the writer surface.

- rule: {"repeated":{"items":{"enum":{"definedOnly":true,"notIn":[0]}}}}

Allowed values (use exactly as shown):

- `azure_key_vault_secret_permission_unspecified` -- Not specified -- invalid; grant an explicit permission.
- `SECRET_GET` -- Read secret values -- the consumer permission.
- `SECRET_LIST` -- List secrets in the vault (names only, not values).
- `SECRET_SET` -- Create or update secret values.
- `SECRET_DELETE` -- Delete secrets (soft delete).
- `SECRET_RECOVER` -- Recover soft-deleted secrets.
- `SECRET_BACKUP` -- Download a protected secret backup.
- `SECRET_RESTORE` -- Restore a secret from backup.
- `SECRET_PURGE` -- Permanently delete (purge) a soft-deleted secret.

### spec.accessPolicies[].certificatePermissions

`[]enum`

Certificate permissions granted to the principal. CERTIFICATE_GET
covers consumers; CERTIFICATE_CREATE/CERTIFICATE_IMPORT plus the
issuer-management permissions belong to certificate operators.

- rule: {"repeated":{"items":{"enum":{"definedOnly":true,"notIn":[0]}}}}

Allowed values (use exactly as shown):

- `azure_key_vault_certificate_permission_unspecified` -- Not specified -- invalid; grant an explicit permission.
- `CERTIFICATE_GET` -- Read certificates (public part and policy) -- the consumer permission.
- `CERTIFICATE_LIST` -- List certificates in the vault.
- `CERTIFICATE_UPDATE` -- Update a certificate's attributes or policy.
- `CERTIFICATE_CREATE` -- Create (enroll) new certificates.
- `CERTIFICATE_IMPORT` -- Import existing certificate bundles.
- `CERTIFICATE_DELETE` -- Delete certificates (soft delete).
- `CERTIFICATE_RECOVER` -- Recover soft-deleted certificates.
- `CERTIFICATE_BACKUP` -- Download a protected certificate backup.
- `CERTIFICATE_RESTORE` -- Restore a certificate from backup.
- `CERTIFICATE_MANAGE_CONTACTS` -- Manage the vault's certificate contacts (expiry notifications).
- `CERTIFICATE_MANAGE_ISSUERS` -- Manage certificate-authority issuer configurations.
- `CERTIFICATE_GET_ISSUERS` -- Read a specific issuer configuration.
- `CERTIFICATE_LIST_ISSUERS` -- List issuer configurations.
- `CERTIFICATE_SET_ISSUERS` -- Create or update issuer configurations.
- `CERTIFICATE_DELETE_ISSUERS` -- Delete issuer configurations.
- `CERTIFICATE_PURGE` -- Permanently delete (purge) a soft-deleted certificate.

### spec.accessPolicies[].storagePermissions

`[]enum`

Managed-storage-account permissions granted to the principal (the
legacy Key Vault storage-key-rotation feature; rarely used today).

- rule: {"repeated":{"items":{"enum":{"definedOnly":true,"notIn":[0]}}}}

Allowed values (use exactly as shown):

- `azure_key_vault_storage_permission_unspecified` -- Not specified -- invalid; grant an explicit permission.
- `STORAGE_GET` -- Read managed storage account definitions.
- `STORAGE_LIST` -- List managed storage accounts.
- `STORAGE_SET` -- Create or update managed storage account definitions.
- `STORAGE_UPDATE` -- Update managed storage account attributes.
- `STORAGE_DELETE` -- Delete managed storage account definitions (soft delete).
- `STORAGE_RECOVER` -- Recover soft-deleted managed storage accounts.
- `STORAGE_BACKUP` -- Download a protected managed-storage backup.
- `STORAGE_RESTORE` -- Restore a managed storage account from backup.
- `STORAGE_PURGE` -- Permanently delete (purge) a soft-deleted managed storage account.
- `STORAGE_REGENERATE_KEY` -- Regenerate the managed storage account's keys.
- `STORAGE_GET_SAS` -- Read a SAS definition.
- `STORAGE_LIST_SAS` -- List SAS definitions.
- `STORAGE_SET_SAS` -- Create or update SAS definitions.
- `STORAGE_DELETE_SAS` -- Delete SAS definitions.

### spec.enabledForDeployment

`bool`

Whether Azure Virtual Machines may retrieve certificates stored as
secrets from this vault (the VM deployment integration). Azure's
default is false.

### spec.enabledForDiskEncryption

`bool`

Whether Azure Disk Encryption may retrieve secrets and unwrap keys
from this vault (the legacy in-guest ADE flow; modern server-side
encryption uses disk encryption sets instead). Azure's default is
false.

### spec.enabledForTemplateDeployment

`bool`

Whether Azure Resource Manager (template deployments) may retrieve
secrets from this vault. Azure's default is false.

### spec.publicNetworkAccessEnabled

`bool` · optional (explicit presence)

Whether the vault accepts connections from the public internet.
Azure's default is true. Setting false takes the vault fully private
-- reachable only through private endpoints. For a public vault
restricted to known networks, keep this true and use network_acls
instead. Updatable in place.

- default: `true`

### spec.purgeProtectionEnabled

`bool`

Whether purge protection is on: a deleted vault (or key/certificate
inside it) can then only be recovered -- never permanently purged --
until the soft-delete retention window expires. Azure's default is
false. Turn it on for production vaults, and ALWAYS for vaults whose
keys encrypt other resources (customer-managed keys refuse to enroll
against a vault without it in many services). Irreversible: once
enabled it cannot be disabled, and destroying the vault then means
waiting out the retention window before the name frees up.

### spec.softDeleteRetentionDays

`int32` · optional (explicit presence)

How many days deleted vaults and vault objects remain recoverable
(soft delete). 7-90; unspecified applies Azure's default of 90. Can
only be set at creation -- changing it replaces the vault.

- default: `90`
- rule: {"int32":{"lte":90,"gte":7}}

### spec.networkAcls

`AzureKeyVaultNetworkAcls`

Network access rules for a PUBLIC vault: a default action, a bypass
carve-out for trusted Microsoft services, an IP allowlist, and a VNet
subnet allowlist. The middle ground between "open to the internet"
and "private endpoints only". Updatable in place.

### spec.networkAcls.defaultAction

`enum` · required

What happens to requests matching no ip_rules or subnet rule: ALLOW
leaves the vault open (the rule set is effectively off); DENY is the
allowlist posture. ARM requires a value whenever network rules are
configured.

- rule: {"required":true}

Allowed values (use exactly as shown):

- `azure_key_vault_network_acls_default_action_unspecified` -- Not specified -- invalid; ARM requires an explicit choice when network rules are configured.
- `ALLOW` -- Requests matching no rules are allowed (the rule set is effectively off).
- `DENY` -- Requests matching no rules are denied -- the allowlist posture.

### spec.networkAcls.bypass

`enum`

Whether trusted Microsoft services (Azure Backup, Disk Encryption,
Azure Monitor, ...) may reach the vault even when default_action is
DENY. Unspecified applies Azure's default (AZURE_SERVICES) -- the
pragmatic choice that keeps first-party integrations working. NONE
closes even that door.

Allowed values (use exactly as shown):

- `azure_key_vault_network_acls_bypass_unspecified` -- Not specified: Azure's default (AzureServices).
- `AZURE_SERVICES` -- Trusted Microsoft services may reach the vault despite network rules.
- `NONE` -- No bypass: network rules apply to everything, first-party services included.

### spec.networkAcls.ipRules

`[]string`

Public IPv4 addresses or CIDR ranges allowed to reach the vault
(office egress, VPN gateways, CI runners). E.g. "203.0.113.0/24" or
"198.51.100.42".

- rule: {"repeated":{"maxItems":"200"}}

### spec.networkAcls.virtualNetworkSubnetIds

`[]string | valueFrom`

Virtual-network subnets allowed to reach the vault over service
endpoints, by subnet ARM ID. The referenced subnets must have the
"Microsoft.KeyVault" service endpoint enabled. Defaults to
referencing an AzureSubnet's subnet_id output in composed
environments.

- references: AzureSubnet (`status.outputs.subnet_id`)
- rule: {"repeated":{"maxItems":"100"}}
- rule: write as {value: <literal>} or {valueFrom: {kind: AzureSubnet, name: <that resource's name>, fieldPath: status.outputs.subnet_id}} -- a bare string does not parse

### spec.tags

`map<string, string>`

Free-form tags applied to the vault, merged over the Planton-derived
resource tags (organization, environment, resource id); a user tag
with the same key wins. Tags are Azure's governance surface -- Azure
Policy enforces them and Microsoft Cost Management groups by them.
Updatable in place.

## Outputs

Reference an output from another manifest as `valueFrom: {kind: AzureKeyVault, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.key_vault_id` | `string` | The vault's ARM resource ID: /subscriptions/{subscription}/resourceGroups/{rg}/providers/Microsoft.KeyVault/vaults/{name} The reference every child object (AzureKeyVaultKey, AzureKeyVaultCertificate) and every vault-scoped grant (AzureRoleAssignment on an RBAC vault) targets. |
| `status.outputs.key_vault_name` | `string` | The vault's name. |
| `status.outputs.vault_uri` | `string` | The vault's data-plane URI: https://{name}.vault.azure.net/. Applications and SDKs address keys, secrets, and certificates through this endpoint. |
| `status.outputs.tenant_id` | `string` | The Azure AD tenant the vault authenticates against (the deploying credential's tenant). |
| `status.outputs.resource_group_name` | `string` | The resource group the vault was created in. |

## References

Fields that can point at another resource's outputs:

| Field | Kind | Output |
|---|---|---|
| `spec.resourceGroup` | AzureResourceGroup | `status.outputs.resource_group_name` |
| `spec.accessPolicies[].objectId` | AzureUserAssignedIdentity | `status.outputs.principal_id` |
| `spec.networkAcls.virtualNetworkSubnetIds` | AzureSubnet | `status.outputs.subnet_id` |

## Referenced By

Fields on other kinds that can point at this resource:

| Kind | Field | Reads |
|---|---|---|
| AzureAksCluster | `spec.serviceMeshProfile.certificateAuthority.keyVaultId` | `status.outputs.key_vault_id` |
| AzureKeyVaultCertificate | `spec.keyVaultId` | `status.outputs.key_vault_id` |
| AzureKeyVaultKey | `spec.keyVaultId` | `status.outputs.key_vault_id` |
| AzureVirtualMachine | `spec.secrets[].keyVaultId` | `status.outputs.key_vault_id` |
| AzureVirtualMachineScaleSet | `spec.extensions[].protectedSettingsFromKeyVault.sourceVaultId` | `status.outputs.key_vault_id` |
| AzureVirtualMachineScaleSet | `spec.secrets[].keyVaultId` | `status.outputs.key_vault_id` |
| KubernetesClusterSecretStore | `spec.config.azureKeyVault.vaultUrl` | `status.outputs.vault_uri` |
| KubernetesSecretStore | `spec.config.azureKeyVault.vaultUrl` | `status.outputs.vault_uri` |

## See Also

- [Overview](../README.md)
