# Standard Web Application Insights

This preset creates a workspace-based Application Insights resource for a web application at full telemetry fidelity -- the default APM shape. Its `connection_string` output is what Function Apps, Web Apps, and Container Apps reference to wire up telemetry.

## When to Use

- Any web application or HTTP service (regardless of language -- WEB is also right for OpenTelemetry-instrumented workloads)
- Development and staging environments where full-fidelity telemetry is worth the volume
- The first Application Insights resource in an environment

## Key Configuration Choices

- **WEB application type** -- shapes the portal's APM experiences; the right choice for anything serving HTTP
- **Workspace binding** (`workspaceId`) -- telemetry lands in the referenced Log Analytics Workspace; the binding is repointable but never removable
- **Full sampling (default 100%)** -- every request/dependency/exception is collected; move to the production-sampled preset when volume costs bite

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `my-observability-rg` | Resource group holding the resource | `AzureResourceGroup` status outputs |
| `my-platform-logs` | Workspace storing the telemetry | `AzureLogAnalyticsWorkspace` status outputs |
| `my-web-app-insights` | Resource name | Your naming convention |

## Related Presets

- **02-production-sampled** -- Cost-controlled production shape with sampling and a lower data cap
