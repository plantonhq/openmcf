---
title: "Throttled Stages"
description: "This preset covers one REST API stage with plan-wide throttle ceilings and a tighter cap on `GET /search` — the shape for protecting an expensive method without starving the rest of the API."
type: "preset"
rank: "02"
presetSlug: "02-throttled-stages"
componentSlug: "rest-api-usage-plan"
componentTitle: "REST API Usage Plan"
provider: "aws"
icon: "package"
order: 2
---

# Throttled Stages

This preset covers one REST API stage with plan-wide throttle ceilings
and a tighter cap on `GET /search` — the shape for protecting an
expensive method without starving the rest of the API.

## When to Use

- APIs whose search (or similar) path is the costly one
- Partner traffic that should be rate-limited independently of a quota

## What You Get

- Plan-wide 50 rps / burst 100
- `/search/GET` capped at 5 rps / burst 10
- One enabled API key (`partner-app`)

## Customize

- Change `methodThrottles.path` to the method you need to protect
  (`{path}/{METHOD}`)
- Add a `quota` block if you also want a daily/weekly/monthly cap
- Add more `apiStages` when one plan covers several APIs
