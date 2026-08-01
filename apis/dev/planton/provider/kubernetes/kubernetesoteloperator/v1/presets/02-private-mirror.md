# Private-mirror preset

The air-gapped posture: both image seams mirrored, because they are
DIFFERENT seams. `imageRegistry` replaces only the registry part of
the operator's own manager image (the path stays the upstream
`open-telemetry/opentelemetry-operator/opentelemetry-operator`) — the
one image this component's pods pull. `defaultCollectorImage` mirrors
what the operator INJECTS into `OpenTelemetryCollector` CRs that
declare no image — collector pods pull that one, fleet-wide. Mirroring
the operator without mirroring the collector default leaves every
collector pod reaching for ghcr.io.

Keep the tag on `defaultCollectorImage`: the chart renders the
operator's `--collector-image` flag only when both repository and tag
are present (a repository-only value still renders — it deep-merges
with the chart's default tag — but pinning the tag keeps the mirror
and the flag explicit).

`imagePullSecrets` names Secrets that must ALREADY exist in the
namespace — the component references them, it does not create them.
`replicas: 2` gives a warm standby (leader election picks one active
reconciler), and `serviceMonitorEnabled` exposes the operator's own
metrics to Prometheus — it requires the monitoring.coreos.com CRDs on
the cluster (`KubernetesKubePrometheusStack` first).

PREREQUISITE: a `KubernetesCertManager` on the cluster, as with every
install of this component — the webhook serving certificate and the
retained CRDs' conversion trust both ride it. An unpullable mirror
image fails the install's readiness wait (atomic, 600s) at apply time.

Change first: the `mirror.example.com` registry and the `mirror-pull`
Secret name; keep the collector image tag paired with the operator
version (chart 0.120.0 = operator v0.156.0).

See [02-private-mirror.yaml](./02-private-mirror.yaml) for the
manifest.
