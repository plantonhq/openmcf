---
title: "Presets"
description: "Ready-to-deploy configuration presets for GHA Runner Scale Set Controller"
type: "preset-list"
componentSlug: "gha-runner-scale-set-controller"
componentTitle: "GHA Runner Scale Set Controller"
provider: "kubernetes"
icon: "package"
order: 200
presets:
  - slug: "01-standard"
    rank: "01"
    title: "Standard controller preset"
    excerpt: "One cluster-wide controller in its own namespace — the shape almost every cluster wants. It installs the runner CRDs and the manager; runner fleets are declared separately (one..."
  - slug: "02-production"
    rank: "02"
    title: "Production controller preset"
    excerpt: "The controller hardened for a fleet the business depends on: a hot standby behind leader election, `eventual` update strategy (controller upgrades wait for running jobs instead of overprovisioning..."
---

# GHA Runner Scale Set Controller Presets

Ready-to-deploy configuration presets for GHA Runner Scale Set Controller. Each preset is a complete manifest you can copy, customize, and deploy.
