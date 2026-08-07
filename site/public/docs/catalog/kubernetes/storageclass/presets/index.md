---
title: "Presets"
description: "Ready-to-deploy configuration presets for StorageClass"
type: "preset-list"
componentSlug: "storageclass"
componentTitle: "StorageClass"
provider: "kubernetes"
icon: "package"
order: 200
presets:
  - slug: "01-aws-ebs-gp3"
    rank: "01"
    title: "AWS EBS gp3"
    excerpt: "This preset creates the standard general-purpose SSD class for EKS: encrypted gp3 volumes through the AWS EBS CSI driver, provisioned in the zone the consuming pod schedules into, and expandable..."
  - slug: "02-gcp-pd-ssd"
    rank: "02"
    title: "GCP PD SSD"
    excerpt: "This preset creates the performance SSD class for GKE: SSD persistent disks through the GCE Persistent Disk CSI driver, provisioned in the zone the consuming pod schedules into, and expandable after..."
  - slug: "03-azure-premium"
    rank: "03"
    title: "Azure Premium"
    excerpt: "This preset creates the premium SSD class for AKS: Premium SSD managed disks (locally redundant) through the Azure Disk CSI driver, provisioned in the zone the consuming pod schedules into, and..."
---

# StorageClass Presets

Ready-to-deploy configuration presets for StorageClass. Each preset is a complete manifest you can copy, customize, and deploy.
