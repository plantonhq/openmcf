# Terraform Examples - GCP Cert Manager Cert

This document provides examples of using the Terraform implementation directly.

## Basic Setup

### 1. Initialize

```bash
cd apis/project/planton/provider/gcp/gcpcertmanagercert/v1/iac/tf
terraform init
```

### 2. Create terraform.tfvars

```hcl
metadata = {
  name = "my-cert"
  id   = "cert-001"
  org  = "my-org"
  env = {
    id = "production"
  }
}

spec = {
  gcp_project_id         = "my-gcp-project"
  primary_domain_name    = "example.com"
  alternate_domain_names = []
  cloud_dns_zone_id = {
    value = "example-com-zone"
  }
  certificate_type  = 0
  validation_method = "DNS"
}
```

### 3. Deploy

```bash
terraform plan
terraform apply
```

## Example 1: Single Domain Certificate Manager

**terraform.tfvars**:
```hcl
metadata = {
  name = "single-domain-cert"
  id   = "cert-single-001"
  org  = "acme-corp"
  env = {
    id = "production"
  }
}

spec = {
  gcp_project_id         = "acme-production"
  primary_domain_name    = "acme.com"
  alternate_domain_names = []
  cloud_dns_zone_id = {
    value = "acme-com-zone"
  }
  certificate_type  = 0  # MANAGED
  validation_method = "DNS"
}
```

## Example 2: Multi-Domain Certificate

**terraform.tfvars**:
```hcl
metadata = {
  name = "multi-domain-cert"
  id   = "cert-multi-001"
  org  = "acme-corp"
  env = {
    id = "production"
  }
}

spec = {
  gcp_project_id      = "acme-production"
  primary_domain_name = "acme.com"
  alternate_domain_names = [
    "www.acme.com",
    "api.acme.com",
    "blog.acme.com"
  ]
  cloud_dns_zone_id = {
    value = "acme-com-zone"
  }
  certificate_type  = 0
  validation_method = "DNS"
}
```

## Example 3: Wildcard Certificate

**terraform.tfvars**:
```hcl
metadata = {
  name = "wildcard-cert"
  id   = "cert-wildcard-001"
  org  = "acme-corp"
  env = {
    id = "production"
  }
}

spec = {
  gcp_project_id      = "acme-production"
  primary_domain_name = "*.acme.com"
  alternate_domain_names = [
    "acme.com"  # Include apex domain too
  ]
  cloud_dns_zone_id = {
    value = "acme-com-zone"
  }
  certificate_type  = 0
  validation_method = "DNS"
}
```

## Example 4: Load Balancer Certificate

**terraform.tfvars**:
```hcl
metadata = {
  name = "lb-cert"
  id   = "cert-lb-001"
  org  = "acme-corp"
  env = {
    id = "production"
  }
}

spec = {
  gcp_project_id      = "acme-production"
  primary_domain_name = "lb.acme.com"
  alternate_domain_names = [
    "www.lb.acme.com"
  ]
  cloud_dns_zone_id = {
    value = "acme-com-zone"
  }
  certificate_type  = 1  # LOAD_BALANCER
  validation_method = "DNS"
}
```

## Example 5: Development Environment

**terraform.tfvars**:
```hcl
metadata = {
  name = "dev-cert"
  id   = "cert-dev-001"
  org  = "acme-corp"
  env = {
    id = "development"
  }
}

spec = {
  gcp_project_id      = "acme-development"
  primary_domain_name = "*.dev.acme.com"
  alternate_domain_names = [
    "dev.acme.com"
  ]
  cloud_dns_zone_id = {
    value = "dev-acme-com-zone"
  }
  certificate_type  = 0
  validation_method = "DNS"
}
```

## Example 6: Complex Multi-Wildcard

**terraform.tfvars**:
```hcl
metadata = {
  name = "complex-cert"
  id   = "cert-complex-001"
  org  = "acme-corp"
  env = {
    id = "production"
  }
}

spec = {
  gcp_project_id      = "acme-production"
  primary_domain_name = "*.services.acme.com"
  alternate_domain_names = [
    "services.acme.com",
    "*.api.acme.com",
    "api.acme.com"
  ]
  cloud_dns_zone_id = {
    value = "acme-com-zone"
  }
  certificate_type  = 0
  validation_method = "DNS"
}
```

## Using with Backend Configuration

### GCS Backend

Create `backend.tf`:

```hcl
terraform {
  backend "gcs" {
    bucket = "acme-terraform-state"
    prefix = "gcp-cert/production"
  }
}
```

Then:

```bash
terraform init
terraform apply
```

### Remote State with Workspaces

```bash
# Create workspaces for different environments
terraform workspace new production
terraform workspace new staging
terraform workspace new development

# Switch between environments
terraform workspace select production
terraform apply -var-file="production.tfvars"

terraform workspace select staging
terraform apply -var-file="staging.tfvars"
```

## Viewing Outputs

```bash
# View all outputs
terraform output

# View specific output
terraform output certificate-id
terraform output certificate-name

# Output in JSON format
terraform output -json
```

## Importing Existing Certificates

If you have existing certificates:

```bash
# Import Certificate Manager certificate
terraform import google_certificate_manager_certificate.cert[0] \
  projects/my-project/locations/global/certificates/my-cert

# Import Load Balancer certificate
terraform import google_compute_managed_ssl_certificate.lb_cert[0] \
  projects/my-project/global/sslCertificates/my-cert
```

## Environment-Specific Configurations

### production.tfvars

```hcl
metadata = {
  name = "prod-cert"
  id   = "cert-prod-001"
  org  = "acme"
  env  = { id = "production" }
}

spec = {
  gcp_project_id      = "acme-prod"
  primary_domain_name = "acme.com"
  alternate_domain_names = ["www.acme.com"]
  cloud_dns_zone_id   = { value = "acme-com-zone" }
  certificate_type    = 0
  validation_method   = "DNS"
}
```

### staging.tfvars

```hcl
metadata = {
  name = "staging-cert"
  id   = "cert-staging-001"
  org  = "acme"
  env  = { id = "staging" }
}

spec = {
  gcp_project_id      = "acme-staging"
  primary_domain_name = "staging.acme.com"
  alternate_domain_names = []
  cloud_dns_zone_id   = { value = "staging-acme-com-zone" }
  certificate_type    = 0
  validation_method   = "DNS"
}
```

Deploy with:

```bash
terraform apply -var-file="production.tfvars"
# or
terraform apply -var-file="staging.tfvars"
```

## Validation

After applying, verify:

```bash
# Check Terraform state
terraform show

# Verify in GCP
gcloud certificate-manager certificates list --project=my-project
gcloud dns record-sets list --zone=my-zone --project=my-project

# Check specific certificate
gcloud certificate-manager certificates describe my-cert \
  --location=global --project=my-project
```

## Clean Up

```bash
# Destroy all resources
terraform destroy

# Destroy with auto-approve (be careful!)
terraform destroy -auto-approve

# Destroy specific resource
terraform destroy -target=google_certificate_manager_certificate.cert
```

## Troubleshooting

### State Lock Issues

```bash
# Force unlock if needed
terraform force-unlock <lock-id>
```

### Refresh State

```bash
terraform refresh
```

### Import Drift

```bash
terraform plan  # Shows drift
terraform apply # Reconciles drift
```

## Best Practices

1. **Always use tfvars files** for different environments
2. **Store state remotely** in GCS
3. **Use workspaces** for environment separation
4. **Run plan before apply** to review changes
5. **Tag resources** using metadata for tracking
6. **Version control** your .tf files
7. **Gitignore** *.tfvars with sensitive data

## Next Steps

- See [README.md](README.md) for architecture details
- Check [overview.md](../overview.md) for Pulumi comparison
- Use OpenMCF CLI for production deployments

