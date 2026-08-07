---
title: "Weighted Canary Split"
description: "Send a small share of production traffic to a canary backend service while the stable backend keeps the majority — the standard GCP mechanism for blue/green and canary rollouts at the URL map layer."
type: "preset"
rank: "02"
presetSlug: "02-weighted-canary"
componentSlug: "url-map-on-google-cloud"
componentTitle: "URL Map on Google Cloud"
provider: "gcp"
icon: "package"
order: 2
---

# Weighted Canary Split

Send a small share of production traffic to a canary backend service while the stable backend keeps the majority — the standard GCP mechanism for blue/green and canary rollouts at the URL map layer.

## When to Use

- Rolling out a new backend service revision without a DNS cutover
- Validating a canary Cloud Run service behind the same global load balancer VIP
- Gradually shifting weight from stable (900) to canary (100) and beyond

## Remix Notes

- Weights are relative: 900/100 means ~90% stable, ~10% canary; GCP normalizes by the sum of weights in the split.
- Set weight to `0` to drain a backend from the split without removing it from the configuration.
- Pair with `tests[]` asserting both backends are reachable from expected paths before raising canary weight.
- `urlRewrite` inside `routeAction` can rewrite paths before forwarding — useful when canary backends expect a different path prefix.
