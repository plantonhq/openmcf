# Custom Cipher Allowlist

Hand-pick the exact cipher suites the load balancer may negotiate — for security reviews that specify an allowlist rather than a named profile.

## When to Use

- A security team hands you an explicit cipher allowlist
- Interop with a client population whose supported suites are known exactly
- Removing one specific suite that RESTRICTED still allows

## Key Configuration Choices

- **`profile: CUSTOM`** — required for (and only valid with) `customFeatures`; the two are validated together before deploy
- **`customFeatures`** — IANA-style suite names from GCP's supported set; unknown names are rejected at deploy time
- **TLS 1.3 suites are not listable** — GCP always enables them regardless of this list; the allowlist governs TLS 1.2 and below

## Placeholders to Replace

| Placeholder | Description | Where to Find |
|---|---|---|
| `<gcp-project-id>` | GCP project ID where the policy will live | GCP Console or `GcpProject` outputs |
| `customFeatures` entries | Your approved cipher suites | GCP's SSL policy feature documentation / your security review |

## Remix Notes

- The suites shown are the four TLS 1.2 AEAD suites with forward secrecy — a sensible starting allowlist
- `customFeatures` updates in place: adding or removing a suite is a one-field change applied fleet-wide
- Compare the `enabled_features` output against your review to prove the allowlist took effect

## Related Presets

- **01-modern-tls12** — Predefined profile when an allowlist is not required
- **02-restricted-strict** — Strictest predefined profile
