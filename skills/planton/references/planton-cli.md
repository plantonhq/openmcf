# The planton CLI — Lookups and Pipeline Diagnosis

The `planton` CLI is your window into the user's Planton: what exists, what
deployed, what failed and why. Everything here is read-only unless marked
otherwise. `chart build` has its own contract (`build-contract.md`); schema
grounding has `component-grounding.md`. Machine-readable output: most
snapshot commands take `-o json` (protobuf field names) or `-o yaml`.
On the platform-tools arm (no CLI — see SKILL.md, "Know your instruments")
every lookup below has a platform tool twin (`list_*`/`get_*` on charts,
projects, pipelines, stack jobs, connections): same questions, same
diagnosis order, tool-for-command.

## Context — whose org and environment am I in?

```
planton context get                      # effective org + environment
planton context set --org <org> -e <env> # change it (a mutation — confirm first)
```

On local desktop instances the context auto-resolves to the seeded org/env —
commands work with zero setup. Most commands also accept `--org` / `-e`
inline for one-off scoping.

## Grounding lookups (session start, before proposing architecture)

```
planton chart list                        # charts in the catalog (org + official)
planton infra-project list -o json        # the org's projects (what has been deployed)
planton infra-project get <id-or-name>    # one project, full record
planton env list                          # environments in the org
planton search connections                # provider connections (which clouds/clusters)
planton connection authorization list     # which connections which envs may use
planton secret list -o json               # managed secrets ("env" field = scope)
planton variable list -o json             # managed variables ("env" field = scope)
```

The secret/variable lists ground `$var`/`$secret` references before you write
them (`config-references.md`); `-o json` matters because only the JSON
records carry each entry's `env`, which decides the reference form. Creating
one (`planton secret set` / `planton variable set`) is a mutation — one
confirmation, same as any other.

Read `infra-project list` results with an eye on `env` and
`infra_project_source` — a chart-sourced project in env `dev` means that
chart's resources exist in dev. Caveat: the list rides the search index, so a
project created seconds ago may briefly be missing; `infra-project get`
by id/name is the direct read.

## The deployment chain — ids you will meet

`InfraChart → InfraProject (infproj_…) → InfraPipeline (infpipe_…) → per-node
CloudResource → stack job (sj_…)`. Diagnosis walks this chain downward
(the model is explained in `deployment-model.md`).

## Diagnosing a failed deploy (the four-step workflow)

```
# 1. Which pipeline? (newest first)
planton infra-pipeline list <project-id-or-name>
#    or skip the lookup: planton infra-project last-pipeline <project-id>

# 2. What happened? One-shot snapshot, per-node results:
planton infra-pipeline status <infpipe_id>
#    → table: NODE | KIND | STATUS | RESULT | STACK_JOB | REASON
#    -o json for the full record (every node's execution, stack-job ids)

# 3. Why did the node fail? Engine logs, one hop:
planton infra-pipeline logs <infpipe_id>              # auto-targets the failed node
planton infra-pipeline logs <infpipe_id> --node <slug> # a specific node (slug from status)

# 4. Deeper: the stack job record carries the full error detail:
planton stack-job list <cloud-resource-id> -o json
planton get stack-job <sj_id> -o json
planton stack-job logs <sj_id>            # tail the engine output directly
```

Typical reading of `status` output: a node with result `failed` and a
`REASON` mentioning "No provider connection available" or
"kubernetes-provider-connection … not found" is the wiring class — fix per
`kubernetes-on-cluster.md` / `issue-catalog.md`. A failure naming cloud
concepts (subnets, IAM, quotas) is an infrastructure class — the engine logs
from step 3 carry the provider's own error, and `deployment-model.md`
explains how to map it back to the module and spec field. A failure saying a
resource **already exists** is the orphaned-resource class (an earlier run
created it, then died before recording it in state) — repair by importing,
never by delete-and-retry: read `state-import.md`.

## Watching a running deploy (humans; agents prefer snapshots)

```
planton follow <infpipe_id|sj_id> --plain   # auto-detects by id prefix
planton infra-pipeline stream-status <id>   # status lines until terminal
```

Streams block until the pipeline finishes — in a session, prefer polling
`status` between other work over holding a stream open.

## Generic escape hatches

Any resource, any kind: `planton get <kind> <id> -o json` (e.g.
`get infra-pipeline`, `get cloud-resource`, `get stack-job`) and
`planton search by-resource-kind <kind>` for kinds without a noun-scoped
list. The noun-scoped verbs above are the preferred, discoverable path.
