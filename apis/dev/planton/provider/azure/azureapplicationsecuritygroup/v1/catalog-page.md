# Azure Application Security Group

Creates an Azure Application Security Group -- a named grouping of network interfaces that network security group rules target by name instead of by IP address. Micro-segmentation policy follows the workload as instances scale and change addresses.

## What Gets Created

When you deploy an AzureApplicationSecurityGroup resource, Planton provisions:

- **Application Security Group** — an `azurerm_application_security_group` in the specified region and resource group

The group holds no members. Membership is declared from the member side: `AzureNetworkInterface` and `AzureVirtualMachineScaleSet` list the ASGs they join, and `AzureNetworkSecurityGroup` security rules reference source/destination ASGs.

## Prerequisites

- **Azure credentials** configured via environment variables or Planton provider config
- **A resource group** to create the group in (an `AzureResourceGroup` in composed environments)
- **Network write rights**: `Microsoft.Network/applicationSecurityGroups/write` (Network Contributor, Contributor, or Owner)

## Quick Start

Create a file `application-security-group.yaml`:

```yaml
apiVersion: azure.planton.dev/v1
kind: AzureApplicationSecurityGroup
metadata:
  name: web-tier
  labels:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: prod.AzureApplicationSecurityGroup.web-tier
spec:
  region: eastus
  resourceGroup:
    value: network-rg
  name: web-tier
```

Deploy:

```shell
planton apply -f application-security-group.yaml
```

After deployment, read `status.outputs.application_security_group_id` for the ARM ID to reference from network interfaces and NSG rules.

## Configuration Reference

### Required Fields

| Field | Type | Description | Validation |
|-------|------|-------------|------------|
| `region` | `string` | Azure region; must match every network interface that joins the group. | Required |
| `resourceGroup` | `StringValueOrRef` | Resource group name. Defaults to referencing an `AzureResourceGroup`'s name output. | Required |
| `name` | `string` | Group name, unique within the resource group. | Required, 1-80 chars, Azure naming rules |

### Optional Fields

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `tags` | `map(string)` | `{}` | User tags, merged over Planton-derived tags (user wins on collision). The only field that updates in place. |

## Examples

### Three-Tier Micro-Segmentation

Group workloads by role so security rules express intent, not addresses:

```yaml
apiVersion: azure.planton.dev/v1
kind: AzureApplicationSecurityGroup
metadata:
  name: app-tier
  labels:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: prod.AzureApplicationSecurityGroup.app-tier
spec:
  region: eastus
  resourceGroup:
    valueFrom:
      name: network-rg
  name: app-tier
```

Join a network interface to the group:

```yaml
apiVersion: azure.planton.dev/v1
kind: AzureNetworkInterface
metadata:
  name: app-vm-nic
spec:
  applicationSecurityGroupIds:
    - valueFrom:
        name: app-tier
```

Target the group in an NSG rule:

```yaml
spec:
  securityRules:
    - name: allow-web-to-app
      direction: INBOUND
      access: ALLOW
      protocol: TCP
      priority: 100
      destinationPortRange: "8080"
      sourceApplicationSecurityGroupIds:
        - valueFrom:
            name: web-tier
      destinationApplicationSecurityGroupIds:
        - valueFrom:
            name: app-tier
```

## Stack Outputs

After deployment, the following outputs are available in `status.outputs`:

| Output | Type | Description |
|--------|------|-------------|
| `application_security_group_id` | `string` | Full ARM ID -- referenced by `AzureNetworkInterface`, `AzureVirtualMachineScaleSet`, and `AzureNetworkSecurityGroup` rules |
| `application_security_group_name` | `string` | The group's name as deployed |

## Related Components

- [AzureNetworkInterface](/docs/catalog/azure/network-interface) — joins ASGs via `applicationSecurityGroupIds`
- [AzureNetworkSecurityGroup](/docs/catalog/azure/network-security-group) — targets ASGs in `source`/`destinationApplicationSecurityGroupIds`
- [AzureVirtualMachineScaleSet](/docs/catalog/azure/virtual-machine-scale-set) — declares ASG membership on instance IP configurations
- [AzureResourceGroup](/docs/catalog/azure/azureresourcegroup) — provides the resource group for placement
