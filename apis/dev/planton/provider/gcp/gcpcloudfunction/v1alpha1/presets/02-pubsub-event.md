# Pub/Sub event processor

An event-driven worker consuming a Pub/Sub topic ([GcpPubSubTopic](/docs/catalog/gcp/gcppubsubtopic)) through Eventarc. Ingress is locked to internal traffic — nothing about an event consumer needs a public endpoint.

`RETRY_POLICY_RETRY` gives at-least-once delivery with exponential backoff, so the handler must be idempotent (duplicate events WILL arrive). `maxInstanceCount` is the backpressure ceiling: a topic storm scales the function only this far, protecting whatever the handler writes to downstream.
