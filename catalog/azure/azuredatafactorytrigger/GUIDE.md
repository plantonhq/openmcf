# Azure Data Factory Trigger -- Operational Guide

Judgment calls that matter when you run Data Factory triggers in production.

## A started trigger is a running system

`activated: true` (the default) means the trigger starts evaluating the moment it deploys -- a schedule trigger with an omitted `start_time` starts from NOW, and every firing is a real, billed pipeline run. Deploy new triggers with an explicit future `start_time` (or `activated: false`), verify the target pipeline's Debug run, then let it go live. The trigger's state survives your intent: it is started or stopped in Azure, not in your manifest history.

## Updates bounce the trigger -- plan for the gap

Azure forbids editing a started trigger, so every update stops it, applies the change, and starts it again. Schedule firings due DURING that gap are skipped, not queued -- for a once-a-day trigger updated at the wrong moment, that is a missed day. Update triggers outside their firing windows, or reconcile the skipped window manually.

## Schedule for clocks, tumbling window for data

A schedule trigger says "run at 02:00"; a tumbling window trigger says "process 00:00-01:00, then 01:00-02:00, ..." and passes each window's bounds to the pipeline (`@trigger().outputs.windowStartTime`). If the pipeline's work is a function of a time range -- every incremental load -- use tumbling windows: you get backfill by pointing `start_time` at history, per-window retry, `max_concurrency` rate-limiting, and window dependencies. Use schedule triggers for genuinely clock-shaped work (report generation, cache warms).

## Backfills are deliberate, bounded operations

A tumbling window `start_time` in the past runs EVERY window between then and now. Before pointing at history: set `max_concurrency` to what the sink tables tolerate (the 50 default fans out fifty parallel runs), set the pipeline's own `concurrency` guard, and give the run `retry` so transient failures do not strand single windows in a sea of green ones.

## Blob event triggers have a second resource behind them

The blob event type creates an Event Grid subscription on the storage account behind the scenes (hence the Microsoft.EventGrid provider requirement). Two operational consequences: the storage account's event pipeline is now part of your trigger's failure domain, and path filters (`blob_path_begins_with`/`ends_with`) are your only volume guard -- an over-broad filter on a busy account fires the pipeline for every blob. Filter to the narrowest path that works, and use `ignore_empty_blobs` when upstream tools write zero-byte markers.

## Self-dependencies serialize windows

A tumbling window dependency entry with no `trigger_name` makes each window wait for the PREVIOUS window of the same trigger -- the lever for loads that must apply strictly in order. Azure requires a negative `offset` for self-dependencies (e.g. `-24:00:00` for daily windows). Without it, windows run in parallel up to `max_concurrency`, which is faster and fine for independent partitions.
