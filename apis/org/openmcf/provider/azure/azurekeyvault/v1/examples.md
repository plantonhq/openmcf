# Azure Key Vault Examples

This document provides comprehensive examples for the `AzureKeyVault` API resource, demonstrating various secret management scenarios in Microsoft Azure.

## Table of Contents

1. [Minimal Configuration](#minimal-configuration)
2. [Development Environment](#development-environment)
3. [Production Environment](#production-environment)
4. [Enterprise Compliance Configuration](#enterprise-compliance-configuration)
5. [Network-Isolated Configuration](#network-isolated-configuration)

---

## Minimal Configuration

The simplest possible configuration with only required fields.

```yaml
apiVersion: azure.openmcf.org/v1
kind: AzureKeyVault
metadata:
  name: myapp-kv
spec:
  region: eastus
  resourceGroup: myapp-rg
  secretNames:
    - api-key
```

**What Gets Created**:
- Standard SKU Key Vault
- RBAC authorization enabled
- Purge protection enabled
- 90-day soft delete retention
- Network ACLs default to Deny with Azure Services bypass

---

## Development Environment

Configuration optimized for development with relaxed security for easier access.

```yaml
apiVersion: azure.openmcf.org/v1
kind: AzureKeyVault
metadata:
  name: myapp-dev-kv
  org: mycompany
  env: development
spec:
  region: eastus
  resourceGroup: dev-rg
  sku: STANDARD
  enableRbacAuthorization: true
  enablePurgeProtection: false  # Allow deletion for cleanup
  softDeleteRetentionDays: 7   # Minimum retention
  networkAcls:
    defaultAction: DENY
    bypassAzureServices: true
    ipRules:
      - "203.0.113.0/24"  # Office IP range
      - "198.51.100.42"   # VPN gateway
  secretNames:
    - database-password
    - api-key
    - jwt-secret
```

**Use Case**: Development environment where secrets need to be accessible from office IPs and VPN.

---

## Production Environment

Standard production configuration with enhanced security.

```yaml
apiVersion: azure.openmcf.org/v1
kind: AzureKeyVault
metadata:
  name: myapp-prod-kv
  id: kv-prod-001
  org: mycompany
  env: production
spec:
  region: eastus
  resourceGroup: prod-security-rg
  sku: STANDARD
  enableRbacAuthorization: true
  enablePurgeProtection: true
  softDeleteRetentionDays: 90
  networkAcls:
    defaultAction: DENY
    bypassAzureServices: true
    ipRules:
      - "203.0.113.0/24"  # CI/CD runner IPs
    virtualNetworkSubnetIds:
      - "/subscriptions/12345678-1234-1234-1234-123456789012/resourceGroups/prod-network-rg/providers/Microsoft.Network/virtualNetworks/prod-vnet/subnets/app-subnet"
  secretNames:
    - database-connection-string
    - redis-password
    - api-key
    - jwt-secret
    - encryption-key
    - sendgrid-api-key
```

**Use Case**: Production environment with VNet integration and IP restrictions.

**Security Features**:
- Purge protection prevents accidental permanent deletion
- Maximum soft delete retention
- Network restricted to specific VNet and IP ranges
- RBAC for fine-grained access control

---

## Enterprise Compliance Configuration

Maximum security configuration for regulated industries (finance, healthcare, government).

```yaml
apiVersion: azure.openmcf.org/v1
kind: AzureKeyVault
metadata:
  name: enterprise-compliance-kv
  id: kv-compliance-prod
  org: enterprise
  env: production
spec:
  region: eastus
  resourceGroup: compliance-security-rg
  sku: PREMIUM  # HSM-backed for FIPS 140-2 Level 3
  enableRbacAuthorization: true
  enablePurgeProtection: true
  softDeleteRetentionDays: 90
  networkAcls:
    defaultAction: DENY
    bypassAzureServices: true
    ipRules: []  # No public IP access
    virtualNetworkSubnetIds:
      - "/subscriptions/12345678-1234-1234-1234-123456789012/resourceGroups/prod-network-rg/providers/Microsoft.Network/virtualNetworks/prod-vnet/subnets/secure-app-subnet"
  secretNames:
    - pii-encryption-key
    - database-master-key
    - api-signing-key
    - audit-log-key
```

**Use Case**: Highly regulated environment requiring HSM-backed keys and maximum security.

**Compliance Features**:
- Premium SKU with HSM (FIPS 140-2 Level 3)
- No public internet access
- Only accessible from secure VNet subnet
- Suitable for PCI-DSS, HIPAA, FedRAMP

---

## Network-Isolated Configuration

Configuration for applications running entirely within Azure Virtual Networks.

```yaml
apiVersion: azure.openmcf.org/v1
kind: AzureKeyVault
metadata:
  name: vnet-isolated-kv
  org: mycompany
  env: production
spec:
  region: westus2
  resourceGroup: prod-network-rg
  sku: STANDARD
  enableRbacAuthorization: true
  enablePurgeProtection: true
  softDeleteRetentionDays: 90
  networkAcls:
    defaultAction: DENY
    bypassAzureServices: true
    ipRules: []  # No public access
    virtualNetworkSubnetIds:
      - "/subscriptions/12345678-1234-1234-1234-123456789012/resourceGroups/prod-network-rg/providers/Microsoft.Network/virtualNetworks/prod-vnet/subnets/aks-subnet"
      - "/subscriptions/12345678-1234-1234-1234-123456789012/resourceGroups/prod-network-rg/providers/Microsoft.Network/virtualNetworks/prod-vnet/subnets/app-service-subnet"
  secretNames:
    - aks-database-password
    - appservice-api-key
    - shared-encryption-key
```

**Use Case**: Applications running in AKS or App Service that access vault via VNet integration.

**Note**: For true private endpoint isolation, deploy a Private Endpoint separately after vault creation.

---

## Post-Deployment: Setting Secret Values

**IMPORTANT**: This API creates placeholder secrets with empty values. Set actual values using:

### Azure CLI

```bash
az keyvault secret set \
  --vault-name myapp-prod-kv \
  --name database-password \
  --value "ActualSecurePassword123!"
```

### Azure Portal

1. Navigate to Key Vault → Secrets
2. Click the secret name
3. Click "+ New Version"
4. Enter the secret value
5. Click "Create"

### Azure SDK (Python)

```python
from azure.identity import DefaultAzureCredential
from azure.keyvault.secrets import SecretClient

credential = DefaultAzureCredential()
client = SecretClient(vault_url="https://myapp-prod-kv.vault.azure.net/", credential=credential)

client.set_secret("database-password", "ActualSecurePassword123!")
```

---

## Granting Access with RBAC

After vault creation, grant applications access using Azure RBAC:

### For Managed Identity

```bash
# Get the managed identity's object ID
IDENTITY_ID=$(az webapp identity show \
  --name myapp \
  --resource-group prod-rg \
  --query principalId -o tsv)

# Grant secrets read access
az role assignment create \
  --role "Key Vault Secrets User" \
  --assignee $IDENTITY_ID \
  --scope /subscriptions/{sub-id}/resourceGroups/prod-security-rg/providers/Microsoft.KeyVault/vaults/myapp-prod-kv
```

### For Service Principal

```bash
az role assignment create \
  --role "Key Vault Secrets User" \
  --assignee {service-principal-object-id} \
  --scope /subscriptions/{sub-id}/resourceGroups/prod-security-rg/providers/Microsoft.KeyVault/vaults/myapp-prod-kv
```

### Common RBAC Roles

- **Key Vault Administrator**: Full management (for admins only)
- **Key Vault Secrets Officer**: Manage secrets (for DevOps/SecOps)
- **Key Vault Secrets User**: Read secrets (for applications)
- **Key Vault Crypto User**: Use keys for cryptographic operations

---

## Integrating with Applications

### Azure App Service

Use Key Vault references in application settings:

```bash
az webapp config appsettings set \
  --name myapp \
  --resource-group prod-rg \
  --settings \
    ConnectionString="@Microsoft.KeyVault(SecretUri=https://myapp-prod-kv.vault.azure.net/secrets/database-password)"
```

### Azure Kubernetes Service (AKS)

Use the CSI Secret Store Driver:

```yaml
apiVersion: secrets-store.csi.x-k8s.io/v1
kind: SecretProviderClass
metadata:
  name: azure-keyvault-secrets
spec:
  provider: azure
  parameters:
    usePodIdentity: "true"
    keyvaultName: "myapp-prod-kv"
    objects: |
      array:
        - |
          objectName: database-password
          objectType: secret
    tenantId: "your-tenant-id"
```

---

## Best Practices Summary

1. **Use RBAC** for access control (not legacy access policies)
2. **Enable purge protection** in production
3. **Use Premium SKU** for compliance requirements (HSM-backed)
4. **Restrict network access** - never leave vaults open to public internet
5. **Set actual secret values separately** - never hardcode in IaC
6. **Use managed identities** for application authentication
7. **Enable audit logging** to Azure Monitor
8. **Rotate secrets regularly** - Key Vault versions secrets automatically
9. **Separate vaults per environment** (dev, staging, prod)
10. **Tag resources** for cost allocation and governance

---

## Troubleshooting

### Cannot Access Vault

**Issue**: "Forbidden" error when accessing vault

**Solutions**:
1. Verify RBAC role assignment exists
2. Check network ACLs allow your source IP/VNet
3. Ensure managed identity is properly configured

### Vault Name Conflicts

**Issue**: Vault name already taken

**Solution**: Use a unique suffix or check for soft-deleted vaults:

```bash
az keyvault list-deleted
az keyvault recover --name myapp-prod-kv  # Recover if accidentally deleted
```

### Secret Not Found

**Issue**: Secret exists but returns "not found"

**Solutions**:
1. Verify secret was created (may take a few seconds)
2. Check RBAC permissions include read access
3. Ensure using correct vault URI

---

## Summary

These examples demonstrate the full range of configurations supported by the `AzureKeyVault` API resource:

- **Minimal**: Quick setup for testing
- **Development**: Relaxed security for faster iteration
- **Production**: Balanced security and accessibility
- **Enterprise**: Maximum security for compliance
- **Network-Isolated**: VNet-only access patterns

All configurations follow Azure security best practices while adapting to different operational requirements.

