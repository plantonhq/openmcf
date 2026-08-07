---
title: "Host and Path Fan-Out"
description: "The classic global external Application Load Balancer routing table: one host rule sends `www.example.com` traffic into a path matcher that longest-prefix matches `/api/*` and `/static/*` to..."
type: "preset"
rank: "01"
presetSlug: "01-host-path-fanout"
componentSlug: "url-map-on-google-cloud"
componentTitle: "URL Map on Google Cloud"
provider: "gcp"
icon: "package"
order: 1
---

# Host and Path Fan-Out

The classic global external Application Load Balancer routing table: one host rule sends `www.example.com` traffic into a path matcher that longest-prefix matches `/api/*` and `/static/*` to different backends, with a catch-all default for everything else.

## When to Use

- One domain serving both a dynamic API and static assets from different backends
- Splitting `/api/*` to a Cloud Run or instance-group backend while `/static/*` goes to a backend bucket with CDN
- A single URL map shared by a target HTTPS proxy in front of a global anycast IP

## Remix Notes

- Reference `GcpBackendService` and `GcpBackendBucket` resources via `valueFrom` on `defaultService`, `pathRules[].service`, and bucket targets instead of literal self-links.
- Add `tests[]` entries for each critical path — GCP evaluates them at update time and blocks a change that would silently break routing.
- For many hosts, add more `hostRules` entries pointing at different path matchers (e.g. a separate matcher for an admin subdomain).
