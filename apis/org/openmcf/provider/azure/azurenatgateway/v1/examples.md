# Azure NAT Gateway Examples

## Minimal Configuration

```yaml
apiVersion: azure.openmcf.org/v1
kind: AzureNatGateway
metadata:
  name: dev-nat
spec:
  subnetId:
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
  subnetId:
    value: "/subscriptions/{sub-id}/resourceGroups/prod-rg/providers/Microsoft.Network/virtualNetworks/prod-vnet/subnets/nodes-subnet"
  idleTimeoutMinutes: 10
  publicIpPrefixLength: 28  # 16 IPs for scale
  tags:
    cost-center: "infrastructure"
```

## Development Configuration

```yaml
apiVersion: azure.openmcf.org/v1
kind: AzureNatGateway
metadata:
  name: dev-nat
  env: development
spec:
  subnetId:
    value: "/subscriptions/{sub-id}/resourceGroups/dev-rg/providers/Microsoft.Network/virtualNetworks/dev-vnet/subnets/nodes-subnet"
  idleTimeoutMinutes: 4  # Default, faster connection cleanup
```

