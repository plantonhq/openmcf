---
title: "JupyterHub"
description: "JupyterHub deployment documentation"
icon: "package"
order: 100
componentName: "kubernetesjupyterhub"
---

# JupyterHub

The multi-user notebook platform. Every teammate signs in and gets
their own JupyterLab server with its own persistent home volume —
spawned on demand, culled when idle, scheduled to keep your cluster
bill honest.

## Highlights

- **Secured by default** — the chart's open-door default (any
  username, no password) never ships; sign-in is a generated shared
  password or your identity provider (GitHub, Google, any OIDC —
  Keycloak composes naturally), with secrets riding environment
  indirection, never rendered values.
- **A real machine menu** — per-user CPU/memory guarantees and limits,
  per-user home volumes, and a spawn-time profile list ("Small",
  "GPU workstation", "Big-memory ETL") backed by any notebook image.
- **Capacity that scales down** — the packing user-scheduler, warm
  placeholder pods and install-time image pre-pulling keep spawns
  instant while autoscalers reclaim idle nodes.
- **Composable state** — hub state on a PVC by default or in a
  composed KubernetesPostgres; exposure composes over the exported
  front-door Service handle.

## Operational notes

User home volumes (`claim-<username>` PVCs) are created at runtime by
the hub and deliberately survive destroy — they are your users' work.
The idle culler stops abandoned servers after an hour by default; the
active-server limit caps the fleet.
