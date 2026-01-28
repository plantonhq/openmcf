# Azure Key Vault - Pulumi Examples

Examples showing how to use the Azure Key Vault Pulumi module with various configurations.

## Minimal Configuration

```yaml
apiVersion: azure.openmcf.org/v1
kind: AzureKeyVault
metadata:
  name: myapp-dev-kv
spec:
  region: eastus
  resource_group: dev-rg
  secret_names:
    - api-key
```

## Production Configuration

```yaml
apiVersion: azure.openmcf.org/v1
kind: AzureKeyVault
metadata:
  name: myapp-prod-kv
  org: mycompany
  env: production
spec:
  region: eastus
  resource_group: prod-security-rg
  sku: PREMIUM
  enable_rbac_authorization: true
  enable_purge_protection: true
  soft_delete_retention_days: 90
  network_acls:
    default_action: DENY
    bypass_azure_services: true
    ip_rules:
      - "203.0.113.0/24"
  secret_names:
    - database-connection-string
    - api-key
    - jwt-secret
```

## Development Configuration

```yaml
apiVersion: azure.openmcf.org/v1
kind: AzureKeyVault
metadata:
  name: myapp-dev-kv
  env: development
spec:
  region: eastus
  resource_group: dev-rg
  sku: STANDARD
  enable_rbac_authorization: true
  enable_purge_protection: false
  soft_delete_retention_days: 7
  network_acls:
    default_action: DENY
    bypass_azure_services: true
    ip_rules:
      - "203.0.113.0/24"
  secret_names:
    - dev-api-key
```
