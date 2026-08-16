# Content Safety Baseline

This preset creates a guardrail covering all six harmful-content
categories (high strength on the severe ones, medium on insults and
misconduct), prompt-attack detection on inputs, and the AWS-managed
profanity list — with one published version (`prod`) for production
pinning.

## When to Use

- The starting point for any customer-facing assistant or chat surface
- Teams that want AWS's managed harm taxonomy enforced before writing any
  custom policy

## Key Configuration Choices

- **PROMPT_ATTACK output strength is NONE** — AWS's contract: jailbreak
  detection applies to what users send, not what models answer.
- **A `prod` version is published immediately** so consumers can pin a
  number from day one; the draft stays free for iteration.
- **No denied topics or PII handling yet** — add `topicPolicy` and
  `sensitiveInformationPolicy` as your policy needs sharpen (see the
  pii-redaction preset).

## After Deployment

Point agents or application inference configs at `guardrail_id` +
`version_numbers["prod"]` from the component outputs.
