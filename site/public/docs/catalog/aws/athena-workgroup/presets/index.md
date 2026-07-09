---
title: "Presets"
description: "Ready-to-deploy configuration presets for Athena Workgroup"
type: "preset-list"
componentSlug: "athena-workgroup"
componentTitle: "Athena Workgroup"
provider: "aws"
icon: "package"
order: 200
presets:
  - slug: "01-basic-sql-workgroup"
    rank: "01"
    title: "Preset: Basic SQL Workgroup"
    excerpt: "A minimal Athena workgroup for interactive SQL analytics with query results stored in S3. All governance defaults apply — configuration enforcement is enabled, CloudWatch metrics are published."
  - slug: "02-encrypted-production"
    rank: "02"
    title: "Preset: Encrypted Production Workgroup"
    excerpt: "A production-grade Athena workgroup with SSE_KMS encryption, cost controls, and strict configuration enforcement."
  - slug: "03-spark-workgroup"
    rank: "03"
    title: "Preset: Spark Workgroup"
    excerpt: "An Athena workgroup configured for Apache Spark workloads (PySpark notebooks and Spark SQL). Requires an IAM execution role with appropriate permissions."
  - slug: "04-managed-results-zero-bucket"
    rank: "04"
    title: "Preset: Managed Results, Zero Bucket"
    excerpt: "An Athena workgroup that stores query results in AWS-managed storage — no S3 bucket to create, secure, or lifecycle. Results are retained for 24 hours and retrieved through Athena APIs..."
---

# Athena Workgroup Presets

Ready-to-deploy configuration presets for Athena Workgroup. Each preset is a complete manifest you can copy, customize, and deploy.
