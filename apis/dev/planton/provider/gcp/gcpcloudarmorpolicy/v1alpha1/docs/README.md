# GcpCloudArmorPolicy — Research Documentation

This document captures the design notes for the GcpCloudArmorPolicy component: background on Google Cloud Armor, the modeled surface, and the deliberately unmodeled edges.

---

## Google Cloud Armor Overview

Google Cloud Armor is a distributed denial-of-service (DDoS) and web application firewall (WAF) service that protects applications and APIs behind Google Cloud HTTP(S) load balancers, Cloud CDN, and backend buckets. It operates at the edge of Google’s network, filtering traffic before it reaches backends.

Cloud Armor policies consist of:

1. **Policy-level configuration** — Type, adaptive protection, advanced options
2. **Rules** — Prioritized list of match conditions and actions
3. **Default rule** — Always present at priority 2147483647; typically "allow all"

Rules are evaluated in ascending priority order (lowest number first). The first matching rule’s action is applied. If no rule matches, the default rule applies.

---

## Policy Types and Their Use Cases

| Type | Attachments | Use Case |
|------|-------------|----------|
| **CLOUD_ARMOR** | Backend services of HTTP(S) load balancers | Full WAF, rate limiting, redirects, headers; backend security |
| **CLOUD_ARMOR_EDGE** | Cloud CDN, backend buckets | IP/geo rules only; edge-level filtering |
| **CLOUD_ARMOR_INTERNAL_SERVICE** | Internal Traffic Director | Limited feature set; internal service mesh |

The policy type is immutable after creation. Choosing the wrong type requires recreating the policy.

---

## Rule Evaluation Order and Priority System

- Priority is an integer from 0 to 2147483647.
- Lower values = higher priority (evaluated first).
- Each rule must have a unique priority.
- Priority 2147483647 is reserved for the default rule.

**Evaluation flow**: For each request, Cloud Armor evaluates rules from lowest priority number to highest. The first rule whose match condition is satisfied triggers that rule’s action. If no rule matches, the default rule (2147483647) applies.

---

## Match Conditions

### IP-Based Matching

- **versioned_expr**: `SRC_IPS_V1` (only supported value)
- **src_ip_ranges**: List of CIDR ranges (max 10 per rule)
- Use `["*"]` to match all IPs

Example: Allow office and VPN ranges, deny rest.

### CEL Expression Matching

Common Expression Language (CEL) provides flexible matching on request attributes:

- **origin.region_code** — Client country/region
- **origin.ip** — Client IP
- **request.path** — URL path
- **request.headers['X-Custom']** — Header values
- **inIpRange(origin.ip, '1.2.3.0/24')** — IP in CIDR

Example: `origin.region_code == 'US'` for geo-allowlist; `request.path.matches('/api/.*')` for path-based rules.

IP-based and CEL are mutually exclusive per rule.

---

## Rate Limiting: Throttle vs Rate-Based Ban

### Throttle

When traffic exceeds the threshold, the configured `exceed_action` is applied (e.g., deny(429), redirect). Below the threshold, traffic is allowed.

### Rate-Based Ban

Two thresholds:

1. **rate_limit_threshold** — When exceeded, apply `exceed_action`
2. **ban_threshold** — When exceeded again, ban the source for `ban_duration_sec` seconds

Useful for escalating from throttle to full block.

### enforce_on_key Options

Determines how requests are grouped for counting:

| Key | Behavior |
|-----|----------|
| `ALL` | Single counter for all matched traffic |
| `IP` | Per source IP |
| `HTTP_HEADER` | Per value of a header (set `enforce_on_key_name`) |
| `XFF_IP` | Per IP from X-Forwarded-For |
| `HTTP_COOKIE` | Per cookie value |
| `HTTP_PATH` | Per URL path |
| `SNI` | Per TLS Server Name Indication |
| `REGION_CODE` | Per client country/region |

---

## Adaptive Protection (CAAP)

Cloud Armor Adaptive Protection (CAAP) uses machine learning to detect Layer 7 DDoS and anomalous traffic. When enabled:

- Traffic patterns are analyzed in real time
- Anomalies generate alerts
- Optional auto-mitigation creates adaptive rules

**rule_visibility**: `STANDARD` (default) or `PREMIUM` (requires Managed Protection Plus).

---

## Preconfigured WAF Rules (OWASP ModSecurity Core Rule Set)

Cloud Armor includes preconfigured WAF rules based on the OWASP ModSecurity Core Rule Set. Common rule set identifiers:

- **sqli-v33-stable** — SQL injection
- **xss-v33-stable** — Cross-site scripting
- **rce-v33-stable** — Remote code execution
- **lfi-v33-stable** — Local file inclusion

Without exclusions, these rules can cause false positives (e.g., SQL-like content in search, HTML in rich text, GraphQL syntax).

---

## WAF Exclusions and False Positive Handling

WAF exclusions carve out specific request components from rule evaluation:

- **target_rule_set** — Which OWASP rule set to exclude from
- **target_rule_ids** — Optional list of specific rule IDs
- **request_headers** — Exclude headers from WAF inspection
- **request_cookies** — Exclude cookies
- **request_uris** — Exclude URI paths
- **request_query_params** — Exclude query parameters

Each exclusion field uses **operator** (EQUALS, STARTS_WITH, ENDS_WITH, CONTAINS, EQUALS_ANY) and **value**.

Example: Exclude `?search=` from SQLi rules so search boxes with "SELECT" are not blocked.

---

## Advanced Options

### JSON Parsing

- **DISABLED** — No JSON body inspection
- **STANDARD** — Parse JSON for WAF rules
- **STANDARD_WITH_GRAPHQL** — Parse JSON and GraphQL

Needed when WAF rules inspect request bodies.

### Logging

- **NORMAL** — Standard logging
- **VERBOSE** — Matched rule and request details

### IP Resolution

**user_ip_request_headers** — Custom headers to use for client IP when traffic passes through a CDN or proxy (e.g., `X-Forwarded-For`, `CF-Connecting-IP`).

---

## Scope

### Modeled

The component models the released provider resource surface:

1. **Core policy** — Name, type, description; the ambient-project contract
2. **Rules** — All actions (allow, deny, redirect, throttle, rate_based_ban), preview mode
3. **Match** — IP-based and CEL expressions, with reCAPTCHA site-key options (`expr_options`) for token-evaluating expressions
4. **Rate limiting** — Thresholds, ban escalation, single `enforce_on_key` or composite `enforce_on_key_configs`, exceed redirect
5. **Redirect** — EXTERNAL_302 and GOOGLE_RECAPTCHA, with the policy-level `recaptcha_options_config` site key
6. **Header injection** — Custom headers for matching requests
7. **Preconfigured WAF** — Exclusions for all four field types (headers, cookies, URIs, query params)
8. **Adaptive Protection** — Layer 7 DDoS defense, rule visibility, and per-granularity `threshold_configs` with auto-deploy tuning
9. **Advanced options** — JSON parsing (with `json_custom_config` content types), log level, user IP headers

### The Default Rule Contract

Every Cloud Armor policy carries a default rule at priority 2147483647.
Creating a policy with NO rules lets the API add a default "allow all" rule
automatically; providing ANY rules requires the set to include that default
explicitly — the spec enforces this pre-deploy, mirroring the API's own
rejection. Choose the default deliberately: allowlist policies end in a
deny-all default, WAF policies in an allow default.

### Deliberately Unmodeled

- **Labels** — `google_compute_security_policy` has no labels attribute on
  the released `google ~> 6.0` provider line; neither engine attempts
  attribution (a one-engine label would break cross-engine parity).
- **`request_body_inspection_size`** — absent from the released provider
  line (a newer-major-only surface); revisit when the provider line moves.
- **Custom ModSecurity WAF rules** — beyond the preconfigured rule sets;
  expressed through CEL expressions where needed.

---

## References

- [Cloud Armor overview](https://cloud.google.com/armor/docs/overview)
- [Security policies](https://cloud.google.com/armor/docs/security-policy-overview)
- [Preconfigured WAF rules](https://cloud.google.com/armor/docs/waf-rules)
- [CEL expressions](https://cloud.google.com/armor/docs/rules-language-reference)
- [Adaptive Protection](https://cloud.google.com/armor/docs/adaptive-protection-overview)
