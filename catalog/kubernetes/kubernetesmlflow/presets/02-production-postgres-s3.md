# Production — PostgreSQL backend, S3 artifacts

The durable, team-scale shape: experiments, runs, metrics and the model
registry in a composed KubernetesPostgres; models and datasets in a
composed KubernetesSeaweedFs bucket; two stateless server replicas
behind one Service; experiments private to their creators until shared;
and scheduled garbage collection so deleted runs eventually free their
storage.

The references do all the wiring: the database host resolves to the
Postgres read-write Service and its credential to the
operator-maintained application Secret; the artifact endpoint and keys
resolve to the SeaweedFS S3 endpoint and its credentials Secret. The
module composes the connection URI from them at deploy time — nothing
password-shaped is ever written into this manifest or any rendered
resource. Clients still need only the tracking URI and their MLflow
login: the server proxies all artifact traffic, so S3 credentials never
reach notebooks or pipelines.

What the preset expects from the composed side: the KubernetesPostgres
declares `mlflow` as its bootstrap database (owner `mlflow`), and the
KubernetesSeaweedFs declares the `mlflow-artifacts` bucket.

Change first: sizing under `server.resources` for your real query load
(experiment comparisons are memory-hungry), and `gc.older_than` if 30
days is the wrong undo window for your team.

See [02-production-postgres-s3.yaml](./02-production-postgres-s3.yaml)
for the manifest.
