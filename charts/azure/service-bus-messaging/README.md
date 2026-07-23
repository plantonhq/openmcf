# Azure Service Bus Messaging

The asynchronous integration backbone: dead-letter-postured command queues for point-to-point work dispatch, a SQL-filtered topic fan-out for domain events, least-privilege credentials scoped per entity, and an alert that fires on the first dead-lettered message. Correct messaging topology is subtle — this chart deploys it right and teaches why.

## Who this is for

An integration or platform team decoupling services: producers that must not block on consumers, events that multiple services react to independently, and the operational safety net (DLQs, alerts, scoped credentials) that separates a production message bus from a demo. Deploy it, hand each service its own credential, and the topology is already the one you would have arrived at after the first incident.

## Architecture

```
  producers ── send-only key ──▶ [orders queue]  ── listen-only key ──▶ orders service
              (per queue)       [billing queue]                        billing service
                                     │ 10 attempts / 14d TTL
                                     ▼
                                  DLQ (per queue) ──▶ metric alert ▶ ops email

  publishers ── send-only key ──▶ (events topic) ──▶ [order-events sub]  SQL: eventType = 'order.created'
                                        │           [audit sub]         SQL: 1=1 (everything)
                                        └── shared listen-only key; isolation by filter
```

Two messaging shapes, deliberately separate:

- **Command queues** — one per consuming service. Competing consumers, ten delivery attempts, expired messages park in the DLQ. Each queue gets its own send-only and listen-only SAS rule, so a leaked producer credential cannot consume and one service's credential cannot touch another's queue.
- **Event topic** — publishers write once; the broker fans out. The filtered subscription receives only what its SQL filter admits (evaluated broker-side against message properties); the audit subscription's `1=1` rule states the catch-all intent explicitly — Service Bus's auto-created `$Default` rule cannot be declared, and once custom rules exist only they govern delivery.

## Resources

| Kind | Name | Purpose |
| --- | --- | --- |
| AzureResourceGroup | `{env}-messaging` | One container for the estate |
| AzureLogAnalyticsWorkspace | `{env}-messaging-logs` | Broker operational logs for post-alert forensics |
| AzureMonitorDiagnosticSetting | `{env}-messaging-bus-diag` | Namespace logs + metrics into the workspace |
| AzureServiceBusNamespace | `{env}-bus` | The broker (STANDARD; PREMIUM toggle) |
| AzureServiceBusQueue | `{env}-{name}-queue` (per entry) | Command queues, DLQ-postured |
| AzureServiceBusAuthorizationRule | `{env}-{name}-send` / `-listen` (per queue) | Least-privilege credentials |
| AzureServiceBusTopic | `{env}-{topic}-topic` | The domain-event fan-out point |
| AzureServiceBusSubscription | `{env}-{sub}-sub` | SQL-filtered consumer |
| AzureServiceBusSubscription | `{env}-{topic}-audit-sub` | Unconditional copy of every event |
| AzureServiceBusAuthorizationRule | `{env}-{topic}-send` / `-listen` | Topic-scoped publisher/subscriber credentials |
| AzureMonitorActionGroup | `{env}-messaging-ops` | Alert routing |
| AzureMonitorMetricAlert | `{env}-messaging-dlq-alert` | Fires on the first dead-lettered message, namespace-wide |

## Parameters

| Parameter | Description | Default | Must change |
| --- | --- | --- | --- |
| `region` | Azure region | `centralus` | |
| `servicebus_namespace_name` | Globally unique namespace name | `my-integration-bus` | yes |
| `premium_enabled` | PREMIUM tier (dedicated capacity, CMK/VNet-rules surface) | `false` | |
| `command_queue_names` | One command queue per consumer service | `orders`, `billing` | |
| `events_topic_name` | The domain-event topic | `events` | |
| `filtered_subscription_name` | The SQL-filtered subscription | `order-events` | |
| `filter_expression` | Broker-side SQL filter on message properties | `eventType = 'order.created'` | |
| `duplicate_detection_enabled` | Drop repeated MessageIds within 10 minutes (fixed at creation) | `false` | |
| `ops_email` | Where the DLQ alert lands | `ops@example.com` | yes |
| `log_retention_days` | Workspace retention | `30` | |

## After deploying

1. **Hand out credentials, narrowest first** — each authorization rule's `primary_connection_string` output is the credential for exactly one entity and one direction. Producers get `{queue}-send`; consumers get `{queue}-listen`; nothing gets the namespace root key.
2. **Publish with properties** — the SQL filter evaluates user properties, so publishers set them explicitly (e.g. `eventType`, `region`, `amount`). A publisher that omits the property never matches the filter.
3. **Confirm the alert path** — Azure sends a test notification when the action group deploys; the first real dead-letter fires within the alert's 5-minute window.

## Day 2

- **New consumer service**: add its queue name to `command_queue_names` — queue plus both credentials appear on the next deploy.
- **New event consumer**: add an `AzureServiceBusSubscription` with its own SQL rule; publishers change nothing.
- **Keyless workloads**: workloads running on Azure compute with managed identities should prefer Entra data-plane roles (Azure Service Bus Data Sender / Data Receiver) over these SAS rules — the SAS credentials exist for the external and legacy producers that cannot.
- **Premium posture**: with `premium_enabled`, the namespace also models CMK and VNet network rules — configure them on the namespace resource directly when compliance requires.

---

© Planton. Licensed under [Apache-2.0](https://github.com/plantonhq/planton/blob/main/LICENSE).
