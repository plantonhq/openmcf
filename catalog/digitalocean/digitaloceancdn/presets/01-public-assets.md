# Public Assets CDN

This preset fronts an assets bucket with DigitalOcean's edge network and a one-day cache TTL -- the standard shape for fingerprinted build artifacts (hashed JS/CSS bundles, versioned images) that never change in place.

## When to Use

- Static site or app assets whose filenames change on every release
- Any public bucket content where global load speed matters

## Key Configuration Choices

- **Origin by reference** -- the endpoint follows the bucket resource; a chart deploys the bucket and its CDN together.
- **Long TTL (86400)** -- safe exactly because fingerprinted files never mutate; content updated in place needs a much shorter TTL.
- **No custom domain** -- consumers use the endpoint output directly; add the certificate + custom-domain pair when a branded URL matters.

## What You Get

A free edge distribution serving the bucket worldwide, with bandwidth drawn from the Spaces subscription's existing allowance.
