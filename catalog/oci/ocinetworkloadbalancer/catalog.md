# Network Load Balancer on OCI

Deploys an Oracle Cloud Infrastructure Network Load Balancer (Layer 4) with backend sets, backends, and listeners for high-performance TCP, UDP, and mixed-protocol traffic distribution. The NLB provides fully elastic bandwidth, source IP preservation, tuple-based load balancing, and health checking across HTTP, HTTPS, TCP, UDP, and DNS protocols. Integrates with Planton's Provider Connections for OCI credential management and ValueFromRef for wiring to compartments, subnets, and security groups.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **Network Load Balancer** -- the L4 load balancer in the specified compartment and subnet with configurable public/private mode, source IP preservation, symmetric hashing, and optional reserved IP addresses
- **Backend Sets** -- one per entry in `backendSets`; each defines a tuple-based load balancing policy, health checker, and failover configuration
- **Backends** -- one per entry in each backend set's `backends` list; identified by IP address, compute instance OCID, or both
- **Listeners** -- one per entry in `listeners`; each binds a port and Layer 4 protocol (TCP, UDP, TCP+UDP, or ANY) to a backend set
- **Freeform Tags** -- resource metadata tags (organization, environment, resource kind, resource ID) applied to the NLB

## Before You Deploy

### Planton Setup

- **OCI Provider Connection** -- an active connection in the Connect module with credentials for the target OCI tenancy. Map it as the default for your environment, or specify it explicitly when creating the Cloud Resource.
- **Planton Runner** -- required when using Runner-based credential delivery. Not needed for inline credentials.

### OCI Tenancy

- A compartment to place the NLB in. Provide the compartment OCID directly or reference an OciCompartment Cloud Resource via ValueFromRef.
- A subnet for the NLB. Unlike the Application Load Balancer which accepts multiple subnets, the NLB deploys into a single subnet. The subnet determines whether the NLB receives public or private IP addresses. Provide the subnet OCID directly or reference an OciSubnet Cloud Resource via ValueFromRef. Changing the subnet forces recreation.
- Backend server IP addresses or compute instance OCIDs for traffic targets.

## Deploy

### Console

Open the deployment store, find **Network Load Balancer on OCI**, and click **Deploy**. The creation wizard walks you through preset selection, environment and connection configuration, and spec fields. Start from the **Public TCP** preset in the [Presets](#presets) tab to pre-populate a public NLB with source IP preservation and HTTP health checks.

### CLI

```yaml
apiVersion: oci.planton.dev/v1
kind: OciNetworkLoadBalancer
metadata:
  name: app-nlb
  org: acme-corp
  env: prod
spec:
  compartmentId:
    value: "ocid1.compartment.oc1..example"
  subnetId:
    value: "ocid1.subnet.oc1..example"
  backendSets:
    - name: tcp-backends
      policy: five_tuple
      healthChecker:
        protocol: tcp
        port: 8080
      backends:
        - ipAddress: "10.0.1.10"
          port: 8080
        - ipAddress: "10.0.1.11"
          port: 8080
  listeners:
    - name: tcp-listener
      port: 80
      protocol: tcp
      defaultBackendSetName: tcp-backends
```

```shell
planton apply -f nlb.yaml
```

This creates a public NLB listening on port 80 (TCP) distributing traffic to two backends on port 8080 using five-tuple hashing. TCP health checks verify backend availability. Source IP preservation and failover features are not configured.

### InfraChart

When deploying as part of a multi-resource environment, use ValueFromRef to wire the NLB to a compartment, subnet, and security group deployed in the same InfraPipeline:

```yaml
spec:
  compartmentId:
    valueFrom:
      kind: OciCompartment
      name: networking
      fieldPath: status.outputs.compartmentId
  subnetId:
    valueFrom:
      kind: OciSubnet
      name: public-subnet
      fieldPath: status.outputs.subnetId
  networkSecurityGroupIds:
    - valueFrom:
        kind: OciSecurityGroup
        name: nlb-nsg
        fieldPath: status.outputs.networkSecurityGroupId
```

The InfraPipeline resolves the dependency graph, deploys the compartment, subnet, and security group first, then provisions the NLB with the resolved values.

## Key Configuration

These are the most important decisions when configuring a network load balancer. Explore the full field reference in the [API Explorer](#api-explorer) tab.

**Load balancing policy** -- NLB uses tuple-based hashing: `five_tuple` (source IP, source port, destination IP, destination port, protocol) for maximum distribution, `three_tuple` (source IP, destination IP, protocol) for moderate affinity, or `two_tuple` (source IP, destination IP) for strongest client affinity. Choose based on whether your workload benefits from session stickiness.

**Source IP preservation** -- Set `isPreserveSourceDestination: true` to forward packets with the original client IP and destination intact. Automatically enables skip-source-destination-check on the NLB's VNIC. Essential for firewalls, intrusion detection systems, and applications that need the true client IP. Enable `isSymmetricHashEnabled` alongside it to ensure return traffic follows the same path without requiring backend SNAT.

**Health check protocol** -- The NLB supports HTTP, HTTPS, TCP, UDP, and DNS health check protocols -- broader than the Application Load Balancer. Use DNS health checks for DNS server backends. Use HTTP health checks with `urlPath` and `returnCode` for application backends. Use TCP for services without HTTP endpoints.

**Failover behavior** -- Enable `isFailOpen` to keep distributing traffic to all backends when every backend is unhealthy (prevents total outage at the cost of degraded responses). Enable `isInstantFailoverEnabled` to immediately redirect existing connections when a backend goes down, rather than waiting for timeout. Add `isInstantFailoverTcpResetEnabled` to send TCP RST to clients for immediate reconnection signaling.

**Public vs private** -- Set `isPrivate: true` for internal-only NLBs (VCN traffic). Public NLBs receive ephemeral public IPs by default; use `reservedIps` to assign pre-created reserved public IPs for stable IP addresses. Changing `isPrivate` after creation forces recreation.

## Outputs and Dependencies

### What This Component Consumes

| Dependency | Field | ValueFromRef Path |
|------------|-------|-------------------|
| **OciCompartment** | `compartmentId` | `status.outputs.compartmentId` |
| **OciSubnet** | `subnetId` | `status.outputs.subnetId` |
| **OciSecurityGroup** (optional) | `networkSecurityGroupIds` | `status.outputs.networkSecurityGroupId` |

### What This Component Provides

After provisioning, `status.outputs` contains values that downstream Cloud Resources can consume via ValueFromRef:

| Output | Description | Common Downstream Use |
|--------|-------------|----------------------|
| `networkLoadBalancerId` | OCID of the network load balancer | Monitoring, IAM policy scoping, resource management |
| `ipAddresses` | Comma-separated IP addresses assigned to the NLB (public and/or private) | DNS record targets, client configuration, firewall rules |

## Common Patterns

Browse the [Presets](#presets) tab for ready-to-deploy configurations.

**Public TCP load balancer** -- A public NLB with source IP preservation, five-tuple hashing, HTTP health checks, and NSG bindings. The standard configuration for internet-facing TCP services that need the true client IP at the backend. Start from the **Public TCP** preset.

**Private internal load balancer** -- A private NLB for distributing traffic between internal VCN services. TCP health checks on the application port with no source IP preservation. Start from the **Private Internal** preset.

**Development** -- A minimal public NLB with a single backend set, TCP health checks, and no backends configured. Suitable for development environments where backends are added dynamically. Start from the **Development** preset.

## Works With

- [**Compartment on OCI**](/cloud-catalog/oci-compartment) -- provides the compartment that scopes this NLB
- [**Subnet on OCI**](/cloud-catalog/oci-subnet) -- provides the subnet where the NLB is deployed
- [**Network Security Group on OCI**](/cloud-catalog/oci-security-group) -- provides network security rules applied to the NLB