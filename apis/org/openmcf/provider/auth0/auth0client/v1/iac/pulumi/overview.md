# Auth0Client Pulumi Module Architecture

## Overview

This document describes the architecture and design of the Auth0Client Pulumi module.

## Architecture Diagram

```
┌─────────────────────────────────────────────────────────────┐
│                     Stack Input                              │
│  (Auth0ClientStackInput from YAML manifest)                  │
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
│  3. Create client resource                                   │
│  4. Export outputs                                           │
└─────────────────────────────────────────────────────────────┘
        │                     │                       │
        ▼                     ▼                       ▼
┌───────────────┐   ┌─────────────────┐   ┌──────────────────┐
│  locals.go    │   │   client.go     │   │   outputs.go     │
│               │   │                 │   │                  │
│ Compute local │   │ Create Auth0    │   │ Export stack     │
│ values from   │   │ client based    │   │ outputs:         │
│ spec:         │   │ on type:        │   │ - id             │
│ - name        │   │ - spa           │   │ - client_id      │
│ - app_type    │   │ - native        │   │ - client_secret  │
│ - callbacks   │   │ - regular_web   │   │ - name           │
│ - grant_types │   │ - non_inter.    │   │ - app_type       │
│               │   │                 │   │                  │
└───────────────┘   └─────────────────┘   └──────────────────┘
```

## Module Components

### main.go (Entrypoint)

The entrypoint is responsible for:
1. Creating the Pulumi runtime context
2. Loading the `Auth0ClientStackInput` from environment
3. Invoking the module's `Resources` function

### module/main.go

The module orchestrator:
1. Initializes local values from stack input
2. Creates the Auth0 provider with either:
   - Explicit credentials from `provider_config`
   - Default provider (environment variables)
3. Calls `createClient` to create the resource
4. Calls `exportOutputs` to export stack outputs

### module/locals.go

Computes derived values from the stack input:
- Client name from metadata
- Application type
- OAuth settings (callbacks, logout URLs, grant types)
- JWT configuration
- Refresh token settings
- Mobile configuration (iOS/Android)

### module/client.go

Creates the Auth0 client (application) resource:
1. Builds base client arguments (name, app_type, description)
2. Adds URL configurations (callbacks, logout URLs, origins)
3. Adds OAuth settings (grant types, OIDC conformant)
4. Adds optional configurations:
   - JWT configuration
   - Refresh token settings
   - Mobile settings
   - Native social login
   - OIDC backchannel logout
5. Creates the `auth0.Client` resource

### module/outputs.go

Exports client information to Pulumi stack outputs:
- Client ID (public OAuth identifier)
- Client Secret (for confidential clients)
- Application name and type
- Signing keys
- Token endpoint auth method

## Data Flow

```
YAML Manifest
     │
     ▼
Auth0ClientStackInput (protobuf)
     │
     ▼
Locals struct (computed values)
     │
     ├──────────────────┐
     ▼                  ▼
Auth0 Provider    Client Args
     │                  │
     └────────┬─────────┘
              ▼
       auth0.Client
              │
              ▼
     Stack Outputs
```

## Application Type Mapping

| Application Type | OAuth Flow | Can Store Secret | Use Case |
|-----------------|------------|------------------|----------|
| `native` | Authorization Code + PKCE | No | Mobile/Desktop apps |
| `spa` | Authorization Code + PKCE | No | Single Page Apps |
| `regular_web` | Authorization Code | Yes | Server-side web apps |
| `non_interactive` | Client Credentials | Yes | M2M/Backend services |

## Configuration Mapping

| Spec Field | Pulumi Field | Notes |
|------------|--------------|-------|
| `application_type` | `AppType` | Required, one of 4 values |
| `callbacks` | `Callbacks` | Allowed callback URLs |
| `allowed_logout_urls` | `AllowedLogoutUrls` | Logout redirect URLs |
| `web_origins` | `WebOrigins` | CORS origins for SPAs |
| `grant_types` | `GrantTypes` | OAuth grant types |
| `jwt_configuration.alg` | `JwtConfiguration.Alg` | RS256, HS256, PS256 |
| `refresh_token.rotation_type` | `RefreshToken.RotationType` | rotating/non-rotating |
| `mobile.ios.team_id` | `Mobile.Ios.TeamId` | Apple Team ID |
| `mobile.android.app_package_name` | `Mobile.Android.AppPackageName` | Android package |

## Error Handling

Errors are propagated using `github.com/pkg/errors`:
- Provider creation failures
- Client creation failures
- Missing required configuration

All errors include context about what operation failed.

## Resource Dependencies

The module creates a single primary resource:
- `auth0.Client` - The OAuth application

The Auth0 provider is created as an explicit provider to ensure credential isolation.

## Extension Points

To add new client features:
1. Add field to `spec.proto`
2. Add local extraction in `locals.go`
3. Add argument mapping in `client.go`
4. Update output exports if needed in `outputs.go`
5. Update documentation and examples


