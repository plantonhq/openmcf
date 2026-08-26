---
title: "How to Configure Branch Deployments and Tag Releases"
date: "2026-04-02"
author:
  - name: "Planton Team"
    title: "Platform Engineering"
    bio: "Helping teams deploy infrastructure and services without the DevOps bottleneck"
tags:
  - "service-hub"
  - "branches"
  - "deployment"
  - "environments"
  - "pipeline"
category: "service-hub"
excerpt: "Control which Git events trigger pipelines, which environments receive deployments, and how to gate production releases with approval workflows."
---

# How to Configure Branch Deployments and Tag Releases

This tutorial shows you how to control the relationship between Git activity and deployments in Planton. You will learn how to configure which branches trigger pipelines, how to deploy to multiple environments from a single branch, how to gate production deployments with manual approval, and how to enable pipelines for pull requests and Git tags.

If you completed [How to Deploy Your First Service with Zero-Config CI/CD](/tutorials/how-to-deploy-your-first-service-with-zero-config-cicd), you have a Service that deploys to a single environment when you push to `main`. This tutorial builds on that foundation and covers the configuration options for more sophisticated deployment workflows.

> **Note**: The Planton web console provides a guided interface for configuring
> pipeline branches, deployment environments, and approval gates. This tutorial
> uses the CLI/YAML approach for stability and reproducibility. The console UI
> evolves frequently — always check it for the latest experience.

## What You Will Learn

- How `pipelineBranches` and `deploymentEnvironments` work together -- and why they are orthogonal
- How to deploy to multiple environments from a single branch with manual approval gates
- How to enable build-only and full-deploy pipelines for pull requests
- How to set up tag-based release workflows with glob pattern matching
- How to manually trigger pipelines for specific branches and commits

## Prerequisites

- [ ] A working Service deployed through Planton (see [How to Deploy Your First Service with Zero-Config CI/CD](/tutorials/how-to-deploy-your-first-service-with-zero-config-cicd))
- [ ] At least two environments configured in your Planton organization (e.g., `dev` and `production`)
- [ ] Kustomize overlays for each target environment in your repository (e.g., `_kustomize/overlays/dev/` and `_kustomize/overlays/production/`)
- [ ] The `planton` CLI installed and authenticated (`planton auth login`)

## How Pipeline Triggers and Deployment Targets Work

Two independent concepts control your CI/CD flow: **pipeline branches** determine which Git branches trigger builds, and **deployment environments** determine where artifacts get deployed. These are orthogonal -- every triggered pipeline deploys to all configured environments, not specific branches to specific environments. For more on pipeline configuration, see the [pipelines documentation](/docs/ci-cd/pipelines).

## Step 1: Deploy to Multiple Environments

If you followed [the first Service tutorial](/tutorials/how-to-deploy-your-first-service-with-zero-config-cicd), your Service deploys to a single environment. To add a second environment, you need two things: a kustomize overlay for that environment and an update to your Service manifest.

Create a production overlay in your repository at `_kustomize/overlays/production/`:

**`_kustomize/overlays/production/kustomization.yaml`**

```yaml
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization

resources:
  - ../../base

patches:
  - path: service.yaml
```

**`_kustomize/overlays/production/service.yaml`**

```yaml
apiVersion: kubernetes.planton.dev/v1alpha1
kind: KubernetesDeployment
metadata:
  name: my-service
  env: production
spec:
  availability:
    minReplicas: 3
  container:
    app:
      resources:
        requests:
          cpu: 250m
          memory: 512Mi
        limits:
          cpu: 2000m
          memory: 2Gi
```

Commit and push these files to your repository.

Then update your Service manifest to include both environments:

```yaml
deploymentEnvironments:
  - dev
  - production
```

Apply the change:

```bash
planton apply -f service.yaml
```

From this point forward, every push to `main` triggers a pipeline that deploys to `dev` first, then to `production`. The deployments happen sequentially -- `production` starts only after `dev` completes successfully. If the `dev` deployment fails, `production` is skipped.

## Step 2: Add a Manual Approval Gate

Deploying to production on every push without human review is appropriate for some teams but not for others. Manual approval gates let you pause the pipeline before a specific environment and wait for explicit approval.

For services using the git-based kustomize path (the default), manual gates are configured through your organization's **Promotion Policy** -- not in the Service manifest itself. The Promotion Policy defines the order in which environments are promoted and which transitions require approval.

A Promotion Policy consists of edges between environments. Each edge can have `manualApproval` set to `true`, which means the pipeline pauses before deploying to the destination environment:

```
dev → staging (automatic)
staging → production (manual approval required)
```

When a pipeline reaches an environment that requires manual approval, the deployment task enters an `awaiting_approval` status. The pipeline does not proceed until someone approves or rejects the deployment.

### Checking pipeline status

To see whether a pipeline is waiting for approval:

```bash
planton service last-pipeline --service my-service
```

Look for a deployment task with `requires_manual_gate: true` and status `awaiting_approval` in the output.

You can also stream the pipeline status in real time:

```bash
planton pipeline stream-status <pipeline-id>
```

### Approving or rejecting a deployment

To approve the deployment and let the pipeline continue:

```bash
planton pipeline resolve-manual-gate <pipeline-id> <deployment-task-name> yes
```

To reject the deployment and stop the pipeline:

```bash
planton pipeline resolve-manual-gate <pipeline-id> <deployment-task-name> no
```

The `deployment-task-name` is the name of the specific deployment task that is awaiting approval. You can find it in the `planton service last-pipeline` output.

The Promotion Policy is configured at the organization level through the Planton web console. The details of Promotion Policy configuration are outside the scope of this tutorial, but the effect on your pipelines is straightforward: when an environment transition requires approval, the pipeline pauses and waits.

## Step 3: Enable Pull Request Builds

By default, pull request commits do not trigger any pipeline activity. You can enable build-only pipelines for pull requests to validate that the code builds successfully before merging.

Add this to your Service's pipeline configuration:

```yaml
pipelineConfiguration:
  enablePullRequestBuilds: true
```

Apply the change:

```bash
planton apply -f service.yaml
```

With this setting, when a pull request is opened (or updated with new commits) against a branch in `pipelineBranches`, a pipeline runs that:

- Clones the repository at the PR commit
- Builds the container image (Buildpacks or Dockerfile)
- Pushes the image to the registry tagged with the commit SHA
- **Skips the deploy stage entirely** -- no environments are affected

This is useful for CI validation: you know the code builds before anyone reviews or merges the pull request. GitHub check runs will show the build status directly on the PR.

The PR's **target branch** (the branch the PR would merge into) must be listed in `pipelineBranches` for the pipeline to trigger. For example, if `pipelineBranches: [main]`, then PRs targeting `main` trigger builds, but PRs targeting `develop` do not.

## Step 4: Enable Pull Request Preview Deployments

If build validation is not enough and you want full preview environments for each pull request, enable preview deployments:

```yaml
pipelineConfiguration:
  enablePullRequestDeployments: true
```

With this setting, pull request pipelines run the full pipeline including the deploy stage. Each PR gets its own deployment, allowing reviewers to test the changes in a live environment before merging.

There are important considerations before enabling this:

- **Resource cost**: Every active PR consumes cloud resources for its deployment. If your team has many concurrent PRs, this can become expensive.
- **Environment management**: Preview deployments create Cloud Resources that persist until the pipeline is rerun or the resources are manually cleaned up. Automatic cleanup when a PR is closed is not currently implemented -- you will need to manage preview environment lifecycle yourself.
- **Both flags**: If both `enablePullRequestBuilds` and `enablePullRequestDeployments` are true, the deployment behavior (`enablePullRequestDeployments`) takes precedence. You do not need to set both -- `enablePullRequestDeployments: true` implies that the build also runs.

## Step 5: Configure Tag-Based Releases

Tag pipelines are independent of branch and pull request pipelines. They are triggered when a Git tag is pushed to the repository, and they use the **Git tag name** as the container image tag instead of the commit SHA. This makes tag-based releases useful for versioned release workflows.

### Build validation for tags

If you want to verify that a tagged release builds successfully without deploying it:

```yaml
pipelineConfiguration:
  enableTagBuilds: true
  tagPatterns:
    - "v*"
```

When a tag matching the pattern is pushed (e.g., `v1.0.0`, `v2.3.1-beta`), a pipeline runs that builds the image and tags it with the Git tag name:

- Branch pipeline image: `myapp:a1b2c3d4e5f6` (commit SHA)
- Tag pipeline image: `myapp:v1.0.0` (Git tag name)

The deploy stage is skipped. This is useful for creating release candidate images that you can deploy manually or through a separate process.

### Full release deployment from tags

If you want tags to trigger both build and deployment:

```yaml
pipelineConfiguration:
  enableTagDeployments: true
  tagPatterns:
    - "v*"
    - "release-*"
```

This triggers the full pipeline including the deploy stage. The same `deploymentEnvironments` filter applies -- the tag pipeline deploys to all listed environments, with manual approval gates if configured.

### Tag pattern matching

`tagPatterns` uses glob-style matching:

- `v*` matches `v1.0.0`, `v2.1.3-beta`, `v0.0.1-rc.1`
- `release-*` matches `release-2026-04-02`, `release-sprint-47`
- Empty `tagPatterns` means **all tags** trigger pipelines when tag builds or deployments are enabled

You can combine multiple patterns:

```yaml
tagPatterns:
  - "v*"
  - "release-*"
```

A tag that matches **any** pattern triggers the pipeline.

## Step 6: Manually Trigger a Pipeline

You can trigger a pipeline from the CLI without pushing a commit. This is useful for re-running a build on a specific branch or testing pipeline behavior.

```bash
planton service run-pipeline --service my-service --branch main
```

This triggers a pipeline using the latest commit on `main`. To use a specific commit:

```bash
planton service run-pipeline --service my-service --branch main --commit a1b2c3d4
```

To rerun a previously executed pipeline:

```bash
planton service rerun-pipeline <pipeline-id>
```

Manual triggers bypass the webhook flow but go through the same pipeline execution path -- the same build, kustomize, and deploy stages run.

## Common Patterns and Tips

### Build-only services

If your service should be built and packaged but not deployed -- for example, a shared library image or a service whose deployment is managed separately -- disable the deploy stage:

```yaml
pipelineConfiguration:
  disableDeployments: true
```

The pipeline still runs: it clones the repository, builds the image, and pushes it to the registry. But the deploy stage is skipped entirely. This is the pattern used by several of Planton's own internal services.

### Pausing all automation

To temporarily stop all pipeline activity for a service without deleting the webhook or changing branch configuration:

```yaml
pipelineConfiguration:
  disablePipelines: true
```

When `disablePipelines` is true, Git commits do not trigger pipelines and manual pipeline triggers are also disabled. Set it back to `false` (or remove the field) to resume automation.

### Watching pipelines by environment

To see only pipelines that deployed to a specific environment:

```bash
planton service pipelines --service my-service --filter-envs production
```

This is useful when you have multiple environments and want to check the deployment history for one of them.

### How trigger paths interact with branch filtering

If your Service uses `triggerPaths` (for monorepo setups), both conditions must be met for a pipeline to run:

1. The push must be to a branch in `pipelineBranches`
2. At least one changed file must match a pattern in `triggerPaths` (or be under `projectRoot`)

A push to `main` that only changes files outside the service's `projectRoot` and `triggerPaths` will not trigger a pipeline, even though `main` is in `pipelineBranches`.

### Canceling a running pipeline

To cancel a pipeline that is currently running:

```bash
planton pipeline cancel <pipeline-id>
```

If the build stage is running, the build is terminated immediately. If the deploy stage is running, the current deployment task completes its in-flight operation and remaining tasks are skipped.

## What to Do Next

- **Switch to self-managed pipelines** for custom build steps like SonarQube analysis, security scanning, or multi-stage builds. See [How to Switch to Self-Managed Tekton Pipelines](/tutorials/how-to-switch-to-self-managed-tekton-pipelines) (coming soon).
- **Manage secrets and variables** to inject environment-specific configuration into your deployments. See [How to Store and Reference Credentials in a Provider Connection](/tutorials/how-to-store-and-reference-credentials-in-a-provider-connection) (coming soon).
