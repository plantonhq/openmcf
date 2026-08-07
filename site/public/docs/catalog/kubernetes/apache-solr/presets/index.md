---
title: "Presets"
description: "Ready-to-deploy configuration presets for Apache Solr"
type: "preset-list"
componentSlug: "apache-solr"
componentTitle: "Apache Solr"
provider: "kubernetes"
icon: "package"
order: 200
presets:
  - slug: "01-dev-single-node"
    rank: "01"
    title: "Dev single node preset"
    excerpt: "The smallest declarable SolrCloud that actually serves: one Solr node, a single-member provided ZooKeeper, ephemeral storage, and no authentication. For developers and CI who need the real SolrCloud..."
  - slug: "02-production-cloud"
    rank: "02"
    title: "Production cloud preset"
    excerpt: "The production topology: three Solr nodes on persistent fast storage, a three-member ZooKeeper quorum with its own persistent volumes, basic-auth security bootstrapped by the operator, shard-aware..."
  - slug: "03-s3-backups"
    rank: "03"
    title: "S3 backups preset"
    excerpt: "A production-shaped cluster (three nodes, quorum ZooKeeper, persistent storage, basic auth) with an S3 backup repository registered from day one. Backup and restore run as SolrBackup operations..."
---

# Apache Solr Presets

Ready-to-deploy configuration presets for Apache Solr. Each preset is a complete manifest you can copy, customize, and deploy.
