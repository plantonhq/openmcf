# Filtered Webhook

This preset delivers only high-value order events to a partner's HTTPS endpoint -- the filtered-push shape: platform-side filtering keeps noise off the wire, and retry tuning plus dead-lettering acknowledge that a webhook's availability is not yours to control.

## When to Use

- Partner/third-party integrations reachable only over HTTPS
- Handlers that should see a narrow slice of a busy topic

## Key Configuration Choices

- **The endpoint must be live before deploy** -- Event Grid's create-time validation handshake fails otherwise; in charts, sequence the handler first
- **Subject prefix + numeric filter compose** -- conditions AND together; the whole filter shares Azure's 25-value budget
- **Retry tuned below the defaults** -- 10 attempts over 4 hours beats 30 attempts over 24 for latency-sensitive integrations; exhausted events dead-letter instead of dropping
- **CloudEvents delivery schema** -- the partner sees the vendor-neutral envelope regardless of the topic's input shape (create-only)

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<your-eventgrid-topic>` | The Planton name of your `AzureEventgridTopic` resource | Planton console |
| `https://handler.example.com/events` | The partner's HTTPS endpoint | Their integration docs |
| `<your-storage-account>` | The Planton name of your `AzureStorageAccount` resource (dead-letter home) | Planton console |
| `data.amount` | The payload field the numeric filter reads | Your event schema |
