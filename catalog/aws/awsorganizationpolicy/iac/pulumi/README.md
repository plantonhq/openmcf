# AwsOrganizationPolicy — Pulumi module

Manages one Organizations policy (`organizations.Policy`) with its
folded attachments (`organizations.PolicyAttachment`, one per spec
entry keyed by resolved target, parented to the policy).

Module facts worth knowing before editing:

- **`Name` renders `spec.policy_name`** — the explicit name field
  (policy names allow spaces `metadata.name` cannot carry).
- **`Type` is sent only on an explicit choice** (unset =
  SERVICE_CONTROL_POLICY, the provider default) and forces
  replacement.
- **`Content` is `json.Marshal(spec.Content.AsMap())`** — the spec
  carries the document structured (the catalog's uniform
  policy-document idiom).
- **Attachment targets arrive resolved** — the module reads each
  entry's `TargetId.GetValue()`; the resolved target names the
  resource and matches the Terraform module's for_each key.
- **`SkipDestroy` is deliberately not sent** on either resource —
  destroy means detach and delete (the recorded apply-behavior
  exclusion).

Outputs mirror the Terraform module key-for-key: `policy_id` (the
import ID) and `arn`.
