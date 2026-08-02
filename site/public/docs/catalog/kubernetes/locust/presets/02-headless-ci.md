---
title: "Headless CI gate"
description: "The automated Locust: no web UI at all — the test starts the moment the pods are up, runs 200 simulated users for ten minutes against `target_host`, and the master's logs carry the summary..."
type: "preset"
rank: "02"
presetSlug: "02-headless-ci"
componentSlug: "locust"
componentTitle: "Locust"
provider: "kubernetes"
icon: "package"
order: 2
---

# Headless CI gate

The automated Locust: no web UI at all — the test starts the moment
the pods are up, runs 200 simulated users for ten minutes against
`target_host`, and the master's logs carry the summary. The run shape
(users, spawn rate, duration) rides Locust's own environment
variables, so a pipeline changes intensity without touching the
script.

Four workers generate the load; scale `workers.replicas` with the
target's size. With no web UI there is nothing to protect, so the
login machinery never deploys on this shape.

Pair it with a deploy step in a pipeline: apply, wait for completion,
read the stats from the master's logs (or keep a non-headless twin
around for interactive exploration of the same script).
