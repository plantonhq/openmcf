# AzureContainerAppEnvironmentDaprComponent

Register a Dapr component (state store, pub/sub broker, secret store, or binding) on a Container App Environment, consumable by Dapr-enabled apps in the environment.

## Overview

Dapr components are the pluggable backends behind Dapr's building blocks. A component is registered once on the environment with its type (`state.azure.blobstorage`, `pubsub.azure.servicebus`, ...), configuration metadata, and secrets -- and every Dapr-enabled Container App whose `dapr.app_id` appears in `scopes` can use it through the Dapr API (an empty scopes list exposes it to all Dapr-enabled apps).

## Key Features

- **Any Dapr component type**: state stores, pub/sub brokers, secret stores, input/output bindings -- in Dapr's own dotted notation
- **Secret-safe configuration**: connection strings and keys travel as component secrets referenced from metadata by `secret_name`, never inlined
- **Composable metadata values**: a metadata value can reference another resource's output -- the keyless-auth entries are the canonical case (`azureClientId` tracking a managed identity's `client_id`), removing broker credentials from the deployment entirely
- **App scoping**: expose the component only to the named `dapr.app_id`s -- scope production components deliberately
- **Fail-loud initialisation**: `ignore_errors` defaults to false so a broken backend fails at sidecar startup, not on first use

## When to Use

- Give microservices a shared state store (blob storage, Redis, Cosmos DB) through Dapr's state API
- Wire pub/sub messaging between Container Apps (Service Bus, Kafka) without SDK lock-in
- Schedule work with the cron binding or bridge external systems through I/O bindings

## Spec Highlights

| Field | Notes |
| --- | --- |
| `component_name` | What app code passes to the Dapr API. Max 60 lowercase alphanumerics/hyphens. ForceNew |
| `component_type` | Dapr's dotted type notation (e.g. `state.azure.blobstorage`). ForceNew |
| `version` | "v1" for virtually all stable components |
| `metadata[]` | Component configuration; each entry carries a value (literal or a reference to another resource's output) XOR a `secret_name` reference |
| `secrets[]` | Named secret values referenced from metadata |
| `scopes[]` | The `dapr.app_id`s allowed to use the component; empty = all Dapr-enabled apps |

## Outputs

| Output | Purpose |
| --- | --- |
| `dapr_component_id` | The component's ARM ID |
| `component_name` | What apps pass to the Dapr API |

## Quick Example

```yaml
apiVersion: azure.planton.dev/v1alpha1
kind: AzureContainerAppEnvironmentDaprComponent
metadata:
  name: statestore
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
      value:
        value: mystorageaccount
    - name: containerName
      value:
        value: dapr-state
    - name: accountKey
      secretName: account-key
  scopes:
    - my-app
```

---

© Planton. Licensed under [Apache-2.0](https://github.com/plantonhq/planton/blob/main/LICENSE).
