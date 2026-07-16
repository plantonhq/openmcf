# Weighted Canary Origin

This preset adds a low-weight origin beside the main backend: same
priority, ~5% of traffic. Ramping the canary is a weight change -- an
in-place update -- and rolling back is deleting one resource.

## When to Use

- Canary releases: a new build takes a trickle of real traffic before
  the fleet moves
- Gradual regional cutovers and A/B backend experiments

## Key Configuration Choices

- **Same `priority` as the main origin** -- weights distribute traffic
  only WITHIN a priority tier; a canary at priority 2 would receive
  nothing until every priority-1 origin failed (that is the
  active/passive shape, not the canary shape)
- **`weight: 50` vs the main origin's 950** -- roughly 5%; weights are
  relative, so any ratio works
- **The origin is the rollback unit** -- deleting this resource
  restores 100% to the main origin without touching it

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<origin-group-resource-name>` | The AzureFrontDoorOriginGroup's Planton resource name | Your Front Door composition |
| `originName` (example value) | 2-90 chars -- rename to your convention | Your naming convention |
| `<canary-backend-hostname>` | The canary deployment's hostname | Your canary environment |
