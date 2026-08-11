# Entra-Auth Private Endpoint

This preset creates the hardened posture: Microsoft Entra token authentication (no scoring secret exists at all) with public network access disabled, for internal services scoring over private networking.

## When to Use

- Internal applications that already carry Entra identities
- Environments where standing scoring keys are a compliance finding
- Workspaces behind Private Link where scoring must not traverse the public internet

## Key Configuration Choices

- **`authMode: AADToken`** -- callers obtain a token for the ML audience with their own identity (managed identity or service principal) and send it as the bearer; nothing to store, nothing to rotate.
- **`publicNetworkAccessEnabled: false`** -- the scoring address answers only through the workspace's private endpoints; verify the caller's network path before shifting production traffic.
- **`identity: SYSTEM_ASSIGNED`** -- unchanged from the everyday shape; the endpoint's own identity still needs its registry and storage grants.

## After Deployment

Grant callers a data-plane RBAC role on the endpoint (e.g. `AzureML Data Scientist` or a custom scoring role), attach a deployment named `blue`, and validate scoring from inside the private network first.
