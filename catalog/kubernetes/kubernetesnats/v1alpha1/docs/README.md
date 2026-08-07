# Kubernetes NATS — design notes

## Grain

One resource = one NATS system from the official `nats` chart (the
served index lives under nats-io.github.io/k8s/helm/charts/; chart and
server move in version lockstep). The release is named after
`metadata.name` and `fullnameOverride` pins every child name
(`<name>`, `<name>-headless`, `<name>-box`), so the exported outputs
are deterministic. Names are capped at 50 characters — the longest
derived child name adds 13 characters and Helm truncates at 63; both
engines fail loudly instead of letting the chart truncate silently.

The server is the grain boundary: streams, consumers and KV buckets
are data-plane objects created by clients, deliberately NOT modeled as
spec fields — declaring them here would couple many independent
lifecycles to one resource and turn every stream change into an infra
deploy.

## The config seam

The chart deliberately ships no typed auth surface — its architecture
is per-listener blocks (`cluster`/`jetstream`/`nats`/`websocket`/
`mqtt`/`leafnodes`) with `merge`/`patch` escapes into raw nats.conf.
The typed spec renders INTO that seam: the auth model lands under
`config.merge` (authorization users XOR accounts — CEL-enforced,
because the server rejects a config mixing both), listener toggles
land in their blocks, and the plain-in-cluster websocket arm carries
`no_tls: true` through its merge — the server refuses a websocket
listener without TLS unless no_tls is explicit.

## The credential design

Chart config values render into a ConfigMap, so a naive rendering of
declared users would put passwords on the API server in plaintext.
Instead: the module generates one password per username into the
`<name>-auth` Secret (one key per username — the paired-credential
contract clients read), wires each as a secretKeyRef-backed env var on
the server container, and the rendered config carries `$NATS_PW_<i>`
references (unquoted via the chart's `<< >>` syntax) the server
resolves from environment at load. Env-var names are index-based
(deterministic and env-safe for any username); password RESOURCES are
keyed by username on both engines, so reordering users in the spec
never rotates or swaps anyone's credential and renaming a user is
honestly a new credential. Generation-shape arguments are ignored
after creation — an imported credential never silently regenerates;
rotation stays an explicit verb.

Generated passwords are LETTERS ONLY (length compensates the smaller
alphabet). This is a server contract, not a style choice: nats-server
re-parses a resolved env reference through its own config parser, so a
password containing digits at the front or structural characters
(`-`, `#`, `$`, braces, quotes) crashes every server at config load
with a variable-reference parse error (verified live). A pure-letter
value is the only shape the parser can never misread.

## JetStream posture

ON by default (the chart's raw default is off) and rendered explicitly
either way, so both engines state the posture: a file-store PVC
(size/class from the spec), an optional capped memory tier, and
per-account JetStream enablement when accounts are declared.

JetStream is driven over core NATS subjects: clients PUBLISH requests
to `$JS.API.>` and acknowledgements to `$JS.ACK.>`, and receive
responses on `_INBOX.>` subscriptions. A user with a publish
ALLOWLIST is therefore fenced from JetStream unless the allowlist
includes those subjects — and because the server silently drops denied
publishes, a fenced user's stream operations time out client-side
instead of failing loudly (verified live). Grant `$JS.API.>` and
`$JS.ACK.>` to every allowlisted user that works with streams;
`publish_deny` rules can then carve out narrower fences where needed
(deny wins on overlap).

## Cross-engine parity

Both engines render byte-identical chart values from one shared
resolution order: typed values first, the `helm_values` escape hatch
merged over them with Helm `-f` semantics, and `fullnameOverride`
re-pinned LAST — the one deliberate exception to the escape hatch's
last-word contract, because every exported output derives from the
fullname. Service type/annotations and scheduling ride the chart's
merge-patch seams (`service.merge`, `podTemplate.merge`) — the chart
has no first-class values for them.

## Deliberate exclusions

Supercluster gateways, the JWT/operator auth mode and its resolver,
per-listener TLS beyond the client listener, and raw nats.conf keys —
reachable through `helm_values` (the chart's `config.merge` carries
arbitrary server config), never the primary interface. A declarative
stream-as-resource surface (the NACK controller's Stream/Consumer/
KeyValue CRDs) is a separate operator concern this kind deliberately
does not bundle.

---

© Planton. Licensed under [Apache-2.0](https://github.com/plantonhq/planton/blob/main/LICENSE).
