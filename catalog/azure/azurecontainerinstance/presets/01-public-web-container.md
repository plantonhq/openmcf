# Public Web Container

This preset runs one Linux container with a public IP, a stable DNS name, and a liveness probe -- the simplest always-on service shape, no cluster to manage.

## When to Use

- A small API, webhook receiver, or demo app that needs a public endpoint today
- Trying an image in a real Azure network context before committing to Container Apps or AKS
- Anything where the group IS the service -- one unit, restarted by Azure via the default "Always" policy

## Key Configuration Choices

- **`dnsNameLabel`** -- the group becomes `{label}.{region}.azurecontainer.io`; the label must be unique per region
- **`dnsNameLabelReusePolicy: TenantReuse`** -- stricter than the default "Unsecure", which lets ANYONE claim the label after this group releases it (dangling-DNS takeover)
- **`livenessProbe`** -- restarts the container in place when it stops answering; keep the probe cheap
- **Ports omitted from `exposedPorts`** -- the group exposes every container port by default; add `exposedPorts` only to narrow the set

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<your-resource-group>` | The Planton name of your `AzureResourceGroup` resource | Planton console (or replace `valueFrom` with `value:` and a literal group name) |
| `hello-web-acme` | The DNS label -- unique per region | Your naming convention |
| `eastus` | The Azure region | Your region strategy |

## Related Presets

- **Private VNet Worker** -- the subnet-joined posture with no public surface
- **Run-Once Job** -- the Never-restart batch shape
