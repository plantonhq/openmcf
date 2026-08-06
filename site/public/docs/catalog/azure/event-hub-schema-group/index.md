---
title: "Event Hub Schema Group"
description: "Event Hub Schema Group deployment documentation"
icon: "package"
order: 100
componentName: "azureeventhubschemagroup"
---

# Azure Event Hub Schema Group

Creates a schema group in an Event Hubs namespace's schema registry -- a named collection of event schemas with a shared format (Avro or JSON) and a compatibility policy that makes schema evolution safe.

## What Gets Created

When you deploy an AzureEventHubSchemaGroup resource, Planton provisions:

- **Schema Group** -- an `azurerm_eventhub_namespace_schema_group` on the referenced namespace, with your chosen serialization format and evolution policy

## Prerequisites

- **Azure credentials** configured via environment variables or Planton provider config
- **An AzureEventHubNamespace** at STANDARD tier or above (referenced through `namespaceId`) -- BASIC namespaces reject schema groups

## Quick Start

Create a file `schema-group.yaml`:

```yaml
apiVersion: azure.planton.dev/v1alpha1
kind: AzureEventHubSchemaGroup
metadata:
  name: telemetry-schemas
  annotations:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: dev.AzureEventHubSchemaGroup.telemetry-schemas
spec:
  namespaceId:
    valueFrom:
      kind: AzureEventHubNamespace
      name: telemetry-hubs
      fieldPath: status.outputs.namespace_id
  schemaGroupName: telemetry-schemas
  schemaCompatibility: BACKWARD
  schemaType: AVRO
```

Deploy:

```shell
planton apply -f schema-group.yaml
```

Choose the compatibility policy deliberately: BACKWARD (new schemas read old data; upgrade consumers first) is the standard for analytics pipelines; FORWARD (old schemas read new data; upgrade producers first) suits producer-led evolution; NONE skips checking entirely. Every field is fixed at creation -- Azure exposes no mutable properties on a schema group, so any change replaces it and drops the schemas registered inside.

## Key Outputs

| Output | Purpose |
|--------|---------|
| `schema_group_id` | The Azure Resource Manager ID of the group |
| `schema_group_name` | What schema-registry serializers address at runtime, alongside the namespace hostname |

## Related Resources

- [Azure Event Hub Namespace](/docs/catalog/azure/event-hub-namespace) -- the namespace whose registry holds the group (STANDARD+)
- [Azure Event Hub](/docs/catalog/azure/event-hub) -- the streams whose events reference the registered schemas
