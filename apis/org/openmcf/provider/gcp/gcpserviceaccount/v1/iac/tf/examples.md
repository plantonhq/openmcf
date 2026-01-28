# GcpServiceAccount Terraform Examples

This document provides Terraform-specific examples for using the `GcpServiceAccount` module. These examples demonstrate
how to use the module in your Terraform configurations to create and manage GCP service accounts with various
configurations.

---

## Example 1: Minimal Service Account (No Key, Default Project)

Creates a service account in the specified project without generating a key or granting any IAM roles.

### Terraform Configuration

```hcl
module "analytics_sa" {
  source = "path/to/gcpserviceaccount/v1/iac/tf"

  metadata = {
    name = "analytics-sa"
  }

  spec = {
    service_account_id = "analytics-sa"
    project_id = {
      value = "my-project-123"
    }
  }
}

output "analytics_sa_email" {
  value = module.analytics_sa.email
}
```

### Explanation

- **service_account_id**: The unique identifier for the service account (6-30 characters)
- **project_id.value**: The GCP project where the service account will be created (using literal value)
- **create_key**: Not specified, defaults to `false` (no key generated)
- **Output**: The service account email can be used for IAM bindings or workload attachments

---

## Example 2: Service Account with Project-Level Roles and Key

Creates a service account in a specific project, generates a JSON key, and binds two project-level roles.

### Terraform Configuration

```hcl
module "logging_writer_sa" {
  source = "path/to/gcpserviceaccount/v1/iac/tf"

  metadata = {
    name = "logging-writer-sa"
  }

  spec = {
    service_account_id = "logging-writer"
    project_id = {
      value = "my-app-prod-1234"
    }
    create_key = true

    project_iam_roles = [
      "roles/logging.logWriter",
      "roles/monitoring.metricWriter"
    ]
  }
}

output "logging_writer_email" {
  value = module.logging_writer_sa.email
}

output "logging_writer_key" {
  value     = module.logging_writer_sa.key_base64
  sensitive = true
}
```

### Explanation

- **project_id.value**: Literal project ID value
- **create_key**: Set to `true` to generate a JSON private key
- **project_iam_roles**: List of IAM roles to grant at the project level
- **Outputs**:
    - `email`: The service account email
    - `key_base64`: The base64-encoded JSON key (marked as sensitive)

### Security Note

The generated key should be stored securely in a secret management system (e.g., Google Secret Manager, HashiCorp Vault)
and rotated regularly. Consider using keyless authentication patterns (Workload Identity) when possible.

---

## Example 3: Service Account with Organization-Level Roles

Creates a service account with both project-level and organization-level roles, without generating a key. Useful when
the service account must operate across multiple projects in the same organization.

### Terraform Configuration

```hcl
module "org_auditor_sa" {
  source = "path/to/gcpserviceaccount/v1/iac/tf"

  metadata = {
    name = "org-auditor-sa"
  }

  spec = {
    service_account_id = "org-auditor"
    project_id = {
      value = "shared-infra-5678"
    }
    org_id     = "123456789012"
    create_key = false

    project_iam_roles = [
      "roles/viewer"
    ]

    org_iam_roles = [
      "roles/resourcemanager.organizationViewer"
    ]
  }
}

output "org_auditor_email" {
  value = module.org_auditor_sa.email
}
```

### Explanation

- **project_id.value**: Literal project ID value
- **org_id**: The organization ID (required when using `org_iam_roles`)
- **project_iam_roles**: Roles granted at the project level
- **org_iam_roles**: Roles granted at the organization level
- **create_key**: Explicitly set to `false` to avoid key generation

### Use Case

This pattern is ideal for:

- Cross-project auditing and monitoring
- Organization-wide security scanning
- Centralized logging and metrics collection
- Infrastructure-as-Code automation across projects

---

## Example 4: Complete Configuration with Metadata

Demonstrates a full configuration using all available metadata fields for comprehensive resource tagging and
organization.

### Terraform Configuration

```hcl
module "ci_cd_sa" {
  source = "path/to/gcpserviceaccount/v1/iac/tf"

  metadata = {
    name = "ci-cd-automation"
    id   = "ci-cd-sa-prod"
    org  = "engineering"
    env  = "production"

    labels = {
      "team"        = "platform"
      "purpose"     = "ci-cd"
      "cost-center" = "engineering-ops"
    }

    tags = ["automation", "ci-cd", "production"]

    version = {
      id      = "v1.2.0"
      message = "Updated IAM roles for deployment automation"
    }
  }

  spec = {
    service_account_id = "ci-cd-automation"
    project_id = {
      value = "platform-prod-5678"
    }
    create_key = true

    project_iam_roles = [
      "roles/storage.admin",
      "roles/container.developer",
      "roles/logging.logWriter"
    ]
  }
}

output "ci_cd_sa_email" {
  value = module.ci_cd_sa.email
}

output "ci_cd_sa_key" {
  value     = module.ci_cd_sa.key_base64
  sensitive = true
}
```

### Explanation

- **metadata.labels**: Custom labels for resource organization and cost tracking
- **metadata.tags**: Tags for categorization and filtering
- **metadata.version**: Version tracking for infrastructure changes
- **Multiple IAM roles**: Grant multiple permissions in a single configuration

---

## Example 5: Using with Terraform Workspaces

Demonstrates how to use the module with Terraform workspaces to manage different environments.

### Terraform Configuration

```hcl
locals {
  env_config = {
    dev = {
      project_id = "my-app-dev-1234"
      roles = [
        "roles/logging.logWriter"
      ]
    }
    staging = {
      project_id = "my-app-staging-5678"
      roles = [
        "roles/logging.logWriter",
        "roles/monitoring.metricWriter"
      ]
    }
    prod = {
      project_id = "my-app-prod-9012"
      roles = [
        "roles/logging.logWriter",
        "roles/monitoring.metricWriter",
        "roles/storage.objectViewer"
      ]
    }
  }

  current_env = local.env_config[terraform.workspace]
}

module "app_service_account" {
  source = "path/to/gcpserviceaccount/v1/iac/tf"

  metadata = {
    name = "app-service-${terraform.workspace}"
    env  = terraform.workspace
  }

  spec = {
    service_account_id = "app-service-${terraform.workspace}"
    project_id = {
      value = local.current_env.project_id
    }
    create_key = terraform.workspace != "prod" # No keys in production

    project_iam_roles = local.current_env.roles
  }
}

output "service_account_email" {
  value = module.app_service_account.email
}
```

### Explanation

- **Workspace-specific configuration**: Different settings per environment
- **Dynamic IAM roles**: More roles in higher environments
- **Conditional key creation**: Keys only in non-production environments
- **Environment tagging**: Automatic environment labels via metadata

---

## Best Practices

### 1. Avoid Key Generation When Possible

```hcl
# ✅ Recommended: Use Workload Identity instead of keys
spec = {
  service_account_id = "my-app-sa"
  project_id = {
    value = var.project_id
  }
  create_key = false  # No key generated
}

# ❌ Avoid: Creating keys unless absolutely necessary
spec = {
  service_account_id = "my-app-sa"
  project_id = {
    value = var.project_id
  }
  create_key = true  # Only if legacy systems require it
}
```

### 2. Use Least Privilege Principle

```hcl
# ✅ Grant only the minimum required roles
project_iam_roles = [
  "roles/logging.logWriter",        # Only what's needed
  "roles/storage.objectViewer"      # Read-only when possible
]

# ❌ Avoid overly broad permissions
project_iam_roles = [
  "roles/editor"  # Too broad, grants unnecessary access
]
```

### 3. Store Keys Securely

```hcl
# ✅ Store key in Secret Manager
resource "google_secret_manager_secret" "sa_key" {
  secret_id = "app-service-account-key"
  
  replication {
    auto {}
  }
}

resource "google_secret_manager_secret_version" "sa_key_version" {
  secret      = google_secret_manager_secret.sa_key.id
  secret_data = module.app_sa.key_base64
}

# ❌ Never output keys to console or logs
output "key" {
  value = module.app_sa.key_base64  # This should always be sensitive = true
}
```

### 4. Document Service Account Purpose

```hcl
# ✅ Use descriptive names and metadata
metadata = {
  name = "logging-aggregator-prod"
  
  labels = {
    "purpose"     = "centralized-logging"
    "team"        = "sre"
    "environment" = "production"
  }
}

spec = {
  service_account_id = "logging-aggregator-prod"
  project_id = {
    value = var.project_id
  }
  # ...
}
```

---

## Validation

Before applying your Terraform configuration:

```bash
# Initialize Terraform
terraform init

# Validate the configuration
terraform validate

# Preview changes
terraform plan

# Apply changes
terraform apply
```

---

## Troubleshooting

### Error: service_account_id must be between 6 and 30 characters

**Solution**: Ensure your `service_account_id` meets GCP's naming requirements:

- Length: 6-30 characters
- Characters: lowercase letters, digits, hyphens
- Must start with a letter

### Error: org_id must be specified when org_iam_roles is not empty

**Solution**: When specifying `org_iam_roles`, you must also provide `org_id`:

```hcl
spec = {
  project_id = {
    value = "my-project-123"
  }
  org_id = "123456789012"  # Required
  org_iam_roles = [
    "roles/resourcemanager.organizationViewer"
  ]
}
```

### Error: Invalid IAM role

**Solution**: Verify the role name is valid. GCP roles follow the format `roles/<service>.<role>`. Check the [GCP IAM
roles documentation](https://cloud.google.com/iam/docs/understanding-roles) for valid role names.

---

## Additional Resources

- [GCP Service Account Documentation](https://cloud.google.com/iam/docs/service-accounts)
- [GCP IAM Roles](https://cloud.google.com/iam/docs/understanding-roles)
- [Workload Identity Federation](https://cloud.google.com/iam/docs/workload-identity-federation)
- [Service Account Best Practices](https://cloud.google.com/iam/docs/best-practices-service-accounts)
