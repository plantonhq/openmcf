---
title: "Pinned Version with Configuration"
description: "This preset pins an add-on to an exact version and carries custom configuration -- the controlled-upgrade pattern for fleets that must run byte-identical clusters."
type: "preset"
rank: "03"
presetSlug: "03-pinned-version"
componentSlug: "eks-addon"
componentTitle: "EKS Addon"
provider: "aws"
icon: "package"
order: 3
---

# Pinned Version with Configuration

This preset pins an add-on to an exact version and carries custom
configuration -- the controlled-upgrade pattern for fleets that must
run byte-identical clusters.

## When to Use

- Fleets where every cluster must run the same audited add-on version
- Add-ons whose defaults need tuning (coredns replica count, vpc-cni
  prefix delegation, ...)
- Change-controlled environments where upgrades happen on a schedule,
  not when AWS moves its default

## Key Configuration Choices

- **`addonVersion` pinned** -- the add-on stays put until the pin
  moves; bumping it rolls the add-on's own pods (never the nodes).
  Empty would follow the AWS default for the cluster's Kubernetes
  version instead.
- **`configurationValues` as JSON** -- passed through to the add-on's
  own schema unmodified; validate against
  `aws eks describe-addon-configuration` output
- **`resolveConflictsOnUpdate: OVERWRITE`** -- the pinned configuration
  wins over out-of-band edits at every update

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<addon-resource-name>` | Name for this add-on resource (e.g. `platform-coredns`) | Your naming convention |
| `<aws-region>` | AWS region code (e.g., `us-west-2`) | Your deployment region |
| `<cluster-resource-name>` | Name of the AwsEksCluster resource | Your cluster manifest's `metadata.name` |
| `addonVersion` value | Replace the example pin with your audited version | `aws eks describe-addon-versions --addon-name coredns` |

## Common Additions

- `preserve: true` if deleting this resource should leave the software
  running as self-managed

## Related Presets

- **01-core-networking** -- AWS-default versions that never go stale
- **02-ebs-csi-pod-identity** -- a storage add-on with its own IAM
  identity
