# Path Rewrite + Cookie Canary

This preset creates a delivery policy that decouples the public URL
surface from the backend's real path layout (a rewrite the client never
sees) and routes cookie-flagged requests to a canary origin group -- an
edge-level canary release that needs no route or DNS change.

## When to Use

- API versioning where the public path (`/v1/...`) differs from the
  backend's mount point (`/api/v1/...`)
- Canary or blue/green releases driven by a cookie your application (or
  test tooling) sets -- traffic moves per request, not per deployment

## Key Configuration Choices

- **`urlRewrite` rather than `urlRedirect`** -- the transformation is
  internal; the client keeps the public URL (the spec forbids both on
  one rule because they contradict)
- **`preserveUnmatchedPath: true`** -- `/v1/users` becomes
  `/api/v1/users`, not just `/api/v1`
- **The canary override pairs `originGroupId` with
  `forwardingProtocol`** -- the spec requires the pair (the protocol
  qualifies the overriding origin)
- **`cacheBehavior: DISABLED` on the canary** -- canary responses must
  never be cached for non-canary users

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<front-door-profile-resource-name>` | The AzureFrontDoorProfile's Planton resource name | Your Front Door composition |
| `<canary-origin-group-resource-name>` | The canary AzureFrontDoorOriginGroup's Planton resource name | Your canary backend composition |

## Downstream Wiring

Attach to the route serving the API endpoint via `ruleSetIds`. The
canary origin group is a normal `AzureFrontDoorOriginGroup` with its own
origins -- the override only redirects matched traffic to it.
