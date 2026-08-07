# Standard Istio Mesh (Sidecar)

This preset installs the Istio control plane with the classic sidecar data plane and chart-default sizing. istiod autoscales via the chart's HPA (min 1 / max 5 at 80% CPU).

## When to Use

- You want mTLS, traffic policy, and mesh observability with the most mature, widely deployed data plane architecture
- You are starting a mesh and have no reason to deviate from defaults

## Key Configuration Choices

- **Sidecar mode** (`dataplaneMode: sidecar`) — namespaces opt in to injection with the `istio-injection=enabled` label; pods created after that get a proxy container
- **Namespace** (`istio-system`) — the conventional control-plane namespace, created by the resource
- **Chart-default sizing** — the HPA owns istiod's replica count; proxy defaults apply per injected pod

## Placeholders to Replace

No placeholders — this preset is directly deployable with sensible defaults.

## Related Presets

- **02-ambient** — sidecar-less data plane (per-node ztunnel + optional waypoints)
- **03-production-sidecar** — production hardening: HA istiod, egress lockdown, node-level CNI agent
