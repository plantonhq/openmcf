# Training Data Import Link

Read-side association: an S3 training-data prefix appears as `/datasets` on the Lustre file system, tracks the bucket continuously, and hydrates the existing contents at creation.

## When to Use

- ML training corpora stored in S3 that GPU jobs read at Lustre speed
- Analytics datasets landed in S3 by upstream pipelines
- Any S3-first data flow where compute needs POSIX access to the current bucket state

## What It Configures

- **`/datasets` → `s3://.../datasets/`** — one directory subtree, one repository
- **Full auto-import** — `NEW`, `CHANGED`, and `DELETED` objects keep the namespace in sync
- **Batch import at creation** — the existing bucket contents appear immediately (metadata only; file data hydrates on first read)
- **No auto-export** — a read-side link; nothing written on the file system flows back to S3

## What to Customize

- Replace placeholders: `<aws-region>`, `<training-data-bucket>`, and the file-system reference name
- Narrow `auto_import_events` (e.g., `[NEW]`) if deletions in the bucket should not remove files
- Raise `imported_file_chunk_size` for very large files striped across many disks
