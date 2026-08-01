# Dev-mode preset

OpenBao with zero ceremony: dev mode auto-initializes and auto-unseals
at startup, the root token is literally `root`, and all data lives in
memory. Port-forward the `openbao-dev` Service and start writing
secrets immediately — no `bao operator init`, no unseal keys, no PVC.

This exists for exactly one purpose: evaluating OpenBao and building
against its API without the seal lifecycle. NEVER put real secrets in
it — everything is lost on every pod restart, and the root token sits
in plain text in the pod spec, readable by anyone who can get pods in
the namespace. Dev mode also drops ServiceAccount annotations (a chart
behavior), so cloud workload-identity seams do not apply here.

Know the lifecycle before you graduate: a REAL OpenBao server starts
uninitialized and SEALED, reports NotReady by design until you run
`bao operator init` and unseal it, and comes back SEALED after every
restart unless an auto-unseal backend is configured. The
production-ha preset carries that shape.

Change first: nothing. When you outgrow evaluation, switch presets
rather than hardening this one.

See [01-dev-mode.yaml](./01-dev-mode.yaml) for the manifest.
