# AWS CodeBuild Project — Architecture and Design

## Overview

AWS CodeBuild is a fully managed build service: it provisions a fresh,
isolated build environment for every build, runs the phases declared in a
buildspec, streams logs, and tears the environment down. There are no build
servers to patch or scale, and billing is per build minute.

A **project** is the build configuration unit — where the source comes from,
what container the build runs in, where output artifacts land, and what
operational guardrails (timeouts, concurrency, retries) apply. The project
itself is a metadata-only control-plane object: creating one provisions
nothing until a build starts.

## Architecture

### Build Execution Model

```
Trigger (webhook / CodePipeline / StartBuild)
     │
     ▼
┌───────────────────┐
│  CodeBuild queue   │ ← queued_timeout caps waiting
└─────┬─────────────┘
      │  fresh container per build
      ▼
┌────────────────────────────┐
│  Build environment          │ ← environment.type + compute_type + image
│  (on-demand or fleet)       │ ← fleet_arn joins reserved capacity
│                             │
│  install → pre_build →      │ ← buildspec phases
│  build   → post_build       │ ← build_timeout caps the whole run
└─────┬───────────┬──────────┘
      │           │
      │           └── logs ──→ CloudWatch Logs and/or S3
      │
      └── artifacts ──→ S3 (primary + up to 12 secondary) or CodePipeline
```

### Source Model

A project has exactly one **primary source** and up to 12 **secondary
sources**. Each secondary source carries a `source_identifier`; CodeBuild
checks it out at `$CODEBUILD_SRC_DIR_<identifier>` so a build can compose an
application repository with, say, a shared build-tooling repository. Version
pins per secondary source live in `secondary_source_versions`.

Source authorization follows one of three paths:

1. **CodeConnections** (`auth.type: CODECONNECTIONS`) — the modern path: a
   connection object grants repository access with no stored tokens.
2. **Secrets Manager** (`auth.type: SECRETS_MANAGER`) — a secret holds the
   provider token.
3. **Account-level OAuth** (`auth.type: OAUTH` or no auth block) — the legacy
   account-wide grant.

### Compute Model

Three compute families, selected by `environment.type` + `compute_type`:

| Family | Types | Traits |
|--------|-------|--------|
| On-demand containers | LINUX/ARM/Windows/GPU containers | Default; per-build isolation, cold start seconds-to-minutes |
| Lambda compute | `*_LAMBDA_CONTAINER` + `BUILD_LAMBDA_*` | Fastest start; no privileged mode, no timeouts (AWS caps runs), no caching |
| Reserved fleets | `*_EC2`, `MAC_ARM`, `ATTRIBUTE_BASED_COMPUTE`, `CUSTOM_INSTANCE_TYPE` | Always-warm machines a project joins via `fleet_arn`; required for macOS |

For Docker-image builds, two accelerators exist: `privileged_mode` (a
per-build Docker daemon) and `docker_server` (a **persistent** dedicated
Docker server whose layer state survives across builds — dramatically faster
repeat image builds).

### Webhook Model

The folded webhook registers the project with its source provider. Filter
groups are OR'd; filters within a group are AND'd. Beyond per-repository
webhooks, `scope_configuration` widens coverage to a GitHub organization,
the whole connection (GITHUB_GLOBAL), or a GitLab group — the shape runner
projects (`build_type: RUNNER_BUILDKITE_BUILD`) and org-wide CI use.

`manual_creation` inverts the registration: AWS mints the payload URL and
HMAC secret (exported as `webhook_payload_url` / `webhook_secret` outputs)
and the operator wires the repository webhook by hand — required for GitHub
Enterprise, useful wherever the connection lacks repository admin rights.

`pull_request_build_policy` gates PR-triggered builds behind an approval
comment — the protection that keeps untrusted fork code away from CI
credentials in open-source repositories.

### Cross-Account Access

The folded `resource_policy` attaches a resource-based IAM policy to the
project, the mechanism by which a central CI account starts builds in this
account. One document per project, keyed by the project ARN.

## Design Decisions

1. **Webhook folds; credentials, report groups, and fleets do not.** The
   webhook is 1:1 with its project and meaningless in isolation. Source
   credentials are an account/region-wide store superseded by the per-source
   `auth` block; report groups and fleets are shared resources referenced by
   many projects — the project composes with a fleet through the
   `environment.fleet_arn` reference.
2. **The project name derives from `metadata.name`.** CodeBuild project
   names are create-time-immutable; both IaC engines key the resource off
   the same basis so state and cloud identity never diverge by engine.
3. **Lambda-compute honesty.** AWS ignores build/queued timeouts and forbids
   privileged mode on Lambda compute. The spec rejects those combinations at
   validation time instead of letting AWS silently drop them.
4. **`host_kernel` is excluded** until it ships in released provider
   versions of both IaC engines.

## Operational Notes

- **IAM eventual consistency**: creating the service role and project in the
  same deploy can hit "not authorized to perform" transients; the providers
  retry for up to 2 minutes, which is the only wait this resource has.
- **Public projects** (`project_visibility: PUBLIC_READ`) expose build
  results, logs, and artifacts world-readably through the
  `public_project_alias`; the `resource_access_role` is what CodeBuild uses
  to read the logs/artifacts it re-exposes. Never enable it for projects
  whose logs may leak secrets.
- **Badges** (`badge_enabled`) only exist for repository sources — not
  CODEPIPELINE, S3, or NO_SOURCE.
- **Secondary artifacts** need matching `secondary-artifacts` sections in
  the buildspec; the `artifact_identifier` is the join key.
