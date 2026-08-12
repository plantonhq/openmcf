# Overview

The **AzureDataFactoryDataFlow** component deploys one data flow inside an Azure Data Factory (AzureDataFactory) -- a visually-designed data transformation (filter, join, aggregate, reshape) that executes on Data Factory's managed Spark runtime when a pipeline runs it. One kind covers both of Azure's forms: the runnable mapping data flow and the reusable flowlet snippet (`flowlet: true`) that other data flows embed.

## Purpose

- **Transformation as configuration**: the data flow script -- the Studio's "Script" view artifact -- travels in the manifest, so transformations author in the Studio and ship through the catalog without translation.
- **Reuse through flowlets**: common cleanup/conformance logic ships once as a flowlet and embeds by reference in every flow that needs it.
- **Teams own their flows**: many data flows per factory, each with its own lifecycle -- the factory's workspace posture stays platform-owned.

## Key Features

- Full azurerm v5 surface across BOTH provider forms: script/script lines, named sources/sinks/transformations with dataset, linked service, schema linked service, and flowlet bindings, rejected-row routing on sinks, annotations, and the Studio folder.
- Chart-ready: `data_factory_id` defaults its reference to AzureDataFactory's ID output; flowlet references inside sources/sinks/transformations default to another AzureDataFactoryDataFlow's name output (the flowlet form); the `data_flow_id`/`data_flow_name` outputs are what pipelines' Execute Data Flow activities and other flows' embeds reference.
- Honest to Azure's model: mapping data flows require sources and sinks; flowlets may omit them (the embedding flow supplies them); both forms share one name namespace inside the factory.

## Use Cases

- **Lakehouse conformance**: a mapping data flow reading raw landing data, applying types/dedup/joins, writing curated tables -- run by a pipeline's Execute Data Flow activity.
- **Shared cleanup logic**: a flowlet encapsulating the organization's standard column-scrubbing rules, embedded at the source of every ingest flow.
- **Quarantine routing**: a sink with rejected-row routing sending schema-violating rows to a quarantine store instead of failing the run.

## Future Enhancements

- Datasets and linked services arrive as their own kinds -- the name references here upgrade in place to typed references when they exist.
