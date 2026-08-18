# AWS ECR Registry Settings

The container registry's control plane for one region: what scans your images and when, where they replicate, which upstream registries pull through transparently, what auto-created repositories look like, and which CI roles stay out of pull-time metrics.

## What Gets Managed

- The registry permissions policy (cross-account replication and sharing grants).
- Scanning: the basic scanner or Amazon Inspector, with per-repository-pattern frequency rules.
- Replication rules to other regions or accounts.
- Pull-through cache rules (Docker Hub, ghcr.io, registry.k8s.io, another account's ECR) with their credentials.
- Repository creation templates stamped onto auto-created repositories.
- Account-level toggles: scanner version, blob mounting, registry policy scope.
- Pull-time update exclusions for automation principals.

## Before You Deploy

### Planton Setup

- An organization and environment.
- An AWS provider connection with ECR permissions.

### AWS Prerequisites

- For credentialed cache upstreams (Docker Hub): a Secrets Manager secret named under `ecr-pullthroughcache/`.
- For cross-account replication: the destination registry's policy must allow ecr:ReplicateImage from this account.

## After You Deploy

- Pulls of `{registry_url}/{prefix}/...` transparently fetch from the cached upstream and store here.
- Repositories created by replication, cache pulls, or first pushes carry the matching template's settings.

## Common Changes

- Everything except cache-rule and template prefixes updates in place (prefixes are fixed-for-life keys).
- One instance per region — a second instance targeting the same region fights over the same registry objects.
