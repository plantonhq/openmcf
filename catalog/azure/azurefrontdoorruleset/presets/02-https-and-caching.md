# HTTPS Upgrade + Tiered Caching

This preset creates a three-rule delivery policy: permanent HTTPS
upgrade at the edge, aggressive caching for immutable static assets
(ignoring tracking parameters), and caching turned off for API paths.

## When to Use

- A web application serving both static assets and an authenticated API
  through one endpoint
- When campaign links (utm_* parameters) fragment your cache -- the
  IGNORE_SPECIFIED behavior collapses them onto one entry

## Key Configuration Choices

- **`behaviorOnMatch: STOP` on the redirect** -- a redirected request
  should not accumulate further actions; the HTTPS retry gets them
- **`OVERRIDE_ALWAYS` + 7-day duration** for fingerprinted assets --
  right when filenames change on deploy; use `HONOR_ORIGIN` if your
  backend already sends precise Cache-Control
- **`DISABLED` on `/api/`** -- an explicit cache decision (the override
  always makes one; there is no leave-it-alone value)
- **`transforms: [LOWERCASE]`** on extensions -- catches `.CSS` and
  `.Js` uploads without a second rule

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<front-door-profile-resource-name>` | The AzureFrontDoorProfile's Planton resource name | Your Front Door composition |

## Downstream Wiring

Attach to routes via `ruleSetIds`; the HTTPS-upgrade rule makes the
route-level `httpsRedirectEnabled` redundant for routes that attach this
set (both are safe together -- the rule wins by running at the edge
first).
