# Public-Facing NIC with NIC-Level NSG

This preset creates a network interface whose primary configuration is fronted by a referenced `AzurePublicIp` and filtered by a NIC-level `AzureNetworkSecurityGroup`. It is the shape for a single internet-facing VM -- a bastion host, a demo box, an appliance's outside arm -- where the workload itself owns its public exposure.

## When to Use

- Bastion/jump hosts reached directly over the internet
- Single-VM services that genuinely need their own public address (fleets belong behind a load balancer instead)
- Network virtual appliances' external interfaces

## Key Configuration Choices

- **`publicIpAddressId` reference** -- resolves to the address's `public_ip_id` output. The address is a first-class resource: visible in the graph, allowlistable, and reusable if the NIC is replaced
- **NIC-level NSG** -- an internet-exposed workload should carry its own filtering rather than relying on rules shared with subnet neighbors; when both subnet and NIC NSGs are attached, inbound traffic must pass BOTH
- **Dynamic private address** -- the public address is the stable contact point; the private address rarely needs pinning
- **Consider the alternatives first** -- for anything serving real traffic, an `AzureLoadBalancer` or Application Gateway in front of private NICs is the production posture; this preset is for workloads that are legitimately one machine

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<azure-region>` | Azure region (must match the virtual network's region) | Your regional deployment strategy |
| `<your-resource-group-name>` | Name of the resource group | Azure portal or `AzureResourceGroup` status outputs |
| `<your-subnet-resource-name>` | Planton metadata name of the `AzureSubnet` | Your subnet resource |
| `<your-public-ip-resource-name>` | Planton metadata name of the `AzurePublicIp` | Your public IP resource |
| `<your-nsg-resource-name>` | Planton metadata name of the `AzureNetworkSecurityGroup` | Your NSG resource |
