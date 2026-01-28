# Azure NAT Gateway - Pulumi Examples

## Minimal Configuration

```yaml
apiVersion: azure.openmcf.org/v1
kind: AzureNatGateway
metadata:
  name: dev-nat
spec:
  subnet_id:
    value: "/subscriptions/{sub-id}/resourceGroups/dev-rg/providers/Microsoft.Network/virtualNetworks/dev-vnet/subnets/nodes-subnet"
```

## Production with IP Prefix

```yaml
apiVersion: azure.openmcf.org/v1
kind: AzureNatGateway
metadata:
  name: prod-nat
  org: mycompany
  env: production
spec:
  subnet_id:
    value: "/subscriptions/{sub-id}/resourceGroups/prod-rg/providers/Microsoft.Network/virtualNetworks/prod-vnet/subnets/nodes-subnet"
  idle_timeout_minutes: 10
  public_ip_prefix_length: 28
  tags:
    cost-center: "infrastructure"
```
