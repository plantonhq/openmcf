# Kubernetes Loki — design notes

## Grain

One resource = one Loki Helm release (chart `loki`, grafana-community
index). The release is named after `metadata.name` and `fullnameOverride`
pins every child name (gateway Service, Loki Services, memberlist) to it,
so several Loki instances coexist in one cluster and the exported outputs
are deterministic.

## The composition seam

Loki sits between a shipper and a reader:

- **In:** a `KubernetesOtelCollector` (daemonset, cluster-logs pipeline)
  exports to the `gateway_endpoint` / `otlp_push_endpoint`.
- **Out:** a `KubernetesGrafana` `loki` datasource reads the
  `gateway_endpoint`; the ruler fires alerts at a
  `KubernetesKubePrometheusStack` Alertmanager.

## Secret discipline

Object-store credentials and gateway tenant material are references to
existing Secrets. Credentials ride `defaults.extraEnv` secretKeyRefs and
`${VAR}` config placeholders expanded at process start
(`-config.expand-env=true`); tenant passwords are bcrypt hashes the
gateway consumes directly. No credential ever lands in rendered chart
values.

## Cache footprint — the one knob small clusters cannot skip

The chart's memcached defaults are production-scale, and it requests
container memory at 1.2× the allocated size: the default chunks cache
(8192MB allocated) renders a **9830Mi** memory request. On a node with
less than ~10Gi allocatable the cache pod stays Pending forever, and
because the install is atomic the WHOLE release rolls back after the
timeout — the failure reads as "Loki won't install," not "cache too
big" (verified live). On any small or dev cluster set
`caching.chunks_cache_memory_mb` (and `results_cache_memory_mb`)
explicitly; 128–1024MB serves light query loads. The results cache
default (1024MB → ~1229Mi requested) usually fits; the chunks cache is
the one that will not.

## Cross-engine parity

The Terraform and Pulumi modules render byte-identical chart values from
the same typed spec. The chart mixes image forms — the Loki/gateway/canary
images are split registry+repository (overridden by `global.imageRegistry`)
while the memcached caches are the combined docker-library form (re-pointed
explicitly) — and both engines handle both.

## Deliberate exclusions

The microservices (Distributed) deployment mode, transitional migration
modes, the bundled MinIO subchart (deprecated by the chart itself), and
per-component overrides beyond the typed surface (bloom filters, zone-aware
rollouts) — all reachable through `helm_values`, none the primary
interface.

---

© Planton. Licensed under [Apache-2.0](https://github.com/plantonhq/planton/blob/main/LICENSE).
