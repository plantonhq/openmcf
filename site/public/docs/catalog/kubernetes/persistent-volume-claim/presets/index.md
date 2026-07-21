---
title: "Presets"
description: "Ready-to-deploy configuration presets for Persistent Volume Claim"
type: "preset-list"
componentSlug: "persistent-volume-claim"
componentTitle: "Persistent Volume Claim"
provider: "kubernetes"
icon: "package"
order: 200
presets:
  - slug: "01-dynamic-default-class"
    rank: "01"
    title: "Dynamic Default Class"
    excerpt: "This preset creates the simplest useful claim: a name and a size, provisioned through the cluster's DEFAULT StorageClass. Omitting `storage_class_name` is itself the request — the claim receives..."
  - slug: "02-pinned-class"
    rank: "02"
    title: "Pinned Class"
    excerpt: "This preset is the production claim shape: the StorageClass is named explicitly instead of inherited from the cluster default. The default class is a moving target — a cluster upgrade or a platform..."
  - slug: "03-shared-rwm"
    rank: "03"
    title: "Shared RWM"
    excerpt: "This preset creates a volume that many pods — across many nodes — mount read-write simultaneously: `ReadWriteMany`. It is the shape for shared asset directories, upload buffers, and legacy..."
  - slug: "04-clone-from-claim"
    rank: "04"
    title: "Clone From Claim"
    excerpt: "This preset creates a new volume populated as a COPY of an existing PersistentVolumeClaim's data, instead of provisioned empty. The `data_source` names a same-namespace source claim; the CSI driver..."
---

# Persistent Volume Claim Presets

Ready-to-deploy configuration presets for Persistent Volume Claim. Each preset is a complete manifest you can copy, customize, and deploy.
