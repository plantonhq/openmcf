# Subnet on OCI

Deploys an Oracle Cloud Infrastructure subnet within a VCN with configurable CIDR block, public/private access controls, DNS label, optional IPv6, and an optional custom route table created inline from route rules. The subnet integrates with Planton's Provider Connections for OCI credential management and supports ValueFromRef wiring to compartments and VCNs.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **Subnet** -- a regional or AD-specific subnet within the specified VCN, with the configured CIDR block, DNS label, public IP prohibition settings, and security list associations
- **Route Table** -- created only when `routeRules` is populated; a custom route table with the specified routing rules is attached to the subnet. When `routeTableId` is provided instead, that existing table is referenced. When neither is provided, the VCN's default route table is used.
- **Freeform Tags** -- resource metadata tags (organization, environment, resource kind, resource ID) applied to the subnet and route table

## Before You Deploy

### Planton Setup

- **OCI Provider Connection** -- an active connection in the Connect module with credentials for the target OCI tenancy. Map it as the default for your environment, or specify it explicitly when creating the Cloud Resource.
- **Planton Runner** -- required when using Runner-based credential delivery. Not needed for inline credentials.

### OCI Tenancy

- A VCN with at least one CIDR block that encompasses the subnet's CIDR. Provide the VCN OCID directly or reference an OciVcn Cloud Resource via ValueFromRef.
- A compartment to place the subnet in. Provide the compartment OCID directly or reference an OciCompartment Cloud Resource via ValueFromRef.
- For private subnets: a NAT Gateway OCID (from the VCN) for the outbound internet route rule, and a Service Gateway OCID for the OCI service route rule.
- For public subnets: an Internet Gateway OCID (from the VCN) for the default route rule.

## Deploy

### Console

Open the deployment store, find **Subnet on OCI**, and click **Deploy**. The creation wizard walks you through preset selection, environment and connection configuration, and spec fields. Start from the **Private** preset in the [Presets](#presets) tab to pre-populate a private subnet with NAT and Service Gateway routing.

### CLI

```yaml
apiVersion: oci.planton.dev/v1
kind: OciSubnet
metadata:
  name: private-app-subnet
  org: acme-corp
  env: prod
spec:
  compartmentId:
    value: "ocid1.compartment.oc1..example"
  vcnId:
    value: "ocid1.vcn.oc1..example"
  cidrBlock: 10.0.1.0/24
  displayName: private-app-subnet
  dnsLabel: privapp
  prohibitPublicIpOnVnic: true
  prohibitInternetIngress: true
```

```shell
planton apply -f subnet.yaml
```

This creates a private subnet with no public IP assignment and no inbound internet traffic. No custom route rules are configured, so the VCN's default route table is used.

### InfraChart

When deploying as part of a multi-resource environment, use ValueFromRef to wire the subnet to a compartment and VCN deployed in the same InfraPipeline:

```yaml
spec:
  compartmentId:
    valueFrom:
      kind: OciCompartment
      name: networking
      fieldPath: status.outputs.compartmentId
  vcnId:
    valueFrom:
      kind: OciVcn
      name: production-vcn
      fieldPath: status.outputs.vcnId
```

The InfraPipeline resolves the dependency graph, deploys the compartment and VCN first, then provisions the subnet with the resolved values.

## Key Configuration

These are the most important decisions when configuring a subnet. Explore the full field reference in the [API Explorer](#api-explorer) tab.

**Public vs private** -- Set `prohibitPublicIpOnVnic: true` and `prohibitInternetIngress: true` for private subnets (worker nodes, databases, internal services). Set both to `false` for public subnets (load balancers, bastion hosts). These are the primary controls for subnet network exposure.

**Route table strategy** -- Provide inline `routeRules` to create a custom route table owned by this subnet (the typical pattern for private subnets needing NAT and Service Gateway routes). Alternatively, reference an existing route table via `routeTableId`. The two fields are mutually exclusive. When neither is provided, the VCN's default route table is used.

**Regional vs AD-specific** -- Omit `availabilityDomain` for a regional subnet that spans all availability domains (recommended for most workloads). Set it to a specific AD name (e.g., `Iocq:US-ASHBURN-AD-1`) only when workload placement requires AD affinity.

**DNS label** -- Set `dnsLabel` to enable subnet-level DNS. Combined with the VCN's DNS label, instances get FQDNs in the form `<instance>.<subnet-dns-label>.<vcn-dns-label>.oraclevcn.com`.

## Outputs and Dependencies

### What This Component Consumes

| Dependency | Field | ValueFromRef Path |
|------------|-------|-------------------|
| **OciCompartment** | `compartmentId` | `status.outputs.compartmentId` |
| **OciVcn** | `vcnId` | `status.outputs.vcnId` |

### What This Component Provides

After provisioning, `status.outputs` contains values that downstream Cloud Resources can consume via ValueFromRef:

| Output | Description | Common Downstream Use |
|--------|-------------|----------------------|
| `subnet_id` | OCID of the subnet | OciComputeInstance, OciContainerEngineCluster, OciApplicationLoadBalancer, OciMysqlDbSystem |
| `subnet_domain_name` | FQDN of the subnet (e.g., `subnet1.vcn1.oraclevcn.com`) | DNS configuration, service discovery |
| `virtual_router_ip` | IP address of the virtual router in the subnet | Network debugging, static route configuration |
| `virtual_router_mac` | MAC address of the virtual router | Network debugging |
| `route_table_id` | OCID of the associated route table (custom, referenced, or VCN default) | Network auditing, related resource lookups |

## Common Patterns

Browse the [Presets](#presets) tab for ready-to-deploy configurations.

**Private subnet** -- A subnet with public IP prohibited, internet ingress blocked, and inline route rules directing traffic through a NAT Gateway (outbound internet) and Service Gateway (OCI services). The standard configuration for production workloads. Start from the **Private** preset.

**Public subnet** -- A subnet that allows public IP assignment and internet ingress, with a route rule directing all traffic through the VCN's Internet Gateway. Used for load balancers, bastion hosts, and API gateways. Start from the **Public** preset.

**Development subnet** -- A minimal public subnet using the VCN's default route table with no custom route rules. Suitable for development and testing where simplicity matters more than network segmentation. Start from the **Development** preset.

## Works With

- [**OCI Compartment**](/cloud-catalog/oci-compartment) -- provides the compartment that scopes this subnet
- [**OCI VCN**](/cloud-catalog/oci-vcn) -- provides the virtual network that contains this subnet