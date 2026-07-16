---
title: "Presets"
description: "Ready-to-deploy configuration presets for Artifact Registry Repository"
type: "preset-list"
componentSlug: "artifact-registry-repository"
componentTitle: "Artifact Registry Repository"
provider: "gcp"
icon: "package"
order: 200
presets:
  - slug: "01-docker-standard"
    rank: "01"
    title: "Standard Docker Repository"
    excerpt: "The workhorse of container delivery: a private Docker repository with immutable tags, self-cleaning storage, and an additive pull grant for the runtime service account."
  - slug: "02-remote-docker-hub-cache"
    rank: "02"
    title: "Remote Docker Hub Cache"
    excerpt: "A pull-through cache of Docker Hub: base images are fetched once, cached in your project, and served from GCP thereafter — insulating builds and node pools from upstream outages and anonymous-pull..."
  - slug: "03-virtual-aggregation"
    rank: "03"
    title: "Virtual Aggregation Endpoint"
    excerpt: "One pull URL for everything: a virtual repository that serves first-party images from the team's standard repository and falls through to the Docker Hub cache for everything else — consumers never..."
---

# Artifact Registry Repository Presets

Ready-to-deploy configuration presets for Artifact Registry Repository. Each preset is a complete manifest you can copy, customize, and deploy.
