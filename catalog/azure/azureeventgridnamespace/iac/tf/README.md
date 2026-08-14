# AzureEventgridNamespace Terraform Module

## Overview

Creates an Azure Event Grid namespace -- the capacity-scaled hub of the newer Event Grid: CloudEvents namespace topics and an optional MQTT broker behind one set of regional endpoints, sized in throughput units.

## Resources Created

- `azurerm_eventgrid_namespace` -- the namespace (capacity, network posture, optional managed identity, optional MQTT topic-spaces configuration)

## Variables

The generated `variables.tf` mirrors the proto contract:

- `metadata` -- Planton resource metadata (name, org, env, labels, tags)
- `spec` -- the AzureEventgridNamespaceSpec fields; the resource group, identity, and MQTT route-topic references arrive as resolved literals

## Outputs

- `namespace_id` -- the namespace's ARM resource ID (the target an AzureEventgridNamespaceTopic's `namespace_id` references)
- `namespace_name` -- the namespace's name
- `identity_principal_id` -- the system-assigned identity's principal (empty when no system-assigned identity)

## Behavior Notes

- **The MQTT block is create-only**: the topic-spaces configuration cannot be added, removed, or changed after create -- the namespace is replaced instead. Block presence is the enable switch (the provider hardcodes state Enabled).
- **The SKU is always "Standard"** (Azure's only value today) and **the inbound rule action is always "Allow"** -- both deliberately not part of the spec; the module sends them explicitly.
- **Platform defaults always sent**: capacity 1 TU (`coalesce`), public network access enabled, MQTT session dials 1/1 -- the rendered plan states every value.
- **Billing**: per throughput unit per hour -- a namespace is NOT free at rest.
