# Cluster Standard

This preset installs the KEDA engine cluster-wide with the upstream
defaults: the operator, the external-metrics API server, and the admission
webhooks in a dedicated `keda` namespace, watching ScaledObjects in ALL
namespaces, with the CRDs installed by the release but kept on uninstall.
KEDA is one installation per cluster — it registers the cluster-wide
`v1beta1.external.metrics.k8s.io` APIService, a singleton.

## When to Use

- Any cluster that needs event-driven autoscaling (queue depth, stream
  lag, cron schedules, cloud metrics — or scale-to-zero, which plain HPA
  cannot do) and has no cloud-identity or HA requirements yet
- The 30-second choice: this is the standard first KEDA installation

## Key Configuration Choices

- **`namespace: keda` + `createNamespace: true`** — the upstream
  convention, in a namespace this resource creates and owns
- **`crds.install: true` + `keepOnUninstall: true`** (spec defaults, made
  explicit) — uninstalling the release does NOT cascade-delete every
  ScaledObject/ScaledJob/TriggerAuthentication in the cluster; the keep
  annotation is the guard rail
- **`watchNamespace` empty** — cluster-wide watching, the normal posture
- **Single replicas, chart-default resources, operator-managed
  certificates, Prometheus off** — the zero-dependency baseline; harden
  with **03-ha-production** when KEDA becomes load-bearing

The scaling declarations themselves (ScaledObject, ScaledJob,
TriggerAuthentication) are per-workload custom resources deployed
alongside the workloads they scale — this component installs the engine.

## Placeholders to Replace

No placeholders — this preset is directly deployable.

## Related Presets

- **02-eks-irsa-scalers** — EKS clusters whose scalers read AWS metric
  sources (SQS, CloudWatch, Kinesis) via IRSA
- **03-ha-production** — HA replicas, resources, cert-manager certificates,
  Prometheus telemetry
