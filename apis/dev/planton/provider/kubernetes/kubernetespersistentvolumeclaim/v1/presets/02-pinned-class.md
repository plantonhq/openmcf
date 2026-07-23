# Pinned Class

This preset is the production claim shape: the StorageClass is named explicitly instead of inherited from the cluster default. The default class is a moving target — a cluster upgrade or a platform change can swap it — while a pinned class makes the claim's performance tier, reclaim policy, and expandability an explicit, reviewable choice.

## When to Use

- Any production claim — databases, queues, anything whose storage tier matters
- Pinning a claim to a platform-defined class (e.g. one created with the KubernetesStorageClass presets: `aws-ebs-gp3`, `gcp-pd-ssd`, `azure-premium`)
- Claims that must be expandable: pin to a class with `allow_volume_expansion: true` so a full disk is an edit, not a migration

## Key Configuration Choices

- **`storage_class_name.value`** — a literal class name. In a Planton chart, this can instead be a reference to a `KubernetesStorageClass` resource (`valueFrom`), resolving through the class's exported `storage_class_name` output so the class and claim deploy in one run
- **`access_modes: [ReadWriteOnce]`** — stated explicitly for review, though it is also the default; one node at a time, the mode every block driver supports
- **`storage_request: 50Gi`** — sized for growth. **Expansion note**: raising this later works only if the pinned class allows volume expansion (grow only — Kubernetes never shrinks a volume); on a non-expandable class, outgrowing the disk means a data migration
- **Pending until consumed is normal** — the recommended cloud classes bind on first consumer, so the claim stays **Pending until a pod uses it — correct behavior, not an error**

## Placeholders to Replace

| Placeholder | Description | Where to Find |
|---|---|---|
| `<your-namespace>` | The namespace the claim lives in | Your namespace management |
| `<your-storage-class>` | The StorageClass name — e.g. a class created with the KubernetesStorageClass presets | `kubectl get storageclass`, or the class resource's outputs |

The size `50Gi` is a working example — replace it with the actual requirement.

## Related Presets

- **01-dynamic-default-class** — the simpler shape when the cluster default is acceptable
- **03-shared-rwm** — a volume shared by several pods
- **04-clone-from-claim** — start from a copy of an existing claim's data
