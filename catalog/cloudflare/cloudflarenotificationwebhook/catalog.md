# Cloudflare Notification Webhook

A webhook destination for Cloudflare alerting: the HTTPS endpoint notification policies deliver to. Real create, update, delete.

## What Gets Created

When you deploy this resource, the IaC module provisions:

- **Webhook destination** -- one `cloudflare_notification_policy_webhooks`

## Prerequisites

- **A Cloudflare API token** with Notifications → Write
- **A live HTTPS endpoint** (a vendor incoming-webhook URL or your own receiver)

## Quick Start

```yaml
apiVersion: cloudflare.planton.dev/v1alpha1
kind: CloudflareNotificationWebhook
metadata:
  name: ops-slack
  org: acme-corp
  env: prod
spec:
  accountId: "0a1b2c3d4e5f6a7b8c9d0e1f2a3b4c5d"
  name: ops-slack
  url: "https://hooks.slack.com/services/T00000000/B00000000/XXXXXXXXXXXXXXXX"
```

```shell
planton apply -f notification-webhook.yaml
```

## Configuration Reference

### Required Fields

| Field | Type | Description | Validation |
|-------|------|-------------|------------|
| `accountId` | string | The Cloudflare account. | Required, 32-hex; replaces on change. |
| `name` | string | Destination name. | Required. |
| `url` | string | Endpoint alerts are POSTed to. | Required. |

### Optional Fields

| Field | Type | Description | Validation |
|-------|------|-------------|------------|
| `secret` | string | Shared secret or vendor API key. | Sensitive; write-only -- never returned by any read. |

## Destroy Semantics

Real delete. Every notification policy that referenced the destination silently stops delivering on that channel.

## Stack Outputs

| Output | Type | Description |
|--------|------|-------------|
| `webhookId` | string | The Cloudflare-assigned UUID policies reference |
| `type` | string | The destination kind Cloudflare inferred from the URL (datadog, discord, feishu, gchat, generic, opsgenie, slack, splunk) |

## Related Components

- [Cloudflare Notification Policy](/docs/catalog/cloudflare/cloudflarenotificationpolicy) -- the alert rules that deliver here
- [Cloudflare Logpush Job](/docs/catalog/cloudflare/cloudflarelogpushjob) -- log delivery, the record-level counterpart
