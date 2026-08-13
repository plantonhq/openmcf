# Secured Hub (Routing Intent)

This preset creates a secured hub: routing intent steers both Internet-bound and private traffic through an Azure Firewall deployed in this hub. Every spoke and branch transiting the hub is inspected -- the Virtual WAN version of a centralized egress/inspection point.

## When to Use

- Regulated estates where all inter-network traffic must pass inspection
- Centralized egress: spokes have no public IPs and reach the internet only through the hub firewall

## Key Configuration Choices

- **The firewall must live in THIS hub** -- routing intent rejects a next hop outside the hub; deploy the `AzureFirewall` into the hub first and reference its ARM ID (`status.outputs.firewall_id`)
- **Routing intent takes over the hub's routing policy** -- per-connection route-table customization and routing intent are mutually exclusive on ARM's side; pick one model per hub
- **Two policies, one appliance** -- the common shape is an Internet policy and a PrivateTraffic policy at the same firewall; drop the Internet policy to inspect only east-west traffic

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<azure-region>` | The hub's region | Your regional deployment strategy |
| `<your-resource-group-name>` | Name of the resource group | Azure portal or `AzureResourceGroup` status outputs |
| `<your-virtual-wan-arm-id>` | ARM ID of the WAN this hub belongs to | `AzureVirtualWan` status outputs (`virtual_wan_id`), or reference it with valueFrom |
| `<your-hub-firewall-arm-id>` | ARM ID of the Azure Firewall deployed in this hub | `AzureFirewall` status outputs (`firewall_id`), or reference it with valueFrom |
