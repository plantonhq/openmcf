# GcpCloudSqlDatabase Guide

The judgment this guide protects: databases are the unit of application
tenancy on a shared instance — cheap to create, expensive to drop
carelessly. Model one database node per application and let each carry
its own lifecycle; never funnel schema decisions through the instance.

## Charset and collation are forever (mostly)

MySQL takes any supported pair — choose `utf8mb4` and never think about
it again. PostgreSQL databases are born UTF8 with an OS-locale collation;
SQL Server ignores charset entirely and takes a SQL Server collation
name. The API validates the combination at deploy time, so a wrong pair
fails the apply, not the application. Changing collation later is a
dump-and-reload in every engine — decide once.

## The drop that will not drop

PostgreSQL refuses to drop a database while clients hold connections —
which turns a routine teardown into a hung destroy at exactly the wrong
moment. `deletionPolicy: ABANDON` is the documented answer: the node
leaves IaC management, the database keeps serving, and the actual drop
happens on your schedule with connections drained. `PREVENT` is the
standing guard for databases whose loss would be an incident.

## Destroy stance

Dropping the resource drops the database AND its data — the instance's
backups and PITR are the only recovery path, so the database node's real
safety net is configured on its `GcpCloudSql` parent. Verify the
instance's backup posture before trusting any database teardown.

## On the diagram

A leaf under its `GcpCloudSql` instance (referenced by `instance`),
usually beside the `GcpCloudSqlUser` that owns its schema. A database
node without a paired user is often a manifest that left credentials to
hand management.

## Pairs well with

- `GcpCloudSql` — the instance; its backup posture is this database's
  recovery story.
- `GcpCloudSqlUser` — one per application database; grant inside the
  database via migrations.
