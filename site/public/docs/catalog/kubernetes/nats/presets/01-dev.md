---
title: "Dev preset"
description: "The smallest useful NATS: one server, JetStream on. A single server is a complete JetStream deployment for dev — streams, consumers, KV and object stores all work, and the file-store volume means a..."
type: "preset"
rank: "01"
presetSlug: "01-dev"
componentSlug: "nats"
componentTitle: "NATS"
provider: "kubernetes"
icon: "package"
order: 1
---

# Dev preset

The smallest useful NATS: one server, JetStream on. A single server is
a complete JetStream deployment for dev — streams, consumers, KV and
object stores all work, and the file-store volume means a pod restart
loses nothing. Connect in-cluster through the exported
`client_endpoint`; `kubectl exec` into the nats-box pod for a
pre-wired `nats` CLI (create streams, publish, debug).

Know what the empty seams mean: without `auth`, ANY client that can
reach the Service connects with full access — an in-cluster trust
posture, deliberate for dev, never for anything reachable from
outside. Without `cluster`, there is one server: fine for
availability-tolerant dev work, and replicated (R3) streams are
impossible by definition.

Change first: declare `auth.users` the moment a second team shares the
cluster (passwords are module-generated into the `<name>-auth` Secret
— nothing to invent or store); enable `cluster` with 3 replicas before
anything depends on availability.

See [01-dev.yaml](./01-dev.yaml) for the manifest.
