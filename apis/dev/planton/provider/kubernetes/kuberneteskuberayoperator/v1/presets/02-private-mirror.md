# Private-mirror preset

The air-gapped posture for the operator: `imageRegistry` replaces only
the registry part of the operator's own image (the path stays the
upstream `kuberay/operator`; the default registry is quay.io) — the
one image THIS component's pods pull.

Know the seam boundary: Ray CLUSTER images are a different seam
entirely — each `KubernetesRayCluster` declares its own `image`
(default `rayproject/ray:<rayVersion>`, Docker Hub), and mirroring the
operator does nothing for them. Mirror the Ray image on each cluster
declaration.

`serviceMonitorEnabled` exposes the operator's control-plane metrics
to Prometheus — it requires the monitoring.coreos.com CRDs on the
cluster (`KubernetesKubePrometheusStack` first; the install FAILS
without them, by upstream design rather than a silent skip).

Change first: the `mirror.example.com` registry; drop
`serviceMonitorEnabled` if Prometheus is not on the cluster yet.

See [02-private-mirror.yaml](./02-private-mirror.yaml) for the
manifest.
