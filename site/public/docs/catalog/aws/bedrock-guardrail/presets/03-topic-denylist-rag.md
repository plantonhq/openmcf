---
title: "Topic Denylist for RAG"
description: "This preset creates a guardrail for retrieval-augmented assistants: two denied topics (legal and medical advice) on the STANDARD cross-lingual tier, plus contextual grounding thresholds that reject..."
type: "preset"
rank: "03"
presetSlug: "03-topic-denylist-rag"
componentSlug: "bedrock-guardrail"
componentTitle: "Bedrock Guardrail"
provider: "aws"
icon: "package"
order: 3
---

# Topic Denylist for RAG

This preset creates a guardrail for retrieval-augmented assistants: two
denied topics (legal and medical advice) on the STANDARD cross-lingual
tier, plus contextual grounding thresholds that reject answers not
supported by the retrieved sources (GROUNDING 0.75) or off-question
(RELEVANCE 0.5).

## When to Use

- Knowledge-base assistants that must stay inside their document corpus
- Products whose legal posture forbids professional-advice territory

## Key Configuration Choices

- **Topic definitions are classifier briefs.** Write them the way you
  would brief a human reviewer; the examples sharpen the boundary.
- **Grounding thresholds only fire when sources are supplied** at
  invocation (RAG flows). 0.75 is a strong starting point — raise it
  toward 0.99 for strict factuality, lower it if too many valid answers
  are rejected.
- **STANDARD tier** extends topic evaluation across languages — the right
  default for user-facing products in 2026.

## After Deployment

Tune thresholds against real traffic in detect-only trials before
tightening; grounding scores are visible in guardrail traces.
