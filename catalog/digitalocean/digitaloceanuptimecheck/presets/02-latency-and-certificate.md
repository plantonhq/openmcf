# API Latency and Certificate Expiry

This preset layers two quality signals on one https probe: a latency rule that fires when responses stay slow, and an `ssl_expiry` rule that warns three weeks before the certificate lapses -- the renewal safety net that turns an outage into a calendar item.

## When to Use

- APIs with latency expectations, not just up/down status
- Any https endpoint whose certificate renewal you want watched from outside your own tooling

## Key Configuration Choices

- **`threshold: 800` milliseconds over `5m`** -- fires on sustained slowness, not one slow request; tune to your latency budget.
- **`threshold: 21` days for ssl_expiry** -- at least your renewal pipeline's worst-case turnaround; an alert at zero days is a post-mortem.
- **Two rules, one check** -- each alert row is its own object with its own channels; add a `down_global` row for availability paging.

## What You Get

One probe carrying both rules, destroyed together with the check, with `check_id` exported.
