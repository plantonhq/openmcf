# AzureDataFactoryPipeline Terraform Module

## Overview

Creates one pipeline inside an Azure Data Factory -- an ordered set of activities that executes as a whole when triggered. The activities travel as raw JSON; Azure owns the activity schema.

## Resources Created

- `azurerm_data_factory_pipeline` -- the pipeline (activities JSON, parameters, variables, concurrency, folder, elapsed-time metric)

## Variables

The generated `variables.tf` mirrors the proto contract:

- `metadata` -- Planton resource metadata (name, org, env, labels, tags)
- `spec` -- the AzureDataFactoryPipelineSpec fields; the factory reference arrives as a resolved literal

## Outputs

- `pipeline_id` -- the pipeline's ARM resource ID (`{factory_id}/pipelines/{name}`)
- `pipeline_name` -- the pipeline's name (what triggers and run-API calls reference)

## Behavior Notes

- **`activities_json` is the raw "activities" ARRAY** from the Data Factory Studio's Code view. Invalid JSON or unknown activity shapes fail at deploy time; JSON key ordering is normalized (never shows as drift).
- **Parameters and variables are String-typed on this surface** -- the provider sends and reads only string-typed entries. Declare other-typed parameters inside the activities JSON.
- **No tags**: pipelines are ARM sub-resources of the factory and expose none.
