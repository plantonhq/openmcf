# S3 Data-Events Trail

This preset records WHO WROTE WHAT into one S3 bucket — object-level
write auditing scoped tightly enough to keep the per-event data-event
bill deliberate.

## When to Use

- Auditing writes to a sensitive bucket (a data lake, a compliance
  archive)
- Alongside a management-events trail (data events are a separate
  scope)

## What You Get

- Every non-read S3 object call under the chosen prefix, as CloudTrail
  data events
- Nothing else — no management events, no other buckets, no reads

## Customize

- Replace `<data-lake-bucket>` with the audited bucket (the trailing
  `/` scopes to its objects); add more `startsWith` entries for more
  buckets
- Flip `readOnly` equals to `["true"]` (or drop the selector) to
  include reads — expect order-of-magnitude more events
- The delivery bucket needs the same `cloudtrail.amazonaws.com`
  policy as the compliance preset

## Composing

Pair with the compliance-audit-trail preset: management events ride
the free first copy there, while this trail carries only the metered
data events.
