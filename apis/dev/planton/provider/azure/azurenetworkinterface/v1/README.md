# Overview

The **Azure Network Interface API Resource** provides a consistent and standardized interface for deploying and managing Azure Network Interfaces (NICs) within our infrastructure. The NIC is the attachment point that gives a virtual machine its presence in a subnet -- a first-class resource in Azure's own model, where a VM references its NICs rather than containing them.

## Purpose

We developed this API resource to make the network side of virtual machines explicit and composable. Modeling the NIC as its own resource -- exactly as Azure does -- enables:

- **Honest Composition**: An AzureVirtualMachine consumes this NIC's `network_interface_id` output; a VM can carry several NICs (management + data planes, appliance arms)
- **Per-Workload Security**: A NIC-level network security group filters one workload's traffic specifically, layering with the subnet-level NSG
- **Member-Side Load-Balancer Membership**: Each IP configuration declares the load-balancer backend pools it joins and the inbound NAT rules it completes -- Azure's own attachment model, referencing the load balancer's name-keyed ID outputs
- **Multi-Address Flexibility**: Multiple IP configurations serve dual-stack (IPv4 + IPv6) and multi-IP scenarios on a single interface
- **Appliance Support**: Static addressing, IP forwarding, and preview NVA acceleration cover network-virtual-appliance patterns

## Key Features

- **Consistent Interface**: Aligns with our existing APIs for deploying cloud infrastructure across multiple providers
- **Composition by Reference**: Each IP configuration deploys into a referenced AzureSubnet and may front a referenced AzurePublicIp; the NIC-level NSG references an AzureNetworkSecurityGroup
- **ARM-Faithful Validation**: Spec-level rules enforce ARM's contracts up front -- when multiple IP configurations are declared, the first must be primary; static allocation requires a pinned address; auxiliary mode and SKU are set together
- **Association Resources**: NSG attachment, application-security-group memberships, load-balancer pool and NAT-rule memberships, and Application Gateway pool memberships are all realized as separate association resources (Azure's own model), so joining and leaving never touches the NIC itself
- **Accelerated Networking**: SR-IOV support for dramatically lower latency on every supported VM size
- **Appliance Controls**: IP forwarding for network virtual appliances and preview auxiliary NVA acceleration for connection-rate-bound workloads

## Use Cases

- **Virtual Machine Attachment**: The standard network attachment for every AzureVirtualMachine -- private-only, dynamically addressed, accelerated
- **Load-Balanced Backends**: NICs that join an AzureLoadBalancer's backend pools (and complete its single-target inbound NAT rules) by referencing the load balancer's name-keyed outputs, e.g. `status.outputs.backend_pool_ids.web`
- **Public-Facing Single VMs**: Bastion hosts and appliance outside arms fronted by a referenced public IP with NIC-level filtering
- **Network Virtual Appliances**: Firewalls and routers with static next-hop addresses and IP forwarding enabled
- **Dual-Stack Workloads**: IPv4 and IPv6 configurations side by side on one interface
- **Multi-IP Servers**: Many addresses on one NIC for per-site TLS or NAT pools
- **Workload-Group Security**: Application security group memberships so NSG rules target groups instead of IP ranges

## Future Enhancements

Future updates will include:

- **First-Class Application Security Groups**: ASG memberships currently take plain ARM IDs; a dedicated ASG resource kind will make them referenceable
- **Referenceable Application Gateway Pools**: Application Gateway pool memberships currently take plain ARM IDs; they become referenceable once the Application Gateway exports per-pool IDs
- **Monitoring Integration**: Built-in Azure Monitor metrics for NIC throughput and packet drops
- **Comprehensive Documentation**: Expanded troubleshooting guides for accelerated-networking eligibility and appliance routing
