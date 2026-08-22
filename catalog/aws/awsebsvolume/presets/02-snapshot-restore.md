# Snapshot Restore

Rebuild a volume from a snapshot in a DIFFERENT zone — the only way volumes "move". Size comes from the snapshot; the full initialization rate (300 MiB/s) hydrates blocks eagerly so first reads hit full performance instead of lazy-loading from S3.
