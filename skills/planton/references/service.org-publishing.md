---
title: Publishing Pipelines and Tasks for the Organization — The Records, the Publish Check, Consuming by Name, and Every Refusal's Fix
description: How an organization shares one pipeline or one Task across many services exactly the way the platform compiles it — the TektonPipeline and TektonTask record shapes as copyable YAML, the name law (record name equals the Tekton document name), publish order, the publish check that refuses a record before it is written and every refusal's fix, consuming a published pipeline with one line and a published task with a plain taskRef, what a run stamps and what an update changes, the delete guards, and how to list and read what the organization has published. Read when someone asks to share a pipeline or a task across services, publish or update an organization pipeline, point a service at a published pipeline, explain a publish refusal or a delete refusal, or find what pipelines the organization has.
---

# Publishing Pipelines and Tasks for the Organization

Authoring one service's pipeline is `references/service.pipeline-authoring.md`; this file is the step after it — the same Tekton, published once as a record, consumed by name from every service. You are the developer's only teacher here: they will say "share storefront's pipeline with the other Node services", and you write the records, publish them in the right order, and switch each service with one line. Nothing is fetched from git at run time; a published record's content IS the compile source.

## The two records — copyable shapes

A **TektonTask** is exactly one `tekton.dev/v1` Task document. A **TektonPipeline** is one `tekton.dev/v1` Pipeline plus any Task documents it bundles (`---`-separated, exactly as `kubectl apply -f` accepts). Both are organization-owned (`metadata.org`); there is no platform variant — the platform's own pipelines and tasks are the open-source catalog, never records.

```yaml
# eslint-check.task.yaml
apiVersion: service-hub.planton.ai/v1alpha1
kind: TektonTask
metadata: {name: eslint-check, org: acme}
spec:
  description: ESLint over the checked-out source
  overviewMarkdown: |
    Params: none. Workspaces: `source`. Runs `npm ci && npx eslint .`
  yamlContent: |
    apiVersion: tekton.dev/v1
    kind: Task
    metadata: {name: eslint-check}
    spec:
      workspaces: [{name: source}]
      steps:
        - name: lint
          image: node:20-alpine@sha256:fb4cd12c85ee03686f6af5362a0b0d56d50c58a04632e6c0fb8363f609372293
          workingDir: $(workspaces.source.path)
          script: npm ci && npx eslint .
```

```yaml
# node-service-ci.pipeline.yaml
apiVersion: service-hub.planton.ai/v1alpha1
kind: TektonPipeline
metadata: {name: node-service-ci, org: acme}
spec:
  description: Clone, lint, build with BuildKit, render kustomize — every Node service
  overviewMarkdown: |
    ## What your service must supply
    - `spec.build.registry` — this pipeline pushes an image (declares `image-name`).
    - `build.tektonPipeline.params.node-version` — the Node major to lint with.
  yamlContent: |
    apiVersion: tekton.dev/v1
    kind: Pipeline
    metadata: {name: node-service-ci}
    spec:
      params:
        - {name: git-url, type: string}
        - {name: git-revision, type: string}
        - {name: git-branch, type: string, default: ""}
        - {name: project-root, type: string, default: "."}
        - {name: sparse-checkout-directories, type: string, default: ""}
        - {name: owner-identifier-label-key, type: string, default: ""}
        - {name: owner-identifier-label-value, type: string, default: ""}
        - {name: kustomize-manifests-config-map-name, type: string}
        - {name: kustomize-base-directory, type: string, default: "_kustomize"}
        - {name: image-name, type: string}
        - {name: node-version, type: string}
      workspaces: [{name: source}]
      tasks:
        - name: git-checkout
          taskRef: {name: git-clone}          # platform catalog
          params:
            - {name: url, value: "$(params.git-url)"}
            - {name: revision, value: "$(params.git-revision)"}
          workspaces: [{name: output, workspace: source}]
        - name: lint
          runAfter: [git-checkout]
          taskRef: {name: eslint-check}       # the organization's published Task
          workspaces: [{name: source, workspace: source}]
        # build-image and kustomize-build exactly as in the repository worked example
```

**The name law**: `metadata.name` on the record must equal the `metadata.name` inside the Tekton document (the Pipeline, or the one Task), and it must be a Tekton name — lowercase letters, digits, hyphens, at most 63 characters. A Tekton name is also the record's slug, so `planton get TektonPipeline node-service-ci`, `build.tektonPipeline.pipeline: node-service-ci`, and a `taskRef: {name: eslint-check}` are all the same one name. Tasks bundled inside a pipeline record are named freely.

`description` is one line for search results. `overviewMarkdown` is where you write what a consumer must do — the params the service must pass, whether it must name a registry; it is what `planton get` and the assistant read.

## Publish — one apply each, tasks first

`planton apply -f eslint-check.task.yaml`, then `planton apply -f node-service-ci.pipeline.yaml`. A pipeline that references a Task the organization has not published yet is refused (`unresolved_task_ref`) — publish the Task first, or bundle it as a second document in the pipeline record. Re-publishing is the same apply with the changed content.

**The publish check runs before anything is written**, at every door — `planton apply`, the assistant's `apply_tekton_pipeline` / `apply_tekton_task` tools, and the API itself. It is the dispatch compiler on the record's content: the document shape, the name law, task references resolved against the bundled Tasks → the organization's published Tasks → the platform catalog, the platform's parameter contract (every always-supplied param must be declared), workspaces (`source` only), result references, the secret tripwire, and Tekton's own admission rules. A refusal is one sentence carrying every problem: `TektonPipeline 'node-service-ci' would not compile, 2 problem(s): unresolved_task_ref(eslint-check): …; undeclared_param(git-branch): …`. Fix them all, apply again.

## Consume — one line, or one taskRef

- **A published pipeline**: in the service's `service.yaml`, `spec.build.tektonPipeline: {pipeline: node-service-ci}` plus whatever `overviewMarkdown` asks for (`registry`, `params: {node-version: "22"}`). Verified when the service is applied: a name no record carries is refused there — `build.tektonPipeline.pipeline 'node-service-cl' names no published TektonPipeline in organization 'acme'. Publish it first (planton apply -f <pipeline-record>.yaml), or name an existing one (planton search --kind TektonPipeline).` — never at the first build.
- **A published task**: `taskRef: {name: eslint-check}` in any repository pipeline or published pipeline of the organization; it resolves after the pipeline's own bundled Tasks and before the platform catalog, and the run record shows the rung.

What a consuming run stamps: `spec.resolvedPipeline.source: org_published`, `pin: <record id>@<record spec updated-at seconds>` — the exact content state that compiled. **Updating a record** changes what the next FRESH run compiles; reruns replay their stamp unchanged (see `references/service.reading-a-run.md`).

## Delete — refused while anything depends on it

`planton delete TektonPipeline node-service-ci` is refused while any service in the organization names it: the refusal lists the services — repoint each `build.tektonPipeline` first. `planton delete TektonTask eslint-check` is refused while any PUBLISHED pipeline references it, naming them. Repository pipelines that reference a deleted Task learn at their next compile (`unresolved_task_ref`) — the honest limit; check `planton search` for repository consumers you know of before deleting.

## Listing and reading what the organization has

`planton search --kind TektonPipeline` / `--kind TektonTask` lists the organization's published records (the search index carries `description`); `planton get TektonPipeline <name>` shows the record with its `overviewMarkdown` and `yamlContent`. Over MCP: `get_tekton_pipeline` / `get_tekton_task` (by ID or org + name) and search by kind. There is no catalog page — published content is read where it is used.

## Asked for → what you do

- **"Share this pipeline with our other services"**: copy the repository pipeline's Pipeline document into a `TektonPipeline` record named as the Pipeline; move each repository Task it uses into a `TektonTask` record (or bundle them as documents); write `overviewMarkdown` listing what consumers must supply; publish tasks, then the pipeline; switch each service to `tektonPipeline: {pipeline: <name>}`; delete the repository copies.
- **"Let each service choose its Node version"**: declare `node-version` without a default in the published Pipeline; document it in `overviewMarkdown`; each service sets `tektonPipeline.params.node-version`. A service that forgets fails its first build with `missing_required_param(node-version)`.
- **"Every service should use our clone step"**: publish a `TektonTask` named `git-clone` — it shadows the catalog's for every pipeline in the organization (bundled and repository Tasks still win over it).
- **"Change the shared pipeline"**: edit the record's `yamlContent`, `planton apply -f` — the publish check re-runs; consuming services pick it up at their next push.

## Every refusal → the fix

The compile verdicts and their fixes are in `references/service.pipeline-authoring.md`; publishing adds these and sharpens two:

- `invalid_record_name` (subject: the record name) — make `metadata.name` a Tekton name (lowercase, digits, hyphens; ≤63) and rename the document to match.
- `record_name_mismatch` (subject: the record name) — set the record's `metadata.name` and the document's `metadata.name` to the same value.
- `unexpected_document_kind` on a TektonTask (subject: `<name> document 2`) — a Task record holds exactly one document; publish each Task as its own record.
- `unresolved_task_ref` at publish — publish the referenced `TektonTask` first, bundle the Task as a document, or fix the name.
- `undeclared_param` at publish (subject: a contract param) — the platform supplies it to every build; declare it under `spec.params` (a default is fine).
- Delete refusals (`FAILED_PRECONDITION`) name the dependents — repoint or edit them, then retry.
