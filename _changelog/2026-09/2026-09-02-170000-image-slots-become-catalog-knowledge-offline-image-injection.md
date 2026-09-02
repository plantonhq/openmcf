# Image slots become catalog knowledge — offline `--image` injection, the Action's `image:` in both modes, and the deploy-from-GitHub-Actions docs

## What changed

- **The image slot is now declared where the kind lives.** Two new field
  options in `shared/options/options.proto`: `artifact_image_slot` (on a
  spec's container-carrying field, its value naming the dotted subpath from
  each element to the image) and `artifact_version_slot` (the field
  deployment pipelines stamp from the git branch). The three artifact
  receiver kinds carry them: `GcpCloudRun` and `AwsEcsTaskDefinition`
  (`containers` → `image`, blank-fill semantics) and
  `KubernetesDeployment` (`container` → `app.image`, the repo+tag split,
  plus `version` as the version slot). The injection semantics derive from
  the annotated field's SHAPE — repeated is blank-fill (authored sidecar
  images are never touched; leaving the image blank IS the authoring
  contract), singular is unconditional, a repo+tag message receives the
  reference split (tag and digest grammars both parse).
- **New package `pkg/artifactslot`** — the one authored implementation of
  those semantics, walking the annotations from descriptors: parse
  (`host/path:tag`, `host/path@digest`, registry ports handled), inject,
  version sanitization (branch → the version grammar), structured
  injection reporting for the caller's own rendering. Never prints. The
  conformance pin asserts the receiver kinds carry exactly the expected
  slots, so an annotation regression fails here; hosted control planes
  keep their authored injectors and hold them to the same annotations
  with their own agreement test.
- **The Deploy Action's `image:` input now works in BOTH modes.** Offline
  previously refused it pointing at the node-addressed `set:` road; with
  slot injection real, the offline lane passes `image:` through as
  `--image` — the flagship input means the same thing in both modes, and
  the mode-switch table simplifies to "image stays exactly as it is" in
  both directions. `set:` remains the field-exact escape hatch. On an
  older installed CLI without the offline arm, the CLI's own refusal
  surfaces verbatim inside the log group with the verdict outside —
  honest degradation, captured live.
- **The docs gain the missing page**: `Deploy from GitHub Actions`
  (`site/public/docs/ci-cd/deploy-from-github-actions.md`) — the
  ten-minute offline path first (blank image slot, provider OIDC action,
  remote state, one workflow step), the exit-code truths, the one-table
  switch to connected. `pipelines.md`'s GitHub section points at it.
  `cli.md` and `open-source.md` were truth-checked and needed nothing.
- **The offline-deploy skill reference re-trued in the same change**:
  `--image` teaches the slot injection (and the honest no-slot refusal),
  the workflow step teaches `image:` instead of a `set:` line, and the
  local verification road teaches `--preflight-only` (exit 0 wall-clean /
  2 refused — pass the same `--image`/`--set` CI will use, so the wall
  verifies exactly what would deploy).

## Why

The image slot was the last piece of deploy knowledge that lived ONLY in
the hosted control plane's code — so offline deploys made users spell a
field path for the one flag every CI job uses. Declaring the slot in the
kind's own proto puts the knowledge where every consumer already reads
(the same single-authored-home move `default_kind` made for references),
and deriving the semantics from the field's shape means one implementation
serves every current and future receiver kind.

## Verification

- `pkg/artifactslot` suite green: blank-fill with an authored sidecar
  untouched, the all-authored no-op, the ECS fill, the Kubernetes split +
  sanitized version stamp, the digest grammar, the no-slot honest no-op,
  reference parsing (ports, digests, bare repos), version sanitization
  (case folding, illegal runs, edge hyphens, the 30-char truncation), and
  the receiver-kind conformance pin.
- Proto regen via the chunked lane (the full `make protos` hit the
  documented lint-plugin download hang; `buf format` + `buf lint` +
  `buf generate --path` over the four touched directories completed in
  seconds) with the mechanical coverage check: all four edited proto
  directories show regenerated `.pb.go` files. Gazelle regenerated the
  new package's BUILD file.
- The Action's scripts: `bash -n` clean; the gate journey "offline with
  image" now passes (mode=offline); the offline step's stub-CLI run shows
  `--image` and `--set` passed through with the group closing and the
  verdict outside; the older-CLI degradation captured against the real
  installed binary.
- `go test ./pkg/skills/...` green after the skill-reference edits; the
  site's lint and full production build run against the new docs page.
- The platform CLI's offline `--image`/`--preflight-only` wiring consumes
  this release through its pin, with its own tests and captures there.
