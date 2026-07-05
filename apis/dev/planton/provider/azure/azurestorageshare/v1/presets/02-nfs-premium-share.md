# NFS Premium Share

This preset creates an NFS v4.1 share on a premium FileStorage account
-- SSD-backed provisioned performance with real POSIX semantics (hard
links, symlinks, chmod) for Linux workloads that SMB semantics would
break.

## When to Use

- Linux applications that need POSIX file semantics (databases' shared
  config, legacy Unix apps, ML training data)
- Latency-sensitive file workloads that outgrow standard-tier
  performance

## Key Configuration Choices

- **`enabledProtocol: NFS`** -- fixed at creation; requires the parent
  to be a premium FileStorage account (`accountKind: FILE_STORAGE`,
  `accountTier: PREMIUM`)
- **`accessTier: PREMIUM`** -- the only legal tier on FileStorage
  accounts
- **`quotaGb: 100`** -- the FileStorage floor; premium bills provisioned
  capacity whether used or not, and performance scales with size
- **Network reachability**: NFS mounts work only over private paths --
  plan a private endpoint or VNet rules on the account

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<filestorage-account-resource-name>` | A PREMIUM FileStorage-kind AzureStorageAccount's Planton resource name | Your storage composition |
| `<share-name>` | 3-63 lowercase letters/digits/hyphens | Your naming convention |
| `<workload>` | What mounts this volume | Your workload inventory |
