# Auth0 Action

Deploys an Auth0 Action -- a custom Node.js function that executes at a specific point in the Auth0 authentication pipeline. Supports 10 trigger types including post-login, pre-registration, and credentials-exchange. Integrates with Planton's Auth0 Provider Connection for credential management and optional trigger binding for immediate activation.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **Auth0 Action** -- a versioned Node.js function registered in Auth0 with the specified trigger, source code, runtime, npm dependencies, and encrypted secrets
- **Auth0 Trigger Binding** -- created only when `triggerBinding` is configured, attaches the deployed action to its trigger flow so it executes during the corresponding Auth0 pipeline stage

## Before You Deploy

### Planton Setup

- **Auth0 Provider Connection** -- an active connection in the Connect module with Auth0 domain, client ID, and client secret. Map it as the default for your environment, or specify it explicitly when creating the Cloud Resource.
- **Planton Runner** -- required when using Runner-based credential delivery. Not needed for inline credential authentication.

### Auth0 Account

- **An Auth0 tenant** with the Actions feature enabled.
- **Management API permissions** -- the client credentials in the Provider Connection must have scopes to create and deploy actions (`create:actions`, `update:actions`, `read:actions`, `read:triggers`, `update:triggers`).

## Deploy

### Console

Open the deployment store, find **Auth0 Action**, and click **Deploy**. The creation wizard walks you through trigger selection, action code, runtime configuration, and optional dependencies and secrets. Start from the **Post-Login Custom Claims** preset in the [Presets](#presets) tab to pre-populate a working configuration.

### CLI

Create a manifest and apply it:

```yaml
apiVersion: auth0.planton.dev/v1
kind: Auth0Action
metadata:
  name: enrich-tokens
  org: acme-corp
  env: prod
spec:
  supportedTrigger:
    id: post-login
    version: v3
  code: |
    exports.onExecutePostLogin = async (event, api) => {
      api.idToken.setCustomClaim('https://myapp/roles', event.authorization?.roles || []);
    };
  deploy: true
  triggerBinding:
    displayName: Enrich Tokens
```

```shell
planton apply -f auth0-action.yaml
```

This creates a post-login action, deploys it as an immutable version, and binds it to the post-login trigger flow. No dependencies or secrets are configured.

## Key Configuration

These are the most important decisions when configuring an Auth0 Action. Explore the full field reference in the [API Explorer](#api-explorer) tab.

- **Trigger selection** -- The `supportedTrigger.id` field determines when the action executes. Use `post-login` for token enrichment and access control, `pre-user-registration` for signup validation, `credentials-exchange` for M2M flow customization. Each trigger exposes different event and API objects to the action code.

- **Runtime version** -- The `runtime` field sets the Node.js version. Use `node22` for new actions. If omitted, Auth0 assigns a default based on the trigger version. Actions on `node18` should be migrated before the runtime reaches end-of-life.

- **Deploy and bind** -- Setting `deploy` to true creates an immutable version on every apply. Adding `triggerBinding` attaches the action to its trigger flow. Omit `triggerBinding` when trigger ordering is managed separately or when staging code without activation.

- **Secrets management** -- The `secrets` array makes sensitive values available as `event.secrets.<name>` in action code. Auth0 requires full secret management -- if you manage secrets here, include ALL secrets for this action. Omitting a previously-set secret deletes it.

- **Dependencies** -- The `dependencies` array specifies npm packages installed when Auth0 builds the action. Pin versions explicitly (e.g., `1.7.0`) rather than using ranges to ensure reproducible builds.

## Outputs and Dependencies

### What This Component Consumes

This component has no foreign key dependencies.

### What This Component Provides

After provisioning, `status.outputs` contains values that downstream Cloud Resources can consume via ValueFromRef:

| Output | Description | Common Downstream Use |
|--------|-------------|----------------------|
| `id` | Auth0-assigned action identifier | Trigger binding references, API operations |
| `name` | Action name from metadata | Audit logs, monitoring dashboards |
| `version_id` | Deployed version identifier (when `deploy` is true) | Version tracking, rollback references |
| `runtime` | Resolved Node.js runtime version | Runtime compatibility checks |

## Common Patterns

Browse the [Presets](#presets) tab for ready-to-deploy configurations.

- **Post-login token enrichment** -- add custom claims (roles, email, organization) to ID and access tokens after authentication. Start from the **Post-Login Custom Claims** preset.
- **Registration domain allowlist** -- restrict signups to approved email domains using a secret-driven allowlist checked during pre-registration. Start from the **Pre-Registration Domain Allowlist** preset.
- **M2M audit logging** -- log client credentials exchanges to an external endpoint for compliance and monitoring, with npm dependencies and bearer token authentication. Start from the **Credentials Exchange Audit Log** preset.

## Works With

This component operates independently and does not reference other components.