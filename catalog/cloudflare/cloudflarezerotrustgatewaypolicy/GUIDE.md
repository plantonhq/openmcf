# CloudflareZeroTrustGatewayPolicy guide

The judgment this guide protects you from: a policy without `enabled: true` deploys and filters nothing, and the current Terraform provider re-plans every policy forever even when you did everything right.

## Known upstream blocker: the plan never converges

At provider v5.23.0 (and v5.24.0 -- the newest release), every `cloudflare_zero_trust_gateway_policy` shows a pending in-place update on every refreshed plan, forever: the resource's computed attributes ship without state-preserving plan modifiers ([upstream issue #7106](https://github.com/cloudflare/terraform-provider-cloudflare/issues/7106)), and Cloudflare's Gateway API echoes only the `rule_settings` fields you sent, so the provider re-plans the rest as "(known after apply)" on every run. Measured live 2026-08-26 on every configuration shape -- `rule_settings` absent, empty, and populated. Applies succeed and the repeated update is a no-op write, so the policy WORKS; but plan-gated pipelines will see a permanent diff, and no configuration or `ignore_changes` incantation clears it (both measured). Pulumi surfaces the same drift only under `pulumi refresh`. Until a provider release fixes #7106, treat the perpetual diff as upstream noise, not configuration error.

## Enabled defaults to false

Cloudflare's Gateway policy `enabled` default is false. The spec models it as proto3 optional so unset is visible, but the API still treats a missing value as disabled. Write `enabled: true` on every policy you intend to enforce. A code review that does not see that field should bounce the change.

## Precedence is yours to manage

`precedence` is optional and has no default. Omit it and Cloudflare assigns a number; two policies created that way can evaluate in an order you did not choose. Set it explicitly. Lower runs earlier.

## Always-sent empty rule_settings

The module always emits `rule_settings` -- an empty object when the spec configures nothing. That is the provider's own workaround (its `teamsruleconfigminimal` fixture comments "Explicitly set empty rule_settings to prevent API drift"). Do not try to omit the block in a forked module; the next plan will fight you.

## add_headers and override_ips drift on first apply

At provider v5.23.0, a policy whose `rule_settings` contains `add_headers` or `override_ips` shows computed-field drift **even on the first apply**. The provider's own migration tests expect a non-empty plan for those shapes. Planton's Cloudflare harness asserts apply idempotency, so those fields stay out of live proof scenarios and belong in offline plans until the defect clears. `block_page`, `check_session`, `block_reason`, and `dns_resolvers` do not have that defect.

## Wirefilter expressions are reformatted

`traffic`, `identity`, and `device_posture` are wirefilter strings. The API rewrites them before storing (whitespace, function spelling, list syntax). If a plan shows drift on an expression you did not change, the stored form is the API's, not yours -- copy the planned value back into the spec. Fighting the formatter produces a perpetual diff.

## Entitlement-gated actions fail at apply

`isolate` / `noisolate` need the Browser Isolation add-on. `egress` needs dedicated egress IPs. The apply fails on an account that lacks the entitlement; nothing is billed or upgraded through this component. Keep those actions out of a free-plan account's live path.

## Filter is singular here

Cloudflare models `filters` as a list whose description says it can only contain a single value -- and then does not enforce the size. This spec is a singular `filter`; the module wraps it to `[filter]`. A two-element list is not expressible, which is the point: the API is the one that would misbehave, and we do not offer the footgun.

## Destroy vs disable

Destroy is a real delete. The schema carries a computed `deleted_at` the provider never reads, so a GET after delete may 404 or return a tombstone. Prefer `enabled: false` when you need a reversible off-switch -- especially for a high-precedence block that you might have to put back in a hurry.

## Pairs well with

- [CloudflareZeroTrustList](../cloudflarezerotrustlist/README.md) -- put shared domains/IPs in a list and reference the list ID from `traffic`.
- [CloudflareZeroTrustTunnelVirtualNetwork](../cloudflarezerotrusttunnelvirtualnetwork/README.md) -- `rule_settings.dns_resolvers.*.vnet_id`.
- [CloudflareRuleset](../cloudflareruleset/README.md) -- website WAF, not employee Gateway traffic. Different product, different kind.
