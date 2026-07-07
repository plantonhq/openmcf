# AWS WAF Family Depth Pass

**Date**: July 7, 2026
**Type**: Feature
**Components**: API Definitions, Terraform Generator, Pulumi CLI Integration, E2E Harness

## Summary

Rebuilt `AwsWafWebAcl` to model the full WAFv2 statement language as a typed recursive tree, forged first-class `AwsWafIpSet` (361) and `AwsWafRegexPatternSet` (362), hardened the Terraform variables generator for recursive proto messages, and wired ALB/CloudFront composability via `web_acl_arn`. Both IaC engines serialize rules to canonical AWS `rule_json` for parity. Live dual-engine E2E lanes are deferred pending AWS SSO credential refresh; offline gate is green.

## Problem Statement

The original WAF web ACL component used a hybrid four-statement model plus a Struct escape hatch — adequate for demos but not 90/10 coverage. IP sets and regex pattern sets were documentation-only references. The Terraform variables generator hung on recursive proto descriptors. ALB and CloudFront used inconsistent WAF association field names.

## Solution

### Full statement tree on AwsWafWebAcl

- Recursive `AwsWafWebAclStatement` oneof covering managed groups (with ATP/ACFP/Bot Control configs), rate limiting (custom aggregation keys), set references, geo/byte/SQLi/XSS/size/regex/label/ASN matching, and AND/OR/NOT composition
- Top-level CAPTCHA/challenge immunity, association config, data protection, logging
- `application_integration_url` stack output
- Terraform module builds `rule_json` via three-level statement unroll; Pulumi module mirrors in Go (`rules.go` + unit tests)

### New composable set kinds

- **AwsWafIpSet** — CIDR-only addresses, IPV4/IPV6 family, REGIONAL/CLOUDFRONT scope with CEL coupling to us-east-1
- **AwsWafRegexPatternSet** — shared regex library (min 1 expression), PCRE-subset honesty documented

### Framework: recursive proto → Terraform types

- Path-scoped `visited` cycle detection in `variablestf.go`
- `TFFreeFormList` / `TFFreeFormMap` when heterogeneous `any` leaves would break Terraform list/map typing
- Regression test using `descriptorpb.DescriptorProto`
- Forge rule 013 uplift documenting the pattern

### Composability seams

- `AwsAlb.spec.web_acl_arn` → `aws_wafv2_web_acl_association`
- `AwsCloudFront.spec.web_acl_id` renamed to `web_acl_arn` (breaking)

## Validation

- Unit tests: generator recursion, all three WAF spec tests, Pulumi `rules_test.go`, outputs conformance
- `tofu validate` green on all three TF modules (v6 floor)
- `planton validate-refs --check` green
- Pulumi entrypoints build; E2E verify package builds
- Site catalog regenerated (`waf-ip-set`, `waf-regex-pattern-set`, updated `waf-web-acl`)
- Live E2E: **deferred** — AWS SSO token expired (`AWS_PROFILE=planton-aws-e2e`); profiles record re-run instructions (~6 lanes, each ≤15 min)

## Breaking Changes

- `AwsWafWebAcl` spec reshaped (recursive statement tree; preset YAML updated)
- `AwsCloudFront.web_acl_id` → `web_acl_arn`
