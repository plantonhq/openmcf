---
title: "Presets"
description: "Ready-to-deploy configuration presets for Machine Learning Datastore"
type: "preset-list"
componentSlug: "machine-learning-datastore"
componentTitle: "Machine Learning Datastore"
provider: "azure"
icon: "package"
order: 200
presets:
  - slug: "01-blob-training-data"
    rank: "01"
    title: "Blob Training Data"
    excerpt: "This preset registers a blob-container datastore under workspace-identity auth -- the credential-free posture: no key or SAS in the manifest, one role assignment on the container instead."
  - slug: "02-datalake-filesystem"
    rank: "02"
    title: "Data Lake Filesystem"
    excerpt: "This preset registers a Data Lake Gen2 filesystem datastore with service-principal auth -- the lakehouse access shape: a Gen2 filesystem on a hierarchical-namespace account, read through an app..."
---

# Machine Learning Datastore Presets

Ready-to-deploy configuration presets for Machine Learning Datastore. Each preset is a complete manifest you can copy, customize, and deploy.
