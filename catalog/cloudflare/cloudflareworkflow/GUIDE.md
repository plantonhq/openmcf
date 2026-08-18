# CloudflareWorkflow guide

The judgment this guide protects you from: the workflow NAME is the identity, create is an upsert on that name, and the registration is only the binding -- the script, the class, and the running instances each live their own lives.

## Create is a PUT: name collisions adopt, never fail

Cloudflare has no separate create endpoint for workflows -- registration writes to `PUT .../workflows/<name>`. Registering a name that already exists silently adopts and OVERWRITES the existing workflow's binding (class, script, retention, schedules). Two teams using the same workflow name in one account will fight over one registration and neither apply will ever error. Namespace workflow names deliberately (team or system prefixes), exactly as you would Worker script names.

## The registration is not the program

This resource binds a name to a class exported by an ALREADY-DEPLOYED Worker script. Deploy order is script first, workflow second -- the `script_name` reference on a `CloudflareWorker` gives you that ordering for free. The class must subclass `WorkflowEntrypoint` and be exported by name; a typo in `class_name` may register fine and fail only when the first instance runs (verify against your script's exports, not the apply result).

## Retention takes two forms -- pick one and stay with it

`error_retention` / `success_retention` accept integer milliseconds ("86400000") or duration expressions ("1 day"). Cloudflare normalizes internally and the provider suppresses equivalent-value diffs, but mixed forms across manifests make review harder than it needs to be. Choose the duration-expression form for anything a human reads.

## Deletion is real but the API never says 404

Deleted workflows keep answering GET with a non-zero `is_deleted` marker. Tooling that probes for existence (including our own E2E verifier) reads the marker, not the status code. Do not be surprised by dashboard-API responses for workflows you deleted -- the record is a tombstone, not a survivor.

## Instances outlive your apply

Destroying the registration does not await or terminate running instances -- they belong to the platform. If a workflow must drain before decommissioning, stop its triggers first (remove `schedules`, stop creating instances from your Worker), let instances finish, then destroy.

## Pairs well with

- [CloudflareWorker](../cloudflareworker/README.md) -- the script that exports the workflow class; its `workflows` bindings reference this registration by name.
- [CloudflareQueue](../cloudflarequeue/README.md) -- queues feed Workers that create workflow instances.
