# AzureFirewallPolicy

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `azure.planton.dev/v1alpha1`

**AzureFirewallPolicySpec** defines an Azure Firewall Policy: the
reusable rule-and-inspection document that Azure Firewall instances
enforce. The policy carries everything about WHAT the firewall does --
threat intelligence posture, DNS proxying, TLS inspection, intrusion
detection, SNAT behavior -- while the firewall instance (AzureFirewall)
carries WHERE it runs (subnet, public IPs, zones). One policy can be
attached to many firewalls across regions, so security policy is
authored once and enforced everywhere.

Rules themselves live in AzureFirewallPolicyRuleCollectionGroup children
(many per policy, each an independently-deployable ordered document);
the policy is the container and the inspection/posture surface.

Policies support inheritance: a policy may reference a base policy
(base_policy_id) whose rules and settings it extends -- the standard
enterprise pattern is a global baseline policy owned by the security
team, with per-application child policies layering on top.

**Tiers**: BASIC (SMB-scale, curated feature subset), STANDARD (the
production default -- full rule engine, threat intelligence), PREMIUM
(adds TLS inspection, IDPS intrusion detection, URL filtering, web
categories). The tier is fixed at creation -- changing it replaces the
policy -- and must match the tier of every firewall that attaches it.

**ForceNew fields**: `name`, `region`, `resource_group`, `sku`.
Everything else updates in place.

## Example

```yaml
apiVersion: azure.planton.dev/v1alpha1
kind: AzureFirewallPolicy
metadata:
  name: test-firewall-policy
spec:
  region: eastus
  resourceGroup:
    value: test-rg
  name: egress-baseline
  sku: PREMIUM
  threatIntelligenceMode: DENY
  threatIntelligenceAllowlist:
    ipAddresses:
      - "203.0.113.7"
    fqdns:
      - "partner.example.com"
  dns:
    servers:
      - "10.0.0.4"
    proxyEnabled: true
  intrusionDetection:
    mode: IDPS_ALERT
    signatureOverrides:
      - id: "2024897"
        state: IDPS_OFF
    privateRanges:
      - "10.0.0.0/8"
    trafficBypass:
      - name: trusted-backup-flow
        description: backup traffic already encrypted end to end
        protocol: TCP
        sourceAddresses:
          - "10.0.1.0/24"
        destinationAddresses:
          - "10.0.2.10"
        destinationPorts:
          - "8443"
  identity:
    type: USER_ASSIGNED
    userAssignedIdentityIds:
      - value: /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/test-rg/providers/Microsoft.ManagedIdentity/userAssignedIdentities/fw-tls
  tlsCertificate:
    keyVaultSecretId:
      value: https://test-kv.vault.azure.net/secrets/egress-ca
    name: egress-ca
  insights:
    enabled: true
    defaultLogAnalyticsWorkspaceId:
      value: /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/test-rg/providers/Microsoft.OperationalInsights/workspaces/central-law
    retentionInDays: 30
  explicitProxy:
    enabled: true
    httpPort: 8087
    httpsPort: 8088
  sqlRedirectAllowed: true
  privateIpRanges:
    - "10.0.0.0/8"
    - "100.64.0.0/10"
  autoLearnPrivateRangesEnabled: true
  tags:
    purpose: hack-test
```

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.region` | `string` | yes |  |  |
| `spec.resourceGroup` | `string \| valueFrom` | yes |  | AzureResourceGroup (`status.outputs.resource_group_name`) |
| `spec.name` | `string` | yes |  |  |
| `spec.sku` | `enum` |  |  |  |
| `spec.basePolicyId` | `string \| valueFrom` |  |  | AzureFirewallPolicy (`status.outputs.firewall_policy_id`) |
| `spec.threatIntelligenceMode` | `enum` |  |  |  |
| `spec.threatIntelligenceAllowlist` | `AzureFirewallPolicyThreatIntelligenceAllowlist` |  |  |  |
| `spec.threatIntelligenceAllowlist.ipAddresses` | `[]string` |  |  |  |
| `spec.threatIntelligenceAllowlist.fqdns` | `[]string` |  |  |  |
| `spec.dns` | `AzureFirewallPolicyDns` |  |  |  |
| `spec.dns.servers` | `[]string` |  |  |  |
| `spec.dns.proxyEnabled` | `bool` |  |  |  |
| `spec.intrusionDetection` | `AzureFirewallPolicyIntrusionDetection` |  |  |  |
| `spec.intrusionDetection.mode` | `enum` |  |  |  |
| `spec.intrusionDetection.signatureOverrides` | `[]AzureFirewallPolicyIdpsSignatureOverride` |  |  |  |
| `spec.intrusionDetection.signatureOverrides[].id` | `string` |  |  |  |
| `spec.intrusionDetection.signatureOverrides[].state` | `enum` |  |  |  |
| `spec.intrusionDetection.privateRanges` | `[]string` |  |  |  |
| `spec.intrusionDetection.trafficBypass` | `[]AzureFirewallPolicyIdpsTrafficBypass` |  |  |  |
| `spec.intrusionDetection.trafficBypass[].name` | `string` | yes |  |  |
| `spec.intrusionDetection.trafficBypass[].description` | `string` |  |  |  |
| `spec.intrusionDetection.trafficBypass[].protocol` | `enum` | yes |  |  |
| `spec.intrusionDetection.trafficBypass[].sourceAddresses` | `[]string` |  |  |  |
| `spec.intrusionDetection.trafficBypass[].sourceIpGroups` | `[]string \| valueFrom` |  |  | AzureIpGroup (`status.outputs.ip_group_id`) |
| `spec.intrusionDetection.trafficBypass[].destinationAddresses` | `[]string` |  |  |  |
| `spec.intrusionDetection.trafficBypass[].destinationIpGroups` | `[]string \| valueFrom` |  |  | AzureIpGroup (`status.outputs.ip_group_id`) |
| `spec.intrusionDetection.trafficBypass[].destinationPorts` | `[]string` |  |  |  |
| `spec.identity` | `AzureFirewallPolicyIdentity` |  |  |  |
| `spec.identity.type` | `enum` |  |  |  |
| `spec.identity.userAssignedIdentityIds` | `[]string \| valueFrom` |  |  | AzureUserAssignedIdentity (`status.outputs.identity_id`) |
| `spec.tlsCertificate` | `AzureFirewallPolicyTlsCertificate` |  |  |  |
| `spec.tlsCertificate.keyVaultSecretId` | `string \| valueFrom` | yes |  | AzureKeyVaultCertificate (`status.outputs.versionless_secret_id`) |
| `spec.tlsCertificate.name` | `string` | yes |  |  |
| `spec.insights` | `AzureFirewallPolicyInsights` |  |  |  |
| `spec.insights.enabled` | `bool` | yes |  |  |
| `spec.insights.defaultLogAnalyticsWorkspaceId` | `string \| valueFrom` | yes |  | AzureLogAnalyticsWorkspace (`status.outputs.workspace_id`) |
| `spec.insights.retentionInDays` | `int32` |  |  |  |
| `spec.insights.logAnalyticsWorkspaces` | `[]AzureFirewallPolicyInsightsWorkspace` |  |  |  |
| `spec.insights.logAnalyticsWorkspaces[].workspaceId` | `string \| valueFrom` | yes |  | AzureLogAnalyticsWorkspace (`status.outputs.workspace_id`) |
| `spec.insights.logAnalyticsWorkspaces[].firewallLocation` | `string` | yes |  |  |
| `spec.explicitProxy` | `AzureFirewallPolicyExplicitProxy` |  |  |  |
| `spec.explicitProxy.enabled` | `bool` |  |  |  |
| `spec.explicitProxy.httpPort` | `int32` |  |  |  |
| `spec.explicitProxy.httpsPort` | `int32` |  |  |  |
| `spec.explicitProxy.enablePacFile` | `bool` |  |  |  |
| `spec.explicitProxy.pacFilePort` | `int32` |  |  |  |
| `spec.explicitProxy.pacFile` | `string` |  |  |  |
| `spec.sqlRedirectAllowed` | `bool` |  |  |  |
| `spec.privateIpRanges` | `[]string` |  |  |  |
| `spec.autoLearnPrivateRangesEnabled` | `bool` |  |  |  |
| `spec.tags` | `map<string, string>` |  |  |  |

## Field Details

### spec.region

`string` · required

The Azure region the policy lives in, e.g. "eastus", "westeurope".
The policy is a regional resource but can be attached to firewalls in
any region. Changing the region replaces the policy.

- rule: {"required":true,"string":{"minLen":"1"}}

### spec.resourceGroup

`string | valueFrom` · required

The Azure resource group the policy is created in. Can be a literal
resource-group name or a reference to an AzureResourceGroup's name
output.

- references: AzureResourceGroup (`status.outputs.resource_group_name`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: AzureResourceGroup, name: <that resource's name>, fieldPath: status.outputs.resource_group_name}} -- a bare string does not parse

### spec.name

`string` · required

The policy's name, unique within the resource group. 2-80 characters;
must begin with a letter or number, end with a letter, number, or
underscore, and may contain only letters, numbers, underscores,
periods, or hyphens. Changing the name replaces the policy (and every
rule collection group and firewall attachment referencing it must be
re-pointed).

- rule: Firewall policy names are 2-80 characters, start with a letter or number, end with a letter, number, or underscore, and may contain alphanumerics, underscores, periods, and hyphens
- rule: {"required":true,"string":{"maxLen":"80"}}

### spec.sku

`enum`

The policy's pricing/capability tier. Fixed at creation (changing it
replaces the policy) and must match the tier of every firewall the
policy attaches to. Unspecified deploys STANDARD, the production
default.

- rule: {"enum":{"definedOnly":true}}

Allowed values (use exactly as shown):

- `azure_firewall_policy_sku_unspecified` -- Not specified -- deploys STANDARD, the production default.
- `BASIC` -- SMB-scale tier: L3-L7 filtering and threat-intelligence alerting only (Basic always alerts; it cannot deny on threat intelligence or run IDPS/TLS inspection). Requires the attached firewall to carry a management IP configuration.
- `STANDARD` -- The production default: full rule engine, threat intelligence alert-and-deny, DNS proxy, FQDN filtering in network rules.
- `PREMIUM` -- Everything in STANDARD plus TLS inspection, IDPS (signature-based intrusion detection and prevention), URL filtering, and web categories.

### spec.basePolicyId

`string | valueFrom`

The parent policy this one inherits from -- rules and settings of the
base apply beneath this policy's own. References an
AzureFirewallPolicy's ARM id. The base policy must live in the same
region and carry the same tier. The standard enterprise shape: a
security-team-owned global baseline, extended by per-application
child policies.

- references: AzureFirewallPolicy (`status.outputs.firewall_policy_id`)
- rule: write as {value: <literal>} or {valueFrom: {kind: AzureFirewallPolicy, name: <that resource's name>, fieldPath: status.outputs.firewall_policy_id}} -- a bare string does not parse

### spec.threatIntelligenceMode

`enum`

How the firewall acts on traffic matching Microsoft's threat
intelligence feed (known-malicious IPs and domains). ALERT (the
default) logs matches; DENY logs and blocks; OFF disables the feed.
Unspecified deploys ALERT -- the safe observability-first posture.

- rule: {"enum":{"definedOnly":true}}

Allowed values (use exactly as shown):

- `azure_firewall_policy_threat_intelligence_mode_unspecified` -- Not specified -- deploys ALERT, Azure's default.
- `ALERT` -- Log traffic that matches known-malicious IPs/domains, but let it flow. The observability-first default.
- `DENY` -- Log and block traffic matching the feed.
- `OFF` -- Disable threat-intelligence-based filtering entirely.

### spec.threatIntelligenceAllowlist

`AzureFirewallPolicyThreatIntelligenceAllowlist`

Traffic that threat intelligence must never flag: trusted addresses
and FQDNs exempted from the feed (e.g. a partner endpoint that trips
a false positive). At least one address or FQDN must be listed when
the block is present.

- rule: The threat intelligence allowlist must list at least one IP address or FQDN -- remove the block if there is nothing to exempt

### spec.threatIntelligenceAllowlist.ipAddresses

`[]string`

IP addresses and CIDR ranges the feed must never flag,
e.g. "203.0.113.7" or "10.0.0.0/16".

- rule: {"repeated":{"items":{"string":{"minLen":"1"}}}}

### spec.threatIntelligenceAllowlist.fqdns

`[]string`

Fully-qualified domain names the feed must never flag,
e.g. "partner.example.com".

- rule: {"repeated":{"items":{"string":{"minLen":"1"}}}}

### spec.dns

`AzureFirewallPolicyDns`

DNS settings for firewalls attached to this policy: custom upstream
servers and/or DNS proxying. DNS proxy is required for FQDN-based
network rules to resolve consistently (the firewall must see the same
answers clients do -- pointing clients at the firewall's proxy
guarantees it).

### spec.dns.servers

`[]string`

Custom upstream DNS servers (IPv4) the firewall resolves against.
When empty, Azure's default resolver is used.

- rule: {"repeated":{"items":{"string":{"minLen":"1"}}}}

### spec.dns.proxyEnabled

`bool`

Run the firewall as a DNS proxy: clients point their DNS at the
firewall's private IP, and the firewall resolves upstream. Required
for FQDN-based network rules to behave deterministically -- without
it, the firewall and the client can resolve the same name to
different addresses and the rule silently misses.

### spec.intrusionDetection

`AzureFirewallPolicyIntrusionDetection`

Premium-tier IDPS (intrusion detection and prevention): signature-based
detection of known attacks in transit, with per-signature overrides,
private-range scoping, and bypass lists for trusted flows. Requires
sku PREMIUM.

### spec.intrusionDetection.mode

`enum`

The IDPS engine mode: OFF (signatures loaded but not evaluated),
ALERT (log matched signatures), DENY (log and block). Individual
signatures can override this via signature_overrides.

- rule: {"enum":{"definedOnly":true}}

Allowed values (use exactly as shown):

- `azure_firewall_policy_intrusion_detection_state_unspecified` -- Not specified -- the engine mode is not sent and Azure's default (Off) applies.
- `IDPS_OFF` -- Signatures are not evaluated.
- `IDPS_ALERT` -- Matched signatures are logged; traffic flows.
- `IDPS_DENY` -- Matched signatures are logged and the traffic is blocked.

### spec.intrusionDetection.signatureOverrides

`[]AzureFirewallPolicyIdpsSignatureOverride`

Per-signature exceptions to the engine mode -- e.g. run the engine in
DENY but set a false-positive-prone signature to ALERT, or OFF.

### spec.intrusionDetection.signatureOverrides[].id

`string`

The numeric IDPS signature id to override (from the Azure Firewall
IDPS signature catalog, e.g. "2024897").

### spec.intrusionDetection.signatureOverrides[].state

`enum`

The state this signature runs in, overriding the engine mode.

- rule: {"enum":{"definedOnly":true}}

Allowed values (use exactly as shown):

- `azure_firewall_policy_intrusion_detection_state_unspecified` -- Not specified -- the engine mode is not sent and Azure's default (Off) applies.
- `IDPS_OFF` -- Signatures are not evaluated.
- `IDPS_ALERT` -- Matched signatures are logged; traffic flows.
- `IDPS_DENY` -- Matched signatures are logged and the traffic is blocked.

### spec.intrusionDetection.privateRanges

`[]string`

The ranges IDPS treats as the private (internal) side of traffic
direction classification. Defaults to the RFC 1918 ranges when empty.

- rule: {"repeated":{"items":{"string":{"minLen":"1"}}}}

### spec.intrusionDetection.trafficBypass

`[]AzureFirewallPolicyIdpsTrafficBypass`

Trusted flows IDPS must not inspect -- match by protocol, source,
destination, and ports. Use for latency-sensitive or already-encrypted
flows that signatures cannot meaningfully inspect.

### spec.intrusionDetection.trafficBypass[].name

`string` · required

A name for this bypass entry (shows in diagnostics).

- rule: {"required":true,"string":{"minLen":"1"}}

### spec.intrusionDetection.trafficBypass[].description

`string`

Why this flow is bypassed -- operator documentation.

### spec.intrusionDetection.trafficBypass[].protocol

`enum` · required

The protocol of the bypassed flow.

- rule: {"required":true,"enum":{"definedOnly":true,"notIn":[0]}}

Allowed values (use exactly as shown):

- `azure_firewall_policy_idps_bypass_protocol_unspecified` -- Not specified -- invalid; every bypass must declare its protocol.
- `ANY` -- Match any protocol.
- `TCP` -- Match TCP flows.
- `UDP` -- Match UDP flows.
- `ICMP` -- Match ICMP traffic.

### spec.intrusionDetection.trafficBypass[].sourceAddresses

`[]string`

Source IP addresses/CIDRs of the bypassed flow.

- rule: {"repeated":{"items":{"string":{"minLen":"1"}}}}

### spec.intrusionDetection.trafficBypass[].sourceIpGroups

`[]string | valueFrom`

Source IP Groups of the bypassed flow -- references to AzureIpGroup
ARM ids, the reusable alternative to literal source_addresses.

- references: AzureIpGroup (`status.outputs.ip_group_id`)
- rule: write as {value: <literal>} or {valueFrom: {kind: AzureIpGroup, name: <that resource's name>, fieldPath: status.outputs.ip_group_id}} -- a bare string does not parse

### spec.intrusionDetection.trafficBypass[].destinationAddresses

`[]string`

Destination IP addresses/CIDRs of the bypassed flow.

- rule: {"repeated":{"items":{"string":{"minLen":"1"}}}}

### spec.intrusionDetection.trafficBypass[].destinationIpGroups

`[]string | valueFrom`

Destination IP Groups of the bypassed flow -- references to
AzureIpGroup ARM ids.

- references: AzureIpGroup (`status.outputs.ip_group_id`)
- rule: write as {value: <literal>} or {valueFrom: {kind: AzureIpGroup, name: <that resource's name>, fieldPath: status.outputs.ip_group_id}} -- a bare string does not parse

### spec.intrusionDetection.trafficBypass[].destinationPorts

`[]string`

Destination ports of the bypassed flow, e.g. "443" or "*".

- rule: {"repeated":{"items":{"string":{"minLen":"1"}}}}

### spec.identity

`AzureFirewallPolicyIdentity`

The policy's managed identity. Required in practice for TLS
inspection: Azure Firewall uses a USER_ASSIGNED identity to read the
CA certificate from Key Vault (grant the identity Key Vault secret
read access -- "Key Vault Secrets User" under RBAC).

- rule: List at least one user-assigned identity id when the identity type includes USER_ASSIGNED
- rule: user_assigned_identity_ids is only used when the identity type includes USER_ASSIGNED -- change the type or remove the ids

### spec.identity.type

`enum`

The identity model: SYSTEM_ASSIGNED (Azure creates and rotates a
service principal bound to the policy's lifecycle), USER_ASSIGNED
(bring identities from user_assigned_identity_ids, shareable across
resources), or SYSTEM_AND_USER_ASSIGNED (both). TLS inspection
requires a user-assigned identity.

- rule: {"enum":{"definedOnly":true,"notIn":[0]}}

Allowed values (use exactly as shown):

- `azure_firewall_policy_identity_type_unspecified` -- Not specified -- invalid; declare the identity model when the block is present.
- `SYSTEM_ASSIGNED` -- Azure creates and rotates a service principal bound to the policy's lifecycle.
- `USER_ASSIGNED` -- Bring user-assigned identities (shareable across resources). The model TLS inspection requires.
- `SYSTEM_AND_USER_ASSIGNED` -- Both a system-assigned principal and user-assigned identities.

### spec.identity.userAssignedIdentityIds

`[]string | valueFrom`

The user-assigned identities to attach -- required when (and only
meaningful when) type includes USER_ASSIGNED. Each entry references
an AzureUserAssignedIdentity's ARM id.

- references: AzureUserAssignedIdentity (`status.outputs.identity_id`)
- rule: write as {value: <literal>} or {valueFrom: {kind: AzureUserAssignedIdentity, name: <that resource's name>, fieldPath: status.outputs.identity_id}} -- a bare string does not parse

### spec.tlsCertificate

`AzureFirewallPolicyTlsCertificate`

Premium-tier TLS inspection: the intermediate CA certificate the
firewall uses to decrypt, inspect, and re-encrypt outbound TLS.
Requires sku PREMIUM and a user-assigned identity with read access to
the Key Vault secret.

### spec.tlsCertificate.keyVaultSecretId

`string | valueFrom` · required

The Key Vault SECRET identifier of the CA certificate (the
certificate's secret face, carrying the private key). Defaults to an
AzureKeyVaultCertificate's versionless secret id so the firewall
follows certificate renewals automatically; reference the versioned
secret_id explicitly to pin one version. The policy's user-assigned
identity must have read access to this secret.

- references: AzureKeyVaultCertificate (`status.outputs.versionless_secret_id`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: AzureKeyVaultCertificate, name: <that resource's name>, fieldPath: status.outputs.versionless_secret_id}} -- a bare string does not parse

### spec.tlsCertificate.name

`string` · required

A display name for the certificate authority as it appears on the
policy (conventionally the Key Vault certificate's name).

- rule: {"required":true,"string":{"minLen":"1"}}

### spec.insights

`AzureFirewallPolicyInsights`

Firewall Policy Analytics: route the policy's traffic analysis into a
Log Analytics workspace, with optional per-region workspace routing
for multi-region firewall fleets.

### spec.insights.enabled

`bool` · required · optional (explicit presence)

Whether policy analytics is on. An explicit false keeps the
workspace wiring in place with analysis paused.

- rule: {"required":true}

### spec.insights.defaultLogAnalyticsWorkspaceId

`string | valueFrom` · required

The default Log Analytics workspace analytics flow into -- references
an AzureLogAnalyticsWorkspace's ARM id.

- references: AzureLogAnalyticsWorkspace (`status.outputs.workspace_id`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: AzureLogAnalyticsWorkspace, name: <that resource's name>, fieldPath: status.outputs.workspace_id}} -- a bare string does not parse

### spec.insights.retentionInDays

`int32` · optional (explicit presence)

How many days of analytics to retain (0 uses the workspace's own
retention).

- rule: {"int32":{"gte":0}}

### spec.insights.logAnalyticsWorkspaces

`[]AzureFirewallPolicyInsightsWorkspace`

Per-region workspace routing for multi-region firewall fleets: a
firewall in `firewall_location` logs to `workspace_id` instead of the
default workspace.

### spec.insights.logAnalyticsWorkspaces[].workspaceId

`string | valueFrom` · required

The Log Analytics workspace for this region -- references an
AzureLogAnalyticsWorkspace's ARM id.

- references: AzureLogAnalyticsWorkspace (`status.outputs.workspace_id`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: AzureLogAnalyticsWorkspace, name: <that resource's name>, fieldPath: status.outputs.workspace_id}} -- a bare string does not parse

### spec.insights.logAnalyticsWorkspaces[].firewallLocation

`string` · required

The Azure region whose firewalls log to this workspace,
e.g. "westeurope".

- rule: {"required":true,"string":{"minLen":"1"}}

### spec.explicitProxy

`AzureFirewallPolicyExplicitProxy`

Explicit proxy: expose the firewall as a forward proxy that clients
address directly (by proxy setting or PAC file) instead of routing
through it transparently. Useful when routing to the firewall cannot
be guaranteed (e.g. traffic from networks Azure does not control).

### spec.explicitProxy.enabled

`bool`

Whether explicit proxy is on.

### spec.explicitProxy.httpPort

`int32` · optional (explicit presence)

The port the proxy listens on for HTTP, e.g. 8087. The azurerm
provider enforces an upper bound of 35536 on proxy ports (its
published validation), so values above that are rejected at deploy
time even though they parse as ports.

- rule: {"int32":{"lte":35536,"gte":0}}

### spec.explicitProxy.httpsPort

`int32` · optional (explicit presence)

The port the proxy listens on for HTTPS, e.g. 8088.

- rule: {"int32":{"lte":35536,"gte":0}}

### spec.explicitProxy.enablePacFile

`bool`

Serve clients a PAC (proxy auto-configuration) file instead of a
static proxy setting. Pair with pac_file_port and pac_file.

### spec.explicitProxy.pacFilePort

`int32` · optional (explicit presence)

The port the PAC file is served on, e.g. 8089.

- rule: {"int32":{"lte":35536,"gte":0}}

### spec.explicitProxy.pacFile

`string`

The SAS URL of the PAC file blob the firewall serves to clients.
Azure requires the port and file together when PAC is enabled.

### spec.sqlRedirectAllowed

`bool`

Allow SQL redirect traffic (ports 11000-11999 and 14000-14999) to be
filtered by FQDN in network rules. Required when Azure SQL clients
use the default Redirect connection policy through the firewall;
leaving it off restricts SQL FQDN filtering to the proxy port 1433.

### spec.privateIpRanges

`[]string`

The address ranges the firewall treats as PRIVATE (traffic to them is
never SNATed). Each entry is a CIDR or a single IPv4 address. When
unset, Azure defaults to the IANA RFC 1918 ranges. Replace the list
when your organization routes public ranges internally (e.g. carrier
NAT space) or wants SNAT for specific private destinations.

- rule: {"repeated":{"items":{"string":{"minLen":"1"}}}}

### spec.autoLearnPrivateRangesEnabled

`bool`

Let Azure learn the private ranges automatically from route tables
and VNet address space instead of (or in addition to) the static
private_ip_ranges list. Azure only records "Enabled" -- turning it
back off is done by omission on the wire.

### spec.tags

`map<string, string>`

Free-form tags applied to the policy, merged over the Planton-derived
resource tags (organization, environment, resource id); a user tag
with the same key wins.

## Validation Rules

- `intrusion_detection_requires_premium`: Intrusion detection (IDPS) is an Azure Firewall Premium feature -- set sku to PREMIUM to use it
- `tls_certificate_requires_premium`: TLS inspection is an Azure Firewall Premium feature -- set sku to PREMIUM to use tls_certificate

## Outputs

Reference an output from another manifest as `valueFrom: {kind: AzureFirewallPolicy, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.firewall_policy_id` | `string` | The Azure Resource Manager ID of the firewall policy. This is the composition seam: rule collection groups nest under it (firewall_policy_id), firewalls attach it (firewall_policy_id), and child policies inherit from it (base_policy_id). Format: /subscriptions/{sub}/resourceGroups/{rg}/providers/Microsoft.Network/firewallPolicies/{name} |
| `status.outputs.firewall_policy_name` | `string` | The name of the firewall policy resource. |
| `status.outputs.identity_principal_id` | `string` | The principal (object) id of the policy's system-assigned managed identity -- grant it Key Vault secret read access when TLS inspection rides the system identity. Empty when no system-assigned identity is configured. |

## References

Fields that can point at another resource's outputs:

| Field | Kind | Output |
|---|---|---|
| `spec.resourceGroup` | AzureResourceGroup | `status.outputs.resource_group_name` |
| `spec.basePolicyId` | AzureFirewallPolicy | `status.outputs.firewall_policy_id` |
| `spec.intrusionDetection.trafficBypass[].sourceIpGroups` | AzureIpGroup | `status.outputs.ip_group_id` |
| `spec.intrusionDetection.trafficBypass[].destinationIpGroups` | AzureIpGroup | `status.outputs.ip_group_id` |
| `spec.identity.userAssignedIdentityIds` | AzureUserAssignedIdentity | `status.outputs.identity_id` |
| `spec.tlsCertificate.keyVaultSecretId` | AzureKeyVaultCertificate | `status.outputs.versionless_secret_id` |
| `spec.insights.defaultLogAnalyticsWorkspaceId` | AzureLogAnalyticsWorkspace | `status.outputs.workspace_id` |
| `spec.insights.logAnalyticsWorkspaces[].workspaceId` | AzureLogAnalyticsWorkspace | `status.outputs.workspace_id` |

## Referenced By

Fields on other kinds that can point at this resource:

| Kind | Field | Reads |
|---|---|---|
| AzureFirewall | `spec.firewallPolicyId` | `status.outputs.firewall_policy_id` |
| AzureFirewallPolicy | `spec.basePolicyId` | `status.outputs.firewall_policy_id` |
| AzureFirewallPolicyRuleCollectionGroup | `spec.firewallPolicyId` | `status.outputs.firewall_policy_id` |

## See Also

- [Overview](../README.md)
