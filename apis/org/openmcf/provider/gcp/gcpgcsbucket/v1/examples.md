# GCP Cloud Storage Bucket Examples

This document provides practical examples of GCS bucket configurations using the `GcpGcsBucket` API resource. Each example demonstrates common deployment patterns with explanations of key design decisions.

## Table of Contents

- [Basic Examples](#basic-examples)
  - [Minimal Private Bucket](#minimal-private-bucket)
  - [Regional Bucket with Labels](#regional-bucket-with-labels)
- [Cross-Resource Reference Examples](#cross-resource-reference-examples)
  - [Using valueFrom for Project Reference](#using-valuefrom-for-project-reference)
- [Access Control Examples](#access-control-examples)
  - [Public Read Bucket](#public-read-bucket)
  - [Service Account Access](#service-account-access)
  - [Conditional IAM Access](#conditional-iam-access)
- [Data Protection Examples](#data-protection-examples)
  - [Versioned Bucket with Lifecycle](#versioned-bucket-with-lifecycle)
  - [Compliance Archive Bucket](#compliance-archive-bucket)
- [Cost Optimization Examples](#cost-optimization-examples)
  - [Development Bucket with Cleanup](#development-bucket-with-cleanup)
  - [Storage Tiering Strategy](#storage-tiering-strategy)
- [Advanced Examples](#advanced-examples)
  - [Static Website Hosting](#static-website-hosting)
  - [CORS-Enabled Bucket](#cors-enabled-bucket)
  - [CMEK-Encrypted Bucket](#cmek-encrypted-bucket)
  - [Multi-Region Production Bucket](#multi-region-production-bucket)

---

## Basic Examples

### Minimal Private Bucket

A simple private bucket with default settings. Suitable for internal application storage.

```yaml
apiVersion: gcp.openmcf.org/v1
kind: GcpGcsBucket
metadata:
  name: my-app-data-prod
spec:
  gcpProjectId:
    value: my-gcp-project-123
  location: us-east1
  bucketName: my-app-data-prod
  uniformBucketLevelAccessEnabled: true
```

**Key Points:**
- Minimal configuration using only required fields
- UBLA enabled for simplified IAM-only access control
- Defaults to STANDARD storage class
- No public access (private by default)
- Uses literal value for `gcpProjectId`

---

### Regional Bucket with Labels

A regional bucket co-located with GKE cluster, with custom labels for cost tracking.

```yaml
apiVersion: gcp.openmcf.org/v1
kind: GcpGcsBucket
metadata:
  name: gke-workload-storage
spec:
  gcpProjectId:
    value: my-gcp-project-123
  location: us-central1
  bucketName: gke-workload-storage
  uniformBucketLevelAccessEnabled: true
  storageClass: STANDARD
  gcpLabels:
    team: platform-engineering
    cost-center: engineering-prod
    application: data-processing
    environment: production
```

**Key Points:**
- Regional placement co-located with compute resources (minimize latency & egress costs)
- Custom labels for cost allocation and governance
- STANDARD class appropriate for frequently accessed data
- Labels enable filtering in billing reports

---

## Cross-Resource Reference Examples

### Using valueFrom for Project Reference

Reference a GcpProject resource instead of hardcoding the project ID. This enables dynamic dependencies between resources.

```yaml
apiVersion: gcp.openmcf.org/v1
kind: GcpGcsBucket
metadata:
  name: app-storage-prod
spec:
  gcpProjectId:
    valueFrom:
      kind: GcpProject
      name: main-project
      fieldPath: status.outputs.project_id
  location: us-east1
  bucketName: app-storage-prod
  uniformBucketLevelAccessEnabled: true
  storageClass: STANDARD
  versioningEnabled: true
```

**Key Points:**
- `gcpProjectId` references another `GcpProject` resource named "main-project"
- The `fieldPath` specifies which output field to use
- Enables infrastructure composition and dependency management
- Ideal for multi-resource deployments where project is provisioned separately

### With Environment Context

When referencing resources in a specific environment:

```yaml
apiVersion: gcp.openmcf.org/v1
kind: GcpGcsBucket
metadata:
  name: data-lake-staging
spec:
  gcpProjectId:
    valueFrom:
      kind: GcpProject
      env: staging
      name: data-platform-project
      fieldPath: status.outputs.project_id
  location: us-central1
  bucketName: data-lake-staging
  uniformBucketLevelAccessEnabled: true
  storageClass: STANDARD
```

**Key Points:**
- The `env` field specifies the environment context for the reference
- Useful in multi-environment deployments
- Ensures correct project resolution across environments

---

## Access Control Examples

### Public Read Bucket

A bucket with public read access for serving open datasets or public website assets.

```yaml
apiVersion: gcp.openmcf.org/v1
kind: GcpGcsBucket
metadata:
  name: public-open-dataset
spec:
  gcpProjectId:
    value: my-gcp-project-123
  location: US  # Multi-region for global access
  bucketName: public-open-dataset
  uniformBucketLevelAccessEnabled: true
  storageClass: STANDARD
  iamBindings:
    - role: "roles/storage.objectViewer"
      members:
        - "allUsers"
  publicAccessPrevention: "inherited"  # Allow public access
```

**Key Points:**
- Explicit IAM binding granting `objectViewer` to `allUsers`
- Multi-region location for global content delivery
- `public_access_prevention: inherited` to allow public access
- UBLA enabled (required for IAM-only access control)
- Clear, auditable access configuration

**Security Note:** Never use this pattern for sensitive data. Consider Cloud CDN + Load Balancer for production websites.

---

### Service Account Access

Grant specific service accounts access to the bucket for application workloads.

```yaml
apiVersion: gcp.openmcf.org/v1
kind: GcpGcsBucket
metadata:
  name: app-backend-storage
spec:
  gcpProjectId:
    value: my-gcp-project-123
  location: us-east1
  bucketName: app-backend-storage
  uniformBucketLevelAccessEnabled: true
  versioningEnabled: true
  iamBindings:
    # Application backend needs read/write access
    - role: "roles/storage.objectAdmin"
      members:
        - "serviceAccount:backend-api@my-gcp-project-123.iam.gserviceaccount.com"
    # Analytics pipeline needs read-only access
    - role: "roles/storage.objectViewer"
      members:
        - "serviceAccount:analytics-etl@my-gcp-project-123.iam.gserviceaccount.com"
```

**Key Points:**
- Principle of least privilege (different roles for different needs)
- Service accounts identified by full email address
- Multiple IAM bindings for different access patterns
- Versioning enabled to protect against accidental deletion

---

### Conditional IAM Access

Use IAM conditions for time-based or attribute-based access control.

```yaml
apiVersion: gcp.openmcf.org/v1
kind: GcpGcsBucket
metadata:
  name: sensitive-project-data
spec:
  gcpProjectId:
    value: my-gcp-project-123
  location: us-east1
  bucketName: sensitive-project-data
  uniformBucketLevelAccessEnabled: true
  iamBindings:
    # Grant access only during business hours
    - role: "roles/storage.objectViewer"
      members:
        - "group:contractors@example.com"
      condition: |
        request.time.getHours() >= 9 && request.time.getHours() <= 17
    # Grant access only from specific IP ranges
    - role: "roles/storage.objectViewer"
      members:
        - "group:external-auditors@example.com"
      condition: |
        origin.ip in ["203.0.113.0/24", "198.51.100.0/24"]
```

**Key Points:**
- IAM conditions enable dynamic access control
- Time-based restrictions for contractor access
- IP-based restrictions for external parties
- Conditions evaluated at request time (no additional infrastructure)

**Note:** IAM condition expressions use Common Expression Language (CEL).

---

## Data Protection Examples

### Versioned Bucket with Lifecycle

Enable versioning with lifecycle rules to prevent unbounded storage growth.

```yaml
apiVersion: gcp.openmcf.org/v1
kind: GcpGcsBucket
metadata:
  name: critical-app-data
spec:
  gcpProjectId:
    value: my-gcp-project-123
  location: us-east1
  bucketName: critical-app-data
  uniformBucketLevelAccessEnabled: true
  storageClass: STANDARD
  versioningEnabled: true
  lifecycleRules:
    # Delete noncurrent versions after 30 days
    - action:
        type: "Delete"
      condition:
        numNewerVersions: 5
    # Delete noncurrent versions older than 90 days regardless of count
    - action:
        type: "Delete"
      condition:
        ageDays: 90
        isLive: false
```

**Key Points:**
- Versioning protects against accidental deletion/overwrite
- Lifecycle rules prevent unbounded storage costs
- Keep last 5 versions OR 90 days, whichever is shorter
- Deletes noncurrent versions automatically

**Cost Impact:** Without lifecycle rules, versioning can double storage costs over time.

---

### Compliance Archive Bucket

A bucket with retention policy for regulatory compliance (e.g., FINRA, HIPAA).

```yaml
apiVersion: gcp.openmcf.org/v1
kind: GcpGcsBucket
metadata:
  name: financial-records-archive
spec:
  gcpProjectId:
    value: my-gcp-project-123
  location: us-east1
  bucketName: financial-records-archive
  uniformBucketLevelAccessEnabled: true
  storageClass: ARCHIVE
  versioningEnabled: true
  retentionPolicy:
    retentionPeriodSeconds: 220752000  # 7 years (FINRA requirement)
    isLocked: false  # Lock after initial validation
  lifecycleRules:
    # Delete objects after 7 years + 30 day buffer
    - action:
        type: "Delete"
      condition:
        ageDays: 2585  # 7 years + 30 days
  iamBindings:
    # Write-only access for record ingestion
    - role: "roles/storage.objectCreator"
      members:
        - "serviceAccount:records-ingest@my-gcp-project-123.iam.gserviceaccount.com"
    # Read-only access for compliance team
    - role: "roles/storage.objectViewer"
      members:
        - "group:compliance-team@example.com"
```

**Key Points:**
- ARCHIVE storage class for lowest storage cost
- Retention policy enforces WORM (Write Once, Read Many) compliance
- Objects cannot be deleted during 7-year retention period
- Lifecycle rule cleans up after retention expires
- Separate roles for ingestion vs. auditing

**Critical:** `is_locked: false` allows testing. Set to `true` in production (irreversible!).

---

## Cost Optimization Examples

### Development Bucket with Cleanup

Aggressive lifecycle policies for development/staging environments to minimize costs.

```yaml
apiVersion: gcp.openmcf.org/v1
kind: GcpGcsBucket
metadata:
  name: dev-ephemeral-storage
spec:
  gcpProjectId:
    value: my-gcp-project-123
  location: us-east1
  bucketName: dev-ephemeral-storage
  uniformBucketLevelAccessEnabled: true
  storageClass: STANDARD
  versioningEnabled: true
  lifecycleRules:
    # Delete all objects after 30 days
    - action:
        type: "Delete"
      condition:
        ageDays: 30
    # Delete noncurrent versions after 7 days
    - action:
        type: "Delete"
      condition:
        numNewerVersions: 2
  gcpLabels:
    environment: development
    auto-cleanup: enabled
```

**Key Points:**
- Automatic cleanup after 30 days (safe for ephemeral dev data)
- Keep only 2 versions (balance protection vs. cost)
- STANDARD class appropriate for active development
- Labels indicate auto-cleanup for cost accountability

**Cost Savings:** Can reduce storage costs by 90% compared to production buckets.

---

### Storage Tiering Strategy

Automatically transition objects to cheaper storage classes based on age.

```yaml
apiVersion: gcp.openmcf.org/v1
kind: GcpGcsBucket
metadata:
  name: tiered-backup-storage
spec:
  gcpProjectId:
    value: my-gcp-project-123
  location: us-east1
  bucketName: tiered-backup-storage
  uniformBucketLevelAccessEnabled: true
  storageClass: STANDARD  # Initial storage class for new objects
  versioningEnabled: true
  lifecycleRules:
    # Transition to NEARLINE after 30 days (infrequent access)
    - action:
        type: "SetStorageClass"
        storageClass: NEARLINE
      condition:
        ageDays: 30
        matchesStorageClass:
          - STANDARD
    # Transition to COLDLINE after 90 days (quarterly access)
    - action:
        type: "SetStorageClass"
        storageClass: COLDLINE
      condition:
        ageDays: 90
        matchesStorageClass:
          - NEARLINE
    # Transition to ARCHIVE after 365 days (yearly access)
    - action:
        type: "SetStorageClass"
        storageClass: ARCHIVE
      condition:
        ageDays: 365
        matchesStorageClass:
          - COLDLINE
    # Delete after 7 years
    - action:
        type: "Delete"
      condition:
        ageDays: 2555  # 7 years
```

**Key Points:**
- Automatic cost optimization through storage class transitions
- Matches actual access patterns (frequent → infrequent → rare → delete)
- `matches_storage_class` prevents double-transitions
- Significant cost savings for long-lived data

**Cost Comparison (per GB per month):**
- STANDARD: $0.020 | NEARLINE: $0.010 | COLDLINE: $0.004 | ARCHIVE: $0.0012

---

## Advanced Examples

### Static Website Hosting

Host a static website directly from GCS (development/internal use only).

```yaml
apiVersion: gcp.openmcf.org/v1
kind: GcpGcsBucket
metadata:
  name: static-website-dev
spec:
  gcpProjectId:
    value: my-gcp-project-123
  location: US  # Multi-region for better availability
  bucketName: static-website-dev
  uniformBucketLevelAccessEnabled: true
  storageClass: STANDARD
  website:
    mainPageSuffix: "index.html"
    notFoundPage: "404.html"
  iamBindings:
    - role: "roles/storage.objectViewer"
      members:
        - "allUsers"
```

**Key Points:**
- Website configuration enables index page and custom 404
- Public access via IAM binding
- Multi-region for availability

**Production Alternative:** Use Cloud CDN + Load Balancer for:
- HTTPS support with custom domains
- CDN caching (reduced egress costs)
- Better access logging and monitoring
- DDoS protection

---

### CORS-Enabled Bucket

Enable CORS for buckets accessed from web browsers (e.g., font hosting, direct uploads).

```yaml
apiVersion: gcp.openmcf.org/v1
kind: GcpGcsBucket
metadata:
  name: user-upload-storage
spec:
  gcpProjectId:
    value: my-gcp-project-123
  location: us-east1
  bucketName: user-upload-storage
  uniformBucketLevelAccessEnabled: true
  storageClass: STANDARD
  corsRules:
    # Allow uploads from web application
    - methods:
        - "GET"
        - "POST"
        - "PUT"
      origins:
        - "https://app.example.com"
        - "https://staging.example.com"
      responseHeaders:
        - "Content-Type"
        - "x-goog-acl"
      maxAgeSeconds: 3600  # Cache preflight for 1 hour
  iamBindings:
    # Allow authenticated users to upload
    - role: "roles/storage.objectCreator"
      members:
        - "allAuthenticatedUsers"
```

**Key Points:**
- CORS enables direct browser uploads (no backend proxy)
- Restrict origins to your application domains
- `max_age_seconds` reduces preflight requests
- `objectCreator` allows uploads but not listing/deletion

**Security:** Use signed URLs or Workload Identity Federation for production instead of `allAuthenticatedUsers`.

---

### CMEK-Encrypted Bucket

Use customer-managed encryption keys for compliance requirements.

```yaml
apiVersion: gcp.openmcf.org/v1
kind: GcpGcsBucket
metadata:
  name: highly-sensitive-data
spec:
  gcpProjectId:
    value: my-gcp-project-123
  location: us-east1
  bucketName: highly-sensitive-data
  uniformBucketLevelAccessEnabled: true
  storageClass: STANDARD
  versioningEnabled: true
  encryption:
    kmsKeyName: "projects/my-gcp-project-123/locations/us-east1/keyRings/production-keys/cryptoKeys/bucket-encryption-key"
  publicAccessPrevention: "enforced"  # Prevent accidental public access
  gcpLabels:
    data-classification: highly-sensitive
    encryption: cmek
    compliance: sox-hipaa
```

**Key Points:**
- CMEK provides audit trail of key usage via Cloud KMS
- Key must exist before bucket creation (module doesn't create it)
- Key location must match or encompass bucket location
- Public access prevention enforced for extra security

**Prerequisites:**
1. Create KMS key ring and key
2. Grant GCS service account `roles/cloudkms.cryptoKeyEncrypterDecrypter` on the key

**Service Account Format:** `service-PROJECT_NUMBER@gs-project-accounts.iam.gserviceaccount.com`

---

### Multi-Region Production Bucket

High-availability bucket for globally distributed applications.

```yaml
apiVersion: gcp.openmcf.org/v1
kind: GcpGcsBucket
metadata:
  name: global-cdn-content
spec:
  gcpProjectId:
    value: my-gcp-project-123
  location: US  # Multi-region (auto-replication across US regions)
  bucketName: global-cdn-content
  uniformBucketLevelAccessEnabled: true
  storageClass: STANDARD
  versioningEnabled: true
  lifecycleRules:
    # Keep last 10 versions
    - action:
        type: "Delete"
      condition:
        numNewerVersions: 10
    # Delete noncurrent versions after 90 days
    - action:
        type: "Delete"
      condition:
        ageDays: 90
        isLive: false
  iamBindings:
    # CDN origin fetch access
    - role: "roles/storage.objectViewer"
      members:
        - "serviceAccount:cdn-backend@my-gcp-project-123.iam.gserviceaccount.com"
    # Content management team write access
    - role: "roles/storage.objectAdmin"
      members:
        - "group:content-managers@example.com"
  gcpLabels:
    application: cdn
    tier: production
    availability: multi-region
```

**Key Points:**
- Multi-region for automatic cross-region replication
- Versioning protects against accidental content changes
- Lifecycle rules prevent unbounded version growth
- Separate read/write access for CDN vs. management

**Cost Consideration:** Multi-region storage costs ~20% more than regional, but eliminates single-region failure scenarios.

---

## Best Practices Summary

### Security
1. ✅ Always enable UBLA (`uniform_bucket_level_access_enabled: true`)
2. ✅ Use explicit IAM bindings instead of legacy ACLs
3. ✅ Set `public_access_prevention: "enforced"` for private buckets
4. ✅ Use service accounts with minimal permissions
5. ✅ Enable versioning for critical data

### Cost Optimization
1. ✅ Always configure lifecycle rules to prevent unbounded storage growth
2. ✅ Use storage class transitions for infrequently accessed data
3. ✅ Delete noncurrent versions automatically
4. ✅ Use regional buckets co-located with compute resources
5. ✅ Leverage labels for cost tracking and accountability

### Reliability
1. ✅ Enable versioning for important data
2. ✅ Use regional buckets for latency-sensitive workloads
3. ✅ Use multi-region buckets for high availability
4. ✅ Implement retention policies for compliance requirements
5. ✅ Use CMEK for regulatory compliance

### Avoid These Anti-Patterns
1. ❌ Using boolean "is_public" flags (unclear, not auditable)
2. ❌ Disabling UBLA (complex dual-permission model)
3. ❌ Enabling versioning without lifecycle cleanup
4. ❌ Using ARCHIVE class without understanding minimum storage duration
5. ❌ Multi-region placement when data is only accessed from one region

---

## Further Reading

- [Component README](README.md) - Overview and features
- [Pulumi Module README](iac/pulumi/README.md) - Implementation details
- [Pulumi Module Overview](iac/pulumi/overview.md) - Architecture and design
- [Research Document](docs/README.md) - Comprehensive analysis of GCS patterns
- [GCS Best Practices (Google Cloud)](https://cloud.google.com/storage/docs/best-practices)
- [Storage Classes Documentation](https://cloud.google.com/storage/docs/storage-classes)
- [Lifecycle Management Guide](https://cloud.google.com/storage/docs/lifecycle)

---

## Need Help?

For additional examples or specific use cases, consult:
1. The [research document](docs/README.md) for design rationale
2. The [audit report](docs/audit/) for component completeness
3. The OpenMCF community forums
