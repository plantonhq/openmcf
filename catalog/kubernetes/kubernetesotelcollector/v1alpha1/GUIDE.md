# KubernetesOtelCollector Guide

The judgment this guide carries: the collector's value is its pipeline
config and the MODE you run it in — and two composition needs (credentials
and cluster-read RBAC) are easy to get wrong in ways that leak secrets or
silently collect nothing.

## Pick the mode by what you are collecting

`deployment` (default) is the scalable gateway/fan-in collector;
`daemonset` is per-node collection (log files, host and kubelet metrics —
this is how cluster logs reach a KubernetesLoki, paired with hostPath
`volumes`); `statefulset` is for stable identities (target allocator,
persistent queues); `sidecar` injects into annotated pods and creates no
standalone workload (the field doc on [reference.md](reference.md)). The
mode is an architecture decision — a gateway where you needed a daemonset
collects none of the per-node telemetry you wanted.

## Never inline credentials; grant RBAC for cluster receivers

Two traps the spec's own docs call out: load exporter credentials as env
from existing Secrets and reference them as `${env:VAR}` in the config, so
tokens never land in the rendered ConfigMap. And receivers that read
cluster state (k8s_events, kubeletstats, k8s_cluster, enriched filelog)
need RBAC beyond the default ServiceAccount — compose a
KubernetesServiceAccount + KubernetesRbac and set `serviceAccount`, or the
receiver silently collects nothing.

## Operator prerequisite

KubernetesOtelOperator is the registry prerequisite, watching cluster-wide
([operator-prerequisite pattern](../../../_patterns/operator-prerequisite.md);
its [guide](../../kubernetesoteloperator/v1alpha1/GUIDE.md) has the watch
judgment). Wire `spec.namespace` to a dedicated KubernetesNamespace, not
`createNamespace: true`
([namespace-ownership pattern](../../../_patterns/namespace-ownership.md)).

## Pairs well with

- KubernetesOtelOperator — required.
- KubernetesLoki / KubernetesTempo — the log and trace backends collector
  pipelines export to.
- KubernetesServiceAccount + KubernetesRbac — the permissions cluster-read
  receivers need.
