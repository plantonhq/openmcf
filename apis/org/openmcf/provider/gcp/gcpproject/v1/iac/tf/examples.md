# Terraform Examples - GCP Project

This document provides Terraform HCL examples for using the GCP Project module.

## Important Notes

- **project_id** is now a required field and must be specified in the spec
- The **add_suffix** field (optional, defaults to false) controls whether a random 3-character suffix is appended to the project_id
- When add_suffix is false (default), the project_id is used as-is
- When add_suffix is true, a random suffix like "-xyz" is appended for uniqueness

---

## Example 1: Minimal Development Project

```hcl
module "dev_sandbox" {
  source = "../../tf"

  metadata = {
    name = "dev-sandbox"
    org  = "example-org"
    env  = "dev"
  }

  spec = {
    project_id              = "dev-sandbox-proj"
    parent_type             = "folder"
    parent_id               = "123456789012"
    billing_account_id      = "ABCDEF-123456-ABCDEF"
    labels = {
      env  = "dev"
      team = "research"
    }
    disable_default_network = false  # Keep default network for quick iteration
    enabled_apis = [
      "compute.googleapis.com",
      "storage.googleapis.com"
    ]
    owner_member = "alice@example.com"
  }
}

output "dev_project_id" {
  value = module.dev_sandbox.project_id
}

output "dev_project_number" {
  value = module.dev_sandbox.project_number
}
```

---

## Example 2: Standard Project with Essential APIs

```hcl
module "staging_api_service" {
  source = "../../tf"

  metadata = {
    name = "staging-api-service"
    id   = "stg-api-svc-001"
    org  = "example-org"
    env  = "staging"
    labels = {
      component = "api-service"
    }
  }

  spec = {
    project_id              = "staging-api-svc"
    parent_type             = "folder"
    parent_id               = "234567890123"
    billing_account_id      = "ABCDEF-123456-ABCDEF"
    labels = {
      env       = "staging"
      team      = "backend"
      component = "api-service"
    }
    disable_default_network = true  # Security best practice
    enabled_apis = [
      "compute.googleapis.com",
      "storage.googleapis.com",
      "container.googleapis.com",
      "logging.googleapis.com",
      "monitoring.googleapis.com",
      "iam.googleapis.com",
      "iamcredentials.googleapis.com"
    ]
    owner_member = "devops-staging@example.com"
  }
}

output "staging_project_id" {
  value       = module.staging_api_service.project_id
  description = "Staging project ID"
}
```

---

## Example 3: Production-Grade Project with Full Security

```hcl
module "prod_payment_processing" {
  source = "../../tf"

  metadata = {
    name = "prod-payment-processing"
    id   = "prd-payment-001"
    org  = "example-org"
    env  = "prod"
    labels = {
      criticality = "high"
      compliance  = "pci-dss"
    }
  }

  spec = {
    project_id              = "prod-payment-proc"
    parent_type             = "folder"
    parent_id               = "345678901234"
    billing_account_id      = "ABCDEF-123456-ABCDEF"
    labels = {
      env         = "prod"
      team        = "e-commerce"
      component   = "payment-processing"
      cost-center = "cc-4510"
      compliance  = "pci-dss"
      criticality = "high"
    }
    disable_default_network = true  # CRITICAL: Disable insecure default network
    delete_protection       = true  # CRITICAL: Prevent accidental project deletion
    enabled_apis = [
      "compute.googleapis.com",
      "storage.googleapis.com",
      "container.googleapis.com",
      "logging.googleapis.com",
      "monitoring.googleapis.com",
      "cloudtrace.googleapis.com",
      "cloudprofiler.googleapis.com",
      "iam.googleapis.com",
      "iamcredentials.googleapis.com",
      "servicenetworking.googleapis.com",
      "dns.googleapis.com",
      "secretmanager.googleapis.com"
    ]
    owner_member = "platform-admins@example.com"
  }
}

output "prod_project_id" {
  value       = module.prod_payment_processing.project_id
  description = "Production payment processing project ID"
  sensitive   = false
}

output "prod_project_number" {
  value       = module.prod_payment_processing.project_number
  description = "Production payment processing project number"
}
```

---

## Example 4: Project Under Organization (Not Folder)

```hcl
module "shared_networking" {
  source = "../../tf"

  metadata = {
    name = "shared-networking"
    org  = "example-org"
    env  = "shared"
  }

  spec = {
    project_id              = "shared-networking"
    parent_type             = "organization"
    parent_id               = "987654321098"
    billing_account_id      = "ABCDEF-123456-ABCDEF"
    labels = {
      env       = "shared"
      team      = "platform"
      component = "networking"
    }
    disable_default_network = true
    enabled_apis = [
      "compute.googleapis.com",
      "dns.googleapis.com",
      "servicenetworking.googleapis.com"
    ]
    owner_member = "network-admins@example.com"
  }
}
```

---

## Example 4b: Project with Random Suffix (add_suffix = true)

This example demonstrates using add_suffix to automatically append a random 3-character 
suffix to the project_id. Useful for testing or temporary projects where uniqueness must be guaranteed.

```hcl
module "test_project" {
  source = "../../tf"

  metadata = {
    name = "test-ephemeral"
    org  = "example-org"
    env  = "test"
  }

  spec = {
    project_id              = "test-proj"
    add_suffix              = true  # Will create project_id like "test-proj-abc"
    parent_type             = "folder"
    parent_id               = "987654321098"
    billing_account_id      = "ABCDEF-123456-ABCDEF"
    labels = {
      env  = "test"
      team = "qa"
    }
    disable_default_network = false
    enabled_apis = [
      "compute.googleapis.com"
    ]
  }
}

output "actual_project_id" {
  value       = module.test_project.project_id
  description = "The actual project ID with suffix (e.g., test-proj-xyz)"
}
```

---

## Example 5: Multiple Projects with for_each

```hcl
locals {
  projects = {
    dev = {
      parent_id     = "123456789012"
      apis          = ["compute.googleapis.com", "storage.googleapis.com"]
      owner         = "group:dev-team@example.com"
      disable_default_network = false
    }
    staging = {
      parent_id     = "234567890123"
      apis          = ["compute.googleapis.com", "storage.googleapis.com", "container.googleapis.com"]
      owner         = "group:staging-team@example.com"
      disable_default_network = true
    }
    prod = {
      parent_id     = "345678901234"
      apis = [
        "compute.googleapis.com",
        "storage.googleapis.com",
        "container.googleapis.com",
        "logging.googleapis.com",
        "monitoring.googleapis.com"
      ]
      owner         = "group:prod-admins@example.com"
      disable_default_network = true
    }
  }
}

module "projects" {
  source   = "../../tf"
  for_each = local.projects

  metadata = {
    name = "my-app-${each.key}"
    org  = "example-org"
    env  = each.key
  }

  spec = {
    parent_type             = "folder"
    parent_id               = each.value.parent_id
    billing_account_id      = "ABCDEF-123456-ABCDEF"
    labels = {
      env       = each.key
      team      = "platform"
      component = "my-app"
    }
    disable_default_network = each.value.disable_default_network
    enabled_apis            = each.value.apis
    owner_member            = each.value.owner
  }
}

output "all_project_ids" {
  value = {
    for env, module in module.projects : env => module.project_id
  }
}
```

---

## Example 6: Data Science Project with BigQuery

```hcl
module "ml_research" {
  source = "../../tf"

  metadata = {
    name = "ml-research"
    org  = "example-org"
    env  = "dev"
    labels = {
      workload = "ml"
    }
  }

  spec = {
    parent_type             = "folder"
    parent_id               = "567890123456"
    billing_account_id      = "ABCDEF-123456-ABCDEF"
    labels = {
      env       = "dev"
      team      = "data-science"
      component = "ml-research"
    }
    disable_default_network = false
    enabled_apis = [
      "compute.googleapis.com",
      "storage.googleapis.com",
      "bigquery.googleapis.com",
      "bigquerystorage.googleapis.com",
      "aiplatform.googleapis.com",
      "notebooks.googleapis.com",
      "ml.googleapis.com",
      "dataflow.googleapis.com"
    ]
    owner_member = "group:data-scientists@example.com"
  }
}
```

---

## Example 7: Service Account as Owner

```hcl
module "ci_automation" {
  source = "../../tf"

  metadata = {
    name = "ci-automation"
    org  = "example-org"
    env  = "shared"
  }

  spec = {
    parent_type             = "folder"
    parent_id               = "678901234567"
    billing_account_id      = "ABCDEF-123456-ABCDEF"
    labels = {
      env       = "shared"
      team      = "platform"
      component = "ci-cd"
    }
    disable_default_network = true
    enabled_apis = [
      "compute.googleapis.com",
      "storage.googleapis.com",
      "cloudbuild.googleapis.com",
      "containerregistry.googleapis.com"
    ]
    owner_member = "serviceAccount:ci-automation@my-seed-project.iam.gserviceaccount.com"
  }
}
```

---

## Example 8: Project with Dynamic APIs from Variable

```hcl
variable "environment" {
  type        = string
  description = "Environment name (dev, staging, prod)"
}

variable "required_apis" {
  type        = list(string)
  description = "List of GCP APIs to enable"
  default = [
    "compute.googleapis.com",
    "storage.googleapis.com"
  ]
}

variable "additional_prod_apis" {
  type        = list(string)
  description = "Additional APIs for production"
  default = [
    "logging.googleapis.com",
    "monitoring.googleapis.com",
    "cloudtrace.googleapis.com"
  ]
}

locals {
  all_apis = var.environment == "prod" ? concat(
    var.required_apis,
    var.additional_prod_apis
  ) : var.required_apis
}

module "dynamic_project" {
  source = "../../tf"

  metadata = {
    name = "my-app-${var.environment}"
    org  = "example-org"
    env  = var.environment
  }

  spec = {
    parent_type             = "folder"
    parent_id               = "789012345678"
    billing_account_id      = "ABCDEF-123456-ABCDEF"
    labels = {
      env  = var.environment
      team = "platform"
    }
    disable_default_network = var.environment != "dev"
    enabled_apis            = local.all_apis
    owner_member            = "group:platform-team@example.com"
  }
}
```

---

## Example 9: Backend Configuration with Remote State

```hcl
terraform {
  backend "gcs" {
    bucket = "my-terraform-state-bucket"
    prefix = "projects/gcp-project"
  }
}

module "prod_project" {
  source = "../../tf"

  metadata = {
    name = "prod-global-app"
    id   = "prd-global-001"
    org  = "example-org"
    env  = "prod"
  }

  spec = {
    parent_type             = "folder"
    parent_id               = "789012345678"
    billing_account_id      = "ABCDEF-123456-ABCDEF"
    labels = {
      env       = "prod"
      team      = "platform"
      component = "global-app"
      region    = "multi-region"
    }
    disable_default_network = true
    enabled_apis = [
      "compute.googleapis.com",
      "storage.googleapis.com",
      "container.googleapis.com",
      "sqladmin.googleapis.com",
      "redis.googleapis.com",
      "servicenetworking.googleapis.com",
      "cloudloadbalancing.googleapis.com",
      "logging.googleapis.com",
      "monitoring.googleapis.com"
    ]
    owner_member = "group:sre-team@example.com"
  }
}

output "project_id" {
  value = module.prod_project.project_id
}
```

---

## Common Patterns

### Pattern 1: Conditional API Enablement

```hcl
locals {
  base_apis = [
    "compute.googleapis.com",
    "storage.googleapis.com"
  ]
  
  observability_apis = var.enable_observability ? [
    "logging.googleapis.com",
    "monitoring.googleapis.com",
    "cloudtrace.googleapis.com"
  ] : []
  
  all_apis = concat(local.base_apis, local.observability_apis)
}

spec = {
  # ...
  enabled_apis = local.all_apis
}
```

### Pattern 2: Environment-Specific Configuration

```hcl
locals {
  env_config = {
    dev = {
      disable_network = false
      owner          = "group:dev-team@example.com"
    }
    prod = {
      disable_network = true
      owner          = "group:sre-team@example.com"
    }
  }
  
  config = local.env_config[var.environment]
}

spec = {
  # ...
  disable_default_network = local.config.disable_network
  owner_member            = local.config.owner
}
```

### Pattern 3: Label Standardization

```hcl
locals {
  standard_labels = {
    managed-by = "terraform"
    cost-center = var.cost_center
    team        = var.team_name
  }
  
  custom_labels = var.additional_labels
  
  all_labels = merge(local.standard_labels, local.custom_labels)
}

spec = {
  # ...
  labels = local.all_labels
}
```

---

## Best Practices

### 1. Use Variables for Environment-Specific Values

```hcl
variable "folder_ids" {
  type = map(string)
  default = {
    dev     = "123456789012"
    staging = "234567890123"
    prod    = "345678901234"
  }
}
```

### 2. Store State Remotely

```hcl
terraform {
  backend "gcs" {
    bucket = "my-terraform-state"
    prefix = "projects"
  }
}
```

### 3. Use Outputs for Cross-Stack References

```hcl
output "project_id" {
  value       = module.my_project.project_id
  description = "The project ID for use in other modules"
}
```

### 4. Group Related Projects

```hcl
module "app_projects" {
  source   = "../../tf"
  for_each = toset(["dev", "staging", "prod"])
  
  # ...configuration...
}
```

---

## Module Outputs Reference

All outputs from the module:

```hcl
output "project_id" {
  description = "The unique project ID (globally unique across all GCP)"
  value       = module.my_project.project_id
}

output "project_number" {
  description = "The numeric identifier of the project (assigned by Google)"
  value       = module.my_project.project_number
}

output "project_name" {
  description = "The display name of the project"
  value       = module.my_project.project_name
}

output "enabled_apis" {
  description = "List of APIs that were enabled on the project"
  value       = module.my_project.enabled_apis
}
```

---

## Troubleshooting

### Issue: "Billing account not found"

**Solution:** Verify the billing account ID format is correct: `XXXXXX-XXXXXX-XXXXXX`

```hcl
billing_account_id = "ABCDEF-123456-ABCDEF"  # Correct format
```

### Issue: "Project ID already exists"

**Solution:** The random suffix should prevent this, but if it happens, run `terraform apply` again to generate a new suffix.

### Issue: "API not enabled"

**Solution:** Ensure you're enabling the required APIs:

```hcl
enabled_apis = [
  "compute.googleapis.com",
  "iam.googleapis.com",
  # Add other required APIs
]
```

---

## Next Steps

After creating projects with Terraform:

1. **Import existing projects** if migrating:
   ```bash
   terraform import module.my_project.google_project.this my-project-id
   ```

2. **Plan before apply** to see what will change:
   ```bash
   terraform plan
   ```

3. **Use workspaces** for managing multiple environments:
   ```bash
   terraform workspace new prod
   ```

4. **Set up CI/CD** for automated project provisioning:
   - Store state in GCS
   - Use Workload Identity for authentication
   - Run `terraform plan` on PRs
   - Auto-apply on merge to main

