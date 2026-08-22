---
title: "Monorepo Support"
description: "Configure multiple services from a single Git repository with per-service project roots, trigger paths, and sparse checkout."
icon: settings
order: 65
tags:
  - Monorepo
  - Git
  - Configuration
  - Service Hub
---

# Monorepo Support

Modern codebases increasingly use monorepos for code sharing and atomic changes. But traditional CI/CD struggles with them: ten services in one repository means a shared library change triggers builds for all ten, even when only two need rebuilding. Clone times balloon because every build fetches the entire repository. Configuring path-based triggers requires deep knowledge of each CI system's quirks.

Service Hub provides built-in monorepo support. Each service is scoped to its own directory within the repository, triggers pipelines only when relevant files change, and clones only the directories it needs.

## Core Settings

Three settings in your service's Git repository configuration control monorepo behavior.

### Project Root

The directory within the repository where the service's code is located. This becomes the working directory for build operations.

```yaml
git_repo:
  project_root: "services/payment-api"
```

If empty, the service uses the repository root. The pipeline's build tasks (Buildpacks, Dockerfile, or self-managed) execute from this directory. The Dockerfile path in the pipeline configuration is relative to the project root.

### Trigger Paths

Glob patterns that trigger pipelines when matched files change, **in addition to** the default rule that any change inside the project root triggers a build.

```yaml
git_repo:
  project_root: "services/payment-api"
  trigger_paths:
    - "packages/shared-lib/**"
    - "proto/**"
```

How trigger paths work:

- Paths are matched against repository-root-relative file paths.
- If no trigger paths are set, pipelines run only when files inside the project root change.
- If trigger paths are set, pipelines run when files match **either** the project root **or** any trigger path.
- To trigger on every commit regardless of changed files, use `["**/*"]`.

Common patterns:

| Pattern | Triggers When |
|---------|--------------|
| `packages/shared-lib/**` | Shared library code changes |
| `proto/**` | Protobuf definitions change |
| `docker/base-image/**` | Base image definitions change |
| `**/*` | Any file in the repository changes |

### Sparse Checkout

Directories to clone during the pipeline's git-clone step. When set, only these directories are fetched instead of the entire repository.

```yaml
git_repo:
  project_root: "services/payment-api"
  sparse_checkout_directories:
    - "services/payment-api"
    - "packages/shared-lib"
    - "proto"
```

Benefits for large monorepos:

- Faster clone times (fewer files to download)
- Less disk space used in pipeline workers
- Only code the service needs is available during the build

If empty, the entire repository is cloned.

## Additional Pipeline Settings for Monorepos

Two additional settings provide flexibility for more complex monorepo layouts:

### Kustomize Base Directory

The path to the kustomize directory relative to the project root. Defaults to `_kustomize`.

This is useful when the project root must be at the repository root (e.g., for Docker build context or yarn workspace resolution) but kustomize overlays are in a service-specific subdirectory:

```yaml
git_repo:
  project_root: ""  # repo root for Docker context
pipeline_configuration:
  kustomize_base_directory: "services/payment-api/_kustomize"
```

### Package Manager Directory

The directory from which to run the package manager install step. Relevant for Cloudflare Worker builds in monorepos using yarn workspaces with `workspace:*` protocol dependencies.

```yaml
git_repo:
  project_root: "workers/edge-router"  # worker code location
pipeline_configuration:
  package_manager_directory: "."  # run yarn install at repo root for workspace resolution
```

If empty, falls back to the project root for dependency installation.

## Practical Patterns

### Standard Monorepo

Multiple services, each in its own directory. Shared code triggers rebuilds for dependent services.

```yaml
# Service 1: User API
git_repo:
  owner_name: acmecorp
  name: backend-monorepo
  project_root: "services/user-api"
  trigger_paths:
    - "packages/shared-lib/**"
  sparse_checkout_directories:
    - "services/user-api"
    - "packages/shared-lib"
pipeline_configuration:
  image_build_method: buildpacks
  image_repository_path: acmecorp/user-api
```

```yaml
# Service 2: Payment API (same repository, different service)
git_repo:
  owner_name: acmecorp
  name: backend-monorepo
  project_root: "services/payment-api"
  trigger_paths:
    - "packages/shared-lib/**"
    - "proto/**"
  sparse_checkout_directories:
    - "services/payment-api"
    - "packages/shared-lib"
    - "proto"
pipeline_configuration:
  image_build_method: dockerfile
  image_repository_path: acmecorp/payment-api
```

A change to `services/user-api/` triggers only the User API pipeline. A change to `packages/shared-lib/` triggers both pipelines.

### Yarn Workspaces

A JavaScript/TypeScript monorepo where services share dependencies through yarn workspaces:

```yaml
git_repo:
  project_root: "apps/web-worker"
pipeline_configuration:
  package_manager_directory: "."  # run yarn install at repo root
  kustomize_base_directory: "apps/web-worker/_kustomize"
```

This ensures `yarn install` runs at the repository root (where `package.json` with workspace definitions lives), while the build targets the worker's subdirectory.

## Configuring in the Web Console

Monorepo settings are available in two places:

**During service creation**, the wizard provides these fields under **Advanced Git Settings**:
- Project root
- Trigger paths (add/remove controls)
- Sparse checkout directories (add/remove controls)

**After creation**, the service's **Settings** tab under **Advanced Settings** provides the same fields as editable controls.

<!-- SCREENSHOT: Monorepo advanced settings
  Page: /orgs/{org}/service/{serviceId} (Settings tab, Advanced Settings section)
  Action: Show the advanced settings with project root, trigger paths, and sparse checkout configured
  Focus: The three monorepo-related fields with example values
  Alt: Service advanced settings showing project root, trigger paths array, and sparse checkout directories array with configured values
-->

## Using the CLI

```bash
# Deploy from a specific project root
planton service deploy --project services/payment-api

# Build kustomize output for inspection
planton service kustomize build
```

## Related Documentation

- [What is a Service?](/docs/ci-cd/what-is-a-service) — Service configuration overview
- [Build Methods](/docs/ci-cd/build-methods) — How builds work within the project root
- [Pipelines](/docs/ci-cd/pipelines) — How trigger paths and branches control pipeline execution
