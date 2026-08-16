# AwsSagemakerEndpoint — Component Guide

Authored operational judgment for the SageMaker endpoint component:
the design decisions behind the spec's shape, and what to know before
running real-time inference in production.

## Design decisions

- **The endpoint's AWS name derives from `metadata.name`.** Clients
  invoke the endpoint by that name; it never changes across
  configuration rolls.
- **The endpoint configuration is folded into the endpoint.** Upstream
  they are two resources, but the configuration is immutable (every
  argument ForceNew) while the endpoint's pointer to it updates in
  place. Keeping them apart would push AWS's roll-and-repoint
  choreography onto every author; folding them lets the modules own it
  — configurations are name-suffixed (`<name>-cfg-…`) and created
  before the old one is destroyed, UpdateEndpoint repoints, and only
  then does the predecessor disappear. The endpoint never references a
  deleted configuration.
- **Variant names default deterministically per position**
  (`variant-0`, `variant-1`, …; `shadow-variant-0` on the shadow
  side). The provider would otherwise mint a random name per plan,
  forcing a configuration roll on every apply.
- **Serverless and instance settings are exclusive per variant.** AWS's
  ServerlessConfig contract — a serverless variant rejects
  `instance_type` and every instance-tuning field; the spec enforces
  the split at manifest time.
- **Shadow testing is exactly one variant on each side.** AWS's own
  shape — the spec rejects other combinations before the API does.
- **Exactly one deployment strategy.** When `deployment` is set, it is
  `blue_green` or `rolling`, never both; canary and linear step sizes
  pair with their routing types (AWS's rules, spec-enforced).

## Running endpoints in production

- **Start serverless, graduate to instances.** A serverless variant
  costs nothing idle and needs no capacity math; move to
  instance-backed variants when latency floors, GPU needs, or
  provisioned-concurrency costs say so.
- **Shape the crossing before you need it.** Configuration rolls are
  routine (every capacity change is one) — give production endpoints a
  `deployment` policy with `auto_rollback_alarm_names` so a bad model
  version backs itself out instead of paging you.
- **Weight, don't flip.** New model versions ride a second variant at
  a small `initial_variant_weight`; a weight of 0 keeps a variant
  deployed but takes no traffic — an instant rollback target.
- **Mind the KMS caveat.** `kms_key_arn` encrypts ML storage volumes,
  but nitro-local-storage families (ml.g5/g6/p4d/p5 and similar)
  encrypt locally by default and reject a custom key — as do
  serverless-only endpoints.
- **Capture from day one.** `data_capture` is the Model Monitor feed;
  flipping it on later is itself a configuration roll, so sample a low
  percentage early rather than retrofitting under incident pressure.
- **Budget bake time.** Blue/green provisions a full parallel fleet
  and holds it through `wait_interval_seconds` per step (plus
  `termination_wait_seconds` after the shift); rolling replaces
  batches in place with no parallel-fleet cost — pick by budget and
  blast radius.

---

© Planton. Licensed under [Apache-2.0](https://github.com/plantonhq/planton/blob/main/LICENSE).
