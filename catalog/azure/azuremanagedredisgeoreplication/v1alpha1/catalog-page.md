# Azure Managed Redis Geo-Replication

Links Azure Managed Redis instances into an active (multi-primary)
geo-replication group: every member accepts writes in its own region
and Azure merges the datasets with conflict-free semantics.

## What Gets Created

When you deploy an AzureManagedRedisGeoReplication resource, Planton
provisions:

- **Geo-replication group membership** -- an
  `azurerm_managed_redis_geo_replication` linking the referenced
  Managed Redis instances into one write-anywhere replica set

## Prerequisites

- **Azure credentials** configured via environment variables or Planton
  provider config
- **Two or more AzureManagedRedis instances** (up to five per group) in
  different regions, each declaring the SAME
  `defaultDatabase.geoReplicationGroupName`, each `BALANCED_B3` or
  larger, with no persistence and only geo-compatible modules
  (RediSearch/RedisJSON)

## Quick Start

Create a file `geo-group.yaml`:

```yaml
apiVersion: azure.planton.dev/v1alpha1
kind: AzureManagedRedisGeoReplication
metadata:
  name: global-cache-group
  annotations:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: dev.AzureManagedRedisGeoReplication.global-cache-group
spec:
  managedRedisId:
    valueFrom:
      kind: AzureManagedRedis
      name: global-cache-east
      fieldPath: status.outputs.managed_redis_id
  linkedManagedRedisIds:
    - valueFrom:
        kind: AzureManagedRedis
        name: global-cache-west
        fieldPath: status.outputs.managed_redis_id
```

Deploy:

```shell
planton pulumi up --manifest geo-group.yaml
```

## Spec Highlights

- **One resource manages the whole group** -- linking is reciprocal;
  never create one per member
- `linked_managed_redis_ids` -- add an ID to link a member, remove one
  to evacuate just that region (it keeps its data)
- **Deleting the resource unlinks all members** -- each becomes an
  independent instance with its own copy of the data

## Stack Outputs

| Output | Description |
| --- | --- |
| `geo_replication_id` | The group's resource ID (the managing cluster's ARM ID -- the group has no ARM object of its own) |
