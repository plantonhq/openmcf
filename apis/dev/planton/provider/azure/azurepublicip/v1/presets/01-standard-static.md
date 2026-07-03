# Standard Static Public IP

This preset creates a zone-redundant Azure Public IP with the Standard SKU and static allocation. Standard with static allocation is the only supported configuration (Azure retired the Basic SKU in September 2025). Zone redundancy across all three availability zones provides the highest availability for production load balancers, application gateways, and NAT gateways.

## When to Use

- Attaching to Azure Load Balancers, Application Gateways, or NAT Gateways
- Any resource that needs a stable, internet-routable IPv4 address
- Production workloads requiring zone-redundant availability

## Key Configuration Choices

- **Zone-redundant** (`zones: ["1", "2", "3"]`) -- Survives the failure of any single availability zone. Use fewer zones only if the region does not support all three
- **Idle timeout** (`idleTimeoutInMinutes: 4`) -- Azure default. Increase for long-lived connections (WebSocket, gRPC streaming); maximum is 30 minutes
- **SKU left unspecified** -- Azure's default (Standard) applies; set `sku: STANDARD_V2` only when the address will attach to a StandardV2 NAT gateway. Allocation is always static

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<azure-region>` | Azure region (must match the resource this IP attaches to) | Your regional deployment strategy |
| `<your-resource-group-name>` | Name of the resource group | Azure portal or `AzureResourceGroup` status outputs |
