---
title: "Presets"
description: "Ready-to-deploy configuration presets for Bedrock Flow"
type: "preset-list"
componentSlug: "bedrock-flow"
componentTitle: "Bedrock Flow"
provider: "aws"
icon: "package"
order: 200
presets:
  - slug: "01-summarize-pipeline"
    rank: "01"
    title: "Summarize Pipeline"
    excerpt: "This preset creates the canonical minimal flow — Input → inline Prompt → Output — summarizing whatever document the caller sends, at temperature 0 on Amazon Nova Micro."
  - slug: "02-classify-and-route"
    rank: "02"
    title: "Classify and Route"
    excerpt: "This preset creates a branching flow: an inline prompt classifies the request, a Condition node routes docs questions through a knowledge-base query (retrieval + generation on Nova Lite) and..."
---

# Bedrock Flow Presets

Ready-to-deploy configuration presets for Bedrock Flow. Each preset is a complete manifest you can copy, customize, and deploy.
