# Kubernetes search and graph: OpenSearch operator + cluster (renamed from Elastic), Solr operator + SolrCloud, and Neo4j — five kinds at full depth

## What changed

- **KubernetesOpenSearchOperator (renamed from KubernetesElasticOperator,
  916)** — the catalog's search-engine operator is now the OpenSearch
  Kubernetes operator, installed from the served `opensearch-operator`
  chart 2.8.0 (matching app 2.8.0 — the newest pairing whose default
  images are all stable releases). The kind was renamed because the
  catalog distributes only permissively licensed software: Elastic's
  ECK operator is Elastic-License-2.0 and Elasticsearch itself is
  ELv2/SSPL-dual (AGPL only in narrow cases), while OpenSearch and its
  operator are Apache-2.0 end to end. The enum number and existing
  resources' identity are preserved; the id prefix is now `k8sosop`.
  The chart TEMPLATES its CRDs with no keep switch, so both modules own
  the CRD lifecycle directly and retain the CRDs on uninstall —
  uninstalling the operator can never cascade-delete every
  OpenSearchCluster in the cluster. Two chart defects are corrected in
  the modules: the optional kube-rbac-proxy sidecar's default image
  points at a registry path that was deleted upstream (re-pointed at
  the maintainer's quay.io repository at the chart's pinned tag), and
  the chart's fullname is not truncated when building child resource
  names (the release fullname is pinned to `metadata.name`, with the
  name budget documented on the spec).
- **KubernetesOpenSearch (renamed from KubernetesElasticsearch, 917)** —
  a full-depth OpenSearchCluster: typed node pools with roles,
  persistence, and per-pool scheduling; TLS on both the transport and
  HTTP layers (generated or user-supplied certificates); OpenSearch
  Dashboards folded into the same kind; keystore entries from Secrets;
  snapshot repositories (s3/gcs/azure); monitoring; additional volumes.
  Validation enforces a floor of TWO cluster-manager-eligible replicas:
  the operator bootstraps a new cluster through a temporary bootstrap
  manager and removes it once the cluster forms, and a lone manager
  node cannot re-form quorum after that handoff — a single-manager
  cluster accepts no writes, permanently. The E2E suite proves
  index-and-search behavior on every lane and full data durability
  through a live node loss (a replicated document is served DURING the
  outage and after recovery), on both engines.
- **KubernetesSolrOperator (920)** — rebuilt from a minimal container
  shell onto the full `solr-operator` chart 0.9.1 surface: watch-scope
  fencing, the bundled zookeeper-operator toggle, mTLS client identity
  for operator-to-Solr calls, metrics and leader election. The chart
  templates its CRDs (SolrCloud, SolrBackup, SolrPrometheusExporter and
  the zookeeper-operator's ZookeeperCluster), so the modules own them
  with keep-on-uninstall, same as the OpenSearch operator.
- **KubernetesSolr (921)** — a full-depth SolrCloud: operator-provided
  or external ZooKeeper, persistent/ephemeral storage, basic-auth
  bootstrap (with the two-secret contract documented: the exported
  basic-auth Secret carries the operator's read-only user; admin
  credentials live in the security-bootstrap Secret), TLS keystores,
  s3/gcs/volume backup repositories, managed rolling updates, scaling
  policies, and the operator's own exposure methods. Both engines
  declare BACKGROUND deletion propagation on the SolrCloud resource:
  a foreground cascade deadlocks against the zookeeper-operator, which
  keeps reconciling the ZooKeeper children while the parent waits on
  its finalizer — with background propagation the operator's own
  teardown logic runs and deletion completes cleanly. The live proof
  creates a collection, indexes a document, and queries it back as the
  admin user.
- **KubernetesNeo4j (922)** — the catalog's graph database, rebuilt at
  full depth on the official `neo4j` Helm chart (2026.6.0): edition
  modeling with the enterprise license-acceptance gate, the NEO4J_AUTH
  Secret contract (the module materializes the auth Secret BEFORE the
  release so the chart never renders a literal password), memory
  tuning through the chart's own server.memory keys, volumes
  (dynamic/volume/selectable), APOC allowlisting, and honest exposure
  (the chart's LoadBalancer default is overridden to ClusterIP unless
  external exposure is requested). The live proof writes a node over
  Bolt, deletes the pod, and reads the node back after recovery.

## Why

Search is a foundational tier of the software catalog: OpenSearch is
the default log/search store for many stacks, Solr remains the
enterprise search workhorse, and a graph database rounds out the
retrieval story. All five kinds now meet the catalog bar: complete
typed specs authored from the pinned upstream sources, dual-engine
modules rendering byte-identical configuration, CEL validation that
encodes live-cluster behavior (not just upstream schema claims), and
end-to-end suites that prove behavior — not just installation — on
both engines.

## Impact

- Users declaring search infrastructure get Apache-2.0-licensed engines
  with no licensing ambiguity, full configuration depth, and specs
  whose field documentation teaches the operational constraints
  (manager quorum floors, secret contracts, name budgets) at the point
  of configuration.
- Existing KubernetesElasticOperator/KubernetesElasticsearch resources
  keep their enum identity through the rename; the API surface is new,
  matching the OpenSearch operator's CRDs.
- Operator uninstalls can no longer destroy workload data: all three
  operator kinds in this set retain their CRDs on uninstall, and
  operator-owned custom resources delete through the operator's own
  teardown path.
