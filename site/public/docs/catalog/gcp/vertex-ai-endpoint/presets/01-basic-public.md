---
title: "Basic Public Endpoint"
description: "The minimal serving surface: a public Vertex AI endpoint in the ambient project with prediction logging sampled into BigQuery."
type: "preset"
rank: "01"
presetSlug: "01-basic-public"
componentSlug: "vertex-ai-endpoint"
componentTitle: "Vertex AI Endpoint"
provider: "gcp"
icon: "package"
order: 1
---

# Basic Public Endpoint

The minimal serving surface: a public Vertex AI endpoint in the ambient
project with prediction logging sampled into BigQuery.

## What this preset creates

An endpoint named `Recommendations Serving` in `us-central1`, reachable
through the shared regional DNS (`us-central1-aiplatform.googleapis.com`)
once a model is deployed to it. Ten percent of prediction requests and
responses are logged to the `ml_logging` BigQuery dataset. The numeric
endpoint ID is derived deterministically from the resource identity, so
re-creating the same manifest always yields the same endpoint reference.

## When to use

- Standard online prediction serving without private-networking needs
- Development and staging endpoints
- Any endpoint where the shared regional DNS is acceptable

## Remix ideas

- Set `dedicatedEndpointEnabled: true` for an isolated DNS name with
  better performance and reliability (see the private-vpc-peered preset).
- Raise `samplingRate` to `1.0` on low-traffic endpoints to capture every
  prediction for debugging.
- Point `bigqueryDestinationUri` at a fully qualified table
  (`bq://project.dataset.table`) to control the table name.
