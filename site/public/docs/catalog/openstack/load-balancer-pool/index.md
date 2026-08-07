---
title: "Load Balancer Pool"
description: "Load Balancer Pool deployment documentation"
icon: "package"
order: 100
componentName: "openstackloadbalancerpool"
---

# OpenStack Load Balancer Pool

Deploys an Octavia backend pool on OpenStack that groups backend members and defines the protocol and load-balancing algorithm for traffic distribution from a listener. The pool supports round-robin, least-connections, and source-IP-based algorithms with optional session persistence via HTTP cookies, application cookies, or source IP hashing. ValueFromRef wiring connects the pool to an OpenStackLoadBalancerListener Cloud Resource in InfraChart deployments.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **Octavia Pool** -- a backend pool resource attached to the specified listener, defining the protocol and load-balancing algorithm used to distribute traffic across members
- **Session Persistence** -- created only when `persistence` is specified; configures sticky sessions via source IP, Octavia-managed HTTP cookie, or application-managed cookie
- **OpenStack Tags** -- user-defined tags applied to the pool for filtering and organization

## Before You Deploy

### OpenStack Account

- **Listener** -- an Octavia listener must exist to attach the pool to. Provide the listener ID directly or reference an OpenStackLoadBalancerListener Cloud Resource via ValueFromRef.
- **Protocol alignment** -- the pool's protocol should match the listener's protocol. An HTTP listener typically uses an HTTP pool; a TCP listener uses a TCP pool.

## Deploy

### Console

Open the deployment store, find **OpenStack Load Balancer Pool**, and click **Deploy**. The creation wizard walks you through preset selection, environment and connection configuration, and spec fields. Start from the **Round-Robin HTTP Pool** preset in the [Presets](#presets) tab to pre-populate a working configuration.

### CLI

Create a manifest and apply it:

```yaml
apiVersion: openstack.planton.dev/v1
kind: OpenStackLoadBalancerPool
metadata:
  name: web-pool
  org: acme-corp
  env: prod
spec:
  listenerId:
    value: "<listener-id>"
  protocol: HTTP
  lbMethod: ROUND_ROBIN
```

```shell
planton apply -f pool.yaml
```

This creates an HTTP pool with round-robin distribution and no session persistence. Add members and a health monitor to complete the traffic path.

### InfraChart

When deploying as part of a multi-resource environment, use ValueFromRef to wire the pool to a listener deployed in the same InfraPipeline:

```yaml
spec:
  listenerId:
    valueFrom:
      kind: OpenStackLoadBalancerListener
      name: http-listener
      fieldPath: status.outputs.listener_id
```

The InfraPipeline resolves the dependency graph, deploys the listener first, then provisions the pool with the resolved listener ID.

## Key Configuration

These are the most important decisions when configuring a pool. Explore the full field reference in the [API Explorer](#api-explorer) tab.

**Load-balancing algorithm** -- `ROUND_ROBIN` distributes traffic equally across members. `LEAST_CONNECTIONS` sends traffic to the member with the fewest active connections, best for long-lived connections. `SOURCE_IP` hashes the client IP for sticky routing without cookies. `SOURCE_IP_PORT` provides finer-grained stickiness by hashing both IP and port.

**Session persistence** -- Leave `persistence` unset for stateless services. For applications with server-side sessions, choose `HTTP_COOKIE` (Octavia manages the cookie automatically), `APP_COOKIE` (your application manages the cookie -- requires `cookieName`), or `SOURCE_IP` (hash-based, no cookies). Only one persistence type is allowed per pool.

**Protocol selection** -- Choose `HTTP`, `HTTPS`, `TCP`, `UDP`, or `PROXY` to match the backend communication protocol. The pool protocol should align with the listener protocol. `PROXY` enables the PROXY protocol for passing client connection metadata to backends that support it.

## Outputs and Dependencies

### What This Component Consumes

| Dependency | Field | ValueFromRef Path |
|------------|-------|-------------------|
| **OpenStackLoadBalancerListener** | `listenerId` | `status.outputs.listener_id` |

### What This Component Provides

After provisioning, `status.outputs` contains values that downstream Cloud Resources can consume via ValueFromRef:

| Output | Description | Common Downstream Use |
|--------|-------------|----------------------|
| `pool_id` | UUID of the pool | Member registration via `poolId`, monitor attachment via `poolId` |
| `name` | Name of the pool | Monitoring labels, resource identification |
| `protocol` | Backend protocol | Diagnostics, configuration auditing |
| `lb_method` | Load-balancing algorithm | Monitoring, capacity planning |
| `region` | OpenStack region where the pool was created | Region-aware downstream resources |

## Common Patterns

Browse the [Presets](#presets) tab for ready-to-deploy configurations.

**Round-Robin HTTP Pool** -- Distributes traffic equally across all healthy members using the round-robin algorithm over HTTP. Best for stateless web applications and REST APIs. Start from the **Round-Robin HTTP Pool** preset.

**Sticky Session Pool** -- Uses round-robin distribution with HTTP cookie-based session persistence. Octavia manages the cookie automatically, routing subsequent requests from the same client to the same backend. Best for applications with server-side session state. Start from the **Sticky Session Pool** preset.

## Works With

- [**OpenStack Load Balancer Listener**](/cloud-catalog/openstack-load-balancer-listener) -- provides the listener ID that this pool serves as the default backend for