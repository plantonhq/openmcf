# Nightly Cron to Lambda

The workhorse: a timezone-aware cron firing a Lambda at 2am Eastern (daylight saving handled by AWS), with bounded retries and a monitored dead-letter queue. The target ARN is a bare polymorphic reference — the `kind:` on the valueFrom is required, not decoration.
