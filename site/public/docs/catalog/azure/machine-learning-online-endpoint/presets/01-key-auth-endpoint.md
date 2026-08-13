---
title: "Key-Auth Endpoint"
description: "This preset creates the everyday online endpoint: static-key authentication, a system-assigned identity for pulling images and models, and all traffic routed to a first deployment named `blue`."
type: "preset"
rank: "01"
presetSlug: "01-key-auth-endpoint"
componentSlug: "machine-learning-online-endpoint"
componentTitle: "Machine Learning Online Endpoint"
provider: "azure"
icon: "package"
order: 1
---

# Key-Auth Endpoint

This preset creates the everyday online endpoint: static-key authentication, a system-assigned identity for pulling images and models, and all traffic routed to a first deployment named `blue`.

## When to Use

- The first endpoint a model-serving application needs
- Teams starting with key-based scoring clients (the simplest caller contract)
- The parent object for a blue/green deployment pair

## Key Configuration Choices

- **`authMode: Key`** -- static keys that never expire; read them with `az ml online-endpoint get-credentials`, store them in your secret manager, and rotate deliberately. Move to `AADToken` when callers can carry Entra tokens.
- **`identity: SYSTEM_ASSIGNED`** -- created with the endpoint; grant it AcrPull and Storage Blob Data Reader between endpoint creation and the first deployment (or switch to a pre-granted user-assigned identity in charts).
- **`traffic: {blue: 100}`** -- names the first deployment before it exists; the deployment attaches by reference and takes the traffic when it provisions.

## After Deployment

Attach an **Azure Machine Learning Online Deployment** named `blue` wiring `endpointId` by reference, then POST to the endpoint's `scoring_uri` output with the key in the `Authorization: Bearer` header.
