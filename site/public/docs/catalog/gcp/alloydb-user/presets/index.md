---
title: "Presets"
description: "Ready-to-deploy configuration presets for AlloyDB User"
type: "preset-list"
componentSlug: "alloydb-user"
componentTitle: "AlloyDB User"
provider: "gcp"
icon: "package"
order: 200
presets:
  - slug: "01-application-user"
    rank: "01"
    title: "Application User (ALLOYDB_BUILT_IN)"
    excerpt: "This preset creates a classic username/password AlloyDB user for an application service."
  - slug: "02-iam-user"
    rank: "02"
    title: "IAM User (ALLOYDB_IAM_USER)"
    excerpt: "This preset maps a GCP IAM principal (here, a service account email) to an AlloyDB database user with no stored password."
---

# AlloyDB User Presets

Ready-to-deploy configuration presets for AlloyDB User. Each preset is a complete manifest you can copy, customize, and deploy.
