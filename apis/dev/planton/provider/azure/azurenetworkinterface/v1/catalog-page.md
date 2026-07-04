# Azure Network Interface

Deploys an Azure Network Interface (NIC) — the attachment point that gives a virtual machine its presence in a subnet. The NIC is a first-class resource in Azure's own model: a VM does not contain its network configuration, it references one or more NICs. The catalog mirrors that honestly — an `AzureVirtualMachine` consumes this NIC's `network_interface_id` output, each IP configuration deploys into a referenced `AzureSubnet` and may front a referenced `AzurePublicIp`, a NIC-level network security group attaches here as the per-workload complement to subnet-level filtering, and load-balancer membership is expressed here, from the member side: each IP configuration lists the backend pools it joins and the inbound NAT rules it completes, referencing the load balancer's name-keyed ID outputs.

## What Gets Created

When you deploy an AzureNetworkInterface resource, Planton provisions:

- **Network Interface** — a `network.NetworkInterface` in the specified region and resource group, carrying the declared IP configurations, DNS settings, accelerated networking, IP forwarding, and (for appliances on enrolled subscriptions) preview auxiliary NVA acceleration
- **NSG Association** — a `NetworkInterfaceSecurityGroupAssociation` when `networkSecurityGroupId` is set, attaching the NIC-level network security group as its own ARM operation so filtering can change without touching the NIC
- **ASG Associations** — a `NetworkInterfaceApplicationSecurityGroupAssociation` for each entry in `applicationSecurityGroupIds`, joining the NIC to workload groups that NSG rules can target
- **Load-Balancer Pool Associations** — a `NetworkInterfaceBackendAddressPoolAssociation` for each entry in an IP configuration's `loadBalancerBackendAddressPoolIds`, joining the configuration to a load balancer's backend pool from the member side (Azure's own model)
- **NAT-Rule Associations** — a `NetworkInterfaceNatRuleAssociation` for each entry in `loadBalancerInboundNatRuleIds`, completing a load balancer's single-target inbound NAT rule by picking this NIC as the receiving instance
- **Application Gateway Pool Associations** — a `NetworkInterfaceApplicationGatewayBackendAddressPoolAssociation` for each entry in `applicationGatewayBackendAddressPoolIds`
- **Azure Tags** — resource metadata tags applied to the NIC for tracking and governance

Nothing else is created here. The NIC does not create its subnet, public IP, or NSG (reference existing ones), and it does not attach itself to a VM (the `AzureVirtualMachine` declares the attachment via its `network_interface_ids`, matching Azure's model where the VM references its NICs).

## Prerequisites

- **Azure credentials** configured via environment variables or Planton provider config
- **An Azure Resource Group** where the NIC will be created (can reference an AzureResourceGroup resource)
- **An AzureSubnet** for each IPv4 configuration's private address — all of a NIC's configurations must address subnets of the same virtual network, in the same region as the NIC
- **Optionally**, an AzurePublicIp to front a configuration and an AzureNetworkSecurityGroup for NIC-level filtering

## Quick Start

Create a file `nic.yaml`:

```yaml
apiVersion: azure.planton.dev/v1
kind: AzureNetworkInterface
metadata:
  name: my-nic
  labels:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: dev.AzureNetworkInterface.my-nic
spec:
  region: eastus
  resourceGroup:
    value: my-rg
  name: app-nic
  ipConfigurations:
    - name: primary
      subnetId:
        valueFrom:
          name: my-subnet
  acceleratedNetworkingEnabled: true
```

Deploy:

```shell
planton apply -f nic.yaml
```

This creates a private-only NIC with one dynamic IPv4 configuration in the referenced subnet — Azure picks a free address and holds it stable for the NIC's lifetime — with accelerated networking enabled (right for every supported VM size). To attach it to a VM, reference this NIC's `status.outputs.network_interface_id` from the `AzureVirtualMachine`'s `networkInterfaceIds`.

## Configuration Reference

### Required Fields

| Field | Type | Description | Validation |
|-------|------|-------------|------------|
| `region` | `string` | Azure region for the NIC (e.g., `eastus`). Must match the region of the virtual network it deploys into and of the VM that will attach it. Changing it replaces the NIC. | Required, minimum length 1 |
| `resourceGroup` | `StringValueOrRef` | Azure Resource Group name. Can reference an AzureResourceGroup resource via `valueFrom`. | Required |
| `name` | `string` | Name of the NIC, unique within the resource group. Changing it replaces the NIC — and detaches it from any VM. | Required, 1-80 chars, Azure naming rules |
| `ipConfigurations` | `object[]` | The NIC's IP configurations — most NICs carry exactly one; multiple serve dual-stack and multi-IP scenarios. When more than one is declared, the first must be marked `primary` (ARM's contract, spec-enforced). | Required, at least 1 |

### IP Configuration Fields

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `name` | `string` | — | A label for this configuration, unique within the NIC (e.g., `primary`, `ipv6`). Required. |
| `subnetId` | `StringValueOrRef` | — | The subnet the private address lives in. Required for IPv4 configurations (ARM's contract); IPv6 configurations inherit the NIC's subnet placement. Defaults to referencing an `AzureSubnet`'s `subnet_id` output. |
| `privateIpAllocation` | `enum` | `DYNAMIC` | `DYNAMIC` lets Azure pick a free address from the subnet (stable for the NIC's lifetime — right for virtually all workloads); `STATIC` pins a specific address for appliances and servers whose IP is configuration elsewhere. |
| `privateIpAddress` | `string` | — | For `STATIC` allocation: the exact address to pin, inside the subnet's range and unassigned. Forbidden for `DYNAMIC` (spec-enforced). |
| `privateIpVersion` | `enum` | `IPV4` | The address family. A dual-stack NIC carries an `IPV4` and an `IPV6` configuration side by side (in a dual-stack subnet). |
| `publicIpAddressId` | `StringValueOrRef` | — | The public IP fronting this configuration — a first-class `AzurePublicIp` so the address is visible in the resource graph, allowlistable, and reusable. Omit for private-only NICs, the production shape behind a load balancer or NAT gateway. |
| `primary` | `bool` | `false` | Whether this is the NIC's primary configuration. With a single configuration ARM treats it as primary automatically; with multiple, the first must be marked (spec-enforced). |
| `gatewayLoadBalancerFrontendIpConfigurationId` | `string` | — | The frontend of a Gateway-SKU load balancer that chains this NIC into a gateway appliance path. A niche service-chaining seam; plain ARM ID. |
| `loadBalancerBackendAddressPoolIds` | `StringValueOrRef[]` | `[]` | Load-balancer backend pools this configuration joins — membership is expressed from the member side in Azure's model. Reference a pool through the load balancer's name-keyed map output, e.g. `valueFrom` fieldPath `status.outputs.backend_pool_ids.web`. Each membership is realized as its own association resource, so joining and leaving pools never touches the NIC. |
| `loadBalancerInboundNatRuleIds` | `StringValueOrRef[]` | `[]` | Single-target inbound NAT rules this configuration completes — the load balancer declares the port forward, the NIC-side association picks the receiving instance. Reference a rule through the load balancer's name-keyed map output, e.g. `valueFrom` fieldPath `status.outputs.nat_rule_ids.ssh-admin`. Realized as association resources. |
| `applicationGatewayBackendAddressPoolIds` | `string[]` | `[]` | Application Gateway backend pools this configuration joins, by plain ARM ID (the Application Gateway does not export per-pool IDs yet). Realized as association resources. |

### Optional Fields

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `dnsServers` | `string[]` | `[]` | DNS servers overriding the virtual network's DNS for this NIC only. Rarely set — prefer configuring DNS on the virtual network; this exists for appliances that need different resolution than their network. |
| `internalDnsNameLabel` | `string` | — | The DNS label other VMs in the same virtual network can resolve this NIC's private IP by (the full name becomes `{label}.{vnet-internal-suffix}`, surfaced in the outputs). Leave unset for IP-only addressing. |
| `acceleratedNetworkingEnabled` | `bool` | `false` | SR-IOV: the NIC bypasses the host's virtual switch for dramatically lower latency and higher packets-per-second. Azure defaults to false, but production NICs on supported VM sizes (most current general-purpose sizes with 2+ vCPUs) should enable it — the constraint is the VM size, not the workload. |
| `ipForwardingEnabled` | `bool` | `false` | Whether the NIC forwards traffic not addressed to it. Enable ONLY on network virtual appliances — a route table pointing at an appliance's IP silently blackholes unless the appliance NIC forwards. |
| `auxiliaryMode` | `enum` | unset | NVA-acceleration auxiliary mode (`ACCELERATED_CONNECTIONS`, `FLOATING`, `MAX_CONNECTIONS`) — a preview Azure feature the subscription must be enrolled in. Unset sends nothing, correct for every non-appliance NIC. Must be set together with `auxiliarySku` (spec-enforced). |
| `auxiliarySku` | `enum` | unset | The auxiliary SKU sizing the NVA acceleration (`A1` smallest to `A8` largest). Must be set together with `auxiliaryMode`. |
| `edgeZone` | `string` | — | Azure Edge Zone pinning for edge-computing workloads. Leave unset for regular regional deployment. Fixed at creation. |
| `networkSecurityGroupId` | `StringValueOrRef` | — | The NSG filtering THIS NIC's traffic — the per-workload complement to the subnet-level NSG. When both are attached, inbound traffic must pass the subnet NSG then the NIC NSG. Omit to rely on subnet-level filtering alone (the common case). Defaults to referencing an `AzureNetworkSecurityGroup`'s `network_security_group_id` output. |
| `applicationSecurityGroupIds` | `string[]` | `[]` | Application security groups this NIC joins, by plain ARM ID, so NSG rules can target workload groups ("web-servers", "databases") instead of IP ranges. |
| `tags` | `map<string, string>` | `{}` | Additional tags applied to the NIC, merged over Planton-derived tags (user wins on collision). |

## Examples

### Standard Private NIC

The shape virtually every VM starts from — one dynamic IPv4 configuration, no public exposure, accelerated networking on:

```yaml
apiVersion: azure.planton.dev/v1
kind: AzureNetworkInterface
metadata:
  name: app-nic
  labels:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: dev.AzureNetworkInterface.app-nic
spec:
  region: eastus
  resourceGroup:
    value: dev-rg
  name: app-nic
  ipConfigurations:
    - name: primary
      subnetId:
        valueFrom:
          name: app-subnet
  acceleratedNetworkingEnabled: true
```

### Public-Facing NIC with NIC-Level NSG

A single internet-facing VM — a bastion host or an appliance's outside arm — fronted by a referenced public IP and carrying its own filtering:

```yaml
apiVersion: azure.planton.dev/v1
kind: AzureNetworkInterface
metadata:
  name: bastion-nic
  labels:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: prod.AzureNetworkInterface.bastion-nic
spec:
  region: eastus
  resourceGroup:
    value: prod-rg
  name: bastion-nic
  ipConfigurations:
    - name: primary
      subnetId:
        valueFrom:
          name: bastion-subnet
      publicIpAddressId:
        valueFrom:
          name: bastion-ip
  acceleratedNetworkingEnabled: true
  networkSecurityGroupId:
    valueFrom:
      name: bastion-nsg
```

### Load-Balanced Backend NIC

A web-tier NIC that joins an `AzureLoadBalancer`'s backend pool and completes its single-target SSH NAT rule — membership expressed from the member side, referencing the load balancer's name-keyed outputs:

```yaml
apiVersion: azure.planton.dev/v1
kind: AzureNetworkInterface
metadata:
  name: web-1-nic
  labels:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: prod.AzureNetworkInterface.web-1-nic
spec:
  region: eastus
  resourceGroup:
    value: prod-rg
  name: web-1-nic
  ipConfigurations:
    - name: primary
      subnetId:
        valueFrom:
          name: web-subnet
      loadBalancerBackendAddressPoolIds:
        - valueFrom:
            name: web-lb
            fieldPath: status.outputs.backend_pool_ids.web
      loadBalancerInboundNatRuleIds:
        - valueFrom:
            name: web-lb
            fieldPath: status.outputs.nat_rule_ids.ssh-web-1
  acceleratedNetworkingEnabled: true
```

Each membership is its own association resource: adding this NIC to another pool, or removing it, is a spec-list change that never touches the NIC itself.

### Network Virtual Appliance NIC

The inside interface of a firewall or router that routes other workloads' traffic — static next-hop address, IP forwarding on:

```yaml
apiVersion: azure.planton.dev/v1
kind: AzureNetworkInterface
metadata:
  name: fw-inside-nic
  labels:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: prod.AzureNetworkInterface.fw-inside-nic
spec:
  region: eastus
  resourceGroup:
    value: hub-rg
  name: fw-inside-nic
  ipConfigurations:
    - name: primary
      subnetId:
        valueFrom:
          name: fw-inside-subnet
      privateIpAllocation: STATIC
      privateIpAddress: 10.0.254.4
  ipForwardingEnabled: true
  acceleratedNetworkingEnabled: true
```

### Dual-Stack NIC

An IPv4 and an IPv6 configuration side by side (in a dual-stack subnet). With multiple configurations, the first must be marked primary:

```yaml
apiVersion: azure.planton.dev/v1
kind: AzureNetworkInterface
metadata:
  name: dual-stack-nic
  labels:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: prod.AzureNetworkInterface.dual-stack-nic
spec:
  region: eastus
  resourceGroup:
    value: prod-rg
  name: dual-stack-nic
  ipConfigurations:
    - name: primary
      subnetId:
        valueFrom:
          name: dual-stack-subnet
      primary: true
    - name: ipv6
      privateIpVersion: IPV6
  acceleratedNetworkingEnabled: true
```

### NIC Attached to a Virtual Machine

The composition the NIC exists for — the VM references the NIC, never contains it:

```yaml
apiVersion: azure.planton.dev/v1
kind: AzureNetworkInterface
metadata:
  name: web-nic
spec:
  region: eastus
  resourceGroup:
    valueFrom:
      name: web-rg
  name: web-nic
  ipConfigurations:
    - name: primary
      subnetId:
        valueFrom:
          name: web-subnet
  acceleratedNetworkingEnabled: true
  internalDnsNameLabel: web-1
  tags:
    tier: web
---
apiVersion: azure.planton.dev/v1
kind: AzureVirtualMachine
metadata:
  name: web-vm
spec:
  # ... region, resourceGroup, name, size, osProfile, osDisk ...
  networkInterfaceIds:
    - valueFrom:
        name: web-nic
```

## Stack Outputs

After deployment, the following outputs are available in `status.outputs`:

| Output | Type | Description |
|--------|------|-------------|
| `network_interface_id` | `string` | Azure Resource Manager ID of the NIC — the primary output; `AzureVirtualMachine`'s `networkInterfaceIds` references it to attach the NIC to a VM |
| `network_interface_name` | `string` | The NIC's name as deployed |
| `private_ip_address` | `string` | The primary configuration's private IP address — what backends, firewall rules, and DNS records key on |
| `private_ip_addresses` | `string[]` | The private IP addresses of ALL configurations, in configuration order (multi-IP and dual-stack NICs carry more than one) |
| `mac_address` | `string` | The NIC's MAC address, populated once attached to a running VM (empty until then) — what license servers and appliance registrations key on |
| `internal_domain_name_suffix` | `string` | The DNS suffix completing `internalDnsNameLabel` into a resolvable VNet-internal FQDN |

## Related Components

- [AzureResourceGroup](/docs/catalog/azure/azureresourcegroup) -- provides the resource group for NIC placement
- [AzureVirtualNetwork](/docs/catalog/azure/azurevirtualnetwork) -- the network whose subnets the NIC's configurations address
- [AzureSubnet](/docs/catalog/azure/azuresubnet) -- each IP configuration's private address lives in a referenced subnet
- [AzurePublicIp](/docs/catalog/azure/azurepublicip) -- the first-class public address that can front an IP configuration
- [AzureNetworkSecurityGroup](/docs/catalog/azure/azurenetworksecuritygroup) -- attaches at the NIC level via `networkSecurityGroupId` for per-workload filtering, layering with subnet-level NSGs
- [AzureLoadBalancer](/docs/catalog/azure/azureloadbalancer) -- the NIC joins its backend pools and completes its inbound NAT rules by referencing the load balancer's name-keyed `backend_pool_ids` and `nat_rule_ids` outputs
- [AzureVirtualMachine](/docs/catalog/azure/azurevirtualmachine) -- consumes the NIC's `network_interface_id` output; a VM references one or more NICs
