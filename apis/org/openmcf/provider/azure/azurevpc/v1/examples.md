# Azure VPC Examples

## Minimal Configuration

```yaml
apiVersion: azure.openmcf.org/v1
kind: AzureVpc
metadata:
  name: dev-vpc
spec:
  addressSpaceCidr: "10.10.0.0/16"
  nodesSubnetCidr: "10.10.0.0/18"
```

## Development Environment

```yaml
apiVersion: azure.openmcf.org/v1
kind: AzureVpc
metadata:
  name: dev-vpc
  env: development
spec:
  addressSpaceCidr: "10.10.0.0/16"
  nodesSubnetCidr: "10.10.0.0/18"
  tags:
    purpose: "development"
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
  addressSpaceCidr: "10.0.0.0/16"
  nodesSubnetCidr: "10.0.0.0/18"
  isNatGatewayEnabled: true
  tags:
    cost-center: "infrastructure"
    compliance: "soc2"
```

## With Private DNS Zone Links

```yaml
apiVersion: azure.openmcf.org/v1
kind: AzureVpc
metadata:
  name: prod-vpc-with-dns
  org: mycompany
  env: production
spec:
  addressSpaceCidr: "10.0.0.0/16"
  nodesSubnetCidr: "10.0.0.0/18"
  isNatGatewayEnabled: true
  dnsPrivateZoneLinks:
    - "/subscriptions/{sub-id}/resourceGroups/{rg}/providers/Microsoft.Network/privateDnsZones/privatelink.database.windows.net"
    - "/subscriptions/{sub-id}/resourceGroups/{rg}/providers/Microsoft.Network/privateDnsZones/privatelink.blob.core.windows.net"
  tags:
    environment: "production"
```

**Use Cases**:
- Basic VNet for development/testing
- Production VNet with NAT Gateway for reliable outbound connectivity
- VNet with Private DNS zones for Azure PaaS services (SQL, Storage)
- Multi-environment network isolation

**NAT Gateway Benefits**:
- Eliminates SNAT port exhaustion
- Provides predictable egress IPs for firewall allow-listing
- Scales to 1M+ ports with public IP prefixes
- Automatic precedence over other outbound methods

