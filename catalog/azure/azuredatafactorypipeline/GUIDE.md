# Azure Data Factory Pipeline -- Operational Guide

Judgment calls that matter when you run Data Factory pipelines in production.

## Author in the Studio, ship through the catalog

The activities JSON is not a place to hand-write orchestration: build the pipeline in the Data Factory Studio (a dev factory is ideal), open the Code view, and copy the "activities" array into the manifest. The catalog deliberately does not re-model Azure's activity schema -- dozens of activity types, each with its own shape, evolving on Azure's schedule -- so the Studio is the authoring surface and the manifest is the shipping surface.

## Deploy-time green is not run-time green

A pipeline deploys successfully the moment its JSON parses and its activity types are known to ARM -- but linked services and datasets the activities reference are validated at RUN time. A pipeline can sit green for weeks and fail on first trigger because a referenced linked service never existed. After deploying, run a Debug execution before wiring triggers.

## Concurrency is a self-protection dial, not a throughput dial

`concurrency` caps simultaneous runs of THIS pipeline. Leave it unset and a backfill trigger can launch dozens of overlapping runs against the same tables; set it to 1 for pipelines that must never overlap (most incremental loads), and remember queued runs wait -- they do not fail.

## Parameterize windows, not environments

Run-time parameters are for what changes per RUN (the date window, the source partition); what changes per ENVIRONMENT (connection strings, account names) belongs in the factory's global parameters and linked services. A pipeline whose JSON embeds environment facts cannot promote between factories -- which defeats the one-definition-many-environments model the factory/pipeline split exists for.
