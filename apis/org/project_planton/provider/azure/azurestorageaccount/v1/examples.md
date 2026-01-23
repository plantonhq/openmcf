# Azure Storage Account Examples

This document provides comprehensive examples for the `AzureStorageAccount` API resource, demonstrating various storage configurations in Microsoft Azure.

## Table of Contents

1. [Minimal Configuration](#minimal-configuration)
2. [Development Environment](#development-environment)
3. [Production Environment](#production-environment)
4. [Static Website Hosting](#static-website-hosting)
5. [Data Lake Configuration](#data-lake-configuration)
6. [Network-Isolated Configuration](#network-isolated-configuration)

---

## Minimal Configuration

The simplest possible configuration with only required fields.

```yaml
apiVersion: azure.project-planton.org/v1
kind: AzureStorageAccount
metadata:
  name: myapp-storage
spec:
  region: eastus
  resourceGroup: myapp-rg
```

**What Gets Created**:
- StorageV2 general-purpose account
- Standard tier (HDD-backed)
- Locally redundant storage (LRS)
- Hot access tier
- HTTPS-only traffic
- TLS 1.2 minimum
- Network ACLs default to Deny with Azure Services bypass
- 7-day soft delete retention

---

## Development Environment

Configuration optimized for development with relaxed security and cost optimization.

```yaml
apiVersion: azure.project-planton.org/v1
kind: AzureStorageAccount
metadata:
  name: myapp-dev-storage
  org: mycompany
  env: development
spec:
  region: eastus
  resourceGroup: dev-rg
  accountKind: STORAGE_V2
  accountTier: STANDARD
  replicationType: LRS
  accessTier: HOT
  enableHttpsTrafficOnly: true
  minTlsVersion: TLS1_2
  networkRules:
    defaultAction: ALLOW  # More permissive for dev
    bypassAzureServices: true
  blobProperties:
    enableVersioning: false
    softDeleteRetentionDays: 7
    containerSoftDeleteRetentionDays: 7
  containers:
    - name: data
      accessType: PRIVATE
    - name: logs
      accessType: PRIVATE
```

**Use Case**: Development environment with easy access and minimal costs.

---

## Production Environment

Standard production configuration with enhanced durability and security.

```yaml
apiVersion: azure.project-planton.org/v1
kind: AzureStorageAccount
metadata:
  name: myapp-prod-storage
  id: azsa-prod-001
  org: mycompany
  env: production
spec:
  region: eastus
  resourceGroup: prod-storage-rg
  accountKind: STORAGE_V2
  accountTier: STANDARD
  replicationType: GRS  # Geo-redundant for disaster recovery
  accessTier: HOT
  enableHttpsTrafficOnly: true
  minTlsVersion: TLS1_2
  networkRules:
    defaultAction: DENY
    bypassAzureServices: true
    ipRules:
      - "203.0.113.0/24"  # CI/CD runner IPs
    virtualNetworkSubnetIds:
      - "/subscriptions/12345678-1234-1234-1234-123456789012/resourceGroups/prod-network-rg/providers/Microsoft.Network/virtualNetworks/prod-vnet/subnets/app-subnet"
  blobProperties:
    enableVersioning: true
    softDeleteRetentionDays: 30
    containerSoftDeleteRetentionDays: 30
  containers:
    - name: application-data
      accessType: PRIVATE
    - name: backups
      accessType: PRIVATE
    - name: logs
      accessType: PRIVATE
```

**Use Case**: Production environment with geo-redundancy, network isolation, and data protection.

**Security Features**:
- Geo-redundant storage for disaster recovery
- Network restricted to VNet and specific IPs
- Versioning enabled for data recovery
- 30-day soft delete retention

---

## Static Website Hosting

Configuration for hosting static websites with public blob access.

```yaml
apiVersion: azure.project-planton.org/v1
kind: AzureStorageAccount
metadata:
  name: mysite-static
  org: mycompany
  env: production
spec:
  region: westus2
  resourceGroup: web-rg
  accountKind: STORAGE_V2
  accountTier: STANDARD
  replicationType: ZRS  # Zone redundancy for availability
  accessTier: HOT
  enableHttpsTrafficOnly: true
  minTlsVersion: TLS1_2
  networkRules:
    defaultAction: ALLOW  # Public access needed for website
    bypassAzureServices: true
  blobProperties:
    enableVersioning: false
    softDeleteRetentionDays: 7
    containerSoftDeleteRetentionDays: 7
  containers:
    - name: $web
      accessType: BLOB  # Public read for static website
    - name: assets
      accessType: BLOB
```

**Use Case**: Static website hosting with public read access for web content.

**Note**: Enable static website hosting in Azure Portal or via Azure CLI after deployment.

---

## Data Lake Configuration

Configuration for Data Lake Storage Gen2 with premium performance.

```yaml
apiVersion: azure.project-planton.org/v1
kind: AzureStorageAccount
metadata:
  name: enterprise-datalake
  org: mycompany
  env: production
spec:
  region: eastus2
  resourceGroup: analytics-rg
  accountKind: STORAGE_V2
  accountTier: STANDARD
  replicationType: GZRS  # Maximum durability
  accessTier: HOT
  enableHttpsTrafficOnly: true
  minTlsVersion: TLS1_2
  networkRules:
    defaultAction: DENY
    bypassAzureServices: true
    virtualNetworkSubnetIds:
      - "/subscriptions/12345678-1234-1234-1234-123456789012/resourceGroups/analytics-rg/providers/Microsoft.Network/virtualNetworks/analytics-vnet/subnets/databricks-subnet"
      - "/subscriptions/12345678-1234-1234-1234-123456789012/resourceGroups/analytics-rg/providers/Microsoft.Network/virtualNetworks/analytics-vnet/subnets/synapse-subnet"
  blobProperties:
    enableVersioning: true
    softDeleteRetentionDays: 90
    containerSoftDeleteRetentionDays: 90
  containers:
    - name: raw
      accessType: PRIVATE
    - name: processed
      accessType: PRIVATE
    - name: curated
      accessType: PRIVATE
```

**Use Case**: Data lake for analytics workloads with hierarchical namespace support.

**Note**: Enable hierarchical namespace in Azure Portal for full Data Lake Gen2 functionality.

---

## Network-Isolated Configuration

Configuration for high-security environments with VNet-only access.

```yaml
apiVersion: azure.project-planton.org/v1
kind: AzureStorageAccount
metadata:
  name: secure-storage
  org: enterprise
  env: production
spec:
  region: eastus
  resourceGroup: secure-rg
  accountKind: STORAGE_V2
  accountTier: STANDARD
  replicationType: GZRS
  accessTier: COOL  # Cost optimization for infrequent access
  enableHttpsTrafficOnly: true
  minTlsVersion: TLS1_2
  networkRules:
    defaultAction: DENY
    bypassAzureServices: true
    ipRules: []  # No public IP access
    virtualNetworkSubnetIds:
      - "/subscriptions/12345678-1234-1234-1234-123456789012/resourceGroups/secure-rg/providers/Microsoft.Network/virtualNetworks/secure-vnet/subnets/aks-subnet"
  blobProperties:
    enableVersioning: true
    softDeleteRetentionDays: 365  # Maximum retention
    containerSoftDeleteRetentionDays: 365
  containers:
    - name: secure-data
      accessType: PRIVATE
    - name: audit-logs
      accessType: PRIVATE
```

**Use Case**: High-security environment with no public internet access.

**Note**: For true private endpoint isolation, deploy a Private Endpoint separately after storage account creation.

---

## CLI Workflows

### Validate Manifest

```bash
project-planton validate --manifest storage.yaml
```

### Deploy with Pulumi

```bash
project-planton pulumi update \
  --manifest storage.yaml \
  --stack myorg/myproject/dev \
  --module-dir apis/org/project_planton/provider/azure/azurestorageaccount/v1/iac/pulumi
```

### Deploy with Terraform/OpenTofu

```bash
project-planton tofu apply \
  --manifest storage.yaml \
  --auto-approve
```

---

## Post-Deployment: Accessing Storage

### Azure CLI - List Containers

```bash
az storage container list \
  --account-name myapp-prod-storage \
  --auth-mode login
```

### Azure CLI - Upload Blob

```bash
az storage blob upload \
  --account-name myapp-prod-storage \
  --container-name data \
  --file ./myfile.txt \
  --name myfile.txt \
  --auth-mode login
```

### Azure SDK (Python)

```python
from azure.identity import DefaultAzureCredential
from azure.storage.blob import BlobServiceClient

credential = DefaultAzureCredential()
client = BlobServiceClient(
    account_url="https://myapp-prod-storage.blob.core.windows.net/",
    credential=credential
)

# List containers
for container in client.list_containers():
    print(container.name)
```

---

## Best Practices Summary

1. **Choose appropriate replication** - LRS for dev, GRS/GZRS for production
2. **Use Cool tier** for infrequently accessed data to reduce costs
3. **Enable soft delete** in production for data recovery
4. **Restrict network access** - Never leave storage open to public internet
5. **Use managed identities** for application authentication
6. **Enable versioning** for critical data
7. **Separate storage accounts per environment** (dev, staging, prod)
8. **Use private endpoints** for high-security scenarios
9. **Monitor costs** - Access tier and replication significantly impact pricing
10. **Enable diagnostic logging** for troubleshooting and compliance

---

## Summary

These examples demonstrate the full range of configurations supported by the `AzureStorageAccount` API resource:

- **Minimal**: Quick setup for testing
- **Development**: Relaxed security for faster iteration
- **Production**: Balanced security and durability
- **Static Website**: Public blob access for web hosting
- **Data Lake**: Analytics workloads with hierarchical namespace
- **Network-Isolated**: Maximum security for compliance

All configurations follow Azure security best practices while adapting to different operational requirements.
