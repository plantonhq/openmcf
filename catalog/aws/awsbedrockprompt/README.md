<p align="center">
  <img src="logo.svg" alt="AWS Bedrock Prompt" width="80"/>
</p>

# AWS Bedrock Prompt

Create and manage [Amazon Bedrock prompts](https://docs.aws.amazon.com/bedrock/latest/userguide/prompt-management.html)
(Prompt Management) — reusable, versionable prompt definitions with one
or more variants, each targeting a model or an agent with a text or chat
template.

## What Gets Created

- **A prompt** carrying your variants (AWS caps a prompt at 3):
  - **Text variants** — a single template string with `{{variable}}`
    placeholders.
  - **Chat variants** — multi-turn messages, system context, tool
    catalogs with tool-choice control, and prompt-caching checkpoints.
  - Each variant targets a **model** (ID, ARN, or inference profile) or
    an **agent** (by alias ARN), with inference settings and optional
    model-specific request fields.
- `default_variant` names which variant executes when the prompt is
  invoked without naming one.

Prompts are free to create — model invocations bill when the prompt
executes.

## The Draft Model

This resource manages the prompt's mutable working draft (version
`DRAFT`). AWS assigns a new internal version string on every update;
numbered production versions are published from the console or API when
you snapshot a draft.

## Prerequisites

- An AWS provider connection in Planton.
- Access to the target models — auto-enabled AWS models work
  immediately; marketplace models need an `AwsBedrockModelAccess`
  agreement first.

## Quick Start

```yaml
apiVersion: aws.planton.dev/v1alpha1
kind: AwsBedrockPrompt
metadata:
  name: summarize
spec:
  region: us-west-2
  variants:
    - name: main
      modelId: amazon.nova-micro-v1:0
      text:
        text: "Summarize the following text in one sentence: {{input}}"
        inputVariables:
          - input
```

## Spec Reference

See [v1alpha1/reference.md](v1alpha1/reference.md) for the full field
reference generated from the spec proto.

---

© Planton. Licensed under [Apache-2.0](https://github.com/plantonhq/planton/blob/main/LICENSE).
