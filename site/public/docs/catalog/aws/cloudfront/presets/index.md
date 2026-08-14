---
title: "Presets"
description: "Ready-to-deploy configuration presets for CloudFront"
type: "preset-list"
componentSlug: "cloudfront"
componentTitle: "CloudFront"
provider: "aws"
icon: "package"
order: 200
presets:
  - slug: "01-s3-static-website"
    rank: "01"
    title: "S3 Static Website with Origin Access Control"
    excerpt: "This preset serves a private S3 bucket through CloudFront using an Origin Access Control (OAC) -- the modern way to front S3. The bucket never becomes public: CloudFront signs its origin requests..."
  - slug: "02-custom-domain-cdn"
    rank: "02"
    title: "Custom Domain CDN with ACM Certificate"
    excerpt: "This preset serves a private S3 bucket through CloudFront on your own domain: an alias (CNAME) backed by an ACM certificate, with an Origin Access Control keeping the bucket private."
  - slug: "03-blue-green-continuous-deployment"
    rank: "03"
    title: "Blue/Green Rollout with Continuous Deployment"
    excerpt: "This preset stages a CloudFront configuration change on real production traffic before promoting it: the primary distribution owns a continuous-deployment policy that routes a weighted slice of..."
---

# CloudFront Presets

Ready-to-deploy configuration presets for CloudFront. Each preset is a complete manifest you can copy, customize, and deploy.
