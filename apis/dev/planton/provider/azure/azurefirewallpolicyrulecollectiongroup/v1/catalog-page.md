# Azure Firewall Policy Rule Collection Group

Creates a rule collection group inside an Azure Firewall Policy -- an ordered document of application, network, and DNAT rule collections. A policy carries many groups (one per team or application), each deployed independently; the policy's firewalls enforce them all.

## What Gets Created

When you deploy an AzureFirewallPolicyRuleCollectionGroup resource, Planton provisions:

- **Rule Collection Group** — an `azurerm_firewall_policy_rule_collection_group` nested under the referenced firewall policy

Azure evaluates groups by priority, collections by priority within a type, and across types always DNAT → network → application. A matching DNAT rule implicitly allows the translated flow.

## Prerequisites

- **Azure credentials** configured via environment variables or Planton provider config
- **A firewall policy** to nest under (an `AzureFirewallPolicy`)
- **Network write rights**: `Microsoft.Network/firewallPolicies/ruleCollectionGroups/write`

## Quick Start

Create a file `rule-collection-group.yaml`:

```yaml
apiVersion: azure.planton.dev/v1
kind: AzureFirewallPolicyRuleCollectionGroup
metadata:
  name: platform-baseline
  labels:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: prod.AzureFirewallPolicyRuleCollectionGroup.platform-baseline
spec:
  firewallPolicyId:
    valueFrom:
      name: egress-baseline
  name: platform-baseline
  priority: 500
  networkRuleCollections:
    - name: core-egress
      priority: 200
      action: ALLOW
      rules:
        - name: allow-dns
          protocols: [UDP]
          sourceAddresses: ["10.0.0.0/16"]
          destinationAddresses: ["8.8.8.8"]
          destinationPorts: ["53"]
```

Deploy:

```shell
planton apply -f rule-collection-group.yaml
```

## Configuration Reference

### Required Fields

| Field | Type | Description | Validation |
|-------|------|-------------|------------|
| `firewallPolicyId` | `StringValueOrRef` | The parent policy's ARM id. Defaults to referencing an `AzureFirewallPolicy`. | Required, ForceNew |
| `name` | `string` | Group name, unique within the policy. | Required, 2-80 chars, ForceNew |
| `priority` | `number` | Evaluation priority among the policy's groups. | Required, 100-65000 |

### Rule Collections

| Field | Rules carry | Action | Notes |
|-------|-------------|--------|-------|
| `applicationRuleCollections` | L7 protocol pairs (HTTP/HTTPS/MSSQL + port), FQDNs, URLs, FQDN tags, web categories, TLS termination, header injection | `ALLOW`/`DENY` | URLs/categories/TLS need a PREMIUM policy |
| `networkRuleCollections` | protocols (ANY/TCP/UDP/ICMP), addresses, IP Groups, FQDNs, ports | `ALLOW`/`DENY` | FQDN destinations need the policy's DNS proxy |
| `natRuleCollections` | TCP/UDP, the firewall public IP + one port entry, translated address XOR FQDN + port | always `Dnat` (not modeled) | a match implicitly allows the flow |

## Examples

### DNAT: Publish a Jumpbox

```yaml
spec:
  firewallPolicyId:
    valueFrom:
      name: egress-baseline
  name: inbound-dnat
  priority: 100
  natRuleCollections:
    - name: inbound
      priority: 100
      rules:
        - name: rdp-to-jumpbox
          protocols: [TCP]
          sourceAddresses: ["*"]
          destinationAddress: "203.0.113.10"
          destinationPorts: ["3389"]
          translatedAddress: "10.0.1.4"
          translatedPort: 3389
```

### IP-Group-Based Egress

```yaml
spec:
  firewallPolicyId:
    valueFrom:
      name: egress-baseline
  name: branch-egress
  priority: 400
  networkRuleCollections:
    - name: branches-to-dc
      priority: 200
      action: ALLOW
      rules:
        - name: branches-to-dc-tcp
          protocols: [TCP]
          sourceIpGroups:
            - valueFrom:
                name: branch-offices
          destinationIpGroups:
            - valueFrom:
                name: on-prem-datacenter
          destinationPorts: ["443", "1433"]
```

## Stack Outputs

| Output | Type | Description |
|--------|------|-------------|
| `rule_collection_group_id` | `string` | Full nested ARM ID of the group |
| `rule_collection_group_name` | `string` | The group's name as deployed |

## Related Components

- [AzureFirewallPolicy](/docs/catalog/azure/firewall-policy) — the parent policy
- [AzureFirewall](/docs/catalog/azure/firewall) — enforces the policy's groups
- [AzureIpGroup](/docs/catalog/azure/ip-group) — reusable address sets for rule sources/destinations
- [AzurePublicIp](/docs/catalog/azure/public-ip) — the DNAT destination address comes from the firewall's public IP
