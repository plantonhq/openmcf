# CloudflareZeroTrustList guide

The judgment this guide protects you from: the list type is a one-way door, and this is not the same object as `CloudflareList`. Pick the type deliberately, and do not mix the two list families.

## Type is immutable: changing it replaces the list

`type` RequiresReplace. A DOMAIN list cannot become an IP list in place -- Cloudflare mints a new list with a new ID. Every Gateway policy or posture rule that referenced the old `list_id` now points at a deleted object. The discipline is the same as identity-provider type: create the new list, retarget the policies, then destroy the old one.

Use the uppercase form (`DOMAIN`, `IP`, `EMAIL`). The API stores type uppercase; lowercase input would round-trip as permanent plan drift.

## Items are a set

Cloudflare treats `items` as a set. Order is not significant and is not preserved. Two applies that only shuffle item order are the same list. Do not depend on position.

This spec requires `value` on every item. The provider leaves `value` optional; an entry with only a description matches nothing and can only be a mistake, so validation rejects it.

## URL-type lists drift forever at v5.23.0

The API normalizes URL values (trailing slashes, scheme casing, and kin). The provider does no normalization and its own URL-type acceptance test expects a non-empty plan. A URL list under an idempotency-gated live lane will false-fail. Prefer DOMAIN or IP lists for managed configuration; if you need URL, write already-normalized values and expect a perpetual diff until the provider catches up.

## This is not CloudflareList

`CloudflareList` / `CloudflareListItem` are the older account-level lists consumed by Rulesets (`rules/lists/`). This kind is the Zero Trust / Gateway list (`gateway/lists/`). They do not share IDs, APIs, or import formats. A Gateway policy cannot reference a Ruleset list, and a Ruleset cannot reference a Zero Trust list. If you are reaching for the wrong one, the policy's error at apply is the tell.

## Destroy breaks referencers

Destroy is a real delete. Policies that referenced `list_id` start failing their list lookup. Update or delete those policies first. Emptying `items` and keeping the list is the reversible alternative when you want the ID to stay stable.

## Pairs well with

- [CloudflareZeroTrustGatewayPolicy](../cloudflarezerotrustgatewaypolicy/README.md) -- the policy that matches this list from `traffic` or `identity`.
- [CloudflareList](../cloudflarelist/README.md) -- the other list family; read the boundary above before you pick.
