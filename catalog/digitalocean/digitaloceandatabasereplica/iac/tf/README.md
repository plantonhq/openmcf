# Terraform Module: DigitalOcean Database Replica

Provisions a single-node read-only replica of a DigitalOcean managed database cluster -- the complete `digitalocean_database_replica` resource surface.

## Resources

| Resource | Purpose |
|---|---|
| `digitalocean_database_replica.replica` | The replica: name, region, size, VPC placement, storage, tags |

## Inputs

Generated `variables.tf` mirrors the `DigitalOceanDatabaseReplicaSpec` proto: `cluster` and optional `vpc` (resolved references -- arrive as literal UUIDs), `replica_name`, `region` (enum slug), `size`, optional `storage_size_mib` (number; rendered as the provider's bare-MiB string), optional `tags`. Authentication uses `digitalocean_token` (sensitive).

## Outputs

Exactly the `DigitalOceanDatabaseReplicaStackOutputs` contract: `replica_id` (the API UUID -- the `uuid` attribute, not the legacy composite state id), `cluster_id`, `replica_name`, `host`, `private_host`, `port`, `database`, `user`, and the secrets `password`, `uri`, `private_uri`.

## Behavior notes

- Only `size`/`storage_size_mib` update in place (a resize waited to "online"); every other change -- including TAGS (create-only upstream) -- replaces the replica.
- `region` and `size` are required by the spec (the provider's omitted-value drift class is unrepresentable).
- Tags are spec tags plus the standard Planton labels, identical to the Pulumi module.
- Import: `terraform import ... <cluster_id>,<replica_name>` (see `iac/import-map.yaml`).
