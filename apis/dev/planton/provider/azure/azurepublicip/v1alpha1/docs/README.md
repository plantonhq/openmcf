# AzurePublicIp: Research & Design Documentation

## 1. What Is an Azure Public IP Address?

An Azure Public IP Address is a static IPv4 or IPv6 address allocated from
Azure's public IP pool. It provides inbound internet connectivity to Azure resources
such as load balancers, application gateways, NAT gateways, VPN gateways, and
virtual machines.

Public IPs are one of Azure's foundational networking primitives. They sit at the
edge of the Azure network and are the entry point for all internet-facing traffic.

### Key Properties

- **SKU**: Standard (Azure's default) or StandardV2, Azure's next-generation SKU
  required for StandardV2 NAT gateway attachment (Basic was retired Sept 2025)
- **Tier**: Regional (the default, correct for virtually everything) or Global,
  a globally-anycast address for cross-region load balancer frontends
- **Allocation**: Static (every current SKU requires static allocation)
- **IP Version**: IPv4 (the default) or IPv6
- **Availability Zones**: zone-redundant, zonal, or non-zonal deployments
- **Prefix Allocation**: an address can be drawn from a reserved Public IP Prefix
  instead of Microsoft's general pool, yielding one contiguous, allowlistable range
- **DNS Integration**: optional domain name label creates a stable FQDN, with an
  optional scope-based reuse policy (hashed labels) as a defense against
  dangling-DNS subdomain takeover
- **Reverse DNS**: an optional reverse FQDN records the PTR name mail servers and
  forward-confirmed-reverse-DNS checks see
- **DDoS Protection**: per-IP protection stance -- inherit from the virtual
  network (default), enabled with a dedicated plan, or explicitly disabled
- **IP Tags**: Azure routing metadata attached to the address itself
  (e.g. `RoutingPreference` for cold-potato vs hot-potato transit)
- **Edge Zones**: the address can be deployed into a metro-local Azure Edge Zone
- **Idle Timeout**: Configurable TCP idle timeout (4-30 minutes)
- **Pricing**: ~$3.65/month per static Standard IPv4 address (as of 2025)

### Basic SKU Retirement

Microsoft retired the Basic SKU for Public IP Addresses on **September 30, 2025**.
All existing Basic IPs were automatically migrated to Standard. New deployments
cannot use Basic. This is why the Basic SKU is deliberately absent from this
Planton component's SKU enum.

Source: [Azure Basic SKU retirement announcement](https://azure.microsoft.com/en-us/updates/upgrade-to-standard-sku-public-ip-addresses-in-azure-by-30-september-2025-basic-sku-will-be-retired/)

## 2. Deployment Landscape

### How People Deploy Public IPs Today

#### Level 0: Azure Portal (Click-Ops)

The Azure Portal provides a GUI for creating Public IPs. Users select the SKU,
allocation method, region, and optional DNS label. This is fine for learning but
creates undocumented infrastructure.

#### Level 1: Azure CLI

```bash
az network public-ip create \
  --name my-pip \
  --resource-group my-rg \
  --location eastus \
  --sku Standard \
  --allocation-method Static \
  --dns-name myapp \
  --zone 1 2 3
```

Simple and scriptable, but lacks state management and drift detection.

#### Level 2: ARM Templates / Bicep

```bicep
resource publicIp 'Microsoft.Network/publicIPAddresses@2023-09-01' = {
  name: 'my-pip'
  location: 'eastus'
  sku: {
    name: 'Standard'
    tier: 'Regional'
  }
  properties: {
    publicIPAllocationMethod: 'Static'
    dnsSettings: {
      domainNameLabel: 'myapp'
    }
  }
  zones: ['1', '2', '3']
}
```

Azure-native IaC with full lifecycle management. Verbose but complete.

#### Level 3: Terraform

```hcl
resource "azurerm_public_ip" "main" {
  name                = "my-pip"
  location            = "eastus"
  resource_group_name = "my-rg"
  allocation_method   = "Static"
  sku                 = "Standard"
  domain_name_label   = "myapp"
  zones               = ["1", "2", "3"]
}
```

The most popular IaC approach for multi-cloud teams. Clean, readable, and
well-supported by the Azure Terraform provider.

#### Level 4: Pulumi

```go
publicIp, _ := network.NewPublicIp(ctx, "my-pip", &network.PublicIpArgs{
    Name:              pulumi.String("my-pip"),
    Location:          pulumi.String("eastus"),
    ResourceGroupName: pulumi.String("my-rg"),
    AllocationMethod:  pulumi.String("Static"),
    Sku:               pulumi.String("Standard"),
    DomainNameLabel:   pulumi.String("myapp"),
    Zones:             pulumi.StringArray{pulumi.String("1"), pulumi.String("2"), pulumi.String("3")},
})
```

Programmatic IaC with type safety and testability.

#### Level 5: Planton (This Component)

```yaml
apiVersion: azure.planton.dev/v1alpha1
kind: AzurePublicIp
metadata:
  name: my-pip
spec:
  region: eastus
  resourceGroup:
    value: my-rg
  name: my-pip
  domainNameLabel: myapp
  zones: ["1", "2", "3"]
```

Declarative, Kubernetes-style API that abstracts Pulumi/Terraform behind a consistent
multi-cloud interface. Enables infra chart composition where Public IPs are referenced
by downstream resources via `StringValueOrRef`.

## 3. 80/20 Analysis: What We Include and What We Skip

### Included

| Feature | Rationale |
|---------|-----------|
| Static allocation | Every current SKU requires it; all production use cases need it |
| SKU (`STANDARD` / `STANDARD_V2`) | Standard is the production default; StandardV2 is required for StandardV2 NAT gateway attachment |
| SKU tier (`REGIONAL` / `GLOBAL`) | Global is the only way to build cross-region load balancer frontends |
| IP version (`IPV4` / `IPV6`) | IPv6 endpoints are a first-class requirement for dual-stack architectures |
| Availability zones | Essential for production resilience |
| Public IP prefix allocation | Draws the address from a reserved, allowlistable contiguous range |
| DNS label + label scope | Stable FQDNs, with hashed-label reuse as a subdomain-takeover defense |
| Reverse FQDN | Reverse-DNS (PTR) for mail servers and FCrDNS checks |
| Idle timeout | Important production tunable for long-lived connections |
| IP tags | Routing metadata such as `RoutingPreference` for transit control |
| DDoS protection mode + plan | Per-IP protection stance, incl. dedicated IP-level protection |
| Edge zones | Metro-local deployments for latency-sensitive workloads |
| Tags | User tags merged over automatic Planton metadata tags |

### Excluded

| Feature | Rationale |
|---------|-----------|
| Basic SKU | Retired by Azure (September 30, 2025); not supported |
| Dynamic allocation | Existed only for the Basic SKU; every current SKU requires static |
| DDoS protection plan creation | Plans are shared, rarely-created governance resources; referenced by ARM ID |
| Consumer association | Load balancers, application gateways, and NAT gateways attach the address from their own specs |

## 4. Downstream Consumers

Public IP addresses are consumed by several Azure resources:

### AzureApplicationGateway
- Uses `public_ip_id` for the frontend IP configuration
- Application Gateway requires a dedicated Standard SKU Public IP

### AzureLoadBalancer
- Uses `public_ip_id` for the frontend IP of public load balancers
- Internal load balancers don't need a Public IP

### AzureNatGateway
- References first-class Public IPs (and Public IP Prefixes) by ARM ID for SNAT,
  keeping egress addresses visible in the resource graph and reusable
- A StandardV2 NAT gateway requires `STANDARD_V2` addresses

### DNS Records
- `ip_address` output can be used to create A records in AzureDnsRecord
- `fqdn` output can be used for CNAME record targets

## 5. Infra Chart Integration

### Enterprise Network Foundation

In the `enterprise-network-foundation` infra chart, Public IPs are Layer 1 resources
created after the resource group:

```
AzureResourceGroup (Layer 0)
├── AzureVirtualNetwork (Layer 1)
│   └── AzureSubnet (Layer 2)
├── AzurePublicIp [gateway] (Layer 1)  <-- THIS RESOURCE
│   └── AzureApplicationGateway (Layer 2) -- references public_ip_id
├── AzurePublicIp [lb] (Layer 1)
│   └── AzureLoadBalancer (Layer 2) -- references public_ip_id
├── AzurePublicIp [nat] (Layer 1)
│   └── AzureNatGateway (Layer 2) -- references public_ip_id
└── AzureLogAnalyticsWorkspace (Layer 1)
```

Each downstream resource references the Public IP via `StringValueOrRef`:

```yaml
publicIpId:
  valueFrom:
    name: gateway-pip
```

## 6. Design Decisions

### Why the SKU Enum Has No Basic Value

Azure retired the Basic SKU on September 30, 2025. There is zero reason to model
a deprecated SKU: it would require validation to prevent selecting it, confuse
users with an option that doesn't work, and create a dead code path. The enum
carries only the SKUs Azure will actually create -- `STANDARD` (the default when
unspecified) and `STANDARD_V2`.

### Why No Allocation Method Field

Every current SKU requires static allocation; dynamic allocation existed only for
the retired Basic SKU, and Azure's API rejects Dynamic with any current SKU.
Including an `allocation_method` field would add a proto enum type with two values
where only one is valid. The IaC modules hardcode Static.

### Why Enum Defaults Defer to Azure

`sku`, `sku_tier`, `ip_version`, `domain_name_label_scope`, and
`ddos_protection_mode` all treat "unspecified" as "let Azure apply its default"
(Standard / Regional / IPv4 / region-unique label / inherit from the network).
A minimal manifest therefore deploys exactly what a minimal `az network
public-ip create` would, and the same manifest deploys identically on both the
Pulumi and Terraform engines.

### Why the DDoS Plan Is a Plain ARM ID

`ddos_protection_plan_id` accepts a raw ARM ID rather than a foreign-key
reference: DDoS protection plans are shared, rarely-created governance
resources managed at the subscription level, not per-workload infrastructure.
Spec validation enforces that the plan is only set when `ddos_protection_mode`
is `ENABLED`.

### Why Include a Domain Name Label Scope

A classic domain name label is region-unique and released when the IP is
deleted -- a dangling CNAME pointing at it becomes a subdomain-takeover vector.
The scope-based reuse policy makes Azure hash the label with the chosen scope
(tenant, subscription, resource group, or no reuse), closing that hole. It is
a cheap field that encodes a meaningful security posture.

### Why Include Idle Timeout

The default idle timeout of 4 minutes is too short for many enterprise workloads.
WebSocket connections, gRPC streams, and database connections through a Public IP
will be terminated after 4 minutes of inactivity. A configurable timeout (up to 30
minutes) is a meaningful production lever that costs nothing to include in the spec.

### Why Include Zones

Availability zones provide resilience against datacenter failures. A zone-redundant
Public IP (`zones: ["1","2","3"]`) survives the loss of an entire availability zone.
This is table stakes for production infrastructure and easy to configure.

## 7. Scope Boundaries

### What This Component Does

- Creates a static Azure Public IP Address on the `STANDARD` or `STANDARD_V2` SKU,
  at the `REGIONAL` or `GLOBAL` tier, as IPv4 or IPv6
- Optionally pins to specific availability zones (zonal or zone-redundant)
- Optionally allocates the address from a reserved Public IP Prefix
- Optionally configures a DNS domain name label (with scope-based reuse) and a
  reverse FQDN
- Configures idle timeout for TCP connection lifecycle
- Optionally sets the per-IP DDoS protection stance and backing plan
- Optionally attaches Azure IP tags (routing metadata) and deploys into an Edge Zone
- Tags the resource with Planton metadata, merged with user tags
- Exports the IP ID, address, FQDN, and name for downstream consumption

### What This Component Does NOT Do

- **NAT Gateway association** -- handled by AzureNatGateway
- **Load Balancer association** -- handled by AzureLoadBalancer
- **Application Gateway association** -- handled by AzureApplicationGateway
- **Public IP Prefix creation** -- handled by AzurePublicIpPrefix; this component
  only references an existing prefix
- **DDoS protection plan creation** -- plans are subscription-level governance
  resources, referenced here by ARM ID
