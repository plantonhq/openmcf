# Ambient Mesh (Sidecar-less)

This preset installs the Istio control plane in ambient mode: no sidecars — a per-node ztunnel DaemonSet carries mTLS and L4 policy, and optional waypoint proxies add L7 where needed. The istio-cni node agent and ztunnel install alongside istiod.

## When to Use

- You want mesh mTLS and L4 policy without a proxy container in every pod (lower steady-state cost, no pod restarts to join the mesh)
- Your L7 needs are selective — waypoints add L7 per namespace or service rather than everywhere

## Key Configuration Choices

- **Ambient mode** (`dataplaneMode: ambient`) — namespaces enroll with the `istio.io/dataplane-mode=ambient` label; switching modes later is a workload-by-workload migration
- **ztunnel sizing** — the chart default request (200m CPU / 512Mi memory) already covers large clusters; scale with cluster size and connection volume
- **NodePort gateways** (`gatewayDefaults.serviceType: NodePort`) — for clusters without a cloud load-balancer controller (bare metal, kind, k3s); drop this on EKS/GKE/AKS to keep the LoadBalancer upstream default

## Placeholders to Replace

No placeholders — this preset is directly deployable with sensible defaults.

## Related Presets

- **01-standard** — the classic sidecar data plane
- **03-production-sidecar** — production hardening on the sidecar architecture
