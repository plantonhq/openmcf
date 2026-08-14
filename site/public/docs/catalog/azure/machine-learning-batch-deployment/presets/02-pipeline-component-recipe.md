---
title: "Pipeline Component Recipe"
description: "This preset creates a pipeline-component batch deployment: each endpoint invocation runs a registered pipeline component as a pipeline job -- multi-step batch work (prepare, score, post-process)..."
type: "preset"
rank: "02"
presetSlug: "02-pipeline-component-recipe"
componentSlug: "machine-learning-batch-deployment"
componentTitle: "Machine Learning Batch Deployment"
provider: "azure"
icon: "package"
order: 2
---

# Pipeline Component Recipe

This preset creates a pipeline-component batch deployment: each endpoint invocation runs a registered pipeline component as a pipeline job -- multi-step batch work (prepare, score, post-process) behind one stable, invocable address.

## When to Use

- Batch work that is a PIPELINE, not a single scoring script
- Turning a registered pipeline into a scheduler-invocable service
- Teams standardizing on components as their unit of reusable ML logic

## Key Configuration Choices

- **`pipelineComponent.componentId`** -- a registered pipeline component version's ARM ID (register with `az ml component create`). Required whenever the block is present: a pipeline recipe with no component runs nothing.
- **`settings.default_compute`** -- where the pipeline's STEPS run. Pipeline recipes do not use the deployment's `computeId`/`resources` (those describe model scoring); omitting `default_compute` fails at job time with a compute-resolution error that masquerades as a permissions problem.
- **`settings` values are strings** -- `"false"`, not `false`; the string-valued bag is the typed engine SDK's own shape, and both deployment engines send it verbatim.

## After Deployment

Submit a job against the endpoint (by this deployment's name, or set it as the default) -- each submission materializes one pipeline job; watch it in Azure ML studio's Jobs view under the pipeline experiment.
