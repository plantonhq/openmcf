# Custom IPsec Policy Tunnel

This preset creates a site-to-site connection with a PINNED IPsec/IKE proposal (DH14 / AES-256 / SHA-256 / PFS2048) -- for on-premises devices whose documentation demands exact algorithms, and for compliance regimes that mandate specific cryptography rather than accepting Azure's negotiated default.

## When to Use

- The device vendor's Azure interop guide specifies exact IKE/IPsec parameters
- Compliance requires documented, fixed cryptography on the tunnel
- Troubleshooting a tunnel stuck in `Connecting` where proposal mismatch is the suspect (pin both ends identically and re-test)

## Key Configuration Choices

- **All six algorithms travel together** -- a partial policy is not a policy (ARM's contract). This preset's set (DHGroup14/AES256/SHA256/PFS2048) is the widely-supported strong baseline
- **GCM pairing rule** -- if you switch IKE encryption to a GCM variant, the IKE integrity value must be the MATCHING GCM value
- **`dpdTimeoutSeconds: 45`** -- Azure's default dead-peer-detection window, made explicit; lower it for faster failover detection on redundant tunnels
- **Basic-SKU gateways reject custom policies** -- the gateway must be VpnGw1AZ or higher
- **`usePolicyBasedTrafficSelectors`** -- add it (with this pinned policy) only when the on-premises device is policy-based; the spec requires the policy when the flag is set

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<azure-region>` | Azure region (must match the gateway) | Your regional deployment strategy |
| `<your-resource-group-name>` | Name of the resource group | Azure portal or `AzureResourceGroup` status outputs |
| `<your-virtual-network-gateway-resource-name>` | Planton metadata name of the `AzureVirtualNetworkGateway` | Your gateway resource |
| `<your-local-network-gateway-resource-name>` | Planton metadata name of the `AzureLocalNetworkGateway` describing this site | Your site description resource |
| `<your-tunnel-psk-secret-name>` | Name of the secret holding the pre-shared key | Your secrets manager |
