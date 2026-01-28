# Terraform Module Examples - GCP Cloud Storage Bucket

This document provides Terraform-specific deployment examples for the GcpGcsBucket component. These examples demonstrate how to use the Terraform module directly with `terraform` CLI.

## Table of Contents

- [Prerequisites](#prerequisites)
- [Basic Deployment](#basic-deployment)
- [Using the Module](#using-the-module)
- [Development Workflow](#development-workflow)
- [Testing and Debugging](#testing-and-debugging)
- [Common Patterns](#common-patterns)
- [Troubleshooting](#troubleshooting)

---

## Prerequisites

### Required Tools

```bash
# Install Terraform
# Visit: https://www.terraform.io/downloads

# Verify installation
terraform version  # Should be v1.0+

# Authenticate with GCP
gcloud auth application-default login
```

### GCP Setup

```bash
# Set your GCP project
export GCP_PROJECT=my-gcp-project-123

# Enable required APIs
gcloud services enable storage-api.googleapis.com \
  --project=${GCP_PROJECT}
```

---

## Basic Deployment

### Directory Structure

```
your-project/
├── main.tf                   # Root module configuration
├── variables.tf              # Input variables
├── outputs.tf                # Output definitions
├── terraform.tfvars          # Variable values
└── backend.tf                # Remote state configuration (optional)
```

### Step 1: Create Terraform Configuration

Create `main.tf`:

```hcl
terraform {
  required_version = ">= 1.0"
  required_providers {
    google = {
      source  = "hashicorp/google"
      version = "~> 5.0"
    }
  }
}

provider "google" {
  project = var.gcp_project_id
  region  = var.gcp_region
}

# Include the GcpGcsBucket module
module "gcs_bucket" {
  source = "../../iac/tf"  # Path to module

  bucket_name                         = var.bucket_name
  gcp_project_id                      = var.gcp_project_id
  location                            = var.location
  uniform_bucket_level_access_enabled = var.uniform_bucket_level_access_enabled
  storage_class                       = var.storage_class
  versioning_enabled                  = var.versioning_enabled
  lifecycle_rules                     = var.lifecycle_rules
  iam_bindings                        = var.iam_bindings
  gcp_labels                          = var.gcp_labels
}
```

Create `variables.tf`:

```hcl
variable "gcp_project_id" {
  description = "GCP project ID"
  type        = string
}

variable "gcp_region" {
  description = "GCP region for provider"
  type        = string
  default     = "us-east1"
}

variable "bucket_name" {
  description = "Name of the GCS bucket"
  type        = string
}

variable "location" {
  description = "Location for the bucket"
  type        = string
}

variable "uniform_bucket_level_access_enabled" {
  description = "Enable uniform bucket-level access"
  type        = bool
  default     = true
}

variable "storage_class" {
  description = "Storage class for the bucket"
  type        = string
  default     = "STANDARD"
}

variable "versioning_enabled" {
  description = "Enable object versioning"
  type        = bool
  default     = false
}

variable "lifecycle_rules" {
  description = "Lifecycle rules for the bucket"
  type = list(object({
    action = object({
      type          = string
      storage_class = optional(string)
    })
    condition = object({
      age_days             = optional(number)
      created_before       = optional(string)
      num_newer_versions   = optional(number)
      matches_storage_class = optional(list(string))
    })
  }))
  default = []
}

variable "iam_bindings" {
  description = "IAM bindings for the bucket"
  type = list(object({
    role      = string
    members   = list(string)
    condition = optional(string)
  }))
  default = []
}

variable "gcp_labels" {
  description = "Labels for the bucket"
  type        = map(string)
  default     = {}
}
```

Create `outputs.tf`:

```hcl
output "bucket_id" {
  description = "The ID of the created bucket"
  value       = module.gcs_bucket.bucket_id
}

output "bucket_url" {
  description = "The URL of the bucket"
  value       = module.gcs_bucket.bucket_url
}
```

Create `terraform.tfvars`:

```hcl
gcp_project_id                      = "my-gcp-project-123"
gcp_region                          = "us-east1"
bucket_name                         = "my-app-storage-prod"
location                            = "us-east1"
uniform_bucket_level_access_enabled = true
storage_class                       = "STANDARD"
versioning_enabled                  = true

lifecycle_rules = [
  {
    action = {
      type = "Delete"
    }
    condition = {
      num_newer_versions = 5
    }
  }
]

iam_bindings = [
  {
    role    = "roles/storage.objectAdmin"
    members = ["serviceAccount:app-backend@my-gcp-project-123.iam.gserviceaccount.com"]
  }
]

gcp_labels = {
  environment = "production"
  team        = "platform"
}
```

### Step 2: Initialize and Deploy

```bash
# Initialize Terraform (download providers)
terraform init

# Validate configuration
terraform validate

# Plan changes
terraform plan

# Apply changes
terraform apply
```

### Step 3: View Outputs

```bash
# Show all outputs
terraform output

# Show specific output
terraform output bucket_id
```

---

## Using the Module

### Direct Module Usage

If you're working directly in the module directory:

```bash
cd iac/tf

# Create terraform.tfvars with your configuration
cat > terraform.tfvars <<EOF
bucket_name = "test-bucket-dev"
gcp_project_id = "my-gcp-project-123"
location = "us-east1"
uniform_bucket_level_access_enabled = true
EOF

# Initialize and apply
terraform init
terraform plan
terraform apply
```

### Module as Source

Reference the module from a remote source:

```hcl
module "gcs_bucket" {
  source = "git::https://github.com/plantonhq/openmcf.git//apis/org/openmcf/provider/gcp/gcpgcsbucket/v1/iac/tf?ref=main"

  bucket_name    = "my-app-storage"
  gcp_project_id = "my-gcp-project-123"
  location       = "us-east1"
  # ... other variables
}
```

---

## Development Workflow

### State Management

Create `backend.tf` for remote state:

```hcl
terraform {
  backend "gcs" {
    bucket = "my-terraform-state-bucket"
    prefix = "gcs-buckets/prod"
  }
}
```

Initialize with backend:

```bash
terraform init -backend-config="bucket=my-terraform-state-bucket"
```

### Workspace Management

Use workspaces for environment separation:

```bash
# Create workspaces
terraform workspace new dev
terraform workspace new staging
terraform workspace new prod

# Switch workspaces
terraform workspace select dev

# List workspaces
terraform workspace list

# Deploy to selected workspace
terraform apply -var-file="env/dev.tfvars"
```

### Variable Files Per Environment

```
your-project/
├── main.tf
├── variables.tf
├── outputs.tf
└── env/
    ├── dev.tfvars
    ├── staging.tfvars
    └── prod.tfvars
```

Deploy to specific environment:

```bash
terraform apply -var-file="env/prod.tfvars"
```

---

## Testing and Debugging

### Enable Debug Logging

```bash
# Enable Terraform debug logging
export TF_LOG=DEBUG
export TF_LOG_PATH=terraform-debug.log

# Enable GCP provider logging
export GOOGLE_APPLICATION_CREDENTIALS=/path/to/service-account-key.json
```

### Validate Configuration

```bash
# Validate syntax
terraform validate

# Format code
terraform fmt -recursive

# Check for issues
terraform plan -out=tfplan
```

### Test with Terraform Console

```bash
# Start interactive console
terraform console

# Test variable interpolation
> var.bucket_name
"my-app-storage-prod"

# Test local values
> local.bucket_labels
```

### Dry Run Testing

```bash
# Create plan file
terraform plan -out=tfplan

# Inspect plan
terraform show tfplan

# Show plan in JSON
terraform show -json tfplan | jq .
```

---

## Common Patterns

### Multi-Environment Deployment

Create environment-specific tfvars:

`env/dev.tfvars`:

```hcl
gcp_project_id = "my-dev-project"
bucket_name    = "app-storage-dev"
location       = "us-east1"

versioning_enabled = false

lifecycle_rules = [
  {
    action = {
      type = "Delete"
    }
    condition = {
      age_days = 30  # Aggressive cleanup for dev
    }
  }
]

gcp_labels = {
  environment = "development"
  auto_cleanup = "enabled"
}
```

`env/prod.tfvars`:

```hcl
gcp_project_id = "my-prod-project"
bucket_name    = "app-storage-prod"
location       = "us-east1"

versioning_enabled = true

lifecycle_rules = [
  {
    action = {
      type = "Delete"
    }
    condition = {
      num_newer_versions = 10
    }
  }
]

gcp_labels = {
  environment = "production"
}
```

Deploy script:

```bash
#!/bin/bash
ENV=$1

if [ -z "$ENV" ]; then
  echo "Usage: ./deploy.sh [dev|staging|prod]"
  exit 1
fi

terraform workspace select $ENV || terraform workspace new $ENV
terraform plan -var-file="env/${ENV}.tfvars" -out=tfplan
terraform apply tfplan
```

### CI/CD Integration

Example GitLab CI configuration:

```yaml
variables:
  TF_VERSION: "1.5.0"
  TF_ROOT: ${CI_PROJECT_DIR}/iac/tf

.terraform:
  image:
    name: hashicorp/terraform:$TF_VERSION
    entrypoint: [""]
  before_script:
    - cd $TF_ROOT
    - echo $GCP_SERVICE_ACCOUNT_KEY | base64 -d > gcp-key.json
    - export GOOGLE_APPLICATION_CREDENTIALS=gcp-key.json
    - terraform init

validate:
  extends: .terraform
  stage: validate
  script:
    - terraform validate
    - terraform fmt -check -recursive

plan:
  extends: .terraform
  stage: plan
  script:
    - terraform plan -var-file="env/${CI_ENVIRONMENT_NAME}.tfvars" -out=tfplan
  artifacts:
    paths:
      - $TF_ROOT/tfplan

apply:
  extends: .terraform
  stage: deploy
  script:
    - terraform apply -input=false tfplan
  when: manual
  only:
    - main
```

### Conditional Resource Creation

Use `count` for conditional resources:

```hcl
# Create public bucket only if public_access is true
variable "public_access" {
  type    = bool
  default = false
}

resource "google_storage_bucket_iam_binding" "public_access" {
  count   = var.public_access ? 1 : 0
  bucket  = google_storage_bucket.bucket.name
  role    = "roles/storage.objectViewer"
  members = ["allUsers"]
}
```

### Dynamic Blocks

Use dynamic blocks for flexible configuration:

```hcl
resource "google_storage_bucket" "bucket" {
  name     = var.bucket_name
  location = var.location

  dynamic "lifecycle_rule" {
    for_each = var.lifecycle_rules
    content {
      action {
        type          = lifecycle_rule.value.action.type
        storage_class = lifecycle_rule.value.action.storage_class
      }
      condition {
        age                   = lifecycle_rule.value.condition.age_days
        created_before        = lifecycle_rule.value.condition.created_before
        num_newer_versions    = lifecycle_rule.value.condition.num_newer_versions
        matches_storage_class = lifecycle_rule.value.condition.matches_storage_class
      }
    }
  }
}
```

---

## Troubleshooting

### Common Errors

#### Error: Bucket Already Exists

```
Error: Error creating bucket: googleapi: Error 409: You already own this bucket.
```

**Solution:** Import existing bucket into state:

```bash
terraform import google_storage_bucket.bucket my-bucket-name
```

#### Error: Provider Configuration Not Found

```
Error: Provider configuration not found
```

**Solution:** Ensure provider block is defined:

```hcl
provider "google" {
  project = var.gcp_project_id
  region  = var.gcp_region
}
```

#### Error: State Lock

```
Error: Error acquiring the state lock
```

**Solution:** Force unlock (use with caution):

```bash
# Get lock ID from error message
terraform force-unlock <LOCK_ID>
```

### Debug Workflow

1. **Check state:**

```bash
terraform state list
terraform state show google_storage_bucket.bucket
```

2. **Refresh state:**

```bash
terraform refresh
```

3. **Check for drift:**

```bash
terraform plan -refresh-only
```

4. **Import existing resources:**

```bash
terraform import google_storage_bucket.bucket gs://bucket-name
```

5. **Remove resource from state:**

```bash
terraform state rm google_storage_bucket.bucket
```

### Performance Optimization

#### Parallelism

```bash
# Increase parallelism (default: 10)
terraform apply -parallelism=20
```

#### State Operations

```bash
# Pull remote state
terraform state pull > backup.tfstate

# Push local state
terraform state push backup.tfstate
```

### Upgrade Module

```bash
# Upgrade providers
terraform init -upgrade

# Show version changes
terraform version
terraform providers
```

---

## Further Reading

- [Terraform Documentation](https://www.terraform.io/docs/)
- [Google Provider Documentation](https://registry.terraform.io/providers/hashicorp/google/latest/docs)
- [Module README](README.md)
- [Component Examples](../../examples.md)

---

## Need Help?

For issues specific to Terraform deployment:
1. Check [Terraform Community](https://discuss.hashicorp.com/)
2. Review [Google Provider Issues](https://github.com/hashicorp/terraform-provider-google/issues)
3. Consult the [component examples](../../examples.md) for configuration patterns


