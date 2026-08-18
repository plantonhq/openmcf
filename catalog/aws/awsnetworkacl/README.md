# AwsNetworkAcl

One network ACL: the stateless subnet-level firewall — ordered allow/deny rules per direction (first match wins) and atomic subnet associations, all managed in-line as the single declarative owner.

## Highlights

- **NACLs can DENY** — the one thing security groups cannot do; rules evaluate in `rule_no` order and unmatched traffic falls through to AWS's invisible catch-all deny.
- **Stateless taught everywhere**: replies are not tracked — the spec, GUIDE, and presets all model the outbound ephemeral-port allow that stateless filtering demands.
- **In-line satellites, single owner**: rules and subnet associations live on the ACL (the standalone rule/association resources are identical payloads that fight the in-line form and record composed); AWS's 32767/32768 catch-alls are unmanageable and rule numbers stop at 32766.
- **Provider truths as CELs**: exactly one address family per rule, all-protocols rules forbid ports, ICMP selectors demand ICMP protocols, all-types pairs with all-codes, unique rule numbers per direction.

## Both Engines

Both modules render the single resource with in-line dynamic rule blocks identically (protocol names pass through — the provider normalizes them to the numbers AWS stores) and export the same outputs: `network_acl_id` (import ID), `network_acl_arn`, `owner_id`.

## Chart Wiring

`vpc_id` → AwsVpc `vpc_id`; `subnet_ids` → AwsSubnet `subnet_id` outputs. Pair with AwsSecurityGroup for the stateful instance-level half of the story.

---

© Planton. Licensed under [Apache-2.0](https://github.com/plantonhq/planton/blob/main/LICENSE).
