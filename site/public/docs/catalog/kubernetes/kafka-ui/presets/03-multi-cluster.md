---
title: "Multi cluster preset"
description: "One console for the whole estate: staging with full console powers (create topics, produce test messages, edit configs) and production locked to observe-only — the per-cluster `read_only` switch is..."
type: "preset"
rank: "03"
presetSlug: "03-multi-cluster"
componentSlug: "kafka-ui"
componentTitle: "Kafka UI"
provider: "kubernetes"
icon: "package"
order: 3
---

# Multi cluster preset

One console for the whole estate: staging with full console powers
(create topics, produce test messages, edit configs) and production
locked to observe-only — the per-cluster `read_only` switch is what
makes a SHARED console safe, because the posture is set where the
risk lives, not globally.

The asymmetry is the teaching point. Staging connects plain and open
because breaking it is cheap; production connects TLS + SCRAM with a
dedicated console credential (a KubernetesKafkaUser by reference —
give it a read-scoped ACL so even a console bug cannot mutate
production data; `read_only` is an app-side switch, defense in depth
rather than the boundary itself).

Cluster `name`s are the console's display and API identifiers — pick
the names your team says out loud. Two replicas because a shared,
depended-on console deserves availability (it is stateless — replicas
are purely that). The login is the single login_form account;
per-person access for a fleet console is the point where OAuth2/LDAP
through `helm_values` earns its complexity.

See [03-multi-cluster.yaml](./03-multi-cluster.yaml) for the
manifest.
