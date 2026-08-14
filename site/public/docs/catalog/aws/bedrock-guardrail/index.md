---
title: "Bedrock Guardrail"
description: "Bedrock Guardrail deployment documentation"
icon: "package"
order: 100
componentName: "awsbedrockguardrail"
---

# AWS Bedrock Guardrail

Content-safety policies for generative AI — filter harmful content, deny
topics, mask PII, and enforce grounding on every model input and output,
with immutable published versions for production pinning.

## What Gets Created

- A Bedrock guardrail carrying the policy families you declare: content
  filters (six harm categories at four strengths), denied topics, word
  filters (managed profanity list + custom words), sensitive-information
  handling (31 PII entity types + custom regexes, each blocked or
  anonymized), and contextual grounding thresholds.
- Optional immutable published versions — one per `versions` entry, each
  AWS-assigned number exported in the `version_numbers` output map keyed
  by your entry name.

Creating a guardrail is free; AWS bills per text unit evaluated when the
guardrail is invoked.

## Before You Deploy

### Planton Setup

- An organization and environment.
- An AWS provider connection with Bedrock permissions
  (`bedrock:CreateGuardrail`, `bedrock:CreateGuardrailVersion` and their
  read/update/delete siblings).

### AWS Account

- Bedrock available in the target region (guardrails are supported in all
  Bedrock commercial regions).
- Optional: a customer-managed KMS key when the guardrail definition must
  be encrypted under your own key (grant Bedrock the standard key-use
  permissions).

## Deploy

### Console

Create the resource from the AWS catalog, pick the region, define at least
one policy family, and deploy.

### CLI

```bash
planton apply -f guardrail.yaml
```

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
  sensitiveInformationPolicy:
    piiEntities:
      - type: EMAIL
        action: ANONYMIZE
  versions:
    - name: prod
```

## Operational Notes

- **Pin versions in production.** The spec edits the mutable DRAFT; agents
  and applications should reference a number from `version_numbers`, so a
  draft edit never changes live behavior unannounced.
- **Trial policies in detect-only mode.** Set a filter's action to `NONE`
  to observe matches in traces without blocking traffic, then flip to
  `BLOCK` when the false-positive rate is acceptable.
- **PROMPT_ATTACK is input-side.** AWS requires its output strength to be
  `NONE` — the filter detects jailbreak attempts in what users send, not
  in what models answer.
- **Version entries are keyed by name, numbered by AWS.** Removing an
  entry deletes that published version unless `keep_on_delete` is set;
  versions in use by agents fail deletion server-side.
