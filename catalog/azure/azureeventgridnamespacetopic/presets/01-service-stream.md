# Service Stream

This preset creates one service's CloudEvents stream inside a shared Event Grid namespace -- the self-service onboarding shape: the platform team owns the namespace, each service owns its topic.

## When to Use

- One stream per publishing service inside an environment's shared hub
- Tenant onboarding where each tenant gets a named stream

## Key Configuration Choices

- **The name is the stream's identity** -- unique within the namespace only (no public hostname); changing it later replaces the topic and drops buffered events
- **Retention 7 days** -- the maximum delivery buffer; tune down for high-volume streams consumed quickly
- **Nothing else to choose** -- Azure fixes the schema (CloudEvents v1.0) and publisher type (Custom) on this resource

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<your-eventgrid-namespace>` | The Planton name of your `AzureEventgridNamespace` resource | Planton console (or replace `valueFrom` with `value:` and the namespace's ARM ID) |
| `orders` | The stream's name inside the namespace | Your service's naming convention |
