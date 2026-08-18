---
title: "Terraform Parity"
description: "Measured parity of the DigitalOcean catalog against the pinned Terraform provider"
icon: "check-circle"
order: 90
---

<!-- GENERATED FILE -- DO NOT EDIT.
     Rendered from the committed provider schemas, the kind registry, the
     Terraform modules, the per-kind provider-parity manifests, the
     dispositions ledger, and the E2E profiles.
     parameters: provider=digitalocean ga-schema=digitalocean
     Regenerate: make generate-provider-parity-report -->

# DigitalOcean Terraform Parity

This catalog is **built for 100% Terraform parity**: every configurable
argument of the pinned Terraform provider is representable through a kind,
and every provider resource carries exactly one recorded disposition --
omission is a decision, never an accident. This page is the measurement,
generated from the same accounting that gates the repository's CI. It makes
no achieved-parity claim: a kind counts as PROVEN only when live end-to-end
runs pass on both IaC engines, and the tables below show exactly how far
that has progressed.

## Measurement baseline

| | |
|---|---|
| Provider schema (parity baseline) | `digitalocean@2.99.1` |
| Kinds in the catalog | 20 |
| Distinct provider resources consumed | 23 |
| Spec fields authored across all kinds | 580 |
| Module pins on `digitalocean` | `~> 2.99` × 20 |

The GA provider is the parity baseline. Capability that exists only in a
secondary channel (for Google, the `google-beta` provider) enters per kind
through an explicitly enumerated admission list, never wholesale.

## Depth: per-kind accounting

Every configurable, non-deprecated provider argument of a kind's consumed
resources must be matched to a spec field, mapped by recorded judgment, or
excluded with a recorded reason -- and every spec field must reach provider
surface. **Accounted** means both directions hold with zero unexplained
gaps. **Proven** means live end-to-end runs passed on both IaC engines.

**20 of 20 kinds are at total accounting; 0 proven live.**

| Kind | Provider args | Matched | Mapped | Excluded | Open gaps | Accounted | Proven |
|---|---|---|---|---|---|---|---|
| DigitalOceanApp | 292 | 0 | 282 | 10 | 0 | ✅ | — |
| DigitalOceanBucket | 28 | 6 | 19 | 3 | 0 | ✅ | — |
| DigitalOceanCertificate | 6 | 0 | 5 | 1 | 0 | ✅ | — |
| DigitalOceanContainerRegistry | 6 | 2 | 3 | 1 | 0 | ✅ | — |
| DigitalOceanDatabaseCluster | 19 | 14 | 5 | 0 | 0 | ✅ | — |
| DigitalOceanDatabaseConnectionPool | 6 | 4 | 2 | 0 | 0 | ✅ | — |
| DigitalOceanDatabaseDb | 2 | 0 | 2 | 0 | 0 | ✅ | — |
| DigitalOceanDatabaseFirewall | 3 | 0 | 2 | 1 | 0 | ✅ | — |
| DigitalOceanDatabaseReplica | 7 | 4 | 3 | 0 | 0 | ✅ | — |
| DigitalOceanDatabaseUser | 7 | 1 | 6 | 0 | 0 | ✅ | — |
| DigitalOceanDnsRecord | 10 | 9 | 1 | 0 | 0 | ✅ | — |
| DigitalOceanDnsZone | 12 | 8 | 3 | 1 | 0 | ✅ | — |
| DigitalOceanDroplet | 21 | 16 | 4 | 1 | 0 | ✅ | — |
| DigitalOceanFirewall | 17 | 2 | 15 | 0 | 0 | ✅ | — |
| DigitalOceanFunction | 292 | 0 | 36 | 256 | 0 | ✅ | — |
| DigitalOceanKubernetesCluster | 47 | 31 | 15 | 1 | 0 | ✅ | — |
| DigitalOceanKubernetesNodePool | 13 | 8 | 5 | 0 | 0 | ✅ | — |
| DigitalOceanLoadBalancer | 46 | 31 | 15 | 0 | 0 | ✅ | — |
| DigitalOceanVolume | 8 | 5 | 3 | 0 | 0 | ✅ | — |
| DigitalOceanVpc | 4 | 2 | 1 | 1 | 0 | ✅ | — |

## Breadth: every GA resource, one disposition

All resources of `digitalocean@2.99.1` land in exactly one class:

| Disposition | Resources | Meaning |
|---|---|---|
| Modeled | 23 | consumed by a kind's Terraform module today |
| IAM-covered | 0 | per-resource IAM member/binding/policy triplets, covered by the owning kinds' additive `iam_members` fields |
| Composed | 20 | capability covered through an existing kind's surface rather than a kind of its own |
| Planned | 17 | judged to be covered by a planned kind or planned composition, not built yet |
| Deferred | 17 | deliberately not offered, each with the recorded reason |
| Excluded as deprecated | 2 | deprecated or superseded provider surface |
| **Total** | **79** | |

## The enumerated record

The full per-resource record, so the accounting above is verifiable
rather than trusted.

### Modeled (23)

| Resource | Consuming kinds |
|---|---|
| `digitalocean_app` | consumed by DigitalOceanApp, DigitalOceanFunction |
| `digitalocean_certificate` | consumed by DigitalOceanCertificate |
| `digitalocean_container_registry` | consumed by DigitalOceanContainerRegistry |
| `digitalocean_container_registry_docker_credentials` | consumed by DigitalOceanContainerRegistry |
| `digitalocean_database_cluster` | consumed by DigitalOceanDatabaseCluster |
| `digitalocean_database_connection_pool` | consumed by DigitalOceanDatabaseConnectionPool |
| `digitalocean_database_db` | consumed by DigitalOceanDatabaseDb |
| `digitalocean_database_firewall` | consumed by DigitalOceanDatabaseFirewall |
| `digitalocean_database_replica` | consumed by DigitalOceanDatabaseReplica |
| `digitalocean_database_user` | consumed by DigitalOceanDatabaseUser |
| `digitalocean_domain` | consumed by DigitalOceanDnsZone |
| `digitalocean_droplet` | consumed by DigitalOceanDroplet |
| `digitalocean_firewall` | consumed by DigitalOceanFirewall |
| `digitalocean_kubernetes_cluster` | consumed by DigitalOceanKubernetesCluster |
| `digitalocean_kubernetes_node_pool` | consumed by DigitalOceanKubernetesNodePool |
| `digitalocean_loadbalancer` | consumed by DigitalOceanLoadBalancer |
| `digitalocean_record` | consumed by DigitalOceanDnsRecord, DigitalOceanDnsZone |
| `digitalocean_spaces_bucket` | consumed by DigitalOceanBucket |
| `digitalocean_spaces_bucket_cors_configuration` | consumed by DigitalOceanBucket |
| `digitalocean_spaces_bucket_logging` | consumed by DigitalOceanBucket |
| `digitalocean_spaces_bucket_policy` | consumed by DigitalOceanBucket |
| `digitalocean_volume` | consumed by DigitalOceanVolume |
| `digitalocean_vpc` | consumed by DigitalOceanVpc |

### Composed (20)

| Resource | Recorded reason |
|---|---|
| `digitalocean_container_registries` | the plural resource is the newer multi-registry API for the same product surface; DigitalOceanContainerRegistry absorbs it as an implementation upgrade rather than a second registry kind |
| `digitalocean_database_advanced_mysql_config` | extended per-cluster MySQL settings singleton -- exactly one per cluster, lifecycle identical to the parent; folds into DigitalOceanDatabaseCluster as per-engine settings |
| `digitalocean_database_advanced_postgresql_config` | extended per-cluster PostgreSQL settings singleton -- exactly one per cluster, lifecycle identical to the parent; folds into DigitalOceanDatabaseCluster as per-engine settings |
| `digitalocean_database_kafka_config` | per-cluster Kafka broker settings singleton -- exactly one per cluster, lifecycle identical to the parent; folds into DigitalOceanDatabaseCluster as per-engine settings |
| `digitalocean_database_logsink_rsyslog` | same ship-cluster-logs concept as the opensearch logsink, different transport; folds into the planned DigitalOceanDatabaseLogsink kind's oneof destination rather than duplicating the kind |
| `digitalocean_database_mongodb_config` | per-cluster MongoDB settings singleton -- exactly one per cluster, lifecycle identical to the parent; folds into DigitalOceanDatabaseCluster as per-engine settings |
| `digitalocean_database_mysql_config` | per-cluster MySQL settings singleton -- exactly one per cluster, lifecycle identical to the parent; folds into DigitalOceanDatabaseCluster as per-engine settings |
| `digitalocean_database_opensearch_config` | per-cluster OpenSearch settings singleton -- exactly one per cluster, lifecycle identical to the parent; folds into DigitalOceanDatabaseCluster as per-engine settings |
| `digitalocean_database_postgresql_config` | per-cluster PostgreSQL settings singleton -- exactly one per cluster, lifecycle identical to the parent; folds into DigitalOceanDatabaseCluster as per-engine settings |
| `digitalocean_database_redis_config` | per-cluster Redis/Caching settings singleton -- exactly one per cluster, lifecycle identical to the parent; folds into DigitalOceanDatabaseCluster as per-engine settings |
| `digitalocean_database_valkey_config` | per-cluster Valkey settings singleton -- exactly one per cluster, lifecycle identical to the parent; folds into DigitalOceanDatabaseCluster as per-engine settings |
| `digitalocean_dedicated_inference_token` | access tokens are many-per-endpoint credential rows tightly coupled to the endpoint they authorize; folds into the planned DigitalOceanDedicatedInference kind |
| `digitalocean_nfs_access_point` | export paths and access policies are many-per-share config rows, meaningless without the share; folds into the planned DigitalOceanNfs kind |
| `digitalocean_nfs_attachment` | additional-VPC attachments are membership rows on the share; folds into the planned DigitalOceanNfs kind |
| `digitalocean_project_resources` | membership assignment rows (URN lists) with no meaning outside the project; the planned DigitalOceanProject kind carries the membership list, and Planton-managed resources prefer declaring their own project FK |
| `digitalocean_reserved_ip_assignment` | attach/detach is a nullable droplet FK on the planned DigitalOceanReservedIp kind, not a separate object |
| `digitalocean_reserved_ipv6` | same reserved-IP concept, v6 flavor; a separate kind would be a bundled duplicate -- folds into the planned DigitalOceanReservedIp kind |
| `digitalocean_reserved_ipv6_assignment` | same reasoning as the v4 assignment -- a nullable droplet FK on the planned DigitalOceanReservedIp kind |
| `digitalocean_uptime_alert` | alerts are many-per-check rows that cannot exist without the check and are always managed alongside it; the planned DigitalOceanUptimeCheck kind carries an alerts list |
| `digitalocean_volume_attachment` | attachment is a nullable droplet FK on DigitalOceanVolume (the droplet kind's volume_ids covers it today); attach/detach is a field change, not a separate object lifecycle worth a kind |

### Planned (17)

| Resource | Recorded reason |
|---|---|
| `digitalocean_cdn` | judged as a planned DigitalOceanCdn kind (CDN endpoint fronting a Spaces bucket origin, certificate reference for custom domains) |
| `digitalocean_custom_image` | judged as a planned DigitalOceanCustomImage kind (imported images are first-class and referenced by droplets and autoscale pools) |
| `digitalocean_database_kafka_schema_registry` | judged as a planned DigitalOceanDatabaseKafkaSchema kind (despite the resource name, each instance manages one registered schema -- many-per-cluster rows managed independently, like topics; FK cluster) |
| `digitalocean_database_kafka_topic` | judged as a planned DigitalOceanDatabaseKafkaTopic kind (topics are many-per-cluster with independent lifecycle and per-topic config; FK cluster) |
| `digitalocean_database_logsink_opensearch` | anchor of a planned DigitalOceanDatabaseLogsink kind modeling ship-cluster-logs-somewhere with a oneof destination config (FK cluster; optionally a Planton-managed OpenSearch cluster) |
| `digitalocean_dedicated_inference` | judged as a planned DigitalOceanDedicatedInference kind (GPU-backed inference endpoints with model deployments and accelerator scaling -- first-class, referenceable infrastructure; FK project) |
| `digitalocean_droplet_autoscale` | judged as a planned DigitalOceanDropletAutoscalePool kind (a pool with its own scaling config and lifecycle -- the closest thing DigitalOcean has to a managed instance group; FK vpc, ssh keys, project) |
| `digitalocean_monitor_alert` | judged as a planned DigitalOceanMonitorAlert kind (standalone alert policy over infrastructure metrics with its own lifecycle and notification config; FK droplets or tag-based targeting) |
| `digitalocean_nfs` | anchor of a planned DigitalOceanNfs kind (the share is the parent object with size, performance tier, and primary VPC; FK vpc) |
| `digitalocean_partner_attachment` | judged as a planned DigitalOceanPartnerAttachment kind (GA but specialist private connectivity via Megaport-class providers with BGP config; first-class because nothing else could honestly absorb it; FK VPCs) |
| `digitalocean_project` | judged as a planned DigitalOceanProject kind (the account's organizational container; nearly every other resource can carry a project reference, so this kind unlocks FK fields across the catalog) |
| `digitalocean_reserved_ip` | anchor of a planned DigitalOceanReservedIp kind carrying an IP version field and an optional droplet attachment (the v4/v6 and create/assign splits are provider-API artifacts, not distinct product concepts; FK droplet, project) |
| `digitalocean_spaces_key` | judged as a planned DigitalOceanSpacesKey kind (access keys have independent create/rotate lifecycles and are the credential workloads actually reference; FK buckets via per-bucket grants) |
| `digitalocean_ssh_key` | judged as a planned DigitalOceanSshKey kind (independent lifecycle, referenced by droplets and autoscale pools; a hard prerequisite for droplet-centric workflows) |
| `digitalocean_uptime_check` | judged as a planned DigitalOceanUptimeCheck kind (a check on an external endpoint is the parent object, referenced by its alerts) |
| `digitalocean_vector_database` | judged as a planned DigitalOceanVectorDatabase kind (Weaviate-powered managed cluster, managed independently from standard database clusters; FK project, vpc) |
| `digitalocean_vpc_peering` | judged as a planned DigitalOceanVpcPeering kind (peerings connect two first-class VPCs and have their own lifecycle; FK two VPCs) |

### Deferred (17)

| Resource | Recorded reason |
|---|---|
| `digitalocean_byoip_prefix` | bring-your-own-IP needs an out-of-band provisioning workflow (signed request, DigitalOcean-side verification, later advertisement) and serves a small compliance-driven audience -- deferred until the workflow can be represented honestly |
| `digitalocean_database_online_migration` | action-style virtual resource (start/stop a replication job that runs up to two weeks), not durable declarative state -- gated on an action/workflow surface |
| `digitalocean_droplet_snapshot` | a point-in-time capture is operationally an action, not durable declared state -- gated with volume_snapshot on an action/snapshot-policy surface |
| `digitalocean_gradientai_agent` | core object of the GradientAI family, deferred wholesale: the provider surface is visibly immature (copy-paste description bugs in its own schemas, churning release to release) -- revisit as one wave when it settles |
| `digitalocean_gradientai_agent_knowledge_base_attachment` | attachment row of a deferred parent -- the GradientAI family is deferred wholesale on documented provider immaturity |
| `digitalocean_gradientai_agent_route` | agent-to-agent routing rows whose schema carries placeholder-grade descriptions -- the GradientAI family is deferred wholesale on documented provider immaturity |
| `digitalocean_gradientai_custom_model` | depends on the deferred agent surface -- the GradientAI family is deferred wholesale on documented provider immaturity |
| `digitalocean_gradientai_function` | agent function attachment, deferred with its parent -- the GradientAI family is deferred wholesale on documented provider immaturity |
| `digitalocean_gradientai_indexing_job_cancel` | action-style (cancel a running indexing job) -- would be deferred even in a mature family; the GradientAI family is deferred wholesale besides |
| `digitalocean_gradientai_knowledge_base` | likely second anchor when the GradientAI family matures; deferred with the family today on documented provider immaturity |
| `digitalocean_gradientai_knowledge_base_data_source` | data-source rows of a deferred parent -- the GradientAI family is deferred wholesale on documented provider immaturity |
| `digitalocean_gradientai_openai_api_key` | credential row for a deferred family -- the GradientAI family is deferred wholesale on documented provider immaturity |
| `digitalocean_nfs_snapshot` | point-in-time capture, action-style -- same reasoning as droplet and volume snapshots; gated on an action/snapshot-policy surface |
| `digitalocean_spaces_bucket_object` | uploading object content is data-plane work, not infrastructure; a declarative kind for file contents is marginal and invites abuse as a deployment mechanism |
| `digitalocean_tag` | DigitalOcean creates tags implicitly when any resource declares them; a standalone name-reservation kind is marginal -- revisit if tag-targeted references (firewalls, load balancers, monitor alerts) prove to need a first-class handle |
| `digitalocean_volume_snapshot` | point-in-time capture, action-style -- gated with droplet_snapshot on an action/snapshot-policy surface |
| `digitalocean_vpc_nat_gateway` | the provider docs mark it currently in Private Preview -- gated on GA; product value is high, promote to P1 the moment the gate lifts |

### Excluded as deprecated (2)

| Resource | Recorded reason |
|---|---|
| `digitalocean_floating_ip` | docs-deprecated in favor of reserved_ip (the schema DeprecationMessage is commented out in provider code, so the checker cannot compute it) |
| `digitalocean_floating_ip_assignment` | docs-deprecated in favor of reserved_ip_assignment (the schema DeprecationMessage is commented out in provider code, so the checker cannot compute it) |
