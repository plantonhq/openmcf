# Push with OIDC Authentication

The authenticated webhook: Pub/Sub POSTs each message to an HTTPS
endpoint with a signed OIDC JWT the receiver can verify.

## What this preset creates

A push subscription that delivers to a Cloud Run service, with the
endpoint wired as a REFERENCE to the `GcpCloudRun` resource — its URL
contains a generated suffix that only exists after the service deploys,
so the reference is both the reliable wiring and the ordering edge (the
subscription deploys after the service). Pub/Sub mints an OIDC token as
the referenced `GcpServiceAccount` and attaches it as the Authorization
header — the receiving service verifies the token instead of trusting
the network.

## Prerequisites

- A `GcpCloudRun` service named `my-app` (or swap the reference for a
  literal HTTPS URL when pushing to an external endpoint).
- A `GcpServiceAccount` named `push-invoker` holding invoke rights on
  the receiving service (e.g. `roles/run.invoker`).
- The identity deploying this subscription needs
  `iam.serviceAccounts.actAs` on that service account.

## Remix ideas

- Add `noWrapper.writeMetadata: true` to deliver the raw payload with
  Pub/Sub metadata as HTTP headers — the shape most third-party webhook
  receivers expect.
- Add a `deadLetterPolicy` so poison messages divert instead of
  retrying forever against a failing endpoint.
