# AWS Bedrock AgentCore Evaluation

Amazon Bedrock AgentCore Evaluations — scoring definitions (LLM judge
or Lambda), repeatable agent test benches, and continuous evaluation
of production sessions sampled from CloudWatch logs.

## What Gets Created

- Evaluators: LLM-as-a-judge with a categorical or numerical scale,
  or a code-based evaluator backed by
  [AWS Lambda](/cloud-catalog/aws-lambda).
- Harnesses: a model (Bedrock, Gemini, or OpenAI), system prompts,
  and optional tools / memory / runtime environment.
- Online evaluation configs: sample production sessions from
  CloudWatch log groups and score them with evaluators.

## Before You Deploy

### Planton Setup

- An organization and environment.
- An AWS provider connection with Bedrock AgentCore control-plane
  permissions.

### AWS Account

- An IAM role trusting `bedrock-agentcore.amazonaws.com` for harnesses
  and online configs ([AWS IAM Role](/cloud-catalog/aws-iam-role)).
- LLM-judge evaluators and Bedrock harnesses need Bedrock model
  access in the region.
- Code-based evaluators need a Lambda the role may invoke.

## Deploy

### Console

Create the resource from the AWS catalog, add an evaluator (or a
harness), and deploy.

### CLI

```bash
planton apply -f evaluation.yaml
```

## After Deploy

- `evaluator_ids` / `harness_ids` / `online_evaluation_config_ids`
  are name-keyed maps. Downstream online configs in other bundles
  can reference a custom evaluator by ID.
- Creating the objects is free. Billing starts when you run an
  evaluation or when an online config samples production traffic.

---

© Planton. Licensed under [Apache-2.0](https://github.com/plantonhq/planton/blob/main/LICENSE).
