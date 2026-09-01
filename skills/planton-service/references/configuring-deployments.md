---
title: Reading and Changing a Service's Deployment Configuration
description: What a service is DECLARED to deploy per environment and who writes that declaration — reading deploy.environments honestly, the two writers (a git kustomize tree vs platform-authored), editing through get-then-apply without destroying authored content, adding and removing target environments, and the deployments on/off switch. Read when someone asks what an environment is configured to run, wants more memory or a new variable on a running service, wants a new target environment, or asks why deployments stopped happening.
---

# Reading and Changing a Service's Deployment Configuration

Every service's per-environment deployment configuration lives in ONE place on the record: `spec.deploy.environments` — one entry per target environment, each carrying the full cloud-resource manifests that deploy there. There are not two kinds of configuration; there are two kinds of WRITER, and knowing which one writes a given service is the first fact to establish before promising any change.

## The two writers, and how to tell them apart

Read the service (`get_service`) and look at `spec.deploy`:

- **`kustomize` present** — the repository's kustomize tree maintains the environments section, and the platform's sync is its ONLY writer. Each entry carries `sync` provenance: which branch, which commit, when. **A direct write to this section does not error — it is silently preserved-over**: apply succeeds and the stored entries stand exactly as they were. Never edit a git-maintained service's environments through apply; the honest answer is "this service's configuration lives in its repository — edit the kustomize tree" (see `references/kustomize-authoring.md`, including taking authorship back with eject).
- **`kustomize` absent** — the section is platform-authored: the console, an agent, or a direct apply writes it. This is the posture the rest of this file is about.

Everything else on the deploy block stays caller-authored on EVERY service, git-maintained included: `disable_deployments` (the pause switch), `hostname` (the serving label), and `branch_deployments` (branch-to-environment mappings).

## Reading an environment's declaration honestly

An entry's `resources` are the exact manifests that deploy — plus two injection facts worth stating when reporting them:

- **A blank image IS the injection slot.** The deploy fills it with the built artifact. On `KubernetesDeployment` the injection is unconditional (a declared image is REPLACED at deploy); on multi-container kinds like `GcpCloudRun`, only blank-image containers receive the artifact (declared images — sidecars — deploy as written).
- **Config references stay unresolved.** `$var/...` and `$secret/...` values in manifests resolve just-in-time at deploy; the declaration never carries secret material.

The console renders all of this on the service page's **Configuration tab**: the writer story with sync provenance, each environment's resources with their facts and exact YAML, a resource diagram when an environment declares two or more resources (wires drawn from the declaration's own `valueFrom` references), and — for platform-authored services — the editors. When a person would rather look than read JSON, point them there.

## What the console can edit (so you never undersell it)

On a platform-authored service the Configuration tab edits ANY declared resource, one editor per shape: the simple single-resource Kubernetes/Cloud Run shape gets the guided dials; every other resource — any catalog kind, any member of a multi-resource set — edits through that kind's own authored form, or as YAML in place for kinds without one. The tab also adds a resource to an environment (the catalog picker honors the org's catalog policy by omission — a policy-disabled kind is simply not offered) and removes one (never the last — removing the environment's configuration is that door). Environment variables edit as three row classes: literal text, secret references through the platform's secret picker (the manifest stores the `$secret/...` pointer, and a key-value secret asks for its key), and locked rows for `valueFrom` wiring, which stays YAML-authored for now. A declaration carrying content outside the kind's schema as this console knows it stays read-only with the reason named — so if a person says the console "won't let them edit" something, that sentence is the reason to relay, not a bug to work around.

## Changing configuration: get, edit surgically, apply

The loop is `get_service` → edit the JSON → `apply_service` with the whole record. Two disciplines make it safe:

1. **Fetch fresh, edit only what you mean.** Apply sends the whole record, so start from a fresh get (never a stale copy) and touch only the fields the person asked about. A manifest is a declaration other hands may have authored — probes, autoscaling, secret variables, wired `valueFrom` references. Changing the memory limit means changing the memory limit, not regenerating the manifest.
2. **Nothing deploys when you save.** Configuration changes take effect on the next deployment to the environment — the next push, or a promote into it. Say this every time: a person who saves and sees nothing change will otherwise read success as failure. (A rollback does NOT pick up new configuration — it re-applies a captured pair; see `references/delivery-verbs.md`.)

Refusals arrive as the server's own sentences and always name the real fix — a serving-hostname collision names both services and the environment; a branch mapped twice names the branch. Relay them verbatim.

## Adding and removing target environments

**Adding** an environment is adding one entry to `deploy.environments` — the slug plus its manifests. Copying an existing environment's resources and re-stamping `metadata.env` (and the `{service}-{env}` name convention, when the source followed it) is the faithful path for any resource shape. Nothing deploys at add time; the next push (or a promote into the new environment) performs the first deployment. Adding an entry can newly collide on the serving hostname — the refusal names it.

**Removing** an entry stops FUTURE deploys to that environment. State the no-teardown fact plainly: **anything already running there keeps running** — removing configuration never destroys resources. Tearing down deployed resources is the service delete cascade's job (`references/delete-cascade.md`), not a configuration edit's.

## The deployments on/off switch

`spec.deploy.disable_deployments: true` pauses deployment execution for the whole service: builds keep running and publishing artifacts, the configuration sync of git-maintained services continues, but nothing deploys and the delivery verbs refuse with a sentence pointing at the service's deploy settings. That switch lives on the console's Configuration tab, and it is caller-authored on every posture — flipping it through get-then-apply works on git-maintained services too. Resuming makes the next push deploy again; nothing deploys at the moment of resuming.

## What a configuration edit can never do

- It cannot change what is currently RUNNING (only the next deployment reads the declaration).
- It cannot touch a git-maintained service's environments (silently preserved-over — say so instead of trying).
- It cannot author `sync` provenance (user-supplied values are discarded) or resolve secrets into the record.
- It is refused entirely on a service that is being deleted.
