# Azure VPC - Pulumi Examples

## Minimal Configuration

```yaml
apiVersion: azure.openmcf.org/v1
kind: AzureVpc
metadata:
  name: dev-vpc
spec:
  address_space_cidr: "10.10.0.0/16"
  nodes_subnet_cidr: "10.10.0.0/18"
```

## Production with NAT Gateway

```yaml
apiVersion: azure.openmcf.org/v1
kind: AzureVpc
metadata:
  name: prod-vpc
  org: mycompany
  env: production
spec:
  address_space_cidr: "10.0.0.0/16"
  nodes_subnet_cidr: "10.0.0.0/18"
  is_nat_gateway_enabled: true
  tags:
    cost-center: "infrastructure"
```

## With Private DNS Zone Links

```yaml
apiVersion: azure.openmcf.org/v1
kind: AzureVpc
metadata:
  name: prod-vpc-dns
spec:
  address_space_cidr: "10.0.0.0/16"
  nodes_subnet_cidr: "10.0.0.0/18"
  is_nat_gateway_enabled: true
  dns_private_zone_links:
    - "/subscriptions/{sub-id}/resourceGroups/{rg}/providers/Microsoft.Network/privateDnsZones/privatelink.database.windows.net"
```
