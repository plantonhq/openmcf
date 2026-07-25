# Dev single-node Loki

The smallest honest Loki: one monolithic replica on a filesystem volume
with the nginx gateway. No object store, no tenancy — a dev loop or a
single-node cluster where logs are convenient but not yet durable at
scale.

**When to use:** local development, demos, a small cluster's own logs.

**When to move on:** for more than one replica, or any real retention
budget, switch to `02-production-scalable` (object storage is required the
moment you leave a single filesystem replica).
