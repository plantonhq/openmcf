---
title: "Team Model Registry"
description: "This preset is a single team's model package group: the named shell a training pipeline registers versioned model packages into for review, approval, and deployment."
type: "preset"
rank: "01"
presetSlug: "01-team-model-registry"
componentSlug: "sagemaker-model-registry"
componentTitle: "SageMaker Model Registry"
provider: "aws"
icon: "package"
order: 1
---

# Team Model Registry

This preset is a single team's model package group: the named shell a
training pipeline registers versioned model packages into for review,
approval, and deployment.

## When to Use

- The first registry group for a team's model family
- Organizing model versions produced by a training pipeline

## What You Get

- A model package group named `team-model-registry` with a description
- `model_package_group_name` / `model_package_group_arn` outputs your
  pipelines register versions against

## Customize

- Write the description carefully before applying — changing it later
  replaces the group (only tags are mutable)
- Create one group per model family — groups are free, and separate
  groups keep approval workflows and lineage legible
- Add a `resourcePolicy` when another account needs to discover or
  register models on the group — it is the one arm that updates in
  place
