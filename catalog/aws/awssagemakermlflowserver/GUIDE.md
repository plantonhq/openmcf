# AwsSagemakerMlflowServer — Component Guide

Authored operational judgment for the MLflow tracking server component:
the design decisions behind the spec's shape, and what to know before
running managed MLflow in production.

## Design decisions

- **The server and the app are two AWS products, not one.** This kind
  is the CLASSIC tracking server: dedicated capacity, sized and billed
  hourly, ~25-minute lifecycle operations. `AwsSagemakerMlflowApp` is
  the SERVERLESS MLflow 3.x successor — billed per use, standalone,
  and NOT a satellite of any tracking server. They share nothing but
  the MLflow protocol; modeling them as one kind would bury that
  split.
- **`mlflow_version` forbids patch-level pins.** AWS normalizes the
  version to `major.minor`, so a `3.0.1` pin would drift forever — the
  spec's pattern rejects it at manifest time. Omitted means AWS picks
  the latest; changing the pin replaces the server.
- **`automatic_model_registration` carries a taught trap.** The
  provider cannot turn it back OFF — a true-to-false change is
  silently not transmitted (an upstream update-guard gap). The spec
  field says so, and the modules always render the value so the intent
  stays visible in the plan.
- **The server's AWS name derives from `metadata.name`** — no rename
  field to drift.
- **`role_arn` defaults to an `AwsIamRole` reference**, and changing
  it replaces the server (provider-enforced) — as does the version
  pin.

## Running a tracking server in production

- **Budget ~25 minutes per lifecycle operation.** Creation and
  deletion each take about that long (provider timeouts are 45m per
  operation, not user-configurable) — plan replacements as
  half-hour-plus events, and remember role and version changes ARE
  replacements.
- **The meter runs from Created onward.** Small is ~$0.6/hour
  (~$430/month) whether anyone logs a run or not. If the team tracks
  intermittently, the serverless `AwsSagemakerMlflowApp` at $0 idle is
  the better fit.
- **Start Small; resize in place.** Size upgrades are a
  maintenance-window style operation, not a replacement — there is no
  penalty for starting at ~25 users and growing.
- **Treat auto-registration as one-way.** Enable
  `automatic_model_registration` only once you are sure — turning it
  off means replacing the server (~50 minutes of lifecycle) or an
  out-of-band API call.
- **Set the maintenance window deliberately.** Omitted means AWS
  picks; `weekly_maintenance_window_start` (UTC `DDD:HH:MM`) puts
  resizes and patching in your quiet hours.

---

© Planton. Licensed under [Apache-2.0](https://github.com/plantonhq/planton/blob/main/LICENSE).
