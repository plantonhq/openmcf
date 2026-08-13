<p align="center">
  <img src="logo.svg" alt="AWS Bedrock Guardrail" width="80"/>
</p>

# AWS Bedrock Guardrail

Create and manage [Amazon Bedrock guardrails](https://docs.aws.amazon.com/bedrock/latest/userguide/guardrails.html) —
content-safety policy sets evaluated on model inputs and outputs, applied
uniformly across foundation models, agents, and knowledge bases.

## What Gets Created

- **A guardrail** with up to five policy families:
  - **Content filters** — strength-tiered detection of sexual content,
    violence, hate, insults, misconduct, and prompt attacks, on text and
    (per filter) images.
  - **Denied topics** — natural-language topic definitions the model must
    not engage with.
  - **Word filters** — the AWS-managed profanity list and/or exact custom
    words and phrases.
  - **Sensitive information** — PII entity types and custom regex
    patterns, each blocked or anonymized (masked).
  - **Contextual grounding** — thresholds that reject responses not
    grounded in supplied source material or not relevant to the query.
- **Published versions** (optional) — immutable numbered snapshots of the
  guardrail for production pinning, one per `versions` entry, with each
  AWS-assigned number exported in the `version_numbers` output map.

Guardrails are free to create — AWS charges per text unit evaluated at
invocation time.

## The Draft-and-Versions Model

A guardrail always has a mutable working draft (version `DRAFT`). Editing
the spec updates the draft in place. Production consumers should pin a
numbered version published through `versions` — draft edits then never
change live behavior until a new version is published and consumers move
to it.

## Prerequisites

- An AWS provider connection in Planton.
- Optional: a customer-managed KMS key (`kms_key_arn`) when the guardrail
  definition itself must be encrypted under your own key.

## Quick Start

```yaml
apiVersion: aws.planton.dev/v1alpha1
kind: AwsBedrockGuardrail
metadata:
  name: assistant-guardrail
spec:
  region: us-west-2
  blockedInputMessaging: Sorry, I can't help with that request.
  blockedOutputsMessaging: Sorry, I can't provide that response.
  contentPolicy:
    filters:
      - type: HATE
        inputStrength: HIGH
        outputStrength: HIGH
      - type: PROMPT_ATTACK
        inputStrength: HIGH
        outputStrength: NONE
  versions:
    - name: prod
```

## Detect-Only Mode

Every filter carries `input_action`/`output_action` (BLOCK or NONE) and
`input_enabled`/`output_enabled` arms. `NONE` detects and reports (visible
in traces and observability) without intervening — the honest way to trial
a new policy against production traffic before enforcing it.

## Spec Reference

See [v1alpha1/reference.md](v1alpha1/reference.md) for the full field
reference generated from the spec proto.
