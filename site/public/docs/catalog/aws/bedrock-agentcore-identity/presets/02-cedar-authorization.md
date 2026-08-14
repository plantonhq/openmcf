---
title: "Cedar Authorization"
description: "This preset stands up a Cedar policy engine with a permit/forbid pair — read tools for everyone, mutations denied off-hours — ready to attach to a gateway in LOG_ONLY mode."
type: "preset"
rank: "02"
presetSlug: "02-cedar-authorization"
componentSlug: "bedrock-agentcore-identity"
componentTitle: "Bedrock AgentCore Identity"
provider: "aws"
icon: "package"
order: 2
---

# Cedar Authorization

This preset stands up a Cedar policy engine with a permit/forbid pair —
read tools for everyone, mutations denied off-hours — ready to attach
to a gateway in LOG_ONLY mode.

## When to Use

- Gateways whose tools mutate real systems
- Compliance regimes that want authorization decisions as reviewable
  policy, not application code

## What You Get

- An engine whose `policy_engine_arn` a gateway's `policyEngine`
  attachment consumes
- Policies validated by Cedar's static analysis at deploy
  (FAIL_ON_ANY_FINDINGS rejects rules that can never match)

## Customize

- Start the gateway attachment in `mode: LOG_ONLY`, read the decision
  logs, then flip to ENFORCE
- Keep policy names stable — they are ForceNew identity and the
  `policy_ids` output keys
