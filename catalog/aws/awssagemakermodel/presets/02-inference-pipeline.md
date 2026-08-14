# Inference Pipeline

This preset chains two containers into a Serial inference pipeline —
raw requests hit the preprocessing container, its output feeds the
scoring container, and one endpoint serves the whole chain.

## When to Use

- Feature transformation that must run inline with every prediction
- Keeping preprocessing and scoring images independently versioned
  while deploying them as one unit

## What You Get

- A 2-container pipeline invoked in sequence
  (`inferenceExecutionMode: Serial`)
- Named hostnames (`preprocess`, `score`) so container logs and
  metrics stay attributable — and so a switch to `Direct` invocation
  needs no rename

## Customize

- Grow the chain up to 15 containers; order in `containers` is
  execution order
- Switch to `inferenceExecutionMode: Direct` to let callers address a
  single container by hostname instead of running the chain
- Any change replaces the model (models are immutable) — roll a new
  one and repoint the endpoint
