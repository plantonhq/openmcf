# The deploy GitHub Action learns the offline mode — one action, two lanes, a few-line switch between them

## What changed

- **`plantonhq/planton/actions/deploy` now serves BOTH deploy postures, with
  the mode inferred from the inputs.** `org` + `audience` present is the
  connected (federated) lane exactly as it shipped — exchange the workflow's
  own OIDC token for a short-lived Planton credential, optionally register,
  deploy, follow honestly. NEITHER present is the new offline lane: the
  action installs the checksum-verified CLI and runs
  `planton service deploy --env <environment> --auto-approve` against the
  checked-out working tree — the repository's own kustomize declaration
  renders locally, the preflight report verifies everything verifiable
  before the first IaC handoff, and the resources deploy to the user's own
  cloud through the open-source engine in the order their references
  declare. No Planton backend anywhere in the offline lane, and no stored
  Planton secrets in either.

- **Half-states refuse before any network call, all problems at once.** A
  new first step decides the mode and enforces each mode's real contract:
  `org` without `audience` (or the reverse) names the exact line to add or
  remove; `service`, `image`, `register: true`, or `follow: false` in
  offline mode each refuse with the offline alternative named (`image` →
  the node-addressed `set` override); a connected deploy missing `service`
  or `image` refuses with the line to add. GitHub does not enforce an action
  input's `required:` at runtime, so this gate is the first REAL input
  enforcement the action has had — in both modes. Detection is defined over
  detectable states only: a composite action cannot distinguish an input's
  default from an explicitly-passed identical value, so only empty-default
  inputs and non-default values count as mode evidence (the reasoning lives
  as a comment law in `action.yml`).

- **New offline inputs**: `set` — node-addressed field overrides, one per
  line, `<kind>/<name>:<fieldPath>=<value>`, the form a multi-document tree
  needs because a bare field path is ambiguous across documents; the classic
  use is injecting the image the job just built. `working-directory` — where
  the service's kustomize declaration lives (default: the repository root).
  State backend configuration deliberately rides the manifests' own
  annotations or the job's `PLANTON_BACKEND_*` env — the action proxies no
  engine knobs, and the preflight report states where state lives.

- **The offline lane forces itself explicitly**: the deploy step exports
  `PLANTON_BACKEND=none` — the CLI's own documented single-invocation
  override — so a self-hosted runner carrying a leftover profile can never
  silently deploy through a backend nobody asked for.

- **The log tells the truth at a glance**: the offline deploy streams inside
  a GitHub log group with the one-line verdict OUTSIDE it, so a failed run's
  summary shows the sentence that matters and the full preflight report and
  engine output are one click away. The CLI's exit codes surface as distinct
  sentences: refused at preflight (nothing ran), a resource failed to deploy
  (completed resources are safe in state; re-running continues as no-ops),
  or success.

- **Cloud credentials in offline mode are deliberately not this action's
  business**: the README shows the providers' own official OIDC actions
  (`aws-actions/configure-aws-credentials`, `google-github-actions/auth`,
  `azure/login`) running before it — trust inherited, never re-earned — and
  the preflight verifies the ambient result honestly.

- **The README leads with the switch**: the mode-upgrade table is the first
  table — offline → connected is adding `org`/`audience`/`service`/`image`
  and dropping `set` + the state env; connected → offline is the reverse.
  The offline quick start leads (it is the zero-setup path) and teaches
  remote state from the first example, because CI runners are ephemeral and
  state must outlive them.

- **The service skill's external-CI reference** gains the one sentence that
  keeps it complete: the same action also runs fully backendless when
  `org`/`audience` are absent.

## Why

Deploying from GitHub Actions to your own cloud has no standard: the world
hand-rolls per-cloud actions or authors Terraform, and the deploy half of CI
is where GitHub's excellence stops. The engine already deploys manifest sets
offline; this closes the last visible gap — a supported, documented,
one-step GitHub surface — while keeping the connected lane byte-for-byte in
behavior and making the switch between the two lanes a few lines in one
`with:` block instead of a different action.

## The gate, on real inputs

Half-state (`org` without `audience`):

```
::error::'org' is set but 'audience' is not — a connected deploy needs both. Add:  audience: <your binding's audience>  — or remove:  org: acme  to deploy offline
```

Offline with connected-only inputs piled up (all problems at once, exit 1):

```
::error::offline mode deploys the checked-out working tree, so 'service' selects nothing — remove:  service: storefront  (or add 'org' and 'audience' to deploy through your Planton backend)
::error::'register' applies the service manifest to a Planton backend, and offline mode has none — remove:  register: true
::error::an offline deploy runs inside this job and is always followed — remove:  follow: false
::error::both modes need the environment to deploy into — add:  environment: <env-slug>
```

The offline deploy step's verdicts, one per exit code, always outside the
closed log group:

```
offline deploy succeeded — every resource in the tree is live in 'prod'
::error::refused at preflight — nothing was handed to an IaC engine. The report in the log group above names every problem and its fix.
::error::a resource failed to deploy — the engine's own error is in the log group above. Completed resources are safe in state: re-running this job re-applies them as no-ops and continues.
```

## Verification

- YAML parses; `bash -n` clean on all six step scripts.
- The mode gate exercised across six journeys against the extracted step
  script: offline happy (mode=offline), connected happy (mode=federated),
  half-state each direction, offline-with-image naming the `set` road, the
  four-problem pileup above, and connected-missing-image — every refusal
  sentence captured above is the script's real output.
- The offline step exercised with a stub CLI at exit 0/1/2: the multiline
  `set` input builds repeated `--set` flags (blank lines skipped, empty
  input adds nothing — bash-3.2-safe expansion), the log group closes on
  every path, the verdict lands outside it, and the step preserves the
  CLI's honest exit code.
- `pkg/skills` suite green after the external-CI reference edit (the same
  validator the definitions release lane packages with).
- The live GitHub Actions round-trip (a real workflow through a real runner
  in both modes) is deliberately deferred to end-to-end verification.
