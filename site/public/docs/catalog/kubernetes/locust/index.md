---
title: "Locust"
description: "Locust deployment documentation"
icon: "package"
order: 100
componentName: "kuberneteslocust"
---

# Locust

The open-source load testing tool: simulate thousands of concurrent
users against your own applications with test behavior written in
plain Python. A master coordinates and serves live charts; a worker
fleet generates the load; your scripts ship as ConfigMaps — no custom
images to build.

## Highlights

- **Load-test what you declare** — point `target_host` at a literal
  URL or another resource's exported endpoint and swarm the services
  you already run on the platform.
- **Secured by default** — upstream's open web UI (able to fire load
  at any reachable host) never ships: the login is on from the first
  apply with a generated credential, delivered through Locust's own
  extension seam — never as pod arguments.
- **Python is the test language** — write user behavior inline with
  supporting modules and pip extras, or reference the ConfigMaps your
  CI already ships; script changes roll the pods automatically.
- **Interactive or headless** — explore with the web UI's live
  request/failure charts, or run headless with users, spawn rate and
  duration declared for CI performance gates.
- **Elastic workers** — fixed count, HPA on CPU, or KEDA scaling on
  the live user count the master reports.

Test credentials ride Secret references into your script's
environment; nothing credential-bearing lands in rendered
configuration.
