# GcpCloudArmorPolicy — Pulumi Implementation

This directory contains the Pulumi implementation for provisioning a GCP
Cloud Armor security policy from the Planton spec. It also enables the
Compute Engine API on the target project.

## File Organization

| File | Purpose |
|------|---------|
| `main.go` | Module entry point; invokes `Resources()` which wires locals, provider, API enablement, and `security_policy` |
| `module/locals.go` | Ambient project fallback, policy name (spec or metadata.name fallback) |
| `module/security_policy.go` | Maps spec to `gcp.compute.SecurityPolicy`; contains rule mapping logic |
| `module/outputs.go` | Output key constants (`policy_id`, `policy_name`, `policy_self_link`, `fingerprint`) |

## Rule Mapping

The spec uses flattened structures; the Pulumi SDK expects nested types. The mapping logic in `security_policy.go` handles:

- **Match**: `versioned_expr` + `src_ip_ranges` → `SecurityPolicyRuleMatchArgs.Config`; `expression` → `SecurityPolicyRuleMatchArgs.Expr`; `expr_options` → the nested reCAPTCHA site-key options
- **Redirect**: Rule-level `redirect_options` and rate-limit `exceed_redirect_options` share the same spec type (`GcpCloudArmorRedirectConfig`) but map to separate SDK types: `SecurityPolicyRuleRedirectOptionsArgs` vs `SecurityPolicyRuleRateLimitOptionsExceedRedirectOptionsArgs`
- **Composite rate-limit keys**: `enforce_on_key_configs` → `SecurityPolicyRuleRateLimitOptionsEnforceOnKeyConfigArgs` (mutually exclusive with the singular `enforce_on_key`, enforced by the spec)
- **Adaptive protection**: `threshold_configs` with per-granularity traffic units map onto the deeply nested `Layer7DdosDefenseConfigThresholdConfig` SDK types
- **WAF exclusions**: The SDK generates distinct types for each field (`RequestHeader`, `RequestCooky`, `RequestUri`, `RequestQueryParam`). The module uses per-field mappers (`mapWafExclusionHeaders`, `mapWafExclusionCookies`, etc.)

## Default Rule Contract

Every Cloud Armor policy carries a default rule at priority 2147483647.
Creating with NO rules lets the API add a default "allow all" rule
automatically; providing ANY rules requires the set to include that default
explicitly — the spec enforces this before the module ever runs.

## Deliberately Unmodeled (parity with Terraform)

- **Labels** — `google_compute_security_policy` has no labels attribute on
  the released `google ~> 6.0` line the Terraform module pins. The bridged
  Pulumi provider would accept them, but that would be intent only one
  engine honors — so neither engine sends labels (PARITY comment in
  `security_policy.go`).
- **`request_body_inspection_size`** — absent from the released provider
  line (a newer-major-only surface); not modeled on either engine.
