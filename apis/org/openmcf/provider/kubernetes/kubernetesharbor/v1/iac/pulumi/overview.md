# KubernetesHarbor Pulumi Module Overview

## Architecture

The KubernetesHarbor Pulumi module deploys Harbor, an open-source cloud-native registry, on a Kubernetes cluster. The module provides a highly flexible deployment architecture that supports both simple development setups and production-grade high-availability configurations.

## Component Structure

The module is organized into the following files:

### Core Files

- **`main.go`** - Entry point that orchestrates the deployment sequence
- **`locals.go`** - Initializes local variables, service names, endpoints, and stack outputs
- **`outputs.go`** - Defines output constant names for stack exports
- **`variables.go`** - Contains deployment constants (ports, Helm chart info)

### Resource Files

- **`harbor.go`** - Main Harbor deployment via Helm chart with comprehensive configuration
- **`ingress_core.go`** - Gateway API ingress for Harbor Core/Portal (web UI and API)
- **`ingress_notary.go`** - Gateway API ingress for Notary service (image signing)

## Architecture Diagram

```
┌─────────────────────────────────────────────────────────────────┐
│                     Kubernetes Cluster                          │
│                                                                 │
│  ┌───────────────────────────────────────────────────────────┐ │
│  │              Namespace: {metadata.name}                   │ │
│  │                                                           │ │
│  │  ┌─────────────┐  ┌─────────────┐  ┌─────────────────┐  │ │
│  │  │   Harbor    │  │   Harbor    │  │     Harbor      │  │ │
│  │  │    Core     │  │   Portal    │  │   Registry      │  │ │
│  │  │   (API)     │  │  (Web UI)   │  │ (Docker/OCI)    │  │ │
│  │  │             │  │             │  │                 │  │ │
│  │  │  Replicas:  │  │  Replicas:  │  │   Replicas:     │  │ │
│  │  │     N       │  │     N       │  │      N          │  │ │
│  │  └─────────────┘  └─────────────┘  └─────────────────┘  │ │
│  │         │                 │                  │           │ │
│  │         └─────────────────┴──────────────────┘           │ │
│  │                           │                               │ │
│  │                   ┌───────▼──────┐                        │ │
│  │                   │   Harbor     │                        │ │
│  │                   │  Jobservice  │                        │ │
│  │                   │  (Background │                        │ │
│  │                   │     Jobs)    │                        │ │
│  │                   │  Replicas: N │                        │ │
│  │                   └───────┬──────┘                        │ │
│  │                           │                               │ │
│  │         ┌─────────────────┴─────────────────┐             │ │
│  │         │                                   │             │ │
│  │  ┌──────▼────────┐               ┌────────▼─────────┐    │ │
│  │  │  PostgreSQL   │               │      Redis       │    │ │
│  │  │   Database    │               │      Cache       │    │ │
│  │  │               │               │                  │    │ │
│  │  │  (Self or     │               │   (Self or       │    │ │
│  │  │   External)   │               │    External)     │    │ │
│  │  └───────────────┘               └──────────────────┘    │ │
│  │                                                           │ │
│  └───────────────────────────────────────────────────────────┘ │
│                                                                 │
│  ┌───────────────────────────────────────────────────────────┐ │
│  │                     Storage Backend                       │ │
│  │  (S3, GCS, Azure Blob, OSS, or Filesystem PVC)            │ │
│  └───────────────────────────────────────────────────────────┘ │
│                                                                 │
│  ┌───────────────────────────────────────────────────────────┐ │
│  │              Gateway API (Optional Ingress)               │ │
│  │                                                           │ │
│  │  Core/Portal: harbor.example.com                         │ │
│  │  Notary: notary.harbor.example.com                       │ │
│  └───────────────────────────────────────────────────────────┘ │
└─────────────────────────────────────────────────────────────────┘
```

## Data Flow

### Image Push Operation

1. **Client authentication** - Docker client authenticates with Harbor Core API
2. **Authorization** - Harbor Core checks RBAC policies
3. **Image data** - Docker client pushes layers to Harbor Registry
4. **Storage** - Harbor Registry stores layers in configured storage backend (S3/GCS/Azure/OSS/Filesystem)
5. **Metadata** - Harbor Core stores image metadata in PostgreSQL
6. **Cache** - Frequently accessed data cached in Redis

### Image Pull Operation

1. **Client authentication** - Docker client authenticates with Harbor Core API
2. **Authorization** - Harbor Core checks pull policies
3. **Cache check** - Harbor checks Redis for cached metadata
4. **Registry request** - Harbor Registry retrieves layers from storage backend
5. **Layer delivery** - Layers streamed to Docker client

### Background Jobs

Harbor Jobservice handles asynchronous operations:
- Image vulnerability scanning
- Image replication between registries
- Garbage collection
- Audit log rotation
- Quota management

## Deployment Modes

### 1. Self-Managed Mode (Development/Testing)

All components deployed within the Kubernetes cluster:

- **PostgreSQL**: Deployed as StatefulSet with persistent storage
- **Redis**: Deployed as StatefulSet with persistent storage
- **Storage**: Filesystem-based PVC (not recommended for HA)

Configuration in `spec.proto`:
```protobuf
database {
  is_external: false
  managed_database {
    container {
      replicas: 1
      persistence_enabled: true
      disk_size: "20Gi"
    }
  }
}

cache {
  is_external: false
  managed_cache {
    container {
      replicas: 1
      persistence_enabled: true
      disk_size: "8Gi"
    }
  }
}

storage {
  type: filesystem
  filesystem {
    disk_size: "100Gi"
  }
}
```

### 2. Hybrid Mode (Recommended for Production)

External database/cache with cloud object storage:

- **PostgreSQL**: Amazon RDS, Cloud SQL, or Azure Database
- **Redis**: ElastiCache, Cloud Memorystore, or Azure Cache
- **Storage**: S3, GCS, or Azure Blob

Configuration example:
```protobuf
database {
  is_external: true
  external_database {
    host: "postgres.abc123.us-west-2.rds.amazonaws.com"
    port: 5432
    username: "harbor"
    password: "secret"
  }
}

cache {
  is_external: true
  external_cache {
    host: "redis.abc123.usw2.cache.amazonaws.com"
    port: 6379
  }
}

storage {
  type: s3
  s3 {
    bucket: "my-harbor-artifacts"
    region: "us-west-2"
    access_key: "AKI..."
    secret_key: "..."
  }
}
```

## Storage Backend Selection

The module supports five storage backends:

### 1. AWS S3 (Recommended for AWS)

- **Use case**: Production deployments on AWS
- **High availability**: Native HA with multi-AZ replication
- **Features**: Server-side encryption, versioning, lifecycle policies
- **Configuration**: Requires bucket, region, and credentials

### 2. Google Cloud Storage (Recommended for GCP)

- **Use case**: Production deployments on GCP
- **High availability**: Native HA with multi-region replication
- **Features**: Server-side encryption, versioning, lifecycle policies
- **Configuration**: Requires bucket and service account key

### 3. Azure Blob Storage (Recommended for Azure)

- **Use case**: Production deployments on Azure
- **High availability**: LRS/ZRS/GRS replication options
- **Features**: Server-side encryption, versioning, lifecycle policies
- **Configuration**: Requires account name, key, and container

### 4. Alibaba Cloud OSS

- **Use case**: Production deployments on Alibaba Cloud
- **High availability**: Multi-zone replication
- **Configuration**: Requires bucket, endpoint, and credentials

### 5. Filesystem (PVC)

- **Use case**: Development, testing, or single-node deployments
- **High availability**: Not suitable for multi-replica Harbor deployments
- **Limitations**: Cannot be shared across pods; single point of failure
- **Configuration**: Requires disk size and optional storage class

## Component Responsibilities

### Harbor Core

- **Authentication/Authorization**: User login, RBAC enforcement
- **Project Management**: Projects, repositories, members
- **Webhook Management**: Event notifications
- **API Gateway**: Central API endpoint for all Harbor operations
- **Configuration Management**: System settings and policies

### Harbor Portal

- **Web Interface**: User-friendly UI for Harbor management
- **Dashboard**: System statistics and monitoring
- **User Management**: User and group administration
- **Repository Browser**: Visual navigation of images and artifacts

### Harbor Registry

- **OCI Distribution**: Docker/OCI image push and pull
- **Layer Storage**: Manages artifact layers and blobs
- **Manifest Handling**: Image manifest validation and storage
- **Content Addressability**: SHA256-based blob storage

### Harbor Jobservice

- **Asynchronous Processing**: Background task execution
- **Vulnerability Scanning**: Integration with Trivy/Clair
- **Replication**: Multi-registry synchronization
- **Garbage Collection**: Reclaim storage from deleted artifacts
- **Quota Enforcement**: Storage and pull request limits

## Database Schema

Harbor uses PostgreSQL with separate databases for different components:

- **registry**: Core Harbor metadata (projects, repositories, users, artifacts)
- **clair**: Vulnerability database (if Clair scanner enabled)
- **notary_server**: Trust data for content trust (if Notary enabled)
- **notary_signer**: Signing keys for content trust (if Notary enabled)

## Cache Strategy

Redis is used for multiple caching layers:

- **Session Cache**: User session data
- **Authentication Cache**: Token validation results
- **API Response Cache**: Frequently accessed API responses
- **Job Queue**: Background job queue management
- **Rate Limiting**: Request rate limiting data

## Ingress Configuration

The module creates two optional Gateway API ingresses:

### Core/Portal Ingress

- **Purpose**: Web UI and API access
- **Hostname**: Configurable (e.g., `harbor.example.com`)
- **Backends**: Routes to both Core and Portal services
- **TLS**: Automatic certificate management via cert-manager

### Notary Ingress

- **Purpose**: Image signing and verification
- **Hostname**: Configurable (e.g., `notary.harbor.example.com`)
- **Backend**: Routes to Notary Server
- **TLS**: Automatic certificate management via cert-manager

## Security Considerations

### 1. Secrets Management

- **Admin Password**: Auto-generated or provided via spec
- **Database Credentials**: Stored as Kubernetes Secrets
- **Storage Credentials**: Stored as Kubernetes Secrets
- **TLS Certificates**: Managed by cert-manager

### 2. Network Policies

- Harbor components communicate over cluster network
- External access controlled via Gateway API
- Storage backend accessed via service endpoints

### 3. RBAC

- Kubernetes RBAC for Harbor pod permissions
- Harbor internal RBAC for user access control
- Project-level access control

## Scalability

### Horizontal Scaling

All Harbor components support horizontal scaling:

- **Core**: Scale replicas for increased API throughput
- **Portal**: Scale replicas for UI performance
- **Registry**: Scale replicas for pull/push performance
- **Jobservice**: Scale replicas for parallel job execution

### Vertical Scaling

Resource limits can be adjusted per component:

```protobuf
core_container {
  replicas: 3
  resources {
    limits {
      cpu: "2000m"
      memory: "4Gi"
    }
    requests {
      cpu: "500m"
      memory: "1Gi"
    }
  }
}
```

## High Availability Considerations

For production HA deployments:

1. **Multiple Replicas**: Deploy 3+ replicas of Core, Portal, Registry, Jobservice
2. **External Database**: Use managed PostgreSQL (RDS/Cloud SQL/Azure Database)
3. **External Cache**: Use managed Redis (ElastiCache/Memorystore/Azure Cache)
4. **Object Storage**: Use S3/GCS/Azure Blob (not filesystem)
5. **Load Balancing**: Gateway API distributes traffic across replicas
6. **Pod Anti-Affinity**: Spread pods across availability zones

## Monitoring and Observability

Stack outputs include:

- **Namespace**: Kubernetes namespace name
- **Service Names**: All Harbor component service names
- **Endpoints**: Internal and external endpoints
- **Port Forward Command**: For local debugging
- **Ingress Hostnames**: Public access URLs

## Custom Configuration

The `helm_values` field allows advanced customization:

```protobuf
helm_values {
  key: "trivy.enabled"
  value: "true"
}
```

Common customizations:
- Enable Trivy scanner
- Enable Notary for image signing
- Configure OIDC authentication
- Set up replication policies
- Configure webhook notifications

## Design Decisions

### Why Helm Chart?

- **Community Standard**: Official Harbor deployment method
- **Comprehensive**: Includes all components and configurations
- **Maintained**: Actively updated by Harbor team
- **Flexible**: Supports all Harbor features and options

### Why Gateway API for Ingress?

- **Modern Standard**: Kubernetes successor to Ingress
- **Rich Features**: Advanced routing, TLS, headers
- **Extensibility**: Provider-specific features
- **Future-proof**: Long-term Kubernetes direction

### Why Separate Database/Cache Configs?

- **Deployment Flexibility**: Choose self-managed or external
- **Cost Optimization**: Use managed services in production, self-hosted in dev
- **Migration Path**: Start simple, migrate to managed services later
- **Multi-tenancy**: Share external services across deployments

### Why Multiple Storage Backends?

- **Cloud Agnostic**: Deploy on any cloud platform
- **Performance**: Use cloud-native storage for best performance
- **Cost**: Optimize storage costs with lifecycle policies
- **Compatibility**: Support S3-compatible providers (MinIO, Ceph)

## Troubleshooting

### Common Issues

1. **Harbor UI not accessible**
   - Check Gateway API/ingress configuration
   - Verify cert-manager is installed
   - Check DNS configuration

2. **Image push fails**
   - Verify storage backend credentials
   - Check storage backend connectivity
   - Review Harbor Core logs

3. **Database connection errors**
   - Verify external database credentials
   - Check network connectivity to external DB
   - Review PostgreSQL logs

4. **Redis connection errors**
   - Verify external Redis credentials
   - Check network connectivity to external Redis
   - Review Redis logs

## References

- [Harbor Official Documentation](https://goharbor.io/docs/)
- [Harbor Helm Chart](https://github.com/goharbor/harbor-helm)
- [Harbor GitHub Repository](https://github.com/goharbor/harbor)
- [Kubernetes Gateway API](https://gateway-api.sigs.k8s.io/)

