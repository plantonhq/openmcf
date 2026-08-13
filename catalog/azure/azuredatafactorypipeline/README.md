# Overview

The **AzureDataFactoryPipeline** component deploys one pipeline inside an Azure Data Factory (AzureDataFactory) -- an ordered set of activities (copy data, run a data flow, call a stored procedure, wait, branch) that executes as a whole when triggered. The activities travel as raw JSON: exactly the "activities" array the Data Factory Studio's Code view shows.

## Purpose

- **One unit of work**: a pipeline is the thing you run, schedule, and monitor -- parameters make one definition serve many windows and environments.
- **Authoring stays honest**: the activities JSON is the same artifact the Studio produces, so pipelines author in the Studio and ship through the catalog without translation.
- **Teams own their pipelines**: many pipelines per factory, each with its own lifecycle -- the factory's workspace posture stays platform-owned.

## Key Features

- Full azurerm v5 surface: activities JSON, run parameters, variables, annotations, concurrency (1-50), Studio folder, and the elapsed-time metric duration.
- Chart-ready: `data_factory_id` defaults its reference to AzureDataFactory's ID output; the `pipeline_id` and `pipeline_name` outputs are what triggers and run-API calls reference.
- JSON key ordering is normalized by both engines -- reordering keys never shows as configuration drift.

## Use Cases

- **Scheduled ingestion**: a copy-activity pipeline parameterized by time window, run by a schedule trigger.
- **Orchestration**: a pipeline of execute-pipeline activities sequencing other pipelines with dependency conditions.
- **Transformation**: run a mapping data flow with compute sized per run.

## Future Enhancements

- Triggers, datasets, linked services, and data flows arrive as their own kinds -- each referencing this pipeline's factory, with triggers naming this pipeline.

---

© Planton. Licensed under [Apache-2.0](https://github.com/plantonhq/planton/blob/main/LICENSE).
