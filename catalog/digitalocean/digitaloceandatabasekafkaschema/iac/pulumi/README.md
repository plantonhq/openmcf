# Pulumi Module: DigitalOcean Database Kafka Schema

Registers one schema subject in a DigitalOcean managed Kafka cluster's schema registry -- the complete `digitalocean_database_kafka_schema_registry` resource surface. Behavioral parity with the Terraform module is the contract.

## Resources

| Resource | Purpose |
|---|---|
| `digitalocean.DatabaseKafkaSchemaRegistry` | The registered subject: name, type, definition |

## Inputs

`DigitalOceanDatabaseKafkaSchemaStackInput`: the target `DigitalOceanDatabaseKafkaSchema` resource and the DigitalOcean provider config (API token).

## Outputs

Exactly the `DigitalOceanDatabaseKafkaSchemaStackOutputs` contract: `cluster_id`, `subject_name`.

## Behavior notes

- ALL arguments are create-only upstream: any change is a replacement in Pulumi too, and it DROPS all previously registered versions of the subject.
- The definition is compared verbatim -- keep the schema string byte-stable.
