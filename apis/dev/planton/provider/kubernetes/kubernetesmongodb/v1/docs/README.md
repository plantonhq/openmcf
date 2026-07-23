# MongoDB on Kubernetes: Choosing the Right Path

## Introduction

For years, the conventional wisdom was clear: don't run stateful databases on Kubernetes. The logic seemed sound—Kubernetes was designed for stateless workloads, and databases need stable identities, persistent storage, and complex orchestration that goes far beyond what standard Kubernetes controllers provide.

That wisdom is now outdated.

Modern Kubernetes has evolved to provide the primitives needed for stateful workloads (StatefulSets, persistent volumes, stable network identities), and more importantly, the **Operator pattern** has emerged to encode database-specific operational knowledge into self-healing controllers. Today, running MongoDB on Kubernetes isn't just viable—when done correctly, it can provide better reliability, easier scaling, and more consistent operations than traditional deployment methods.

But here's the critical insight: not all deployment methods are created equal. The path from "Hello World" to production-grade MongoDB on Kubernetes passes through several maturity levels, each solving different problems. This document explains what those levels are, which approaches are production-ready, and why Planton chose the Percona Operator as the engine behind **KubernetesMongodb**.

## The Evolution: From Anti-Patterns to Production

Understanding how to deploy MongoDB on Kubernetes requires understanding the maturity spectrum—what works, what doesn't, and why.

### Level 0: The Anti-Pattern (Deployments and ReplicaSets)

**What it is:** Using standard Kubernetes `Deployment` or `ReplicaSet` controllers to run MongoDB pods.

**Why it fails:** Kubernetes excels at managing **stateless** applications where pods are fungible—if one fails, the controller simply creates a new one with a fresh identity and empty storage. This is fundamentally incompatible with a stateful database:

1. **Loss of Identity:** A new pod gets a new hostname and IP address. For a MongoDB replica set, this new pod is an unknown entity, not the member that disappeared. The cluster loses quorum and requires manual intervention.

2. **Loss of State:** A new pod gets fresh, empty storage. The data from the failed pod is gone.

This is an **impedance mismatch**: Kubernetes abstracts away state and identity, while a database's entire value *is* its state and identity.

**Verdict:** Never use this approach for any database workload.

### Level 1: The Foundation (StatefulSets)

**What it is:** Using Kubernetes `StatefulSet` controllers, which provide stable network identities (predictable hostnames like `mongo-0`, `mongo-1`) and stable, persistent storage that survives pod restarts.

**What it solves:** StatefulSets fix the anti-pattern by ensuring that when `mongo-0` fails and restarts, it reconnects to its exact same persistent volume with its same hostname. The database cluster recognizes it as the same member returning, not a new entity.

**What it doesn't solve:** StatefulSets are just *primitives*. They don't understand MongoDB's replica set topology, election logic, or operational requirements. This creates a "Day 2" operations gap:

- **No application logic:** StatefulSets don't know how to initialize a replica set, manage elections, or perform graceful primary step-downs.
- **Manual scaling:** Scaling from 3 to 4 pods doesn't automatically run `rs.add()` to add the new member to the MongoDB cluster.
- **Orphaned storage:** StatefulSets don't delete PersistentVolumeClaims when scaled down, leaving expensive volumes stranded.

**Verdict:** Necessary foundation, but insufficient alone. You need automation on top of this.

### Level 2: The "Helm Trap" (Helm Charts)

**What it is:** Using Helm charts (like the popular Bitnami MongoDB chart) to package and deploy MongoDB on Kubernetes.

**What it solves:** Helm simplifies "Day 1" installation by templating all the required Kubernetes YAML (StatefulSets, Services, Secrets, ConfigMaps) into a single, configurable package. A simple `helm install` gives you a running MongoDB cluster in minutes.

**The trap:** This simplicity creates a false sense of production-readiness. Teams achieve a trivial installation and only hit the wall later, when they need "Day 2" operations. Helm is a **client-side, stateless templating tool**—it has no runtime controller and thus cannot automate:

- **Backups:** The Bitnami chart requires manual `mongodump` commands or separate tools like Velero.
- **Failover:** Helm plays no role in runtime operations; MongoDB's native replica set handles this, but without advanced automation.
- **Safe upgrades:** `helm upgrade` just re-templates and applies YAML. It has no application-specific knowledge to perform a safe, rolling upgrade (e.g., upgrade secondaries first, then step down and upgrade the primary).
- **Resource cleanup:** Helm doesn't track PVCs created by StatefulSets, leaving orphaned volumes after `helm uninstall`.

**Verdict:** Acceptable for development, testing, and non-critical applications. Not a production-grade solution for mission-critical workloads. The high operational cost and manual risk make this a technical debt time bomb.

### Level 3: The Production Solution (Kubernetes Operators)

**What it is:** A Kubernetes Operator is a custom controller that extends Kubernetes with application-specific operational knowledge. For MongoDB, this means a controller that watches custom resources (like `kind: PerconaServerMongoDB`) and continuously reconciles the actual state of the cluster with the desired state you declare.

**What it solves:** Everything. An Operator encodes the domain expertise of a human MongoDB administrator into software:

- **Automated scaling:** Change a replica set's size from 3 to 5, and the Operator scales the StatefulSet, waits for new pods to be ready, then automatically joins them to the replica set.
- **Automated backups:** Schedule logical, physical, and incremental backups with point-in-time recovery (PITR) via oplog streaming to S3, GCS, or Azure Blob Storage.
- **Safe upgrades:** The Operator orchestrates rolling upgrades, upgrading secondaries first, then stepping down and upgrading the primary.
- **Self-healing:** If a pod fails, the Operator doesn't just restart it—it ensures the MongoDB cluster recognizes it and maintains quorum.

**The catch:** Not all operators are created equal. Licensing, feature completeness, and openness vary dramatically.

**Verdict:** This is the only production-ready approach for running MongoDB on Kubernetes. The question isn't "Operator or not?"—it's "*Which* Operator?"

## Operator Comparison: The Licensing Minefield

For an open-source infrastructure framework like Planton, choosing the right operator isn't just about features—it's about legal soundness, long-term viability, and avoiding "open core" traps where essential production features are paywalled.

### The Contenders

1. **MongoDB Community Operator** — Official, but incomplete
2. **MongoDB Enterprise Operator** — Official, but commercial
3. **Percona Operator for MongoDB** — 100% open source, feature-complete
4. **KubeDB MongoDB Operator** — Open core with restrictive paywalls

### Feature and Licensing Matrix

| Feature             | MongoDB Community | MongoDB Enterprise | Percona Operator | KubeDB |
|:--------------------|:------------------|:-------------------|:-----------------|:-------|
| **Code License**    | Apache 2.0        | Apache 2.0         | **Apache 2.0**   | Open Core |
| **DB License**      | SSPL ⚠️           | Commercial ⚠️      | **SSPL-Free** ✅ | SSPL ⚠️ |
| **Replica Sets**    | ✅ Free           | ✅ Free            | **✅ Free**      | ✅ Free |
| **Sharding**        | ❌ No             | 💰 Paid            | **✅ Free**      | ✅ Free |
| **Auto Backups**    | ❌ No             | 💰 Paid            | **✅ Free**      | 💰 Paid |
| **PITR**            | ❌ No             | 💰 Paid            | **✅ Free**      | 💰 Paid |
| **Safe Upgrades**   | Manual            | 💰 Paid            | **✅ Free**      | Manual |
| **Storage Scaling** | ❌ No             | 💰 Paid            | **✅ Free**      | 💰 Paid |
| **TLS/SSL**         | ✅ Free           | ✅ Free            | **✅ Free**      | 💰 Paid |
| **GUI**             | No                | Ops Mgr (Paid)     | **PMM (Free)**   | Paid |

### The MongoDB "Open Core Trap"

The MongoDB Community Operator is intentionally crippled—it lacks sharding and automated backups, which are classified as "Enterprise" features. This isn't an accident; it's product segmentation designed to upsell users to the expensive MongoDB Enterprise Advanced subscription.

More concerning for Planton: MongoDB Community Server is licensed under the **Server Side Public License (SSPL)**, a viral license that's not OSI-approved. The SSPL requires any entity offering the software *as a service* to open-source all of its own management and service-delivery code. Planton's KubernetesMongodb API—which abstracts and manages MongoDB deployment—could legally be considered such a service, creating existential licensing risk for the entire framework.

### The KubeDB Dead End

KubeDB represents the worst of both worlds:

1. **Legal Risk:** It deploys standard SSPL-licensed MongoDB, inheriting the full legal exposure.
2. **Extreme Open Core:** The free version is a 30-day trial. Even basic features like TLS/SSL, automated backups, and scaling are paywalled.

This is a non-starter for any serious deployment.

### The Percona Advantage

The Percona Operator for MongoDB is the only solution that avoids all legal and technical traps:

**Legally Sound:**
- The operator code is 100% Apache 2.0 licensed
- It deploys **Percona Server for MongoDB**, a fully compatible, enhanced distribution that is **SSPL-free**
- No commercial licenses, no viral restrictions, no legal risk

**Production-Ready and Free:**
- ✅ Automated backups with PITR via integrated Percona Backup for MongoDB
- ✅ Automated sharding support
- ✅ Safe, orchestrated rolling upgrades (SmartUpdate)
- ✅ Automated storage and compute scaling (horizontal and vertical)
- ✅ Full TLS/SSL support with automated certificate management
- ✅ Integrated monitoring with Percona Monitoring and Management (PMM)

This is not a "community edition" with missing features. This is a complete, enterprise-grade operator with all production capabilities, provided entirely free and open.

## Planton's Choice: Percona Operator

**KubernetesMongodb** is built on the Percona Operator for MongoDB. This is the only choice that aligns with the philosophy of an open-source infrastructure framework:

1. **No Legal Risk:** Completely Apache 2.0 licensed with an SSPL-free database distribution.
2. **No Paywalls:** All production features—PITR backups, sharding, scaling—are free.
3. **Production-Proven:** Used by organizations that require enterprise-grade reliability without enterprise-grade costs.
4. **Best-Practice Defaults:** The spec enforces safe postures (three-member replica sets, TLS, anti-affinity) and makes every unsafe deviation an explicit opt-in.

### How It Works

The engine and the databases are two components, composed deliberately:

1. **KubernetesPerconaMongoOperator** installs the operator from the official Helm chart — a cluster prerequisite deployed first. By default the operator watches its OWN namespace, so databases live beside it (or the watch is widened).
2. Each **KubernetesMongodb** resource renders one `psmdb.percona.com/v1` PerconaServerMongoDB custom resource — replica sets, optional sharding, TLS, users, backups.
3. The Percona Operator takes over, managing the full lifecycle of the MongoDB cluster: elections, failover, rolling upgrades, backup scheduling, user reconciliation.

You get the simplicity of a "Day 1" declarative experience with the power and reliability of a "Day 2" Operator solution—without the licensing traps or operational complexity.

### Example Configurations

Every example below is the real spec surface (snake_case fields; `namespace` accepts a literal `value:` or a KubernetesNamespace reference).

**Development (single member, unsafe by declaration):**

```yaml
apiVersion: kubernetes.planton.dev/v1
kind: KubernetesMongodb
metadata:
  name: dev-mongo
spec:
  namespace:
    value: percona-mongo # where the operator watches
  replica_sets:
    - name: rs0
      size: 1
      storage:
        size: 5Gi
  unsafe:
    replset_size: true # the operator rejects 1-member sets without this
```

**Staging (HA, 3 members, explicit resources):**

```yaml
apiVersion: kubernetes.planton.dev/v1
kind: KubernetesMongodb
metadata:
  name: staging-mongo
spec:
  namespace:
    value: percona-mongo
  replica_sets:
    - name: rs0
      size: 3 # a majority survives one loss — automated failover
      storage:
        size: 50Gi
      resources:
        requests:
          cpu: "1"
          memory: 4Gi
        limits:
          cpu: "2"
          memory: 4Gi # WiredTiger sizes its cache from the limit
```

**Production (HA, backups + PITR, required TLS, declared user):**

```yaml
apiVersion: kubernetes.planton.dev/v1
kind: KubernetesMongodb
metadata:
  name: prod-mongo
spec:
  namespace:
    value: percona-mongo
  replica_sets:
    - name: rs0
      size: 3
      storage:
        size: 500Gi
      resources:
        requests:
          cpu: "4"
          memory: 16Gi
        limits:
          cpu: "8"
          memory: 16Gi
      pod_disruption_budget:
        max_unavailable: 1
      scheduling:
        anti_affinity_topology_key: topology.kubernetes.io/zone
  tls:
    mode: requireTLS
  users:
    - name: app
      roles:
        - name: readWrite
          db: app
  backup:
    storages:
      - name: primary
        s3:
          bucket: my-prod-mongo-backups
          region: us-east-1
          access_keys:
            access_key_id: AKIAEXAMPLE
            secret_access_key: <backup-user-secret-key>
    tasks:
      - name: nightly
        schedule: "0 0 * * *" # daily at midnight, five-field cron
        storage_name: primary
        keep: 5
    pitr:
      enabled: true # continuous oplog archiving between backups
```

## Production Best Practices

Running MongoDB on Kubernetes in production requires more than just correct configuration—it requires codifying best practices into your infrastructure.

### High Availability

**Replica Sets:** Always use at least 3 members in production to maintain quorum if one member fails. The spec enforces this: smaller sets require the explicit `unsafe.replset_size` opt-in.

**Anti-Affinity (Critical):** The operator spreads members across nodes by default (`kubernetes.io/hostname`). For true resilience, set `scheduling.anti_affinity_topology_key: topology.kubernetes.io/zone` so the cluster survives zone-level failures, and declare a `pod_disruption_budget` so voluntary disruptions never break the majority.

**Arbiters over even numbers:** Even member counts waste a vote; declare the `arbiter` block (votes, holds no data) instead of a fourth data member.

### Backup and Disaster Recovery

**Never rely on manual backups.** Production systems require automated, application-consistent backups — the `backup` block with named storages and scheduled tasks.

**Point-in-Time Recovery (PITR):** The gold standard is continuous oplog archiving, which enables recovery to any specific moment (e.g., 5 minutes before a bad query was executed). Set `pitr.enabled: true`; oplog chunks land on the main storage, and PITR needs at least one completed base backup to be meaningful.

**Test your restores.** Backups you've never tested are backups that don't work.

### Security Hardening

1. **Authentication:** Always on — the operator manages the system users (SCRAM) and their `<name>-secrets` Secret; there is no unauthenticated posture.
2. **Authorization (RBAC):** Declare application users with least-privilege `roles`; applications should never connect as the database admin.
3. **Encryption-in-Transit (TLS):** The default is `preferTLS`; production should set `tls.mode: requireTLS`, with certificates issued by a cert-manager (Cluster)Issuer via `tls.issuer` for an organization-trusted chain. Disabling TLS requires the explicit `unsafe.tls` opt-in.
4. **Encryption-at-Rest:** Use a StorageClass that provisions encrypted volumes (e.g., cloud provider encryption) via `storage.storage_class`.
5. **Network Isolation:** Use Kubernetes `NetworkPolicy` to enforce zero-trust networking. MongoDB pods should only accept traffic on port 27017 from authorized application pods.

### Monitoring

Deploy the Prometheus `mongodb_exporter` to scrape metrics. Key metrics to track:

- **Health:** `mongodb_up`
- **Performance:** `mongodb_op_counters_total` (inserts, queries, updates, deletes), query latency
- **Resource Usage:** `mongodb_connections`, `mongodb_memory_usage_bytes`
- **Replication:** Oplog lag, replication health

The Percona Operator integrates natively with Percona Monitoring and Management (PMM) for comprehensive, pre-built dashboards. The fluent-bit log-collector sidecar (the `log_collector` block — off unless declared) ships mongod logs alongside.

### Resource Planning

**Storage:** Always use a StorageClass backed by SSDs. Don't rely on default provisioners that might use slow magnetic disks. Sizes can only grow — the operator applies grows to live PVCs and rejects shrinks, so starting lean is safe.

**CPU and Memory:** *Always* set `resources` with both requests and limits. WiredTiger sizes its cache from the memory limit, and failure to bound memory is the #1 cause of pod evictions and production instability.

### Common Production Pitfalls

1. **Skipping resource requests/limits:** Leads to "noisy neighbor" problems and OOMKilled pods.
2. **The watch-scope surprise:** The operator watches its own namespace by default — a database declared in another namespace is silently never reconciled. Install the operator beside the databases or widen its watch.
3. **Underprovisioning IOPS:** Choosing cheap, slow storage that can't handle database I/O demands.
4. **No backup strategy:** Discovering that backups weren't configured *after* data loss. The spec makes this a visible omission, not a silent default.
5. **The licensing trap:** Choosing MongoDB Community Operator and discovering in production that critical features are paywalled.

## Decision Framework: Which Path Should You Take?

### Choose Helm (Bitnami Chart) If:

- You're running development, testing, or CI/CD environments
- Your workload is non-critical and you're comfortable with manual operations
- You have a dedicated database team to handle backups, scaling, and upgrades

**Trade-off:** "Day 1" simplicity for "Day 2" operational burden and risk.

### Choose an Operator If:

- You're running any production-critical workload
- You want automated backups, safe upgrades, and self-healing
- You value long-term operational efficiency over short-term setup convenience

**Trade-off:** Slightly higher "Day 1" learning curve for massive "Day 2" automation gains.

### Choose Percona Operator If:

- You need a 100% open-source solution with no legal risk
- You need enterprise-grade features (PITR, sharding) without commercial licenses
- You're building infrastructure tooling (like Planton) that manages MongoDB on behalf of users

**Trade-off:** None. This is the complete solution.

### Avoid:

- **Standard Deployments/ReplicaSets:** Not viable for any database workload.
- **MongoDB Community Operator:** Missing critical production features (backups, sharding).
- **KubeDB:** Extreme open-core model with paywalls on basic features, plus SSPL legal risk.

## Migration Considerations

If you're currently running MongoDB with Helm (e.g., Bitnami chart) and want to migrate to an Operator-managed deployment, be aware: **there is no automated migration path.**

Migration requires:

1. Deploying a new, empty cluster using the Percona Operator
2. Running `mongodump` against the old cluster
3. Running `mongorestore` into the new cluster
4. Application cutover (update connection strings, restart services)

This is a high-downtime, high-risk operation. The difficulty and risk of this migration is the single strongest argument against *starting* with Helm. You trade short-term convenience for long-term technical debt.

Planton avoids this trap by starting you on an Operator-based deployment from day one, giving you the simplicity of a Helm-like API with the power of a production-grade Operator underneath.

## Conclusion: The Paradigm Shift

The database-on-Kubernetes paradigm shift is complete. The question is no longer *whether* you can run MongoDB on Kubernetes in production—it's *how you choose to do it*.

Helm charts provide a quick start but leave you with a mountain of manual operational work. "Official" operators from MongoDB Inc. trap you in an open-core licensing model that paywalls essential features and exposes you to legal risk via the SSPL.

The Percona Operator for MongoDB represents a third path: a 100% open-source, legally sound, enterprise-grade solution that automates the full lifecycle of MongoDB on Kubernetes—without paywalls, without commercial licenses, and without compromise.

Planton builds on this foundation: KubernetesPerconaMongoOperator installs the engine, KubernetesMongodb declares each cluster with best practices enforced by default and every unsafe deviation an explicit opt-in. You get the best of both worlds: the ease of "Day 1" and the power of "Day 2."

Start on the right path from the beginning. Your future self will thank you.
