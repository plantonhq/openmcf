# Azure Monitor Data Collection Rule Association -- Operational Guide

Judgment calls that matter when you run association-based onboarding in production.

## The association is the unit of fleet membership

Treat "which machines feed which rules" as association inventory, not rule edits. Onboarding is CREATE an association; offboarding is DESTROY it -- the rule, its destinations, and every other machine's collection are untouched either way. When a chart owns a machine, put the machine's associations in the same chart: the machine and its monitoring share one lifecycle, and tearing the chart down detaches cleanly.

## Name associations for the reader of a fleet listing

Association names live under the machine and appear in every "what feeds this rule" listing. A convention like `<rule-shortname>-assoc` tells an operator at a glance what a machine is bound to; `assoc1` guarantees a future archaeology session. The name is ForceNew -- renaming recreates the association (harmless, seconds, but a plan diff worth understanding).

## Layering rules on one machine is a feature -- with one namespace

A machine carrying the fleet baseline rule PLUS a workload rule is the intended composition. Azure evaluates all of a machine's associations together, and duplicate association NAMES on one machine collide -- which is why per-rule naming (not generic names) matters as fleets layer.

## The agent is a separate concern, on purpose

The association creates successfully on a machine that does not run the Azure Monitor Agent -- that is Azure's design, not an error: configuration and agent rollout are decoupled. Operationally this cuts both ways: an association that "does nothing" usually means no agent, and an agent with no associations collects nothing. Check the machine's extensions blade first when telemetry is missing.

## Endpoint associations are singular and name-locked

A machine carries at most ONE endpoint association, and Azure mandates its name (`configurationAccessEndpoint` -- leave the name unset and the engines apply it). Use it only when configuration access must traverse private networking; machines on open egress do not need one.

## Destroy order in charts: association before target or rule

The association dies automatically when its target machine is deleted (extension-resource semantics), but a clean chart teardown destroys associations explicitly before the rule -- Azure otherwise briefly holds rule deletion while associations reference it. Planton's reverse-dependency destroy order produces this naturally when the association references both by `valueFrom`.
