<p align="center">
  <img src="logo.svg" alt="AWS Bedrock AgentCore Evaluation" width="80"/>
</p>

# AWS Bedrock AgentCore Evaluation

Create and manage [Amazon Bedrock AgentCore Evaluations](https://docs.aws.amazon.com/bedrock-agentcore/latest/devguide/evaluations.html)
— scoring definitions, repeatable agent test benches, and continuous
evaluation of production sessions.

Creating evaluation objects is free; AWS bills per evaluation run
(model tokens for LLM judges, Lambda invocations for code evaluators,
sampled-session scoring for online configs). None of the three arms
requires an agent runtime to exist.

## What Gets Created

- **Evaluators** — an LLM judge with a rating scale, or your own
  Lambda function.
- **Harnesses** — repeatable agent test benches (model, tools,
  prompts) that evaluation runs execute against.
- **Online evaluation configs** — continuous scoring of sampled
  production sessions from CloudWatch logs.

Every arm is optional; author the ones this bundle owns (at least
one).

See [v1alpha1/reference.md](v1alpha1/reference.md) for the full field
reference generated from the spec proto.

---

© Planton. Licensed under [Apache-2.0](https://github.com/plantonhq/planton/blob/main/LICENSE).
