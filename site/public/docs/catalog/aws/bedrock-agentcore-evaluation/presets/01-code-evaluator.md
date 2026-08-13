---
title: "Code Evaluator"
description: "This preset creates a single code-based evaluator backed by a Lambda — the cheapest first evaluation object. Create does not invoke the function."
type: "preset"
rank: "01"
presetSlug: "01-code-evaluator"
componentSlug: "bedrock-agentcore-evaluation"
componentTitle: "Bedrock AgentCore Evaluation"
provider: "aws"
icon: "package"
order: 1
---

# Code Evaluator

This preset creates a single code-based evaluator backed by a Lambda
— the cheapest first evaluation object. Create does not invoke the
function.

## When to Use

- The first AgentCore evaluator in an environment
- Scoring logic you already own in a function
- Accounts that should not depend on Bedrock model access yet

## What You Get

- One TRACE-level evaluator named `custom_score`
- The Lambda ARN composed from the named AwsLambda

## Customize

- Point `lambdaArn` at your scorer function
- Add `timeoutSeconds` when the function needs more than the AWS
  default
- Add an LLM-judge evaluator from the second preset when you are
  ready to score with a model
