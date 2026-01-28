# Create using CLI

Create a YAML file using the example shown below. After the YAML is created, use the command below to apply.

```shell
planton apply -f <yaml-path>
```

# Basic Example

```yaml
apiVersion: azure.openmcf.org/v1
kind: AzureAksNodePool
metadata:
  name: my-azure-aks-node-pool
spec:
  azureProviderConfigId: my-azure-credential-id
```

# Example with Environment Info

```yaml
apiVersion: azure.openmcf.org/v1
kind: AzureAksNodePool
metadata:
  name: my-azure-aks-node-pool
spec:
  azureProviderConfigId: my-azure-credential-id
  environmentInfo:
    envId: production
```

# Example with Stack Job Settings

```yaml
apiVersion: azure.openmcf.org/v1
kind: AzureAksNodePool
metadata:
  name: my-azure-aks-node-pool
spec:
  azureProviderConfigId: my-azure-credential-id
  stackUpdateSettings:
    jobTimeout: 3600
```

# Notes

Since the `spec` is currently empty and the module is not completely implemented, these examples are provided for illustrative purposes. They demonstrate how you would structure your YAML configuration files to create an Azure AKS Cluster using the `AzureAksNodePool` API resource once the module is fully implemented.
