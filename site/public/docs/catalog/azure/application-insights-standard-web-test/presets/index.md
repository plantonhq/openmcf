---
title: "Presets"
description: "Ready-to-deploy configuration presets for Application Insights Standard Web Test"
type: "preset-list"
componentSlug: "application-insights-standard-web-test"
componentTitle: "Application Insights Standard Web Test"
provider: "azure"
icon: "package"
order: 200
presets:
  - slug: "01-endpoint-availability"
    rank: "01"
    title: "Multi-Region Endpoint Availability"
    excerpt: "This preset creates a standard web test that pings a public endpoint every 5 minutes from three regions, asserts a 200 with a valid SSL certificate at least 30 days from expiry, and retries a failed..."
  - slug: "02-content-check"
    rank: "02"
    title: "Response-Content Health Check"
    excerpt: "This preset creates a web test that not only checks for a 200 but also asserts the response body contains an expected healthy payload (for example `\"status\":\"ok\"`). A endpoint can return 200 while..."
---

# Application Insights Standard Web Test Presets

Ready-to-deploy configuration presets for Application Insights Standard Web Test. Each preset is a complete manifest you can copy, customize, and deploy.
