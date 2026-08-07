# AwsWafIpSet

An AWS WAFv2 IP set — a named, reusable collection of IP addresses and CIDR ranges that web ACL rules match requests against.

## What It Is

An IP set centralizes the addresses your WAF rules care about. Office egress ranges, partner integrations, health-checker fleets, and known-bad clients all belong in a set rather than duplicated across rules. Update the set once and every web ACL rule referencing its ARN sees the change immediately — no web ACL redeploy.

The action (allow, block, count, CAPTCHA) lives on the **referencing rule** in the web ACL, not on the set itself. The same set can back an allow rule in one ACL and a block rule in another.

## When to Use It

| Use Case | Description |
|----------|-------------|
| **Office/VPN allow-list** | Permit traffic only from corporate egress ranges on a private API. |
| **Partner integration** | Maintain a shared set of partner IP ranges referenced by multiple web ACLs. |
| **Deny-list** | Block known-abusive ranges without touching the web ACL rule tree. |
| **Health-checker fleet** | Give load-balancer health checks a stable set that ops can update independently. |

## When NOT to Use It

| Need | Use Instead |
|------|-------------|
| **One-off regex on a URI path** | An inline `regex_match` statement on the web ACL, or an [AwsWafRegexPatternSet](../awswafregexpatternset/README.md) when the same patterns back multiple rules. |
| **Rate limiting by IP** | A `rate_based` statement on the web ACL (optionally with custom aggregation keys). |

## Key Facts

- **Scope is create-time immutable.** A REGIONAL set protects resources in its region; a CLOUDFRONT set must be created in **us-east-1** and is referenced from CloudFront distributions.
- **Address family is create-time immutable.** IPv4 and IPv6 require separate sets. Match both families with two sets or one rule with an OR statement.
- **CIDR notation only.** AWS rejects bare addresses — write a single IPv4 host as `203.0.113.44/32` and a single IPv6 host as `2001:db8::1/128`.
- **Empty sets are valid.** A set with no addresses matches nothing — useful as a placeholder referenced before ranges are known.
- **Deletion is blocked while referenced.** AWS refuses to delete a set still referenced by a web ACL rule. Detach the rule first.

## Spec Fields

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `region` | string | **Yes** | AWS region. Must be `us-east-1` when `scope` is `CLOUDFRONT`. |
| `scope` | string | **Yes** | `REGIONAL` or `CLOUDFRONT`. **ForceNew.** |
| `ip_address_version` | string | **Yes** | `IPV4` or `IPV6`. **ForceNew.** |
| `addresses` | string[] | No | CIDR ranges (up to 10,000). Empty list is valid. |
| `description` | string | No | What the set represents (max 256 chars). |

## Outputs

| Field | Type | Description |
|-------|------|-------------|
| `ip_set_arn` | string | The ARN web ACL `ip_set_reference` statements reference. |
| `ip_set_id` | string | AWS-assigned UUID for direct WAFv2 API calls. |
| `ip_set_name` | string | Set name in AWS (from `metadata.name`). |

## Example

```yaml
apiVersion: aws.planton.dev/v1alpha1
kind: AwsWafIpSet
metadata:
  name: office-allowlist
  org: my-org
spec:
  region: us-west-2
  scope: REGIONAL
  ipAddressVersion: IPV4
  addresses:
    - 203.0.113.0/24
    - 198.51.100.44/32
  description: Corporate office egress ranges, owned by NetEng
```

Reference it from a web ACL rule:

```yaml
spec:
  rules:
    - name: allow-office
      priority: 1
      action: allow
      statement:
        ipSetReference:
          arn:
            valueFrom:
              kind: AwsWafIpSet
              name: office-allowlist
              fieldPath: status.outputs.ip_set_arn
```

See docs/README.md for composition patterns and [AwsWafWebAcl](../awswafwebacl/README.md) for the rule tree.

---

© Planton. Licensed under [Apache-2.0](https://github.com/plantonhq/planton/blob/main/LICENSE).
