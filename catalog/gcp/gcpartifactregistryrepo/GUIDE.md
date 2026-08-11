# GcpArtifactRegistryRepo Guide

The judgment this guide protects: a registry is the supply chain's
load-bearing middle — every deploy pulls from it, every build pushes to
it, and its failure mode is "nothing ships". Configure it for the pull
path first.

## Pick the mode for the traffic, not the taxonomy

STANDARD_REPOSITORY is where your own artifacts live. REMOTE_REPOSITORY
is a caching proxy for someone else's registry — it exists to absorb
Docker Hub rate limits and internet flakiness at build time.
VIRTUAL_REPOSITORY federates both behind one URL so CI configures a
single endpoint. The productive pattern is one standard repo per team or
product, one remote repo per external registry you depend on, and a
virtual repo only when consumers genuinely should not know the split.

## Remote upstreams: use commonRepository for anything custom

The well-known public arms (`dockerPublicRepository: DOCKER_HUB`,
`MAVEN_CENTRAL`, `NPMJS`, `PYPI`, apt/yum mirror bases) cover the big
registries. Everything else — GitLab registries, JFrog, another AR
repository, a private mirror — goes through `commonRepository` with a
URI. The provider's per-format `custom_repository` arms are deprecated
in its own schema; this spec never modeled them, and their capability
lives in `commonRepository`. Credentials rotate in place (username +
Secret Manager version); everything else about a remote config is
immutable, so choose the upstream deliberately.

## Cleanup policies: dry-run first, always

A DELETE cleanup policy is an automated artifact shredder; a condition
written against the wrong tag state removes images production still
pulls. Ship new policies with `cleanupPolicyDryRun: true`, watch the
audit log for what WOULD be deleted, then flip it live. KEEP policies
(`mostRecentVersions`) are the guard rail — pair every age-based DELETE
with a most-recent KEEP so a stale-but-current artifact never matches.

## Immutable tags are a supply-chain control

`dockerConfig.immutableTags: true` makes a tag permanently name one
digest — re-pushing `v1.2.3` fails instead of silently changing what
production deploys. Turn it on for release repositories; leave it off
only where mutable channel tags (`latest`, `stable`) are the contract.

## Destroy semantics

`deletionPolicy: DELETE` (the default) removes the repository AND every
artifact in it — new pulls fail immediately; only already-pulled images
on nodes keep running. `PREVENT` is the right posture for any registry
that running workloads reference. `ABANDON` keeps the registry serving
while dropping management — the IAM grants this kind created are
removed with the destroy either way, so re-adopt quickly or pulls that
depended on those grants start failing.
