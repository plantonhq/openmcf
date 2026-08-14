# AzureDataFactoryTrigger Terraform Module

## Overview

Creates one trigger inside an Azure Data Factory -- the instruction that starts pipelines automatically. The spec's variant block (exactly one of `schedule`, `tumbling_window`, `blob_event`, `custom_event`) selects the trigger type.

## Resources Created

Exactly one of:

- `azurerm_data_factory_trigger_schedule` -- wall-clock recurrence (minutes to months, optionally narrowed to specific days/times)
- `azurerm_data_factory_trigger_tumbling_window` -- contiguous, non-overlapping time windows with backfill, retries, and dependencies
- `azurerm_data_factory_trigger_blob_event` -- blob created/deleted events on a storage account, filtered by path
- `azurerm_data_factory_trigger_custom_event` -- events on an Event Grid custom topic, filtered by subject and event type

## Variables

The generated `variables.tf` mirrors the proto contract:

- `metadata` -- Planton resource metadata (name, org, env, labels, tags)
- `spec` -- the AzureDataFactoryTriggerSpec fields; the factory, storage account, Event Grid topic, and pipeline-name references arrive as resolved literals

## Outputs

- `trigger_id` -- the trigger's ARM resource ID (`{factory_id}/triggers/{name}`, the same shape for all four types)
- `trigger_name` -- the trigger's name (what tumbling window dependency entries resolve against)

## Behavior Notes

- **`activated` drives a live Start/Stop lifecycle**: the provider starts the trigger after create when `activated` is true (the platform default, always sent explicitly), STOPS it before every update (Azure forbids editing a started trigger), re-starts it after, and stops it before delete.
- **All four types share one name namespace** inside the factory, so switching variant blocks replaces the trigger (a different provider resource is created at the same ARM address).
- **Schedule pipelines are the modern plural blocks** -- the provider's legacy singular `pipeline_name`/`pipeline_parameters` pair covers the same wire surface and is never sent.
- **Tumbling window triggers drive exactly ONE pipeline** (Azure's own model); a dependency entry without `trigger_name` is a self-dependency.
- **Blob event triggers wire through Event Grid on the storage account** behind the scenes -- the subscription requires the Microsoft.EventGrid resource provider registered.
- **No tags**: triggers are ARM sub-resources of the factory and expose none.
