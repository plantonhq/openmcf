---
title: "Remote Docker Hub Cache"
description: "A pull-through cache of Docker Hub: base images are fetched once, cached in your project, and served from GCP thereafter — insulating builds and node pools from upstream outages and anonymous-pull..."
type: "preset"
rank: "02"
presetSlug: "02-remote-docker-hub-cache"
componentSlug: "artifact-registry-repo"
componentTitle: "Artifact Registry Repo"
provider: "gcp"
icon: "package"
order: 2
---

# Remote Docker Hub Cache

A pull-through cache of Docker Hub: base images are fetched once, cached
in your project, and served from GCP thereafter — insulating builds and
node pools from upstream outages and anonymous-pull rate limits.

## What this preset creates

A `REMOTE_REPOSITORY` whose upstream is Docker Hub. Pull
`nginx:1.27` as
`us-central1-docker.pkg.dev/{project}/dockerhub-cache/nginx:1.27` — the
first pull populates the cache, every later pull is local. A cleanup
policy expires cached artifacts untouched for 30 days.

## Prerequisites

None. Uncomment `upstreamCredentials` (with a Secret Manager secret
holding a Docker Hub access token) to raise the upstream rate limits —
the Artifact Registry service agent needs
`roles/secretmanager.secretAccessor` on that secret.

## Composing hermetic builds

Point every Dockerfile `FROM`, GKE image reference, and CI base-image
pull at the cache path instead of `docker.io`. Combine with a virtual
repository to give consumers one endpoint spanning this cache and your
own standard repositories.

## Remix ideas

- Swap the upstream arm for `mavenPublicRepository: MAVEN_CENTRAL`,
  `npmPublicRepository: NPMJS`, or `pythonPublicRepository: PYPI` (with
  the matching `format`) for language-package caches.
- Use `commonRepository.uri: https://ghcr.io` to cache GitHub's registry
  or any other custom upstream.
