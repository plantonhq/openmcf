# Fenced HA preset

The multi-tenant, standby-backed posture: the operator watches ONLY
the listed namespaces — the chart scopes its RBAC AND the admission
webhook's namespaceSelector to exactly this list, so Flink
declarations outside it are ignored without an error (a missing
namespace here looks like a deployment that never reconciles) — and
`replicas: 2` runs a warm standby behind leader election, which the
module configures automatically (the chart refuses multi-replica
installs without it).

`operatorConfig` entries are Flink's own config key space
(`kubernetes.operator.*`) appended over the chart defaults — and they
become CLUSTER-WIDE defaults for every FlinkDeployment this operator
manages. Per-pipeline configuration belongs on each
`KubernetesFlinkDeployment`, not here.

PREREQUISITE: a `KubernetesCertManager`, as with every webhook-enabled
install of this component (fail-closed webhooks, cert-manager-issued
certificate, no self-signed fallback — the default preset explains
the lifecycle).

Change first: the namespace list — it is the fence, the RBAC scope,
and the webhook scope in one value.

See [02-fenced-ha.yaml](./02-fenced-ha.yaml) for the manifest.
