---
title: "Standard Private NIC"
description: "This preset creates a network interface with one dynamic IPv4 configuration in a referenced `AzureSubnet` and accelerated networking (SR-IOV) enabled. It is the shape virtually every VM starts from:..."
type: "preset"
rank: "01"
presetSlug: "01-standard"
componentSlug: "network-interface"
componentTitle: "Network Interface"
provider: "azure"
icon: "package"
order: 1
---

# Standard Private NIC

This preset creates a network interface with one dynamic IPv4 configuration in a referenced `AzureSubnet` and accelerated networking (SR-IOV) enabled. It is the shape virtually every VM starts from: a private-only NIC whose address Azure assigns and holds stable, with no public exposure.

## When to Use

- The network attachment for a standard `AzureVirtualMachine` (the VM references this NIC's `network_interface_id` output)
- Workloads behind a load balancer or NAT gateway, which need no public IP of their own
- Any VM on a supported size (most current general-purpose sizes with 2+ vCPUs), which should always take accelerated networking

## Key Configuration Choices

- **Dynamic private address** -- Azure picks a free address from the subnet at creation and keeps it for the NIC's lifetime; the assigned address surfaces in the `private_ip_address` output. Pin a static address only when external configuration (firewall rules, appliance peers) depends on the exact IP
- **`acceleratedNetworkingEnabled: true`** -- SR-IOV bypasses the host's virtual switch for dramatically lower latency; the constraint is the VM size, not the workload, so production NICs default to on
- **No public IP** -- the production posture; front the workload with a load balancer, Application Gateway, or NAT gateway instead. Add a `publicIpAddressId` reference to an `AzurePublicIp` only for genuinely internet-facing single VMs
- **No NIC-level NSG** -- subnet-level filtering (the subnet's `networkSecurityGroupId`) covers the common case; attach an NSG here only when one workload needs rules its subnet neighbors must not share

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<azure-region>` | Azure region (must match the virtual network's region) | Your regional deployment strategy |
| `<your-resource-group-name>` | Name of the resource group | Azure portal or `AzureResourceGroup` status outputs |
| `<your-subnet-resource-name>` | Planton metadata name of the `AzureSubnet` | Your subnet resource |
