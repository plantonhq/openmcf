---
title: "Cloud Resource Kinds"
description: "The taxonomy of components in Planton: 360+ resource kinds across 17 cloud providers, each with a unique kind name, provider mapping, and API version"
icon: "database"
order: 30
---

# Cloud Resource Kinds

Every component in Planton has a unique kind name -- `AwsS3Bucket`, `KubernetesPostgres`, `GcpCloudSql`. These kind names are not arbitrary strings. They are entries in the `CloudResourceKind` enum, a Protocol Buffer enum that serves as the canonical registry of everything Planton can deploy.

<!-- AI-AGENT NOTE: The component count below MUST be verified against the authoritative source:
     shared/cloudresourcekind/cloud_resource_kind.proto
     Count only non-test, non-unspecified enum values in the CloudResourceKind enum. -->

The enum currently contains 360+ resource kinds spanning 17 cloud providers.

## The CloudResourceKind Enum

The `CloudResourceKind` enum is defined in `cloud_resource_kind.proto`. Each entry carries metadata that maps the kind to its provider, API version, and an ID prefix:

```protobuf
KubernetesPostgres = 814 [(kind_meta) = {
    provider: kubernetes
    version: v1
    id_prefix: "k8spg"
}];

AwsS3Bucket = 213 [(kind_meta) = {
    provider: aws
    version: v1
    id_prefix: "s3bkt"
}];
```

The `kind_meta` annotation tells the CLI everything it needs: which provider this kind belongs to (determines the IaC module path and provider config type), which API version to use (determines the `apiVersion` field value), and a short ID prefix for resource identification.

## The CloudResourceProvider Enum

Each provider is registered in a separate `CloudResourceProvider` enum with a group name that forms the `apiVersion` domain:

```protobuf
enum CloudResourceProvider {
    aws = 12 [(provider_meta) = {
        group: "aws.planton.dev"
        display_name: "AWS"
    }];
    kubernetes = 19 [(provider_meta) = {
        group: "kubernetes.planton.dev"
        display_name: "Kubernetes"
    }];
}
```

The `group` value directly maps to the `apiVersion` in manifests. A resource with `provider: aws` uses `apiVersion: aws.planton.dev/v1alpha1`. A resource with `provider: kubernetes` uses `apiVersion: kubernetes.planton.dev/v1alpha1`.

## Provider Breakdown

| Provider | Components | apiVersion Domain | Example Kinds |
|----------|-----------|-------------------|---------------|
| **Kubernetes** | 51 | `kubernetes.planton.dev` | KubernetesPostgres, KubernetesRedis, KubernetesDeployment, KubernetesHelmRelease |
| **OpenStack** | 27 | `openstack.planton.dev` | OpenStackInstance, OpenStackNetwork, OpenStackLoadBalancer, OpenStackVolume |
| **AWS** | 25 | `aws.planton.dev` | AwsS3Bucket, AwsEksCluster, AwsRdsInstance, AwsLambda, AwsVpc |
| **GCP** | 19 | `gcp.planton.dev` | GcpCloudSql, GcpGkeCluster, GcpGcsBucket, GcpCloudRun, GcpVpcNetwork |
| **Scaleway** | 19 | `scaleway.planton.dev` | ScalewayInstance, ScalewayKapsuleCluster, ScalewayRdbInstance, ScalewayVpc |
| **DigitalOcean** | 15 | `digital-ocean.planton.dev` | DigitalOceanDroplet, DigitalOceanKubernetesCluster, DigitalOceanDatabaseCluster |
| **Azure** | 10 | `azure.planton.dev` | AzureAksCluster, AzureKeyVault, AzureStorageAccount, AzureVpc |
| **Civo** | 12 | `civo.planton.dev` | CivoKubernetesCluster, CivoDatabase, CivoComputeInstance, CivoVpc |
| **Cloudflare** | 8 | `cloudflare.planton.dev` | CloudflareDnsZone, CloudflareWorker, CloudflareR2Bucket, CloudflareD1Database |
| **Auth0** | 4 | `auth0.planton.dev` | Auth0Client, Auth0Connection, Auth0EventStream, Auth0ResourceServer |
| **OpenFGA** | 3 | `openfga.planton.dev` | OpenFgaStore, OpenFgaAuthorizationModel, OpenFgaRelationshipTuple |
| **Confluent** | 1 | `confluent.planton.dev` | ConfluentKafka |
| **MongoDB Atlas** | 1 | `atlas.planton.dev` | MongodbAtlas |
| **Snowflake** | 1 | `snowflake.planton.dev` | SnowflakeDatabase |

## Naming Convention

Kind names follow a consistent pattern: `{Provider}{Resource}`.

- The provider prefix identifies the cloud platform: `Aws`, `Gcp`, `Azure`, `Kubernetes`, `DigitalOcean`, `Civo`, `Cloudflare`, `OpenStack`, `Scaleway`, `Auth0`, `OpenFga`
- The resource suffix describes what it deploys: `S3Bucket`, `Postgres`, `CloudSql`, `EksCluster`, `Vpc`

This convention eliminates ambiguity. When you see `GcpCloudSql` in a manifest, you know immediately that this is a Google Cloud SQL resource managed through the GCP provider, not a generic SQL database abstraction.

## Enum Band Allocation

Every cloud provider owns a 1,000-wide number band, giving each catalog room to grow to its provider's full surface without ever colliding with a neighbor:

| Band | Provider |
|-------|----------|
| 1-49 | Test/development |
| 50-199 | Third-party services (Confluent, Atlas, Snowflake) |
| 1000-1999 | AWS |
| 2000-2999 | Azure |
| 3000-3999 | GCP |
| 4000-4999 | Kubernetes |
| 5000-5999 | DigitalOcean |
| 6000-6999 | Civo |
| 7000-7999 | Cloudflare |
| 8000-8999 | Auth0 |
| 9000-9999 | OpenFGA |
| 10000-10999 | OpenStack |
| 11000-11999 | Scaleway |
| 12000-12999 | Alibaba Cloud |
| 13000-13999 | OCI |
| 14000-14999 | Hetzner Cloud |

New resources for an existing provider are added within its band. The next new provider takes the next free 1,000-wide band.

## From Kind to Deployment

The kind name is the key that unlocks the entire deployment pipeline. When you run:

```bash
planton pulumi up -f my-resource.yaml
```

The CLI reads the `kind` field from your manifest and uses the `CloudResourceKind` enum to:

1. **Resolve the provider** -- determines which `ProviderConfig` type to use for credentials
2. **Locate the IaC module** -- maps to `catalog/{provider}/{kind}/{version}/iac/pulumi/` or `iac/tf/`, where `{version}` is the kind's current API version from the registry (for example `v1alpha1`)
3. **Load the protobuf schema** -- determines which message type to use for validation
4. **Construct the stack input** -- wraps your manifest and provider config into the IaC input contract

This is why the kind name must exactly match the enum entry. `kubernetespostgres` will not work. `Kubernetes-Postgres` will not work. It must be `KubernetesPostgres`, matching the protobuf const validation:

```protobuf
string kind = 2 [(buf.validate.field).string.const = 'KubernetesPostgres'];
```

## Browsing Available Components

The [Component Catalog](/docs/catalog) provides detailed documentation for every component, organized by provider. Each catalog page includes the component's configuration fields, deployment behavior, and usage examples.

## What's Next

- **[Components](components)** -- The anatomy of what each kind maps to
- **[Manifests](manifests)** -- How to write manifests using these kind names
- **[Validation](validation)** -- How kind and apiVersion values are validated
- **[Component Catalog](/docs/catalog)** -- Browse documentation for all 360+ components
