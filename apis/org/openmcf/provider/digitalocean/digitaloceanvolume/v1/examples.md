# DigitalOcean Volume Examples

This document provides practical examples for creating DigitalOcean Block Storage Volumes using OpenMCF's manifest-based approach.

## Table of Contents

1. [Minimal Volume (Development)](#1-minimal-volume-development)
2. [Production Database Volume with XFS](#2-production-database-volume-with-xfs)
3. [General Purpose Volume with ext4](#3-general-purpose-volume-with-ext4)
4. [Volume from Snapshot (Restore/Clone)](#4-volume-from-snapshot-restoreclone)
5. [Large Storage Volume (Multi-TB)](#5-large-storage-volume-multi-tb)

---

## 1. Minimal Volume (Development)

**Use Case:** Simple development volume with minimal configuration.

**Features:**
- Small size (10 GiB) for testing
- No pre-formatting (manual filesystem creation)
- Basic tagging

```yaml
apiVersion: digital-ocean.openmcf.org/v1
kind: DigitalOceanVolume
metadata:
  name: dev-test-volume
spec:
  volume_name: dev-test-volume
  region: nyc3
  size_gib: 10
  filesystem_type: NONE  # No pre-formatting
  tags:
    - env:dev
    - purpose:testing
```

**Deploy:**
```bash
planton pulumi up --manifest dev-test-volume.yaml
```

**Notes:**
- `filesystem_type: NONE` means the volume is unformatted
- You'll need to manually format and mount the volume after attaching to a Droplet
- Suitable for development and learning

---

## 2. Production Database Volume with XFS

**Use Case:** Production PostgreSQL or MySQL database storage with optimal filesystem.

**Features:**
- XFS filesystem (best for databases)
- Large capacity (500 GiB)
- Production tags
- Descriptive documentation

```yaml
apiVersion: digital-ocean.openmcf.org/v1
kind: DigitalOceanVolume
metadata:
  name: prod-db-data
  labels:
    app: postgres
    environment: production
spec:
  volume_name: prod-pg-data
  description: "PostgreSQL production database data volume"
  region: sfo3
  size_gib: 500
  filesystem_type: XFS  # Optimal for databases
  tags:
    - env:prod
    - service:postgres
    - tier:database
    - backup:daily
```

**Deploy:**
```bash
planton pulumi up --manifest prod-db-volume.yaml
```

**Why XFS for Databases:**
- Superior performance for large files and high I/O workloads
- Efficient online resizing with `xfs_growfs`
- Better concurrent write handling
- Industry standard for PostgreSQL and MySQL

**Mounting on Droplet:**
```bash
# XFS volumes come pre-formatted by DigitalOcean
mkdir -p /mnt/postgres-data
mount -o defaults,nofail,discard,noatime /dev/disk/by-id/scsi-0DO_Volume_prod-pg-data /mnt/postgres-data

# Add to /etc/fstab for persistence
echo "/dev/disk/by-id/scsi-0DO_Volume_prod-pg-data /mnt/postgres-data xfs defaults,nofail,discard,noatime 0 0" >> /etc/fstab

# Set ownership for PostgreSQL
chown -R postgres:postgres /mnt/postgres-data
```

---

## 3. General Purpose Volume with ext4

**Use Case:** Application data, file storage, general-purpose workloads.

**Features:**
- ext4 filesystem (stable, widely compatible)
- Medium size (100 GiB)
- Application-specific tags

```yaml
apiVersion: digital-ocean.openmcf.org/v1
kind: DigitalOceanVolume
metadata:
  name: app-data-volume
spec:
  volume_name: app-data-vol
  description: "Application data storage volume"
  region: nyc1
  size_gib: 100
  filesystem_type: EXT4  # General purpose filesystem
  tags:
    - env:staging
    - app:web-app
    - data:persistent
```

**Deploy:**
```bash
planton pulumi up --manifest app-data-volume.yaml
```

**Why ext4 for General Purpose:**
- Mature and stable
- Excellent for many small files
- Great compatibility with tools and utilities
- Lower CPU overhead than XFS for mixed workloads

**Mounting on Droplet:**
```bash
# ext4 volumes come pre-formatted by DigitalOcean
mkdir -p /mnt/app-data
mount -o defaults,nofail,discard /dev/disk/by-id/scsi-0DO_Volume_app-data-vol /mnt/app-data

# Add to /etc/fstab
echo "/dev/disk/by-id/scsi-0DO_Volume_app-data-vol /mnt/app-data ext4 defaults,nofail,discard 0 0" >> /etc/fstab
```

---

## 4. Volume from Snapshot (Restore/Clone)

**Use Case:** Restore from backup or clone an existing volume.

**Features:**
- Created from existing snapshot
- Inherits filesystem type from snapshot
- Can be larger than source snapshot (for growth)

```yaml
apiVersion: digital-ocean.openmcf.org/v1
kind: DigitalOceanVolume
metadata:
  name: restored-volume
spec:
  volume_name: restored-db-data
  description: "Database volume restored from nightly backup"
  region: nyc3  # Must match snapshot's region
  size_gib: 600  # Can be larger than snapshot (500 GiB), but not smaller
  filesystem_type: XFS  # Should match original volume's filesystem
  snapshot_id: "123456789"  # Replace with actual snapshot ID
  tags:
    - env:prod
    - restored:true
    - source:nightly-backup
```

**Deploy:**
```bash
planton pulumi up --manifest restored-volume.yaml
```

**Important Notes:**
- **Region Lock**: Snapshot and new volume must be in the same region
- **Size Requirement**: Volume must be ≥ snapshot size
- **Filesystem**: Inherits from snapshot, but specify for clarity
- **Snapshots are regional**: Cannot copy between regions (DR limitation)

**Getting Snapshot ID:**
```bash
# List snapshots
doctl compute volume-snapshot list

# Or use DigitalOcean Control Panel
```

---

## 5. Large Storage Volume (Multi-TB)

**Use Case:** Large-scale data storage, archives, media files.

**Features:**
- Maximum practical size (8 TB = 8192 GiB)
- XFS for large file performance
- Archive-specific tags

```yaml
apiVersion: digital-ocean.openmcf.org/v1
kind: DigitalOceanVolume
metadata:
  name: archive-storage
spec:
  volume_name: prod-archive-vol
  description: "Large archive storage for media and backups"
  region: fra1
  size_gib: 8192  # 8 TB (max practical size)
  filesystem_type: XFS  # Best for large files
  tags:
    - env:prod
    - purpose:archive
    - retention:long-term
```

**Deploy:**
```bash
planton pulumi up --manifest archive-storage.yaml
```

**Cost Considerations:**
- DigitalOcean charges $0.10/GB/month
- 8 TB volume = ~$820/month
- Consider object storage (DigitalOcean Spaces) for cheaper archival

**Performance Notes:**
- Single volume IOPS limits apply regardless of size
- For very high I/O needs, consider multiple smaller volumes
- Maximum size is 16,000 GiB (16 TB), but 8 TB is more practical

---

## Advanced Patterns

### Multi-Environment Volumes

Create consistent volumes across dev/staging/prod with different sizes:

```yaml
# dev-volume.yaml
apiVersion: digital-ocean.openmcf.org/v1
kind: DigitalOceanVolume
metadata:
  name: dev-db-volume
spec:
  volume_name: dev-db-data
  region: nyc3
  size_gib: 10  # Small for dev
  filesystem_type: XFS
  tags: ["env:dev"]
---
# staging-volume.yaml
apiVersion: digital-ocean.openmcf.org/v1
kind: DigitalOceanVolume
metadata:
  name: staging-db-volume
spec:
  volume_name: staging-db-data
  region: nyc3
  size_gib: 100  # Medium for staging
  filesystem_type: XFS
  tags: ["env:staging"]
---
# prod-volume.yaml
apiVersion: digital-ocean.openmcf.org/v1
kind: DigitalOceanVolume
metadata:
  name: prod-db-volume
spec:
  volume_name: prod-db-data
  region: nyc3
  size_gib: 500  # Large for prod
  filesystem_type: XFS
  tags: ["env:prod"]
```

### Volume Resizing Workflow

Resizing a volume is a two-step process:

**Step 1: Update manifest and apply**
```yaml
spec:
  size_gib: 1000  # Increased from 500
```

```bash
planton pulumi up --manifest volume.yaml
# Volume is resized at API level
```

**Step 2: Expand filesystem (manual)**
```bash
# SSH into the Droplet
ssh droplet-ip

# For XFS
sudo xfs_growfs /mnt/mount-point

# For ext4
sudo resize2fs /dev/disk/by-id/scsi-0DO_Volume_name

# Verify
df -h
```

**Important:** IaC tools only handle step 1. Step 2 is always manual.

### Snapshot and Restore Pattern

**Create snapshot (manual or automated):**
```bash
# Create snapshot
doctl compute volume-snapshot create \
  --volume-id vol-abc123 \
  --snapshot-name "nightly-backup-2024-01-15"

# Automated with cron
0 2 * * * doctl compute volume-snapshot create --volume-id vol-abc123 --snapshot-name "backup-$(date +\%Y-\%m-\%d)"
```

**Restore from snapshot:**
```yaml
spec:
  snapshot_id: "snapshot-123456"
  size_gib: 500  # Must be >= snapshot size
```

---

## Best Practices

### 1. Always Use Pre-Formatted Filesystems

✅ **Recommended:**
```yaml
filesystem_type: XFS  # or EXT4
```

**Benefits:**
- Simpler mount scripts
- No manual `mkfs` required
- Eliminates risk of accidentally formatting existing data

❌ **Avoid (unless necessary):**
```yaml
filesystem_type: NONE
```

### 2. Choose the Right Filesystem

**Use XFS for:**
- Databases (PostgreSQL, MySQL, MongoDB)
- Large files (videos, backups)
- High I/O workloads

**Use ext4 for:**
- General-purpose applications
- Many small files
- Legacy compatibility needs

### 3. Use Descriptive Names and Tags

✅ **Good:**
```yaml
volume_name: prod-postgres-data-2024
description: "Primary PostgreSQL data volume for production cluster"
tags: ["env:prod", "db:postgres", "cluster:primary", "backup:hourly"]
```

❌ **Avoid:**
```yaml
volume_name: vol1
```

### 4. Plan for Region Constraints

- Volume and Droplet **must be in same region**
- Snapshots **cannot be copied between regions**
- For DR, use manual file-level replication or managed services

### 5. Tag for Cost Allocation

```yaml
tags:
  - env:prod
  - project:customer-portal
  - cost-center:engineering
  - owner:platform-team
```

Enables cost tracking and organization-wide reporting.

---

## Troubleshooting

### Volume Not Visible on Droplet

**Symptom:** `df -h` doesn't show the volume

**Cause:** Volume is attached but not mounted

**Solution:**
```bash
# List block devices
lsblk

# Find your volume (usually scsi-0DO_Volume_name)
ls -l /dev/disk/by-id/

# Mount it
mount /dev/disk/by-id/scsi-0DO_Volume_name /mnt/data
```

### "Cannot Attach Volume" Error

**Cause:** Region mismatch

**Solution:** Volume and Droplet must be in the exact same region:
```bash
# Check Droplet region
doctl compute droplet get DROPLET_ID --format Region

# Check Volume region
doctl compute volume get VOLUME_ID --format Region
```

### Filesystem Resize Not Working

**Cause:** Step 2 (filesystem resize) not completed

**Solution:**
```bash
# After Terraform/Pulumi resize, always run:
# For XFS:
sudo xfs_growfs /mnt/mount-point

# For ext4:
sudo resize2fs /dev/disk/by-id/scsi-0DO_Volume_name
```

### Snapshot Restore Fails

**Cause:** Size too small or region mismatch

**Solutions:**
- Ensure `size_gib` ≥ snapshot size
- Ensure `region` matches snapshot region
- Verify snapshot ID is correct

---

## Next Steps

- Review [docs/README.md](docs/README.md) for architecture and best practices
- Check [iac/pulumi/README.md](iac/pulumi/README.md) for Pulumi deployment details
- Check [iac/tf/README.md](iac/tf/README.md) for Terraform deployment details
- See [hack/manifest.yaml](hack/manifest.yaml) for a test manifest

## Additional Resources

- [DigitalOcean Volumes Documentation](https://docs.digitalocean.com/products/volumes/)
- [Filesystem Comparison Guide](https://docs.digitalocean.com/products/volumes/how-to/format/)
- [Volume Snapshots](https://docs.digitalocean.com/products/volumes/how-to/snapshot/)

