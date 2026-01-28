# KubernetesJob Examples

This document provides practical examples of KubernetesJob configurations for common use cases.

## Basic Job

A minimal job that runs a single command to completion.

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: KubernetesJob
metadata:
  name: hello-world
spec:
  namespace:
    value: default
  image:
    repo: busybox
    tag: latest
  command:
    - echo
    - "Hello, World!"
```

## Database Migration Job

A job that runs database migrations with environment variables and secrets.

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: KubernetesJob
metadata:
  name: db-migration
spec:
  namespace:
    value: production
  image:
    repo: myregistry/migration-runner
    tag: v2.0.0
  resources:
    limits:
      cpu: 1000m
      memory: 2Gi
    requests:
      cpu: 250m
      memory: 512Mi
  env:
    variables:
      DATABASE_HOST:
        value: postgres.production.svc.cluster.local
      DATABASE_NAME:
        value: myapp
    secrets:
      DATABASE_USER:
        secretRef:
          name: db-credentials
          key: username
      DATABASE_PASSWORD:
        secretRef:
          name: db-credentials
          key: password
  backoffLimit: 3
  activeDeadlineSeconds: 1800
  ttlSecondsAfterFinished: 3600
  restartPolicy: Never
  command:
    - python
    - manage.py
    - migrate
```

## Parallel Data Processing Job

A job that processes data in parallel using multiple pods.

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: KubernetesJob
metadata:
  name: parallel-processor
spec:
  namespace:
    value: data-processing
  createNamespace: true
  image:
    repo: myregistry/data-processor
    tag: v1.5.0
  resources:
    limits:
      cpu: 2000m
      memory: 4Gi
    requests:
      cpu: 500m
      memory: 1Gi
  parallelism: 5
  completions: 20
  backoffLimit: 6
  env:
    variables:
      INPUT_BUCKET:
        value: s3://my-bucket/input
      OUTPUT_BUCKET:
        value: s3://my-bucket/output
    secrets:
      AWS_ACCESS_KEY_ID:
        secretRef:
          name: aws-credentials
          key: access-key-id
      AWS_SECRET_ACCESS_KEY:
        secretRef:
          name: aws-credentials
          key: secret-access-key
  command:
    - python
    - /app/process.py
```

## Indexed Job for Partitioned Processing

A job that uses indexed mode to process partitioned data where each pod handles a specific partition.

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: KubernetesJob
metadata:
  name: indexed-processor
spec:
  namespace:
    value: data-processing
  image:
    repo: myregistry/partition-processor
    tag: v1.0.0
  resources:
    limits:
      cpu: 1000m
      memory: 2Gi
    requests:
      cpu: 250m
      memory: 512Mi
  completionMode: Indexed
  completions: 10
  parallelism: 5
  backoffLimit: 3
  env:
    variables:
      TOTAL_PARTITIONS:
        value: "10"
  command:
    - /bin/sh
    - -c
    - "echo Processing partition $JOB_COMPLETION_INDEX of $TOTAL_PARTITIONS && python /app/process_partition.py --partition=$JOB_COMPLETION_INDEX"
```

## Backup Job with Volume Mounts

A job that creates a backup using mounted volumes.

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: KubernetesJob
metadata:
  name: database-backup
spec:
  namespace:
    value: backups
  createNamespace: true
  image:
    repo: postgres
    tag: "15"
  resources:
    limits:
      cpu: 500m
      memory: 1Gi
    requests:
      cpu: 100m
      memory: 256Mi
  env:
    variables:
      PGHOST:
        value: postgres.production.svc.cluster.local
      PGDATABASE:
        value: myapp
    secrets:
      PGUSER:
        secretRef:
          name: db-credentials
          key: username
      PGPASSWORD:
        secretRef:
          name: db-credentials
          key: password
  backoffLimit: 2
  activeDeadlineSeconds: 7200
  ttlSecondsAfterFinished: 604800
  configMaps:
    backup-script: |
      #!/bin/bash
      set -e
      BACKUP_FILE="/backup/backup-$(date +%Y%m%d-%H%M%S).sql.gz"
      echo "Creating backup: $BACKUP_FILE"
      pg_dump | gzip > "$BACKUP_FILE"
      echo "Backup complete: $BACKUP_FILE"
  volumeMounts:
    - name: backup-script
      mountPath: /scripts/backup.sh
      configMap:
        name: backup-script
        key: backup-script
        path: backup.sh
        defaultMode: 493
    - name: backup-storage
      mountPath: /backup
      pvc:
        claimName: backup-pvc
  command:
    - /bin/bash
    - /scripts/backup.sh
```

## ETL Job with Multiple Containers

A job for ETL processing with environment-specific configuration.

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: KubernetesJob
metadata:
  name: etl-daily
spec:
  namespace:
    value: etl
  image:
    repo: myregistry/etl-runner
    tag: v3.2.1
  resources:
    limits:
      cpu: 4000m
      memory: 8Gi
    requests:
      cpu: 1000m
      memory: 2Gi
  env:
    variables:
      ETL_MODE:
        value: daily
      SOURCE_DATABASE:
        value: source-db.internal:5432/analytics
      TARGET_DATABASE:
        value: target-db.internal:5432/warehouse
      LOG_LEVEL:
        value: INFO
    secrets:
      SOURCE_DB_PASSWORD:
        secretRef:
          name: etl-secrets
          key: source-password
      TARGET_DB_PASSWORD:
        secretRef:
          name: etl-secrets
          key: target-password
  backoffLimit: 1
  activeDeadlineSeconds: 21600
  ttlSecondsAfterFinished: 86400
  restartPolicy: Never
  command:
    - python
    - -m
    - etl.main
  args:
    - --date
    - "$(date -d 'yesterday' +%Y-%m-%d)"
    - --full-refresh=false
```

## Cleanup Job

A job that performs cleanup operations.

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: KubernetesJob
metadata:
  name: cleanup-old-data
spec:
  namespace:
    value: maintenance
  image:
    repo: myregistry/cleanup-tool
    tag: v1.0.0
  resources:
    limits:
      cpu: 500m
      memory: 512Mi
    requests:
      cpu: 100m
      memory: 128Mi
  env:
    variables:
      RETENTION_DAYS:
        value: "30"
      DRY_RUN:
        value: "false"
  backoffLimit: 2
  activeDeadlineSeconds: 3600
  ttlSecondsAfterFinished: 600
  command:
    - /app/cleanup
    - --older-than
    - 30d
    - --confirm
```

## Deployment Instructions

Deploy any of these examples using the OpenMCF CLI:

```bash
# Preview changes
openmcf pulumi preview --manifest job.yaml

# Deploy
openmcf pulumi up --manifest job.yaml

# Check status
kubectl get jobs -n <namespace>

# View logs
kubectl logs job/<job-name> -n <namespace>

# Delete when done (if TTL not set)
openmcf pulumi destroy --manifest job.yaml
```
