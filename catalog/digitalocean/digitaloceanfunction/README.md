# DigitalOcean Function

## Overview

**DigitalOceanFunction** deploys serverless functions as an App Platform application with a single functions component. DigitalOcean's Terraform provider has no standalone Functions resource; both Pulumi and Terraform create `digitalocean_app`.

Runtime, memory, timeout, entrypoint, and cron schedules are **not** on this spec. They live in the repo's `project.yml` (inside `sourceDirectory`). App Platform reads that file at deploy time. Putting those knobs on the spec would silently do nothing.

## When to use this kind vs DigitalOceanApp

Use **DigitalOceanFunction** when the functions component is the whole app. Use **DigitalOceanApp** `spec.functions` when functions should share an app with services, workers, or static sites.

## Sources

Set exactly one:

- **git** — public HTTPS clone URL. Use this when the DigitalOcean account has no linked GitHub/GitLab/Bitbucket connection.
- **github / gitlab / bitbucket** — `owner/repo` plus branch. `deployOnPush` needs the matching VCS connection.

Functions are git-only. There is no container-image source.

`sourceDirectory` is required. It is the directory that contains `project.yml` and the packages tree (for the official hello-world sample, that is `packages`).

## Example

```yaml
apiVersion: digital-ocean.planton.dev/v1alpha1
kind: DigitalOceanFunction
metadata:
  name: hello
spec:
  functionName: hello
  region: nyc3
  git:
    repoCloneUrl: https://github.com/digitalocean/sample-functions-nodejs-helloworld.git
    branch: master
  sourceDirectory: packages
```

## Stack outputs

| Output | Description |
|--------|-------------|
| `function_id` | App Platform app UUID that hosts the functions component. Used to import `digitalocean_app`. |
| `https_endpoint` | Public HTTPS URL of the functions HTTP endpoint. |
| `default_hostname` | Default `ondigitalocean.app` hostname. |

## Infrastructure as Code

- [Pulumi module](./iac/pulumi/README.md)
- [Terraform module](./iac/tf/README.md)

How `project.yml` works, and why runtime is not on the spec, is in [GUIDE.md](./GUIDE.md). Field-by-field schema is in [v1alpha1/reference.md](./v1alpha1/reference.md).

---

© Planton. Licensed under [Apache-2.0](https://github.com/plantonhq/planton/blob/main/LICENSE).
