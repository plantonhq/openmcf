# AWS WAF IP Set: Concepts

An IP set is WAFv2's reusable address book. This reference covers what the set actually owns, how web ACL rules consume it, and the constraints that bite operators in production.

## Why IP Sets Exist

Without sets, every allow-list or deny-list rule carries its own copy of the same CIDR ranges. When NetEng adds a new office egress range, every rule must be edited and every web ACL redeployed. IP sets invert that:

- **Central ownership** — one resource holds the addresses; many rules reference one ARN.
- **In-place updates** — changing `addresses` updates the set without touching referencing web ACLs.
- **Scope isolation** — REGIONAL and CLOUDFRONT sets are separate namespaces; a web ACL can only reference sets of its own scope.

The trade-off: the set carries **no action**. Allow, block, count, and CAPTCHA are configured on the referencing rule in the web ACL.

## Design Notes

- **CIDR-only addresses.** AWS rejects bare IPs. A single host is always `/32` (IPv4) or `/128` (IPv6). The spec's CEL enforces the `/nn` suffix at validation time so manifests fail before apply.
- **Dual-stack deployments need two sets.** `ip_address_version` is ForceNew — one set holds IPv4 or IPv6, never both. Match both families with two sets referenced by two rules or one OR statement.
- **Empty sets are intentional.** An empty set matches nothing. Deploy the set and wire rules to its ARN before the address list is finalized — a common pattern for infrastructure-as-code pipelines.
- **Deletion guard.** AWS refuses to delete a set still referenced by a web ACL rule. Remove or retarget the rule first.
- **CLOUDFRONT lives in us-east-1.** CloudFront-scoped WAF resources are global; AWS stores them in the us-east-1 WAF partition regardless of where your distributions originate.

## Composition

| Consumer | What it references | Why |
|----------|-------------------|-----|
| `AwsWafWebAcl` rule `ip_set_reference.arn` | `status.outputs.ip_set_arn` | The statement matches when the request source IP is in the set. |
| Multiple web ACLs | Same `ip_set_arn` | One maintained list shared across applications. |

A typical allow-list pattern pairs a restrictive web ACL default action with an early-priority allow rule:

```yaml
# AwsWafWebAcl excerpt
defaultAction:
  type: block
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

See the [AwsWafWebAcl architecture reference](../../awswafwebacl/v1/docs/README.md) for the full rule tree and association model.
