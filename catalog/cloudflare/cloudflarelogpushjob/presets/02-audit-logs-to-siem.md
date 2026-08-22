# Audit logs to a SIEM

The compliance shape: account audit logs streamed to an external SIEM every 30 seconds, so a record of who changed what survives independently of the Cloudflare account itself. Note the account scope (audit logs are not a zone dataset) and the `ownership_challenge` token: an HTTPS destination must prove ownership first, so deploy once with `generate_ownership_challenge: true`, read the token from the challenge Cloudflare posts, then apply this shape with it filled in.
