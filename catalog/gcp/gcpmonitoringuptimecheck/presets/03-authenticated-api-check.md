# Authenticated API Check

Probes a health endpoint that sits behind HTTP basic auth, asserting
both transport health (2xx over valid TLS) and body truth (the health
JSON actually says ok).

## What it configures

- `httpCheck.authInfo` — basic-auth credentials; the password is a
  managed secret, never plaintext at rest.
- `acceptedResponseStatusCodes: STATUS_CLASS_2XX` — explicit about what
  passes (GCP's default is also 2xx; stating it keeps the contract
  visible next to the auth block).
- A JSON-path content matcher: `$.status` must equal `ok` — a 200
  serving `{"status":"degraded"}` fails the probe.

## Adjust before deploying

- **host / path** — the health endpoint.
- **username / password** — a probe-scoped credential, not a human's;
  supply the password through the platform's secret handling.
- **jsonPath / content** — match your health endpoint's actual contract.

## When to choose something else

For endpoints that authenticate Google identities (Cloud Run with IAM),
`httpCheck.serviceAgentAuthentication` (OIDC_TOKEN) is the keyless form —
no password to rotate.
