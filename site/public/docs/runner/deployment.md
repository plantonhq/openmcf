---
title: "Deployment"
description: "Generate credentials, install, and deploy Planton Runner to Kubernetes, AWS ECS, GCP Cloud Run, or Azure Container Apps"
icon: server
order: 30
tags:
  - Runner
  - Deployment
  - Kubernetes
  - AWS
  - GCP
  - Azure
---

# Deployment

Getting a runner operational involves three phases: registering it with Planton, generating credentials (which produces a single JSON credentials file), and deploying the runner to your infrastructure. This page walks through each phase.

## Generating Credentials

Before you can deploy a runner, you need to register it and generate its credentials. Registration creates the runner's identity in Planton. Credential generation produces a single JSON file containing everything the runner needs to authenticate: mTLS certificates for the secure tunnel connection, an API key for control-plane authentication, and organizational context.

During credential generation, Planton auto-provisions a [service account](/docs/security/authentication-and-authorization#service-accounts) for the runner and mints an API key. The runner uses this key to authenticate with the control plane for secret and variable resolution at runtime.

### Using the CLI

```bash
planton runner generate-credentials prod-runner
```

This generates credentials for the runner named `prod-runner` in your current organization. By default, the credentials JSON is printed to stdout. To save it to a directory:

```bash
planton runner generate-credentials prod-runner --output-dir ./runner-creds
```

This writes a single `credentials.json` file to the specified directory. The file contains the tunnel certificates, the API key, the control-plane endpoint, and the runner's organizational context.

You can also output credentials as a Kubernetes Secret for direct deployment:

```bash
planton runner generate-credentials prod-runner --output-format k8s-secret
```

Each call to `generate-credentials` mints a new API key. Previous keys remain valid, so you can generate multiple credentials files (for example, one per deployment environment) without invalidating existing ones.

### Using the Web Console

Navigate to **Connections** in the sidebar, then click **Deploy Runner**. The deployment wizard walks through:

1. **Name** — Choose a slug for the runner (e.g., `prod-runner`)
2. **Credentials** — Generate and download the credentials file
3. **Target** — Select your deployment target (Kubernetes, AWS, GCP, or Azure)
4. **Deploy Guide** — Follow the target-specific instructions

<!-- SCREENSHOT: Runner deployment wizard
  Page: /resource/connect/runner-registration/create
  Action: Show the wizard at the Target selection step
  Focus: The four deployment target options
  Alt: Runner deployment wizard showing Kubernetes, AWS ECS, GCP Cloud Run, and Azure Container Apps as deployment target options
-->

### Credential Security

The credentials file contains sensitive material: the runner's private key and an API key. Both are generated server-side and returned only once. Planton does not store the private key or the raw API key after generation.

If you lose the credentials file, generate a new one. If you need to invalidate all previous credentials (for example, after a security incident), use `regenerate-credentials`:

```bash
planton runner regenerate-credentials prod-runner --output-dir ./runner-creds
```

Regeneration creates new certificates and a new API key, and revokes all previous API keys for the runner's service account. Any deployed runner using old credentials will lose both tunnel connectivity and control-plane authentication. Plan for a brief connectivity gap during credential rotation.

Store the credentials file securely using your platform's secrets management:

| Deployment Target | Recommended Storage |
|-------------------|-------------------|
| Kubernetes | Kubernetes Secret with a `credentials.json` key |
| AWS ECS | AWS Secrets Manager |
| GCP Cloud Run | GCP Secret Manager |
| Azure Container Apps | Azure Key Vault |

## Installation

To run a runner locally (for development or testing), install the binary:

```bash
planton runner install
```

This downloads the runner binary. To install a specific version:

```bash
planton runner install --version 0.1.5
```

## Starting a Runner Locally

For development or testing, you can start a runner directly from your machine. The runner resolves credentials in this order:

1. `--credentials-file` flag or `PLANTON_RUNNER_CREDENTIALS_FILE` environment variable pointing to the JSON file
2. `PLANTON_RUNNER_CREDENTIALS` environment variable containing the raw JSON inline
3. Named runner from the local credential store (scans `~/.planton/`)

The simplest approach after generating credentials with `--output-dir`:

```bash
planton runner start --credentials-file ./runner-creds/credentials.json
```

Or if the credentials are already in the local store:

```bash
planton runner start prod-runner --org acme
```

The runner loads the credentials, establishes the secure tunnel, authenticates with the control plane, and begins accepting operations.

## Deployment Targets

For production use, deploy the runner as a long-running container. The CLI supports four deployment targets through `planton runner deploy`:

```bash
planton runner deploy prod-runner
```

This command prompts you to select a deployment target and walks through the target-specific configuration. You can also specify `--image-tag` to pin a specific runner version.

### Kubernetes

Deploys the runner using a Helm chart. This is the most common deployment target.

**Prerequisites:**
- `kubectl` configured with access to the target cluster
- Helm 3 installed
- Runner credentials (generated during registration)

**What gets deployed:**
- A Deployment running the runner container
- A Secret containing the `credentials.json` file
- A Kubernetes ServiceAccount (if using runner-delegated auth with Workload Identity or IRSA)

**When to use Kubernetes:**
- Your infrastructure operations target Kubernetes clusters
- You want the runner co-located with the workloads it manages
- You are already running Kubernetes and want to keep operations consolidated

### AWS ECS (Fargate)

Deploys the runner as a Fargate task in Amazon ECS. The runner runs as a single container with no EC2 instances to manage.

**Prerequisites:**
- AWS credentials with permissions to create ECS resources
- A VPC with subnets that have outbound internet access (for the tunnel connection)

**When to use ECS:**
- Your infrastructure is primarily on AWS
- You want a serverless deployment without managing Kubernetes
- You need the runner to have AWS IAM role-based access via ECS task roles

### GCP Cloud Run

Deploys the runner as an always-on Cloud Run service. Cloud Run's minimum instance count is set to ensure the runner stays warm for tunnel connectivity.

**Prerequisites:**
- GCP credentials with permissions to deploy Cloud Run services
- A project with the Cloud Run API enabled

**When to use Cloud Run:**
- Your infrastructure is primarily on GCP
- You want a managed container runtime without Kubernetes overhead
- You need the runner to use GCP Workload Identity for authentication

### Azure Container Apps

Deploys the runner as an always-on Container App. Similar to Cloud Run, minimum replicas are configured to maintain the tunnel connection.

**Prerequisites:**
- Azure credentials with permissions to create Container App resources
- A Container Apps Environment in your target region

**When to use Container Apps:**
- Your infrastructure is primarily on Azure
- You want a managed container runtime with Azure-native identity integration
- You need the runner to use Managed Identity for authentication

## Default Runner Binding

Once a runner is deployed, you can set it as the default for your organization. When a [connection](/docs/connections) does not specify an explicit runner, Planton automatically routes requests to the default runner.

### Setting a Default

```bash
planton runner set-default prod-runner
```

This sets `prod-runner` as the default runner for your current organization.

### Resolution Chain

When Planton needs to route a request through a runner, it resolves which runner to use in this order:

1. **Explicit runner on the connection** — If the connection's configuration specifies a runner, that runner is used.
2. **Organization default** — If no explicit runner is set, Planton looks for a default runner binding at the organization level.
3. **Platform default** — If no organization default exists, Planton falls back to the platform-level default (set by platform operators).

If no runner can be resolved at any level, the request fails with an error indicating that no runner is available.

### Managing Defaults

```bash
# View the effective default (shows resolution: org → platform)
planton runner get-default

# Remove the organization default
planton runner unset-default

# Set a platform-level default (requires platform operator permissions)
planton runner set-default shared-runner --platform

# View the platform default
planton runner get-default --platform
```

## Verifying Connectivity

After deploying a runner, verify that it has connected successfully. The simplest way is to run a Cloud Ops command that routes through the runner:

```bash
planton kubectl get pods --connection my-k8s-connection -n default
```

If the runner is connected and the connection is configured correctly, you will see the pod listing from your cluster. If the runner is not connected, you will receive a connection error indicating the runner is unreachable.

## Runner Management

### Listing Runners

```bash
# List runners in your organization
planton runner list

# Include platform runners
planton runner list --all

# List platform runners only
planton runner list --platform
```

### Viewing Runner Details

```bash
planton runner get prod-runner
```

### Deleting a Runner

```bash
planton runner delete prod-runner
```

This removes the runner registration and its auto-provisioned service account (which cascades to revoke all API keys). Any deployed runner using those credentials will lose both tunnel connectivity and control-plane authentication. The deployed container or service must be cleaned up separately.

## Related Documentation

- [Runner Overview](/docs/runner) — What Runner is and why it exists
- [Security Model](/docs/runner/security-model) — Credential isolation, authentication modes, and trust boundaries
- [Connections](/docs/connections) — Managing connections that route through runners
