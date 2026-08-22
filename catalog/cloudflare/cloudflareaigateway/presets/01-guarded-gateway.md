# Guarded production gateway

A production-shaped gateway: five-minute response caching, a sliding 1000-requests-per-minute limit, exponential retries, gateway authentication required on every call, a prompt guardrail (category P1 blocked), and a $50/day spend cap.

**When to use it:** the first gateway an organization puts in front of its LLM traffic -- cost control and abuse protection from day one.

**What to change:** the `gateway_id` (it becomes the endpoint URL -- pick it like a domain name), the spend cap's `limit`/`window`, and the guardrail categories per your content policy. Every spend rule keeps its own unique `id`.
