# CloudflareZeroTrustAccessPolicy guide

Operational judgment for standalone Access policies. The README covers what each field is; this covers how the pieces interact.

## Standalone means reusable — and independently owned

This kind models the modern, account-level policy that applications attach BY REFERENCE. One "allow employees" policy can guard twenty applications; tightening it once tightens all twenty. The old per-application inline policy is a different provider surface this catalog deliberately skips.

## Decision is the whole ballgame

`allow` grants after rules match; `deny` blocks even when a later policy would allow (deny is evaluated with priority at the application); `non_identity` admits service tokens and other non-human traffic; `bypass` turns Access OFF for matching traffic — use it for health checks and only with rules far narrower than `everyone`. A bypass with a broad include is an open door wearing a policy's name.

## Rules compose the same way groups do

Include is OR, exclude is NOT (and wins), require is AND. Keep a policy's rules about the DECISION and put the population in a referenced Access group: "allow, include group:engineering, require geo:US" reads as intent; five inline email rules do not.

## Deleting is blocked while applications still reference it

Cloudflare refuses to delete a policy with a non-zero application count. Detach it from every application first — in a chart teardown the dependency graph does this for you, but a hand-managed policy needs the detach step remembered.

## Approval flows change the login experience

`approval_required` with approval groups turns login into a request: the user waits until an approver in the configured group says yes (per session, or once when `purpose_justification` collects a reason). Reserve it for genuinely sensitive targets — it adds a human to every session establishment.

## Adopting an existing policy: the first apply re-asserts three flags

Cloudflare's read API omits `approval_required`, `isolation_required`, and `purpose_justification_required`, so an import cannot restore them — the first plan after adopting an existing policy shows an in-place update on exactly those attributes (measured live at provider v5.23.0). The apply is a no-op against Cloudflare's real state: it writes your declared values into local state. Expected, harmless, once.
