---
title: "S3 snapshots preset"
description: "A three-node cluster with the backup story wired end to end: the `repository-s3` plugin installed at node startup, S3 credentials loaded into the OpenSearch keystore from an existing Secret, and an..."
type: "preset"
rank: "03"
presetSlug: "03-s3-snapshots"
componentSlug: "opensearch"
componentTitle: "OpenSearch"
provider: "kubernetes"
icon: "package"
order: 3
---

# S3 snapshots preset

A three-node cluster with the backup story wired end to end: the
`repository-s3` plugin installed at node startup, S3 credentials
loaded into the OpenSearch keystore from an existing Secret, and an
S3 snapshot repository registered on the cluster. Snapshot and
restore calls (and index-management snapshot policies) reference the
repository by name from day one.

The credential model is the teaching point. Static keys belong in
the KEYSTORE — the settings map carries only bucket/region/paths, and
nothing secret ever appears in this manifest. Better still, on EKS
with IRSA (or any workload identity granting bucket access) delete
the `keystore` block entirely: the s3 client falls back to the
ambient AWS identity and no static credential exists at all. Prefer
that keyless path wherever it is available.

Do not use this preset on air-gapped clusters as-is — `plugins_list`
downloads the plugin from the internet at every pod start, and a
failed download crash-loops the node; bake the plugin into a custom
image there. Change first: the bucket, region and `base_path` in the
repository settings, and either create the credentials Secret
(`opensearch-s3-credentials` with `access-key`/`secret-key`) before
applying or delete the keystore block for the keyless path.

See [03-s3-snapshots.yaml](./03-s3-snapshots.yaml) for the manifest.
