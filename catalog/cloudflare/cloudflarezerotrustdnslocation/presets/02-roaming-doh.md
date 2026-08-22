# Roaming DoH location

The roaming shape: DNS-over-HTTPS only, token-gated (source networks can't gate devices that move), with a 5-minute TTL cap so policy changes propagate quickly. Clients embed the `doh_subdomain` output in their resolver URL. Keep `max_ttl` declared forever -- an update that omits it resets the behavior to inherit.
