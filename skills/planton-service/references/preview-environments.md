# Preview Environments: A Real Environment per Pull Request

A pull request on an opted-in service gets its own short-lived **environment** — a real Environment record, not a special kind. It is always a preview OF a base environment ("a preview of dev"): the record carries the base environment's slug, the service, the pull request number, and an expiry instant. Deployed-resource identity on Planton is environment-scoped, so a per-PR environment isolates everything the preview deploys with no renaming anywhere.

## The opt-in

```yaml
spec:
  build:
    triggers:
      pullRequests:
        deploy: true          # implies build; previews per pull request
        previewTtlHours: 168  # optional; how long an untouched preview lives (default 72, max 720)
```

A service that never sets this sees zero new records, zero new surface, zero cost. `build: true` alone builds PR commits without ever minting an environment. `previewTtlHours` tunes the abandoned-PR window per service — every push refreshes the clock, so only a PR nobody touches ages toward teardown; values outside 1–720 hours are refused at save.

## What happens on a pull request

- **Open or push** (opened, synchronize, reopened, ready for review): the build runs, and the platform ensures the PR's preview environment — named `{service}-pr-{number}` (e.g. `storefront-pr-123`) — minting it on the first push and refreshing its expiry on every later one.
- **The base environment** is where the PR's target branch deploys: a branch mapped to an environment previews that environment; a PR against a trigger branch previews the first environment of the promotion walk.
- **The run DEPLOYS into the preview** when the service's kustomize tree authors a `previews/<base-env>` directory (see the authoring section below): the tree renders at the PR's own commit — a PR that changes deploy configuration previews that change — the platform re-stamps the manifests onto the preview environment, fills the blank slots with per-PR values, and applies. Rollout verification then watches the workload come online and stamps the deployment's URLs, exactly as for any environment — so an agent that authored the PR can verify its own preview before asking for review.
- **The preview's address**: when the base environment declares a serving domain, the preview serves at `{preview-env-slug}.{base-domain}` (e.g. `storefront-pr-123.dev.acme.com`), filled into the same blank-hostname slots normal deploys use. Without a serving domain, the deployment record carries the target's native endpoint (a run.app URL, an ALB hostname) discovered by rollout verification.
- **Configuration stays live**: config references name their environment inside the reference itself, so the base environment's variables and secrets serve the preview for its whole life — nothing is copied, and a rotated value is what the preview reads next. The previewed service runs against the base environment's stable neighbors; only the changed service deploys into the preview.
- **Close tears it down** (merged or not): the preview's cloud resources are destroyed first, then its records — its deployment receipts and the environment record itself — while the pull request's builds stay browsable in the run history. A preview whose expiry passes with no close (a laptop closed on a Friday) is torn down by a scheduled sweep, so an abandoned pull request never leaks cloud spend.
- **A push racing the teardown builds without a preview**, with the reason on the run; the next push after the teardown completes mints a fresh preview.

## Authoring the previews tree (kustomize services)

Preview topology is the TEAM's, never the platform's: a `previews/<env>` directory in the `_kustomize` tree — the exact sibling of the local-dev flavors under `dev/<flavor>` — stacks on the base environment's overlay and patches exactly what the team's topology needs:

```
_kustomize/
├── base/
├── overlays/
│   └── dev/                  # env: dev, the base environment's overlay
└── previews/
    └── dev/
        ├── kustomization.yaml   # resources: [../../overlays/dev], patches: [service.yaml]
        └── service.yaml         # the deltas previews of dev deploy with
```

Two teams, two topologies, zero platform opinion:

- **Isolated previews**: patch the workload's `spec.namespace` to blank (and drop `createNamespace`) in the previews tree — the platform fills each preview's namespace with the preview environment's own name and creates it, so every PR is fully isolated and destroyed wholesale.
- **Shared namespace**: leave the namespace as the overlay authored it — all previews land beside the base environment's workloads, coexisting through the workload kind's own version track (the deploy pipeline stamps `spec.version` from the PR branch, so multiple tracks of one app share a namespace with disjoint workloads).

The platform's contribution stays mechanical: manifests rendered from the base overlay legitimately carry the BASE environment's name, and the platform re-stamps exactly that onto the preview (any OTHER environment name still refuses — a manifest naming a third environment is a lie); blank namespaces and blank hostnames fill with per-PR values; authored values are never touched.

Check a tree renders before pushing: `planton service env check --preview dev` (or `pull`/`run`) renders `previews/dev` locally. Per-PR fills happen only on a real preview run, so blank slots render blank locally — that is expected.

A service with the opt-in but NO previews tree still builds every PR; its deploy skips naming the authoring path. Inline-configured services (no kustomize tree) do not have preview deploys yet — the run says so plainly.

## Reading a preview run's deploy state

The run record's deployment stage carries exactly one of these explanations:

| The record says | It means | The working path |
|---|---|---|
| "enable build.triggers.pullRequests.deploy…" | The service never opted in | Set the flag if the user wants previews |
| "no preview environment was born for this run…" | The flag is on, but the mint was refused — the per-service concurrency cap (previews of other open PRs), a missing base environment, or the environment name is already taken by a non-preview environment | Close a PR to free a preview slot, declare deploy environments, or rename the colliding environment |
| "…no preview configuration for base environment…" | The preview exists but the repository authors no `previews/<base-env>` tree | Author the tree (the skip names the exact path) |
| "The previews tree for base environment '…' failed to render…" | The tree exists and kustomize refused it — carried with the kustomize error, failing only the PR's own run (a broken previews tree never fails a main-branch deploy) | Fix the tree; `planton service env check --preview <env>` reproduces locally |
| "…this service's configuration is inline…" | Preview deploys ride the kustomize previews tree; the inline arm does not exist yet | Eject to kustomize to author preview configuration |
| "Preview environment '…' no longer exists (it was torn down…)" | A close or the TTL sweep raced this run | Push again for a fresh preview |

A preview outcome never fails the BUILD: a capped-out or refused preview degrades the run to build-only, with the reason in the delivery log.

## Inspecting previews

`list_environments` and `get_environment` (CLI: the environment read commands) show previews beside durable environments; the `spec.preview` block is the tell — base environment, service, pull request number, expiry. The block is server-managed: create, update, and apply refuse it, so no manifest can disguise a durable environment as a preview or convert one. A preview whose `spec.is_delete_in_progress` is true is being torn down right now and will disappear shortly. Previews never appear in promotion order — promotion walks are derived from each service's own declared environments.

## A preview is never deleted by hand

`delete_environment` (CLI: `planton env delete`) REFUSES a preview, marked or not — the direct delete removes records only, which would orphan the preview's cloud resources. A preview dies only through the platform's own teardown: close the pull request, or let the expiry pass. If a user asks to remove a preview, close its pull request — that IS the delete button.
