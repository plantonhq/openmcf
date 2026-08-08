# Site-to-Site Tunnel

This preset creates the classic site-to-site IPsec connection: the tunnel joining a virtual network gateway to one described on-premises site, with the pre-shared key referenced from a secret. Azure's default IPsec proposal set negotiates with any modern VPN device -- no pinned policy needed.

## When to Use

- Connecting a datacenter or branch office to Azure (one connection per site)
- Completing the three-resource site-to-site story: site description + gateway + this tunnel
- Any tunnel where the on-premises device negotiates standard IKEv2 proposals

## Key Configuration Choices

- **`type: IPSEC`** -- requires `localNetworkGatewayId` (spec-enforced); the site description carries the device endpoint and reachable prefixes
- **Secret-referenced `sharedKey`** -- both tunnel ends must agree on it; never paste keys into manifests. Alternatively OMIT the key entirely and Azure generates a strong one (read it back with `az network vpn-connection shared-key show` and configure the device)
- **No `ipsecPolicy`** -- Azure's default proposal set is broad; pin algorithms only when the device documentation demands it (see the Custom IPsec Policy preset)
- **Provisioned is not connected** -- the deployment succeeds when ARM accepts the tunnel; `Connected` requires the on-premises device to negotiate

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<azure-region>` | Azure region (must match the gateway) | Your regional deployment strategy |
| `<your-resource-group-name>` | Name of the resource group | Azure portal or `AzureResourceGroup` status outputs |
| `<your-virtual-network-gateway-resource-name>` | Planton metadata name of the `AzureVirtualNetworkGateway` | Your gateway resource |
| `<your-local-network-gateway-resource-name>` | Planton metadata name of the `AzureLocalNetworkGateway` describing this site | Your site description resource |
| `<your-tunnel-psk-secret-name>` | Name of the secret holding the pre-shared key | Your secrets manager |
