# Clone From Claim

This preset creates a new volume populated as a COPY of an existing PersistentVolumeClaim's data, instead of provisioned empty. The `data_source` names a same-namespace source claim; the CSI driver performs the clone. It is the shape for spinning up a staging database from production data, testing a risky migration against a real copy, or forking a dataset.

**This preset deploys via the Pulumi engine only.** The Terraform Kubernetes provider cannot express a PVC data source at all, and the Terraform module rejects a manifest with `data_source` set at plan time with a clear precondition error — by design, because silently provisioning an EMPTY volume where a copy was asked for would be worse than failing. Use the Pulumi provisioner annotation on this manifest.

## When to Use

- Staging or test environments seeded from a real dataset
- Rehearsing a schema migration or upgrade against a copy before touching the original
- Forking a dataset so two workloads can diverge from a common starting point

## Key Configuration Choices

- **`data_source.kind: persistent_volume_claim`** — a clone. The other supported kind is `volume_snapshot`, which restores a VolumeSnapshot (`snapshot.storage.k8s.io`) instead
- **Same namespace only** — the source claim must live in the claim's own namespace (cross-namespace sources are feature-gated upstream and deliberately unmodeled)
- **Driver support required** — cloning is a CSI capability; the major cloud block drivers implement it, but verify before building runbooks on it
- **`storage_request` must be at least the source's size** — a clone cannot be smaller than what it copies
- **Same class, typically** — clones are generally created in the same StorageClass as the source; pin the class explicitly to make that true by construction

## Placeholders to Replace

| Placeholder | Description | Where to Find |
|---|---|---|
| `<your-namespace>` | The namespace — both this claim AND the source claim live here | Your namespace management |
| `<your-storage-class>` | The StorageClass — typically the source claim's class | `kubectl get pvc <source> -o jsonpath='{.spec.storageClassName}'` |
| `<source-claim-name>` | The PersistentVolumeClaim to copy | `kubectl get pvc` in the namespace |

The size `50Gi` is a working example — it must be at least the source claim's requested size.

## Related Presets

- **02-pinned-class** — the plain empty-volume shape this preset builds on
- **01-dynamic-default-class** — the simplest claim, default class
