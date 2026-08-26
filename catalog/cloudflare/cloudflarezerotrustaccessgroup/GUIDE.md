# CloudflareZeroTrustAccessGroup guide

Operational judgment for Access groups. The README covers what each field is; this covers how the pieces interact.

## A group is factored-out membership, nothing more

Groups grant nothing by themselves — a group only matters when a policy's rules reference it. Factor criteria into a group when two or more policies would repeat them (the engineering team, the corporate email domains, the office IP ranges); keep one-off criteria inline in the policy. The payoff is that membership evolves in one place and every referencing policy follows.

## Include, exclude, require: OR, NOT, AND

A user matches if they satisfy ANY include rule, is thrown out if they match ANY exclude rule (exclude always wins), and must satisfy EVERY require rule. The classic shape: include the broad population (email domain), exclude the leavers, require the country or device posture. If you find yourself with many require rules, you probably want two groups composed by a policy instead.

## Account scope unless you have a reason

An account-scoped group is reusable by every application in the account; a zone-scoped group only serves that zone's applications. Zone scope exists for delegated-administration setups where a zone team owns its own access criteria. Default to account scope.

## The default-group flag affects apps you haven't created yet

`is_default: true` makes Cloudflare include this group in every NEWLY created application's membership. That is a policy decision about the future, not the present — set it only on the group that genuinely represents "everyone who should reach anything by default."

## IdP-backed rules need the IdP first

Rules like `azure_ad`, `okta`, `gsuite`, and `saml` reference an identity-provider integration by ID — the IdP must already be configured in Zero Trust or evaluation fails closed. Self-contained rules (email, email_domain, ip, geo, certificate) have no such dependency and are the right building blocks early in an account's life.

## Adopting an existing group: the first apply re-asserts `is_default`

Cloudflare's read API omits `is_default`, so an import cannot restore it — the first plan after adopting an existing group shows an in-place update on exactly that attribute (measured live at provider v5.23.0). The apply is a no-op against Cloudflare's real state: it simply writes your declared value into local state. Expected, harmless, once.
