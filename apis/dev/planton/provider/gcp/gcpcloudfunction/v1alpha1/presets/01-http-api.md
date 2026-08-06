# HTTP API — basic

A public HTTPS endpoint for webhooks, small REST APIs, and glue code. The source lives as a versioned zip in GCS ([GcpGcsBucket](/docs/catalog/gcp/gcpgcsbucket)) — shipping a new deploy is uploading a new archive and pointing `object` at it.

`allowUnauthenticated: true` grants `run.invoker` to `allUsers` on the underlying Cloud Run service. For internal APIs, drop it and grant invoker to specific service accounts instead; for latency-sensitive endpoints, raise `minInstanceCount` to keep instances warm.
