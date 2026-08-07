---
title: "Presets"
description: "Ready-to-deploy configuration presets for CloudNativePG Operator"
type: "preset-list"
componentSlug: "cloudnativepg-operator"
componentTitle: "CloudNativePG Operator"
provider: "kubernetes"
icon: "package"
order: 200
presets:
  - slug: "01-standard"
    rank: "01"
    title: "Standard"
    excerpt: "This preset installs the CloudNativePG operator in its standard posture: the operator release alone, cluster-wide watch scope, pinned chart version, sized and prioritized for a production control..."
  - slug: "02-with-backup-plugin"
    rank: "02"
    title: "With Backup Plugin"
    excerpt: "This preset installs the CloudNativePG operator PLUS the Barman Cloud plugin — the backup-capable posture every KubernetesPostgres backup block depends on. The plugin is a separate Helm release..."
---

# CloudNativePG Operator Presets

Ready-to-deploy configuration presets for CloudNativePG Operator. Each preset is a complete manifest you can copy, customize, and deploy.
