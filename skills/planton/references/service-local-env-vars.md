# Running a Service Locally — The Derived Environment

Planton derives a developer's local environment variables from the service's own deploy declaration — the same manifests the deploy lane applies, resolved through the same server seams — so local configuration can never drift from deployed configuration. Read this when someone wants to run their service on their laptop with real config, generate a `.env`, keep local config fresh, or asks why a value resolved the way it did.

## The one idea to hold

There is no second config inventory. Vercel-style platforms make teams maintain env vars in a dashboard beside the deployment; Planton reads the deployment declaration itself and answers "what environment does this service need" from it. When someone asks where to ADD an env var for local dev, the answer is almost always: add it to the deploy declaration (the manifest), because that is the one source — local picks it up on the next resolve.

## The three verbs

- **`planton service env run --env <env> -- <their command>`** — the flagship. Runs THEIR start command (`npm run dev`, `go run .`, whatever) with the resolved environment injected into that process. Secrets live only in that process's memory — no file is written. This is the mode to recommend first.
- **`planton service env pull`** — writes `.env` (and `.dev.vars` when a `wrangler.toml` marks a Cloudflare Worker project) for frameworks that read files natively. Output is sorted and regeneration-stable, written owner-readable-only, and REFUSED when the target file is not gitignored — the refusal names the fix (add it to `.gitignore`, or `--force`). `planton service dot-env` is the same command under its old name.
- **`planton service env check`** — regenerates in memory and compares against the existing `.env`, naming drifted KEYS only (never values), exit non-zero on drift. This is the CI-enforceable freshness gate.

## Where the environment comes from

The WORKING TREE first — uncommitted changes count:

- An inline-declared service: the `service.yaml` in the current directory, the named `--env` entry.
- A kustomize-sourced service: the `_kustomize` tree — `--env <slug>` renders `overlays/<slug>`, `--flavor <name>` renders `dev/<name>`.
- A console-managed service with no local checkout: `--service <slug> --env <env>` reads the record's inline declaration.

**Dev flavors** are the local-development convention for kustomize trees: `_kustomize/dev/<flavor>/` stacks on a real overlay and patches only the laptop deltas (external endpoints instead of cluster-internal ones, an isolation suffix so a laptop process never steals deployed work). Flavors never deploy — anything under `overlays/` IS a deploy environment, which is exactly why flavors live outside it. Someone typing `--env local` from old muscle memory gets a refusal that names the tree's overlays and its flavors.

## Layering, and explaining a value

On top of the derived contract sit exactly two personal layers, later wins:

1. **`.env.local`** in the project directory — the developer's own values, gitignored. This is also the honest remedy when their role cannot read a secret: set a local value there and the failure is answered.
2. **`--set NAME=value`** — explicit one-shot overrides.

`--explain` shows every entry's provenance — which layer won and which reference it resolved from — with secret VALUES masked. Use it to answer "why is this value what it is".

## Reading failures accurately

Resolution failures are per-key and partial results are usable: a developer without permission to read one secret still gets everything else, with a table naming each unresolved reference and its reason. Two things to relay carefully:

- **"does not exist or you are not permitted"** is ONE deliberate wording — the platform never reveals which. Don't speculate; the working paths are asking an admin for read access, or setting a local value in `.env.local`.
- By default nothing is materialized while references stay unresolved (`--allow-missing` proceeds without them). A missing key at runtime is harder to debug than a refusal at generation, so the refusal is the kindness.

ECS task definitions add one honest gap: their provider-managed secrets (AWS Secrets Manager / SSM ARNs in the containers' `secrets` map) are resolved by the ECS agent at task start — no platform seam can fetch them locally. They surface as "not locally resolvable" notes; the local answer is `.env.local`.

## The agent tool

`get_service_env_contract` returns the CONTRACT for one environment — every declared env var and secret with its `$var`/`$secret` reference string and the manifest that declares it — NEVER resolved values, so no secret material passes through an agent conversation. Its answer names its source: the service record's inline declaration, or (for kustomize services) the latest deployment receipt's applied manifests. A hosted agent cannot see a working tree, so uncommitted changes are visible only to the local CLI — when freshness against the checkout matters, direct the person to `planton service env` locally.
