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
    excerpt: "Zip-backed Lambda with Node.js 22.x, composed AwsIamRole reference, and S3 deployment package — function name from metadata.name."
  - slug: "02-container-basic"
    rank: "02"
    title: "Container-Based Lambda Function"
    excerpt: "ECR container image deployment with ARM64 architecture — no runtime or handler fields; image defines the entrypoint."
---

# Lambda Presets

Ready-to-deploy configuration presets for Lambda. Each preset is a complete manifest you can copy, customize, and deploy.
