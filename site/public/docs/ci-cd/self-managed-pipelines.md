---
title: "Self-Managed Pipelines"
description: "Bring your own Tekton pipeline while retaining platform deployment orchestration, status tracking, and manual gates."
icon: pipeline
order: 35
tags:
  - Pipelines
  - Tekton
  - CI/CD
  - Service Hub
---

# Self-Managed Pipelines

Self-managed pipelines let you define your own Tekton pipeline YAML in your repository while Planton handles execution orchestration, deployment sequencing, status streaming, and manual approval gates.

## Platform-Managed vs Self-Managed

The pipeline provider setting on a Service's pipeline configuration determines who controls the build pipeline:

| Provider | Build Pipeline | Deploy Stage | Status Tracking |
|----------|---------------|--------------|-----------------|
| Platform-managed | Planton creates and maintains the pipeline. You choose Buildpacks or Dockerfile. | Planton-managed | Full |
| Self-managed | You write a Tekton pipeline YAML in your repository. | Planton-managed | Full |

In both cases, the deploy stage is managed by Planton. The difference is only in who controls the build pipeline definition.

**Choose self-managed when**: You need custom build steps (test suites, linting, security scanning, multi-stage builds with specific tooling), parallel task execution, or conditional task logic that goes beyond what Buildpacks or Dockerfile builds provide.

## Pipeline File Location

The pipeline YAML file path is configured on the Service's pipeline configuration. If left empty, it defaults to `.planton/pipeline.yaml` in the repository root.

```yaml
pipeline_configuration:
  pipeline_provider: self
  tekton_pipeline_yaml_file: ".planton/pipeline.yaml"  # default
```

The file path is relative to the repository root. This allows shared pipeline definitions across services in a monorepo.

Other common locations:

```yaml
tekton_pipeline_yaml_file: ".planton/my-service-pipeline.yaml"
tekton_pipeline_yaml_file: "ci/pipelines/build.yaml"
tekton_pipeline_yaml_file: ".github/tekton/pipeline.yaml"
```

<!-- SCREENSHOT: Self-managed pipeline configuration
  Page: /{org}/service/{slug} (pipeline configuration section)
  Action: Show the service configuration with pipeline_provider set to "self" and tekton_pipeline_yaml_file field visible
  Focus: The pipeline provider toggle and YAML file path configuration
  Alt: Service pipeline configuration showing self-managed provider selected with Tekton pipeline YAML file path
-->

## Platform-Injected Parameters

When Planton executes a self-managed pipeline, it injects a set of standard parameters that your pipeline can reference. These parameters provide the context your build tasks need to clone the repository, build and tag the image, and produce deployment manifests.

Your Tekton pipeline YAML must declare matching `params` entries for any platform-injected parameters you want to use.

**Git parameters**:
- `git-url` — Clone URL for the repository
- `git-revision` — The commit SHA that triggered the pipeline
- `sparse-checkout-directories` — Directories for sparse checkout (from Service's git_repo configuration)

**Image build parameters**:
- `image-name` — Full image reference including registry host and repository path
- `image-repository-path` — Just the repository path portion (from pipeline configuration)
- `dockerfile-path` — Path to Dockerfile relative to project root

**Kustomize parameters**:
- `kustomize-manifests-config-map-name` — Name of the ConfigMap where kustomize-build output is stored for the deploy stage
- `kustomize-base-directory` — Path to the kustomize base directory

**Tracking parameters**:
- `owner-identifier-label-key` — Label key for resource ownership tracking
- `owner-identifier-label-value` — Label value for resource ownership tracking

## Custom Parameters

Additional parameters can be passed from the Service configuration to your pipeline using the `params` map:

```yaml
pipeline_configuration:
  pipeline_provider: self
  params:
    ENABLE_TESTS: "true"
    DEPLOY_REGION: "us-west-2"
    BUILD_TIMEOUT: "30m"
```

These key-value pairs are injected alongside platform parameters. Your Tekton pipeline YAML must declare matching `params` entries.

## Editing Pipelines in the Web Console

The web console provides a built-in pipeline editor for self-managed pipelines, accessible from the **Pipelines** tab of a service's detail page when `pipeline_provider` is `self`.

The editor provides:

- **Branch selector**: Choose which branch to view and edit pipeline files on.
- **File selector**: Dropdown listing pipeline files found in the repository (within `.planton/` or the configured path).
- **YAML editor**: Syntax-highlighted editor for the pipeline definition.
- **DAG preview**: Visual representation of the pipeline's task graph, rendered from the YAML.
- **Commit flow**: Save changes with a commit message and view the diff before committing.

Changes committed through the editor are pushed directly to the selected branch in the repository.

<!-- SCREENSHOT: Self-managed pipeline editor
  Page: /orgs/{org}/service/{serviceId} (Pipelines tab, self-managed mode)
  Action: Show the pipeline editor with YAML content and DAG preview
  Focus: The YAML editor panel and the DAG visualization side-by-side
  Alt: Self-managed pipeline editor showing Tekton YAML on the left and a visual DAG of tasks on the right
-->

## Pipeline File Management via API

The platform provides APIs for managing pipeline files in the service repository:

- **Update pipeline file** — Commits an updated pipeline file to the repository. Supports optimistic locking via an expected base commit SHA to prevent concurrent edit conflicts.
- **List pipeline files** — Lists pipeline files found in the repository, useful for populating the file selector in the editor.

## Tekton Pipeline and Task Registry

Planton maintains a registry of reusable Tekton pipeline and task definitions. These can be registered at the platform or organization level.

Each registered resource includes:

| Field | Purpose |
|-------|---------|
| `description` | Human-readable description of what the pipeline or task does |
| `git_repo` | Git repository where the source YAML is stored (web URL, clone URL, file path) |
| `yaml_content` | Complete YAML content, stored for immediate use without requiring git access |
| `overview_markdown` | Markdown documentation for the pipeline or task |
| `selector` | Ownership — platform-level (shared across all organizations) or organization-level |

### Registering via CLI

```bash
# Register a Tekton pipeline
planton tekton pipeline register \
  --yaml-file pipelines/docker-build.yaml \
  --name "Docker Build Pipeline" \
  --description "Builds container images using Dockerfile with layer caching" \
  --git-web-url "https://github.com/plantonhq/tekton-hub/blob/main/pipelines/docker-build.yaml" \
  --git-clone-url "https://github.com/plantonhq/tekton-hub.git" \
  --git-file-path "pipelines/docker-build.yaml" \
  --tags "docker,build" \
  --overview-markdown-file docs/docker-build.md

# Register a Tekton task
planton tekton task register \
  --yaml-file tasks/run-tests.yaml \
  --name "Run Tests" \
  --description "Executes test suite with configurable test runner" \
  --git-web-url "https://github.com/plantonhq/tekton-hub/blob/main/tasks/run-tests.yaml" \
  --git-clone-url "https://github.com/plantonhq/tekton-hub.git" \
  --git-file-path "tasks/run-tests.yaml"
```

Add `--platform` to register as a platform-level resource (shared across all organizations).

## Monitoring Self-Managed Pipelines

Self-managed pipelines receive the same status tracking and log streaming as platform-managed pipelines:

```bash
# Stream pipeline status updates
planton pipeline stream-status <pipeline-id>

# Stream pipeline logs
planton pipeline stream-logs <pipeline-id>

# Cancel a running pipeline
planton pipeline cancel <pipeline-id>
```

The web console displays real-time pipeline progress, build logs, and deployment status for both pipeline types.

## Related Documentation

- [What is a Service?](/docs/ci-cd/what-is-a-service) — Service configuration overview
- [Pipelines](/docs/ci-cd/pipelines) — The pipeline execution model
- [Build Methods](/docs/ci-cd/build-methods) — Platform-managed build methods (Buildpacks, Dockerfile)
- [Deployment Targets](/docs/ci-cd/deployment-targets) — Where built artifacts are deployed
