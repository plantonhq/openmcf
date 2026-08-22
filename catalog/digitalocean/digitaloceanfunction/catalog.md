# Function on DigitalOcean

Deploys serverless functions as an App Platform app with a single functions component. There is no standalone Functions resource; both engines create `digitalocean_app`. Runtime, memory, timeout, and schedules live in the repo's `project.yml`, not on this spec.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **App Platform Application** -- a `digitalocean_app` whose name is the Cloud Resource metadata name
- **Functions component** -- one functions component named `spec.functionName`, sourced from the Git remote you declared
- **Environment variables** -- from `spec.envs`; `secret` values are stored in App Platform's secret store

App Platform reads `project.yml` from `sourceDirectory` to set runtime, memory, timeout, entrypoint, and any cron triggers.

## Before You Deploy

### Planton Setup

- **DigitalOcean Provider Connection** -- an active connection in the Connect module with a DigitalOcean API token.

### DigitalOcean Account

- **A Git repository** in DigitalOcean Functions layout: a `project.yml` and a packages tree. Provide a public clone URL (`git`), or a linked GitHub/GitLab/Bitbucket repo.
- **A source directory** that contains `project.yml` (for DigitalOcean's hello-world sample, `packages`).
- **A supported App Platform region** (for example `nyc3`).

## Deploy

### Console

Open the deployment store, find **Function on DigitalOcean**, and click **Deploy**. Start from the **Hello-World Function** preset to deploy DigitalOcean's public Node.js sample.

### CLI

```yaml
apiVersion: digital-ocean.planton.dev/v1alpha1
kind: DigitalOceanFunction
metadata:
  name: hello
  org: acme-corp
  env: prod
spec:
  functionName: hello
  region: nyc3
  git:
    repoCloneUrl: https://github.com/digitalocean/sample-functions-nodejs-helloworld.git
    branch: master
  sourceDirectory: packages
```

```shell
planton apply -f do-function.yaml
```

This clones the public hello-world sample and deploys it as an HTTP function. No GitHub connection is required.

## Key Configuration

**Function name** -- `functionName` is the component name inside the app (max 32 characters). The App Platform app name is the Cloud Resource `metadata.name`.

**Source** -- exactly one of `git`, `github`, `gitlab`, or `bitbucket`. Use `git` with a public clone URL when the account has no VCS connection.

**Source directory** -- required. App Platform looks here for `project.yml`.

**Runtime, memory, timeout, schedules** -- edit `project.yml` in the repo. They are not fields on this spec.

## Outputs and Dependencies

### What This Component Consumes

This component has no foreign key dependencies.

### What This Component Provides

| Output | Description | Common Downstream Use |
|--------|-------------|----------------------|
| `function_id` | App Platform app UUID hosting the functions component | Import, API operations |
| `https_endpoint` | Public HTTPS URL for invoking the function | Webhooks, API clients |
| `default_hostname` | Default `ondigitalocean.app` hostname | DNS, health checks |

## Common Patterns

**Hello-world from public git** -- DigitalOcean's Node.js sample, no GitHub connection. Start from the **Hello-World Function** preset.

**Linked GitHub with deploy-on-push** -- switch the source to `github` with `repo: owner/repo` and `deployOnPush: true` after connecting GitHub in the DigitalOcean control panel.

## Works With

Functions that should share an app with HTTP services or workers belong on [**App on DigitalOcean App Platform**](/cloud-catalog/digital-ocean-app) as `spec.functions`, not on this kind.
