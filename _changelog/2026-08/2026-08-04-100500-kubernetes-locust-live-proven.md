# KubernetesLocust live-proven: the load-testing kind enters the green E2E matrix

## What changed

- **KubernetesLocust proven against a live cluster, both engines.** All
  six scenario-engine lanes green with verifier-output evidence on
  every lane: the web UI's auth gate bounced an anonymous stats read to
  the login (upstream ships the UI open — that posture never deploys
  from this kind), a wrong password was refused while the
  module-generated credential signed in through the platform-managed
  login backend, and a real distributed load test started through the
  master's own REST API drove requests through registered workers at a
  composed in-cluster target with zero failures. The full-surface lane
  additionally proved the pip-install-at-pod-start arm (a package
  fetched from PyPI inside the readiness window), the secret-native
  test environment (the lane's locustfile raises without its
  Secret-delivered token, so a zero-failure swarm is the delivery
  proof), tag filtering, PodDisruptionBudgets and a custom login
  username. The behavioral lane deleted the master pod, waited for a
  UID-verified replacement, and swarmed again with the SAME
  pre-replacement session cookie — the module-generated Flask
  session-signing key is stable by design, so operators stay signed in
  across pod restarts and workers re-register through the stable
  Service. Three blind Terraform import round-trips re-imported every
  module-owned resource — the Helm release, the script and
  login-backend ConfigMaps, the credential Secret, and both generated
  values (the login password and the signing key) imported by value —
  with the follow-up plan proposing no real change. Locust entered
  both Tier-2 engine tiers of the green E2E CI matrix.

- **The import-recipe ledger reconciled for the analytics and ML
  kinds.** The live-proven ledger in `pkg/iac/importmap/README.md`
  gained the rows for Airflow, the Spark/Ray/Flink five, JupyterHub
  and MLflow — each written from its lane records (blind round-trips
  already run and green) — plus the new Locust row, so the record of
  which import recipes are correctness-proven is complete for every
  Kubernetes kind that ships a map.

## Why

A load-testing tool is handed the credentials and the network position
to fire traffic at anything it can reach; "install succeeded" is not
evidence it is safe or that it can actually generate load. These lanes
prove the product loop the way a customer uses it — a locked-by-default
console with a real credential, a real distributed swarm against a real
target, coordination that survives a master crash without logging
anyone out — on both engines, with import recipes proven correct by
blind re-import rather than review.
