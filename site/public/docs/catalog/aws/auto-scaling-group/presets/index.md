---
title: "Presets"
description: "Ready-to-deploy configuration presets for Auto Scaling Group"
type: "preset-list"
componentSlug: "auto-scaling-group"
componentTitle: "Auto Scaling Group"
provider: "aws"
icon: "package"
order: 200
presets:
  - slug: "01-web-service-behind-alb"
    rank: "01"
    title: "Web Service Behind ALB"
    excerpt: "This preset runs a load-balanced web fleet: instances register into an `AwsLbTargetGroup` on launch, ELB health checks replace anything the application health check fails, CPU target tracking holds..."
  - slug: "02-spot-mixed-fleet"
    rank: "02"
    title: "Spot Mixed Fleet"
    excerpt: "This preset runs the classic cost architecture: two guaranteed On-Demand instances as the base, everything above them on Spot (`onDemandPercentageAboveBaseCapacity: 0`), drawn from four..."
  - slug: "03-scheduled-scale"
    rank: "03"
    title: "Scheduled Scale with Warm Pool"
    excerpt: "This preset shapes capacity around a business calendar: four instances minimum during weekday business hours, scale-to-zero overnight, and a warm pool of stopped, pre-initialized instances so the..."
  - slug: "04-reserved-fleet"
    rank: "04"
    title: "Reserved Fleet with Guaranteed Capacity"
    excerpt: "This preset runs a fleet on capacity you have already paid for. The group fills the targeted On-Demand Capacity Reservations first (`capacity-reservations-first`) and falls back to regular On-Demand..."
---

# Auto Scaling Group Presets

Ready-to-deploy configuration presets for Auto Scaling Group. Each preset is a complete manifest you can copy, customize, and deploy.
