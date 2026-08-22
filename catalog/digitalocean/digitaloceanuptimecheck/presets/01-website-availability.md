# Website Availability Check

This preset probes a public site over https from all four DigitalOcean vantage regions and pages only when EVERY region agrees the site is down -- the `down_global` signal that filters out single-region network weather.

## When to Use

- The baseline "is the site up" monitor every public endpoint deserves
- Measuring what anonymous users actually see, from outside your own infrastructure

## Key Configuration Choices

- **All four regions** -- more vantage points make `down_global` sharper and per-region history richer.
- **`down_global` over `down`** -- page humans on global outages; add a separate `down` rule routed to a low-urgency channel if you want early regional signals.
- **https probe** -- also unlocks `ssl_expiry` rules (see the latency-and-certificate preset).

## What You Get

A check with uptime history in the control panel's Monitoring -> Uptime section and its `check_id` exported.
