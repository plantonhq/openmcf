# Standard Remote Users

This preset deploys the classic remote-access shape: one client address pool with split tunneling -- connected devices reach everything the hub reaches over the tunnel, and the internet locally. Scale unit 1 (unset; 500 concurrent connections) and default hub routing.

## When to Use

- The first point-to-site gateway in a hub for a remote workforce
- Teams whose users need internal resources over VPN but not inspected internet egress

## Key Configuration Choices

- **Split tunneling** -- `internetSecurityEnabled` stays off: cheap, fast, invisible to users; only internal traffic rides the tunnel
- **One /24 pool** -- roughly 250 concurrent client addresses; growing it later updates in place, but it must never overlap anything the hub reaches
- **Default routing** -- no `route` block: ARM associates and propagates to the hub's default route table (any-to-any reachability)

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<your-resource-group-name>` | Resource group for the gateway | `AzureResourceGroup` status outputs (`resource_group_name`), or reference it with valueFrom |
| `<your-virtual-hub-arm-id>` | ARM ID of the virtual hub the gateway deploys into | `AzureVirtualHub` status outputs (`virtual_hub_id`), or reference it with valueFrom |
| `<your-vpn-server-configuration-arm-id>` | ARM ID of the authentication policy | `AzureVpnServerConfiguration` status outputs (`vpn_server_configuration_id`), or reference it with valueFrom |

Replace the example pool (`172.16.201.0/24`) with client space your network plan reserves.
