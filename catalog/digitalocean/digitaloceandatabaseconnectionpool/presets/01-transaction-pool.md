# Transaction Pool for One Service

This preset creates the workhorse pool shape: transaction mode, a dedicated service user, pointed at the service's own logical database. Transaction mode hands a server connection out per transaction -- the right default for web applications and APIs.

## When to Use

- Web applications and APIs with many short transactions
- Services whose connection count would otherwise exhaust the cluster's limit
- Pairing with a DigitalOceanDatabaseDb (`orders`) and DigitalOceanDatabaseUser (`orders-service`) of the same service

## Key Configuration Choices

- **`mode: transaction`** -- maximum connection reuse; incompatible with session state (LISTEN/NOTIFY, prepared statements at session scope) -- use the session-mode preset for those.
- **`size: 20`** -- backend connections held open. Bounded by the cluster's connection limit (grows with node size); leave headroom for other pools and direct connections.
- **Dedicated `user`** -- the pool authenticates as the service's own user, keeping per-service credential rotation intact.

## What You Get

A pool endpoint (its own port beside the cluster's) whose `uri`/`private_uri` secret outputs wire straight into the service's configuration. Every field is create-only: plan changes as replacements.
