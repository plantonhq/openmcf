---
title: "Presets"
description: "Ready-to-deploy configuration presets for Lambda"
type: "preset-list"
componentSlug: "lambda"
componentTitle: "Lambda"
provider: "aws"
icon: "package"
order: 200
presets:
  - slug: "01-zip-basic"
    rank: "01"
    title: "Zip-Based Lambda Function"
    excerpt: "This preset creates a Lambda function deployed from a zip archive in S3. The function name comes from `metadata.name` (create-time immutable in AWS). It uses the Node.js 22.x runtime with 256 MB..."
  - slug: "02-container-basic"
    rank: "02"
    title: "Container-Based Lambda Function"
    excerpt: "This preset creates a Lambda function from a container image in ECR. The function name comes from `metadata.name`. Runtime and entrypoint are defined by the image — leave `runtime` and `handler`..."
---

# Lambda Presets

Ready-to-deploy configuration presets for Lambda. Each preset is a complete manifest you can copy, customize, and deploy.
