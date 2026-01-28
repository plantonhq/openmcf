# DigitalOcean App Platform Service Examples

This document provides comprehensive examples for deploying applications on DigitalOcean App Platform using OpenMCF's declarative API.

## Table of Contents

- [Basic Web Service (Git Source)](#basic-web-service-git-source)
- [Web Service with Autoscaling](#web-service-with-autoscaling)
- [Container Image Deployment](#container-image-deployment-from-docr)
- [Background Worker Service](#background-worker-service)
- [Pre-Deployment Job](#pre-deployment-job)
- [Web Service with Custom Domain](#web-service-with-custom-domain)
- [Multi-Environment Configuration](#multi-environment-configuration)
- [Service with Database Integration](#service-with-database-integration)

---

## Basic Web Service (Git Source)

Deploy a simple Node.js web application from a GitHub repository.

```yaml
apiVersion: digital-ocean.openmcf.org/v1
kind: DigitalOceanAppPlatformService
metadata:
  name: simple-api
  labels:
    environment: development
    team: backend
spec:
  serviceName: simple-api
  region: nyc1
  serviceType: web_service
  
  gitSource:
    repoUrl: https://github.com/myorg/simple-api.git
    branch: main
  
  instanceSizeSlug: basic_xxs
  instanceCount: 1
  
  env:
    NODE_ENV: production
    PORT: "8080"
```

**Key characteristics:**
- **Cost**: $5/month (Basic XXS instance)
- **Use case**: Development, prototyping, low-traffic services
- **Deployment**: Auto-deploys on git push to `main` branch

---

## Web Service with Autoscaling

Production-grade service with CPU-based autoscaling.

```yaml
apiVersion: digital-ocean.openmcf.org/v1
kind: DigitalOceanAppPlatformService
metadata:
  name: production-api
  labels:
    environment: production
    team: backend
    criticality: high
spec:
  serviceName: production-api
  region: sfo3
  serviceType: web_service
  
  gitSource:
    repoUrl: https://github.com/myorg/production-api.git
    branch: production
    buildCommand: "npm run build"
    runCommand: "npm start"
  
  instanceSizeSlug: professional_s
  enableAutoscale: true
  minInstanceCount: 2
  maxInstanceCount: 10
  
  env:
    NODE_ENV: production
    LOG_LEVEL: info
    MAX_CONNECTIONS: "100"
```

**Key characteristics:**
- **Cost**: $48-$240/month (2-10 Professional S instances)
- **Use case**: Production APIs with variable traffic
- **HA**: Minimum 2 instances for redundancy
- **Scaling**: Auto-scales at 80% CPU utilization

---

## Container Image Deployment (from DOCR)

Deploy a pre-built Docker image from DigitalOcean Container Registry.

```yaml
apiVersion: digital-ocean.openmcf.org/v1
kind: DigitalOceanAppPlatformService
metadata:
  name: backend-service
spec:
  serviceName: backend-service
  region: nyc3
  serviceType: web_service
  
  imageSource:
    registry: registry.digitalocean.com/myorg
    repository: backend-service
    tag: v1.2.3
  
  instanceSizeSlug: professional_xs
  instanceCount: 3
  
  env:
    ENVIRONMENT: production
    DATABASE_POOL_SIZE: "20"
```

**Key characteristics:**
- **Cost**: $36/month (3 Professional XS instances)
- **Use case**: Immutable production deployments
- **Deployment**: Manual (redeploy by changing `tag`)
- **CI/CD**: Build and push images in CI, then update tag

**Prerequisites:**
- DigitalOcean Container Registry created
- Docker images built and pushed to DOCR
- Image tags follow semantic versioning

---

## Background Worker Service

Deploy a Redis queue worker for async task processing.

```yaml
apiVersion: digital-ocean.openmcf.org/v1
kind: DigitalOceanAppPlatformService
metadata:
  name: email-worker
  labels:
    environment: production
    service-type: worker
spec:
  serviceName: email-worker
  region: tor1
  serviceType: worker
  
  gitSource:
    repoUrl: https://github.com/myorg/workers.git
    branch: main
    runCommand: "python worker.py"
  
  instanceSizeSlug: basic_s
  instanceCount: 2
  
  env:
    WORKER_TYPE: email
    REDIS_URL: ${digitalocean-redis-prod.connection-uri}
    QUEUE_NAME: emails
    WORKER_CONCURRENCY: "10"
```

**Key characteristics:**
- **Cost**: $48/month (2 Basic S instances)
- **Use case**: Background job processing (email sending, image processing, data pipelines)
- **No HTTP traffic**: Workers don't receive external requests
- **Scaling**: Manual scaling (workers don't support autoscaling)

**Common worker patterns:**
- Sidekiq (Ruby)
- Celery (Python)
- Bull/BullMQ (Node.js)
- RQ (Python)

---

## Pre-Deployment Job

Run database migrations before deploying the main application.

```yaml
apiVersion: digital-ocean.openmcf.org/v1
kind: DigitalOceanAppPlatformService
metadata:
  name: db-migration
spec:
  serviceName: db-migration
  region: lon1
  serviceType: job
  
  gitSource:
    repoUrl: https://github.com/myorg/app.git
    branch: main
    runCommand: "npm run migrate"
  
  instanceSizeSlug: basic_xxs
  
  env:
    DATABASE_URL: ${digitalocean-postgres-prod.connection-uri}
    MIGRATION_TIMEOUT: "300"
```

**Key characteristics:**
- **Cost**: Minimal (job runs only during deployments)
- **Use case**: Database migrations, schema updates, data seeding
- **Execution**: Runs before each deployment
- **Failure handling**: Deployment fails if job exits with non-zero code

**Best practices:**
- Keep migrations idempotent
- Set appropriate timeout values
- Test migrations in staging first

---

## Web Service with Custom Domain

Deploy a production service with a custom domain and SSL.

```yaml
apiVersion: digital-ocean.openmcf.org/v1
kind: DigitalOceanAppPlatformService
metadata:
  name: api-myapp
spec:
  serviceName: api-myapp
  region: fra1
  serviceType: web_service
  
  gitSource:
    repoUrl: https://github.com/myorg/myapp-api.git
    branch: main
  
  instanceSizeSlug: professional_m
  instanceCount: 3
  
  env:
    NODE_ENV: production
    ALLOWED_ORIGINS: "https://myapp.com,https://www.myapp.com"
  
  customDomain: api.myapp.com
```

**Key characteristics:**
- **Domain**: `api.myapp.com` (instead of `*.ondigitalocean.app`)
- **SSL**: Automatically provisioned and renewed
- **DNS**: Requires CNAME record pointing to App Platform

**DNS Setup (if using external DNS):**
1. Add CNAME record: `api.myapp.com` → `<app-name>.ondigitalocean.app`
2. Wait for DNS propagation (up to 48 hours)
3. App Platform automatically provisions SSL certificate

**If using DigitalOceanDnsZone resource:**
```yaml
custom_domain:
  kind: DigitalOceanDnsZone
  field_path: spec.domain_name
  value: api.myapp.com
```

---

## Multi-Environment Configuration

Deploy the same application across dev, staging, and production.

### Development Environment

```yaml
apiVersion: digital-ocean.openmcf.org/v1
kind: DigitalOceanAppPlatformService
metadata:
  name: myapp-dev
  labels:
    environment: development
spec:
  serviceName: myapp-dev
  region: nyc1
  serviceType: web_service
  
  gitSource:
    repoUrl: https://github.com/myorg/myapp.git
    branch: develop
  
  instanceSizeSlug: basic_xxs
  instanceCount: 1
  
  env:
    NODE_ENV: development
    DATABASE_URL: ${digitalocean-postgres-dev.connection-uri}
    LOG_LEVEL: debug
```

### Staging Environment

```yaml
apiVersion: digital-ocean.openmcf.org/v1
kind: DigitalOceanAppPlatformService
metadata:
  name: myapp-staging
  labels:
    environment: staging
spec:
  serviceName: myapp-staging
  region: nyc1
  serviceType: web_service
  
  gitSource:
    repoUrl: https://github.com/myorg/myapp.git
    branch: staging
  
  instanceSizeSlug: basic_s
  instanceCount: 2
  
  env:
    NODE_ENV: production
    DATABASE_URL: ${digitalocean-postgres-staging.connection-uri}
    LOG_LEVEL: info
```

### Production Environment

```yaml
apiVersion: digital-ocean.openmcf.org/v1
kind: DigitalOceanAppPlatformService
metadata:
  name: myapp-prod
  labels:
    environment: production
    criticality: high
spec:
  serviceName: myapp-prod
  region: nyc3
  serviceType: web_service
  
  imageSource:
    registry: registry.digitalocean.com/myorg
    repository: myapp
    tag: v1.0.0
  
  instanceSizeSlug: professional_m
  enableAutoscale: true
  minInstanceCount: 3
  maxInstanceCount: 15
  
  env:
    NODE_ENV: production
    DATABASE_URL: ${digitalocean-postgres-prod.connection-uri}
    LOG_LEVEL: warn
  
  customDomain: api.myapp.com
```

**Environment progression strategy:**
1. **Dev**: Git-based auto-deploy from `develop` branch, single instance
2. **Staging**: Git-based from `staging` branch, production-like configuration
3. **Production**: Container image deployment, immutable tags, autoscaling

---

## Service with Database Integration

Deploy a web service connected to a DigitalOcean Managed PostgreSQL database.

```yaml
apiVersion: digital-ocean.openmcf.org/v1
kind: DigitalOceanAppPlatformService
metadata:
  name: orders-api
spec:
  serviceName: orders-api
  region: sgp1
  serviceType: web_service
  
  gitSource:
    repoUrl: https://github.com/myorg/orders-api.git
    branch: main
  
  instanceSizeSlug: professional_xs
  enableAutoscale: true
  minInstanceCount: 2
  maxInstanceCount: 6
  
  env:
    NODE_ENV: production
    DATABASE_URL: ${digitalocean-postgres-orders.connection-uri}
    DATABASE_POOL_MIN: "5"
    DATABASE_POOL_MAX: "20"
    REDIS_URL: ${digitalocean-redis-cache.connection-uri}
    CACHE_TTL: "3600"
```

**Integration with other resources:**
- `${digitalocean-postgres-orders.connection-uri}`: References DigitalOceanPostgres resource
- `${digitalocean-redis-cache.connection-uri}`: References DigitalOceanRedis resource

**Prerequisites:**
- DigitalOcean Managed Database (PostgreSQL) created
- DigitalOcean Managed Redis created
- Databases configured in same VPC as App Platform (for private networking)

---

## Best Practices Summary

### Source Configuration
- **Development**: Use `git_source` with `develop` branch for auto-deploy
- **Staging**: Use `git_source` with `staging` branch
- **Production**: Use `image_source` with semantic version tags (`v1.0.0`)

### Scaling
- **Web services**: Enable autoscaling for variable traffic
- **Workers**: Use fixed `instance_count` based on queue depth
- **Jobs**: Single instance (no scaling needed)

### Instance Sizing
- **Dev/Test**: `basic_xxs` or `basic_xs`
- **Production Low-Traffic**: `basic_s` or `professional_xs`
- **Production High-Traffic**: `professional_s` or higher

### High Availability
- **Always use min_instance_count ≥ 2** for production web services
- **Deploy across multiple regions** for geo-redundancy
- **Use health checks** (automatic in App Platform)

### Security
- **Use environment variables** for secrets (encrypted by DigitalOcean)
- **Rotate credentials regularly** via env var updates
- **Enable custom domains** with automatic SSL
- **Use VPC networking** for database connections

---

## Next Steps

- Review [README.md](./README.md) for field reference and detailed documentation
- Check [iac/pulumi/README.md](./iac/pulumi/README.md) for Pulumi deployment
- Check [iac/tf/README.md](./iac/tf/README.md) for Terraform deployment
- Explore [DigitalOcean App Platform docs](https://docs.digitalocean.com/products/app-platform/)

