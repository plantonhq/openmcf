# AzureDataFactoryTrigger Pulumi Module

## Overview

Creates one trigger inside an Azure Data Factory -- the instruction that starts pipelines automatically. The spec's variant block (exactly one of `schedule`, `tumbling_window`, `blob_event`, `custom_event`) selects the trigger type; the module creates the matching classic-SDK resource.

## Resources Created

Exactly one of:

- `datafactory.TriggerSchedule` -- wall-clock recurrence
- `datafactory.TriggerTumblingWindow` -- contiguous time windows with backfill, retries, and dependencies
- `datafactory.TriggerBlobEvent` -- blob created/deleted events on a storage account
- `datafactory.TriggerCustomEvent` -- events on an Event Grid custom topic

## Outputs

- `trigger_id` -- the trigger's ARM resource ID (`{factory_id}/triggers/{name}`, the same shape for all four types)
- `trigger_name` -- the trigger's name (what tumbling window dependency entries resolve against)

## Behavior Notes

- **`activated` drives a live Start/Stop lifecycle**: the provider starts the trigger after create when `activated` is true (the platform default, always sent explicitly), STOPS it before every update (Azure forbids editing a started trigger), re-starts it after, and stops it before delete.
- **ENGINE-SHAPE**: the bridged SDK pluralizes the recurrence schedule's list-arg names (`DaysOfMonths`, `DaysOfWeeks`, `Monthlies`) -- name differences only; both engines write the same ARM recurrence schedule.
- **Tumbling window triggers drive exactly ONE pipeline** (Azure's own model); a dependency entry without `trigger_name` is a self-dependency.
- **Blob event triggers wire through Event Grid on the storage account** behind the scenes -- the subscription requires the Microsoft.EventGrid resource provider registered.
- **No tags**: triggers are ARM sub-resources of the factory and expose none.
