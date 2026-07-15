---
title: "Container App Environment Dapr Component"
description: "Container App Environment Dapr Component deployment documentation"
icon: "package"
order: 100
componentName: "azurecontainerappenvironmentdaprcomponent"
---

# Azure Container App Environment Dapr Component

Registers a Dapr component -- a state store, pub/sub broker, secret store, or binding backend -- on a Container App Environment, consumable by the Dapr-enabled Container Apps scoped to it.

## What Gets Created

When you deploy an AzureContainerAppEnvironmentDaprComponent resource, Planton provisions:

- **Dapr component registration** -- an `azurerm_container_app_environment_dapr_component` on the referenced environment, carrying the component type, version, metadata, secrets, and app scopes

## Prerequisites

- **Azure credentials** configured via environment variables or Planton provider config
- **An AzureContainerAppEnvironment** to register the component on (referenced through `containerAppEnvironmentId`)
- **The backing service** the component type requires (a storage account for `state.azure.blobstorage`, a Service Bus namespace for `pubsub.azure.servicebus`, ...)

## Quick Start

Create a file `statestore.yaml`:

```yaml
apiVersion: azure.planton.dev/v1
kind: AzureContainerAppEnvironmentDaprComponent
metadata:
  name: statestore
  annotations:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: dev.AzureContainerAppEnvironmentDaprComponent.statestore
spec:
  containerAppEnvironmentId:
    valueFrom:
      kind: AzureContainerAppEnvironment
      name: my-env
      fieldPath: status.outputs.environment_id
  componentName: statestore
  componentType: state.azure.blobstorage
  version: v1
  secrets:
    - name: account-key
      value: "<storage-account-access-key>"
  metadata:
    - name: accountName
      value: mystorageaccount
    - name: containerName
      value: dapr-state
    - name: accountKey
      secretName: account-key
  scopes:
    - my-app
```

Deploy:

```shell
planton apply -f statestore.yaml
```

Apps consume the component through the Dapr API by its `componentName`; only apps whose `dapr.app_id` appears in `scopes` see it (an empty scopes list exposes it to every Dapr-enabled app in the environment).

## Key Outputs

| Output | Purpose |
|--------|---------|
| `dapr_component_id` | The component's ARM ID |
| `component_name` | What application code passes to the Dapr API |

## Related Resources

- [Azure Container App Environment](/docs/catalog/azure/container-app-environment) -- the environment the component is registered on
- [Azure Container App](/docs/catalog/azure/container-app) -- the Dapr-enabled workloads that consume it
