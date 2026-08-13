# Overview

The **AzureEventgridNamespace** component deploys an Azure Event Grid namespace -- the capacity-scaled hub of the newer Event Grid. Where the classic resources (AzureEventgridTopic, AzureEventgridDomain) each own one endpoint, a namespace hosts many CloudEvents namespace topics behind one set of regional endpoints and, optionally, runs an MQTT broker for IoT-style pub/sub -- all sized by throughput units you can change in place.

## Purpose

- **One hub, many streams**: namespace topics (AzureEventgridNamespaceTopic) are created inside the namespace as their own resources -- teams onboard streams without touching the hub.
- **MQTT without servers**: the topic-spaces block turns on a managed MQTT broker (client certificate authentication, session dials, optional routing of every MQTT message into a custom topic).
- **Capacity as a dial**: throughput units (1-40) set the ceiling for all topics and MQTT traffic together and update in place.

## Key Features

- Full azurerm v5 surface: capacity, public network access, inbound IP rules (up to 128 CIDRs), managed identity (including the combined mode), the complete MQTT topic-spaces configuration, tags.
- Chart-ready: `resource_group` defaults its reference to AzureResourceGroup, identity IDs to AzureUserAssignedIdentity, and the MQTT route topic to AzureEventgridTopic; the `namespace_id` output is the wiring edge an AzureEventgridNamespaceTopic's `namespace_id` references.
- The SKU has exactly one legal value ("Standard") and is deliberately not part of the spec -- both engines send it explicitly.

## Use Cases

- **CloudEvents backbone**: one namespace per environment; each service publishes to its own namespace topic.
- **IoT ingestion**: MQTT clients authenticate with client certificates and publish telemetry; route everything into an Event Grid topic for fan-out to non-MQTT consumers.
- **Network-restricted eventing**: disable public access or pin inbound CIDRs while private endpoints carry the traffic.

## Future Enhancements

- Stream wiring lives in AzureEventgridNamespaceTopic -- point its `namespace_id` at this component's ID output.

---

© Planton. Licensed under [Apache-2.0](https://github.com/plantonhq/planton/blob/main/LICENSE).
