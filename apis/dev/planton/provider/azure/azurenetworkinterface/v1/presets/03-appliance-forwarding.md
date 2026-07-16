# Network Virtual Appliance NIC (Forwarding + Static IP)

This preset creates the inside interface of a network virtual appliance -- a firewall, router, or SD-WAN box that routes other workloads' traffic. It pins a static private address (the next-hop IP that route tables reference) and enables IP forwarding (without which Azure silently drops traffic not addressed to the NIC itself).

## When to Use

- Firewall/router/SD-WAN VMs that user-defined routes (an `AzureRouteTable`'s `VirtualAppliance` next hop) steer traffic through
- Egress-inspection appliances in a hub network
- Any VM that must accept and forward packets addressed to other destinations

## Key Configuration Choices

- **`privateIpAllocation: STATIC` + `privateIpAddress`** -- route tables forward to the appliance by exact IP (`nextHopInIpAddress`, referencing this NIC's `private_ip_address` output with an explicit kind, or carrying the static IP as a literal); a dynamic address would break every route pointing at it if the NIC were ever recreated
- **`ipForwardingEnabled: true`** -- the load-bearing flag: an appliance NIC without forwarding blackholes every routed packet, and the failure is silent
- **Accelerated networking** -- appliances are packets-per-second workloads; SR-IOV matters more here than anywhere
- **Auxiliary NVA acceleration** (`auxiliaryMode`/`auxiliarySku`) -- a preview feature for connection-rate-bound appliances on enrolled subscriptions; leave unset otherwise
- **Multi-NIC appliances** -- a real appliance usually carries this inside NIC plus an outside NIC (see the public-facing preset); the VM references both

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<azure-region>` | Azure region (must match the virtual network's region) | Your regional deployment strategy |
| `<your-resource-group-name>` | Name of the resource group | Azure portal or `AzureResourceGroup` status outputs |
| `<your-subnet-resource-name>` | Planton metadata name of the `AzureSubnet` (a dedicated appliance subnet) | Your subnet resource |
| `<appliance-static-ip>` | A free address inside the subnet's range, e.g. `10.0.254.4` | Your IP plan; route tables will reference it |
