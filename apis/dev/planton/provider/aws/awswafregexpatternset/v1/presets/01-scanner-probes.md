# Scanner Probe Patterns

A REGIONAL set of common automated-scanner URI paths. Pair with a web ACL block rule whose `regex_pattern_set_reference` inspects `uriPath`.

## When to Use

- Public-facing web apps that see constant `/wp-admin`, `/.env`, and CMS probe traffic
- Baseline hardening before layering managed rule groups
- SecOps-owned probe lists shared across many services

## What It Configures

- **Three anchored path patterns** — WordPress admin, dotenv leak, phpMyAdmin probes
- **REGIONAL scope** — for ALB and API Gateway front doors

## What to Customize

- Replace `<aws-region>` with your target region
- Extend with your application's known-abuse patterns (stay within the 10-expression default quota)
- Tighten or loosen anchors depending on whether subpaths should match

## Common Pairing

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
