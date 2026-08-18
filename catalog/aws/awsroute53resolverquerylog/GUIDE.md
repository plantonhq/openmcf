# AwsRoute53ResolverQueryLog — Operational Guide

Live-earned judgment lands here as proof runs and adopter operations teach it; the notes below are the forge-time seed.

## The association can fail AFTER a clean apply

Association is asynchronous: AWS accepts it, then tries to write. An unwritable destination (missing S3 bucket policy, Firehose permissions) flips the association to FAILED minutes later with an error code. The E2E verifier asserts ACTIVE for exactly this reason — check association status, not just existence, when logs do not appear.

## Pick the destination by retention economics

CloudWatch Logs for interactive queries and alarms (priciest per GB), S3 for cheap archival and Athena, Firehose for streaming into SIEM pipelines. High-traffic VPCs generate serious query volume — a busy VPC's DNS logs in CloudWatch can out-cost the workloads' own logs.

## Volume control is association-level only

There is no filtering: an associated VPC logs every query. The only dial is which VPCs associate. Log the VPCs under investigation or compliance scope, not the whole estate by reflex.

## Replacement is safe for data

Changing name or destination replaces the configuration, but log data already written stays in the destination — replacement costs continuity (a gap during the swap), never history.

## One VPC, many configs — do not

A VPC can associate to at most one query-log config per account. A second config's association for the same VPC fails; keep one config per logging destination strategy, associated wide.
