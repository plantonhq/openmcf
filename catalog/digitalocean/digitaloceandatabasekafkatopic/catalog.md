# Kafka Topic on DigitalOcean

Creates a topic on a DigitalOcean managed Kafka cluster with the full per-topic configuration surface -- partitions, replication, cleanup policy, retention, and segment tuning. Integrates with Planton's Provider Connections for DigitalOcean API token management; the owning cluster is wired by reference.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **Kafka Topic** -- the named topic on the referenced cluster, with your partition count, replication factor, and tuned configuration

## Before You Deploy

### Planton Setup

- **DigitalOcean Provider Connection** -- an active connection in the Connect module with a DigitalOcean API token. Map it as the default for your environment, or specify it explicitly when creating the Cloud Resource.
- **Kafka Database Cluster** -- a DigitalOceanDatabaseCluster running the `kafka` engine (topics cannot live on other engines).

### DigitalOcean Account

- Nothing beyond the cluster: topics are free API objects on it.

## After You Deploy

Producers and consumers connect through the CLUSTER's connection outputs (host, port, credentials from its users); the topic is addressed by its `topic_name` output. Partition growth applies in place -- plan capacity with the topic's consumers in mind, since partitions can never be removed.
