---
title: "Network Security Group"
description: "Network Security Group deployment documentation"
icon: "package"
order: 100
componentName: "ocisecuritygroup"
---

# Network Security Group on OCI

Deploys an Oracle Cloud Infrastructure Network Security Group (NSG) with bundled ingress and egress security rules. Unlike security lists (which apply at the subnet level), NSGs provide per-VNIC traffic control and are OCI's recommended approach for fine-grained network security. Rules are split into `ingressRules` and `egressRules` with support for TCP, UDP, ICMP, ICMPv6, and all-protocol matching. The component integrates with Planton's Provider Connections for OCI credential management and supports ValueFromRef wiring to compartments and VCNs.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **Network Security Group** -- an NSG in the specified compartment and VCN with a display name and freeform tags
- **Security Rules** -- one `core.NetworkSecurityGroupSecurityRule` per entry in `ingressRules` and `egressRules`, with direction, protocol, source/destination, port ranges, ICMP options, and stateless flag
- **Freeform Tags** -- resource metadata tags (organization, environment, resource kind, resource ID) applied to the NSG

## Before You Deploy

### Planton Setup

- **OCI Provider Connection** -- an active connection in the Connect module with credentials for the target OCI tenancy. Map it as the default for your environment, or specify it explicitly when creating the Cloud Resource.
- **Planton Runner** -- required when using Runner-based credential delivery. Not needed for inline credentials.

### OCI Tenancy

- A VCN to attach the NSG to. Each NSG belongs to a single VCN. Provide the VCN OCID directly or reference an OciVcn Cloud Resource via ValueFromRef. Changing the VCN forces recreation of the NSG.
- A compartment to place the NSG in. Provide the compartment OCID directly or reference an OciCompartment Cloud Resource via ValueFromRef.

## Deploy

### Console

Open the deployment store, find **Network Security Group on OCI**, and click **Deploy**. The creation wizard walks you through preset selection, environment and connection configuration, and spec fields. Start from the **Web Tier** preset in the [Presets](#presets) tab to pre-populate an NSG allowing HTTPS/HTTP inbound with ICMP Path MTU Discovery.

### CLI

```yaml
apiVersion: oci.planton.dev/v1
kind: OciSecurityGroup
metadata:
  name: web-tier-nsg
  org: acme-corp
  env: prod
spec:
  compartmentId:
    value: "ocid1.compartment.oc1..example"
  vcnId:
    value: "ocid1.vcn.oc1..example"
  displayName: web-tier-nsg
  ingressRules:
    - source: 0.0.0.0/0
      sourceType: cidr_block
      protocol: tcp
      description: Allow HTTPS from anywhere
      tcpOptions:
        destinationPortRange:
          min: 443
          max: 443
  egressRules:
    - destination: 0.0.0.0/0
      destinationType: cidr_block
      protocol: all
      description: Allow all outbound traffic
```

```shell
planton apply -f nsg.yaml
```

This creates an NSG with HTTPS inbound from anywhere and unrestricted outbound. HTTP and ICMP rules are not included in this minimal manifest.

### InfraChart

When deploying as part of a multi-resource environment, use ValueFromRef to wire the NSG to a compartment and VCN deployed in the same InfraPipeline:

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

The InfraPipeline resolves the dependency graph, deploys the compartment and VCN first, then provisions the NSG with the resolved values.

## Key Configuration

These are the most important decisions when configuring an NSG. Explore the full field reference in the [API Explorer](#api-explorer) tab.

**Rule direction** -- Rules are split into `ingressRules` (inbound to VNICs) and `egressRules` (outbound from VNICs). This eliminates the error-prone "direction + conditional source/destination" pattern from the raw OCI API.

**Protocol and port selection** -- Set `protocol` to `tcp`, `udp`, `icmp`, `icmpv6`, or `all`. For TCP/UDP rules, specify `tcpOptions.destinationPortRange` or `udpOptions.destinationPortRange` to restrict ports. Set `min == max` for a single port. When port options are omitted, all ports for that protocol are allowed.

**Source and destination types** -- Each rule's source (ingress) or destination (egress) can be a CIDR block (`cidr_block`), an OCI service CIDR label (`service_cidr_block`), or another NSG's OCID (`network_security_group`). NSG-to-NSG references enable micro-segmentation between tiers without hard-coding IP ranges.

**Rule limit** -- OCI enforces a maximum of 120 security rules per NSG (ingress + egress combined). For complex workloads, split rules across multiple NSGs and attach up to 5 NSGs per VNIC.

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
| `network_security_group_id` | OCID of the NSG | OciComputeInstance VNIC attachment, OciContainerEngineCluster, OciApplicationLoadBalancer, OciNetworkLoadBalancer |

## Common Patterns

Browse the [Presets](#presets) tab for ready-to-deploy configurations.

**Web tier** -- An NSG for internet-facing resources allowing HTTPS (443), HTTP (80), and ICMP Path MTU Discovery inbound from anywhere, with unrestricted outbound. Start from the **Web Tier** preset.

**Private backend** -- An NSG for resources reachable only from within the VCN. All protocols and ports are allowed from the VCN CIDR block, with unrestricted outbound. Start from the **Private Backend** preset.

**Development** -- A fully permissive NSG allowing all inbound and outbound traffic for development and testing. Start from the **Development** preset.

## Works With

- [**OCI Compartment**](/cloud-catalog/oci-compartment) -- provides the compartment that scopes this NSG
- [**OCI VCN**](/cloud-catalog/oci-vcn) -- provides the virtual network that this NSG belongs to