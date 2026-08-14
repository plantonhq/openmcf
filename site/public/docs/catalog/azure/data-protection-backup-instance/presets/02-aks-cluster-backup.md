---
title: "AKS Cluster Backup"
description: "This preset protects an AKS cluster's workloads: scheduled cluster backups with the plumbing namespace excluded and persistent-volume snapshots enabled, retained per the referenced Kubernetes policy."
type: "preset"
rank: "02"
presetSlug: "02-aks-cluster-backup"
componentSlug: "data-protection-backup-instance"
componentTitle: "Data Protection Backup Instance"
provider: "azure"
icon: "package"
order: 2
---

# AKS Cluster Backup

This preset protects an AKS cluster's workloads: scheduled cluster backups with the plumbing namespace excluded and persistent-volume snapshots enabled, retained per the referenced Kubernetes policy.

## When to Use

- Production AKS clusters whose workloads (and their persistent volumes) need point-in-time recovery
- Protecting against namespace-level mistakes -- a deleted deployment or PVC restores from the last backup
- Compliance regimes that mandate scheduled backups of containerized workloads

## Key Configuration Choices

- **`volumeSnapshotEnabled: true`** -- without it the backup captures Kubernetes objects only, not the data in persistent volumes
- **`kube-system` excluded** -- cluster plumbing is recreated by AKS, not restored from backup; back up what you deploy
- **This variant is immutable end to end** -- every change, the policy binding included, replaces the instance (the service ships no update path)
- **Cluster-side prerequisites are real**: the AKS Backup extension and a trusted-access role binding to the vault must exist before this deploys -- model them as part of cluster provisioning

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<your-backup-vault>` | The AzureDataProtectionBackupVault holding the backups | The vault component's name |
| `<your-kubernetes-policy>` | An AzureDataProtectionBackupPolicy with the `kubernetesCluster` variant, on the same vault | The policy component's name |
| `<your-aks-cluster>` | The AzureAksCluster being protected | The cluster component's name |
| `<your-snapshot-resource-group>` | The AzureResourceGroup where snapshots are stored | The resource-group component's name |

The instance is free; snapshot and backup storage bill per the policy's retention.
