---
title: "Basic Dev Registry"
description: "This preset creates a Basic-edition Container Registry instance billed pay-as-you-go, with a single auto-created private namespace. It is the right starting point for development and small-team..."
type: "preset"
rank: "01"
presetSlug: "01-basic-dev"
componentSlug: "container-registry"
componentTitle: "Container Registry"
provider: "alicloud"
icon: "package"
order: 1
---

# Basic Dev Registry

This preset creates a Basic-edition Container Registry instance billed pay-as-you-go, with a single auto-created private namespace. It is the right starting point for development and small-team experimentation.

## When to Use

- Development and testing image storage
- Small teams or individual projects with modest image volume
- Getting started with ACR before committing to a subscription

## Key Configuration Choices

- **Basic edition** (`instanceType: Basic`) — the entry tier; enough for dev workloads without enterprise features
- **Pay-as-you-go** (`paymentType: PayAsYouGo`) — no upfront commitment; cost follows usage
- **One private namespace** (`dev`) with `autoCreate: true` — pushing an image to a new repository under this namespace creates the repository automatically
- **Private by default** (`defaultVisibility: PRIVATE`) — no image is publicly pullable unless deliberately made so

## Placeholders to Replace

- `metadata.name` and `instanceName` — your registry's name
- `region` — the AliCloud region closest to your build infrastructure (e.g., `cn-hangzhou`)
