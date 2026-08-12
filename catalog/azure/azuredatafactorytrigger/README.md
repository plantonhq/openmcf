# Overview

The **AzureDataFactoryTrigger** component deploys one trigger inside an Azure Data Factory (AzureDataFactory) -- the instruction that starts pipelines automatically. One kind covers all four of Azure's trigger types as variants: wall-clock schedules, contiguous tumbling windows (the backfill-friendly type), storage blob events, and Event Grid custom events.

## Purpose

- **Automation as configuration**: the schedule that runs production lives in the manifest, reviewed and versioned like everything else -- not clicked together in a portal.
- **One kind, four types**: the variant block declares the trigger type; Azure stores all four in one factory-scoped namespace with one started/stopped lifecycle, and so does the catalog.
- **Windows, not just clocks**: tumbling window triggers carry their window's start/end into each run, retry failed windows, rate-limit backfills, and can depend on other windows -- the type for incremental data loads.

## Key Features

- Full azurerm v5 surface across ALL FOUR provider trigger resources: recurrence schedules down to monthly weekday occurrences, tumbling window delay/concurrency/retry/dependencies (including self-dependency), blob path filters, and custom event subject/type filters.
- Chart-ready: `data_factory_id` defaults its reference to AzureDataFactory's ID output; pipeline references default to AzureDataFactoryPipeline's name output; the blob event's storage account defaults to AzureStorageAccount's ID output; the custom event's topic defaults to AzureEventgridTopic's ID output; a tumbling window dependency defaults to another AzureDataFactoryTrigger's name output.
- The `activated` flag drives Azure's Start/Stop lifecycle honestly: deploy started (the default) or stopped, flip it to pause -- both engines stop before updating (Azure forbids editing a started trigger) and re-start after.

## Use Cases

- **Nightly ingestion**: a schedule trigger firing the ingest pipeline at 02:00 in the factory's time zone.
- **Windowed incremental loads with backfill**: a tumbling window trigger passing each window's start/end to the pipeline -- point `start_time` at history and Azure backfills window by window, rate-limited by `max_concurrency`.
- **File-arrival processing**: a blob event trigger firing the parse pipeline the moment a `.csv` lands under the landing path.
- **Event-driven orchestration**: a custom event trigger firing when an upstream system publishes readiness to an Event Grid topic.

## Future Enhancements

- Datasets, linked services, and integration runtimes arrive as their own kinds, completing the Data Factory family.
