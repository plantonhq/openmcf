# Auth0 Resource Server - Pulumi Module Overview

## Architecture

This module creates Auth0 Resource Server resources using the Pulumi Go SDK.

```mermaid
flowchart TB
    subgraph Input
        SI[StackInput]
        PC[ProviderConfig]
        Target[Auth0ResourceServer]
    end
    
    subgraph Module
        Locals[Initialize Locals]
        Provider[Create Provider]
        RS[Create Resource Server]
        Scopes[Create Scopes]
        Outputs[Export Outputs]
    end
    
    subgraph Auth0
        API[Resource Server]
        Perms[Permissions/Scopes]
    end
    
    SI --> Locals
    SI --> PC
    SI --> Target
    
    PC --> Provider
    Locals --> RS
    Provider --> RS
    RS --> API
    
    Locals --> Scopes
    RS --> Scopes
    Scopes --> Perms
    
    RS --> Outputs
    Scopes --> Outputs
```

## Component Structure

### Entry Point (`main.go`)

Loads stack input and invokes the module:

```go
pulumi.Run(func(ctx *pulumi.Context) error {
    stackInput := &auth0resourceserverv1.Auth0ResourceServerStackInput{}
    if err := stackinput.LoadStackInput(ctx, stackInput); err != nil {
        return err
    }
    return module.Run(ctx, stackInput)
})
```

### Module (`module/main.go`)

Orchestrates resource creation:

1. Validate input
2. Initialize locals
3. Create Auth0 provider
4. Create resource server
5. Create scopes (if defined)
6. Export outputs

### Locals (`module/locals.go`)

Computes derived values from stack input:

- Resource name from metadata
- Display name (spec.name or metadata.name)
- Token configuration
- Access control settings
- Scope list

### Resource Server (`module/resourceserver.go`)

Creates the Auth0 Resource Server resource:

```go
auth0.NewResourceServer(ctx, name, &auth0.ResourceServerArgs{
    Identifier: pulumi.String(identifier),
    Name:       pulumi.String(displayName),
    SigningAlg: pulumi.String(signingAlg),
    // ... other settings
})
```

Also creates scopes using `auth0.NewResourceServerScopes`.

### Outputs (`module/outputs.go`)

Exports stack outputs:

- `id` - Auth0 resource ID
- `identifier` - API audience
- `name` - Display name
- `signing_alg` - Signing algorithm
- Token settings
- Access control flags

## Resource Dependencies

```mermaid
flowchart LR
    Provider[Auth0 Provider]
    RS[Resource Server]
    Scopes[Resource Server Scopes]
    
    Provider --> RS
    RS --> Scopes
    Provider --> Scopes
```

## Configuration Flow

```mermaid
sequenceDiagram
    participant CLI as OpenMCF CLI
    participant Pulumi
    participant Module
    participant Auth0
    
    CLI->>Pulumi: pulumi up
    Pulumi->>Module: Load StackInput
    Module->>Module: Initialize Locals
    Module->>Auth0: Create Provider
    Module->>Auth0: Create Resource Server
    Auth0-->>Module: Resource Server ID
    Module->>Auth0: Create Scopes
    Auth0-->>Module: Scopes Created
    Module->>Pulumi: Export Outputs
    Pulumi-->>CLI: Deployment Complete
```

## Error Handling

The module handles errors at each step:

1. **Input validation**: Ensures required fields are present
2. **Provider creation**: Validates Auth0 credentials
3. **Resource creation**: Handles API errors gracefully
4. **Scope creation**: Ensures resource server exists first

## Testing Strategy

1. **Unit tests**: Validate input transformation in locals
2. **Integration tests**: Deploy to test tenant
3. **Smoke tests**: Use hack manifest for quick validation

## Performance Considerations

- Single API call for resource server creation
- Batch scope creation in one resource
- Minimal state management overhead
