---
title: "NATS"
description: "NATS deployment documentation"
icon: "package"
order: 100
componentName: "kubernetesnats"
---

# NATS

The lightweight, high-speed messaging system. Pub/sub, request/reply
and queue groups measured in microseconds, plus JetStream persistence
— streams, consumers, key-value and object stores — when messages must
survive restarts. One small server binary carries all of it.

## Highlights

- **Persistent by default** — JetStream is ON with a file-store volume
  per server; published messages outlive pod restarts. A single server
  is a complete dev deployment; 3 clustered servers keep replicated
  streams available through a pod loss.
- **Credentials nobody typed** — declare users (flat, or grouped into
  isolated multi-tenant accounts) with subject-level permissions;
  passwords are module-generated into one Secret, wired to the server
  through environment expansion, and never rendered into config or
  values.
- **Every listener, typed** — WebSocket for browsers, MQTT for IoT
  devices bridged into JetStream, leafnodes for edge servers, TLS from
  a cert-manager Secret, and Prometheus metrics with an optional
  PodMonitor.
- **Batteries included** — the nats-box pod ships a pre-wired `nats`
  CLI for creating streams and debugging.
- **Clean lifecycle** — no CRDs, no bundled operators; destroy leaves
  nothing behind.

---

© Planton. Licensed under [Apache-2.0](https://github.com/plantonhq/planton/blob/main/LICENSE).
