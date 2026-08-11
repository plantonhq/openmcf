# Standard Hub Gateway

This preset creates the hub's branch on-ramp with sensible defaults: one scale unit (500 Mbps aggregate), Microsoft-backbone routing preference, and Azure's default BGP settings (ASN 65515). **The gateway bills from creation and takes 30-45 minutes to create** -- deploy it once per hub and grow it in place.

## When to Use

- The first VPN gateway of a new Virtual WAN hub
- Branch estates whose aggregate throughput fits within one scale unit

## Key Configuration Choices

- **Scale unit 1** -- capacity updates in place; start small and grow with measured demand rather than paying for headroom
- **Defaults everywhere else** -- routing preference and BGP asn/peer_weight are fixed at creation; the defaults are right unless a specific branch demands otherwise
- **Region must match the hub's** -- the gateway deploys INTO the hub, and ARM allows one VPN gateway per hub

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<your-resource-group-name>` | Resource group for the gateway | `AzureResourceGroup` status outputs (`resource_group_name`), or reference it with valueFrom |
| `<your-virtual-hub-arm-id>` | ARM ID of the hub | `AzureVirtualHub` status outputs (`virtual_hub_id`), or reference it with valueFrom |
