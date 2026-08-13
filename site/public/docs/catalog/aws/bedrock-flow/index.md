---
title: "Bedrock Flow"
description: "Bedrock Flow deployment documentation"
icon: "package"
order: 100
componentName: "awsbedrockflow"
---

# AWS Bedrock Flow

Visual-builder-class AI pipelines as declarative YAML — wire prompts,
agents, knowledge bases, Lambdas, conditions, and inline code into one
invocable graph, versioned like the rest of your infrastructure.

## What Gets Created

- A Bedrock flow whose definition carries your node graph: Input/Output
  boundaries, inline or referenced prompt nodes, agent delegation (by
  alias), knowledge-base queries with optional generation, Lambda and
  Lex calls, condition-based branching, inline Python, and S3
  retrieval/storage.
- Data connections move typed values between node sockets; conditional
  connections branch on a Condition node's expressions (with a `default`
  else-arm).

Creating a flow is free; the nodes' model, agent, and knowledge-base
invocations bill when the flow runs.

## Before You Deploy

### Planton Setup

- An organization and environment.
- An AWS provider connection with Bedrock agent permissions
  (`bedrock:CreateFlow` and its read/update/delete siblings, plus
  `iam:PassRole` on the execution role).

### AWS Account

- Bedrock flows available in the target region.
- An IAM role trusting `bedrock.amazonaws.com` with invoke permissions
  on everything the nodes reference (models, agent aliases, knowledge
  bases, Lambdas).

## Deploy

### Console

Create the resource from the AWS catalog, pick the region, reference the
execution role, declare the graph, and deploy.

### CLI

```bash
planton apply -f flow.yaml
```

```yaml
apiVersion: aws.planton.dev/v1alpha1
kind: AwsBedrockFlow
metadata:
  name: summarizer
spec:
  region: us-west-2
  executionRoleArn:
    valueFrom:
      kind: AwsIamRole
      name: bedrock-flow-role
      fieldPath: status.outputs.role_arn
  definition:
    nodes:
      - name: FlowInput
        type: Input
        outputs:
          - name: document
            type: String
      - name: Summarize
        type: Prompt
        inputs:
          - name: input
            expression: $.data
            type: String
        outputs:
          - name: modelCompletion
            type: String
        prompt:
          inline:
            modelId: amazon.nova-micro-v1:0
            text:
              text: "Summarize in one sentence: {{input}}"
              inputVariables:
                - input
      - name: FlowOutput
        type: Output
        inputs:
          - name: document
            expression: $.data
            type: String
    connections:
      - name: InToPrompt
        source: FlowInput
        target: Summarize
        data:
          sourceOutput: document
          targetInput: input
      - name: PromptToOut
        source: Summarize
        target: FlowOutput
        data:
          sourceOutput: modelCompletion
          targetInput: document
```

## Operational Notes

- **AWS validates the topology, not just the shapes.** Expect named
  validation classes (UnreachableNode, MismatchedNodeInputType,
  MissingConnectionConfiguration) at create — fix the graph, not the
  module.
- **Condition nodes need a `default` arm** — a conditional edge with
  condition `default` is the else-branch; branches with no ending Output
  node fail validation.
- **Reference prompts and agents as chart blocks.** A prompt node can
  consume an `AwsBedrockPrompt`'s `prompt_arn` output; an agent node
  consumes an `AwsBedrockAgent`'s `alias_arns` entry — the chart orders
  the deployments.
- **`$.data` is the whole-value expression** — the common socket wiring;
  narrower JSONPath expressions select fields from Object-typed sockets.
