# Auth0Connection Pulumi Module Architecture

## Overview

This document describes the architecture and design of the Auth0Connection Pulumi module.

## Architecture Diagram

```
┌─────────────────────────────────────────────────────────────┐
│                     Stack Input                              │
│  (Auth0ConnectionStackInput from YAML manifest)              │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                     main.go (entrypoint)                     │
│  1. Load stack input from environment                        │
│  2. Call module.Resources()                                  │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                   module/main.go                             │
│  1. Initialize locals from stack input                       │
│  2. Create Auth0 provider with credentials                   │
│  3. Create connection resource                               │
│  4. Export outputs                                           │
└─────────────────────────────────────────────────────────────┘
        │                     │                       │
        ▼                     ▼                       ▼
┌───────────────┐   ┌─────────────────┐   ┌──────────────────┐
│  locals.go    │   │  connection.go  │   │   outputs.go     │
│               │   │                 │   │                  │
│ Compute local │   │ Create Auth0    │   │ Export stack     │
│ values from   │   │ connection      │   │ outputs:         │
│ spec:         │   │ based on        │   │ - id             │
│ - name        │   │ strategy:       │   │ - name           │
│ - strategy    │   │ - auth0         │   │ - strategy       │
│ - options     │   │ - social        │   │ - enabled        │
│               │   │ - saml          │   │ - clients        │
│               │   │ - oidc          │   │                  │
│               │   │ - waad          │   │                  │
└───────────────┘   └─────────────────┘   └──────────────────┘
```

## Module Components

### main.go (Entrypoint)

The entrypoint is responsible for:
1. Creating the Pulumi runtime context
2. Loading the `Auth0ConnectionStackInput` from environment
3. Invoking the module's `Resources` function

### module/main.go

The module orchestrator:
1. Initializes local values from stack input
2. Creates the Auth0 provider with either:
   - Explicit credentials from `provider_config`
   - Default provider (environment variables)
3. Calls `createConnection` to create the resource
4. Calls `exportOutputs` to export stack outputs

### module/locals.go

Computes derived values from the stack input:
- Connection name from metadata
- Display name (with fallback to metadata name)
- Strategy-specific option structs
- Client and realm lists

### module/connection.go

Creates the Auth0 connection resource:
1. Builds base connection arguments (name, strategy, display_name)
2. Adds optional fields (enabled_clients, realms, metadata)
3. Builds strategy-specific options via `buildConnectionOptions`
4. Creates the `auth0.Connection` resource

Strategy-specific option building handles:
- Database: password policy, brute force, MFA
- Social: client credentials, scopes
- SAML: endpoints, certificates, attribute mappings
- OIDC: issuer, client credentials, discovery overrides
- Azure AD: domain, tenant, group settings

### module/outputs.go

Exports connection information to Pulumi stack outputs:
- Connection ID and name
- Strategy type
- Enabled status (based on client count)
- Enabled client IDs
- Realms

## Data Flow

```
YAML Manifest
     │
     ▼
Auth0ConnectionStackInput (protobuf)
     │
     ▼
Locals struct (computed values)
     │
     ├──────────────────┐
     ▼                  ▼
Auth0 Provider    Connection Args
     │                  │
     └────────┬─────────┘
              ▼
     auth0.Connection
              │
              ▼
     Stack Outputs
```

## Strategy Option Mapping

| Spec Field | Pulumi Field | Notes |
|------------|--------------|-------|
| `database_options.password_policy` | `options.PasswordPolicy` | String enum |
| `database_options.brute_force_protection` | `options.BruteForceProtection` | Boolean |
| `social_options.client_id` | `options.ClientId` | OAuth client ID |
| `social_options.scopes` | `options.Scopes` | String array |
| `saml_options.sign_in_endpoint` | `options.SignInEndpoint` | URL |
| `saml_options.signing_cert` | `options.SigningCert` | PEM certificate |
| `oidc_options.issuer` | `options.Issuer` | OIDC issuer URL |
| `azure_ad_options.domain` | `options.Domain` | Azure AD domain |

## Error Handling

Errors are propagated using `github.com/pkg/errors`:
- Provider creation failures
- Connection creation failures
- Missing required configuration

All errors include context about what operation failed.

## Resource Dependencies

The module creates a single primary resource:
- `auth0.Connection` - The identity connection

The Auth0 provider is created as an explicit provider to ensure credential isolation.

## Extension Points

To add a new connection strategy:
1. Add strategy value to `spec.proto`
2. Add strategy-specific options message to `spec.proto`
3. Add case to `buildConnectionOptions` in `connection.go`
4. Add local extraction in `locals.go`
5. Update documentation and examples

