---
title: "Container Registries"
description: "Connect container image registries and package repositories so Planton can push build artifacts and resolve private dependencies"
icon: connect
order: 40
tags:
  - Connect
  - Container Registry
  - Docker
  - NPM
  - Maven
---

# Container Registries

When Service Hub builds your code, it produces artifacts — container images for most services, or script bundles for Cloudflare Workers. Those artifacts need to be stored somewhere accessible to the deployment target. Container registry connections tell Planton where to push those artifacts and how to authenticate with the registry.

This page covers three types of registry connections: container image registries (for Docker/OCI images), NPM registries (for JavaScript/TypeScript packages), and Maven repositories (for Java artifacts).

## Container Image Registries

Container image registries store the Docker images that Service Hub produces during builds. Every service that uses Buildpacks or Dockerfile builds needs a container registry connection so the built image can be pushed and later pulled by the deployment target.

Planton supports five container image registry providers:

### GCP Artifact Registry

Google's managed container registry. This is the recommended choice if your deployment targets run on GCP (GKE, Cloud Run).

| Field | Description |
|-------|-------------|
| GCP Project ID | The Google Cloud project that hosts the registry |
| GCP Region | The region where the registry is located |
| Repository Name | The Artifact Registry repository name |
| Service Account Key | Base64-encoded JSON key for a service account with Artifact Registry Writer role |

### AWS Elastic Container Registry (ECR)

Amazon's managed container registry. The natural choice for ECS and EKS deployments.

| Field | Description |
|-------|-------------|
| Account ID | Your 12-digit AWS account number |
| Access Key ID | IAM user access key with ECR permissions |
| Secret Access Key | IAM user secret key |
| Region | The AWS region where the ECR registry is hosted |

### Azure Container Registry

Microsoft's managed container registry. Pairs with AKS and Azure Container Instances deployments.

> **Note**: Azure Container Registry connection support is defined but the configuration interface is still being implemented. Check back for updates.

### JFrog Artifactory

A self-hosted or cloud-hosted universal artifact repository. Useful for organizations that standardize on JFrog for all artifact management.

> **Note**: JFrog Artifactory connection support is defined but the configuration interface is still being implemented. Check back for updates.

### GitHub Container Registry (GHCR)

GitHub's container registry, tightly integrated with GitHub Actions and GitHub Packages.

| Field | Description |
|-------|-------------|
| GitHub Username | Your GitHub username or a machine account username |
| Personal Access Token | A GitHub PAT with `write:packages` scope |

### Connecting via the Web Console

1. Navigate to **Connections** and click the **Docker** card under DevOps Pipeline.
2. **Name your connection**.
3. **Select your registry provider** — GCP Artifact Registry, AWS ECR, Azure Container Registry, JFrog Artifactory, or GitHub Container Registry.
4. **Provide the provider-specific credentials** listed above.
5. **Create the connection**.

<!-- SCREENSHOT: Container registry connection form
  Page: /resource/connect/container-registry-connection/create
  Action: Show the registry provider selector with all five options
  Focus: The provider selection dropdown and provider-specific form fields
  Alt: Docker registry connection form with provider selection showing GCP Artifact Registry, AWS ECR, and other options
-->

### How Container Registries Are Used During Builds

When a Service Hub pipeline runs:

1. The build stage compiles your code and produces a container image.
2. The image is tagged with the commit SHA and pushed to the connected registry.
3. The deployment stage pulls the image from the registry and deploys it to the target environment.

The image path follows the convention: `<registry-host>/<repository-path>:<commit-sha>`. Planton constructs this path automatically from the registry connection configuration.

---

## NPM Registries

NPM registry connections allow Service Hub builds to resolve private JavaScript and TypeScript packages. If your monorepo or service depends on internal packages published to a private NPM registry, you need this connection so the build stage can install dependencies.

Planton supports three NPM registry providers:

| Provider | What You Need |
|----------|--------------|
| **GCP Artifact Registry** | GCP Project ID, region, repository name, service account key |
| **GitHub Packages** | GitHub username, personal access token with `read:packages` scope |
| **JFrog Artifactory** | Configuration interface being implemented |

### Connecting via the Web Console

1. Navigate to **Connections** and click the **NPM** card under DevOps Pipeline.
2. **Name your connection**, select your provider, and provide the credentials.
3. **Create the connection**.

---

## Maven Repositories

Maven repository connections allow Service Hub builds to resolve private Java artifacts. If your services depend on internal libraries published to a private Maven repository, you need this connection.

Planton supports three Maven repository providers:

| Provider | What You Need |
|----------|--------------|
| **GCP Artifact Registry** | GCP Project ID, region, repository name, service account key |
| **GitHub Packages** | GitHub username, personal access token with `read:packages` scope |
| **JFrog Artifactory** | Configuration interface being implemented |

### Connecting via the Web Console

1. Navigate to **Connections** and click the **Maven** card under DevOps Pipeline.
2. **Name your connection**, select your provider, and provide the credentials.
3. **Create the connection**.

---

## Cloudflare Worker Script Storage

Cloudflare Workers are deployed differently from container-based services — they use script bundles stored in R2 buckets rather than container images. If you deploy Cloudflare Workers through Service Hub, you need a Wrangler connection that specifies the R2 bucket for storing script bundles.

This connection appears under the DevOps Pipeline category as **Wrangler** in the Connections page.

| Field | Description |
|-------|-------------|
| R2 Access Key ID | Access key for R2 API operations |
| R2 Secret Access Key | Secret key for R2 API operations |
| R2 Endpoint | Custom endpoint URL (optional) |

---

## Choosing Your Registry

If you're starting fresh and choosing which registry to use, a practical guideline:

- **Deploying to AWS (EKS, ECS)?** Use ECR — it integrates natively with AWS IAM and avoids cross-cloud network transfer costs.
- **Deploying to GCP (GKE, Cloud Run)?** Use GCP Artifact Registry — same benefits for the GCP ecosystem.
- **Deploying to multiple clouds?** Use GitHub Container Registry or JFrog Artifactory for a cloud-neutral registry.
- **Already using JFrog?** Continue using it — Planton connects to it the same way your existing tooling does.

## Related Documentation

- [Connections Overview](/docs/connections) — Understanding the Connect system
- [Git Providers](/docs/connections/git-providers) — Connect source code repositories
- [CI/CD: Build Methods](/docs/ci-cd/build-methods) — How Planton builds container images
