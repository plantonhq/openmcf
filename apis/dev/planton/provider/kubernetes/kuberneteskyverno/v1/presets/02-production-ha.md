# Production HA preset

Kyverno sized for a cluster where admission availability matters. The
admission controller runs three replicas — it sits on the cluster's
WRITE PATH, and with the default fail-closed policy posture an
admission outage blocks every matched resource until the engine
returns; three replicas make that a rolling-restart non-event instead
of an incident. Background scanning gets four workers so drift
reports keep up with a busy cluster, logs go structured for your
pipeline, and every controller ships a ServiceMonitor (requires the
kube-prometheus-stack CRDs — compose KubernetesKubePrometheusStack
first or flip `metrics.serviceMonitor` off).

The requests here are honest floors, not benchmarks: the admission
controller's memory grows with policy count and webhook traffic, the
reports controller's with cluster size. Watch
`kyverno_admission_review_duration_seconds` after the first week and
resize from evidence.

Change first: `config.webhookExcludeNamespaces` for namespaces that
must never wait on policy (CI churn namespaces are the usual case),
and `crds.keepOnUninstall: true` once real policies accumulate — a
reinstall should not delete your policy library.

See [02-production-ha.yaml](./02-production-ha.yaml) for the manifest.
