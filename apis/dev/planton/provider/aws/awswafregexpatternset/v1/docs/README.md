# AWS WAF Regex Pattern Set: Concepts

A regex pattern set is WAFv2's reusable pattern library. This reference covers what the set owns, how web ACL rules consume it, and the PCRE constraints AWS enforces.

## Why Pattern Sets Exist

Inline `regex_match` statements work for one-off patterns. When SecOps maintains a probe block-list referenced by five web ACLs, duplicating the same regexes across rules creates drift. Pattern sets centralize that:

- **Central ownership** — one resource holds the expressions; many rules reference one ARN.
- **In-place updates** — changing `regular_expressions` updates the set without touching referencing web ACLs.
- **OR semantics** — a `regex_pattern_set_reference` matches when **any** expression in the set matches the inspected field.

## Design Notes

- **Minimum one expression.** Unlike IP sets, AWS rejects an empty pattern set. Plan for at least one pattern at create time.
- **PCRE subset.** AWS WAFv2 supports a subset of PCRE. Backreferences (`\1`) and lookaround (`(?=...)`, `(?<=...)`, `(?!...)`, `(?<!...)`) are rejected at create time. Anchors (`^`, `$`), character classes, and quantifiers work as expected.
- **Field targeting lives on the rule.** The set holds patterns only; the referencing statement's `field_to_match` chooses URI path, query string, header, body, etc.
- **Quota awareness.** AWS defaults to 10 expressions per set (adjustable). Each expression is capped at 200 characters.
- **CLOUDFRONT lives in us-east-1.** Same regional constraint as IP sets and web ACLs.

## Composition

| Consumer | What it references | Why |
|----------|-------------------|-----|
| `AwsWafWebAcl` rule `regex_pattern_set_reference.arn` | `status.outputs.regex_pattern_set_arn` | Matches when any set expression hits the configured field. |
| Multiple web ACLs | Same ARN | One maintained probe list shared across applications. |

Pair with a block rule on the URI path:

```yaml
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

See the [AwsWafWebAcl architecture reference](../../awswafwebacl/v1/docs/README.md) for the full statement tree.
