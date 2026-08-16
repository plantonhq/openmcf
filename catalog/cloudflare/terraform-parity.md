---
title: "Terraform Parity"
description: "Measured parity of the CLOUDFLARE catalog against the pinned Terraform provider"
icon: "check-circle"
order: 90
---

<!-- GENERATED FILE -- DO NOT EDIT.
     Rendered from the committed provider schemas, the kind registry, the
     Terraform modules, the per-kind provider-parity manifests, the
     dispositions ledger, and the E2E profiles.
     parameters: provider=cloudflare ga-schema=cloudflare
     Regenerate: make generate-provider-parity-report -->

# CLOUDFLARE Terraform Parity

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
| Provider schema | `aws@6.58.0` |
| Provider schema (parity baseline) | `cloudflare@5.23.0` |
| Kinds in the catalog | 33 |
| Distinct provider resources consumed | 68 |
| Spec fields authored across all kinds | 1379 |
| Module pins on `aws` | `~> 5.0` × 1 |
| Module pins on `cloudflare` | `~> 5.23` × 33 |
| Module pins on `tls` | `~> 4.0` × 1 |

The GA provider is the parity baseline. Capability that exists only in a
secondary channel (for Google, the `google-beta` provider) enters per kind
through an explicitly enumerated admission list, never wholesale.

## Depth: per-kind accounting

Every configurable, non-deprecated provider argument of a kind's consumed
resources must be matched to a spec field, mapped by recorded judgment, or
excluded with a recorded reason -- and every spec field must reach provider
surface. **Accounted** means both directions hold with zero unexplained
gaps. **Proven** means live end-to-end runs passed on both IaC engines.

**33 of 33 kinds are at total accounting; 0 proven live.**

| Kind | Provider args | Matched | Mapped | Excluded | Open gaps | Accounted | Proven |
|---|---|---|---|---|---|---|---|
| CloudflareCacheSettings | 12 | 6 | 6 | 0 | 0 | ✅ | — |
| CloudflareCertificatePack | 7 | 7 | 0 | 0 | 0 | ✅ | — |
| CloudflareCustomHostname | 6 | 5 | 1 | 0 | 0 | ✅ | — |
| CloudflareCustomHostnameFallbackOrigin | 2 | 2 | 0 | 0 | 0 | ✅ | — |
| CloudflareD1Database | 5 | 2 | 3 | 0 | 0 | ✅ | — |
| CloudflareDnsRecord | 12 | 10 | 2 | 0 | 0 | ✅ | — |
| CloudflareDnsZone | 38 | 21 | 12 | 5 | 0 | ✅ | — |
| CloudflareEmailRoutingAddress | 3 | 3 | 0 | 0 | 0 | ✅ | — |
| CloudflareEmailRoutingRule | 8 | 4 | 2 | 2 | 0 | ✅ | — |
| CloudflareEmailRoutingZone | 10 | 3 | 4 | 3 | 0 | ✅ | — |
| CloudflareHyperdriveConfig | 6 | 3 | 3 | 0 | 0 | ✅ | — |
| CloudflareKvNamespace | 2 | 1 | 1 | 0 | 0 | ✅ | — |
| CloudflareList | 5 | 4 | 0 | 1 | 0 | ✅ | — |
| CloudflareListItem | 7 | 5 | 2 | 0 | 0 | ✅ | — |
| CloudflareLoadBalancer | 20 | 11 | 9 | 0 | 0 | ✅ | — |
| CloudflareLoadBalancerMonitor | 17 | 16 | 1 | 0 | 0 | ✅ | — |
| CloudflareLoadBalancerPool | 15 | 9 | 4 | 2 | 0 | ✅ | — |
| CloudflareOriginCaCertificate | 4 | 4 | 0 | 0 | 0 | ✅ | — |
| CloudflarePagesProject | 9 | 4 | 5 | 0 | 0 | ✅ | — |
| CloudflareQueue | 9 | 5 | 2 | 2 | 0 | ✅ | — |
| CloudflareR2Bucket | 34 | 22 | 12 | 0 | 0 | ✅ | — |
| CloudflareRuleset | 7 | 5 | 2 | 0 | 0 | ✅ | — |
| CloudflareTurnstileWidget | 9 | 9 | 0 | 0 | 0 | ✅ | — |
| CloudflareWorker | 40 | 10 | 20 | 10 | 0 | ✅ | — |
| CloudflareWorkersKvPair | 5 | 5 | 0 | 0 | 0 | ✅ | — |
| CloudflareZeroTrustAccessApplication | 39 | 29 | 10 | 0 | 0 | ✅ | — |
| CloudflareZeroTrustAccessGroup | 7 | 4 | 3 | 0 | 0 | ✅ | — |
| CloudflareZeroTrustAccessPolicy | 14 | 8 | 6 | 0 | 0 | ✅ | — |
| CloudflareZeroTrustTunnel | 8 | 5 | 2 | 1 | 0 | ✅ | — |
| CloudflareZeroTrustTunnelRoute | 5 | 5 | 0 | 0 | 0 | ✅ | — |
| CloudflareZeroTrustTunnelVirtualNetwork | 4 | 4 | 0 | 0 | 0 | ✅ | — |
| CloudflareZoneSettings | 16 | 5 | 9 | 2 | 0 | ✅ | — |
| CloudflareZoneTlsSettings | 16 | 6 | 9 | 1 | 0 | ✅ | — |

## Breadth: every GA resource, one disposition

All resources of `cloudflare@5.23.0` land in exactly one class:

| Disposition | Resources | Meaning |
|---|---|---|
| Modeled | 66 | consumed by a kind's Terraform module today |
| IAM-covered | 0 | per-resource IAM member/binding/policy triplets, covered by the owning kinds' additive `iam_members` fields |
| Composed | 0 | capability covered through an existing kind's surface rather than a kind of its own |
| Planned | 139 | judged to be covered by a planned kind or planned composition, not built yet |
| Deferred | 45 | deliberately not offered, each with the recorded reason |
| Excluded as deprecated | 7 | deprecated or superseded provider surface |
| **Total** | **257** | |

## The enumerated record

The full per-resource record, so the accounting above is verifiable
rather than trusted.

### Modeled (66)

| Resource | Consuming kinds |
|---|---|
| `cloudflare_argo_smart_routing` | consumed by CloudflareCacheSettings |
| `cloudflare_argo_tiered_caching` | consumed by CloudflareCacheSettings |
| `cloudflare_certificate_authorities_hostname_associations` | consumed by CloudflareZoneTlsSettings |
| `cloudflare_certificate_pack` | consumed by CloudflareCertificatePack |
| `cloudflare_custom_hostname` | consumed by CloudflareCustomHostname |
| `cloudflare_custom_hostname_fallback_origin` | consumed by CloudflareCustomHostnameFallbackOrigin |
| `cloudflare_d1_database` | consumed by CloudflareD1Database |
| `cloudflare_dns_record` | consumed by CloudflareDnsRecord, CloudflareDnsZone |
| `cloudflare_email_routing_address` | consumed by CloudflareEmailRoutingAddress |
| `cloudflare_email_routing_catch_all` | consumed by CloudflareEmailRoutingZone |
| `cloudflare_email_routing_dns` | consumed by CloudflareEmailRoutingZone |
| `cloudflare_email_routing_rule` | consumed by CloudflareEmailRoutingRule |
| `cloudflare_email_routing_settings` | consumed by CloudflareEmailRoutingZone |
| `cloudflare_hostname_tls_setting` | consumed by CloudflareZoneTlsSettings |
| `cloudflare_hyperdrive_config` | consumed by CloudflareHyperdriveConfig |
| `cloudflare_list` | consumed by CloudflareList |
| `cloudflare_list_item` | consumed by CloudflareListItem |
| `cloudflare_load_balancer` | consumed by CloudflareLoadBalancer |
| `cloudflare_load_balancer_monitor` | consumed by CloudflareLoadBalancerMonitor |
| `cloudflare_load_balancer_pool` | consumed by CloudflareLoadBalancerPool |
| `cloudflare_managed_transforms` | consumed by CloudflareZoneSettings |
| `cloudflare_origin_ca_certificate` | consumed by CloudflareOriginCaCertificate |
| `cloudflare_origin_cloud_region` | consumed by CloudflareZoneSettings |
| `cloudflare_origin_tls_compliance_modes` | consumed by CloudflareZoneTlsSettings |
| `cloudflare_pages_domain` | consumed by CloudflarePagesProject |
| `cloudflare_pages_project` | consumed by CloudflarePagesProject |
| `cloudflare_queue` | consumed by CloudflareQueue |
| `cloudflare_queue_consumer` | consumed by CloudflareQueue |
| `cloudflare_r2_bucket` | consumed by CloudflareR2Bucket |
| `cloudflare_r2_bucket_cors` | consumed by CloudflareR2Bucket |
| `cloudflare_r2_bucket_event_notification` | consumed by CloudflareR2Bucket |
| `cloudflare_r2_bucket_lifecycle` | consumed by CloudflareR2Bucket |
| `cloudflare_r2_bucket_lock` | consumed by CloudflareR2Bucket |
| `cloudflare_r2_custom_domain` | consumed by CloudflareR2Bucket |
| `cloudflare_r2_managed_domain` | consumed by CloudflareR2Bucket |
| `cloudflare_regional_tiered_cache` | consumed by CloudflareCacheSettings |
| `cloudflare_ruleset` | consumed by CloudflareRuleset |
| `cloudflare_tiered_cache` | consumed by CloudflareCacheSettings |
| `cloudflare_total_tls` | consumed by CloudflareZoneTlsSettings |
| `cloudflare_turnstile_widget` | consumed by CloudflareTurnstileWidget |
| `cloudflare_universal_ssl_setting` | consumed by CloudflareZoneTlsSettings |
| `cloudflare_url_normalization_settings` | consumed by CloudflareZoneSettings |
| `cloudflare_waiting_room_settings` | consumed by CloudflareZoneSettings |
| `cloudflare_workers_cron_trigger` | consumed by CloudflareWorker |
| `cloudflare_workers_custom_domain` | consumed by CloudflareWorker |
| `cloudflare_workers_kv` | consumed by CloudflareWorkersKvPair |
| `cloudflare_workers_kv_namespace` | consumed by CloudflareKvNamespace |
| `cloudflare_workers_route` | consumed by CloudflareWorker |
| `cloudflare_workers_script` | consumed by CloudflareWorker |
| `cloudflare_workers_script_subdomain` | consumed by CloudflareWorker |
| `cloudflare_zero_trust_access_application` | consumed by CloudflareZeroTrustAccessApplication |
| `cloudflare_zero_trust_access_group` | consumed by CloudflareZeroTrustAccessGroup |
| `cloudflare_zero_trust_access_policy` | consumed by CloudflareZeroTrustAccessPolicy |
| `cloudflare_zero_trust_tunnel_cloudflared` | consumed by CloudflareZeroTrustTunnel |
| `cloudflare_zero_trust_tunnel_cloudflared_config` | consumed by CloudflareZeroTrustTunnel |
| `cloudflare_zero_trust_tunnel_cloudflared_route` | consumed by CloudflareZeroTrustTunnelRoute |
| `cloudflare_zero_trust_tunnel_cloudflared_virtual_network` | consumed by CloudflareZeroTrustTunnelVirtualNetwork |
| `cloudflare_zone` | consumed by CloudflareDnsZone |
| `cloudflare_zone_auto_origin_tls_kex` | consumed by CloudflareZoneTlsSettings |
| `cloudflare_zone_cache_reserve` | consumed by CloudflareCacheSettings |
| `cloudflare_zone_cache_variants` | consumed by CloudflareCacheSettings |
| `cloudflare_zone_dns_settings` | consumed by CloudflareDnsZone |
| `cloudflare_zone_dnssec` | consumed by CloudflareDnsZone |
| `cloudflare_zone_hold` | consumed by CloudflareDnsZone |
| `cloudflare_zone_setting` | consumed by CloudflareZoneSettings |
| `cloudflare_zone_subscription` | consumed by CloudflareDnsZone |

### Planned (139)

| Resource | Recorded reason |
|---|---|
| `cloudflare_access_rule` | judged as a planned CloudflareIpAccessRule kind (IP/ASN/country allow-block rules; named to avoid collision with the Zero Trust Access namespace) |
| `cloudflare_account` | judged as a planned CloudflareAccount kind (foundational for multi-account and MSP automation) |
| `cloudflare_account_dns_settings` | judged as a planned CloudflareAccountDnsSettings kind (account-level DNS settings singleton with apply/revert semantics) |
| `cloudflare_account_member` | judged as a planned CloudflareAccountMember kind (memberships with roles) |
| `cloudflare_account_subscription` | folds into the planned CloudflareAccount kind (the plan is an attribute of the account) |
| `cloudflare_account_token` | judged as a planned CloudflareAccountApiToken kind (scoped account API tokens with independent rotation lifecycle) |
| `cloudflare_ai_gateway` | judged as a planned CloudflareAiGateway kind (caching, rate limiting, and logging for AI calls) |
| `cloudflare_ai_gateway_dynamic_routing` | folds into the planned CloudflareAiGateway kind (routing configuration of the gateway) |
| `cloudflare_ai_search_instance` | judged as a planned CloudflareAiSearchInstance kind (managed retrieval-augmented search instance) |
| `cloudflare_ai_search_namespace` | judged as a planned CloudflareAiSearchNamespace kind (grouping namespace referenced by search instances) |
| `cloudflare_ai_search_token` | folds into the planned CloudflareAiSearchInstance kind (endpoint access token of the instance) |
| `cloudflare_api_shield` | judged as a planned CloudflareApiShield kind (the zone's API Shield configuration umbrella: auth characteristics, validation defaults) |
| `cloudflare_api_shield_operation` | folds into the planned CloudflareApiShield kind (endpoint registrations are per-zone configuration rows meaningless outside the shield) |
| `cloudflare_api_shield_schema_validation_settings` | folds into the planned CloudflareApiShield kind (zone validation defaults) |
| `cloudflare_api_token` | judged as a planned CloudflareApiToken kind (user-scoped tokens; the API and ownership model differ from account tokens) |
| `cloudflare_authenticated_origin_pulls` | judged as a planned CloudflareAuthenticatedOriginPulls kind (zone and hostname authenticated-origin-pulls enablement) |
| `cloudflare_authenticated_origin_pulls_certificate` | judged as a planned CloudflareAuthenticatedOriginPullsCertificate kind (uploaded client certificate with independent rotation lifecycle) |
| `cloudflare_authenticated_origin_pulls_hostname_certificate` | folds into the planned CloudflareAuthenticatedOriginPullsCertificate kind (hostname-scoped variant of the same certificate upload) |
| `cloudflare_authenticated_origin_pulls_settings` | folds into the planned CloudflareAuthenticatedOriginPulls kind (settings of the enablement surface) |
| `cloudflare_bot_management` | judged as a planned CloudflareBotManagement kind (zone-singleton bot management configuration) |
| `cloudflare_calls_sfu_app` | judged as a planned CloudflareCallsSfuApp kind (realtime SFU application credentials provisioned per application) |
| `cloudflare_calls_turn_app` | judged as a planned CloudflareCallsTurnApp kind (TURN service credentials provisioned per application) |
| `cloudflare_client_certificate` | judged as a planned CloudflareClientCertificate kind (Cloudflare-issued client certificates for API Shield mTLS) |
| `cloudflare_cloud_connector_rules` | judged as a planned CloudflareCloudConnectorRules kind (zone-singleton ordered routing table to cloud storage origins) |
| `cloudflare_connectivity_directory_service` | judged as a planned CloudflareConnectivityDirectoryService kind (service directory entries referencing tunnels) |
| `cloudflare_content_scanning` | judged as a planned CloudflareContentScanning kind (malicious-upload scanning enablement plus expressions) |
| `cloudflare_content_scanning_expression` | folds into the planned CloudflareContentScanning kind (custom expression rows of the parent) |
| `cloudflare_custom_csr` | judged as a planned CloudflareCustomCsr kind (CSR generation referenced during certificate issuance) |
| `cloudflare_custom_origin_trust_store` | judged as a planned CloudflareCustomOriginTrustStore kind (uploaded CA bundle for origin verification) |
| `cloudflare_custom_page_asset` | folds into the planned CloudflareCustomPage kind (uploaded assets exist only to serve custom pages) |
| `cloudflare_custom_pages` | judged as a planned CloudflareCustomPage kind (one instance per error-page identifier with independent lifecycle) |
| `cloudflare_custom_ssl` | judged as a planned CloudflareCustomSslCertificate kind (bring-your-own certificate upload with rotation lifecycle) |
| `cloudflare_dns_firewall` | judged as a planned CloudflareDnsFirewall kind (standalone DNS firewall clusters with their own nameserver IPs and lifecycle) |
| `cloudflare_flagship_app` | judged as a planned CloudflareFlagshipApp kind (feature-flags application container) |
| `cloudflare_flagship_flag` | judged as a planned CloudflareFlagshipFlag kind (flags with independent lifecycle referenced by SDKs) |
| `cloudflare_google_tag_gateway` | judged as a planned CloudflareGoogleTagGateway kind (zone-singleton tag-gateway configuration) |
| `cloudflare_healthcheck` | judged as a planned CloudflareHealthcheck kind (standalone zone health checks, no load balancer required) |
| `cloudflare_image_variant` | judged as a planned CloudflareImagesVariant kind (named transformation preset referenced by delivery URLs) |
| `cloudflare_leaked_credential_check` | judged as a planned CloudflareLeakedCredentialCheck kind (enablement plus custom detection rules in one kind) |
| `cloudflare_leaked_credential_check_rule` | folds into the planned CloudflareLeakedCredentialCheck kind (detection rows are meaningless without the check) |
| `cloudflare_load_balancer_monitor_group` | judged as a planned CloudflareLoadBalancerMonitorGroup kind (monitor groups referenced by pools) |
| `cloudflare_logpush_job` | judged as a planned CloudflareLogpushJob kind (log delivery jobs, the observability backbone) |
| `cloudflare_logpush_ownership_challenge` | folds into the planned CloudflareLogpushJob kind (the challenge is a destination-validation step, not an object) |
| `cloudflare_mtls_certificate` | judged as a planned CloudflareMtlsCertificate kind (CA certificates referenced by Workers mTLS bindings and authenticated origin pulls) |
| `cloudflare_notification_policy` | judged as a planned CloudflareNotificationPolicy kind (alerting policies for every Cloudflare product) |
| `cloudflare_notification_policy_webhooks` | judged as a planned CloudflareNotificationWebhook kind (webhook destinations referenced by policies) |
| `cloudflare_oauth_client` | judged as a planned CloudflareOauthClient kind (OAuth clients against Cloudflare as identity provider) |
| `cloudflare_observatory_scheduled_test` | judged as a planned CloudflareObservatoryScheduledTest kind (scheduled speed tests per page) |
| `cloudflare_organization` | judged as a planned CloudflareOrganization kind (organization hierarchy container) |
| `cloudflare_organization_profile` | folds into the planned CloudflareOrganization kind (the profile is settings of the organization) |
| `cloudflare_pipeline` | judged as a planned CloudflarePipeline kind (ingestion pipeline referencing sinks and streams) |
| `cloudflare_pipeline_sink` | judged as a planned CloudflarePipelineSink kind (pipeline destination, e.g. an R2 bucket) |
| `cloudflare_pipeline_stream` | judged as a planned CloudflarePipelineStream kind (pipeline stream input) |
| `cloudflare_r2_bucket_sippy` | folds into the existing CloudflareR2Bucket kind's planned depth expansion (incremental-migration setting of a bucket) |
| `cloudflare_r2_data_catalog` | folds into the existing CloudflareR2Bucket kind's planned depth expansion (per-bucket data catalog toggle) |
| `cloudflare_schema_validation_operation_settings` | folds into the planned CloudflareApiShield kind (per-operation validation overrides) |
| `cloudflare_schema_validation_schemas` | judged as a planned CloudflareSchemaValidationSchema kind (uploaded OpenAPI schema with its own upload and activation lifecycle) |
| `cloudflare_schema_validation_settings` | folds into the planned CloudflareApiShield kind (v2 zone validation defaults) |
| `cloudflare_secrets_store` | judged as a planned CloudflareSecretsStore kind (account secrets store, parent of secrets) |
| `cloudflare_secrets_store_secret` | judged as a planned CloudflareSecretsStoreSecret kind (secrets referenced by Worker bindings, rotating independently) |
| `cloudflare_share` | judged as a planned CloudflareShare kind (cross-account resource sharing) |
| `cloudflare_share_recipient` | folds into the planned CloudflareShare kind (recipient rows of the share) |
| `cloudflare_share_resource` | folds into the planned CloudflareShare kind (resource rows of the share) |
| `cloudflare_snippet` | judged as a planned CloudflareSnippet kind (edge code unit with its own upload lifecycle) |
| `cloudflare_snippet_rules` | judged as a planned CloudflareSnippetRules kind (the zone's snippet routing table, referencing many snippets and outliving any one of them) |
| `cloudflare_spectrum_application` | judged as a planned CloudflareSpectrumApplication kind (TCP/UDP proxying applications) |
| `cloudflare_stream_key` | judged as a planned CloudflareStreamSigningKey kind (signing key referenced by application code for signed playback URLs) |
| `cloudflare_stream_live_input` | judged as a planned CloudflareStreamLiveInput kind (live ingest endpoint with keys and independent lifecycle) |
| `cloudflare_stream_watermark` | judged as a planned CloudflareStreamWatermarkProfile kind (created once, referenced at every upload) |
| `cloudflare_stream_webhook` | judged as a planned CloudflareStreamWebhook kind (account-level notification webhook singleton) |
| `cloudflare_token_validation_config` | judged as a planned CloudflareTokenValidationConfig kind (JWT credential sets referenced by validation rules) |
| `cloudflare_token_validation_rules` | folds into the planned CloudflareApiShield kind (zone-level JWT rules spanning multiple credential configs) |
| `cloudflare_user_agent_blocking_rule` | judged as a planned CloudflareUserAgentBlockingRule kind (standalone user-agent blocking rule object) |
| `cloudflare_user_group` | judged as a planned CloudflareUserGroup kind (permission groups referenced by policies) |
| `cloudflare_user_group_members` | folds into the planned CloudflareUserGroup kind (membership rows of the group) |
| `cloudflare_waiting_room` | judged as a planned CloudflareWaitingRoom kind (the room itself) |
| `cloudflare_waiting_room_event` | judged as a planned CloudflareWaitingRoomEvent kind (scheduled events created and deleted on their own cadence per launch) |
| `cloudflare_waiting_room_rules` | folds into the planned CloudflareWaitingRoom kind (the room's singleton bypass-rule list) |
| `cloudflare_web_analytics_rule` | folds into the planned CloudflareWebAnalyticsSite kind (path rules of the site) |
| `cloudflare_web_analytics_site` | judged as a planned CloudflareWebAnalyticsSite kind (real-user-monitoring site) |
| `cloudflare_worker` | folds into the existing CloudflareWorker kind's planned depth expansion (the modern script-settings surface of the worker the kind already owns) |
| `cloudflare_worker_version` | folds into the existing CloudflareWorker kind's planned depth expansion (versions are deployment artifacts of the worker, not independent objects) |
| `cloudflare_workers_deployment` | folds into the existing CloudflareWorker kind's planned depth expansion (gradual-deployment configuration of the worker) |
| `cloudflare_workers_for_platforms_dispatch_namespace` | judged as a planned CloudflareWorkersDispatchNamespace kind (Workers for Platforms dispatch namespace referenced by dispatch workers) |
| `cloudflare_workflow` | judged as a planned CloudflareWorkflow kind (durable-execution workflows bound to a worker script) |
| `cloudflare_zero_trust_access_ai_controls_mcp_portal` | judged as a planned CloudflareZeroTrustMcpPortal kind (MCP server portals behind Access) |
| `cloudflare_zero_trust_access_ai_controls_mcp_server` | judged as a planned CloudflareZeroTrustMcpServer kind (registered MCP servers behind Access, pairing with the portal) |
| `cloudflare_zero_trust_access_custom_page` | judged as a planned CloudflareZeroTrustAccessCustomPage kind (branded block and login pages referenced by applications) |
| `cloudflare_zero_trust_access_identity_provider` | judged as a planned CloudflareZeroTrustAccessIdentityProvider kind (identity providers referenced by applications and policies) |
| `cloudflare_zero_trust_access_infrastructure_target` | judged as a planned CloudflareZeroTrustAccessInfrastructureTarget kind (SSH infrastructure targets referenced by infrastructure applications) |
| `cloudflare_zero_trust_access_key_configuration` | folds into the planned CloudflareZeroTrustOrganization kind (key rotation settings are organization-singleton configuration) |
| `cloudflare_zero_trust_access_mtls_certificate` | judged as a planned CloudflareZeroTrustAccessMtlsCertificate kind (client CA certificate with associated hostnames and independent lifecycle) |
| `cloudflare_zero_trust_access_mtls_hostname_settings` | folds into the planned CloudflareZeroTrustAccessMtlsCertificate kind (hostname settings exist in service of the mTLS certificates) |
| `cloudflare_zero_trust_access_service_token` | judged as a planned CloudflareZeroTrustAccessServiceToken kind (machine credentials with independent rotation lifecycle) |
| `cloudflare_zero_trust_access_short_lived_certificate` | folds into the existing CloudflareZeroTrustAccessApplication kind's planned depth expansion (SSH short-lived certificate issuance is a per-application toggle) |
| `cloudflare_zero_trust_access_tag` | judged as a planned CloudflareZeroTrustAccessTag kind (label referenced by name across applications) |
| `cloudflare_zero_trust_device_custom_profile` | judged as a planned CloudflareZeroTrustDeviceCustomProfile kind (targeted WARP profiles, many per account) |
| `cloudflare_zero_trust_device_custom_profile_local_domain_fallback` | folds into the planned CloudflareZeroTrustDeviceCustomProfile kind (fallback-domain rows of the profile) |
| `cloudflare_zero_trust_device_default_profile` | judged as a planned CloudflareZeroTrustDeviceDefaultProfile kind (account-singleton default WARP profile) |
| `cloudflare_zero_trust_device_default_profile_certificates` | folds into the planned CloudflareZeroTrustDeviceDefaultProfile kind (certificate enablement of the profile) |
| `cloudflare_zero_trust_device_default_profile_local_domain_fallback` | folds into the planned CloudflareZeroTrustDeviceDefaultProfile kind (fallback-domain rows of the profile) |
| `cloudflare_zero_trust_device_deployment_groups` | judged as a planned CloudflareZeroTrustDeviceDeploymentGroup kind (WARP client rollout rings with independent membership lifecycle) |
| `cloudflare_zero_trust_device_ip_profile` | judged as a planned CloudflareZeroTrustDeviceIpProfile kind (IP-scoped WARP profile variant) |
| `cloudflare_zero_trust_device_managed_networks` | judged as a planned CloudflareZeroTrustDeviceManagedNetwork kind (network fingerprints referenced by profile targeting) |
| `cloudflare_zero_trust_device_posture_integration` | judged as a planned CloudflareZeroTrustDevicePostureIntegration kind (third-party posture providers with credentials and lifecycle) |
| `cloudflare_zero_trust_device_posture_rule` | judged as a planned CloudflareZeroTrustDevicePostureRule kind (posture checks referenced by Access and Gateway policies) |
| `cloudflare_zero_trust_device_settings` | judged as a planned CloudflareZeroTrustDeviceSettings kind (account-wide device enrollment settings singleton) |
| `cloudflare_zero_trust_device_subnet` | judged as a planned CloudflareZeroTrustDeviceSubnet kind (virtual device subnets referencing virtual networks) |
| `cloudflare_zero_trust_dex_rule` | folds into the planned CloudflareZeroTrustDexTest kind (targeting rules exist in service of the tests) |
| `cloudflare_zero_trust_dex_test` | judged as a planned CloudflareZeroTrustDexTest kind (synthetic experience tests) |
| `cloudflare_zero_trust_dlp_custom_entry` | folds into the planned CloudflareZeroTrustDlpCustomProfile kind (custom detection rows of the profile) |
| `cloudflare_zero_trust_dlp_custom_profile` | judged as a planned CloudflareZeroTrustDlpCustomProfile kind (custom detection profile) |
| `cloudflare_zero_trust_dlp_data_class` | folds into the planned CloudflareZeroTrustDlpSettings kind (classification taxonomy row) |
| `cloudflare_zero_trust_dlp_data_tag` | folds into the planned CloudflareZeroTrustDlpSettings kind (tag taxonomy row) |
| `cloudflare_zero_trust_dlp_data_tag_category` | folds into the planned CloudflareZeroTrustDlpSettings kind (tag taxonomy row) |
| `cloudflare_zero_trust_dlp_dataset` | judged as a planned CloudflareZeroTrustDlpDataset kind (exact-data-match and wordlist datasets with upload lifecycle) |
| `cloudflare_zero_trust_dlp_entry` | folds into the planned CloudflareZeroTrustDlpCustomProfile kind (generic entry rows inside a profile) |
| `cloudflare_zero_trust_dlp_integration_entry` | folds into the planned CloudflareZeroTrustDlpSettings kind (integration-sourced entries managed at account level) |
| `cloudflare_zero_trust_dlp_predefined_entry` | folds into the planned CloudflareZeroTrustDlpPredefinedProfile kind (enablement rows of predefined entries) |
| `cloudflare_zero_trust_dlp_predefined_profile` | judged as a planned CloudflareZeroTrustDlpPredefinedProfile kind (enablement and configuration of predefined profiles) |
| `cloudflare_zero_trust_dlp_sensitivity_group` | folds into the planned CloudflareZeroTrustDlpSettings kind (sensitivity taxonomy row) |
| `cloudflare_zero_trust_dlp_sensitivity_level` | folds into the planned CloudflareZeroTrustDlpSettings kind (sensitivity taxonomy row) |
| `cloudflare_zero_trust_dlp_sensitivity_level_order` | folds into the planned CloudflareZeroTrustDlpSettings kind (ordering of the sensitivity taxonomy) |
| `cloudflare_zero_trust_dlp_settings` | judged as a planned CloudflareZeroTrustDlpSettings kind (account-level DLP configuration including the sensitivity and tag taxonomy) |
| `cloudflare_zero_trust_dns_location` | judged as a planned CloudflareZeroTrustDnsLocation kind (DNS filtering locations with their own endpoints and lifecycle) |
| `cloudflare_zero_trust_gateway_certificate` | judged as a planned CloudflareZeroTrustGatewayCertificate kind (TLS inspection certificates with activation lifecycle) |
| `cloudflare_zero_trust_gateway_logging` | folds into the planned CloudflareZeroTrustGatewaySettings kind (logging redaction toggles are Gateway settings) |
| `cloudflare_zero_trust_gateway_pacfile` | folds into the planned CloudflareZeroTrustGatewaySettings kind (PAC file content is Gateway proxy configuration) |
| `cloudflare_zero_trust_gateway_policy` | judged as a planned CloudflareZeroTrustGatewayPolicy kind (DNS/HTTP/network filtering rules, the core Gateway object) |
| `cloudflare_zero_trust_gateway_proxy_endpoint` | judged as a planned CloudflareZeroTrustGatewayProxyEndpoint kind (proxy endpoints with allowed-IP configuration) |
| `cloudflare_zero_trust_gateway_settings` | judged as a planned CloudflareZeroTrustGatewaySettings kind (account-level Gateway configuration singleton) |
| `cloudflare_zero_trust_list` | judged as a planned CloudflareZeroTrustList kind (domain/IP/serial lists referenced by nearly every Gateway policy) |
| `cloudflare_zero_trust_network_hostname_route` | judged as a planned CloudflareZeroTrustNetworkHostnameRoute kind (hostname-based private network routes referencing tunnels and virtual networks) |
| `cloudflare_zero_trust_organization` | judged as a planned CloudflareZeroTrustOrganization kind (team-domain singleton: login design, session defaults) |
| `cloudflare_zero_trust_risk_behavior` | judged as a planned CloudflareZeroTrustRiskScoring kind (account risk-behavior configuration) |
| `cloudflare_zero_trust_risk_scoring_integration` | folds into the planned CloudflareZeroTrustRiskScoring kind (SIEM/IdP integration exists in service of risk scoring) |
| `cloudflare_zero_trust_tunnel_warp_connector` | judged as a planned CloudflareZeroTrustWarpConnector kind (site-to-site WARP connector tunnels) |
| `cloudflare_zero_trust_tunnel_warp_connector_config` | folds into the planned CloudflareZeroTrustWarpConnector kind (configuration of the connector, mirroring the cloudflared pattern) |
| `cloudflare_zone_lockdown` | judged as a planned CloudflareZoneLockdown kind (URL lockdown rules, many per zone with independent lifecycle) |

### Deferred (45)

| Resource | Recorded reason |
|---|---|
| `cloudflare_account_dns_settings_internal_view` | Enterprise-only entitlement: internal DNS views |
| `cloudflare_address_map` | requires the BYOIP entitlement |
| `cloudflare_byo_ip_prefix` | requires the BYOIP entitlement |
| `cloudflare_cloudforce_one_request` | threat-intelligence request workflow, not provisioned infrastructure |
| `cloudflare_cloudforce_one_request_asset` | asset attachment of the threat-intelligence request workflow, not provisioned infrastructure |
| `cloudflare_cloudforce_one_request_message` | message on a threat-intelligence request ticket, not provisioned infrastructure |
| `cloudflare_cloudforce_one_request_priority` | priority row of the threat-intelligence request workflow, not provisioned infrastructure |
| `cloudflare_dls_prefix_binding` | requires the BYOIP and Data Localization Suite entitlements |
| `cloudflare_dns_zone_transfers_acl` | Enterprise-only entitlement: secondary DNS zone transfers |
| `cloudflare_dns_zone_transfers_incoming` | Enterprise-only entitlement: secondary DNS zone transfers |
| `cloudflare_dns_zone_transfers_outgoing` | Enterprise-only entitlement: secondary DNS zone transfers |
| `cloudflare_dns_zone_transfers_peer` | Enterprise-only entitlement: secondary DNS zone transfers |
| `cloudflare_dns_zone_transfers_tsig` | Enterprise-only entitlement: secondary DNS zone transfers |
| `cloudflare_email_security_block_sender` | Enterprise add-on entitlement: cloud email security |
| `cloudflare_email_security_impersonation_registry` | Enterprise add-on entitlement: cloud email security |
| `cloudflare_email_security_trusted_domains` | Enterprise add-on entitlement: cloud email security |
| `cloudflare_image` | per-image content object; application data rather than provisioned infrastructure |
| `cloudflare_keyless_certificate` | Enterprise-only entitlement: Keyless SSL |
| `cloudflare_logpull_retention` | Enterprise-only entitlement: Logpull API |
| `cloudflare_magic_network_monitoring_configuration` | Enterprise-only entitlement: Magic Network Monitoring |
| `cloudflare_magic_network_monitoring_rule` | Enterprise-only entitlement: Magic Network Monitoring |
| `cloudflare_magic_transit_cf1_site` | Enterprise-only entitlement: Magic Transit |
| `cloudflare_magic_transit_connector` | Enterprise-only entitlement: Magic Transit |
| `cloudflare_magic_transit_site` | Enterprise-only entitlement: Magic Transit |
| `cloudflare_magic_transit_site_acl` | Enterprise-only entitlement: Magic Transit |
| `cloudflare_magic_transit_site_lan` | Enterprise-only entitlement: Magic Transit |
| `cloudflare_magic_transit_site_wan` | Enterprise-only entitlement: Magic Transit |
| `cloudflare_magic_wan_gre_tunnel` | Enterprise-only entitlement: Magic WAN |
| `cloudflare_magic_wan_ipsec_tunnel` | Enterprise-only entitlement: Magic WAN |
| `cloudflare_magic_wan_static_route` | Enterprise-only entitlement: Magic WAN |
| `cloudflare_moq_relay` | experimental Media-over-QUIC relay; revisit when the product stabilizes |
| `cloudflare_page_rule` | legacy surface superseded by the Rulesets engine, which the catalog already models; Cloudflare is migrating users off Page Rules |
| `cloudflare_page_shield_policy` | Enterprise add-on entitlement: Page Shield CSP policies |
| `cloudflare_regional_hostname` | Enterprise add-on entitlement: Regional Services |
| `cloudflare_registrar_domain` | interactive, account-singular domain transfers; not fleet-deployable infrastructure |
| `cloudflare_sso_connector` | Enterprise-only entitlement: dashboard single sign-on |
| `cloudflare_stream` | per-video content object; application data rather than provisioned infrastructure |
| `cloudflare_stream_audio_track` | content-follower of a video; application data rather than provisioned infrastructure |
| `cloudflare_stream_caption_language` | content-follower of a video; application data rather than provisioned infrastructure |
| `cloudflare_stream_download` | content-follower of a video; application data rather than provisioned infrastructure |
| `cloudflare_user` | manages the calling user's own profile; not a deployable infrastructure object |
| `cloudflare_vulnerability_scanner_credential` | nascent scanning surface; revisit on demand |
| `cloudflare_vulnerability_scanner_credential_set` | nascent scanning surface; revisit on demand |
| `cloudflare_vulnerability_scanner_target_environment` | nascent scanning surface; revisit on demand |
| `cloudflare_web3_hostname` | Cloudflare has de-emphasized its Web3 gateways; revisit on demand |

### Excluded as deprecated (7)

| Resource | Recorded reason |
|---|---|
| `cloudflare_api_shield_discovery_operation` | deprecated in the provider schema |
| `cloudflare_api_shield_operation_schema_validation_settings` | deprecated in the provider schema |
| `cloudflare_api_shield_schema` | deprecated in the provider schema |
| `cloudflare_filter` | deprecated in the provider schema |
| `cloudflare_firewall_rule` | deprecated in the provider schema |
| `cloudflare_rate_limit` | deprecated in the provider schema |
| `cloudflare_snippets` | deprecated in the provider schema |
