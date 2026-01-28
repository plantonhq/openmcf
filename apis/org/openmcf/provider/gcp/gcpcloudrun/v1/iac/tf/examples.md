# GCP Cloud Run Terraform Examples

This document provides comprehensive examples of using the GCP Cloud Run Terraform module for various deployment scenarios.

## Table of Contents

1. [Minimal Configuration](#minimal-configuration)
2. [Standard Configuration with Environment Variables](#standard-configuration-with-environment-variables)
3. [Custom DNS Domain Mapping](#custom-dns-domain-mapping)
4. [Private Service with VPC Access](#private-service-with-vpc-access)
5. [High-Traffic Production Service](#high-traffic-production-service)
6. [Multi-Region Deployment](#multi-region-deployment)

---

## Minimal Configuration

The simplest Cloud Run deployment with required fields only.

```hcl
module "minimal_cloudrun" {
  source = "./iac/tf"

  metadata = {
    name = "hello-world"
    env  = "dev"
  }

  spec = {
    # Project ID using StringValueOrRef format
    project_id = {
      value = "my-gcp-project-123"
    }
    region = "us-central1"

    container = {
      image = {
        repo = "gcr.io/cloudrun/hello"
        tag  = "latest"
      }
      cpu    = 1
      memory = 512
      replicas = {
        min = 0
        max = 10
      }
    }
  }
}

output "service_url" {
  value = module.minimal_cloudrun.url
}
```

**Use case:** Quick prototype or demo service with scale-to-zero capability.

---

## Standard Configuration with Environment Variables

A production-ready service with environment variables and secrets.

```hcl
module "api_service" {
  source = "./iac/tf"

  metadata = {
    name = "api-service"
    id   = "api-svc-001"
    org  = "acme-corp"
    env  = "production"
    labels = {
      team    = "backend"
      service = "api"
    }
  }

  spec = {
    # Project ID using StringValueOrRef format
    project_id = {
      value = "acme-prod-123"
    }
    region          = "us-west1"
    service_name    = "acme-api"
    service_account = "api-sa@acme-prod-123.iam.gserviceaccount.com"

    container = {
      image = {
        repo = "us-docker.pkg.dev/acme-prod-123/api/service"
        tag  = "v1.2.3"
      }

      # Environment variables and secrets
      env = {
        variables = {
          ENV           = "production"
          PORT          = "8080"
          LOG_LEVEL     = "info"
          FEATURE_FLAGS = "new_ui,analytics"
        }
        secrets = {
          API_KEY          = "projects/123456/secrets/api-key:latest"
          DATABASE_URL     = "projects/123456/secrets/db-url:1"
          OAUTH_SECRET     = "projects/123456/secrets/oauth:latest"
        }
      }

      port   = 8080
      cpu    = 2
      memory = 1024
      replicas = {
        min = 2  # Always warm instances
        max = 100
      }
    }

    max_concurrency       = 100
    timeout_seconds       = 300
    allow_unauthenticated = true
  }
}

output "api_url" {
  value = module.api_service.url
}

output "service_name" {
  value = module.api_service.service_name
}

output "latest_revision" {
  value = module.api_service.revision
}
```

**Use case:** Production API service with secrets management, custom service account, and predictable performance (min 2 instances).

---

## Production Service with Deletion Protection

Critical production service with deletion protection enabled to prevent accidental deletion.

```hcl
module "critical_service" {
  source = "./iac/tf"

  metadata = {
    name = "payment-processor"
    id   = "pay-prod-001"
    org  = "finance-corp"
    env  = "production"
    labels = {
      tier        = "critical"
      compliance  = "pci-dss"
      cost_center = "payments"
    }
  }

  spec = {
    # Project ID using StringValueOrRef format
    project_id = {
      value = "finance-prod-123"
    }
    region          = "us-central1"
    service_account = "payment-sa@finance-prod-123.iam.gserviceaccount.com"

    container = {
      image = {
        repo = "us-docker.pkg.dev/finance-prod-123/services/payment"
        tag  = "v2.5.1"
      }

      env = {
        variables = {
          ENV              = "production"
          TRANSACTION_MODE = "live"
        }
        secrets = {
          STRIPE_KEY       = "projects/finance-prod-123/secrets/stripe-key:latest"
          DATABASE_URL     = "projects/finance-prod-123/secrets/db-url:latest"
          ENCRYPTION_KEY   = "projects/finance-prod-123/secrets/enc-key:latest"
        }
      }

      port   = 8080
      cpu    = 4
      memory = 4096
      replicas = {
        min = 5  # Always warm for critical service
        max = 100
      }
    }

    max_concurrency       = 80
    timeout_seconds       = 300
    allow_unauthenticated = false

    # Enable deletion protection for critical production service
    # Must be set to false before the service can be deleted
    delete_protection = true
  }
}

output "service_url" {
  value = module.critical_service.url
}

output "service_name" {
  value = module.critical_service.service_name
}
```

**Use case:** Critical production service (e.g., payment processing) where accidental deletion would cause significant business impact. Deletion protection ensures the service cannot be deleted until explicitly disabled.

**Note:** To delete a service with deletion protection enabled, you must first set `delete_protection = false` and apply the change, then destroy the resource.

---

## Custom DNS Domain Mapping

Deploy a service with custom domain mapping and SSL certificate.

```hcl
module "custom_domain_service" {
  source = "./iac/tf"

  metadata = {
    name = "public-api"
    env  = "production"
  }

  spec = {
    # Project ID using StringValueOrRef format
    project_id = {
      value = "my-project-123"
    }
    region = "us-east1"

    container = {
      image = {
        repo = "us-docker.pkg.dev/my-project-123/apps/api"
        tag  = "v2.0.0"
      }
      cpu    = 2
      memory = 2048
      replicas = {
        min = 1
        max = 50
      }
    }

    # Custom DNS configuration
    dns = {
      enabled      = true
      hostnames    = ["api.example.com", "api-v2.example.com"]
      managed_zone = "example-com-zone"
    }

    max_concurrency       = 80
    allow_unauthenticated = true
  }
}

# Note: Ensure DNS managed zone exists and is properly configured
# The module will create the TXT record for domain verification

output "custom_domain_url" {
  value = "https://api.example.com"
}

output "cloudrun_url" {
  value = module.custom_domain_service.url
}
```

**Use case:** Public-facing service with branded custom domain and automatic SSL certificate provisioning.

---

## Private Service with VPC Access

Deploy a private service accessible only within a VPC network.

```hcl
module "internal_service" {
  source = "./iac/tf"

  metadata = {
    name = "internal-processor"
    env  = "production"
  }

  spec = {
    # Project ID using StringValueOrRef format
    project_id = {
      value = "my-project-123"
    }
    region = "us-central1"

    container = {
      image = {
        repo = "us-docker.pkg.dev/my-project-123/internal/processor"
        tag  = "v1.5.0"
      }

      env = {
        variables = {
          DATABASE_HOST = "10.0.0.5"
          REDIS_HOST    = "10.0.0.10"
        }
      }

      cpu    = 4
      memory = 4096
      replicas = {
        min = 3
        max = 20
      }
    }

    # VPC access configuration using StringValueOrRef format
    vpc_access = {
      network = {
        value = "projects/my-project-123/global/networks/private-vpc"
      }
      subnet = {
        value = "projects/my-project-123/regions/us-central1/subnetworks/cloudrun-subnet"
      }
      egress = "PRIVATE_RANGES_ONLY"
    }

    # Internal traffic only
    ingress               = "INGRESS_TRAFFIC_INTERNAL_ONLY"
    allow_unauthenticated = false

    max_concurrency = 100
    timeout_seconds = 600
  }
}

# IAM binding for authorized services
resource "google_cloud_run_v2_service_iam_member" "authorized_service" {
  project  = module.internal_service.service_name
  location = "us-central1"
  name     = module.internal_service.service_name

  role   = "roles/run.invoker"
  member = "serviceAccount:frontend-sa@my-project-123.iam.gserviceaccount.com"
}

output "internal_url" {
  value = module.internal_service.url
}
```

**Use case:** Backend service that needs to access private databases or internal APIs within a VPC.

---

## High-Traffic Production Service

Optimized configuration for high-throughput production workloads.

```hcl
module "high_traffic_service" {
  source = "./iac/tf"

  metadata = {
    name = "video-transcoder"
    id   = "vt-prod-001"
    org  = "media-corp"
    env  = "production"
    labels = {
      tier        = "premium"
      cost_center = "media"
    }
  }

  spec = {
    # Project ID using StringValueOrRef format
    project_id = {
      value = "media-prod-123"
    }
    region          = "us-west2"
    service_account = "transcoder-sa@media-prod-123.iam.gserviceaccount.com"

    container = {
      image = {
        repo = "us-docker.pkg.dev/media-prod-123/services/transcoder"
        tag  = "v3.1.2"
      }

      env = {
        variables = {
          WORKER_THREADS    = "8"
          MAX_FILE_SIZE_MB = "500"
          OUTPUT_FORMAT    = "mp4,webm"
        }
        secrets = {
          STORAGE_KEY      = "projects/media-prod-123/secrets/gcs-key:latest"
          LICENSE_KEY      = "projects/media-prod-123/secrets/encoder-license:latest"
        }
      }

      port = 8080
      cpu  = 4  # Maximum CPU allocation
      memory = 8192  # 8GB memory
      replicas = {
        min = 10  # Always warm for instant response
        max = 200 # Scale to handle traffic spikes
      }
    }

    max_concurrency = 50  # Lower to handle compute-intensive tasks
    timeout_seconds = 3600  # 1 hour for long transcoding jobs

    execution_environment = "EXECUTION_ENVIRONMENT_GEN2"
    allow_unauthenticated = false  # Require authentication
    ingress               = "INGRESS_TRAFFIC_ALL"
  }
}

# Monitoring alert for high instance usage
output "service_name" {
  description = "Use this for Cloud Monitoring alerts"
  value       = module.high_traffic_service.service_name
}

output "service_url" {
  value = module.high_traffic_service.url
}
```

**Use case:** Compute-intensive service with high memory requirements, long-running tasks, and predictable high traffic.

---

## Multi-Region Deployment

Deploy the same service in multiple regions for global availability.

```hcl
# Define regions for deployment
locals {
  regions = {
    us = "us-central1"
    eu = "europe-west1"
    asia = "asia-northeast1"
  }
}

# Deploy to multiple regions
module "global_api" {
  source   = "./iac/tf"
  for_each = local.regions

  metadata = {
    name = "global-api-${each.key}"
    env  = "production"
    labels = {
      region = each.key
    }
  }

  spec = {
    # Project ID using StringValueOrRef format
    project_id = {
      value = "global-prod-123"
    }
    region = each.value

    container = {
      image = {
        repo = "us-docker.pkg.dev/global-prod-123/api/service"
        tag  = "v1.0.0"
      }

      env = {
        variables = {
          REGION = each.key
        }
        secrets = {
          API_KEY = "projects/global-prod-123/secrets/api-key:latest"
        }
      }

      cpu    = 2
      memory = 1024
      replicas = {
        min = 3
        max = 100
      }
    }

    max_concurrency       = 100
    timeout_seconds       = 300
    allow_unauthenticated = true
  }
}

# Output all regional endpoints
output "regional_urls" {
  value = {
    for region, service in module.global_api : region => service.url
  }
}

# Example output:
# regional_urls = {
#   "asia" = "https://global-api-asia-abc123-an.a.run.app"
#   "eu"   = "https://global-api-eu-abc123-ew.a.run.app"
#   "us"   = "https://global-api-us-abc123-uc.a.run.app"
# }
```

**Use case:** Global service requiring low latency access from multiple geographic regions, with the same configuration deployed everywhere.

---

## Testing and Deployment

### Initialize and Validate

```bash
cd iac/tf
terraform init
terraform validate
terraform plan
```

### Deploy

```bash
terraform apply
```

### Verify Deployment

```bash
# Get service URL
terraform output url

# Test the service
curl $(terraform output -raw url)
```

### Clean Up

```bash
terraform destroy
```

---

## Best Practices

1. **Use specific image tags** instead of `latest` for reproducible deployments
2. **Store secrets in Secret Manager** rather than environment variables
3. **Set appropriate min_instances** for production services to avoid cold starts
4. **Configure VPC access** for services that need private resource access
5. **Use service accounts** with minimal required permissions
6. **Enable custom domains** for user-facing services
7. **Set realistic timeouts** based on your workload characteristics
8. **Monitor costs** by tracking instance hours and request counts
9. **Use labels** for resource organization and cost allocation
10. **Test with terraform plan** before applying changes to production

---

## Troubleshooting

### Service fails to start

Check container logs in Cloud Console and verify:
- Image is accessible and properly tagged
- Environment variables and secrets are correctly configured
- Container port matches the service configuration

### Domain mapping not working

Ensure:
- DNS managed zone exists and is correctly configured
- TXT record for verification is created
- Domain ownership is verified in Cloud Console

### VPC access issues

Verify:
- VPC network and subnet exist in the same project and region
- Service account has necessary permissions
- Egress setting matches your network requirements

---

## Additional Resources

- [Cloud Run Documentation](https://cloud.google.com/run/docs)
- [Terraform Google Provider](https://registry.terraform.io/providers/hashicorp/google/latest/docs)
- [Cloud Run Pricing](https://cloud.google.com/run/pricing)
- [Best Practices for Cloud Run](https://cloud.google.com/run/docs/tips)
