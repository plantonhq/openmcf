# Kubernetes Zalando Postgres Operator Examples

This document provides practical examples for deploying the Zalando Postgres Operator using OpenMCF.

## Table of Contents

1. [Basic Operator Deployment](#basic-operator-deployment)
2. [Production Deployment with Backups](#production-deployment-with-backups)
3. [Custom Resource Limits](#custom-resource-limits)
4. [Multi-Environment Setup](#multi-environment-setup)
5. [Creating PostgreSQL Databases](#creating-postgresql-databases)

---

## Basic Operator Deployment

Minimal configuration for deploying the Zalando Postgres Operator without backup configuration.

### Manifest File

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: KubernetesZalandoPostgresOperator
metadata:
  name: postgres-operator
spec:
  target_cluster:
    cluster_name: my-gke-cluster
  namespace:
    value: postgres-operator
  container:
    resources:
      requests:
        cpu: 50m
        memory: 100Mi
      limits:
        cpu: 1000m
        memory: 1Gi
```

### Deployment Commands

```bash
# Save the manifest
cat > postgres-operator.yaml <<EOF
apiVersion: kubernetes.openmcf.org/v1
kind: KubernetesZalandoPostgresOperator
metadata:
  name: postgres-operator
spec:
  target_cluster:
    cluster_name: my-gke-cluster
  namespace:
    value: postgres-operator
  create_namespace: true
  container:
    resources:
      requests:
        cpu: 50m
        memory: 100Mi
      limits:
        cpu: 1000m
        memory: 1Gi
EOF

# Deploy using Pulumi
openmcf pulumi up --manifest postgres-operator.yaml

# Or deploy using Terraform
openmcf terraform apply --manifest postgres-operator.yaml
```

### Verification

```bash
# Check operator deployment
kubectl get deployment -n postgres-operator

# Check operator pods
kubectl get pods -n postgres-operator

# View operator logs
kubectl logs -n postgres-operator deployment/postgres-operator -f

# Verify CRDs are installed
kubectl get crd | grep postgresql
```

Expected output:
```
postgresqls.acid.zalan.do
operatorconfigurations.acid.zalan.do
postgresqls.acid.zalan.do
```

---

## Using Existing Namespace

If you have a pre-existing namespace and don't want the module to create it, set `create_namespace: false`:

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: KubernetesZalandoPostgresOperator
metadata:
  name: postgres-operator
spec:
  target_cluster:
    cluster_name: my-gke-cluster
  namespace:
    value: existing-postgres-namespace
  create_namespace: false  # Use existing namespace
  container:
    resources:
      requests:
        cpu: 50m
        memory: 100Mi
      limits:
        cpu: 1000m
        memory: 1Gi
```

**Important**: The namespace must already exist in the cluster when `create_namespace: false`.

---

## Production Deployment with Backups

Production-ready configuration with Cloudflare R2 backups, custom schedules, and enhanced resource limits.

### Prerequisites

1. **Cloudflare R2 Bucket**: Create an R2 bucket for backups
2. **R2 API Token**: Generate API token with read/write permissions
3. **Account ID**: Find your Cloudflare account ID in the dashboard

### Manifest File

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: KubernetesZalandoPostgresOperator
metadata:
  name: postgres-operator
  labels:
    org: acme-corp
    env: production
spec:
  target_cluster:
    cluster_name: prod-gke-cluster
  namespace:
    value: postgres-operator
  create_namespace: true
  container:
    resources:
      requests:
        cpu: 100m
        memory: 256Mi
      limits:
        cpu: 2000m
        memory: 2Gi
  backup_config:
    r2_config:
      cloudflare_account_id: "a1b2c3d4e5f6"
      bucket_name: "postgres-backups-prod"
      access_key_id: "YOUR_R2_ACCESS_KEY_ID"
      secret_access_key: "YOUR_R2_SECRET_ACCESS_KEY"
    backup_schedule: "0 2 * * *"  # 2 AM daily
    enable_wal_g_backup: true
    enable_wal_g_restore: true
    enable_clone_wal_g_restore: true
```

### Using Environment Variables for Secrets

For better security, use environment variables or secret management tools:

```bash
# Set environment variables
export CLOUDFLARE_ACCOUNT_ID="a1b2c3d4e5f6"
export R2_BUCKET_NAME="postgres-backups-prod"
export R2_ACCESS_KEY_ID="your-access-key-id"
export R2_SECRET_ACCESS_KEY="your-secret-access-key"

# Create manifest with envsubst
cat > postgres-operator-prod.yaml <<EOF
apiVersion: kubernetes.openmcf.org/v1
kind: KubernetesZalandoPostgresOperator
metadata:
  name: postgres-operator
  labels:
    org: acme-corp
    env: production
spec:
  target_cluster:
    cluster_name: prod-gke-cluster
  namespace:
    value: postgres-operator
  container:
    resources:
      requests:
        cpu: 100m
        memory: 256Mi
      limits:
        cpu: 2000m
        memory: 2Gi
  backup_config:
    r2_config:
      cloudflare_account_id: "${CLOUDFLARE_ACCOUNT_ID}"
      bucket_name: "${R2_BUCKET_NAME}"
      access_key_id: "${R2_ACCESS_KEY_ID}"
      secret_access_key: "${R2_SECRET_ACCESS_KEY}"
    backup_schedule: "0 2 * * *"
    enable_wal_g_backup: true
    enable_wal_g_restore: true
    enable_clone_wal_g_restore: true
EOF

# Deploy
openmcf pulumi up --manifest <(envsubst < postgres-operator-prod.yaml)
```

### Verification

```bash
# Check backup ConfigMap
kubectl get configmap -n postgres-operator postgres-pod-backup-config -o yaml

# Check backup Secret
kubectl get secret -n postgres-operator r2-postgres-backup-credentials

# Verify backup configuration
kubectl get configmap -n postgres-operator postgres-pod-backup-config -o jsonpath='{.data}' | jq
```

Expected ConfigMap data:
```json
{
  "AWS_ACCESS_KEY_ID": "your-access-key",
  "AWS_ENDPOINT": "https://a1b2c3d4e5f6.r2.cloudflarestorage.com",
  "AWS_FORCE_PATH_STYLE": "true",
  "AWS_REGION": "auto",
  "AWS_SECRET_ACCESS_KEY": "your-secret-key",
  "BACKUP_SCHEDULE": "0 2 * * *",
  "CLONE_USE_WALG_RESTORE": "true",
  "USE_WALG_BACKUP": "true",
  "USE_WALG_RESTORE": "true",
  "WALG_S3_PREFIX": "s3://postgres-backups-prod/backups/$(SCOPE)/$(PGVERSION)"
}
```

---

## Custom Resource Limits

Examples with different resource configurations for various cluster sizes.

### Small Cluster (Development)

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: KubernetesZalandoPostgresOperator
metadata:
  name: postgres-operator
  labels:
    env: development
spec:
  target_cluster:
    cluster_name: dev-gke-cluster
  namespace:
    value: postgres-operator
  create_namespace: true
  container:
    resources:
      requests:
        cpu: 25m
        memory: 50Mi
      limits:
        cpu: 250m
        memory: 256Mi
```

### Medium Cluster (Staging)

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: KubernetesZalandoPostgresOperator
metadata:
  name: postgres-operator
  labels:
    env: staging
spec:
  target_cluster:
    cluster_name: staging-gke-cluster
  namespace:
    value: postgres-operator
  create_namespace: true
  container:
    resources:
      requests:
        cpu: 100m
        memory: 256Mi
      limits:
        cpu: 1000m
        memory: 1Gi
```

### Large Cluster (Production)

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: KubernetesZalandoPostgresOperator
metadata:
  name: postgres-operator
  labels:
    env: production
spec:
  target_cluster:
    cluster_name: prod-gke-cluster
  namespace:
    value: postgres-operator
  create_namespace: true
  container:
    resources:
      requests:
        cpu: 200m
        memory: 512Mi
      limits:
        cpu: 4000m
        memory: 4Gi
```

---

## Multi-Environment Setup

Managing multiple environments with consistent configuration.

### Directory Structure

```
postgres-operator/
├── base/
│   └── operator.yaml
├── dev/
│   └── kustomization.yaml
├── staging/
│   └── kustomization.yaml
└── prod/
    └── kustomization.yaml
```

### Base Configuration

**base/operator.yaml**
```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: KubernetesZalandoPostgresOperator
metadata:
  name: postgres-operator
spec:
  target_cluster:
    cluster_name: my-gke-cluster
  namespace:
    value: postgres-operator
  create_namespace: true
  container:
    resources:
      requests:
        cpu: 50m
        memory: 100Mi
      limits:
        cpu: 1000m
        memory: 1Gi
```

### Development Override

**dev/operator-dev.yaml**
```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: KubernetesZalandoPostgresOperator
metadata:
  name: postgres-operator
  labels:
    env: development
spec:
  target_cluster:
    cluster_name: dev-gke-cluster
  namespace:
    value: postgres-operator-dev
  create_namespace: true
  container:
    resources:
      requests:
        cpu: 25m
        memory: 50Mi
      limits:
        cpu: 500m
        memory: 512Mi
```

### Production with Backups

**prod/operator-prod.yaml**
```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: KubernetesZalandoPostgresOperator
metadata:
  name: postgres-operator
  labels:
    org: acme-corp
    env: production
spec:
  target_cluster:
    cluster_name: prod-gke-cluster
  namespace:
    value: postgres-operator-prod
  create_namespace: true
  container:
    resources:
      requests:
        cpu: 200m
        memory: 512Mi
      limits:
        cpu: 4000m
        memory: 4Gi
  backup_config:
    r2_config:
      cloudflare_account_id: "${CLOUDFLARE_ACCOUNT_ID}"
      bucket_name: "postgres-backups-prod"
      access_key_id: "${R2_ACCESS_KEY_ID}"
      secret_access_key: "${R2_SECRET_ACCESS_KEY}"
    backup_schedule: "0 2 * * *"
    s3_prefix_template: "backups/$(SCOPE)/$(PGVERSION)"
    enable_wal_g_backup: true
    enable_wal_g_restore: true
    enable_clone_wal_g_restore: true
```

### Deployment Commands

```bash
# Deploy development
openmcf pulumi up --manifest dev/operator-dev.yaml

# Deploy production
openmcf pulumi up --manifest <(envsubst < prod/operator-prod.yaml)
```

---

## Creating PostgreSQL Databases

Once the operator is deployed, you can create PostgreSQL databases using the Zalando CRD.

### Basic Database

```yaml
apiVersion: "acid.zalan.do/v1"
kind: postgresql
metadata:
  name: my-app-db
  namespace: my-app
spec:
  teamId: "my-team"
  volume:
    size: 10Gi
  numberOfInstances: 1
  users:
    myapp:
      - superuser
      - createdb
  databases:
    myapp: myapp
  postgresql:
    version: "15"
```

### High Availability Database

```yaml
apiVersion: "acid.zalan.do/v1"
kind: postgresql
metadata:
  name: my-app-db-ha
  namespace: my-app
  labels:
    env: production
spec:
  teamId: "my-team"
  volume:
    size: 50Gi
    storageClass: fast-ssd
  numberOfInstances: 3  # 1 master + 2 replicas
  users:
    myapp:
      - superuser
      - createdb
    readonly:
      - login
  databases:
    myapp: myapp
  postgresql:
    version: "15"
    parameters:
      max_connections: "200"
      shared_buffers: "2GB"
      effective_cache_size: "6GB"
      maintenance_work_mem: "512MB"
      checkpoint_completion_target: "0.9"
      wal_buffers: "16MB"
      default_statistics_target: "100"
      random_page_cost: "1.1"
      effective_io_concurrency: "200"
      work_mem: "10MB"
      min_wal_size: "1GB"
      max_wal_size: "4GB"
  resources:
    requests:
      cpu: 1000m
      memory: 2Gi
    limits:
      cpu: 4000m
      memory: 8Gi
  enableConnectionPooler: true  # PgBouncer
  enableReplicaConnectionPooler: true
```

### Database with Clone from Backup

```yaml
apiVersion: "acid.zalan.do/v1"
kind: postgresql
metadata:
  name: my-app-db-restored
  namespace: my-app
spec:
  teamId: "my-team"
  volume:
    size: 50Gi
  numberOfInstances: 2
  clone:
    cluster: "my-app-db-ha"
    timestamp: "2024-01-15T10:00:00+00:00"  # Point-in-time recovery
  users:
    myapp:
      - superuser
      - createdb
  databases:
    myapp: myapp
  postgresql:
    version: "15"
```

### Connecting to the Database

```bash
# Get master service
kubectl get svc -n my-app | grep master

# Port-forward to local machine
kubectl port-forward svc/my-app-db -n my-app 5432:5432

# Connect using psql
psql -h localhost -U myapp -d myapp

# Get password from Secret
kubectl get secret myapp.my-app-db.credentials.postgresql.acid.zalan.do \
  -n my-app \
  -o jsonpath='{.data.password}' | base64 -d
```

---

## Advanced Backup Configuration

### Custom S3 Prefix Template

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: KubernetesZalandoPostgresOperator
metadata:
  name: postgres-operator
spec:
  target_cluster:
    cluster_name: prod-gke-cluster
  namespace:
    value: postgres-operator
  create_namespace: true
  container:
    resources:
      requests:
        cpu: 100m
        memory: 256Mi
      limits:
        cpu: 2000m
        memory: 2Gi
  backup_config:
    r2_config:
      cloudflare_account_id: "a1b2c3d4e5f6"
      bucket_name: "postgres-backups-prod"
      access_key_id: "${R2_ACCESS_KEY_ID}"
      secret_access_key: "${R2_SECRET_ACCESS_KEY}"
    backup_schedule: "0 */6 * * *"  # Every 6 hours
    s3_prefix_template: "prod/pg/$(SCOPE)/v$(PGVERSION)/$(CLUSTER)"
    enable_wal_g_backup: true
    enable_wal_g_restore: true
    enable_clone_wal_g_restore: true
```

### Backup Schedules

| Schedule | Description | Cron Expression |
|----------|-------------|-----------------|
| Hourly | Every hour | `0 * * * *` |
| Every 6 hours | 4 times daily | `0 */6 * * *` |
| Daily at 2 AM | Once daily | `0 2 * * *` |
| Every 12 hours | Twice daily | `0 */12 * * *` |
| Weekly Sunday 2 AM | Once weekly | `0 2 * * 0` |

---

## Troubleshooting Examples

### Check Operator Status

```bash
# Operator deployment status
kubectl get deployment -n postgres-operator postgres-operator

# Operator pod status
kubectl get pods -n postgres-operator

# Operator logs
kubectl logs -n postgres-operator deployment/postgres-operator --tail=100

# Check for errors
kubectl logs -n postgres-operator deployment/postgres-operator | grep -i error
```

### Verify Backup Configuration

```bash
# Check if ConfigMap exists
kubectl get configmap -n postgres-operator postgres-pod-backup-config

# View ConfigMap contents
kubectl get configmap -n postgres-operator postgres-pod-backup-config -o yaml

# Check Secret
kubectl get secret -n postgres-operator r2-postgres-backup-credentials

# Test R2 connectivity from a database pod
kubectl exec -n my-app my-app-db-0 -- \
  env AWS_ACCESS_KEY_ID=xxx AWS_SECRET_ACCESS_KEY=yyy \
  aws s3 ls s3://postgres-backups-prod/ --endpoint-url https://xxx.r2.cloudflarestorage.com
```

### Database Creation Issues

```bash
# Check if CRD is installed
kubectl get crd postgresqls.acid.zalan.do

# List all PostgreSQL resources
kubectl get postgresql --all-namespaces

# Describe a specific database
kubectl describe postgresql my-app-db -n my-app

# Check operator logs for specific database
kubectl logs -n postgres-operator deployment/postgres-operator | grep my-app-db

# Check StatefulSet creation
kubectl get statefulset -n my-app

# Check PVC creation
kubectl get pvc -n my-app
```

---

## Complete Production Example

This example combines all best practices for a production deployment:

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: KubernetesZalandoPostgresOperator
metadata:
  name: postgres-operator
  id: pgop-prod-001
  labels:
    org: acme-corp
    env: production
    team: platform
    cost-center: engineering
spec:
  target_cluster:
    cluster_name: prod-gke-cluster
  namespace:
    value: postgres-operator-prod
  create_namespace: true
  container:
    resources:
      requests:
        cpu: 200m
        memory: 512Mi
      limits:
        cpu: 4000m
        memory: 4Gi
  backup_config:
    r2_config:
      cloudflare_account_id: "a1b2c3d4e5f6g7h8i9j0"
      bucket_name: "acme-postgres-backups-prod-us-east"
      access_key_id: "${R2_ACCESS_KEY_ID}"
      secret_access_key: "${R2_SECRET_ACCESS_KEY}"
    backup_schedule: "0 2 * * *"
    s3_prefix_template: "prod/postgres/$(SCOPE)/pg$(PGVERSION)"
    enable_wal_g_backup: true
    enable_wal_g_restore: true
    enable_clone_wal_g_restore: true
```

### Deployment Script

```bash
#!/bin/bash
set -e

# Load secrets from vault or environment
export CLOUDFLARE_ACCOUNT_ID=$(vault kv get -field=account_id secret/postgres/r2)
export R2_ACCESS_KEY_ID=$(vault kv get -field=access_key_id secret/postgres/r2)
export R2_SECRET_ACCESS_KEY=$(vault kv get -field=secret_access_key secret/postgres/r2)

# Deploy operator
echo "Deploying Postgres Operator to production..."
openmcf pulumi up --manifest <(envsubst < postgres-operator-prod.yaml) --stack prod

# Verify deployment
echo "Verifying deployment..."
kubectl wait --for=condition=available --timeout=300s \
  deployment/postgres-operator -n postgres-operator

echo "Postgres Operator deployed successfully!"
echo "Backup configuration:"
kubectl get configmap -n postgres-operator postgres-pod-backup-config -o yaml
```

---

## Next Steps

After deploying the operator:

1. **Create databases** using the Zalando `postgresql` CRD
2. **Configure monitoring** with Prometheus and Grafana
3. **Test backup/restore** procedures
4. **Set up alerts** for operator and database health
5. **Review security** settings and RBAC permissions

For more information, see:
- [README.md](./README.md) - Component documentation
- [docs/README.md](./docs/README.md) - Architecture and design decisions
- [Zalando Postgres Operator Docs](https://postgres-operator.readthedocs.io/)

