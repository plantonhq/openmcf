# MQTT Broker

This preset creates a namespace with the MQTT broker enabled and every MQTT message routed into an Event Grid custom topic -- the IoT ingestion pattern: devices publish over MQTT, downstream consumers subscribe through the classic delivery machinery.

## When to Use

- Device fleets publishing telemetry over MQTT with client-certificate authentication
- Bridging MQTT traffic into Functions, queues, or webhooks via a custom topic's subscriptions

## Key Configuration Choices

- **The MQTT block is create-only** -- it cannot be added to or removed from an existing namespace; this preset exists because the choice must be made up front
- **`ClientCertificateSubject` as the identity source** -- the certificate's subject names the client; pick one convention fleet-wide (changing it later re-maps every device)
- **`routeTopicId` bridges to the classic world** -- the referenced custom topic must live in the same region and use the CloudEvents schema; remove the reference to keep traffic inside the broker
- **The dynamic enrichment stamps the client's name** onto every routed event, so consumers know the publisher

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<your-resource-group>` | The Planton name of your `AzureResourceGroup` resource | Planton console (or replace `valueFrom` with `value:` and a literal group name) |
| `<your-eventgrid-topic>` | The Planton name of the `AzureEventgridTopic` receiving routed MQTT messages | Planton console (or replace `valueFrom` with `value:` and the topic's ARM ID) |
| `my-org-mqtt-broker` | The namespace's name (appears in its hostnames) | Your naming convention -- org-prefixed |

## Related Presets

- **01 CloudEvents Hub** -- the pure-CloudEvents namespace without the MQTT broker
