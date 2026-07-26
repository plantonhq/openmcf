---
title: "Namespace fenced preset"
description: "An operator that reconciles KubernetesSolr resources only in an explicit namespace list, sized for a governed platform: two leader-elected replicas (one active, one warm standby) and declared..."
type: "preset"
rank: "03"
presetSlug: "03-namespace-fenced"
componentSlug: "solr-operator"
componentTitle: "Solr Operator"
provider: "kubernetes"
icon: "package"
order: 3
---

# Namespace fenced preset

An operator that reconciles KubernetesSolr resources only in an
explicit namespace list, sized for a governed platform: two
leader-elected replicas (one active, one warm standby) and declared
resource requests/limits so the pods fit namespace quotas. This is
the multi-tenant posture — several fenced installs can share a
cluster, each owning its own namespaces.

Do not use this when one team administers all Solr on the cluster:
the standard preset's cluster-wide watch is simpler, and a fence you
do not need turns into a debugging trap — a KubernetesSolr declared
outside the watched list is silently ignored, which presents as a
cluster that never comes up. Also remember the bundled
zookeeper-operator still installs here and IS cluster-scoped; on a
cluster that already runs one, combine this preset's fence with the
existing-zookeeper-operator preset's `zookeeper_operator` block.

Change first: the `watch_namespaces` list to the namespaces your Solr
clusters will live in, and drop `replicas` back to 1 if a warm
standby is more operator than your environment needs.

See [03-namespace-fenced.yaml](./03-namespace-fenced.yaml) for the
manifest.
