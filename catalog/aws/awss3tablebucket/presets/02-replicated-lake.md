# Replicated Lake

The compliance shape: KMS-encrypted at rest (remember the maintenance service's key grant) with bucket-level replication of every table to a us-west-2 replica bucket. The replication role needs the s3tables replication trust and both buckets in scope.
