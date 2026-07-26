---
title: "Standard preset"
description: "The one-per-cluster engine install — and the spec is empty, because the release manifest's own defaults ARE the production-standard posture: the operator watching ALL namespaces (the upstream..."
type: "preset"
rank: "01"
presetSlug: "01-standard"
componentSlug: "rabbitmq-operator"
componentTitle: "RabbitMQ Operator"
provider: "kubernetes"
icon: "package"
order: 1
---

# Standard preset

The one-per-cluster engine install — and the spec is empty, because
the release manifest's own defaults ARE the production-standard
posture: the operator watching ALL namespaces (the upstream default),
the pinned operator image, 200m CPU / 500Mi memory for both requests
and limits, one replica. Declare it once; every KubernetesRabbitMq
resource in the cluster is reconciled by this single install.

Nothing here is a placeholder to replace. There is no namespace field
(the manifest installs into its fixed `rabbitmq-system`, and exactly
one install per cluster is the upstream contract) and no version
field (the operator and its CRD schema are pinned to the release the
catalog's typed SDK was generated against).

Two things must be true of the cluster before this applies. First,
cert-manager must be running — the admission webhooks' serving
certificate is a cert-manager Certificate, and without it every
RabbitmqCluster admission fails. Second, and unlike the chart-based
sibling operators: the RabbitmqCluster CRD is part of the applied
manifest and DELETES with this resource — destroying the operator
cascade-deletes every RabbitmqCluster on the cluster, so never
destroy it while KubernetesRabbitMq resources exist.

On multi-tenant clusters where the watch scope should be fenced, use
the namespace-fenced preset instead.

See [01-standard.yaml](./01-standard.yaml) for the manifest.
