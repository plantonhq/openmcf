# AwsEventBridgeBus — Research Documentation

## Overview

Amazon EventBridge is a serverless event bus service that connects applications using events. It ingests events from AWS services, SaaS applications, and custom sources, then routes them to targets based on rules. EventBridge is the evolution of CloudWatch Events, providing the same core functionality with expanded features for SaaS integrations and schema management.

## Architecture

EventBridge operates around three core concepts: **event buses**, **rules**, and **targets**.

**Event buses** receive events. Every AWS account has a default event bus that automatically receives events from AWS services (EC2 state changes, S3 object operations, etc.). Custom event buses isolate application-defined events from the default bus, enabling independent access control, encryption, and monitoring.

**Rules** match incoming events against patterns and route matched events to one or more targets. Rules are attached to a specific bus (`AwsEventBridgeRule` composes onto this bus's `bus_name` output).

**Targets** are AWS services that process events (Lambda, SQS, SNS, Step Functions, Kinesis, etc.), modeled as folded blocks on the rule.

### Custom vs Default Bus

The default bus is shared across all AWS services in an account. Custom buses provide:
- **Isolation**: Application events are separated from AWS service events.
- **Access control**: Fine-grained resource policies per bus.
- **Encryption**: Customer-managed KMS keys for event encryption.
- **Partner integration**: SaaS partners deliver events to dedicated custom buses.

A bus can never be named `default` — both IaC modules fail fast on that name before reaching the AWS API, since the account's built-in bus already exists and rules can target it directly.

### Partner Event Sources

AWS EventBridge integrates with 30+ SaaS providers (Datadog, PagerDuty, Zendesk, etc.) via partner event sources. When you create a partner integration, AWS creates an event source in your account. You then create a custom bus with the same name as the event source, and the partner's events flow to your bus.

## Design Decisions

### Why StringValueOrRef for kms_key_identifier

The KMS key identifier field uses `StringValueOrRef` to enable infra-chart composability. In a typical infra chart, a KMS key is created as a separate resource and its ARN is wired into downstream resources. The `valueFrom` reference creates a dependency edge in the deployment DAG, ensuring the KMS key is provisioned before the bus.

### Why StringValueOrRef for dead_letter_config.arn

The DLQ ARN uses `StringValueOrRef` to enable the common pattern of defining both the bus and its DLQ in the same infra chart. The DLQ (an SQS queue) is deployed first, and the bus's `deadLetterConfig.arn` references the queue's output ARN via `valueFrom`.

### Why string + CEL for log levels (not proto enums)

Log levels (`OFF`, `ERROR`, `INFO`, `TRACE`) and include_detail values (`NONE`, `FULL`) use plain strings with CEL `in` validation rather than protobuf enums. This keeps the values provider-authentic (matching the exact AWS API strings) and avoids proto enum prefix conventions.

### Why the resource policy IS folded

The bus's resource-based policy (`aws_cloudwatch_event_bus_policy`) is a **single-per-bus** setting keyed by the bus name — the classic folded-satellite shape. It is the mechanism behind cross-account and cross-organization event ingestion (`events:PutEvents` grants), and it has no independent lifecycle: it exists to configure exactly one bus. The modules materialize it as its own provider resource when `resource_policy` is set. The granular `aws_cloudwatch_event_permission` resource is deliberately NOT modeled — it is the same policy expressed one statement at a time, and two writers to one policy document fight on every apply; author the full document in `resource_policy` instead.

### Deliberately Omitted

- **Event archive** (`aws_cloudwatch_event_archive`): event retention/replay with its own lifecycle (many archives per bus) — a candidate kind on demand.
- **Connections + API destinations** (`aws_cloudwatch_event_connection` / `_api_destination`): the outbound HTTP-integration surface with its own auth lifecycle — a candidate kind pair on demand.
- **Global endpoints** (`aws_cloudwatch_event_endpoint`): multi-region failover routing across two buses — a separate DR product surface.
- **Schema discovery**: EventBridge Schema Registry integration. Separate concern.

## Terraform Provider Reference

The primary Terraform resource is `aws_cloudwatch_event_bus` from the `hashicorp/aws` provider (the 6.x line — `kms_key_identifier`, `dead_letter_config`, and `log_config` all require it).

Key attributes:
- `name` (ForceNew) — bus name is immutable
- `event_source_name` (ForceNew) — partner source is immutable
- `kms_key_identifier` — KMS key ARN, key ID, key alias, or key alias ARN
- `dead_letter_config` — nested block with `arn` (SQS queue ARN)
- `log_config` — nested block with `level` and `include_detail`
- `description` — up to 512 characters

Related resources:
- `aws_cloudwatch_event_bus_policy` — resource-based policy (folded as `resource_policy`)
- `aws_cloudwatch_event_rule` — rules that match and route events (separate component: AwsEventBridgeRule)
- `aws_cloudwatch_event_target` — targets attached to rules (folded on the rule)

## Pulumi Resource Reference

The Pulumi resources are `cloudwatch.EventBus` and `cloudwatch.EventBusPolicy` from `pulumi-aws/sdk/v7/go/aws/cloudwatch`. Input properties map directly to Terraform attributes with camelCase naming. The bus name and ARN are the primary outputs.
