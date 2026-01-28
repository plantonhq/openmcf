# Harbor Kubernetes - Pulumi Module Examples

This document provides examples of using the KubernetesHarbor Pulumi module directly in your Pulumi programs.

## Namespace Management

All Harbor deployments require a namespace. Control namespace management with the `CreateNamespace` flag:

- **`CreateNamespace: true`** - Component creates and manages the namespace with appropriate labels
- **`CreateNamespace: false`** - Deploy into an existing namespace (must exist before deployment)

**Managed namespace example:**
```go
Spec: &kubernetesharborv1.KubernetesHarborSpec{
    Namespace: &foreignkeyv1.StringValueOrRef{
        LiteralOrRef: &foreignkeyv1.StringValueOrRef_Value{
            Value: "harbor-prod",
        },
    },
    CreateNamespace: true,  // Component creates the namespace
    // ... rest of spec
}
```

**Existing namespace example:**
```go
Spec: &kubernetesharborv1.KubernetesHarborSpec{
    Namespace: &foreignkeyv1.StringValueOrRef{
        LiteralOrRef: &foreignkeyv1.StringValueOrRef_Value{
            Value: "container-registry",
        },
    },
    CreateNamespace: false,  // Must exist before deployment
    // ... rest of spec
}
```

## Prerequisites

- Pulumi CLI installed
- Go 1.21 or later
- Access to a Kubernetes cluster
- kubectl configured with cluster access

## Setup

Create a new Pulumi Go project:

```bash
mkdir my-harbor-deployment
cd my-harbor-deployment
pulumi new go
```

Add the required dependencies to your `go.mod`:

```go
require (
    github.com/plantonhq/openmcf v1.0.0
    github.com/pulumi/pulumi-kubernetes/sdk/v4 v4.0.0
    github.com/pulumi/pulumi/sdk/v3 v3.0.0
)
```

## Example 1: Basic Development Deployment

This example deploys Harbor with self-managed PostgreSQL and Redis, suitable for development and testing.

```go
package main

import (
    kubernetesharborv1 "github.com/plantonhq/openmcf/apis/org/openmcf/provider/kubernetes/kubernetesharbor/v1"
    "github.com/plantonhq/openmcf/apis/org/openmcf/provider/kubernetes/kubernetesharbor/v1/iac/pulumi/module"
    "github.com/plantonhq/openmcf/apis/org/openmcf/shared"
    "github.com/plantonhq/openmcf/apis/org/openmcf/provider/kubernetes"
    "github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
    pulumi.Run(func(ctx *pulumi.Context) error {
        // Define Harbor configuration
        harborSpec := &kubernetesharborv1.KubernetesHarborSpec{
            TargetCluster: &kubernetes.KubernetesClusterSelector{
                ClusterName: "my-gke-cluster",
            },
            Namespace: &foreignkeyv1.StringValueOrRef{
                LiteralOrRef: &foreignkeyv1.StringValueOrRef_Value{
                    Value: "harbor-dev",
                },
            },
            CreateNamespace: true,
            CoreContainer: &kubernetesharborv1.KubernetesHarborContainer{
                Replicas: 1,
                Resources: &kubernetes.ContainerResources{
                    Limits: &kubernetes.CpuMemory{
                        Cpu:    "1000m",
                        Memory: "2Gi",
                    },
                    Requests: &kubernetes.CpuMemory{
                        Cpu:    "200m",
                        Memory: "512Mi",
                    },
                },
            },
            PortalContainer: &kubernetesharborv1.KubernetesHarborContainer{
                Replicas: 1,
                Resources: &kubernetes.ContainerResources{
                    Limits: &kubernetes.CpuMemory{
                        Cpu:    "500m",
                        Memory: "512Mi",
                    },
                    Requests: &kubernetes.CpuMemory{
                        Cpu:    "100m",
                        Memory: "256Mi",
                    },
                },
            },
            RegistryContainer: &kubernetesharborv1.KubernetesHarborContainer{
                Replicas: 1,
                Resources: &kubernetes.ContainerResources{
                    Limits: &kubernetes.CpuMemory{
                        Cpu:    "1000m",
                        Memory: "2Gi",
                    },
                    Requests: &kubernetes.CpuMemory{
                        Cpu:    "200m",
                        Memory: "512Mi",
                    },
                },
            },
            JobserviceContainer: &kubernetesharborv1.KubernetesHarborContainer{
                Replicas: 1,
                Resources: &kubernetes.ContainerResources{
                    Limits: &kubernetes.CpuMemory{
                        Cpu:    "1000m",
                        Memory: "1Gi",
                    },
                    Requests: &kubernetes.CpuMemory{
                        Cpu:    "100m",
                        Memory: "256Mi",
                    },
                },
            },
            Database: &kubernetesharborv1.KubernetesHarborDatabaseConfig{
                IsExternal: false,
                ManagedDatabase: &kubernetesharborv1.KubernetesHarborManagedPostgresql{
                    Container: &kubernetesharborv1.KubernetesHarborPostgresqlContainer{
                        Replicas:           1,
                        PersistenceEnabled: true,
                        DiskSize:           "20Gi",
                        Resources: &kubernetes.ContainerResources{
                            Limits: &kubernetes.CpuMemory{
                                Cpu:    "1000m",
                                Memory: "2Gi",
                            },
                            Requests: &kubernetes.CpuMemory{
                                Cpu:    "200m",
                                Memory: "512Mi",
                            },
                        },
                    },
                },
            },
            Cache: &kubernetesharborv1.KubernetesHarborCacheConfig{
                IsExternal: false,
                ManagedCache: &kubernetesharborv1.KubernetesHarborManagedRedis{
                    Container: &kubernetesharborv1.KubernetesHarborRedisContainer{
                        Replicas:           1,
                        PersistenceEnabled: true,
                        DiskSize:           "8Gi",
                        Resources: &kubernetes.ContainerResources{
                            Limits: &kubernetes.CpuMemory{
                                Cpu:    "500m",
                                Memory: "512Mi",
                            },
                            Requests: &kubernetes.CpuMemory{
                                Cpu:    "100m",
                                Memory: "256Mi",
                            },
                        },
                    },
                },
            },
            Storage: &kubernetesharborv1.KubernetesHarborStorageConfig{
                Type: kubernetesharborv1.KubernetesHarborStorageType_filesystem,
                Filesystem: &kubernetesharborv1.KubernetesHarborFilesystemStorage{
                    DiskSize: "100Gi",
                },
            },
        }

        // Create Harbor resource
        harbor := &kubernetesharborv1.KubernetesHarbor{
            ApiVersion: "kubernetes.openmcf.org/v1",
            Kind:       "KubernetesHarbor",
            Metadata: &shared.CloudResourceMetadata{
                Name: "dev-harbor",
            },
            Spec: harborSpec,
        }

        // Create stack input
        stackInput := &kubernetesharborv1.KubernetesHarborStackInput{
            Target: harbor,
            // ProviderConfig would be set here in production
        }

        // Deploy Harbor
        if err := module.Resources(ctx, stackInput); err != nil {
            return err
        }

        return nil
    })
}
```

## Example 2: Production HA with AWS S3

This example shows a production-ready high-availability deployment using external AWS services.

```go
package main

import (
    kubernetesharborv1 "github.com/plantonhq/openmcf/apis/org/openmcf/provider/kubernetes/kubernetesharbor/v1"
    "github.com/plantonhq/openmcf/apis/org/openmcf/provider/kubernetes/kubernetesharbor/v1/iac/pulumi/module"
    "github.com/plantonhq/openmcf/apis/org/openmcf/shared"
    "github.com/plantonhq/openmcf/apis/org/openmcf/provider/kubernetes"
    foreignkeyv1 "github.com/plantonhq/openmcf/apis/org/openmcf/shared/foreignkey/v1"
    "github.com/plantonhq/openmcf/apis/org/openmcf/shared/cloudresourcekind"
    "github.com/pulumi/pulumi/sdk/v3/go/pulumi"
    "github.com/pulumi/pulumi/sdk/v3/go/pulumi/config"
)

func main() {
    pulumi.Run(func(ctx *pulumi.Context) error {
        cfg := config.New(ctx, "")

        // Get secrets from Pulumi config
        dbPassword := cfg.RequireSecret("db-password")
        redisPassword := cfg.RequireSecret("redis-password")
        awsAccessKey := cfg.RequireSecret("aws-access-key")
        awsSecretKey := cfg.RequireSecret("aws-secret-key")

        // Get configuration values
        dbHost := cfg.Require("db-host")
        redisHost := cfg.Require("redis-host")
        s3Bucket := cfg.Require("s3-bucket")
        s3Region := cfg.Require("s3-region")

        harborSpec := &kubernetesharborv1.KubernetesHarborSpec{
            TargetCluster: &kubernetes.KubernetesClusterSelector{
                ClusterName: "my-eks-cluster",
            },
            Namespace: &foreignkeyv1.StringValueOrRef{
                LiteralOrRef: &foreignkeyv1.StringValueOrRef_Value{
                    Value: "harbor-prod",
                },
            },
            CreateNamespace: true,
            // Harbor Core - 3 replicas for HA
            CoreContainer: &kubernetesharborv1.KubernetesHarborContainer{
                Replicas: 3,
                Resources: &kubernetes.ContainerResources{
                    Limits: &kubernetes.CpuMemory{
                        Cpu:    "2000m",
                        Memory: "4Gi",
                    },
                    Requests: &kubernetes.CpuMemory{
                        Cpu:    "500m",
                        Memory: "1Gi",
                    },
                },
            },
            // Harbor Portal - 2 replicas
            PortalContainer: &kubernetesharborv1.KubernetesHarborContainer{
                Replicas: 2,
                Resources: &kubernetes.ContainerResources{
                    Limits: &kubernetes.CpuMemory{
                        Cpu:    "1000m",
                        Memory: "1Gi",
                    },
                    Requests: &kubernetes.CpuMemory{
                        Cpu:    "200m",
                        Memory: "512Mi",
                    },
                },
            },
            // Harbor Registry - 3 replicas for high throughput
            RegistryContainer: &kubernetesharborv1.KubernetesHarborContainer{
                Replicas: 3,
                Resources: &kubernetes.ContainerResources{
                    Limits: &kubernetes.CpuMemory{
                        Cpu:    "2000m",
                        Memory: "4Gi",
                    },
                    Requests: &kubernetes.CpuMemory{
                        Cpu:    "500m",
                        Memory: "1Gi",
                    },
                },
            },
            // Harbor Jobservice - 2 replicas for parallel jobs
            JobserviceContainer: &kubernetesharborv1.KubernetesHarborContainer{
                Replicas: 2,
                Resources: &kubernetes.ContainerResources{
                    Limits: &kubernetes.CpuMemory{
                        Cpu:    "2000m",
                        Memory: "2Gi",
                    },
                    Requests: &kubernetes.CpuMemory{
                        Cpu:    "200m",
                        Memory: "512Mi",
                    },
                },
            },
            // External PostgreSQL (AWS RDS)
            Database: &kubernetesharborv1.KubernetesHarborDatabaseConfig{
                IsExternal: true,
                ExternalDatabase: &kubernetesharborv1.KubernetesHarborExternalPostgresql{
                    Host:     dbHost,
                    Port:     pulumi.Int32Ptr(5432),
                    Username: "harbor",
                    Password: dbPassword,
                    UseSSL:   true,
                },
            },
            // External Redis (AWS ElastiCache)
            Cache: &kubernetesharborv1.KubernetesHarborCacheConfig{
                IsExternal: true,
                ExternalCache: &kubernetesharborv1.KubernetesHarborExternalRedis{
                    Host:     redisHost,
                    Port:     pulumi.Int32Ptr(6379),
                    Password: redisPassword,
                },
            },
            // S3 Storage
            Storage: &kubernetesharborv1.KubernetesHarborStorageConfig{
                Type: kubernetesharborv1.KubernetesHarborStorageType_s3,
                S3: &kubernetesharborv1.KubernetesHarborS3Storage{
                    Bucket:         s3Bucket,
                    Region:         s3Region,
                    AccessKey:      awsAccessKey,
                    SecretKey:      awsSecretKey,
                    RegionEndpoint: false,
                    Encrypt:        true,
                    Secure:         true,
                },
            },
            // Ingress for external access
            Ingress: &kubernetesharborv1.KubernetesHarborIngress{
                Core: &kubernetesharborv1.KubernetesHarborIngressEndpoint{
                    Enabled:  true,
                    Hostname: "harbor.example.com",
                },
                Notary: &kubernetesharborv1.KubernetesHarborIngressEndpoint{
                    Enabled:  true,
                    Hostname: "notary.harbor.example.com",
                },
            },
        }

        harbor := &kubernetesharborv1.KubernetesHarbor{
            ApiVersion: "kubernetes.openmcf.org/v1",
            Kind:       "KubernetesHarbor",
            Metadata: &shared.CloudResourceMetadata{
                Name: "prod-harbor",
                Org:  "my-company",
                Env:  "production",
            },
            Spec: harborSpec,
        }

        stackInput := &kubernetesharborv1.KubernetesHarborStackInput{
            Target: harbor,
        }

        if err := module.Resources(ctx, stackInput); err != nil {
            return err
        }

        return nil
    })
}
```

Run with:

```bash
pulumi config set db-host harbor-db.xxxx.us-west-2.rds.amazonaws.com
pulumi config set redis-host harbor-redis.xxxx.use1.cache.amazonaws.com
pulumi config set s3-bucket my-harbor-artifacts
pulumi config set s3-region us-west-2
pulumi config set --secret db-password <your-db-password>
pulumi config set --secret redis-password <your-redis-password>
pulumi config set --secret aws-access-key <your-aws-access-key>
pulumi config set --secret aws-secret-key <your-aws-secret-key>
pulumi up
```

## Example 3: GCP with Cloud SQL and GCS

This example demonstrates deployment on Google Cloud Platform.

```go
package main

import (
    kubernetesharborv1 "github.com/plantonhq/openmcf/apis/org/openmcf/provider/kubernetes/kubernetesharbor/v1"
    "github.com/plantonhq/openmcf/apis/org/openmcf/provider/kubernetes/kubernetesharbor/v1/iac/pulumi/module"
    "github.com/plantonhq/openmcf/apis/org/openmcf/shared"
    "github.com/plantonhq/openmcf/apis/org/openmcf/provider/kubernetes"
    foreignkeyv1 "github.com/plantonhq/openmcf/apis/org/openmcf/shared/foreignkey/v1"
    "github.com/plantonhq/openmcf/apis/org/openmcf/shared/cloudresourcekind"
    "github.com/pulumi/pulumi/sdk/v3/go/pulumi"
    "github.com/pulumi/pulumi/sdk/v3/go/pulumi/config"
)

func main() {
    pulumi.Run(func(ctx *pulumi.Context) error {
        cfg := config.New(ctx, "")

        gcsSA := cfg.RequireSecret("gcs-service-account-key")
        dbPassword := cfg.RequireSecret("db-password")
        redisPassword := cfg.RequireSecret("redis-password")

        harborSpec := &kubernetesharborv1.KubernetesHarborSpec{
            TargetCluster: &kubernetes.KubernetesClusterSelector{
                ClusterName: "my-gke-cluster",
            },
            Namespace: &foreignkeyv1.StringValueOrRef{
                LiteralOrRef: &foreignkeyv1.StringValueOrRef_Value{
                    Value: "harbor-gcp",
                },
            },
            CreateNamespace: true,
            CoreContainer: &kubernetesharborv1.KubernetesHarborContainer{
                Replicas: 2,
                Resources: &kubernetes.ContainerResources{
                    Limits: &kubernetes.CpuMemory{
                        Cpu:    "2000m",
                        Memory: "4Gi",
                    },
                    Requests: &kubernetes.CpuMemory{
                        Cpu:    "500m",
                        Memory: "1Gi",
                    },
                },
            },
            PortalContainer: &kubernetesharborv1.KubernetesHarborContainer{
                Replicas: 2,
            },
            RegistryContainer: &kubernetesharborv1.KubernetesHarborContainer{
                Replicas: 2,
                Resources: &kubernetes.ContainerResources{
                    Limits: &kubernetes.CpuMemory{
                        Cpu:    "2000m",
                        Memory: "4Gi",
                    },
                    Requests: &kubernetes.CpuMemory{
                        Cpu:    "500m",
                        Memory: "1Gi",
                    },
                },
            },
            JobserviceContainer: &kubernetesharborv1.KubernetesHarborContainer{
                Replicas: 2,
            },
            // Cloud SQL PostgreSQL
            Database: &kubernetesharborv1.KubernetesHarborDatabaseConfig{
                IsExternal: true,
                ExternalDatabase: &kubernetesharborv1.KubernetesHarborExternalPostgresql{
                    Host:     cfg.Require("cloudsql-host"),
                    Port:     pulumi.Int32Ptr(5432),
                    Username: "harbor",
                    Password: dbPassword,
                    UseSSL:   true,
                },
            },
            // Cloud Memorystore Redis
            Cache: &kubernetesharborv1.KubernetesHarborCacheConfig{
                IsExternal: true,
                ExternalCache: &kubernetesharborv1.KubernetesHarborExternalRedis{
                    Host:     cfg.Require("memorystore-host"),
                    Port:     pulumi.Int32Ptr(6379),
                    Password: redisPassword,
                },
            },
            // Google Cloud Storage
            Storage: &kubernetesharborv1.KubernetesHarborStorageConfig{
                Type: kubernetesharborv1.KubernetesHarborStorageType_gcs,
                Gcs: &kubernetesharborv1.KubernetesHarborGcsStorage{
                    Bucket:  cfg.Require("gcs-bucket"),
                    KeyData: gcsSA,
                },
            },
            Ingress: &kubernetesharborv1.KubernetesHarborIngress{
                Core: &kubernetesharborv1.KubernetesHarborIngressEndpoint{
                    Enabled:  true,
                    Hostname: "harbor.example.com",
                },
            },
        }

        harbor := &kubernetesharborv1.KubernetesHarbor{
            ApiVersion: "kubernetes.openmcf.org/v1",
            Kind:       "KubernetesHarbor",
            Metadata: &shared.CloudResourceMetadata{
                Name: "gcp-harbor",
            },
            Spec: harborSpec,
        }

        stackInput := &kubernetesharborv1.KubernetesHarborStackInput{
            Target: harbor,
        }

        if err := module.Resources(ctx, stackInput); err != nil {
            return err
        }

        return nil
    })
}
```

## Example 4: Advanced Configuration with Trivy Scanner

This example shows how to enable additional Harbor features like Trivy vulnerability scanner using helm_values.

```go
package main

import (
    kubernetesharborv1 "github.com/plantonhq/openmcf/apis/org/openmcf/provider/kubernetes/kubernetesharbor/v1"
    "github.com/plantonhq/openmcf/apis/org/openmcf/provider/kubernetes/kubernetesharbor/v1/iac/pulumi/module"
    "github.com/plantonhq/openmcf/apis/org/openmcf/shared"
    "github.com/plantonhq/openmcf/apis/org/openmcf/provider/kubernetes"
    foreignkeyv1 "github.com/plantonhq/openmcf/apis/org/openmcf/shared/foreignkey/v1"
    "github.com/plantonhq/openmcf/apis/org/openmcf/shared/cloudresourcekind"
    "github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
    pulumi.Run(func(ctx *pulumi.Context) error {
        harborSpec := &kubernetesharborv1.KubernetesHarborSpec{
            TargetCluster: &kubernetes.KubernetesClusterSelector{
                ClusterName: "my-gke-cluster",
            },
            Namespace: &foreignkeyv1.StringValueOrRef{
                LiteralOrRef: &foreignkeyv1.StringValueOrRef_Value{
                    Value: "harbor-advanced",
                },
            },
            CreateNamespace: true,
            CoreContainer: &kubernetesharborv1.KubernetesHarborContainer{
                Replicas: 2,
            },
            PortalContainer: &kubernetesharborv1.KubernetesHarborContainer{
                Replicas: 2,
            },
            RegistryContainer: &kubernetesharborv1.KubernetesHarborContainer{
                Replicas: 2,
            },
            JobserviceContainer: &kubernetesharborv1.KubernetesHarborContainer{
                Replicas: 2,
            },
            Database: &kubernetesharborv1.KubernetesHarborDatabaseConfig{
                IsExternal: false,
                ManagedDatabase: &kubernetesharborv1.KubernetesHarborManagedPostgresql{
                    Container: &kubernetesharborv1.KubernetesHarborPostgresqlContainer{
                        Replicas:           1,
                        PersistenceEnabled: true,
                        DiskSize:           "20Gi",
                    },
                },
            },
            Cache: &kubernetesharborv1.KubernetesHarborCacheConfig{
                IsExternal: false,
                ManagedCache: &kubernetesharborv1.KubernetesHarborManagedRedis{
                    Container: &kubernetesharborv1.KubernetesHarborRedisContainer{
                        Replicas:           1,
                        PersistenceEnabled: true,
                        DiskSize:           "8Gi",
                    },
                },
            },
            Storage: &kubernetesharborv1.KubernetesHarborStorageConfig{
                Type: kubernetesharborv1.KubernetesHarborStorageType_s3,
                S3: &kubernetesharborv1.KubernetesHarborS3Storage{
                    Bucket:    "my-harbor-bucket",
                    Region:    "us-west-2",
                    AccessKey: "...",
                    SecretKey: "...",
                },
            },
            // Enable Trivy scanner and other features via helm_values
            HelmValues: map[string]string{
                "trivy.enabled":           "true",
                "notary.enabled":          "true",
                "chartmuseum.enabled":     "false",
                "metrics.enabled":         "true",
                "metrics.serviceMonitor": "true",
            },
        }

        harbor := &kubernetesharborv1.KubernetesHarbor{
            ApiVersion: "kubernetes.openmcf.org/v1",
            Kind:       "KubernetesHarbor",
            Metadata: &shared.CloudResourceMetadata{
                Name: "advanced-harbor",
            },
            Spec: harborSpec,
        }

        stackInput := &kubernetesharborv1.KubernetesHarborStackInput{
            Target: harbor,
        }

        if err := module.Resources(ctx, stackInput); err != nil {
            return err
        }

        return nil
    })
}
```

## Example 5: MinIO S3-Compatible Storage

This example shows using MinIO as an S3-compatible storage backend.

```go
package main

import (
    kubernetesharborv1 "github.com/plantonhq/openmcf/apis/org/openmcf/provider/kubernetes/kubernetesharbor/v1"
    "github.com/plantonhq/openmcf/apis/org/openmcf/provider/kubernetes/kubernetesharbor/v1/iac/pulumi/module"
    "github.com/plantonhq/openmcf/apis/org/openmcf/shared"
    "github.com/plantonhq/openmcf/apis/org/openmcf/provider/kubernetes"
    foreignkeyv1 "github.com/plantonhq/openmcf/apis/org/openmcf/shared/foreignkey/v1"
    "github.com/plantonhq/openmcf/apis/org/openmcf/shared/cloudresourcekind"
    "github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
    pulumi.Run(func(ctx *pulumi.Context) error {
        harborSpec := &kubernetesharborv1.KubernetesHarborSpec{
            TargetCluster: &kubernetes.KubernetesClusterSelector{
                ClusterName: "my-gke-cluster",
            },
            Namespace: &foreignkeyv1.StringValueOrRef{
                LiteralOrRef: &foreignkeyv1.StringValueOrRef_Value{
                    Value: "harbor-minio",
                },
            },
            CreateNamespace: true,
            CoreContainer:       &kubernetesharborv1.KubernetesHarborContainer{Replicas: 1},
            PortalContainer:     &kubernetesharborv1.KubernetesHarborContainer{Replicas: 1},
            RegistryContainer:   &kubernetesharborv1.KubernetesHarborContainer{Replicas: 1},
            JobserviceContainer: &kubernetesharborv1.KubernetesHarborContainer{Replicas: 1},
            Database: &kubernetesharborv1.KubernetesHarborDatabaseConfig{
                IsExternal: false,
                ManagedDatabase: &kubernetesharborv1.KubernetesHarborManagedPostgresql{
                    Container: &kubernetesharborv1.KubernetesHarborPostgresqlContainer{
                        Replicas:           1,
                        PersistenceEnabled: true,
                        DiskSize:           "20Gi",
                    },
                },
            },
            Cache: &kubernetesharborv1.KubernetesHarborCacheConfig{
                IsExternal: false,
                ManagedCache: &kubernetesharborv1.KubernetesHarborManagedRedis{
                    Container: &kubernetesharborv1.KubernetesHarborRedisContainer{
                        Replicas:           1,
                        PersistenceEnabled: true,
                        DiskSize:           "8Gi",
                    },
                },
            },
            // MinIO S3-compatible storage
            Storage: &kubernetesharborv1.KubernetesHarborStorageConfig{
                Type: kubernetesharborv1.KubernetesHarborStorageType_s3,
                S3: &kubernetesharborv1.KubernetesHarborS3Storage{
                    Bucket:      "harbor",
                    Region:      "us-east-1", // MinIO doesn't use regions but field is required
                    AccessKey:   "minio-access-key",
                    SecretKey:   "minio-secret-key",
                    EndpointUrl: "http://minio.minio-system.svc.cluster.local:9000",
                    Secure:      false, // Use true if MinIO has TLS
                },
            },
        }

        harbor := &kubernetesharborv1.KubernetesHarbor{
            ApiVersion: "kubernetes.openmcf.org/v1",
            Kind:       "KubernetesHarbor",
            Metadata: &shared.CloudResourceMetadata{
                Name: "minio-harbor",
            },
            Spec: harborSpec,
        }

        stackInput := &kubernetesharborv1.KubernetesHarborStackInput{
            Target: harbor,
        }

        if err := module.Resources(ctx, stackInput); err != nil {
            return err
        }

        return nil
    })
}
```

## Stack Outputs

All examples export the following outputs that you can access after deployment:

```go
// Access outputs in your Pulumi program
ctx.Export("harbor-namespace", pulumi.String("..."))
ctx.Export("harbor-core-service", pulumi.String("..."))
ctx.Export("harbor-portal-service", pulumi.String("..."))
ctx.Export("harbor-registry-service", pulumi.String("..."))
ctx.Export("harbor-jobservice-service", pulumi.String("..."))
ctx.Export("internal-core-endpoint", pulumi.String("..."))
ctx.Export("internal-registry-endpoint", pulumi.String("..."))
ctx.Export("port-forward-command", pulumi.String("..."))
ctx.Export("external-hostname", pulumi.String("..."))
```

View outputs:

```bash
pulumi stack output harbor-namespace
pulumi stack output external-hostname
pulumi stack output port-forward-command
```

## Accessing Harbor

### Via Port-Forward (Development)

```bash
# Get the port-forward command from stack outputs
kubectl port-forward -n <namespace> svc/<core-service-name> 8080:80

# Open Harbor UI
open http://localhost:8080

# Default credentials
Username: admin
Password: Harbor12345
```

### Via Ingress (Production)

If ingress is enabled, access Harbor at:

```
https://harbor.example.com
```

## Best Practices

1. **Use Pulumi Secrets for Sensitive Data**
   ```bash
   pulumi config set --secret db-password <password>
   pulumi config set --secret aws-secret-key <key>
   ```

2. **Separate Stacks for Environments**
   ```bash
   pulumi stack init dev
   pulumi stack init staging
   pulumi stack init prod
   ```

3. **Use Config Files**
   Create `Pulumi.dev.yaml`, `Pulumi.prod.yaml` etc.:
   ```yaml
   config:
     harbor:db-host: dev-db.example.com
     harbor:s3-bucket: dev-harbor-artifacts
   ```

4. **Version Control**
   - Commit Pulumi code to git
   - Use `.gitignore` for `Pulumi.*.yaml` files containing secrets
   - Use Pulumi ESC for centralized secrets management

## Troubleshooting

### Check Deployment Status

```go
// Add to your Pulumi program
ctx.Export("deployment-status", pulumi.String("deployed"))
```

### View Resources

```bash
pulumi stack
pulumi stack output
kubectl get all -n <namespace>
```

### Debug Helm Values

```bash
helm get values <release-name> -n <namespace>
```

## Additional Resources

- [Pulumi Documentation](https://www.pulumi.com/docs/)
- [Harbor Documentation](https://goharbor.io/docs/)
- [Pulumi Kubernetes Provider](https://www.pulumi.com/registry/packages/kubernetes/)

