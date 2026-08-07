# Virtual Cloud Network on OCI

Deploys an Oracle Cloud Infrastructure Virtual Cloud Network (VCN) with multiple CIDR blocks, DNS resolution, optional IPv6, and bundled gateway resources -- Internet Gateway, NAT Gateway, and Service Gateway. Each gateway is controlled by a boolean toggle and only created when enabled. The VCN integrates with Planton's Provider Connections for OCI credential management and supports ValueFromRef wiring to compartments.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **VCN** -- the virtual network in the specified compartment with one or more IPv4 CIDR blocks, DNS label, and optional IPv6 prefix
- **Internet Gateway** -- created only when `isInternetGatewayEnabled` is `true`; provides direct inbound and outbound internet access for public subnets
- **NAT Gateway** -- created only when `isNatGatewayEnabled` is `true`; allows private subnets to initiate outbound internet connections without exposing them to inbound traffic
- **Service Gateway** -- created only when `isServiceGatewayEnabled` is `true`; provides private access to OCI services (Object Storage, Container Registry) without traffic leaving the Oracle backbone. Automatically configured for all services in Oracle Services Network.
- **Freeform Tags** -- resource metadata tags (organization, environment, resource kind, resource ID) applied to the VCN and all gateway resources

## Before You Deploy

### Planton Setup

- **OCI Provider Connection** -- an active connection in the Connect module with credentials for the target OCI tenancy. Map it as the default for your environment, or specify it explicitly when creating the Cloud Resource.
- **Planton Runner** -- required when using Runner-based credential delivery. Not needed for inline credentials.

### OCI Tenancy

- A compartment to place the VCN in. Provide the compartment OCID directly or reference an OciCompartment Cloud Resource via ValueFromRef.
- A CIDR plan that avoids overlap with existing VCNs if you plan to use DRG peering, Site-to-Site VPN, or FastConnect.

## Deploy

### Console

Open the deployment store, find **Virtual Cloud Network on OCI**, and click **Deploy**. The creation wizard walks you through preset selection, environment and connection configuration, and spec fields. Start from the **Standard Public-Private** preset in the [Presets](#presets) tab to pre-populate a production-ready VCN with all three gateways.

### CLI

```yaml
apiVersion: oci.planton.dev/v1
kind: OciVcn
metadata:
  name: production-vcn
  org: acme-corp
  env: prod
spec:
  compartmentId:
    value: "ocid1.compartment.oc1..example"
  cidrBlocks:
    - 10.0.0.0/16
  displayName: production-vcn
  dnsLabel: prodvcn
  isInternetGatewayEnabled: true
  isNatGatewayEnabled: true
  isServiceGatewayEnabled: true
```

```shell
planton apply -f vcn.yaml
```

This creates a VCN with a /16 CIDR, DNS resolution, and all three gateways enabled. IPv6 is not configured.

### InfraChart

When deploying as part of a multi-resource environment, use ValueFromRef to wire the VCN to a compartment deployed in the same InfraPipeline:

```yaml
spec:
  compartmentId:
    valueFrom:
      kind: OciCompartment
      name: networking
      fieldPath: status.outputs.compartmentId
```

The InfraPipeline resolves the dependency graph, deploys the compartment first, then provisions the VCN with the resolved compartment OCID.

## Key Configuration

These are the most important decisions when configuring a VCN. Explore the full field reference in the [API Explorer](#api-explorer) tab.

**CIDR blocks** -- Provide at least one IPv4 CIDR block in `cidrBlocks`. Each block must be between /16 and /30. OCI supports multiple non-overlapping CIDRs per VCN (e.g., `["10.0.0.0/16", "172.16.0.0/16"]`). Plan CIDR ranges to avoid overlap with other VCNs if you intend to use DRG peering.

**DNS label** -- Set `dnsLabel` to enable VCN-internal DNS resolution. Instances get hostnames in the form `<instance>.<subnet-dns-label>.<vcn-dns-label>.oraclevcn.com`. Must be alphanumeric, start with a letter, and be at most 15 characters. Omit only if you do not need VCN-internal DNS.

**Gateway selection** -- Enable `isInternetGatewayEnabled` for public subnets that need direct internet access. Enable `isNatGatewayEnabled` for private subnets that need outbound internet (patching, image pulls). Enable `isServiceGatewayEnabled` for private access to OCI services over the Oracle backbone. NAT and Service Gateways incur hourly charges -- omit them in development environments.

**IPv6** -- Set `isIpv6Enabled: true` to allocate an Oracle-assigned /56 IPv6 GUA prefix. IPv6 can be enabled later but cannot be disabled once on.

## Outputs and Dependencies

### What This Component Consumes

| Dependency | Field | ValueFromRef Path |
|------------|-------|-------------------|
| **OciCompartment** | `compartmentId` | `status.outputs.compartmentId` |

### What This Component Provides

After provisioning, `status.outputs` contains values that downstream Cloud Resources can consume via ValueFromRef:

| Output | Description | Common Downstream Use |
|--------|-------------|----------------------|
| `vcn_id` | OCID of the VCN | OciSubnet, OciSecurityGroup, OciContainerEngineCluster, OciApplicationLoadBalancer |
| `default_route_table_id` | OCID of the VCN's default route table | Subnets that use the default route table instead of a custom one |
| `default_security_list_id` | OCID of the VCN's default security list | Subnets that use the default security list |
| `default_dhcp_options_id` | OCID of the VCN's default DHCP options | Subnets that use the default DHCP options |
| `internet_gateway_id` | OCID of the Internet Gateway (empty when disabled) | OciSubnet route rules targeting internet traffic |
| `nat_gateway_id` | OCID of the NAT Gateway (empty when disabled) | OciSubnet route rules for private subnet outbound internet |
| `service_gateway_id` | OCID of the Service Gateway (empty when disabled) | OciSubnet route rules for private OCI service access |

## Common Patterns

Browse the [Presets](#presets) tab for ready-to-deploy configurations.

**Standard public-private** -- A production VCN with all three gateways enabled: Internet Gateway for public subnets, NAT Gateway for private outbound, Service Gateway for private OCI service access. Covers 80%+ of production deployments. Start from the **Standard Public-Private** preset.

**Private only** -- A security-hardened VCN with no Internet Gateway. NAT and Service Gateways provide outbound connectivity, but no resources are directly reachable from the internet. Start from the **Private Only** preset.

**Development** -- A minimal-cost VCN with only an Internet Gateway. NAT and Service Gateways are omitted to avoid hourly charges. Start from the **Development** preset.

## Works With

- [**OCI Compartment**](/cloud-catalog/oci-compartment) -- provides the compartment that scopes this VCN and its gateways