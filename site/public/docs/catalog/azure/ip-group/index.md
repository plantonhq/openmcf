---
title: "IP Group"
description: "IP Group deployment documentation"
icon: "package"
order: 100
componentName: "azureipgroup"
---

# Azure IP Group

Creates an Azure IP Group -- a named, reusable set of IP addresses and CIDR ranges that Azure Firewall and Firewall Policy rules reference by ARM id instead of repeating literal address lists. Update the group once and every rule that references it follows.

## What Gets Created

When you deploy an AzureIpGroup resource, Planton provisions:

- **IP Group** — an `azurerm_ip_group` in the specified region and resource group, carrying the address set

The group is passive. Consumption is declared from the rule's side: `AzureFirewallPolicyRuleCollectionGroup` rules list source/destination IP Groups, and `AzureFirewallPolicy` intrusion-detection traffic bypasses reference them the same way.

## Prerequisites

- **Azure credentials** configured via environment variables or Planton provider config
- **A resource group** to create the group in (an `AzureResourceGroup` in composed environments)
- **Network write rights**: `Microsoft.Network/ipGroups/write` (Network Contributor, Contributor, or Owner)

## Quick Start

Create a file `ip-group.yaml`:

```yaml
apiVersion: azure.planton.dev/v1
kind: AzureIpGroup
metadata:
  name: branch-offices
  labels:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: prod.AzureIpGroup.branch-offices
spec:
  region: eastus
  resourceGroup:
    value: network-rg
  name: branch-offices
  cidrs:
    - "198.51.100.0/24"
    - "203.0.113.0/24"
```

Deploy:

```shell
planton apply -f ip-group.yaml
```

After deployment, read `status.outputs.ip_group_id` for the ARM ID to reference from firewall policy rules.

## Configuration Reference

### Required Fields

| Field | Type | Description | Validation |
|-------|------|-------------|------------|
| `region` | `string` | Azure region. IP Groups are regional but referable from policies in any region. | Required |
| `resourceGroup` | `StringValueOrRef` | Resource group name. Defaults to referencing an `AzureResourceGroup`'s name output. | Required |
| `name` | `string` | Group name, unique within the resource group. | Required, 1-80 chars, Azure naming rules |

### Optional Fields

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `cidrs` | `list(string)` | `[]` | IP addresses ("203.0.113.7") and CIDR blocks ("10.10.0.0/16"), up to 5,000 entries. Updates in place. |
| `tags` | `map(string)` | `{}` | User tags, merged over Planton-derived tags (user wins on collision). |

## Examples

### Branch-Office Address Set

Curate the addresses once; reference them from every rule that means
"the branches":

```yaml
apiVersion: azure.planton.dev/v1
kind: AzureIpGroup
metadata:
  name: branch-offices
  labels:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: prod.AzureIpGroup.branch-offices
spec:
  region: eastus
  resourceGroup:
    valueFrom:
      name: network-rg
  name: branch-offices
  cidrs:
    - "198.51.100.0/24"
    - "203.0.113.0/24"
```

Reference it from a firewall policy network rule:

```yaml
spec:
  networkRuleCollections:
    - name: allow-branches-to-dc
      priority: 200
      action: ALLOW
      rules:
        - name: branches-to-dc-any
          protocols: [ANY]
          sourceIpGroups:
            - valueFrom:
                name: branch-offices
          destinationAddresses: ["10.0.0.0/8"]
          destinationPorts: ["*"]
```

## Stack Outputs

After deployment, the following outputs are available in `status.outputs`:

| Output | Type | Description |
|--------|------|-------------|
| `ip_group_id` | `string` | Full ARM ID -- referenced by `AzureFirewallPolicyRuleCollectionGroup` rules and `AzureFirewallPolicy` IDPS bypasses |
| `ip_group_name` | `string` | The group's name as deployed |

## Related Components

- [AzureFirewallPolicyRuleCollectionGroup](/docs/catalog/azure/firewall-policy-rule-collection-group) — references IP Groups in rule source/destination fields
- [AzureFirewallPolicy](/docs/catalog/azure/firewall-policy) — references IP Groups in intrusion-detection traffic bypasses
- [AzureFirewall](/docs/catalog/azure/firewall) — enforces the policies whose rules target IP Groups
- [AzureResourceGroup](/docs/catalog/azure/resource-group) — provides the resource group for placement
