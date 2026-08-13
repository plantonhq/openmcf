---
title: "Presets"
description: "Ready-to-deploy configuration presets for S3 Bucket"
type: "preset-list"
componentSlug: "s3-bucket"
componentTitle: "S3 Bucket"
provider: "aws"
icon: "package"
order: 200
presets:
  - slug: "01-private-encrypted"
    rank: "01"
    title: "Private Encrypted Bucket"
    excerpt: "This preset creates a fully private, versioned S3 bucket encrypted with a customer-managed KMS key. It is the right starting point for application data, backups, and anything holding real information."
  - slug: "02-public-static-website"
    rank: "02"
    title: "Public Static Website Bucket"
    excerpt: "This preset creates a bucket that serves a static website directly over S3's website endpoint. Public access is granted deliberately and precisely: the policy grants read on objects, and only the two..."
  - slug: "03-log-archive-lifecycle"
    rank: "03"
    title: "Log Archive with Lifecycle Tiering"
    excerpt: "This preset creates a private bucket purpose-built as a log/archive destination: objects tier down through cheaper storage classes as they age and are deleted after a year, with no manual..."
  - slug: "04-governed-data-lake"
    rank: "04"
    title: "Governed Data Lake"
    excerpt: "This preset creates a private, versioned data-lake bucket with the full governance toolkit switched on: scheduled inventory reports, storage-class analytics, request metrics, and S3 Metadata tables —..."
---

# S3 Bucket Presets

Ready-to-deploy configuration presets for S3 Bucket. Each preset is a complete manifest you can copy, customize, and deploy.
