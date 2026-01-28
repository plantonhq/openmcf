# DigitalOcean Spaces Bucket Examples

This document provides comprehensive examples for creating and managing DigitalOcean Spaces buckets using OpenMCF's declarative API.

## Table of Contents

- [Minimal Private Bucket](#minimal-private-bucket)
- [Public Bucket for Static Website](#public-bucket-for-static-website)
- [Versioned Bucket for Backups](#versioned-bucket-for-backups)
- [Production Media Bucket](#production-media-bucket-with-tags)
- [Multi-Region Setup](#multi-region-deployment-strategy)
- [Development vs Production](#development-vs-production-configuration)

---

## Minimal Private Bucket

The simplest configuration with only required fields - creates a private bucket for application data.

```yaml
apiVersion: digital-ocean.openmcf.org/v1
kind: DigitalOceanBucket
metadata:
  name: app-data
spec:
  bucketName: app-data
  region: nyc3
```

**Key characteristics:**
- **Access**: Private (default)
- **Versioning**: Disabled (default)
- **Cost**: $5/month (includes 250GB storage, 1TB transfer)
- **Use case**: Application data, user uploads, temporary files

**Accessing the bucket:**
```bash
# Using AWS CLI with Spaces endpoint
aws --endpoint-url https://nyc3.digitaloceanspaces.com s3 ls s3://app-data/

# Upload a file
aws --endpoint-url https://nyc3.digitaloceanspaces.com s3 cp myfile.txt s3://app-data/
```

---

## Public Bucket for Static Website

Deploy a public bucket for hosting a static website with CDN acceleration.

```yaml
apiVersion: digital-ocean.openmcf.org/v1
kind: DigitalOceanBucket
metadata:
  name: company-website
  labels:
    environment: production
    purpose: static-website
spec:
  bucketName: company-website
  region: nyc3
  accessControl: PUBLIC_READ
  tags:
    - website
    - cdn-enabled
    - public
```

**Key characteristics:**
- **Access**: Public read (anyone can access objects via URL)
- **CDN**: Automatically enabled for public buckets
- **Cost**: $5/month base (additional if > 250GB or 1TB transfer)
- **Use case**: Static websites, public downloads, open-source packages

**Accessing content:**
```
https://company-website.nyc3.digitaloceanspaces.com/index.html
https://company-website.nyc3.cdn.digitaloceanspaces.com/index.html (CDN)
```

**Upload your site:**
```bash
# Upload entire directory
aws --endpoint-url https://nyc3.digitaloceanspaces.com s3 sync ./dist s3://company-website/ --acl public-read
```

---

## Versioned Bucket for Backups

Enable versioning to protect critical data from accidental deletion or overwrites.

```yaml
apiVersion: digital-ocean.openmcf.org/v1
kind: DigitalOceanBucket
metadata:
  name: database-backups
  labels:
    environment: production
    data-classification: critical
spec:
  bucketName: prod-db-backups
  region: sfo3
  accessControl: PRIVATE
  versioningEnabled: true
  tags:
    - backups
    - versioned
    - retention-30days
```

**Key characteristics:**
- **Access**: Private (backups should not be public)
- **Versioning**: Enabled (protects against accidental deletes)
- **Region**: SFO3 (West Coast for low-latency backups)
- **Use case**: Database backups, critical document storage, disaster recovery

**Important notes:**
- Versioning **cannot be disabled** once enabled (only suspended)
- Storage costs increase as versions accumulate
- Use lifecycle policies to expire old versions (configure via S3 API)

**Example backup workflow:**
```bash
# Daily backup (creates new version)
pg_dump mydb | gzip | aws --endpoint-url https://sfo3.digitaloceanspaces.com s3 cp - s3://prod-db-backups/daily/$(date +%Y%m%d).sql.gz

# Restore from specific version
aws --endpoint-url https://sfo3.digitaloceanspaces.com s3api list-object-versions --bucket prod-db-backups --prefix daily/
aws --endpoint-url https://sfo3.digitaloceanspaces.com s3api get-object --bucket prod-db-backups --key daily/20250101.sql.gz --version-id <VERSION_ID> backup.sql.gz
```

---

## Production Media Bucket with Tags

A production-ready configuration for user-generated media with comprehensive tagging.

```yaml
apiVersion: digital-ocean.openmcf.org/v1
kind: DigitalOceanBucket
metadata:
  name: prod-user-media
  labels:
    environment: production
    team: platform
    cost-center: engineering
    compliance: gdpr
spec:
  bucketName: prod-user-media
  region: ams3
  accessControl: PUBLIC_READ
  versioningEnabled: false
  tags:
    - production
    - user-media
    - cdn-enabled
    - europe
    - gdpr-compliant
```

**Key characteristics:**
- **Access**: Public read (user profile pictures, uploaded images)
- **Region**: Amsterdam (EU region for GDPR compliance)
- **Tags**: Comprehensive tagging for organization and cost allocation
- **Use case**: User uploads, profile pictures, product images

**Integration example** (application code):
```javascript
// Upload user profile picture
const AWS = require('aws-sdk');
const spacesEndpoint = new AWS.Endpoint('ams3.digitaloceanspaces.com');
const s3 = new AWS.S3({
  endpoint: spacesEndpoint,
  accessKeyId: process.env.SPACES_KEY,
  secretAccessKey: process.env.SPACES_SECRET
});

await s3.upload({
  Bucket: 'prod-user-media',
  Key: `users/${userId}/profile.jpg`,
  Body: imageBuffer,
  ACL: 'public-read',
  ContentType: 'image/jpeg'
}).promise();

// Access via CDN
const cdnUrl = `https://prod-user-media.ams3.cdn.digitaloceanspaces.com/users/${userId}/profile.jpg`;
```

---

## Multi-Region Deployment Strategy

Deploy buckets in multiple regions for redundancy and global performance.

### US East (Primary)

```yaml
apiVersion: digital-ocean.openmcf.org/v1
kind: DigitalOceanBucket
metadata:
  name: app-media-us-east
  labels:
    region-group: us
    region-role: primary
spec:
  bucketName: app-media-us-east
  region: nyc3
  accessControl: PUBLIC_READ
  versioningEnabled: true
  tags:
    - production
    - us-east
    - primary
```

### US West (Secondary)

```yaml
apiVersion: digital-ocean.openmcf.org/v1
kind: DigitalOceanBucket
metadata:
  name: app-media-us-west
  labels:
    region-group: us
    region-role: secondary
spec:
  bucketName: app-media-us-west
  region: sfo3
  accessControl: PUBLIC_READ
  versioningEnabled: true
  tags:
    - production
    - us-west
    - secondary
```

### Europe

```yaml
apiVersion: digital-ocean.openmcf.org/v1
kind: DigitalOceanBucket
metadata:
  name: app-media-europe
  labels:
    region-group: europe
spec:
  bucketName: app-media-europe
  region: ams3
  accessControl: PUBLIC_READ
  versioningEnabled: true
  tags:
    - production
    - europe
    - gdpr
```

**Strategy:**
- **Geo-routing**: Application logic routes users to nearest bucket based on location
- **Replication**: Use `rclone sync` or custom scripts to replicate content across regions
- **CDN**: Each bucket has its own CDN endpoint for edge caching

**Note**: DigitalOcean Spaces doesn't have native cross-region replication. Use external tools.

---

## Development vs Production Configuration

### Development Environment

```yaml
apiVersion: digital-ocean.openmcf.org/v1
kind: DigitalOceanBucket
metadata:
  name: dev-app-uploads
  labels:
    environment: development
spec:
  bucketName: dev-app-uploads
  region: nyc3
  accessControl: PRIVATE
  versioningEnabled: false
  tags:
    - development
    - temporary
```

**Characteristics:**
- Private access (test data shouldn't be public)
- No versioning (saves costs for throwaway data)
- Minimal configuration
- Can be destroyed and recreated easily

### Staging Environment

```yaml
apiVersion: digital-ocean.openmcf.org/v1
kind: DigitalOceanBucket
metadata:
  name: staging-app-uploads
  labels:
    environment: staging
spec:
  bucketName: staging-app-uploads
  region: nyc3
  accessControl: PRIVATE
  versioningEnabled: true
  tags:
    - staging
    - production-like
```

**Characteristics:**
- Production-like configuration (versioning enabled)
- Private access (staging data is sensitive)
- Used for testing backup/restore workflows

### Production Environment

```yaml
apiVersion: digital-ocean.openmcf.org/v1
kind: DigitalOceanBucket
metadata:
  name: prod-app-uploads
  labels:
    environment: production
    criticality: high
    backup-enabled: "true"
spec:
  bucketName: prod-app-uploads
  region: sfo3
  accessControl: PRIVATE
  versioningEnabled: true
  tags:
    - production
    - critical
    - versioned
    - backup-enabled
```

**Characteristics:**
- Versioning enabled for data protection
- Comprehensive tagging for governance
- Private access with controlled signed URLs
- Production region selection based on user base

---

## Best Practices Summary

### Naming Conventions
- Include environment: `dev-`, `staging-`, `prod-`
- Include purpose: `-media`, `-backups`, `-logs`
- Example: `prod-user-media`, `staging-app-backups`

### Access Control
- Default to `PRIVATE` for security
- Use `PUBLIC_READ` only for truly public content (websites, downloads)
- Generate signed URLs for temporary public access

### Versioning Strategy
- **Enable for**: Backups, archives, critical data
- **Skip for**: Logs, temporary files, cached content
- **Remember**: Cannot disable once enabled

### Tagging Strategy
Use consistent tags across all buckets:
- `environment`: dev, staging, production
- `purpose`: media, backups, logs, static-site
- `cdn-enabled`: true, false
- `team`: owning team name
- `cost-center`: for billing allocation

### Region Selection
- Choose region closest to primary users
- Consider data residency requirements (EU data in EU regions)
- Plan multi-region strategy for global applications

---

## Next Steps

- Review [README.md](./README.md) for detailed field descriptions and best practices
- Check [iac/pulumi/README.md](./iac/pulumi/README.md) for Pulumi deployment
- Check [iac/tf/README.md](./iac/tf/README.md) for Terraform deployment
- Explore [DigitalOcean Spaces documentation](https://docs.digitalocean.com/products/spaces/)

