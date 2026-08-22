# CloudflareKvNamespace guide

Operational judgment for Workers KV namespaces. The README covers what each field is; this covers how the pieces interact.

## Titles are account-unique identity

The namespace title doubles as its human identity: creating a second namespace with the same title in the account fails. Name namespaces for their purpose (`app-config-prod`), not their consumer — several Workers can bind one namespace.

## KV is eventually consistent, by design

A written value can take up to ~60 seconds to be visible from every edge location. KV suits read-heavy configuration and content; it is the wrong store for counters, locks, or anything needing read-after-write consistency (use Durable Objects or D1 for those).

## Deleting the namespace deletes every entry

Destroy is total and immediate — there is no recycle bin. Entries seeded through `CloudflareWorkersKvPair` are recorded in IaC and can be re-created; runtime-written data is simply gone. Treat a namespace holding runtime data as a stateful resource when planning teardown.

## Seed configuration here, write data at runtime

The clean split: infrastructure seeds slow-changing entries (feature flags, per-environment endpoints) as `CloudflareWorkersKvPair` resources, and the application writes high-churn data through its KV binding. IaC-managing hot data guarantees perpetual drift.
