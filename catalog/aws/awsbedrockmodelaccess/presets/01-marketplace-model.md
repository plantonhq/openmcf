# Marketplace Model

This preset accepts the public offer for a third-party marketplace model
(Cohere Command R) — the simple access shape: no use-case form,
pay-per-token pricing with no fixed fee, revocable by destroying the
component.

## When to Use

- Enabling any marketplace model that carries a public offer (probe with
  `aws bedrock list-foundation-model-agreement-offers` — auto-enabled
  models reject the call and need no agreement)
- The access building block a chart deploys before agents and profiles

## Key Configuration Choices

- **Only the model id is declared** — the module resolves the offer token
  at deploy time (tokens are short-lived and never belong in a manifest).
- **Region matters**: agreements are regional; deploy one instance per
  region that needs the model.

## After Deployment

Downstream components (agents, inference profiles, throughput purchases)
should reference the `model_id` output for deploy ordering.
