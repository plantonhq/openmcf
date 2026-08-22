# DigitalOcean Function -- Operational Guide

Judgment calls that matter when you run Functions on App Platform.

## There is no Functions Terraform resource

Both engines create `digitalocean_app` with one functions component. `function_id` is that app's UUID. Import uses the same `digitalocean_app` id format as DigitalOceanApp (`{app_id}`), derived from the `function_id` stack output.

## Runtime does not belong on the spec

`project.yml` inside `sourceDirectory` is the source of truth for:

- runtime (`nodejs:18`, `python:3.11`, …)
- memory and timeout
- web vs. non-web (`web: true/false`)
- cron schedules
- per-function entrypoints

App Platform reads that file at deploy time. Terraform and Pulumi cannot set those knobs on `digitalocean_app`. A spec field for them would look configurable and do nothing — so they are omitted on purpose.

To change runtime or add a schedule, edit `project.yml` and redeploy.

## `sourceDirectory` is load-bearing

It must point at the directory that contains `project.yml` and the packages tree. For DigitalOcean's `sample-functions-nodejs-helloworld` that is `packages`, not the repo root. A wrong directory produces a failed App Platform build, not a spec-validation error.

## Pick `git` unless GitHub is actually connected

`github.deployOnPush` requires the DigitalOcean account to have GitHub connected in the control panel. Without that connection the deploy fails. A public clone URL on `git` works with no extra setup and is the right default for new accounts and for E2E.

## Function name vs app name

`spec.functionName` is the component name (max 32 characters). The App Platform **app** name is `metadata.name`. They can differ. Import and the verifier key off the app UUID (`function_id`), not the component name.

## Environment variables

Use `envs[].plaintext` for ordinary values and `envs[].secret` for credentials. Secrets are stored in App Platform's secret store (`type = SECRET` on the provider). There is no separate `secretEnvironmentVariables` map.

## When this kind is the wrong shape

If the functions component should sit next to an HTTP service, a worker, or a static site, use DigitalOceanApp and put the functions component in `spec.functions`. DigitalOceanFunction is the one-component app.
