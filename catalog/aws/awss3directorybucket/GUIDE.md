# AwsS3DirectoryBucket — Operational Guide

Live-earned judgment lands here as proof runs and adopter operations teach it; the notes below are the forge-time seed.

## Zone IDs, never zone letters

`use1-az4` is the same physical zone in every AWS account; `us-east-1a` is a per-account shuffle. Express buckets take the ID form only — and co-location with compute means matching the compute subnet's zone ID (`aws ec2 describe-subnets` shows AvailabilityZoneId), not its letter.

## One zone means one zone

Express One Zone data has NO cross-AZ replica: a zone impairment takes the data offline, and a zone loss can lose it. It is a speed tier for reconstructible data (training shards, caches, intermediates) — the system of record stays in regional S3.

## The session API is the auth model

Express buckets authenticate through `s3express:CreateSession`, not per-object IAM — SDKs handle it transparently, but IAM policies must grant CreateSession on the bucket (a policy written for regular S3 object ARNs silently fails).

## Directory semantics change listing habits

Objects are organized by real directories, listing is prefix-constrained, and there is no lifecycle/versioning/replication surface — code written against regular S3's flat-namespace conventions may need its listing paths revisited.

## The name is infrastructure

The full `{base}--{zone_id}--x-s3` name encodes the zone: consumers hardcoding it survive zone-id changes only through the `bucket_name` output. Wire the output, never the literal.
