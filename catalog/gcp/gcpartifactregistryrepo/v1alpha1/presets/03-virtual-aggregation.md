# Virtual Aggregation Endpoint

One pull URL for everything: a virtual repository that serves first-party
images from the team's standard repository and falls through to the
Docker Hub cache for everything else — consumers never need to know where
an image actually lives.

## What this preset creates

A `VIRTUAL_REPOSITORY` aggregating two upstreams by priority. A pull of
`us-central1-docker.pkg.dev/{project}/docker-all/api-server:v1` resolves
from `app-images` (priority 100); a pull of `.../docker-all/nginx:1.27`
falls through to `dockerhub-cache` (priority 50). Higher priority wins
when both upstreams hold the same package.

## Prerequisites

- The `app-images` standard repository (preset 01) and the
  `dockerhub-cache` remote repository (preset 02), both Docker-format and
  in the same location. The references resolve each upstream's
  `repository_path` output — the exact value the upstream policy consumes.

## Composing consumption

Point every cluster, service, and CI pull at the virtual endpoint and
manage what it serves by editing upstream policies — mutable in place, so
adding a team repository or re-prioritizing never touches consumers.

## Remix ideas

- Add a per-team standard repository as a third upstream with its own
  priority band.
- Virtual repositories cannot be pushed to — CI keeps pushing to the
  standard repository directly; only pulls go through the aggregate.
