# Core Networking Add-on

This preset adopts one of the cluster's core networking add-ons
(vpc-cni, coredns, kube-proxy) into AWS-managed lifecycle -- upgrades
and configuration flow through the EKS control plane instead of
whatever bootstrap copy the cluster started with.

## When to Use

- Every production cluster: the core trio should be managed, not
  bootstrap leftovers
- Migrating a self-managed add-on install to managed lifecycle
- Standardizing add-on versions across a fleet of clusters

## Key Configuration Choices

- **`resolveConflictsOnCreate: OVERWRITE`** -- clusters created with
  bootstrap self-managed add-ons (the AWS default) already run copies
  of the core trio; OVERWRITE adopts them instead of failing with
  `ConfigurationConflict`
- **`resolveConflictsOnUpdate: OVERWRITE`** -- out-of-band kubectl
  edits are restored to managed values on the next update; use
  `PRESERVE` instead if operators intentionally hand-tune
- **No `addonVersion`** -- follows the AWS default for the cluster's
  Kubernetes version, so the manifest never goes stale; pin a version
  for byte-identical clusters

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<addon-resource-name>` | Name for this add-on resource (e.g. `platform-vpc-cni`) | Your naming convention |
| `<aws-region>` | AWS region code (e.g., `us-west-2`) | Your deployment region |
| `<cluster-resource-name>` | Name of the AwsEksCluster resource | Your cluster manifest's `metadata.name` |

## Common Additions

- `configurationValues` for add-on-specific tuning (each add-on
  publishes its own JSON schema)
- One resource per core add-on -- deploy this preset three times
  (vpc-cni, coredns, kube-proxy)

## Related Presets

- **02-ebs-csi-pod-identity** -- a storage add-on with its own IAM
  identity
- **03-pinned-version** -- version-pinned add-on for controlled
  upgrades
