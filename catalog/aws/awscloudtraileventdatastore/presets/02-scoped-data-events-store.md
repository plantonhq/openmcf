# Scoped Data Events Store

This preset creates a cost-disciplined store for one question — who
wrote what into one bucket — on fixed-retention pricing with a 90-day
queryable window.

## When to Use

- Data-access auditing for a specific sensitive bucket
- Short-horizon investigations where 7-year retention is waste

## What You Get

- Only S3 object WRITES under the named bucket, nothing else ingested
- A 90-day queryable window on FIXED_RETENTION_PRICING (cheaper
  ingest, retention billed)
- Multi-region capture of the scoped events (the AWS default)

## Customize

- Point the `resources.ARN` prefix at your bucket (keep the trailing
  `/`)
- Drop the `readOnly` condition to also ingest reads (bigger bill)
- Fixed pricing can later move to extendable, never the reverse —
  pick deliberately
