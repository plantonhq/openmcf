# AWS WAF IP Set

Deploys an AWS WAFv2 IP set — a named, reusable collection of IP addresses and CIDR ranges. Web ACL rules reference the set by ARN; update the addresses once and every referencing rule sees the change without redeploying the web ACL.

## What Gets Created

When you deploy an AwsWafIpSet resource, Planton provisions:

- **WAFv2 IP Set** — an `aws_wafv2_ip_set` resource with the configured scope, address family, and CIDR entries

## Prerequisites

- **AWS credentials** configured via environment variables or Planton provider config
- **Appropriate IAM permissions** for `wafv2:*` operations
- **us-east-1 region** when using `CLOUDFRONT` scope

## Quick Start

Create a file `ip-set.yaml`:

```yaml
apiVersion: aws.planton.dev/v1
kind: AwsWafIpSet
metadata:
  name: office-allowlist
  labels:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: dev.AwsWafIpSet.office-allowlist
spec:
  region: us-west-2
  scope: REGIONAL
  ipAddressVersion: IPV4
  addresses:
    - 203.0.113.0/24
    - 198.51.100.44/32
  description: Corporate office egress ranges
```

Deploy:

```shell
planton apply -f ip-set.yaml
```

## Configuration Reference

### Required Fields

| Field | Type | Description | Validation |
|-------|------|-------------|------------|
| `region` | `string` | AWS region. Use `us-east-1` for CLOUDFRONT scope. | Required; non-empty |
| `scope` | `string` | `REGIONAL` or `CLOUDFRONT` | Required; ForceNew |
| `ipAddressVersion` | `string` | `IPV4` or `IPV6` | Required; ForceNew |

### Optional Fields

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `addresses` | `string[]` | `[]` | CIDR ranges (up to 10,000). Empty list matches nothing. Each entry must include a `/nn` suffix. |
| `description` | `string` | — | Human-readable description (max 256 characters). |

`scope` and `ipAddressVersion` are create-time immutable — changing either replaces the set.

## Examples

### Partner Integration Allow-List

```yaml
apiVersion: aws.planton.dev/v1
kind: AwsWafIpSet
metadata:
  name: partner-apis
spec:
  region: eu-west-1
  scope: REGIONAL
  ipAddressVersion: IPV4
  addresses:
    - 192.0.2.0/24
    - 198.51.100.0/24
  description: Partner API integration egress ranges
```

### CloudFront Global Deny-List

```yaml
apiVersion: aws.planton.dev/v1
kind: AwsWafIpSet
metadata:
  name: blocked-clients
spec:
  region: us-east-1
  scope: CLOUDFRONT
  ipAddressVersion: IPV4
  addresses:
    - 203.0.113.55/32
  description: Known abusive clients, SecOps maintained
```

## Stack Outputs

After deployment, the following outputs are available in `status.outputs`:

| Output | Type | Description |
|--------|------|-------------|
| `ip_set_arn` | `string` | IP set ARN for web ACL `ip_set_reference` statements. |
| `ip_set_id` | `string` | AWS-assigned UUID. |
| `ip_set_name` | `string` | Set name in AWS. |

## Related Components

- [AwsWafWebAcl](/docs/catalog/aws/waf-web-acl) — references IP sets in `ip_set_reference` rules
- [AwsAlb](/docs/catalog/aws/alb) — associates a REGIONAL web ACL via `web_acl_arn`
- [AwsCloudFront](/docs/catalog/aws/cloudfront) — associates a CLOUDFRONT web ACL via `web_acl_arn`
