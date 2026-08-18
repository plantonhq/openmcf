---
title: "How to Connect a Container Registry to Planton"
date: "2026-04-02"
author:
  - name: "Planton Team"
    title: "Platform Engineering"
    bio: "Helping teams deploy infrastructure and services without the DevOps bottleneck"
tags:
  - "container-registry"
  - "connect"
  - "docker"
  - "ecr"
  - "artifact-registry"
  - "getting-started"
category: "connect"
excerpt: "Connect your container registry so Planton can push built images during CI and pull them during deployment."
---

# How to Connect a Container Registry to Planton

When Planton builds your service from source code, the result is a container image that needs to be stored somewhere. Your organization brings its own registry -- Google Cloud Artifact Registry, AWS Elastic Container Registry, or GitHub Container Registry -- and a Container Registry Connection provides the authentication Planton needs to push images during CI and pull them during deployment.

This tutorial walks you through creating a Container Registry Connection for each supported provider.

> **Note**: The Planton web console provides a guided creation wizard for container registry connections. Since the console is actively evolving, this tutorial focuses on the CLI/YAML path which remains stable. Always check the console for the latest experience.

## What You Will Learn

- What a container registry connection is and how Service Hub pipelines use it
- How to store registry credentials as organization-level secrets in Planton
- How to create a Container Registry Connection for GCP Artifact Registry, AWS ECR, or GitHub Container Registry
- How to verify the connection from the CLI

## Prerequisites

- [ ] A Planton account with an organization already created
- [ ] The `planton` CLI installed and authenticated (run `planton auth login` if you have not done this yet)
- [ ] One of the following, depending on which registry you use:
  - **GCP Artifact Registry**: A GCP project with Artifact Registry enabled and a service account key with Artifact Registry Writer role
  - **AWS ECR**: An AWS account with an ECR repository and an IAM user with ECR push/pull permissions
  - **GitHub Container Registry**: A GitHub account with a personal access token that has `write:packages` scope

## How Registry Connections Work

A container registry connection tells Planton where to push images that pipelines build and where Kubernetes clusters pull images from for deployment. This connection stores registry credentials as secrets and references them by slug. For more details, see the [container registry connections documentation](/docs/connections/container-registries).

### Supported Providers

| Provider | Status | Registry Hostname Pattern |
|---|---|---|
| GCP Artifact Registry | Supported | `<region>-docker.pkg.dev` |
| AWS ECR | Supported | `<account-id>.dkr.ecr.<region>.amazonaws.com` |
| GitHub Container Registry | Supported | `ghcr.io` |

Choose the path below that matches your registry provider.

## Path A: GCP Artifact Registry

### Step 1: Create a GCP Service Account

In your GCP project, create a service account with the **Artifact Registry Writer** role. This role allows both pushing images during builds and pulling images during deployments.

```bash
gcloud iam service-accounts create planton-registry \
  --display-name="Planton Container Registry"

gcloud projects add-iam-policy-binding YOUR_PROJECT_ID \
  --member="serviceAccount:planton-registry@YOUR_PROJECT_ID.iam.gserviceaccount.com" \
  --role="roles/artifactregistry.writer"
```

Download the service account key as a JSON file:

```bash
gcloud iam service-accounts keys create sa-key.json \
  --iam-account=planton-registry@YOUR_PROJECT_ID.iam.gserviceaccount.com
```

### Step 2: Store the Service Account Key as a Planton Secret

Use the `--from-file` flag to store the JSON key file content as an organization-level secret:

```bash
planton secret set gcr-sa-key --from-file value=./sa-key.json
```

After storing the secret, delete the local key file:

```bash
rm sa-key.json
```

### Step 3: Create the Container Registry Connection

Create a file named `container-registry.yaml`:

```yaml
apiVersion: connect.planton.ai/v1
kind: ContainerRegistryConnection
metadata:
  name: gcp-registry
  org: your-org
spec:
  provider: gcp_artifact_registry
  gcp_artifact_registry:
    gcp_project_id:
      value: "your-gcp-project-id"
    gcp_region:
      value: "us-central1"
    service_account_key:
      secret: gcr-sa-key
    gcp_artifact_registry_repo_name:
      value: "platform-images"
```

Replace the placeholder values:

- `metadata.name`: A descriptive name for this connection (becomes the slug you reference in Services)
- `metadata.org`: Your Planton organization slug
- `gcp_project_id.value`: Your GCP project ID
- `gcp_region.value`: The region where your Artifact Registry repository is located
- `service_account_key.secret`: The slug of the secret you created in Step 2
- `gcp_artifact_registry_repo_name.value`: The name of your Docker repository in Artifact Registry

### Step 4: Apply the Connection

```bash
planton apply -f container-registry.yaml
```

## Path B: AWS Elastic Container Registry

### Step 1: Create an IAM User with ECR Permissions

In your AWS account, create an IAM user with programmatic access and attach a policy that grants ECR permissions:

```json
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": [
        "ecr:GetAuthorizationToken",
        "ecr:BatchGetImage",
        "ecr:GetDownloadUrlForLayer",
        "ecr:BatchCheckLayerAvailability",
        "ecr:PutImage",
        "ecr:InitiateLayerUpload",
        "ecr:UploadLayerPart",
        "ecr:CompleteLayerUpload"
      ],
      "Resource": "*"
    }
  ]
}
```

Note the **Access Key ID** and **Secret Access Key** from the user creation.

### Step 2: Store Credentials as Planton Secrets

```bash
planton secret set ecr-access-key value=AKIAIOSFODNN7EXAMPLE
planton secret set ecr-secret-key value=wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY
```

### Step 3: Create the Container Registry Connection

Create a file named `container-registry.yaml`:

```yaml
apiVersion: connect.planton.ai/v1
kind: ContainerRegistryConnection
metadata:
  name: aws-registry
  org: your-org
spec:
  provider: aws_elastic_container_registry
  aws_elastic_container_registry:
    account_id:
      value: "123456789012"
    access_key_id:
      secret: ecr-access-key
    secret_access_key:
      secret: ecr-secret-key
    region:
      value: "us-east-1"
```

Replace the placeholder values:

- `metadata.name`: A descriptive name for this connection
- `metadata.org`: Your Planton organization slug
- `account_id.value`: Your 12-digit AWS account ID
- `access_key_id.secret`: The slug of the access key secret
- `secret_access_key.secret`: The slug of the secret key secret
- `region.value`: The AWS region where your ECR repository is located

### Step 4: Apply the Connection

```bash
planton apply -f container-registry.yaml
```

## Path C: GitHub Container Registry

### Step 1: Create a Personal Access Token

In GitHub, navigate to **Settings > Developer settings > Personal access tokens** and create a token with:

- **`write:packages`** scope (includes `read:packages`)
- **`delete:packages`** scope (optional, if you want Planton to clean up images)

Note the token value. You will not be able to see it again after creation.

### Step 2: Store the Token as a Planton Secret

```bash
planton secret set ghcr-pat value=ghp_xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
```

### Step 3: Create the Container Registry Connection

Create a file named `container-registry.yaml`:

```yaml
apiVersion: connect.planton.ai/v1
kind: ContainerRegistryConnection
metadata:
  name: github-registry
  org: your-org
spec:
  provider: github_container_registry
  github_container_registry:
    github_username:
      value: "your-github-username"
    personal_access_token:
      secret: ghcr-pat
```

Replace the placeholder values:

- `metadata.name`: A descriptive name for this connection
- `metadata.org`: Your Planton organization slug
- `github_username.value`: Your GitHub username or organization name
- `personal_access_token.secret`: The slug of the secret you created in Step 2

### Step 4: Apply the Connection

```bash
planton apply -f container-registry.yaml
```

## Verifying Your Connection

After applying any of the paths above, verify the connection was created:

```bash
planton get container-registry-connection gcp-registry
```

Replace `gcp-registry` with the name you used in your manifest.

The output shows the connection resource with the provider type, region, and secret references. Actual secret values are never displayed -- only the secret slugs.

## What to Do Next

- **Connect your GitHub repositories** to Planton if you have not already: [How to Connect Your GitHub Account to Planton](/tutorials/how-to-connect-your-github-account-to-planton)
- **Deploy your first backend service** with zero-config CI/CD: [How to Deploy Your First Service with Zero-Config CI/CD](/tutorials/how-to-deploy-your-first-service-with-zero-config-cicd)
