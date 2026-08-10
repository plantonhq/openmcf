# Adopt Default Bucket

Bring the project's built-in `_Default` bucket under declarative
management and raise its retention from GCP's 30-day default — the
smallest change that stops silent log expiry in its tracks.

## What it configures

- `bucketId: _Default` — ADOPTS the existing built-in bucket (creating
  a config whose name matches an existing bucket patches it rather than
  failing; this is the designed path, not a trick).
- `retentionDays: 90` — three months instead of GCP's 30 days.
- `deletionPolicy: ABANDON` — `_Default` is undeletable by API decree; a
  destroy can only un-manage it, and ABANDON says so honestly instead
  of pretending DELETE means anything here.

## Adjust before deploying

- **retentionDays** — every project's `_Default` receives whatever no
  sink explicitly routes elsewhere; longer retention here is the
  broadest (and bluntest) cost lever in Cloud Logging. Price it before
  choosing 365.

## When to choose something else

For a dedicated bucket that specific sinks route into (rather than the
catch-all), start from the **Compliance Retention** preset — leave
`_Default` short and cheap, and give the important streams their own
home.
