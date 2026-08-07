---
title: "Load Balancer"
description: "Load Balancer deployment documentation"
icon: "package"
order: 100
componentName: "hetznercloudloadbalancer"
---

# Hetzner Cloud Load Balancer

Deploys a fully configured load balancer on Hetzner Cloud with services (listeners), backend targets, health checks, and optional private network attachment. Supports HTTP, HTTPS with TLS termination, and TCP pass-through protocols. Backend targets can be specific servers (via server ID), dynamically discovered servers (via label selectors), or external IP addresses. HTTPS services support Hetzner-managed Let's Encrypt certificates, uploaded certificates, sticky sessions, and HTTP-to-HTTPS redirection.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **Load Balancer** -- an `hcloud_load_balancer` resource with the specified type, location, and algorithm
- **Services** -- one `hcloud_load_balancer_service` per entry in the services list, each defining a listener protocol, ports, and optional HTTP configuration
- **Targets** -- one `hcloud_load_balancer_target` per entry across server targets, label selector targets, and IP targets
- **Network Attachment** (optional) -- an `hcloud_load_balancer_network` resource connecting the LB to a private network when the `network` block is specified

## Before You Deploy

### Hetzner Cloud Account

- **A Hetzner Cloud account** with an active project and API token.
- **Backend servers** -- at least one server, label selector, or external IP to serve as a target.
- **Certificates** (for HTTPS) -- a HetznerCloudCertificate (uploaded or managed) for TLS termination.

## Deploy

### Console

Open the deployment store, find **Hetzner Cloud Load Balancer**, and click **Deploy**. The creation wizard walks you through preset selection, environment and connection configuration, and spec fields including services, targets, and networking.

### CLI

Create a manifest and apply it:

```yaml
apiVersion: hetzner-cloud.planton.dev/v1
kind: HetznerCloudLoadBalancer
metadata:
  name: web-lb
  org: acme-corp
  env: prod
spec:
  loadBalancerType: lb11
  location: fsn1
  algorithm: round_robin
  services:
    - protocol: http
      listenPort: 80
      destinationPort: 80
  serverTargets:
    - serverId:
        value: "12345678"
```

```shell
planton apply -f hetznercloud-load-balancer.yaml
```

This creates an lb11 load balancer in Falkenstein with an HTTP service on port 80 and one server target. A Stack Job tracks the provisioning in real time.

### InfraChart

When deploying as part of a load-balanced application, use ValueFromRef to wire servers, certificates, and networks:

```yaml
spec:
  loadBalancerType: lb11
  location: fsn1
  services:
    - protocol: https
      http:
        certificateIds:
          - valueFrom:
              kind: HetznerCloudCertificate
              name: web-cert
              fieldPath: status.outputs.certificate_id
        redirectHttp: true
  serverTargets:
    - serverId:
        valueFrom:
          kind: HetznerCloudServer
          name: web-1
          fieldPath: status.outputs.server_id
      usePrivateIp: true
  network:
    networkId:
      valueFrom:
        kind: HetznerCloudNetwork
        name: main-vpc
        fieldPath: status.outputs.network_id
```

The InfraPipeline resolves the dependency graph, provisioning the network, servers, and certificate before the load balancer.

## Key Configuration

These are the most important decisions when configuring a load balancer. Explore the full field reference in the [API Explorer](#api-explorer) tab.

**Load balancer type** -- The `loadBalancerType` field selects the capacity tier: lb11 (25 targets, 10k connections/s), lb21 (75 targets, 20k connections/s), or lb31 (150 targets, 40k connections/s). Can be resized in place.

**Algorithm** -- The `algorithm` field selects traffic distribution: `round_robin` (even distribution) or `least_connections` (fewest active connections).

**Services** -- Each service defines a listener with a `protocol` (http, https, tcp), `listenPort`, `destinationPort`, and optional HTTP configuration (sticky sessions, certificates, HTTP-to-HTTPS redirect). At least one service is required.

**Targets** -- Backend servers can be added by server ID (`serverTargets`), dynamically by label selector (`labelSelectorTargets`), or as external IPs (`ipTargets`). Use `usePrivateIp: true` to route traffic over a private network.

**Network** -- The `network` block optionally attaches the LB to a private Hetzner Cloud network for private backend communication. Requires a network with a subnet in the LB's network zone.

## Outputs and Dependencies

### What This Component Consumes

| Dependency | Field | ValueFromRef Path |
|------------|-------|-------------------|
| **HetznerCloudServer** | `serverTargets[].serverId` | `status.outputs.server_id` |
| **HetznerCloudCertificate** | `services[].http.certificateIds` | `status.outputs.certificate_id` |
| **HetznerCloudNetwork** | `network.networkId` | `status.outputs.network_id` |

### What This Component Provides

After provisioning, `status.outputs` contains values that downstream Cloud Resources can consume via ValueFromRef:

| Output | Description | Common Downstream Use |
|--------|-------------|----------------------|
| `load_balancer_id` | Hetzner Cloud numeric ID of the load balancer | API operations, monitoring |
| `ipv4_address` | Public IPv4 address of the load balancer | DNS record configuration, application endpoints |
| `ipv6_address` | Public IPv6 address of the load balancer | DNS record configuration |

## Common Patterns

Browse the [Presets](#presets) tab for ready-to-deploy configurations.

**HTTPS web application** -- An lb11 with HTTPS service, managed Let's Encrypt certificate, HTTP-to-HTTPS redirect, and server targets over a private network.

**Private internal LB** -- A load balancer attached to a private network with the public interface disabled, serving internal microservices.

**TCP pass-through** -- A load balancer with TCP services for non-HTTP protocols (databases, custom TCP services).

## Works With

- [**Hetzner Cloud Server**](/cloud-catalog/hetznercloud-server) -- servers are added as backend targets
- [**Hetzner Cloud Certificate**](/cloud-catalog/hetznercloud-certificate) -- certificates enable HTTPS termination
- [**Hetzner Cloud Network**](/cloud-catalog/hetznercloud-network) -- private network attachment for internal traffic
- [**Hetzner Cloud DNS Zone**](/cloud-catalog/hetznercloud-dns-zone) -- DNS records point to the LB's public IP
