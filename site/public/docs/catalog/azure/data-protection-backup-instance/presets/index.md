---
title: "Presets"
description: "Ready-to-deploy configuration presets for Data Protection Backup Instance"
type: "preset-list"
componentSlug: "data-protection-backup-instance"
componentTitle: "Data Protection Backup Instance"
provider: "azure"
icon: "package"
order: 200
presets:
  - slug: "01-disk-backup"
    rank: "01"
    title: "Disk Backup"
    excerpt: "This preset protects one managed disk with the modern Data Protection service: scheduled incremental snapshots, kept in a dedicated snapshot resource group, retained per the referenced disk policy."
  - slug: "02-aks-cluster-backup"
    rank: "02"
    title: "AKS Cluster Backup"
    excerpt: "This preset protects an AKS cluster's workloads: scheduled cluster backups with the plumbing namespace excluded and persistent-volume snapshots enabled, retained per the referenced Kubernetes policy."
  - slug: "03-blob-backup"
    rank: "03"
    title: "Blob Backup"
    excerpt: "This preset protects a storage account's blob data with vault-tier backups of named containers (plus whatever operational-tier protection the referenced blob policy configures)."
---

# Data Protection Backup Instance Presets

Ready-to-deploy configuration presets for Data Protection Backup Instance. Each preset is a complete manifest you can copy, customize, and deploy.
