# Azure Front Door WAF: Transform Wire-Casing Fix and Review Uplift

**Date**: July 8, 2026
**Type**: Bug Fix
**Components**: Azure Provider, API Definitions, Testing Framework

## Summary

A full-depth review of the Front Door edge-security pair (`AzureFrontDoorFirewallPolicy` 488, `AzureFrontDoorSecurityPolicy` 489) caught one apply-blocking defect: both IaC engines mapped the `URL_DECODE`/`URL_ENCODE` custom-rule transforms to `"URLDecode"`/`"URLEncode"`, but ARM's canonical values are `"UrlDecode"`/`"UrlEncode"` -- the SDK's Go *identifiers* are `URLDecode`/`URLEncode` while their *string values* carry the lower casing, and the provider validates case-sensitively. Any custom rule using those two transforms failed at plan time on both engines. Fixed in both modules and live-proven; the review's smaller findings (a mislabeled spec test, missing boundary tests, missing preview-status spec comments) are closed in the same pass.

## Problem Statement / Motivation

The transform maps in both modules copied the SDK constant *identifiers* instead of their *values*:

```go
// go-azure-sdk frontdoor/2025-03-01 constants.go -- identifier vs value:
TransformTypeURLDecode TransformType = "UrlDecode"
TransformTypeURLEncode TransformType = "UrlEncode"
```

The provider's schema validates transforms with a case-sensitive `StringInSlice`, so `"URLDecode"` was rejected before ARM was ever reached. The defect escaped the original live E2E because no scenario, preset, or hack manifest happened to use those two transforms -- the mapping rows were dead in every validation path.

## Solution / What's New

- **Both transform maps corrected** to `"UrlDecode"`/`"UrlEncode"`:
  - `apis/dev/planton/provider/azure/azurefrontdoorfirewallpolicy/v1/iac/pulumi/module/locals.go`
  - `apis/dev/planton/provider/azure/azurefrontdoorfirewallpolicy/v1/iac/tf/locals.tf`
  - Both maps now carry a comment warning that the SDK identifier casing is NOT the wire value.
- **The coverage hole closed**: the minimal E2E scenario and the hack manifest now include `URL_DECODE` in a transform list, so offline plans and every future live run exercise the previously dead mapping rows.
- **Spec-test uplift** (firewall policy 45 → 57 cases, security policy 11 → 13):
  - The "anomaly-scoring overrides on a 2.x set" test now actually asserts `OVERRIDE_ANOMALY_SCORING` valid (it previously only used `OVERRIDE_LOG` -- the gate's central positive branch was untested).
  - New boundary cases: 128-char policy name, block-status 990/999 valid + 989 invalid, challenge lifetimes at the 5/1440 edges (valid) and both reject sides for both fields, CAPTCHA action valid on PREMIUM, `DefaultRuleSet preview-0.1` valid, `Microsoft_DefaultRuleSet 0.x` invalid, >10 match conditions, >600 match values, >100 scrubbing rules.
  - Security policy: exactly-500 domains valid (the Premium cap edge), present-but-empty `StringValueOrRef` rejected.
- **Preview-status notes added to the spec** (comment-only; stubs regenerated): the JS_CHALLENGE/CAPTCHA custom-rule actions, both challenge-lifetime fields, and log scrubbing now carry Azure's PREVIEW flag -- making the existing docs/README claim true.
- **E2E YAML made timeless**: the bot-manager `Bot`-prefix fact now lives as an inline comment next to the rule ID in the premium scenario instead of session narration in the profile.

## Validation (what ran and passed)

- Spec tests both kinds ALL PASS (57 + 13); targeted Go builds of both module trees.
- Offline `planton tofu plan` with the URL_DECODE-bearing hack manifest -- the plan now renders `"UrlDecode"` and passes the provider's case-sensitive validator (it would have failed before the fix).
- **Live re-run** (test subscription `8158df85-…`): firewall-policy minimal scenario green on both engines (Pulumi 123s, Terraform 208s), with the `UrlDecode` transform applied and destroyed against real ARM. Zero-orphan sweep clean.
- Stub regen was comment-only (verified: no non-comment diff lines in `spec.pb.go`).

## Impact

- Custom rules using the URL-decode/encode transforms -- a standard evasion-hardening pattern (`LOWERCASE` + `URL_DECODE`) -- now deploy on both engines instead of failing plan validation.
- The vendor-constant trap (SDK identifier casing vs wire value) is now documented at both mapping sites, and the transform rows can no longer silently regress: they are exercised by the offline plan gate and the minimal live scenario.

## Related Work

- Azure Front Door edge security: firewall and security policy kinds (2026-07-08) -- the pair this review hardened
- The same-day vendor-catalog-ID lesson in `e2e/README.md` (look up, never infer) is the sibling failure class: this one is vendor-constant *casing* rather than catalog IDs

---

**Status**: ✅ Production Ready
