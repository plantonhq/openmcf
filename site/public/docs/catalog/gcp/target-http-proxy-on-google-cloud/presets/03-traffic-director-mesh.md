---
title: "Traffic Director Mesh Frontend"
description: "A target HTTP proxy for Traffic Director (service mesh): `proxyBind: true` binds the proxy to the mesh's private IPs instead of Google's edge, and the forwarding rule that references it uses the..."
type: "preset"
rank: "03"
presetSlug: "03-traffic-director-mesh"
componentSlug: "target-http-proxy-on-google-cloud"
componentTitle: "Target HTTP Proxy on Google Cloud"
provider: "gcp"
icon: "package"
order: 3
---

# Traffic Director Mesh Frontend

A target HTTP proxy for Traffic Director (service mesh): `proxyBind: true` binds the proxy to the mesh's private IPs instead of Google's edge, and the forwarding rule that references it uses the `INTERNAL_SELF_MANAGED` scheme.

## When to Use

- xDS-driven service mesh routing where sidecar or gateway proxies consume the configuration
- Advanced traffic management (weighted canaries, header routing) for east-west traffic inside a VPC

## Remix Notes

- The paired `GcpGlobalForwardingRule` must set `loadBalancingScheme: INTERNAL_SELF_MANAGED`, a VPC `network`, and typically the VIP `0.0.0.0` with `metadataFilters` scoping which xDS clients receive it.
- `proxyBind` is immutable — flipping it later recreates the proxy.
- The URL map's `routeRules` with header/metadata matching is where mesh routing policy actually lives.
