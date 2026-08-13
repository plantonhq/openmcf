# Local Metro Circuit

This preset creates a LOCAL-tier circuit: connectivity to the Azure regions in the circuit's own metro only, with NO egress fees. When your facility sits near the Azure region you use, this is the cost-efficient shape.

## When to Use

- Your datacenter/colocation is in the same metro as your primary Azure region
- Heavy egress workloads where the no-egress-fee economics dominate
- You do not need reach beyond the local metro's regions

## Key Configuration Choices

- **Reach is the trade** -- LOCAL peers only with the metro's Azure regions; traffic to other regions needs a STANDARD/PREMIUM circuit or another path
- **LOCAL is offered on 1 Gbps and up** at most locations -- check your provider's offering for the site
- **Unlimited pairs naturally** -- LOCAL's pricing already folds egress in; unlimited keeps the bill flat
- **Same lifecycle rules as any circuit** -- billing from creation, bandwidth ratchets upward only, provider handoff via the service key

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<azure-region>` | The ARM metadata region | Your regional deployment strategy |
| `<your-resource-group-name>` | Name of the resource group | Azure portal or `AzureResourceGroup` status outputs |
| `<your-connectivity-provider>` | The provider's exact Azure-listed name | `az network express-route list-service-providers` |
| `<your-peering-location>` | The provider's cross-connect site (in-metro) | Your provider's ExpressRoute order |
