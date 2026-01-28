# CivoVolume Examples

This document provides real-world examples of `CivoVolume` configurations for common use cases. Each example includes the YAML manifest and explains when and why to use that configuration.

## Table of Contents

1. [Minimal Development Volume](#1-minimal-development-volume)
2. [Production Database Volume](#2-production-database-volume)
3. [Multi-Environment Setup](#3-multi-environment-setup)
4. [Application Logs Storage](#4-application-logs-storage)
5. [Large Data Warehouse Volume](#5-large-data-warehouse-volume)
6. [Kubernetes StatefulSet Volume](#6-kubernetes-statefulset-volume)
7. [Shared Data Volume for Multiple Pods](#7-shared-data-volume-for-multiple-pods)

---

## 1. Minimal Development Volume

**Use Case**: Quick development/testing environment for small databases or file storage.

**Scenario**: A developer needs a persistent volume to test PostgreSQL locally on a Civo instance.

```yaml
apiVersion: civo.openmcf.org/v1
kind: CivoVolume
metadata:
  name: dev-postgres-data
  description: Development PostgreSQL data volume
spec:
  volumeName: dev-postgres-data
  region: FRA1
  sizeGib: 10
```

**Key Points**:
- Small 10 GiB size (sufficient for development)
- Frankfurt region (European developers)
- No filesystem type specified (defaults to `NONE` - format manually)
- No tags (quick provisioning)
- Minimal configuration for rapid testing

**Post-Deployment**:
```bash
# Attach volume to instance via Civo CLI
civo volume attach dev-postgres-data <instance-id>

# SSH into instance
ssh root@<instance-ip>

# Format and mount
mkfs.ext4 /dev/vdb
mkdir -p /var/lib/postgresql/data
mount /dev/vdb /var/lib/postgresql/data

# Make persistent
echo "/dev/vdb /var/lib/postgresql/data ext4 defaults,nofail 0 2" >> /etc/fstab
```

---

## 2. Production Database Volume

**Use Case**: Critical production database requiring high capacity and XFS filesystem.

**Scenario**: Production PostgreSQL database with 1 TB storage, daily backups, and high availability.

```yaml
apiVersion: civo.openmcf.org/v1
kind: CivoVolume
metadata:
  name: prod-db-data
  description: Production PostgreSQL data volume
spec:
  volumeName: prod-db-data
  region: LON1
  sizeGib: 1000
  filesystemType: XFS
  tags:
    - env:prod
    - criticality:high
    - backup:daily
    - app:postgresql
```

**Key Points**:
- 1 TB (1000 GiB) for production scale
- London region (UK data residency)
- `filesystemType: XFS` - Better for large volumes and high concurrency (informational - format manually)
- Tags for governance and compliance tracking
- Backup strategy indicated in tags

**Post-Deployment**:
```bash
# Attach volume
civo volume attach prod-db-data <db-instance-id>

# SSH and format as XFS
ssh root@<db-instance-ip>
mkfs.xfs /dev/vdb
mkdir -p /var/lib/postgresql/14/main
mount /dev/vdb /var/lib/postgresql/14/main

# Make persistent
echo "/dev/vdb /var/lib/postgresql/14/main xfs defaults,nofail 0 2" >> /etc/fstab

# Set up daily backups to object storage
cat > /root/backup.sh <<'EOF'
#!/bin/bash
DATE=$(date +%Y%m%d)
pg_dumpall | gzip > /tmp/backup-$DATE.sql.gz
aws s3 cp /tmp/backup-$DATE.sql.gz s3://prod-db-backups/ --endpoint-url https://objectstore.civo.com
rm /tmp/backup-$DATE.sql.gz
EOF

chmod +x /root/backup.sh

# Add to cron (daily at 2 AM)
echo "0 2 * * * /root/backup.sh" | crontab -
```

---

## 3. Multi-Environment Setup

**Use Case**: Separate volumes for dev, staging, and production environments.

**Scenario**: A SaaS application needs isolated storage for each environment's MySQL database.

### Development

```yaml
apiVersion: civo.openmcf.org/v1
kind: CivoVolume
metadata:
  name: myapp-dev-mysql
  description: Development MySQL data volume
spec:
  volumeName: myapp-dev-mysql
  region: NYC1
  sizeGib: 20
  filesystemType: EXT4
  tags:
    - env:dev
    - team:backend
    - app:mysql
```

### Staging

```yaml
apiVersion: civo.openmcf.org/v1
kind: CivoVolume
metadata:
  name: myapp-staging-mysql
  description: Staging MySQL data volume
spec:
  volumeName: myapp-staging-mysql
  region: NYC1
  sizeGib: 100
  filesystemType: EXT4
  tags:
    - env:staging
    - team:backend
    - app:mysql
```

### Production

```yaml
apiVersion: civo.openmcf.org/v1
kind: CivoVolume
metadata:
  name: myapp-prod-mysql
  description: Production MySQL data volume
spec:
  volumeName: myapp-prod-mysql
  region: NYC1
  sizeGib: 500
  filesystemType: XFS
  tags:
    - env:prod
    - criticality:high
    - team:backend
    - app:mysql
    - backup:hourly
```

**Key Points**:
- Same region (NYC1) for consistent latency
- Dev: 20 GiB (small dataset)
- Staging: 100 GiB (production-like data)
- Prod: 500 GiB (full production scale) with XFS for performance
- Tags differentiate environments for cost tracking and access policies

---

## 4. Application Logs Storage

**Use Case**: Centralized log storage volume for application servers.

**Scenario**: Web application logs aggregated to a shared volume, retained for 30 days before rotation.

```yaml
apiVersion: civo.openmcf.org/v1
kind: CivoVolume
metadata:
  name: app-logs-volume
  description: Application logs storage with 30-day retention
spec:
  volumeName: app-logs-volume
  region: FRA1
  sizeGib: 200
  filesystemType: EXT4
  tags:
    - env:prod
    - data-type:logs
    - retention:30-days
    - team:devops
```

**Post-Deployment**:
```bash
# Attach and mount
civo volume attach app-logs-volume <instance-id>
ssh root@<instance-ip>
mkfs.ext4 /dev/vdb
mkdir -p /var/log/app
mount /dev/vdb /var/log/app
echo "/dev/vdb /var/log/app ext4 defaults,nofail 0 2" >> /etc/fstab

# Configure log rotation (delete logs older than 30 days)
cat > /etc/logrotate.d/app-logs <<'EOF'
/var/log/app/*.log {
    daily
    rotate 30
    compress
    delaycompress
    missingok
    notifempty
    create 0644 root root
}
EOF
```

**Key Points**:
- 200 GiB for log accumulation
- EXT4 for stable log writing
- Frankfurt region (GDPR compliance)
- Tag-based retention policy documentation
- Logrotate for automatic cleanup

---

## 5. Large Data Warehouse Volume

**Use Case**: Multi-terabyte volume for analytics and data warehousing.

**Scenario**: ClickHouse data warehouse requiring 8 TB of storage for historical analytics.

```yaml
apiVersion: civo.openmcf.org/v1
kind: CivoVolume
metadata:
  name: clickhouse-data-warehouse
  description: ClickHouse data warehouse volume (8 TB)
spec:
  volumeName: clickhouse-data-warehouse
  region: LON1
  sizeGib: 8000
  filesystemType: XFS
  tags:
    - env:prod
    - team:data-analytics
    - app:clickhouse
    - criticality:high
```

**Key Points**:
- 8 TB (8000 GiB) for large-scale analytics
- XFS for better performance with large files and high concurrency
- London region for European data residency
- High criticality due to business intelligence dependencies

**Post-Deployment**:
```bash
# Attach and format as XFS
civo volume attach clickhouse-data-warehouse <clickhouse-instance-id>
ssh root@<clickhouse-instance-ip>
mkfs.xfs /dev/vdb
mkdir -p /var/lib/clickhouse
mount /dev/vdb /var/lib/clickhouse
echo "/dev/vdb /var/lib/clickhouse xfs defaults,nofail 0 2" >> /etc/fstab

# Optimize XFS mount options for ClickHouse
mount -o remount,noatime,nodiratime /var/lib/clickhouse
sed -i 's|/dev/vdb /var/lib/clickhouse xfs defaults,nofail 0 2|/dev/vdb /var/lib/clickhouse xfs defaults,noatime,nodiratime,nofail 0 2|' /etc/fstab
```

---

## 6. Kubernetes StatefulSet Volume

**Use Case**: Persistent storage for StatefulSet pods (e.g., MongoDB replica set).

**Scenario**: MongoDB replica set with 3 pods, each requiring a dedicated 100 GiB volume.

**Note**: For Kubernetes, you typically use `PersistentVolumeClaim` with the `civo-volume` StorageClass for dynamic provisioning. However, this example shows pre-provisioning volumes with OpenMCF for more control.

### Pre-Provision Volumes

```yaml
# Volume for mongodb-0
apiVersion: civo.openmcf.org/v1
kind: CivoVolume
metadata:
  name: mongodb-data-0
  description: MongoDB replica set volume for pod 0
spec:
  volumeName: mongodb-data-0
  region: NYC1
  sizeGib: 100
  filesystemType: XFS
  tags:
    - env:prod
    - app:mongodb
    - pod:mongodb-0
---
# Volume for mongodb-1
apiVersion: civo.openmcf.org/v1
kind: CivoVolume
metadata:
  name: mongodb-data-1
  description: MongoDB replica set volume for pod 1
spec:
  volumeName: mongodb-data-1
  region: NYC1
  sizeGib: 100
  filesystemType: XFS
  tags:
    - env:prod
    - app:mongodb
    - pod:mongodb-1
---
# Volume for mongodb-2
apiVersion: civo.openmcf.org/v1
kind: CivoVolume
metadata:
  name: mongodb-data-2
  description: MongoDB replica set volume for pod 2
spec:
  volumeName: mongodb-data-2
  region: NYC1
  sizeGib: 100
  filesystemType: XFS
  tags:
    - env:prod
    - app:mongodb
    - pod:mongodb-2
```

### Reference in Kubernetes PersistentVolume

After provisioning, reference the volume IDs in PersistentVolume manifests:

```yaml
apiVersion: v1
kind: PersistentVolume
metadata:
  name: mongodb-data-0
spec:
  capacity:
    storage: 100Gi
  accessModes:
    - ReadWriteOnce
  csi:
    driver: civo.csi.k8s.io
    volumeHandle: <volume-id-from-stack-outputs>
```

**Key Points**:
- Three separate volumes for StatefulSet pods
- 100 GiB each for MongoDB data
- XFS for database performance
- Tags track pod affinity
- CSI driver handles attachment and formatting automatically

---

## 7. Shared Data Volume for Multiple Pods

**Use Case**: Read-only shared data volume for multiple application pods.

**Scenario**: ML model files stored on a volume, mounted read-only by multiple inference pods.

**Important**: Civo Volumes only support single-attach (`ReadWriteOnce`). For true shared storage (`ReadWriteMany`), use Civo Object Storage or NFS.

**Workaround**: Attach the volume to a dedicated "data server" pod and use NFS to share it:

### Step 1: Create the Data Volume

```yaml
apiVersion: civo.openmcf.org/v1
kind: CivoVolume
metadata:
  name: ml-models-volume
  description: Shared ML model storage
spec:
  volumeName: ml-models-volume
  region: FRA1
  sizeGib: 50
  filesystemType: EXT4
  tags:
    - env:prod
    - team:ml
    - data-type:models
```

### Step 2: Deploy NFS Server Pod (Kubernetes)

```yaml
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: ml-models-pvc
spec:
  volumeName: ml-models-volume
  accessModes:
    - ReadWriteOnce
  resources:
    requests:
      storage: 50Gi
---
apiVersion: v1
kind: Pod
metadata:
  name: nfs-server
  labels:
    app: nfs-server
spec:
  containers:
  - name: nfs-server
    image: itsthenetwork/nfs-server-alpine:latest
    ports:
    - name: nfs
      containerPort: 2049
    securityContext:
      privileged: true
    volumeMounts:
    - name: data
      mountPath: /nfsshare
  volumes:
  - name: data
    persistentVolumeClaim:
      claimName: ml-models-pvc
---
apiVersion: v1
kind: Service
metadata:
  name: nfs-server
spec:
  ports:
  - port: 2049
    name: nfs
  selector:
    app: nfs-server
```

### Step 3: Mount NFS in Application Pods

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: ml-inference-pod-1
spec:
  containers:
  - name: inference
    image: ml-inference:latest
    volumeMounts:
    - name: models
      mountPath: /models
      readOnly: true
  volumes:
  - name: models
    nfs:
      server: nfs-server
      path: /
      readOnly: true
```

**Key Points**:
- Civo Volumes are single-attach only
- NFS server pod bridges the gap for shared read-only access
- ML inference pods mount NFS read-only
- 50 GiB sufficient for model files
- Production should use managed NFS or object storage for true shared storage

---

## Best Practices Summary

1. **Naming Convention**: Use consistent patterns (`<app>-<env>-<purpose>`)
2. **Region Selection**: Collocate with compute resources for performance
3. **Filesystem Choice**: 
   - EXT4 for general use (stable, well-tested)
   - XFS for large volumes (>1 TB) or high concurrency
4. **Size Planning**: Start small, expand as needed (you can't shrink)
5. **Tags**: Use for environment tracking, cost allocation, and compliance
6. **Backups**: Application-level backups to object storage (no native snapshots)
7. **Kubernetes**: Use `civo-volume` StorageClass for dynamic provisioning

---

## Testing Your Volume

After provisioning any volume, verify it works:

```bash
# 1. Attach volume via Civo CLI
civo volume attach <volume-name> <instance-id>

# 2. SSH into instance
ssh root@<instance-ip>

# 3. Identify the device
lsblk
# Look for the new device (e.g., vdb)

# 4. Format
mkfs.ext4 /dev/vdb

# 5. Mount
mkdir -p /data
mount /dev/vdb /data

# 6. Test write
echo "Test $(date)" > /data/test.txt

# 7. Test read
cat /data/test.txt

# 8. Test persistence
df -h /data

# 9. Make mount persistent
echo "/dev/vdb /data ext4 defaults,nofail 0 2" >> /etc/fstab

# 10. Verify /etc/fstab
cat /etc/fstab

# 11. Test reboot persistence
reboot
# After reboot, SSH back in and verify:
df -h /data
cat /data/test.txt
```

---

## Advanced Use Cases

### RAID0 for Extreme IOPS

If you need higher IOPS than a single volume provides, use LVM to stripe multiple volumes:

```bash
# Create 4 volumes (e.g., 250 GiB each)
# Attach all 4 to the same instance

# Create RAID0 stripe
apt-get install lvm2
pvcreate /dev/vdb /dev/vdc /dev/vdd /dev/vde
vgcreate data_vg /dev/vdb /dev/vdc /dev/vdd /dev/vde
lvcreate -l 100%FREE -i 4 -I 64 -n data_lv data_vg

# Format and mount
mkfs.xfs /dev/data_vg/data_lv
mkdir -p /data
mount /dev/data_vg/data_lv /data

# Make persistent
echo "/dev/data_vg/data_lv /data xfs defaults,nofail 0 2" >> /etc/fstab
```

**Result**: 1 TB volume (4 x 250 GiB) with 4x the IOPS of a single volume.

**Warning**: RAID0 has no redundancy. If one volume fails, all data is lost. Use only for ephemeral high-performance workloads.

---

## Related Documentation

- **API Reference**: [README.md](./README.md)
- **Research**: [docs/README.md](./docs/README.md)
- **Pulumi Module**: [iac/pulumi/](./iac/pulumi/)
- **Terraform Module**: [iac/tf/](./iac/tf/)
- **Civo Documentation**: [Civo Volumes](https://www.civo.com/docs/compute/instance-volumes)

