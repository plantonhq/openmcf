---
title: "Presets"
description: "Ready-to-deploy configuration presets for Helm Release"
type: "preset-list"
componentSlug: "helm-release"
componentTitle: "Helm Release"
provider: "kubernetes"
icon: "package"
order: 200
presets:
  - slug: "01-https-repo-chart"
    rank: "01"
    title: "HTTPS Repo Chart"
    excerpt: "This preset installs [podinfo](https://github.com/stefanprodan/podinfo) 6.9.2 from its HTTPS Helm repository, with a `values_yaml` block overriding the chart's defaults. It is the baseline shape for..."
  - slug: "02-oci-registry-chart"
    rank: "02"
    title: "OCI Registry Chart"
    excerpt: "This preset installs [podinfo](https://github.com/stefanprodan/podinfo) 6.9.2 from an OCI registry (`oci://ghcr.io/stefanprodan/charts`) — the same chart as preset 01, pulled the other way charts are..."
  - slug: "03-private-repo-with-secrets"
    rank: "03"
    title: "Private Repo With Secrets"
    excerpt: "This preset is the production shape for a chart from a private repository: repository credentials for the pull, `set_sensitive` for a secret chart value, and the two lifecycle knobs (`atomic`,..."
---

# Helm Release Presets

Ready-to-deploy configuration presets for Helm Release. Each preset is a complete manifest you can copy, customize, and deploy.
