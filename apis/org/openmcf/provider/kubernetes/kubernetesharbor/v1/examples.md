# Harbor Kubernetes API - Example Configurations

## Namespace Management

All Harbor deployments require a namespace. Control namespace management with the `create_namespace` flag:

- **`create_namespace: true`** - Component creates and manages the namespace with appropriate labels
- **`create_namespace: false`** - Deploy into an existing namespace (must exist before deployment)

**Managed namespace example:**
```yaml
spec:
  namespace:
    value: "harbor-prod"
  create_namespace: true  # Component creates the namespace
```

**Existing namespace example:**
```yaml
spec:
  namespace:
    value: "container-registry"
  create_namespace: false  # Must exist before deployment
```

---

## Example w/ Basic Configuration (Evaluation)

### Create Using CLI

Create a YAML file using the example shown below. After the YAML is created, use the following command to apply it:

```shell
planton apply -f <yaml-path>
```

### Basic Example

This basic example demonstrates a minimal configuration for deploying Harbor Kubernetes for evaluation or development purposes, using default settings with self-managed PostgreSQL and Redis.

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: HarborKubernetes
metadata:
  name: basic-harbor
spec:
  targetCluster:
    clusterName: my-gke-cluster
  namespace:
    value: basic-harbor
  create_namespace: true
  coreContainer:
    replicas: 1
    resources:
      requests:
        cpu: 200m
        memory: 512Mi
      limits:
        cpu: 1000m
        memory: 2Gi
  portalContainer:
    replicas: 1
    resources:
      requests:
        cpu: 100m
        memory: 256Mi
      limits:
        cpu: 500m
        memory: 512Mi
  registryContainer:
    replicas: 1
    resources:
      requests:
        cpu: 200m
        memory: 512Mi
      limits:
        cpu: 1000m
        memory: 2Gi
  jobserviceContainer:
    replicas: 1
    resources:
      requests:
        cpu: 100m
        memory: 256Mi
      limits:
        cpu: 1000m
        memory: 1Gi
  database:
    isExternal: false
    managedDatabase:
      container:
        replicas: 1
        isPersistenceEnabled: true
        diskSize: 20Gi
        resources:
          requests:
            cpu: 200m
            memory: 512Mi
          limits:
            cpu: 1000m
            memory: 2Gi
  cache:
    isExternal: false
    managedCache:
      container:
        replicas: 1
        isPersistenceEnabled: true
        diskSize: 8Gi
        resources:
          requests:
            cpu: 100m
            memory: 256Mi
          limits:
            cpu: 500m
            memory: 512Mi
  storage:
    type: filesystem
    filesystem:
      diskSize: 100Gi
```

---

## Example w/ Production Configuration (High Availability with AWS S3)

This example demonstrates a production-ready deployment with external managed PostgreSQL (AWS RDS), external Redis (ElastiCache), and S3 object storage for high availability.

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: HarborKubernetes
metadata:
  name: production-harbor
spec:
  targetCluster:
    clusterName: my-gke-cluster
  namespace:
    value: production-harbor
  create_namespace: true
  coreContainer:
    replicas: 2
    resources:
      requests:
        cpu: 500m
        memory: 1Gi
      limits:
        cpu: 2000m
        memory: 4Gi
  portalContainer:
    replicas: 2
    resources:
      requests:
        cpu: 200m
        memory: 512Mi
      limits:
        cpu: 1000m
        memory: 1Gi
  registryContainer:
    replicas: 2
    resources:
      requests:
        cpu: 500m
        memory: 1Gi
      limits:
        cpu: 2000m
        memory: 4Gi
  jobserviceContainer:
    replicas: 2
    resources:
      requests:
        cpu: 200m
        memory: 512Mi
      limits:
        cpu: 1000m
        memory: 2Gi
  database:
    isExternal: true
    externalDatabase:
      host: harbor-db.xxxxxxxxxxxx.us-west-2.rds.amazonaws.com
      port: 5432
      username: harbor
      password: ${HARBOR_DB_PASSWORD}
      coreDatabase: registry
      clairDatabase: clair
      notaryServerDatabase: notary_server
      notarySignerDatabase: notary_signer
      useSsl: true
  cache:
    isExternal: true
    externalCache:
      host: harbor-redis.xxxx.use1.cache.amazonaws.com
      port: 6379
      password: ${REDIS_PASSWORD}
      useSentinel: true
      sentinelMasterSet: mymaster
  storage:
    type: s3
    s3:
      bucket: my-harbor-artifacts
      region: us-west-2
      accessKey: ${AWS_ACCESS_KEY_ID}
      secretKey: ${AWS_SECRET_ACCESS_KEY}
      regionEndpoint: false
      encrypt: true
      secure: true
```

---

## Example w/ Google Cloud Storage and Cloud SQL

This example demonstrates a production deployment on Google Cloud Platform using GCS for storage and Cloud SQL for PostgreSQL.

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: HarborKubernetes
metadata:
  name: harbor-gcp
spec:
  targetCluster:
    clusterName: my-gke-cluster
  namespace:
    value: harbor-gcp
  create_namespace: true
  coreContainer:
    replicas: 2
    resources:
      requests:
        cpu: 500m
        memory: 1Gi
      limits:
        cpu: 2000m
        memory: 4Gi
  portalContainer:
    replicas: 2
  registryContainer:
    replicas: 2
    resources:
      requests:
        cpu: 500m
        memory: 1Gi
      limits:
        cpu: 2000m
        memory: 4Gi
  jobserviceContainer:
    replicas: 2
  database:
    isExternal: true
    externalDatabase:
      host: 10.0.0.3
      port: 5432
      username: harbor
      password: ${CLOUDSQL_PASSWORD}
      useSsl: false
  cache:
    isExternal: true
    externalCache:
      host: 10.0.0.4
      port: 6379
      password: ${REDIS_PASSWORD}
  storage:
    type: gcs
    gcs:
      bucket: my-harbor-artifacts-bucket
      keyData: ${GCS_SERVICE_ACCOUNT_KEY_BASE64}
      chunkSize: 5242880
```

---

## Example w/ Azure Blob Storage

This example demonstrates deploying Harbor with Azure Blob Storage and Azure Database for PostgreSQL.

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: HarborKubernetes
metadata:
  name: harbor-azure
spec:
  targetCluster:
    clusterName: my-aks-cluster
  namespace:
    value: harbor-azure
  create_namespace: true
  coreContainer:
    replicas: 2
  portalContainer:
    replicas: 2
  registryContainer:
    replicas: 2
  jobserviceContainer:
    replicas: 2
  database:
    isExternal: true
    externalDatabase:
      host: harbor-db.postgres.database.azure.com
      port: 5432
      username: harbor@harbor-db
      password: ${AZURE_POSTGRES_PASSWORD}
      useSsl: true
  cache:
    isExternal: true
    externalCache:
      host: harbor-redis.redis.cache.windows.net
      port: 6380
      password: ${AZURE_REDIS_PASSWORD}
  storage:
    type: azure
    azure:
      accountName: myharborstorageacct
      accountKey: ${AZURE_STORAGE_KEY}
      container: harbor-artifacts
```

---

## Example w/ Ingress Configuration

This example demonstrates how to configure ingress for both Harbor UI and Notary service with custom hostnames.

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: HarborKubernetes
metadata:
  name: harbor-with-ingress
spec:
  kubernetesProviderConfigId: my-cluster-credential-id
  coreContainer:
    replicas: 2
  portalContainer:
    replicas: 2
  registryContainer:
    replicas: 2
  jobserviceContainer:
    replicas: 2
  database:
    isExternal: false
    managedDatabase:
      container:
        replicas: 1
        isPersistenceEnabled: true
        diskSize: 50Gi
  cache:
    isExternal: false
    managedCache:
      container:
        replicas: 1
        isPersistenceEnabled: true
        diskSize: 10Gi
  storage:
    type: s3
    s3:
      bucket: harbor-storage
      region: us-east-1
      accessKey: ${AWS_ACCESS_KEY}
      secretKey: ${AWS_SECRET_KEY}
      encrypt: true
      secure: true
  ingress:
    core:
      enabled: true
      hostname: harbor.example.com
    notary:
      enabled: true
      hostname: notary.harbor.example.com
```

---

## Example w/ Minimal Resources (Development)

This example uses minimal resource allocations suitable for local development or resource-constrained environments.

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: HarborKubernetes
metadata:
  name: harbor-dev
spec:
  targetCluster:
    clusterName: my-gke-cluster
  namespace:
    value: harbor-dev
  create_namespace: true
  coreContainer:
    replicas: 1
    resources:
      requests:
        cpu: 100m
        memory: 256Mi
      limits:
        cpu: 500m
        memory: 1Gi
  portalContainer:
    replicas: 1
    resources:
      requests:
        cpu: 50m
        memory: 128Mi
      limits:
        cpu: 200m
        memory: 256Mi
  registryContainer:
    replicas: 1
    resources:
      requests:
        cpu: 100m
        memory: 256Mi
      limits:
        cpu: 500m
        memory: 1Gi
  jobserviceContainer:
    replicas: 1
    resources:
      requests:
        cpu: 50m
        memory: 128Mi
      limits:
        cpu: 200m
        memory: 512Mi
  database:
    isExternal: false
    managedDatabase:
      container:
        replicas: 1
        isPersistenceEnabled: false
        resources:
          requests:
            cpu: 100m
            memory: 256Mi
          limits:
            cpu: 500m
            memory: 1Gi
  cache:
    isExternal: false
    managedCache:
      container:
        replicas: 1
        isPersistenceEnabled: false
        resources:
          requests:
            cpu: 50m
            memory: 128Mi
          limits:
            cpu: 200m
            memory: 256Mi
  storage:
    type: filesystem
    filesystem:
      diskSize: 20Gi
```

---

## Example w/ External MinIO (S3-Compatible)

This example demonstrates using MinIO (S3-compatible object storage) as the storage backend.

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: HarborKubernetes
metadata:
  name: harbor-minio
spec:
  targetCluster:
    clusterName: my-gke-cluster
  namespace:
    value: harbor-minio
  create_namespace: true
  coreContainer:
    replicas: 2
  portalContainer:
    replicas: 2
  registryContainer:
    replicas: 2
  jobserviceContainer:
    replicas: 2
  database:
    isExternal: false
    managedDatabase:
      container:
        replicas: 1
        isPersistenceEnabled: true
        diskSize: 30Gi
  cache:
    isExternal: false
    managedCache:
      container:
        replicas: 1
        isPersistenceEnabled: true
        diskSize: 10Gi
  storage:
    type: s3
    s3:
      bucket: harbor
      region: us-east-1
      accessKey: minio-access-key
      secretKey: minio-secret-key
      endpointUrl: http://minio.storage.svc.cluster.local:9000
      secure: false
      regionEndpoint: false
```

---

## Example w/ Custom Helm Values (Trivy Scanner)

This example demonstrates enabling and configuring the Trivy vulnerability scanner through Helm values.

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: HarborKubernetes
metadata:
  name: harbor-with-trivy
spec:
  targetCluster:
    clusterName: my-gke-cluster
  namespace:
    value: harbor-with-trivy
  create_namespace: true
  coreContainer:
    replicas: 2
  portalContainer:
    replicas: 2
  registryContainer:
    replicas: 2
  jobserviceContainer:
    replicas: 2
  database:
    isExternal: false
    managedDatabase:
      container:
        replicas: 1
        isPersistenceEnabled: true
        diskSize: 50Gi
  cache:
    isExternal: false
    managedCache:
      container:
        replicas: 1
        isPersistenceEnabled: true
        diskSize: 10Gi
  storage:
    type: s3
    s3:
      bucket: harbor-prod
      region: us-west-2
      accessKey: ${AWS_ACCESS_KEY}
      secretKey: ${AWS_SECRET_KEY}
      encrypt: true
  helmValues:
    # Enable Trivy scanner
    trivy.enabled: "true"
    trivy.replicas: "2"
    trivy.resources.requests.cpu: "200m"
    trivy.resources.requests.memory: "512Mi"
    trivy.resources.limits.cpu: "1000m"
    trivy.resources.limits.memory: "2Gi"
    # Configure Trivy database update
    trivy.gitHubToken: "${GITHUB_TOKEN}"
```

---

## Example w/ OIDC Authentication (via Helm Values)

This example demonstrates configuring OIDC authentication for Harbor using Helm values.

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: HarborKubernetes
metadata:
  name: harbor-oidc
spec:
  targetCluster:
    clusterName: my-gke-cluster
  namespace:
    value: harbor-oidc
  create_namespace: true
  coreContainer:
    replicas: 2
  portalContainer:
    replicas: 2
  registryContainer:
    replicas: 2
  jobserviceContainer:
    replicas: 2
  database:
    isExternal: true
    externalDatabase:
      host: postgres.database.svc
      username: harbor
      password: ${DB_PASSWORD}
  cache:
    isExternal: true
    externalCache:
      host: redis.cache.svc
      password: ${REDIS_PASSWORD}
  storage:
    type: s3
    s3:
      bucket: harbor-artifacts
      region: us-east-1
      accessKey: ${AWS_ACCESS_KEY}
      secretKey: ${AWS_SECRET_KEY}
      encrypt: true
  ingress:
    core:
      enabled: true
      hostname: harbor.company.com
  helmValues:
    # Configure OIDC authentication
    core.oidcClientId: "harbor-client"
    core.oidcClientSecret: "${OIDC_CLIENT_SECRET}"
    core.oidcEndpoint: "https://auth.company.com"
    core.oidcScope: "openid,email,profile"
    core.oidcName: "Company SSO"
    core.oidcAutoOnboard: "true"
    core.oidcGroupsClaim: "groups"
```

---

## Example w/ Replication Policy (via Helm Values)

This example demonstrates configuring replication to another Harbor instance for disaster recovery.

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: HarborKubernetes
metadata:
  name: harbor-primary
spec:
  targetCluster:
    clusterName: my-gke-cluster
  namespace:
    value: harbor-primary
  create_namespace: true
  coreContainer:
    replicas: 2
  portalContainer:
    replicas: 2
  registryContainer:
    replicas: 2
  jobserviceContainer:
    replicas: 2
  database:
    isExternal: true
    externalDatabase:
      host: primary-db.us-west-2.rds.amazonaws.com
      username: harbor
      password: ${DB_PASSWORD}
  cache:
    isExternal: true
    externalCache:
      host: primary-redis.cache.amazonaws.com
      password: ${REDIS_PASSWORD}
  storage:
    type: s3
    s3:
      bucket: harbor-us-west-2
      region: us-west-2
      accessKey: ${AWS_ACCESS_KEY}
      secretKey: ${AWS_SECRET_KEY}
      encrypt: true
  ingress:
    core:
      enabled: true
      hostname: harbor-primary.company.com
  helmValues:
    # Configure metrics for monitoring
    metrics.enabled: "true"
    metrics.serviceMonitor.enabled: "true"
```

---

## Example w/ Notary and Content Trust

This example demonstrates enabling Notary for Docker Content Trust and image signing.

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: HarborKubernetes
metadata:
  name: harbor-notary
spec:
  targetCluster:
    clusterName: my-gke-cluster
  namespace:
    value: harbor-notary
  create_namespace: true
  coreContainer:
    replicas: 2
  portalContainer:
    replicas: 2
  registryContainer:
    replicas: 2
  jobserviceContainer:
    replicas: 2
  database:
    isExternal: false
    managedDatabase:
      container:
        replicas: 1
        isPersistenceEnabled: true
        diskSize: 50Gi
  cache:
    isExternal: false
    managedCache:
      container:
        replicas: 1
        isPersistenceEnabled: true
        diskSize: 10Gi
  storage:
    type: s3
    s3:
      bucket: harbor-secure
      region: us-east-1
      accessKey: ${AWS_ACCESS_KEY}
      secretKey: ${AWS_SECRET_KEY}
      encrypt: true
  ingress:
    core:
      enabled: true
      hostname: harbor.secure.company.com
    notary:
      enabled: true
      hostname: notary.harbor.secure.company.com
  helmValues:
    # Enable Notary for content trust
    notary.enabled: "true"
    notary.server.replicas: "2"
    notary.signer.replicas: "2"
```

---

## Important Notes

### Storage Backend Selection

For production deployments, **always use external object storage** (S3, GCS, Azure Blob) instead of filesystem storage:

- **Filesystem**: Limited to single Registry pod (ReadWriteOnce PVC). Not suitable for HA.
- **S3/GCS/Azure**: Enables multi-replica Registry deployments for high availability.
- **Performance**: Object storage provides better scalability for large image volumes.
- **Cost**: Lifecycle policies can automatically move cold data to cheaper storage tiers.

### Database and Cache HA

For production deployments:

- **External PostgreSQL**: Use managed services (AWS RDS, Google Cloud SQL, Azure Database) with automatic backups and HA
- **External Redis**: Use Redis Sentinel or managed services (ElastiCache, Memorystore) for cache HA
- **Self-Managed**: Suitable only for development and testing environments

### Ingress and TLS

Ensure your Kubernetes cluster has:
- **Ingress Controller**: nginx-ingress or other compatible controller
- **cert-manager**: For automatic TLS certificate provisioning from Let's Encrypt
- **DNS Configuration**: Point hostnames to ingress load balancer

```yaml
# Example cert-manager annotation (add via helmValues)
helmValues:
  expose.ingress.annotations.cert-manager\.io/cluster-issuer: "letsencrypt-prod"
```

### Resource Planning

The resource allocations in these examples are starting points. Monitor actual usage and adjust based on:
- Number of concurrent users and CI/CD pipelines
- Image push/pull frequency and sizes
- Number of vulnerability scans
- Replication job volumes

### Security Best Practices

1. **Use Secrets**: Never hardcode credentials in manifests
2. **Enable SSL**: Always use SSL for external databases and Redis
3. **Enable Scanning**: Configure Trivy or Clair for vulnerability scanning
4. **RBAC**: Implement project-based access control
5. **Content Trust**: Enable Notary for critical production images
6. **Network Policies**: Restrict network access to Harbor components

### Helm Values Reference

For advanced configuration options not exposed in the spec, use the `helmValues` field. Refer to:
- [Harbor Helm Chart Values](https://github.com/goharbor/harbor-helm/blob/main/values.yaml)
- [Harbor Installation Guide](https://goharbor.io/docs/latest/install-config/)


