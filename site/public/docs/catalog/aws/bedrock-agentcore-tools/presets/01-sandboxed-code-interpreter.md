---
title: "Sandboxed Code Interpreter"
description: "This preset gives an agent the safest possible compute: a code interpreter with NO network access and no AWS credentials — model-written code can calculate, parse, and plot, and nothing else."
type: "preset"
rank: "01"
presetSlug: "01-sandboxed-code-interpreter"
componentSlug: "bedrock-agentcore-tools"
componentTitle: "Bedrock AgentCore Tools"
provider: "aws"
icon: "package"
order: 1
---

# Sandboxed Code Interpreter

This preset gives an agent the safest possible compute: a code
interpreter with NO network access and no AWS credentials — model-written
code can calculate, parse, and plot, and nothing else.

## When to Use

- Any agent that runs model-generated code (the default posture)
- Data analysis, file transformation, chart generation

## What You Get

- An isolated sandbox per session, billed only while sessions run
- Zero egress: exfiltration by generated code is structurally impossible

## Customize

- Switch to `mode: PUBLIC` with an `executionRoleArn` only when the code
  genuinely needs the internet or AWS APIs — and scope that role tightly
- Add `certificates` (Secrets Manager) when the code calls
  mTLS-protected internal endpoints
