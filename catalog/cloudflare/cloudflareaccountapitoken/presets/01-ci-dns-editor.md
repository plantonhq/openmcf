---
display_name: CI DNS Editor
---

# CI DNS editor

The pipeline-credential shape: DNS write across the account's zones, usable only from the CI egress range. Note the nested `subresources` form -- it scopes the grant to the account's ZONES rather than to the whole account, which is the narrower of the two shapes Cloudflare accepts. Fetch the permission-group UUID once with `GET /accounts/{account_id}/tokens/permission_groups?name=DNS%20Write` and pin it. Capture the `value` output into a managed secret on the first apply; it is never retrievable again.
