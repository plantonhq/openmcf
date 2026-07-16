# Azure Public IP

Deploys an Azure Public IP Address -- a static, internet-routable address -- in a specified region and resource group. The component covers the full azurerm surface: SKU (`STANDARD` / `STANDARD_V2`) and tier (`REGIONAL` / `GLOBAL`) selection, IPv4 or IPv6, availability zones, allocation from a reserved Public IP Prefix, DNS labels with scope-based reuse, reverse DNS, per-IP DDoS protection, Azure IP tags, edge zones, and idle timeout tuning. Public IPs created by this component are referenced by downstream resources such as load balancers, application gateways, and NAT gateways via their resource ID.

## What Gets Created

When you deploy an AzurePublicIp resource, Planton provisions:

- **Public IP Address** — an `azurerm_public_ip` / `network.PublicIp` resource in the specified region and resource group, always with static allocation (dynamic allocation existed only for the Basic SKU, retired September 2025; every current SKU requires static)
- **DNS A Record** — when `domainNameLabel` is set, Azure creates an A record at `{label}.{region}.cloudapp.azure.com` pointing to the allocated IP
- **Azure Tags** — resource metadata tags applied to the Public IP for tracking and governance, merged with any user-supplied `tags`

## Prerequisites

- **Azure credentials** configured via environment variables or Planton provider config
- **An Azure Resource Group** where the Public IP will be created (can reference an AzureResourceGroup resource)
- **Region selection** — the Public IP must be in the same region as the resource it will be attached to (load balancer, application gateway, NAT gateway, etc.)
- **An AzurePublicIpPrefix** (optional) — only when allocating the address from a reserved prefix via `publicIpPrefixId`

## Quick Start

Create a file `publicip.yaml`:

```yaml
apiVersion: azure.planton.dev/v1
kind: AzurePublicIp
metadata:
  name: my-pip
  annotations:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: dev.AzurePublicIp.my-pip
spec:
  region: eastus
  resourceGroup:
    value: my-rg
  name: my-pip
```

Deploy:

```shell
planton apply -f publicip.yaml
```

This creates a static public IP on Azure's defaults: Standard SKU, Regional tier, IPv4, a 4-minute idle timeout, no zone preference, and no DNS label.

## Configuration Reference

### Required Fields

| Field | Type | Description | Validation |
|-------|------|-------------|------------|
| `region` | `string` | Azure region for the Public IP (e.g., `eastus`, `westeurope`). Must match the region of the resource it will be attached to. | Required, minimum length 1 |
| `resourceGroup` | `StringValueOrRef` | Azure Resource Group name. Can reference an AzureResourceGroup resource via `valueFrom`. | Required |
| `name` | `string` | Name of the Public IP resource. Must be unique within the resource group. Allowed characters: alphanumeric, underscores, hyphens, and periods. Must start with a letter or number and end with a letter, number, or underscore. | Required, 1–80 characters |

### Optional Fields

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `sku` | `enum` | `STANDARD` | `STANDARD` or `STANDARD_V2` (Azure's next-generation SKU, required to attach the address to a StandardV2 NAT gateway; not valid with the `GLOBAL` tier). Fixed at creation. |
| `skuTier` | `enum` | `REGIONAL` | `REGIONAL` or `GLOBAL`. `GLOBAL` exists solely for cross-region load balancer frontends and requires the `STANDARD` SKU. Fixed at creation. |
| `ipVersion` | `enum` | `IPV4` | `IPV4` or `IPV6`. Fixed at creation. |
| `zones` | `string[]` | `[]` | Availability zones for the Public IP. Valid values: `"1"`, `"2"`, `"3"`. Use `["1", "2", "3"]` for zone-redundant (recommended for production), `["1"]` for zonal, or omit for non-zonal. Zone support depends on the Azure region. Fixed at creation. |
| `publicIpPrefixId` | `StringValueOrRef` | -- | ARM ID of a Public IP Prefix to allocate the address from instead of Microsoft's general pool. Can reference an AzurePublicIpPrefix resource via `valueFrom`. Fixed at creation. |
| `domainNameLabel` | `string` | `""` | DNS label that creates an A record at `{label}.{region}.cloudapp.azure.com`. Must start with a lowercase letter, end with a letter or digit, contain only lowercase letters, digits, and hyphens, and be 3–63 characters. Unique within the region unless a `domainNameLabelScope` is set. |
| `domainNameLabelScope` | `enum` | -- | Scope-based reuse policy for the DNS label: `TENANT_REUSE`, `SUBSCRIPTION_REUSE`, `RESOURCE_GROUP_REUSE`, or `NO_REUSE`. Azure hashes the label with the chosen scope (a defense against dangling-DNS subdomain takeover). Requires `domainNameLabel`. |
| `reverseFqdn` | `string` | `""` | A fully qualified domain name that resolves to this address, recorded as its reverse-DNS (PTR) name. The forward record must exist before setting this. |
| `idleTimeoutInMinutes` | `int` | `4` | Idle timeout in minutes for TCP connections. Higher values suit long-lived connections (WebSocket, gRPC streaming, database connections). Range: 4–30. |
| `ipTags` | `map<string,string>` | `{}` | Azure IP tags — routing metadata attached to the address itself (e.g. `RoutingPreference: Internet`), not governance tags. Only specific tag/value pairs are permitted by Azure. Fixed at creation. |
| `ddosProtectionMode` | `enum` | inherit | `ENABLED` (dedicated IP-level protection; pair with `ddosProtectionPlanId`) or `DISABLED` (opt out even when the network is protected). Unset inherits from the virtual network's DDoS plan, if any. |
| `ddosProtectionPlanId` | `string` | `""` | ARM ID of the DDoS protection plan backing IP-level protection. Only valid when `ddosProtectionMode` is `ENABLED`. |
| `edgeZone` | `string` | `""` | Deploy the address into an Azure Edge Zone (e.g. `losangeles`) instead of the main region. Fixed at creation. |
| `tags` | `map<string,string>` | `{}` | Free-form tags merged over the Planton-derived resource tags; a user tag with the same key wins. |

## Examples

### Basic Public IP

A minimal Public IP for development or testing:

```yaml
apiVersion: azure.planton.dev/v1
kind: AzurePublicIp
metadata:
  name: dev-pip
  annotations:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: dev.AzurePublicIp.dev-pip
spec:
  region: eastus
  resourceGroup:
    value: dev-rg
  name: dev-pip
```

### Public IP with DNS Label

A Public IP with a DNS label for a stable domain name, hashed with a tenant-scoped reuse policy so the label survives safely across environments:

```yaml
apiVersion: azure.planton.dev/v1
kind: AzurePublicIp
metadata:
  name: api-pip
  annotations:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: staging.AzurePublicIp.api-pip
spec:
  region: westeurope
  resourceGroup:
    value: staging-rg
  name: api-pip
  domainNameLabel: my-api-staging
  domainNameLabelScope: TENANT_REUSE
```

After deployment, the Public IP is reachable at a stable FQDN under `westeurope.cloudapp.azure.com`.

### Zone-Redundant Production Public IP

A production Public IP spread across all three availability zones with an extended idle timeout for long-lived connections:

```yaml
apiVersion: azure.planton.dev/v1
kind: AzurePublicIp
metadata:
  name: prod-lb-pip
  annotations:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: prod.AzurePublicIp.prod-lb-pip
spec:
  region: eastus
  resourceGroup:
    value: prod-rg
  name: prod-lb-pip
  domainNameLabel: prod-api
  zones:
    - "1"
    - "2"
    - "3"
  idleTimeoutInMinutes: 15
```

### StandardV2 Address from a Reserved Prefix

A next-generation SKU address drawn from an AzurePublicIpPrefix, for attachment to a StandardV2 NAT gateway (partners allowlist the whole prefix range once):

```yaml
apiVersion: azure.planton.dev/v1
kind: AzurePublicIp
metadata:
  name: nat-egress-pip
  annotations:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: prod.AzurePublicIp.nat-egress-pip
spec:
  region: eastus
  resourceGroup:
    value: prod-rg
  name: nat-egress-pip
  sku: STANDARD_V2
  publicIpPrefixId:
    valueFrom:
      name: prod-snat-prefix
```

### DDoS-Protected IPv6 Frontend

An IPv6 address with dedicated IP-level DDoS protection:

```yaml
apiVersion: azure.planton.dev/v1
kind: AzurePublicIp
metadata:
  name: v6-frontend-pip
  annotations:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: prod.AzurePublicIp.v6-frontend-pip
spec:
  region: eastus
  resourceGroup:
    value: prod-rg
  name: v6-frontend-pip
  ipVersion: IPV6
  zones:
    - "1"
    - "2"
    - "3"
  ddosProtectionMode: ENABLED
  ddosProtectionPlanId: /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/governance-rg/providers/Microsoft.Network/ddosProtectionPlans/org-ddos-plan
```

### Using Foreign Key References

Reference a Planton-managed resource group instead of hardcoding the name:

```yaml
apiVersion: azure.planton.dev/v1
kind: AzurePublicIp
metadata:
  name: ref-pip
  annotations:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: prod.AzurePublicIp.ref-pip
spec:
  region: eastus
  resourceGroup:
    valueFrom:
      name: my-rg
  name: ref-pip
  zones:
    - "1"
    - "2"
    - "3"
```

## Stack Outputs

After deployment, the following outputs are available in `status.outputs`:

| Output | Type | Description |
|--------|------|-------------|
| `public_ip_id` | `string` | Azure Resource Manager ID of the Public IP. This is the primary output referenced by downstream resources (AzureApplicationGateway, AzureLoadBalancer, AzureNatGateway) via `valueFrom`. |
| `ip_address` | `string` | The allocated address itself. Static for the resource's lifetime — the value that lands in DNS records and partner allowlists. |
| `fqdn` | `string` | The Azure-managed FQDN (`{label}.{region}.cloudapp.azure.com`). Only populated when `domainNameLabel` is set; empty otherwise. |
| `public_ip_name` | `string` | Name of the Public IP resource. |

## Related Components

- [AzureResourceGroup](/docs/catalog/azure/azureresourcegroup) — provides the resource group for Public IP placement
- [AzurePublicIpPrefix](/docs/catalog/azure/azurepublicipprefix) — provides the reserved range Public IPs can be allocated from
- [AzureLoadBalancer](/docs/catalog/azure/azureloadbalancer) — attaches a Public IP as a frontend IP configuration
- [AzureApplicationGateway](/docs/catalog/azure/azureapplicationgateway) — attaches a Public IP for HTTP/HTTPS ingress
- [AzureNatGateway](/docs/catalog/azure/azurenatgateway) — attaches a Public IP for outbound NAT
