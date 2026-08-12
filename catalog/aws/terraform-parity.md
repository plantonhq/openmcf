---
title: "Terraform Parity"
description: "Measured parity of the AWS catalog against the pinned Terraform provider"
icon: "check-circle"
order: 90
---

<!-- GENERATED FILE -- DO NOT EDIT.
     Rendered from the committed provider schemas, the kind registry, the
     Terraform modules, the per-kind provider-parity manifests, the
     dispositions ledger, and the E2E profiles.
     parameters: provider=aws ga-schema=aws
     Regenerate: make generate-provider-parity-report -->

# AWS Terraform Parity

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
| Provider schema (parity baseline) | `aws@6.58.0` |
| Provider schema | `google@7.43.0` |
| Provider schema | `google-beta@7.43.0` |
| Kinds in the catalog | 112 |
| Distinct provider resources consumed | 249 |
| Spec fields authored across all kinds | 4474 |
| Module pins on `aws` | `~> 6.58` × 112 |

The GA provider is the parity baseline. Capability that exists only in a
secondary channel (for Google, the `google-beta` provider) enters per kind
through an explicitly enumerated admission list, never wholesale.

## Depth: per-kind accounting

Every configurable, non-deprecated provider argument of a kind's consumed
resources must be matched to a spec field, mapped by recorded judgment, or
excluded with a recorded reason -- and every spec field must reach provider
surface. **Accounted** means both directions hold with zero unexplained
gaps. **Proven** means live end-to-end runs passed on both IaC engines.

**78 of 112 kinds are at total accounting; 82 proven live.**

| Kind | Provider args | Matched | Mapped | Excluded | Open gaps | Accounted | Proven |
|---|---|---|---|---|---|---|---|
| AwsAlb | 71 | 17 | 12 | 42 | 0 | ✅ | ✅ pulumi, terraform |
| AwsAppRunnerAutoScalingConfiguration | 9 | 5 | 1 | 3 | 0 | ✅ | — |
| AwsAppRunnerObservabilityConfiguration | 5 | 2 | 0 | 3 | 0 | ✅ | ✅ pulumi, terraform |
| AwsAppRunnerService | 55 | 6 | 38 | 11 | 0 | ✅ | — |
| AwsAppRunnerVpcConnector | 6 | 1 | 2 | 3 | 0 | ✅ | ✅ pulumi, terraform |
| AwsAthenaWorkgroup | 37 | 4 | 0 | 0 | 59 | ❌ | ✅ pulumi, terraform |
| AwsAutoScalingGroup | 217 | 54 | 146 | 17 | 0 | ✅ | ✅ pulumi, terraform |
| AwsBatchComputeEnvironment | 33 | 20 | 6 | 7 | 0 | ✅ | ✅ pulumi, terraform |
| AwsBatchJobDefinition | 67 | 12 | 49 | 6 | 0 | ✅ | ✅ pulumi, terraform |
| AwsBatchJobQueue | 12 | 5 | 5 | 2 | 0 | ✅ | ✅ pulumi, terraform |
| AwsBatchSchedulingPolicy | 8 | 1 | 4 | 3 | 0 | ✅ | ✅ pulumi, terraform |
| AwsCertManagerCert | 18 | 7 | 7 | 4 | 0 | ✅ | ✅ pulumi, terraform |
| AwsClientVpn | 47 | 28 | 0 | 0 | 32 | ❌ | — |
| AwsCloudFront | 126 | 32 | 88 | 6 | 0 | ✅ | ✅ pulumi, terraform |
| AwsCloudwatchAlarm | 39 | 24 | 12 | 3 | 0 | ✅ | — |
| AwsCloudwatchCompositeAlarm | 13 | 10 | 0 | 3 | 0 | ✅ | ✅ pulumi, terraform |
| AwsCloudwatchLogGroup | 100 | 19 | 70 | 11 | 0 | ✅ | — |
| AwsCodeBuildProject | 115 | 97 | 11 | 7 | 0 | ✅ | ✅ pulumi, terraform |
| AwsCodePipeline | 84 | 4 | 75 | 5 | 0 | ✅ | ✅ pulumi, terraform |
| AwsCognitoIdentityProvider | 7 | 6 | 1 | 0 | 0 | ✅ | — |
| AwsCognitoResourceServer | 6 | 4 | 2 | 0 | 0 | ✅ | ✅ pulumi, terraform |
| AwsCognitoUserPool | 122 | 59 | 52 | 11 | 0 | ✅ | — |
| AwsCognitoUserPoolClient | 56 | 30 | 23 | 3 | 0 | ✅ | — |
| AwsDocumentDb | 74 | 40 | 8 | 26 | 0 | ✅ | ✅ pulumi, terraform |
| AwsDynamodb | 71 | 27 | 0 | 0 | 82 | ❌ | ✅ pulumi, terraform |
| AwsEc2Instance | 92 | 42 | 43 | 7 | 0 | ✅ | ✅ pulumi, terraform |
| AwsEcrRepo | 17 | 5 | 0 | 0 | 25 | ❌ | ✅ pulumi, terraform |
| AwsEcsCluster | 80 | 7 | 64 | 9 | 0 | ✅ | ✅ pulumi, terraform |
| AwsEcsService | 98 | 53 | 34 | 11 | 0 | ✅ | ✅ pulumi, terraform |
| AwsEcsTaskDefinition | 56 | 13 | 22 | 21 | 0 | ✅ | ✅ pulumi, terraform |
| AwsEfsAccessPoint | 11 | 9 | 0 | 0 | 2 | ❌ | ✅ pulumi, terraform |
| AwsEgressOnlyInternetGateway | 4 | 2 | 0 | 0 | 2 | ❌ | ✅ pulumi, terraform |
| AwsEksAccessEntry | 14 | 9 | 0 | 5 | 0 | ✅ | ✅ pulumi, terraform |
| AwsEksAddon | 14 | 10 | 2 | 2 | 0 | ✅ | ✅ pulumi, terraform |
| AwsEksCluster | 38 | 8 | 20 | 10 | 0 | ✅ | ✅ pulumi, terraform |
| AwsEksFargateProfile | 9 | 4 | 2 | 3 | 0 | ✅ | ✅ pulumi, terraform |
| AwsEksNodeGroup | 43 | 26 | 12 | 5 | 0 | ✅ | ✅ pulumi, terraform |
| AwsElasticFileSystem | 34 | 13 | 0 | 0 | 35 | ❌ | ✅ pulumi, terraform |
| AwsElasticIp | 12 | 4 | 0 | 0 | 8 | ❌ | ✅ pulumi, terraform |
| AwsElasticacheUser | 13 | 6 | 0 | 7 | 0 | ✅ | ✅ pulumi, terraform |
| AwsElasticacheUserGroup | 6 | 3 | 0 | 3 | 0 | ✅ | ✅ pulumi, terraform |
| AwsEventBridgeBus | 20 | 12 | 3 | 5 | 0 | ✅ | — |
| AwsEventBridgeRule | 65 | 53 | 5 | 7 | 0 | ✅ | ✅ pulumi, terraform |
| AwsFsxDataRepositoryAssociation | 11 | 7 | 2 | 2 | 0 | ✅ | — |
| AwsFsxLustreFileSystem | 34 | 29 | 3 | 2 | 0 | ✅ | — |
| AwsFsxOntapFileSystem | 21 | 18 | 1 | 2 | 0 | ✅ | — |
| AwsFsxOntapStorageVirtualMachine | 14 | 6 | 6 | 2 | 0 | ✅ | — |
| AwsFsxOntapVolume | 34 | 31 | 0 | 3 | 0 | ✅ | — |
| AwsFsxOpenzfsFileSystem | 35 | 31 | 2 | 2 | 0 | ✅ | — |
| AwsFsxWindowsFileSystem | 34 | 28 | 2 | 4 | 0 | ✅ | — |
| AwsGlobalAccelerator | 28 | 3 | 0 | 0 | 48 | ❌ | ✅ pulumi, terraform |
| AwsGlueCatalogDatabase | 15 | 10 | 0 | 0 | 7 | ❌ | ✅ pulumi, terraform |
| AwsHttpApiDomain | 26 | 7 | 11 | 8 | 0 | ✅ | ✅ pulumi, terraform |
| AwsHttpApiGateway | 91 | 42 | 11 | 38 | 0 | ✅ | ✅ pulumi, terraform |
| AwsHttpApiVpcLink | 6 | 3 | 0 | 3 | 0 | ✅ | ✅ pulumi, terraform |
| AwsIamInstanceProfile | 6 | 2 | 0 | 4 | 0 | ✅ | ✅ pulumi, terraform |
| AwsIamOidcProvider | 5 | 3 | 0 | 0 | 3 | ❌ | ✅ pulumi, terraform |
| AwsIamPolicy | 8 | 2 | 1 | 5 | 0 | ✅ | ✅ pulumi, terraform |
| AwsIamRole | 16 | 5 | 0 | 0 | 20 | ❌ | ✅ pulumi, terraform |
| AwsIamUser | 15 | 3 | 4 | 8 | 0 | ✅ | ✅ pulumi, terraform |
| AwsInternetGateway | 4 | 2 | 0 | 0 | 2 | ❌ | ✅ pulumi, terraform |
| AwsKinesisFirehose | 337 | 1 | 268 | 68 | 0 | ✅ | ✅ pulumi, terraform |
| AwsKinesisStream | 17 | 8 | 0 | 0 | 12 | ❌ | ✅ pulumi, terraform |
| AwsKinesisStreamConsumer | 8 | 3 | 0 | 0 | 6 | ❌ | ✅ pulumi, terraform |
| AwsKmsKey | 29 | 15 | 6 | 8 | 0 | ✅ | ✅ pulumi, terraform |
| AwsLambda | 108 | 51 | 30 | 27 | 0 | ✅ | ✅ pulumi, terraform |
| AwsLambdaEventSourceMapping | 45 | 11 | 31 | 3 | 0 | ✅ | ✅ pulumi, terraform |
| AwsLaunchTemplate | 139 | 99 | 26 | 14 | 0 | ✅ | ✅ pulumi, terraform |
| AwsLbListener | 76 | 13 | 59 | 4 | 0 | ✅ | ✅ pulumi, terraform |
| AwsLbListenerRule | 61 | 3 | 55 | 3 | 0 | ✅ | ✅ pulumi, terraform |
| AwsLbTargetGroup | 48 | 34 | 7 | 7 | 0 | ✅ | ✅ pulumi, terraform |
| AwsMemcachedElasticache | 48 | 21 | 3 | 24 | 0 | ✅ | ✅ pulumi, terraform |
| AwsMemorydbAcl | 6 | 2 | 0 | 0 | 4 | ❌ | — |
| AwsMemorydbCluster | 46 | 31 | 0 | 0 | 18 | ❌ | — |
| AwsMemorydbUser | 7 | 4 | 0 | 0 | 3 | ❌ | — |
| AwsMskCluster | 50 | 9 | 0 | 0 | 74 | ❌ | partial: pulumi, terraform |
| AwsMskServerlessCluster | 7 | 1 | 0 | 0 | 8 | ❌ | partial: pulumi, terraform |
| AwsMwaaEnvironment | 38 | 32 | 0 | 0 | 9 | ❌ | partial: pulumi, terraform |
| AwsNatGateway | 15 | 10 | 3 | 2 | 0 | ✅ | ✅ pulumi, terraform |
| AwsNeptuneCluster | 70 | 34 | 0 | 0 | 52 | ❌ | ✅ pulumi, terraform |
| AwsNlb | 68 | 10 | 10 | 48 | 0 | ✅ | ✅ pulumi, terraform |
| AwsOpenSearchDomain | 77 | 51 | 0 | 0 | 48 | ❌ | partial: pulumi, terraform |
| AwsPlantonRunner | 216 | 9 | 0 | 0 | 216 | ❌ | ✅ pulumi, terraform |
| AwsRdsCluster | 143 | 84 | 19 | 40 | 0 | ✅ | ✅ pulumi, terraform |
| AwsRdsInstance | 119 | 65 | 20 | 34 | 0 | ✅ | ✅ pulumi, terraform |
| AwsRedisElasticache | 66 | 40 | 12 | 14 | 0 | ✅ | ✅ pulumi, terraform |
| AwsRedshiftCluster | 67 | 40 | 0 | 0 | 39 | ❌ | ✅ pulumi, terraform |
| AwsRedshiftServerlessNamespace | 15 | 10 | 0 | 0 | 5 | ❌ | ✅ pulumi, terraform |
| AwsRedshiftServerlessWorkgroup | 17 | 12 | 0 | 0 | 7 | ❌ | ✅ pulumi, terraform |
| AwsRoute53DnsRecord | 25 | 7 | 17 | 1 | 0 | ✅ | — |
| AwsRoute53HealthCheck | 23 | 17 | 2 | 4 | 0 | ✅ | — |
| AwsRoute53Zone | 17 | 5 | 5 | 7 | 0 | ✅ | — |
| AwsS3Bucket | 204 | 24 | 102 | 78 | 0 | ✅ | ✅ pulumi, terraform |
| AwsS3ObjectSet | 28 | 3 | 0 | 0 | 46 | ❌ | ✅ pulumi, terraform |
| AwsSagemakerDomain | 159 | 93 | 0 | 0 | 126 | ❌ | ✅ pulumi, terraform |
| AwsSecurityGroup | 10 | 4 | 0 | 0 | 26 | ❌ | ✅ pulumi, terraform |
| AwsServerlessElasticache | 19 | 12 | 4 | 3 | 0 | ✅ | ✅ pulumi, terraform |
| AwsSesConfigurationSet | 27 | 7 | 0 | 0 | 36 | ❌ | ✅ pulumi, terraform |
| AwsSesEmailIdentity | 19 | 9 | 0 | 0 | 18 | ❌ | ✅ pulumi, terraform |
| AwsSnsSubscription | 13 | 11 | 2 | 0 | 0 | ✅ | — |
| AwsSnsTopic | 33 | 10 | 18 | 5 | 0 | ✅ | — |
| AwsSqsQueue | 20 | 12 | 4 | 4 | 0 | ✅ | — |
| AwsStepFunction | 16 | 5 | 0 | 0 | 17 | ❌ | ✅ pulumi, terraform |
| AwsSubnet | 33 | 22 | 0 | 11 | 0 | ✅ | ✅ pulumi, terraform |
| AwsTransitGateway | 14 | 12 | 0 | 2 | 0 | ✅ | ✅ pulumi, terraform |
| AwsTransitGatewayRouteTable | 27 | 11 | 6 | 10 | 0 | ✅ | — |
| AwsTransitGatewayVpcAttachment | 12 | 8 | 2 | 2 | 0 | ✅ | ✅ pulumi, terraform |
| AwsVpc | 39 | 16 | 14 | 9 | 0 | ✅ | ✅ pulumi, terraform |
| AwsVpcEndpoint | 23 | 17 | 0 | 0 | 10 | ❌ | ✅ pulumi, terraform |
| AwsWafIpSet | 9 | 5 | 0 | 4 | 0 | ✅ | ✅ pulumi, terraform |
| AwsWafRegexPatternSet | 8 | 3 | 1 | 4 | 0 | ✅ | ✅ pulumi, terraform |
| AwsWafWebAcl | 9727 | 7 | 9685 | 35 | 0 | ✅ | ✅ pulumi, terraform |

## Breadth: every GA resource, one disposition

All resources of `aws@6.58.0` land in exactly one class:

| Disposition | Resources | Meaning |
|---|---|---|
| Modeled | 249 | consumed by a kind's Terraform module today |
| IAM-covered | 0 | per-resource IAM member/binding/policy triplets, covered by the owning kinds' additive `iam_members` fields |
| Composed | 12 | capability covered through an existing kind's surface rather than a kind of its own |
| Planned | 769 | judged to be covered by a planned kind or planned composition, not built yet |
| Deferred | 532 | deliberately not offered, each with the recorded reason |
| Excluded as deprecated | 129 | deprecated or superseded provider surface |
| **Total** | **1691** | |

## The enumerated record

The full per-resource record, so the accounting above is verifiable
rather than trusted.

### Modeled (249)

| Resource | Consuming kinds |
|---|---|
| `aws_acm_certificate` | consumed by AwsCertManagerCert |
| `aws_acm_certificate_validation` | consumed by AwsCertManagerCert |
| `aws_apigatewayv2_api` | consumed by AwsHttpApiGateway |
| `aws_apigatewayv2_api_mapping` | consumed by AwsHttpApiDomain |
| `aws_apigatewayv2_authorizer` | consumed by AwsHttpApiGateway |
| `aws_apigatewayv2_domain_name` | consumed by AwsHttpApiDomain |
| `aws_apigatewayv2_integration` | consumed by AwsHttpApiGateway |
| `aws_apigatewayv2_route` | consumed by AwsHttpApiGateway |
| `aws_apigatewayv2_routing_rule` | consumed by AwsHttpApiDomain |
| `aws_apigatewayv2_stage` | consumed by AwsHttpApiGateway |
| `aws_apigatewayv2_vpc_link` | consumed by AwsHttpApiVpcLink |
| `aws_appautoscaling_policy` | consumed by AwsEcsService |
| `aws_appautoscaling_target` | consumed by AwsEcsService |
| `aws_apprunner_auto_scaling_configuration_version` | consumed by AwsAppRunnerAutoScalingConfiguration |
| `aws_apprunner_custom_domain_association` | consumed by AwsAppRunnerService |
| `aws_apprunner_default_auto_scaling_configuration_version` | consumed by AwsAppRunnerAutoScalingConfiguration |
| `aws_apprunner_observability_configuration` | consumed by AwsAppRunnerObservabilityConfiguration |
| `aws_apprunner_service` | consumed by AwsAppRunnerService |
| `aws_apprunner_vpc_connector` | consumed by AwsAppRunnerVpcConnector |
| `aws_apprunner_vpc_ingress_connection` | consumed by AwsAppRunnerService |
| `aws_athena_workgroup` | consumed by AwsAthenaWorkgroup |
| `aws_autoscaling_group` | consumed by AwsAutoScalingGroup |
| `aws_autoscaling_lifecycle_hook` | consumed by AwsAutoScalingGroup |
| `aws_autoscaling_notification` | consumed by AwsAutoScalingGroup |
| `aws_autoscaling_policy` | consumed by AwsAutoScalingGroup |
| `aws_autoscaling_schedule` | consumed by AwsAutoScalingGroup |
| `aws_batch_compute_environment` | consumed by AwsBatchComputeEnvironment |
| `aws_batch_job_definition` | consumed by AwsBatchJobDefinition |
| `aws_batch_job_queue` | consumed by AwsBatchJobQueue |
| `aws_batch_scheduling_policy` | consumed by AwsBatchSchedulingPolicy |
| `aws_cloudfront_continuous_deployment_policy` | consumed by AwsCloudFront |
| `aws_cloudfront_distribution` | consumed by AwsCloudFront |
| `aws_cloudfront_monitoring_subscription` | consumed by AwsCloudFront |
| `aws_cloudfront_origin_access_control` | consumed by AwsCloudFront |
| `aws_cloudwatch_composite_alarm` | consumed by AwsCloudwatchCompositeAlarm |
| `aws_cloudwatch_event_archive` | consumed by AwsEventBridgeBus |
| `aws_cloudwatch_event_bus` | consumed by AwsEventBridgeBus |
| `aws_cloudwatch_event_bus_policy` | consumed by AwsEventBridgeBus |
| `aws_cloudwatch_event_rule` | consumed by AwsEventBridgeRule |
| `aws_cloudwatch_event_target` | consumed by AwsEventBridgeRule |
| `aws_cloudwatch_log_data_protection_policy` | consumed by AwsCloudwatchLogGroup |
| `aws_cloudwatch_log_group` | consumed by AwsCloudwatchLogGroup, AwsEcsTaskDefinition, AwsPlantonRunner |
| `aws_cloudwatch_log_index_policy` | consumed by AwsCloudwatchLogGroup |
| `aws_cloudwatch_log_metric_filter` | consumed by AwsCloudwatchLogGroup |
| `aws_cloudwatch_log_stream` | consumed by AwsCloudwatchLogGroup |
| `aws_cloudwatch_log_subscription_filter` | consumed by AwsCloudwatchLogGroup |
| `aws_cloudwatch_log_transformer` | consumed by AwsCloudwatchLogGroup |
| `aws_cloudwatch_metric_alarm` | consumed by AwsCloudwatchAlarm |
| `aws_codebuild_project` | consumed by AwsCodeBuildProject |
| `aws_codebuild_resource_policy` | consumed by AwsCodeBuildProject |
| `aws_codebuild_webhook` | consumed by AwsCodeBuildProject |
| `aws_codepipeline` | consumed by AwsCodePipeline |
| `aws_cognito_identity_provider` | consumed by AwsCognitoIdentityProvider |
| `aws_cognito_log_delivery_configuration` | consumed by AwsCognitoUserPool |
| `aws_cognito_resource_server` | consumed by AwsCognitoResourceServer |
| `aws_cognito_risk_configuration` | consumed by AwsCognitoUserPool, AwsCognitoUserPoolClient |
| `aws_cognito_user_group` | consumed by AwsCognitoUserPool |
| `aws_cognito_user_pool` | consumed by AwsCognitoUserPool |
| `aws_cognito_user_pool_client` | consumed by AwsCognitoUserPoolClient |
| `aws_cognito_user_pool_domain` | consumed by AwsCognitoUserPool |
| `aws_db_instance` | consumed by AwsRdsInstance |
| `aws_db_instance_role_association` | consumed by AwsRdsInstance |
| `aws_db_option_group` | consumed by AwsRdsInstance |
| `aws_db_parameter_group` | consumed by AwsRdsInstance |
| `aws_db_subnet_group` | consumed by AwsRdsCluster, AwsRdsInstance |
| `aws_docdb_cluster` | consumed by AwsDocumentDb |
| `aws_docdb_cluster_instance` | consumed by AwsDocumentDb |
| `aws_docdb_cluster_parameter_group` | consumed by AwsDocumentDb |
| `aws_docdb_subnet_group` | consumed by AwsDocumentDb |
| `aws_dynamodb_contributor_insights` | consumed by AwsDynamodb |
| `aws_dynamodb_kinesis_streaming_destination` | consumed by AwsDynamodb |
| `aws_dynamodb_resource_policy` | consumed by AwsDynamodb |
| `aws_dynamodb_table` | consumed by AwsDynamodb |
| `aws_ec2_client_vpn_authorization_rule` | consumed by AwsClientVpn |
| `aws_ec2_client_vpn_endpoint` | consumed by AwsClientVpn |
| `aws_ec2_client_vpn_network_association` | consumed by AwsClientVpn |
| `aws_ec2_client_vpn_route` | consumed by AwsClientVpn |
| `aws_ec2_transit_gateway` | consumed by AwsTransitGateway |
| `aws_ec2_transit_gateway_default_route_table_association` | consumed by AwsTransitGatewayRouteTable |
| `aws_ec2_transit_gateway_default_route_table_propagation` | consumed by AwsTransitGatewayRouteTable |
| `aws_ec2_transit_gateway_prefix_list_reference` | consumed by AwsTransitGatewayRouteTable |
| `aws_ec2_transit_gateway_route` | consumed by AwsTransitGatewayRouteTable |
| `aws_ec2_transit_gateway_route_table` | consumed by AwsTransitGatewayRouteTable |
| `aws_ec2_transit_gateway_route_table_association` | consumed by AwsTransitGatewayRouteTable |
| `aws_ec2_transit_gateway_route_table_propagation` | consumed by AwsTransitGatewayRouteTable |
| `aws_ec2_transit_gateway_vpc_attachment` | consumed by AwsTransitGatewayVpcAttachment |
| `aws_ecr_lifecycle_policy` | consumed by AwsEcrRepo |
| `aws_ecr_repository` | consumed by AwsEcrRepo |
| `aws_ecr_repository_policy` | consumed by AwsEcrRepo |
| `aws_ecs_capacity_provider` | consumed by AwsEcsCluster |
| `aws_ecs_cluster` | consumed by AwsEcsCluster, AwsPlantonRunner |
| `aws_ecs_cluster_capacity_providers` | consumed by AwsEcsCluster |
| `aws_ecs_service` | consumed by AwsEcsService, AwsPlantonRunner |
| `aws_ecs_task_definition` | consumed by AwsEcsTaskDefinition, AwsPlantonRunner |
| `aws_efs_access_point` | consumed by AwsEfsAccessPoint |
| `aws_efs_backup_policy` | consumed by AwsElasticFileSystem |
| `aws_efs_file_system` | consumed by AwsElasticFileSystem |
| `aws_efs_file_system_policy` | consumed by AwsElasticFileSystem |
| `aws_efs_mount_target` | consumed by AwsElasticFileSystem |
| `aws_efs_replication_configuration` | consumed by AwsElasticFileSystem |
| `aws_egress_only_internet_gateway` | consumed by AwsEgressOnlyInternetGateway |
| `aws_eip` | consumed by AwsElasticIp |
| `aws_eks_access_entry` | consumed by AwsEksAccessEntry |
| `aws_eks_access_policy_association` | consumed by AwsEksAccessEntry |
| `aws_eks_addon` | consumed by AwsEksAddon |
| `aws_eks_cluster` | consumed by AwsEksCluster |
| `aws_eks_fargate_profile` | consumed by AwsEksFargateProfile |
| `aws_eks_node_group` | consumed by AwsEksNodeGroup |
| `aws_elasticache_cluster` | consumed by AwsMemcachedElasticache |
| `aws_elasticache_parameter_group` | consumed by AwsMemcachedElasticache, AwsRedisElasticache |
| `aws_elasticache_replication_group` | consumed by AwsRedisElasticache |
| `aws_elasticache_serverless_cache` | consumed by AwsServerlessElasticache |
| `aws_elasticache_subnet_group` | consumed by AwsMemcachedElasticache, AwsRedisElasticache |
| `aws_elasticache_user` | consumed by AwsElasticacheUser |
| `aws_elasticache_user_group` | consumed by AwsElasticacheUserGroup |
| `aws_fsx_data_repository_association` | consumed by AwsFsxDataRepositoryAssociation |
| `aws_fsx_lustre_file_system` | consumed by AwsFsxLustreFileSystem |
| `aws_fsx_ontap_file_system` | consumed by AwsFsxOntapFileSystem |
| `aws_fsx_ontap_storage_virtual_machine` | consumed by AwsFsxOntapStorageVirtualMachine |
| `aws_fsx_ontap_volume` | consumed by AwsFsxOntapVolume |
| `aws_fsx_openzfs_file_system` | consumed by AwsFsxOpenzfsFileSystem |
| `aws_fsx_windows_file_system` | consumed by AwsFsxWindowsFileSystem |
| `aws_globalaccelerator_accelerator` | consumed by AwsGlobalAccelerator |
| `aws_globalaccelerator_endpoint_group` | consumed by AwsGlobalAccelerator |
| `aws_globalaccelerator_listener` | consumed by AwsGlobalAccelerator |
| `aws_glue_catalog_database` | consumed by AwsGlueCatalogDatabase |
| `aws_iam_access_key` | consumed by AwsIamUser |
| `aws_iam_instance_profile` | consumed by AwsIamInstanceProfile |
| `aws_iam_openid_connect_provider` | consumed by AwsIamOidcProvider |
| `aws_iam_policy` | consumed by AwsIamPolicy |
| `aws_iam_role` | consumed by AwsIamRole, AwsPlantonRunner |
| `aws_iam_role_policy` | consumed by AwsIamRole, AwsPlantonRunner |
| `aws_iam_role_policy_attachment` | consumed by AwsIamRole, AwsPlantonRunner |
| `aws_iam_user` | consumed by AwsIamUser |
| `aws_iam_user_policy` | consumed by AwsIamUser |
| `aws_iam_user_policy_attachment` | consumed by AwsIamUser |
| `aws_instance` | consumed by AwsEc2Instance |
| `aws_internet_gateway` | consumed by AwsInternetGateway |
| `aws_kinesis_firehose_delivery_stream` | consumed by AwsKinesisFirehose |
| `aws_kinesis_resource_policy` | consumed by AwsKinesisStream, AwsKinesisStreamConsumer |
| `aws_kinesis_stream` | consumed by AwsKinesisStream |
| `aws_kinesis_stream_consumer` | consumed by AwsKinesisStreamConsumer |
| `aws_kms_alias` | consumed by AwsKmsKey |
| `aws_kms_grant` | consumed by AwsKmsKey |
| `aws_kms_key` | consumed by AwsKmsKey |
| `aws_lambda_alias` | consumed by AwsLambda |
| `aws_lambda_event_source_mapping` | consumed by AwsLambdaEventSourceMapping |
| `aws_lambda_function` | consumed by AwsLambda |
| `aws_lambda_function_event_invoke_config` | consumed by AwsLambda |
| `aws_lambda_function_recursion_config` | consumed by AwsLambda |
| `aws_lambda_function_scaling_config` | consumed by AwsLambda |
| `aws_lambda_function_url` | consumed by AwsLambda |
| `aws_lambda_permission` | consumed by AwsLambda |
| `aws_lambda_provisioned_concurrency_config` | consumed by AwsLambda |
| `aws_lambda_runtime_management_config` | consumed by AwsLambda |
| `aws_launch_template` | consumed by AwsLaunchTemplate |
| `aws_lb` | consumed by AwsAlb, AwsNlb |
| `aws_lb_listener` | consumed by AwsLbListener |
| `aws_lb_listener_certificate` | consumed by AwsLbListener |
| `aws_lb_listener_rule` | consumed by AwsLbListenerRule |
| `aws_lb_target_group` | consumed by AwsLbTargetGroup |
| `aws_lb_target_group_attachment` | consumed by AwsLbTargetGroup |
| `aws_memorydb_acl` | consumed by AwsMemorydbAcl |
| `aws_memorydb_cluster` | consumed by AwsMemorydbCluster |
| `aws_memorydb_parameter_group` | consumed by AwsMemorydbCluster |
| `aws_memorydb_subnet_group` | consumed by AwsMemorydbCluster |
| `aws_memorydb_user` | consumed by AwsMemorydbUser |
| `aws_msk_cluster` | consumed by AwsMskCluster |
| `aws_msk_cluster_policy` | consumed by AwsMskCluster |
| `aws_msk_configuration` | consumed by AwsMskCluster |
| `aws_msk_serverless_cluster` | consumed by AwsMskServerlessCluster |
| `aws_msk_single_scram_secret_association` | consumed by AwsMskCluster |
| `aws_mwaa_environment` | consumed by AwsMwaaEnvironment |
| `aws_nat_gateway` | consumed by AwsNatGateway |
| `aws_neptune_cluster` | consumed by AwsNeptuneCluster |
| `aws_neptune_cluster_instance` | consumed by AwsNeptuneCluster |
| `aws_neptune_cluster_parameter_group` | consumed by AwsNeptuneCluster |
| `aws_neptune_subnet_group` | consumed by AwsNeptuneCluster |
| `aws_opensearch_domain` | consumed by AwsOpenSearchDomain |
| `aws_rds_cluster` | consumed by AwsRdsCluster |
| `aws_rds_cluster_activity_stream` | consumed by AwsRdsCluster |
| `aws_rds_cluster_endpoint` | consumed by AwsRdsCluster |
| `aws_rds_cluster_instance` | consumed by AwsRdsCluster |
| `aws_rds_cluster_parameter_group` | consumed by AwsRdsCluster |
| `aws_rds_cluster_role_association` | consumed by AwsRdsCluster |
| `aws_redshift_cluster` | consumed by AwsRedshiftCluster |
| `aws_redshift_logging` | consumed by AwsRedshiftCluster |
| `aws_redshift_parameter_group` | consumed by AwsRedshiftCluster |
| `aws_redshift_snapshot_copy` | consumed by AwsRedshiftCluster |
| `aws_redshift_subnet_group` | consumed by AwsRedshiftCluster |
| `aws_redshiftserverless_namespace` | consumed by AwsRedshiftServerlessNamespace |
| `aws_redshiftserverless_workgroup` | consumed by AwsRedshiftServerlessWorkgroup |
| `aws_route53_health_check` | consumed by AwsRoute53HealthCheck |
| `aws_route53_hosted_zone_dnssec` | consumed by AwsRoute53Zone |
| `aws_route53_key_signing_key` | consumed by AwsRoute53Zone |
| `aws_route53_query_log` | consumed by AwsRoute53Zone |
| `aws_route53_record` | consumed by AwsAlb, AwsCertManagerCert, AwsNlb, AwsRoute53DnsRecord |
| `aws_route53_zone` | consumed by AwsRoute53Zone |
| `aws_route_table` | consumed by AwsSubnet |
| `aws_route_table_association` | consumed by AwsSubnet |
| `aws_s3_bucket` | consumed by AwsS3Bucket |
| `aws_s3_bucket_abac` | consumed by AwsS3Bucket |
| `aws_s3_bucket_accelerate_configuration` | consumed by AwsS3Bucket |
| `aws_s3_bucket_acl` | consumed by AwsS3Bucket |
| `aws_s3_bucket_analytics_configuration` | consumed by AwsS3Bucket |
| `aws_s3_bucket_cors_configuration` | consumed by AwsS3Bucket |
| `aws_s3_bucket_intelligent_tiering_configuration` | consumed by AwsS3Bucket |
| `aws_s3_bucket_inventory` | consumed by AwsS3Bucket |
| `aws_s3_bucket_lifecycle_configuration` | consumed by AwsS3Bucket |
| `aws_s3_bucket_logging` | consumed by AwsS3Bucket |
| `aws_s3_bucket_metadata_configuration` | consumed by AwsS3Bucket |
| `aws_s3_bucket_metric` | consumed by AwsS3Bucket |
| `aws_s3_bucket_notification` | consumed by AwsS3Bucket |
| `aws_s3_bucket_object_lock_configuration` | consumed by AwsS3Bucket |
| `aws_s3_bucket_ownership_controls` | consumed by AwsS3Bucket |
| `aws_s3_bucket_policy` | consumed by AwsS3Bucket |
| `aws_s3_bucket_public_access_block` | consumed by AwsS3Bucket |
| `aws_s3_bucket_replication_configuration` | consumed by AwsS3Bucket |
| `aws_s3_bucket_request_payment_configuration` | consumed by AwsS3Bucket |
| `aws_s3_bucket_server_side_encryption_configuration` | consumed by AwsS3Bucket |
| `aws_s3_bucket_versioning` | consumed by AwsS3Bucket |
| `aws_s3_bucket_website_configuration` | consumed by AwsS3Bucket |
| `aws_s3_object` | consumed by AwsS3ObjectSet |
| `aws_sagemaker_domain` | consumed by AwsSagemakerDomain |
| `aws_secretsmanager_secret` | consumed by AwsPlantonRunner |
| `aws_secretsmanager_secret_version` | consumed by AwsPlantonRunner |
| `aws_security_group` | consumed by AwsPlantonRunner, AwsSecurityGroup |
| `aws_sesv2_configuration_set` | consumed by AwsSesConfigurationSet |
| `aws_sesv2_configuration_set_event_destination` | consumed by AwsSesConfigurationSet |
| `aws_sesv2_email_identity` | consumed by AwsSesEmailIdentity |
| `aws_sesv2_email_identity_feedback_attributes` | consumed by AwsSesEmailIdentity |
| `aws_sesv2_email_identity_mail_from_attributes` | consumed by AwsSesEmailIdentity |
| `aws_sesv2_email_identity_policy` | consumed by AwsSesEmailIdentity |
| `aws_sfn_state_machine` | consumed by AwsStepFunction |
| `aws_sns_topic` | consumed by AwsSnsTopic |
| `aws_sns_topic_data_protection_policy` | consumed by AwsSnsTopic |
| `aws_sns_topic_subscription` | consumed by AwsSnsSubscription |
| `aws_sqs_queue` | consumed by AwsSqsQueue |
| `aws_subnet` | consumed by AwsSubnet |
| `aws_vpc` | consumed by AwsVpc |
| `aws_vpc_encryption_control` | consumed by AwsVpc |
| `aws_vpc_endpoint` | consumed by AwsVpcEndpoint |
| `aws_vpc_ipv4_cidr_block_association` | consumed by AwsVpc |
| `aws_vpc_ipv6_cidr_block_association` | consumed by AwsVpc |
| `aws_wafv2_ip_set` | consumed by AwsWafIpSet |
| `aws_wafv2_regex_pattern_set` | consumed by AwsWafRegexPatternSet |
| `aws_wafv2_web_acl` | consumed by AwsWafWebAcl |
| `aws_wafv2_web_acl_association` | consumed by AwsAlb, AwsAppRunnerService |
| `aws_wafv2_web_acl_logging_configuration` | consumed by AwsWafWebAcl |

### Composed (12)

| Resource | Recorded reason |
|---|---|
| `aws_autoscaling_attachment` | covered by AwsAutoScalingGroup: target_groups registers ALB/NLB target groups and traffic_sources covers Classic ELBs -- the standalone attachment is the imperative pattern for a group that owns its attachments |
| `aws_autoscaling_traffic_source_attachment` | covered by AwsAutoScalingGroup.traffic_sources -- the standalone attachment is the imperative pattern for a group that owns its traffic sources |
| `aws_cloudwatch_event_permission` | per-statement delivery (PutPermission with a StatementId) of the same bus policy AwsEventBridgeBus.spec.resource_policy models as one whole document via aws_cloudwatch_event_bus_policy -- mixing them fights over one policy (the bus-policy delete issues RemoveAllPermissions, wiping permission-managed statements) |
| `aws_elasticache_user_group_association` | covered by AwsElasticacheUserGroup's declarative user_ids membership -- the standalone one-user-at-a-time association is the imperative alternative to membership the group already owns (the autoscaling-attachment class); it attaches users to GROUPS, not to the Redis kind the prior reason named |
| `aws_kms_key_policy` | covered by AwsKmsKey spec.policy -- the standalone resource is the detached-management pattern for keys owned elsewhere |
| `aws_nat_gateway_eip_association` | covered by AwsNatGateway's existing EIP fields -- secondary_allocation_ids (zonal public gateways) and availability_zone_addresses[].allocation_ids (regional gateways) declare the same associations declaratively; the standalone association resource exists for imperatively attaching EIPs to gateways not owned by the same configuration, an anti-pattern for a kind that owns its gateway |
| `aws_sns_topic_policy` | out-of-band delivery of the topic's Policy attribute, which AwsSnsTopic.spec.policy models inline on aws_sns_topic -- both write the same single attribute (the standalone resource cannot even remove a policy on delete; it writes back a synthesized owner-default document, per the provider source) |
| `aws_sqs_queue_policy` | out-of-band delivery of the queue's Policy attribute, which AwsSqsQueue.spec.policy models inline on aws_sqs_queue -- both mechanisms upsert the same single queue attribute and fight over it when mixed (the attribute-splitter class) |
| `aws_sqs_queue_redrive_allow_policy` | out-of-band delivery of the queue's RedriveAllowPolicy attribute, which AwsSqsQueue.spec.redrive_allow_policy models inline (typed redrivePermission + sourceQueueArns) -- same single-attribute upsert, same mixing conflict (the attribute-splitter class) |
| `aws_sqs_queue_redrive_policy` | out-of-band delivery of the queue's RedrivePolicy attribute, which AwsSqsQueue.spec.dead_letter_config models inline (chart-wired deadLetterTargetArn + maxReceiveCount) -- same single-attribute upsert, same mixing conflict (the attribute-splitter class) |
| `aws_wafv2_web_acl_rule` | covered by AwsWafWebAcl.spec.rules -- this satellite manages a single rule of an existing web ACL out-of-band, an alternative delivery mechanism for the same statement grammar the kind models inline in full; mixing out-of-band rules with an ACL whose rules are declared inline fights over one rule set |
| `aws_wafv2_web_acl_rule_group_association` | covered by AwsWafWebAcl.spec.rules (the rule_group_reference and managed_rule_group arms with rule_action_overrides) -- this satellite injects a group-reference rule into an existing web ACL out-of-band; the kind models the same attachment inline, and mixing the two fights over one rule set |

### Planned (769)

| Resource | Recorded reason |
|---|---|
| `aws_accessanalyzer_analyzer` | judged as a planned AwsIamAccessAnalyzer kind (analyzer with archive rules folding in) |
| `aws_accessanalyzer_archive_rule` | judged as a planned AwsIamAccessAnalyzer kind (analyzer with archive rules folding in) |
| `aws_account_alternate_contact` | account-level contact and region-enablement settings fold into the planned AwsOrganizationAccount kind |
| `aws_account_primary_contact` | account-level contact and region-enablement settings fold into the planned AwsOrganizationAccount kind |
| `aws_account_region` | account-level contact and region-enablement settings fold into the planned AwsOrganizationAccount kind |
| `aws_acmpca_certificate` | judged as a planned AwsPrivateCa kind (certificate authority, CA certificate, issued certificates, permission, policy) |
| `aws_acmpca_certificate_authority` | judged as a planned AwsPrivateCa kind (certificate authority, CA certificate, issued certificates, permission, policy) |
| `aws_acmpca_certificate_authority_certificate` | judged as a planned AwsPrivateCa kind (certificate authority, CA certificate, issued certificates, permission, policy) |
| `aws_acmpca_permission` | judged as a planned AwsPrivateCa kind (certificate authority, CA certificate, issued certificates, permission, policy) |
| `aws_acmpca_policy` | judged as a planned AwsPrivateCa kind (certificate authority, CA certificate, issued certificates, permission, policy) |
| `aws_ami` | judged as a planned AwsAmi kind (image, copy, launch permissions; instance capture folds in) |
| `aws_ami_copy` | judged as a planned AwsAmi kind (image, copy, launch permissions; instance capture folds in) |
| `aws_ami_from_instance` | judged as a planned AwsAmi kind (image, copy, launch permissions; instance capture folds in) |
| `aws_ami_launch_permission` | judged as a planned AwsAmi kind (image, copy, launch permissions; instance capture folds in) |
| `aws_amplify_app` | judged as a planned AwsAmplifyApp kind (app with branches, domain associations, webhooks, backend environments) |
| `aws_amplify_backend_environment` | judged as a planned AwsAmplifyApp kind (app with branches, domain associations, webhooks, backend environments) |
| `aws_amplify_branch` | judged as a planned AwsAmplifyApp kind (app with branches, domain associations, webhooks, backend environments) |
| `aws_amplify_domain_association` | judged as a planned AwsAmplifyApp kind (app with branches, domain associations, webhooks, backend environments) |
| `aws_amplify_webhook` | judged as a planned AwsAmplifyApp kind (app with branches, domain associations, webhooks, backend environments) |
| `aws_api_gateway_account` | judged as a planned AwsRestApiGateway kind (REST v1 API: resources, methods, integrations, deployments, stages, authorizers, models, validators, policies and account settings fold in) |
| `aws_api_gateway_api_key` | judged as a planned AwsRestApiUsagePlan kind (usage plans, plan keys, API keys) |
| `aws_api_gateway_authorizer` | judged as a planned AwsRestApiGateway kind (REST v1 API: resources, methods, integrations, deployments, stages, authorizers, models, validators, policies and account settings fold in) |
| `aws_api_gateway_base_path_mapping` | judged as a planned AwsRestApiDomain kind (custom domains with base-path mappings and access associations) |
| `aws_api_gateway_client_certificate` | judged as a planned AwsRestApiGateway kind (REST v1 API: resources, methods, integrations, deployments, stages, authorizers, models, validators, policies and account settings fold in) |
| `aws_api_gateway_deployment` | judged as a planned AwsRestApiGateway kind (REST v1 API: resources, methods, integrations, deployments, stages, authorizers, models, validators, policies and account settings fold in) |
| `aws_api_gateway_documentation_part` | judged as a planned AwsRestApiGateway kind (REST v1 API: resources, methods, integrations, deployments, stages, authorizers, models, validators, policies and account settings fold in) |
| `aws_api_gateway_documentation_version` | judged as a planned AwsRestApiGateway kind (REST v1 API: resources, methods, integrations, deployments, stages, authorizers, models, validators, policies and account settings fold in) |
| `aws_api_gateway_domain_name` | judged as a planned AwsRestApiDomain kind (custom domains with base-path mappings and access associations) |
| `aws_api_gateway_domain_name_access_association` | judged as a planned AwsRestApiDomain kind (custom domains with base-path mappings and access associations) |
| `aws_api_gateway_gateway_response` | judged as a planned AwsRestApiGateway kind (REST v1 API: resources, methods, integrations, deployments, stages, authorizers, models, validators, policies and account settings fold in) |
| `aws_api_gateway_integration` | judged as a planned AwsRestApiGateway kind (REST v1 API: resources, methods, integrations, deployments, stages, authorizers, models, validators, policies and account settings fold in) |
| `aws_api_gateway_integration_response` | judged as a planned AwsRestApiGateway kind (REST v1 API: resources, methods, integrations, deployments, stages, authorizers, models, validators, policies and account settings fold in) |
| `aws_api_gateway_method` | judged as a planned AwsRestApiGateway kind (REST v1 API: resources, methods, integrations, deployments, stages, authorizers, models, validators, policies and account settings fold in) |
| `aws_api_gateway_method_response` | judged as a planned AwsRestApiGateway kind (REST v1 API: resources, methods, integrations, deployments, stages, authorizers, models, validators, policies and account settings fold in) |
| `aws_api_gateway_method_settings` | judged as a planned AwsRestApiGateway kind (REST v1 API: resources, methods, integrations, deployments, stages, authorizers, models, validators, policies and account settings fold in) |
| `aws_api_gateway_model` | judged as a planned AwsRestApiGateway kind (REST v1 API: resources, methods, integrations, deployments, stages, authorizers, models, validators, policies and account settings fold in) |
| `aws_api_gateway_request_validator` | judged as a planned AwsRestApiGateway kind (REST v1 API: resources, methods, integrations, deployments, stages, authorizers, models, validators, policies and account settings fold in) |
| `aws_api_gateway_resource` | judged as a planned AwsRestApiGateway kind (REST v1 API: resources, methods, integrations, deployments, stages, authorizers, models, validators, policies and account settings fold in) |
| `aws_api_gateway_rest_api` | judged as a planned AwsRestApiGateway kind (REST v1 API: resources, methods, integrations, deployments, stages, authorizers, models, validators, policies and account settings fold in) |
| `aws_api_gateway_rest_api_policy` | judged as a planned AwsRestApiGateway kind (REST v1 API: resources, methods, integrations, deployments, stages, authorizers, models, validators, policies and account settings fold in) |
| `aws_api_gateway_rest_api_put` | judged as a planned AwsRestApiGateway kind (REST v1 API: resources, methods, integrations, deployments, stages, authorizers, models, validators, policies and account settings fold in) |
| `aws_api_gateway_stage` | judged as a planned AwsRestApiGateway kind (REST v1 API: resources, methods, integrations, deployments, stages, authorizers, models, validators, policies and account settings fold in) |
| `aws_api_gateway_usage_plan` | judged as a planned AwsRestApiUsagePlan kind (usage plans, plan keys, API keys) |
| `aws_api_gateway_usage_plan_key` | judged as a planned AwsRestApiUsagePlan kind (usage plans, plan keys, API keys) |
| `aws_api_gateway_vpc_link` | judged as a planned AwsRestApiGateway kind (REST v1 API: resources, methods, integrations, deployments, stages, authorizers, models, validators, policies and account settings fold in) |
| `aws_apigatewayv2_deployment` | judged into the planned AwsWebSocketApiGateway kind (WebSocket APIs are their own protocol surface -- AwsHttpApiGateway records that boundary in its spec); explicit deployments are the WebSocket publish mechanism, while HTTP APIs deploy via the stage's auto_deploy, already modeled on AwsHttpApiGateway |
| `aws_apigatewayv2_integration_response` | judged into the planned AwsWebSocketApiGateway kind -- integration responses exist only in the WebSocket request/response model; HTTP APIs have no integration-response concept (they transform via response_parameters, already modeled on AwsHttpApiGateway) |
| `aws_apigatewayv2_model` | judged into the planned AwsWebSocketApiGateway kind -- request models are consumed via route request_models/model_selection_expression, WebSocket-only route surface; HTTP APIs never reference models |
| `aws_apigatewayv2_route_response` | judged into the planned AwsWebSocketApiGateway kind -- route responses implement WebSocket two-way communication (only the $default route response is supported); HTTP APIs have no route-response concept |
| `aws_appautoscaling_scheduled_action` | application-autoscaling scheduled actions fold into the consuming kinds' scaling surface (AwsDynamodb, AwsEcsService) |
| `aws_appconfig_application` | judged as a planned AwsAppConfig kind (application, environments, profiles, versions, deployment strategies, extensions) |
| `aws_appconfig_configuration_profile` | judged as a planned AwsAppConfig kind (application, environments, profiles, versions, deployment strategies, extensions) |
| `aws_appconfig_deployment` | judged as a planned AwsAppConfig kind (application, environments, profiles, versions, deployment strategies, extensions) |
| `aws_appconfig_deployment_strategy` | judged as a planned AwsAppConfig kind (application, environments, profiles, versions, deployment strategies, extensions) |
| `aws_appconfig_environment` | judged as a planned AwsAppConfig kind (application, environments, profiles, versions, deployment strategies, extensions) |
| `aws_appconfig_extension` | judged as a planned AwsAppConfig kind (application, environments, profiles, versions, deployment strategies, extensions) |
| `aws_appconfig_extension_association` | judged as a planned AwsAppConfig kind (application, environments, profiles, versions, deployment strategies, extensions) |
| `aws_appconfig_hosted_configuration_version` | judged as a planned AwsAppConfig kind (application, environments, profiles, versions, deployment strategies, extensions) |
| `aws_apprunner_connection` | future AwsAppRunnerConnection kind: account-scoped GitHub/Bitbucket authorization shared by many services and referenced by AwsAppRunnerService code sources (connection_arn); its one-time OAuth handshake belongs to that kind's own lifecycle story (the codebuild source-credential class) |
| `aws_appsync_api` | judged as a planned AwsAppSyncApi kind (GraphQL and Events APIs: data sources, resolvers, functions, types, channel namespaces, caches, keys, domains) |
| `aws_appsync_api_cache` | judged as a planned AwsAppSyncApi kind (GraphQL and Events APIs: data sources, resolvers, functions, types, channel namespaces, caches, keys, domains) |
| `aws_appsync_api_key` | judged as a planned AwsAppSyncApi kind (GraphQL and Events APIs: data sources, resolvers, functions, types, channel namespaces, caches, keys, domains) |
| `aws_appsync_channel_namespace` | judged as a planned AwsAppSyncApi kind (GraphQL and Events APIs: data sources, resolvers, functions, types, channel namespaces, caches, keys, domains) |
| `aws_appsync_datasource` | judged as a planned AwsAppSyncApi kind (GraphQL and Events APIs: data sources, resolvers, functions, types, channel namespaces, caches, keys, domains) |
| `aws_appsync_domain_name` | judged as a planned AwsAppSyncApi kind (GraphQL and Events APIs: data sources, resolvers, functions, types, channel namespaces, caches, keys, domains) |
| `aws_appsync_domain_name_api_association` | judged as a planned AwsAppSyncApi kind (GraphQL and Events APIs: data sources, resolvers, functions, types, channel namespaces, caches, keys, domains) |
| `aws_appsync_function` | judged as a planned AwsAppSyncApi kind (GraphQL and Events APIs: data sources, resolvers, functions, types, channel namespaces, caches, keys, domains) |
| `aws_appsync_graphql_api` | judged as a planned AwsAppSyncApi kind (GraphQL and Events APIs: data sources, resolvers, functions, types, channel namespaces, caches, keys, domains) |
| `aws_appsync_resolver` | judged as a planned AwsAppSyncApi kind (GraphQL and Events APIs: data sources, resolvers, functions, types, channel namespaces, caches, keys, domains) |
| `aws_appsync_source_api_association` | judged as a planned AwsAppSyncApi kind (GraphQL and Events APIs: data sources, resolvers, functions, types, channel namespaces, caches, keys, domains) |
| `aws_appsync_type` | judged as a planned AwsAppSyncApi kind (GraphQL and Events APIs: data sources, resolvers, functions, types, channel namespaces, caches, keys, domains) |
| `aws_athena_data_catalog` | federated data catalogs fold into the existing AwsAthenaWorkgroup kind as its spec deepens |
| `aws_backup_framework` | judged as a planned AwsBackupPlan kind (plans, selections, frameworks, report and restore-testing plans) |
| `aws_backup_global_settings` | judged as a planned AwsBackupVault kind (vaults incl. air-gapped, lock, policy, notifications; account-wide backup settings fold in) |
| `aws_backup_logically_air_gapped_vault` | judged as a planned AwsBackupVault kind (vaults incl. air-gapped, lock, policy, notifications; account-wide backup settings fold in) |
| `aws_backup_plan` | judged as a planned AwsBackupPlan kind (plans, selections, frameworks, report and restore-testing plans) |
| `aws_backup_region_settings` | judged as a planned AwsBackupVault kind (vaults incl. air-gapped, lock, policy, notifications; account-wide backup settings fold in) |
| `aws_backup_report_plan` | judged as a planned AwsBackupPlan kind (plans, selections, frameworks, report and restore-testing plans) |
| `aws_backup_restore_testing_plan` | judged as a planned AwsBackupPlan kind (plans, selections, frameworks, report and restore-testing plans) |
| `aws_backup_restore_testing_selection` | judged as a planned AwsBackupPlan kind (plans, selections, frameworks, report and restore-testing plans) |
| `aws_backup_selection` | judged as a planned AwsBackupPlan kind (plans, selections, frameworks, report and restore-testing plans) |
| `aws_backup_vault` | judged as a planned AwsBackupVault kind (vaults incl. air-gapped, lock, policy, notifications; account-wide backup settings fold in) |
| `aws_backup_vault_lock_configuration` | judged as a planned AwsBackupVault kind (vaults incl. air-gapped, lock, policy, notifications; account-wide backup settings fold in) |
| `aws_backup_vault_notifications` | judged as a planned AwsBackupVault kind (vaults incl. air-gapped, lock, policy, notifications; account-wide backup settings fold in) |
| `aws_backup_vault_policy` | judged as a planned AwsBackupVault kind (vaults incl. air-gapped, lock, policy, notifications; account-wide backup settings fold in) |
| `aws_bcmdataexports_export` | billing data exports fold into the planned cost-reporting kinds (with aws_cur_report_definition) |
| `aws_bedrock_custom_model` | judged as a planned AwsBedrockCustomModel kind |
| `aws_bedrock_foundation_model_agreement` | judged as a planned AwsBedrockModelAccess kind (foundation-model agreements, use cases, invocation logging) |
| `aws_bedrock_guardrail` | judged as a planned AwsBedrockGuardrail kind (guardrail with versions) |
| `aws_bedrock_guardrail_version` | judged as a planned AwsBedrockGuardrail kind (guardrail with versions) |
| `aws_bedrock_inference_profile` | judged as a planned AwsBedrockInferenceProfile kind |
| `aws_bedrock_model_invocation_logging_configuration` | judged as a planned AwsBedrockModelAccess kind (foundation-model agreements, use cases, invocation logging) |
| `aws_bedrock_provisioned_model_throughput` | judged as a planned AwsBedrockProvisionedThroughput kind |
| `aws_bedrock_use_case_for_model_access` | judged as a planned AwsBedrockModelAccess kind (foundation-model agreements, use cases, invocation logging) |
| `aws_bedrockagent_agent` | judged as a planned AwsBedrockAgent kind (agent with action groups, aliases, collaborators, knowledge-base associations) |
| `aws_bedrockagent_agent_action_group` | judged as a planned AwsBedrockAgent kind (agent with action groups, aliases, collaborators, knowledge-base associations) |
| `aws_bedrockagent_agent_alias` | judged as a planned AwsBedrockAgent kind (agent with action groups, aliases, collaborators, knowledge-base associations) |
| `aws_bedrockagent_agent_collaborator` | judged as a planned AwsBedrockAgent kind (agent with action groups, aliases, collaborators, knowledge-base associations) |
| `aws_bedrockagent_agent_knowledge_base_association` | judged as a planned AwsBedrockAgent kind (agent with action groups, aliases, collaborators, knowledge-base associations) |
| `aws_bedrockagent_data_source` | judged as a planned AwsBedrockKnowledgeBase kind (knowledge base with data sources) |
| `aws_bedrockagent_flow` | judged as a planned AwsBedrockFlow kind |
| `aws_bedrockagent_knowledge_base` | judged as a planned AwsBedrockKnowledgeBase kind (knowledge base with data sources) |
| `aws_bedrockagent_prompt` | judged as a planned AwsBedrockPrompt kind |
| `aws_bedrockagentcore_agent_runtime` | judged as a planned AwsBedrockAgentCoreRuntime kind (agent runtimes with endpoints; evaluation and resource-policy companions fold in) |
| `aws_bedrockagentcore_agent_runtime_endpoint` | judged as a planned AwsBedrockAgentCoreRuntime kind (agent runtimes with endpoints; evaluation and resource-policy companions fold in) |
| `aws_bedrockagentcore_api_key_credential_provider` | judged as a planned AwsBedrockAgentCoreIdentity kind (workload identities, credential providers, token vault, policy engine) |
| `aws_bedrockagentcore_browser` | judged as a planned AwsBedrockAgentCoreTools kind (browsers with profiles, code interpreters) |
| `aws_bedrockagentcore_browser_profile` | judged as a planned AwsBedrockAgentCoreTools kind (browsers with profiles, code interpreters) |
| `aws_bedrockagentcore_code_interpreter` | judged as a planned AwsBedrockAgentCoreTools kind (browsers with profiles, code interpreters) |
| `aws_bedrockagentcore_evaluator` | judged as a planned AwsBedrockAgentCoreRuntime kind (agent runtimes with endpoints; evaluation and resource-policy companions fold in) |
| `aws_bedrockagentcore_gateway` | judged as a planned AwsBedrockAgentCoreGateway kind (gateways with targets) |
| `aws_bedrockagentcore_gateway_target` | judged as a planned AwsBedrockAgentCoreGateway kind (gateways with targets) |
| `aws_bedrockagentcore_harness` | judged as a planned AwsBedrockAgentCoreRuntime kind (agent runtimes with endpoints; evaluation and resource-policy companions fold in) |
| `aws_bedrockagentcore_memory` | judged as a planned AwsBedrockAgentCoreMemory kind (memories with strategies) |
| `aws_bedrockagentcore_memory_strategy` | judged as a planned AwsBedrockAgentCoreMemory kind (memories with strategies) |
| `aws_bedrockagentcore_oauth2_credential_provider` | judged as a planned AwsBedrockAgentCoreIdentity kind (workload identities, credential providers, token vault, policy engine) |
| `aws_bedrockagentcore_online_evaluation_config` | judged as a planned AwsBedrockAgentCoreRuntime kind (agent runtimes with endpoints; evaluation and resource-policy companions fold in) |
| `aws_bedrockagentcore_policy` | judged as a planned AwsBedrockAgentCoreIdentity kind (workload identities, credential providers, token vault, policy engine) |
| `aws_bedrockagentcore_policy_engine` | judged as a planned AwsBedrockAgentCoreIdentity kind (workload identities, credential providers, token vault, policy engine) |
| `aws_bedrockagentcore_resource_policy` | judged as a planned AwsBedrockAgentCoreRuntime kind (agent runtimes with endpoints; evaluation and resource-policy companions fold in) |
| `aws_bedrockagentcore_token_vault_cmk` | judged as a planned AwsBedrockAgentCoreIdentity kind (workload identities, credential providers, token vault, policy engine) |
| `aws_bedrockagentcore_workload_identity` | judged as a planned AwsBedrockAgentCoreIdentity kind (workload identities, credential providers, token vault, policy engine) |
| `aws_budgets_budget` | judged as a planned AwsBudget kind (budgets with actions) |
| `aws_budgets_budget_action` | judged as a planned AwsBudget kind (budgets with actions) |
| `aws_ce_anomaly_monitor` | judged as a planned AwsCostAnomalyMonitor kind (monitors with subscriptions) |
| `aws_ce_anomaly_subscription` | judged as a planned AwsCostAnomalyMonitor kind (monitors with subscriptions) |
| `aws_ce_cost_allocation_tag` | judged as a planned AwsCostCategory kind (cost categories with allocation tags) |
| `aws_ce_cost_category` | judged as a planned AwsCostCategory kind (cost categories with allocation tags) |
| `aws_cloudfront_anycast_ip_list` | account-scoped shared Anycast static IP list referenced by many distributions; judged as a planned AwsCloudFrontAnycastIpList kind -- AwsCloudFront models the attachment (anycast_ip_list_id) today |
| `aws_cloudfront_cache_policy` | account-scoped shared cache policy referenced by many behaviors across distributions; judged as a planned AwsCloudFrontCachePolicy kind -- AwsCloudFront behaviors model the attachment (cache_policy_id) today |
| `aws_cloudfront_connection_function` | judged as a planned AwsCloudFrontFunction kind (edge functions incl. connection functions and the key-value store); AwsCloudFront models the attachment (connection_function_id) today |
| `aws_cloudfront_field_level_encryption_config` | account-scoped shared field-level-encryption configuration referenced by many behaviors; judged as a planned AwsCloudFrontFieldLevelEncryption kind together with its profile resource -- AwsCloudFront behaviors model the attachment (field_level_encryption_id) today |
| `aws_cloudfront_field_level_encryption_profile` | account-scoped shared field-level-encryption profile (the public-key-to-field binding the config consumes); judged as part of the planned AwsCloudFrontFieldLevelEncryption kind |
| `aws_cloudfront_function` | judged as a planned AwsCloudFrontFunction kind (edge functions incl. connection functions and the key-value store); AwsCloudFront behaviors model the attachment (function_associations) today |
| `aws_cloudfront_key_group` | account-scoped signed-URL/signed-cookie key group referenced by many behaviors (the private-content mechanism, not edge-function surface); judged as a planned AwsCloudFrontKeyGroup kind together with its public keys -- AwsCloudFront behaviors model the attachment (trusted_key_group_ids) today |
| `aws_cloudfront_key_value_store` | judged as part of the planned AwsCloudFrontFunction kind (key-value stores and keys) |
| `aws_cloudfront_origin_request_policy` | account-scoped shared origin-request policy referenced by many behaviors across distributions; judged as a planned AwsCloudFrontOriginRequestPolicy kind -- AwsCloudFront behaviors model the attachment (origin_request_policy_id) today |
| `aws_cloudfront_public_key` | account-scoped public key consumed by key groups (signed URLs) and field-level encryption profiles; judged as part of the planned AwsCloudFrontKeyGroup kind |
| `aws_cloudfront_realtime_log_config` | account-scoped shared real-time log configuration (Kinesis-backed) referenced by many behaviors; judged as a planned AwsCloudFrontRealtimeLogConfig kind -- AwsCloudFront behaviors model the attachment (realtime_log_config_arn) today |
| `aws_cloudfront_response_headers_policy` | account-scoped shared response-headers policy referenced by many behaviors across distributions; judged as a planned AwsCloudFrontResponseHeadersPolicy kind -- AwsCloudFront behaviors model the attachment (response_headers_policy_id) today |
| `aws_cloudfront_trust_store` | account-scoped shared CA-bundle trust store referenced by many distributions' viewer-mTLS configurations; judged as a planned AwsCloudFrontTrustStore kind -- AwsCloudFront models the attachment (viewer_mtls.trust_store_id) today |
| `aws_cloudfront_vpc_origin` | account-scoped provisioned VPC origin shareable across distributions (and cross-account); judged as a planned AwsCloudFrontVpcOrigin kind -- AwsCloudFront models the attachment (origins[].vpc_origin.vpc_origin_id) today |
| `aws_cloudfrontkeyvaluestore_key` | judged as part of the planned AwsCloudFrontFunction kind (key-value stores and keys) |
| `aws_cloudfrontkeyvaluestore_keys_exclusive` | judged as part of the planned AwsCloudFrontFunction kind (key-value stores and keys) |
| `aws_cloudtrail` | judged as a planned AwsCloudTrail kind (trails, event data stores, organization delegated admin) |
| `aws_cloudtrail_event_data_store` | judged as a planned AwsCloudTrail kind (trails, event data stores, organization delegated admin) |
| `aws_cloudtrail_organization_delegated_admin_account` | judged as a planned AwsCloudTrail kind (trails, event data stores, organization delegated admin) |
| `aws_cloudwatch_alarm_mute_rule` | judged as a planned AwsCloudwatchAlarmMuteRule kind: one rule mutes up to 100 alarms on a recurring schedule -- multi-alarm scope, not per-alarm configuration, so a fold onto AwsCloudwatchAlarm would misrepresent it (the CloudFront-companion class) |
| `aws_cloudwatch_dashboard` | judged as a planned AwsCloudwatchDashboard kind |
| `aws_cloudwatch_event_api_destination` | judged as a planned AwsEventbridgeApiDestination kind (API destinations with connections) |
| `aws_cloudwatch_event_connection` | judged as a planned AwsEventbridgeApiDestination kind (API destinations with connections) |
| `aws_cloudwatch_log_account_policy` | judged as a planned AwsCloudwatchLogAccountPolicy kind: an account-wide policy per type (data protection / subscription filter / field index / transformer) optionally scoped by selection criteria -- standalone account-scoped configuration, never a per-group satellite |
| `aws_cloudwatch_log_anomaly_detector` | judged as a planned AwsCloudwatchLogAnomalyDetector kind: one detector trains over a LIST of log groups (log_group_arn_list) -- multi-parent scope, not a single group's satellite (the CloudFront-companion class) |
| `aws_cloudwatch_log_delivery` | judged as a planned AwsCloudWatchLogDelivery kind (vended-log deliveries, sources, destinations and their policies) |
| `aws_cloudwatch_log_delivery_destination` | judged as a planned AwsCloudWatchLogDelivery kind (vended-log deliveries, sources, destinations and their policies) |
| `aws_cloudwatch_log_delivery_destination_policy` | judged as a planned AwsCloudWatchLogDelivery kind (vended-log deliveries, sources, destinations and their policies) |
| `aws_cloudwatch_log_delivery_source` | judged as a planned AwsCloudWatchLogDelivery kind (vended-log deliveries, sources, destinations and their policies) |
| `aws_cloudwatch_log_destination` | judged as a planned AwsCloudWatchLogDelivery kind (vended-log deliveries, sources, destinations and their policies) |
| `aws_cloudwatch_log_destination_policy` | judged as a planned AwsCloudWatchLogDelivery kind (vended-log deliveries, sources, destinations and their policies) |
| `aws_cloudwatch_log_resource_policy` | judged as a planned AwsCloudwatchLogResourcePolicy kind: account-scoped named policies that grant AWS services log-delivery permission (Route53 query logging, EventBridge, ...) -- standalone glue many kinds need; nothing in the catalog renders it today, so composed would overstate coverage |
| `aws_cloudwatch_metric_stream` | judged as a planned AwsCloudwatchMetricStream kind |
| `aws_codeartifact_domain` | judged as a planned AwsCodeArtifact kind (domains and repositories with permission policies) |
| `aws_codeartifact_domain_permissions_policy` | judged as a planned AwsCodeArtifact kind (domains and repositories with permission policies) |
| `aws_codeartifact_repository` | judged as a planned AwsCodeArtifact kind (domains and repositories with permission policies) |
| `aws_codeartifact_repository_permissions_policy` | judged as a planned AwsCodeArtifact kind (domains and repositories with permission policies) |
| `aws_codebuild_fleet` | a future AwsCodeBuildFleet kind: reserved-capacity fleets are account-scoped shared compute pools with their own scaling/VPC/compute-configuration surface, referenced by many projects -- AwsCodeBuildProject models the attachment (environment.fleet_arn) and records the fleet itself as deliberately excluded |
| `aws_codebuild_source_credential` | a future AwsCodeBuildSourceCredential kind: the account/region-wide git credential store (one per server type, shared by every project; secret-bearing token) -- AwsCodeBuildProject records it as deliberately excluded and carries per-source auth via source.auth |
| `aws_codeconnections_connection` | judged as a planned AwsCodeConnection kind (connections and hosts to external SCM providers) |
| `aws_codeconnections_host` | judged as a planned AwsCodeConnection kind (connections and hosts to external SCM providers) |
| `aws_codedeploy_app` | judged as a planned AwsCodeDeployApp kind (applications with deployment groups and configs) |
| `aws_codedeploy_deployment_config` | judged as a planned AwsCodeDeployApp kind (applications with deployment groups and configs) |
| `aws_codedeploy_deployment_group` | judged as a planned AwsCodeDeployApp kind (applications with deployment groups and configs) |
| `aws_cognito_identity_pool` | judged as a planned AwsCognitoIdentityPool kind (identity pools with roles attachments and provider principal tags) |
| `aws_cognito_identity_pool_provider_principal_tag` | judged as a planned AwsCognitoIdentityPool kind (identity pools with roles attachments and provider principal tags) |
| `aws_cognito_identity_pool_roles_attachment` | judged as a planned AwsCognitoIdentityPool kind (identity pools with roles attachments and provider principal tags) |
| `aws_cognito_managed_login_branding` | judged as a planned AwsCognitoManagedLoginBranding kind: client-scoped managed-login design assets (smithy-JSON settings plus up to 40 binary asset payloads with whole-set replacement semantics) -- a distinct artifact lifecycle, not pool or client configuration |
| `aws_cognito_managed_user_pool_client` | judged as a planned AwsCognitoManagedUserPoolClient kind: ADOPTS an app client AWS created on the account's behalf (discovers by name pattern, updates in place, no-op delete) -- an adoption lifecycle distinct from AwsCognitoUserPoolClient's create-own lifecycle (the transit-gateway accepter class) |
| `aws_config_aggregate_authorization` | judged as a planned AwsConfigRecorder kind (recorders, delivery channels, retention, aggregators) |
| `aws_config_config_rule` | judged as a planned AwsConfigRule kind (managed/custom/organization rules, remediation, conformance packs) |
| `aws_config_configuration_aggregator` | judged as a planned AwsConfigRecorder kind (recorders, delivery channels, retention, aggregators) |
| `aws_config_configuration_recorder` | judged as a planned AwsConfigRecorder kind (recorders, delivery channels, retention, aggregators) |
| `aws_config_configuration_recorder_status` | judged as a planned AwsConfigRecorder kind (recorders, delivery channels, retention, aggregators) |
| `aws_config_conformance_pack` | judged as a planned AwsConfigRule kind (managed/custom/organization rules, remediation, conformance packs) |
| `aws_config_delivery_channel` | judged as a planned AwsConfigRecorder kind (recorders, delivery channels, retention, aggregators) |
| `aws_config_organization_conformance_pack` | judged as a planned AwsConfigRule kind (managed/custom/organization rules, remediation, conformance packs) |
| `aws_config_organization_custom_policy_rule` | judged as a planned AwsConfigRule kind (managed/custom/organization rules, remediation, conformance packs) |
| `aws_config_organization_custom_rule` | judged as a planned AwsConfigRule kind (managed/custom/organization rules, remediation, conformance packs) |
| `aws_config_organization_managed_rule` | judged as a planned AwsConfigRule kind (managed/custom/organization rules, remediation, conformance packs) |
| `aws_config_remediation_configuration` | judged as a planned AwsConfigRule kind (managed/custom/organization rules, remediation, conformance packs) |
| `aws_config_retention_configuration` | judged as a planned AwsConfigRecorder kind (recorders, delivery channels, retention, aggregators) |
| `aws_controltower_baseline` | judged as a planned AwsControlTower kind (controls, baselines, landing zones) |
| `aws_controltower_control` | judged as a planned AwsControlTower kind (controls, baselines, landing zones) |
| `aws_controltower_landing_zone` | judged as a planned AwsControlTower kind (controls, baselines, landing zones) |
| `aws_cur_report_definition` | judged as part of the planned cost-reporting kinds (with BCM data exports) |
| `aws_customer_gateway` | judged as a planned AwsSiteToSiteVpn kind (customer/VPN gateways, connections, routes, concentrators, route propagation) |
| `aws_datasync_agent` | judged as a planned AwsDataSyncLocation kind (agents and every location variant) |
| `aws_datasync_location_azure_blob` | judged as a planned AwsDataSyncLocation kind (agents and every location variant) |
| `aws_datasync_location_efs` | judged as a planned AwsDataSyncLocation kind (agents and every location variant) |
| `aws_datasync_location_fsx_lustre_file_system` | judged as a planned AwsDataSyncLocation kind (agents and every location variant) |
| `aws_datasync_location_fsx_ontap_file_system` | judged as a planned AwsDataSyncLocation kind (agents and every location variant) |
| `aws_datasync_location_fsx_openzfs_file_system` | judged as a planned AwsDataSyncLocation kind (agents and every location variant) |
| `aws_datasync_location_fsx_windows_file_system` | judged as a planned AwsDataSyncLocation kind (agents and every location variant) |
| `aws_datasync_location_hdfs` | judged as a planned AwsDataSyncLocation kind (agents and every location variant) |
| `aws_datasync_location_nfs` | judged as a planned AwsDataSyncLocation kind (agents and every location variant) |
| `aws_datasync_location_object_storage` | judged as a planned AwsDataSyncLocation kind (agents and every location variant) |
| `aws_datasync_location_s3` | judged as a planned AwsDataSyncLocation kind (agents and every location variant) |
| `aws_datasync_location_smb` | judged as a planned AwsDataSyncLocation kind (agents and every location variant) |
| `aws_datasync_task` | judged as a planned AwsDataSyncTask kind |
| `aws_db_event_subscription` | judged as a planned AwsRdsEventSubscription kind -- an event subscription binds an SNS topic to account-wide source sets (all instances, all clusters, or explicit lists) with its own lifecycle, never satellite to one database |
| `aws_db_instance_automated_backups_replication` | judged as a planned AwsRdsAutomatedBackupsReplication kind -- the replication is created in the DESTINATION region, a cross-region lifecycle outside the single-region instance kind (the replica-key class) |
| `aws_db_proxy` | judged as a planned AwsRdsProxy kind (proxies, endpoints, target groups, targets) |
| `aws_db_proxy_default_target_group` | judged as a planned AwsRdsProxy kind (proxies, endpoints, target groups, targets) |
| `aws_db_proxy_endpoint` | judged as a planned AwsRdsProxy kind (proxies, endpoints, target groups, targets) |
| `aws_db_proxy_target` | judged as a planned AwsRdsProxy kind (proxies, endpoints, target groups, targets) |
| `aws_directory_service_conditional_forwarder` | judged as a planned AwsManagedAd kind (directories with forwarders, log subscriptions, RADIUS, regions, sharing, trusts) |
| `aws_directory_service_directory` | judged as a planned AwsManagedAd kind (directories with forwarders, log subscriptions, RADIUS, regions, sharing, trusts) |
| `aws_directory_service_log_subscription` | judged as a planned AwsManagedAd kind (directories with forwarders, log subscriptions, RADIUS, regions, sharing, trusts) |
| `aws_directory_service_radius_settings` | judged as a planned AwsManagedAd kind (directories with forwarders, log subscriptions, RADIUS, regions, sharing, trusts) |
| `aws_directory_service_region` | judged as a planned AwsManagedAd kind (directories with forwarders, log subscriptions, RADIUS, regions, sharing, trusts) |
| `aws_directory_service_shared_directory` | judged as a planned AwsManagedAd kind (directories with forwarders, log subscriptions, RADIUS, regions, sharing, trusts) |
| `aws_directory_service_shared_directory_accepter` | judged as a planned AwsManagedAd kind (directories with forwarders, log subscriptions, RADIUS, regions, sharing, trusts) |
| `aws_directory_service_trust` | judged as a planned AwsManagedAd kind (directories with forwarders, log subscriptions, RADIUS, regions, sharing, trusts) |
| `aws_dlm_lifecycle_policy` | DLM lifecycle policies fold into the planned AwsEbsSnapshot kind |
| `aws_dms_certificate` | judged as a planned AwsDmsReplication kind family (replication instances/configs, endpoints, tasks, subnet groups, certificates, event subscriptions) |
| `aws_dms_endpoint` | judged as a planned AwsDmsReplication kind family (replication instances/configs, endpoints, tasks, subnet groups, certificates, event subscriptions) |
| `aws_dms_event_subscription` | judged as a planned AwsDmsReplication kind family (replication instances/configs, endpoints, tasks, subnet groups, certificates, event subscriptions) |
| `aws_dms_replication_config` | judged as a planned AwsDmsReplication kind family (replication instances/configs, endpoints, tasks, subnet groups, certificates, event subscriptions) |
| `aws_dms_replication_instance` | judged as a planned AwsDmsReplication kind family (replication instances/configs, endpoints, tasks, subnet groups, certificates, event subscriptions) |
| `aws_dms_replication_subnet_group` | judged as a planned AwsDmsReplication kind family (replication instances/configs, endpoints, tasks, subnet groups, certificates, event subscriptions) |
| `aws_dms_replication_task` | judged as a planned AwsDmsReplication kind family (replication instances/configs, endpoints, tasks, subnet groups, certificates, event subscriptions) |
| `aws_dms_s3_endpoint` | judged as a planned AwsDmsReplication kind family (replication instances/configs, endpoints, tasks, subnet groups, certificates, event subscriptions) |
| `aws_docdb_event_subscription` | judged as a planned AwsDocumentDbEventSubscription kind -- an event subscription binds an SNS topic to account-wide source sets with its own lifecycle, never satellite to one cluster (the AwsRdsEventSubscription precedent) |
| `aws_docdb_global_cluster` | judged as a planned AwsDocumentDbGlobalCluster kind -- the cross-region umbrella clusters join through spec.global_cluster_identifier (the AwsRdsGlobalCluster precedent) |
| `aws_docdbelastic_cluster` | judged as a planned AwsDocumentDbElasticCluster kind -- elastic clusters are a different architecture (sharded, shard-capacity-sized, no instances list), not the instance-based cluster this kind models |
| `aws_dsql_cluster` | judged as a planned AwsAuroraDsql kind (clusters with peering) |
| `aws_dsql_cluster_peering` | judged as a planned AwsAuroraDsql kind (clusters with peering) |
| `aws_dx_bgp_peer` | judged as a planned AwsDxVirtualInterface kind (private/public/transit and hosted virtual interfaces, accepters, BGP peers) |
| `aws_dx_connection` | judged as a planned AwsDxConnection kind (connections, hosted connections, associations, MACsec keys) |
| `aws_dx_connection_association` | judged as a planned AwsDxConnection kind (connections, hosted connections, associations, MACsec keys) |
| `aws_dx_connection_confirmation` | judged as a planned AwsDxConnection kind (connections, hosted connections, associations, MACsec keys) |
| `aws_dx_gateway` | judged as a planned AwsDxGateway kind (gateways with associations and proposals) |
| `aws_dx_gateway_association` | judged as a planned AwsDxGateway kind (gateways with associations and proposals) |
| `aws_dx_gateway_association_proposal` | judged as a planned AwsDxGateway kind (gateways with associations and proposals) |
| `aws_dx_hosted_connection` | judged as a planned AwsDxConnection kind (connections, hosted connections, associations, MACsec keys) |
| `aws_dx_hosted_private_virtual_interface` | judged as a planned AwsDxVirtualInterface kind (private/public/transit and hosted virtual interfaces, accepters, BGP peers) |
| `aws_dx_hosted_private_virtual_interface_accepter` | judged as a planned AwsDxVirtualInterface kind (private/public/transit and hosted virtual interfaces, accepters, BGP peers) |
| `aws_dx_hosted_public_virtual_interface` | judged as a planned AwsDxVirtualInterface kind (private/public/transit and hosted virtual interfaces, accepters, BGP peers) |
| `aws_dx_hosted_public_virtual_interface_accepter` | judged as a planned AwsDxVirtualInterface kind (private/public/transit and hosted virtual interfaces, accepters, BGP peers) |
| `aws_dx_hosted_transit_virtual_interface` | judged as a planned AwsDxVirtualInterface kind (private/public/transit and hosted virtual interfaces, accepters, BGP peers) |
| `aws_dx_hosted_transit_virtual_interface_accepter` | judged as a planned AwsDxVirtualInterface kind (private/public/transit and hosted virtual interfaces, accepters, BGP peers) |
| `aws_dx_lag` | judged as a planned AwsDxLag kind |
| `aws_dx_macsec_key_association` | judged as a planned AwsDxConnection kind (connections, hosted connections, associations, MACsec keys) |
| `aws_dx_private_virtual_interface` | judged as a planned AwsDxVirtualInterface kind (private/public/transit and hosted virtual interfaces, accepters, BGP peers) |
| `aws_dx_public_virtual_interface` | judged as a planned AwsDxVirtualInterface kind (private/public/transit and hosted virtual interfaces, accepters, BGP peers) |
| `aws_dx_transit_virtual_interface` | judged as a planned AwsDxVirtualInterface kind (private/public/transit and hosted virtual interfaces, accepters, BGP peers) |
| `aws_dynamodb_global_secondary_index` | table companion surface (GSIs, replicas); folds into the existing AwsDynamodb kind as its spec deepens |
| `aws_dynamodb_table_replica` | table companion surface (GSIs, replicas); folds into the existing AwsDynamodb kind as its spec deepens |
| `aws_ebs_default_kms_key` | account-wide EBS security toggles fold into the planned AwsEc2AccountSettings kind |
| `aws_ebs_encryption_by_default` | account-wide EBS security toggles fold into the planned AwsEc2AccountSettings kind |
| `aws_ebs_fast_snapshot_restore` | judged as a planned AwsEbsSnapshot kind (snapshots, copies, imports, fast restore, volume permissions) |
| `aws_ebs_snapshot` | judged as a planned AwsEbsSnapshot kind (snapshots, copies, imports, fast restore, volume permissions) |
| `aws_ebs_snapshot_block_public_access` | account-wide EBS security toggles fold into the planned AwsEc2AccountSettings kind |
| `aws_ebs_snapshot_copy` | judged as a planned AwsEbsSnapshot kind (snapshots, copies, imports, fast restore, volume permissions) |
| `aws_ebs_snapshot_import` | judged as a planned AwsEbsSnapshot kind (snapshots, copies, imports, fast restore, volume permissions) |
| `aws_ebs_volume` | judged as a planned AwsEbsVolume kind (volumes, copies, attachments) |
| `aws_ebs_volume_copy` | judged as a planned AwsEbsVolume kind (volumes, copies, attachments) |
| `aws_ec2_allowed_images_settings` | account-wide EC2 toggles fold into the planned AwsEc2AccountSettings kind |
| `aws_ec2_availability_zone_group` | account-wide EC2 toggles fold into the planned AwsEc2AccountSettings kind |
| `aws_ec2_capacity_block_reservation` | judged as a planned AwsEc2CapacityReservation kind |
| `aws_ec2_capacity_reservation` | judged as a planned AwsEc2CapacityReservation kind |
| `aws_ec2_default_credit_specification` | account-wide EC2 toggles fold into the planned AwsEc2AccountSettings kind |
| `aws_ec2_image_block_public_access` | account-wide EC2 toggles fold into the planned AwsEc2AccountSettings kind |
| `aws_ec2_instance_connect_endpoint` | judged as a planned AwsEc2InstanceConnectEndpoint kind (re-judged 2026-08-11 from the blanket fold-into-AwsEc2Instance reason): a subnet-resident endpoint with its own lifecycle serving MANY instances -- SSH/RDP reachability infrastructure, never per-instance configuration |
| `aws_ec2_instance_metadata_defaults` | account-wide EC2 toggles fold into the planned AwsEc2AccountSettings kind |
| `aws_ec2_managed_prefix_list` | judged as a planned AwsManagedPrefixList kind (lists with entries) |
| `aws_ec2_managed_prefix_list_entry` | judged as a planned AwsManagedPrefixList kind (lists with entries) |
| `aws_ec2_serial_console_access` | account-wide EC2 toggles fold into the planned AwsEc2AccountSettings kind |
| `aws_ec2_subnet_cidr_reservation` | subnet CIDR reservations fold into the existing AwsSubnet kind as its spec deepens |
| `aws_ec2_transit_gateway_connect` | Connect (GRE) attachment for SD-WAN appliance integration over a transport attachment -- a standalone attachment type with its own lifecycle, judged as a future AwsTransitGatewayConnect kind (never a fold: the existing kinds model the hub, its VPC spokes, and routing domains) |
| `aws_ec2_transit_gateway_connect_peer` | BGP peer inside a Connect attachment -- connect-scoped satellite with no independent lifecycle; folds into the future AwsTransitGatewayConnect kind as its peers list |
| `aws_ec2_transit_gateway_multicast_domain` | multicast domain on a Transit Gateway -- a standalone routing surface with its own lifecycle and sharing model, judged as a future AwsTransitGatewayMulticastDomain kind |
| `aws_ec2_transit_gateway_multicast_domain_association` | subnet/attachment membership of a multicast domain -- domain-scoped satellite; folds into the future AwsTransitGatewayMulticastDomain kind |
| `aws_ec2_transit_gateway_multicast_group_member` | ENI registration as a multicast group member -- domain-scoped satellite; folds into the future AwsTransitGatewayMulticastDomain kind |
| `aws_ec2_transit_gateway_multicast_group_source` | ENI registration as a multicast group source -- domain-scoped satellite; folds into the future AwsTransitGatewayMulticastDomain kind |
| `aws_ec2_transit_gateway_peering_attachment` | cross-region/cross-account hub-to-hub peering -- a standalone attachment type whose requester/accepter halves live in different deployment scopes, judged as a future AwsTransitGatewayPeeringAttachment kind |
| `aws_ec2_transit_gateway_peering_attachment_accepter` | the accepter-side half of transit-gateway peering; designed together with the future AwsTransitGatewayPeeringAttachment kind |
| `aws_ec2_transit_gateway_policy_table` | dynamic-routing policy table (used with Connect/Cloud WAN peering) -- a standalone table with its own lifecycle, judged as a future AwsTransitGatewayPolicyTable kind |
| `aws_ec2_transit_gateway_policy_table_association` | attachment membership of a policy table -- table-scoped satellite; folds into the future AwsTransitGatewayPolicyTable kind |
| `aws_ec2_transit_gateway_vpc_attachment_accepter` | accepter-side half of a cross-account VPC attachment on a RAM-shared gateway -- runs in the sharing account's deployment scope, judged as a future AwsTransitGatewayVpcAttachmentAccepter kind alongside RAM-share modeling |
| `aws_ecr_account_setting` | judged as a planned AwsEcrRegistrySettings kind (registry policy, scanning, replication, pull-through caches, creation templates, account settings) |
| `aws_ecr_pull_through_cache_rule` | judged as a planned AwsEcrRegistrySettings kind (registry policy, scanning, replication, pull-through caches, creation templates, account settings) |
| `aws_ecr_pull_time_update_exclusion` | judged as a planned AwsEcrRegistrySettings kind (registry policy, scanning, replication, pull-through caches, creation templates, account settings) |
| `aws_ecr_registry_policy` | judged as a planned AwsEcrRegistrySettings kind (registry policy, scanning, replication, pull-through caches, creation templates, account settings) |
| `aws_ecr_registry_scanning_configuration` | judged as a planned AwsEcrRegistrySettings kind (registry policy, scanning, replication, pull-through caches, creation templates, account settings) |
| `aws_ecr_replication_configuration` | judged as a planned AwsEcrRegistrySettings kind (registry policy, scanning, replication, pull-through caches, creation templates, account settings) |
| `aws_ecr_repository_creation_template` | judged as a planned AwsEcrRegistrySettings kind (registry policy, scanning, replication, pull-through caches, creation templates, account settings) |
| `aws_ecs_account_setting_default` | account-level ECS setting defaults (a per-account singleton, not per-cluster surface); planned as its own account-posture admission, not a fold into AwsEcsCluster |
| `aws_ecs_daemon` | the ECS daemon workload paradigm (provider 6.50); planned alongside aws_ecs_daemon_task_definition as its own admission when daemon demand lands -- DAEMON-strategy services remain covered by AwsEcsService |
| `aws_ecs_daemon_task_definition` | the ECS daemon workload paradigm (provider 6.50); planned alongside aws_ecs_daemon as its own admission when daemon demand lands |
| `aws_ecs_express_gateway_service` | ECS Express -- a distinct gateway-fronted service paradigm with its own lifecycle; planned as its own kind admission, not a fold into AwsEcsService |
| `aws_ecs_task_set` | task sets exist only under the EXTERNAL deployment controller; folds into AwsEcsService if the external-controller workflow is ever admitted (its native and blue/green controllers are fully modeled) |
| `aws_eip_association` | EIP companion surface (associations, domain names); folds into the existing AwsElasticIp kind as its spec deepens |
| `aws_eip_domain_name` | EIP companion surface (associations, domain names); folds into the existing AwsElasticIp kind as its spec deepens |
| `aws_eks_capability` | EKS managed capabilities (ACK, KRO, Argo CD; provider 6.25.0) with their own configuration tree (IAM Identity Center wiring, RBAC role mappings, network access); a cluster-scoped companion with independent lifecycle -- its own kind (AwsEksCapability) on the EKS family's composition pattern (addon/fargate-profile/access-entry precedent), not a cluster fold-in |
| `aws_eks_identity_provider_config` | associates an external OIDC identity provider with a cluster for user authentication; an associable cluster companion with independent (create-only) lifecycle -- its own kind (AwsEksIdentityProviderConfig) on the EKS family's composition pattern, not a cluster fold-in |
| `aws_eks_pod_identity_association` | maps any Kubernetes service account to an IAM role via EKS Pod Identity; the per-ADDON arm is already modeled on AwsEksAddon (pod_identity_associations), and the standalone workload-scoped association is its own kind (AwsEksPodIdentityAssociation) -- per-workload IAM identity is chart-wiring surface, not cluster or addon configuration |
| `aws_elasticache_global_replication_group` | cross-region global datastore with its own lifecycle (suffix-named, anchored to a primary replication group, upgraded independently) -- a future AwsElasticacheGlobalReplicationGroup kind; AwsRedisElasticache already models the secondary-side join path (global_replication_group_id) |
| `aws_emr_block_public_access_configuration` | judged as a planned AwsEmrCluster kind (clusters, instance groups/fleets, managed scaling, security configs, block-public-access) |
| `aws_emr_cluster` | judged as a planned AwsEmrCluster kind (clusters, instance groups/fleets, managed scaling, security configs, block-public-access) |
| `aws_emr_instance_fleet` | judged as a planned AwsEmrCluster kind (clusters, instance groups/fleets, managed scaling, security configs, block-public-access) |
| `aws_emr_instance_group` | judged as a planned AwsEmrCluster kind (clusters, instance groups/fleets, managed scaling, security configs, block-public-access) |
| `aws_emr_managed_scaling_policy` | judged as a planned AwsEmrCluster kind (clusters, instance groups/fleets, managed scaling, security configs, block-public-access) |
| `aws_emr_security_configuration` | judged as a planned AwsEmrCluster kind (clusters, instance groups/fleets, managed scaling, security configs, block-public-access) |
| `aws_emr_studio` | judged as a planned AwsEmrStudio kind (studios with session mappings) |
| `aws_emr_studio_session_mapping` | judged as a planned AwsEmrStudio kind (studios with session mappings) |
| `aws_emrcontainers_job_template` | EMR on EKS folds into the planned AwsEmrCluster kind family |
| `aws_emrcontainers_virtual_cluster` | EMR on EKS folds into the planned AwsEmrCluster kind family |
| `aws_emrserverless_application` | EMR Serverless folds into the planned AwsEmrCluster kind family |
| `aws_flow_log` | VPC companion surface (flow logs, IPv6 associations, block-public-access, encryption control, route management); folds into the existing AwsVpc kind as its spec deepens |
| `aws_fsx_openzfs_volume` | judged as a planned AwsFsxOpenzfsVolume kind (child-volume hierarchy, NFS exports, quotas, clone-from-snapshot origin) -- standalone lifecycle on the AwsFsxOntapVolume precedent; the AwsFsxOpenzfsFileSystem spec records child volumes as out of its scope |
| `aws_fsx_s3_access_point_attachment` | attaches an S3 access point to an OpenZFS volume -- satellite of the planned AwsFsxOpenzfsVolume kind; judged with that kind's family |
| `aws_globalaccelerator_cross_account_attachment` | cross-account attachments fold into the existing AwsGlobalAccelerator kind as its spec deepens |
| `aws_glue_catalog` | Glue catalog-level surface (policies, encryption, security configs); folds into the existing AwsGlueCatalogDatabase kind as its spec deepens |
| `aws_glue_catalog_table` | judged as a planned AwsGlueCatalogTable kind (tables, partitions, indexes, optimizers) |
| `aws_glue_catalog_table_optimizer` | judged as a planned AwsGlueCatalogTable kind (tables, partitions, indexes, optimizers) |
| `aws_glue_classifier` | judged as a planned AwsGlueCrawler kind (crawlers with classifiers) |
| `aws_glue_connection` | judged as a planned AwsGlueConnection kind |
| `aws_glue_crawler` | judged as a planned AwsGlueCrawler kind (crawlers with classifiers) |
| `aws_glue_data_catalog_encryption_settings` | Glue catalog-level surface (policies, encryption, security configs); folds into the existing AwsGlueCatalogDatabase kind as its spec deepens |
| `aws_glue_job` | judged as a planned AwsGlueJob kind (jobs, triggers, workflows) |
| `aws_glue_partition` | judged as a planned AwsGlueCatalogTable kind (tables, partitions, indexes, optimizers) |
| `aws_glue_partition_index` | judged as a planned AwsGlueCatalogTable kind (tables, partitions, indexes, optimizers) |
| `aws_glue_registry` | judged as a planned AwsGlueSchemaRegistry kind (registries with schemas) |
| `aws_glue_resource_policy` | Glue catalog-level surface (policies, encryption, security configs); folds into the existing AwsGlueCatalogDatabase kind as its spec deepens |
| `aws_glue_schema` | judged as a planned AwsGlueSchemaRegistry kind (registries with schemas) |
| `aws_glue_security_configuration` | Glue catalog-level surface (policies, encryption, security configs); folds into the existing AwsGlueCatalogDatabase kind as its spec deepens |
| `aws_glue_trigger` | judged as a planned AwsGlueJob kind (jobs, triggers, workflows) |
| `aws_glue_workflow` | judged as a planned AwsGlueJob kind (jobs, triggers, workflows) |
| `aws_grafana_license_association` | judged as a planned AwsManagedGrafana kind (workspaces with licenses, SAML, API keys, service accounts) |
| `aws_grafana_role_association` | judged as a planned AwsManagedGrafana kind (workspaces with licenses, SAML, API keys, service accounts) |
| `aws_grafana_workspace` | judged as a planned AwsManagedGrafana kind (workspaces with licenses, SAML, API keys, service accounts) |
| `aws_grafana_workspace_api_key` | judged as a planned AwsManagedGrafana kind (workspaces with licenses, SAML, API keys, service accounts) |
| `aws_grafana_workspace_saml_configuration` | judged as a planned AwsManagedGrafana kind (workspaces with licenses, SAML, API keys, service accounts) |
| `aws_grafana_workspace_service_account` | judged as a planned AwsManagedGrafana kind (workspaces with licenses, SAML, API keys, service accounts) |
| `aws_grafana_workspace_service_account_token` | judged as a planned AwsManagedGrafana kind (workspaces with licenses, SAML, API keys, service accounts) |
| `aws_guardduty_detector` | judged as a planned AwsGuardDuty kind (detectors, features, filters, IP/threat-intel sets, malware protection, publishing destinations, organization administration) |
| `aws_guardduty_detector_feature` | judged as a planned AwsGuardDuty kind (detectors, features, filters, IP/threat-intel sets, malware protection, publishing destinations, organization administration) |
| `aws_guardduty_filter` | judged as a planned AwsGuardDuty kind (detectors, features, filters, IP/threat-intel sets, malware protection, publishing destinations, organization administration) |
| `aws_guardduty_invite_accepter` | judged as a planned AwsGuardDuty kind (detectors, features, filters, IP/threat-intel sets, malware protection, publishing destinations, organization administration) |
| `aws_guardduty_ipset` | judged as a planned AwsGuardDuty kind (detectors, features, filters, IP/threat-intel sets, malware protection, publishing destinations, organization administration) |
| `aws_guardduty_malware_protection_plan` | judged as a planned AwsGuardDuty kind (detectors, features, filters, IP/threat-intel sets, malware protection, publishing destinations, organization administration) |
| `aws_guardduty_member` | judged as a planned AwsGuardDuty kind (detectors, features, filters, IP/threat-intel sets, malware protection, publishing destinations, organization administration) |
| `aws_guardduty_member_detector_feature` | judged as a planned AwsGuardDuty kind (detectors, features, filters, IP/threat-intel sets, malware protection, publishing destinations, organization administration) |
| `aws_guardduty_organization_admin_account` | judged as a planned AwsGuardDuty kind (detectors, features, filters, IP/threat-intel sets, malware protection, publishing destinations, organization administration) |
| `aws_guardduty_organization_configuration` | judged as a planned AwsGuardDuty kind (detectors, features, filters, IP/threat-intel sets, malware protection, publishing destinations, organization administration) |
| `aws_guardduty_organization_configuration_feature` | judged as a planned AwsGuardDuty kind (detectors, features, filters, IP/threat-intel sets, malware protection, publishing destinations, organization administration) |
| `aws_guardduty_publishing_destination` | judged as a planned AwsGuardDuty kind (detectors, features, filters, IP/threat-intel sets, malware protection, publishing destinations, organization administration) |
| `aws_guardduty_threatintelset` | judged as a planned AwsGuardDuty kind (detectors, features, filters, IP/threat-intel sets, malware protection, publishing destinations, organization administration) |
| `aws_iam_account_alias` | judged as a planned AwsIamAccountSettings kind (alias, password policy, STS preferences, organization features) |
| `aws_iam_account_password_policy` | judged as a planned AwsIamAccountSettings kind (alias, password policy, STS preferences, organization features) |
| `aws_iam_group` | judged as a planned AwsIamGroup kind (groups, memberships, group policies and attachments) |
| `aws_iam_group_membership` | judged as a planned AwsIamGroup kind (groups, memberships, group policies and attachments) |
| `aws_iam_group_policies_exclusive` | judged as a planned AwsIamGroup kind (groups, memberships, group policies and attachments) |
| `aws_iam_group_policy` | judged as a planned AwsIamGroup kind (groups, memberships, group policies and attachments) |
| `aws_iam_group_policy_attachment` | judged as a planned AwsIamGroup kind (groups, memberships, group policies and attachments) |
| `aws_iam_group_policy_attachments_exclusive` | judged as a planned AwsIamGroup kind (groups, memberships, group policies and attachments) |
| `aws_iam_organizations_features` | judged as a planned AwsIamAccountSettings kind (alias, password policy, STS preferences, organization features) |
| `aws_iam_role_policies_exclusive` | role companion surface (exclusive policy lockdown, service-linked roles); folds into the existing AwsIamRole kind as its spec deepens |
| `aws_iam_role_policy_attachments_exclusive` | role companion surface (exclusive policy lockdown, service-linked roles); folds into the existing AwsIamRole kind as its spec deepens |
| `aws_iam_saml_provider` | judged as a planned AwsIamSamlProvider kind |
| `aws_iam_security_token_service_preferences` | judged as a planned AwsIamAccountSettings kind (alias, password policy, STS preferences, organization features) |
| `aws_iam_service_linked_role` | role companion surface (exclusive policy lockdown, service-linked roles); folds into the existing AwsIamRole kind as its spec deepens |
| `aws_iam_service_specific_credential` | user-scoped service credentials (Bedrock API keys, Cassandra, CodeCommit); folds into AwsIamUser with the Bedrock wave, which brings the demand signal (API keys), the per-credential secret-delivery design, and the live fixtures its proof needs |
| `aws_iam_user_group_membership` | judged as a planned AwsIamGroup kind (groups, memberships, group policies and attachments) |
| `aws_identitystore_group` | Identity Center directory objects (users, groups, memberships) fold into the planned AwsIdentityCenterAssignment kind family |
| `aws_identitystore_group_membership` | Identity Center directory objects (users, groups, memberships) fold into the planned AwsIdentityCenterAssignment kind family |
| `aws_identitystore_user` | Identity Center directory objects (users, groups, memberships) fold into the planned AwsIdentityCenterAssignment kind family |
| `aws_imagebuilder_component` | judged as a planned AwsImagePipeline kind (pipelines, recipes, components, infrastructure/distribution configs, lifecycle policies, workflows) |
| `aws_imagebuilder_container_recipe` | judged as a planned AwsImagePipeline kind (pipelines, recipes, components, infrastructure/distribution configs, lifecycle policies, workflows) |
| `aws_imagebuilder_distribution_configuration` | judged as a planned AwsImagePipeline kind (pipelines, recipes, components, infrastructure/distribution configs, lifecycle policies, workflows) |
| `aws_imagebuilder_image` | judged as a planned AwsImagePipeline kind (pipelines, recipes, components, infrastructure/distribution configs, lifecycle policies, workflows) |
| `aws_imagebuilder_image_pipeline` | judged as a planned AwsImagePipeline kind (pipelines, recipes, components, infrastructure/distribution configs, lifecycle policies, workflows) |
| `aws_imagebuilder_image_recipe` | judged as a planned AwsImagePipeline kind (pipelines, recipes, components, infrastructure/distribution configs, lifecycle policies, workflows) |
| `aws_imagebuilder_infrastructure_configuration` | judged as a planned AwsImagePipeline kind (pipelines, recipes, components, infrastructure/distribution configs, lifecycle policies, workflows) |
| `aws_imagebuilder_lifecycle_policy` | judged as a planned AwsImagePipeline kind (pipelines, recipes, components, infrastructure/distribution configs, lifecycle policies, workflows) |
| `aws_imagebuilder_workflow` | judged as a planned AwsImagePipeline kind (pipelines, recipes, components, infrastructure/distribution configs, lifecycle policies, workflows) |
| `aws_inspector2_delegated_admin_account` | judged as a planned AwsInspector kind (enabler, delegated admin, filters, organization configuration) |
| `aws_inspector2_enabler` | judged as a planned AwsInspector kind (enabler, delegated admin, filters, organization configuration) |
| `aws_inspector2_filter` | judged as a planned AwsInspector kind (enabler, delegated admin, filters, organization configuration) |
| `aws_inspector2_member_association` | judged as a planned AwsInspector kind (enabler, delegated admin, filters, organization configuration) |
| `aws_inspector2_organization_configuration` | judged as a planned AwsInspector kind (enabler, delegated admin, filters, organization configuration) |
| `aws_internet_gateway_attachment` | gateway attachments fold into the existing AwsInternetGateway kind as its spec deepens |
| `aws_key_pair` | judged as a planned AwsKeyPair kind (re-judged 2026-08-11 from the blanket fold-into-AwsEc2Instance reason): account-level public-key material referenced by many instances and launch templates -- AwsEc2Instance.key_name consumes a key pair by name but cannot create the material |
| `aws_keyspaces_keyspace` | judged as a planned AwsKeyspaces kind (keyspaces with tables) |
| `aws_keyspaces_table` | judged as a planned AwsKeyspaces kind (keyspaces with tables) |
| `aws_kinesis_account_settings` | account-level Kinesis settings fold into the existing AwsKinesisStream kind family's account composition |
| `aws_kinesisanalyticsv2_application` | judged as a planned AwsManagedFlink kind (applications with snapshots) |
| `aws_kinesisanalyticsv2_application_snapshot` | judged as a planned AwsManagedFlink kind (applications with snapshots) |
| `aws_kms_external_key` | judged as a planned AwsKmsExternalKey kind (BYOK imported key material -- its own creation lifecycle with sensitive-material handling, per the provider's own resource split) |
| `aws_kms_replica_external_key` | judged as a planned AwsKmsReplicaKey kind (the external-material arm of replica keys) |
| `aws_kms_replica_key` | judged as a planned AwsKmsReplicaKey kind (a replica is its own regional resource referencing a primary key ARN -- a cross-region lifecycle outside the single-region key kind) |
| `aws_lakeformation_data_cells_filter` | judged as a planned AwsLakeFormation kind (data-lake settings, permissions, resources, LF-tags, opt-ins, Identity Center config) |
| `aws_lakeformation_data_lake_settings` | judged as a planned AwsLakeFormation kind (data-lake settings, permissions, resources, LF-tags, opt-ins, Identity Center config) |
| `aws_lakeformation_identity_center_configuration` | judged as a planned AwsLakeFormation kind (data-lake settings, permissions, resources, LF-tags, opt-ins, Identity Center config) |
| `aws_lakeformation_lf_tag` | judged as a planned AwsLakeFormation kind (data-lake settings, permissions, resources, LF-tags, opt-ins, Identity Center config) |
| `aws_lakeformation_lf_tag_expression` | judged as a planned AwsLakeFormation kind (data-lake settings, permissions, resources, LF-tags, opt-ins, Identity Center config) |
| `aws_lakeformation_opt_in` | judged as a planned AwsLakeFormation kind (data-lake settings, permissions, resources, LF-tags, opt-ins, Identity Center config) |
| `aws_lakeformation_permissions` | judged as a planned AwsLakeFormation kind (data-lake settings, permissions, resources, LF-tags, opt-ins, Identity Center config) |
| `aws_lakeformation_resource` | judged as a planned AwsLakeFormation kind (data-lake settings, permissions, resources, LF-tags, opt-ins, Identity Center config) |
| `aws_lakeformation_resource_lf_tag` | judged as a planned AwsLakeFormation kind (data-lake settings, permissions, resources, LF-tags, opt-ins, Identity Center config) |
| `aws_lakeformation_resource_lf_tags` | judged as a planned AwsLakeFormation kind (data-lake settings, permissions, resources, LF-tags, opt-ins, Identity Center config) |
| `aws_lambda_capacity_provider` | judged as a planned AwsLambdaCapacityProvider kind: account-scoped shared infrastructure (own VPC config, scaling policies, KMS key, tags, 30-minute delete) that many functions reference by ARN -- AwsLambda models the function-side attachment (managed_instances.capacity_provider_arn) and upgrades it to a reference when the kind ships |
| `aws_lambda_code_signing_config` | judged as a planned AwsLambdaCodeSigningConfig kind: a shareable trust policy (allowed signing profiles + enforcement mode) many functions reference -- AwsLambda already models the function-side attachment (code_signing_config_arn) and upgrades it to a reference when the kind ships |
| `aws_lambda_layer_version` | judged as a planned AwsLambdaLayer kind (layer versions with permissions) |
| `aws_lambda_layer_version_permission` | judged as a planned AwsLambdaLayer kind (layer versions with permissions) |
| `aws_lb_trust_store` | standalone mTLS CA-bundle store with its own ARN, name, tags, revocation sub-resources, and many-listener sharing -- a future AwsLbTrustStore kind (the listener's trust_store_arn reference gains its default kind then); never a fold into AwsLbListener, whose spec records trust stores as deliberately not modeled |
| `aws_lb_trust_store_revocation` | certificate revocation lists attached to a trust store -- folds into the future AwsLbTrustStore kind as its revocations surface (pure sub-resource of the store, no cross-store identity) |
| `aws_macie2_account` | judged as a planned AwsMacie kind (account enablement, classification jobs, custom identifiers, findings filters, organization administration) |
| `aws_macie2_classification_export_configuration` | judged as a planned AwsMacie kind (account enablement, classification jobs, custom identifiers, findings filters, organization administration) |
| `aws_macie2_classification_job` | judged as a planned AwsMacie kind (account enablement, classification jobs, custom identifiers, findings filters, organization administration) |
| `aws_macie2_custom_data_identifier` | judged as a planned AwsMacie kind (account enablement, classification jobs, custom identifiers, findings filters, organization administration) |
| `aws_macie2_findings_filter` | judged as a planned AwsMacie kind (account enablement, classification jobs, custom identifiers, findings filters, organization administration) |
| `aws_macie2_invitation_accepter` | judged as a planned AwsMacie kind (account enablement, classification jobs, custom identifiers, findings filters, organization administration) |
| `aws_macie2_member` | judged as a planned AwsMacie kind (account enablement, classification jobs, custom identifiers, findings filters, organization administration) |
| `aws_macie2_organization_admin_account` | judged as a planned AwsMacie kind (account enablement, classification jobs, custom identifiers, findings filters, organization administration) |
| `aws_macie2_organization_configuration` | judged as a planned AwsMacie kind (account enablement, classification jobs, custom identifiers, findings filters, organization administration) |
| `aws_main_route_table_association` | VPC companion surface (flow logs, IPv6 associations, block-public-access, encryption control, route management); folds into the existing AwsVpc kind as its spec deepens |
| `aws_memorydb_multi_region_cluster` | multi-region clusters fold into the existing AwsMemorydbCluster kind as its spec deepens |
| `aws_mq_broker` | judged as a planned AwsMqBroker kind (brokers with configurations) |
| `aws_mq_configuration` | judged as a planned AwsMqBroker kind (brokers with configurations) |
| `aws_msk_topic` | MSK companion surface (topics, VPC connections); folds into the existing AwsMskCluster kind as its spec deepens |
| `aws_msk_vpc_connection` | MSK companion surface (topics, VPC connections); folds into the existing AwsMskCluster kind as its spec deepens |
| `aws_mskconnect_connector` | judged as a planned AwsMskConnect kind (connectors, custom plugins, worker configurations) |
| `aws_mskconnect_custom_plugin` | judged as a planned AwsMskConnect kind (connectors, custom plugins, worker configurations) |
| `aws_mskconnect_worker_configuration` | judged as a planned AwsMskConnect kind (connectors, custom plugins, worker configurations) |
| `aws_neptune_cluster_endpoint` | Neptune companion surface (endpoints, event subscriptions, parameter groups, global clusters); folds into the existing AwsNeptuneCluster kind as its spec deepens |
| `aws_neptune_event_subscription` | Neptune companion surface (endpoints, event subscriptions, parameter groups, global clusters); folds into the existing AwsNeptuneCluster kind as its spec deepens |
| `aws_neptune_global_cluster` | Neptune companion surface (endpoints, event subscriptions, parameter groups, global clusters); folds into the existing AwsNeptuneCluster kind as its spec deepens |
| `aws_neptune_parameter_group` | Neptune companion surface (endpoints, event subscriptions, parameter groups, global clusters); folds into the existing AwsNeptuneCluster kind as its spec deepens |
| `aws_network_acl` | judged as a planned AwsNetworkAcl kind (ACLs, associations, rules) |
| `aws_network_acl_association` | judged as a planned AwsNetworkAcl kind (ACLs, associations, rules) |
| `aws_network_acl_rule` | judged as a planned AwsNetworkAcl kind (ACLs, associations, rules) |
| `aws_network_interface` | judged as a planned AwsNetworkInterface kind (re-judged 2026-08-11 from the blanket fold-into-AwsEc2Instance reason): a standalone ENI lifecycle attachable across instances -- AwsEc2Instance consumes one by reference (primary_network_interface_id) and creates launch-scoped ones inline (secondary_network_interfaces), neither of which is the durable standalone ENI; the provider's own aws_instance.network_interface deprecation points at the attachment resource as the sanctioned companion |
| `aws_network_interface_attachment` | attachment satellite of the planned AwsNetworkInterface kind (re-judged 2026-08-11): binds an existing ENI to an instance at a device index -- the provider's sanctioned replacement for the deprecated aws_instance.network_interface block |
| `aws_network_interface_permission` | cross-account permission satellite of the planned AwsNetworkInterface kind (re-judged 2026-08-11): grants another account attach rights on an ENI |
| `aws_network_interface_sg_attachment` | security-group satellite of the planned AwsNetworkInterface kind (re-judged 2026-08-11): attaches a security group to an existing ENI without owning either side |
| `aws_networkfirewall_firewall` | judged as a planned AwsNetworkFirewall kind family (firewalls, policies, rule groups, logging, TLS inspection) |
| `aws_networkfirewall_firewall_policy` | judged as a planned AwsNetworkFirewall kind family (firewalls, policies, rule groups, logging, TLS inspection) |
| `aws_networkfirewall_firewall_transit_gateway_attachment_accepter` | judged as a planned AwsNetworkFirewall kind family (firewalls, policies, rule groups, logging, TLS inspection) |
| `aws_networkfirewall_logging_configuration` | judged as a planned AwsNetworkFirewall kind family (firewalls, policies, rule groups, logging, TLS inspection) |
| `aws_networkfirewall_resource_policy` | judged as a planned AwsNetworkFirewall kind family (firewalls, policies, rule groups, logging, TLS inspection) |
| `aws_networkfirewall_rule_group` | judged as a planned AwsNetworkFirewall kind family (firewalls, policies, rule groups, logging, TLS inspection) |
| `aws_networkfirewall_tls_inspection_configuration` | judged as a planned AwsNetworkFirewall kind family (firewalls, policies, rule groups, logging, TLS inspection) |
| `aws_networkfirewall_vpc_endpoint_association` | judged as a planned AwsNetworkFirewall kind family (firewalls, policies, rule groups, logging, TLS inspection) |
| `aws_oam_link` | judged as a planned AwsObservabilityLink kind (CloudWatch cross-account sinks and links) |
| `aws_oam_sink` | judged as a planned AwsObservabilityLink kind (CloudWatch cross-account sinks and links) |
| `aws_oam_sink_policy` | judged as a planned AwsObservabilityLink kind (CloudWatch cross-account sinks and links) |
| `aws_opensearch_authorize_vpc_endpoint_access` | domain companion surface (policies, SAML, VPC endpoints, packages); folds into the existing AwsOpenSearchDomain kind as its spec deepens |
| `aws_opensearch_domain_policy` | domain companion surface (policies, SAML, VPC endpoints, packages); folds into the existing AwsOpenSearchDomain kind as its spec deepens |
| `aws_opensearch_domain_saml_options` | domain companion surface (policies, SAML, VPC endpoints, packages); folds into the existing AwsOpenSearchDomain kind as its spec deepens |
| `aws_opensearch_package` | domain companion surface (policies, SAML, VPC endpoints, packages); folds into the existing AwsOpenSearchDomain kind as its spec deepens |
| `aws_opensearch_package_association` | domain companion surface (policies, SAML, VPC endpoints, packages); folds into the existing AwsOpenSearchDomain kind as its spec deepens |
| `aws_opensearch_vpc_endpoint` | domain companion surface (policies, SAML, VPC endpoints, packages); folds into the existing AwsOpenSearchDomain kind as its spec deepens |
| `aws_opensearchserverless_access_policy` | judged as a planned AwsOpenSearchServerlessCollection kind (collections with security/access/lifecycle policies and VPC endpoints) |
| `aws_opensearchserverless_collection` | judged as a planned AwsOpenSearchServerlessCollection kind (collections with security/access/lifecycle policies and VPC endpoints) |
| `aws_opensearchserverless_collection_group` | judged as a planned AwsOpenSearchServerlessCollection kind (collections with security/access/lifecycle policies and VPC endpoints) |
| `aws_opensearchserverless_lifecycle_policy` | judged as a planned AwsOpenSearchServerlessCollection kind (collections with security/access/lifecycle policies and VPC endpoints) |
| `aws_opensearchserverless_security_config` | judged as a planned AwsOpenSearchServerlessCollection kind (collections with security/access/lifecycle policies and VPC endpoints) |
| `aws_opensearchserverless_security_policy` | judged as a planned AwsOpenSearchServerlessCollection kind (collections with security/access/lifecycle policies and VPC endpoints) |
| `aws_opensearchserverless_vpc_endpoint` | judged as a planned AwsOpenSearchServerlessCollection kind (collections with security/access/lifecycle policies and VPC endpoints) |
| `aws_organizations_account` | judged as a planned AwsOrganizationAccount kind |
| `aws_organizations_aws_service_access` | judged as a planned AwsOrganization kind (organization, service access, delegated administrators, resource policies) |
| `aws_organizations_delegated_administrator` | judged as a planned AwsOrganization kind (organization, service access, delegated administrators, resource policies) |
| `aws_organizations_organization` | judged as a planned AwsOrganization kind (organization, service access, delegated administrators, resource policies) |
| `aws_organizations_organizational_unit` | judged as a planned AwsOrganizationalUnit kind |
| `aws_organizations_policy` | judged as a planned AwsOrganizationPolicy kind (policies with attachments) |
| `aws_organizations_policy_attachment` | judged as a planned AwsOrganizationPolicy kind (policies with attachments) |
| `aws_organizations_resource_policy` | judged as a planned AwsOrganization kind (organization, service access, delegated administrators, resource policies) |
| `aws_organizations_tag` | judged as a planned AwsOrganization kind (organization, service access, delegated administrators, resource policies) |
| `aws_pipes_pipe` | judged as a planned AwsEventbridgePipe kind |
| `aws_placement_group` | judged as a planned AwsPlacementGroup kind (re-judged 2026-08-11 from the blanket fold-into-AwsEc2Instance reason): an account-level placement strategy (cluster/partition/spread) consumed by many instances and ASGs -- AwsEc2Instance.placement references a group by name/id but cannot create one |
| `aws_prometheus_alert_manager_definition` | judged as a planned AwsManagedPrometheus kind (workspaces with alert-manager/rule definitions, scrapers, anomaly detectors, logging) |
| `aws_prometheus_anomaly_detector` | judged as a planned AwsManagedPrometheus kind (workspaces with alert-manager/rule definitions, scrapers, anomaly detectors, logging) |
| `aws_prometheus_query_logging_configuration` | judged as a planned AwsManagedPrometheus kind (workspaces with alert-manager/rule definitions, scrapers, anomaly detectors, logging) |
| `aws_prometheus_resource_policy` | judged as a planned AwsManagedPrometheus kind (workspaces with alert-manager/rule definitions, scrapers, anomaly detectors, logging) |
| `aws_prometheus_rule_group_namespace` | judged as a planned AwsManagedPrometheus kind (workspaces with alert-manager/rule definitions, scrapers, anomaly detectors, logging) |
| `aws_prometheus_scraper` | judged as a planned AwsManagedPrometheus kind (workspaces with alert-manager/rule definitions, scrapers, anomaly detectors, logging) |
| `aws_prometheus_scraper_logging_configuration` | judged as a planned AwsManagedPrometheus kind (workspaces with alert-manager/rule definitions, scrapers, anomaly detectors, logging) |
| `aws_prometheus_workspace` | judged as a planned AwsManagedPrometheus kind (workspaces with alert-manager/rule definitions, scrapers, anomaly detectors, logging) |
| `aws_prometheus_workspace_configuration` | judged as a planned AwsManagedPrometheus kind (workspaces with alert-manager/rule definitions, scrapers, anomaly detectors, logging) |
| `aws_ram_permission` | judged as a planned AwsResourceShare kind (shares with principal/resource associations, accepters, permissions) |
| `aws_ram_principal_association` | judged as a planned AwsResourceShare kind (shares with principal/resource associations, accepters, permissions) |
| `aws_ram_resource_association` | judged as a planned AwsResourceShare kind (shares with principal/resource associations, accepters, permissions) |
| `aws_ram_resource_share` | judged as a planned AwsResourceShare kind (shares with principal/resource associations, accepters, permissions) |
| `aws_ram_resource_share_accepter` | judged as a planned AwsResourceShare kind (shares with principal/resource associations, accepters, permissions) |
| `aws_ram_resource_share_associations_exclusive` | judged as a planned AwsResourceShare kind (shares with principal/resource associations, accepters, permissions) |
| `aws_ram_sharing_with_organization` | judged as a planned AwsResourceShare kind (shares with principal/resource associations, accepters, permissions) |
| `aws_rds_global_cluster` | judged as a planned AwsRdsGlobalCluster kind |
| `aws_rds_integration` | judged as a planned AwsRdsZeroEtlIntegration kind (pairs with the Redshift-side integration) |
| `aws_redshift_authentication_profile` | cluster companion surface (auth profiles, IAM roles, endpoint access/authorization, event subscriptions, scheduled actions, snapshot schedules/copies, usage limits, resource policies, namespace registration); folds into the existing AwsRedshiftCluster kind as its spec deepens |
| `aws_redshift_cluster_iam_roles` | cluster companion surface (auth profiles, IAM roles, endpoint access/authorization, event subscriptions, scheduled actions, snapshot schedules/copies, usage limits, resource policies, namespace registration); folds into the existing AwsRedshiftCluster kind as its spec deepens |
| `aws_redshift_cluster_snapshot` | cluster companion surface (auth profiles, IAM roles, endpoint access/authorization, event subscriptions, scheduled actions, snapshot schedules/copies, usage limits, resource policies, namespace registration); folds into the existing AwsRedshiftCluster kind as its spec deepens |
| `aws_redshift_data_share_authorization` | judged as a planned AwsRedshiftDataShare kind (share authorizations and consumer associations) |
| `aws_redshift_data_share_consumer_association` | judged as a planned AwsRedshiftDataShare kind (share authorizations and consumer associations) |
| `aws_redshift_endpoint_access` | cluster companion surface (auth profiles, IAM roles, endpoint access/authorization, event subscriptions, scheduled actions, snapshot schedules/copies, usage limits, resource policies, namespace registration); folds into the existing AwsRedshiftCluster kind as its spec deepens |
| `aws_redshift_endpoint_authorization` | cluster companion surface (auth profiles, IAM roles, endpoint access/authorization, event subscriptions, scheduled actions, snapshot schedules/copies, usage limits, resource policies, namespace registration); folds into the existing AwsRedshiftCluster kind as its spec deepens |
| `aws_redshift_event_subscription` | cluster companion surface (auth profiles, IAM roles, endpoint access/authorization, event subscriptions, scheduled actions, snapshot schedules/copies, usage limits, resource policies, namespace registration); folds into the existing AwsRedshiftCluster kind as its spec deepens |
| `aws_redshift_integration` | judged as part of the planned AwsRdsZeroEtlIntegration kind (Redshift-side zero-ETL integration) |
| `aws_redshift_namespace_registration` | cluster companion surface (auth profiles, IAM roles, endpoint access/authorization, event subscriptions, scheduled actions, snapshot schedules/copies, usage limits, resource policies, namespace registration); folds into the existing AwsRedshiftCluster kind as its spec deepens |
| `aws_redshift_resource_policy` | cluster companion surface (auth profiles, IAM roles, endpoint access/authorization, event subscriptions, scheduled actions, snapshot schedules/copies, usage limits, resource policies, namespace registration); folds into the existing AwsRedshiftCluster kind as its spec deepens |
| `aws_redshift_scheduled_action` | cluster companion surface (auth profiles, IAM roles, endpoint access/authorization, event subscriptions, scheduled actions, snapshot schedules/copies, usage limits, resource policies, namespace registration); folds into the existing AwsRedshiftCluster kind as its spec deepens |
| `aws_redshift_snapshot_copy_grant` | cluster companion surface (auth profiles, IAM roles, endpoint access/authorization, event subscriptions, scheduled actions, snapshot schedules/copies, usage limits, resource policies, namespace registration); folds into the existing AwsRedshiftCluster kind as its spec deepens |
| `aws_redshift_snapshot_schedule` | cluster companion surface (auth profiles, IAM roles, endpoint access/authorization, event subscriptions, scheduled actions, snapshot schedules/copies, usage limits, resource policies, namespace registration); folds into the existing AwsRedshiftCluster kind as its spec deepens |
| `aws_redshift_snapshot_schedule_association` | cluster companion surface (auth profiles, IAM roles, endpoint access/authorization, event subscriptions, scheduled actions, snapshot schedules/copies, usage limits, resource policies, namespace registration); folds into the existing AwsRedshiftCluster kind as its spec deepens |
| `aws_redshift_usage_limit` | cluster companion surface (auth profiles, IAM roles, endpoint access/authorization, event subscriptions, scheduled actions, snapshot schedules/copies, usage limits, resource policies, namespace registration); folds into the existing AwsRedshiftCluster kind as its spec deepens |
| `aws_redshiftserverless_custom_domain_association` | serverless companion surface (custom domains, endpoint access, resource policies, usage limits); folds into the existing AwsRedshiftServerlessNamespace/AwsRedshiftServerlessWorkgroup kinds as their specs deepen |
| `aws_redshiftserverless_endpoint_access` | serverless companion surface (custom domains, endpoint access, resource policies, usage limits); folds into the existing AwsRedshiftServerlessNamespace/AwsRedshiftServerlessWorkgroup kinds as their specs deepen |
| `aws_redshiftserverless_resource_policy` | serverless companion surface (custom domains, endpoint access, resource policies, usage limits); folds into the existing AwsRedshiftServerlessNamespace/AwsRedshiftServerlessWorkgroup kinds as their specs deepen |
| `aws_redshiftserverless_usage_limit` | serverless companion surface (custom domains, endpoint access, resource policies, usage limits); folds into the existing AwsRedshiftServerlessNamespace/AwsRedshiftServerlessWorkgroup kinds as their specs deepen |
| `aws_route` | VPC companion surface (flow logs, IPv6 associations, block-public-access, encryption control, route management); folds into the existing AwsVpc kind as its spec deepens |
| `aws_route53_cidr_collection` | account-scoped container of CIDR blocks for IP-based routing, referenced by many records across zones -- an owning collection, never one record's satellite; judged as a planned AwsRoute53CidrCollection kind (AwsRoute53DnsRecord.spec.routing_policy.cidr already composes onto it by collection id + location name) |
| `aws_route53_cidr_location` | a named group of CIDR blocks INSIDE a collection (collection-scoped child, the only updatable part); folds into the planned AwsRoute53CidrCollection kind as its entries |
| `aws_route53_delegation_set` | account-scoped reusable delegation set: one fixed four-name-server set shared by MANY zones (white-label/vanity name servers, bulk migrations) -- an account object zones reference, never one zone's satellite; judged as a planned AwsRoute53DelegationSet kind (AwsRoute53Zone.spec.delegation_set_id already composes onto it by id) |
| `aws_route53_resolver_config` | per-VPC resolver settings; fold into the existing AwsVpc kind as its spec deepens |
| `aws_route53_resolver_dnssec_config` | per-VPC resolver settings; fold into the existing AwsVpc kind as its spec deepens |
| `aws_route53_resolver_endpoint` | judged as a planned AwsRoute53ResolverEndpoint kind (endpoints, forwarding rules, rule associations) |
| `aws_route53_resolver_firewall_config` | judged as a planned AwsRoute53ResolverFirewall kind (configs, domain lists, rules, rule groups, associations) |
| `aws_route53_resolver_firewall_domain_list` | judged as a planned AwsRoute53ResolverFirewall kind (configs, domain lists, rules, rule groups, associations) |
| `aws_route53_resolver_firewall_rule` | judged as a planned AwsRoute53ResolverFirewall kind (configs, domain lists, rules, rule groups, associations) |
| `aws_route53_resolver_firewall_rule_group` | judged as a planned AwsRoute53ResolverFirewall kind (configs, domain lists, rules, rule groups, associations) |
| `aws_route53_resolver_firewall_rule_group_association` | judged as a planned AwsRoute53ResolverFirewall kind (configs, domain lists, rules, rule groups, associations) |
| `aws_route53_resolver_query_log_config` | judged as a planned AwsRoute53ResolverQueryLog kind (configs with associations) |
| `aws_route53_resolver_query_log_config_association` | judged as a planned AwsRoute53ResolverQueryLog kind (configs with associations) |
| `aws_route53_resolver_rule` | judged as a planned AwsRoute53ResolverEndpoint kind (endpoints, forwarding rules, rule associations) |
| `aws_route53_resolver_rule_association` | judged as a planned AwsRoute53ResolverEndpoint kind (endpoints, forwarding rules, rule associations) |
| `aws_route53_vpc_association_authorization` | zone-owner side of the CROSS-ACCOUNT private-zone handshake (authorizes another account's VPC to associate; runs with the zone owner's credentials); same-account associations are modeled on AwsRoute53Zone.spec.vpc_associations -- judged as a planned AwsRoute53VpcAssociationAuthorization kind (the cross-account-half class, the TGW vpc-attachment-accepter precedent) |
| `aws_route53_zone_association` | VPC-owner side of the CROSS-ACCOUNT private-zone handshake (associates a foreign account's authorized zone; runs with the VPC owner's credentials); same-account associations are modeled on AwsRoute53Zone.spec.vpc_associations -- judged as a planned AwsRoute53ZoneAssociation kind pairing with AwsRoute53VpcAssociationAuthorization |
| `aws_route53domains_delegation_signer_record` | judged as a planned AwsRoute53Domain kind (registered domains, delegation signer records) |
| `aws_route53domains_domain` | judged as a planned AwsRoute53Domain kind (registered domains, delegation signer records) |
| `aws_route53domains_registered_domain` | judged as a planned AwsRoute53Domain kind (registered domains, delegation signer records) |
| `aws_s3_access_point` | judged as a planned AwsS3AccessPoint kind |
| `aws_s3_account_public_access_block` | judged as a planned AwsS3AccountPublicAccessBlock kind |
| `aws_s3_directory_bucket` | judged as a planned AwsS3DirectoryBucket kind (S3 Express One Zone) |
| `aws_s3_object_copy` | object copies fold into the existing AwsS3ObjectSet kind as its spec deepens |
| `aws_s3control_access_grant` | judged as a planned AwsS3AccessGrants kind (instances, locations, grants, policies) |
| `aws_s3control_access_grants_instance` | judged as a planned AwsS3AccessGrants kind (instances, locations, grants, policies) |
| `aws_s3control_access_grants_instance_resource_policy` | judged as a planned AwsS3AccessGrants kind (instances, locations, grants, policies) |
| `aws_s3control_access_grants_location` | judged as a planned AwsS3AccessGrants kind (instances, locations, grants, policies) |
| `aws_s3control_access_point_policy` | access-point policy surface; folds into the planned AwsS3AccessPoint kind |
| `aws_s3control_directory_bucket_access_point_scope` | access-point policy surface; folds into the planned AwsS3AccessPoint kind |
| `aws_s3control_multi_region_access_point` | judged as a planned AwsS3MultiRegionAccessPoint kind (access points, policies, routes) |
| `aws_s3control_multi_region_access_point_policy` | judged as a planned AwsS3MultiRegionAccessPoint kind (access points, policies, routes) |
| `aws_s3control_multi_region_access_point_routes` | judged as a planned AwsS3MultiRegionAccessPoint kind (access points, policies, routes) |
| `aws_s3control_storage_lens_configuration` | judged as a planned AwsS3StorageLens kind |
| `aws_s3tables_namespace` | judged as a planned AwsS3TableBucket kind (table buckets, namespaces, tables, policies, replication) |
| `aws_s3tables_table` | judged as a planned AwsS3TableBucket kind (table buckets, namespaces, tables, policies, replication) |
| `aws_s3tables_table_bucket` | judged as a planned AwsS3TableBucket kind (table buckets, namespaces, tables, policies, replication) |
| `aws_s3tables_table_bucket_policy` | judged as a planned AwsS3TableBucket kind (table buckets, namespaces, tables, policies, replication) |
| `aws_s3tables_table_bucket_replication` | judged as a planned AwsS3TableBucket kind (table buckets, namespaces, tables, policies, replication) |
| `aws_s3tables_table_policy` | judged as a planned AwsS3TableBucket kind (table buckets, namespaces, tables, policies, replication) |
| `aws_s3tables_table_replication` | judged as a planned AwsS3TableBucket kind (table buckets, namespaces, tables, policies, replication) |
| `aws_s3vectors_index` | judged as a planned AwsS3VectorBucket kind (vector buckets, indexes, policies) |
| `aws_s3vectors_vector_bucket` | judged as a planned AwsS3VectorBucket kind (vector buckets, indexes, policies) |
| `aws_s3vectors_vector_bucket_policy` | judged as a planned AwsS3VectorBucket kind (vector buckets, indexes, policies) |
| `aws_sagemaker_app` | Studio domain companion surface (apps, spaces, user profiles, lifecycle configs, portfolio status); folds into the existing AwsSagemakerDomain kind as its spec deepens |
| `aws_sagemaker_app_image_config` | Studio domain companion surface (apps, spaces, user profiles, lifecycle configs, portfolio status); folds into the existing AwsSagemakerDomain kind as its spec deepens |
| `aws_sagemaker_endpoint` | judged as a planned AwsSagemakerEndpoint kind (endpoints with configurations) |
| `aws_sagemaker_endpoint_configuration` | judged as a planned AwsSagemakerEndpoint kind (endpoints with configurations) |
| `aws_sagemaker_feature_group` | judged as a planned AwsSagemakerFeatureGroup kind |
| `aws_sagemaker_image` | judged as a planned AwsSagemakerImage kind (images with versions) |
| `aws_sagemaker_image_version` | judged as a planned AwsSagemakerImage kind (images with versions) |
| `aws_sagemaker_mlflow_app` | judged as a planned AwsSagemakerMlflowServer kind (tracking servers with apps) |
| `aws_sagemaker_mlflow_tracking_server` | judged as a planned AwsSagemakerMlflowServer kind (tracking servers with apps) |
| `aws_sagemaker_model` | judged as a planned AwsSagemakerModel kind |
| `aws_sagemaker_model_package_group` | judged as a planned AwsSagemakerModelRegistry kind (package groups with policies) |
| `aws_sagemaker_model_package_group_policy` | judged as a planned AwsSagemakerModelRegistry kind (package groups with policies) |
| `aws_sagemaker_notebook_instance` | judged as a planned AwsSagemakerNotebookInstance kind (instances with lifecycle configurations) |
| `aws_sagemaker_notebook_instance_lifecycle_configuration` | judged as a planned AwsSagemakerNotebookInstance kind (instances with lifecycle configurations) |
| `aws_sagemaker_pipeline` | judged as a planned AwsSagemakerPipeline kind |
| `aws_sagemaker_servicecatalog_portfolio_status` | Studio domain companion surface (apps, spaces, user profiles, lifecycle configs, portfolio status); folds into the existing AwsSagemakerDomain kind as its spec deepens |
| `aws_sagemaker_space` | Studio domain companion surface (apps, spaces, user profiles, lifecycle configs, portfolio status); folds into the existing AwsSagemakerDomain kind as its spec deepens |
| `aws_sagemaker_studio_lifecycle_config` | Studio domain companion surface (apps, spaces, user profiles, lifecycle configs, portfolio status); folds into the existing AwsSagemakerDomain kind as its spec deepens |
| `aws_sagemaker_user_profile` | Studio domain companion surface (apps, spaces, user profiles, lifecycle configs, portfolio status); folds into the existing AwsSagemakerDomain kind as its spec deepens |
| `aws_scheduler_schedule` | judged as a planned AwsEventbridgeScheduler kind (schedules with groups) |
| `aws_scheduler_schedule_group` | judged as a planned AwsEventbridgeScheduler kind (schedules with groups) |
| `aws_secretsmanager_secret_policy` | judged as a planned AwsSecretsManagerSecret kind (secret with versions, rotation, policy, tags) |
| `aws_secretsmanager_secret_rotation` | judged as a planned AwsSecretsManagerSecret kind (secret with versions, rotation, policy, tags) |
| `aws_secretsmanager_tag` | judged as a planned AwsSecretsManagerSecret kind (secret with versions, rotation, policy, tags) |
| `aws_security_group_rule` | security-group rule resources; fold into the existing AwsSecurityGroup kind as its spec deepens |
| `aws_securityhub_account` | judged as a planned AwsSecurityHub kind (account enablement v1/v2, organization configuration, standards, insights, action targets, automation rules, aggregators, connectors) |
| `aws_securityhub_account_v2` | judged as a planned AwsSecurityHub kind (account enablement v1/v2, organization configuration, standards, insights, action targets, automation rules, aggregators, connectors) |
| `aws_securityhub_action_target` | judged as a planned AwsSecurityHub kind (account enablement v1/v2, organization configuration, standards, insights, action targets, automation rules, aggregators, connectors) |
| `aws_securityhub_aggregator_v2` | judged as a planned AwsSecurityHub kind (account enablement v1/v2, organization configuration, standards, insights, action targets, automation rules, aggregators, connectors) |
| `aws_securityhub_automation_rule` | judged as a planned AwsSecurityHub kind (account enablement v1/v2, organization configuration, standards, insights, action targets, automation rules, aggregators, connectors) |
| `aws_securityhub_automation_rule_v2` | judged as a planned AwsSecurityHub kind (account enablement v1/v2, organization configuration, standards, insights, action targets, automation rules, aggregators, connectors) |
| `aws_securityhub_configuration_policy` | judged as a planned AwsSecurityHub kind (account enablement v1/v2, organization configuration, standards, insights, action targets, automation rules, aggregators, connectors) |
| `aws_securityhub_configuration_policy_association` | judged as a planned AwsSecurityHub kind (account enablement v1/v2, organization configuration, standards, insights, action targets, automation rules, aggregators, connectors) |
| `aws_securityhub_connector_v2` | judged as a planned AwsSecurityHub kind (account enablement v1/v2, organization configuration, standards, insights, action targets, automation rules, aggregators, connectors) |
| `aws_securityhub_finding_aggregator` | judged as a planned AwsSecurityHub kind (account enablement v1/v2, organization configuration, standards, insights, action targets, automation rules, aggregators, connectors) |
| `aws_securityhub_insight` | judged as a planned AwsSecurityHub kind (account enablement v1/v2, organization configuration, standards, insights, action targets, automation rules, aggregators, connectors) |
| `aws_securityhub_invite_accepter` | judged as a planned AwsSecurityHub kind (account enablement v1/v2, organization configuration, standards, insights, action targets, automation rules, aggregators, connectors) |
| `aws_securityhub_member` | judged as a planned AwsSecurityHub kind (account enablement v1/v2, organization configuration, standards, insights, action targets, automation rules, aggregators, connectors) |
| `aws_securityhub_organization_admin_account` | judged as a planned AwsSecurityHub kind (account enablement v1/v2, organization configuration, standards, insights, action targets, automation rules, aggregators, connectors) |
| `aws_securityhub_organization_configuration` | judged as a planned AwsSecurityHub kind (account enablement v1/v2, organization configuration, standards, insights, action targets, automation rules, aggregators, connectors) |
| `aws_securityhub_product_subscription` | judged as a planned AwsSecurityHub kind (account enablement v1/v2, organization configuration, standards, insights, action targets, automation rules, aggregators, connectors) |
| `aws_securityhub_standards_control` | judged as a planned AwsSecurityHub kind (account enablement v1/v2, organization configuration, standards, insights, action targets, automation rules, aggregators, connectors) |
| `aws_securityhub_standards_control_association` | judged as a planned AwsSecurityHub kind (account enablement v1/v2, organization configuration, standards, insights, action targets, automation rules, aggregators, connectors) |
| `aws_securityhub_standards_subscription` | judged as a planned AwsSecurityHub kind (account enablement v1/v2, organization configuration, standards, insights, action targets, automation rules, aggregators, connectors) |
| `aws_service_discovery_http_namespace` | judged as a planned AwsCloudMapNamespace kind (HTTP/DNS namespaces, services, instances) |
| `aws_service_discovery_instance` | judged as a planned AwsCloudMapNamespace kind (HTTP/DNS namespaces, services, instances) |
| `aws_service_discovery_private_dns_namespace` | judged as a planned AwsCloudMapNamespace kind (HTTP/DNS namespaces, services, instances) |
| `aws_service_discovery_public_dns_namespace` | judged as a planned AwsCloudMapNamespace kind (HTTP/DNS namespaces, services, instances) |
| `aws_service_discovery_service` | judged as a planned AwsCloudMapNamespace kind (HTTP/DNS namespaces, services, instances) |
| `aws_servicecatalog_budget_resource_association` | judged as a planned AwsServiceCatalogPortfolio kind (portfolios, shares, principal associations, constraints, budgets, tag options, service actions, organizations access) |
| `aws_servicecatalog_constraint` | judged as a planned AwsServiceCatalogPortfolio kind (portfolios, shares, principal associations, constraints, budgets, tag options, service actions, organizations access) |
| `aws_servicecatalog_organizations_access` | judged as a planned AwsServiceCatalogPortfolio kind (portfolios, shares, principal associations, constraints, budgets, tag options, service actions, organizations access) |
| `aws_servicecatalog_portfolio` | judged as a planned AwsServiceCatalogPortfolio kind (portfolios, shares, principal associations, constraints, budgets, tag options, service actions, organizations access) |
| `aws_servicecatalog_portfolio_share` | judged as a planned AwsServiceCatalogPortfolio kind (portfolios, shares, principal associations, constraints, budgets, tag options, service actions, organizations access) |
| `aws_servicecatalog_principal_portfolio_association` | judged as a planned AwsServiceCatalogPortfolio kind (portfolios, shares, principal associations, constraints, budgets, tag options, service actions, organizations access) |
| `aws_servicecatalog_product` | judged as a planned AwsServiceCatalogProduct kind (products, portfolio associations, provisioning artifacts, provisioned products) |
| `aws_servicecatalog_product_portfolio_association` | judged as a planned AwsServiceCatalogProduct kind (products, portfolio associations, provisioning artifacts, provisioned products) |
| `aws_servicecatalog_provisioned_product` | judged as a planned AwsServiceCatalogProduct kind (products, portfolio associations, provisioning artifacts, provisioned products) |
| `aws_servicecatalog_provisioning_artifact` | judged as a planned AwsServiceCatalogProduct kind (products, portfolio associations, provisioning artifacts, provisioned products) |
| `aws_servicecatalog_service_action` | judged as a planned AwsServiceCatalogPortfolio kind (portfolios, shares, principal associations, constraints, budgets, tag options, service actions, organizations access) |
| `aws_servicecatalog_tag_option` | judged as a planned AwsServiceCatalogPortfolio kind (portfolios, shares, principal associations, constraints, budgets, tag options, service actions, organizations access) |
| `aws_servicecatalog_tag_option_resource_association` | judged as a planned AwsServiceCatalogPortfolio kind (portfolios, shares, principal associations, constraints, budgets, tag options, service actions, organizations access) |
| `aws_servicequotas_auto_management` | judged as a planned AwsServiceQuota kind (quota requests, templates, auto-management) |
| `aws_servicequotas_service_quota` | judged as a planned AwsServiceQuota kind (quota requests, templates, auto-management) |
| `aws_servicequotas_template` | judged as a planned AwsServiceQuota kind (quota requests, templates, auto-management) |
| `aws_servicequotas_template_association` | judged as a planned AwsServiceQuota kind (quota requests, templates, auto-management) |
| `aws_ses_active_receipt_rule_set` | judged as a planned AwsSesReceiptRuleSet kind (inbound email: filters, rules, rule sets, activation) |
| `aws_ses_receipt_filter` | judged as a planned AwsSesReceiptRuleSet kind (inbound email: filters, rules, rule sets, activation) |
| `aws_ses_receipt_rule` | judged as a planned AwsSesReceiptRuleSet kind (inbound email: filters, rules, rule sets, activation) |
| `aws_ses_receipt_rule_set` | judged as a planned AwsSesReceiptRuleSet kind (inbound email: filters, rules, rule sets, activation) |
| `aws_sesv2_account_suppression_attributes` | account-level SES attributes (suppression, VDM); fold into the existing AwsSesConfigurationSet kind family's account composition |
| `aws_sesv2_account_vdm_attributes` | account-level SES attributes (suppression, VDM); fold into the existing AwsSesConfigurationSet kind family's account composition |
| `aws_sesv2_dedicated_ip_assignment` | judged as a planned AwsSesDedicatedIpPool kind (pools with IP assignments) |
| `aws_sesv2_dedicated_ip_pool` | judged as a planned AwsSesDedicatedIpPool kind (pools with IP assignments) |
| `aws_sfn_activity` | Step Functions companion surface (activities, aliases); folds into the existing AwsStepFunction kind as its spec deepens |
| `aws_sfn_alias` | Step Functions companion surface (activities, aliases); folds into the existing AwsStepFunction kind as its spec deepens |
| `aws_shield_application_layer_automatic_response` | judged as a planned AwsShieldAdvanced kind (subscription, protections, protection groups, DRT access, proactive engagement) |
| `aws_shield_drt_access_log_bucket_association` | judged as a planned AwsShieldAdvanced kind (subscription, protections, protection groups, DRT access, proactive engagement) |
| `aws_shield_drt_access_role_arn_association` | judged as a planned AwsShieldAdvanced kind (subscription, protections, protection groups, DRT access, proactive engagement) |
| `aws_shield_proactive_engagement` | judged as a planned AwsShieldAdvanced kind (subscription, protections, protection groups, DRT access, proactive engagement) |
| `aws_shield_protection` | judged as a planned AwsShieldAdvanced kind (subscription, protections, protection groups, DRT access, proactive engagement) |
| `aws_shield_protection_group` | judged as a planned AwsShieldAdvanced kind (subscription, protections, protection groups, DRT access, proactive engagement) |
| `aws_shield_protection_health_check_association` | judged as a planned AwsShieldAdvanced kind (subscription, protections, protection groups, DRT access, proactive engagement) |
| `aws_shield_subscription` | judged as a planned AwsShieldAdvanced kind (subscription, protections, protection groups, DRT access, proactive engagement) |
| `aws_snapshot_create_volume_permission` | judged as a planned AwsEbsSnapshot kind (snapshots, copies, imports, fast restore, volume permissions) |
| `aws_spot_datafeed_subscription` | account-wide EC2 toggles fold into the planned AwsEc2AccountSettings kind |
| `aws_ssm_activation` | SSM account-level plumbing (hybrid activations, resource data syncs, service settings); folds into the planned SSM kind family led by AwsSsmDocument |
| `aws_ssm_association` | judged as a planned AwsSsmDocument kind (documents with associations) |
| `aws_ssm_default_patch_baseline` | judged as a planned AwsSsmPatchBaseline kind (baselines, groups, defaults) |
| `aws_ssm_document` | judged as a planned AwsSsmDocument kind (documents with associations) |
| `aws_ssm_maintenance_window` | judged as a planned AwsSsmMaintenanceWindow kind (windows, targets, tasks) |
| `aws_ssm_maintenance_window_target` | judged as a planned AwsSsmMaintenanceWindow kind (windows, targets, tasks) |
| `aws_ssm_maintenance_window_task` | judged as a planned AwsSsmMaintenanceWindow kind (windows, targets, tasks) |
| `aws_ssm_parameter` | judged as a planned AwsSsmParameter kind |
| `aws_ssm_patch_baseline` | judged as a planned AwsSsmPatchBaseline kind (baselines, groups, defaults) |
| `aws_ssm_patch_group` | judged as a planned AwsSsmPatchBaseline kind (baselines, groups, defaults) |
| `aws_ssm_resource_data_sync` | SSM account-level plumbing (hybrid activations, resource data syncs, service settings); folds into the planned SSM kind family led by AwsSsmDocument |
| `aws_ssm_service_setting` | SSM account-level plumbing (hybrid activations, resource data syncs, service settings); folds into the planned SSM kind family led by AwsSsmDocument |
| `aws_ssoadmin_account_assignment` | judged as a planned AwsIdentityCenterAssignment kind |
| `aws_ssoadmin_application` | judged as a planned AwsIdentityCenterApplication kind (applications, access scopes, assignments, trusted token issuers) |
| `aws_ssoadmin_application_access_scope` | judged as a planned AwsIdentityCenterApplication kind (applications, access scopes, assignments, trusted token issuers) |
| `aws_ssoadmin_application_assignment` | judged as a planned AwsIdentityCenterApplication kind (applications, access scopes, assignments, trusted token issuers) |
| `aws_ssoadmin_application_assignment_configuration` | judged as a planned AwsIdentityCenterApplication kind (applications, access scopes, assignments, trusted token issuers) |
| `aws_ssoadmin_customer_managed_policy_attachment` | judged as a planned AwsIdentityCenterPermissionSet kind (permission sets, policy attachments, boundaries; instance-level settings fold in) |
| `aws_ssoadmin_customer_managed_policy_attachments_exclusive` | judged as a planned AwsIdentityCenterPermissionSet kind (permission sets, policy attachments, boundaries; instance-level settings fold in) |
| `aws_ssoadmin_instance_access_control_attributes` | judged as a planned AwsIdentityCenterPermissionSet kind (permission sets, policy attachments, boundaries; instance-level settings fold in) |
| `aws_ssoadmin_managed_policy_attachment` | judged as a planned AwsIdentityCenterPermissionSet kind (permission sets, policy attachments, boundaries; instance-level settings fold in) |
| `aws_ssoadmin_managed_policy_attachments_exclusive` | judged as a planned AwsIdentityCenterPermissionSet kind (permission sets, policy attachments, boundaries; instance-level settings fold in) |
| `aws_ssoadmin_permission_set` | judged as a planned AwsIdentityCenterPermissionSet kind (permission sets, policy attachments, boundaries; instance-level settings fold in) |
| `aws_ssoadmin_permission_set_inline_policy` | judged as a planned AwsIdentityCenterPermissionSet kind (permission sets, policy attachments, boundaries; instance-level settings fold in) |
| `aws_ssoadmin_permissions_boundary_attachment` | judged as a planned AwsIdentityCenterPermissionSet kind (permission sets, policy attachments, boundaries; instance-level settings fold in) |
| `aws_ssoadmin_region` | judged as a planned AwsIdentityCenterPermissionSet kind (permission sets, policy attachments, boundaries; instance-level settings fold in) |
| `aws_ssoadmin_trusted_token_issuer` | judged as a planned AwsIdentityCenterApplication kind (applications, access scopes, assignments, trusted token issuers) |
| `aws_storagegateway_cache` | judged as a planned AwsStorageGateway kind (gateways, file shares, volumes, tape pools, caches) |
| `aws_storagegateway_cached_iscsi_volume` | judged as a planned AwsStorageGateway kind (gateways, file shares, volumes, tape pools, caches) |
| `aws_storagegateway_file_system_association` | judged as a planned AwsStorageGateway kind (gateways, file shares, volumes, tape pools, caches) |
| `aws_storagegateway_gateway` | judged as a planned AwsStorageGateway kind (gateways, file shares, volumes, tape pools, caches) |
| `aws_storagegateway_nfs_file_share` | judged as a planned AwsStorageGateway kind (gateways, file shares, volumes, tape pools, caches) |
| `aws_storagegateway_smb_file_share` | judged as a planned AwsStorageGateway kind (gateways, file shares, volumes, tape pools, caches) |
| `aws_storagegateway_stored_iscsi_volume` | judged as a planned AwsStorageGateway kind (gateways, file shares, volumes, tape pools, caches) |
| `aws_storagegateway_tape_pool` | judged as a planned AwsStorageGateway kind (gateways, file shares, volumes, tape pools, caches) |
| `aws_storagegateway_upload_buffer` | judged as a planned AwsStorageGateway kind (gateways, file shares, volumes, tape pools, caches) |
| `aws_storagegateway_working_storage` | judged as a planned AwsStorageGateway kind (gateways, file shares, volumes, tape pools, caches) |
| `aws_synthetics_canary` | judged as a planned AwsCloudwatchSynthetics kind (canaries, groups, associations) |
| `aws_synthetics_group` | judged as a planned AwsCloudwatchSynthetics kind (canaries, groups, associations) |
| `aws_synthetics_group_association` | judged as a planned AwsCloudwatchSynthetics kind (canaries, groups, associations) |
| `aws_transfer_access` | judged as a planned AwsTransferServer kind (servers, users, SSH keys, access, host keys, workflows, web apps, tags) |
| `aws_transfer_agreement` | judged as a planned AwsTransferConnector kind (connectors, agreements, profiles, certificates) |
| `aws_transfer_certificate` | judged as a planned AwsTransferConnector kind (connectors, agreements, profiles, certificates) |
| `aws_transfer_connector` | judged as a planned AwsTransferConnector kind (connectors, agreements, profiles, certificates) |
| `aws_transfer_host_key` | judged as a planned AwsTransferServer kind (servers, users, SSH keys, access, host keys, workflows, web apps, tags) |
| `aws_transfer_profile` | judged as a planned AwsTransferConnector kind (connectors, agreements, profiles, certificates) |
| `aws_transfer_server` | judged as a planned AwsTransferServer kind (servers, users, SSH keys, access, host keys, workflows, web apps, tags) |
| `aws_transfer_ssh_key` | judged as a planned AwsTransferServer kind (servers, users, SSH keys, access, host keys, workflows, web apps, tags) |
| `aws_transfer_tag` | judged as a planned AwsTransferServer kind (servers, users, SSH keys, access, host keys, workflows, web apps, tags) |
| `aws_transfer_user` | judged as a planned AwsTransferServer kind (servers, users, SSH keys, access, host keys, workflows, web apps, tags) |
| `aws_transfer_web_app` | judged as a planned AwsTransferServer kind (servers, users, SSH keys, access, host keys, workflows, web apps, tags) |
| `aws_transfer_web_app_customization` | judged as a planned AwsTransferServer kind (servers, users, SSH keys, access, host keys, workflows, web apps, tags) |
| `aws_transfer_workflow` | judged as a planned AwsTransferServer kind (servers, users, SSH keys, access, host keys, workflows, web apps, tags) |
| `aws_verifiedaccess_endpoint` | judged as a planned AwsVerifiedAccess kind (instances, trust providers, groups, endpoints, logging) |
| `aws_verifiedaccess_group` | judged as a planned AwsVerifiedAccess kind (instances, trust providers, groups, endpoints, logging) |
| `aws_verifiedaccess_instance` | judged as a planned AwsVerifiedAccess kind (instances, trust providers, groups, endpoints, logging) |
| `aws_verifiedaccess_instance_logging_configuration` | judged as a planned AwsVerifiedAccess kind (instances, trust providers, groups, endpoints, logging) |
| `aws_verifiedaccess_instance_trust_provider_attachment` | judged as a planned AwsVerifiedAccess kind (instances, trust providers, groups, endpoints, logging) |
| `aws_verifiedaccess_trust_provider` | judged as a planned AwsVerifiedAccess kind (instances, trust providers, groups, endpoints, logging) |
| `aws_volume_attachment` | judged as a planned AwsEbsVolume kind (volumes, copies, attachments) |
| `aws_vpc_block_public_access_exclusion` | judged as a planned AwsVpcBlockPublicAccess kind (the per-VPC/per-subnet exclusion arm of the regional block-public-access singleton) |
| `aws_vpc_block_public_access_options` | judged as a planned AwsVpcBlockPublicAccess kind (a REGIONAL account-level singleton -- its id is the region -- paired with per-VPC/subnet exclusions; never per-VPC surface) |
| `aws_vpc_dhcp_options` | judged as a planned AwsVpcDhcpOptions kind (option sets with associations) |
| `aws_vpc_dhcp_options_association` | judged as a planned AwsVpcDhcpOptions kind (option sets with associations) |
| `aws_vpc_endpoint_connection_accepter` | judged as a planned AwsVpcEndpointService kind (PrivateLink provider side: services, allowed principals, DNS verification, connection accepters/notifications) |
| `aws_vpc_endpoint_connection_notification` | judged as a planned AwsVpcEndpointService kind (PrivateLink provider side: services, allowed principals, DNS verification, connection accepters/notifications) |
| `aws_vpc_endpoint_policy` | endpoint companion surface (policies, private DNS, route table/security group/subnet associations); folds into the existing AwsVpcEndpoint kind as its spec deepens |
| `aws_vpc_endpoint_private_dns` | endpoint companion surface (policies, private DNS, route table/security group/subnet associations); folds into the existing AwsVpcEndpoint kind as its spec deepens |
| `aws_vpc_endpoint_route_table_association` | endpoint companion surface (policies, private DNS, route table/security group/subnet associations); folds into the existing AwsVpcEndpoint kind as its spec deepens |
| `aws_vpc_endpoint_security_group_association` | endpoint companion surface (policies, private DNS, route table/security group/subnet associations); folds into the existing AwsVpcEndpoint kind as its spec deepens |
| `aws_vpc_endpoint_service` | judged as a planned AwsVpcEndpointService kind (PrivateLink provider side: services, allowed principals, DNS verification, connection accepters/notifications) |
| `aws_vpc_endpoint_service_allowed_principal` | judged as a planned AwsVpcEndpointService kind (PrivateLink provider side: services, allowed principals, DNS verification, connection accepters/notifications) |
| `aws_vpc_endpoint_service_private_dns_verification` | judged as a planned AwsVpcEndpointService kind (PrivateLink provider side: services, allowed principals, DNS verification, connection accepters/notifications) |
| `aws_vpc_endpoint_subnet_association` | endpoint companion surface (policies, private DNS, route table/security group/subnet associations); folds into the existing AwsVpcEndpoint kind as its spec deepens |
| `aws_vpc_ipam` | judged as a planned AwsVpcIpam kind (IPAM, scopes, pools, CIDR allocations, resource discovery, organization admin) |
| `aws_vpc_ipam_organization_admin_account` | judged as a planned AwsVpcIpam kind (IPAM, scopes, pools, CIDR allocations, resource discovery, organization admin) |
| `aws_vpc_ipam_pool` | judged as a planned AwsVpcIpam kind (IPAM, scopes, pools, CIDR allocations, resource discovery, organization admin) |
| `aws_vpc_ipam_pool_cidr` | judged as a planned AwsVpcIpam kind (IPAM, scopes, pools, CIDR allocations, resource discovery, organization admin) |
| `aws_vpc_ipam_pool_cidr_allocation` | judged as a planned AwsVpcIpam kind (IPAM, scopes, pools, CIDR allocations, resource discovery, organization admin) |
| `aws_vpc_ipam_preview_next_cidr` | judged as a planned AwsVpcIpam kind (IPAM, scopes, pools, CIDR allocations, resource discovery, organization admin) |
| `aws_vpc_ipam_resource_discovery` | judged as a planned AwsVpcIpam kind (IPAM, scopes, pools, CIDR allocations, resource discovery, organization admin) |
| `aws_vpc_ipam_resource_discovery_association` | judged as a planned AwsVpcIpam kind (IPAM, scopes, pools, CIDR allocations, resource discovery, organization admin) |
| `aws_vpc_ipam_scope` | judged as a planned AwsVpcIpam kind (IPAM, scopes, pools, CIDR allocations, resource discovery, organization admin) |
| `aws_vpc_peering_connection` | judged as a planned AwsVpcPeering kind (peering connections, accepters, options) |
| `aws_vpc_peering_connection_accepter` | judged as a planned AwsVpcPeering kind (peering connections, accepters, options) |
| `aws_vpc_peering_connection_options` | judged as a planned AwsVpcPeering kind (peering connections, accepters, options) |
| `aws_vpc_security_group_egress_rule` | security-group rule resources; fold into the existing AwsSecurityGroup kind as its spec deepens |
| `aws_vpc_security_group_ingress_rule` | security-group rule resources; fold into the existing AwsSecurityGroup kind as its spec deepens |
| `aws_vpc_security_group_rules_exclusive` | security-group rule resources; fold into the existing AwsSecurityGroup kind as its spec deepens |
| `aws_vpc_security_group_vpc_association` | security-group rule resources; fold into the existing AwsSecurityGroup kind as its spec deepens |
| `aws_vpclattice_access_log_subscription` | judged as a planned AwsVpcLatticeServiceNetwork kind (service networks, associations, auth/resource policies, resource gateways/configurations, log subscriptions, domain verifications) |
| `aws_vpclattice_auth_policy` | judged as a planned AwsVpcLatticeServiceNetwork kind (service networks, associations, auth/resource policies, resource gateways/configurations, log subscriptions, domain verifications) |
| `aws_vpclattice_domain_verification` | judged as a planned AwsVpcLatticeServiceNetwork kind (service networks, associations, auth/resource policies, resource gateways/configurations, log subscriptions, domain verifications) |
| `aws_vpclattice_listener` | judged as a planned AwsVpcLatticeService kind (services with listeners and rules) |
| `aws_vpclattice_listener_rule` | judged as a planned AwsVpcLatticeService kind (services with listeners and rules) |
| `aws_vpclattice_resource_configuration` | judged as a planned AwsVpcLatticeServiceNetwork kind (service networks, associations, auth/resource policies, resource gateways/configurations, log subscriptions, domain verifications) |
| `aws_vpclattice_resource_gateway` | judged as a planned AwsVpcLatticeServiceNetwork kind (service networks, associations, auth/resource policies, resource gateways/configurations, log subscriptions, domain verifications) |
| `aws_vpclattice_resource_policy` | judged as a planned AwsVpcLatticeServiceNetwork kind (service networks, associations, auth/resource policies, resource gateways/configurations, log subscriptions, domain verifications) |
| `aws_vpclattice_service` | judged as a planned AwsVpcLatticeService kind (services with listeners and rules) |
| `aws_vpclattice_service_network` | judged as a planned AwsVpcLatticeServiceNetwork kind (service networks, associations, auth/resource policies, resource gateways/configurations, log subscriptions, domain verifications) |
| `aws_vpclattice_service_network_resource_association` | judged as a planned AwsVpcLatticeServiceNetwork kind (service networks, associations, auth/resource policies, resource gateways/configurations, log subscriptions, domain verifications) |
| `aws_vpclattice_service_network_service_association` | judged as a planned AwsVpcLatticeServiceNetwork kind (service networks, associations, auth/resource policies, resource gateways/configurations, log subscriptions, domain verifications) |
| `aws_vpclattice_service_network_vpc_association` | judged as a planned AwsVpcLatticeServiceNetwork kind (service networks, associations, auth/resource policies, resource gateways/configurations, log subscriptions, domain verifications) |
| `aws_vpclattice_target_group` | judged as a planned AwsVpcLatticeTargetGroup kind (target groups with attachments) |
| `aws_vpclattice_target_group_attachment` | judged as a planned AwsVpcLatticeTargetGroup kind (target groups with attachments) |
| `aws_vpn_concentrator` | judged as a planned AwsSiteToSiteVpn kind (customer/VPN gateways, connections, routes, concentrators, route propagation) |
| `aws_vpn_connection` | judged as a planned AwsSiteToSiteVpn kind (customer/VPN gateways, connections, routes, concentrators, route propagation) |
| `aws_vpn_connection_route` | judged as a planned AwsSiteToSiteVpn kind (customer/VPN gateways, connections, routes, concentrators, route propagation) |
| `aws_vpn_gateway` | judged as a planned AwsSiteToSiteVpn kind (customer/VPN gateways, connections, routes, concentrators, route propagation) |
| `aws_vpn_gateway_attachment` | judged as a planned AwsSiteToSiteVpn kind (customer/VPN gateways, connections, routes, concentrators, route propagation) |
| `aws_vpn_gateway_route_propagation` | judged as a planned AwsSiteToSiteVpn kind (customer/VPN gateways, connections, routes, concentrators, route propagation) |
| `aws_wafv2_api_key` | judged as a planned AwsWafApiKey kind: a scope-level API key for the CAPTCHA/challenge client-side JS integration, create/delete-only with a sensitive key output serving up to 5 token domains -- a standalone lifecycle that does not fold into AwsWafWebAcl (whose 2026-08-11 depth closure completed the spec deepening the prior blanket reason waited on) |
| `aws_wafv2_rule_group` | judged as a planned AwsWafRuleGroup kind |
| `aws_xray_encryption_config` | judged as a planned AwsXraySettings kind (sampling rules, groups, indexing, encryption, trace destinations, resource policies) |
| `aws_xray_group` | judged as a planned AwsXraySettings kind (sampling rules, groups, indexing, encryption, trace destinations, resource policies) |
| `aws_xray_indexing_rule` | judged as a planned AwsXraySettings kind (sampling rules, groups, indexing, encryption, trace destinations, resource policies) |
| `aws_xray_resource_policy` | judged as a planned AwsXraySettings kind (sampling rules, groups, indexing, encryption, trace destinations, resource policies) |
| `aws_xray_sampling_rule` | judged as a planned AwsXraySettings kind (sampling rules, groups, indexing, encryption, trace destinations, resource policies) |
| `aws_xray_trace_segment_destination` | judged as a planned AwsXraySettings kind (sampling rules, groups, indexing, encryption, trace destinations, resource policies) |

### Deferred (532)

| Resource | Recorded reason |
|---|---|
| `aws_appfabric_app_authorization` | SaaS-connector vertical (AppFabric); deferred pending demand |
| `aws_appfabric_app_authorization_connection` | SaaS-connector vertical (AppFabric); deferred pending demand |
| `aws_appfabric_app_bundle` | SaaS-connector vertical (AppFabric); deferred pending demand |
| `aws_appfabric_ingestion` | SaaS-connector vertical (AppFabric); deferred pending demand |
| `aws_appfabric_ingestion_destination` | SaaS-connector vertical (AppFabric); deferred pending demand |
| `aws_appflow_connector_profile` | SaaS data-flow vertical (AppFlow); deferred pending demand |
| `aws_appflow_flow` | SaaS data-flow vertical (AppFlow); deferred pending demand |
| `aws_appintegrations_data_integration` | Amazon Connect app-integration surface; deferred pending demand |
| `aws_appintegrations_event_integration` | Amazon Connect app-integration surface; deferred pending demand |
| `aws_applicationinsights_application` | CloudWatch Application Insights enrollment; deferred pending demand |
| `aws_apprunner_deployment` | imperative deployment trigger, a poor declarative fit; deployments belong to the service lifecycle |
| `aws_appstream_directory_config` | application-streaming vertical (AppStream); deferred pending demand |
| `aws_appstream_fleet` | application-streaming vertical (AppStream); deferred pending demand |
| `aws_appstream_fleet_stack_association` | application-streaming vertical (AppStream); deferred pending demand |
| `aws_appstream_image_builder` | application-streaming vertical (AppStream); deferred pending demand |
| `aws_appstream_stack` | application-streaming vertical (AppStream); deferred pending demand |
| `aws_appstream_user` | application-streaming vertical (AppStream); deferred pending demand |
| `aws_appstream_user_stack_association` | application-streaming vertical (AppStream); deferred pending demand |
| `aws_arcregionswitch_plan` | Route 53 ARC region/zonal shift controls; deferred pending demand |
| `aws_arczonalshift_autoshift_observer_notification_status` | Route 53 ARC region/zonal shift controls; deferred pending demand |
| `aws_arczonalshift_zonal_autoshift_configuration` | Route 53 ARC region/zonal shift controls; deferred pending demand |
| `aws_athena_capacity_reservation` | Athena capacity reservations and saved/prepared queries are workload content, not infrastructure shape; deferred |
| `aws_athena_named_query` | Athena capacity reservations and saved/prepared queries are workload content, not infrastructure shape; deferred |
| `aws_athena_prepared_statement` | Athena capacity reservations and saved/prepared queries are workload content, not infrastructure shape; deferred |
| `aws_auditmanager_account_registration` | compliance-audit vertical (Audit Manager); deferred pending demand |
| `aws_auditmanager_assessment` | compliance-audit vertical (Audit Manager); deferred pending demand |
| `aws_auditmanager_assessment_delegation` | compliance-audit vertical (Audit Manager); deferred pending demand |
| `aws_auditmanager_assessment_report` | compliance-audit vertical (Audit Manager); deferred pending demand |
| `aws_auditmanager_control` | compliance-audit vertical (Audit Manager); deferred pending demand |
| `aws_auditmanager_framework` | compliance-audit vertical (Audit Manager); deferred pending demand |
| `aws_auditmanager_framework_share` | compliance-audit vertical (Audit Manager); deferred pending demand |
| `aws_auditmanager_organization_admin_account_registration` | compliance-audit vertical (Audit Manager); deferred pending demand |
| `aws_autoscaling_group_tag` | tags auto-scaling groups OWNED BY ANOTHER SERVICE (e.g. EKS managed node groups' ASGs); AwsAutoScalingGroup tags its own group natively with propagate-at-launch -- revisit from the EKS family if node-group tag-propagation demand appears |
| `aws_bedrock_evaluation_job` | model evaluation jobs are imperative run-once workloads, a poor declarative fit; deferred |
| `aws_billing_view` | billing views are console/reporting configuration; deferred pending demand |
| `aws_chatbot_slack_channel_configuration` | chat-channel notification integrations (AWS Chatbot); deferred pending demand |
| `aws_chatbot_teams_channel_configuration` | chat-channel notification integrations (AWS Chatbot); deferred pending demand |
| `aws_chimesdkmediapipelines_media_insights_pipeline_configuration` | Chime SDK media/voice vertical; deferred pending demand |
| `aws_chimesdkvoice_global_settings` | Chime SDK media/voice vertical; deferred pending demand |
| `aws_chimesdkvoice_sip_media_application` | Chime SDK media/voice vertical; deferred pending demand |
| `aws_chimesdkvoice_sip_rule` | Chime SDK media/voice vertical; deferred pending demand |
| `aws_chimesdkvoice_voice_profile_domain` | Chime SDK media/voice vertical; deferred pending demand |
| `aws_cleanrooms_collaboration` | Clean Rooms collaboration vertical; deferred pending demand |
| `aws_cleanrooms_configured_table` | Clean Rooms collaboration vertical; deferred pending demand |
| `aws_cleanrooms_membership` | Clean Rooms collaboration vertical; deferred pending demand |
| `aws_cloudcontrolapi_resource` | generic Cloud Control meta-resource; a typed catalog models resources directly rather than through an untyped escape hatch |
| `aws_cloudformation_stack` | driving CloudFormation stacks through the catalog's own IaC engines is meta-IaC recursion, contrary to catalog doctrine |
| `aws_cloudformation_stack_instances` | driving CloudFormation stacks through the catalog's own IaC engines is meta-IaC recursion, contrary to catalog doctrine |
| `aws_cloudformation_stack_set` | driving CloudFormation stacks through the catalog's own IaC engines is meta-IaC recursion, contrary to catalog doctrine |
| `aws_cloudformation_stack_set_instance` | driving CloudFormation stacks through the catalog's own IaC engines is meta-IaC recursion, contrary to catalog doctrine |
| `aws_cloudformation_type` | driving CloudFormation stacks through the catalog's own IaC engines is meta-IaC recursion, contrary to catalog doctrine |
| `aws_cloudfront_connection_group` | CloudFront multi-tenant distribution set; deferred pending demand |
| `aws_cloudfront_distribution_tenant` | CloudFront multi-tenant distribution set; deferred pending demand |
| `aws_cloudfront_multitenant_distribution` | CloudFront multi-tenant distribution set; deferred pending demand |
| `aws_cloudhsm_v2_cluster` | dedicated CloudHSM clusters; KMS custom key stores cover the mainstream need; deferred |
| `aws_cloudhsm_v2_hsm` | dedicated CloudHSM clusters; KMS custom key stores cover the mainstream need; deferred |
| `aws_cloudwatch_contributor_insight_rule` | Contributor Insights rules; deferred pending demand |
| `aws_cloudwatch_contributor_managed_insight_rule` | Contributor Insights rules; deferred pending demand |
| `aws_cloudwatch_event_endpoint` | EventBridge global endpoints (multi-region failover); deferred pending demand |
| `aws_cloudwatch_log_s3_table_integration_source` | attaches a data source to an S3 Tables log INTEGRATION ARN whose parent integration has no provider resource (console/API-created) -- not log-group-scoped despite the former blanket; deferred until the S3 Tables integration surface is modelable end to end |
| `aws_cloudwatch_log_storage_tier_policy` | regional account-level SINGLETON (one storage-tier policy per region, 6.58.0-new) -- account posture, not log-group configuration (the VPC block-public-access class); deferred pending demand |
| `aws_cloudwatch_otel_enrichment` | OpenTelemetry enrichment is newly introduced account telemetry plumbing; deferred pending maturity |
| `aws_cloudwatch_query_definition` | saved Logs Insights queries are operator content (the data-plane class, like Cognito users) -- not infrastructure shape; deferred pending demand |
| `aws_codebuild_report_group` | standalone report-export destination referenced from buildspec report sections (never project configuration); niche CI-analytics surface, deferred pending demand |
| `aws_codecatalyst_dev_environment` | CodeCatalyst vertical; deferred pending demand |
| `aws_codecatalyst_project` | CodeCatalyst vertical; deferred pending demand |
| `aws_codecatalyst_source_repository` | CodeCatalyst vertical; deferred pending demand |
| `aws_codeguruprofiler_profiling_group` | CodeGuru profiling/review services; deferred pending demand |
| `aws_codegurureviewer_repository_association` | CodeGuru profiling/review services; deferred pending demand |
| `aws_codepipeline_custom_action_type` | account-level custom action-type registration with an independent lifecycle (the AwsCodePipeline spec's recorded exclusion); pipeline actions reference custom types by owner/provider/version today -- a future kind on demand |
| `aws_codepipeline_webhook` | the legacy V1 webhook trigger mechanism (per-pipeline GitHub-OAuth/HMAC push endpoint); V2 pipelines use native CodeConnections triggers, which AwsCodePipeline models in full -- the spec records webhooks as deliberately excluded; revisit only on V1-pipeline demand |
| `aws_cognito_user` | Cognito user records are data-plane identity content, not infrastructure |
| `aws_cognito_user_in_group` | Cognito user records are data-plane identity content, not infrastructure |
| `aws_cognito_user_pool_ui_customization` | CSS/image customization for the CLASSIC hosted UI -- AWS's current path is managed login (v2) with its branding designer; the legacy skin surface has no demand signal (revisit on classic-hosted-UI demand) |
| `aws_comprehend_document_classifier` | Comprehend custom NLP models are imperative training artifacts; deferred |
| `aws_comprehend_entity_recognizer` | Comprehend custom NLP models are imperative training artifacts; deferred |
| `aws_computeoptimizer_enrollment_status` | Compute Optimizer is a read-mostly advisory service; deferred |
| `aws_computeoptimizer_recommendation_preferences` | Compute Optimizer is a read-mostly advisory service; deferred |
| `aws_connect_bot_association` | contact-center vertical (Amazon Connect); deferred pending demand |
| `aws_connect_contact_flow` | contact-center vertical (Amazon Connect); deferred pending demand |
| `aws_connect_contact_flow_module` | contact-center vertical (Amazon Connect); deferred pending demand |
| `aws_connect_hours_of_operation` | contact-center vertical (Amazon Connect); deferred pending demand |
| `aws_connect_instance` | contact-center vertical (Amazon Connect); deferred pending demand |
| `aws_connect_instance_storage_config` | contact-center vertical (Amazon Connect); deferred pending demand |
| `aws_connect_lambda_function_association` | contact-center vertical (Amazon Connect); deferred pending demand |
| `aws_connect_phone_number` | contact-center vertical (Amazon Connect); deferred pending demand |
| `aws_connect_phone_number_contact_flow_association` | contact-center vertical (Amazon Connect); deferred pending demand |
| `aws_connect_queue` | contact-center vertical (Amazon Connect); deferred pending demand |
| `aws_connect_quick_connect` | contact-center vertical (Amazon Connect); deferred pending demand |
| `aws_connect_routing_profile` | contact-center vertical (Amazon Connect); deferred pending demand |
| `aws_connect_security_profile` | contact-center vertical (Amazon Connect); deferred pending demand |
| `aws_connect_user` | contact-center vertical (Amazon Connect); deferred pending demand |
| `aws_connect_user_hierarchy_group` | contact-center vertical (Amazon Connect); deferred pending demand |
| `aws_connect_user_hierarchy_structure` | contact-center vertical (Amazon Connect); deferred pending demand |
| `aws_connect_vocabulary` | contact-center vertical (Amazon Connect); deferred pending demand |
| `aws_costoptimizationhub_enrollment_status` | Cost Optimization Hub is a read-mostly advisory service; deferred |
| `aws_costoptimizationhub_preferences` | Cost Optimization Hub is a read-mostly advisory service; deferred |
| `aws_customerprofiles_domain` | Connect Customer Profiles vertical; deferred pending demand |
| `aws_customerprofiles_profile` | Connect Customer Profiles vertical; deferred pending demand |
| `aws_dataexchange_data_set` | Data Exchange marketplace surface; deferred pending demand |
| `aws_dataexchange_event_action` | Data Exchange marketplace surface; deferred pending demand |
| `aws_dataexchange_revision` | Data Exchange marketplace surface; deferred pending demand |
| `aws_dataexchange_revision_assets` | Data Exchange marketplace surface; deferred pending demand |
| `aws_datazone_asset_type` | DataZone data-governance portal; deferred pending demand |
| `aws_datazone_domain` | DataZone data-governance portal; deferred pending demand |
| `aws_datazone_environment` | DataZone data-governance portal; deferred pending demand |
| `aws_datazone_environment_blueprint_configuration` | DataZone data-governance portal; deferred pending demand |
| `aws_datazone_environment_profile` | DataZone data-governance portal; deferred pending demand |
| `aws_datazone_form_type` | DataZone data-governance portal; deferred pending demand |
| `aws_datazone_glossary` | DataZone data-governance portal; deferred pending demand |
| `aws_datazone_glossary_term` | DataZone data-governance portal; deferred pending demand |
| `aws_datazone_project` | DataZone data-governance portal; deferred pending demand |
| `aws_datazone_user_profile` | DataZone data-governance portal; deferred pending demand |
| `aws_dax_cluster` | DAX caching layer; deferred pending demand |
| `aws_dax_parameter_group` | DAX caching layer; deferred pending demand |
| `aws_dax_subnet_group` | DAX caching layer; deferred pending demand |
| `aws_db_cluster_snapshot` | database snapshot operations are imperative backup actions, a poor declarative fit; deferred |
| `aws_db_snapshot` | database snapshot operations are imperative backup actions, a poor declarative fit; deferred |
| `aws_db_snapshot_copy` | database snapshot operations are imperative backup actions, a poor declarative fit; deferred |
| `aws_default_network_acl` | adopting AWS's implicit default resources into management is an anti-pattern; the catalog provisions explicit resources instead |
| `aws_default_route_table` | adopting AWS's implicit default resources into management is an anti-pattern; the catalog provisions explicit resources instead |
| `aws_default_security_group` | adopting AWS's implicit default resources into management is an anti-pattern; the catalog provisions explicit resources instead |
| `aws_default_subnet` | adopting AWS's implicit default resources into management is an anti-pattern; the catalog provisions explicit resources instead |
| `aws_default_vpc` | adopting AWS's implicit default resources into management is an anti-pattern; the catalog provisions explicit resources instead |
| `aws_default_vpc_dhcp_options` | adopting AWS's implicit default resources into management is an anti-pattern; the catalog provisions explicit resources instead |
| `aws_detective_graph` | Detective security-graph vertical; deferred pending demand |
| `aws_detective_invitation_accepter` | Detective security-graph vertical; deferred pending demand |
| `aws_detective_member` | Detective security-graph vertical; deferred pending demand |
| `aws_detective_organization_admin_account` | Detective security-graph vertical; deferred pending demand |
| `aws_detective_organization_configuration` | Detective security-graph vertical; deferred pending demand |
| `aws_devicefarm_device_pool` | Device Farm testing vertical; deferred pending demand |
| `aws_devicefarm_instance_profile` | Device Farm testing vertical; deferred pending demand |
| `aws_devicefarm_network_profile` | Device Farm testing vertical; deferred pending demand |
| `aws_devicefarm_project` | Device Farm testing vertical; deferred pending demand |
| `aws_devicefarm_test_grid_project` | Device Farm testing vertical; deferred pending demand |
| `aws_devicefarm_upload` | Device Farm testing vertical; deferred pending demand |
| `aws_devopsguru_event_sources_config` | DevOps Guru advisory service; deferred pending demand |
| `aws_devopsguru_notification_channel` | DevOps Guru advisory service; deferred pending demand |
| `aws_devopsguru_resource_collection` | DevOps Guru advisory service; deferred pending demand |
| `aws_devopsguru_service_integration` | DevOps Guru advisory service; deferred pending demand |
| `aws_docdb_cluster_snapshot` | database snapshot operations are imperative backup actions, a poor declarative fit; deferred |
| `aws_drs_replication_configuration_template` | Elastic Disaster Recovery vertical; deferred pending demand |
| `aws_dynamodb_table_export` | table exports are imperative one-shot operations, a poor declarative fit; deferred |
| `aws_dynamodb_table_item` | table items and tag helpers are data-plane content and imperative helpers, not infrastructure |
| `aws_dynamodb_tag` | table items and tag helpers are data-plane content and imperative helpers, not infrastructure |
| `aws_ec2_carrier_gateway` | Wavelength carrier gateways; deferred pending demand |
| `aws_ec2_fleet` | EC2 Fleet and dedicated hosts; deferred pending demand |
| `aws_ec2_host` | EC2 Fleet and dedicated hosts; deferred pending demand |
| `aws_ec2_instance_state` | imperative state/tag helpers, not durable infrastructure |
| `aws_ec2_local_gateway_route` | Outposts local-gateway surface; deferred pending demand |
| `aws_ec2_local_gateway_route_table` | Outposts local-gateway surface; deferred pending demand |
| `aws_ec2_local_gateway_route_table_virtual_interface_group_association` | Outposts local-gateway surface; deferred pending demand |
| `aws_ec2_local_gateway_route_table_vpc_association` | Outposts local-gateway surface; deferred pending demand |
| `aws_ec2_network_insights_access_scope` | Network Insights analyses are imperative diagnostics, a poor declarative fit; deferred |
| `aws_ec2_network_insights_analysis` | Network Insights analyses are imperative diagnostics, a poor declarative fit; deferred |
| `aws_ec2_network_insights_path` | Network Insights analyses are imperative diagnostics, a poor declarative fit; deferred |
| `aws_ec2_secondary_network` | EC2 secondary network/subnet surface; deferred pending demand |
| `aws_ec2_secondary_subnet` | EC2 secondary network/subnet surface; deferred pending demand |
| `aws_ec2_tag` | imperative state/tag helpers, not durable infrastructure |
| `aws_ec2_traffic_mirror_filter` | traffic-mirroring surface; deferred pending demand |
| `aws_ec2_traffic_mirror_filter_rule` | traffic-mirroring surface; deferred pending demand |
| `aws_ec2_traffic_mirror_session` | traffic-mirroring surface; deferred pending demand |
| `aws_ec2_traffic_mirror_target` | traffic-mirroring surface; deferred pending demand |
| `aws_ec2_transit_gateway_metering_policy` | Transit Gateway metering policy (provider 6.37.0) -- cross-account data-processing payer attribution for large multi-account estates; deferred pending demand signals |
| `aws_ec2_transit_gateway_metering_policy_entry` | per-rule entry of the deferred metering policy; folds into that surface whenever demand revives it |
| `aws_ecrpublic_repository` | public ECR gallery repositories; deferred pending demand |
| `aws_ecrpublic_repository_policy` | public ECR gallery repositories; deferred pending demand |
| `aws_ecs_tag` | imperative state/tag helpers, not durable infrastructure |
| `aws_elastic_beanstalk_application` | Elastic Beanstalk is strategically displaced by App Runner and ECS kinds; not deprecated by AWS, deliberately not offered |
| `aws_elastic_beanstalk_application_version` | Elastic Beanstalk is strategically displaced by App Runner and ECS kinds; not deprecated by AWS, deliberately not offered |
| `aws_elastic_beanstalk_configuration_template` | Elastic Beanstalk is strategically displaced by App Runner and ECS kinds; not deprecated by AWS, deliberately not offered |
| `aws_elastic_beanstalk_environment` | Elastic Beanstalk is strategically displaced by App Runner and ECS kinds; not deprecated by AWS, deliberately not offered |
| `aws_elasticache_reserved_cache_node` | reserved nodes are a billing purchase commitment, not infrastructure |
| `aws_finspace_kx_cluster` | FinSpace kdb capital-markets vertical; deferred pending demand |
| `aws_finspace_kx_database` | FinSpace kdb capital-markets vertical; deferred pending demand |
| `aws_finspace_kx_dataview` | FinSpace kdb capital-markets vertical; deferred pending demand |
| `aws_finspace_kx_environment` | FinSpace kdb capital-markets vertical; deferred pending demand |
| `aws_finspace_kx_scaling_group` | FinSpace kdb capital-markets vertical; deferred pending demand |
| `aws_finspace_kx_user` | FinSpace kdb capital-markets vertical; deferred pending demand |
| `aws_finspace_kx_volume` | FinSpace kdb capital-markets vertical; deferred pending demand |
| `aws_fis_experiment_template` | Fault Injection Service chaos experiments; deferred pending demand |
| `aws_fis_target_account_configuration` | Fault Injection Service chaos experiments; deferred pending demand |
| `aws_fms_admin_account` | Firewall Manager organization-scale policy administration; deferred pending demand |
| `aws_fms_policy` | Firewall Manager organization-scale policy administration; deferred pending demand |
| `aws_fms_resource_set` | Firewall Manager organization-scale policy administration; deferred pending demand |
| `aws_fsx_backup` | one-shot point-in-time backup of a file system or volume (only source ids are configurable) -- an imperative backup action, a poor declarative fit (the aws_db_snapshot class); the FSx kinds model AUTOMATIC backups (retention, start time, final-backup controls) in-spec |
| `aws_fsx_file_cache` | FSx File Cache (HPC burst cache); deferred pending demand |
| `aws_fsx_openzfs_snapshot` | one-shot point-in-time snapshot of an OpenZFS volume (name + volume id only) -- an imperative snapshot action, a poor declarative fit (the aws_db_snapshot class); unlike EBS there is no snapshot-management family to justify a kind |
| `aws_gamelift_alias` | game-server hosting vertical (GameLift); deferred pending demand |
| `aws_gamelift_build` | game-server hosting vertical (GameLift); deferred pending demand |
| `aws_gamelift_fleet` | game-server hosting vertical (GameLift); deferred pending demand |
| `aws_gamelift_game_server_group` | game-server hosting vertical (GameLift); deferred pending demand |
| `aws_gamelift_game_session_queue` | game-server hosting vertical (GameLift); deferred pending demand |
| `aws_gamelift_script` | game-server hosting vertical (GameLift); deferred pending demand |
| `aws_glacier_vault` | standalone Glacier vaults; S3 lifecycle transitions (covered by AwsS3Bucket) are the mainstream path; deferred |
| `aws_glacier_vault_lock` | standalone Glacier vaults; S3 lifecycle transitions (covered by AwsS3Bucket) are the mainstream path; deferred |
| `aws_globalaccelerator_custom_routing_accelerator` | custom-routing accelerators; deferred pending demand |
| `aws_globalaccelerator_custom_routing_endpoint_group` | custom-routing accelerators; deferred pending demand |
| `aws_globalaccelerator_custom_routing_listener` | custom-routing accelerators; deferred pending demand |
| `aws_glue_data_quality_ruleset` | Glue ML transforms, data-quality rulesets and UDFs are workload content; deferred |
| `aws_glue_ml_transform` | Glue ML transforms, data-quality rulesets and UDFs are workload content; deferred |
| `aws_glue_user_defined_function` | Glue ML transforms, data-quality rulesets and UDFs are workload content; deferred |
| `aws_iam_outbound_web_identity_federation` | human-credential MFA devices and outbound identity federation; deferred pending demand |
| `aws_iam_user_login_profile` | human console-login credentials (PGP-encrypted passwords) -- the human-credential class alongside aws_iam_virtual_mfa_device; the catalog provisions machine identities; revisit on demand |
| `aws_iam_user_policies_exclusive` | exclusive-set reconciliation (purges out-of-band inline policies at apply, no-op delete) -- engine-workflow surface, not resource configuration; AwsIamUser already declares its full intended policy set |
| `aws_iam_user_policy_attachments_exclusive` | exclusive-set reconciliation (purges out-of-band managed-policy attachments at apply, no-op delete) -- engine-workflow surface, not resource configuration; AwsIamUser already declares its full intended policy set |
| `aws_iam_user_ssh_key` | CodeCommit SSH public keys -- the consuming service closed to new customers in 2024; revisit on demand |
| `aws_iam_virtual_mfa_device` | human-credential MFA devices and outbound identity federation; deferred pending demand |
| `aws_internetmonitor_monitor` | Internet Monitor observability; deferred pending demand |
| `aws_invoicing_invoice_unit` | invoice-unit billing configuration; deferred pending demand |
| `aws_iot_authorizer` | IoT fleet-management vertical; deferred pending demand |
| `aws_iot_billing_group` | IoT fleet-management vertical; deferred pending demand |
| `aws_iot_ca_certificate` | IoT fleet-management vertical; deferred pending demand |
| `aws_iot_certificate` | IoT fleet-management vertical; deferred pending demand |
| `aws_iot_domain_configuration` | IoT fleet-management vertical; deferred pending demand |
| `aws_iot_event_configurations` | IoT fleet-management vertical; deferred pending demand |
| `aws_iot_indexing_configuration` | IoT fleet-management vertical; deferred pending demand |
| `aws_iot_logging_options` | IoT fleet-management vertical; deferred pending demand |
| `aws_iot_policy` | IoT fleet-management vertical; deferred pending demand |
| `aws_iot_policy_attachment` | IoT fleet-management vertical; deferred pending demand |
| `aws_iot_provisioning_template` | IoT fleet-management vertical; deferred pending demand |
| `aws_iot_role_alias` | IoT fleet-management vertical; deferred pending demand |
| `aws_iot_thing` | IoT fleet-management vertical; deferred pending demand |
| `aws_iot_thing_group` | IoT fleet-management vertical; deferred pending demand |
| `aws_iot_thing_group_membership` | IoT fleet-management vertical; deferred pending demand |
| `aws_iot_thing_principal_attachment` | IoT fleet-management vertical; deferred pending demand |
| `aws_iot_thing_type` | IoT fleet-management vertical; deferred pending demand |
| `aws_iot_topic_rule` | IoT fleet-management vertical; deferred pending demand |
| `aws_iot_topic_rule_destination` | IoT fleet-management vertical; deferred pending demand |
| `aws_ivs_channel` | Interactive Video Service vertical; deferred pending demand |
| `aws_ivs_playback_key_pair` | Interactive Video Service vertical; deferred pending demand |
| `aws_ivs_recording_configuration` | Interactive Video Service vertical; deferred pending demand |
| `aws_ivschat_logging_configuration` | Interactive Video Service vertical; deferred pending demand |
| `aws_ivschat_room` | Interactive Video Service vertical; deferred pending demand |
| `aws_kendra_data_source` | Kendra enterprise search; the Bedrock knowledge-base path covers retrieval-augmented generation; deferred |
| `aws_kendra_experience` | Kendra enterprise search; the Bedrock knowledge-base path covers retrieval-augmented generation; deferred |
| `aws_kendra_faq` | Kendra enterprise search; the Bedrock knowledge-base path covers retrieval-augmented generation; deferred |
| `aws_kendra_index` | Kendra enterprise search; the Bedrock knowledge-base path covers retrieval-augmented generation; deferred |
| `aws_kendra_query_suggestions_block_list` | Kendra enterprise search; the Bedrock knowledge-base path covers retrieval-augmented generation; deferred |
| `aws_kendra_thesaurus` | Kendra enterprise search; the Bedrock knowledge-base path covers retrieval-augmented generation; deferred |
| `aws_kinesis_video_stream` | Kinesis Video Streams vertical; deferred pending demand |
| `aws_kms_ciphertext` | ciphertext is an imperative encrypt operation, a poor declarative fit |
| `aws_kms_custom_key_store` | CloudHSM-backed custom key stores; deferred pending demand |
| `aws_lambda_invocation` | function invocation is an imperative operation, a poor declarative fit |
| `aws_lexv2models_bot` | conversational bot building (Lex v2); deferred pending demand |
| `aws_lexv2models_bot_locale` | conversational bot building (Lex v2); deferred pending demand |
| `aws_lexv2models_bot_version` | conversational bot building (Lex v2); deferred pending demand |
| `aws_lexv2models_intent` | conversational bot building (Lex v2); deferred pending demand |
| `aws_lexv2models_slot` | conversational bot building (Lex v2); deferred pending demand |
| `aws_lexv2models_slot_type` | conversational bot building (Lex v2); deferred pending demand |
| `aws_licensemanager_association` | License Manager tracking; deferred pending demand |
| `aws_licensemanager_grant` | License Manager tracking; deferred pending demand |
| `aws_licensemanager_grant_accepter` | License Manager tracking; deferred pending demand |
| `aws_licensemanager_license_configuration` | License Manager tracking; deferred pending demand |
| `aws_lightsail_bucket` | Lightsail is a simplified-VPS product; the catalog's kinds are themselves the simplification layer over first-class services |
| `aws_lightsail_bucket_access_key` | Lightsail is a simplified-VPS product; the catalog's kinds are themselves the simplification layer over first-class services |
| `aws_lightsail_bucket_resource_access` | Lightsail is a simplified-VPS product; the catalog's kinds are themselves the simplification layer over first-class services |
| `aws_lightsail_certificate` | Lightsail is a simplified-VPS product; the catalog's kinds are themselves the simplification layer over first-class services |
| `aws_lightsail_container_service` | Lightsail is a simplified-VPS product; the catalog's kinds are themselves the simplification layer over first-class services |
| `aws_lightsail_container_service_deployment_version` | Lightsail is a simplified-VPS product; the catalog's kinds are themselves the simplification layer over first-class services |
| `aws_lightsail_database` | Lightsail is a simplified-VPS product; the catalog's kinds are themselves the simplification layer over first-class services |
| `aws_lightsail_disk` | Lightsail is a simplified-VPS product; the catalog's kinds are themselves the simplification layer over first-class services |
| `aws_lightsail_disk_attachment` | Lightsail is a simplified-VPS product; the catalog's kinds are themselves the simplification layer over first-class services |
| `aws_lightsail_distribution` | Lightsail is a simplified-VPS product; the catalog's kinds are themselves the simplification layer over first-class services |
| `aws_lightsail_domain` | Lightsail is a simplified-VPS product; the catalog's kinds are themselves the simplification layer over first-class services |
| `aws_lightsail_domain_entry` | Lightsail is a simplified-VPS product; the catalog's kinds are themselves the simplification layer over first-class services |
| `aws_lightsail_instance` | Lightsail is a simplified-VPS product; the catalog's kinds are themselves the simplification layer over first-class services |
| `aws_lightsail_instance_public_ports` | Lightsail is a simplified-VPS product; the catalog's kinds are themselves the simplification layer over first-class services |
| `aws_lightsail_key_pair` | Lightsail is a simplified-VPS product; the catalog's kinds are themselves the simplification layer over first-class services |
| `aws_lightsail_lb` | Lightsail is a simplified-VPS product; the catalog's kinds are themselves the simplification layer over first-class services |
| `aws_lightsail_lb_attachment` | Lightsail is a simplified-VPS product; the catalog's kinds are themselves the simplification layer over first-class services |
| `aws_lightsail_lb_certificate` | Lightsail is a simplified-VPS product; the catalog's kinds are themselves the simplification layer over first-class services |
| `aws_lightsail_lb_certificate_attachment` | Lightsail is a simplified-VPS product; the catalog's kinds are themselves the simplification layer over first-class services |
| `aws_lightsail_lb_https_redirection_policy` | Lightsail is a simplified-VPS product; the catalog's kinds are themselves the simplification layer over first-class services |
| `aws_lightsail_lb_stickiness_policy` | Lightsail is a simplified-VPS product; the catalog's kinds are themselves the simplification layer over first-class services |
| `aws_lightsail_static_ip` | Lightsail is a simplified-VPS product; the catalog's kinds are themselves the simplification layer over first-class services |
| `aws_lightsail_static_ip_attachment` | Lightsail is a simplified-VPS product; the catalog's kinds are themselves the simplification layer over first-class services |
| `aws_location_geofence_collection` | Location Service maps/geofencing vertical; deferred pending demand |
| `aws_location_map` | Location Service maps/geofencing vertical; deferred pending demand |
| `aws_location_place_index` | Location Service maps/geofencing vertical; deferred pending demand |
| `aws_location_route_calculator` | Location Service maps/geofencing vertical; deferred pending demand |
| `aws_location_tracker` | Location Service maps/geofencing vertical; deferred pending demand |
| `aws_location_tracker_association` | Location Service maps/geofencing vertical; deferred pending demand |
| `aws_m2_application` | mainframe modernization vertical; deferred pending demand |
| `aws_m2_deployment` | mainframe modernization vertical; deferred pending demand |
| `aws_m2_environment` | mainframe modernization vertical; deferred pending demand |
| `aws_mailmanager_rule_set` | SES Mail Manager rule sets; deferred pending demand |
| `aws_mailmanager_traffic_policy` | SES Mail Manager rule sets; deferred pending demand |
| `aws_media_convert_queue` | broadcast media services vertical; deferred pending demand |
| `aws_media_package_channel` | broadcast media services vertical; deferred pending demand |
| `aws_media_packagev2_channel_group` | broadcast media services vertical; deferred pending demand |
| `aws_medialive_channel` | broadcast media services vertical; deferred pending demand |
| `aws_medialive_input` | broadcast media services vertical; deferred pending demand |
| `aws_medialive_input_security_group` | broadcast media services vertical; deferred pending demand |
| `aws_medialive_multiplex` | broadcast media services vertical; deferred pending demand |
| `aws_medialive_multiplex_program` | broadcast media services vertical; deferred pending demand |
| `aws_memorydb_snapshot` | database snapshot operations are imperative backup actions, a poor declarative fit; deferred |
| `aws_msk_replicator` | MSK cross-cluster replicators; deferred pending demand |
| `aws_neptune_cluster_snapshot` | database snapshot operations are imperative backup actions, a poor declarative fit; deferred |
| `aws_neptunegraph_graph` | Neptune Analytics graphs; deferred pending demand |
| `aws_networkflowmonitor_monitor` | network flow/path monitors; deferred pending demand |
| `aws_networkflowmonitor_scope` | network flow/path monitors; deferred pending demand |
| `aws_networkmanager_attachment_accepter` | Cloud WAN / global network manager is large-enterprise WAN topology; deferred pending demand |
| `aws_networkmanager_attachment_routing_policy_label` | Cloud WAN / global network manager is large-enterprise WAN topology; deferred pending demand |
| `aws_networkmanager_connect_attachment` | Cloud WAN / global network manager is large-enterprise WAN topology; deferred pending demand |
| `aws_networkmanager_connect_peer` | Cloud WAN / global network manager is large-enterprise WAN topology; deferred pending demand |
| `aws_networkmanager_connection` | Cloud WAN / global network manager is large-enterprise WAN topology; deferred pending demand |
| `aws_networkmanager_core_network` | Cloud WAN / global network manager is large-enterprise WAN topology; deferred pending demand |
| `aws_networkmanager_core_network_policy_attachment` | Cloud WAN / global network manager is large-enterprise WAN topology; deferred pending demand |
| `aws_networkmanager_customer_gateway_association` | Cloud WAN / global network manager is large-enterprise WAN topology; deferred pending demand |
| `aws_networkmanager_device` | Cloud WAN / global network manager is large-enterprise WAN topology; deferred pending demand |
| `aws_networkmanager_dx_gateway_attachment` | Cloud WAN / global network manager is large-enterprise WAN topology; deferred pending demand |
| `aws_networkmanager_global_network` | Cloud WAN / global network manager is large-enterprise WAN topology; deferred pending demand |
| `aws_networkmanager_link` | Cloud WAN / global network manager is large-enterprise WAN topology; deferred pending demand |
| `aws_networkmanager_link_association` | Cloud WAN / global network manager is large-enterprise WAN topology; deferred pending demand |
| `aws_networkmanager_prefix_list_association` | Cloud WAN / global network manager is large-enterprise WAN topology; deferred pending demand |
| `aws_networkmanager_site` | Cloud WAN / global network manager is large-enterprise WAN topology; deferred pending demand |
| `aws_networkmanager_site_to_site_vpn_attachment` | Cloud WAN / global network manager is large-enterprise WAN topology; deferred pending demand |
| `aws_networkmanager_transit_gateway_connect_peer_association` | Cloud WAN / global network manager is large-enterprise WAN topology; deferred pending demand |
| `aws_networkmanager_transit_gateway_peering` | Cloud WAN / global network manager is large-enterprise WAN topology; deferred pending demand |
| `aws_networkmanager_transit_gateway_registration` | Cloud WAN / global network manager is large-enterprise WAN topology; deferred pending demand |
| `aws_networkmanager_transit_gateway_route_table_attachment` | Cloud WAN / global network manager is large-enterprise WAN topology; deferred pending demand |
| `aws_networkmanager_vpc_attachment` | Cloud WAN / global network manager is large-enterprise WAN topology; deferred pending demand |
| `aws_networkmonitor_monitor` | network flow/path monitors; deferred pending demand |
| `aws_networkmonitor_probe` | network flow/path monitors; deferred pending demand |
| `aws_notifications_channel_association` | AWS User Notifications console configuration; deferred pending demand |
| `aws_notifications_event_rule` | AWS User Notifications console configuration; deferred pending demand |
| `aws_notifications_managed_notification_account_contact_association` | AWS User Notifications console configuration; deferred pending demand |
| `aws_notifications_managed_notification_additional_channel_association` | AWS User Notifications console configuration; deferred pending demand |
| `aws_notifications_notification_configuration` | AWS User Notifications console configuration; deferred pending demand |
| `aws_notifications_notification_hub` | AWS User Notifications console configuration; deferred pending demand |
| `aws_notifications_organizational_unit_association` | AWS User Notifications console configuration; deferred pending demand |
| `aws_notifications_organizations_access` | AWS User Notifications console configuration; deferred pending demand |
| `aws_notificationscontacts_email_contact` | AWS User Notifications console configuration; deferred pending demand |
| `aws_observabilityadmin_centralization_rule_for_organization` | organization telemetry administration is newly introduced; deferred pending maturity |
| `aws_observabilityadmin_s3_table_integration` | organization telemetry administration is newly introduced; deferred pending maturity |
| `aws_observabilityadmin_telemetry_enrichment` | organization telemetry administration is newly introduced; deferred pending maturity |
| `aws_observabilityadmin_telemetry_evaluation` | organization telemetry administration is newly introduced; deferred pending maturity |
| `aws_observabilityadmin_telemetry_evaluation_for_organization` | organization telemetry administration is newly introduced; deferred pending maturity |
| `aws_observabilityadmin_telemetry_pipeline` | organization telemetry administration is newly introduced; deferred pending maturity |
| `aws_observabilityadmin_telemetry_rule` | organization telemetry administration is newly introduced; deferred pending maturity |
| `aws_observabilityadmin_telemetry_rule_for_organization` | organization telemetry administration is newly introduced; deferred pending maturity |
| `aws_odb_cloud_autonomous_vm_cluster` | Oracle Database@AWS vertical; deferred pending demand |
| `aws_odb_cloud_exadata_infrastructure` | Oracle Database@AWS vertical; deferred pending demand |
| `aws_odb_cloud_vm_cluster` | Oracle Database@AWS vertical; deferred pending demand |
| `aws_odb_network` | Oracle Database@AWS vertical; deferred pending demand |
| `aws_odb_network_peering_connection` | Oracle Database@AWS vertical; deferred pending demand |
| `aws_opensearch_application` | OpenSearch UI applications; deferred pending demand |
| `aws_opensearch_inbound_connection_accepter` | cross-cluster search connections; deferred pending demand |
| `aws_opensearch_outbound_connection` | cross-cluster search connections; deferred pending demand |
| `aws_osis_pipeline` | OpenSearch Ingestion pipelines; deferred pending demand |
| `aws_osis_pipeline_endpoint` | OpenSearch Ingestion pipelines; deferred pending demand |
| `aws_osis_resource_policy` | OpenSearch Ingestion pipelines; deferred pending demand |
| `aws_outposts_capacity_task` | Outposts capacity operations; deferred pending demand |
| `aws_paymentcryptography_key` | payment-HSM cryptography vertical; deferred pending demand |
| `aws_paymentcryptography_key_alias` | payment-HSM cryptography vertical; deferred pending demand |
| `aws_pinpointsmsvoicev2_configuration_set` | End User Messaging SMS/voice channels; deferred pending demand |
| `aws_pinpointsmsvoicev2_event_destination` | End User Messaging SMS/voice channels; deferred pending demand |
| `aws_pinpointsmsvoicev2_opt_out_list` | End User Messaging SMS/voice channels; deferred pending demand |
| `aws_pinpointsmsvoicev2_phone_number` | End User Messaging SMS/voice channels; deferred pending demand |
| `aws_pinpointsmsvoicev2_pool` | End User Messaging SMS/voice channels; deferred pending demand |
| `aws_qbusiness_application` | Q Business assistant applications; deferred pending demand |
| `aws_quicksight_account_settings` | QuickSight dashboards/analyses/datasets are BI content authoring, not infrastructure; revisit on customer demand |
| `aws_quicksight_account_subscription` | QuickSight dashboards/analyses/datasets are BI content authoring, not infrastructure; revisit on customer demand |
| `aws_quicksight_analysis` | QuickSight dashboards/analyses/datasets are BI content authoring, not infrastructure; revisit on customer demand |
| `aws_quicksight_custom_permissions` | QuickSight dashboards/analyses/datasets are BI content authoring, not infrastructure; revisit on customer demand |
| `aws_quicksight_dashboard` | QuickSight dashboards/analyses/datasets are BI content authoring, not infrastructure; revisit on customer demand |
| `aws_quicksight_data_set` | QuickSight dashboards/analyses/datasets are BI content authoring, not infrastructure; revisit on customer demand |
| `aws_quicksight_data_source` | QuickSight dashboards/analyses/datasets are BI content authoring, not infrastructure; revisit on customer demand |
| `aws_quicksight_folder` | QuickSight dashboards/analyses/datasets are BI content authoring, not infrastructure; revisit on customer demand |
| `aws_quicksight_folder_membership` | QuickSight dashboards/analyses/datasets are BI content authoring, not infrastructure; revisit on customer demand |
| `aws_quicksight_group` | QuickSight dashboards/analyses/datasets are BI content authoring, not infrastructure; revisit on customer demand |
| `aws_quicksight_group_membership` | QuickSight dashboards/analyses/datasets are BI content authoring, not infrastructure; revisit on customer demand |
| `aws_quicksight_iam_policy_assignment` | QuickSight dashboards/analyses/datasets are BI content authoring, not infrastructure; revisit on customer demand |
| `aws_quicksight_ingestion` | QuickSight dashboards/analyses/datasets are BI content authoring, not infrastructure; revisit on customer demand |
| `aws_quicksight_ip_restriction` | QuickSight dashboards/analyses/datasets are BI content authoring, not infrastructure; revisit on customer demand |
| `aws_quicksight_key_registration` | QuickSight dashboards/analyses/datasets are BI content authoring, not infrastructure; revisit on customer demand |
| `aws_quicksight_namespace` | QuickSight dashboards/analyses/datasets are BI content authoring, not infrastructure; revisit on customer demand |
| `aws_quicksight_refresh_schedule` | QuickSight dashboards/analyses/datasets are BI content authoring, not infrastructure; revisit on customer demand |
| `aws_quicksight_role_custom_permission` | QuickSight dashboards/analyses/datasets are BI content authoring, not infrastructure; revisit on customer demand |
| `aws_quicksight_role_membership` | QuickSight dashboards/analyses/datasets are BI content authoring, not infrastructure; revisit on customer demand |
| `aws_quicksight_template` | QuickSight dashboards/analyses/datasets are BI content authoring, not infrastructure; revisit on customer demand |
| `aws_quicksight_template_alias` | QuickSight dashboards/analyses/datasets are BI content authoring, not infrastructure; revisit on customer demand |
| `aws_quicksight_theme` | QuickSight dashboards/analyses/datasets are BI content authoring, not infrastructure; revisit on customer demand |
| `aws_quicksight_user` | QuickSight dashboards/analyses/datasets are BI content authoring, not infrastructure; revisit on customer demand |
| `aws_quicksight_user_custom_permission` | QuickSight dashboards/analyses/datasets are BI content authoring, not infrastructure; revisit on customer demand |
| `aws_quicksight_vpc_connection` | QuickSight dashboards/analyses/datasets are BI content authoring, not infrastructure; revisit on customer demand |
| `aws_rbin_rule` | Recycle Bin retention rules; deferred pending demand |
| `aws_rds_certificate` | regional account singleton (sets the account's default CA certificate; its id IS the region) -- account-wide operational posture, never instance or cluster surface; per-database CA choice is already modeled as ca_cert_identifier |
| `aws_rds_cluster_snapshot_copy` | snapshot/export operations and custom engine versions are imperative lifecycle actions; deferred |
| `aws_rds_custom_db_engine_version` | snapshot/export operations and custom engine versions are imperative lifecycle actions; deferred |
| `aws_rds_export_task` | snapshot/export operations and custom engine versions are imperative lifecycle actions; deferred |
| `aws_rds_instance_state` | imperative instance start/stop state helper, not durable infrastructure |
| `aws_rds_reserved_instance` | reserved instances are a billing purchase commitment, not infrastructure |
| `aws_redshift_idc_application` | Identity Center application wiring and partner integrations; deferred pending demand |
| `aws_redshift_partner` | Identity Center application wiring and partner integrations; deferred pending demand |
| `aws_redshiftdata_statement` | imperative SQL statement execution, a poor declarative fit |
| `aws_redshiftserverless_snapshot` | database snapshot operations are imperative backup actions, a poor declarative fit; deferred |
| `aws_rekognition_collection` | computer-vision vertical (Rekognition); deferred pending demand |
| `aws_rekognition_project` | computer-vision vertical (Rekognition); deferred pending demand |
| `aws_rekognition_stream_processor` | computer-vision vertical (Rekognition); deferred pending demand |
| `aws_resiliencehub_resiliency_policy` | Resilience Hub assessment policies; deferred pending demand |
| `aws_resiliencehubv2_policy` | Resilience Hub assessment policies; deferred pending demand |
| `aws_resourceexplorer2_index` | Resource Explorer indexing; deferred pending demand |
| `aws_resourceexplorer2_view` | Resource Explorer indexing; deferred pending demand |
| `aws_resourcegroups_group` | resource-group organization helpers; deferred pending demand |
| `aws_resourcegroups_resource` | resource-group organization helpers; deferred pending demand |
| `aws_rolesanywhere_profile` | IAM Roles Anywhere certificate-based access; deferred pending demand |
| `aws_rolesanywhere_trust_anchor` | IAM Roles Anywhere certificate-based access; deferred pending demand |
| `aws_route53_records_exclusive` | exclusive-management lockdown of a zone's COMPLETE record set (declares the authoritative set, prunes anything undeclared, no-op delete); the catalog's per-record AwsRoute53DnsRecord model composes freely with out-of-band records -- the whole-zone lockdown is the exclusive-management class (the IAM exclusive-lockdown precedent), deferred pending demand |
| `aws_route53_traffic_policy` | Route 53 Traffic Flow policies are the legacy routing-policy path; deferred |
| `aws_route53_traffic_policy_instance` | Route 53 Traffic Flow policies are the legacy routing-policy path; deferred |
| `aws_route53profiles_association` | Route 53 Profiles DNS-config sharing; deferred pending demand |
| `aws_route53profiles_profile` | Route 53 Profiles DNS-config sharing; deferred pending demand |
| `aws_route53profiles_resource_association` | Route 53 Profiles DNS-config sharing; deferred pending demand |
| `aws_route53recoverycontrolconfig_cluster` | Route 53 Application Recovery Controller readiness/control surface; deferred pending demand |
| `aws_route53recoverycontrolconfig_control_panel` | Route 53 Application Recovery Controller readiness/control surface; deferred pending demand |
| `aws_route53recoverycontrolconfig_routing_control` | Route 53 Application Recovery Controller readiness/control surface; deferred pending demand |
| `aws_route53recoverycontrolconfig_safety_rule` | Route 53 Application Recovery Controller readiness/control surface; deferred pending demand |
| `aws_route53recoveryreadiness_cell` | Route 53 Application Recovery Controller readiness/control surface; deferred pending demand |
| `aws_route53recoveryreadiness_readiness_check` | Route 53 Application Recovery Controller readiness/control surface; deferred pending demand |
| `aws_route53recoveryreadiness_recovery_group` | Route 53 Application Recovery Controller readiness/control surface; deferred pending demand |
| `aws_route53recoveryreadiness_resource_set` | Route 53 Application Recovery Controller readiness/control surface; deferred pending demand |
| `aws_rum_app_monitor` | CloudWatch RUM monitoring; deferred pending demand |
| `aws_rum_metrics_destination` | CloudWatch RUM monitoring; deferred pending demand |
| `aws_s3control_bucket` | S3 on Outposts buckets; deferred pending demand |
| `aws_s3control_bucket_lifecycle_configuration` | S3 on Outposts buckets; deferred pending demand |
| `aws_s3control_bucket_policy` | S3 on Outposts buckets; deferred pending demand |
| `aws_s3control_object_lambda_access_point` | S3 Object Lambda access points; deferred pending demand |
| `aws_s3control_object_lambda_access_point_policy` | S3 Object Lambda access points; deferred pending demand |
| `aws_s3files_access_point` | the new S3 file-system service; deferred pending maturity and demand |
| `aws_s3files_file_system` | the new S3 file-system service; deferred pending maturity and demand |
| `aws_s3files_file_system_policy` | the new S3 file-system service; deferred pending maturity and demand |
| `aws_s3files_mount_target` | the new S3 file-system service; deferred pending maturity and demand |
| `aws_s3files_synchronization_configuration` | the new S3 file-system service; deferred pending maturity and demand |
| `aws_s3outposts_endpoint` | S3 on Outposts endpoints; deferred pending demand |
| `aws_sagemaker_algorithm` | SageMaker workload content (training/tuning/labeling jobs, monitoring schedules, model cards, workforces, devices, hubs, algorithms, code repositories, flow definitions, projects) is imperative or data-plane surface; deferred |
| `aws_sagemaker_code_repository` | SageMaker workload content (training/tuning/labeling jobs, monitoring schedules, model cards, workforces, devices, hubs, algorithms, code repositories, flow definitions, projects) is imperative or data-plane surface; deferred |
| `aws_sagemaker_data_quality_job_definition` | SageMaker workload content (training/tuning/labeling jobs, monitoring schedules, model cards, workforces, devices, hubs, algorithms, code repositories, flow definitions, projects) is imperative or data-plane surface; deferred |
| `aws_sagemaker_device` | SageMaker workload content (training/tuning/labeling jobs, monitoring schedules, model cards, workforces, devices, hubs, algorithms, code repositories, flow definitions, projects) is imperative or data-plane surface; deferred |
| `aws_sagemaker_device_fleet` | SageMaker workload content (training/tuning/labeling jobs, monitoring schedules, model cards, workforces, devices, hubs, algorithms, code repositories, flow definitions, projects) is imperative or data-plane surface; deferred |
| `aws_sagemaker_flow_definition` | SageMaker workload content (training/tuning/labeling jobs, monitoring schedules, model cards, workforces, devices, hubs, algorithms, code repositories, flow definitions, projects) is imperative or data-plane surface; deferred |
| `aws_sagemaker_hub` | SageMaker workload content (training/tuning/labeling jobs, monitoring schedules, model cards, workforces, devices, hubs, algorithms, code repositories, flow definitions, projects) is imperative or data-plane surface; deferred |
| `aws_sagemaker_hub_content_reference` | SageMaker workload content (training/tuning/labeling jobs, monitoring schedules, model cards, workforces, devices, hubs, algorithms, code repositories, flow definitions, projects) is imperative or data-plane surface; deferred |
| `aws_sagemaker_human_task_ui` | SageMaker workload content (training/tuning/labeling jobs, monitoring schedules, model cards, workforces, devices, hubs, algorithms, code repositories, flow definitions, projects) is imperative or data-plane surface; deferred |
| `aws_sagemaker_hyper_parameter_tuning_job` | SageMaker workload content (training/tuning/labeling jobs, monitoring schedules, model cards, workforces, devices, hubs, algorithms, code repositories, flow definitions, projects) is imperative or data-plane surface; deferred |
| `aws_sagemaker_labeling_job` | SageMaker workload content (training/tuning/labeling jobs, monitoring schedules, model cards, workforces, devices, hubs, algorithms, code repositories, flow definitions, projects) is imperative or data-plane surface; deferred |
| `aws_sagemaker_model_card` | SageMaker workload content (training/tuning/labeling jobs, monitoring schedules, model cards, workforces, devices, hubs, algorithms, code repositories, flow definitions, projects) is imperative or data-plane surface; deferred |
| `aws_sagemaker_model_card_export_job` | SageMaker workload content (training/tuning/labeling jobs, monitoring schedules, model cards, workforces, devices, hubs, algorithms, code repositories, flow definitions, projects) is imperative or data-plane surface; deferred |
| `aws_sagemaker_monitoring_schedule` | SageMaker workload content (training/tuning/labeling jobs, monitoring schedules, model cards, workforces, devices, hubs, algorithms, code repositories, flow definitions, projects) is imperative or data-plane surface; deferred |
| `aws_sagemaker_project` | SageMaker workload content (training/tuning/labeling jobs, monitoring schedules, model cards, workforces, devices, hubs, algorithms, code repositories, flow definitions, projects) is imperative or data-plane surface; deferred |
| `aws_sagemaker_training_job` | SageMaker workload content (training/tuning/labeling jobs, monitoring schedules, model cards, workforces, devices, hubs, algorithms, code repositories, flow definitions, projects) is imperative or data-plane surface; deferred |
| `aws_sagemaker_workforce` | SageMaker workload content (training/tuning/labeling jobs, monitoring schedules, model cards, workforces, devices, hubs, algorithms, code repositories, flow definitions, projects) is imperative or data-plane surface; deferred |
| `aws_sagemaker_workteam` | SageMaker workload content (training/tuning/labeling jobs, monitoring schedules, model cards, workforces, devices, hubs, algorithms, code repositories, flow definitions, projects) is imperative or data-plane surface; deferred |
| `aws_savingsplans_savings_plan` | savings plans are a billing purchase commitment, not infrastructure |
| `aws_schemas_discoverer` | EventBridge schema registry content; deferred pending demand |
| `aws_schemas_registry` | EventBridge schema registry content; deferred pending demand |
| `aws_schemas_registry_policy` | EventBridge schema registry content; deferred pending demand |
| `aws_schemas_schema` | EventBridge schema registry content; deferred pending demand |
| `aws_securitylake_aws_log_source` | Security Lake data-lake vertical; deferred pending demand |
| `aws_securitylake_custom_log_source` | Security Lake data-lake vertical; deferred pending demand |
| `aws_securitylake_data_lake` | Security Lake data-lake vertical; deferred pending demand |
| `aws_securitylake_subscriber` | Security Lake data-lake vertical; deferred pending demand |
| `aws_securitylake_subscriber_notification` | Security Lake data-lake vertical; deferred pending demand |
| `aws_serverlessapplicationrepository_cloudformation_stack` | Serverless Application Repository deployments are CloudFormation-mediated; contrary to catalog doctrine; deferred |
| `aws_servicecatalogappregistry_application` | AppRegistry application metadata; deferred pending demand |
| `aws_servicecatalogappregistry_attribute_group` | AppRegistry application metadata; deferred pending demand |
| `aws_servicecatalogappregistry_attribute_group_association` | AppRegistry application metadata; deferred pending demand |
| `aws_sesv2_contact_list` | SES tenants and contact lists; deferred pending demand |
| `aws_sesv2_tenant` | SES tenants and contact lists; deferred pending demand |
| `aws_sesv2_tenant_resource_association` | SES tenants and contact lists; deferred pending demand |
| `aws_signer_signing_job` | code-signing profiles and jobs; deferred pending demand |
| `aws_signer_signing_profile` | code-signing profiles and jobs; deferred pending demand |
| `aws_signer_signing_profile_permission` | code-signing profiles and jobs; deferred pending demand |
| `aws_sns_platform_application` | mobile-push platform applications and SMS preferences; deferred pending demand |
| `aws_sns_sms_preferences` | mobile-push platform applications and SMS preferences; deferred pending demand |
| `aws_ssmcontacts_contact` | incident-management and quick-setup surface; deferred pending demand |
| `aws_ssmcontacts_contact_channel` | incident-management and quick-setup surface; deferred pending demand |
| `aws_ssmcontacts_plan` | incident-management and quick-setup surface; deferred pending demand |
| `aws_ssmcontacts_rotation` | incident-management and quick-setup surface; deferred pending demand |
| `aws_ssmincidents_replication_set` | incident-management and quick-setup surface; deferred pending demand |
| `aws_ssmincidents_response_plan` | incident-management and quick-setup surface; deferred pending demand |
| `aws_ssmquicksetup_configuration_manager` | incident-management and quick-setup surface; deferred pending demand |
| `aws_timestreaminfluxdb_db_cluster` | time-series database vertical (Timestream); deferred pending demand |
| `aws_timestreaminfluxdb_db_instance` | time-series database vertical (Timestream); deferred pending demand |
| `aws_timestreamquery_scheduled_query` | time-series database vertical (Timestream); deferred pending demand |
| `aws_timestreamwrite_database` | time-series database vertical (Timestream); deferred pending demand |
| `aws_timestreamwrite_table` | time-series database vertical (Timestream); deferred pending demand |
| `aws_transcribe_language_model` | speech-to-text vocabularies and models; deferred pending demand |
| `aws_transcribe_medical_vocabulary` | speech-to-text vocabularies and models; deferred pending demand |
| `aws_transcribe_vocabulary` | speech-to-text vocabularies and models; deferred pending demand |
| `aws_transcribe_vocabulary_filter` | speech-to-text vocabularies and models; deferred pending demand |
| `aws_uxc_account_customizations` | account console customization; deferred pending demand |
| `aws_verifiedpermissions_identity_source` | fine-grained authorization stores (Verified Permissions); deferred pending demand |
| `aws_verifiedpermissions_policy` | fine-grained authorization stores (Verified Permissions); deferred pending demand |
| `aws_verifiedpermissions_policy_store` | fine-grained authorization stores (Verified Permissions); deferred pending demand |
| `aws_verifiedpermissions_policy_template` | fine-grained authorization stores (Verified Permissions); deferred pending demand |
| `aws_verifiedpermissions_schema` | fine-grained authorization stores (Verified Permissions); deferred pending demand |
| `aws_vpc_network_performance_metric_subscription` | account-scoped AZ-pair network performance monitoring (aggregate-latency metric subscriptions between availability zones) -- never per-VPC surface; revisit on demand |
| `aws_vpc_route_server` | VPC route servers; deferred pending demand |
| `aws_vpc_route_server_endpoint` | VPC route servers; deferred pending demand |
| `aws_vpc_route_server_peer` | VPC route servers; deferred pending demand |
| `aws_vpc_route_server_propagation` | VPC route servers; deferred pending demand |
| `aws_vpc_route_server_vpc_association` | VPC route servers; deferred pending demand |
| `aws_workmail_default_domain` | hosted email vertical (WorkMail); deferred pending demand |
| `aws_workmail_domain` | hosted email vertical (WorkMail); deferred pending demand |
| `aws_workmail_group` | hosted email vertical (WorkMail); deferred pending demand |
| `aws_workmail_organization` | hosted email vertical (WorkMail); deferred pending demand |
| `aws_workmail_user` | hosted email vertical (WorkMail); deferred pending demand |
| `aws_workspaces_connection_alias` | desktop streaming vertical (WorkSpaces); deferred pending demand |
| `aws_workspaces_directory` | desktop streaming vertical (WorkSpaces); deferred pending demand |
| `aws_workspaces_ip_group` | desktop streaming vertical (WorkSpaces); deferred pending demand |
| `aws_workspaces_pool` | desktop streaming vertical (WorkSpaces); deferred pending demand |
| `aws_workspaces_workspace` | desktop streaming vertical (WorkSpaces); deferred pending demand |
| `aws_workspacesweb_browser_settings` | secure-browser portal vertical (WorkSpaces Web); deferred pending demand |
| `aws_workspacesweb_browser_settings_association` | secure-browser portal vertical (WorkSpaces Web); deferred pending demand |
| `aws_workspacesweb_data_protection_settings` | secure-browser portal vertical (WorkSpaces Web); deferred pending demand |
| `aws_workspacesweb_data_protection_settings_association` | secure-browser portal vertical (WorkSpaces Web); deferred pending demand |
| `aws_workspacesweb_identity_provider` | secure-browser portal vertical (WorkSpaces Web); deferred pending demand |
| `aws_workspacesweb_ip_access_settings` | secure-browser portal vertical (WorkSpaces Web); deferred pending demand |
| `aws_workspacesweb_ip_access_settings_association` | secure-browser portal vertical (WorkSpaces Web); deferred pending demand |
| `aws_workspacesweb_network_settings` | secure-browser portal vertical (WorkSpaces Web); deferred pending demand |
| `aws_workspacesweb_network_settings_association` | secure-browser portal vertical (WorkSpaces Web); deferred pending demand |
| `aws_workspacesweb_portal` | secure-browser portal vertical (WorkSpaces Web); deferred pending demand |
| `aws_workspacesweb_session_logger` | secure-browser portal vertical (WorkSpaces Web); deferred pending demand |
| `aws_workspacesweb_session_logger_association` | secure-browser portal vertical (WorkSpaces Web); deferred pending demand |
| `aws_workspacesweb_trust_store` | secure-browser portal vertical (WorkSpaces Web); deferred pending demand |
| `aws_workspacesweb_trust_store_association` | secure-browser portal vertical (WorkSpaces Web); deferred pending demand |
| `aws_workspacesweb_user_access_logging_settings` | secure-browser portal vertical (WorkSpaces Web); deferred pending demand |
| `aws_workspacesweb_user_access_logging_settings_association` | secure-browser portal vertical (WorkSpaces Web); deferred pending demand |
| `aws_workspacesweb_user_settings` | secure-browser portal vertical (WorkSpaces Web); deferred pending demand |
| `aws_workspacesweb_user_settings_association` | secure-browser portal vertical (WorkSpaces Web); deferred pending demand |

### Excluded as deprecated (129)

| Resource | Recorded reason |
|---|---|
| `aws_alb` | legacy aws_alb* alias of the aws_lb* resource set already consumed by the load-balancer kinds |
| `aws_alb_listener` | legacy aws_alb* alias of the aws_lb* resource set already consumed by the load-balancer kinds |
| `aws_alb_listener_certificate` | legacy aws_alb* alias of the aws_lb* resource set already consumed by the load-balancer kinds |
| `aws_alb_listener_rule` | legacy aws_alb* alias of the aws_lb* resource set already consumed by the load-balancer kinds |
| `aws_alb_target_group` | legacy aws_alb* alias of the aws_lb* resource set already consumed by the load-balancer kinds |
| `aws_alb_target_group_attachment` | legacy aws_alb* alias of the aws_lb* resource set already consumed by the load-balancer kinds |
| `aws_app_cookie_stickiness_policy` | Classic Load Balancer surface, superseded by ALB/NLB (covered by AwsAlb and AwsNlb) |
| `aws_appmesh_gateway_route` | App Mesh end-of-support announced by AWS (September 2026); superseded by ECS Service Connect and VPC Lattice |
| `aws_appmesh_mesh` | App Mesh end-of-support announced by AWS (September 2026); superseded by ECS Service Connect and VPC Lattice |
| `aws_appmesh_route` | App Mesh end-of-support announced by AWS (September 2026); superseded by ECS Service Connect and VPC Lattice |
| `aws_appmesh_virtual_gateway` | App Mesh end-of-support announced by AWS (September 2026); superseded by ECS Service Connect and VPC Lattice |
| `aws_appmesh_virtual_node` | App Mesh end-of-support announced by AWS (September 2026); superseded by ECS Service Connect and VPC Lattice |
| `aws_appmesh_virtual_router` | App Mesh end-of-support announced by AWS (September 2026); superseded by ECS Service Connect and VPC Lattice |
| `aws_appmesh_virtual_service` | App Mesh end-of-support announced by AWS (September 2026); superseded by ECS Service Connect and VPC Lattice |
| `aws_athena_database` | superseded by the Glue catalog database, covered by AwsGlueCatalogDatabase |
| `aws_autoscalingplans_scaling_plan` | AutoScaling Plans is the legacy predictive-scaling path, superseded by ASG predictive scaling policies |
| `aws_bedrockagentcore_registry` | deprecated in the provider schema |
| `aws_chime_voice_connector` | Amazon Chime service retired by AWS (2026); voice-connector successors live in the Chime SDK voice family |
| `aws_chime_voice_connector_group` | Amazon Chime service retired by AWS (2026); voice-connector successors live in the Chime SDK voice family |
| `aws_chime_voice_connector_logging` | Amazon Chime service retired by AWS (2026); voice-connector successors live in the Chime SDK voice family |
| `aws_chime_voice_connector_origination` | Amazon Chime service retired by AWS (2026); voice-connector successors live in the Chime SDK voice family |
| `aws_chime_voice_connector_streaming` | Amazon Chime service retired by AWS (2026); voice-connector successors live in the Chime SDK voice family |
| `aws_chime_voice_connector_termination` | Amazon Chime service retired by AWS (2026); voice-connector successors live in the Chime SDK voice family |
| `aws_chime_voice_connector_termination_credentials` | Amazon Chime service retired by AWS (2026); voice-connector successors live in the Chime SDK voice family |
| `aws_cloud9_environment_ec2` | Cloud9 is closed to new customers by AWS |
| `aws_cloud9_environment_membership` | Cloud9 is closed to new customers by AWS |
| `aws_cloudfront_origin_access_identity` | superseded by origin access control, already consumed by AwsCloudFront |
| `aws_cloudsearch_domain` | CloudSearch is legacy, superseded by OpenSearch (covered by AwsOpenSearchDomain) |
| `aws_cloudsearch_domain_service_access_policy` | CloudSearch is legacy, superseded by OpenSearch (covered by AwsOpenSearchDomain) |
| `aws_codecommit_approval_rule_template` | CodeCommit is closed to new customers by AWS (2024) |
| `aws_codecommit_approval_rule_template_association` | CodeCommit is closed to new customers by AWS (2024) |
| `aws_codecommit_repository` | CodeCommit is closed to new customers by AWS (2024) |
| `aws_codecommit_trigger` | CodeCommit is closed to new customers by AWS (2024) |
| `aws_codestarconnections_connection` | codestarconnections is the superseded name of codeconnections (planned AwsCodeConnection kind) |
| `aws_codestarconnections_host` | codestarconnections is the superseded name of codeconnections (planned AwsCodeConnection kind) |
| `aws_codestarnotifications_notification_rule` | CodeStar is retired by AWS; its notification rules went with it |
| `aws_datapipeline_pipeline` | Data Pipeline is closed to new customers by AWS; superseded by Glue/MWAA orchestration |
| `aws_datapipeline_pipeline_definition` | Data Pipeline is closed to new customers by AWS; superseded by Glue/MWAA orchestration |
| `aws_dynamodb_global_table` | the 2017 global-table resource is superseded by table replicas (folds into AwsDynamodb) |
| `aws_elasticsearch_domain` | the legacy Elasticsearch service is superseded by OpenSearch (covered by AwsOpenSearchDomain) |
| `aws_elasticsearch_domain_policy` | the legacy Elasticsearch service is superseded by OpenSearch (covered by AwsOpenSearchDomain) |
| `aws_elasticsearch_domain_saml_options` | the legacy Elasticsearch service is superseded by OpenSearch (covered by AwsOpenSearchDomain) |
| `aws_elasticsearch_vpc_endpoint` | the legacy Elasticsearch service is superseded by OpenSearch (covered by AwsOpenSearchDomain) |
| `aws_elastictranscoder_pipeline` | deprecated in the provider schema |
| `aws_elastictranscoder_preset` | deprecated in the provider schema |
| `aws_elb` | Classic Load Balancer surface, superseded by ALB/NLB (covered by AwsAlb and AwsNlb) |
| `aws_elb_attachment` | Classic Load Balancer surface, superseded by ALB/NLB (covered by AwsAlb and AwsNlb) |
| `aws_evidently_feature` | deprecated in the provider schema |
| `aws_evidently_launch` | deprecated in the provider schema |
| `aws_evidently_project` | deprecated in the provider schema |
| `aws_evidently_segment` | deprecated in the provider schema |
| `aws_glue_dev_endpoint` | Glue dev endpoints are deprecated in favor of interactive sessions |
| `aws_iam_policy_attachment` | the whole-list attachment resource is a documented footgun, superseded by per-principal attachments (consumed) |
| `aws_iam_server_certificate` | IAM-hosted certificates are legacy, superseded by ACM (covered by AwsCertManagerCert) |
| `aws_iam_signing_certificate` | IAM-hosted certificates are legacy, superseded by ACM (covered by AwsCertManagerCert) |
| `aws_inspector_assessment_target` | Inspector Classic is superseded by Inspector v2 (planned AwsInspector kind) |
| `aws_inspector_assessment_template` | Inspector Classic is superseded by Inspector v2 (planned AwsInspector kind) |
| `aws_inspector_resource_group` | Inspector Classic is superseded by Inspector v2 (planned AwsInspector kind) |
| `aws_kinesis_analytics_application` | deprecated in the provider schema |
| `aws_launch_configuration` | launch configurations are deprecated by AWS; launch templates (AwsLaunchTemplate) are the successor |
| `aws_lb_cookie_stickiness_policy` | Classic Load Balancer surface, superseded by ALB/NLB (covered by AwsAlb and AwsNlb) |
| `aws_lb_ssl_negotiation_policy` | Classic Load Balancer surface, superseded by ALB/NLB (covered by AwsAlb and AwsNlb) |
| `aws_lex_bot` | Lex v1 is superseded by Lex v2 (deferred as a bot-building vertical) |
| `aws_lex_bot_alias` | Lex v1 is superseded by Lex v2 (deferred as a bot-building vertical) |
| `aws_lex_intent` | Lex v1 is superseded by Lex v2 (deferred as a bot-building vertical) |
| `aws_lex_slot_type` | Lex v1 is superseded by Lex v2 (deferred as a bot-building vertical) |
| `aws_load_balancer_backend_server_policy` | Classic Load Balancer surface, superseded by ALB/NLB (covered by AwsAlb and AwsNlb) |
| `aws_load_balancer_listener_policy` | Classic Load Balancer surface, superseded by ALB/NLB (covered by AwsAlb and AwsNlb) |
| `aws_load_balancer_policy` | Classic Load Balancer surface, superseded by ALB/NLB (covered by AwsAlb and AwsNlb) |
| `aws_media_store_container` | deprecated in the provider schema |
| `aws_media_store_container_policy` | deprecated in the provider schema |
| `aws_msk_scram_secret_association` | the whole-list SCRAM association is superseded by the consumed single-secret association resource |
| `aws_pinpoint_adm_channel` | Pinpoint engagement end-of-support announced by AWS (October 2026); superseded by SES and End User Messaging |
| `aws_pinpoint_apns_channel` | Pinpoint engagement end-of-support announced by AWS (October 2026); superseded by SES and End User Messaging |
| `aws_pinpoint_apns_sandbox_channel` | Pinpoint engagement end-of-support announced by AWS (October 2026); superseded by SES and End User Messaging |
| `aws_pinpoint_apns_voip_channel` | Pinpoint engagement end-of-support announced by AWS (October 2026); superseded by SES and End User Messaging |
| `aws_pinpoint_apns_voip_sandbox_channel` | Pinpoint engagement end-of-support announced by AWS (October 2026); superseded by SES and End User Messaging |
| `aws_pinpoint_app` | Pinpoint engagement end-of-support announced by AWS (October 2026); superseded by SES and End User Messaging |
| `aws_pinpoint_baidu_channel` | Pinpoint engagement end-of-support announced by AWS (October 2026); superseded by SES and End User Messaging |
| `aws_pinpoint_email_channel` | deprecated in the provider schema |
| `aws_pinpoint_email_template` | deprecated in the provider schema |
| `aws_pinpoint_event_stream` | deprecated in the provider schema |
| `aws_pinpoint_gcm_channel` | Pinpoint engagement end-of-support announced by AWS (October 2026); superseded by SES and End User Messaging |
| `aws_pinpoint_sms_channel` | Pinpoint engagement end-of-support announced by AWS (October 2026); superseded by SES and End User Messaging |
| `aws_proxy_protocol_policy` | Classic Load Balancer surface, superseded by ALB/NLB (covered by AwsAlb and AwsNlb) |
| `aws_qldb_ledger` | QLDB is deprecated by AWS (2025) |
| `aws_qldb_stream` | QLDB is deprecated by AWS (2025) |
| `aws_rds_shard_group` | Aurora Limitless shard group -- AWS discontinued Aurora Limitless Database (the cluster kind excludes cluster_scalability_type on the same judgment); superseded surface not flagged in the schema |
| `aws_redshift_hsm_client_certificate` | CloudHSM Classic surface is legacy; superseded by KMS-based encryption (covered by AwsRedshiftCluster) |
| `aws_redshift_hsm_configuration` | CloudHSM Classic surface is legacy; superseded by KMS-based encryption (covered by AwsRedshiftCluster) |
| `aws_s3_bucket_object` | deprecated in the provider schema |
| `aws_ses_configuration_set` | SES v1 identity/configuration surface is superseded by the SESv2 resources consumed by the SES kinds |
| `aws_ses_domain_dkim` | SES v1 identity/configuration surface is superseded by the SESv2 resources consumed by the SES kinds |
| `aws_ses_domain_identity` | SES v1 identity/configuration surface is superseded by the SESv2 resources consumed by the SES kinds |
| `aws_ses_domain_identity_verification` | SES v1 identity/configuration surface is superseded by the SESv2 resources consumed by the SES kinds |
| `aws_ses_domain_mail_from` | SES v1 identity/configuration surface is superseded by the SESv2 resources consumed by the SES kinds |
| `aws_ses_email_identity` | SES v1 identity/configuration surface is superseded by the SESv2 resources consumed by the SES kinds |
| `aws_ses_event_destination` | SES v1 identity/configuration surface is superseded by the SESv2 resources consumed by the SES kinds |
| `aws_ses_identity_notification_topic` | SES v1 identity/configuration surface is superseded by the SESv2 resources consumed by the SES kinds |
| `aws_ses_identity_policy` | SES v1 identity/configuration surface is superseded by the SESv2 resources consumed by the SES kinds |
| `aws_ses_template` | SES v1 identity/configuration surface is superseded by the SESv2 resources consumed by the SES kinds |
| `aws_spot_fleet_request` | spot request resources are superseded by launch-template market options and EC2 Fleet |
| `aws_spot_instance_request` | superseded surface by AWS's own guidance (the provider carries NO deprecation marker at 6.58.0 -- recorded 2026-08-11): AWS directs Spot workloads to launch-template/instance market options (AwsEc2Instance.market_type + spot_options, AwsLaunchTemplate) and EC2 Fleet; the standalone request resource duplicates aws_instance's whole schema with everything ForceNew |
| `aws_swf_domain` | SWF is superseded by Step Functions (covered by AwsStepFunction) |
| `aws_waf_byte_match_set` | WAF Classic (global) is superseded by WAFv2 (covered by AwsWafWebAcl) |
| `aws_waf_geo_match_set` | WAF Classic (global) is superseded by WAFv2 (covered by AwsWafWebAcl) |
| `aws_waf_ipset` | WAF Classic (global) is superseded by WAFv2 (covered by AwsWafWebAcl) |
| `aws_waf_rate_based_rule` | WAF Classic (global) is superseded by WAFv2 (covered by AwsWafWebAcl) |
| `aws_waf_regex_match_set` | WAF Classic (global) is superseded by WAFv2 (covered by AwsWafWebAcl) |
| `aws_waf_regex_pattern_set` | WAF Classic (global) is superseded by WAFv2 (covered by AwsWafWebAcl) |
| `aws_waf_rule` | WAF Classic (global) is superseded by WAFv2 (covered by AwsWafWebAcl) |
| `aws_waf_rule_group` | WAF Classic (global) is superseded by WAFv2 (covered by AwsWafWebAcl) |
| `aws_waf_size_constraint_set` | WAF Classic (global) is superseded by WAFv2 (covered by AwsWafWebAcl) |
| `aws_waf_sql_injection_match_set` | WAF Classic (global) is superseded by WAFv2 (covered by AwsWafWebAcl) |
| `aws_waf_web_acl` | WAF Classic (global) is superseded by WAFv2 (covered by AwsWafWebAcl) |
| `aws_waf_xss_match_set` | WAF Classic (global) is superseded by WAFv2 (covered by AwsWafWebAcl) |
| `aws_wafregional_byte_match_set` | WAF Classic (regional) is superseded by WAFv2 (covered by AwsWafWebAcl) |
| `aws_wafregional_geo_match_set` | WAF Classic (regional) is superseded by WAFv2 (covered by AwsWafWebAcl) |
| `aws_wafregional_ipset` | WAF Classic (regional) is superseded by WAFv2 (covered by AwsWafWebAcl) |
| `aws_wafregional_rate_based_rule` | WAF Classic (regional) is superseded by WAFv2 (covered by AwsWafWebAcl) |
| `aws_wafregional_regex_match_set` | WAF Classic (regional) is superseded by WAFv2 (covered by AwsWafWebAcl) |
| `aws_wafregional_regex_pattern_set` | WAF Classic (regional) is superseded by WAFv2 (covered by AwsWafWebAcl) |
| `aws_wafregional_rule` | WAF Classic (regional) is superseded by WAFv2 (covered by AwsWafWebAcl) |
| `aws_wafregional_rule_group` | WAF Classic (regional) is superseded by WAFv2 (covered by AwsWafWebAcl) |
| `aws_wafregional_size_constraint_set` | WAF Classic (regional) is superseded by WAFv2 (covered by AwsWafWebAcl) |
| `aws_wafregional_sql_injection_match_set` | WAF Classic (regional) is superseded by WAFv2 (covered by AwsWafWebAcl) |
| `aws_wafregional_web_acl` | WAF Classic (regional) is superseded by WAFv2 (covered by AwsWafWebAcl) |
| `aws_wafregional_web_acl_association` | WAF Classic (regional) is superseded by WAFv2 (covered by AwsWafWebAcl) |
| `aws_wafregional_xss_match_set` | WAF Classic (regional) is superseded by WAFv2 (covered by AwsWafWebAcl) |
