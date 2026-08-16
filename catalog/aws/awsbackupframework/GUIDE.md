# AwsBackupFramework — Component Guide

Authored operational judgment for the Backup Audit Manager framework
component: the design decisions behind the spec's shape, and what to
know before operating frameworks in production.

## Design decisions

- **Its own kind, not a plan arm.** A framework has zero schema edges
  to backup plans — it audits the ACCOUNT's backup posture and deploys
  with zero plans. Folding it into AwsBackupPlan would have made one
  plan's lifecycle own an account-wide audit surface.
- **The AWS name is an explicit spec field** (`framework_name`): AWS
  forbids hyphens in framework names (letter first, then letters,
  digits, underscores), which is stricter than metadata.name
  conventions. The CEL rule rejects invalid names at validate time
  instead of at AWS.
- **Control names are for_each keys** on both engines — the spec
  forbids duplicates so entries never silently collapse.
- **The scope tag map is capped at ONE pair by CEL** — the provider
  documents the single-pair limit and AWS rejects more.

## Operating frameworks in production

- **The Config dependency is real and silent.** Framework deployment
  needs an ACTIVE Config recorder recording the backup types; without
  one, deployment lands `FAILED` — and the provider treats FAILED as a
  completed apply (its waiter accepts the state). Check
  `deployment_status` in the console after first deploy, not just the
  apply result.
- **One framework per compliance regime.** Controls evaluate
  account-wide by default; use per-control scopes to narrow to tagged
  or typed subsets rather than multiplying frameworks.
- **Control evaluations materialize as Config rules** (named after the
  framework) and bill as Config rule evaluations — the framework's
  cost lever is how many controls × how many resources they evaluate.
- **Deployment status cycles on every edit** (`UPDATE_IN_PROGRESS` for
  up to a few minutes); concurrent edits are rejected with conflicts
  the provider retries through.

---

© Planton. Licensed under [Apache-2.0](https://github.com/plantonhq/planton/blob/main/LICENSE).
