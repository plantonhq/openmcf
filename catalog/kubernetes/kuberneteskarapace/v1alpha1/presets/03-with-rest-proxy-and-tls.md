# REST proxy and TLS preset

The hardened single-replica shape: the registry API served over HTTPS
from a cert-manager-issued Secret, HTTP Basic authentication from a
hot-reloaded authfile Secret, and the Kafka REST-proxy role deployed
alongside for schema-aware produce/consume over HTTP.

Why `replicas: 1` is part of the preset, not an oversight: with
`server_tls`, follower replicas forward writes to the leader at its
advertised POD IP — an address no DNS-name certificate covers, so
multi-replica TLS serving breaks follower forwarding. The spec
carries this caveat on the field. When you need both TLS and
availability, run the production-ha preset's plain-HTTP replicas
behind TLS terminated at an Ingress/Gateway.

The composition seams:

- **`server_tls.secret_name`** is the cert-manager seam — a
  KubernetesCertificate's output Secret (`tls.crt`/`tls.key`, the
  defaults). The probes switch to HTTPS automatically and skip
  certificate verification, so Service-SAN certificates work.
- **`http_authentication.basic`** mounts a Karapace authfile (a JSON
  users/permissions document) from a Secret; rotating credentials is
  a Secret update — the engine hot-reloads, no restart. `/_health`
  stays open by engine design, so probes keep working.
- **The REST proxy** is its own Deployment (`<name>-rest`) with its
  own sizing and endpoint output; it resolves schemas through the
  registry over HTTPS (the scheme follows `server_tls`) but serves
  plain HTTP itself.

See
[03-with-rest-proxy-and-tls.yaml](./03-with-rest-proxy-and-tls.yaml)
for the manifest.
