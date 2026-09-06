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

Container image registries store the Docker images that Service Hub produces during builds. Every service that uses Buildpacks or Dockerfile builds needs a container registry connection so the built image can be pushed.

### Two Ways a Registry Connection Signs In

A registry connection describes how Planton signs in to a registry; it never holds the result. It signs in one of two ways:

- **A connection you already trust** — the registry names a connection your organization already has: a GitHub connection for GHCR, an AWS connection for ECR, a GCP connection for Artifact Registry, an Azure connection for ACR. The connection stores no key. At the moment a build pushes, the runner derives a short-lived registry token from that connection's own credential — the sign-in your laptop already holds, or the cloud identity a hosted runner already has. Sign out of the trusted connection and the registry stops working the same moment.
- **Stored keys** — the registry references a secret you created: a service account key, IAM access keys, a service principal, a personal access token. The runner reads the secret at push time. This is the arm for JFrog Artifactory, and for any registry that sits outside the connections you have.

The Container Registry card on the Connections page opens on **You Already Trust These** — one card per registry the connections in your organization can reach — and **Use This Registry** writes a connection that stores nothing. **Set Up Manually** reaches the wizard, whose **Credential** step asks the same question: **A Connection You Already Trust** or **Stored Keys**.

Planton supports five container image registry providers:

### GCP Artifact Registry

Google's managed container registry. This is the recommended choice if your deployment targets run on GCP (GKE, Cloud Run).

| Field | Description |
|-------|-------------|
| GCP Project ID | The Google Cloud project that hosts the registry |
| GCP Region | The region where the registry is located |
| Repository Name | The Artifact Registry repository name |
| Credential | A GCP connection you already trust, or a service account key (JSON) with the Artifact Registry Writer role, stored as a secret |

### AWS Elastic Container Registry (ECR)

Amazon's managed container registry. The natural choice for ECS and EKS deployments.

| Field | Description |
|-------|-------------|
| Account ID | Your 12-digit AWS account number |
| Region | The AWS region where the ECR registry is hosted |
| Credential | An AWS connection you already trust, or an IAM access key pair with ECR push permissions, stored as secrets |

### Azure Container Registry (ACR)

Microsoft's managed container registry. Pairs with AKS and Azure Container Apps deployments.

| Field | Description |
|-------|-------------|
| Login Server | The registry's login server, `<name>.azurecr.io` |
| Credential | An Azure connection you already trust, or a service principal (tenant, client ID, client secret) with the AcrPush role, stored as secrets |

### JFrog Artifactory

A self-hosted or cloud-hosted universal artifact repository. Useful for organizations that standardize on JFrog for all artifact management.

| Field | Description |
|-------|-------------|
| Server URL | The Artifactory instance, e.g. `https://mycompany.jfrog.io/artifactory` |
| Repository Key | The Docker repository within Artifactory |
| Username | The Artifactory user or service account |
| Access Token | An Artifactory access token, stored as a secret |

### GitHub Container Registry (GHCR)

GitHub's container registry, tightly integrated with GitHub Packages.

| Field | Description |
|-------|-------------|
| Credential | The GitHub connection you already trust — the sign-in on the machine running Planton, which needs the `write:packages` scope — or a GitHub username with a personal access token carrying `write:packages`, stored as a secret |

### Connecting via the Web Console

1. Navigate to **Connections** and click the **Container Registry** card under Artifact Stores.
2. **You Already Trust These** lists the registries your connections reach — *GHCR as priya-dev · via GitHub connection github-account-priya-dev*, *ECR us-east-1 for 123456789012 · via AWS connection aws-profile-acme-dev*. Fill in anything the connection cannot know (a region, a repository name) on the card, and click **Use This Registry**. Planton writes the connection and proves it: the runner derives the credential once and the registry accepts it, or you read the exact sentence saying why not.
3. Or choose **Set Up Manually**: pick the provider, answer **Credential** (**A Connection You Already Trust** or **Stored Keys**), and provide the fields listed above.
4. On the completed screen — and on the connection's page any time after — **Verify This Registry** asks the registry to accept the derived credential right now.

<!-- SCREENSHOT: You Already Trust These
  Page: /orgs/{org}/connections/create/container-registry
  Action: Show the registry cards the organization's trusted connections reach, with Use This Registry
  Focus: The GHCR card via a GitHub connection and the ECR card via an AWS connection
  Alt: Container registry create page listing the registries the organization's trusted connections reach, one confirm card each
-->

### Connecting via the CLI

```bash
# Which registries do my connections reach?
planton connect registry detect

# GHCR through the GitHub sign-in on this machine, without the confirmation prompt
planton connect registry detect --from github-account-priya-dev --yes

# ECR in a region through an AWS connection
planton connect registry detect --from aws-profile-acme-dev --region us-east-1
```

The verb writes the same credential-free connection the card writes, then proves it and reads **Registry Connection Ready — ghcr.io accepted the credential derived from GitHub connection github-account-priya-dev**.

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
- **Already connected a cloud or GitHub?** Its registry is one click away with nothing to paste — GHCR through your GitHub sign-in, ECR / Artifact Registry / ACR through the cloud connection you already have.

## Related Documentation

- [Connections Overview](/docs/connections) — Understanding the Connect system
- [Git Providers](/docs/connections/git-providers) — Connect source code repositories
- [CI/CD: Build Methods](/docs/ci-cd/build-methods) — How Planton builds container images
