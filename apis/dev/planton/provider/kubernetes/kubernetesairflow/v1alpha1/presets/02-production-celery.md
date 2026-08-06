# Production preset — Celery + git-sync + KEDA

The production shape: a Celery worker fleet that KEDA scales on real
queue depth (polling the metadata database for queued tasks — and back
to ZERO between runs), DAGs synced from your Git repository on every
component (a merged PR is a deployed pipeline), PgBouncer pooling the
many short-lived connections Airflow opens, two API servers and two
schedulers for availability, and task logs on a shared volume.

The database composes a KubernetesPostgres named `airflow-db` in the
same namespace (bootstrap database `airflow`, owner `airflow`); the
module composes every connection Secret from the referenced
credentials at deploy time — this manifest and the rendered chart
values never carry a password. The bundled Redis broker fits most
installs (upstream pins the last BSD-licensed 7.2 image line); a
composed KubernetesValkey is the licensing-clean shared-broker path —
swap the `broker` arm and nothing else changes.

Know the trade-offs: two schedulers coordinate through database row
locks (PostgreSQL handles this well — it is the reason postgres is the
recommended backend); `logging.persistence` needs a ReadWriteMany
storage class once components scale past one replica — on clusters
without one, point `logging.opensearch` at a composed
KubernetesOpenSearch instead and drop the volume.

Change first: worker sizing to your heaviest task (every Celery worker
runs up to 16 tasks concurrently by default), `keda.max_replicas` to
your cluster's headroom, and the admin bootstrap user's email.

See [02-production-celery.yaml](./02-production-celery.yaml) for the
manifest.
