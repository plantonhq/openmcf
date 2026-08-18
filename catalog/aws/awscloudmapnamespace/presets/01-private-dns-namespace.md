# Private DNS Namespace

The ECS-style shape: `api.corp.internal` for tasks to register into (wire the service ARN to the ECS service's registry), plus a statically registered CNAME so the database resolves by the same private domain.
