# Office Allow-List

A REGIONAL IPv4 set carrying a /24 office range and a /32 single-host entry — the two CIDR shapes AWS accepts. Pair with a web ACL whose default action is block and an early-priority allow rule referencing this set's ARN.

## When to Use

- Private APIs that should accept traffic only from corporate egress
- Staging environments gated to office/VPN ranges before public launch
- Any allow-list maintained by NetEng independently of application teams

## What It Configures

- **REGIONAL scope** — for ALB, API Gateway, App Runner, and other in-region resources
- **IPv4 family** — deploy a second set with `IPV6` when dual-stack egress matters
- **Two CIDR entries** — a subnet range plus a single-host /32 (WAF never accepts bare addresses)

## What to Customize

- Replace `<aws-region>` with the region your protected resources live in
- Swap the example RFC 5737 ranges for your real office/VPN CIDRs
- Add more entries up to AWS's 10,000-address quota

## Common Pairing

Reference from [AwsWafWebAcl](../../awswafwebacl/v1alpha1/README.md):

```yaml
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
