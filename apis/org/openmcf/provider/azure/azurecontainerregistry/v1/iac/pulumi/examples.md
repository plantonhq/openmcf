# Azure Container Registry - Pulumi Examples

Examples of deploying Azure Container Registry using this Pulumi module.

## Basic Registry

Standard SKU registry for production workloads.

```yaml
apiVersion: azure.openmcf.org/v1
kind: AzureContainerRegistry
metadata:
  name: basic-acr
spec:
  region: eastus
  registryName: mycompanyacr
  sku: STANDARD
  adminUserEnabled: false
```

## Premium with Geo-Replication

Multi-region deployment for global access.

```yaml
apiVersion: azure.openmcf.org/v1
kind: AzureContainerRegistry
metadata:
  name: global-acr
spec:
  region: eastus
  registryName: globalacr
  sku: PREMIUM
  adminUserEnabled: false
  geoReplicationRegions:
    - westeurope
    - southeastasia
```

## Development Registry

Basic SKU with admin user for local development.

```yaml
apiVersion: azure.openmcf.org/v1
kind: AzureContainerRegistry
metadata:
  name: dev-acr
spec:
  region: eastus
  registryName: devacr123
  sku: BASIC
  adminUserEnabled: true
```

## Deploying

All examples above can be deployed using:

```shell
planton apply -f <manifest-file>.yaml
```

## Accessing the Registry

After deployment, authenticate and push images:

```shell
# Authenticate
az acr login --name mycompanyacr

# Tag and push
docker tag myapp:latest mycompanyacr.azurecr.io/myapp:latest
docker push mycompanyacr.azurecr.io/myapp:latest

# Pull
docker pull mycompanyacr.azurecr.io/myapp:latest
```

## For More Examples

See the [main examples documentation](../../examples.md) for comprehensive scenarios including:
- Production configurations
- Multi-region deployments
- Development setups
- Integration with AKS
