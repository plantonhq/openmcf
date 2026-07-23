---
title: "Presets"
description: "Ready-to-deploy configuration presets for Velero"
type: "preset-list"
componentSlug: "velero"
componentTitle: "Velero"
provider: "kubernetes"
icon: "package"
order: 200
presets:
  - slug: "01-eks-s3-irsa"
    rank: "01"
    title: "EKS S3 IRSA"
    excerpt: "This preset installs Velero on EKS in the standard production posture: backups land in an S3 bucket reached keylessly via IRSA, volume data is captured through CSI snapshots (the modern path with the..."
  - slug: "02-gke-gcs-workload-identity"
    rank: "02"
    title: "GKE GCS Workload Identity"
    excerpt: "This preset installs Velero on GKE with backups landing in a Google Cloud Storage bucket reached keylessly via Workload Identity, and the kopia-based node-agent deployed for file-system backup of..."
  - slug: "03-minio-self-contained"
    rank: "03"
    title: "MinIO Self-Contained"
    excerpt: "This preset installs Velero against an S3-COMPATIBLE object store — an in-cluster MinIO by default — with no cloud account involved anywhere: the store is reached by endpoint URL with path-style..."
---

# Velero Presets

Ready-to-deploy configuration presets for Velero. Each preset is a complete manifest you can copy, customize, and deploy.
