---
title: "S3 backups preset"
description: "A production-shaped cluster (three nodes, quorum ZooKeeper, persistent storage, basic auth) with an S3 backup repository registered from day one. Backup and restore run as SolrBackup operations..."
type: "preset"
rank: "03"
presetSlug: "03-s3-backups"
componentSlug: "apache-solr"
componentTitle: "Apache Solr"
provider: "kubernetes"
icon: "package"
order: 3
---

# S3 backups preset

A production-shaped cluster (three nodes, quorum ZooKeeper,
persistent storage, basic auth) with an S3 backup repository
registered from day one. Backup and restore run as SolrBackup
operations against the repository name; `base_location` chroots this
cluster into its own prefix so one bucket can serve several clusters.
The Solr module that does the S3 I/O loads automatically because the
repository is declared — `solr_modules` stays for other modules only.

The credential model is the teaching point. Static keys ride
references to an existing Secret — nothing secret appears in this
manifest — and on EKS with IRSA (or any workload identity granting
bucket access) you should delete the `credentials` block entirely:
the S3 client falls back to the ambient AWS identity and no static
credential exists anywhere. Prefer that keyless path wherever it is
available.

Do not use this preset expecting to bolt the repository on later
without cost: repositories are part of the SolrCloud spec, so adding
or changing one is a rolling restart — declare the backup target
up front. Change first: the bucket, region and `base_location`, and
either create the `solr-s3-credentials` Secret before applying or
delete the credentials block for the keyless path.

See [03-s3-backups.yaml](./03-s3-backups.yaml) for the manifest.
