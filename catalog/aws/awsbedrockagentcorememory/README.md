<p align="center">
  <img src="logo.svg" alt="AWS Bedrock AgentCore Memory" width="80"/>
</p>

# AWS Bedrock AgentCore Memory

Create and manage [Amazon Bedrock AgentCore memories](https://docs.aws.amazon.com/bedrock-agentcore/latest/devguide/memory.html)
— managed stores that give agents short-term memory (raw session
events) and long-term memory (structured records extracted from those
events: facts, summaries, preferences, episodes).

## What Gets Created

- **A memory** with a short-term retention window
  (`event_expiry_days`, 7–365) for raw session events.
- **Strategies** (folded satellites) — long-term extraction pipelines:
  - built-in types: **SEMANTIC** (facts), **SUMMARIZATION**,
    **USER_PREFERENCE**, **EPISODIC** (with reflection);
  - **CUSTOM**: override the extraction/consolidation prompts and
    models of a built-in shape.
  - Every strategy declares `namespace_templates` partitioning its
    records (e.g. `/facts/{actorId}`).
- Optional: customer-managed KMS encryption, indexed metadata keys for
  filtered retrieval, and Kinesis delivery of records as they are
  written.

## Strategy Writes Are Serialized

Strategy attach/detach goes through the parent memory's update API — AWS
serializes them per memory, and operations can take tens of minutes.
Both engines order strategy changes after the memory; plan accordingly.

See [v1alpha1/reference.md](v1alpha1/reference.md) for the full field
reference generated from the spec proto.

---

© Planton. Licensed under [Apache-2.0](https://github.com/plantonhq/planton/blob/main/LICENSE).
