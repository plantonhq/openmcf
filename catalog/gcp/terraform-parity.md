---
title: "Terraform Parity"
description: "Measured parity of the GCP catalog against the pinned Terraform provider"
icon: "check-circle"
order: 90
---

<!-- GENERATED FILE -- DO NOT EDIT.
     Rendered from the committed provider schemas, the kind registry, the
     Terraform modules, the per-kind provider-parity manifests, the
     dispositions ledger, and the E2E profiles.
     parameters: provider=gcp ga-schema=google
     Regenerate: make generate-provider-parity-report -->

# GCP Terraform Parity

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
| Provider schema (parity baseline) | `google@7.43.0` |
| Kinds in the catalog | 79 |
| Distinct provider resources consumed | 92 |
| Spec fields authored across all kinds | 2956 |
| Module pins on `google` | `~> 7.43` × 79 |

The GA provider is the parity baseline. Capability that exists only in a
secondary channel (for Google, the `google-beta` provider) enters per kind
through an explicitly enumerated admission list, never wholesale.

## Depth: per-kind accounting

Every configurable, non-deprecated provider argument of a kind's consumed
resources must be matched to a spec field, mapped by recorded judgment, or
excluded with a recorded reason -- and every spec field must reach provider
surface. **Accounted** means both directions hold with zero unexplained
gaps. **Proven** means live end-to-end runs passed on both IaC engines.

**71 of 79 kinds are at total accounting; 10 proven live.**

| Kind | Provider args | Matched | Mapped | Excluded | Open gaps | Accounted | Proven |
|---|---|---|---|---|---|---|---|
| GcpAddress | 16 | 14 | 2 | 0 | 0 | ✅ | — |
| GcpAlloydbCluster | 84 | 19 | 0 | 0 | 89 | ❌ | ✅ pulumi, terraform |
| GcpAlloydbInstance | 35 | 16 | 0 | 0 | 27 | ❌ | ✅ pulumi, terraform |
| GcpAlloydbUser | 13 | 5 | 0 | 0 | 9 | ❌ | ✅ pulumi, terraform |
| GcpArtifactRegistryRepo | 52 | 33 | 12 | 7 | 0 | ✅ | — |
| GcpBackendBucket | 29 | 22 | 6 | 1 | 0 | ✅ | — |
| GcpBackendService | 116 | 91 | 22 | 3 | 0 | ✅ | — |
| GcpBigQueryDataset | 39 | 33 | 6 | 0 | 0 | ✅ | — |
| GcpBigQueryTable | 98 | 86 | 12 | 0 | 0 | ✅ | — |
| GcpBigtableInstance | 19 | 6 | 13 | 0 | 0 | ✅ | — |
| GcpBigtableTable | 23 | 8 | 14 | 1 | 0 | ✅ | — |
| GcpCertManagerCert | 12 | 10 | 2 | 0 | 0 | ✅ | — |
| GcpCertManagerDnsAuthorization | 8 | 6 | 2 | 0 | 0 | ✅ | — |
| GcpCloudArmorPolicy | 61 | 10 | 51 | 0 | 0 | ✅ | — |
| GcpCloudComposerEnvironment | 81 | 3 | 66 | 12 | 0 | ✅ | — |
| GcpCloudComposerUserWorkloadsConfigMap | 6 | 4 | 2 | 0 | 0 | ✅ | — |
| GcpCloudComposerUserWorkloadsSecret | 6 | 4 | 2 | 0 | 0 | ✅ | — |
| GcpCloudFunction | 65 | 43 | 15 | 7 | 0 | ✅ | — |
| GcpCloudRun | 129 | 27 | 90 | 12 | 0 | ✅ | — |
| GcpCloudRunJob | 74 | 10 | 61 | 3 | 0 | ✅ | — |
| GcpCloudSchedulerJob | 32 | 29 | 3 | 0 | 0 | ✅ | — |
| GcpCloudSql | 147 | 40 | 97 | 10 | 0 | ✅ | — |
| GcpCloudSqlDatabase | 6 | 4 | 2 | 0 | 0 | ✅ | — |
| GcpCloudSqlUser | 14 | 10 | 2 | 2 | 0 | ✅ | — |
| GcpCloudTasksQueue | 29 | 22 | 7 | 0 | 0 | ✅ | — |
| GcpComputeDisk | 36 | 18 | 14 | 4 | 0 | ✅ | — |
| GcpComputeInstance | 124 | 47 | 64 | 13 | 0 | ✅ | — |
| GcpDataprocAutoscalingPolicy | 16 | 15 | 1 | 0 | 0 | ✅ | — |
| GcpDataprocCluster | 148 | 77 | 52 | 19 | 0 | ✅ | — |
| GcpDnsRecord | 49 | 43 | 6 | 0 | 0 | ✅ | — |
| GcpDnsZone | 23 | 18 | 2 | 3 | 0 | ✅ | — |
| GcpFilestoreInstance | 36 | 14 | 22 | 0 | 0 | ✅ | — |
| GcpFirestoreBackupSchedule | 5 | 4 | 1 | 0 | 0 | ✅ | — |
| GcpFirestoreDatabase | 15 | 11 | 4 | 0 | 0 | ✅ | — |
| GcpFirestoreIndex | 17 | 16 | 1 | 0 | 0 | ✅ | — |
| GcpFirewallRule | 20 | 12 | 0 | 0 | 13 | ❌ | — |
| GcpGcsBucket | 64 | 30 | 29 | 5 | 0 | ✅ | — |
| GcpGkeCluster | 537 | 61 | 137 | 339 | 0 | ✅ | — |
| GcpGkeNodePool | 184 | 127 | 51 | 6 | 0 | ✅ | — |
| GcpGkeWorkloadIdentityBinding | 6 | 3 | 0 | 3 | 0 | ✅ | ✅ pulumi, terraform |
| GcpGlobalAddress | 11 | 9 | 2 | 0 | 0 | ✅ | — |
| GcpGlobalForwardingRule | 23 | 18 | 4 | 1 | 0 | ✅ | — |
| GcpHealthCheck | 105 | 12 | 0 | 0 | 136 | ❌ | ✅ pulumi, terraform |
| GcpIamCustomRole | 7 | 6 | 1 | 0 | 0 | ✅ | — |
| GcpKmsKey | 12 | 10 | 2 | 0 | 0 | ✅ | — |
| GcpKmsKeyIamMember | 6 | 6 | 0 | 0 | 0 | ✅ | ✅ pulumi, terraform |
| GcpKmsKeyRing | 3 | 1 | 2 | 0 | 0 | ✅ | — |
| GcpManagedSslCertificate | 6 | 2 | 3 | 1 | 0 | ✅ | — |
| GcpMemorystoreInstance | 38 | 29 | 7 | 2 | 0 | ✅ | — |
| GcpProject | 14 | 5 | 5 | 4 | 0 | ✅ | — |
| GcpProjectIamMember | 6 | 5 | 1 | 0 | 0 | ✅ | — |
| GcpPubSubSchema | 5 | 3 | 2 | 0 | 0 | ✅ | — |
| GcpPubSubSubscription | 44 | 41 | 3 | 0 | 0 | ✅ | — |
| GcpPubSubTopic | 44 | 41 | 3 | 0 | 0 | ✅ | — |
| GcpRedisInstance | 32 | 24 | 6 | 2 | 0 | ✅ | — |
| GcpRegionNetworkEndpointGroup | 18 | 16 | 2 | 0 | 0 | ✅ | — |
| GcpRouterNat | 54 | 28 | 21 | 5 | 0 | ✅ | — |
| GcpServerlessVpcConnector | 13 | 9 | 2 | 2 | 0 | ✅ | — |
| GcpServiceAccount | 26 | 11 | 5 | 10 | 0 | ✅ | — |
| GcpServiceAccountIamMember | 6 | 6 | 0 | 0 | 0 | ✅ | ✅ pulumi, terraform |
| GcpServiceConnectionPolicy | 12 | 9 | 3 | 0 | 0 | ✅ | — |
| GcpServiceNetworkingConnection | 5 | 5 | 0 | 0 | 0 | ✅ | — |
| GcpSpannerBackupSchedule | 10 | 7 | 3 | 0 | 0 | ✅ | — |
| GcpSpannerDatabase | 12 | 10 | 2 | 0 | 0 | ✅ | — |
| GcpSpannerInstance | 28 | 21 | 7 | 0 | 0 | ✅ | — |
| GcpSslCertificate | 19 | 9 | 4 | 6 | 0 | ✅ | — |
| GcpSslPolicy | 17 | 13 | 4 | 0 | 0 | ✅ | — |
| GcpSubnetwork | 34 | 17 | 0 | 0 | 22 | ❌ | ✅ pulumi, terraform |
| GcpTargetHttpProxy | 7 | 5 | 2 | 0 | 0 | ✅ | — |
| GcpTargetHttpsProxy | 14 | 12 | 2 | 0 | 0 | ✅ | — |
| GcpUrlMap | 333 | 74 | 259 | 0 | 0 | ✅ | — |
| GcpVertexAiDeployedIndex | 16 | 12 | 4 | 0 | 0 | ✅ | — |
| GcpVertexAiEndpoint | 24 | 7 | 0 | 0 | 23 | ❌ | ✅ pulumi, terraform |
| GcpVertexAiIndex | 17 | 5 | 12 | 0 | 0 | ✅ | — |
| GcpVertexAiIndexEndpoint | 13 | 10 | 3 | 0 | 0 | ✅ | — |
| GcpVertexAiNotebook | 46 | 8 | 34 | 4 | 0 | ✅ | — |
| GcpVpcNetwork | 17 | 10 | 6 | 1 | 0 | ✅ | — |
| GcpWorkloadIdentityPool | 16 | 11 | 0 | 0 | 6 | ❌ | ✅ pulumi, terraform |
| GcpWorkloadIdentityPoolProvider | 16 | 15 | 1 | 0 | 0 | ✅ | — |

## Breadth: every GA resource, one disposition

All resources of `google@7.43.0` land in exactly one class:

| Disposition | Resources | Meaning |
|---|---|---|
| Modeled | 92 | consumed by a kind's Terraform module today |
| IAM-covered | 409 | per-resource IAM member/binding/policy triplets, covered by the owning kinds' additive `iam_members` fields |
| Composed | 2 | capability covered through an existing kind's surface rather than a kind of its own |
| Planned | 89 | judged to be covered by a planned kind or planned composition, not built yet |
| Deferred | 665 | deliberately not offered, each with the recorded reason |
| Excluded as deprecated | 76 | deprecated or superseded provider surface |
| **Total** | **1333** | |

## The enumerated record

The full per-resource record, so the accounting above is verifiable
rather than trusted.

### Modeled (92)

| Resource | Consuming kinds |
|---|---|
| `google_alloydb_cluster` | consumed by GcpAlloydbCluster |
| `google_alloydb_instance` | consumed by GcpAlloydbCluster, GcpAlloydbInstance |
| `google_alloydb_user` | consumed by GcpAlloydbUser |
| `google_artifact_registry_repository` | consumed by GcpArtifactRegistryRepo |
| `google_artifact_registry_repository_iam_member` | consumed by GcpArtifactRegistryRepo |
| `google_bigquery_dataset` | consumed by GcpBigQueryDataset |
| `google_bigquery_table` | consumed by GcpBigQueryTable |
| `google_bigtable_gc_policy` | consumed by GcpBigtableTable |
| `google_bigtable_instance` | consumed by GcpBigtableInstance |
| `google_bigtable_table` | consumed by GcpBigtableTable |
| `google_certificate_manager_certificate` | consumed by GcpCertManagerCert |
| `google_certificate_manager_dns_authorization` | consumed by GcpCertManagerDnsAuthorization |
| `google_cloud_run_service_iam_member` | consumed by GcpCloudFunction |
| `google_cloud_run_v2_job` | consumed by GcpCloudRunJob |
| `google_cloud_run_v2_service` | consumed by GcpCloudRun |
| `google_cloud_run_v2_service_iam_member` | consumed by GcpCloudRun |
| `google_cloud_scheduler_job` | consumed by GcpCloudSchedulerJob |
| `google_cloud_tasks_queue` | consumed by GcpCloudTasksQueue |
| `google_cloudfunctions2_function` | consumed by GcpCloudFunction |
| `google_composer_environment` | consumed by GcpCloudComposerEnvironment |
| `google_composer_user_workloads_config_map` | consumed by GcpCloudComposerUserWorkloadsConfigMap |
| `google_composer_user_workloads_secret` | consumed by GcpCloudComposerUserWorkloadsSecret |
| `google_compute_address` | consumed by GcpAddress |
| `google_compute_backend_bucket` | consumed by GcpBackendBucket |
| `google_compute_backend_bucket_signed_url_key` | consumed by GcpBackendBucket |
| `google_compute_backend_service` | consumed by GcpBackendService |
| `google_compute_backend_service_signed_url_key` | consumed by GcpBackendService |
| `google_compute_disk` | consumed by GcpComputeDisk |
| `google_compute_firewall` | consumed by GcpFirewallRule |
| `google_compute_global_address` | consumed by GcpGlobalAddress |
| `google_compute_global_forwarding_rule` | consumed by GcpGlobalForwardingRule |
| `google_compute_health_check` | consumed by GcpHealthCheck |
| `google_compute_instance` | consumed by GcpComputeInstance |
| `google_compute_managed_ssl_certificate` | consumed by GcpManagedSslCertificate |
| `google_compute_network` | consumed by GcpVpcNetwork |
| `google_compute_region_health_check` | consumed by GcpHealthCheck |
| `google_compute_region_network_endpoint_group` | consumed by GcpRegionNetworkEndpointGroup |
| `google_compute_region_ssl_certificate` | consumed by GcpSslCertificate |
| `google_compute_region_ssl_policy` | consumed by GcpSslPolicy |
| `google_compute_router` | consumed by GcpRouterNat |
| `google_compute_router_nat` | consumed by GcpRouterNat |
| `google_compute_security_policy` | consumed by GcpCloudArmorPolicy |
| `google_compute_ssl_certificate` | consumed by GcpSslCertificate |
| `google_compute_ssl_policy` | consumed by GcpSslPolicy |
| `google_compute_subnetwork` | consumed by GcpSubnetwork |
| `google_compute_target_http_proxy` | consumed by GcpTargetHttpProxy |
| `google_compute_target_https_proxy` | consumed by GcpTargetHttpsProxy |
| `google_compute_url_map` | consumed by GcpUrlMap |
| `google_container_cluster` | consumed by GcpGkeCluster |
| `google_container_node_pool` | consumed by GcpGkeNodePool |
| `google_dataproc_autoscaling_policy` | consumed by GcpDataprocAutoscalingPolicy |
| `google_dataproc_cluster` | consumed by GcpDataprocCluster |
| `google_dns_managed_zone` | consumed by GcpDnsZone |
| `google_dns_record_set` | consumed by GcpDnsRecord |
| `google_filestore_instance` | consumed by GcpFilestoreInstance |
| `google_firestore_backup_schedule` | consumed by GcpFirestoreBackupSchedule |
| `google_firestore_database` | consumed by GcpFirestoreDatabase |
| `google_firestore_index` | consumed by GcpFirestoreIndex |
| `google_iam_workload_identity_pool` | consumed by GcpWorkloadIdentityPool |
| `google_iam_workload_identity_pool_provider` | consumed by GcpWorkloadIdentityPoolProvider |
| `google_kms_crypto_key` | consumed by GcpKmsKey |
| `google_kms_crypto_key_iam_member` | consumed by GcpKmsKeyIamMember |
| `google_kms_key_ring` | consumed by GcpKmsKeyRing |
| `google_memorystore_instance` | consumed by GcpMemorystoreInstance |
| `google_network_connectivity_service_connection_policy` | consumed by GcpServiceConnectionPolicy |
| `google_organization_iam_member` | consumed by GcpServiceAccount |
| `google_project` | consumed by GcpProject |
| `google_project_iam_custom_role` | consumed by GcpIamCustomRole |
| `google_project_iam_member` | consumed by GcpProjectIamMember, GcpServiceAccount |
| `google_project_service` | consumed by GcpAddress, GcpAlloydbCluster, GcpAlloydbInstance, GcpAlloydbUser, GcpArtifactRegistryRepo, GcpBackendBucket, GcpBackendService, GcpBigQueryDataset, GcpBigQueryTable, GcpBigtableInstance, GcpBigtableTable, GcpCertManagerCert, GcpCertManagerDnsAuthorization, GcpCloudArmorPolicy, GcpCloudComposerEnvironment, GcpCloudFunction, GcpCloudRun, GcpCloudRunJob, GcpCloudSchedulerJob, GcpCloudSql, GcpCloudTasksQueue, GcpComputeDisk, GcpComputeInstance, GcpDataprocAutoscalingPolicy, GcpDataprocCluster, GcpDnsRecord, GcpDnsZone, GcpFilestoreInstance, GcpFirestoreBackupSchedule, GcpFirestoreDatabase, GcpFirestoreIndex, GcpGcsBucket, GcpGkeCluster, GcpGkeNodePool, GcpGlobalAddress, GcpGlobalForwardingRule, GcpHealthCheck, GcpKmsKey, GcpKmsKeyRing, GcpManagedSslCertificate, GcpMemorystoreInstance, GcpProject, GcpPubSubSchema, GcpPubSubSubscription, GcpPubSubTopic, GcpRedisInstance, GcpRegionNetworkEndpointGroup, GcpRouterNat, GcpServerlessVpcConnector, GcpServiceConnectionPolicy, GcpServiceNetworkingConnection, GcpSpannerBackupSchedule, GcpSpannerDatabase, GcpSpannerInstance, GcpSslCertificate, GcpSslPolicy, GcpSubnetwork, GcpTargetHttpProxy, GcpTargetHttpsProxy, GcpUrlMap, GcpVertexAiEndpoint, GcpVertexAiIndex, GcpVertexAiIndexEndpoint, GcpVertexAiNotebook, GcpVpcNetwork |
| `google_pubsub_schema` | consumed by GcpPubSubSchema |
| `google_pubsub_subscription` | consumed by GcpPubSubSubscription |
| `google_pubsub_topic` | consumed by GcpPubSubTopic |
| `google_redis_instance` | consumed by GcpRedisInstance |
| `google_service_account` | consumed by GcpServiceAccount |
| `google_service_account_iam_member` | consumed by GcpGkeWorkloadIdentityBinding, GcpServiceAccountIamMember |
| `google_service_account_key` | consumed by GcpServiceAccount |
| `google_service_networking_connection` | consumed by GcpServiceNetworkingConnection |
| `google_spanner_backup_schedule` | consumed by GcpSpannerBackupSchedule |
| `google_spanner_database` | consumed by GcpSpannerDatabase |
| `google_spanner_instance` | consumed by GcpSpannerInstance |
| `google_sql_database` | consumed by GcpCloudSqlDatabase |
| `google_sql_database_instance` | consumed by GcpCloudSql |
| `google_sql_user` | consumed by GcpCloudSqlUser |
| `google_storage_bucket` | consumed by GcpGcsBucket |
| `google_storage_bucket_iam_member` | consumed by GcpGcsBucket |
| `google_vertex_ai_endpoint` | consumed by GcpVertexAiEndpoint |
| `google_vertex_ai_index` | consumed by GcpVertexAiIndex |
| `google_vertex_ai_index_endpoint` | consumed by GcpVertexAiIndexEndpoint |
| `google_vertex_ai_index_endpoint_deployed_index` | consumed by GcpVertexAiDeployedIndex |
| `google_vpc_access_connector` | consumed by GcpServerlessVpcConnector |
| `google_workbench_instance` | consumed by GcpVertexAiNotebook |

### IAM-covered (409)

| Resource | Detail |
|---|---|
| `google_access_context_manager_access_policy_iam_binding` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_access_context_manager_access_policy_iam_member` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_access_context_manager_access_policy_iam_policy` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_apigee_environment_iam_binding` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_apigee_environment_iam_member` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_apigee_environment_iam_policy` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_artifact_registry_repository_iam_binding` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_artifact_registry_repository_iam_policy` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_beyondcorp_security_gateway_application_iam_binding` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_beyondcorp_security_gateway_application_iam_member` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_beyondcorp_security_gateway_application_iam_policy` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_beyondcorp_security_gateway_iam_binding` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_beyondcorp_security_gateway_iam_member` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_beyondcorp_security_gateway_iam_policy` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_biglake_iceberg_catalog_iam_binding` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_biglake_iceberg_catalog_iam_member` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_biglake_iceberg_catalog_iam_policy` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_biglake_iceberg_namespace_iam_binding` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_biglake_iceberg_namespace_iam_member` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_biglake_iceberg_namespace_iam_policy` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_biglake_iceberg_table_iam_binding` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_biglake_iceberg_table_iam_member` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_biglake_iceberg_table_iam_policy` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_bigquery_analytics_hub_data_exchange_iam_binding` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_bigquery_analytics_hub_data_exchange_iam_member` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_bigquery_analytics_hub_data_exchange_iam_policy` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_bigquery_analytics_hub_listing_iam_binding` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_bigquery_analytics_hub_listing_iam_member` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_bigquery_analytics_hub_listing_iam_policy` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_bigquery_connection_iam_binding` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_bigquery_connection_iam_member` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_bigquery_connection_iam_policy` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_bigquery_datapolicy_data_policy_iam_binding` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_bigquery_datapolicy_data_policy_iam_member` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_bigquery_datapolicy_data_policy_iam_policy` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_bigquery_datapolicyv2_data_policy_iam_binding` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_bigquery_datapolicyv2_data_policy_iam_member` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_bigquery_datapolicyv2_data_policy_iam_policy` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_bigquery_dataset_iam_binding` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_bigquery_dataset_iam_member` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_bigquery_dataset_iam_policy` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_bigquery_routine_iam_binding` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_bigquery_routine_iam_member` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_bigquery_routine_iam_policy` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_bigquery_table_iam_binding` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_bigquery_table_iam_member` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_bigquery_table_iam_policy` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_bigtable_instance_iam_binding` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_bigtable_instance_iam_member` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_bigtable_instance_iam_policy` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_bigtable_table_iam_binding` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_bigtable_table_iam_member` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_bigtable_table_iam_policy` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_billing_account_iam_binding` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_billing_account_iam_member` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_billing_account_iam_policy` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_binary_authorization_attestor_iam_binding` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_binary_authorization_attestor_iam_member` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_binary_authorization_attestor_iam_policy` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_cloud_run_service_iam_binding` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_cloud_run_service_iam_policy` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_cloud_run_v2_job_iam_binding` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_cloud_run_v2_job_iam_member` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_cloud_run_v2_job_iam_policy` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_cloud_run_v2_service_iam_binding` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_cloud_run_v2_service_iam_policy` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_cloud_run_v2_worker_pool_iam_binding` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_cloud_run_v2_worker_pool_iam_member` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_cloud_run_v2_worker_pool_iam_policy` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_cloud_tasks_queue_iam_binding` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_cloud_tasks_queue_iam_member` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_cloud_tasks_queue_iam_policy` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_cloudbuildv2_connection_iam_binding` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_cloudbuildv2_connection_iam_member` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_cloudbuildv2_connection_iam_policy` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_clouddeploy_custom_target_type_iam_binding` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_clouddeploy_custom_target_type_iam_member` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_clouddeploy_custom_target_type_iam_policy` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_clouddeploy_delivery_pipeline_iam_binding` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_clouddeploy_delivery_pipeline_iam_member` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_clouddeploy_delivery_pipeline_iam_policy` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_clouddeploy_target_iam_binding` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_clouddeploy_target_iam_member` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_clouddeploy_target_iam_policy` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_cloudfunctions2_function_iam_binding` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_cloudfunctions2_function_iam_member` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_cloudfunctions2_function_iam_policy` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_cloudfunctions_function_iam_binding` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_cloudfunctions_function_iam_member` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_cloudfunctions_function_iam_policy` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_colab_runtime_template_iam_binding` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_colab_runtime_template_iam_member` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_colab_runtime_template_iam_policy` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_compute_disk_iam_binding` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_compute_disk_iam_member` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_compute_disk_iam_policy` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_compute_firewall_policy_iam_binding` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_compute_firewall_policy_iam_member` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_compute_firewall_policy_iam_policy` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_compute_image_iam_binding` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_compute_image_iam_member` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_compute_image_iam_policy` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_compute_instance_iam_binding` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_compute_instance_iam_member` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_compute_instance_iam_policy` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_compute_instance_template_iam_binding` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_compute_instance_template_iam_member` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_compute_instance_template_iam_policy` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_compute_instant_snapshot_iam_binding` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_compute_instant_snapshot_iam_member` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_compute_instant_snapshot_iam_policy` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_compute_network_firewall_policy_iam_binding` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_compute_network_firewall_policy_iam_member` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_compute_network_firewall_policy_iam_policy` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_compute_region_disk_iam_binding` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_compute_region_disk_iam_member` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_compute_region_disk_iam_policy` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_compute_region_instant_snapshot_iam_binding` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_compute_region_instant_snapshot_iam_member` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_compute_region_instant_snapshot_iam_policy` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_compute_region_network_firewall_policy_iam_binding` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_compute_region_network_firewall_policy_iam_member` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_compute_region_network_firewall_policy_iam_policy` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_compute_snapshot_iam_binding` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_compute_snapshot_iam_member` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_compute_snapshot_iam_policy` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_compute_storage_pool_iam_binding` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_compute_storage_pool_iam_member` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_compute_storage_pool_iam_policy` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_compute_subnetwork_iam_binding` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_compute_subnetwork_iam_member` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_compute_subnetwork_iam_policy` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_container_analysis_note_iam_binding` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_container_analysis_note_iam_member` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_container_analysis_note_iam_policy` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_data_catalog_entry_group_iam_binding` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_data_catalog_entry_group_iam_member` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_data_catalog_entry_group_iam_policy` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_data_catalog_policy_tag_iam_binding` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_data_catalog_policy_tag_iam_member` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_data_catalog_policy_tag_iam_policy` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_data_catalog_tag_template_iam_binding` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_data_catalog_tag_template_iam_member` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_data_catalog_tag_template_iam_policy` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_data_catalog_taxonomy_iam_binding` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_data_catalog_taxonomy_iam_member` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_data_catalog_taxonomy_iam_policy` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_data_fusion_instance_iam_binding` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_data_fusion_instance_iam_member` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_data_fusion_instance_iam_policy` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_dataplex_aspect_type_iam_binding` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_dataplex_aspect_type_iam_member` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_dataplex_aspect_type_iam_policy` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_dataplex_asset_iam_binding` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_dataplex_asset_iam_member` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_dataplex_asset_iam_policy` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_dataplex_data_product_iam_binding` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_dataplex_data_product_iam_member` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_dataplex_data_product_iam_policy` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_dataplex_datascan_iam_binding` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_dataplex_datascan_iam_member` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_dataplex_datascan_iam_policy` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_dataplex_entry_group_iam_binding` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_dataplex_entry_group_iam_member` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_dataplex_entry_group_iam_policy` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_dataplex_entry_type_iam_binding` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_dataplex_entry_type_iam_member` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_dataplex_entry_type_iam_policy` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_dataplex_glossary_iam_binding` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_dataplex_glossary_iam_member` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_dataplex_glossary_iam_policy` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_dataplex_lake_iam_binding` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_dataplex_lake_iam_member` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_dataplex_lake_iam_policy` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_dataplex_task_iam_binding` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_dataplex_task_iam_member` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_dataplex_task_iam_policy` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_dataplex_zone_iam_binding` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_dataplex_zone_iam_member` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_dataplex_zone_iam_policy` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_dataproc_autoscaling_policy_iam_binding` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_dataproc_autoscaling_policy_iam_member` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_dataproc_autoscaling_policy_iam_policy` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_dataproc_cluster_iam_binding` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_dataproc_cluster_iam_member` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_dataproc_cluster_iam_policy` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_dataproc_job_iam_binding` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_dataproc_job_iam_member` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_dataproc_job_iam_policy` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_dataproc_metastore_database_iam_binding` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_dataproc_metastore_database_iam_member` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_dataproc_metastore_database_iam_policy` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_dataproc_metastore_federation_iam_binding` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_dataproc_metastore_federation_iam_member` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_dataproc_metastore_federation_iam_policy` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_dataproc_metastore_service_iam_binding` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_dataproc_metastore_service_iam_member` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_dataproc_metastore_service_iam_policy` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_dataproc_metastore_table_iam_binding` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_dataproc_metastore_table_iam_member` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_dataproc_metastore_table_iam_policy` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_discovery_engine_search_engine_iam_binding` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_discovery_engine_search_engine_iam_member` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_discovery_engine_search_engine_iam_policy` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_dns_managed_zone_iam_binding` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_dns_managed_zone_iam_member` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_dns_managed_zone_iam_policy` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_endpoints_service_consumers_iam_binding` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_endpoints_service_consumers_iam_member` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_endpoints_service_consumers_iam_policy` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_endpoints_service_iam_binding` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_endpoints_service_iam_member` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_endpoints_service_iam_policy` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_folder_iam_binding` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_folder_iam_member` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_folder_iam_policy` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_gemini_repository_group_iam_binding` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_gemini_repository_group_iam_member` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_gemini_repository_group_iam_policy` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_gke_backup_backup_plan_iam_binding` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_gke_backup_backup_plan_iam_member` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_gke_backup_backup_plan_iam_policy` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_gke_backup_restore_plan_iam_binding` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_gke_backup_restore_plan_iam_member` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_gke_backup_restore_plan_iam_policy` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_gke_hub_feature_iam_binding` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_gke_hub_feature_iam_member` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_gke_hub_feature_iam_policy` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_gke_hub_membership_iam_binding` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_gke_hub_membership_iam_member` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_gke_hub_membership_iam_policy` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_gke_hub_scope_iam_binding` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_gke_hub_scope_iam_member` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_gke_hub_scope_iam_policy` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_healthcare_consent_store_iam_binding` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_healthcare_consent_store_iam_member` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_healthcare_consent_store_iam_policy` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_healthcare_dataset_iam_binding` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_healthcare_dataset_iam_member` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_healthcare_dataset_iam_policy` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_healthcare_dicom_store_iam_binding` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_healthcare_dicom_store_iam_member` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_healthcare_dicom_store_iam_policy` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_healthcare_fhir_store_iam_binding` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_healthcare_fhir_store_iam_member` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_healthcare_fhir_store_iam_policy` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_healthcare_hl7_v2_store_iam_binding` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_healthcare_hl7_v2_store_iam_member` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_healthcare_hl7_v2_store_iam_policy` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_iam_workforce_pool_iam_binding` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_iam_workforce_pool_iam_member` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_iam_workforce_pool_iam_policy` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_iam_workload_identity_pool_iam_binding` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_iam_workload_identity_pool_iam_member` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_iam_workload_identity_pool_iam_policy` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_iap_agent_registry_agent_iam_binding` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_iap_agent_registry_agent_iam_member` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_iap_agent_registry_agent_iam_policy` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_iap_agent_registry_endpoint_iam_binding` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_iap_agent_registry_endpoint_iam_member` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_iap_agent_registry_endpoint_iam_policy` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_iap_agent_registry_iam_binding` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_iap_agent_registry_iam_member` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_iap_agent_registry_iam_policy` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_iap_agent_registry_mcp_server_iam_binding` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_iap_agent_registry_mcp_server_iam_member` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_iap_agent_registry_mcp_server_iam_policy` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_iap_app_engine_service_iam_binding` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_iap_app_engine_service_iam_member` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_iap_app_engine_service_iam_policy` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_iap_app_engine_version_iam_binding` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_iap_app_engine_version_iam_member` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_iap_app_engine_version_iam_policy` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_iap_location_web_iam_binding` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_iap_location_web_iam_member` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_iap_location_web_iam_policy` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_iap_tunnel_dest_group_iam_binding` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_iap_tunnel_dest_group_iam_member` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_iap_tunnel_dest_group_iam_policy` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_iap_tunnel_iam_binding` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_iap_tunnel_iam_member` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_iap_tunnel_iam_policy` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_iap_tunnel_instance_iam_binding` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_iap_tunnel_instance_iam_member` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_iap_tunnel_instance_iam_policy` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_iap_web_backend_service_iam_binding` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_iap_web_backend_service_iam_member` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_iap_web_backend_service_iam_policy` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_iap_web_cloud_run_service_iam_binding` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_iap_web_cloud_run_service_iam_member` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_iap_web_cloud_run_service_iam_policy` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_iap_web_forwarding_rule_service_iam_binding` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_iap_web_forwarding_rule_service_iam_member` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_iap_web_forwarding_rule_service_iam_policy` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_iap_web_iam_binding` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_iap_web_iam_member` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_iap_web_iam_policy` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_iap_web_region_backend_service_iam_binding` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_iap_web_region_backend_service_iam_member` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_iap_web_region_backend_service_iam_policy` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_iap_web_region_forwarding_rule_service_iam_binding` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_iap_web_region_forwarding_rule_service_iam_member` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_iap_web_region_forwarding_rule_service_iam_policy` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_iap_web_type_app_engine_iam_binding` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_iap_web_type_app_engine_iam_member` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_iap_web_type_app_engine_iam_policy` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_iap_web_type_compute_iam_binding` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_iap_web_type_compute_iam_member` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_iap_web_type_compute_iam_policy` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_kms_crypto_key_iam_binding` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_kms_crypto_key_iam_policy` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_kms_ekm_connection_iam_binding` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_kms_ekm_connection_iam_member` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_kms_ekm_connection_iam_policy` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_kms_key_ring_iam_binding` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_kms_key_ring_iam_member` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_kms_key_ring_iam_policy` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_logging_log_view_iam_binding` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_logging_log_view_iam_member` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_logging_log_view_iam_policy` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_network_connectivity_hub_iam_binding` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_network_connectivity_hub_iam_member` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_network_connectivity_hub_iam_policy` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_network_security_address_group_iam_binding` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_network_security_address_group_iam_member` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_network_security_address_group_iam_policy` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_notebooks_instance_iam_binding` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_notebooks_instance_iam_member` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_notebooks_instance_iam_policy` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_notebooks_runtime_iam_binding` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_notebooks_runtime_iam_member` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_notebooks_runtime_iam_policy` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_organization_iam_binding` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_organization_iam_policy` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_privateca_ca_pool_iam_binding` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_privateca_ca_pool_iam_member` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_privateca_ca_pool_iam_policy` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_privateca_certificate_template_iam_binding` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_privateca_certificate_template_iam_member` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_privateca_certificate_template_iam_policy` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_project_iam_binding` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_project_iam_policy` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_pubsub_schema_iam_binding` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_pubsub_schema_iam_member` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_pubsub_schema_iam_policy` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_pubsub_subscription_iam_binding` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_pubsub_subscription_iam_member` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_pubsub_subscription_iam_policy` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_pubsub_topic_iam_binding` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_pubsub_topic_iam_member` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_pubsub_topic_iam_policy` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_scc_source_iam_binding` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_scc_source_iam_member` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_scc_source_iam_policy` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_scc_v2_organization_source_iam_binding` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_scc_v2_organization_source_iam_member` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_scc_v2_organization_source_iam_policy` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_secret_manager_regional_secret_iam_binding` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_secret_manager_regional_secret_iam_member` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_secret_manager_regional_secret_iam_policy` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_secret_manager_secret_iam_binding` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_secret_manager_secret_iam_member` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_secret_manager_secret_iam_policy` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_secure_source_manager_instance_iam_binding` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_secure_source_manager_instance_iam_member` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_secure_source_manager_instance_iam_policy` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_secure_source_manager_repository_iam_binding` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_secure_source_manager_repository_iam_member` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_secure_source_manager_repository_iam_policy` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_service_account_iam_binding` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_service_account_iam_policy` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_service_directory_namespace_iam_binding` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_service_directory_namespace_iam_member` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_service_directory_namespace_iam_policy` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_service_directory_service_iam_binding` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_service_directory_service_iam_member` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_service_directory_service_iam_policy` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_sourcerepo_repository_iam_binding` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_sourcerepo_repository_iam_member` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_sourcerepo_repository_iam_policy` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_spanner_database_iam_binding` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_spanner_database_iam_member` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_spanner_database_iam_policy` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_spanner_instance_iam_binding` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_spanner_instance_iam_member` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_spanner_instance_iam_policy` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_storage_bucket_iam_binding` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_storage_bucket_iam_policy` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_storage_managed_folder_iam_binding` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_storage_managed_folder_iam_member` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_storage_managed_folder_iam_policy` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_tags_tag_key_iam_binding` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_tags_tag_key_iam_member` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_tags_tag_key_iam_policy` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_tags_tag_value_iam_binding` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_tags_tag_value_iam_member` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_tags_tag_value_iam_policy` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_vertex_ai_reasoning_engine_iam_binding` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_vertex_ai_reasoning_engine_iam_member` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_vertex_ai_reasoning_engine_iam_policy` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_workbench_instance_iam_binding` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_workbench_instance_iam_member` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_workbench_instance_iam_policy` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_workstations_workstation_config_iam_binding` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_workstations_workstation_config_iam_member` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_workstations_workstation_config_iam_policy` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_workstations_workstation_iam_binding` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_workstations_workstation_iam_member` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |
| `google_workstations_workstation_iam_policy` | per-resource IAM triplet, covered by the owning kind's additive iam_members field |

### Composed (2)

| Resource | Recorded reason |
|---|---|
| `google_bigquery_dataset_access` | the GcpBigQueryDataset spec models dataset access entries directly on the dataset, which is this standalone resource's entire surface |
| `google_project_iam_member_remove` | declarative member removal is inherent to the additive iam_members reconciliation on the IAM member kinds (GcpProjectIamMember); a dedicated removal escape hatch is redundant |

### Planned (89)

| Resource | Recorded reason |
|---|---|
| `google_certificate_manager_certificate_issuance_config` | planned composition into the existing GcpCertManagerCert kind (trust and issuance configuration) |
| `google_certificate_manager_certificate_map` | planned kind GcpCertificateMap |
| `google_certificate_manager_certificate_map_entry` | composes into the planned GcpCertificateMap kind |
| `google_certificate_manager_trust_config` | planned composition into the existing GcpCertManagerCert kind (trust and issuance configuration) |
| `google_cloud_run_domain_mapping` | planned composition into the existing GcpCloudRun kind (custom-domain mapping) |
| `google_compute_autoscaler` | planned kind GcpComputeMig (instance template, group manager, autoscaler, and per-instance configuration composed, zonal and regional) |
| `google_compute_bulk_per_instance_config` | planned kind GcpComputeMig (instance template, group manager, autoscaler, and per-instance configuration composed, zonal and regional) |
| `google_compute_forwarding_rule` | planned kind GcpRegionalLoadBalancer (internal/regional load-balancer facade; the global path is modeled by the consumed global LB resources) |
| `google_compute_global_network_endpoint` | composes into the planned load-balancer facade kinds (GcpHttpsLoadBalancer, GcpRegionalLoadBalancer) as backend endpoint-group blocks |
| `google_compute_global_network_endpoint_group` | composes into the planned load-balancer facade kinds (GcpHttpsLoadBalancer, GcpRegionalLoadBalancer) as backend endpoint-group blocks |
| `google_compute_instance_group` | composes into the planned load-balancer facade kinds (unmanaged instance groups as backends) |
| `google_compute_instance_group_manager` | planned kind GcpComputeMig (instance template, group manager, autoscaler, and per-instance configuration composed, zonal and regional) |
| `google_compute_instance_group_membership` | composes into the planned load-balancer facade kinds (unmanaged instance groups as backends) |
| `google_compute_instance_group_named_port` | composes into the planned load-balancer facade kinds (unmanaged instance groups as backends) |
| `google_compute_instance_template` | planned kind GcpComputeMig (instance template, group manager, autoscaler, and per-instance configuration composed, zonal and regional) |
| `google_compute_network_endpoint` | composes into the planned load-balancer facade kinds (GcpHttpsLoadBalancer, GcpRegionalLoadBalancer) as backend endpoint-group blocks |
| `google_compute_network_endpoint_group` | composes into the planned load-balancer facade kinds (GcpHttpsLoadBalancer, GcpRegionalLoadBalancer) as backend endpoint-group blocks |
| `google_compute_network_endpoints` | composes into the planned load-balancer facade kinds (GcpHttpsLoadBalancer, GcpRegionalLoadBalancer) as backend endpoint-group blocks |
| `google_compute_per_instance_config` | planned kind GcpComputeMig (instance template, group manager, autoscaler, and per-instance configuration composed, zonal and regional) |
| `google_compute_region_autoscaler` | planned kind GcpComputeMig (instance template, group manager, autoscaler, and per-instance configuration composed, zonal and regional) |
| `google_compute_region_backend_service` | planned kind GcpRegionalLoadBalancer (internal/regional load-balancer facade; the global path is modeled by the consumed global LB resources) |
| `google_compute_region_instance_group_manager` | planned kind GcpComputeMig (instance template, group manager, autoscaler, and per-instance configuration composed, zonal and regional) |
| `google_compute_region_instance_template` | planned kind GcpComputeMig (instance template, group manager, autoscaler, and per-instance configuration composed, zonal and regional) |
| `google_compute_region_network_endpoint` | planned kind GcpRegionalLoadBalancer (internal/regional load-balancer facade; the global path is modeled by the consumed global LB resources) |
| `google_compute_region_per_instance_config` | planned kind GcpComputeMig (instance template, group manager, autoscaler, and per-instance configuration composed, zonal and regional) |
| `google_compute_region_resize_request` | planned kind GcpComputeMig (instance template, group manager, autoscaler, and per-instance configuration composed, zonal and regional) |
| `google_compute_region_target_http_proxy` | planned kind GcpRegionalLoadBalancer (internal/regional load-balancer facade; the global path is modeled by the consumed global LB resources) |
| `google_compute_region_target_https_proxy` | planned kind GcpRegionalLoadBalancer (internal/regional load-balancer facade; the global path is modeled by the consumed global LB resources) |
| `google_compute_region_target_tcp_proxy` | planned kind GcpRegionalLoadBalancer (internal/regional load-balancer facade; the global path is modeled by the consumed global LB resources) |
| `google_compute_region_url_map` | planned kind GcpRegionalLoadBalancer (internal/regional load-balancer facade; the global path is modeled by the consumed global LB resources) |
| `google_compute_resize_request` | planned kind GcpComputeMig (instance template, group manager, autoscaler, and per-instance configuration composed, zonal and regional) |
| `google_compute_target_grpc_proxy` | planned kind GcpHttpsLoadBalancer (global load-balancer facade; TCP/SSL/gRPC proxy variants) |
| `google_compute_target_ssl_proxy` | planned kind GcpHttpsLoadBalancer (global load-balancer facade; TCP/SSL/gRPC proxy variants) |
| `google_compute_target_tcp_proxy` | planned kind GcpHttpsLoadBalancer (global load-balancer facade; TCP/SSL/gRPC proxy variants) |
| `google_eventarc_channel` | composes into the planned GcpEventarcTrigger kind |
| `google_eventarc_enrollment` | composes into the planned GcpEventarcMessageBus kind (Eventarc Advanced) |
| `google_eventarc_google_api_source` | composes into the planned GcpEventarcMessageBus kind (Eventarc Advanced) |
| `google_eventarc_google_channel_config` | composes into the planned GcpEventarcTrigger kind |
| `google_eventarc_message_bus` | planned kind GcpEventarcMessageBus (Eventarc Advanced) |
| `google_eventarc_pipeline` | composes into the planned GcpEventarcMessageBus kind (Eventarc Advanced) |
| `google_eventarc_trigger` | planned kind GcpEventarcTrigger |
| `google_iam_oauth_client` | planned kind GcpIamOauthClient |
| `google_iam_oauth_client_credential` | composes into the planned GcpIamOauthClient kind |
| `google_iap_settings` | composes into the planned GcpIapOauth kind (Identity-Aware Proxy configuration) |
| `google_iap_tunnel_dest_group` | composes into the planned GcpIapOauth kind (Identity-Aware Proxy configuration) |
| `google_identity_platform_config` | planned kind GcpIdentityPlatformConfig |
| `google_identity_platform_default_supported_idp_config` | composes into the planned GcpIdentityPlatformConfig and GcpIdentityPlatformTenant kinds (identity-provider configuration) |
| `google_identity_platform_inbound_saml_config` | composes into the planned GcpIdentityPlatformConfig and GcpIdentityPlatformTenant kinds (identity-provider configuration) |
| `google_identity_platform_oauth_idp_config` | composes into the planned GcpIdentityPlatformConfig and GcpIdentityPlatformTenant kinds (identity-provider configuration) |
| `google_identity_platform_tenant` | planned kind GcpIdentityPlatformTenant |
| `google_identity_platform_tenant_default_supported_idp_config` | composes into the planned GcpIdentityPlatformConfig and GcpIdentityPlatformTenant kinds (identity-provider configuration) |
| `google_identity_platform_tenant_inbound_saml_config` | composes into the planned GcpIdentityPlatformConfig and GcpIdentityPlatformTenant kinds (identity-provider configuration) |
| `google_identity_platform_tenant_oauth_idp_config` | composes into the planned GcpIdentityPlatformConfig and GcpIdentityPlatformTenant kinds (identity-provider configuration) |
| `google_logging_billing_account_bucket_config` | planned kind GcpLogBucket (bucket configs across scopes, with log views composed) |
| `google_logging_billing_account_exclusion` | composes into the planned GcpLoggingSink kind (sink exclusions) |
| `google_logging_billing_account_sink` | planned kind GcpLoggingSink (one kind with project/folder/organization/billing-account scope) |
| `google_logging_folder_bucket_config` | planned kind GcpLogBucket (bucket configs across scopes, with log views composed) |
| `google_logging_folder_exclusion` | composes into the planned GcpLoggingSink kind (sink exclusions) |
| `google_logging_folder_settings` | composes into the planned GcpLogBucket kind (scope-level logging settings) |
| `google_logging_folder_sink` | planned kind GcpLoggingSink (one kind with project/folder/organization/billing-account scope) |
| `google_logging_linked_dataset` | composes into the planned GcpLogBucket kind (BigQuery linked datasets) |
| `google_logging_log_view` | planned kind GcpLogBucket (bucket configs across scopes, with log views composed) |
| `google_logging_metric` | planned kind GcpLogMetric |
| `google_logging_organization_bucket_config` | planned kind GcpLogBucket (bucket configs across scopes, with log views composed) |
| `google_logging_organization_exclusion` | composes into the planned GcpLoggingSink kind (sink exclusions) |
| `google_logging_organization_settings` | composes into the planned GcpLogBucket kind (scope-level logging settings) |
| `google_logging_organization_sink` | planned kind GcpLoggingSink (one kind with project/folder/organization/billing-account scope) |
| `google_logging_project_bucket_config` | planned kind GcpLogBucket (bucket configs across scopes, with log views composed) |
| `google_logging_project_exclusion` | composes into the planned GcpLoggingSink kind (sink exclusions) |
| `google_logging_project_sink` | planned kind GcpLoggingSink (one kind with project/folder/organization/billing-account scope) |
| `google_monitoring_alert_policy` | planned kind GcpMonitoringAlertPolicy |
| `google_monitoring_custom_service` | planned kind GcpMonitoringSlo (covers the service, custom-service, and SLO resources) |
| `google_monitoring_dashboard` | planned kind GcpMonitoringDashboard |
| `google_monitoring_group` | composes into the planned GcpMonitoringAlertPolicy kind (group-scoped alerting) |
| `google_monitoring_monitored_project` | composes into the planned monitoring kinds (metrics-scope management) |
| `google_monitoring_notification_channel` | planned kind GcpMonitoringNotificationChannel |
| `google_monitoring_service` | planned kind GcpMonitoringSlo (covers the service, custom-service, and SLO resources) |
| `google_monitoring_slo` | planned kind GcpMonitoringSlo (covers the service, custom-service, and SLO resources) |
| `google_monitoring_uptime_check_config` | planned kind GcpMonitoringUptimeCheck |
| `google_secret_manager_regional_secret` | planned kind GcpSecretManagerSecret models regional secrets as a location flag |
| `google_secret_manager_regional_secret_version` | planned kind GcpSecretManagerSecret models regional secrets as a location flag |
| `google_secret_manager_secret` | planned kind GcpSecretManagerSecret (secret with versions composed) |
| `google_secret_manager_secret_version` | planned kind GcpSecretManagerSecret (secret with versions composed) |
| `google_storage_bucket_object` | planned composition into the existing GcpGcsBucket kind (bucket companions: objects, folders, HMAC keys, notifications) |
| `google_storage_folder` | planned composition into the existing GcpGcsBucket kind (bucket companions: objects, folders, HMAC keys, notifications) |
| `google_storage_hmac_key` | planned composition into the existing GcpGcsBucket kind (bucket companions: objects, folders, HMAC keys, notifications) |
| `google_storage_managed_folder` | planned composition into the existing GcpGcsBucket kind (bucket companions: objects, folders, HMAC keys, notifications) |
| `google_storage_notification` | planned composition into the existing GcpGcsBucket kind (bucket companions: objects, folders, HMAC keys, notifications) |
| `google_workflows_workflow` | planned kind GcpWorkflow |

### Deferred (665)

| Resource | Recorded reason |
|---|---|
| `google_access_context_manager_access_level` | VPC Service Controls surface judged as an access-policy/access-level/service-perimeter kind family; deferred pending demand |
| `google_access_context_manager_access_level_condition` | VPC Service Controls surface judged as an access-policy/access-level/service-perimeter kind family; deferred pending demand |
| `google_access_context_manager_access_levels` | VPC Service Controls surface judged as an access-policy/access-level/service-perimeter kind family; deferred pending demand |
| `google_access_context_manager_access_policy` | VPC Service Controls surface judged as an access-policy/access-level/service-perimeter kind family; deferred pending demand |
| `google_access_context_manager_authorized_orgs_desc` | VPC Service Controls surface judged as an access-policy/access-level/service-perimeter kind family; deferred pending demand |
| `google_access_context_manager_egress_policy` | VPC Service Controls surface judged as an access-policy/access-level/service-perimeter kind family; deferred pending demand |
| `google_access_context_manager_gcp_user_access_binding` | VPC Service Controls surface judged as an access-policy/access-level/service-perimeter kind family; deferred pending demand |
| `google_access_context_manager_ingress_policy` | VPC Service Controls surface judged as an access-policy/access-level/service-perimeter kind family; deferred pending demand |
| `google_access_context_manager_service_perimeter` | VPC Service Controls surface judged as an access-policy/access-level/service-perimeter kind family; deferred pending demand |
| `google_access_context_manager_service_perimeter_dry_run_egress_policy` | VPC Service Controls surface judged as an access-policy/access-level/service-perimeter kind family; deferred pending demand |
| `google_access_context_manager_service_perimeter_dry_run_ingress_policy` | VPC Service Controls surface judged as an access-policy/access-level/service-perimeter kind family; deferred pending demand |
| `google_access_context_manager_service_perimeter_dry_run_resource` | VPC Service Controls surface judged as an access-policy/access-level/service-perimeter kind family; deferred pending demand |
| `google_access_context_manager_service_perimeter_egress_policy` | VPC Service Controls surface judged as an access-policy/access-level/service-perimeter kind family; deferred pending demand |
| `google_access_context_manager_service_perimeter_ingress_policy` | VPC Service Controls surface judged as an access-policy/access-level/service-perimeter kind family; deferred pending demand |
| `google_access_context_manager_service_perimeter_resource` | VPC Service Controls surface judged as an access-policy/access-level/service-perimeter kind family; deferred pending demand |
| `google_access_context_manager_service_perimeters` | VPC Service Controls surface judged as an access-policy/access-level/service-perimeter kind family; deferred pending demand |
| `google_active_directory_domain` | Managed Microsoft AD judged as a managed-domain kind (with domain trusts composed); deferred pending demand |
| `google_active_directory_domain_trust` | Managed Microsoft AD judged as a managed-domain kind (with domain trusts composed); deferred pending demand |
| `google_agent_identity_auth_provider` | agent-platform identity surfaces are emerging and pre-consolidation; deferred pending demand |
| `google_agent_registry_binding` | agent-platform registry surfaces are emerging and pre-consolidation; deferred pending demand |
| `google_agent_registry_service` | agent-platform registry surfaces are emerging and pre-consolidation; deferred pending demand |
| `google_alloydb_backup` | judged to fold into the existing GcpAlloyDbCluster kind's spec (on-demand backups); the composition is not built |
| `google_apigee_addons_config` | Apigee API management is a named niche family (an eventual ~8-10 kind family: organization, environment, environment group, instance, proxies/products/developers/apps); deferred pending demand |
| `google_apigee_api` | Apigee API management is a named niche family (an eventual ~8-10 kind family: organization, environment, environment group, instance, proxies/products/developers/apps); deferred pending demand |
| `google_apigee_api_deployment` | Apigee API management is a named niche family (an eventual ~8-10 kind family: organization, environment, environment group, instance, proxies/products/developers/apps); deferred pending demand |
| `google_apigee_api_product` | Apigee API management is a named niche family (an eventual ~8-10 kind family: organization, environment, environment group, instance, proxies/products/developers/apps); deferred pending demand |
| `google_apigee_app_group` | Apigee API management is a named niche family (an eventual ~8-10 kind family: organization, environment, environment group, instance, proxies/products/developers/apps); deferred pending demand |
| `google_apigee_control_plane_access` | Apigee API management is a named niche family (an eventual ~8-10 kind family: organization, environment, environment group, instance, proxies/products/developers/apps); deferred pending demand |
| `google_apigee_data_collector` | Apigee API management is a named niche family (an eventual ~8-10 kind family: organization, environment, environment group, instance, proxies/products/developers/apps); deferred pending demand |
| `google_apigee_datastore` | Apigee API management is a named niche family (an eventual ~8-10 kind family: organization, environment, environment group, instance, proxies/products/developers/apps); deferred pending demand |
| `google_apigee_developer` | Apigee API management is a named niche family (an eventual ~8-10 kind family: organization, environment, environment group, instance, proxies/products/developers/apps); deferred pending demand |
| `google_apigee_developer_app` | Apigee API management is a named niche family (an eventual ~8-10 kind family: organization, environment, environment group, instance, proxies/products/developers/apps); deferred pending demand |
| `google_apigee_dns_zone` | Apigee API management is a named niche family (an eventual ~8-10 kind family: organization, environment, environment group, instance, proxies/products/developers/apps); deferred pending demand |
| `google_apigee_endpoint_attachment` | Apigee API management is a named niche family (an eventual ~8-10 kind family: organization, environment, environment group, instance, proxies/products/developers/apps); deferred pending demand |
| `google_apigee_env_keystore` | Apigee API management is a named niche family (an eventual ~8-10 kind family: organization, environment, environment group, instance, proxies/products/developers/apps); deferred pending demand |
| `google_apigee_env_references` | Apigee API management is a named niche family (an eventual ~8-10 kind family: organization, environment, environment group, instance, proxies/products/developers/apps); deferred pending demand |
| `google_apigee_envgroup` | Apigee API management is a named niche family (an eventual ~8-10 kind family: organization, environment, environment group, instance, proxies/products/developers/apps); deferred pending demand |
| `google_apigee_envgroup_attachment` | Apigee API management is a named niche family (an eventual ~8-10 kind family: organization, environment, environment group, instance, proxies/products/developers/apps); deferred pending demand |
| `google_apigee_environment` | Apigee API management is a named niche family (an eventual ~8-10 kind family: organization, environment, environment group, instance, proxies/products/developers/apps); deferred pending demand |
| `google_apigee_environment_addons_config` | Apigee API management is a named niche family (an eventual ~8-10 kind family: organization, environment, environment group, instance, proxies/products/developers/apps); deferred pending demand |
| `google_apigee_environment_api_revision_deployment` | Apigee API management is a named niche family (an eventual ~8-10 kind family: organization, environment, environment group, instance, proxies/products/developers/apps); deferred pending demand |
| `google_apigee_environment_debugmask` | Apigee API management is a named niche family (an eventual ~8-10 kind family: organization, environment, environment group, instance, proxies/products/developers/apps); deferred pending demand |
| `google_apigee_environment_keyvaluemaps` | Apigee API management is a named niche family (an eventual ~8-10 kind family: organization, environment, environment group, instance, proxies/products/developers/apps); deferred pending demand |
| `google_apigee_environment_keyvaluemaps_entries` | Apigee API management is a named niche family (an eventual ~8-10 kind family: organization, environment, environment group, instance, proxies/products/developers/apps); deferred pending demand |
| `google_apigee_flowhook` | Apigee API management is a named niche family (an eventual ~8-10 kind family: organization, environment, environment group, instance, proxies/products/developers/apps); deferred pending demand |
| `google_apigee_instance` | Apigee API management is a named niche family (an eventual ~8-10 kind family: organization, environment, environment group, instance, proxies/products/developers/apps); deferred pending demand |
| `google_apigee_instance_attachment` | Apigee API management is a named niche family (an eventual ~8-10 kind family: organization, environment, environment group, instance, proxies/products/developers/apps); deferred pending demand |
| `google_apigee_keystores_aliases_key_cert_file` | Apigee API management is a named niche family (an eventual ~8-10 kind family: organization, environment, environment group, instance, proxies/products/developers/apps); deferred pending demand |
| `google_apigee_keystores_aliases_pkcs12` | Apigee API management is a named niche family (an eventual ~8-10 kind family: organization, environment, environment group, instance, proxies/products/developers/apps); deferred pending demand |
| `google_apigee_keystores_aliases_self_signed_cert` | Apigee API management is a named niche family (an eventual ~8-10 kind family: organization, environment, environment group, instance, proxies/products/developers/apps); deferred pending demand |
| `google_apigee_nat_address` | Apigee API management is a named niche family (an eventual ~8-10 kind family: organization, environment, environment group, instance, proxies/products/developers/apps); deferred pending demand |
| `google_apigee_organization` | Apigee API management is a named niche family (an eventual ~8-10 kind family: organization, environment, environment group, instance, proxies/products/developers/apps); deferred pending demand |
| `google_apigee_security_action` | Apigee API management is a named niche family (an eventual ~8-10 kind family: organization, environment, environment group, instance, proxies/products/developers/apps); deferred pending demand |
| `google_apigee_security_feedback` | Apigee API management is a named niche family (an eventual ~8-10 kind family: organization, environment, environment group, instance, proxies/products/developers/apps); deferred pending demand |
| `google_apigee_security_monitoring_condition` | Apigee API management is a named niche family (an eventual ~8-10 kind family: organization, environment, environment group, instance, proxies/products/developers/apps); deferred pending demand |
| `google_apigee_security_profile_v2` | Apigee API management is a named niche family (an eventual ~8-10 kind family: organization, environment, environment group, instance, proxies/products/developers/apps); deferred pending demand |
| `google_apigee_sharedflow` | Apigee API management is a named niche family (an eventual ~8-10 kind family: organization, environment, environment group, instance, proxies/products/developers/apps); deferred pending demand |
| `google_apigee_sharedflow_deployment` | Apigee API management is a named niche family (an eventual ~8-10 kind family: organization, environment, environment group, instance, proxies/products/developers/apps); deferred pending demand |
| `google_apigee_space` | Apigee API management is a named niche family (an eventual ~8-10 kind family: organization, environment, environment group, instance, proxies/products/developers/apps); deferred pending demand |
| `google_apigee_sync_authorization` | Apigee API management is a named niche family (an eventual ~8-10 kind family: organization, environment, environment group, instance, proxies/products/developers/apps); deferred pending demand |
| `google_apigee_target_server` | Apigee API management is a named niche family (an eventual ~8-10 kind family: organization, environment, environment group, instance, proxies/products/developers/apps); deferred pending demand |
| `google_apihub_api_hub_instance` | API Hub is a new service and pairs with the Apigee deferral; deferred pending demand |
| `google_apihub_curation` | API Hub is a new service and pairs with the Apigee deferral; deferred pending demand |
| `google_apihub_host_project_registration` | API Hub is a new service and pairs with the Apigee deferral; deferred pending demand |
| `google_apihub_plugin` | API Hub is a new service and pairs with the Apigee deferral; deferred pending demand |
| `google_apihub_plugin_instance` | API Hub is a new service and pairs with the Apigee deferral; deferred pending demand |
| `google_apihub_runtime_project_attachment` | API Hub governance attachments are a specialty; deferred pending demand |
| `google_apikeys_key` | judged to deserve a GcpApiKey kind (Maps/Firebase restricted API keys); deferred pending demand |
| `google_apphub_application` | App Hub is an organizing overlay with low IaC demand today; deferred |
| `google_apphub_boundary` | App Hub is an organizing overlay with low IaC demand today; deferred |
| `google_apphub_service` | App Hub is an organizing overlay with low IaC demand today; deferred |
| `google_apphub_service_project_attachment` | App Hub is an organizing overlay with low IaC demand today; deferred |
| `google_apphub_workload` | App Hub is an organizing overlay with low IaC demand today; deferred |
| `google_artifact_registry_project_config` | judged to fold into the existing GcpArtifactRegistryRepo kind's spec (project-level config and cleanup rules); the composition is not built |
| `google_artifact_registry_rule` | judged to fold into the existing GcpArtifactRegistryRepo kind's spec (project-level config and cleanup rules); the composition is not built |
| `google_assured_workloads_workload` | Assured Workloads compliance folders are a regulated-industry specialty; deferred |
| `google_backup_dr_backup_plan` | Backup and DR judged as backup-vault and backup-plan kinds with companions composed; deferred pending demand |
| `google_backup_dr_backup_plan_association` | Backup and DR judged as backup-vault and backup-plan kinds with companions composed; deferred pending demand |
| `google_backup_dr_backup_vault` | Backup and DR judged as backup-vault and backup-plan kinds with companions composed; deferred pending demand |
| `google_backup_dr_management_server` | Backup and DR judged as backup-vault and backup-plan kinds with companions composed; deferred pending demand |
| `google_backup_dr_restore_workload` | Backup and DR judged as backup-vault and backup-plan kinds with companions composed; deferred pending demand |
| `google_backup_dr_service_config` | Backup and DR judged as backup-vault and backup-plan kinds with companions composed; deferred pending demand |
| `google_beyondcorp_app_connection` | BeyondCorp app connectors are a zero-trust specialty; deferred |
| `google_beyondcorp_app_connector` | BeyondCorp app connectors are a zero-trust specialty; deferred |
| `google_beyondcorp_app_gateway` | BeyondCorp app connectors are a zero-trust specialty; deferred |
| `google_beyondcorp_security_gateway` | BeyondCorp app connectors are a zero-trust specialty; deferred |
| `google_beyondcorp_security_gateway_application` | BeyondCorp app connectors are a zero-trust specialty; deferred |
| `google_biglake_catalog` | BigLake catalog metadata is emerging; revisit with lakehouse demand |
| `google_biglake_database` | BigLake catalog metadata is emerging; revisit with lakehouse demand |
| `google_biglake_iceberg_catalog` | BigLake catalog metadata is emerging; revisit with lakehouse demand |
| `google_biglake_iceberg_namespace` | BigLake catalog metadata is emerging; revisit with lakehouse demand |
| `google_biglake_iceberg_table` | BigLake catalog metadata is emerging; revisit with lakehouse demand |
| `google_biglake_table` | BigLake catalog metadata is emerging; revisit with lakehouse demand |
| `google_bigquery_analytics_hub_data_exchange` | Analytics Hub judged as data-exchange and listing kinds (subscriptions composed); deferred pending demand |
| `google_bigquery_analytics_hub_listing` | Analytics Hub judged as data-exchange and listing kinds (subscriptions composed); deferred pending demand |
| `google_bigquery_analytics_hub_listing_subscription` | Analytics Hub judged as data-exchange and listing kinds (subscriptions composed); deferred pending demand |
| `google_bigquery_analytics_hub_query_template` | Analytics Hub judged as data-exchange and listing kinds (subscriptions composed); deferred pending demand |
| `google_bigquery_bi_reservation` | BigQuery Reservations judged as one reservation kind (assignments, commitments, and BI reservation composed); deferred pending demand |
| `google_bigquery_capacity_commitment` | BigQuery Reservations judged as one reservation kind (assignments, commitments, and BI reservation composed); deferred pending demand |
| `google_bigquery_connection` | judged to deserve a GcpBigQueryConnection kind (prerequisite for BigLake, federated queries, and remote functions); deferred pending demand |
| `google_bigquery_data_transfer_config` | judged to deserve a GcpBigQueryDataTransfer kind (scheduled queries and SaaS ingestion); deferred pending demand |
| `google_bigquery_datapolicyv2_data_policy` | BigQuery data policies judged as a data-policy kind on the v2 API; deferred pending demand |
| `google_bigquery_job` | BigQuery jobs are imperative one-shot operations, a poor declarative fit |
| `google_bigquery_reservation` | BigQuery Reservations judged as one reservation kind (assignments, commitments, and BI reservation composed); deferred pending demand |
| `google_bigquery_reservation_assignment` | BigQuery Reservations judged as one reservation kind (assignments, commitments, and BI reservation composed); deferred pending demand |
| `google_bigquery_reservation_group` | BigQuery Reservations judged as one reservation kind (assignments, commitments, and BI reservation composed); deferred pending demand |
| `google_bigquery_routine` | judged to fold into the existing GcpBigQueryDataset kind's spec (stored routines); the composition is not built |
| `google_bigquery_row_access_policy` | judged to fold into the existing GcpBigQueryTable kind's spec (row-level access policies); the composition is not built |
| `google_bigtable_app_profile` | judged to fold into the existing GcpBigtableInstance kind's spec (app profiles); the composition is not built |
| `google_bigtable_authorized_view` | judged to fold into the existing GcpBigtableTable kind's spec (views and schema bundles); the composition is not built |
| `google_bigtable_logical_view` | judged to fold into the existing GcpBigtableTable kind's spec (views and schema bundles); the composition is not built |
| `google_bigtable_materialized_view` | judged to fold into the existing GcpBigtableTable kind's spec (views and schema bundles); the composition is not built |
| `google_bigtable_schema_bundle` | judged to fold into the existing GcpBigtableTable kind's spec (views and schema bundles); the composition is not built |
| `google_billing_budget` | judged to deserve a GcpBillingBudget kind (FinOps staple); deferred pending demand |
| `google_billing_project_info` | judged to fold into the existing GcpProject kind's spec (billing-account link); the composition is not built |
| `google_billing_subaccount` | billing subaccounts serve resellers; deferred |
| `google_binary_authorization_attestor` | Binary Authorization judged as policy and attestor kinds; deferred pending demand |
| `google_binary_authorization_policy` | Binary Authorization judged as policy and attestor kinds; deferred pending demand |
| `google_ces_agent` | Customer Engagement Suite conversational-agent authoring is a specialty family; deferred pending demand |
| `google_ces_app` | Customer Engagement Suite conversational-agent authoring is a specialty family; deferred pending demand |
| `google_ces_app_root_agent_association` | Customer Engagement Suite conversational-agent authoring is a specialty family; deferred pending demand |
| `google_ces_app_version` | Customer Engagement Suite conversational-agent authoring is a specialty family; deferred pending demand |
| `google_ces_deployment` | Customer Engagement Suite conversational-agent authoring is a specialty family; deferred pending demand |
| `google_ces_example` | Customer Engagement Suite conversational-agent authoring is a specialty family; deferred pending demand |
| `google_ces_guardrail` | Customer Engagement Suite conversational-agent authoring is a specialty family; deferred pending demand |
| `google_ces_tool` | Customer Engagement Suite conversational-agent authoring is a specialty family; deferred pending demand |
| `google_ces_toolset` | Customer Engagement Suite conversational-agent authoring is a specialty family; deferred pending demand |
| `google_chronicle_big_query_export` | Chronicle SecOps is a named niche family (an eventual ~8 kind family); deferred pending demand |
| `google_chronicle_custom_list` | Chronicle SecOps is a named niche family (an eventual ~8 kind family); deferred pending demand |
| `google_chronicle_dashboard_chart` | Chronicle SecOps is a named niche family (an eventual ~8 kind family); deferred pending demand |
| `google_chronicle_data_access_label` | Chronicle SecOps is a named niche family (an eventual ~8 kind family); deferred pending demand |
| `google_chronicle_data_access_scope` | Chronicle SecOps is a named niche family (an eventual ~8 kind family); deferred pending demand |
| `google_chronicle_data_export` | Chronicle SecOps is a named niche family (an eventual ~8 kind family); deferred pending demand |
| `google_chronicle_data_table` | Chronicle SecOps is a named niche family (an eventual ~8 kind family); deferred pending demand |
| `google_chronicle_data_table_row` | Chronicle SecOps is a named niche family (an eventual ~8 kind family); deferred pending demand |
| `google_chronicle_environment` | Chronicle SecOps is a named niche family (an eventual ~8 kind family); deferred pending demand |
| `google_chronicle_environment_group` | Chronicle SecOps is a named niche family (an eventual ~8 kind family); deferred pending demand |
| `google_chronicle_feed` | Chronicle SecOps is a named niche family (an eventual ~8 kind family); deferred pending demand |
| `google_chronicle_findings_refinement` | Chronicle SecOps is a named niche family (an eventual ~8 kind family); deferred pending demand |
| `google_chronicle_findings_refinement_deployment` | Chronicle SecOps is a named niche family (an eventual ~8 kind family); deferred pending demand |
| `google_chronicle_native_dashboard` | Chronicle SecOps is a named niche family (an eventual ~8 kind family); deferred pending demand |
| `google_chronicle_parser` | Chronicle SecOps is a named niche family (an eventual ~8 kind family); deferred pending demand |
| `google_chronicle_parser_extension` | Chronicle SecOps is a named niche family (an eventual ~8 kind family); deferred pending demand |
| `google_chronicle_reference_list` | Chronicle SecOps is a named niche family (an eventual ~8 kind family); deferred pending demand |
| `google_chronicle_retrohunt` | Chronicle SecOps is a named niche family (an eventual ~8 kind family); deferred pending demand |
| `google_chronicle_rule` | Chronicle SecOps is a named niche family (an eventual ~8 kind family); deferred pending demand |
| `google_chronicle_rule_deployment` | Chronicle SecOps is a named niche family (an eventual ~8 kind family); deferred pending demand |
| `google_chronicle_watchlist` | Chronicle SecOps is a named niche family (an eventual ~8 kind family); deferred pending demand |
| `google_cloud_asset_folder_feed` | Cloud Asset feeds judged as one feed kind with project/folder/organization scope; deferred pending demand |
| `google_cloud_asset_organization_feed` | Cloud Asset feeds judged as one feed kind with project/folder/organization scope; deferred pending demand |
| `google_cloud_asset_project_feed` | Cloud Asset feeds judged as one feed kind with project/folder/organization scope; deferred pending demand |
| `google_cloud_identity_group` | Cloud Identity groups judged as a group kind (memberships composed) for Google Groups IAM; deferred pending demand |
| `google_cloud_identity_group_membership` | Cloud Identity groups judged as a group kind (memberships composed) for Google Groups IAM; deferred pending demand |
| `google_cloud_ids_endpoint` | Cloud IDS is superseded in practice by NGFW intrusion features; deferred |
| `google_cloud_quotas_quota_adjuster_settings` | quota preferences are rarely IaC-managed; deferred |
| `google_cloud_quotas_quota_preference` | quota preferences are rarely IaC-managed; deferred |
| `google_cloud_run_v2_worker_pool` | judged to deserve a GcpCloudRunWorkerPool kind (the third Cloud Run runtime shape beside service and job); deferred pending demand |
| `google_cloud_security_compliance_cloud_control` | Compliance Manager frameworks are org-governance surface; deferred pending demand |
| `google_cloud_security_compliance_framework` | Compliance Manager frameworks are org-governance surface; deferred pending demand |
| `google_cloud_security_compliance_framework_deployment` | Compliance Manager frameworks are org-governance surface; deferred pending demand |
| `google_cloud_support_support_event_subscription` | Cloud Support event subscriptions are operational tooling, not provisioned infrastructure; deferred |
| `google_cloudbuild_trigger` | Cloud Build judged as trigger and worker-pool kinds; deferred pending demand |
| `google_cloudbuild_worker_pool` | Cloud Build judged as trigger and worker-pool kinds; deferred pending demand |
| `google_cloudbuildv2_connection` | Cloud Build v2 judged as a connection kind (repositories composed); deferred pending demand |
| `google_cloudbuildv2_repository` | Cloud Build v2 judged as a connection kind (repositories composed); deferred pending demand |
| `google_clouddeploy_automation` | Cloud Deploy judged as delivery-pipeline and target kinds (automations and policies composed); deferred pending demand |
| `google_clouddeploy_custom_target_type` | Cloud Deploy judged as delivery-pipeline and target kinds (automations and policies composed); deferred pending demand |
| `google_clouddeploy_delivery_pipeline` | Cloud Deploy judged as delivery-pipeline and target kinds (automations and policies composed); deferred pending demand |
| `google_clouddeploy_deploy_policy` | Cloud Deploy judged as delivery-pipeline and target kinds (automations and policies composed); deferred pending demand |
| `google_clouddeploy_target` | Cloud Deploy judged as delivery-pipeline and target kinds (automations and policies composed); deferred pending demand |
| `google_clouddomains_registration` | domain registration through IaC is rare; deferred |
| `google_colab_notebook_execution` | Colab Enterprise runtimes are a specialty; Workbench covers the mainstream need |
| `google_colab_runtime` | Colab Enterprise runtimes are a specialty; Workbench covers the mainstream need |
| `google_colab_runtime_template` | Colab Enterprise runtimes are a specialty; Workbench covers the mainstream need |
| `google_colab_schedule` | Colab Enterprise runtimes are a specialty; Workbench covers the mainstream need |
| `google_compute_attached_disk` | judged to fold into the existing GcpComputeInstance kind's spec (attached disks, from-template creation, and instance settings); the composition is not built |
| `google_compute_cross_site_network` | Cross-Site Interconnect (cross-site networks, wire groups) is specialty networking; deferred pending demand |
| `google_compute_disk_async_replication` | judged to fold into the existing GcpComputeDisk kind's spec (snapshots, regional disks, resource policies, and async replication); the composition is not built |
| `google_compute_disk_resource_policy_attachment` | judged to fold into the existing GcpComputeDisk kind's spec (snapshots, regional disks, resource policies, and async replication); the composition is not built |
| `google_compute_external_vpn_gateway` | HA VPN judged as one kind (HA gateway, external gateway, and tunnels composed); deferred pending demand |
| `google_compute_firewall_policy` | firewall policies judged as hierarchical and network firewall-policy kinds (regional as a flag; rules and associations composed); deferred pending demand |
| `google_compute_firewall_policy_association` | firewall policies judged as hierarchical and network firewall-policy kinds (regional as a flag; rules and associations composed); deferred pending demand |
| `google_compute_firewall_policy_rule` | firewall policies judged as hierarchical and network firewall-policy kinds (regional as a flag; rules and associations composed); deferred pending demand |
| `google_compute_firewall_policy_with_rules` | firewall policies judged as hierarchical and network firewall-policy kinds (regional as a flag; rules and associations composed); deferred pending demand |
| `google_compute_global_vm_extension_policy` | VM extension policies are emerging fleet tooling; deferred pending demand |
| `google_compute_ha_vpn_gateway` | HA VPN judged as one kind (HA gateway, external gateway, and tunnels composed); deferred pending demand |
| `google_compute_image` | judged to deserve a compute-image kind (golden images); deferred pending demand |
| `google_compute_instance_from_template` | judged to fold into the existing GcpComputeInstance kind's spec (attached disks, from-template creation, and instance settings); the composition is not built |
| `google_compute_instance_settings` | judged to fold into the existing GcpComputeInstance kind's spec (attached disks, from-template creation, and instance settings); the composition is not built |
| `google_compute_instant_snapshot` | judged to fold into the existing GcpComputeDisk kind's spec (snapshots, regional disks, resource policies, and async replication); the composition is not built |
| `google_compute_interconnect` | physical interconnects, sole-tenancy nodes, capacity reservations, packet mirroring, and public IP prefixes are niche surfaces; deferred |
| `google_compute_interconnect_attachment` | judged to deserve an interconnect-attachment kind; deferred pending demand |
| `google_compute_interconnect_attachment_group` | physical interconnects, sole-tenancy nodes, capacity reservations, packet mirroring, and public IP prefixes are niche surfaces; deferred |
| `google_compute_interconnect_group` | physical interconnects, sole-tenancy nodes, capacity reservations, packet mirroring, and public IP prefixes are niche surfaces; deferred |
| `google_compute_network_attachment` | judged to fold into the existing GcpVpcNetwork kind's spec (static routes and network attachments); the composition is not built |
| `google_compute_network_firewall_policy` | firewall policies judged as hierarchical and network firewall-policy kinds (regional as a flag; rules and associations composed); deferred pending demand |
| `google_compute_network_firewall_policy_association` | firewall policies judged as hierarchical and network firewall-policy kinds (regional as a flag; rules and associations composed); deferred pending demand |
| `google_compute_network_firewall_policy_rule` | firewall policies judged as hierarchical and network firewall-policy kinds (regional as a flag; rules and associations composed); deferred pending demand |
| `google_compute_network_firewall_policy_with_rules` | firewall policies judged as hierarchical and network firewall-policy kinds (regional as a flag; rules and associations composed); deferred pending demand |
| `google_compute_network_peering` | judged as VPC-peering and shared-VPC kinds; deferred pending demand |
| `google_compute_network_peering_routes_config` | judged as VPC-peering and shared-VPC kinds; deferred pending demand |
| `google_compute_node_group` | physical interconnects, sole-tenancy nodes, capacity reservations, packet mirroring, and public IP prefixes are niche surfaces; deferred |
| `google_compute_node_template` | physical interconnects, sole-tenancy nodes, capacity reservations, packet mirroring, and public IP prefixes are niche surfaces; deferred |
| `google_compute_organization_security_policy` | organization security policies are org-admin surface (hierarchical firewall policies are the project-reachable path); deferred pending demand |
| `google_compute_organization_security_policy_association` | organization security policies are org-admin surface (hierarchical firewall policies are the project-reachable path); deferred pending demand |
| `google_compute_organization_security_policy_rule` | organization security policies are org-admin surface (hierarchical firewall policies are the project-reachable path); deferred pending demand |
| `google_compute_packet_mirroring` | physical interconnects, sole-tenancy nodes, capacity reservations, packet mirroring, and public IP prefixes are niche surfaces; deferred |
| `google_compute_preview_feature` | preview-feature opt-ins are not durable infrastructure; deferred |
| `google_compute_project_cloud_armor_tier` | judged to fold into the existing GcpProject kind's spec (project-level compute defaults and metadata); the composition is not built |
| `google_compute_project_default_network_tier` | judged to fold into the existing GcpProject kind's spec (project-level compute defaults and metadata); the composition is not built |
| `google_compute_project_metadata` | judged to fold into the existing GcpProject kind's spec (project-level compute defaults and metadata); the composition is not built |
| `google_compute_project_metadata_item` | judged to fold into the existing GcpProject kind's spec (project-level compute defaults and metadata); the composition is not built |
| `google_compute_public_advertised_prefix` | physical interconnects, sole-tenancy nodes, capacity reservations, packet mirroring, and public IP prefixes are niche surfaces; deferred |
| `google_compute_public_delegated_prefix` | physical interconnects, sole-tenancy nodes, capacity reservations, packet mirroring, and public IP prefixes are niche surfaces; deferred |
| `google_compute_region_commitment` | physical interconnects, sole-tenancy nodes, capacity reservations, packet mirroring, and public IP prefixes are niche surfaces; deferred |
| `google_compute_region_composite_health_check` | composite health-check aggregation is emerging load-balancing surface; deferred pending demand |
| `google_compute_region_disk` | judged to fold into the existing GcpComputeDisk kind's spec (snapshots, regional disks, resource policies, and async replication); the composition is not built |
| `google_compute_region_disk_resource_policy_attachment` | judged to fold into the existing GcpComputeDisk kind's spec (snapshots, regional disks, resource policies, and async replication); the composition is not built |
| `google_compute_region_health_aggregation_policy` | composite health-check aggregation is emerging load-balancing surface; deferred pending demand |
| `google_compute_region_health_source` | composite health-check aggregation is emerging load-balancing surface; deferred pending demand |
| `google_compute_region_instant_snapshot` | judged to fold into the existing GcpComputeDisk kind's spec (snapshots, regional disks, resource policies, and async replication); the composition is not built |
| `google_compute_region_network_firewall_policy` | firewall policies judged as hierarchical and network firewall-policy kinds (regional as a flag; rules and associations composed); deferred pending demand |
| `google_compute_region_network_firewall_policy_association` | firewall policies judged as hierarchical and network firewall-policy kinds (regional as a flag; rules and associations composed); deferred pending demand |
| `google_compute_region_network_firewall_policy_rule` | firewall policies judged as hierarchical and network firewall-policy kinds (regional as a flag; rules and associations composed); deferred pending demand |
| `google_compute_region_network_firewall_policy_with_rules` | firewall policies judged as hierarchical and network firewall-policy kinds (regional as a flag; rules and associations composed); deferred pending demand |
| `google_compute_region_security_policy` | judged to fold into the existing GcpCloudArmorPolicy kind's spec (standalone and regional rules); the composition is not built |
| `google_compute_region_security_policy_rule` | judged to fold into the existing GcpCloudArmorPolicy kind's spec (standalone and regional rules); the composition is not built |
| `google_compute_reservation` | physical interconnects, sole-tenancy nodes, capacity reservations, packet mirroring, and public IP prefixes are niche surfaces; deferred |
| `google_compute_resource_policy` | judged to fold into the existing GcpComputeDisk kind's spec (snapshots, regional disks, resource policies, and async replication); the composition is not built |
| `google_compute_resource_policy_attachment` | judged to fold into the existing GcpComputeDisk kind's spec (snapshots, regional disks, resource policies, and async replication); the composition is not built |
| `google_compute_rollout_plan` | compute rollout plans are emerging release tooling; deferred pending demand |
| `google_compute_route` | judged to fold into the existing GcpVpcNetwork kind's spec (static routes and network attachments); the composition is not built |
| `google_compute_router_interface` | judged to fold into the existing GcpRouterNat kind and the router family's specs; the composition is not built |
| `google_compute_router_named_set` | judged to fold into the existing GcpRouterNat kind and the router family's specs; the composition is not built |
| `google_compute_router_nat_address` | judged to fold into the existing GcpRouterNat kind and the router family's specs; the composition is not built |
| `google_compute_router_peer` | judged to fold into the existing GcpRouterNat kind and the router family's specs; the composition is not built |
| `google_compute_router_route_policy` | judged to fold into the existing GcpRouterNat kind and the router family's specs; the composition is not built |
| `google_compute_security_policy_rule` | judged to fold into the existing GcpCloudArmorPolicy kind's spec (standalone and regional rules); the composition is not built |
| `google_compute_service_attachment` | judged to deserve a Private Service Connect service-attachment kind; deferred pending demand |
| `google_compute_shared_vpc_host_project` | judged as VPC-peering and shared-VPC kinds; deferred pending demand |
| `google_compute_shared_vpc_service_project` | judged as VPC-peering and shared-VPC kinds; deferred pending demand |
| `google_compute_snapshot` | judged to fold into the existing GcpComputeDisk kind's spec (snapshots, regional disks, resource policies, and async replication); the composition is not built |
| `google_compute_snapshot_settings` | judged to fold into the existing GcpComputeDisk kind's spec (snapshots, regional disks, resource policies, and async replication); the composition is not built |
| `google_compute_storage_pool` | judged to deserve a Hyperdisk storage-pool kind; deferred pending demand |
| `google_compute_target_instance` | physical interconnects, sole-tenancy nodes, capacity reservations, packet mirroring, and public IP prefixes are niche surfaces; deferred |
| `google_compute_vpn_tunnel` | HA VPN judged as one kind (HA gateway, external gateway, and tunnels composed); deferred pending demand |
| `google_compute_wire_group` | Cross-Site Interconnect (cross-site networks, wire groups) is specialty networking; deferred pending demand |
| `google_compute_zone_vm_extension_policy` | VM extension policies are emerging fleet tooling; deferred pending demand |
| `google_config_deployment` | Infrastructure Manager deployments are themselves IaC orchestration, not catalog surface; deferred |
| `google_contact_center_insights_analysis_rule` | Contact Center AI Insights is a specialty; deferred |
| `google_contact_center_insights_assessment_rule` | Contact Center AI Insights is a specialty; deferred |
| `google_contact_center_insights_auto_labeling_rule` | Contact Center AI Insights is a specialty; deferred |
| `google_contact_center_insights_encryption_spec` | Contact Center AI Insights is a specialty; deferred |
| `google_contact_center_insights_qa_question` | Contact Center AI Insights is a specialty; deferred |
| `google_contact_center_insights_qa_scorecard` | Contact Center AI Insights is a specialty; deferred |
| `google_contact_center_insights_qa_scorecard_revision` | Contact Center AI Insights is a specialty; deferred |
| `google_contact_center_insights_view` | Contact Center AI Insights is a specialty; deferred |
| `google_container_analysis_note` | judged to fold into a Binary Authorization attestor kind when admitted (attestation notes and occurrences); deferred |
| `google_container_analysis_occurrence` | judged to fold into a Binary Authorization attestor kind when admitted (attestation notes and occurrences); deferred |
| `google_container_attached_cluster` | attached and multi-cloud GKE clusters (AWS/Azure) are a specialty; deferred |
| `google_container_aws_cluster` | attached and multi-cloud GKE clusters (AWS/Azure) are a specialty; deferred |
| `google_container_aws_node_pool` | attached and multi-cloud GKE clusters (AWS/Azure) are a specialty; deferred |
| `google_container_azure_client` | attached and multi-cloud GKE clusters (AWS/Azure) are a specialty; deferred |
| `google_container_azure_cluster` | attached and multi-cloud GKE clusters (AWS/Azure) are a specialty; deferred |
| `google_container_azure_node_pool` | attached and multi-cloud GKE clusters (AWS/Azure) are a specialty; deferred |
| `google_data_fusion_instance` | judged to deserve a GcpDataFusionInstance kind; deferred pending demand |
| `google_data_lineage_config` | data lineage config folds into the Dataplex governance family judgment; deferred pending demand |
| `google_data_loss_prevention_deidentify_template` | Cloud DLP judged as template and job-trigger kinds (stored info types and discovery configs composed); deferred pending demand |
| `google_data_loss_prevention_discovery_config` | Cloud DLP judged as template and job-trigger kinds (stored info types and discovery configs composed); deferred pending demand |
| `google_data_loss_prevention_inspect_template` | Cloud DLP judged as template and job-trigger kinds (stored info types and discovery configs composed); deferred pending demand |
| `google_data_loss_prevention_job_trigger` | Cloud DLP judged as template and job-trigger kinds (stored info types and discovery configs composed); deferred pending demand |
| `google_data_loss_prevention_stored_info_type` | Cloud DLP judged as template and job-trigger kinds (stored info types and discovery configs composed); deferred pending demand |
| `google_data_pipeline_pipeline` | Data Pipelines is a thin wrapper with low usage; deferred |
| `google_database_migration_service_connection_profile` | Database Migration Service is episodic migration tooling, not steady-state infrastructure; deferred |
| `google_database_migration_service_migration_job` | Database Migration Service is episodic migration tooling, not steady-state infrastructure; deferred |
| `google_database_migration_service_private_connection` | Database Migration Service is episodic migration tooling, not steady-state infrastructure; deferred |
| `google_dataflow_job` | judged to deserve a GcpDataflowJob kind (long-running streaming pipelines); deferred pending demand |
| `google_dataform_folder` | folds into the Dataform repository kind judged alongside the beta-only core resource (see the beta admission list); deferred until that admission |
| `google_dataform_team_folder` | folds into the Dataform repository kind judged alongside the beta-only core resource (see the beta admission list); deferred until that admission |
| `google_dataplex_aspect_type` | Dataplex data governance judged as a ~6 kind family (lake, zone, asset, datascan, catalog entries, glossary); deferred pending demand |
| `google_dataplex_asset` | Dataplex data governance judged as a ~6 kind family (lake, zone, asset, datascan, catalog entries, glossary); deferred pending demand |
| `google_dataplex_data_product` | Dataplex data governance judged as a ~6 kind family (lake, zone, asset, datascan, catalog entries, glossary); deferred pending demand |
| `google_dataplex_data_product_data_asset` | Dataplex data governance judged as a ~6 kind family (lake, zone, asset, datascan, catalog entries, glossary); deferred pending demand |
| `google_dataplex_datascan` | Dataplex data governance judged as a ~6 kind family (lake, zone, asset, datascan, catalog entries, glossary); deferred pending demand |
| `google_dataplex_entry` | Dataplex data governance judged as a ~6 kind family (lake, zone, asset, datascan, catalog entries, glossary); deferred pending demand |
| `google_dataplex_entry_group` | Dataplex data governance judged as a ~6 kind family (lake, zone, asset, datascan, catalog entries, glossary); deferred pending demand |
| `google_dataplex_entry_link` | Dataplex data governance judged as a ~6 kind family (lake, zone, asset, datascan, catalog entries, glossary); deferred pending demand |
| `google_dataplex_entry_type` | Dataplex data governance judged as a ~6 kind family (lake, zone, asset, datascan, catalog entries, glossary); deferred pending demand |
| `google_dataplex_glossary` | Dataplex data governance judged as a ~6 kind family (lake, zone, asset, datascan, catalog entries, glossary); deferred pending demand |
| `google_dataplex_glossary_category` | Dataplex data governance judged as a ~6 kind family (lake, zone, asset, datascan, catalog entries, glossary); deferred pending demand |
| `google_dataplex_glossary_term` | Dataplex data governance judged as a ~6 kind family (lake, zone, asset, datascan, catalog entries, glossary); deferred pending demand |
| `google_dataplex_lake` | Dataplex data governance judged as a ~6 kind family (lake, zone, asset, datascan, catalog entries, glossary); deferred pending demand |
| `google_dataplex_metadata_feed` | Dataplex data governance judged as a ~6 kind family (lake, zone, asset, datascan, catalog entries, glossary); deferred pending demand |
| `google_dataplex_task` | Dataplex data governance judged as a ~6 kind family (lake, zone, asset, datascan, catalog entries, glossary); deferred pending demand |
| `google_dataplex_zone` | Dataplex data governance judged as a ~6 kind family (lake, zone, asset, datascan, catalog entries, glossary); deferred pending demand |
| `google_dataproc_batch` | Dataproc batches and jobs are imperative run-once workloads; deferred |
| `google_dataproc_gdc_application_environment` | Dataproc on GDC (Distributed Cloud) is a niche; deferred |
| `google_dataproc_gdc_service_instance` | Dataproc on GDC (Distributed Cloud) is a niche; deferred |
| `google_dataproc_gdc_spark_application` | Dataproc on GDC (Distributed Cloud) is a niche; deferred |
| `google_dataproc_job` | Dataproc batches and jobs are imperative run-once workloads; deferred |
| `google_dataproc_metastore_federation` | judged to deserve a GcpDataprocMetastore kind (federations composed); deferred pending demand |
| `google_dataproc_metastore_service` | judged to deserve a GcpDataprocMetastore kind (federations composed); deferred pending demand |
| `google_dataproc_session_template` | judged to fold into the existing Dataproc kinds' specs (workflow and session templates); the composition is not built |
| `google_dataproc_workflow_template` | judged to fold into the existing Dataproc kinds' specs (workflow and session templates); the composition is not built |
| `google_datastream_connection_profile` | Datastream judged as a stream kind (connection profiles and private connections composed) for CDC; deferred pending demand |
| `google_datastream_private_connection` | Datastream judged as a stream kind (connection profiles and private connections composed) for CDC; deferred pending demand |
| `google_datastream_stream` | Datastream judged as a stream kind (connection profiles and private connections composed) for CDC; deferred pending demand |
| `google_developer_connect_account_connector` | Developer Connect is new; revisit with Cloud Build v2 adoption |
| `google_developer_connect_connection` | Developer Connect is new; revisit with Cloud Build v2 adoption |
| `google_developer_connect_git_repository_link` | Developer Connect is new; revisit with Cloud Build v2 adoption |
| `google_developer_connect_insights_config` | Developer Connect is new; revisit with Cloud Build v2 adoption |
| `google_dialogflow_cx_agent` | Dialogflow CX conversational AI is a specialty (an eventual ~5 kind family); deferred pending demand |
| `google_dialogflow_cx_entity_type` | Dialogflow CX conversational AI is a specialty (an eventual ~5 kind family); deferred pending demand |
| `google_dialogflow_cx_environment` | Dialogflow CX conversational AI is a specialty (an eventual ~5 kind family); deferred pending demand |
| `google_dialogflow_cx_flow` | Dialogflow CX conversational AI is a specialty (an eventual ~5 kind family); deferred pending demand |
| `google_dialogflow_cx_generative_settings` | Dialogflow CX conversational AI is a specialty (an eventual ~5 kind family); deferred pending demand |
| `google_dialogflow_cx_generator` | Dialogflow CX conversational AI is a specialty (an eventual ~5 kind family); deferred pending demand |
| `google_dialogflow_cx_intent` | Dialogflow CX conversational AI is a specialty (an eventual ~5 kind family); deferred pending demand |
| `google_dialogflow_cx_page` | Dialogflow CX conversational AI is a specialty (an eventual ~5 kind family); deferred pending demand |
| `google_dialogflow_cx_playbook` | Dialogflow CX conversational AI is a specialty (an eventual ~5 kind family); deferred pending demand |
| `google_dialogflow_cx_security_settings` | Dialogflow CX conversational AI is a specialty (an eventual ~5 kind family); deferred pending demand |
| `google_dialogflow_cx_test_case` | Dialogflow CX conversational AI is a specialty (an eventual ~5 kind family); deferred pending demand |
| `google_dialogflow_cx_tool` | Dialogflow CX conversational AI is a specialty (an eventual ~5 kind family); deferred pending demand |
| `google_dialogflow_cx_tool_version` | Dialogflow CX conversational AI is a specialty (an eventual ~5 kind family); deferred pending demand |
| `google_dialogflow_cx_version` | Dialogflow CX conversational AI is a specialty (an eventual ~5 kind family); deferred pending demand |
| `google_dialogflow_cx_webhook` | Dialogflow CX conversational AI is a specialty (an eventual ~5 kind family); deferred pending demand |
| `google_discovery_engine_acl_config` | Vertex AI Search (Discovery Engine) is growing but a specialty; first candidate to promote out of deferral |
| `google_discovery_engine_assistant` | Vertex AI Search (Discovery Engine) is growing but a specialty; first candidate to promote out of deferral |
| `google_discovery_engine_chat_engine` | Vertex AI Search (Discovery Engine) is growing but a specialty; first candidate to promote out of deferral |
| `google_discovery_engine_cmek_config` | Vertex AI Search (Discovery Engine) is growing but a specialty; first candidate to promote out of deferral |
| `google_discovery_engine_control` | Vertex AI Search (Discovery Engine) is growing but a specialty; first candidate to promote out of deferral |
| `google_discovery_engine_data_connector` | Vertex AI Search (Discovery Engine) is growing but a specialty; first candidate to promote out of deferral |
| `google_discovery_engine_data_store` | Vertex AI Search (Discovery Engine) is growing but a specialty; first candidate to promote out of deferral |
| `google_discovery_engine_license_config` | Vertex AI Search (Discovery Engine) is growing but a specialty; first candidate to promote out of deferral |
| `google_discovery_engine_recommendation_engine` | Vertex AI Search (Discovery Engine) is growing but a specialty; first candidate to promote out of deferral |
| `google_discovery_engine_schema` | Vertex AI Search (Discovery Engine) is growing but a specialty; first candidate to promote out of deferral |
| `google_discovery_engine_search_engine` | Vertex AI Search (Discovery Engine) is growing but a specialty; first candidate to promote out of deferral |
| `google_discovery_engine_serving_config` | Vertex AI Search (Discovery Engine) is growing but a specialty; first candidate to promote out of deferral |
| `google_discovery_engine_sitemap` | Vertex AI Search (Discovery Engine) is growing but a specialty; first candidate to promote out of deferral |
| `google_discovery_engine_target_site` | Vertex AI Search (Discovery Engine) is growing but a specialty; first candidate to promote out of deferral |
| `google_discovery_engine_user_store` | Vertex AI Search (Discovery Engine) is growing but a specialty; first candidate to promote out of deferral |
| `google_discovery_engine_widget_config` | Vertex AI Search (Discovery Engine) is growing but a specialty; first candidate to promote out of deferral |
| `google_dns_policy` | judged to fold into the existing GcpVpcNetwork kind's spec (per-network DNS server policy); the composition is not built |
| `google_dns_response_policy` | judged to deserve a GcpDnsResponsePolicy kind (DNS firewall, rules composed); deferred pending demand |
| `google_dns_response_policy_rule` | judged to deserve a GcpDnsResponsePolicy kind (DNS firewall, rules composed); deferred pending demand |
| `google_document_ai_processor` | Document AI processors are a specialty; deferred |
| `google_document_ai_processor_default_version` | Document AI processors are a specialty; deferred |
| `google_document_ai_schema` | Document AI processors are a specialty; deferred |
| `google_edgecontainer_cluster` | Distributed Cloud Edge is a niche; deferred |
| `google_edgecontainer_node_pool` | Distributed Cloud Edge is a niche; deferred |
| `google_edgecontainer_vpn_connection` | Distributed Cloud Edge is a niche; deferred |
| `google_edgenetwork_interconnect_attachment` | Distributed Cloud Edge is a niche; deferred |
| `google_edgenetwork_network` | Distributed Cloud Edge is a niche; deferred |
| `google_edgenetwork_subnet` | Distributed Cloud Edge is a niche; deferred |
| `google_essential_contacts_contact` | judged to deserve a GcpEssentialContact kind (organization hygiene); deferred pending demand |
| `google_filestore_backup` | judged to fold into the existing GcpFilestoreInstance kind's spec (backups and snapshots); the composition is not built |
| `google_filestore_snapshot` | judged to fold into the existing GcpFilestoreInstance kind's spec (backups and snapshots); the composition is not built |
| `google_firebase_app_check_app_attest_config` | App Check folds into the Firebase kinds that enter through the beta admission list (the Firebase family is the standing seed); deferred until that admission |
| `google_firebase_app_check_debug_token` | App Check folds into the Firebase kinds that enter through the beta admission list (the Firebase family is the standing seed); deferred until that admission |
| `google_firebase_app_check_device_check_config` | App Check folds into the Firebase kinds that enter through the beta admission list (the Firebase family is the standing seed); deferred until that admission |
| `google_firebase_app_check_play_integrity_config` | App Check folds into the Firebase kinds that enter through the beta admission list (the Firebase family is the standing seed); deferred until that admission |
| `google_firebase_app_check_recaptcha_enterprise_config` | App Check folds into the Firebase kinds that enter through the beta admission list (the Firebase family is the standing seed); deferred until that admission |
| `google_firebase_app_check_recaptcha_v3_config` | App Check folds into the Firebase kinds that enter through the beta admission list (the Firebase family is the standing seed); deferred until that admission |
| `google_firebase_app_check_resource_policy` | App Check folds into the Firebase kinds that enter through the beta admission list (the Firebase family is the standing seed); deferred until that admission |
| `google_firebase_app_check_service_config` | App Check folds into the Firebase kinds that enter through the beta admission list (the Firebase family is the standing seed); deferred until that admission |
| `google_firebase_app_hosting_backend` | Firebase App Hosting judged as a backend kind (builds, traffic, and domains composed); deferred pending demand |
| `google_firebase_app_hosting_build` | Firebase App Hosting judged as a backend kind (builds, traffic, and domains composed); deferred pending demand |
| `google_firebase_app_hosting_default_domain` | Firebase App Hosting judged as a backend kind (builds, traffic, and domains composed); deferred pending demand |
| `google_firebase_app_hosting_domain` | Firebase App Hosting judged as a backend kind (builds, traffic, and domains composed); deferred pending demand |
| `google_firebase_app_hosting_traffic` | Firebase App Hosting judged as a backend kind (builds, traffic, and domains composed); deferred pending demand |
| `google_firebase_data_connect_service` | Firebase Data Connect is new; deferred |
| `google_firebase_remote_config_remote_config` | Remote Config folds into the Firebase kinds that enter through the beta admission list (the Firebase family is the standing seed); deferred until that admission |
| `google_firebaserules_release` | judged to fold into the Firestore database and Firebase storage kinds (security rules); the composition is not built |
| `google_firebaserules_ruleset` | judged to fold into the Firestore database and Firebase storage kinds (security rules); the composition is not built |
| `google_firestore_document` | Firestore documents are data-plane content, not infrastructure |
| `google_firestore_field` | judged to fold into the existing GcpFirestoreIndex and GcpFirestoreDatabase kinds' specs (single-field index configuration); the composition is not built |
| `google_firestore_user_creds` | Firestore user credentials are data-plane auth material, not infrastructure; deferred |
| `google_folder` | judged to deserve a GcpFolder kind (resource hierarchy); deferred pending demand |
| `google_folder_access_approval_settings` | Access Approval settings judged as one kind with project/folder/organization scope; deferred pending demand |
| `google_folder_iam_audit_config` | IAM audit-config surface judged to fold into the project/folder/organization kinds when admitted; not expressible through the additive iam_members pattern today |
| `google_gemini_code_repository_index` | Gemini Code Assist admin settings are a niche; deferred |
| `google_gemini_code_tools_setting` | Gemini Code Assist admin settings are a niche; deferred |
| `google_gemini_code_tools_setting_binding` | Gemini Code Assist admin settings are a niche; deferred |
| `google_gemini_data_sharing_with_google_setting` | Gemini Code Assist admin settings are a niche; deferred |
| `google_gemini_data_sharing_with_google_setting_binding` | Gemini Code Assist admin settings are a niche; deferred |
| `google_gemini_gemini_gcp_enablement_setting` | Gemini Code Assist admin settings are a niche; deferred |
| `google_gemini_gemini_gcp_enablement_setting_binding` | Gemini Code Assist admin settings are a niche; deferred |
| `google_gemini_logging_setting` | Gemini Code Assist admin settings are a niche; deferred |
| `google_gemini_logging_setting_binding` | Gemini Code Assist admin settings are a niche; deferred |
| `google_gemini_release_channel_setting` | Gemini Code Assist admin settings are a niche; deferred |
| `google_gemini_release_channel_setting_binding` | Gemini Code Assist admin settings are a niche; deferred |
| `google_gemini_repository_group` | Gemini Code Assist admin settings are a niche; deferred |
| `google_gke_backup_backup_channel` | Backup for GKE judged as backup-plan and restore-plan kinds (channels composed); deferred pending demand |
| `google_gke_backup_backup_plan` | Backup for GKE judged as backup-plan and restore-plan kinds (channels composed); deferred pending demand |
| `google_gke_backup_restore_channel` | Backup for GKE judged as backup-plan and restore-plan kinds (channels composed); deferred pending demand |
| `google_gke_backup_restore_plan` | Backup for GKE judged as backup-plan and restore-plan kinds (channels composed); deferred pending demand |
| `google_gke_hub_feature` | GKE fleet management judged as fleet, feature, scope, and membership kinds; deferred pending demand |
| `google_gke_hub_feature_membership` | GKE fleet management judged as fleet, feature, scope, and membership kinds; deferred pending demand |
| `google_gke_hub_fleet` | GKE fleet management judged as fleet, feature, scope, and membership kinds; deferred pending demand |
| `google_gke_hub_membership` | GKE fleet management judged as fleet, feature, scope, and membership kinds; deferred pending demand |
| `google_gke_hub_membership_binding` | GKE fleet management judged as fleet, feature, scope, and membership kinds; deferred pending demand |
| `google_gke_hub_namespace` | GKE fleet management judged as fleet, feature, scope, and membership kinds; deferred pending demand |
| `google_gke_hub_rollout_sequence` | GKE fleet management judged as fleet, feature, scope, and membership kinds; deferred pending demand |
| `google_gke_hub_scope` | GKE fleet management judged as fleet, feature, scope, and membership kinds; deferred pending demand |
| `google_gke_hub_scope_rbac_role_binding` | GKE fleet management judged as fleet, feature, scope, and membership kinds; deferred pending demand |
| `google_gkeonprem_bare_metal_admin_cluster` | GKE on-prem (bare metal and VMware) clusters are a niche class; deferred |
| `google_gkeonprem_bare_metal_cluster` | GKE on-prem (bare metal and VMware) clusters are a niche class; deferred |
| `google_gkeonprem_bare_metal_node_pool` | GKE on-prem (bare metal and VMware) clusters are a niche class; deferred |
| `google_gkeonprem_vmware_admin_cluster` | GKE on-prem (bare metal and VMware) clusters are a niche class; deferred |
| `google_gkeonprem_vmware_cluster` | GKE on-prem (bare metal and VMware) clusters are a niche class; deferred |
| `google_gkeonprem_vmware_node_pool` | GKE on-prem (bare metal and VMware) clusters are a niche class; deferred |
| `google_healthcare_consent_store` | Cloud Healthcare (FHIR/DICOM/HL7) is an industry vertical (an eventual ~5 kind family); deferred pending demand |
| `google_healthcare_dataset` | Cloud Healthcare (FHIR/DICOM/HL7) is an industry vertical (an eventual ~5 kind family); deferred pending demand |
| `google_healthcare_dicom_store` | Cloud Healthcare (FHIR/DICOM/HL7) is an industry vertical (an eventual ~5 kind family); deferred pending demand |
| `google_healthcare_fhir_store` | Cloud Healthcare (FHIR/DICOM/HL7) is an industry vertical (an eventual ~5 kind family); deferred pending demand |
| `google_healthcare_hl7_v2_store` | Cloud Healthcare (FHIR/DICOM/HL7) is an industry vertical (an eventual ~5 kind family); deferred pending demand |
| `google_healthcare_pipeline_job` | Cloud Healthcare (FHIR/DICOM/HL7) is an industry vertical (an eventual ~5 kind family); deferred pending demand |
| `google_healthcare_workspace` | Cloud Healthcare (FHIR/DICOM/HL7) is an industry vertical (an eventual ~5 kind family); deferred pending demand |
| `google_hypercomputecluster_cluster` | Hypercompute Cluster AI-supercomputing is a specialty; deferred pending demand |
| `google_iam_deny_policy` | judged to deserve a GcpIamDenyPolicy kind (core enterprise governance); deferred pending demand |
| `google_iam_folders_policy_binding` | judged to deserve a GcpPrincipalAccessBoundaryPolicy kind (scope bindings composed); deferred pending demand |
| `google_iam_organizations_policy_binding` | judged to deserve a GcpPrincipalAccessBoundaryPolicy kind (scope bindings composed); deferred pending demand |
| `google_iam_principal_access_boundary_policy` | judged to deserve a GcpPrincipalAccessBoundaryPolicy kind (scope bindings composed); deferred pending demand |
| `google_iam_projects_policy_binding` | judged to deserve a GcpPrincipalAccessBoundaryPolicy kind (scope bindings composed); deferred pending demand |
| `google_iam_workforce_pool` | workforce identity federation judged as workforce-pool and provider kinds (provider keys composed); deferred pending demand |
| `google_iam_workforce_pool_provider` | workforce identity federation judged as workforce-pool and provider kinds (provider keys composed); deferred pending demand |
| `google_iam_workforce_pool_provider_key` | workforce identity federation judged as workforce-pool and provider kinds (provider keys composed); deferred pending demand |
| `google_iam_workforce_pool_provider_scim_tenant` | workforce identity federation judged as workforce-pool and provider kinds (provider keys composed); deferred pending demand |
| `google_iam_workforce_pool_provider_scim_token` | workforce identity federation judged as workforce-pool and provider kinds (provider keys composed); deferred pending demand |
| `google_iam_workload_identity_pool_managed_identity` | judged to fold into the existing GcpWorkloadIdentityPool kind's spec (namespaces and managed identities); the composition is not built |
| `google_iam_workload_identity_pool_namespace` | judged to fold into the existing GcpWorkloadIdentityPool kind's spec (namespaces and managed identities); the composition is not built |
| `google_integration_connectors_connection` | Application Integration is a niche; deferred |
| `google_integration_connectors_endpoint_attachment` | Application Integration is a niche; deferred |
| `google_integration_connectors_managed_zone` | Application Integration is a niche; deferred |
| `google_integrations_auth_config` | Application Integration is a niche; deferred |
| `google_integrations_client` | Application Integration is a niche; deferred |
| `google_kms_autokey_config` | KMS Autokey judged as a GcpKmsAutokey kind (autokey config and key handles composed); deferred pending demand |
| `google_kms_crypto_key_version` | judged to fold into the existing GcpKmsKey kind's spec (key versions); the composition is not built |
| `google_kms_ekm_connection` | external key manager connections are a specialty; deferred |
| `google_kms_key_handle` | KMS Autokey judged as a GcpKmsAutokey kind (autokey config and key handles composed); deferred pending demand |
| `google_kms_key_ring_import_job` | judged to fold into the existing GcpKmsKeyRing kind's spec (import jobs); the composition is not built |
| `google_kms_project_autokey_config` | KMS Autokey judged as a GcpKmsAutokey kind (autokey config and key handles composed); deferred pending demand |
| `google_kms_secret_ciphertext` | secret ciphertext is an imperative encrypt operation, a poor declarative fit |
| `google_license_manager_configuration` | License Manager is a specialty; deferred pending demand |
| `google_logging_log_scope` | log scopes are console conveniences, rarely IaC-managed; deferred |
| `google_logging_saved_query` | saved queries are console artifacts, not provisioned infrastructure; deferred |
| `google_looker_instance` | judged to deserve a GcpLookerInstance kind; deferred pending demand |
| `google_lustre_instance` | Managed Lustre is an HPC niche; deferred |
| `google_managed_kafka_acl` | Managed Kafka judged as a cluster kind (topics and ACLs composed); deferred pending demand |
| `google_managed_kafka_cluster` | Managed Kafka judged as a cluster kind (topics and ACLs composed); deferred pending demand |
| `google_managed_kafka_connect_cluster` | Managed Kafka judged as a cluster kind (topics and ACLs composed); deferred pending demand |
| `google_managed_kafka_connector` | Managed Kafka judged as a cluster kind (topics and ACLs composed); deferred pending demand |
| `google_managed_kafka_topic` | Managed Kafka judged as a cluster kind (topics and ACLs composed); deferred pending demand |
| `google_memcache_instance` | Memorystore Memcached is fading relative to Redis/Valkey; deferred |
| `google_memorystore_instance_desired_user_created_endpoints` | judged to fold into the existing GcpMemorystoreInstance kind's spec (user-created endpoint connections); the composition is not built |
| `google_migration_center_assets_export_job` | Migration Center assessment tooling is episodic; deferred |
| `google_migration_center_discovery_client` | Migration Center assessment tooling is episodic; deferred |
| `google_migration_center_group` | Migration Center assessment tooling is episodic; deferred |
| `google_migration_center_import_data_file` | Migration Center assessment tooling is episodic; deferred |
| `google_migration_center_import_job` | Migration Center assessment tooling is episodic; deferred |
| `google_migration_center_preference_set` | Migration Center assessment tooling is episodic; deferred |
| `google_migration_center_report` | Migration Center assessment tooling is episodic; deferred |
| `google_migration_center_report_config` | Migration Center assessment tooling is episodic; deferred |
| `google_migration_center_settings` | Migration Center assessment tooling is episodic; deferred |
| `google_migration_center_source` | Migration Center assessment tooling is episodic; deferred |
| `google_model_armor_floorsetting` | Model Armor prompt-safety templates are new; deferred |
| `google_model_armor_template` | Model Armor prompt-safety templates are new; deferred |
| `google_monitoring_metric_descriptor` | metric descriptors are rarely hand-managed; deferred |
| `google_netapp_active_directory` | NetApp Volumes judged as storage-pool, volume, backup-vault, and backup-policy kinds (companions composed); deferred pending demand |
| `google_netapp_backup` | NetApp Volumes judged as storage-pool, volume, backup-vault, and backup-policy kinds (companions composed); deferred pending demand |
| `google_netapp_backup_policy` | NetApp Volumes judged as storage-pool, volume, backup-vault, and backup-policy kinds (companions composed); deferred pending demand |
| `google_netapp_backup_vault` | NetApp Volumes judged as storage-pool, volume, backup-vault, and backup-policy kinds (companions composed); deferred pending demand |
| `google_netapp_host_group` | NetApp Volumes judged as storage-pool, volume, backup-vault, and backup-policy kinds (companions composed); deferred pending demand |
| `google_netapp_kmsconfig` | NetApp Volumes judged as storage-pool, volume, backup-vault, and backup-policy kinds (companions composed); deferred pending demand |
| `google_netapp_storage_pool` | NetApp Volumes judged as storage-pool, volume, backup-vault, and backup-policy kinds (companions composed); deferred pending demand |
| `google_netapp_volume` | NetApp Volumes judged as storage-pool, volume, backup-vault, and backup-policy kinds (companions composed); deferred pending demand |
| `google_netapp_volume_quota_rule` | NetApp Volumes judged as storage-pool, volume, backup-vault, and backup-policy kinds (companions composed); deferred pending demand |
| `google_netapp_volume_replication` | NetApp Volumes judged as storage-pool, volume, backup-vault, and backup-policy kinds (companions composed); deferred pending demand |
| `google_netapp_volume_snapshot` | NetApp Volumes judged as storage-pool, volume, backup-vault, and backup-policy kinds (companions composed); deferred pending demand |
| `google_network_connectivity_destination` | Network Connectivity Center judged as hub, spoke, and policy-based-route kinds (groups composed); deferred pending demand |
| `google_network_connectivity_gateway_advertised_route` | Network Connectivity Center judged as hub, spoke, and policy-based-route kinds (groups composed); deferred pending demand |
| `google_network_connectivity_group` | Network Connectivity Center judged as hub, spoke, and policy-based-route kinds (groups composed); deferred pending demand |
| `google_network_connectivity_hub` | Network Connectivity Center judged as hub, spoke, and policy-based-route kinds (groups composed); deferred pending demand |
| `google_network_connectivity_internal_range` | judged to fold into the existing GcpVpcNetwork kind's spec (internal ranges); the composition is not built |
| `google_network_connectivity_multicloud_data_transfer_config` | Network Connectivity Center judged as hub, spoke, and policy-based-route kinds (groups composed); deferred pending demand |
| `google_network_connectivity_policy_based_route` | Network Connectivity Center judged as hub, spoke, and policy-based-route kinds (groups composed); deferred pending demand |
| `google_network_connectivity_regional_endpoint` | regional endpoints are new Network Connectivity surface; deferred |
| `google_network_connectivity_spoke` | Network Connectivity Center judged as hub, spoke, and policy-based-route kinds (groups composed); deferred pending demand |
| `google_network_connectivity_transport` | Network Connectivity Center judged as hub, spoke, and policy-based-route kinds (groups composed); deferred pending demand |
| `google_network_management_connectivity_test` | connectivity tests are diagnostics, not infrastructure; deferred |
| `google_network_management_organization_vpc_flow_logs_config` | organization-scoped VPC flow-logs config is org-admin surface; deferred pending demand |
| `google_network_management_vpc_flow_logs_config` | judged to fold into the VPC network family's specs (VPC flow-logs configuration); the composition is not built |
| `google_network_security_address_group` | NGFW and TLS policy surface judged as firewall-endpoint, security-profile, address-group, and TLS-policy kinds; deferred pending demand |
| `google_network_security_authz_policy` | judged to fold into the load-balancer facade kinds when they are admitted (authorization and backend-authentication policies); deferred |
| `google_network_security_backend_authentication_config` | judged to fold into the load-balancer facade kinds when they are admitted (authorization and backend-authentication policies); deferred |
| `google_network_security_client_tls_policy` | NGFW and TLS policy surface judged as firewall-endpoint, security-profile, address-group, and TLS-policy kinds; deferred pending demand |
| `google_network_security_dns_threat_detector` | NGFW and TLS policy surface judged as firewall-endpoint, security-profile, address-group, and TLS-policy kinds; deferred pending demand |
| `google_network_security_firewall_endpoint` | NGFW and TLS policy surface judged as firewall-endpoint, security-profile, address-group, and TLS-policy kinds; deferred pending demand |
| `google_network_security_firewall_endpoint_association` | NGFW and TLS policy surface judged as firewall-endpoint, security-profile, address-group, and TLS-policy kinds; deferred pending demand |
| `google_network_security_gateway_security_policy` | Secure Web Gateway and TLS inspection are a specialty; deferred |
| `google_network_security_gateway_security_policy_rule` | Secure Web Gateway and TLS inspection are a specialty; deferred |
| `google_network_security_intercept_deployment` | packet intercept and mirroring v2 surface is new; deferred |
| `google_network_security_intercept_deployment_group` | packet intercept and mirroring v2 surface is new; deferred |
| `google_network_security_intercept_endpoint_group` | packet intercept and mirroring v2 surface is new; deferred |
| `google_network_security_intercept_endpoint_group_association` | packet intercept and mirroring v2 surface is new; deferred |
| `google_network_security_mirroring_deployment` | packet intercept and mirroring v2 surface is new; deferred |
| `google_network_security_mirroring_deployment_group` | packet intercept and mirroring v2 surface is new; deferred |
| `google_network_security_mirroring_endpoint` | NGFW and TLS policy surface judged as firewall-endpoint, security-profile, address-group, and TLS-policy kinds; deferred pending demand |
| `google_network_security_mirroring_endpoint_group` | packet intercept and mirroring v2 surface is new; deferred |
| `google_network_security_mirroring_endpoint_group_association` | packet intercept and mirroring v2 surface is new; deferred |
| `google_network_security_security_profile` | NGFW and TLS policy surface judged as firewall-endpoint, security-profile, address-group, and TLS-policy kinds; deferred pending demand |
| `google_network_security_security_profile_group` | NGFW and TLS policy surface judged as firewall-endpoint, security-profile, address-group, and TLS-policy kinds; deferred pending demand |
| `google_network_security_server_tls_policy` | NGFW and TLS policy surface judged as firewall-endpoint, security-profile, address-group, and TLS-policy kinds; deferred pending demand |
| `google_network_security_tls_inspection_policy` | Secure Web Gateway and TLS inspection are a specialty; deferred |
| `google_network_security_ull_mirroring_collector` | NGFW and TLS policy surface judged as firewall-endpoint, security-profile, address-group, and TLS-policy kinds; deferred pending demand |
| `google_network_security_ull_mirroring_collector_rule` | NGFW and TLS policy surface judged as firewall-endpoint, security-profile, address-group, and TLS-policy kinds; deferred pending demand |
| `google_network_security_ull_mirroring_engine` | NGFW and TLS policy surface judged as firewall-endpoint, security-profile, address-group, and TLS-policy kinds; deferred pending demand |
| `google_network_security_url_lists` | Secure Web Gateway and TLS inspection are a specialty; deferred |
| `google_network_services_agent_gateway` | agent gateways are emerging network-services surface; deferred pending demand |
| `google_network_services_authz_extension` | judged to fold into the load-balancer facade kinds when they are admitted (service extensions and callouts); deferred |
| `google_network_services_edge_cache_keyset` | Media CDN is allowlist-only; deferred |
| `google_network_services_edge_cache_origin` | Media CDN is allowlist-only; deferred |
| `google_network_services_edge_cache_service` | Media CDN is allowlist-only; deferred |
| `google_network_services_endpoint_policy` | Cloud Service Mesh resources; the Kubernetes-native path is preferred today; deferred |
| `google_network_services_gateway` | Cloud Service Mesh resources; the Kubernetes-native path is preferred today; deferred |
| `google_network_services_grpc_route` | Cloud Service Mesh resources; the Kubernetes-native path is preferred today; deferred |
| `google_network_services_http_route` | Cloud Service Mesh resources; the Kubernetes-native path is preferred today; deferred |
| `google_network_services_lb_edge_extension` | judged to fold into the load-balancer facade kinds when they are admitted (service extensions and callouts); deferred |
| `google_network_services_lb_route_extension` | judged to fold into the load-balancer facade kinds when they are admitted (service extensions and callouts); deferred |
| `google_network_services_lb_traffic_extension` | judged to fold into the load-balancer facade kinds when they are admitted (service extensions and callouts); deferred |
| `google_network_services_mesh` | Cloud Service Mesh resources; the Kubernetes-native path is preferred today; deferred |
| `google_network_services_multicast_consumer_association` | Multicast networking is specialty; deferred pending demand |
| `google_network_services_multicast_domain` | Multicast networking is specialty; deferred pending demand |
| `google_network_services_multicast_domain_activation` | Multicast networking is specialty; deferred pending demand |
| `google_network_services_multicast_domain_group` | Multicast networking is specialty; deferred pending demand |
| `google_network_services_multicast_group_consumer_activation` | Multicast networking is specialty; deferred pending demand |
| `google_network_services_multicast_group_producer_activation` | Multicast networking is specialty; deferred pending demand |
| `google_network_services_multicast_group_range` | Multicast networking is specialty; deferred pending demand |
| `google_network_services_multicast_group_range_activation` | Multicast networking is specialty; deferred pending demand |
| `google_network_services_multicast_producer_association` | Multicast networking is specialty; deferred pending demand |
| `google_network_services_tcp_route` | Cloud Service Mesh resources; the Kubernetes-native path is preferred today; deferred |
| `google_network_services_tls_route` | Cloud Service Mesh resources; the Kubernetes-native path is preferred today; deferred |
| `google_network_services_wasm_plugin` | judged to fold into the load-balancer facade kinds when they are admitted (service extensions and callouts); deferred |
| `google_observability_trace_scope` | observability scopes are console organization, not provisioned infrastructure; deferred |
| `google_oracle_database_autonomous_database` | Oracle Database at Google Cloud is a named niche family (an eventual ~6 kind family); deferred pending demand |
| `google_oracle_database_cloud_exadata_infrastructure` | Oracle Database at Google Cloud is a named niche family (an eventual ~6 kind family); deferred pending demand |
| `google_oracle_database_cloud_exadata_infrastructure_exascale_config` | Oracle Database at Google Cloud is a named niche family (an eventual ~6 kind family); deferred pending demand |
| `google_oracle_database_cloud_vm_cluster` | Oracle Database at Google Cloud is a named niche family (an eventual ~6 kind family); deferred pending demand |
| `google_oracle_database_db_system` | Oracle Database at Google Cloud is a named niche family (an eventual ~6 kind family); deferred pending demand |
| `google_oracle_database_exadb_vm_cluster` | Oracle Database at Google Cloud is a named niche family (an eventual ~6 kind family); deferred pending demand |
| `google_oracle_database_exascale_db_storage_vault` | Oracle Database at Google Cloud is a named niche family (an eventual ~6 kind family); deferred pending demand |
| `google_oracle_database_goldengate_connection` | Oracle Database at Google Cloud is a named niche family (an eventual ~6 kind family); deferred pending demand |
| `google_oracle_database_goldengate_connection_assignment` | Oracle Database at Google Cloud is a named niche family (an eventual ~6 kind family); deferred pending demand |
| `google_oracle_database_goldengate_deployment` | Oracle Database at Google Cloud is a named niche family (an eventual ~6 kind family); deferred pending demand |
| `google_oracle_database_odb_network` | Oracle Database at Google Cloud is a named niche family (an eventual ~6 kind family); deferred pending demand |
| `google_oracle_database_odb_subnet` | Oracle Database at Google Cloud is a named niche family (an eventual ~6 kind family); deferred pending demand |
| `google_org_policy_custom_constraint` | Organization Policy judged as policy and custom-constraint kinds (core enterprise governance); deferred pending demand |
| `google_org_policy_policy` | Organization Policy judged as policy and custom-constraint kinds (core enterprise governance); deferred pending demand |
| `google_organization_access_approval_settings` | Access Approval settings judged as one kind with project/folder/organization scope; deferred pending demand |
| `google_organization_iam_audit_config` | IAM audit-config surface judged to fold into the project/folder/organization kinds when admitted; not expressible through the additive iam_members pattern today |
| `google_organization_iam_custom_role` | judged to fold into the existing GcpIamCustomRole kind's spec (organization scope); the composition is not built |
| `google_os_config_os_policy_assignment` | OS Config judged as OS-policy-assignment and patch-deployment kinds; deferred pending demand |
| `google_os_config_patch_deployment` | OS Config judged as OS-policy-assignment and patch-deployment kinds; deferred pending demand |
| `google_os_config_v2_policy_orchestrator` | policy orchestrators fold into the OS policy kind when it is admitted; deferred |
| `google_os_config_v2_policy_orchestrator_for_folder` | policy orchestrators fold into the OS policy kind when it is admitted; deferred |
| `google_os_config_v2_policy_orchestrator_for_organization` | policy orchestrators fold into the OS policy kind when it is admitted; deferred |
| `google_os_login_ssh_public_key` | per-user SSH keys are user data, not platform infrastructure; deferred |
| `google_parallelstore_instance` | Parallelstore is an HPC niche; deferred |
| `google_parameter_manager_parameter` | Parameter Manager judged as one parameter kind with versions composed and regional variants as a location flag (the Secret Manager pattern); deferred pending demand |
| `google_parameter_manager_parameter_version` | Parameter Manager judged as one parameter kind with versions composed and regional variants as a location flag (the Secret Manager pattern); deferred pending demand |
| `google_parameter_manager_regional_parameter` | Parameter Manager judged as one parameter kind with versions composed and regional variants as a location flag (the Secret Manager pattern); deferred pending demand |
| `google_parameter_manager_regional_parameter_version` | Parameter Manager judged as one parameter kind with versions composed and regional variants as a location flag (the Secret Manager pattern); deferred pending demand |
| `google_privateca_ca_pool` | Private CA judged as a CA-pool kind (pool with certificate authorities and leaf issuance composed) and a certificate-template kind; deferred pending demand |
| `google_privateca_certificate` | Private CA judged as a CA-pool kind (pool with certificate authorities and leaf issuance composed) and a certificate-template kind; deferred pending demand |
| `google_privateca_certificate_authority` | Private CA judged as a CA-pool kind (pool with certificate authorities and leaf issuance composed) and a certificate-template kind; deferred pending demand |
| `google_privateca_certificate_template` | Private CA judged as a CA-pool kind (pool with certificate authorities and leaf issuance composed) and a certificate-template kind; deferred pending demand |
| `google_privileged_access_manager_entitlement` | judged to deserve a just-in-time privileged-access entitlement kind; deferred pending demand |
| `google_project_access_approval_settings` | Access Approval settings judged as one kind with project/folder/organization scope; deferred pending demand |
| `google_project_default_service_accounts` | judged to fold into the existing GcpProject kind's spec (default service-account posture); the composition is not built |
| `google_project_iam_audit_config` | IAM audit-config surface judged to fold into the project/folder/organization kinds when admitted; not expressible through the additive iam_members pattern today |
| `google_project_usage_export_bucket` | judged to fold into the existing GcpProject kind's spec (compute usage-export bucket); the composition is not built |
| `google_public_ca_external_account_key` | ACME external account keys are a niche; deferred |
| `google_recaptcha_enterprise_key` | judged to deserve a GcpRecaptchaKey kind; deferred pending demand |
| `google_redis_cluster` | judged to deserve a GcpRedisCluster kind (Memorystore cluster tier, user-created connections composed); deferred pending demand |
| `google_redis_cluster_user_created_connections` | judged to deserve a GcpRedisCluster kind (Memorystore cluster tier, user-created connections composed); deferred pending demand |
| `google_resource_manager_capability` | Resource Manager capabilities are org-admin toggles; deferred pending demand |
| `google_resource_manager_lien` | judged to fold into the existing GcpProject kind's spec (liens); the composition is not built |
| `google_scc_management_folder_security_health_analytics_custom_module` | SCC Management custom modules judged as SHA and ETD custom-module kinds with scope selectors; deferred pending demand |
| `google_scc_management_organization_event_threat_detection_custom_module` | SCC Management custom modules judged as SHA and ETD custom-module kinds with scope selectors; deferred pending demand |
| `google_scc_management_organization_security_health_analytics_custom_module` | SCC Management custom modules judged as SHA and ETD custom-module kinds with scope selectors; deferred pending demand |
| `google_scc_management_project_security_health_analytics_custom_module` | SCC Management custom modules judged as SHA and ETD custom-module kinds with scope selectors; deferred pending demand |
| `google_scc_v2_folder_mute_config` | SCC v2 judged as notification-config, mute-config, and BigQuery-export kinds with project/folder/organization scope; deferred pending demand |
| `google_scc_v2_folder_notification_config` | SCC v2 judged as notification-config, mute-config, and BigQuery-export kinds with project/folder/organization scope; deferred pending demand |
| `google_scc_v2_folder_scc_big_query_export` | SCC v2 judged as notification-config, mute-config, and BigQuery-export kinds with project/folder/organization scope; deferred pending demand |
| `google_scc_v2_organization_mute_config` | SCC v2 judged as notification-config, mute-config, and BigQuery-export kinds with project/folder/organization scope; deferred pending demand |
| `google_scc_v2_organization_notification_config` | SCC v2 judged as notification-config, mute-config, and BigQuery-export kinds with project/folder/organization scope; deferred pending demand |
| `google_scc_v2_organization_scc_big_query_export` | SCC v2 judged as notification-config, mute-config, and BigQuery-export kinds with project/folder/organization scope; deferred pending demand |
| `google_scc_v2_organization_source` | SCC v2 judged as notification-config, mute-config, and BigQuery-export kinds with project/folder/organization scope; deferred pending demand |
| `google_scc_v2_project_mute_config` | SCC v2 judged as notification-config, mute-config, and BigQuery-export kinds with project/folder/organization scope; deferred pending demand |
| `google_scc_v2_project_notification_config` | SCC v2 judged as notification-config, mute-config, and BigQuery-export kinds with project/folder/organization scope; deferred pending demand |
| `google_scc_v2_project_scc_big_query_export` | SCC v2 judged as notification-config, mute-config, and BigQuery-export kinds with project/folder/organization scope; deferred pending demand |
| `google_secure_source_manager_branch_rule` | Secure Source Manager is a niche; deferred |
| `google_secure_source_manager_hook` | Secure Source Manager is a niche; deferred |
| `google_secure_source_manager_instance` | Secure Source Manager is a niche; deferred |
| `google_secure_source_manager_repository` | Secure Source Manager is a niche; deferred |
| `google_securityposture_posture` | Security Posture deployments are new; deferred |
| `google_securityposture_posture_deployment` | Security Posture deployments are new; deferred |
| `google_service_directory_endpoint` | Service Directory judged as namespace, service, and endpoint kinds; modest usage today; deferred pending demand |
| `google_service_directory_namespace` | Service Directory judged as namespace, service, and endpoint kinds; modest usage today; deferred pending demand |
| `google_service_directory_service` | Service Directory judged as namespace, service, and endpoint kinds; modest usage today; deferred pending demand |
| `google_service_networking_peered_dns_domain` | judged to fold into the existing GcpServiceNetworkingConnection kind's spec (peered DNS domains and VPC service controls); the composition is not built |
| `google_service_networking_vpc_service_controls` | judged to fold into the existing GcpServiceNetworkingConnection kind's spec (peered DNS domains and VPC service controls); the composition is not built |
| `google_site_verification_owner` | site verification is a niche; deferred |
| `google_site_verification_web_resource` | site verification is a niche; deferred |
| `google_spanner_instance_config` | judged to fold into the existing GcpSpannerInstance kind's spec (custom instance configs and partitions); the composition is not built |
| `google_spanner_instance_partition` | judged to fold into the existing GcpSpannerInstance kind's spec (custom instance configs and partitions); the composition is not built |
| `google_sql_provision_script` | Cloud SQL provision scripts are data-plane bootstrap, not infrastructure shape; deferred |
| `google_sql_source_representation_instance` | source representation instances serve episodic migrations; deferred |
| `google_sql_ssl_cert` | judged to fold into the existing GcpCloudSql kind's spec (client SSL certificates); the composition is not built |
| `google_storage_anywhere_cache` | Anywhere Cache is a niche; deferred |
| `google_storage_batch_operations_job` | batch object operations are imperative; deferred |
| `google_storage_control_folder_intelligence_config` | storage intelligence configs are new; deferred |
| `google_storage_control_organization_intelligence_config` | storage intelligence configs are new; deferred |
| `google_storage_control_project_intelligence_config` | storage intelligence configs are new; deferred |
| `google_storage_insights_dataset_config` | storage inventory reports are a niche; deferred |
| `google_storage_insights_report_config` | storage inventory reports are a niche; deferred |
| `google_storage_transfer_agent_pool` | judged to deserve a GcpStorageTransferJob kind (agent pools composed); deferred pending demand |
| `google_storage_transfer_job` | judged to deserve a GcpStorageTransferJob kind (agent pools composed); deferred pending demand |
| `google_tags_location_tag_binding` | resource tags judged as tag-key and tag-value kinds with bindings composed onto target kinds; deferred pending demand |
| `google_tags_tag_binding` | resource tags judged as tag-key and tag-value kinds with bindings composed onto target kinds; deferred pending demand |
| `google_tags_tag_key` | resource tags judged as tag-key and tag-value kinds with bindings composed onto target kinds; deferred pending demand |
| `google_tags_tag_value` | resource tags judged as tag-key and tag-value kinds with bindings composed onto target kinds; deferred pending demand |
| `google_transcoder_job` | media transcoding jobs are imperative and a niche; deferred |
| `google_transcoder_job_template` | media transcoding jobs are imperative and a niche; deferred |
| `google_vector_search_collection` | standalone Vector Search is new; judged to fold into the Vertex AI family when it stabilizes; deferred pending demand |
| `google_vector_search_data_object` | standalone Vector Search is new; judged to fold into the Vertex AI family when it stabilizes; deferred pending demand |
| `google_vector_search_index` | standalone Vector Search is new; judged to fold into the Vertex AI family when it stabilizes; deferred pending demand |
| `google_vertex_ai_cache_config` | Vertex AI datasets, tensorboards, and RAG engine configuration are specialty surfaces; deferred |
| `google_vertex_ai_dataset` | Vertex AI datasets, tensorboards, and RAG engine configuration are specialty surfaces; deferred |
| `google_vertex_ai_deployment_resource_pool` | judged to fold into the existing GcpVertexAiEndpoint kind's spec (deployment resource pools and Model Garden deployments); the composition is not built |
| `google_vertex_ai_endpoint_with_model_garden_deployment` | judged to fold into the existing GcpVertexAiEndpoint kind's spec (deployment resource pools and Model Garden deployments); the composition is not built |
| `google_vertex_ai_feature_group` | Vertex AI feature platform judged as feature-group and feature-online-store kinds (features and feature views composed); deferred pending demand |
| `google_vertex_ai_feature_group_feature` | Vertex AI feature platform judged as feature-group and feature-online-store kinds (features and feature views composed); deferred pending demand |
| `google_vertex_ai_feature_online_store` | Vertex AI feature platform judged as feature-group and feature-online-store kinds (features and feature views composed); deferred pending demand |
| `google_vertex_ai_feature_online_store_featureview` | Vertex AI feature platform judged as feature-group and feature-online-store kinds (features and feature views composed); deferred pending demand |
| `google_vertex_ai_persistent_resource` | Vertex AI datasets, tensorboards, and RAG engine configuration are specialty surfaces; deferred |
| `google_vertex_ai_rag_engine_config` | Vertex AI datasets, tensorboards, and RAG engine configuration are specialty surfaces; deferred |
| `google_vertex_ai_reasoning_engine` | Vertex AI datasets, tensorboards, and RAG engine configuration are specialty surfaces; deferred |
| `google_vertex_ai_semantic_governance_policy_engine` | Vertex AI datasets, tensorboards, and RAG engine configuration are specialty surfaces; deferred |
| `google_vertex_ai_tensorboard` | Vertex AI datasets, tensorboards, and RAG engine configuration are specialty surfaces; deferred |
| `google_vertex_ai_tensorboard_experiment` | Vertex AI datasets, tensorboards, and RAG engine configuration are specialty surfaces; deferred |
| `google_vertex_ai_tensorboard_run` | Vertex AI datasets, tensorboards, and RAG engine configuration are specialty surfaces; deferred |
| `google_vmwareengine_cluster` | VMware Engine is a named niche family (an eventual ~5 kind family); deferred pending demand |
| `google_vmwareengine_datastore` | VMware Engine is a named niche family (an eventual ~5 kind family); deferred pending demand |
| `google_vmwareengine_external_access_rule` | VMware Engine is a named niche family (an eventual ~5 kind family); deferred pending demand |
| `google_vmwareengine_external_address` | VMware Engine is a named niche family (an eventual ~5 kind family); deferred pending demand |
| `google_vmwareengine_network` | VMware Engine is a named niche family (an eventual ~5 kind family); deferred pending demand |
| `google_vmwareengine_network_peering` | VMware Engine is a named niche family (an eventual ~5 kind family); deferred pending demand |
| `google_vmwareengine_network_policy` | VMware Engine is a named niche family (an eventual ~5 kind family); deferred pending demand |
| `google_vmwareengine_private_cloud` | VMware Engine is a named niche family (an eventual ~5 kind family); deferred pending demand |
| `google_vmwareengine_subnet` | VMware Engine is a named niche family (an eventual ~5 kind family); deferred pending demand |
| `google_workload_identity_service_agent` | service agents are auto-provisioned by service enablement; deferred |
| `google_workstations_workstation` | Cloud Workstations judged as workstation-cluster (with config) and workstation kinds; deferred pending demand |
| `google_workstations_workstation_cluster` | Cloud Workstations judged as workstation-cluster (with config) and workstation kinds; deferred pending demand |
| `google_workstations_workstation_config` | Cloud Workstations judged as workstation-cluster (with config) and workstation kinds; deferred pending demand |

### Excluded as deprecated (76)

| Resource | Recorded reason |
|---|---|
| `google_app_engine_application` | App Engine is a legacy PaaS surface superseded by Cloud Run (containers) and Firebase App Hosting (web apps) |
| `google_app_engine_application_url_dispatch_rules` | App Engine is a legacy PaaS surface superseded by Cloud Run (containers) and Firebase App Hosting (web apps) |
| `google_app_engine_domain_mapping` | App Engine is a legacy PaaS surface superseded by Cloud Run (containers) and Firebase App Hosting (web apps) |
| `google_app_engine_firewall_rule` | App Engine is a legacy PaaS surface superseded by Cloud Run (containers) and Firebase App Hosting (web apps) |
| `google_app_engine_flexible_app_version` | App Engine is a legacy PaaS surface superseded by Cloud Run (containers) and Firebase App Hosting (web apps) |
| `google_app_engine_service_network_settings` | App Engine is a legacy PaaS surface superseded by Cloud Run (containers) and Firebase App Hosting (web apps) |
| `google_app_engine_service_split_traffic` | App Engine is a legacy PaaS surface superseded by Cloud Run (containers) and Firebase App Hosting (web apps) |
| `google_app_engine_standard_app_version` | App Engine is a legacy PaaS surface superseded by Cloud Run (containers) and Firebase App Hosting (web apps) |
| `google_bigquery_datapolicy_data_policy` | superseded by the BigQuery Data Policy v2 API |
| `google_blockchain_node_engine_blockchain_nodes` | Blockchain Node Engine is sunset by Google |
| `google_cloud_run_service` | first-generation Knative-serving Cloud Run resource, superseded by the v2 API (modeled by GcpCloudRun) |
| `google_cloudbuild_bitbucket_server_config` | first-generation Bitbucket Server integration, superseded by Cloud Build v2 connections |
| `google_cloudfunctions_function` | first-generation Cloud Functions, superseded by Cloud Functions 2nd gen (modeled by GcpCloudFunction) |
| `google_compute_http_health_check` | legacy HTTP/HTTPS health checks, superseded by google_compute_health_check (modeled by GcpHealthCheck) |
| `google_compute_https_health_check` | legacy HTTP/HTTPS health checks, superseded by google_compute_health_check (modeled by GcpHealthCheck) |
| `google_compute_target_pool` | legacy network load-balancing target pools, superseded by regional backend services |
| `google_compute_vpn_gateway` | Classic VPN gateway, superseded by HA VPN |
| `google_container_registry` | deprecated in the provider schema |
| `google_data_catalog_entry` | deprecated in the provider schema |
| `google_data_catalog_entry_group` | deprecated in the provider schema |
| `google_data_catalog_policy_tag` | Data Catalog is deprecated by Google, superseded by Dataplex Catalog |
| `google_data_catalog_tag` | deprecated in the provider schema |
| `google_data_catalog_tag_template` | deprecated in the provider schema |
| `google_data_catalog_taxonomy` | Data Catalog is deprecated by Google, superseded by Dataplex Catalog |
| `google_deployment_manager_deployment` | Deployment Manager is legacy and sunsetting; competing IaC deployment engines are out of scope by design |
| `google_dialogflow_agent` | Dialogflow ES is superseded by Dialogflow CX; the ES API is in maintenance |
| `google_dialogflow_conversation_profile` | Dialogflow ES is superseded by Dialogflow CX; the ES API is in maintenance |
| `google_dialogflow_encryption_spec` | Dialogflow ES is superseded by Dialogflow CX; the ES API is in maintenance |
| `google_dialogflow_entity_type` | Dialogflow ES is superseded by Dialogflow CX; the ES API is in maintenance |
| `google_dialogflow_environment` | Dialogflow ES is superseded by Dialogflow CX; the ES API is in maintenance |
| `google_dialogflow_fulfillment` | Dialogflow ES is superseded by Dialogflow CX; the ES API is in maintenance |
| `google_dialogflow_generator` | Dialogflow ES is superseded by Dialogflow CX; the ES API is in maintenance |
| `google_dialogflow_intent` | Dialogflow ES is superseded by Dialogflow CX; the ES API is in maintenance |
| `google_dialogflow_sip_trunk` | Dialogflow ES is superseded by Dialogflow CX; the ES API is in maintenance |
| `google_dialogflow_version` | Dialogflow ES is superseded by Dialogflow CX; the ES API is in maintenance |
| `google_document_ai_warehouse_document_schema` | Document AI Warehouse is deprecated by Google |
| `google_document_ai_warehouse_location` | Document AI Warehouse is deprecated by Google |
| `google_endpoints_service` | Cloud Endpoints with ESP is a legacy surface; Google recommends API Gateway |
| `google_folder_organization_policy` | legacy organization-policy resources superseded by the Organization Policy API (google_org_policy_policy) |
| `google_iam_access_boundary_policy` | superseded by principal access boundary policies (google_iam_principal_access_boundary_policy) |
| `google_iap_brand` | deprecated in the provider schema |
| `google_iap_client` | deprecated in the provider schema |
| `google_ml_engine_model` | deprecated in the provider schema |
| `google_network_services_service_binding` | deprecated in the provider schema |
| `google_notebooks_environment` | deprecated in the provider schema |
| `google_notebooks_instance` | deprecated in the provider schema |
| `google_notebooks_runtime` | deprecated in the provider schema |
| `google_organization_policy` | legacy organization-policy resources superseded by the Organization Policy API (google_org_policy_policy) |
| `google_project_organization_policy` | legacy organization-policy resources superseded by the Organization Policy API (google_org_policy_policy) |
| `google_pubsub_lite_reservation` | deprecated in the provider schema |
| `google_pubsub_lite_subscription` | deprecated in the provider schema |
| `google_pubsub_lite_topic` | deprecated in the provider schema |
| `google_scc_event_threat_detection_custom_module` | Security Command Center v1 API surface, superseded by the SCC v2 and SCC Management APIs (equivalents exist for every resource) |
| `google_scc_folder_custom_module` | Security Command Center v1 API surface, superseded by the SCC v2 and SCC Management APIs (equivalents exist for every resource) |
| `google_scc_folder_notification_config` | Security Command Center v1 API surface, superseded by the SCC v2 and SCC Management APIs (equivalents exist for every resource) |
| `google_scc_folder_scc_big_query_export` | Security Command Center v1 API surface, superseded by the SCC v2 and SCC Management APIs (equivalents exist for every resource) |
| `google_scc_mute_config` | Security Command Center v1 API surface, superseded by the SCC v2 and SCC Management APIs (equivalents exist for every resource) |
| `google_scc_notification_config` | Security Command Center v1 API surface, superseded by the SCC v2 and SCC Management APIs (equivalents exist for every resource) |
| `google_scc_organization_custom_module` | Security Command Center v1 API surface, superseded by the SCC v2 and SCC Management APIs (equivalents exist for every resource) |
| `google_scc_organization_scc_big_query_export` | Security Command Center v1 API surface, superseded by the SCC v2 and SCC Management APIs (equivalents exist for every resource) |
| `google_scc_project_custom_module` | Security Command Center v1 API surface, superseded by the SCC v2 and SCC Management APIs (equivalents exist for every resource) |
| `google_scc_project_notification_config` | Security Command Center v1 API surface, superseded by the SCC v2 and SCC Management APIs (equivalents exist for every resource) |
| `google_scc_project_scc_big_query_export` | Security Command Center v1 API surface, superseded by the SCC v2 and SCC Management APIs (equivalents exist for every resource) |
| `google_scc_source` | Security Command Center v1 API surface, superseded by the SCC v2 and SCC Management APIs (equivalents exist for every resource) |
| `google_scc_v2_organization_scc_big_query_exports` | deprecated in the provider schema |
| `google_sourcerepo_repository` | Cloud Source Repositories is closed to new customers and sunsetting |
| `google_storage_bucket_access_control` | legacy per-object ACL model, superseded by uniform bucket-level IAM (the GcpGcsBucket kind's iam_members) |
| `google_storage_bucket_acl` | legacy per-object ACL model, superseded by uniform bucket-level IAM (the GcpGcsBucket kind's iam_members) |
| `google_storage_default_object_access_control` | legacy per-object ACL model, superseded by uniform bucket-level IAM (the GcpGcsBucket kind's iam_members) |
| `google_storage_default_object_acl` | legacy per-object ACL model, superseded by uniform bucket-level IAM (the GcpGcsBucket kind's iam_members) |
| `google_storage_object_access_control` | legacy per-object ACL model, superseded by uniform bucket-level IAM (the GcpGcsBucket kind's iam_members) |
| `google_storage_object_acl` | legacy per-object ACL model, superseded by uniform bucket-level IAM (the GcpGcsBucket kind's iam_members) |
| `google_vertex_ai_featurestore` | legacy Vertex AI Featurestore, superseded by feature groups and feature online stores |
| `google_vertex_ai_featurestore_entitytype` | legacy Vertex AI Featurestore, superseded by feature groups and feature online stores |
| `google_vertex_ai_featurestore_entitytype_feature` | legacy Vertex AI Featurestore, superseded by feature groups and feature online stores |
| `google_vertex_ai_schedule` | deprecated in the provider schema |
