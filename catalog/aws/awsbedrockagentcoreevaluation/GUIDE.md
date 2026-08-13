# AwsBedrockAgentCoreEvaluation — Component Guide

Authored operational judgment for the AgentCore Evaluations component:
the design decisions behind the spec's shape, and what to know before
scoring agents in production.

## Design decisions

- **One bundle, three arms.** Evaluators, harnesses, and online
  configs are AWS-standalone resources, but a team provisions them
  together — name-keyed collections, each arm optional (at least one,
  CEL-guarded).
- **Harness model vendors are exactly-one.** Bedrock, Gemini, and
  OpenAI are mutually exclusive; the spec rejects a silent first-arm
  pick the provider would otherwise make.
- **Rating scales are exactly-one.** Categorical and numerical are
  alternatives; mixing them is a manifest error, not an AWS surprise.
- **Runtime environment is optional and computed.** A harness does
  not need an agent runtime; when omitted, AWS assigns the default.
  The explicit arm is there for the day you pin one.
- **Float inference parameters are import-normalized.** The harness
  temperature / top-p / max-tokens class is the same float32
  import-drift family as other AgentCore kinds — the import catalog
  writes normalized attributes so a no-op apply stays a no-op.

## Scoring agents in production

- **Start with a code evaluator.** Create does not invoke the Lambda,
  so the first evaluator is fixture-cheap and independent of Bedrock
  model access.
- **LLM judges need model access.** A create that validates the
  judge's model will fail in accounts that have not enabled that
  foundation model — the same payment-instrument class as other
  Bedrock kinds.
- **Online configs sample production, they do not create it.** Point
  them at log groups AgentCore observability actually writes to;
  sampling an empty group scores nothing and still bills the config's
  idle path only when AWS starts runs.
- **Builtin evaluator IDs** (`Builtin.Helpfulness` and siblings) are
  AWS-owned; custom IDs come from this bundle's `evaluator_ids` map.
  Do not invent a builtin name.

---

© Planton. Licensed under [Apache-2.0](https://github.com/plantonhq/planton/blob/main/LICENSE).
