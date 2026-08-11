# GcpCloudArmorPolicy Guide

The judgment this guide protects: a WAF policy is production traffic
filtering — a wrong rule blocks customers, a missing one admits attacks.
Preview mode and the priority ladder are what make changes safe; use
both, every time.

## Preview before enforce, always

Every rule supports `preview: true` — matched traffic is logged but the
action is not enforced. Ship new deny/throttle rules in preview, read
the logs for a representative window, then flip preview off. This is the
single highest-value habit with Cloud Armor: false positives found in
logs cost nothing; found in production they are an outage.

## The priority ladder is the architecture

Rules evaluate lowest priority number first, first match wins, and the
default rule at 2147483647 catches everything else. Leave deliberate
gaps (1000, 2000, 3000...) so emergency rules can slot between existing
ones. Creating with NO rules lets the API add an allow-all default;
providing ANY rules requires including that default explicitly — decide
its action consciously (allow-by-default perimeter vs deny-by-default
allowlist), because it IS your security posture.

## Rate limiting needs its key thought through

`throttle` and `rate_based_ban` are only as good as `enforceOnKey`: IP
alone punishes everyone behind a corporate NAT; HTTP_HEADER/HTTP_COOKIE
keys need the header to actually exist on abusive traffic. The
conform/exceed action pair is validated to provider truth (conform is
always `allow`). Start with generous thresholds in preview and tighten
from measured traffic, not intuition.

## WAF exclusions carve, they do not disable

Preconfigured WAF rules (SQLi, XSS...) false-positive on legitimate rich
content. The `exclusions` block removes a specific cookie, header, query
param, or URI from inspection for a named rule set — a scalpel. Reaching
for "remove the rule entirely" when one search field trips it is the
wrong tool; exclude that field instead. `requestBodyInspectionSize`
bounds how deep the WAF reads bodies (8KB default, 64KB max) — larger
catches deeper payloads at higher processing cost, and bodies beyond the
limit pass uninspected either way.

## Adaptive Protection is a second pair of eyes

Layer-7 DDoS detection learns baselines and can auto-deploy mitigations
via `thresholdConfigs` (confidence, impacted-baseline, expiration).
Start with detection only (`enable` without auto-deploy thresholds) and
graduate to auto-deploy once you trust its verdicts on your traffic.

## Teardown discipline

Detaching the policy from a backend service leaves that traffic
UNFILTERED — deletion is a security event, not a cleanup. GCP refuses to
delete a policy still attached, and `PREVENT` also covers the window
after detachment. `ABANDON` keeps the policy enforcing while dropping
management.
