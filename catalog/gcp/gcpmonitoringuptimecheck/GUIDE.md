# GcpMonitoringUptimeCheck Guide

The judgment this guide protects: an uptime check that exists but pages
nobody is availability theater. The check is HALF of a monitor — the
alert policy filtering on its `uptime_check_id` is the other half, and
the two ship together or the green dashboard lies.

## The check-plus-alert contract

`uptime_check_passed` is just a metric until a GcpMonitoringAlertPolicy
watches it. The canonical pairing: a threshold condition on
`monitoring.googleapis.com/uptime_check/check_passed` with a
`metric.label.check_id` filter on this check's `uptime_check_id` output,
`COMPARISON_GT` on a fraction-true reduction. Wire the id via valueFrom —
hand-copied check ids rot when the check is recreated.

## validate_ssl is the certificate monitor

`useSsl` alone accepts an EXPIRED certificate — the TLS handshake is
performed but not judged. `validateSsl: true` is what turns the probe
into the certificate-expiry monitor most teams assume they have. The
spec defaults both to false because GCP does; production presets flip
both deliberately.

## Content matchers catch the lying 200

Error pages served with status 200 are the classic silent outage. A
`CONTAINS_STRING` matcher on a known-good token (or a JSON-path matcher
on a health endpoint's `$.status`) fails the probe on body truth, not
transport truth. All matchers must pass — they AND together.

## Region math

`selectedRegions` must cover at least three checker locations or GCP
rejects the config. Leaving it EMPTY probes from all regions — more
coverage, no config to maintain, and the default every preset keeps.

## The recreate trap

`uptime_check_id` changes when the check is recreated, and alert-policy
filters that hard-code the old id keep evaluating against a metric that
no longer updates — the alert goes permanently green. valueFrom wiring
heals this automatically; hand-wired filters do not.
