---
title: "Starter Notebook"
description: "This preset gives a data scientist a ready Jupyter workstation on the cheapest current-generation instance (~$0.05/hour), with the everyday Python stack installed once at creation."
type: "preset"
rank: "01"
presetSlug: "01-starter-notebook"
componentSlug: "sagemaker-notebook-instance"
componentTitle: "SageMaker Notebook Instance"
provider: "aws"
icon: "package"
order: 1
---

# Starter Notebook

This preset gives a data scientist a ready Jupyter workstation on the
cheapest current-generation instance (~$0.05/hour), with the everyday
Python stack installed once at creation.

## When to Use

- The first notebook for exploration and prototyping
- Lightweight data work that doesn't need a GPU

## What You Get

- An `ml.t3.medium` instance with a 20 GB ML storage volume
- A one-time `onCreate` bootstrap installing pandas, scikit-learn, and
  matplotlib — run as root, well inside AWS's 5-minute script limit

## Customize

- Move installs that must survive platform patches into `onStart` (it
  runs on every start, including the first) — keep it fast or push
  long work to the background
- Grow `volumeSizeGb` any time (it updates in place); shrinking
  replaces the instance, so size generously up front
- Most other changes stop, update, and restart the instance — budget
  several minutes and batch them
