# Compliance Evidence

This preset creates the auditor-facing stream: control-compliance
results from your Backup Audit Manager frameworks, delivered daily to
S3.

## When to Use

- Scheduled evidence for compliance regimes (the framework evaluates;
  this documents)
- Tracking control drift over time from the report history

## What You Get

- A daily CONTROL_COMPLIANCE_REPORT over the referenced framework,
  under the `compliance/` prefix (CSV by default)
- Framework references by component — add more `frameworkArns` entries
  as frameworks multiply (a report plan is many-to-many by design)

## Customize

- Swap to `RESOURCE_COMPLIANCE_REPORT` for per-resource findings
  instead of per-control rollups (a template change REPLACES the
  report plan)
- Add `formats: ["CSV", "JSON"]` when pipelines consume the evidence
  too
