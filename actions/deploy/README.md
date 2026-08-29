# Planton Deploy Action

Deploy to Planton from GitHub Actions with **zero stored secrets**. The action exchanges the workflow's own OIDC token for a short-lived Planton credential, optionally registers the service from its manifest, deploys your container image into one environment, and fails the job unless the deployment verifiably succeeded.

## Quick start

```yaml
jobs:
  deploy:
    runs-on: ubuntu-latest
    permissions:
      id-token: write   # lets the job mint its OIDC token
      contents: read
    steps:
      - uses: actions/checkout@v4

      # ... build and push your image ...

      - uses: plantonhq/planton/actions/deploy@main
        with:
          org: acme
          service: checkout-api
          environment: prod
          image: ghcr.io/acme/checkout-api@${{ steps.build.outputs.digest }}
          audience: https://acme.example/planton   # must match your binding's audience
```

## One-time setup (keyless trust)

The keyless exchange works because your organization registered a **workload identity binding** — the trust rule "workflows from this repository may act as this service account". Someone with manage-access on the organization sets it up once:

```yaml
# workload-identity-binding.yaml — applied with: planton apply -f workload-identity-binding.yaml
apiVersion: iam.planton.ai/v1alpha1
kind: WorkloadIdentityBinding
metadata:
  org: acme
  name: checkout-api-github
spec:
  serviceAccount: ci-deployer          # an existing service account — the audit trail's name
  audience: https://acme.example/planton
  githubActions:
    repository: acme/checkout-api      # matched exactly against the token's own claim
```

There is no API key to create, rotate, or leak. Each workflow run proves itself with GitHub's own token; the exchanged credential expires in minutes and can act only on the service this repository's token names.

## Inputs

| Input | Required | Default | Description |
|---|---|---|---|
| `org` | yes | | Planton organization slug. |
| `service` | yes | | Service slug or id to deploy. |
| `environment` | yes | | Environment to deploy into (must be among the service's declared deploy environments). |
| `image` | yes | | Container image reference (`host/path:tag` or `@digest`). |
| `audience` | yes | | The OIDC token audience — exactly the value your binding declares. |
| `register` | no | `false` | Apply the service manifest before deploying. The manifest's declared repository is **proven** against this workflow's token, never trusted as typed. |
| `service-file` | no | `service.yaml` | Manifest path, used with `register: true`. |
| `follow` | no | `true` | Wait for the run and fail the job honestly (see below). Set `false` to fire-and-forget. |
| `cli-version` | no | `latest` | Planton CLI version to install. |

## Outputs

| Output | Description |
|---|---|
| `run-id` | The delivery run this deploy started (`svcpipe_...`). |

## What makes the job red

With `follow: true` (the default), green means something:

- The **run failed** (a resource failed to apply) → the job fails with the failed resources named.
- The run succeeded but **rollout verification reported `failed`** — Planton watched the workload on your own cloud and it demonstrably never came online (crash loops, unreachable endpoints) → the job fails with each failed check and the provider's own diagnostic words.
- Verification was honestly **unverifiable** (no declared endpoint, no runner to probe from) → the job passes with a note; honesty is not an alarm.

## Behavior worth knowing

- **Protected environments**: a deploy into a protected environment pauses at that environment's approval gate, and the deployer can never be the approver. With `follow: true` the job waits for the approval — bound your job with `timeout-minutes` if that wait should not be open-ended.
- **Inline deploy configuration only**: the deploy verb renders manifests from the service's inline deploy declaration. A service using the kustomize path refuses with the working alternatives named (the git-push lane deploys it fully; promote moves an existing deployment).
- **Provenance rides along**: the action passes the workflow's commit and branch, so the run header and history read honestly for an image Planton never built.
- **Runners**: Linux and macOS runners, x64 and arm64. The CLI download is checksum-verified against the release's published checksums before it is ever executed.
