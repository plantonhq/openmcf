# Private Peering

This preset configures private peering -- the routing configuration virtually every ExpressRoute deployment needs. Your VNets' private address space flows over the circuit, and an EXPRESS_ROUTE-type virtual network gateway connects VNets to it.

## When to Use

- Connecting VNets to on-premises over the circuit (the standard hybrid path)
- The middle link of the chain: circuit → this peering → ExpressRoute gateway → gateway connection

## Key Configuration Choices

- **The circuit must be PROVISIONED** -- ARM rejects peering configuration while the provider handoff is incomplete; watch the circuit's `service_provider_provisioning_state` output first
- **Addressing comes from the provider's handoff document** -- the VLAN id and both /30s are configured identically on the provider's side; your router takes each /30's first usable address
- **`peerAsn` can stay 0** -- Azure records the ASN your router presents; declare it only when the provider requires it up front
- **One private peering per circuit** -- the type is the ARM identity; a "second" one would overwrite this one

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<your-resource-group-name>` | The circuit's resource group | Azure portal or `AzureResourceGroup` status outputs |
| `<your-circuit-name>` | The parent circuit's name | `AzureExpressRouteCircuit` status outputs (`express_route_circuit_name`) |
| `<your-primary-slash-30>` / `<your-secondary-slash-30>` | The BGP session /30s, one per physical link | Your provider's handoff document |
