# Cross-Region DR Copy

Disaster recovery in one resource: copy the primary region's snapshot into the DR region, re-encrypted under the DR region's own KMS key (copying is the ONLY way snapshots change keys). Restore volumes from it the day the primary region has a bad day.
