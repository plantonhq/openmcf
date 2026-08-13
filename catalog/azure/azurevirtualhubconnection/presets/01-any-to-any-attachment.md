# Any-to-Any Attachment

This preset creates the simplest attachment: the spoke joins the hub with ARM's default routing -- associated with and propagating to the hub's built-in default route table, so every connected network reaches every other. The connection is free; transit through the hub is what bills.

## When to Use

- Trusted estates where full mutual reachability is the intent
- The first spokes of a new WAN, before topology hardening begins

## Key Configuration Choices

- **No routing block** -- unset means ARM's default (any-to-any through the default table); add the block only to change topology
- **Internet security off** -- the spoke keeps its own internet egress; flip it on only when a hub firewall handles egress
- **Address spaces must not overlap** -- ARM validates the spoke against the hub and every already-connected network at attach time

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<your-virtual-hub-arm-id>` | ARM ID of the hub | `AzureVirtualHub` status outputs (`virtual_hub_id`), or reference it with valueFrom |
| `<your-virtual-network-arm-id>` | ARM ID of the spoke VNet | `AzureVirtualNetwork` status outputs (`virtual_network_id`), or reference it with valueFrom |
