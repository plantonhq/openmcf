# Standard preset

The default install: one cluster-wide operator watching every
namespace, the stable 2.8.0 chart/operator pairing, chart-default
manager sizing. This is the shape for a cluster you administer —
install once, then declare KubernetesOpenSearch resources anywhere.

Do not use this on a shared cluster where you only own a slice of the
namespaces: the cluster-wide watch needs ClusterRoleBindings, and a
second cluster-wide operator install would fight this one over the
same resources. Use the namespace-scoped preset there instead.

The first thing to change is nothing — the value of this preset is
that the defaults are already the recommended posture. If your
platform requires it, add `node_selector`/`tolerations` to place the
operator pod, or `resources` to satisfy a namespace quota.

See [01-standard.yaml](./01-standard.yaml) for the manifest.
