# AWS Lambda Layer

Ship shared code once: libraries, custom runtimes, and config files your Lambda functions attach by ARN instead of bundling into every deployment package. One layer version feeds any number of functions — including other accounts' or your whole organization's.

## What Gets Managed

- The layer version: its S3 archive, runtime/architecture compatibility metadata, and license info.
- The share grants: who else (an account, everyone, or everyone-in-your-organization) may use the version.

## Before You Deploy

### Planton Setup

- An organization and environment.
- An AWS provider connection with Lambda permissions.

### AWS Prerequisites

- The layer content zip in an S3 bucket in the same region, laid out the runtimes expect (`python/`, `nodejs/node_modules/`, ...). Lambda copies it at publish, so the object only needs to exist during deploy.

## After You Deploy

- Functions attach the `layer_version_arn` output; new function deployments pick the layer content up under `/opt`.
- Shared accounts can attach the version immediately — the grant is active as soon as the apply finishes.

## Common Changes

- Publish a new version: change the archive key (or bump `source_code_hash`) — the old version is replaced in state; set `skip_destroy` to keep it available for functions still pinning it.
- Add or remove share grants (in-place per statement).
