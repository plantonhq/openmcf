---
title: "Load Balancer Listener"
description: "Load Balancer Listener deployment documentation"
icon: "package"
order: 100
componentName: "openstackloadbalancerlistener"
---

# OpenStack Load Balancer Listener

Deploys an Octavia listener on OpenStack that binds a protocol and port to a load balancer, accepting incoming client traffic and forwarding it to a backend pool. The listener supports HTTP, HTTPS passthrough, TLS termination, TCP, and UDP protocols with optional CIDR-based access control, connection limits, and header insertion. ValueFromRef wiring connects the listener to an OpenStackLoadBalancer Cloud Resource in InfraChart deployments.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **Octavia Listener** -- a listener resource attached to the specified load balancer, accepting traffic on the configured protocol and port
- **TLS Configuration** -- created only when protocol is `TERMINATED_HTTPS`; references a Barbican TLS secret container for certificate and private key
- **Header Insertion** -- created only when `insertHeaders` is specified; injects headers like `X-Forwarded-For` and `X-Forwarded-Proto` into HTTP requests before forwarding to backends
- **CIDR Allow List** -- created only when `allowedCidrs` is specified; restricts which source IP ranges can reach the listener
- **OpenStack Tags** -- user-defined tags applied to the listener for filtering and organization

## Before You Deploy

### OpenStack Account

- **Load Balancer** -- an Octavia load balancer must exist to attach the listener to. Provide the load balancer ID directly or reference an OpenStackLoadBalancer Cloud Resource via ValueFromRef.
- **TLS certificate** (for TERMINATED_HTTPS) -- store the certificate and private key in Barbican as a TLS secret container. Note the container URI for the `defaultTlsContainerRef` field.

## Deploy

### Console

Open the deployment store, find **OpenStack Load Balancer Listener**, and click **Deploy**. The creation wizard walks you through preset selection, environment and connection configuration, and spec fields. Start from the **HTTP Listener** preset in the [Presets](#presets) tab to pre-populate a working configuration.

### CLI

Create a manifest and apply it:

```yaml
apiVersion: openstack.planton.dev/v1
kind: OpenStackLoadBalancerListener
metadata:
  name: http-listener
  org: acme-corp
  env: prod
spec:
  loadbalancerId:
    value: "<loadbalancer-id>"
  protocol: HTTP
  protocolPort: 80
```

```shell
planton apply -f listener.yaml
```

This creates an HTTP listener on port 80 with no connection limit, no CIDR restrictions, and no header insertion. Attach a pool to begin routing traffic to backends.

### InfraChart

When deploying as part of a multi-resource environment, use ValueFromRef to wire the listener to a load balancer deployed in the same InfraPipeline:

```yaml
spec:
  loadbalancerId:
    valueFrom:
      kind: OpenStackLoadBalancer
      name: web-lb
      fieldPath: status.outputs.loadbalancer_id
```

The InfraPipeline resolves the dependency graph, deploys the load balancer first, then provisions the listener with the resolved load balancer ID.

## Key Configuration

These are the most important decisions when configuring a listener. Explore the full field reference in the [API Explorer](#api-explorer) tab.

**Protocol selection** -- Choose `HTTP` for unencrypted Layer 7 traffic, `TERMINATED_HTTPS` for TLS termination at the load balancer, `HTTPS` for encrypted passthrough, `TCP` for Layer 4 services like databases, or `UDP` for datagram-based services. Changing the protocol requires recreating the listener.

**TLS termination** -- When using `TERMINATED_HTTPS`, the `defaultTlsContainerRef` field is required and must point to a Barbican secret container with the certificate and private key. Add `insertHeaders` with `X-Forwarded-For` and `X-Forwarded-Proto` so backends can see the original client IP and protocol.

**Connection limits** -- Leave `connectionLimit` unset to use Octavia's default. Set a specific value to cap the maximum concurrent connections, or `-1` for explicitly unlimited connections.

**Access control** -- Use `allowedCidrs` to restrict the listener to specific source IP ranges. When empty, all traffic is allowed. This provides network-level access control without requiring external security group rules.

## Outputs and Dependencies

### What This Component Consumes

| Dependency | Field | ValueFromRef Path |
|------------|-------|-------------------|
| **OpenStackLoadBalancer** | `loadbalancerId` | `status.outputs.loadbalancer_id` |

### What This Component Provides

After provisioning, `status.outputs` contains values that downstream Cloud Resources can consume via ValueFromRef:

| Output | Description | Common Downstream Use |
|--------|-------------|----------------------|
| `listener_id` | UUID of the listener | Pool attachment via `listenerId` |
| `name` | Name of the listener | Monitoring labels, resource identification |
| `protocol` | Protocol the listener accepts | Downstream configuration validation |
| `protocol_port` | Port the listener accepts traffic on | DNS SRV records, client configuration |
| `region` | OpenStack region where the listener was created | Region-aware downstream resources |

## Common Patterns

Browse the [Presets](#presets) tab for ready-to-deploy configurations.

**HTTP Listener** -- Accepts unencrypted HTTP traffic on port 80. Suitable for internal services, HTTP-to-HTTPS redirects, and development environments. Start from the **HTTP Listener** preset.

**HTTPS with TLS Termination** -- Terminates TLS at the load balancer on port 443, forwarding decrypted traffic to backends as plain HTTP. Inserts `X-Forwarded-For` and `X-Forwarded-Proto` headers for backend visibility. Start from the **HTTPS Listener with TLS Termination** preset.

**TCP Passthrough** -- Passes raw TCP traffic to backends without protocol-level processing. Use for databases, message queues, gRPC, or end-to-end TLS. Start from the **TCP Passthrough Listener** preset.

## Works With

- [**OpenStack Load Balancer**](/cloud-catalog/openstack-load-balancer) -- provides the load balancer ID that this listener attaches to