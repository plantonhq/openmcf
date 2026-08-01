---
title: "OpenBao"
description: "OpenBao deployment documentation"
icon: "package"
order: 100
componentName: "kubernetesopenbao"
---

# OpenBao

The open-source secrets manager: secure secret storage, dynamic
secrets, encryption as a service, leasing and revocation — governed
by the Linux Foundation as the MPL-2.0 fork of Vault, deployed from
the official chart with the seal lifecycle taught instead of papered
over.

## Highlights

- **The seal lifecycle, honest** — fresh servers are sealed and
  NotReady by design until you initialize and unseal them (runtime
  API operations that stay yours); auto-unseal removes the restart
  toil, never the one-time initialization.
- **One mode, typed** — dev, standalone, or HA with integrated Raft
  and module-synthesized `retry_join` for every peer (the chart alone
  ships none — a multi-replica cluster never forms without it).
- **Auto-unseal, keyless-first** — AWS KMS, GCP Cloud KMS, Azure Key
  Vault, or a central Transit engine; static credentials ride a
  module-owned Secret as environment variables, never the config
  ConfigMap.
- **TLS as the composite it is** — listener files, certificate Secret
  mount, and every derived URL and probe switch together; the chart's
  lone `tlsDisable` flag would otherwise produce a plaintext server
  addressed as https.
- **Raft disaster recovery built in** — the snapshot-agent CronJob
  ships scheduled snapshots to any S3-compatible bucket; the
  cluster-wide injector webhook stays off by default, and metrics are
  a deliberate choice, not a surprise.

---

© Planton. Licensed under [Apache-2.0](https://github.com/plantonhq/planton/blob/main/LICENSE).
