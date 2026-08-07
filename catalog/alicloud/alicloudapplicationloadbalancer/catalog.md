# AliCloud Application Load Balancer

Deploys an Alibaba Cloud Application Load Balancer (ALB) with bundled server groups and listeners. ALB is a modern Layer 7 load balancer for HTTP, HTTPS, and QUIC traffic, offering advanced routing, health checking, and optional WAF integration. The ALB, server groups, and listeners are deployed as a single atomic unit because an ALB without at least one server group and listener is non-functional.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **ALB** -- an `alicloud_alb_load_balancer` resource placed across at least two Availability Zones for high availability, with configurable edition, address type, and access logging
- **Server Groups** -- one `alicloud_alb_server_group` per entry in `serverGroups`, each with its own health check, protocol, scheduling algorithm, and stickiness configuration. Groups are created empty -- backend membership is managed externally by ACK ingress controllers or manual attachment
- **Listeners** -- one `alicloud_alb_listener` per entry in `listeners`, each binding a port/protocol to a server group via default forwarding action

## Before You Deploy

### Planton Setup

- **AliCloud Provider Connection** -- an active connection in the Connect module with credentials for the target Alibaba Cloud account. Map it as the default for your environment, or specify it explicitly when creating the Cloud Resource.
- **Planton Runner** -- required when using Runner-based credential delivery.

### Alibaba Cloud Account

- **An existing VPC** with at least two VSwitches in different Availability Zones -- ALB requires multi-AZ deployment.
- **A server certificate** (for HTTPS listeners) -- obtain from Alibaba Cloud Certificate Management Service (CAS).
- **An SLS log project and log store** (optional) -- for access log delivery.

## Deploy

### Console

Open the deployment store, find **AliCloud Application Load Balancer**, and click **Deploy**. The creation wizard walks you through preset selection, environment and connection configuration, and spec fields including VPC, zone mappings, server groups, and listeners.

### CLI

Create a manifest and apply it:

```yaml
apiVersion: alicloud.planton.dev/v1
kind: AliCloudApplicationLoadBalancer
metadata:
  name: web-alb
  org: acme-corp
  env: prod
spec:
  region: cn-hangzhou
  vpcId:
    value: vpc-bp1234567890
  zoneMappings:
    - zoneId: cn-hangzhou-a
      vswitchId:
        value: vsw-zone-a
    - zoneId: cn-hangzhou-b
      vswitchId:
        value: vsw-zone-b
  serverGroups:
    - name: web-backend
      healthCheckConfig:
        healthCheckEnabled: true
        healthCheckPath: /health
  listeners:
    - listenerPort: 80
      listenerProtocol: HTTP
      defaultActionServerGroupName: web-backend
```

```shell
planton apply -f alicloud-alb.yaml
```

This creates an internet-facing ALB with an HTTP listener forwarding to the web-backend server group. A Stack Job tracks the provisioning in real time.

### InfraChart

When deploying as part of an application stack, use ValueFromRef to wire VPC and VSwitch dependencies:

```yaml
spec:
  region: cn-hangzhou
  vpcId:
    valueFrom:
      kind: AliCloudVpc
      name: platform-vpc
      fieldPath: status.outputs.vpc_id
  zoneMappings:
    - zoneId: cn-hangzhou-a
      vswitchId:
        valueFrom:
          kind: AliCloudVswitch
          name: app-vswitch-a
          fieldPath: status.outputs.vswitch_id
    - zoneId: cn-hangzhou-b
      vswitchId:
        valueFrom:
          kind: AliCloudVswitch
          name: app-vswitch-b
          fieldPath: status.outputs.vswitch_id
```

The InfraPipeline resolves the dependency graph and provisions VPC and VSwitches before the ALB.

## Key Configuration

These are the most important decisions when configuring an ALB. Explore the full field reference in the [API Explorer](#api-explorer) tab.

**Address type** -- The `addressType` field determines whether the ALB is internet-facing ("Internet", default) or VPC-internal ("Intranet"). This affects DNS resolution and accessibility.

**Edition** -- The `loadBalancerEdition` field selects the feature tier. "Standard" (default) supports advanced routing. "StandardWithWaf" adds integrated WAF protection. "Basic" is limited to basic L7 load balancing.

**Server group protocol** -- Each server group's `protocol` field controls how the ALB communicates with backends: "HTTP" (default), "HTTPS" (end-to-end encryption), or "GRPC" (for gRPC services).

**Health checks** -- Each server group requires a `healthCheckConfig`. Configure the path, interval, and thresholds to match your application's health endpoint.

**HTTPS listeners** -- HTTPS listeners require a `certificateId` from CAS and support TLS security policy selection via `securityPolicyId` for compliance requirements.

## Outputs and Dependencies

### What This Component Consumes

| Dependency | Field | ValueFromRef Path |
|------------|-------|-------------------|
| **AliCloudVpc** | `vpcId` | `status.outputs.vpc_id` |
| **AliCloudVswitch** | `zoneMappings[].vswitchId` | `status.outputs.vswitch_id` |

### What This Component Provides

After provisioning, `status.outputs` contains values that downstream Cloud Resources can consume via ValueFromRef:

| Output | Description | Common Downstream Use |
|--------|-------------|----------------------|
| `load_balancer_id` | ALB instance ID (e.g., alb-xxxxx) | Monitoring, forwarding rules |
| `dns_name` | Auto-assigned DNS name for the ALB | CNAME target for custom domain DNS records |
| `server_group_ids` | Map of server group names to IDs | Backend attachment, advanced forwarding rules |

## Common Patterns

Browse the [Presets](#presets) tab for ready-to-deploy configurations.

**Internet HTTP** -- A public-facing ALB with an HTTP listener and basic health checks. The simplest starting point for web applications. Start from the **Internet HTTP** preset.

**HTTPS production** -- A StandardWithWaf ALB with HTTPS listener, strict TLS policy, access logging, and session stickiness. Start from the **HTTPS Production** preset.

**Internal gRPC** -- A VPC-internal ALB with a gRPC server group using least-connections scheduling. Start from the **Internal gRPC** preset.

## Works With

- [**AliCloud VPC**](/cloud-catalog/ali-cloud-vpc) -- the VPC this ALB belongs to
- [**AliCloud VSwitch**](/cloud-catalog/ali-cloud-vswitch) -- provides zone-specific IP allocation for ALB nodes
- [**AliCloud Security Group**](/cloud-catalog/ali-cloud-security-group) -- network security rules for backend instances
- [**AliCloud ECS Instance**](/cloud-catalog/ali-cloud-ecs-instance) -- backend targets for server groups
- [**AliCloud Kubernetes Cluster**](/cloud-catalog/ali-cloud-kubernetes-cluster) -- ACK ingress controllers manage ALB server group membership
- [**AliCloud DNS Record**](/cloud-catalog/ali-cloud-dns-record) -- CNAME records pointing to the ALB's dns_name
