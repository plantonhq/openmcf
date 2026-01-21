---
title: "State Backends"
description: "Configure state storage for Pulumi, OpenTofu, and Terraform using manifest labels for automatic backend detection"
icon: "database"
order: 25
---

# State Backends

State backends store the state of your infrastructure, tracking what resources exist and their current configuration. Project Planton supports automatic backend detection from manifest labels for all three provisioners: Pulumi, OpenTofu, and Terraform.

---

## Overview

Infrastructure as Code tools need to track what they've deployed. This tracking happens through "state" - a record of resources, their properties, and their relationships. Where this state is stored is the "backend."

**Why backends matter:**

- **Collaboration**: Teams need shared state to avoid conflicts
- **Persistence**: State survives beyond your local machine
- **Locking**: Prevents simultaneous modifications
- **History**: Enables rollbacks and auditing

Project Planton automatically detects backend configuration from labels in your manifest, eliminating the need for separate backend configuration files.

---

## Quick Reference

| Provisioner | Backend Type Label | Backend Object/Stack Label |
|-------------|-------------------|---------------------------|
| **Terraform** | `terraform.project-planton.org/backend.type` | `terraform.project-planton.org/backend.object` |
| **OpenTofu** | `tofu.project-planton.org/backend.type` | `tofu.project-planton.org/backend.object` |
| **Pulumi** | N/A (uses Pulumi Cloud or local) | `pulumi.project-planton.org/stack.name` |

---

## Pulumi State

Pulumi stores state either in Pulumi Cloud (default) or locally. The stack name label identifies where your state is stored.

### Stack Name Label

```yaml
apiVersion: kubernetes.project-planton.org/v1
kind: KubernetesPostgres
metadata:
  name: app-database
  labels:
    project-planton.org/provisioner: pulumi
    pulumi.project-planton.org/stack.name: prod.PostgresKubernetes.app-database
spec:
  container:
    replicas: 1
```

### Stack Name Format

The stack name follows the pattern: `<environment>.<project>.<stack>`

- **environment**: Deployment environment (prod, staging, dev)
- **project**: Pulumi project name (usually matches the kind)
- **stack**: Unique identifier for this deployment

### Pulumi Backend Options

**1. Pulumi Cloud (Recommended)**

```bash
# Login to Pulumi Cloud
pulumi login

# State is automatically stored in Pulumi Cloud
project-planton apply -f database.yaml
```

**2. Local Backend**

```bash
# Use local filesystem for state
pulumi login --local

# State stored in ~/.pulumi/
project-planton apply -f database.yaml
```

**3. Self-hosted Backend (S3, GCS, Azure)**

```bash
# Use S3 for state
pulumi login s3://my-pulumi-state-bucket

# Or GCS
pulumi login gs://my-pulumi-state-bucket

# Or Azure Blob
pulumi login azblob://my-container
```

---

## OpenTofu / Terraform State

OpenTofu and Terraform use a backend configuration to store state. Project Planton reads this configuration from manifest labels.

### Label Format

Each provisioner uses its own label prefix:

**For Terraform:**
```yaml
metadata:
  labels:
    project-planton.org/provisioner: terraform
    terraform.project-planton.org/backend.type: "s3"
    terraform.project-planton.org/backend.object: "bucket-name/path/to/state.tfstate"
```

**For OpenTofu:**
```yaml
metadata:
  labels:
    project-planton.org/provisioner: tofu
    tofu.project-planton.org/backend.type: "gcs"
    tofu.project-planton.org/backend.object: "bucket-name/path/to/state.tfstate"
```

### Backward Compatibility

For backward compatibility, OpenTofu also accepts `terraform.project-planton.org/*` labels if the provisioner-specific `tofu.project-planton.org/*` labels are not present:

```yaml
metadata:
  labels:
    project-planton.org/provisioner: tofu
    # Legacy labels - still work with OpenTofu
    terraform.project-planton.org/backend.type: "s3"
    terraform.project-planton.org/backend.object: "my-bucket/state.tfstate"
```

However, we recommend using the provisioner-specific labels for clarity.

---

## Supported Backend Types

### Amazon S3

Store state in an S3 bucket with optional DynamoDB locking.

```yaml
apiVersion: aws.project-planton.org/v1
kind: AwsVpc
metadata:
  name: production-vpc
  labels:
    project-planton.org/provisioner: terraform
    terraform.project-planton.org/backend.type: "s3"
    terraform.project-planton.org/backend.object: "my-terraform-state/vpc/production.tfstate"
spec:
  cidrBlock: 10.0.0.0/16
  region: us-west-2
```

**Object format:** `bucket-name/path/to/state.tfstate`

**Prerequisites:**
- S3 bucket must exist
- IAM permissions: `s3:GetObject`, `s3:PutObject`, `s3:DeleteObject`, `s3:ListBucket`
- Optional: DynamoDB table for state locking (configured via environment)

---

### Google Cloud Storage

Store state in a GCS bucket.

```yaml
apiVersion: gcp.project-planton.org/v1
kind: GkeCluster
metadata:
  name: staging-cluster
  labels:
    project-planton.org/provisioner: tofu
    tofu.project-planton.org/backend.type: "gcs"
    tofu.project-planton.org/backend.object: "my-gcs-state-bucket/gke/staging-cluster"
spec:
  projectId: my-gcp-project
  region: us-central1
```

**Object format:** `bucket-name/prefix/path`

**Prerequisites:**
- GCS bucket must exist
- IAM permissions: `storage.objects.get`, `storage.objects.create`, `storage.objects.delete`, `storage.objects.list`

---

### Azure Storage

Store state in Azure Blob Storage.

```yaml
apiVersion: azure.project-planton.org/v1
kind: AzureAksCluster
metadata:
  name: production-aks
  labels:
    project-planton.org/provisioner: terraform
    terraform.project-planton.org/backend.type: "azurerm"
    terraform.project-planton.org/backend.object: "tfstate-container/aks/production"
spec:
  location: eastus
  nodeCount: 3
```

**Object format:** `container-name/path/to/state`

**Prerequisites:**
- Storage account and container must exist
- Storage account name configured via environment
- IAM permissions: Storage Blob Data Contributor

---

### Local Backend

Store state on the local filesystem. **Not recommended for production or team use.**

```yaml
apiVersion: kubernetes.project-planton.org/v1
kind: MicroserviceKubernetes
metadata:
  name: test-service
  labels:
    project-planton.org/provisioner: tofu
    tofu.project-planton.org/backend.type: "local"
    tofu.project-planton.org/backend.object: "/tmp/test-service.tfstate"
spec:
  replicas: 1
```

**Use cases:**
- Local development
- Testing
- Single-user scenarios

**Limitations:**
- No state locking
- Not shareable between team members
- Lost if machine is wiped

---

## Complete Examples

### Example 1: AWS Infrastructure with S3 Backend

```yaml
apiVersion: aws.project-planton.org/v1
kind: AwsRdsInstance
metadata:
  name: app-database
  labels:
    # Provisioner selection
    project-planton.org/provisioner: terraform
    
    # Backend configuration
    terraform.project-planton.org/backend.type: "s3"
    terraform.project-planton.org/backend.object: "company-terraform-state/rds/app-database/production"
spec:
  engine: postgres
  engineVersion: "15"
  instanceClass: db.t3.medium
  allocatedStorage: 100
  region: us-west-2
```

**Deploy:**
```bash
project-planton apply -f database.yaml
```

---

### Example 2: GCP Infrastructure with GCS Backend (OpenTofu)

```yaml
apiVersion: gcp.project-planton.org/v1
kind: GcpCloudRun
metadata:
  name: api-service
  labels:
    # Provisioner selection
    project-planton.org/provisioner: tofu
    
    # Backend configuration (OpenTofu-specific)
    tofu.project-planton.org/backend.type: "gcs"
    tofu.project-planton.org/backend.object: "company-tofu-state/cloud-run/api-service/prod"
spec:
  projectId: my-gcp-project
  region: us-central1
  image: gcr.io/my-project/api:latest
```

**Deploy:**
```bash
project-planton apply -f api-service.yaml
```

---

### Example 3: Multi-Environment with Kustomize

**Base manifest** (`base/database.yaml`):
```yaml
apiVersion: aws.project-planton.org/v1
kind: AwsRdsInstance
metadata:
  name: app-database
  labels:
    project-planton.org/provisioner: terraform
    terraform.project-planton.org/backend.type: "s3"
    # Object path will be patched per environment
spec:
  engine: postgres
  instanceClass: db.t3.small
```

**Production overlay** (`overlays/prod/kustomization.yaml`):
```yaml
patches:
  - patch: |
      - op: replace
        path: /metadata/labels/terraform.project-planton.org~1backend.object
        value: "terraform-state/rds/production"
      - op: replace
        path: /spec/instanceClass
        value: db.t3.large
    target:
      kind: AwsRdsInstance
```

**Deploy:**
```bash
project-planton apply --kustomize-dir . --overlay prod
```

---

### Example 4: Pulumi with Stack Name

```yaml
apiVersion: kubernetes.project-planton.org/v1
kind: KubernetesPostgres
metadata:
  name: analytics-db
  labels:
    # Provisioner selection
    project-planton.org/provisioner: pulumi
    
    # Stack name for state identification
    pulumi.project-planton.org/stack.name: production.PostgresKubernetes.analytics-db
spec:
  container:
    replicas: 3
    resources:
      limits:
        cpu: 2000m
        memory: 8Gi
```

**Deploy:**
```bash
# Ensure you're logged into Pulumi
pulumi login

# Deploy
project-planton apply -f analytics-db.yaml
```

---

## Best Practices

### 1. Use Consistent Naming

Establish a naming convention for state paths:

```
<bucket>/<resource-type>/<resource-name>/<environment>
```

Example: `terraform-state/rds/app-database/production`

### 2. Separate State by Environment

Use different buckets or prefixes for different environments:

```yaml
# Production
terraform.project-planton.org/backend.object: "prod-terraform-state/vpc/main"

# Staging
terraform.project-planton.org/backend.object: "staging-terraform-state/vpc/main"

# Development
terraform.project-planton.org/backend.object: "dev-terraform-state/vpc/main"
```

### 3. Enable Versioning

Always enable versioning on your state bucket:

**S3:**
```bash
aws s3api put-bucket-versioning \
  --bucket my-terraform-state \
  --versioning-configuration Status=Enabled
```

**GCS:**
```bash
gsutil versioning set on gs://my-terraform-state
```

### 4. Enable Encryption

Encrypt state at rest:

**S3:**
```bash
aws s3api put-bucket-encryption \
  --bucket my-terraform-state \
  --server-side-encryption-configuration '{
    "Rules": [{"ApplyServerSideEncryptionByDefault": {"SSEAlgorithm": "AES256"}}]
  }'
```

### 5. Restrict Access

Implement least-privilege access to state files:

- Use IAM roles/policies
- Enable bucket logging
- Consider using separate accounts for production state

### 6. Use State Locking

For S3, configure DynamoDB for state locking:

```bash
aws dynamodb create-table \
  --table-name terraform-state-lock \
  --attribute-definitions AttributeName=LockID,AttributeType=S \
  --key-schema AttributeName=LockID,KeyType=HASH \
  --billing-mode PAY_PER_REQUEST
```

---

## Troubleshooting

### "Backend configuration required"

**Error:** Backend type is specified but object is missing.

**Solution:** Both `backend.type` and `backend.object` labels must be specified together:

```yaml
labels:
  terraform.project-planton.org/backend.type: "s3"
  terraform.project-planton.org/backend.object: "bucket/path/state.tfstate"  # Required!
```

---

### "Access Denied" to State Bucket

**Error:** Permission denied when accessing state.

**Solutions:**
1. Verify IAM permissions on the bucket
2. Check credential configuration
3. Ensure bucket exists in the expected region
4. Verify bucket policy allows access

---

### State Lock Timeout

**Error:** Unable to acquire state lock.

**Solutions:**
1. Check if another operation is running
2. Force unlock if previous operation crashed:
   ```bash
   terraform force-unlock <LOCK_ID>
   ```
3. Verify DynamoDB table permissions (for S3 backend)

---

### Wrong Provisioner Labels

**Error:** Backend not detected, using local backend.

**Solution:** Ensure labels match the provisioner:

```yaml
# For Terraform
project-planton.org/provisioner: terraform
terraform.project-planton.org/backend.type: "s3"
terraform.project-planton.org/backend.object: "..."

# For OpenTofu
project-planton.org/provisioner: tofu
tofu.project-planton.org/backend.type: "s3"
tofu.project-planton.org/backend.object: "..."
```

---

## Related Documentation

- [Unified Commands](/docs/cli/unified-commands) - Using apply, destroy, init, plan, refresh
- [Credentials Guide](/docs/guides/credentials) - Setting up cloud provider credentials
- [Kustomize Integration](/docs/guides/kustomize) - Multi-environment deployments
- [OpenTofu Commands](/docs/cli/tofu-commands) - OpenTofu-specific details
- [CLI Reference](/docs/cli/cli-reference) - Complete command reference

---

## Getting Help

**Questions?** [GitHub Discussions](https://github.com/plantonhq/project-planton/discussions)

**Found a bug?** [Open an issue](https://github.com/plantonhq/project-planton/issues)
