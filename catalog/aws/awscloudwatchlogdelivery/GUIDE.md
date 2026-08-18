# AwsCloudwatchLogDelivery — Operational Guide

Live-earned judgment lands here as proof runs and adopter operations teach it; the notes below are the forge-time seed.

## One delivery per (source, destination-type)

AWS accepts at most one delivery from a source to each destination TYPE — S3 plus Firehose is fine, two S3 destinations from one source is a ConflictException. Fan out to multiple buckets via Firehose or replicate downstream.

## The CloudFront prefix is AWS's, not yours

For CloudFront sources AWS prepends `AWSLogs/{account-id}/CloudFront/` to the S3 suffix path server-side. Configure only your own segment in `suffix_path`; the provider strips the AWS-added prefix on reads, so state never fights it.

## The legacy arm's policy outlives its resource

Destroying only the cross-account destination's POLICY is a provider no-op — the policy stays attached at AWS. Destroying the destination itself removes everything. The kind folds both, so a normal teardown is clean; remember this only when hand-pruning.

## The first cross-account create is eventually consistent

CloudWatch Logs validates it can write a test message through the role at PutDestination time; a freshly created role fails that for up to two minutes while the `logs.amazonaws.com` trust propagates. Both engines retry through the window — a slow first apply is normal.

## Destination policies need a second account to mean anything

The vended destination's `policy` grants OTHER accounts `logs:CreateDelivery`. Same-account pipelines never need it; set it only on shared central destinations.
