# Azure Container Registry - Pulumi Module Overview

## Purpose

This Pulumi module automates the deployment of Azure Container Registry (ACR) instances, providing a production-ready container image storage solution for Azure-based workloads.

## What It Creates

### Core Resources

1. **Resource Group**
   - Dedicated resource group for the registry
   - Named: `rg-<registry_name>`
   - Located in the specified region

2. **Container Registry**
   - Azure Container Registry instance
   - Configurable SKU (Basic, Standard, Premium)
   - Admin user optional (disabled by default)
   - Public or private network access

3. **Geo-Replications** (Premium SKU only)
   - Regional replicas in specified regions
   - Automatic traffic routing to nearest replica
   - Zone-redundant storage per replica

## Architecture

```
┌─────────────────────────────────────────────────────┐
│  Resource Group (rg-<registry_name>)                │
│                                                      │
│  ┌────────────────────────────────────────────┐    │
│  │  Container Registry (<registry_name>.azurecr.io) │
│  │  - SKU: Basic/Standard/Premium              │    │
│  │  - Admin user: Enabled/Disabled             │    │
│  │  - Network: Public/Private                  │    │
│  └────────────────────────────────────────────┘    │
│                                                      │
│  ┌─ Geo-Replications (Premium only) ─────────┐    │
│  │  Replica in westeurope                      │    │
│  │  Replica in southeastasia                   │    │
│  │  Replica in australiaeast                   │    │
│  └─────────────────────────────────────────────┘    │
└─────────────────────────────────────────────────────┘
```

## SKU Comparison

| Feature | Basic | Standard | Premium |
|---------|-------|----------|---------|
| Storage | 10 GB | 100 GB | 500 GB |
| Throughput | ~1K pulls/min | ~3K pulls/min | ~10K pulls/min |
| Geo-replication | ❌ | ❌ | ✅ |
| Private Link | ❌ | ❌ | ✅ |
| Content Trust | ❌ | ❌ | ✅ |
| Cost (monthly) | ~$5 | ~$20 | ~$20 + $50/replica |

## Integration with AKS

ACR integrates seamlessly with Azure Kubernetes Service:

1. **Managed Identity**: AKS can automatically pull from ACR using its managed identity
2. **Attach ACR**: Use `az aks update --attach-acr` to grant AcrPull permissions
3. **No Secrets**: No need to manage image pull secrets in Kubernetes

## Network Security

By default, the registry is publicly accessible. For enhanced security:

- **Premium SKU**: Enable Private Link to restrict access to private endpoints
- **Firewall Rules**: Configure IP allowlists
- **Service Endpoints**: Restrict access to specific VNets

## Authentication Methods

1. **Azure AD**: Use `az acr login` (preferred for users)
2. **Managed Identity**: Automatic for AKS and other Azure services
3. **Service Principal**: For CI/CD pipelines
4. **Admin User**: Simple username/password (dev/test only)

## Lifecycle Management

Post-deployment considerations:

- **Retention Policies**: Auto-delete untagged manifests after N days
- **Image Purge**: Use ACR Tasks to clean up old versions
- **Vulnerability Scanning**: Enable Defender for Cloud integration
- **Backup**: ACR data is automatically replicated within Azure

## Deployment Workflow

1. Define AzureContainerRegistry manifest
2. Planton orchestrator invokes this Pulumi module
3. Module creates resource group and registry
4. For Premium: Creates geo-replications
5. Outputs registry login server and resource ID
6. Users can immediately start pushing images

## References

- [Main Component README](../../README.md)
- [Research Documentation](../../docs/README.md)
- [Azure ACR Documentation](https://docs.microsoft.com/en-us/azure/container-registry/)
