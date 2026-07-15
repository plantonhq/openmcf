# AWS WAF Regex Pattern Set

Deploys an AWS WAFv2 regex pattern set — a named, reusable collection of regular expressions. Web ACL rules reference the set by ARN; update the patterns once and every referencing rule sees the change without redeploying the web ACL.

## What Gets Created

When you deploy an AwsWafRegexPatternSet resource, Planton provisions:

- **WAFv2 Regex Pattern Set** — an `aws_wafv2_regex_pattern_set` resource with the configured scope and regular expressions

## Prerequisites

- **AWS credentials** configured via environment variables or Planton provider config
- **Appropriate IAM permissions** for `wafv2:*` operations
- **us-east-1 region** when using `CLOUDFRONT` scope
- **Valid PCRE patterns** — AWS rejects backreferences and lookaround assertions

## Quick Start

Create a file `pattern-set.yaml`:

```yaml
apiVersion: aws.planton.dev/v1
kind: AwsWafRegexPatternSet
metadata:
  name: scanner-probes
  annotations:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: dev.AwsWafRegexPatternSet.scanner-probes
spec:
  region: us-west-2
  scope: REGIONAL
  regularExpressions:
    - ^/wp-admin/.*
    - ^/\.env$
  description: Common scanner probe paths
```

Deploy:

```shell
planton apply -f pattern-set.yaml
```

## Configuration Reference

### Required Fields

| Field | Type | Description | Validation |
|-------|------|-------------|------------|
| `region` | `string` | AWS region. Use `us-east-1` for CLOUDFRONT scope. | Required; non-empty |
| `scope` | `string` | `REGIONAL` or `CLOUDFRONT` | Required; ForceNew |
| `regularExpressions` | `string[]` | Regex patterns (min 1; each 1–200 chars) | Required |

### Optional Fields

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `description` | `string` | — | Human-readable description (max 256 characters). |

`scope` is create-time immutable. AWS's default quota is 10 expressions per set.

## Examples

### Admin Path Protection

```yaml
apiVersion: aws.planton.dev/v1
kind: AwsWafRegexPatternSet
metadata:
  name: internal-admin-paths
spec:
  region: eu-west-1
  scope: REGIONAL
  regularExpressions:
    - ^/internal/admin/.*
    - ^/debug/.*
  description: Internal admin and debug URL patterns
```

### CloudFront Global Scanner Block

```yaml
apiVersion: aws.planton.dev/v1
kind: AwsWafRegexPatternSet
metadata:
  name: global-scanner-probes
spec:
  region: us-east-1
  scope: CLOUDFRONT
  regularExpressions:
    - ^/wp-login\.php
    - ^/xmlrpc\.php
  description: WordPress scanner probes on CloudFront distributions
```

## Stack Outputs

After deployment, the following outputs are available in `status.outputs`:

| Output | Type | Description |
|--------|------|-------------|
| `regex_pattern_set_arn` | `string` | Pattern set ARN for web ACL `regex_pattern_set_reference` statements. |
| `regex_pattern_set_id` | `string` | AWS-assigned UUID. |
| `regex_pattern_set_name` | `string` | Set name in AWS. |

## Related Components

- [AwsWafWebAcl](/docs/catalog/aws/waf-web-acl) — references pattern sets in `regex_pattern_set_reference` rules
- [AwsWafIpSet](/docs/catalog/aws/waf-ip-set) — the sibling reusable set for IP/CIDR matching
