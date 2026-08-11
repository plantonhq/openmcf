# Behind a Load Balancer

This preset publishes a service running behind a Standard internal load balancer -- the classic Private Link shape. Consumers in other virtual networks (or tenants) connect through private endpoints using the service's generated alias; their traffic arrives source-NATed through the NAT configuration, so address spaces never meet.

## When to Use

- Publishing an internal API, database tier, or platform service to other teams' VNets without peering
- Cross-subscription or cross-tenant service delivery without opening anything publicly
- Any service already fronted by a Standard internal load balancer

## Key Configuration Choices

- **The NAT subnet's policy flag** -- the subnet must have `privateLinkServiceNetworkPoliciesEnabled: false`; ARM rejects the create otherwise, and nothing catches it offline
- **One NAT configuration suffices** -- it funds ~64k concurrent flows per consumer endpoint; add more (up to 8) only for genuinely high fan-in
- **Visibility defaults closed** -- with no `visibilitySubscriptionIds`, only your own tenant's RBAC can discover the service; add UUIDs (or `"*"`) deliberately
- **The frontend set is fixed at creation** -- moving to a different load balancer is a replace

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<azure-region>` | Azure region (must match the load balancer's) | Your regional deployment strategy |
| `<your-resource-group-name>` | Name of the resource group | Azure portal or `AzureResourceGroup` status outputs |
| `<your-pls-subnet-arm-id>` | The policies-disabled subnet's ARM ID | `AzureSubnet` status outputs (`subnet_id`) |
| `<your-lb-frontend-arm-id>` | The Standard LB frontend's ARM ID | `AzureLoadBalancer` status outputs (`frontend_ip_configuration_ids.<name>`) |
