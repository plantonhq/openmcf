# Kafka Schema on DigitalOcean

Registers a schema subject (Avro, JSON Schema, or Protobuf) in a DigitalOcean managed Kafka cluster's schema registry, so producers and consumers agree on message structure. Integrates with Planton's Provider Connections for DigitalOcean API token management; the owning cluster is wired by reference.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **Schema Registry Subject** -- the named subject on the referenced cluster's registry, carrying your schema definition

## Before You Deploy

### Planton Setup

- **DigitalOcean Provider Connection** -- an active connection in the Connect module with a DigitalOcean API token. Map it as the default for your environment, or specify it explicitly when creating the Cloud Resource.
- **Kafka Database Cluster** -- a DigitalOceanDatabaseCluster running the `kafka` engine (the registry exists only on Kafka clusters).

### DigitalOcean Account

- Nothing beyond the cluster: schema subjects are free API objects on it.

## After You Deploy

Producers and consumers fetch the schema through the cluster's registry endpoint using the `subject_name` output. Plan schema EVOLUTION carefully before depending on this resource: any change to the definition replaces the subject and drops its prior versions.
