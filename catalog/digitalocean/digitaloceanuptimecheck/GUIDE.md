# DigitalOcean Uptime Check -- Operational Guide

What experience with this component teaches that the field reference cannot.

## down vs down_global: page on global, watch regional

A `down` alert fires when ANY vantage region cannot reach the target -- which includes that region's own network weather. A `down_global` alert fires only when ALL regions agree the target is unreachable. Page humans on `down_global`; route `down` to a low-urgency channel if you want early regional signals. Probing from more regions makes the distinction sharper -- declare at least two.

## Always declare regions (the spec makes you)

DigitalOcean happily defaults the region set when omitted -- and the provider then reads the default back into state and tries to remove it on every subsequent plan, forever. That is why `regions` is required here even though the API can choose. Pick the regions your users actually connect from.

## latency needs a threshold; ssl_expiry wants one

A latency rule without a threshold would be sent as a silent zero -- an always-firing alert -- so validation requires it (in milliseconds). `ssl_expiry`'s threshold is DAYS before certificate expiry; give it at least your renewal pipeline's worst-case turnaround (14+ days), because an alert at zero days is a post-mortem, not a warning.

## ping targets are hosts, http(s) targets are URLs

`https://www.example.com` for http/https probes; `www.example.com` (or an IP) for ping. DigitalOcean enforces the pairing at request time -- the spec documents it rather than guessing your intent.

## Alert rules live and die with the check

Deleting the check deletes every alert rule under it -- there is nothing to clean up, and nothing survives to alert on a target you stopped probing. Renaming an alert row REPLACES that row (new id, fresh alert history on DigitalOcean's side); the check itself renames in place.

## Slack webhooks are credentials

The webhook URL lets anyone post to your channel. The spec marks it sensitive and both provisioners keep it out of plain-text state rendering -- prefer secret references over literals in committed manifests.

## What is deliberately NOT here

Metric alerts on droplets/load balancers/databases (that is the DigitalOceanMonitorAlert kind); authenticated probes, custom headers, and response-body assertions (DigitalOcean's Uptime API has none of them); and standalone alert objects pointing at existing checks -- the mutable parent id upstream is a corruption class, so rules are declared on their check, period.
