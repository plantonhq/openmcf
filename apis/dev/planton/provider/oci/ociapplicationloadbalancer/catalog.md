# Application Load Balancer on OCI

Deploys an Oracle Cloud Infrastructure Application Load Balancer (Layer 7) with backend sets, backends, listeners, TLS certificates, virtual hostnames, and rule sets for HTTP/HTTPS/gRPC traffic distribution. The ALB supports flexible bandwidth scaling, SSL termination and backend re-encryption, cookie-based session persistence, host-based routing, and rule-based request manipulation (redirects, header injection, access control). Integrates with Planton's Provider Connections for OCI credential management and ValueFromRef for wiring to compartments, subnets, and security groups.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **Load Balancer** -- the L7 load balancer in the specified compartment and subnets with configurable shape (flexible or fixed bandwidth), public/private mode, delete protection, and optional reserved IP addresses
- **Backend Sets** -- one per entry in `backendSets`; each defines a load balancing policy (round robin, least connections, IP hash), health checker, optional SSL re-encryption, session persistence, and connection limits
- **Backends** -- one per entry in each backend set's `backends` list; identified by IP address with configurable weight, backup, drain, and offline states
- **Listeners** -- one per entry in `listeners`; each binds a port and protocol (HTTP, HTTP/2, TCP, gRPC) to a backend set with optional SSL termination, virtual hostname filtering, and rule set application
- **Certificates** -- one per entry in `certificates`; PEM-encoded TLS certificates referenced by listener and backend set SSL configurations
- **Hostnames** -- one per entry in `hostnames`; virtual hostname definitions for host-based routing via listener hostname bindings
- **Rule Sets** -- one per entry in `ruleSets`; groups of rules for HTTP redirects, header add/remove/extend, access control by HTTP method, and IP-based connection limits
- **Freeform Tags** -- resource metadata tags (organization, environment, resource kind, resource ID) applied to the load balancer

## Before You Deploy

### Planton Setup

- **OCI Provider Connection** -- an active connection in the Connect module with credentials for the target OCI tenancy. Map it as the default for your environment, or specify it explicitly when creating the Cloud Resource.
- **Planton Runner** -- required when using Runner-based credential delivery. Not needed for inline credentials.

### OCI Tenancy

- A compartment to place the load balancer in. Provide the compartment OCID directly or reference an OciCompartment Cloud Resource via ValueFromRef.
- One or more subnets for the load balancer. For regional high availability, provide subnets in two different availability domains. Changing subnets after creation forces recreation. Provide subnet OCIDs directly or reference OciSubnet Cloud Resources via ValueFromRef.
- OCI Certificate Service certificate OCIDs (for HTTPS listeners) or PEM-encoded certificates defined inline in the manifest.

## Deploy

### Console

Open the deployment store, find **Application Load Balancer on OCI**, and click **Deploy**. The creation wizard walks you through preset selection, environment and connection configuration, and spec fields. Start from the **Internet-Facing HTTPS** preset in the [Presets](#presets) tab to pre-populate a public ALB with SSL termination and HTTP-to-HTTPS redirect.

### CLI

```yaml
apiVersion: oci.planton.dev/v1
kind: OciApplicationLoadBalancer
metadata:
  name: web-lb
  org: acme-corp
  env: prod
spec:
  compartmentId:
    value: "ocid1.compartment.oc1..example"
  shape: flexible
  shapeDetails:
    minimumBandwidthInMbps: 10
    maximumBandwidthInMbps: 100
  subnetIds:
    - value: "ocid1.subnet.oc1..ad1-example"
    - value: "ocid1.subnet.oc1..ad2-example"
  isPrivate: false
  backendSets:
    - name: web-backend
      policy: round_robin
      healthChecker:
        protocol: http
        port: 80
        urlPath: /health
        returnCode: 200
  listeners:
    - name: http
      port: 80
      protocol: http
      defaultBackendSetName: web-backend
```

```shell
planton apply -f alb.yaml
```

This creates a public flexible-shape ALB across two subnets with a single backend set using round-robin policy and HTTP health checks. No SSL termination, hostnames, or rule sets are configured. Backends are omitted and can be added after creation.

### InfraChart

When deploying as part of a multi-resource environment, use ValueFromRef to wire the ALB to a compartment, subnets, and security groups deployed in the same InfraPipeline:

```yaml
spec:
  compartmentId:
    valueFrom:
      kind: OciCompartment
      name: networking
      fieldPath: status.outputs.compartmentId
  subnetIds:
    - valueFrom:
        kind: OciSubnet
        name: public-subnet-ad1
        fieldPath: status.outputs.subnetId
    - valueFrom:
        kind: OciSubnet
        name: public-subnet-ad2
        fieldPath: status.outputs.subnetId
  networkSecurityGroupIds:
    - valueFrom:
        kind: OciSecurityGroup
        name: lb-nsg
        fieldPath: status.outputs.networkSecurityGroupId
```

The InfraPipeline resolves the dependency graph, deploys the compartment, subnets, and security group first, then provisions the ALB with the resolved values.

## Key Configuration

These are the most important decisions when configuring an application load balancer. Explore the full field reference in the [API Explorer](#api-explorer) tab.

**Shape and bandwidth** -- Use `shape: flexible` with `shapeDetails` to configure minimum and maximum bandwidth in Mbps (10-8000). The ALB scales within this range based on traffic. Fixed shapes (`100Mbps`, `400Mbps`, `8000Mbps`) are deprecated but accepted for backward compatibility.

**SSL termination** -- Configure `sslConfiguration` on a listener with certificate references (either OCI Certificate Service OCIDs via `certificateIds` or inline PEM certificates via `certificateName`). Specify TLS protocol versions (`TLSv1.2`, `TLSv1.3`) and cipher suite. For end-to-end encryption, add `sslConfiguration` to the backend set for re-encryption between the ALB and backends.

**Session persistence** -- Choose between load-balancer-managed cookies (`lbCookieSessionPersistence` -- the ALB injects and tracks a cookie) or application-managed cookies (`appCookieSessionPersistence` -- the ALB reads an existing application cookie). The two are mutually exclusive per backend set. Omit both for stateless distribution.

**Rule sets** -- Attach rule sets to listeners for HTTP redirects (e.g., HTTP-to-HTTPS via `redirect` action with response code 301), header manipulation (add/remove/extend request or response headers), access control by HTTP method, and IP-based connection limiting. Rules in a set are applied in order.

**Virtual hostname routing** -- Define `hostnames` with FQDN values and reference them in listener `hostnameNames` to route requests based on the HTTP Host header. This enables a single ALB to serve multiple domains, each routing to different backend sets via separate listeners.

## Outputs and Dependencies

### What This Component Consumes

| Dependency | Field | ValueFromRef Path |
|------------|-------|-------------------|
| **OciCompartment** | `compartmentId` | `status.outputs.compartmentId` |
| **OciSubnet** | `subnetIds` | `status.outputs.subnetId` |
| **OciSecurityGroup** (optional) | `networkSecurityGroupIds` | `status.outputs.networkSecurityGroupId` |

### What This Component Provides

After provisioning, `status.outputs` contains values that downstream Cloud Resources can consume via ValueFromRef:

| Output | Description | Common Downstream Use |
|--------|-------------|----------------------|
| `loadBalancerId` | OCID of the load balancer | Monitoring, IAM policy scoping, resource management |
| `ipAddresses` | Comma-separated IP addresses assigned to the ALB (public and/or private) | DNS record targets (OciDnsRecord), client configuration |

## Common Patterns

Browse the [Presets](#presets) tab for ready-to-deploy configurations.

**Internet-facing HTTPS** -- A public flexible-shape ALB with HTTPS listener (TLS 1.2/1.3), HTTP-to-HTTPS redirect rule set, HTTP health checks, delete protection, and NSG bindings. Multi-subnet deployment for regional HA. Start from the **Internet-Facing HTTPS** preset.

**Internal HTTP** -- A private ALB for distributing traffic between internal VCN services on port 80. Least-connections policy for balanced workload distribution. No SSL, no rule sets. Start from the **Internal HTTP** preset.

**Development** -- A minimal public ALB with the lowest flexible bandwidth (10 Mbps), a single backend set with TCP health checks, and no backends. Suitable for development and testing. Start from the **Development** preset.

## Works With

- [**Compartment on OCI**](/cloud-catalog/oci-compartment) -- provides the compartment that scopes this load balancer
- [**Subnet on OCI**](/cloud-catalog/oci-subnet) -- provides the subnets where the ALB is deployed for regional HA
- [**Network Security Group on OCI**](/cloud-catalog/oci-security-group) -- provides network security rules applied to the ALB