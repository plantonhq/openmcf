---
title: "Web load test"
description: "The interactive Locust: a master with the web UI (login ON — the credential lives in the `load-test-auth` Secret; upstream's open, anyone-can-start-tests UI never ships) and two fixed workers running..."
type: "preset"
rank: "01"
presetSlug: "01-web-load-test"
componentSlug: "locust"
componentTitle: "Locust"
provider: "kubernetes"
icon: "package"
order: 1
---

# Web load test

The interactive Locust: a master with the web UI (login ON — the
credential lives in the `load-test-auth` Secret; upstream's open,
anyone-can-start-tests UI never ships) and two fixed workers running
an inline locustfile against the service you point `target_host` at.

Port-forward the exported service (or compose an exposure kind), sign
in with `locust` and the generated password, choose user count and
spawn rate, and watch live request/failure charts as the swarm runs.

Grow it by sizing `workers.replicas` (each worker is roughly one CPU
core of load generation), adding supporting modules under
`inline.lib_files`, or handing test credentials to your script with
`env_from_secrets` — they arrive as environment variables your Python
reads.
