---
title: "Presets"
description: "Ready-to-deploy configuration presets for Job"
type: "preset-list"
componentSlug: "job"
componentTitle: "Job"
provider: "kubernetes"
icon: "package"
order: 200
presets:
  - slug: "01-database-migration"
    rank: "01"
    title: "Database Migration Job"
    excerpt: "This preset runs a schema migration to completion: one pod, one success required, a small retry budget, and automatic cleanup a day after the Job finishes. It is the canonical one-shot Job pattern —..."
  - slug: "02-parallel-batch"
    rank: "02"
    title: "Parallel Batch Job (Indexed)"
    excerpt: "This preset partitions work across 10 numbered completions, running 3 pods at a time. Each pod receives its completion index (0–9) through the `batch.kubernetes.io/job-completion-index` annotation..."
  - slug: "03-resilient-batch"
    rank: "03"
    title: "Resilient Batch Job (Failure Policy)"
    excerpt: "This preset classifies pod failures before they consume the retry budget. Two rules, evaluated in order: an application-signaled unrecoverable error (exit code 42) fails the Job immediately, and..."
---

# Job Presets

Ready-to-deploy configuration presets for Job. Each preset is a complete manifest you can copy, customize, and deploy.
