---
title: "Presets"
description: "Ready-to-deploy configuration presets for Locust"
type: "preset-list"
componentSlug: "locust"
componentTitle: "Locust"
provider: "kubernetes"
icon: "package"
order: 200
presets:
  - slug: "01-web-load-test"
    rank: "01"
    title: "Web load test"
    excerpt: "The interactive Locust: a master with the web UI (login ON — the credential lives in the `load-test-auth` Secret; upstream's open, anyone-can-start-tests UI never ships) and two fixed workers running..."
  - slug: "02-headless-ci"
    rank: "02"
    title: "Headless CI gate"
    excerpt: "The automated Locust: no web UI at all — the test starts the moment the pods are up, runs 200 simulated users for ten minutes against `target_host`, and the master's logs carry the summary. The run..."
---

# Locust Presets

Ready-to-deploy configuration presets for Locust. Each preset is a complete manifest you can copy, customize, and deploy.
