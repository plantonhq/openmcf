# Hardened Baseline

This preset is the strict-posture template for compliance-sensitive
workloads: IMDSv2 with the hop limit locked to the host, root volume
encrypted with a customer-managed KMS key (revocable, auditable), stop and
termination protection, and automatic recovery from hardware impairment.

## When to Use

- Long-lived, stateful, or compliance-scoped instances (databases,
  license servers, security tooling)
- Org-wide golden templates where the security posture must be explicit
  rather than inherited from account defaults

## Key Configuration Choices

- **`httpPutResponseHopLimit: 1`** -- metadata never leaves the host;
  raise to 2 only when containers legitimately need instance credentials
- **Customer-managed KMS key by reference** -- key rotation, revocation,
  and cross-account policy stay in your hands; composes with `AwsKmsKey`
- **`disableApiStop` / `disableApiTermination`** -- API-level guards for
  pet-like instances. NOTE: an auto-scaling group cannot scale in
  termination-protected instances -- use the group's
  `protectFromScaleIn` for fleet members instead
- **`instanceMetadataTags: enabled`** -- on-instance agents read tags
  without `ec2:DescribeTags` permission

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<fleet-name>` | Name for the launch template | Your workload's name |
| `<aws-region>` | AWS region code (e.g., `us-east-1`) | Your deployment region |
| `ami-<your-ami-id>` | AMI to boot from (a hardened/golden image) | Your image pipeline output |
| `<instance-profile-resource-name>` | Name of the AwsIamInstanceProfile resource | Your instance-profile manifest's `metadata.name` |
| `<security-group-resource-name>` | Name of the AwsSecurityGroup resource | Your security-group manifest's `metadata.name` |
| `<kms-key-resource-name>` | Name of the AwsKmsKey resource | Your KMS-key manifest's `metadata.name` |

## Common Additions

- `enclaveEnabled: true` for Nitro Enclaves secret-processing (mutually
  exclusive with hibernation)
- `cpuOptions.amdSevSnp: enabled` on AMD types for memory encryption
- `hibernationEnabled: true` for stateful workloads that resume rather
  than reboot

## Related Presets

- **01-web-server** -- the standard auto-scaled web fleet baseline
- **02-spot-flexible** -- attribute-based Spot capacity
