# Terraform Module: Cloudflare R2 Bucket

This directory contains the Terraform module for deploying Cloudflare R2 buckets.

## Overview

The Terraform module provisions Cloudflare R2 buckets with S3-compatible object storage and zero egress fees. R2 is simpler than AWS S3 by design—no versioning, no bucket policies—optimized for storing and serving content efficiently.

## Module Structure

```
iac/tf/
├── README.md        # This file - deployment guide
├── variables.tf     # Input variables
├── provider.tf      # Cloudflare provider configuration
├── locals.tf        # Local variables and computed values
├── main.tf          # Resource definitions
└── outputs.tf       # Output values
```

## Prerequisites

1. **Cloudflare Account**:
   - Active Cloudflare account
   - R2 enabled (free tier available)
   - Cloudflare API token with permissions:
     - R2: Edit
     - Account Settings: Read

2. **Required Environment Variables**:
   ```bash
   export CLOUDFLARE_API_TOKEN="your-api-token-here"
   ```

3. **Terraform CLI**:
   ```bash
   # macOS
   brew install terraform

   # Linux
   wget https://releases.hashicorp.com/terraform/1.6.0/terraform_1.6.0_linux_amd64.zip
   unzip terraform_1.6.0_linux_amd64.zip
   sudo mv terraform /usr/local/bin/

   # Verify installation
   terraform version
   ```

## Usage

### Step 1: Create Terraform Configuration

Create a `main.tf` file using this module:

```hcl
terraform {
  required_providers {
    cloudflare = {
      source  = "cloudflare/cloudflare"
      version = "~> 4.0"
    }
  }
}

provider "cloudflare" {
  api_token = var.cloudflare_api_token
}

module "r2_bucket" {
  source = "./iac/tf"

  metadata = {
    name = "media-bucket"
    labels = {
      env     = "production"
      team    = "platform"
    }
  }

  spec = {
    bucket_name        = "myapp-media-assets"
    account_id         = "a1b2c3d4e5f6g7h8i9j0k1l2m3n4o5p6"  # Your account ID
    location           = 3  # WEUR (Western Europe)
    public_access      = true
    versioning_enabled = false
  }
}

output "bucket_name" {
  value = module.r2_bucket.bucket_name
}
```

### Step 2: Initialize Terraform

```bash
cd iac/tf

# Initialize providers and modules
terraform init
```

### Step 3: Plan Deployment

```bash
# Preview changes
terraform plan
```

**Expected Output**:

```
Terraform will perform the following actions:

  # cloudflare_r2_bucket.main will be created
  + resource "cloudflare_r2_bucket" "main" {
      + id         = (known after apply)
      + name       = "myapp-media-assets"
      + account_id = "a1b2c3d4e5f6g7h8i9j0k1l2m3n4o5p6"
      + location   = "WEUR"
    }

Plan: 1 to add, 0 to change, 0 to destroy.
```

### Step 4: Apply Configuration

```bash
# Apply changes
terraform apply

# Auto-approve (for CI/CD)
terraform apply -auto-approve
```

### Step 5: Verify Deployment

```bash
# View outputs
terraform output

# Test bucket access
aws s3 ls --endpoint-url https://<account-id>.r2.cloudflarestorage.com
```

## Input Variables

### Required Variables

#### `metadata` (object)

Metadata for the R2 bucket resource.

```hcl
metadata = {
  name = "media-bucket"  # Required
}
```

#### `spec` (object)

R2 bucket specification.

**Required fields**:

- **`bucket_name`** (string): DNS-compatible bucket name (3-63 characters)
  ```hcl
  bucket_name = "myapp-assets"
  ```

- **`account_id`** (string): Cloudflare account ID (32 hex characters)
  ```hcl
  account_id = "a1b2c3d4e5f6g7h8i9j0k1l2m3n4o5p6"
  ```

- **`location`** (number): Primary region for the bucket
  - `0` = auto (unspecified)
  - `1` = WNAM (Western North America)
  - `2` = ENAM (Eastern North America)
  - `3` = WEUR (Western Europe)
  - `4` = EEUR (Eastern Europe)
  - `5` = APAC (Asia-Pacific)
  - `6` = OC (Oceania)
  ```hcl
  location = 3  # WEUR
  ```

**Optional fields**:

- **`public_access`** (bool): Enable public access via r2.dev subdomain - Default: `false`
  ```hcl
  public_access = true
  ```
  **Note**: Terraform provider doesn't yet support toggling this. Must enable manually via Dashboard.

- **`versioning_enabled`** (bool): Enable object versioning - Default: `false`
  ```hcl
  versioning_enabled = false
  ```
  **Note**: R2 does not support versioning. This field is ignored.

## Outputs

| Output | Type | Description |
|--------|------|-------------|
| `bucket_name` | string | The name of the R2 bucket |
| `bucket_id` | string | The ID of the R2 bucket |
| `account_id` | string | The Cloudflare account ID |
| `location` | string | The location hint for the bucket |

Access outputs:

```bash
# View all outputs
terraform output

# Get specific output
terraform output bucket_name
```

Use outputs in other modules:

```hcl
module "other_module" {
  source = "./other-module"
  
  bucket_name = module.r2_bucket.bucket_name
}
```

## Examples

### Example 1: Basic Private Bucket

```hcl
module "private_bucket" {
  source = "./iac/tf"

  metadata = {
    name = "app-data"
  }

  spec = {
    bucket_name  = "myapp-private-data"
    account_id   = "a1b2c3d4e5f6g7h8i9j0k1l2m3n4o5p6"
    location     = 2  # ENAM
    public_access = false
  }
}
```

### Example 2: Public CDN Bucket

```hcl
module "cdn_bucket" {
  source = "./iac/tf"

  metadata = {
    name = "cdn-assets"
  }

  spec = {
    bucket_name  = "public-cdn-assets"
    account_id   = "a1b2c3d4e5f6g7h8i9j0k1l2m3n4o5p6"
    location     = 3  # WEUR
    public_access = true
  }
}
```

### Example 3: Multi-Region

```hcl
# US bucket
module "bucket_us" {
  source = "./iac/tf"

  metadata = {
    name = "assets-us"
  }

  spec = {
    bucket_name = "myapp-assets-us"
    account_id  = "a1b2c3d4e5f6g7h8i9j0k1l2m3n4o5p6"
    location    = 2  # ENAM
  }
}

# EU bucket
module "bucket_eu" {
  source = "./iac/tf"

  metadata = {
    name = "assets-eu"
  }

  spec = {
    bucket_name = "myapp-assets-eu"
    account_id  = "a1b2c3d4e5f6g7h8i9j0k1l2m3n4o5p6"
    location    = 3  # WEUR
  }
}
```

## State Management

### Local State (Development Only)

By default, Terraform stores state locally in `terraform.tfstate`:

```bash
terraform apply  # Creates terraform.tfstate
```

**Warning**: Do not use local state for production. State files contain sensitive data.

### Remote State (Recommended)

#### S3 Backend

```hcl
terraform {
  backend "s3" {
    bucket         = "my-terraform-state"
    key            = "cloudflare/r2-bucket/terraform.tfstate"
    region         = "us-east-1"
    encrypt        = true
    dynamodb_table = "terraform-lock"
  }
}
```

#### Terraform Cloud

```hcl
terraform {
  cloud {
    organization = "my-org"
    workspaces {
      name = "cloudflare-r2-bucket"
    }
  }
}
```

## Updating the Bucket

Modify your Terraform configuration and re-run:

```bash
terraform plan   # Preview changes
terraform apply  # Apply changes
```

### Common Updates

**Change location**:

```hcl
location = 5  # APAC instead of WEUR
```

**Enable public access**:

```hcl
public_access = true
```

**Note**: Changing bucket_name requires bucket replacement (destroy + create).

## Destroying the Bucket

```bash
# Preview what will be deleted
terraform plan -destroy

# Confirm and delete all resources
terraform destroy
```

**Warning**: This permanently deletes the bucket and all objects. Ensure data is backed up.

## CI/CD Integration

### GitHub Actions Example

```yaml
name: Deploy R2 Bucket

on:
  push:
    branches: [main]

jobs:
  terraform:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      
      - uses: hashicorp/setup-terraform@v2
        with:
          terraform_version: 1.6.0
      
      - name: Terraform Init
        run: terraform init
        working-directory: iac/tf
        env:
          CLOUDFLARE_API_TOKEN: ${{ secrets.CLOUDFLARE_API_TOKEN }}
      
      - name: Terraform Plan
        run: terraform plan
        working-directory: iac/tf
        env:
          CLOUDFLARE_API_TOKEN: ${{ secrets.CLOUDFLARE_API_TOKEN }}
      
      - name: Terraform Apply
        run: terraform apply -auto-approve
        working-directory: iac/tf
        env:
          CLOUDFLARE_API_TOKEN: ${{ secrets.CLOUDFLARE_API_TOKEN }}
```

## Troubleshooting

### Common Issues

**Issue**: `Error: authentication error - invalid API token`

**Solution**: Verify `CLOUDFLARE_API_TOKEN` environment variable:
```bash
echo $CLOUDFLARE_API_TOKEN
```

---

**Issue**: `Error: bucket already exists`

**Solution**: Bucket names must be unique within your account. Choose a different name.

---

**Issue**: Public access not working

**Solution**: The Terraform provider doesn't yet support toggling r2.dev public URLs. Enable manually via Cloudflare Dashboard.

---

**Issue**: `Error: resource already exists`

**Solution**: Import existing resource:
```bash
terraform import cloudflare_r2_bucket.main <bucket-id>
```

### Debug Mode

Enable detailed logging:

```bash
export TF_LOG=DEBUG
terraform apply
```

### View Resource State

```bash
# List all resources in state
terraform state list

# Show details of a specific resource
terraform state show cloudflare_r2_bucket.main
```

## Best Practices

1. **Use remote state backend** (S3, Terraform Cloud) for production
2. **Store API tokens in secrets** (GitHub Secrets, AWS Secrets Manager)
3. **Separate environments** (dev, staging, prod) using workspaces or separate state files
4. **Enable state locking** (DynamoDB for S3 backend) to prevent concurrent modifications
5. **Review plans carefully** before applying to production
6. **Tag resources** using `metadata.labels` for cost tracking and organization

## Terraform Workspaces

Manage multiple environments with workspaces:

```bash
# Create dev workspace
terraform workspace new dev
terraform apply

# Create prod workspace
terraform workspace new prod
terraform apply

# Switch between workspaces
terraform workspace select dev
terraform workspace select prod
```

Each workspace maintains separate state.

## Validation

Validate configuration before applying:

```bash
# Check syntax
terraform fmt -check

# Validate configuration
terraform validate

# Security scan (requires tfsec)
tfsec .
```

## Limitations

### Public Access

The Cloudflare Terraform provider does not yet expose a field for toggling the r2.dev public URL. When `public_access: true` is specified:
- Terraform creates the bucket
- Public access must be enabled manually via Cloudflare Dashboard or API
- See: https://developers.cloudflare.com/r2/buckets/public-buckets/

### Versioning

R2 does not support object versioning. The `versioning_enabled` field is ignored.

## Additional Resources

- [Cloudflare Terraform Provider Docs](https://registry.terraform.io/providers/cloudflare/cloudflare/latest/docs)
- [Cloudflare R2 API Docs](https://developers.cloudflare.com/api/operations/r2-create-bucket)
- [Component README](../../README.md) - User-facing documentation
- [Examples](../../examples.md) - Complete usage examples

## Support

For issues or questions:
1. Check [Common Issues](#common-issues) above
2. Review [Component README](../../README.md)
3. Consult Cloudflare and Terraform official documentation

---

**Ready to deploy?** Run `terraform init && terraform apply` to get started!

