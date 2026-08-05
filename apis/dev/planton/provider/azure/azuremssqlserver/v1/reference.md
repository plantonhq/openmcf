# AzureMssqlServer

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `azure.planton.dev/v1`

**AzureMssqlServerSpec** defines the configuration for creating an Azure
SQL Database logical server: the administrative container that carries
the login endpoint, authentication (SQL and/or Microsoft Entra),
networking posture, transparent-data-encryption key, auditing, and
Microsoft Defender settings for every database and elastic pool created
on it.

**Architecture**: unlike PostgreSQL and MySQL Flexible Servers (physical
servers with compute and storage), Azure SQL uses a **logical server +
database** model. The server has no compute of its own -- each
AzureMssqlDatabase carries its own SKU and storage and each
AzureMssqlElasticPool carries shared compute, so databases and pools are
first-class kinds referencing this server's `server_id` output rather
than fields folded into this spec.

**Network access**: Azure SQL does not support VNet delegation. Public
reachability is `public_network_access_enabled` with `firewall_rules`
(IPv4 allowlist) and `virtual_network_rules` (subnet service-endpoint
allowlist); private connectivity is exclusively AzurePrivateEndpoint.
Outbound restriction (`outbound_network_restriction_enabled` +
`outbound_firewall_rules`) governs where the server may reach OUT to,
e.g. for elastic queries.

**Authentication** is SQL logins, Microsoft Entra, or both -- matching
ARM's real contract: at least one of `administrator_login` or an
`azuread_administrator` with `azuread_authentication_only` must be
configured. With Entra-only auth the SQL credentials are omitted
entirely.

**Fixed at creation** (changing replaces the server): `server_name`,
`region`, `resource_group`, `version`, and `administrator_login` once
set. Two update-time contracts ARM enforces that validation cannot
express: `minimum_tls_version` cannot be removed once set, and the
administrator password cannot change while
`azuread_authentication_only` is true.

## Example

```yaml
apiVersion: azure.planton.dev/v1
kind: AzureMssqlServer
metadata:
  name: test-mssql
spec:
  region: eastus
  resourceGroup:
    value: test-rg
  serverName: test-mssql-server
  version: "12.0"
  administratorLogin: sqladmin
  administratorPassword:
    value: P@ssw0rd1234!
  # Exercises the Entra administrator block alongside SQL auth (mixed
  # mode) plus the tenant passthrough.
  azureadAdministrator:
    loginUsername: dba-team@contoso.com
    objectId:
      value: 11111111-2222-3333-4444-555555555555
    azureadAuthenticationOnly: false
  # Exercises the identity-type enum mapping, the primary-identity pin,
  # and the TDE customer-managed key.
  identity:
    type: SYSTEM_AND_USER_ASSIGNED
    identityIds:
      - value: /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/test-rg/providers/Microsoft.ManagedIdentity/userAssignedIdentities/sql-uai
  primaryUserAssignedIdentityId:
    value: /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/test-rg/providers/Microsoft.ManagedIdentity/userAssignedIdentities/sql-uai
  transparentDataEncryptionKeyVaultKeyId:
    value: https://test-vault.vault.azure.net/keys/sql-tde/0123456789abcdef0123456789abcdef
  # Exercises the connection-policy enum mapping.
  connectionPolicy: REDIRECT
  minimumTlsVersion: "1.2"
  publicNetworkAccessEnabled: true
  # Exercises the outbound restriction + FQDN rule sub-resources.
  outboundNetworkRestrictionEnabled: true
  outboundFirewallRules:
    - peer.database.windows.net
  expressVulnerabilityAssessmentEnabled: true
  firewallRules:
    - name: allow-azure-services
      startIpAddress: "0.0.0.0"
      endIpAddress: "0.0.0.0"
    - name: allow-office
      startIpAddress: "203.0.113.0"
      endIpAddress: "203.0.113.255"
  # Exercises the subnet service-endpoint allowlist sub-resource.
  virtualNetworkRules:
    - name: allow-app-subnet
      subnetId:
        value: /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/test-rg/providers/Microsoft.Network/virtualNetworks/test-vnet/subnets/app
  # Exercises the folded auditing singleton with blob export.
  extendedAuditing:
    storageEndpoint: https://testauditlogs.blob.core.windows.net
    storageAccountAccessKey:
      value: dGVzdC1rZXk=
    retentionInDays: 90
    logMonitoringEnabled: true
  # Exercises the Defender policy with the alert-type enum mapping.
  securityAlertPolicy:
    state: ENABLED
    disabledAlerts:
      - UNSAFE_ACTION
    emailAccountAdmins: true
    emailAddresses:
      - secops@contoso.com
    retentionDays: 30
  tags:
    team: data
```

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.region` | `string` | yes |  |  |
| `spec.resourceGroup` | `string \| valueFrom` | yes |  | AzureResourceGroup (`status.outputs.resource_group_name`) |
| `spec.serverName` | `string` | yes |  |  |
| `spec.version` | `string` |  | `12.0` |  |
| `spec.administratorLogin` | `string` |  |  |  |
| `spec.administratorPassword` | `string \| valueFrom` (sensitive) |  |  |  |
| `spec.azureadAdministrator` | `AzureMssqlServerAzureadAdministrator` |  |  |  |
| `spec.azureadAdministrator.loginUsername` | `string` | yes |  |  |
| `spec.azureadAdministrator.objectId` | `string \| valueFrom` | yes |  | AzureUserAssignedIdentity (`status.outputs.principal_id`) |
| `spec.azureadAdministrator.tenantId` | `string` |  |  |  |
| `spec.azureadAdministrator.azureadAuthenticationOnly` | `bool` |  |  |  |
| `spec.identity` | `AzureMssqlServerIdentity` |  |  |  |
| `spec.identity.type` | `enum` | yes |  |  |
| `spec.identity.identityIds` | `[]string \| valueFrom` |  |  | AzureUserAssignedIdentity (`status.outputs.identity_id`) |
| `spec.primaryUserAssignedIdentityId` | `string \| valueFrom` |  |  | AzureUserAssignedIdentity (`status.outputs.identity_id`) |
| `spec.transparentDataEncryptionKeyVaultKeyId` | `string \| valueFrom` |  |  | AzureKeyVaultKey (`status.outputs.key_id`) |
| `spec.connectionPolicy` | `enum` |  |  |  |
| `spec.minimumTlsVersion` | `string` |  | `1.2` |  |
| `spec.publicNetworkAccessEnabled` | `bool` |  | `true` |  |
| `spec.outboundNetworkRestrictionEnabled` | `bool` |  |  |  |
| `spec.outboundFirewallRules` | `[]string` |  |  |  |
| `spec.expressVulnerabilityAssessmentEnabled` | `bool` |  |  |  |
| `spec.firewallRules` | `[]AzureMssqlServerFirewallRule` |  |  |  |
| `spec.firewallRules[].name` | `string` | yes |  |  |
| `spec.firewallRules[].startIpAddress` | `string` | yes |  |  |
| `spec.firewallRules[].endIpAddress` | `string` | yes |  |  |
| `spec.virtualNetworkRules` | `[]AzureMssqlServerVirtualNetworkRule` |  |  |  |
| `spec.virtualNetworkRules[].name` | `string` | yes |  |  |
| `spec.virtualNetworkRules[].subnetId` | `string \| valueFrom` | yes |  | AzureSubnet (`status.outputs.subnet_id`) |
| `spec.virtualNetworkRules[].ignoreMissingVnetServiceEndpoint` | `bool` |  |  |  |
| `spec.extendedAuditing` | `AzureMssqlServerExtendedAuditing` |  |  |  |
| `spec.extendedAuditing.storageEndpoint` | `string` |  |  |  |
| `spec.extendedAuditing.storageAccountAccessKey` | `string \| valueFrom` (sensitive) |  |  |  |
| `spec.extendedAuditing.storageAccountAccessKeyIsSecondary` | `bool` |  |  |  |
| `spec.extendedAuditing.retentionInDays` | `int32` |  |  |  |
| `spec.extendedAuditing.logMonitoringEnabled` | `bool` |  | `true` |  |
| `spec.extendedAuditing.storageAccountSubscriptionId` | `string` |  |  |  |
| `spec.extendedAuditing.predicateExpression` | `string` |  |  |  |
| `spec.extendedAuditing.auditActionsAndGroups` | `[]string` |  |  |  |
| `spec.securityAlertPolicy` | `AzureMssqlServerSecurityAlertPolicy` |  |  |  |
| `spec.securityAlertPolicy.state` | `enum` | yes |  |  |
| `spec.securityAlertPolicy.disabledAlerts` | `[]enum` |  |  |  |
| `spec.securityAlertPolicy.emailAccountAdmins` | `bool` |  |  |  |
| `spec.securityAlertPolicy.emailAddresses` | `[]string` |  |  |  |
| `spec.securityAlertPolicy.retentionDays` | `int32` |  |  |  |
| `spec.securityAlertPolicy.storageEndpoint` | `string` |  |  |  |
| `spec.securityAlertPolicy.storageAccountAccessKey` | `string \| valueFrom` (sensitive) |  |  |  |
| `spec.tags` | `map<string, string>` |  |  |  |

## Field Details

### spec.region

`string` · required

The Azure region where the logical server is created (e.g. "eastus",
"westeurope"). Changing the region replaces the server.

- rule: {"required":true,"string":{"minLen":"1"}}

### spec.resourceGroup

`string | valueFrom` · required

The Azure resource group the server is created in. Can be a literal
resource-group name or a reference to an AzureResourceGroup's name
output. Changing it replaces the server.

- references: AzureResourceGroup (`status.outputs.resource_group_name`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: AzureResourceGroup, name: <that resource's name>, fieldPath: status.outputs.resource_group_name}} -- a bare string does not parse

### spec.serverName

`string` · required

The server's name: 1-63 lowercase letters, digits, and hyphens, not
starting or ending with a hyphen -- and GLOBALLY unique across Azure,
because it becomes the server's DNS name
({name}.database.windows.net, the fqdn output). Changing the name
replaces the server.

- rule: {"required":true,"string":{"minLen":"1","maxLen":"63","pattern":"^[0-9a-z]([-0-9a-z]*[0-9a-z])?$"}}

### spec.version

`string` · optional (explicit presence)

The SQL Server version identifier. Unspecified applies "12.0" -- the
only version real deployments use (it tracks the current Azure SQL
engine; "2.0" is the SQL Server 2005-era legacy identifier kept for
API completeness). Fixed at creation.

- default: `12.0`
- rule: version must be one of: 2.0, 12.0

### spec.administratorLogin

`string`

The administrator login for SQL authentication. Required unless the
server is Entra-only (azuread_administrator.azuread_authentication_only
= true). Azure reserves names like "admin", "administrator", "sa",
"root", "dbmanager", "loginmanager", "dbo", "guest", and "public".
The login is fixed once set.

- rule: administrator_login cannot be a reserved name (admin, administrator, sa, root, dbmanager, loginmanager, dbo, guest, public)

### spec.administratorPassword

`string | valueFrom` · sensitive

The administrator password for SQL authentication (8-128 characters
from at least three of: uppercase, lowercase, digits, special
characters). Can be a literal value or a reference to another
resource's output. Required with administrator_login; rotates in
place -- but ARM rejects a password change while
azuread_authentication_only is true.

- rule: write as {value: <literal>} or {valueFrom: {kind: <Kind>, name: <that resource's name>, fieldPath: status.outputs.<output>}} -- a bare string does not parse

### spec.azureadAdministrator

`AzureMssqlServerAzureadAdministrator`

The server's Microsoft Entra (Azure AD) administrator -- the
principal (use a group to admit a team) that can administer the
server with directory tokens. Enables Entra authentication; with
azuread_authentication_only it becomes the ONLY authentication
mechanism and the SQL credentials must be omitted.

### spec.azureadAdministrator.loginUsername

`string` · required

The administrator principal's display name as it appears in Entra
(e.g. "dba-team@contoso.com" for a user, the group name for a
group). Used as the login name for Entra connections.

- rule: {"required":true,"string":{"minLen":"1"}}

### spec.azureadAdministrator.objectId

`string | valueFrom` · required

The directory object ID of the Entra principal being granted the
administrator role. The default reference points at an
AzureUserAssignedIdentity's principal_id output so a deployment
identity can administer the server; for a user or group supply the
directory object ID literally.

- references: AzureUserAssignedIdentity (`status.outputs.principal_id`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: AzureUserAssignedIdentity, name: <that resource's name>, fieldPath: status.outputs.principal_id}} -- a bare string does not parse

### spec.azureadAdministrator.tenantId

`string` · optional (explicit presence)

The Entra tenant of the administrator principal. Leave unset to use
the deploying credential's tenant -- the correct value for virtually
every deployment.

- rule: {"string":{"uuid":true}}

### spec.azureadAdministrator.azureadAuthenticationOnly

`bool`

When true, Entra becomes the ONLY authentication mechanism: SQL
logins are disabled server-wide and the SQL admin credentials must
be omitted from the spec. Azure's recommended posture for
passwordless estates.

### spec.identity

`AzureMssqlServerIdentity`

The server's managed identity. Required for the
transparent-data-encryption customer-managed key (the identity
unwraps the key) and usable by databases for their own CMK.

- rule: identity_ids are required for USER_ASSIGNED / SYSTEM_AND_USER_ASSIGNED and must be empty for SYSTEM_ASSIGNED

### spec.identity.type

`enum` · required

The identity flavor. SYSTEM_ASSIGNED is created and rotated by
Azure with the server's lifecycle; USER_ASSIGNED attaches
pre-existing identities (required when Key Vault grants must exist
BEFORE the server does, e.g. for TDE CMK).

- rule: {"required":true}

Allowed values (use exactly as shown):

- `azure_mssql_server_identity_type_unspecified` -- Not specified -- invalid; choose an explicit type when the identity block is present.
- `SYSTEM_ASSIGNED` -- Azure creates and rotates an identity tied to the server's lifecycle.
- `USER_ASSIGNED` -- Attach pre-existing user-assigned identities (identity_ids).
- `SYSTEM_AND_USER_ASSIGNED` -- Both a system-assigned identity and the attached user-assigned identities.

### spec.identity.identityIds

`[]string | valueFrom`

The user-assigned identities to attach, by ARM ID. Required (and
only allowed) when the type includes USER_ASSIGNED.

- references: AzureUserAssignedIdentity (`status.outputs.identity_id`)
- rule: write as {value: <literal>} or {valueFrom: {kind: AzureUserAssignedIdentity, name: <that resource's name>, fieldPath: status.outputs.identity_id}} -- a bare string does not parse

### spec.primaryUserAssignedIdentityId

`string | valueFrom`

Which user-assigned identity (by ARM ID) is the server's PRIMARY
identity -- the one ARM uses for Key Vault access when several are
attached. Must be one of identity.identity_ids. Required when the
identity type includes USER_ASSIGNED.

- references: AzureUserAssignedIdentity (`status.outputs.identity_id`)
- rule: write as {value: <literal>} or {valueFrom: {kind: AzureUserAssignedIdentity, name: <that resource's name>, fieldPath: status.outputs.identity_id}} -- a bare string does not parse

### spec.transparentDataEncryptionKeyVaultKeyId

`string | valueFrom`

Server-level transparent-data-encryption customer-managed key: every
database's TDE protector defaults to this Key Vault key instead of a
service-managed key. Takes a VERSIONED key ID (ARM pins the exact
version at the server level), so the default reference points at an
AzureKeyVaultKey's key_id output. Requires an identity with
wrap/unwrap access on the key's vault.

- references: AzureKeyVaultKey (`status.outputs.key_id`)
- rule: write as {value: <literal>} or {valueFrom: {kind: AzureKeyVaultKey, name: <that resource's name>, fieldPath: status.outputs.key_id}} -- a bare string does not parse

### spec.connectionPolicy

`enum`

How clients establish connections. Unspecified applies Azure's
Default policy (Redirect inside Azure, Proxy from outside). REDIRECT
gives the lowest latency but needs ports 11000-11999 open in
addition to 1433; PROXY funnels everything through the gateway on
1433.

Allowed values (use exactly as shown):

- `azure_mssql_server_connection_policy_unspecified` -- Not specified: Azure's Default policy (Redirect for Azure-internal clients, Proxy for external ones).
- `DEFAULT` -- Redirect inside Azure, Proxy from outside -- Azure's default.
- `PROXY` -- All connections funnel through the regional gateway on port 1433. Simplest firewalling, highest latency.
- `REDIRECT` -- Clients are redirected to the database node after login -- the lowest latency, but requires ports 11000-11999 open alongside 1433.

### spec.minimumTlsVersion

`string` · optional (explicit presence)

The minimum TLS version for client connections. Azure requires
"1.2" (the only accepted value on current API versions -- older
protocol floors are retired). Unspecified applies "1.2". ARM rejects
removing the floor once set.

- default: `1.2`
- rule: minimum_tls_version must be 1.2 (older TLS floors are retired by Azure)

### spec.publicNetworkAccessEnabled

`bool` · optional (explicit presence)

Whether the server is reachable over the public internet. When true,
firewall_rules and virtual_network_rules control who connects; when
false, only private endpoints reach the server. Unspecified applies
Azure's default of true.

- default: `true`

### spec.outboundNetworkRestrictionEnabled

`bool`

Whether OUTBOUND connections from the server (elastic queries,
linked external tables) are restricted to the FQDNs in
outbound_firewall_rules. Azure's default is false (unrestricted).

### spec.outboundFirewallRules

`[]string`

The FQDNs the server may reach out to while outbound restriction is
enabled. Each entry is its own ARM outbound-firewall-rule
sub-resource.

- rule: {"repeated":{"items":{"string":{"minLen":"1"}}}}

### spec.expressVulnerabilityAssessmentEnabled

`bool`

Whether Microsoft Defender's express vulnerability assessment runs
on the server (agentless SQL scanning without a storage account).
Azure's default is false.

### spec.firewallRules

`[]AzureMssqlServerFirewallRule`

Public-endpoint firewall allowlist: each rule admits a contiguous
IPv4 range. Only meaningful while public network access is enabled.
The special rule 0.0.0.0-0.0.0.0 admits Azure-internal services only
(not the internet).

### spec.firewallRules[].name

`string` · required

The rule's name, unique within the server. E.g. "allow-office",
"allow-azure-services".

- rule: {"required":true,"string":{"minLen":"1","maxLen":"128"}}

### spec.firewallRules[].startIpAddress

`string` · required

The first IPv4 address of the admitted range (inclusive). Use
0.0.0.0 for both start and end to admit Azure-internal services.

- rule: {"required":true,"string":{"ipv4":true}}

### spec.firewallRules[].endIpAddress

`string` · required

The last IPv4 address of the admitted range (inclusive). Equal to
start_ip_address for a single-address rule.

- rule: {"required":true,"string":{"ipv4":true}}

### spec.virtualNetworkRules

`[]AzureMssqlServerVirtualNetworkRule`

Subnet allowlist: each rule admits one subnet through its
Microsoft.Sql service endpoint -- the classic (non-Private-Link) way
to keep traffic on the Azure backbone. Only meaningful while public
network access is enabled.

### spec.virtualNetworkRules[].name

`string` · required

The rule's name, unique within the server. E.g. "allow-app-subnet".

- rule: {"required":true,"string":{"minLen":"1","maxLen":"128"}}

### spec.virtualNetworkRules[].subnetId

`string | valueFrom` · required

The subnet to admit, by ARM ID. The subnet needs the Microsoft.Sql
service endpoint enabled (see AzureSubnet.service_endpoints) unless
ignore_missing_vnet_service_endpoint is set.

- references: AzureSubnet (`status.outputs.subnet_id`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: AzureSubnet, name: <that resource's name>, fieldPath: status.outputs.subnet_id}} -- a bare string does not parse

### spec.virtualNetworkRules[].ignoreMissingVnetServiceEndpoint

`bool`

Create the rule even while the subnet's Microsoft.Sql service
endpoint is not yet enabled (traffic flows once it is). Azure's
default is false.

### spec.extendedAuditing

`AzureMssqlServerExtendedAuditing`

Server-level SQL auditing: writes audit events for every database on
the server to blob storage and/or Azure Monitor. Presence enables
auditing; omit the block to leave auditing off.

- rule: storage_account_access_key requires storage_endpoint

### spec.extendedAuditing.storageEndpoint

`string`

The blob-storage endpoint audit logs are written to (e.g.
https://<account>.blob.core.windows.net). Omit for
Azure-Monitor-only auditing.

- rule: storage_endpoint must be an https:// URL

### spec.extendedAuditing.storageAccountAccessKey

`string | valueFrom` · sensitive

The storage account's access key for the audit-log container.
Required with storage_endpoint unless the server's managed identity
holds the Storage Blob Data Contributor role on the account.

- rule: write as {value: <literal>} or {valueFrom: {kind: <Kind>, name: <that resource's name>, fieldPath: status.outputs.<output>}} -- a bare string does not parse

### spec.extendedAuditing.storageAccountAccessKeyIsSecondary

`bool`

Whether the access key belongs to the storage account's secondary
key slot (documents which key rotates). Azure's default is false.

### spec.extendedAuditing.retentionInDays

`int32` · optional (explicit presence)

How many days audit logs are retained in blob storage. 0 (Azure's
default) retains them indefinitely.

- rule: {"int32":{"lte":3285,"gte":0}}

### spec.extendedAuditing.logMonitoringEnabled

`bool` · optional (explicit presence)

Whether audit events are also sent to Azure Monitor (surfaced
through the server's diagnostic settings as SQLSecurityAuditEvents).
Azure's default is true.

- default: `true`

### spec.extendedAuditing.storageAccountSubscriptionId

`string` · optional (explicit presence)

The subscription of the audit storage account, when it lives in a
different subscription than the server.

- rule: {"string":{"uuid":true}}

### spec.extendedAuditing.predicateExpression

`string`

A T-SQL predicate that filters which audit records are captured
(e.g. "statement <> 'select 1'"). Omit to capture everything the
action groups select.

### spec.extendedAuditing.auditActionsAndGroups

`[]string`

The audit action groups (and actions) to capture. Omit for Azure's
default set: SUCCESSFUL_DATABASE_AUTHENTICATION_GROUP,
FAILED_DATABASE_AUTHENTICATION_GROUP, BATCH_COMPLETED_GROUP.

- rule: {"repeated":{"items":{"string":{"minLen":"1"}}}}

### spec.securityAlertPolicy

`AzureMssqlServerSecurityAlertPolicy`

Microsoft Defender for SQL threat detection at the server scope:
alerting on SQL injection, anomalous access, data exfiltration, and
unsafe actions. Presence configures the policy; omit the block to
leave Defender's advanced threat protection unconfigured.

- rule: storage_endpoint and storage_account_access_key must be set together

### spec.securityAlertPolicy.state

`enum` · required

Whether the policy is enforced. ENABLED turns Defender's advanced
threat protection on for every database on the server.

- rule: {"required":true}

Allowed values (use exactly as shown):

- `azure_mssql_server_security_alert_policy_state_unspecified` -- Not specified -- invalid; choose an explicit state when the policy block is present.
- `ENABLED` -- Threat detection is on.
- `DISABLED` -- Threat detection is configured but off.

### spec.securityAlertPolicy.disabledAlerts

`[]enum`

Alert classes to suppress. Each entry disables one detector.

- rule: {"repeated":{"items":{"enum":{"definedOnly":true,"notIn":[0]}}}}

Allowed values (use exactly as shown):

- `azure_mssql_server_security_alert_type_unspecified` -- Not specified -- invalid as a disabled_alerts entry.
- `SQL_INJECTION` -- Potential SQL injection attacks.
- `SQL_INJECTION_VULNERABILITY` -- Application code vulnerable to SQL injection.
- `ACCESS_ANOMALY` -- Logins from unusual locations, principals, or applications.
- `DATA_EXFILTRATION` -- Unusually large result sets or data movement.
- `UNSAFE_ACTION` -- High-privilege or destructive statements.

### spec.securityAlertPolicy.emailAccountAdmins

`bool`

Whether alert emails also go to the subscription administrators.
Azure's default is false.

### spec.securityAlertPolicy.emailAddresses

`[]string`

Additional email addresses that receive alerts.

- rule: {"repeated":{"items":{"string":{"minLen":"3"}}}}

### spec.securityAlertPolicy.retentionDays

`int32` · optional (explicit presence)

How many days threat-detection audit records are retained in the
export storage account. 0 (Azure's default) retains indefinitely.

- rule: {"int32":{"gte":0}}

### spec.securityAlertPolicy.storageEndpoint

`string`

The blob-storage endpoint threat-detection audit records are
exported to. Set together with storage_account_access_key; omit both
to keep alerts without a storage export.

### spec.securityAlertPolicy.storageAccountAccessKey

`string | valueFrom` · sensitive

The access key of the export storage account. Required with
storage_endpoint.

- rule: write as {value: <literal>} or {valueFrom: {kind: <Kind>, name: <that resource's name>, fieldPath: status.outputs.<output>}} -- a bare string does not parse

### spec.tags

`map<string, string>`

Free-form tags applied to the server, merged over the
Planton-derived resource tags (organization, environment, resource
id); a user tag with the same key wins. Tags are Azure's governance
surface -- Azure Policy enforces them and Microsoft Cost Management
groups by them.

## Validation Rules

- `mssql_at_least_one_auth`: configure administrator_login (SQL auth) and/or an azuread_administrator with azuread_authentication_only (Entra-only auth)
- `mssql_sql_auth_credentials_pair`: administrator_login and administrator_password must be set together
- `mssql_aad_only_forbids_sql_credentials`: an Entra-only server (azuread_authentication_only = true) must omit administrator_login and administrator_password
- `mssql_primary_uai_requires_identity`: primary_user_assigned_identity_id requires an identity block whose type includes USER_ASSIGNED
- `mssql_tde_cmk_requires_identity`: transparent_data_encryption_key_vault_key_id requires an identity block (the identity unwraps the key)
- `mssql_outbound_rules_require_restriction`: outbound_firewall_rules require outbound_network_restriction_enabled = true

## Outputs

Reference an output from another manifest as `valueFrom: {kind: AzureMssqlServer, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.server_id` | `string` | The Azure Resource Manager ID of the logical server. Format: /subscriptions/{sub}/resourceGroups/{rg}/providers/Microsoft.Sql/servers/{name} Referenced by AzureMssqlDatabase.server_id, AzureMssqlElasticPool.server_id, and AzurePrivateEndpoint (private_connection_resource_id). |
| `status.outputs.server_name` | `string` | The name of the logical server. |
| `status.outputs.fqdn` | `string` | The server's fully qualified domain name ({name}.database.windows.net) -- the connection-string host. With private endpoints it resolves privately through the privatelink.database.windows.net zone. |
| `status.outputs.administrator_login` | `string` | The administrator login, echoed so applications can construct connection strings without duplicating the value. Empty on an Entra-only server (SQL logins disabled). |
| `status.outputs.identity_principal_id` | `string` | The principal (directory object) ID of the server's system-assigned managed identity -- the subject for AzureRoleAssignment grants. Empty unless the identity type includes SYSTEM_ASSIGNED. |

## References

Fields that can point at another resource's outputs:

| Field | Kind | Output |
|---|---|---|
| `spec.resourceGroup` | AzureResourceGroup | `status.outputs.resource_group_name` |
| `spec.azureadAdministrator.objectId` | AzureUserAssignedIdentity | `status.outputs.principal_id` |
| `spec.identity.identityIds` | AzureUserAssignedIdentity | `status.outputs.identity_id` |
| `spec.primaryUserAssignedIdentityId` | AzureUserAssignedIdentity | `status.outputs.identity_id` |
| `spec.transparentDataEncryptionKeyVaultKeyId` | AzureKeyVaultKey | `status.outputs.key_id` |
| `spec.virtualNetworkRules[].subnetId` | AzureSubnet | `status.outputs.subnet_id` |

## Referenced By

Fields on other kinds that can point at this resource:

| Kind | Field | Reads |
|---|---|---|
| AzureMssqlDatabase | `spec.serverId` | `status.outputs.server_id` |
| AzureMssqlElasticPool | `spec.serverId` | `status.outputs.server_id` |
| AzureMssqlFailoverGroup | `spec.serverId` | `status.outputs.server_id` |
| AzureMssqlFailoverGroup | `spec.partnerServers[].serverId` | `status.outputs.server_id` |

## See Also

- [Overview](./README.md)
- [Design notes](./docs/README.md)
