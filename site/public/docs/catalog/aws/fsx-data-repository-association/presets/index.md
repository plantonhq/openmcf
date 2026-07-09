---
title: "Presets"
description: "Ready-to-deploy configuration presets for FSx Data Repository Association"
type: "preset-list"
componentSlug: "fsx-data-repository-association"
componentTitle: "FSx Data Repository Association"
provider: "aws"
icon: "package"
order: 200
presets:
  - slug: "01-training-data-import"
    rank: "01"
    title: "Training Data Import Link"
    excerpt: "Read-side association: an S3 training-data prefix appears as `/datasets` on the Lustre file system, tracks the bucket continuously, and hydrates the existing contents at creation."
  - slug: "02-results-export"
    rank: "02"
    title: "Results Export Link"
    excerpt: "Write-side association: everything jobs write under `/output` on the Lustre file system flows back to an S3 results prefix automatically."
---

# FSx Data Repository Association Presets

Ready-to-deploy configuration presets for FSx Data Repository Association. Each preset is a complete manifest you can copy, customize, and deploy.
