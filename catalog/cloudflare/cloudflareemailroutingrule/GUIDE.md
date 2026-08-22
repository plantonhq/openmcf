# CloudflareEmailRoutingRule guide

Operational judgment for routing rules. The README covers what each field is; this covers how the pieces interact.

## Rules need an enabled zone and verified destinations

A rule only exists inside a zone whose Email Routing is ENABLED (`CloudflareEmailRoutingZone`) — the API rejects rules otherwise. Every `forward` destination must be a VERIFIED `CloudflareEmailRoutingAddress`: Cloudflare rejects forwarding to unverified mailboxes, and verification happens only through the emailed link, out-of-band of any deploy.

## Actions are a list, and order matters

One rule can forward to mailboxes AND hand the message to an Email Worker — the common "human inbox + automation" shape (see the multi-action preset). `drop` stands alone by nature: combining it with delivery actions is contradictory even though the API does not statically reject it.

## Matchers: literal vs all

- A `literal` matcher targets exactly one recipient (`field: to` is the only supported field). A message matches the rule if ANY matcher matches.
- An `all` matcher makes the rule catch-all-shaped — but prefer the zone's real catch-all (`CloudflareEmailRoutingZone.catchAll`) for fallback policy: it is evaluated after every rule by design, while an `all` rule competes on priority.

## Priority discipline

Lower runs first; `0` is fine for a single rule. With several rules, space priorities (10, 20, 30…) so later inserts need no renumbering, and keep specific-recipient rules at lower numbers than any broad rule.
