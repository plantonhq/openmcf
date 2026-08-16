# Public Static Site Assets

This preset creates a public-read Spaces bucket for a website's static assets, with a CORS rule letting the site's browser code fetch them cross-origin and cache preflight responses for an hour.

## When to Use

- Hosting images, scripts, stylesheets, fonts, or downloads a website serves directly from Spaces
- Single-page applications fetching assets from the bucket's endpoint on a different origin
- Anything meant to be read anonymously; keep private data in a PRIVATE bucket instead

## Key Configuration Choices

- **Public-read ACL** (`accessControl: PUBLIC_READ`) -- anyone can GET objects; only credentialed clients can write.
- **CORS rule** (`corsRules`) -- allows `GET`/`HEAD` from the site's origin, lets browsers read the `ETag` header for cache validation, and caches the preflight for 3600 seconds. CORS is managed through DigitalOcean's standalone CORS-configuration resource, so rules round-trip with real drift detection.
- **No versioning** -- static assets are rebuilt and re-uploaded by pipelines; version history is rarely worth the storage. Add `versioningEnabled: true` if you overwrite in place and want rollback.

## Placeholders to Replace

- `metadata.name` / `bucketName` -- your bucket's name. It must be globally unique within the region (across ALL DigitalOcean customers) and DNS-compatible.
- `corsRules[0].allowedOrigins` -- your website's exact origin (scheme + host); use `"*"` only for truly public assets.
- `region` -- a Spaces-capable region slug (ams3, atl1, blr1, fra1, lon1, nyc3, sfo2, sfo3, sgp1, syd1, tor1).
