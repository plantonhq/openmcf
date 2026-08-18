# AwsNetworkAcl — Pulumi module (Go)

Manages one network ACL (`ec2.NetworkAcl`) with its rules and subnet associations in-line.

Module facts worth knowing before editing:

- **`VpcId` replaces the ACL** — everything else updates in place.
- **In-line rules and `SubnetIds` are the single declarative owner** — the standalone `ec2.NetworkAclRule` / `ec2.NetworkAclAssociation` resources are identical payloads and fight this form; this module never uses them.
- **Protocols are stored as numbers at AWS** — the provider normalizes names ("tcp" → 6) in its rule hash, so names never cause perpetual diffs.
- **AWS's catch-all rules (32767 IPv4 / 32768 IPv6) are invisible** — they can be neither configured nor deleted; spec rules stay at 1–32766.
- **Subnet association is replace, not attach** — a subnet always belongs to exactly one ACL; listing it here atomically replaces its previous association, removing it hands it back to the VPC's default NACL.
- **Destroy stomps associations first** — the provider removes all subnet associations (even externally created ones) before deleting the ACL, retrying on DependencyViolation.

Outputs mirror the Terraform module key-for-key: `network_acl_id` (import ID), `network_acl_arn`, `owner_id`.
