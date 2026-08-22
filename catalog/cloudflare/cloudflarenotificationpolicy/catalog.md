# Cloudflare Notification Policy

A notification policy: one Cloudflare alert type delivered to email, PagerDuty, and webhook destinations, optionally narrowed by filters. Real create, update, delete.

## What Gets Created

When you deploy this resource, the IaC module provisions:

- **Notification policy** -- one `cloudflare_notification_policy`

## Prerequisites

- **A Cloudflare API token** with Notifications → Write
- **A plan carrying the chosen alert type** (several families are Business or Enterprise only)
- For PagerDuty delivery, the PagerDuty integration connected in the dashboard

## Quick Start

```yaml
apiVersion: cloudflare.planton.dev/v1alpha1
kind: CloudflareNotificationPolicy
metadata:
  name: tunnel-health
  org: acme-corp
  env: prod
spec:
  accountId: "0a1b2c3d4e5f6a7b8c9d0e1f2a3b4c5d"
  name: tunnel-health
  alertType: tunnel_health_event
  mechanisms:
    emails:
      - oncall@example.com
```

```shell
planton apply -f notification-policy.yaml
```

## Configuration Reference

### Required Fields

| Field | Type | Description | Validation |
|-------|------|-------------|------------|
| `accountId` | string | The Cloudflare account. | Required, 32-hex; replaces on change. |
| `name` | string | Policy name. | Required. |
| `alertType` | string | The event class watched. | One of 69 provider values (enum-walled). |
| `mechanisms` | object | Delivery destinations. | Required; at least one of emails, pagerdutyIds, webhookIds. |

### Optional Fields

| Field | Type | Description | Validation |
|-------|------|-------------|------------|
| `description` | string | Why the policy exists. | Free-form. |
| `enabled` | bool | Whether the policy fires. | Cloudflare's default is enabled. |
| `alertInterval` | string | Refire floor, e.g. `30m`. | Honored per alert type. |
| `filters` | object | 43 list fields narrowing events. | `incidentImpact` (4 values) and `trafficExclusions` (`security_events`) enum-walled; the rest free-form per the provider. |

## Destroy Semantics

Real delete. The alerts stop; destinations and every other policy are untouched.

## Stack Outputs

| Output | Type | Description |
|--------|------|-------------|
| `policyId` | string | The Cloudflare-assigned UUID of the policy |

## Related Components

- [Cloudflare Notification Webhook](/docs/catalog/cloudflare/cloudflarenotificationwebhook) -- webhook destinations this policy references
- [Cloudflare Healthcheck](/docs/catalog/cloudflare/cloudflarehealthcheck) -- origin probes behind health-check alerts
- [Cloudflare Logpush Job](/docs/catalog/cloudflare/cloudflarelogpushjob) -- log export, and a subject of the failing-job alert
