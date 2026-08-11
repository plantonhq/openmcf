# Performance Routing

This preset deploys the everyday multi-region front: a Performance profile that answers each user's DNS lookup with the lowest-latency healthy endpoint, probed over HTTPS on a real health path.

## When to Use

- Multi-region deployments where users should land on the nearest region
- Any latency-sensitive service fronted by regional endpoints (Azure or external)

## Key Configuration Choices

- **Performance routing** -- Azure's latency map picks the closest endpoint per caller; endpoints added later (AzureTrafficManagerEndpoint) need locations only for external/nested types
- **HTTPS probe on `/healthz`** -- proves the SERVICE, not just the socket; scope `expectedStatusCodeRanges` if your health path answers anything but 200
- **60-second TTL** -- the failover clock: clients keep cached answers this long after Traffic Manager stops handing an endpoint out

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<your-resource-group-name>` | Name of the resource group | Azure portal or `AzureResourceGroup` status outputs |
| `<your-globally-unique-dns-label>` | The trafficmanager.net label -- globally unique across ALL of Azure | Your organization's naming convention (e.g. `contoso-app-prod`) |
