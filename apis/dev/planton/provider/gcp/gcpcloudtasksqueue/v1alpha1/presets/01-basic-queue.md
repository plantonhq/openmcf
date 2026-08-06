# Basic Queue

The minimal task queue: a named queue in the ambient project with
GCP-managed defaults for dispatch rate and retries.

## What this preset creates

A Cloud Tasks queue named `background-jobs` in `us-central1`. Producers
enqueue tasks against it (each task carrying its own HTTP target), and
Cloud Tasks dispatches them at the GCP default rate with the GCP default
retry policy.

## When to use

- Simple background job processing where the defaults are enough
- Development and testing environments
- Low-volume task dispatch

## Remix ideas

- Add `rateLimits` and `retryConfig` to protect a downstream service
  (see the rate-limited-processing preset).
- Add `httpTarget` with an OIDC token to centralize auth and routing at
  the queue level (see the secure-cloud-run-target preset).
- Add `stackdriverLoggingConfig` with a `samplingRatio` for operational
  visibility into dispatch activity.
