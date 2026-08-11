# Standard with Tunneling

This preset deploys a Standard host shaped for engineering teams: native-client sessions from a local terminal (`az network bastion ssh/rdp/tunnel`), file transfer, IP-based connections across peerings, and headroom at 4 scale units.

## When to Use

- Teams who live in a terminal and want real ssh/scp ergonomics instead of a browser tab
- Hub networks whose Bastion serves peered spokes (IP connect reaches what Bastion cannot enumerate)
- Environments where session concurrency outgrows Basic's fixed capacity

## Key Configuration Choices

- **Tunneling + file copy + IP connect** are Standard/Premium features -- validation rejects them on Basic before any cloud call
- **`scaleUnits: 4`** carries ~80 concurrent sessions; scaling is an in-place update, so start modest and grow with usage
- **Clipboard stays on** (the default) -- disable `copyPasteEnabled` in exfiltration-sensitive environments

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<your-resource-group>` | The AzureResourceGroup the host is created in | The resource-group component's name |
| `<your-bastion-subnet>` | The AzureSubnet whose ARM name is `AzureBastionSubnet` | The subnet component's name |
| `<your-bastion-pip>` | The AzurePublicIp the host binds exclusively | The address component's name |

Standard bills hourly plus per-scale-unit. If session recording or private-only deployment is plausibly ahead, start at Premium -- downgrades replace the host.
