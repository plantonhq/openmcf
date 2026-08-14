# AWS Bedrock Prompt

Prompts as managed infrastructure — reusable, reviewable prompt
definitions with variants for A/B comparison, chat templates with tools,
and model or agent execution targets, versioned in Bedrock Prompt
Management instead of scattered through application code.

## What Gets Created

- A Bedrock prompt with up to 3 variants: text templates with
  `{{variables}}`, or chat templates carrying multi-turn messages,
  system context, tool catalogs (with auto/any/forced tool choice), and
  prompt-caching checkpoints.
- Each variant targets a model (ID/ARN/inference profile) or a Bedrock
  agent alias, with its own inference settings, metadata annotations,
  and model-specific request fields.

Creating a prompt is free; the target model bills per invocation when
the prompt executes.

## Before You Deploy

### Planton Setup

- An organization and environment.
- An AWS provider connection with Bedrock agent permissions
  (`bedrock:CreatePrompt` and its read/update/delete siblings).

### AWS Account

- Bedrock Prompt Management available in the target region.
- Access to the variant models — auto-enabled AWS models need nothing;
  marketplace models need an agreement (`AwsBedrockModelAccess`) first.

## Deploy

### Console

Create the resource from the AWS catalog, pick the region, define at
least one variant, and deploy.

### CLI

```bash
planton apply -f prompt.yaml
```

```yaml
apiVersion: aws.planton.dev/v1alpha1
kind: AwsBedrockPrompt
metadata:
  name: support-answer
spec:
  region: us-west-2
  defaultVariant: chat
  variants:
    - name: text
      modelId: amazon.nova-micro-v1:0
      text:
        text: "Answer concisely: {{question}}"
        inputVariables:
          - question
    - name: chat
      modelId: amazon.nova-lite-v1:0
      chat:
        system:
          - text: You are a concise support assistant.
        messages:
          - role: user
            text: "Answer the question: {{question}}"
        inputVariables:
          - question
```

## Operational Notes

- **Variants are the A/B mechanism.** Keep the candidate formulations as
  named variants on ONE prompt and flip `default_variant` — consumers
  invoking the prompt by ID follow the default without redeploying.
- **This resource manages the DRAFT.** Publish numbered versions from
  the console/API when a formulation graduates; the draft keeps
  evolving underneath them.
- **Cache points need model support.** `cache_point: true` marks a
  prompt-caching checkpoint — verify the target model supports caching
  before relying on the cost savings.
- **Agent-targeted variants reference an ALIAS.** Read the ARN from the
  agent's `alias_arns` output map — the chart orders the deployments.
