# Default preset

The standard operator install: the pinned chart (0.120.0 = operator
v0.156.0) into its own `otel-operator` namespace, watching every
namespace on the cluster. Nothing else is set — the chart creates its
own self-signed cert-manager Issuer for the webhook certificate, the
right choice for almost everyone (the certificate only needs to be
trusted by the API server, which cert-manager's CA injection handles).

PREREQUISITE: a `KubernetesCertManager` on the cluster. cert-manager
issues and rotates the operator's webhook serving certificate, and it
keeps the module-owned CRDs' conversion trust current through its CA
injector — the CRDs are retained past the operator's lifetime, so
their trust must be maintained by a running reconciler, not a
certificate embedded once at install time. Without a working
cert-manager the install fails its readiness wait (atomic, 600s) — by
design, rather than surfacing later as collectors that never
reconcile.

What the install owns: the four opentelemetry.io CRDs
(`opentelemetrycollectors`, `instrumentations`, `opampbridges`,
`targetallocators`), applied outside the Helm release and retained on
destroy — removing the operator never deletes the fleet's
declarations. One operator per cluster is the grain; declare
collectors with `KubernetesOtelCollector` against it.

Keep `metadata.name` at 30 characters or fewer — the chart derives a
`<name>-controller-manager-service-cert` Secret and Kubernetes caps
names at 63; the modules fail loudly over the budget.

Change first: nothing, usually. Reach for `replicas: 2` (a warm
standby behind leader election) and `serviceMonitorEnabled` once
Prometheus is on the cluster.

See [01-default.yaml](./01-default.yaml) for the manifest.
