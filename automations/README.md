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
