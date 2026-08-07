# Poison-Queue Companion

This preset declares the dead-letter companion of a work queue. Azure
Functions moves messages that exhaust their retries to `{queue}-poison`
by NAMING CONVENTION -- if nothing declares that queue, the runtime
creates it implicitly and it becomes unowned, unmonitored
infrastructure.

## When to Use

- Alongside every work queue consumed by a Functions queue trigger
- Any worker framework that adopts the `-poison` suffix convention

## Key Configuration Choices

- **The name IS the contract** -- `{work-queue-name}-poison`, exactly;
  the runtime routes by string match
- **`metadata.companion-of`** -- records the pairing so operators see
  the relationship when browsing the account
- **Monitor THIS queue's depth** -- a growing poison queue is the
  failure signal the work queue hides

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<storage-account-resource-name>` | The AzureStorageAccount's Planton resource name | Your storage composition |
| `order-processing` | The work queue this companion serves | Your queue composition |
