---
title: "Single cluster readonly preset"
description: "The safe first console: one Kafka cluster wired in with `read_only` on, so the console can browse topics, messages, and consumer lag but cannot create, delete, produce, or edit anything — an app-side..."
type: "preset"
rank: "01"
presetSlug: "01-single-cluster-readonly"
componentSlug: "kafka-ui"
componentTitle: "Kafka UI"
provider: "kubernetes"
icon: "package"
order: 1
---

# Single cluster readonly preset

The safe first console: one Kafka cluster wired in with `read_only`
on, so the console can browse topics, messages, and consumer lag but
cannot create, delete, produce, or edit anything — an app-side
switch that makes the console harmless to point at production.

The posture is deliberately paired: no `auth` block (the console is
open to anyone who can reach the Service) AND no exposure (ClusterIP,
nothing composed on top). Those two only balance together — reach the
console through the exported port-forward command. The moment an
Ingress or Gateway route enters the picture, add the full-stack
preset's `login_form` account first.

`read_only` is not a substitute for least-privilege credentials: on
an unauthenticated dev listener it is the only guard, which is fine
for observing; against SASL-secured clusters, pair it with a
scoped KubernetesKafkaUser as in the multi-cluster preset.

See
[01-single-cluster-readonly.yaml](./01-single-cluster-readonly.yaml)
for the manifest.
