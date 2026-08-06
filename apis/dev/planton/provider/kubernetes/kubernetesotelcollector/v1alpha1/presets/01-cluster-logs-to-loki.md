# Cluster-logs-to-Loki preset

The per-node log pipeline: daemonset mode puts one collector on every
node, the filelog receiver tails every container's log files under
`/var/log/pods` (the standard `container` operator parses the runtime
format and extracts the Kubernetes metadata from the file path), the
`k8sattributes` processor enriches records with pod/namespace/workload
attributes, and the `otlphttp` exporter ships everything to a Loki
gateway's `/otlp` route. The standard `otlp` receiver rides along so
node-local applications can push their own telemetry too — and so the
exported OTLP endpoints stay valid.

PREREQUISITE: a `KubernetesOtelOperator` on the cluster. The
`k8sattributes` processor reads cluster state, which needs RBAC beyond
the operator's default ServiceAccount — compose a
`KubernetesServiceAccount` + `KubernetesRbac` for the
`otel-logs-collector` account this preset names in `serviceAccount`.

The volumes are the daemonset log-collection pattern: `/var/log/pods`
mounted read-only from the host (the receiver only reads), plus a
writable hostPath at `/var/lib/otelcol-checkpoints` for the filelog
receiver's checkpoint state, so a restarted collector resumes where it
left off. The control-plane toleration covers every node — remove it
if control-plane logs should stay uncollected.

The `podSecurityContext.runAsUser: 0` is load-bearing, not a
convenience: container runtimes write pod log files readable only by
root, and the default collector image runs as a non-root user that
cannot open them — without it the filelog receiver reports permission
errors and ships nothing.

Sizing discipline: the `memory_limiter` (400 MiB limit, 100 MiB spike)
sheds load visibly instead of OOMing — if you add container resources,
keep the memory limit and `limit_mib` in agreement.

Change first: the `otlphttp` endpoint's host — point it at your
`KubernetesLoki`'s exported gateway service (Loki ingests OTLP at the
gateway's `/otlp` route).

See [01-cluster-logs-to-loki.yaml](./01-cluster-logs-to-loki.yaml) for
the manifest.
