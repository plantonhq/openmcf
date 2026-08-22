# DigitalOcean Monitor Alert -- Operational Guide

What experience with this component teaches that the field reference cannot.

## Prefer tags over id lists for droplet fleets

An id-targeted policy watches exactly the droplets listed -- replacements and autoscaled additions are NOT covered until the manifest changes. A tag-targeted policy tracks membership automatically: every droplet carrying the tag is watched the moment it exists. Use id references for singular pets, tags for fleets.

## The metric name IS the contract -- copy it exactly

Metric names are DigitalOcean's raw API paths with their inconsistencies intact: droplet CPU is `v1/insights/droplet/cpu` (no `_utilization_percent` suffix, unlike memory and disk), and database metrics live under `v1/dbaas/alerts/` with `_alerts` suffixes. Validation carries the exact 28-value list, so a typo fails at validation -- but read the error's list rather than guessing the spelling.

## Thresholds are float32 upstream

DigitalOcean stores `value` as a 32-bit float. More than 7 significant digits silently truncate -- 99.999999 becomes 100 by the time it evaluates. Round thresholds are also easier for humans at 3 AM.

## One policy per symptom, not per droplet

Policies accept many targets; a CPU policy covering the whole web fleet beats ten identical per-droplet policies. Split policies when the THRESHOLD differs (databases at 80%, batch workers at 95%), not per target.

## Slack webhooks are credentials

The webhook URL lets anyone post to your channel. The spec marks it sensitive and both provisioners keep it out of plain-text state rendering -- treat it with the same care in your manifest storage (prefer secret references over literals in committed manifests).

## Disabling beats deleting

`enabled: false` keeps the policy defined but silent -- ideal for maintenance windows or pre-staging alerts before a service carries traffic. Deleting the policy loses nothing but its UUID, though: recreating it is cheap and the manifest is the source of truth.

## What is deliberately NOT here

Uptime probing of external endpoints (that is the DigitalOceanUptimeCheck kind); Kubernetes and App Platform alerting (DigitalOcean exposes no monitor-alert metrics for them at the pinned provider); and PagerDuty/webhook delivery beyond email and Slack (the API supports only these two).
