# Availability SLO

Three nines of good requests over a rolling 30 days, on a custom
service created in the same apply — the standard first SLO for a
service whose truth lives in its own metrics.

## What it configures

- A `customService` container named by the manifest — no dependency on
  GCP auto-detection.
- A good/total ratio from two log-based metrics (the GcpLogMetric
  pairing); GCP derives the bad count.
- `deletionPolicy: PREVENT` — burn-rate alerts will reference this SLO;
  deleting it silently breaks them.

## Adjust before deploying

- **The two filters** — point them at metrics that actually exist
  (log-based metrics chart as
  `logging.googleapis.com/user/{metric_name}`). Prefer good+total over
  bad-counting: timeouts never write a bad event.
- **goal** — 0.999 is a starting point; set it from measured history,
  not aspiration (an SLO you already violate is noise from day one).

## When to choose something else

For "fast enough" rather than "up", start from the **Latency SLO**
preset — availability and latency are separate objectives on purpose.
