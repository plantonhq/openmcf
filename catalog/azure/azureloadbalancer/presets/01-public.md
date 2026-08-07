# Public Load Balancer

This preset creates a public (internet-facing) Azure Load Balancer: one public frontend, a `web` backend pool, an HTTP health probe, and TCP rules for ports 80 and 443 with TCP reset enabled. This is the standard configuration for distributing inbound internet traffic across VMs, scale-set instances, or appliances.

## When to Use

- Internet-facing web applications that need Layer 4 load balancing across multiple backend instances
- Public APIs or services that require high availability with health-checked traffic distribution
- The inbound entry point in front of a virtual machine scale set (the scale set's ip_configurations reference this LB's pool ID)

## Key Configuration Choices

- **Public frontend** (`frontendIpConfigurations[0].publicIpAddressId`) -- references a first-class Standard SKU `AzurePublicIp`; the address's zone redundancy comes from that resource
- **One backend pool named `web`** -- membership is expressed from the member side: a NIC ip_configuration or scale-set network profile references `status.outputs.backendPoolIds.web`
- **HTTP health probe** -- GETs `/healthz` every 15 seconds; two consecutive failures remove a backend from rotation
- **TCP reset enabled on rules** -- clients learn immediately when an idle connection is dropped instead of discovering a dead socket on the next write

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<azure-region>` | Azure region (must match backend resources) | Your regional deployment strategy |
| `<your-resource-group-name>` | Name of the resource group | Azure portal or `AzureResourceGroup` status outputs |
| `<your-lb-name>` | Name for the load balancer (unique within resource group) | Your naming convention |
| `<public-ip-resource-id>` | Full ARM resource ID of a Standard SKU public IP | `AzurePublicIp` status outputs |

## Related Presets

- **02-internal** -- private VNet load balancing without internet exposure
- **03-outbound-and-nat** -- explicit SNAT egress and admin port forwarding on a public LB
