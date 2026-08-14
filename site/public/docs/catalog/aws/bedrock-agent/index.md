---
title: "Bedrock Agent"
description: "Bedrock Agent deployment documentation"
icon: "package"
order: 100
componentName: "awsbedrockagent"
---

# AWS Bedrock Agent

A foundation-model-powered assistant with tools, knowledge retrieval,
memory, and multi-agent delegation — the application layer of the Bedrock
AI stack, deployed declaratively with its action groups, aliases,
collaborators, and knowledge-base associations folded into one component.

## What Gets Created

- A Bedrock agent on the model you name (model ID, model ARN, or
  inference-profile ID/ARN), with optional guardrail attachment,
  session-summary memory, and per-step prompt-template overrides.
- Its satellites, keyed by your stable entry names: action groups
  (Lambda-backed or return-control tools, or reserved capabilities like
  `AMAZON.UserInput`), immutable serving aliases (each snapshots the
  draft into a numbered version), collaborator delegations, and
  knowledge-base associations.

Creating an agent is free; AWS bills per model invocation when the agent
serves traffic.

## Before You Deploy

### Planton Setup

- An organization and environment.
- An AWS provider connection with Bedrock agent permissions
  (`bedrock:CreateAgent` and its satellite/read/update/delete siblings,
  plus `iam:PassRole` on the agent's role).

### AWS Account

- Bedrock agents available in the target region.
- An IAM role trusting `bedrock.amazonaws.com` with model-invocation
  permissions (and Lambda invoke permissions when action groups use
  Lambda executors).
- Access to the foundation model — auto-enabled AWS models work
  immediately; marketplace models need an agreement
  (`AwsBedrockModelAccess`) first.

## Deploy

### Console

Create the resource from the AWS catalog, pick the region and model,
reference the agent role, write the instructions, and deploy.

### CLI

```bash
planton apply -f agent.yaml
```

```yaml
apiVersion: aws.planton.dev/v1alpha1
kind: AwsBedrockAgent
metadata:
  name: support-agent
spec:
  region: us-west-2
  foundationModel: amazon.nova-micro-v1:0
  agentResourceRoleArn:
    valueFrom:
      kind: AwsIamRole
      name: bedrock-agent-role
      fieldPath: status.outputs.role_arn
  instruction: >-
    You are a customer support agent. Answer order questions using the
    order tools and keep responses short and accurate.
  actionGroups:
    - name: orders
      executor:
        returnControl: true
      functionSchema:
        functions:
          - name: get_order
            parameters:
              - name: order_id
                type: string
                required: true
  aliases:
    - name: live
```

## Operational Notes

- **Serve through aliases, never the draft.** Every spec edit changes the
  DRAFT; an alias pins the numbered snapshot taken at its creation, so
  production behavior only changes when you create or re-point an alias.
- **Satellites re-prepare the agent.** Each action-group, collaborator,
  or association change triggers an AWS prepare cycle; the platform waits
  for PREPARED before the deploy completes.
- **Collaborators need supervisor mode.** Set `agent_collaboration` to
  SUPERVISOR or SUPERVISOR_ROUTER before adding `collaborators` — and
  reference the collaborating agent's ALIAS (its `alias_arns` output),
  not the agent itself.
- **Reserved action groups take no schema.** `AMAZON.UserInput`,
  `AMAZON.CodeInterpreter`, and the ANTHROPIC computer-use signatures are
  AWS-defined capabilities — name the signature and nothing else.
