# AwsSesAccountSettings — Pulumi module (Go)

Manages the region's SES account object — the suppression list
(`sesv2.AccountSuppressionAttributes`) and VDM posture
(`sesv2.AccountVdmAttributes`) — a settings singleton.

Module facts worth knowing before editing:

- **Arms render only when present.** An omitted spec arm leaves the
  account's current setting untouched; that omission is deliberate
  configuration, not a gap.
- **An empty suppression list is a real posture** (auto-suppression
  explicitly off).
- **Destroy semantics differ per arm.** Suppression PERSISTS after
  destroy (upstream delete is a no-op — the last-applied reasons stay
  in effect); VDM is reset to DISABLED. The E2E profile's
  verify-absent asserts exactly this asymmetry.
- **VDM sub-toggles are presence-typed** — unset renders nothing, set
  maps to ENABLED/DISABLED FeatureStatus strings.
- **No tags.** Neither upstream resource carries a tags argument.

Outputs mirror the Terraform module key-for-key: `account_id`.
