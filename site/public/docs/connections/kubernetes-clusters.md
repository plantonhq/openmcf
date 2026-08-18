---
title: "Kubernetes Clusters"
description: "Connect external Kubernetes clusters for service deployments, infrastructure management, and Cloud Ops access"
icon: kubernetes
order: 60
tags:
  - Connect
  - Kubernetes
  - GKE
  - EKS
  - AKS
  - DOKS
---

# Kubernetes Clusters

Planton can create Kubernetes clusters for you through Infra Hub (GKE, EKS, AKS, DOKS). But many organizations already have clusters running — production clusters that predate Planton, clusters managed by a separate platform team, or clusters in environments that Planton doesn't manage directly.

Kubernetes cluster connections let you bring those existing clusters into Planton. Once connected, you can deploy services to them through Service Hub, manage workloads through Cloud Ops, and include them in your environment authorization model alongside cloud provider credentials.

## When to Use Kubernetes Cluster Connections

- **Existing production clusters** — You have clusters already running and want to deploy services to them through Planton without recreating them.
- **Hybrid cloud** — Some clusters are managed by Planton, others are managed externally. Connecting external clusters gives you a unified deployment surface.
- **Migration** — You're moving to Planton gradually and want to deploy new services to existing clusters before migrating the cluster management itself.
- **Multi-cluster architectures** — Your workloads span multiple clusters, some of which Planton created and some of which it didn't.

## Supported Providers

Kubernetes cluster connections support four managed Kubernetes providers. Each requires different credentials depending on how the provider handles authentication.

### Google Kubernetes Engine (GKE)

The most complete implementation. GKE connections use a service account key to authenticate.

| Field | Description |
|-------|-------------|
| Cluster Endpoint | The API server endpoint of your GKE cluster |
| Cluster CA Data | Base64-encoded certificate authority data for TLS verification |
| Service Account Key | Base64-encoded JSON key for a GCP service account with Kubernetes Engine access |

To find these values:

1. In the Google Cloud Console, navigate to your GKE cluster's details page.
2. The **Endpoint** and **Cluster CA certificate** are on the cluster details page.
3. Create a GCP service account with the `Kubernetes Engine Developer` or `Kubernetes Engine Admin` role and generate a JSON key.

### DigitalOcean Kubernetes (DOKS)

DOKS connections use a kubeconfig file.

| Field | Description |
|-------|-------------|
| Kubeconfig | The kubeconfig content for your DOKS cluster |

To get your kubeconfig:

1. In the DigitalOcean control panel, navigate to your Kubernetes cluster.
2. Download the kubeconfig file from the cluster details page, or use the CLI: `doctl kubernetes cluster kubeconfig show my-cluster`.

### Amazon EKS

> **Note**: EKS cluster connection support is defined in the API but the credential configuration is not yet implemented. EKS clusters created by Planton through Infra Hub work automatically — this connection type is for externally managed EKS clusters. Check back for updates.

### Azure AKS

> **Note**: AKS cluster connection support is defined in the API but the credential configuration is not yet implemented. AKS clusters created by Planton through Infra Hub work automatically — this connection type is for externally managed AKS clusters. Check back for updates.

## Connecting via the Web Console

1. Navigate to **Connections** and click the **Kubernetes** card under Infrastructure.
2. **Name your connection** — use a name that identifies the cluster (e.g., "prod-gke-us-east", "legacy-doks-cluster").
3. **Select the provider** — GKE or DigitalOcean DOKS.
4. **Provide the credentials** listed above for your provider.
5. **Create the connection**.

<!-- SCREENSHOT: Kubernetes cluster connection form
  Page: /resource/connect/kubernetes-cluster-credential/create
  Action: Show the provider selector with GKE selected and credentials form visible
  Focus: Provider selection and GKE-specific credential fields
  Alt: Kubernetes cluster connection form showing GKE provider selected with endpoint, CA data, and service account key fields
-->

## Authentication Modes

Like cloud provider connections, Kubernetes cluster connections support inline and runner-delegated authentication:

- **Inline** — Provide the cluster credentials directly (endpoint, CA cert, service account key or kubeconfig). Simplest option.
- **Runner-delegated** — A [Planton Runner](/docs/runner) deployed with access to the cluster handles authentication. Useful when the cluster is in a private network and credentials should not leave the network perimeter.

## How Kubernetes Connections Are Used

Once connected, external clusters become deployment targets:

- **Service Hub** — When creating a service deployment target, you can select a connected external cluster alongside clusters that Planton created through Infra Hub.
- **Cloud Ops** — You can browse pods, stream logs, and exec into containers on connected clusters through the Cloud Ops interface, as long as a Runner with access to the cluster is configured.
- **Infra Hub** — Cloud resources with Kubernetes deployment components (Helm charts, operators, custom resources) can target connected clusters.

## Practical Guidance

### Cluster Naming

Name connections by cluster identity and purpose, not by how they were created:

- `prod-gke-us-central1` — identifies the cluster's role, provider, and region
- `staging-doks-nyc1` — clear and specific
- `legacy-app-cluster` — useful during migration

### Credential Scope

The credentials you provide should have the minimum permissions needed for your use case:

- **Service deployments only**: The service account needs permissions to create and manage Deployments, Services, ConfigMaps, Secrets, and Ingress resources in the target namespaces.
- **Cloud Ops access**: Additionally needs pod list, log read, and exec permissions.
- **Full Infra Hub management**: Needs cluster-admin or equivalent broad permissions.

### Keep Credentials Current

Kubernetes credentials (especially kubeconfig tokens and service account keys) have expiration policies. Monitor connection health and rotate credentials before they expire to avoid deployment failures.

## Related Documentation

- [Connections Overview](/docs/connections) — Understanding the Connect system
- [Cloud Providers](/docs/connections/cloud-providers) — Connect cloud provider accounts for creating new clusters
- [CI/CD: Deployment Targets](/docs/ci-cd/deployment-targets) — How services target Kubernetes clusters
- [Operations](/docs/operations) — Runtime operations on connected clusters
