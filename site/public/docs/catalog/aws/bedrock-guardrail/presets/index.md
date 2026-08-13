---
title: "Presets"
description: "Ready-to-deploy configuration presets for Bedrock Guardrail"
type: "preset-list"
componentSlug: "bedrock-guardrail"
componentTitle: "Bedrock Guardrail"
provider: "aws"
icon: "package"
order: 200
presets:
  - slug: "01-content-safety-baseline"
    rank: "01"
    title: "Content Safety Baseline"
    excerpt: "This preset creates a guardrail covering all six harmful-content categories (high strength on the severe ones, medium on insults and misconduct), prompt-attack detection on inputs, and the..."
  - slug: "02-pii-redaction"
    rank: "02"
    title: "PII Redaction"
    excerpt: "This preset creates a guardrail focused on sensitive-information handling: contact details are masked (ANONYMIZE — the model still answers, with `{NAME}`-style placeholders), while payment data,..."
  - slug: "03-topic-denylist-rag"
    rank: "03"
    title: "Topic Denylist for RAG"
    excerpt: "This preset creates a guardrail for retrieval-augmented assistants: two denied topics (legal and medical advice) on the STANDARD cross-lingual tier, plus contextual grounding thresholds that reject..."
---

# Bedrock Guardrail Presets

Ready-to-deploy configuration presets for Bedrock Guardrail. Each preset is a complete manifest you can copy, customize, and deploy.
