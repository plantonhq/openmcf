# Terraform Module: DigitalOcean Database Kafka Schema

Registers one schema subject in a DigitalOcean managed Kafka cluster's schema registry -- the complete `digitalocean_database_kafka_schema_registry` resource surface.

## Resources

| Resource | Purpose |
|---|---|
| `digitalocean_database_kafka_schema_registry.schema` | The registered subject: name, type, definition |

## Inputs

Generated `variables.tf` mirrors the `DigitalOceanDatabaseKafkaSchemaSpec` proto: `cluster` (flattened reference string), `subject_name`, `schema_type`, `schema`. Authentication uses `digitalocean_token` (sensitive).

## Outputs

Exactly the `DigitalOceanDatabaseKafkaSchemaStackOutputs` contract: `cluster_id`, `subject_name`.

## Behavior notes

- ALL arguments are create-only (the resource has no update function): any change is destroy+recreate and DROPS all previously registered versions of the subject.
- The definition is compared verbatim -- a whitespace-only reformat is a replacement.
- Import: excluded -- the provider's importer is defective at the pin (see `iac/import-map.yaml`).
