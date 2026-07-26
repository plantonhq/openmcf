---
title: "Namespace-fenced preset"
description: "The multi-tenant posture: the operator watches only the namespaces that are allowed to hold RabbitMQ clusters. The upstream default is the opposite — empty `watch_namespaces` watches EVERYTHING..."
type: "preset"
rank: "02"
presetSlug: "02-namespace-fenced"
componentSlug: "rabbitmq-operator"
componentTitle: "RabbitMQ Operator"
provider: "kubernetes"
icon: "package"
order: 2
---

# Namespace-fenced preset

The multi-tenant posture: the operator watches only the namespaces
that are allowed to hold RabbitMQ clusters. The upstream default is
the opposite — empty `watch_namespaces` watches EVERYTHING (unlike
chart-scoped operators such as the Altinity operator, whose defaults
fence to the install namespace) — so on a shared cluster where teams
must not declare brokers wherever they please, the fence is an
explicit choice.

The list renders as the operator's `OPERATOR_SCOPE_NAMESPACE`
environment variable (comma-separated). Every namespace that will
hold KubernetesRabbitMq resources must be covered — a fenced operator
silently ignores RabbitmqCluster resources everywhere else, which
presents as clusters that never reconcile.

Note what the fence does NOT change: this is still one install per
cluster (the admission webhooks are cluster-scoped singletons with
fixed names, and the install still lives in the fixed
`rabbitmq-system` namespace), and the RabbitmqCluster CRD still
deletes with this resource. The fence is about where clusters may
live, not about per-team operator copies.

The first thing to change is the namespace list — make it exactly the
namespaces where RabbitMQ clusters belong.

See [02-namespace-fenced.yaml](./02-namespace-fenced.yaml) for the
manifest.
