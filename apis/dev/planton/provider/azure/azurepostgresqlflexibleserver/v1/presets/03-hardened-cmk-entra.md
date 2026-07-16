# Hardened Server: Entra-Only Auth + Customer-Managed Key

This preset creates the compliance posture: PostgreSQL password
authentication is disabled entirely (no admin password exists to leak,
rotate, or audit), administration flows through a Microsoft Entra group,
and the server's data is encrypted with a Key Vault key you own instead of
a Microsoft-managed key.

## When to Use

- Regulated workloads whose compliance regime requires customer-managed
  encryption keys (CMK) and centralized identity
- Organizations eliminating static database credentials in favor of Entra
  tokens

## Key Configuration Choices

- **`password_auth_enabled: false`** -- the admin login/password pair must
  be OMITTED (Azure rejects it); clients connect with Entra access tokens
  and the `aad_administrators` group manages roles inside PostgreSQL
- **CMK through a user-assigned identity** -- the identity must hold
  wrap/unwrap access on the key's vault BEFORE the server is created
  (compose an `AzureRoleAssignment` granting "Key Vault Crypto Service
  Encryption User" at the vault scope); the key's vault must have purge
  protection enabled
- **Versionless key reference** -- key rotations propagate automatically;
  pin a versioned ID only when the compliance regime demands an immutable
  key version
- **Audit-friendly parameters** -- connection logging is pinned on

## Prerequisite Wiring

The CMK identity's vault grant, composed so it exists before the server:

```yaml
apiVersion: azure.planton.dev/v1
kind: AzureRoleAssignment
metadata:
  name: postgres-cmk-unwrap
spec:
  scope:
    valueFrom:
      kind: AzureKeyVault
      name: platform-vault
      fieldPath: status.outputs.key_vault_id
  roleDefinitionName: Key Vault Crypto Service Encryption User
  principalId:
    valueFrom:
      name: <cmk-identity-resource-name>
```

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<resource-group-name>` | The resource group to create the server in | The resource group's `status.outputs.resource_group_name` |
| `myorg-hardened-postgres` | 3-63 lowercase chars, globally unique | Your naming convention |
| `<entra-group-object-id>` | The admin group's directory object ID | Entra ID -> Groups -> Overview |
| `<entra-group-name>` | The group's display name (the PostgreSQL role name) | Entra ID -> Groups |
| `<cmk-identity-resource-name>` | The AzureUserAssignedIdentity resource | Your identity manifests |
| `<cmk-key-resource-name>` | The AzureKeyVaultKey resource | Your key manifests |
| `<cost-center>` | Your org's cost-attribution tag value | Your tagging convention |
