---
title: "Presets"
description: "Ready-to-deploy configuration presets for ConfigMap"
type: "preset-list"
componentSlug: "configmap"
componentTitle: "ConfigMap"
provider: "kubernetes"
icon: "package"
order: 200
presets:
  - slug: "01-app-config"
    rank: "01"
    title: "Application Configuration"
    excerpt: "This preset creates a ConfigMap carrying a typical application's configuration: a couple of scalar settings (consumed as environment variables via `envFrom`/`configMapKeyRef`) and a properties file..."
  - slug: "02-immutable-versioned"
    rank: "02"
    title: "Immutable Versioned Configuration"
    excerpt: "This preset creates an immutable ConfigMap with a versioned name (`app-config-v1`). Configuration changes ship as a NEW ConfigMap (`app-config-v2`) plus a workload update pointing at the new name —..."
  - slug: "03-binary-payload"
    rank: "03"
    title: "Binary Payload"
    excerpt: "This preset creates a ConfigMap that carries a binary entry (`binaryData`, base64-encoded) alongside a regular text entry (`data`). Binary entries are for payloads that are not valid UTF-8 — icons,..."
---

# ConfigMap Presets

Ready-to-deploy configuration presets for ConfigMap. Each preset is a complete manifest you can copy, customize, and deploy.
