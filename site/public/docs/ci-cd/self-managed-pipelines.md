---
title: "Self-Managed Pipelines"
description: "Bring your own Tekton pipeline — and the Tasks it references — in your repository, compiled by Planton at dispatch into a self-contained run with platform deployment, status tracking, and gates."
icon: pipeline
order: 35
tags:
  - Pipelines
  - Tekton
  - CI/CD
  - Service Hub
---

# Self-Managed Pipelines

A self-managed pipeline is plain Tekton YAML in your repository. Planton compiles it at dispatch into one self-contained run — every task reference resolved inline, the platform's parameter contract verified, the compiled definition stamped on the run record — and then handles what it always handles: webhook triggering, credentials, status streaming, the deploy stage, and approval gates.

## Platform track or your own pipeline

A service's `spec.build` selects exactly one builder:

| Builder | Who writes the pipeline | Deploy stage |
|---|---|---|
| `dockerfile` | Planton's release-pinned dockerfile track | Planton |
| `buildpacks` | Planton's release-pinned buildpacks track | Planton |
| `tektonPipeline` | You, as Tekton YAML in your repository (or an organization-published pipeline) | Planton |

Choose your own pipeline when you need build steps the platform tracks do not have: test suites, linting, security scanning, a build tool like Bazel, parallel or conditional tasks.

## The pipeline file

By default the pipeline lives at `.planton/pipeline.yaml` under the service's project root — beside the service's Dockerfile and its `_kustomize` tree. Selecting it is one line:

```yaml
spec:
  build:
    tektonPipeline: {}
    registry: my-registry
    imageRepositoryPath: acme/storefront
```

To keep it elsewhere, name the path (relative to the project root):

```yaml
spec:
  build:
    tektonPipeline:
      yamlFile: ci/build.yaml
```

The file is a `tekton.dev/v1` `Pipeline` document. It is read at the exact commit being built — never from a branch.

## Tasks beside the pipeline

A pipeline references tasks by their plain Tekton name:

```yaml
tasks:
  - name: lint
    taskRef:
      name: eslint-check
```

The compiler resolves a plain reference from three places, most specific first:

1. **Your repository** — a `Task` document in a `tasks/` directory beside the pipeline file (`.planton/tasks/*.yaml` for the default location), or a second document in the pipeline file itself (`---`-separated, exactly as `kubectl apply -f` accepts).
2. **Your organization** — a published `TektonTask` record with that name.
3. **The platform catalog** — `git-clone`, `buildkit-daemonless`, `buildpacks`, `kustomize-build` (their definitions live in the open-source repository under `cicd/tekton/tasks/`).

A repository task with the same name as a catalog task wins — your file is what runs, and the run record shows it. A `taskRef` that uses a remote resolver (`resolver: git`, `hub`, `bundles`) is refused: compiled pipelines are self-contained, so nothing is fetched at run time.

A repository task is an ordinary Tekton `Task`:

```yaml
# .planton/tasks/eslint-check.yaml
apiVersion: tekton.dev/v1
kind: Task
metadata:
  name: eslint-check
spec:
  workspaces:
    - name: source
  steps:
    - name: lint
      image: node:20-alpine@sha256:fb4cd12c85ee03686f6af5362a0b0d56d50c58a04632e6c0fb8363f609372293
      workingDir: $(workspaces.source.path)
      script: |
        npm ci && npx eslint .
```

## The parameter contract

Planton supplies one set of parameters to every build, and the compiler verifies it against your pipeline's `spec.params` in both directions before dispatch: a supplied parameter your pipeline does not declare is a compile error, and a declared parameter with no default that nothing supplies is a compile error. Your pipeline must therefore declare every one of these (defaults are fine):

| Parameter | Value |
|---|---|
| `git-url` | The repository clone URL |
| `git-revision` | The commit SHA being built |
| `git-branch` | The branch name (empty for a tag build) |
| `project-root` | The service's project root within the repository |
| `sparse-checkout-directories` | Comma-separated sparse-checkout patterns (may be empty) |
| `owner-identifier-label-key`, `owner-identifier-label-value` | Stamp these on everything the pipeline creates — log attribution selects by them |
| `kustomize-manifests-config-map-name` | Where the deploy stage reads the rendered overlays |
| `kustomize-base-directory` | The service's kustomize base directory |
| `image-name` | Present when the build names a registry: `<registryHost>/<imageRepositoryPath>:<commitSha>` |
| `dockerfile-config-map-name`, `build-context`, `dockerfile-path` | Dockerfile builds only |

Plus anything you put under `build.tektonPipeline.params` (free-form key/value pairs; they cannot shadow the platform's names).

### Trigger facts — supplied only if you declare them

Three more parameters describe the run, and the platform supplies them **only when your pipeline declares them**, so an existing pipeline never has to change:

| Parameter | Value |
|---|---|
| `git-tag` | The tag name for a tag build; empty otherwise |
| `pull-request-number` | The pull request number for a pull-request build; empty otherwise |
| `trigger-type` | `webhook`, `manual`, or `rerun` |

Gate a task with Tekton's own `when` — no platform grammar to learn:

```yaml
spec:
  params:
    - name: git-tag
      type: string
      default: ""
  tasks:
    - name: sbom
      when:
        - input: "$(params.git-tag)"
          operator: notin
          values: [""]
      taskRef:
        name: generate-sbom
```

## Workspaces

The platform binds one pipeline-level workspace: `source`, the cloned repository. Declare it and pass it to your tasks. A pipeline that requires any other workspace is refused at compile time — declare it `optional: true` or drop it. Secrets are always runtime references (secret names, mounts), never values in a definition.

## Validate before you push

```bash
planton service pipeline validate .planton/pipeline.yaml \
  --param git-tag=v1.4.0 -o json
```

The same compiler that runs at dispatch runs locally: it discovers the tasks beside the file, resolves catalog references, checks the parameter contract, and reports every problem in one pass (`--task <name>=<path>` supplies a task the way an organization-published record would). The JSON carries `valid`, `source`, `pin`, `compiler_version`, the resolved tasks and where each came from, and every verdict; the exit code is 1 on any verdict. A pipeline that validates clean compiles clean.

## Compile verdicts

| Code | Meaning | Fix |
|---|---|---|
| `invalid_yaml` | A document is not valid YAML | Fix the YAML; the message names the file |
| `invalid_pipeline_schema` | The document is not a `tekton.dev/v1` Pipeline, or fails Tekton's own validation | Correct `apiVersion`/`kind`, or the field Tekton names |
| `multiple_pipeline_documents` | The source contains more than one Pipeline | Keep one Pipeline per source; extra documents may only be Tasks |
| `unexpected_document_kind` | A document beside the pipeline is neither a Pipeline nor a Task | Remove it, or move it out of the pipeline's `tasks/` directory |
| `duplicate_task_definition` | Two repository documents define a Task with the same name | Rename one |
| `unresolved_task_ref` | A `taskRef` names nothing in your repository, your organization, or the catalog | Add the Task beside the pipeline, publish it, or fix the name |
| `resolver_ref_unsupported` | A `taskRef` uses a remote resolver | Reference the task by name; inline it or add it beside the pipeline |
| `undeclared_param` | The platform supplies a parameter the pipeline does not declare | Declare it under `spec.params` |
| `missing_required_param` | A declared parameter has no default and nothing supplies it | Give it a default, or set it under `build.tektonPipeline.params` |
| `unbindable_workspace` | The pipeline requires a workspace the platform does not bind | Use `source`, or mark the workspace `optional: true` |
| `dangling_result_reference` | `$(tasks.X.results.Y)` names a task or result that does not exist | Fix the task or result name |
| `literal_secret_env` | An env var looks like a literal secret value | Use a secret reference or a mounted workspace |
| `repo_file_unreadable` | The pipeline file could not be read at the built commit | Check `yamlFile` and that the file exists at that commit |

## What a run records

Every run stamps the compiled definition on its record: the exact YAML that executed, `source: repo`, and the commit SHA as its pin. A rerun replays those bytes byte-identically. The run view and `planton service pipeline get` show them; see [Pipelines](/docs/ci-cd/pipelines) for the run model.

## Organization-published pipelines

An organization can publish a pipeline once and reuse it across services: `build.tektonPipeline.pipeline: <name>` compiles from the published `TektonPipeline` record's content at dispatch, and published `TektonTask` records resolve for plain references the same way. Publishing is one `planton apply -f` of the record.

## Related Documentation

- [Pipelines](/docs/ci-cd/pipelines) — the run model, stages, and gates
- [Build Methods](/docs/ci-cd/build-methods) — the platform's dockerfile and buildpacks tracks
- [Monorepo Support](/docs/ci-cd/monorepo-support) — project roots and trigger paths
