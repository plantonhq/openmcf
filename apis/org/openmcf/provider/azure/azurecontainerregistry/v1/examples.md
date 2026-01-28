# Azure Container Registry Examples

This document provides comprehensive examples for deploying Azure Container Registry (ACR) using the AzureContainerRegistry API resource.

## Table of Contents

- [Minimal Example](#minimal-example)
- [Production Standard SKU](#production-standard-sku)
- [Premium SKU with Geo-Replication](#premium-sku-with-geo-replication)
- [Development with Admin User](#development-with-admin-user)
- [Global Multi-Region Deployment](#global-multi-region-deployment)

---

## Minimal Example

The simplest ACR deployment with defaults.

```yaml
apiVersion: azure.openmcf.org/v1
kind: AzureContainerRegistry
metadata:
  name: minimal-acr
spec:
  region: eastus
  registryName: minimalacr123
```

**Deploy:**
```shell
planton apply -f minimal-acr.yaml
```

**What you get:**
- Standard SKU registry
- Admin user disabled (secure default)
- Public network access (can be restricted post-deployment)
- Login server: `minimalacr123.azurecr.io`

**Push an image:**
```shell
az acr login --name minimalacr123
docker tag myapp:latest minimalacr123.azurecr.io/myapp:latest
docker push minimalacr123.azurecr.io/myapp:latest
```

---

## Production Standard SKU

Standard SKU for most production workloads.

```yaml
apiVersion: azure.openmcf.org/v1
kind: AzureContainerRegistry
metadata:
  name: production-acr
  labels:
    environment: production
    team: platform
spec:
  region: eastus
  registryName: mycompanyacr
  sku: STANDARD
  adminUserEnabled: false  # Secure: use Azure AD/managed identities
```

**Deploy:**
```shell
planton apply -f production-acr.yaml
```

**What you get:**
- Standard SKU: 100 GB storage, ~3,000 pulls/min
- Secure: admin user disabled, Azure AD authentication only
- Cost: ~$20/month
- Perfect for single-region production deployments

**Authenticate with Azure AD:**
```shell
az acr login --name mycompanyacr
# Or from AKS using managed identity (automatic)
```

---

## Premium SKU with Geo-Replication

Premium SKU with multi-region replication for global deployments.

```yaml
apiVersion: azure.openmcf.org/v1
kind: AzureContainerRegistry
metadata:
  name: global-acr
  labels:
    environment: production
    scope: global
spec:
  region: eastus  # Primary region
  registryName: globalcompanyacr
  sku: PREMIUM
  adminUserEnabled: false
  geoReplicationRegions:
    - westeurope      # European customers
    - southeastasia   # APAC customers
    - australiaeast   # Australia/NZ customers
```

**Deploy:**
```shell
planton apply -f global-acr.yaml
```

**What you get:**
- Premium SKU: 500 GB storage, ~10,000 pulls/min
- Geo-replicated to 4 regions total (primary + 3 replicas)
- Automatic traffic routing to nearest replica
- Lower latency for global Kubernetes deployments
- Cost: ~$20/month base + ~$50/month per replica = ~$170/month

**Benefits:**
- Reduced cross-region bandwidth costs
- Lower latency for image pulls
- Regional redundancy

---

## Development with Admin User

Basic SKU for development with admin user enabled.

```yaml
apiVersion: azure.openmcf.org/v1
kind: AzureContainerRegistry
metadata:
  name: dev-acr
spec:
  region: eastus
  registryName: devacr123
  sku: BASIC
  adminUserEnabled: true  # Convenience for local development
```

**Deploy:**
```shell
planton apply -f dev-acr.yaml
```

**What you get:**
- Basic SKU: 10 GB storage (cost-effective for dev)
- Admin user enabled for simple docker login
- Cost: ~$5/month

**Login with admin credentials:**
```shell
# Get credentials
az acr credential show --name devacr123

# Docker login
docker login devacr123.azurecr.io \
  --username <admin-username> \
  --password <admin-password>
```

**⚠️ Warning:** Admin user is a security anti-pattern. Use only for development, never in production.

---

## Global Multi-Region Deployment

Enterprise deployment with comprehensive geo-replication.

```yaml
apiVersion: azure.openmcf.org/v1
kind: AzureContainerRegistry
metadata:
  name: enterprise-acr
  labels:
    environment: production
    compliance: high
    team: platform-engineering
spec:
  region: eastus  # Primary (US East Coast)
  registryName: enterpriseacr
  sku: PREMIUM
  adminUserEnabled: false
  geoReplicationRegions:
    - westus          # US West Coast
    - westeurope      # Europe
    - northeurope     # Europe North
    - southeastasia   # Singapore
    - japaneast       # Japan
    - australiaeast   # Australia
```

**Deploy:**
```shell
planton apply -f enterprise-acr.yaml
```

**What you get:**
- Premium SKU with 7 total regions (primary + 6 replicas)
- Global coverage across continents
- Traffic automatically routed to nearest replica
- Enterprise-grade throughput and storage

**Cost analysis:**
- Base: $20/month
- 6 replicas × $50/month = $300/month
- **Total: ~$320/month**

**Best for:**
- Global SaaS platforms
- Multi-region Kubernetes deployments
- Organizations with compliance requirements for data residency

---

## Deployment Tips

### Authenticating to ACR

**From Azure CLI:**
```shell
az acr login --name <registry-name>
```

**From AKS (automatic with managed identity):**
```shell
# Attach ACR to AKS cluster
az aks update \
  --name my-aks-cluster \
  --resource-group rg-my-aks-cluster \
  --attach-acr <registry-name>
```

**From CI/CD (GitHub Actions):**
```yaml
- name: Login to ACR
  uses: azure/docker-login@v1
  with:
    login-server: myacr.azurecr.io
    username: ${{ secrets.AZURE_CLIENT_ID }}
    password: ${{ secrets.AZURE_CLIENT_SECRET }}
```

### Verifying Registry

```shell
# Show registry details
az acr show --name <registry-name> --query "{loginServer:loginServer, sku:sku.name}"

# List images
az acr repository list --name <registry-name>

# List tags
az acr repository show-tags --name <registry-name> --repository myapp
```

### Pushing Images

```shell
# Tag image
docker tag myapp:v1.0 <registry-name>.azurecr.io/myapp:v1.0

# Push image
docker push <registry-name>.azurecr.io/myapp:v1.0

# Pull image
docker pull <registry-name>.azurecr.io/myapp:v1.0
```

---

## Common Configurations

### Basic SKU (Dev/Test)
```yaml
sku: BASIC
adminUserEnabled: true  # For convenience
```
- Cost: ~$5/month
- Storage: 10 GB
- Throughput: ~1,000 pulls/min

### Standard SKU (Production)
```yaml
sku: STANDARD
adminUserEnabled: false  # Secure
```
- Cost: ~$20/month
- Storage: 100 GB
- Throughput: ~3,000 pulls/min

### Premium SKU (Enterprise)
```yaml
sku: PREMIUM
adminUserEnabled: false
geoReplicationRegions:
  - westeurope
  - southeastasia
```
- Cost: ~$20/month + $50/month per replica
- Storage: 500 GB
- Throughput: ~10,000 pulls/min
- Geo-replication enabled
- Private Link support

---

## Best Practices

✅ **Use Standard or Premium for production**
- Basic is only for dev/test
- Standard provides adequate throughput for most workloads

✅ **Disable admin user in production**
- Use Azure AD, managed identities, or service principals
- Admin user bypasses Azure AD auditing

✅ **Enable geo-replication for global deployments**
- Reduces cross-region bandwidth costs
- Lower latency for image pulls
- Regional redundancy

✅ **Implement image lifecycle policies**
- Delete untagged manifests after 7-30 days
- Keep only last N versions of each image

✅ **Enable vulnerability scanning**
- Integrate Microsoft Defender for Cloud
- Scan images on push for CVEs

---

## Next Steps

After deploying your container registry:

1. **Authenticate**: Use `az acr login` or configure managed identity
2. **Push images**: Tag and push your container images
3. **Integrate with AKS**: Attach registry to your Kubernetes clusters
4. **Set up CI/CD**: Configure Azure DevOps or GitHub Actions to build and push images
5. **Enable scanning**: Integrate Defender for Cloud for vulnerability detection
6. **Configure retention**: Set up lifecycle policies to manage storage costs

## Support

For issues or questions:
- Check the [research documentation](./docs/README.md) for architecture details
- Review the [Pulumi module documentation](./iac/pulumi/README.md)
- Review the [Terraform module documentation](./iac/tf/README.md)
- Consult the [Azure ACR Documentation](https://docs.microsoft.com/en-us/azure/container-registry/)

