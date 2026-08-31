# Preview Environments: A Real Environment per Pull Request

A pull request on an opted-in service gets its own short-lived **environment** — a real Environment record, not a special kind. It is always a preview OF a base environment ("a preview of dev"): the record carries the base environment's slug, the service, the pull request number, and an expiry instant. Deployed-resource identity on Planton is environment-scoped, so a per-PR environment isolates everything the preview deploys with no renaming anywhere.

## The opt-in

```yaml
spec:
  build:
    triggers:
      pullRequests:
        deploy: true   # implies build; previews per pull request
```

A service that never sets this sees zero new records, zero new surface, zero cost. `build: true` alone builds PR commits without ever minting an environment.

## What happens on a pull request today

- **Open or push** (opened, synchronize, reopened, ready for review): the build runs, and the platform ensures the PR's preview environment — named `{service}-pr-{number}` (e.g. `storefront-pr-123`) — minting it on the first push and refreshing its expiry on every later one. Only an abandoned pull request ages toward its expiry.
- **The base environment** is where the PR's target branch deploys: a branch mapped to an environment previews that environment; a PR against a trigger branch previews the first environment of the promotion walk.
- **Configuration stays live**: config references name their environment inside the reference itself, so the base environment's variables and secrets serve the preview for its whole life — nothing is copied, and a rotated value is what the preview reads next.
- **The preview deploy itself is not yet open.** A preview run currently completes at its build terminal with an explained skip naming its preview environment. Relay that honestly — never suggest the preview URL exists yet.

## Reading a preview run's deploy state

The run record's deployment stage carries exactly one of three explanations:

| The record says | It means | The working path |
|---|---|---|
| "enable build.triggers.pullRequests.deploy…" | The service never opted in | Set the flag if the user wants previews |
| "no preview environment was born for this run…" | The flag is on, but the mint was refused — the per-service concurrency cap (previews of other open PRs), a missing base environment, or the environment name is already taken by a non-preview environment | Close a PR to free a preview slot, declare deploy environments, or rename the colliding environment |
| "Preview environment '…' is ready, but preview deploys are not yet available" | The environment record exists and is visible everywhere environments are | Nothing to fix — this is the platform's current truth |

A preview outcome never fails the BUILD: a capped-out or refused preview degrades the run to build-only, with the reason in the delivery log.

## Inspecting previews

`list_environments` and `get_environment` (CLI: the environment read commands) show previews beside durable environments; the `spec.preview` block is the tell — base environment, service, pull request number, expiry. The block is server-managed: create, update, and apply refuse it, so no manifest can disguise a durable environment as a preview or convert one. Previews never appear in promotion order — promotion walks are derived from each service's own declared environments.
