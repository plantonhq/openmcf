# SES Configuration Sets: The Policy Layer of Email Sending

## Introduction

Amazon SES separates WHO may send (identities) from HOW mail is sent (configuration sets). A configuration set is the named group of sending rules — TLS posture, IP pool, tracking, suppression, deliverability dashboards, and event publishing — that identities reference as their default and that any SendEmail call can name explicitly. Because many identities share one set, an organization defines its delivery posture once and every sender inherits it.

## Why Event Destinations Are the Heart of the Kind

SES is a black box until events flow out of it. Each named event destination streams a chosen slice of email events — SEND, REJECT, BOUNCE, COMPLAINT, DELIVERY, OPEN, CLICK, RENDERING_FAILURE, DELIVERY_DELAY, SUBSCRIPTION — into exactly one AWS destination:

- **CloudWatch** — events become metrics, dimensioned by message tags, headers, or link tags. The zero-infrastructure way to alarm on bounce rate.
- **EventBridge** — events land on the account's default bus for rule-based routing. AWS only supports the default bus for SES publishing today; rules can forward anywhere from there.
- **Kinesis Firehose** — the durable analytics path into S3/Redshift/OpenSearch. SES assumes an IAM role (trust principal `ses.amazonaws.com`) with `firehose:PutRecordBatch`.
- **SNS** — one message per event; the classic bounce/complaint feedback-loop wiring.
- **Pinpoint** — engagement events into a Pinpoint project (literal ARN; Pinpoint is not modeled in this catalog).

Bounce/complaint feedback loops — the difference between a healthy sender reputation and a suspended SES account — are built from exactly these destinations.

## Design Notes

- **Event destinations are folded, not a kind.** Each is an AWS sub-resource keyed by (set, name), set-scoped in lifecycle, many-per-set, and never referenced by anything else — the per-name materialization class. Both engines create one provider resource per named entry so destinations add/remove independently.
- **`enabled` defaults to true in the catalog.** AWS's own default is FALSE — a created-but-silent destination is a common source of missing events. The modules always send the value explicitly.
- **Suppression is a tri-state.** An absent `suppressedReasons` list inherits the account-level suppression configuration; an explicit list overrides it. The modules only emit the block when reasons are listed.
- **The set name is `metadata.name`** in both engines; `configuration_set_name` is exported as the output-backed join key identities reference.

## 90/10 Coverage Notes

The full `aws_sesv2_configuration_set` + `aws_sesv2_configuration_set_event_destination` surface is modeled: delivery options (TLS policy, max delivery seconds, sending pool), reputation metrics, the sending kill switch, suppression overrides, custom tracking domains with HTTPS policy, VDM dashboard/guardian overrides, and all five event-destination arms.

## Deferred Surface (recorded reasons)

- **Dedicated IP pools + assignments** (`aws_sesv2_dedicated_ip_pool`, `aws_sesv2_dedicated_ip_assignment`) — a paid capacity surface with its own lifecycle; the set's `sendingPoolName` arm composes by name with zero rework when a pool exists.
- **Contact lists** (`aws_sesv2_contact_list`) — marketing contact management, a data-plane surface, not sending infrastructure.
- **Tenants** (`aws_sesv2_tenant` + resource associations) — the new multi-tenant sending-isolation surface; revisit on concrete pull.
- **Account-level VDM / suppression attributes** (`aws_sesv2_account_vdm_attributes`, `aws_sesv2_account_suppression_attributes`) — account-scoped regional singletons, not graph resources.
- **Classic SES receiving** (receipt rule sets, receipt rules, filters) — email RECEIVING is a separate product surface from sending.
- **Templates** (`aws_ses_template`) — a classic-SES-only resource; SESv2 templates have no provider resource today.
