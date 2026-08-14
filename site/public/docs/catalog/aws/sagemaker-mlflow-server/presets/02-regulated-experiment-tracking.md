---
title: "Regulated Experiment Tracking"
description: "This preset pins down every variable a governed ML platform cares about: a `Medium` server on a pinned MLflow version, maintenance in a declared quiet window, and every logged model auto-registered..."
type: "preset"
rank: "02"
presetSlug: "02-regulated-experiment-tracking"
componentSlug: "sagemaker-mlflow-server"
componentTitle: "SageMaker MLflow Server"
provider: "aws"
icon: "package"
order: 2
---

# Regulated Experiment Tracking

This preset pins down every variable a governed ML platform cares
about: a `Medium` server on a pinned MLflow version, maintenance in a
declared quiet window, and every logged model auto-registered into the
SageMaker Model Registry for audit.

## When to Use

- Platforms where model lineage must land in the Model Registry
  automatically, not by convention
- Teams that need a pinned MLflow version and predictable maintenance
  timing

## What You Get

- A `Medium` server (~50 users) pinned to MLflow `3.0` — the pin is
  `major.minor` because AWS normalizes away the patch
- Automatic model registration: every model logged to MLflow appears
  in the SageMaker Model Registry
- Maintenance held to Sundays 03:00 UTC

## Customize

- `automaticModelRegistration: true` is effectively one-way — the
  provider silently drops a true-to-false change, so turning it off
  later means replacing the server (~50 minutes of lifecycle) or an
  out-of-band API call; enable it deliberately
- Changing `mlflowVersion` replaces the server — treat version bumps
  as scheduled ~25-minute-each create/delete events
- Shift `weeklyMaintenanceWindowStart` (UTC `DDD:HH:MM`) to your own
  quiet hours
