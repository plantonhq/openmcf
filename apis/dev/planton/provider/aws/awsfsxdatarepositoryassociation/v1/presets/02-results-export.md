# Results Export Link

Write-side association: everything jobs write under `/output` on the Lustre file system flows back to an S3 results prefix automatically.

## When to Use

- Model checkpoints and artifacts that downstream systems consume from S3
- Processed datasets handed off to S3-based pipelines
- Durable capture of results from scratch-style processing directories

## What It Configures

- **`/output` → `s3://.../output/`** — one directory subtree, one repository
- **Auto-export of new and changed files** — results land in S3 as they are written (deletions on the file system do NOT delete objects)
- **No auto-import** — the file system does not track bucket-side changes

## What to Customize

- Replace placeholders: `<aws-region>`, `<results-bucket>`, and the file-system reference name
- Add `DELETED` to `auto_export_events` if file-system deletions should propagate to the bucket
- Pair with the training-data import preset on the same file system for a full in/out pipeline
