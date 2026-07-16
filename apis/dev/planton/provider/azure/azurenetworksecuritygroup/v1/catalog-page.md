# Azure Network Security Group

Deploys an Azure Network Security Group (NSG) with priority-ordered security rules that control inbound and outbound traffic for Azure resources. The component bundles the NSG with its security rules because an NSG without rules relies entirely on Azure defaults, making the rules the substance of the resource.

## What Gets Created

When you deploy an AzureNetworkSecurityGroup resource, Planton provisions:

- **Network Security Group** — a `network.NetworkSecurityGroup` resource in the specified region and resource group, acting as a stateful firewall for associated subnets or NICs
- **Security Rules** — one security rule per entry in `securityRules`, managed as part of the NSG; rules update in place and take effect immediately for everything the group guards
- **Azure Tags** — Planton-derived resource tags (organization, environment, resource id) merged with any user-supplied `tags`, applied to the NSG for tracking and governance

The component does not create subnet-to-NSG associations. A subnet declares which NSG guards it (AzureSubnet's `networkSecurityGroupId`), keeping the NSG lifecycle independent of any particular subnet or NIC.

Azure automatically creates implicit default rules in every NSG (priorities 65000-65500) that allow VNet-to-VNet traffic, allow Azure Load Balancer probes, and deny all other inbound traffic. User-defined rules (priorities 100-4096) are evaluated before these defaults.

## Prerequisites

- **Azure credentials** configured via environment variables or Planton provider config
- **An Azure Resource Group** where the NSG will be created (can reference an AzureResourceGroup resource)
- **Network planning** — understand the traffic flows to allow or deny before defining security rules

## Quick Start

Create a file `nsg.yaml`:

```yaml
apiVersion: azure.planton.dev/v1
kind: AzureNetworkSecurityGroup
metadata:
  name: my-nsg
  annotations:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: dev.AzureNetworkSecurityGroup.my-nsg
spec:
  region: eastus
  resourceGroup:
    value: my-rg
  name: my-nsg
  securityRules:
    - name: allow-https-inbound
      priority: 100
      direction: INBOUND
      access: ALLOW
      protocol: TCP
      destinationPortRange: "443"
```

Deploy:

```shell
planton apply -f nsg.yaml
```

This creates an NSG with a single rule allowing inbound HTTPS traffic from any source. All other inbound traffic is handled by Azure's implicit default rules (VNet-to-VNet allowed, everything else denied).

## Configuration Reference

### Required Fields

| Field | Type | Description | Validation |
|-------|------|-------------|------------|
| `region` | `string` | Azure region for the NSG (e.g., `eastus`, `westeurope`). Must match the region of resources it will be associated with. | Required, minimum length 1 |
| `resourceGroup` | `StringValueOrRef` | Azure Resource Group name. Can reference an AzureResourceGroup resource via `valueFrom`. | Required |
| `name` | `string` | Name of the Network Security Group. Must be unique within the resource group. Allowed characters: alphanumeric, underscores, hyphens, periods. Must start with alphanumeric. | Required, 1-80 characters |

### Optional Fields

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `securityRules` | `AzureNetworkSecurityGroupRule[]` | `[]` | Security rules defining allowed or denied traffic flows. Rules are evaluated in priority order (lowest number first). An NSG with no rules relies on Azure defaults: allow VNet-to-VNet, allow load balancer probes, deny all other inbound, allow all outbound. |
| `tags` | `map<string, string>` | `{}` | Free-form tags applied to the NSG, merged over the Planton-derived resource tags (organization, environment, resource id); a user tag with the same key wins. |

### Security Rule Fields

Each entry in `securityRules` supports the following fields:

| Field | Type | Default | Required | Description |
|-------|------|---------|----------|-------------|
| `name` | `string` | — | Yes | Unique name within the NSG. Use descriptive names like `allow-https-inbound`. 1-80 characters. |
| `description` | `string` | — | No | Human-readable description of the rule's purpose. Maximum 140 characters. |
| `priority` | `int` | — | Yes | Evaluation priority. Lower numbers are evaluated first. Range: 100-4096. Use increments of 10 or 100 to leave room for future rules. |
| `direction` | `enum` | — | Yes | Traffic direction. Values: `INBOUND`, `OUTBOUND`. |
| `access` | `enum` | — | Yes | Access decision when the rule matches. Values: `ALLOW`, `DENY`. |
| `protocol` | `enum` | — | Yes | Network protocol. Values: `ANY`, `TCP`, `UDP`, `ICMP`, `AH`, `ESP`. |
| `sourcePortRange` | `string` | — | No | Source port (`443`), range (`1024-65535`), or `*` for any. Unset means any — the right choice for most rules, since source ports are typically ephemeral. Never combined with `sourcePortRanges`. |
| `sourcePortRanges` | `string[]` | `[]` | No | Multiple source ports/ranges in one rule. Never combined with `sourcePortRange`. |
| `destinationPortRange` | `string` | — | Yes* | Destination port, range (`1024-65535`), or `*` for any. Examples: `22` (SSH), `80` (HTTP), `443` (HTTPS). Exactly one of `destinationPortRange` or `destinationPortRanges` must be set. |
| `destinationPortRanges` | `string[]` | `[]` | Yes* | Multiple destination ports/ranges in one rule (e.g. `["80", "443"]`). Exactly one of `destinationPortRange` or `destinationPortRanges` must be set. |
| `sourceAddressPrefix` | `string` | — | No | Source CIDR, IP, Azure service tag (`VirtualNetwork`, `Internet`), or `*`. Service tags and `*` only work in this singular form. At most one source addressing style may be set; all unset means any. |
| `sourceAddressPrefixes` | `string[]` | `[]` | No | Multiple source CIDRs or IPs. Service tags are not supported in this field. At most one source addressing style may be set. |
| `sourceApplicationSecurityGroupIds` | `string[]` | `[]` | No | Source as application security group membership — plain ARM IDs, up to 10. At most one source addressing style may be set. |
| `destinationAddressPrefix` | `string` | — | No | Destination CIDR, IP, Azure service tag, or `*`. Service tags and `*` only work in this singular form. At most one destination addressing style may be set; all unset means any. |
| `destinationAddressPrefixes` | `string[]` | `[]` | No | Multiple destination CIDRs or IPs. Service tags are not supported in this field. At most one destination addressing style may be set. |
| `destinationApplicationSecurityGroupIds` | `string[]` | `[]` | No | Destination as application security group membership — plain ARM IDs, up to 10. At most one destination addressing style may be set. |

\* Exactly one of `destinationPortRange` or `destinationPortRanges` is required per rule (use `"*"` for any port).

## Examples

### Allow HTTPS Only

A minimal NSG that allows inbound HTTPS and denies everything else (via Azure defaults):

```yaml
apiVersion: azure.planton.dev/v1
kind: AzureNetworkSecurityGroup
metadata:
  name: web-nsg
  annotations:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: dev.AzureNetworkSecurityGroup.web-nsg
spec:
  region: eastus
  resourceGroup:
    value: dev-rg
  name: web-nsg
  securityRules:
    - name: allow-https-inbound
      priority: 100
      direction: INBOUND
      access: ALLOW
      protocol: TCP
      destinationPortRange: "443"
```

### Web Tier with HTTP and HTTPS

An NSG for a web tier that allows both HTTP and HTTPS inbound from the internet:

```yaml
apiVersion: azure.planton.dev/v1
kind: AzureNetworkSecurityGroup
metadata:
  name: web-tier-nsg
  annotations:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: prod.AzureNetworkSecurityGroup.web-tier-nsg
spec:
  region: eastus
  resourceGroup:
    value: prod-rg
  name: web-tier-nsg
  securityRules:
    - name: allow-https-inbound
      priority: 100
      direction: INBOUND
      access: ALLOW
      protocol: TCP
      destinationPortRange: "443"
      sourceAddressPrefix: Internet
    - name: allow-http-inbound
      priority: 200
      direction: INBOUND
      access: ALLOW
      protocol: TCP
      destinationPortRange: "80"
      sourceAddressPrefix: Internet
    - name: deny-all-inbound
      priority: 4096
      direction: INBOUND
      access: DENY
      protocol: ANY
      destinationPortRange: "*"
      description: Explicit deny-all as a safety net
```

### Application Tier with Restricted Sources

An NSG for an application tier that only accepts traffic from the web tier subnet and allows SSH from a bastion host:

```yaml
apiVersion: azure.planton.dev/v1
kind: AzureNetworkSecurityGroup
metadata:
  name: app-tier-nsg
  annotations:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: prod.AzureNetworkSecurityGroup.app-tier-nsg
spec:
  region: eastus
  resourceGroup:
    value: prod-rg
  name: app-tier-nsg
  securityRules:
    - name: allow-web-to-app
      priority: 100
      direction: INBOUND
      access: ALLOW
      protocol: TCP
      destinationPortRange: "8080"
      sourceAddressPrefix: "10.0.1.0/24"
      description: Allow traffic from web tier subnet
    - name: allow-ssh-from-bastion
      priority: 200
      direction: INBOUND
      access: ALLOW
      protocol: TCP
      destinationPortRange: "22"
      sourceAddressPrefix: "10.0.255.4"
      description: Allow SSH from bastion host
    - name: deny-all-inbound
      priority: 4096
      direction: INBOUND
      access: DENY
      protocol: ANY
      destinationPortRange: "*"
```

### Data Tier with Multiple Source Ranges

An NSG for a data tier that allows database traffic from multiple application subnets using plural address prefixes:

```yaml
apiVersion: azure.planton.dev/v1
kind: AzureNetworkSecurityGroup
metadata:
  name: data-tier-nsg
  annotations:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: prod.AzureNetworkSecurityGroup.data-tier-nsg
spec:
  region: westeurope
  resourceGroup:
    value: prod-rg
  name: data-tier-nsg
  securityRules:
    - name: allow-postgres-from-app-subnets
      priority: 100
      direction: INBOUND
      access: ALLOW
      protocol: TCP
      destinationPortRange: "5432"
      sourceAddressPrefixes:
        - "10.0.2.0/24"
        - "10.0.3.0/24"
        - "10.0.4.0/24"
      description: Allow PostgreSQL from all app subnets
    - name: deny-all-inbound
      priority: 4096
      direction: INBOUND
      access: DENY
      protocol: ANY
      destinationPortRange: "*"
    - name: deny-internet-outbound
      priority: 4096
      direction: OUTBOUND
      access: DENY
      protocol: ANY
      destinationPortRange: "*"
      destinationAddressPrefix: Internet
      description: Prevent data tier from reaching the internet
```

### Using Foreign Key References

Reference an Planton-managed resource group instead of hardcoding the name:

```yaml
apiVersion: azure.planton.dev/v1
kind: AzureNetworkSecurityGroup
metadata:
  name: ref-nsg
  annotations:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: prod.AzureNetworkSecurityGroup.ref-nsg
spec:
  region: eastus
  resourceGroup:
    valueFrom:
      kind: AzureResourceGroup
      name: my-rg
      field: status.outputs.resource_group_name
  name: ref-nsg
  securityRules:
    - name: allow-https-inbound
      priority: 100
      direction: INBOUND
      access: ALLOW
      protocol: TCP
      destinationPortRange: "443"
```

## Stack Outputs

After deployment, the following outputs are available in `status.outputs`:

| Output | Type | Description |
|--------|------|-------------|
| `network_security_group_id` | `string` | Azure Resource Manager ID of the Network Security Group. Format: `/subscriptions/{sub}/resourceGroups/{rg}/providers/Microsoft.Network/networkSecurityGroups/{name}`. Referenced by AzureSubnet's `networkSecurityGroupId` to attach the group to a subnet. |
| `network_security_group_name` | `string` | Name of the Network Security Group |

## Related Components

- [AzureResourceGroup](/docs/catalog/azure/azureresourcegroup) — provides the resource group for NSG placement
- [AzureVirtualNetwork](/docs/catalog/azure/azurevirtualnetwork) — provides the virtual network and subnets that NSGs are associated with
- [AzureSubnet](/docs/catalog/azure/azuresubnet) — NSGs are associated with subnets to filter traffic at the subnet level
- [AzureAksCluster](/docs/catalog/azure/azureakscluster) — AKS node pool subnets often require NSGs for controlling cluster traffic
