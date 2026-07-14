# Azure Firewall Policy

Creates an Azure Firewall Policy -- the reusable rule-and-inspection document Azure Firewall instances enforce. Author security policy once (threat intelligence, DNS proxying, TLS inspection, IDPS, SNAT posture) and attach it to firewalls anywhere; rules nest under the policy as independently-deployed rule collection groups.

## What Gets Created

When you deploy an AzureFirewallPolicy resource, Planton provisions:

- **Firewall Policy** — an `azurerm_firewall_policy` in the specified region and resource group, carrying the tier, inspection settings, and posture

Rules are separate `AzureFirewallPolicyRuleCollectionGroup` resources; enforcement is a separate `AzureFirewall` that references this policy.

## Prerequisites

- **Azure credentials** configured via environment variables or Planton provider config
- **A resource group** to create the policy in (an `AzureResourceGroup` in composed environments)
- **Network write rights**: `Microsoft.Network/firewallPolicies/write` (Network Contributor, Contributor, or Owner)
- For TLS inspection (Premium): an `AzureKeyVaultCertificate` holding the CA, and an `AzureUserAssignedIdentity` with secret read access on it

## Quick Start

Create a file `firewall-policy.yaml`:

```yaml
apiVersion: azure.planton.dev/v1
kind: AzureFirewallPolicy
metadata:
  name: egress-baseline
  labels:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: prod.AzureFirewallPolicy.egress-baseline
spec:
  region: eastus
  resourceGroup:
    value: network-rg
  name: egress-baseline
  threatIntelligenceMode: DENY
  dns:
    proxyEnabled: true
```

Deploy:

```shell
planton apply -f firewall-policy.yaml
```

After deployment, read `status.outputs.firewall_policy_id` for the ARM ID that rule collection groups and firewalls reference.

## Configuration Reference

### Required Fields

| Field | Type | Description | Validation |
|-------|------|-------------|------------|
| `region` | `string` | Azure region. The policy is regional but attachable to firewalls anywhere. | Required |
| `resourceGroup` | `StringValueOrRef` | Resource group name. Defaults to referencing an `AzureResourceGroup`'s name output. | Required |
| `name` | `string` | Policy name, unique within the resource group. | Required, 2-80 chars, Azure naming rules |

### Optional Fields

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `sku` | `enum` | `STANDARD` | `BASIC`, `STANDARD`, or `PREMIUM`. Fixed at creation; must match attached firewalls' tier. |
| `basePolicyId` | `StringValueOrRef` | -- | Parent policy for inheritance. |
| `threatIntelligenceMode` | `enum` | `ALERT` | `ALERT`, `DENY`, or `OFF`. |
| `threatIntelligenceAllowlist` | `object` | -- | IPs/FQDNs the feed must never flag (at least one entry). |
| `dns` | `object` | -- | Custom upstream servers and/or DNS proxying (required for FQDN network rules). |
| `intrusionDetection` | `object` | -- | IDPS engine mode, signature overrides, private ranges, traffic bypasses. **PREMIUM only** (validated). |
| `identity` | `object` | -- | Managed identity; TLS inspection requires USER_ASSIGNED. |
| `tlsCertificate` | `object` | -- | The CA certificate's Key Vault secret reference + display name. **PREMIUM only** (validated). |
| `insights` | `object` | -- | Policy analytics into Log Analytics (default workspace + per-region routing). |
| `explicitProxy` | `object` | -- | Forward-proxy ports and PAC file (ports capped at 35536). |
| `sqlRedirectAllowed` | `bool` | `false` | Allow SQL redirect ports (11000-11999, 14000-14999) in FQDN rules. |
| `privateIpRanges` | `list(string)` | IANA private | Ranges never SNATed. |
| `autoLearnPrivateRangesEnabled` | `bool` | `false` | Learn SNAT-exempt ranges automatically. |
| `tags` | `map(string)` | `{}` | User tags, merged over Planton-derived tags. |

## Examples

### Premium Policy with TLS Inspection and IDPS

```yaml
apiVersion: azure.planton.dev/v1
kind: AzureFirewallPolicy
metadata:
  name: inspected-egress
  labels:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: prod.AzureFirewallPolicy.inspected-egress
spec:
  region: eastus
  resourceGroup:
    valueFrom:
      name: network-rg
  name: inspected-egress
  sku: PREMIUM
  threatIntelligenceMode: DENY
  dns:
    proxyEnabled: true
  identity:
    type: USER_ASSIGNED
    userAssignedIdentityIds:
      - valueFrom:
          name: fw-tls-identity
  tlsCertificate:
    keyVaultSecretId:
      valueFrom:
        name: egress-ca-cert
    name: egress-ca
  intrusionDetection:
    mode: IDPS_DENY
```

### Child Policy Inheriting a Baseline

```yaml
spec:
  region: eastus
  resourceGroup:
    valueFrom:
      name: network-rg
  name: payments-app-policy
  basePolicyId:
    valueFrom:
      name: egress-baseline
```

## Stack Outputs

| Output | Type | Description |
|--------|------|-------------|
| `firewall_policy_id` | `string` | Full ARM ID -- referenced by `AzureFirewallPolicyRuleCollectionGroup`, `AzureFirewall`, and child policies |
| `firewall_policy_name` | `string` | The policy's name as deployed |
| `identity_principal_id` | `string` | System-assigned identity principal (empty without one) |

## Related Components

- [AzureFirewallPolicyRuleCollectionGroup](/docs/catalog/azure/firewall-policy-rule-collection-group) — the rules, nested under this policy
- [AzureFirewall](/docs/catalog/azure/firewall) — the data plane enforcing this policy
- [AzureIpGroup](/docs/catalog/azure/ip-group) — reusable address sets referenced by rules and IDPS bypasses
- [AzureKeyVaultCertificate](/docs/catalog/azure/key-vault-certificate) — the TLS inspection CA
- [AzureUserAssignedIdentity](/docs/catalog/azure/user-assigned-identity) — reads the CA from Key Vault
- [AzureLogAnalyticsWorkspace](/docs/catalog/azure/log-analytics-workspace) — receives policy analytics
