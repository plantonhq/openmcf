# App on DigitalOcean App Platform

Deploys a full App Platform application — HTTP services, workers, jobs, static sites, functions, and in-app databases — from Git or from a container image. Integrates with Planton's Provider Connections for DigitalOcean credential management and ValueFromRef for VPC, DNS zone, and database-cluster wiring.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **App Platform Application** -- a `digitalocean_app` resource whose spec lists every component you declared
- **HTTP services** -- created from `spec.services`; receive traffic with automatic HTTPS
- **Workers** -- created from `spec.workers`; long-running processes without HTTP ingress
- **Jobs** -- created from `spec.jobs`; run around a deployment (`pre_deploy`, `post_deploy`, or `failed_deploy`)
- **Static sites, functions, in-app databases** -- created from the matching spec lists when present
- **Domains and ingress** -- created from `spec.domains` and `spec.ingress` when present
- **Environment variables** -- from `spec.envs` and per-component `envs`; `secret` values are stored in App Platform's secret store

## Before You Deploy

### Planton Setup

- **DigitalOcean Provider Connection** -- an active connection in the Connect module with a DigitalOcean API token. Map it as the default for your environment, or specify it explicitly when creating the Cloud Resource.

### DigitalOcean Account

- **A source** for each component: a public Git clone URL, a linked GitHub/GitLab/Bitbucket repo, or a container image (Docker Hub, GHCR, or DigitalOcean Container Registry).
- **A supported App Platform region** (for example `nyc3`, `sfo3`, `fra1`).

## Deploy

### Console

Open the deployment store, find **App on DigitalOcean App Platform**, and click **Deploy**. Start from the **Git Source Web App** preset to build from a repository, or **Container Image App** to run a pre-built image.

### CLI

```yaml
apiVersion: digital-ocean.planton.dev/v1alpha1
kind: DigitalOceanApp
metadata:
  name: my-web-app
  org: acme-corp
  env: prod
spec:
  appName: my-web-app
  region: nyc3
  services:
    - name: web
      git:
        repoCloneUrl: https://github.com/digitalocean/sample-nodejs.git
        branch: main
      instanceSizeSlug: basic-xxs
      instanceCount: 1
```

```shell
planton apply -f app.yaml
```

This creates a single-instance web app in NYC3 built from the sample Node.js repository.

## Key Configuration

**App name** -- `appName` is 2–32 characters and unique in the DigitalOcean account.

**Components** -- add at least one of `services`, `workers`, `jobs`, `staticSites`, or `functions`. An empty app cannot deploy.

**Source** -- each component sets exactly one of `git`, `github`, `gitlab`, `bitbucket`, or (services/workers/jobs) `image`. Use `git` with a public clone URL when the account has no VCS connection.

**Instance size** -- `instanceSizeSlug` is a free-form string (`basic-xxs`, `professional-s`, …). New sizes work without a catalog change.

**Autoscaling** -- set `autoscaling` on a service or worker and leave `instanceCount` unset. Autoscaling is not available on jobs.

## Outputs and Dependencies

### What This Component Consumes

| Dependency | Field | ValueFromRef Path |
|------------|-------|-------------------|
| **DigitalOceanVpc** (optional) | `vpc` | `status.outputs.vpc_id` |
| **DigitalOceanDnsZone** (optional) | `domains[].zone` | `status.outputs.zone_name` |
| **DigitalOceanDatabaseCluster** (optional) | `databases[].clusterName` | `spec.cluster_name` |

VPC placement is wired by Terraform. The Pulumi SDK at v4.49.0 cannot set it; Pulumi fails loudly if `vpc` is set.

### What This Component Provides

| Output | Description | Common Downstream Use |
|--------|-------------|----------------------|
| `app_id` | App Platform application UUID | Import, API operations |
| `default_hostname` | Default `ondigitalocean.app` hostname | DNS CNAME targets |
| `live_url` | Public URL including protocol | Application integrations |
| `live_domain` | Live hostname without scheme | Certificates, DNS |
| `active_deployment_id` | Currently live deployment UUID | Deployment tracking |

## Common Patterns

**Git source web app** -- one HTTP service built from a public repository. Start from the **Git Source Web App** preset.

**Container image app** -- one HTTP service from Docker Hub, GHCR, or DOCR, with optional autoscaling. Start from the **Container Image App** preset.

## Works With

- [**VPC on DigitalOcean**](/cloud-catalog/digital-ocean-vpc) -- optional egress placement (`spec.vpc`)
- [**DNS Zone on DigitalOcean**](/cloud-catalog/digital-ocean-dns-zone) -- custom domains (`spec.domains`)
- [**Database Cluster on DigitalOcean**](/cloud-catalog/digital-ocean-database-cluster) -- in-app database `clusterName`
- [**Function on DigitalOcean**](/cloud-catalog/digital-ocean-function) -- standalone functions app when the functions component should not share this app
