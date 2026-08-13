# Oracle Enterprise Edition (BYOL)

This preset creates a Multi-AZ Oracle EE instance under bring-your-own-license, with the two engine-configuration surfaces the instance owns itself: a parameter group built from inline parameters and an option group activating STATSPACK and Transparent Data Encryption. Character sets are pinned at create time and the master password is AWS-managed.

## When to Use

- Enterprise Oracle workloads moving onto RDS under an existing Oracle license agreement
- Regulated data that requires TDE at-rest encryption inside the database engine (in addition to storage-level encryption)
- Teams that want Oracle engine tunables and options as reviewable manifest surface instead of console clicks

## Key Configuration Choices

- **Instance-owned groups** (`parameters`, `options`) -- the module manages a dedicated parameter group and option group named after the instance, deriving the family (`oracle-ee-19`) and major version from the pinned engine version. Bringing existing groups by name (`parameterGroupName` / `optionGroupName`) is the mutually exclusive alternative.
- **TDE and STATSPACK** -- option-group options; TDE is Oracle EE's at-rest encryption for regulated schemas, STATSPACK the classic performance-diagnostics pack. Options with ports or settings declare them per entry.
- **License model** (`bring-your-own-license`) -- Oracle EE's only RDS model; `oracle-se2` also offers `license-included`.
- **Character sets** (`AL32UTF8` / `AL16UTF16`) -- create-time only; changing them later means a new instance.
- **Multi-AZ + deletion protection + named final snapshot** -- the availability and safety posture for anything holding enterprise data.

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `subnet-replace-with-private-az1` | Private subnet in the first Availability Zone | `AwsSubnet` status outputs or the AWS VPC console |
| `subnet-replace-with-private-az2` | Private subnet in the second Availability Zone | `AwsSubnet` status outputs or the AWS VPC console |
| `sg-replace-with-database-sg` | Security group allowing port 1521 from the application tier | `AwsSecurityGroup` status outputs or the AWS EC2 console |

## Related Presets

- **01-postgresql-production** -- The open-source production baseline
- **05-mysql-s3-migration** -- Migrating self-managed MySQL data onto RDS
