<p align="center">
  <img src="logo.svg" alt="AWS SSM Maintenance Window" width="80"/>
</p>

# AWS SSM Maintenance Window

Manage an [AWS Systems Manager maintenance window](https://docs.aws.amazon.com/systems-manager/latest/userguide/systems-manager-maintenance.html)
— a recurring window of time plus the targets registered with it and
the tasks that execute inside it.

## What Gets Managed

- **The window** (`metadata.name` is the window name; AWS identifies
  it as `mw-...`): its cron/rate **schedule** with timezone and
  offset, **duration** and **cutoff** (cutoff < duration — a CEL
  rule), an **enabled** tri-state (unset = enabled), activation
  dates, and the unassociated-targets allowance.
- **Targets** folded as name-keyed registrations: instance or
  resource-group selectors. AWS-generated registration IDs land in
  the `target_ids` output map.
- **Tasks** folded as name-keyed registrations: **RUN_COMMAND /
  AUTOMATION / LAMBDA / STEP_FUNCTIONS** with what-to-run as a
  value-or-reference (`taskArn`), priority, rate controls (targeted
  tasks only — an AWS rule), cutoff behavior, and the type-matched
  **invocation** union (exactly one arm — CEL-enforced). Task IDs
  land in the `task_ids` output map.

See [v1alpha1/reference.md](v1alpha1/reference.md) for the full field
reference generated from the spec proto.

---

© Planton. Licensed under [Apache-2.0](https://github.com/plantonhq/planton/blob/main/LICENSE).
