---
title: "Infrastructure Overview"
description: "Fleet-level host health for Compute Engine workloads: CPU, memory, disk ops, and network egress on one page — the page you open when \"everything feels slow\" and no single service is obviously at..."
type: "preset"
rank: "02"
presetSlug: "02-infrastructure-overview"
componentSlug: "monitoring-dashboard"
componentTitle: "Monitoring Dashboard"
provider: "gcp"
icon: "package"
order: 2
---

# Infrastructure Overview

Fleet-level host health for Compute Engine workloads: CPU, memory, disk
ops, and network egress on one page — the page you open when "everything
feels slow" and no single service is obviously at fault.

## What it configures

- Four grid widgets over `gce_instance` resources, aligned at 60s.
- `deletionPolicy: PREVENT` — this is the posture for a team's primary
  operational view; a destroy must be deliberate.

## Adjust before deploying

- **Memory widget needs the Ops Agent** — `agent.googleapis.com/*`
  metrics only exist on instances running it. Without the agent, the
  chart is empty (not an error).
- **Scope the filters** to your fleet with `resource.labels.zone` or
  instance labels; a project-wide page across teams' VMs answers no
  one's question.

## When to choose something else

For request-level service health (traffic, errors, latency), start from
the **Golden Signals** preset — host health and service health are
different pages on purpose.
