# Terraform Examples for GCP DNS Zone

This document provides Terraform-specific examples for deploying the GCP DNS Zone module. These examples demonstrate how to use Terraform to manage Google Cloud DNS zones and records using the OpenMCF GCP DNS Zone module.

## Table of Contents

1. [Basic Usage with terraform.tfvars](#basic-usage-with-terraformtfvars)
2. [Minimal Configuration](#minimal-configuration)
3. [Production Configuration with Multiple Records](#production-configuration-with-multiple-records)
4. [Using Terraform Module](#using-terraform-module)
5. [Multi-Environment Setup](#multi-environment-setup)
6. [Integration with Other Terraform Resources](#integration-with-other-terraform-resources)
7. [Advanced: Dynamic Records from External Data](#advanced-dynamic-records-from-external-data)

---

## Basic Usage with terraform.tfvars

The simplest way to use this module is by creating a `terraform.tfvars` file with your configuration.

### Directory Structure

```
my-dns-config/
├── main.tf
├── terraform.tfvars
└── outputs.tf
```

### main.tf

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

module "gcp_dns_zone" {
  source = "path/to/gcpdnszone/v1/iac/tf"

  metadata = var.metadata
  spec     = var.spec
}

variable "metadata" {
  description = "Metadata for the resource"
  type = object({
    name   = string
    id     = optional(string)
    org    = optional(string)
    env    = optional(string)
    labels = optional(map(string))
    tags   = optional(list(string))
  })
}

variable "spec" {
  description = "Specification for the DNS zone"
  type = object({
    # project_id uses StringValueOrRef pattern
    project_id = object({
      value = string
    })
    iam_service_accounts = optional(list(string), [])
    records = optional(list(object({
      record_type = string
      name        = string
      values      = list(string)
      ttl_seconds = optional(number, 60)
    })), [])
  })
}

output "zone_id" {
  description = "The managed zone ID"
  value       = module.gcp_dns_zone.zone_id
}

output "zone_name" {
  description = "The managed zone name"
  value       = module.gcp_dns_zone.zone_name
}

output "nameservers" {
  description = "The nameservers for the managed zone"
  value       = module.gcp_dns_zone.nameservers
}
```

### terraform.tfvars

```hcl
metadata = {
  name = "example.com"
  env  = "production"
  org  = "mycompany"
  labels = {
    team        = "platform"
    cost-center = "infrastructure"
  }
}

spec = {
  # project_id uses StringValueOrRef pattern - wrap literal values in {value: "..."}
  project_id = {
    value = "my-gcp-project"
  }
  
  iam_service_accounts = [
    "cert-manager@my-gcp-project.iam.gserviceaccount.com",
    "external-dns@my-gcp-project.iam.gserviceaccount.com"
  ]
  
  records = [
    {
      record_type = "A"
      name        = "example.com."
      values      = ["104.198.14.52"]
      ttl_seconds = 300
    },
    {
      record_type = "CNAME"
      name        = "www.example.com."
      values      = ["example.com."]
      ttl_seconds = 300
    },
    {
      record_type = "TXT"
      name        = "example.com."
      values      = ["google-site-verification=abc123xyz789"]
      ttl_seconds = 3600
    }
  ]
}
```

### Deployment Commands

```bash
# Initialize Terraform
terraform init

# Review the execution plan
terraform plan

# Apply the configuration
terraform apply

# View outputs
terraform output

# Destroy resources (when needed)
terraform destroy
```

---

## Minimal Configuration

A minimal Terraform configuration for creating a DNS zone without any records.

### minimal.tfvars

```hcl
metadata = {
  name = "minimal-example.com"
}

spec = {
  project_id = {
    value = "my-gcp-project"
  }
}
```

**Usage**:
```bash
terraform apply -var-file="minimal.tfvars"
```

This creates a public DNS zone with DNSSEC enabled and no initial DNS records. Records can be added later by external tools like external-dns or manually through the GCP console.

---

## Production Configuration with Multiple Records

A comprehensive production configuration demonstrating various DNS record types.

### production.tfvars

```hcl
metadata = {
  name = "mycompany.io"
  env  = "production"
  org  = "mycompany"
  labels = {
    environment = "production"
    team        = "platform-engineering"
    managed-by  = "terraform"
    component   = "dns"
  }
  tags = ["production", "dns", "public-zone"]
}

spec = {
  project_id = {
    value = "mycompany-production"
  }
  
  iam_service_accounts = [
    "cert-manager@mycompany-production.iam.gserviceaccount.com",
    "external-dns@mycompany-production.iam.gserviceaccount.com"
  ]
  
  records = [
    # Root domain A record
    {
      record_type = "A"
      name        = "mycompany.io."
      values      = ["35.201.85.123"]
      ttl_seconds = 300
    },
    
    # IPv6 support
    {
      record_type = "AAAA"
      name        = "mycompany.io."
      values      = ["2600:1901:0:b::1"]
      ttl_seconds = 300
    },
    
    # WWW subdomain
    {
      record_type = "CNAME"
      name        = "www.mycompany.io."
      values      = ["mycompany.io."]
      ttl_seconds = 300
    },
    
    # API endpoint
    {
      record_type = "A"
      name        = "api.mycompany.io."
      values      = ["35.201.85.124"]
      ttl_seconds = 300
    },
    
    # MX records for email
    {
      record_type = "MX"
      name        = "mycompany.io."
      values      = [
        "10 aspmx.l.google.com.",
        "20 alt1.aspmx.l.google.com.",
        "30 alt2.aspmx.l.google.com."
      ]
      ttl_seconds = 3600
    },
    
    # SPF record
    {
      record_type = "TXT"
      name        = "mycompany.io."
      values      = ["v=spf1 include:_spf.google.com ~all"]
      ttl_seconds = 3600
    },
    
    # DMARC policy
    {
      record_type = "TXT"
      name        = "_dmarc.mycompany.io."
      values      = ["v=DMARC1; p=quarantine; rua=mailto:dmarc@mycompany.io; pct=100"]
      ttl_seconds = 3600
    },
    
    # Domain verification
    {
      record_type = "TXT"
      name        = "mycompany.io."
      values      = ["google-site-verification=AbC123XyZ789"]
      ttl_seconds = 3600
    }
  ]
}
```

---

## Using Terraform Module

Organize your Terraform code using modules for reusability and maintainability.

### Project Structure

```
infrastructure/
├── modules/
│   └── gcp-dns-zone/
│       ├── main.tf
│       ├── variables.tf
│       └── outputs.tf
├── environments/
│   ├── production/
│   │   ├── main.tf
│   │   └── terraform.tfvars
│   └── staging/
│       ├── main.tf
│       └── terraform.tfvars
└── terraform.tfvars
```

### modules/gcp-dns-zone/main.tf

```hcl
module "gcp_dns_zone" {
  source = "path/to/gcpdnszone/v1/iac/tf"

  metadata = var.metadata
  spec     = var.spec
}
```

### modules/gcp-dns-zone/variables.tf

```hcl
variable "metadata" {
  description = "Metadata for the DNS zone"
  type = object({
    name   = string
    env    = optional(string)
    org    = optional(string)
    labels = optional(map(string))
  })
}

variable "spec" {
  description = "DNS zone specification"
  type = object({
    # project_id uses StringValueOrRef pattern
    project_id = object({
      value = string
    })
    iam_service_accounts = optional(list(string), [])
    records = optional(list(object({
      record_type = string
      name        = string
      values      = list(string)
      ttl_seconds = optional(number, 60)
    })), [])
  })
}
```

### modules/gcp-dns-zone/outputs.tf

```hcl
output "zone_id" {
  description = "The managed zone ID"
  value       = module.gcp_dns_zone.zone_id
}

output "zone_name" {
  description = "The managed zone name"
  value       = module.gcp_dns_zone.zone_name
}

output "nameservers" {
  description = "List of nameservers for the zone"
  value       = module.gcp_dns_zone.nameservers
}
```

### environments/production/main.tf

```hcl
terraform {
  required_version = ">= 1.0"
  
  backend "gcs" {
    bucket = "mycompany-terraform-state"
    prefix = "dns/production"
  }
}

module "production_dns" {
  source = "../../modules/gcp-dns-zone"

  metadata = var.metadata
  spec     = var.spec
}

output "production_nameservers" {
  description = "Nameservers for production DNS zone"
  value       = module.production_dns.nameservers
}
```

---

## Multi-Environment Setup

Manage multiple environments (dev, staging, production) with environment-specific configurations.

### environments/dev/terraform.tfvars

```hcl
metadata = {
  name = "dev.mycompany.io"
  env  = "dev"
  org  = "mycompany"
  labels = {
    environment = "dev"
    team        = "platform"
  }
}

spec = {
  project_id = {
    value = "mycompany-dev"
  }
  
  iam_service_accounts = [
    "cert-manager@mycompany-dev.iam.gserviceaccount.com"
  ]
  
  records = [
    {
      record_type = "A"
      name        = "dev.mycompany.io."
      values      = ["34.120.45.10"]
      ttl_seconds = 60  # Shorter TTL for dev environment
    }
  ]
}
```

### environments/staging/terraform.tfvars

```hcl
metadata = {
  name = "staging.mycompany.io"
  env  = "staging"
  org  = "mycompany"
  labels = {
    environment = "staging"
    team        = "platform"
  }
}

spec = {
  project_id = {
    value = "mycompany-staging"
  }
  
  iam_service_accounts = [
    "cert-manager@mycompany-staging.iam.gserviceaccount.com",
    "external-dns@mycompany-staging.iam.gserviceaccount.com"
  ]
  
  records = [
    {
      record_type = "A"
      name        = "staging.mycompany.io."
      values      = ["34.120.45.20"]
      ttl_seconds = 300
    }
  ]
}
```

### environments/production/terraform.tfvars

```hcl
metadata = {
  name = "mycompany.io"
  env  = "production"
  org  = "mycompany"
  labels = {
    environment = "production"
    team        = "platform"
    compliance  = "required"
  }
}

spec = {
  project_id = {
    value = "mycompany-production"
  }
  
  iam_service_accounts = [
    "cert-manager@mycompany-production.iam.gserviceaccount.com",
    "external-dns@mycompany-production.iam.gserviceaccount.com"
  ]
  
  records = [
    {
      record_type = "A"
      name        = "mycompany.io."
      values      = ["35.201.85.123"]
      ttl_seconds = 300
    },
    {
      record_type = "CNAME"
      name        = "www.mycompany.io."
      values      = ["mycompany.io."]
      ttl_seconds = 300
    }
  ]
}
```

### Workspace-Based Deployment

```bash
# Create workspaces for each environment
terraform workspace new dev
terraform workspace new staging
terraform workspace new production

# Deploy to dev
terraform workspace select dev
terraform apply -var-file="environments/dev/terraform.tfvars"

# Deploy to staging
terraform workspace select staging
terraform apply -var-file="environments/staging/terraform.tfvars"

# Deploy to production
terraform workspace select production
terraform apply -var-file="environments/production/terraform.tfvars"
```

---

## Integration with Other Terraform Resources

Demonstrate integration with other GCP resources and dynamic configuration.

### main.tf

```hcl
# Create GKE cluster
resource "google_container_cluster" "primary" {
  name     = "production-gke"
  location = "us-central1"
  project  = var.gcp_project_id
  
  # ... other GKE configuration ...
}

# Get the load balancer IP after GKE creates the Ingress
data "google_compute_address" "ingress_ip" {
  name    = "gke-ingress-ip"
  project = var.gcp_project_id
  region  = "us-central1"
  
  depends_on = [google_container_cluster.primary]
}

# Create DNS zone with A record pointing to the load balancer
module "dns_zone" {
  source = "path/to/gcpdnszone/v1/iac/tf"

  metadata = {
    name = "myapp.io"
    env  = "production"
    labels = {
      component = "dns"
      cluster   = google_container_cluster.primary.name
    }
  }

  spec = {
    project_id = {
      value = var.gcp_project_id
    }
    
    iam_service_accounts = [
      google_service_account.cert_manager.email,
      google_service_account.external_dns.email
    ]
    
    records = [
      {
        record_type = "A"
        name        = "myapp.io."
        values      = [data.google_compute_address.ingress_ip.address]
        ttl_seconds = 300
      },
      {
        record_type = "CNAME"
        name        = "www.myapp.io."
        values      = ["myapp.io."]
        ttl_seconds = 300
      }
    ]
  }
}

# Create service accounts for DNS automation
resource "google_service_account" "cert_manager" {
  account_id   = "cert-manager"
  display_name = "cert-manager DNS automation"
  project      = var.gcp_project_id
}

resource "google_service_account" "external_dns" {
  account_id   = "external-dns"
  display_name = "external-dns automation"
  project      = var.gcp_project_id
}

# Output nameservers for registrar configuration
output "nameservers" {
  description = "Configure these nameservers at your domain registrar"
  value       = module.dns_zone.nameservers
}

output "zone_id" {
  description = "DNS Zone ID"
  value       = module.dns_zone.zone_id
}
```

---

## Advanced: Dynamic Records from External Data

Generate DNS records dynamically based on external data sources or computations.

### main.tf

```hcl
# Read list of backend services from a JSON file
locals {
  backend_services = jsondecode(file("${path.module}/backend-services.json"))
  
  # Generate A records for each backend service
  backend_dns_records = [
    for service in local.backend_services : {
      record_type = "A"
      name        = "${service.subdomain}.myapp.io."
      values      = [service.ip_address]
      ttl_seconds = 300
    }
  ]
  
  # Static records
  static_records = [
    {
      record_type = "A"
      name        = "myapp.io."
      values      = ["35.201.85.100"]
      ttl_seconds = 300
    },
    {
      record_type = "TXT"
      name        = "myapp.io."
      values      = ["v=spf1 include:_spf.google.com ~all"]
      ttl_seconds = 3600
    }
  ]
  
  # Combine static and dynamic records
  all_records = concat(local.static_records, local.backend_dns_records)
}

module "dns_zone" {
  source = "path/to/gcpdnszone/v1/iac/tf"

  metadata = {
    name = "myapp.io"
    env  = "production"
  }

  spec = {
    project_id = {
      value = var.gcp_project_id
    }
    records    = local.all_records
  }
}
```

### backend-services.json

```json
[
  {
    "name": "api",
    "subdomain": "api",
    "ip_address": "35.201.85.101"
  },
  {
    "name": "admin",
    "subdomain": "admin",
    "ip_address": "35.201.85.102"
  },
  {
    "name": "monitoring",
    "subdomain": "metrics",
    "ip_address": "35.201.85.103"
  }
]
```

This approach is useful when:
- DNS records are generated from configuration management systems
- You have a large number of similar records (e.g., customer subdomains)
- Records need to be synchronized from an external source of truth

---

## Using Terraform Remote State

Share DNS zone outputs across multiple Terraform configurations.

### dns-infrastructure/main.tf (DNS Zone Configuration)

```hcl
terraform {
  backend "gcs" {
    bucket = "mycompany-terraform-state"
    prefix = "dns/production"
  }
}

module "dns_zone" {
  source = "path/to/gcpdnszone/v1/iac/tf"

  metadata = {
    name = "mycompany.io"
  }

  spec = {
    project_id = {
      value = "mycompany-production"
    }
    records = [
      {
        record_type = "A"
        name        = "mycompany.io."
        values      = ["35.201.85.123"]
        ttl_seconds = 300
      }
    ]
  }
}

output "zone_name" {
  value = module.dns_zone.zone_name
}

output "nameservers" {
  value = module.dns_zone.nameservers
}
```

### application-infrastructure/main.tf (Application Configuration)

```hcl
# Reference DNS zone from remote state
data "terraform_remote_state" "dns" {
  backend = "gcs"
  config = {
    bucket = "mycompany-terraform-state"
    prefix = "dns/production"
  }
}

# Use nameservers from DNS zone in other resources
resource "google_dns_record_set" "app_record" {
  # This is just an example; typically external-dns would manage app records
  name         = "app.mycompany.io."
  type         = "A"
  ttl          = 300
  managed_zone = data.terraform_remote_state.dns.outputs.zone_name
  rrdatas      = [google_compute_address.app.address]
}

output "configured_nameservers" {
  description = "Nameservers configured for this domain"
  value       = data.terraform_remote_state.dns.outputs.nameservers
}
```

---

## Terraform Best Practices

### 1. Use Remote State

Always use remote state (GCS, S3, Terraform Cloud) for production environments:

```hcl
terraform {
  backend "gcs" {
    bucket  = "mycompany-terraform-state"
    prefix  = "dns/production"
    
    # Enable state locking
    encryption_key = "your-kms-key"
  }
}
```

### 2. Use Variables for Sensitive Data

Never hardcode sensitive values. Use variables and external secret management:

```hcl
variable "gcp_project_id" {
  description = "GCP project ID"
  type        = string
  sensitive   = true
}
```

### 3. Version Lock Your Providers

Always specify provider versions to ensure reproducibility:

```hcl
terraform {
  required_version = ">= 1.5"
  
  required_providers {
    google = {
      source  = "hashicorp/google"
      version = "~> 5.0"
    }
  }
}
```

### 4. Use Meaningful Names and Labels

Always include descriptive metadata:

```hcl
metadata = {
  name = "production.mycompany.io"
  env  = "production"
  labels = {
    team        = "platform-engineering"
    cost-center = "infrastructure"
    managed-by  = "terraform"
    component   = "dns"
  }
}
```

### 5. Document Outputs

Provide clear descriptions for outputs:

```hcl
output "nameservers" {
  description = "Add these nameservers to your domain registrar"
  value       = module.dns_zone.nameservers
}
```

---

## Troubleshooting

### Issue: "Error creating ManagedZone"

**Solution**: Verify the GCP project ID is correct and you have appropriate permissions:

```bash
gcloud projects list
gcloud auth application-default login
```

### Issue: "Nameservers not propagating"

**Solution**: 
1. Retrieve nameservers from Terraform output:
   ```bash
   terraform output nameservers
   ```
2. Add these nameservers to your domain registrar
3. Wait up to 48 hours for full propagation (typically 5-10 minutes)
4. Verify with:
   ```bash
   dig NS yourdomain.com
   ```

### Issue: "IAM binding errors"

**Solution**: Ensure service accounts exist before referencing them:

```hcl
# Create service account first
resource "google_service_account" "cert_manager" {
  account_id = "cert-manager"
  project    = var.gcp_project_id
}

# Then reference in DNS module
module "dns_zone" {
  source = "..."
  
  spec = {
    iam_service_accounts = [
      google_service_account.cert_manager.email
    ]
  }
  
  depends_on = [google_service_account.cert_manager]
}
```

---

## Summary

These Terraform examples demonstrate:

- **Basic Usage**: Simple terraform.tfvars configuration
- **Modular Approach**: Reusable modules for multiple environments
- **Integration**: Combining DNS with other GCP resources
- **Dynamic Configuration**: Generating records from external data
- **Best Practices**: Remote state, versioning, and proper organization

For more general DNS configuration examples (not Terraform-specific), see the main [examples.md](../../examples.md) in the component root directory.

