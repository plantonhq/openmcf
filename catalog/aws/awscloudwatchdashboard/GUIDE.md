# AwsCloudwatchDashboard — Operational Guide

Live-earned judgment lands here as proof runs and adopter operations teach it; the notes below are the forge-time seed.

## Design the body in the console, own it in the manifest

The fastest authoring loop: build the dashboard visually in the CloudWatch console, open Actions → View/edit source, and paste the JSON (as YAML) into `spec.dashboard_body`. From then on the manifest is the source of truth — every apply is an idempotent PutDashboard upsert.

## The body diffs semantically

AWS normalizes the document server-side and both engines diff it as JSON, so key order and whitespace never show as drift. If a plan shows a body change you did not make, someone edited the dashboard in the console — the apply restores the manifest's truth.

## Metric widgets need no resources to exist

A metric widget charting a non-existent metric renders an empty graph, never an error — dashboards can ship ahead of the services they observe (chart-friendly: deploy the dashboard with the stack, graphs populate as traffic arrives).

## Cost

The first three dashboards per account are free; beyond that ~$3/month each, prorated hourly. Minutes-lived test dashboards cost effectively nothing.
