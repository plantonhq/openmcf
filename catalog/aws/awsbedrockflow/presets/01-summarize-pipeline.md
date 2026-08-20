# Summarize Pipeline

This preset creates the canonical minimal flow — Input → inline Prompt →
Output — summarizing whatever document the caller sends, at temperature 0
on Amazon Nova Micro.

## When to Use

- The starting skeleton for ANY flow: deploy it, invoke it, then grow
  the graph a node at a time against AWS's strict topology validation
- Single-step model tasks (summarize, classify, extract) that want
  flow-style invocation and versioning

## What You Get

- A three-node graph with typed String sockets and two data connections
- An inline prompt (no separate prompt resource to manage) with one
  `{{input}}` variable

## Customize

- Swap the inline prompt for a Prompt Management resource:
  `prompt.promptArn` referencing an `AwsBedrockPrompt`'s output
- Add a Condition node after the prompt to branch on the completion
- Raise `temperature` for creative tasks
