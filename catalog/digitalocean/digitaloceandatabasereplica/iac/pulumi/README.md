# Pulumi Module: DigitalOcean Database Replica

Provisions a single-node read-only replica of a DigitalOcean managed database cluster -- the complete `digitalocean_database_replica` resource surface, at 100% behavioral parity with the Terraform module (same arguments, same outputs).

## Layout

- `main.go` -- entrypoint (`package main`), loads the stack input and calls the module
- `module/main.go` -- orchestration: locals, provider, resource
- `module/database_replica.go` -- the `DatabaseReplica` resource and output exports
- `module/locals.go` -- metadata handle + the standard Planton label map (rendered as tags, identical to Terraform)
- `module/outputs.go` -- output key constants (the `DigitalOceanDatabaseReplicaStackOutputs` contract)

## Behavior notes

- `replica_id` exports the SDK's `Uuid` attribute (the API UUID); the resource's own ID() is a legacy composite string.
- Tags are create-only upstream -- a retag replaces the replica; the module documents it where the tags are built.
- `storage_size_mib` renders as the provider's bare-MiB string from the spec's number.
