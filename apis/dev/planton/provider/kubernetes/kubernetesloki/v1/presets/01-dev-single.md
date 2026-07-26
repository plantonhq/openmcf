# Dev single-node Loki

The smallest honest Loki: one monolithic replica on a filesystem volume
with the nginx gateway. No object store, no tenancy — a dev loop or a
single-node cluster where logs are convenient but not yet durable at
scale.

The preset sizes the memcached caches explicitly because the chart's
defaults are production-scale: unset, the chunks cache alone requests
9830Mi of memory (8192MB allocated × the chart's 1.2 request factor),
which never schedules on a small node — and the atomic install then rolls
the whole release back (verified live). Cache memory is the one knob a
small cluster cannot leave to defaults.

**When to use:** local development, demos, a small cluster's own logs.

**When to move on:** for more than one replica, or any real retention
budget, switch to `02-production-scalable` (object storage is required the
moment you leave a single filesystem replica).
