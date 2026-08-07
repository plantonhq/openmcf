# AwsWafRegexPatternSet

An AWS WAFv2 regex pattern set — a named, reusable collection of regular expressions that web ACL rules match request components against.

## What It Is

A pattern set centralizes regexes that many rules share — scanner probes, admin-path patterns, API-version prefixes. Update the set once and every web ACL rule referencing its ARN sees the change immediately.

A `regex_pattern_set_reference` statement matches when **any** expression in the set matches the inspected component (URI path, query string, header, body, etc.). For a one-off regex no other rule needs, an inline `regex_match` statement on the web ACL is simpler — reach for a pattern set when the same expressions back multiple rules or need independent ownership.

## When to Use It

| Use Case | Description |
|----------|-------------|
| **Scanner probe blocking** | Block common `/wp-admin`, `/.env`, and `/phpmyadmin` probes with one maintained list. |
| **Admin path protection** | Centralize internal admin URL patterns referenced by multiple web ACLs. |
| **API abuse patterns** | Share version-guessing or credential-stuffing path regexes across services. |

## When NOT to Use It

| Need | Use Instead |
|------|-------------|
| **Match client IP addresses** | [AwsWafIpSet](../awswafipset/README.md) |
| **One regex used once** | Inline `regex_match` on the web ACL |

## Key Facts

- **At least one expression required.** Unlike IP sets, AWS rejects an empty pattern set.
- **PCRE subset.** AWS supports a subset of PCRE — no backreferences (`\1`) and no lookaround (`(?=...)`, `(?<=...)`). Invalid patterns fail at AWS create time.
- **Scope is create-time immutable.** CLOUDFRONT sets must be created in **us-east-1**.
- **Default quota: 10 expressions per set** (adjustable via Service Quotas). Each expression is up to 200 characters.

## Spec Fields

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `region` | string | **Yes** | AWS region. Must be `us-east-1` when `scope` is `CLOUDFRONT`. |
| `scope` | string | **Yes** | `REGIONAL` or `CLOUDFRONT`. **ForceNew.** |
| `regular_expressions` | string[] | **Yes** | Regex patterns (1–10 by default quota; each ≤200 chars). |
| `description` | string | No | What the patterns match (max 256 chars). |

## Outputs

| Field | Type | Description |
|-------|------|-------------|
| `regex_pattern_set_arn` | string | The ARN web ACL `regex_pattern_set_reference` statements reference. |
| `regex_pattern_set_id` | string | AWS-assigned UUID for direct WAFv2 API calls. |
| `regex_pattern_set_name` | string | Set name in AWS (from `metadata.name`). |

## Example

```yaml
apiVersion: aws.planton.dev/v1alpha1
kind: AwsWafRegexPatternSet
metadata:
  name: scanner-probes
  org: my-org
spec:
  region: us-west-2
  scope: REGIONAL
  regularExpressions:
    - ^/wp-admin/.*
    - ^/\.env$
    - ^/phpmyadmin/.*
  description: Common scanner probe paths, owned by SecEng
```

Reference it from a web ACL rule:

```yaml
spec:
  rules:
    - name: block-scanner-probes
      priority: 10
      action: block
      statement:
        regexPatternSetReference:
          arn:
            valueFrom:
              kind: AwsWafRegexPatternSet
              name: scanner-probes
              fieldPath: status.outputs.regex_pattern_set_arn
          fieldToMatch:
            uriPath: {}
```

See docs/README.md and [AwsWafWebAcl](../awswafwebacl/README.md).

---

© Planton. Licensed under [Apache-2.0](https://github.com/plantonhq/planton/blob/main/LICENSE).
