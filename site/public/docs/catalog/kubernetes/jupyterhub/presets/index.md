---
title: "Presets"
description: "Ready-to-deploy configuration presets for JupyterHub"
type: "preset-list"
componentSlug: "jupyterhub"
componentTitle: "JupyterHub"
provider: "kubernetes"
icon: "package"
order: 200
presets:
  - slug: "01-team-notebooks"
    rank: "01"
    title: "Team notebooks — shared-password sign-in"
    excerpt: "The fastest way to a real multi-user notebook platform: every teammate signs in with their own username and ONE shared password (module-generated into the `team-notebooks-auth` Secret — the chart's..."
  - slug: "02-production-oidc"
    rank: "02"
    title: "Production — OIDC sign-in on PostgreSQL"
    excerpt: "The durable, org-scale shape: sign-in delegates to your identity provider over OIDC (a KubernetesKeycloak realm's endpoints slot straight in — Okta, Auth0 and Dex work identically), hub state lives..."
---

# JupyterHub Presets

Ready-to-deploy configuration presets for JupyterHub. Each preset is a complete manifest you can copy, customize, and deploy.
