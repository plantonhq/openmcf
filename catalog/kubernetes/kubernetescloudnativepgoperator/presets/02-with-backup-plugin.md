# With Backup Plugin

This preset installs the CloudNativePG operator PLUS the Barman Cloud
plugin — the backup-capable posture every KubernetesPostgres backup
block depends on. The plugin is a separate Helm release
(`plugin-barman-cloud`) in the same namespace, installed after the
operator; its internal TLS is issued by cert-manager, so cert-manager
(KubernetesCertManager) must be on the cluster before this deploys.

## When to Use

- Any cluster whose databases will declare backup blocks — WAL archiving
  and scheduled base backups to S3/GCS/Azure-Blob/S3-compatible stores
- The production default: databases without continuous backups are a
  deliberate exception, not a posture

## Key Configuration Choices

- **`barman_cloud_plugin.enabled: true`** — CloudNativePG delegates
  object-store backups to this plugin (its built-in object-store support
  is deprecated upstream); one plugin installation serves every
  database's backup blocks. WHERE backups land is declared per database
  in `KubernetesPostgres.spec.backup`
- **cert-manager is a hard prerequisite** — the plugin chart renders
  cert-manager Certificate resources unconditionally; without
  cert-manager the release fails to install and rolls back cleanly
  (atomic)
- **Independent version pins** — the plugin chart (0.7.0 = plugin
  v0.13.0) versions separately from the operator chart (0.29.0 =
  operator 1.30.0); each release carries its own pin
- **Plugin resources declared** — like the operator chart, the plugin
  chart ships no requests/limits by default
- **Everything else as 01-standard** — sized operator, cluster-wide
  watch, `system-cluster-critical` priority

## Placeholders to Replace

None — this preset deploys as-is (cert-manager must already be on the
cluster).

## Related Presets

- **01-standard** — the operator alone, for clusters that do not need
  object-store backups yet
