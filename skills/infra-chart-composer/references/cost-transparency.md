# Cost Transparency — a Standing Duty

Most people composing charts here are solo practitioners and small teams
paying the cloud bill themselves. Cost is on their mind even when it is not
in their words. **Surface the cost picture of every architecture you propose
— without being asked.** Nobody else combines "composes the infrastructure"
with "tells me what it will cost and how to pay less"; this duty is a large
part of why the composer is loved.

## When and how

- **At the explain-after (Phase 4a)** — or at the plan, when the user asked
  to review one first: alongside the built (or proposed) resource list, give
  the monthly picture — the total order of magnitude and the two or three
  resources that dominate it. One short block, not a rate card:

  > Rough monthly cost: ~$150. The EKS control plane (~$73) and the NAT
  > gateway (~$32 + data) dominate; the rest is nodes and storage.

- **At every costly choice**: when a decision moves the bill meaningfully
  (NAT per-AZ vs single, instance sizes, multi-AZ databases, provisioned vs
  on-demand), state the delta as part of recommending — "per-AZ NAT adds
  ~$65/month; for a dev environment I'd use one."
- **At the finish**: the chart summary repeats the cost picture and names
  the params the user can turn when they want it cheaper.

## Honesty rules

- Numbers are **estimates** — say "roughly" / "~" and note they vary by
  region and usage. Never present an estimate as a quote.
- Know the always-on charges cold (they dominate small setups): EKS control
  plane ~$73/mo; NAT gateway ~$32/mo + data processing; ALB/NLB ~$16-22/mo +
  usage; RDS/OpenSearch/ElastiCache instances bill hourly regardless of
  traffic; EIPs bill when idle. Per-request/per-GB services (S3, Lambda,
  DynamoDB on-demand, Route 53) round to ~$0 at solo-practitioner scale —
  say so, it is good news.
- When precision matters (the user asks, or a choice hinges on it), verify
  against current pricing rather than reciting from memory — the AWS pricing
  pages or `aws pricing` API (read-only) — and say what you checked.

## Saving recommendations

Always pair the number with the lever. The classics that fit chart params:

- Single NAT gateway for non-production (the fleet's charts expose this).
- Right-size nodes and start with fewer; scaling up later is a param change.
- Public EKS endpoint for dev (private needs a standing runner — itself a
  cost); flip to private for production.
- Spot/ARM (Graviton) instance types where the workload tolerates them.
- Turn off what the motive does not need: multi-AZ databases, provisioned
  IOPS, per-AZ redundancy in a sandbox.

Frame savings against the user's motive (see `discovery.md`): production
resilience is worth paying for; a learning sandbox is not.
