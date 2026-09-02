# Forward AND Process with a Worker

One rule, two actions: forward matched mail to a verified inbox AND hand a copy to an Email Worker (e.g. to file a ticket or post to a channel). Actions apply in order.

## When to use

- A human needs the mail in their inbox while automation processes it too.
- Archiving mail to a system of record while still delivering it.

## Key choices

- `actions` is a list — this is the Cloudflare API's real shape, so one rule can combine `forward` and `worker` (a `drop` action stands alone by nature).
- Each `forward` destination must be a verified `CloudflareEmailRoutingAddress`; the `worker` must be an Email Worker.

## Placeholders

| Placeholder | Description |
|---|---|
| `<zone-name>` | Name of the CloudflareDnsZone |
| `<rule-name>` | Descriptive rule name |
| `<matched-address>` | The recipient address to match (e.g. support@example.com) |
| `<destination-email>` | The verified destination mailbox |
| `<email-worker-name>` | The CloudflareWorker handling the message |
