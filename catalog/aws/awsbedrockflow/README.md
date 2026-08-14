<p align="center">
  <img src="logo.svg" alt="AWS Bedrock Flow" width="80"/>
</p>

# AWS Bedrock Flow

Create and manage [Amazon Bedrock flows](https://docs.aws.amazon.com/bedrock/latest/userguide/flows.html) —
node graphs that orchestrate prompts, agents, knowledge bases, Lambda
functions, and control-flow logic into one invocable generative-AI
pipeline.

## What Gets Created

- **A flow** whose `definition` declares the graph:
  - **Nodes** — one per step: Input/Output, inline or referenced Prompts,
    Agents (by alias), KnowledgeBase queries, Lambda calls, Lex
    classification, Conditions, inline Python code, and S3
    retrieval/storage.
  - **Connections** — directed edges: `data` edges move values between
    node sockets; `conditional` edges activate when a Condition node's
    named condition fires.

Flows are free to create — the nodes' model, agent, and knowledge-base
invocations bill at runtime.

## Graph Validation

The spec validates shapes (names, socket types, exactly one
configuration arm per node class); AWS validates the TOPOLOGY server-side
at create/update — unreachable nodes, socket type mismatches, and missing
connections fail with named validation classes. Every flow needs one
Input node and at least one Output node.

## Prerequisites

- An AWS provider connection in Planton.
- An IAM role trusting `bedrock.amazonaws.com`
  (`execution_role_arn`) with invoke permissions on the models, agents,
  knowledge bases, and Lambdas the nodes reference.

## Quick Start

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

## Spec Reference

See [v1alpha1/reference.md](v1alpha1/reference.md) for the full field
reference generated from the spec proto.

---

© Planton. Licensed under [Apache-2.0](https://github.com/plantonhq/planton/blob/main/LICENSE).
