# Planton Assistant Automations

This tree is the source of truth for the automation definitions the
Planton Assistant runs while nobody is watching. An automation is
published data — authored here, versioned by release, published to the
AI engine by the platform's publish lane — never a database row. What an
organization stores is only its own on/off switch per automation; the
definition itself is identical on every deployment, and the settings
page an org reads renders the exact document the engine enforces.

## Structure contract

```
automations/
  <automation-slug>.yaml   # one file per automation; the file name
                           # (without extension) must equal the slug
  README.md
```

Each file is one automation definition: slug, display name,
description, class (read or mutation), trigger bindings, budget, pinned
model, daily run cap default, and the engine workflow document. The
schema lives in the Planton platform and is enforced by every consumer
of this tree; the laws that matter when authoring:

- **No budget, no publish.** Every definition carries a positive
  per-run cost ceiling (`budget.maxCostMicros`). The exceeded policy is
  not authorable — the publish lane always applies
  terminate-on-exceeded, so a definition cannot soften its own failure
  mode.
- **No pinned model, no publish.** An unattended run never floats on a
  default model that can change underneath its budget. Declare the
  model once (`pinnedModel`); the publish lane injects it into every
  model-bearing task.
- **No declared class, no publish.** `automationClass` states the
  consent contract: `automation_class_read` automations only read and
  report; `automation_class_mutation` automations may propose changes,
  and every proposal waits for a human's per-action approval on the
  platform — an org's allow-switch never waives that.
- **At least one trigger binding**, naming a real platform resource
  kind and event type. Unknown names refuse loudly at parse time.
- **Slugs are permanent.** Org switch records adopt an automation by
  its slug; rename the display name freely, never the slug.

## The trigger predicate grammar (CEL)

A trigger binding's optional `predicate` is a [CEL](https://cel.dev)
expression evaluated by the platform's dispatcher against the event
resource's document, bound as the variable `resource` in its
protobuf-JSON form — field paths are camelCase-tolerant JSON paths and
enum values read as their NAMES, as strings:

```
resource.status.progress.status == "completed" && resource.status.progress.result == "failed"
```

The predicate must evaluate to a boolean. A predicate that does not
compile refuses to publish; a predicate that errors at evaluation (a
path the document lacks) fails closed — the run simply does not
dispatch. Predicates see the state the EVENT announced, never a
re-read, and they must tolerate refires: events carry no previous
state, so a re-saved terminal record fires its binding again (the
dispatcher's deterministic run identity makes refires idempotent).

## Runtime values in workflow tasks

Task configs reference per-run values the platform threads into each
execution using the engine's runtime placeholder grammar — no space
after `${`, namespaced by sensitivity:

- `${.env_vars.KEY}` for plain values
- `${.secrets.KEY}` for secret values (encrypted at rest, redacted in logs)

The platform threads these names into every run: `PLANTON_API_BASE_URL`
(the platform API), `PLANTON_RUN_CREDENTIAL` (the run's short-lived
read-only credential — always reference it through `${.secrets....}`),
`PLANTON_TRIGGER_RESOURCE_ID` (the resource whose event triggered the
run), and `PLANTON_ENGINE_ORG` (the org's engine workspace). A plain
`${KEY}` matches nothing and passes through as literal text.

Two grammars, one hazard: a space after `${` makes a jq expression
(`${ . | tojson }`), evaluated over the task's input; no space makes a
runtime placeholder (`${.env_vars.KEY}`). A task's input is the previous
task's output — that chaining, not any implicit context, is how a fetch
task's document reaches the task after it. And an `agent_call` passes
the agent ONLY its resolved `message`: facts the agent must see are
embedded in the message (jq's `tojson` renders an object as JSON text),
never assumed present.

## The finding output contract

An investigation's conclusion becomes a finding record the organization
reads — a card with a title, a one-glance summary, and a full markdown
body. The `agent_call` that produces it declares the engine's structured
output contract so the answer is extracted and validated, never parsed
from prose:

```yaml
output:
  schema:
    type: object
    required: [title, summary, content]
    properties:
      title: { type: string }     # one line naming the root cause
      summary: { type: string }   # one paragraph, readable alone in a feed
      content: { type: string }   # the full finding, markdown
  on_invalid: ON_INVALID_RETRY
  max_retries: 1
```

Three authoring rules ride it: the engine validates a schema SUBSET
(`type`, `required`, `properties`, `enum`) — put length guidance in the
message, never in schema constraints that would silently not be
enforced; each `on_invalid` retry re-runs the full agent call (real
spend inside the run's budget), so keep `max_retries` at 1; and keep the
whole answer bounded (roughly 2000 words) — oversized task outputs are
truncated or offloaded by the engine, and a bounded answer stays inline
where the platform reads it directly.

## Contributing

Pull requests proposing new automations are welcome — this tree is
public precisely so anyone can read, and propose, what the assistant
can do. Hold one distinction while authoring: a skill is knowledge (its
worst failure is bad advice in a conversation); an automation is
capability, running unattended with a scoped credential and a budget.
Every merged automation is therefore reviewed like platform code —
credential scope, budget, trigger blast radius, and consent class all
get a security review before a definition ships, regardless of author.

## How releases carry this content

Every semver release of this repository packages each automation
definition beside the skill archives and records it in
`definitions-manifest.json` with a SHA-256 checksum. The platform's
publish lane downloads a release, verifies every checksum and the
schema laws above, and only then places the workflow on the AI engine —
a definition that fails any law refuses to publish, loudly.
