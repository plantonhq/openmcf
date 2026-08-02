---
title: "Federated warehouse"
description: "The production posture: a composed PostgreSQL warehouse queryable through Trino (and JOIN-able against any other catalog), autoscaled workers (up to 6 on CPU) that drain running queries before..."
type: "preset"
rank: "02"
presetSlug: "02-federated-warehouse"
componentSlug: "trino"
componentTitle: "Trino"
provider: "kubernetes"
icon: "package"
order: 2
---

# Federated warehouse

The production posture: a composed PostgreSQL warehouse queryable
through Trino (and JOIN-able against any other catalog), autoscaled
workers (up to 6 on CPU) that drain running queries before
termination, percent-based JVM heaps that follow the container limits,
and JMX metrics flowing to the Prometheus operator.

The catalog references a `KubernetesPostgres` named `warehouse-pg` —
its read-write Service becomes the connection host and its
application-user Secret supplies the password, which reaches Trino as
an environment variable referenced through `${ENV:...}` (Trino's own
secrets mechanism); nothing credential-bearing renders into any
ConfigMap. The samples are off: only declared data sources exist.

Point BI tools (a KubernetesSuperset composes naturally) at the
exported coordinator endpoint with the generated admin credential.
