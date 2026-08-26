---
title: "Deployment Environments"
description: "Control which environments your services deploy to — selective filtering, branch-based targeting, and overlay interaction."
icon: deployment
order: 50
tags:
  - Deployment
  - Environments
  - Service Hub
---

# Deployment Environments

When a service has Kustomize overlays for dev, staging, QA, and production, the pipeline deploys to all of them by default. That is the right behavior for most services. But sometimes you need to be selective — a canary service should only reach staging, a legacy migration service targets specific environments, or branch-based workflows mean each branch should deploy to one environment and skip the rest.

The deployment environments feature gives you that control without modifying your repository structure. You configure which environments a service deploys to, and the pipeline skips the rest — even though the overlays exist in Git.

<!-- VIDEO: Deployment environments walkthrough
  Source: Cloudflare Stream
  Alt: Video walkthrough showing deployment environments configuration and filtering
-->

## How It Works

The deployment environments setting is a filter on overlay directories. It does not create environments or modify your repository — it controls which existing overlays produce deployment tasks during the pipeline's deploy stage.

### Default Behavior

When no deployment environments are configured (the default), the pipeline deploys to **every** overlay directory found in `_kustomize/overlays/`, except the `local` overlay which is always skipped.

### Selective Deployment

To deploy to specific environments only, list them explicitly in the service configuration:

```yaml
deployment_environments:
  - dev
  - stage
```

With this configuration, even if the repository contains dev, stage, uat, and production overlays, only dev and stage produce deployment tasks. The uat and production overlays are ignored.

<!-- SCREENSHOT: Environment list
  Page: /orgs/{org}/environments
  Action: Show the environments list with at least dev, staging, and production environments
  Focus: The environment list with names and resource counts
  Alt: Environment management page showing development, staging, and production environments
-->

## Configuring Deployment Environments

### Web Console

The deployment environments selector is in the **Settings** tab on the service detail page. By default, the setting shows "Service is configured to deploy to all environments found under the `_kustomize` directory in the repository."

Click **Edit** to open the configuration modal — a checkbox list of all environments in your organization. Select only the environments you want this service to deploy to. Leaving all unchecked restores the default behavior (deploy to all).

![Service showing deployment environments configuration](https://assets.planton.ai/site/images/docs/ci-cd/deployment-environments/service-showing-deployment-environments.png)

### CLI

Update the deployment environments in your service YAML and apply:

```yaml
deployment_environments:
  - dev
  - stage
```

```bash
planton apply -f service.yaml
```

## Branch-Based Targeting

For more granular control, you can configure branch-based deployment directly in your Kustomize overlay manifests using the `planton.ai/git-branch` label. When present on an overlay, this label determines whether that overlay deploys based on the Git branch that triggered the pipeline.

### Adding the Label

Add `planton.ai/git-branch` to the `metadata.labels` section of an overlay's service manifest:

```yaml
# _kustomize/overlays/dev/service.yaml
apiVersion: kubernetes.planton.dev/v1alpha1
kind: KubernetesDeployment
metadata:
  name: my-service
  env: dev
  labels:
    planton.ai/git-branch: dev-*
spec:
  # ... configuration
```

### Pattern Matching

The label value supports glob patterns:

| Pattern | Matches |
|---------|---------|
| `main` | Exact match — only the `main` branch |
| `dev-*` | Wildcard suffix — `dev-feature`, `dev-hotfix`, `dev-123` |
| `feature/*` | Wildcard prefix — `feature/login`, `feature/checkout` |
| `release-*-hotfix` | Multiple wildcards — `release-1.0-hotfix`, `release-2.1-hotfix` |
| `"*"` | Any branch |

Matching is **case-sensitive**: `dev` does not match `Dev`.

### Precedence Rules

When determining whether an overlay produces a deployment task:

1. **Branch label first** — If `planton.ai/git-branch` exists on the overlay manifest, the branch must match the pattern. If it does not match, the overlay is skipped regardless of other settings.
2. **Deployment environments second** — If the label is absent, the overlay is subject to the service-level deployment environments filter.
3. **Default** — If neither is configured, the overlay deploys.

This means you can mix approaches: some overlays use branch labels for repository-level control, while others rely on the service-level deployment environments setting.

### Example: Branch-to-Environment Mapping

```yaml
# _kustomize/overlays/dev/service.yaml
metadata:
  env: dev
  labels:
    planton.ai/git-branch: dev-*

---
# _kustomize/overlays/staging/service.yaml
metadata:
  env: staging
  labels:
    planton.ai/git-branch: main

---
# _kustomize/overlays/prod/service.yaml
metadata:
  env: prod
  labels:
    planton.ai/git-branch: release-*
```

With this setup:
- Push to `dev-feature` → deploys only to dev
- Push to `main` → deploys only to staging
- Push to `release-1.0` → deploys only to production

## Use Cases

### Branch-Based Deployments Without Merge Conflicts

When teams use separate branches for different environments, the `_kustomize` directories frequently cause merge conflicts. The deployment environments feature solves this:

1. Keep all overlay directories in all branches.
2. Create separate services for each branch-environment combination.
3. Configure each service to deploy only to its target environment.

When branches merge, the overlay directories do not conflict because every branch carries the full set — the deployment environment setting controls which ones actually execute.

### Progressive Rollout

Start a new service with limited environments and expand as confidence grows:

```yaml
# Week 1: Dev only
deployment_environments:
  - dev

# Week 2: Add staging
deployment_environments:
  - dev
  - stage

# Week 3: Full rollout — clear the list to deploy everywhere
deployment_environments: []
```

### Environment-Specific Services

Some services should never deploy everywhere:

```yaml
# Debug tooling — dev only
deployment_environments:
  - dev

# Internal monitoring — non-production
deployment_environments:
  - dev
  - stage
```

## How It Appears in Pipelines

When a pipeline runs, the deploy stage only creates tasks for environments that pass the filter. The pipeline detail page shows which deployment tasks were created — skipped environments do not appear.

![Pipeline showing deployment tasks](https://assets.planton.ai/site/images/docs/ci-cd/deployment-environments/pipeline-detail-showing-the-deployment-task.png)

If your service is configured to deploy only to dev and stage, only those two environments appear as deployment tasks, even if the repository contains additional overlays.

![Deployment environment selector](https://assets.planton.ai/site/images/docs/ci-cd/deployment-environments/deployment-environment-selector-model-with-dev-and-stage-selected.png)

## Troubleshooting

### Overlay Not Deploying

- **Check the `planton.ai/git-branch` label** — If the label exists, its pattern must match the triggering branch. Matching is case-sensitive and uses glob syntax (`*` for any characters, `?` for a single character).
- **Check the deployment environments list** — The overlay directory name must appear in the list. Names must match exactly.
- **Check the `local` overlay** — The `local` overlay is always skipped. It exists for `.env` file generation only.

### All Environments Deploying Despite Configuration

- Verify that the deployment environments list is not empty — an empty list means "deploy to all."
- Confirm the configuration was saved successfully in the service settings.

## Related Documentation

- [Deployment Stage](/docs/ci-cd/deployment-stage) — How manifests are resolved and deployed, including the Kustomize model
- [Deployment Targets](/docs/ci-cd/deployment-targets) — Supported platforms and Git-based vs inline configuration
- [What is a Service?](/docs/ci-cd/what-is-a-service) — Service configuration overview
- [Pipelines](/docs/ci-cd/pipelines) — The pipeline execution model
