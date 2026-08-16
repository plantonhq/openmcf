# AwsDlmLifecyclePolicy

One Data Lifecycle Manager policy — account-level, tag-targeted snapshot and AMI automation (create, retain, archive, copy cross-region, share, deprecate) as a default-XOR-custom mode union. The policy references no specific volume: it acts on whatever carries the target tags when a schedule fires.

## Highlights

- **Backup automation as one declarative object**: AWS's simplified default policy ("snapshot everything daily, keep two weeks") or the full custom engine — up to four named schedules with cron/interval triggers, per-schedule retention, archive tiering, cross-region copies, FSR, sharing, and AMI deprecation.
- **The provider's mode walls are CELs**: default-policy dials conflict with schedules at the provider; the spec's arm union plus the event-based rules make every illegal combination fail at validation instead of at AWS.
- **Two derived arguments, honestly recorded**: the provider's `default_policy` marker and `policy_language` derive from the configured arm in both modules — spec fields for them could contradict the mode (recorded parity exclusions).
- **Event-based policies included**: react to snapshots shared INTO the account (`shareSnapshot`) by copying them — the cross-account backup-ingestion pattern.

## Both Engines

Both modules derive the mode arguments identically and export the same outputs: `policy_id` (import ID), `policy_arn`.

## Chart Wiring

`execution_role_arn` references an AwsIamRole (DLM's trust is `dlm.amazonaws.com`). The policy pairs naturally with AwsEbsVolume in charts: tag the volumes, target the tags — nothing else to wire.
