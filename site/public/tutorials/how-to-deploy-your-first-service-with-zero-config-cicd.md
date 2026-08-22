---
title: "How to Deploy Your First Service with Zero-Config CI/CD"
date: "2026-04-02"
author:
  - name: "Planton Team"
    title: "Platform Engineering"
    bio: "Helping teams deploy infrastructure and services without the DevOps bottleneck"
tags:
  - "service-hub"
  - "ci-cd"
  - "buildpacks"
  - "getting-started"
  - "pipeline"
  - "kustomize"
category: "service-hub"
excerpt: "Push code to GitHub and get a running deployment on Kubernetes -- no Dockerfile, no pipeline YAML, no CI/CD configuration."
---

# How to Deploy Your First Service with Zero-Config CI/CD

This tutorial takes you from a GitHub repository containing application code to a running deployment on Kubernetes. You will not write a Dockerfile. You will not author pipeline YAML. You will not configure a CI/CD system. Planton handles all of that.

You will do two things: tell Planton where your code lives and where you want it deployed. The platform auto-detects your application's language using Cloud Native Buildpacks, builds a container image, and deploys it to Kubernetes through a fully managed pipeline. When you push code to your default branch, the entire cycle repeats automatically.

> **Note**: The Planton web console provides a guided creation wizard for Services.
> This tutorial uses the CLI/YAML approach for stability and reproducibility.
> The console UI evolves frequently -- always check it for the latest experience.

> **A note on field casing**: Planton manifests accept both camelCase and snake_case
> field names (standard protobuf JSON mapping). This tutorial uses camelCase, matching
> the convention used by Planton presets and ops manifests. Either style works.

## What You Will Learn

- Create a Service resource that links your GitHub repository to Planton's CI/CD system
- Structure kustomize overlays that define where and how your application deploys
- Monitor a platform-managed pipeline from image build through deployment
- Verify your deployed service is running on Kubernetes
- Trigger a new deployment by pushing a code change

## Prerequisites

- [ ] A Planton organization with at least one environment configured
- [ ] A GitHub connection (see [How to Connect Your GitHub Account to Planton](/tutorials/how-to-connect-your-github-account-to-planton))
- [ ] A container registry connection (see [How to Connect a Container Registry to Planton](/tutorials/how-to-connect-a-container-registry-to-planton))
- [ ] A Kubernetes cluster accessible from your Planton organization
- [ ] A Git repository with a [Cloud Native Buildpacks](https://paketo.io/)-compatible application (Node.js, Go, Java, Python, Ruby, .NET, PHP, or other Paketo-supported language)
- [ ] The `planton` CLI installed and authenticated (`planton auth login`)

## How It Works

A Planton Service connects your Git repository to a CI/CD pipeline that builds container images and deploys them to Kubernetes using kustomize overlays. You define the Service in a YAML manifest pointing to your repo, and Planton handles the rest -- building on every push, exporting manifests, and applying them to your target cluster. For more on how Services, pipelines, and kustomize overlays work together, see the [CI/CD documentation](/docs/ci-cd/what-is-a-service).

## Step 1: Prepare Your Repository

Your application repository needs a `_kustomize/` directory that tells Planton how to deploy your service. You can scaffold this automatically with the CLI or create it by hand.

### Option A: Scaffold with the CLI

If you have the `planton` CLI and are inside your Git repository, run:

```bash
planton service kustomize init --new
```

The command prompts you to choose a deployment platform (select **KubernetesDeployment**) and an environment. It generates the base and overlay structure with sensible defaults.

After running, you will see:

```text
_kustomize/
├── base/
│   ├── kustomization.yaml
│   └── service.yaml
└── overlays/
    └── <your-chosen-env>/
        ├── kustomization.yaml
        └── service.yaml
```

Review the generated files and adjust the values (container port, resource limits) to match your application before continuing to Step 2.

### Option B: Create the Structure Manually

Create the following directory structure in your repository root:

```text
_kustomize/
├── base/
│   ├── kustomization.yaml
│   └── service.yaml
└── overlays/
    └── dev/
        ├── kustomization.yaml
        └── service.yaml
```

**`_kustomize/base/kustomization.yaml`**

```yaml
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization

resources:
  - service.yaml
```

**`_kustomize/base/service.yaml`**

This is a `KubernetesDeployment` manifest -- a Planton Resource that defines how your application runs on Kubernetes. The base contains configuration shared across all environments.

```yaml
apiVersion: kubernetes.planton.dev/v1
kind: KubernetesDeployment
metadata:
  name: my-service
  org: your-org
spec:
  version: main
  container:
    app:
      image:
        repo: ghcr.io/your-github-org/your-repo
      resources:
        requests:
          cpu: 50m
          memory: 100Mi
        limits:
          cpu: 1000m
          memory: 1Gi
      ports:
        - name: rest
          containerPort: 8080
          networkProtocol: TCP
          appProtocol: http
          servicePort: 80
  availability:
    minReplicas: 1
```

A few things to note about this manifest:

- **`image.repo`** is the container registry path where your built image will be stored. Set this to match the registry and path you configured in your container registry connection. The image **tag** is not specified here -- the pipeline automatically sets it to the Git commit SHA when deploying.
- **`containerPort`** should match the port your application listens on. The example uses 8080; adjust to match your application.
- **`resources`** sets the CPU and memory requests and limits for your container. Start with conservative values and adjust based on actual usage.

**`_kustomize/overlays/dev/kustomization.yaml`**

```yaml
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization

resources:
  - ../../base

patches:
  - path: service.yaml
```

**`_kustomize/overlays/dev/service.yaml`**

The overlay patches only what differs from the base. At minimum, it sets `metadata.env` to identify the target Planton environment:

```yaml
apiVersion: kubernetes.planton.dev/v1
kind: KubernetesDeployment
metadata:
  name: my-service
  env: dev
```

The `env` field tells Planton which environment this deployment targets. The platform resolves the Kubernetes cluster and namespace from your environment's configuration -- you do not need to specify a target cluster in the manifest.

### Verify Locally

Before proceeding, confirm that your kustomize structure is valid:

```bash
kubectl kustomize _kustomize/overlays/dev/
```

This command merges the overlay with the base and prints the resulting manifest. Verify that the output contains your base configuration with `metadata.env: dev` applied.

Commit and push the `_kustomize/` directory to your repository.

## Step 2: Create the Service Manifest

Create a file named `service.yaml` (this file lives outside your application repository -- it is a Planton API resource, not part of your application code):

```yaml
apiVersion: service-hub.planton.ai/v1
kind: Service
metadata:
  name: my-service
  org: your-org
spec:
  description: My first backend service on Planton
  gitRepo:
    gitConnection: your-github-connection
    ownerName: your-github-org
    name: your-repo
    defaultBranch: main
  packageType: container_image
  containerRegistry: your-registry-connection
  pipelineConfiguration:
    pipelineProvider: platform
    imageBuildMethod: buildpacks
    imageRepositoryPath: your-github-org/your-repo
    pipelineBranches:
      - main
  deploymentEnvironments:
    - dev
```

Replace the placeholder values with your actual configuration:

| Field | What to set | Where to find it |
|-------|-------------|------------------|
| `metadata.name` | A name for your service (lowercase, hyphens allowed) | You choose this |
| `metadata.org` | Your Planton organization slug | `planton get organization` |
| `gitRepo.gitConnection` | The slug of your GitHub connection | `planton get github-connection` |
| `gitRepo.ownerName` | Your GitHub organization or username | Your GitHub account |
| `gitRepo.name` | Your repository name | Your GitHub repository |
| `containerRegistry` | The slug of your container registry connection | `planton get container-registry-connection` |
| `pipelineConfiguration.imageRepositoryPath` | The path within your registry where images are pushed | Must match `image.repo` in your kustomize base |
| `deploymentEnvironments` | Which overlay directories to deploy | Must match directory names under `_kustomize/overlays/` |

Here is what each section does:

**`gitRepo`**: Links this Service to your GitHub repository. `gitConnection` is the slug of the GitHub connection you created in the [GitHub connection tutorial](/tutorials/how-to-connect-your-github-account-to-planton). Planton uses this connection to clone your code, register a webhook for push notifications, and report build status back to GitHub.

**`packageType: container_image`**: Tells Planton that this service produces a container image. This requires both a `containerRegistry` connection (where to push the image) and an `imageRepositoryPath` (the path within that registry).

**`pipelineConfiguration`**: Controls how the pipeline runs.
- `pipelineProvider: platform` means Planton manages the entire pipeline -- you do not write any pipeline YAML.
- `imageBuildMethod: buildpacks` means Planton uses Cloud Native Buildpacks (specifically the `paketobuildpacks/builder-jammy-base` builder) to detect your language and build a container image. No Dockerfile required.
- `pipelineBranches` lists the branches that trigger a pipeline when pushed to. Only pushes to branches in this list trigger the pipeline.

**`deploymentEnvironments`**: Filters which kustomize overlays the deploy stage processes. In this example, only the `dev` overlay is deployed. If you leave this field empty, the platform processes all overlays found under `_kustomize/overlays/`.

## Step 3: Apply the Service and Watch the Initial Pipeline

Deploy the Service:

```bash
planton apply -f service.yaml
```

When a Service is created for the first time, two things happen automatically:

1. Planton registers a webhook on your GitHub repository so future pushes trigger pipelines
2. An initial pipeline is triggered immediately using the latest commit on your default branch

Retrieve the pipeline ID:

```bash
planton service last-pipeline --service my-service
```

This returns the ID of the most recent pipeline for your service. Use it to watch the pipeline in real time:

```bash
planton follow <pipeline-id>
```

The `follow` command streams pipeline progress to your terminal. You will see each stage execute:

1. **git-checkout**: Clones your repository at the latest commit
2. **build-image** and **kustomize-build** (running in parallel):
   - `build-image` uses Buildpacks to detect your language, compile your application, and push a container image tagged with the commit SHA
   - `kustomize-build` merges each overlay with the base, base64-encodes the result, and stores it in a Kubernetes ConfigMap
3. **deploy**: Takes the kustomize output, applies the built image tag, and creates a Cloud Resource deployment for each environment listed in `deploymentEnvironments`

If you prefer to stream only the build logs (useful for debugging build failures), use:

```bash
planton service pipeline stream-logs <pipeline-id>
```

## Step 4: Verify the Deployment

Once the pipeline completes, verify that your service is deployed:

```bash
planton get service my-service -o yaml
```

In the output, look for the `status.envDeploymentMap` section. It shows which environments have active deployments and the Cloud Resource ID for each:

```yaml
status:
  envDeploymentMap:
    dev:
      resourceKind: KubernetesDeployment
      resourceId: cr_01abc...
```

To see the deployment history for your service:

```bash
planton service deployments my-service
```

This lists all deployments with their timestamps, commit SHAs, and statuses.

## Step 5: Push a Code Change and Watch the Pipeline

Make a change to your application code, commit it, and push to the branch listed in `pipelineBranches`:

```bash
git add .
git commit -m "update: my first change"
git push origin main
```

The webhook fires, and Planton starts a new pipeline. Retrieve and follow it:

```bash
planton service last-pipeline --service my-service
planton follow <pipeline-id>
```

The new pipeline builds a fresh container image tagged with the new commit SHA and deploys it to the `dev` environment. Each push to `main` produces a new deployment with a traceable image tag.

## Common Patterns and Tips

### Switching from Buildpacks to Dockerfile

If your application needs a custom build process, switch to a Dockerfile by changing two fields in the Service manifest:

```yaml
pipelineConfiguration:
  imageBuildMethod: dockerfile
  dockerfilePath: Dockerfile    # relative to projectRoot; defaults to "Dockerfile"
```

Everything else -- the pipeline, the kustomize structure, the deployment flow -- remains the same.

### Monorepo Support

If your service lives in a subdirectory of a monorepo, use these `gitRepo` fields:

```yaml
gitRepo:
  gitConnection: your-github-connection
  ownerName: your-github-org
  name: your-monorepo
  defaultBranch: main
  projectRoot: services/my-service
  triggerPaths:
    - services/my-service
  sparseCheckoutDirectories:
    - services/my-service
```

- `projectRoot` tells the pipeline where to find `_kustomize/` and where to run the build
- `triggerPaths` limits pipeline triggers to pushes that touch files under the specified paths
- `sparseCheckoutDirectories` speeds up the git clone by only checking out the directories the pipeline needs

### Enabling Ingress

To expose your service externally, add ingress configuration to your kustomize overlay:

```yaml
# _kustomize/overlays/dev/service.yaml
apiVersion: kubernetes.planton.dev/v1
kind: KubernetesDeployment
metadata:
  name: my-service
  env: dev
spec:
  ingress:
    enabled: true
    hostname: my-service-dev.your-domain.com
  container:
    app:
      ports:
        - name: rest
          containerPort: 8080
          networkProtocol: TCP
          appProtocol: http
          servicePort: 80
          isIngressPort: true
```

Note the `isIngressPort: true` flag on the port -- this tells the ingress controller which port to route traffic to.

### Adding a Production Environment

Create a new overlay directory and update the Service:

1. Create `_kustomize/overlays/production/kustomization.yaml` and `_kustomize/overlays/production/service.yaml` with production-specific settings (higher resource limits, more replicas, production secrets)

2. Add `production` to `deploymentEnvironments` in your Service manifest:

```yaml
deploymentEnvironments:
  - dev
  - production
```

3. Apply the updated Service manifest:

```bash
planton apply -f service.yaml
```

Future pipelines will deploy to both `dev` and `production` environments.

### Overriding the Kustomize Directory

By default, the pipeline looks for `_kustomize/` at the root of your `projectRoot`. If your kustomize directory is elsewhere, set it in the pipeline configuration:

```yaml
pipelineConfiguration:
  kustomizeBaseDirectory: deploy/kustomize
```

## What to Do Next

- **[How to Deploy Redis on Kubernetes](/tutorials/how-to-deploy-redis-on-kubernetes)** -- provision backing infrastructure for your service
- **Configure branch-based deployments** -- deploy feature branches to isolated environments (Tutorial H, coming soon)
- **Manage secrets and variables** -- inject sensitive configuration into your deployments (Tutorial J, coming soon)
- **Switch to self-managed pipelines** -- take full control of your pipeline with custom Tekton YAML (Tutorial I, coming soon)
