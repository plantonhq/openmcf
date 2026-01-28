# CivoBucket Examples

This document provides real-world examples of `CivoBucket` configurations for common use cases. Each example includes the YAML manifest and explains when and why to use that configuration.

## Table of Contents

1. [Minimal Development Bucket](#1-minimal-development-bucket)
2. [Production Backups with Versioning](#2-production-backups-with-versioning)
3. [Multi-Environment Setup](#3-multi-environment-setup)
4. [Public Assets Bucket](#4-public-assets-bucket)
5. [Application Logs Storage](#5-application-logs-storage)
6. [Machine Learning Dataset Storage](#6-machine-learning-dataset-storage)
7. [Multi-Region Deployment](#7-multi-region-deployment)

---

## 1. Minimal Development Bucket

**Use Case**: Quick development/testing environment where versioning and tagging aren't critical.

**Scenario**: A developer needs a temporary bucket to test S3 integration in a microservice.

```yaml
apiVersion: civo.openmcf.org/v1
kind: CivoBucket
metadata:
  name: myapp-dev-scratch
  description: Development scratch bucket for testing
spec:
  bucketName: myapp-dev-scratch
  region: FRA1
```

**Key Points**:
- No versioning (default: `false`)
- No tags
- Frankfurt region (European developers)
- Private by default
- Minimal configuration for rapid provisioning

**Post-Deployment**:
```bash
# Get credentials from stack outputs
export AWS_ACCESS_KEY_ID=<access-key>
export AWS_SECRET_ACCESS_KEY=<secret-key>

# Test upload
echo "Hello Civo" > test.txt
aws s3 cp test.txt s3://myapp-dev-scratch/ --endpoint-url https://objectstore.civo.com
```

---

## 2. Production Backups with Versioning

**Use Case**: Critical production database backups requiring versioning for data protection.

**Scenario**: PostgreSQL nightly backups with 180-day retention, versioning to recover from accidental overwrites.

```yaml
apiVersion: civo.openmcf.org/v1
kind: CivoBucket
metadata:
  name: acme-prod-db-backups
  description: Production PostgreSQL backups for ACME Corp
spec:
  bucketName: acme-prod-db-backups
  region: LON1
  versioningEnabled: true
  tags:
    - env:prod
    - criticality:high
    - retention:180-days
    - compliance:required
```

**Key Points**:
- `versioningEnabled: true` - Protects against accidental deletion/overwrite
- London region (UK data residency)
- Tags for governance and compliance tracking
- Retention policy indicated in tags (apply lifecycle rule via S3 API)

**Post-Deployment**:
```bash
# Configure versioning via S3 API
aws s3api put-bucket-versioning \
  --bucket acme-prod-db-backups \
  --versioning-configuration Status=Enabled \
  --endpoint-url https://objectstore.civo.com

# Set lifecycle rule to expire old versions after 30 days
cat > lifecycle.json <<EOF
{
  "Rules": [
    {
      "Id": "expire-old-versions",
      "Status": "Enabled",
      "NoncurrentVersionExpiration": {
        "NoncurrentDays": 30
      }
    },
    {
      "Id": "delete-old-backups",
      "Status": "Enabled",
      "Expiration": {
        "Days": 180
      }
    }
  ]
}
EOF

aws s3api put-bucket-lifecycle-configuration \
  --bucket acme-prod-db-backups \
  --lifecycle-configuration file://lifecycle.json \
  --endpoint-url https://objectstore.civo.com
```

---

## 3. Multi-Environment Setup

**Use Case**: Separate buckets for dev, staging, and production environments.

**Scenario**: A SaaS application needs isolated storage for each environment.

### Development

```yaml
apiVersion: civo.openmcf.org/v1
kind: CivoBucket
metadata:
  name: saasapp-dev-storage
  description: Development environment storage
spec:
  bucketName: saasapp-dev-storage
  region: NYC1
  versioningEnabled: false
  tags:
    - env:dev
    - team:backend
    - cost-center:engineering
```

### Staging

```yaml
apiVersion: civo.openmcf.org/v1
kind: CivoBucket
metadata:
  name: saasapp-staging-storage
  description: Staging environment storage
spec:
  bucketName: saasapp-staging-storage
  region: NYC1
  versioningEnabled: true
  tags:
    - env:staging
    - team:backend
    - cost-center:engineering
```

### Production

```yaml
apiVersion: civo.openmcf.org/v1
kind: CivoBucket
metadata:
  name: saasapp-prod-storage
  description: Production environment storage
spec:
  bucketName: saasapp-prod-storage
  region: NYC1
  versioningEnabled: true
  tags:
    - env:prod
    - team:backend
    - criticality:high
    - cost-center:product
```

**Key Points**:
- Same region (NYC1) for consistent latency
- Dev: No versioning (cost optimization)
- Staging/Prod: Versioning enabled (data protection)
- Tags differentiate environments for cost tracking and access policies

---

## 4. Public Assets Bucket

**Use Case**: Serving static website assets (images, CSS, JavaScript) publicly.

**Scenario**: Marketing website needs a CDN-backed bucket for static assets.

```yaml
apiVersion: civo.openmcf.org/v1
kind: CivoBucket
metadata:
  name: marketing-static-assets
  description: Public assets for marketing website
spec:
  bucketName: marketing-static-assets
  region: LON1
  versioningEnabled: false
  tags:
    - env:prod
    - public:true
    - team:marketing
```

**Post-Deployment** (Configure Public Access):
```bash
# Make bucket publicly readable
cat > public-policy.json <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Sid": "PublicRead",
      "Effect": "Allow",
      "Principal": "*",
      "Action": "s3:GetObject",
      "Resource": "arn:aws:s3:::marketing-static-assets/*"
    }
  ]
}
EOF

aws s3api put-bucket-policy \
  --bucket marketing-static-assets \
  --policy file://public-policy.json \
  --endpoint-url https://objectstore.civo.com

# Upload assets
aws s3 sync ./website-assets/ s3://marketing-static-assets/ \
  --acl public-read \
  --endpoint-url https://objectstore.civo.com
```

**Key Points**:
- No versioning (assets are immutable, use cache-busting filenames)
- `public:true` tag signals public access policy
- Post-deployment S3 policy for public reads
- Ideal for static websites, CDN origins, or public downloads

---

## 5. Application Logs Storage

**Use Case**: Centralized log storage with automatic expiration.

**Scenario**: Kubernetes application logs aggregated to object storage, retained for 30 days.

```yaml
apiVersion: civo.openmcf.org/v1
kind: CivoBucket
metadata:
  name: app-logs-archive
  description: Centralized application logs with 30-day retention
spec:
  bucketName: app-logs-archive
  region: FRA1
  versioningEnabled: false
  tags:
    - env:prod
    - data-type:logs
    - retention:30-days
    - team:devops
```

**Post-Deployment** (Configure Lifecycle):
```bash
# Auto-delete logs after 30 days
cat > logs-lifecycle.json <<EOF
{
  "Rules": [
    {
      "Id": "expire-old-logs",
      "Status": "Enabled",
      "Expiration": {
        "Days": 30
      }
    }
  ]
}
EOF

aws s3api put-bucket-lifecycle-configuration \
  --bucket app-logs-archive \
  --lifecycle-configuration file://logs-lifecycle.json \
  --endpoint-url https://objectstore.civo.com
```

**Key Points**:
- No versioning (logs are immutable, write-once)
- Lifecycle policy for automatic expiration (cost optimization)
- Frankfurt region (compliance with European data residency)
- Tag-based retention policy documentation

---

## 6. Machine Learning Dataset Storage

**Use Case**: Large datasets for ML training with versioning for reproducibility.

**Scenario**: Data science team needs versioned datasets for model training and experiments.

```yaml
apiVersion: civo.openmcf.org/v1
kind: CivoBucket
metadata:
  name: ml-training-datasets
  description: Versioned ML training datasets
spec:
  bucketName: ml-training-datasets
  region: NYC1
  versioningEnabled: true
  tags:
    - env:prod
    - team:data-science
    - data-type:ml-datasets
    - criticality:medium
```

**Usage Pattern**:
```bash
# Upload dataset with version tagging
aws s3 cp customer-churn-v2.csv s3://ml-training-datasets/datasets/ \
  --metadata version=v2.0,commit=abc123 \
  --endpoint-url https://objectstore.civo.com

# Retrieve specific version for reproducibility
aws s3api get-object \
  --bucket ml-training-datasets \
  --key datasets/customer-churn-v2.csv \
  --version-id <version-id> \
  --endpoint-url https://objectstore.civo.com \
  local-copy.csv
```

**Key Points**:
- Versioning for dataset reproducibility (critical for ML experiments)
- NYC region (collocated with GPU-enabled Civo Kubernetes cluster)
- Large dataset support (no per-request fees)
- Tag-based organization for dataset lineage

---

## 7. Multi-Region Deployment

**Use Case**: Disaster recovery with cross-region replication.

**Scenario**: Critical production data replicated from London (primary) to Frankfurt (backup).

### Primary (London)

```yaml
apiVersion: civo.openmcf.org/v1
kind: CivoBucket
metadata:
  name: acme-prod-primary
  description: Primary production bucket (London)
spec:
  bucketName: acme-prod-primary
  region: LON1
  versioningEnabled: true
  tags:
    - env:prod
    - role:primary
    - dr-region:FRA1
    - criticality:critical
```

### Backup (Frankfurt)

```yaml
apiVersion: civo.openmcf.org/v1
kind: CivoBucket
metadata:
  name: acme-prod-backup
  description: Backup production bucket (Frankfurt)
spec:
  bucketName: acme-prod-backup
  region: FRA1
  versioningEnabled: true
  tags:
    - env:prod
    - role:backup
    - dr-region:LON1
    - criticality:critical
```

**Replication Setup**:
```bash
# Use rclone for cross-region replication (Civo offers free inter-region transfer)
rclone sync civo-lon1:acme-prod-primary civo-fra1:acme-prod-backup \
  --transfers 8 \
  --checkers 16 \
  --log-file replication.log
```

**Key Points**:
- Both buckets have versioning for data integrity
- Tags document disaster recovery pairing
- Civo's free inter-region transfer makes replication cost-effective
- Use scheduled jobs (Kubernetes CronJob) for continuous replication

---

## Best Practices Summary

1. **Naming Convention**: Use consistent naming patterns (`<app>-<env>-<purpose>`)
2. **Region Selection**: Collocate with compute resources for performance
3. **Versioning**: Enable for critical data (backups, datasets), disable for immutable data (logs, assets)
4. **Tags**: Use for environment tracking, cost allocation, and compliance
5. **Lifecycle Policies**: Configure via S3 API for automatic data expiration
6. **Security**: Default to private, only expose specific prefixes/objects when necessary
7. **Cost Optimization**: Compress data, use lifecycle rules, right-size allocations

---

## Testing Your Bucket

After provisioning any bucket, verify it works:

```bash
# Set credentials (from stack outputs)
export AWS_ACCESS_KEY_ID=<access-key>
export AWS_SECRET_ACCESS_KEY=<secret-key>
export ENDPOINT_URL=https://objectstore.civo.com

# Test write
echo "Test $(date)" > test.txt
aws s3 cp test.txt s3://<bucket-name>/ --endpoint-url $ENDPOINT_URL

# Test read
aws s3 cp s3://<bucket-name>/test.txt - --endpoint-url $ENDPOINT_URL

# Test list
aws s3 ls s3://<bucket-name>/ --endpoint-url $ENDPOINT_URL

# Clean up
aws s3 rm s3://<bucket-name>/test.txt --endpoint-url $ENDPOINT_URL
```

---

## Advanced Configurations

For advanced S3 features not exposed in the Protobuf API, use the AWS CLI or SDK after provisioning:

### Enable CORS

```bash
cat > cors.json <<EOF
{
  "CORSRules": [
    {
      "AllowedOrigins": ["https://example.com"],
      "AllowedMethods": ["GET", "PUT"],
      "AllowedHeaders": ["*"],
      "MaxAgeSeconds": 3000
    }
  ]
}
EOF

aws s3api put-bucket-cors \
  --bucket <bucket-name> \
  --cors-configuration file://cors.json \
  --endpoint-url https://objectstore.civo.com
```

### Enable Logging

```bash
# Create a separate bucket for logs
aws s3api put-bucket-logging \
  --bucket <source-bucket> \
  --bucket-logging-status '{
    "LoggingEnabled": {
      "TargetBucket": "<logs-bucket>",
      "TargetPrefix": "logs/"
    }
  }' \
  --endpoint-url https://objectstore.civo.com
```

### Generate Presigned URLs

```python
import boto3
from botocore.client import Config

s3 = boto3.client(
    's3',
    endpoint_url='https://objectstore.civo.com',
    aws_access_key_id='<access-key>',
    aws_secret_access_key='<secret-key>',
    config=Config(signature_version='s3v4'),
)

# Generate presigned URL (valid for 1 hour)
url = s3.generate_presigned_url(
    'get_object',
    Params={'Bucket': 'my-bucket', 'Key': 'file.pdf'},
    ExpiresIn=3600
)
print(url)
```

---

## Related Documentation

- **API Reference**: [README.md](./README.md)
- **Research**: [docs/README.md](./docs/README.md)
- **Pulumi Module**: [iac/pulumi/](./iac/pulumi/)
- **Civo Documentation**: [Civo Object Storage](https://www.civo.com/docs/object-stores)

