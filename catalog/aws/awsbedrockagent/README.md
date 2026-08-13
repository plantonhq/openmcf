<p align="center">
  <img src="logo.svg" alt="AWS Bedrock Agent" width="80"/>
</p>

# AWS Bedrock Agent

Create and manage [Amazon Bedrock agents](https://docs.aws.amazon.com/bedrock/latest/userguide/agents.html) —
foundation-model-powered assistants that reason over requests, call tools,
retrieve knowledge, and (in multi-agent mode) delegate to collaborator
agents.

## What Gets Created

- **An agent** powered by the foundation model you name, operating through
  the IAM role you reference, with optional guardrail, conversation
  memory, and prompt-template overrides.
- **Action groups** — the agent's tools, one per `action_groups` entry:
  Lambda-backed or return-control operations described by an OpenAPI or
  function schema, or reserved AWS capabilities (user-input elicitation,
  code interpreter).
- **Aliases** — immutable serving endpoints, one per `aliases` entry. Each
  alias snapshots the working draft into a numbered version at creation.
- **Collaborators** — other Bedrock agents this supervisor delegates to,
  referenced through their alias ARNs.
- **Knowledge-base associations** — the knowledge bases the agent queries
  for retrieval-augmented answers.

Agents are free to create — AWS bills per model invocation at runtime.

## The Draft-and-Alias Model

Everything attaches to the agent's mutable working version (`DRAFT`):
action groups, collaborators, and knowledge-base associations all edit the
draft, and after every change the platform lets AWS "prepare" (compile)
the agent for serving. Creating an alias snapshots the assembled draft
into an immutable numbered version — live traffic goes through aliases,
so a draft edit never changes serving behavior until a new alias (or
re-pointed alias) picks it up.

## Prerequisites

- An AWS provider connection in Planton.
- An IAM role trusting `bedrock.amazonaws.com`
  (`agent_resource_role_arn`) with model-invocation permissions.
- The foundation model accessible in the target region — auto-enabled AWS
  models (Amazon Nova, Mistral, Meta) work immediately; marketplace
  models need an `AwsBedrockModelAccess` agreement first.

## Quick Start

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
  aliases:
    - name: live
```

## Multi-Agent Collaboration

Set `agent_collaboration: SUPERVISOR` (planning) or `SUPERVISOR_ROUTER`
(direct routing) and add `collaborators` — each names another agent's
alias ARN (read it from that agent's `alias_arns` output map) and tells
the supervisor when to delegate.

## Spec Reference

See [v1alpha1/reference.md](v1alpha1/reference.md) for the full field
reference generated from the spec proto.

---

© Planton. Licensed under [Apache-2.0](https://github.com/plantonhq/planton/blob/main/LICENSE).
