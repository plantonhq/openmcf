---
title: "How to Switch to Self-Managed Tekton Pipelines"
date: "2026-04-02"
author:
  - name: "Planton Team"
    title: "Platform Engineering"
    bio: "Helping teams deploy infrastructure and services without the DevOps bottleneck"
tags:
  - "service-hub"
  - "tekton"
  - "self-managed"
  - "pipeline"
  - "custom-build"
category: "service-hub"
excerpt: "Take control of your CI/CD pipeline by writing Tekton YAML in your repository -- add SonarQube analysis, unit tests, or any custom build step while Planton handles triggering, credentials, and deployment."
---

# How to Switch to Self-Managed Tekton Pipelines

If you followed [How to Deploy Your First Service with Zero-Config CI/CD](/tutorials/how-to-deploy-your-first-service-with-zero-config-cicd), you have a Service running with a platform-managed pipeline. The platform builds your image, exports kustomize manifests, and deploys to your target cluster -- all without any pipeline configuration in your repository.

That works well until you need something the platform pipeline does not provide. You need SonarQube analysis before every production deploy. Or you need to run unit tests in the pipeline. Or your monorepo requires a custom build tool like Bazel. Platform-managed pipelines handle a fixed set of build steps, and they cannot be extended.

Self-managed pipelines solve this. You write a [Tekton](https://tekton.dev) pipeline YAML file in your repository, and Planton executes it on every commit. You get complete control over what your pipeline does -- while Planton continues to handle everything else: Git webhook triggering, credential provisioning, pipeline monitoring, and deployment.

This tutorial walks you through switching an existing platform-managed Service to a self-managed pipeline, using SonarQube code quality analysis as the motivating example.

> **Note**: The Planton web console provides a guided interface for configuring
> pipeline providers. This tutorial uses the CLI/YAML approach for stability and
> reproducibility. The console UI evolves frequently -- always check it for the
> latest experience.

## What You Will Learn

- When to switch to self-managed pipelines -- and when to stay with platform-managed
- What the platform still handles for you and what you now own
- How to write a Tekton pipeline YAML that satisfies the platform's contract
- How to add a SonarQube analysis step as a custom build task
- How to update your Service manifest to use the new pipeline
- How to verify that your self-managed pipeline builds, analyzes, and deploys correctly

## Prerequisites

- [ ] A working Service deployed through Planton with a platform-managed pipeline (see [How to Deploy Your First Service with Zero-Config CI/CD](/tutorials/how-to-deploy-your-first-service-with-zero-config-cicd))
- [ ] A Dockerfile in your repository (if you were using Buildpacks with platform-managed pipelines, you will need to create one -- self-managed pipelines give you full control over the build, and the Tekton Hub provides a BuildKit task for Dockerfile-based builds)
- [ ] A [SonarCloud](https://sonarcloud.io) account with a project configured for your repository (free for public repositories)
- [ ] A Kubernetes Secret named `sonar-token` in the `planton-cloud-pipelines` namespace containing your SonarCloud token (coordinate with your platform administrator to create this)
- [ ] The `planton` CLI installed and authenticated (`planton auth login`)
- [ ] Basic familiarity with [Tekton Pipeline](https://tekton.dev/docs/pipelines/) concepts (Pipeline, Task, PipelineRun)

## When and Why to Switch

Platform-managed pipelines handle everything automatically but offer limited customization. Switch to self-managed when you need custom build steps (code analysis, security scanning, integration tests) that the platform pipeline does not support. Your self-managed Tekton pipeline runs on the same Runner infrastructure and must follow a platform contract -- specific parameters, workspace layout, and task naming conventions. For the full contract specification, see the [self-managed pipelines documentation](/docs/ci-cd/self-managed-pipelines).

## Step 1: Write the Tekton Pipeline YAML

Create a file at `.planton/pipeline.yaml` in the root of your repository. This is the default path the platform looks for when `pipelineProvider` is set to `self`. You can use a different path by setting `tektonPipelineYamlFile` in your Service manifest.

Here is a complete pipeline that clones your repository, runs SonarQube analysis, builds a Docker image with BuildKit, and exports kustomize manifests for deployment:

```yaml
apiVersion: tekton.dev/v1
kind: Pipeline
metadata:
  name: service-pipeline-with-sonar
  namespace: planton-cloud-pipelines
spec:
  description: >
    Self-managed pipeline that adds SonarQube analysis before building and
    deploying a containerized service through Planton.

  params:
    - name: git-url
      type: string
      description: Repository clone URL (injected by platform)
    - name: git-revision
      type: string
      description: Git commit SHA (injected by platform)
      default: main
    - name: project-root
      type: string
      description: Service root relative to repository root (injected by platform)
      default: "."
    - name: sparse-checkout-directories
      type: string
      description: Comma-separated paths for sparse checkout (injected by platform)
      default: ""
    - name: image-name
      type: string
      description: Full image reference with registry, path, and tag (injected by platform)
    - name: kustomize-manifests-config-map-name
      type: string
      description: ConfigMap name for kustomize output (injected by platform)
    - name: kustomize-base-directory
      type: string
      description: Kustomize directory relative to project root (injected by platform)
      default: _kustomize
    - name: setup-package-credentials
      type: string
      description: Legacy flag for package credential setup (injected by platform)
      default: "false"
    - name: dockerfile-config-map-name
      type: string
      description: ConfigMap name for Dockerfile export (injected by platform)
    - name: owner-identifier-label-key
      type: string
      description: Label key for pipeline event routing (injected by platform)
      default: ""
    - name: owner-identifier-label-value
      type: string
      description: Label value for pipeline event routing (injected by platform)
      default: ""
    - name: sonar-project-key
      type: string
      description: SonarCloud project key (from Service params)
    - name: sonar-host-url
      type: string
      description: SonarCloud host URL (from Service params)
      default: "https://sonarcloud.io"

  workspaces:
    - name: source
      description: Workspace where source code is cloned
    - name: package-credentials
      description: Workspace for package manager credentials (provisioned by platform)

  tasks:
    # ──────────────────────────────────────────────────────────────────────
    # 1. Clone source from Git
    # ──────────────────────────────────────────────────────────────────────
    - name: git-checkout
      taskRef:
        resolver: git
        params:
          - name: url
            value: "https://github.com/plantonhq/tekton-hub.git"
          - name: revision
            value: "main"
          - name: pathInRepo
            value: "tasks/git-clone.yaml"
      params:
        - name: url
          value: $(params.git-url)
        - name: revision
          value: $(params.git-revision)
        - name: deleteExisting
          value: "true"
        - name: sparseCheckoutDirectories
          value: $(params.sparse-checkout-directories)
      workspaces:
        - name: output
          workspace: source

    # ──────────────────────────────────────────────────────────────────────
    # 2. SonarQube analysis (custom step -- this is why you switched)
    # ──────────────────────────────────────────────────────────────────────
    - name: sonar-analysis
      runAfter: [git-checkout]
      taskSpec:
        params:
          - name: project-key
          - name: host-url
          - name: project-root
        workspaces:
          - name: source
        steps:
          - name: scan
            image: sonarsource/sonar-scanner-cli:latest
            workingDir: $(workspaces.source.path)/$(params.project-root)
            env:
              - name: SONAR_TOKEN
                valueFrom:
                  secretKeyRef:
                    name: sonar-token
                    key: token
            script: |
              #!/usr/bin/env bash
              set -euo pipefail
              sonar-scanner \
                -Dsonar.host.url=$(params.host-url) \
                -Dsonar.projectKey=$(params.project-key)
      params:
        - name: project-key
          value: $(params.sonar-project-key)
        - name: host-url
          value: $(params.sonar-host-url)
        - name: project-root
          value: $(params.project-root)
      workspaces:
        - name: source
          workspace: source

    # ──────────────────────────────────────────────────────────────────────
    # 3. Build and push container image with BuildKit
    # ──────────────────────────────────────────────────────────────────────
    - name: build-image
      runAfter: [sonar-analysis]
      taskRef:
        resolver: git
        params:
          - name: url
            value: "https://github.com/plantonhq/tekton-hub.git"
          - name: revision
            value: "main"
          - name: pathInRepo
            value: "tasks/buildkit.yaml"
      params:
        - name: image
          value: $(params.image-name)
        - name: contextDir
          value: $(params.project-root)
        - name: dockerfile-config-map-name
          value: $(params.dockerfile-config-map-name)
        - name: owner-identifier-label-key
          value: $(params.owner-identifier-label-key)
        - name: owner-identifier-label-value
          value: $(params.owner-identifier-label-value)
      workspaces:
        - name: source
          workspace: source

    # ──────────────────────────────────────────────────────────────────────
    # 4. Export kustomize manifests for deployment
    # ──────────────────────────────────────────────────────────────────────
    - name: kustomize-build
      runAfter: [git-checkout]
      taskRef:
        resolver: git
        params:
          - name: url
            value: "https://github.com/plantonhq/tekton-hub.git"
          - name: revision
            value: "main"
          - name: pathInRepo
            value: "tasks/kustomize-build.yaml"
      params:
        - name: config-map-name
          value: $(params.kustomize-manifests-config-map-name)
        - name: config-map-namespace
          value: planton-cloud-pipelines
        - name: project-root
          value: $(params.project-root)
        - name: kustomize-base-directory
          value: $(params.kustomize-base-directory)
        - name: owner-identifier-label-key
          value: $(params.owner-identifier-label-key)
        - name: owner-identifier-label-value
          value: $(params.owner-identifier-label-value)
      workspaces:
        - name: source
          workspace: source
```

### How the Tasks Connect

The task dependency graph looks like this:

```
git-checkout ──┬──→ sonar-analysis ──→ build-image
               │
               └──→ kustomize-build
```

- `sonar-analysis` runs after `git-checkout` and blocks `build-image`. If the SonarQube scan fails, the image is not built and the pipeline fails. This makes SonarQube a quality gate.
- `kustomize-build` runs in parallel with the analysis and build chain. It processes your kustomize overlays and exports the resulting manifests to a ConfigMap. This is independent of the image build.
- The deployment stage (managed by the platform, not by your pipeline) waits for the entire pipeline to complete, then reads the kustomize ConfigMap and the built image to deploy your service.

### Understanding the Pipeline YAML

**Reusing Tekton Hub tasks**: Three of the four tasks (`git-checkout`, `build-image`, `kustomize-build`) reference tasks from the [Planton Tekton Hub](https://github.com/plantonhq/tekton-hub) via the Git resolver. These tasks are tested and maintained by the platform team. You do not need to write them from scratch. Use `revision: "main"` to pick up updates automatically, or pin to a specific commit SHA for stability.

**Inline custom tasks**: The `sonar-analysis` task uses `taskSpec` instead of `taskRef`. This means the task definition is inline -- it lives inside your pipeline YAML. This is the simplest way to add custom steps. For reusable tasks, you can extract them to separate YAML files in your repository and reference them via the Git resolver.

**SonarQube authentication**: The `sonar-analysis` step reads the `SONAR_TOKEN` environment variable from a Kubernetes Secret named `sonar-token`. This secret must exist in the pipeline namespace (`planton-cloud-pipelines`). Coordinate with your platform administrator to create it. Do not put secrets in pipeline parameters -- parameters are visible in pipeline logs.

## Step 2: Update the Service Manifest

With the pipeline YAML in your repository, update your Service manifest to tell the platform to use it.

If your current Service manifest looks like this (from Tutorial G):

```yaml
apiVersion: service-hub.planton.ai/v1
kind: Service
metadata:
  name: my-service
  org: your-org
spec:
  gitRepo:
    gitConnection: your-github-connection
    cloneUrl: https://github.com/your-org/your-repo.git
    defaultBranch: main
    gitRepoProvider: github
    name: your-repo
    ownerName: your-org
  packageType: container_image
  containerRegistry: your-registry-connection
  pipelineConfiguration:
    pipelineProvider: platform
    imageBuildMethod: buildpacks
    imageRepositoryPath: your-org/your-repo
    pipelineBranches:
      - main
  deploymentEnvironments:
    - dev
```

Change it to:

```yaml
apiVersion: service-hub.planton.ai/v1
kind: Service
metadata:
  name: my-service
  org: your-org
spec:
  gitRepo:
    gitConnection: your-github-connection
    cloneUrl: https://github.com/your-org/your-repo.git
    defaultBranch: main
    gitRepoProvider: github
    name: your-repo
    ownerName: your-org
  packageType: container_image
  containerRegistry: your-registry-connection
  pipelineConfiguration:
    pipelineProvider: self
    imageRepositoryPath: your-org/your-repo
    pipelineBranches:
      - main
    params:
      sonar-project-key: "your-org_your-repo"
      sonar-host-url: "https://sonarcloud.io"
  deploymentEnvironments:
    - dev
```

Three things changed:

1. **`pipelineProvider: self`** (was `platform`): This tells the platform to load your pipeline YAML from the repository instead of using a platform-managed pipeline from the Tekton Hub.

2. **`imageBuildMethod` removed**: This field is only applicable for platform-managed pipelines. With self-managed pipelines, you control the build method in your pipeline YAML. The platform ignores this field when `pipelineProvider` is `self`.

3. **`params` added**: The `params` map injects custom key-value pairs into your Tekton PipelineRun. Your pipeline YAML declares matching parameter names (`sonar-project-key`, `sonar-host-url`) and passes them to the `sonar-analysis` task. This is how you configure custom pipeline behavior without modifying the pipeline YAML itself.

The `tektonPipelineYamlFile` field is not shown because `.planton/pipeline.yaml` is the default path. If your pipeline YAML is at a different location (for example, `ci/pipelines/build.yaml` in a monorepo), set it explicitly:

```yaml
pipelineConfiguration:
  pipelineProvider: self
  tektonPipelineYamlFile: ci/pipelines/build.yaml
```

## Step 3: Apply the Updated Service

Apply the updated Service manifest:

```bash
planton apply -f service.yaml
```

This updates the Service configuration. The next push to a branch in `pipelineBranches` will use your self-managed pipeline.

## Step 4: Trigger and Monitor the Pipeline

Commit the `.planton/pipeline.yaml` file (and your Dockerfile, if you created one) and push to your default branch:

```bash
git add .planton/pipeline.yaml Dockerfile
git commit -m "Switch to self-managed pipeline with SonarQube analysis"
git push origin main
```

The push triggers a pipeline through the same Git webhook mechanism as before. Monitor it:

```bash
planton service last-pipeline --service my-service
```

This shows the pipeline metadata, including the pipeline ID. Stream the status:

```bash
planton follow <pipeline-id>
```

You will see your four tasks progress: `git-checkout`, `sonar-analysis`, `build-image`, and `kustomize-build`. If the SonarQube scan fails (for example, your code does not meet quality thresholds), `build-image` never runs and the pipeline fails.

To see detailed logs for a specific task:

```bash
planton service pipeline stream-logs <pipeline-id>
```

## Verifying Your Deployment

After the pipeline succeeds, verify that the deployment completed:

```bash
planton get service my-service -o yaml
```

Check the deployment status in the output. The deployment flow is identical to platform-managed pipelines: the platform reads the kustomize ConfigMap created by your `kustomize-build` task and applies the manifests to the target cluster.

Verify your SonarQube results by checking your SonarCloud dashboard at `https://sonarcloud.io/project/overview?id=your-org_your-repo`.

## Common Patterns and Tips

### Adding Unit Tests

Insert a test task between `git-checkout` and `build-image`:

```yaml
- name: unit-tests
  runAfter: [git-checkout]
  taskSpec:
    params:
      - name: project-root
    workspaces:
      - name: source
    steps:
      - name: test
        image: node:20
        workingDir: $(workspaces.source.path)/$(params.project-root)
        script: |
          #!/usr/bin/env bash
          set -euo pipefail
          npm ci
          npm test
  params:
    - name: project-root
      value: $(params.project-root)
  workspaces:
    - name: source
      workspace: source
```

Then update `build-image` to run after `unit-tests` instead of (or in addition to) `sonar-analysis`:

```yaml
- name: build-image
  runAfter: [sonar-analysis, unit-tests]
```

### Conditional Task Execution with Params

Use a Service-level param as a feature flag to conditionally run tasks:

```yaml
# In the Service manifest
params:
  run-sonar: "true"

# In the pipeline YAML
- name: sonar-analysis
  when:
    - input: $(params.run-sonar)
      operator: in
      values: ["true"]
  runAfter: [git-checkout]
  # ... rest of task spec
```

This lets you disable SonarQube analysis without modifying the pipeline YAML -- change the param in the Service manifest and apply.

### Custom Pipeline Path for Monorepos

If your repository contains multiple services, each service can have its own pipeline YAML:

```yaml
# Service A
pipelineConfiguration:
  pipelineProvider: self
  tektonPipelineYamlFile: services/api/.planton/pipeline.yaml

# Service B
pipelineConfiguration:
  pipelineProvider: self
  tektonPipelineYamlFile: services/worker/.planton/pipeline.yaml
```

### Version Pinning Tekton Hub Tasks

For production-critical services, pin task references to a specific commit instead of `main`:

```yaml
taskRef:
  resolver: git
  params:
    - name: url
      value: "https://github.com/plantonhq/tekton-hub.git"
    - name: revision
      value: "a1b2c3d4e5f6789..."
    - name: pathInRepo
      value: "tasks/buildkit.yaml"
```

Use `main` during development for automatic updates. Pin to a commit SHA when stability matters more than getting the latest improvements.

### Build-Only Mode

Combine self-managed pipelines with `disableDeployments: true` for services that only need image builds (no kustomize, no deployment):

```yaml
pipelineConfiguration:
  pipelineProvider: self
  disableDeployments: true
```

Your pipeline still runs the build tasks, but the platform skips the deployment stage entirely.

### Switching Back to Platform-Managed

If you decide self-managed pipelines are more maintenance than you need, switching back is straightforward:

```yaml
pipelineConfiguration:
  pipelineProvider: platform
  imageBuildMethod: dockerfile
```

Remove the `params` block (or leave it -- the platform ignores custom params for platform-managed pipelines). Your `.planton/pipeline.yaml` file can stay in the repository; it is ignored when `pipelineProvider` is `platform`.

## What to Do Next

- **[How to Store and Reference Credentials in a Provider Connection](/tutorials/how-to-store-and-reference-credentials-in-a-provider-connection)** -- learn how Planton's config-manager handles sensitive values across services and environments (coming soon)
- **[How to Configure Branch Deployments and Tag Releases](/tutorials/how-to-configure-branch-deployments-and-tag-releases)** -- configure PR builds, tag releases, and multi-environment deployment workflows (all pipeline configuration options work the same with self-managed pipelines)
- **[Tekton documentation](https://tekton.dev/docs/)** -- the official Tekton reference for Pipeline and Task authoring
- **[Planton Tekton Hub](https://github.com/plantonhq/tekton-hub)** -- browse available tasks (`tasks/`) and platform pipelines (`pipelines/`) for reference
