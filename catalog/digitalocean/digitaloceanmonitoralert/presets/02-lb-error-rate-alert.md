# Load Balancer 5xx Error-Rate Alert

This preset pages when a load balancer's 5xx error rate climbs -- the user-facing "the site is broken" signal, wired to the balancer by reference and delivered to both email and a Slack incidents channel.

## When to Use

- Any production load balancer: 5xx rate is the closest built-in metric to user pain
- Pairing with a droplet CPU policy so cause (hot backends) and effect (errors) page together

## Key Configuration Choices

- **Balancer by reference** (`valueFrom`) -- wires to a DigitalOceanLoadBalancer in the same chart or environment; swap in a literal UUID for an existing balancer.
- **`value: 5` percent over `5m`** -- catches real degradation without paging on a single bad request; tune to traffic volume.
- **Slack webhook URL is a secret** -- replace the placeholder with your real webhook; both provisioners keep it out of plain-text state.

## What You Get

One policy on the referenced balancer, its `alert_id` exported for reference.
