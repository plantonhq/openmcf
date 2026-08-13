# Support Agent with Tools

This preset creates a customer-support agent on Amazon Nova Micro with a
return-control order-lookup tool, the reserved `AMAZON.UserInput`
capability (so the agent can ask clarifying questions), and one `live`
alias serving the assembled agent.

## When to Use

- The starting point for a tool-using assistant whose operations your
  application executes (return-control: the agent hands the tool call
  back to you)
- Teams that want the alias-versioned serving model from day one

## What You Get

- An agent with clear instructions, a 15-minute idle session window, and
  two action groups: your `orders` function tool and AWS's user-input
  elicitation capability
- A `live` alias — the immutable snapshot your application invokes; edit
  the spec freely and re-alias when ready to ship

## Customize

- Swap `returnControl: true` for a Lambda executor
  (`executor.lambda` referencing your function) to run tools server-side
- Add `guardrail` to attach a Bedrock guardrail at a pinned version
- Add `memory` to carry session summaries across conversations
