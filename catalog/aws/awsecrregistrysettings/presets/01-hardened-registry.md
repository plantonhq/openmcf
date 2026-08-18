# Hardened Registry Posture

The security baseline: Inspector-backed enhanced scanning re-scans production images as new CVEs publish (everything else scans on push), account settings pinned to current-generation values, and the CI role excluded from pull-time metrics so lifecycle policies expire by REAL usage. Inspector bills per scanned image — the prod-only continuous rule keeps that spend proportional.
