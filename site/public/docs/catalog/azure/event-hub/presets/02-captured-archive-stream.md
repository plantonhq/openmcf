---
title: "Captured Archive Stream"
description: "This preset creates a stream with capture enabled: every event is archived to Blob Storage in Avro format automatically -- audit trails, replay beyond the retention window, and batch analytics feed..."
type: "preset"
rank: "02"
presetSlug: "02-captured-archive-stream"
componentSlug: "event-hub"
componentTitle: "Event Hub"
provider: "azure"
icon: "package"
order: 2
---

# Captured Archive Stream

This preset creates a stream with capture enabled: every event is
archived to Blob Storage in Avro format automatically -- audit trails,
replay beyond the retention window, and batch analytics feed straight
from the archive.

## When to Use

- Compliance/audit streams that must keep events beyond the hub's
  retention window
- Lambda-style architectures: real-time consumers on the hub, batch
  jobs on the capture archives

## Key Configuration Choices

- **Size-or-interval cadence** -- an archive window closes at 5 minutes
  OR 300 MB (Azure's default), whichever comes first
- **`skipEmptyArchives: true`** -- saves storage on sparse streams;
  leave false when downstream batch jobs expect a continuous file
  cadence
- **Service-managed SAS auth (the default)** -- no identity setup; for
  a keyless posture set `storageAuthenticationType: SYSTEM_ASSIGNED`
  (grant the namespace's principal Storage Blob Data Contributor) or
  `USER_ASSIGNED` with `storageAuthenticationId`
- **The all-nine-tokens `archiveNameFormat`** -- Azure requires every
  placeholder; the order and surrounding text are yours

## Values to Customize

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `my-capture-account` | The AzureStorageAccount holding archives | Your storage composition |
| `my-capture-container` | The AzureStorageContainer archives land in | Same composition |
| `audit-events` | The hub name | Your stream taxonomy |
