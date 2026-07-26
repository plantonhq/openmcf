# Kubernetes Tempo — design notes

## Grain

One resource = one Tempo Helm release (monolithic `tempo` chart,
grafana-community index). The release is named after `metadata.name` and
`fullnameOverride` pins the Service to it, so the exported outputs are
deterministic and several Tempo instances coexist in one cluster.

## The composition seam

- **In:** applications or a `KubernetesOtelCollector` (traces pipeline)
  send OTLP to the `otlp_grpc_endpoint` / `otlp_http_endpoint`.
- **Out:** a `KubernetesGrafana` `tempo` datasource reads the
  `http_endpoint`; the metrics generator remote-writes to a
  `KubernetesKubePrometheusStack` Prometheus.

## Persistence posture

The chart defaults to emptyDir; this component defaults to a
PersistentVolume so traces survive restarts. `ephemeral: true` restores
the chart's throwaway posture. With an object-storage backend, replicas
share the backend and the local volume holds only the WAL.

## Retention units

Tempo parses `compactor.compaction.block_retention` with Go's
`time.ParseDuration`, which accepts only `ns/us/ms/s/m/h` — the spec's CEL
rejects a day unit (`d`) so a manifest can never render a value Tempo
fails to parse at startup.

## Cross-engine parity

The Terraform and Pulumi modules render byte-identical chart values. The
tempo image is split registry+repository (overridden by
`global.imageRegistry`); the tempo-query sidecar is the combined
docker-library form (re-pointed explicitly under an image-registry
override).

## Deliberate exclusions

The `tempo-distributed` microservices chart, and per-receiver / search /
tenant-override tuning beyond the typed surface — reachable through
`helm_values`, never the primary interface.

---

© Planton. Licensed under [Apache-2.0](https://github.com/plantonhq/planton/blob/main/LICENSE).
