# OciDbSystem

Managed Oracle Database on OCI Virtual Machine or Bare Metal infrastructure, deployed and configured through Planton declarative manifests.

## Overview

OciDbSystem provisions an Oracle Cloud Infrastructure Database System — the compute, storage, DB Home, and initial database instance — as a single declarative resource. It maps directly to the `oci_database_db_system` Terraform/Pulumi resource and exposes the full creation-time configuration surface while keeping the manifest concise for common use cases.

This component handles fresh-creation scenarios only (`source=NONE`). Clone and restore workflows (`source=DB_BACKUP`, `DATABASE`, `DB_SYSTEM`) require different field sets and are excluded from this version.

## Purpose

- Provide a single YAML manifest that creates a fully operational Oracle Database environment (DB System + DB Home + Database) without writing procedural IaC code.
- Enable composability with other Planton OCI components via `StringValueOrRef` foreign key references for compartment, subnet, NSG, and KMS resources.
- Expose Oracle-specific knobs (edition, license model, storage layout, RAC, maintenance windows) as flat configuration fields rather than nested Pulumi/Terraform blocks.

## Key Features

- **Single-node and 2-node RAC** — set `nodeCount` to 1 or 2, with fault domain distribution for high availability.
- **Flexible compute shapes** — supports VM Standard, VM Flex, and Bare Metal Dense I/O shapes with configurable CPU core count.
- **Oracle Database editions** — Standard, Enterprise, Enterprise High Performance, and Enterprise Extreme Performance.
- **Licensing flexibility** — choose between BYOL (Bring Your Own License) and License Included.
- **Storage configuration** — control data storage size, data-vs-recovery percentage (BM), disk redundancy level (BM), storage volume performance mode, and storage management strategy (ASM or LVM).
- **Automatic backups** — configure backup windows and retention periods through the nested `dbBackupConfig` block.
- **KMS encryption** — encrypt data at rest at the DB System level and/or per-database via TDE with OCI Vault integration.
- **Network security groups** — attach NSGs to both the primary and backup VNICs.
- **Maintenance windows** — schedule patching with custom preferences including months, weeks, days, hours, patching mode (rolling/nonrolling), and lead time.
- **Diagnostic telemetry** — toggle diagnostic events, health monitoring, and incident log collection.
- **Freeform tags** — automatically applied from metadata (resource kind, ID, organization, environment, custom labels).

## Critical Constraints

- **Fresh creation only** — this component supports `source=NONE`. Cloning from backups, databases, or other DB Systems is not supported.
- **Mutually exclusive version fields** — exactly one of `dbHome.dbVersion` or `dbHome.databaseSoftwareImageId` must be provided. The proto validation rejects manifests where both or neither are set.
- **Immutable fields** — many fields force recreation when changed: `availabilityDomain`, `subnetId`, `hostname`, `databaseEdition`, `nodeCount`, `clusterName`, `faultDomains`, `domain`, `dataStoragePercentage`, `diskRedundancy`, `timeZone`, `sparseDiskgroup`, `storageVolumePerformanceMode`, `privateIp`, `dbHome.dbVersion`, `dbHome.databaseSoftwareImageId`, `dbHome.database.dbName`, `dbHome.database.characterSet`, `dbHome.database.ncharacterSet`, `dbHome.database.pdbName`, `dbHome.database.dbDomain`, and the TDE fields.
- **BM-only fields** — `dataStoragePercentage`, `diskRedundancy`, and `sparseDiskgroup` apply only to Bare Metal shapes. Setting them on VM shapes is a no-op or causes errors.
- **RAC prerequisites** — 2-node RAC requires `enterprise_edition_extreme_performance`, a `backupSubnetId`, and appropriate fault domain placement.
- **Admin password not returned** — `adminPassword` is write-only; the OCI API does not return it after creation.
- **Deprecated fields excluded** — `dbWorkload` (deprecated by Oracle November 2023) and Exadata-specific fields (`computeModel`, `computeCount`) are not exposed.

## Use Cases

1. **Development database** — single-node VM with Standard Edition, minimal storage, no backups. Tear down and recreate as needed.
2. **Staging environment** — Enterprise Edition with automatic backups and a pluggable database for application testing.
3. **Production RAC** — 2-node cluster with Extreme Performance edition, fault domain distribution, NSG isolation, KMS encryption, and a quarterly rolling maintenance window.
4. **BYOL migration** — bring existing Oracle licenses to OCI by setting `licenseModel: bring_your_own_license` and selecting the appropriate edition.

## Production Features

| Feature | How to enable |
|---------|---------------|
| High availability | Set `nodeCount: 2` with `enterprise_edition_extreme_performance` and distribute across `faultDomains` |
| Automatic backups | Set `dbHome.database.dbBackupConfig.autoBackupEnabled: true` with window and retention |
| Encryption at rest | Set `kmsKeyId` at the DB System level and/or `dbHome.database.kmsKeyId` with `vaultId` for TDE |
| Network isolation | Attach NSGs via `nsgIds` and `backupNetworkNsgIds` |
| Maintenance control | Configure `maintenanceWindowDetails` with `custom_preference`, patching mode, and schedule |
| Diagnostic telemetry | Enable via `dataCollectionOptions` sub-fields |
| Storage performance | Set `storageVolumePerformanceMode: high_performance` for lower latency |

---

© Planton. Licensed under [Apache-2.0](https://github.com/plantonhq/planton/blob/main/LICENSE).
