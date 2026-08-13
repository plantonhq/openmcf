<p align="center">
  <img src="logo.svg" alt="AWS Bedrock AgentCore Runtime" width="80"/>
</p>

# AWS Bedrock AgentCore Runtime

Create and manage [Amazon Bedrock AgentCore agent runtimes](https://docs.aws.amazon.com/bedrock-agentcore/latest/devguide/agents-tools-runtime.html)
— serverless, session-isolated execution environments that host YOUR
agent code (any framework: LangGraph, CrewAI, Strands, plain
Python/Node) behind AWS-managed scaling, identity, and networking.

## What Gets Created

- **An agent runtime** hosting one immutable artifact per version:
  - a **container image** (any language, ECR URI), or
  - an **AWS-managed code bundle** (Python/Node source in S3 plus an
    entrypoint — no image to build).
- **Named endpoints** (folded satellites): each floats on the latest
  runtime version or pins an explicit one — DEFAULT/staging/prod traffic
  splits over one runtime. AWS also maintains an implicit DEFAULT
  endpoint.
- **A resource policy** (optional) on the runtime's own ARN — grant
  other accounts or services permission to invoke it.
- Optional: VPC networking, inbound JWT authorization (OIDC), filesystem
  mounts (EFS, S3 access points, per-session scratch), environment
  variables, and session lifecycle tuning.

Creating a runtime is free — AWS bills per-second for CPU/memory only
while sessions execute.

## The Version Model

Every configuration change creates a new runtime VERSION in place;
switching the artifact arm (code ↔ container) replaces the runtime.
Endpoints decide when live traffic moves: a floating endpoint tracks the
latest version, a pinned one serves its version until re-pointed.

## Naming

AWS's runtime-name charset (a letter, then letters/digits/underscore —
no hyphens) is stricter than platform naming conventions, so the name is
the explicit `spec.runtime_name` field, validated at manifest time.

See [v1alpha1/reference.md](v1alpha1/reference.md) for the full field
reference generated from the spec proto.

---

© Planton. Licensed under [Apache-2.0](https://github.com/plantonhq/planton/blob/main/LICENSE).
