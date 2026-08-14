---
title: "Presets"
description: "Ready-to-deploy configuration presets for Machine Learning Compute Instance"
type: "preset-list"
componentSlug: "machine-learning-compute-instance"
componentTitle: "Machine Learning Compute Instance"
provider: "azure"
icon: "package"
order: 200
presets:
  - slug: "01-personal-dev-instance"
    rank: "01"
    title: "Personal Dev Instance"
    excerpt: "This preset creates a personal workstation assigned to a specific team member -- the admin-provisions-for-the-team shape: one general-purpose VM, locked to its owner, with a system identity for..."
  - slug: "02-private-instance"
    rank: "02"
    title: "Private Instance"
    excerpt: "This preset creates a VNet-placed workstation with no public IP -- the hardened posture for estates where personal compute must stay off the internet."
---

# Machine Learning Compute Instance Presets

Ready-to-deploy configuration presets for Machine Learning Compute Instance. Each preset is a complete manifest you can copy, customize, and deploy.
