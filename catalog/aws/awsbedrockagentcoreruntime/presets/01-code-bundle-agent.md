# Code-Bundle Agent

This preset hosts a Python agent from an S3 code bundle on the managed
runtime — no container image to build or maintain — with public egress
and one floating `live` endpoint.

## When to Use

- The fastest path from agent source code to a managed, scalable host
- Python/Node agents whose dependencies fit the managed base runtime

## What You Get

- A serverless, session-isolated runtime billing per second only while
  sessions execute
- A `live` endpoint that tracks the latest version — every spec change
  rolls forward automatically

## Customize

- Switch to `artifact.container.imageUri` for full control of the
  execution environment (any language; replaces the runtime)
- Pin `endpoints[].agentRuntimeVersion` to hold production on a known
  version while a floating endpoint tracks the draft
- Add `customJwtAuthorizer` to admit OIDC bearer tokens alongside IAM
