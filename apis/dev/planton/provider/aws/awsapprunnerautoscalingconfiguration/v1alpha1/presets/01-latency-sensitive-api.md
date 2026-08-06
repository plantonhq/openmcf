# Latency-Sensitive API

A scaling posture for user-facing APIs where tail latency matters more than baseline cost: three warm instances absorb traffic without cold starts, and a lowered concurrency ceiling gives every instance headroom.

## When to Use

- Production APIs on the critical path of user interactions
- Services whose p99 latency degrades when instances run near saturation

## What It Configures

- **`minSize: 3`** — three instances stay warm at all times (billed for memory only while idle), so bursts never wait on a cold start
- **`maxConcurrency: 50`** — instances scale out at half the default load, trading instance count for per-request headroom
- **`maxSize: 15`** — the cost ceiling during spikes

## What to Customize

- Replace `<aws-region>` with your region
- Reference from services via `autoScalingConfigurationArn` — one configuration tunes every service that adopts it
- Remember every value change registers a NEW revision; referencing services roll to it on their next deployment
