# Durable pipelines preset

Argo Workflows with its two durability seams filled: an S3-compatible
artifact store so steps pass files and archived logs outlive pods, and
a Postgres archive so run history outlives the Workflow CRs
themselves. With history safe in the database, the retention policy
keeps only a working set of CRs in the cluster — the UI still shows
everything, served from the archive.

The composition is the point: the artifact endpoint, the credentials
Secret name, and the database host are reference fields, so this
preset's literals become `value_from` references at a
KubernetesSeaweedFs and a KubernetesPostgres resource in a real chart
— the credential key-name defaults already match the SeaweedFS
generated `-s3-secret` (its admin pair), so the store's Secret
composes untouched, and none of those credentials ever ride this
manifest. A generic Secret keyed `accesskey`/`secretkey` (the argo
chart's documented example shape) works by setting the two key-name
fields explicitly.

Change first: on a cloud store, drop `insecure`, drop the declared
Secret and set `use_ambient_credentials` with IRSA/workload identity
annotated on the runner ServiceAccount — the keyless posture. And
create the `argo_archive` database before first boot: the controller
creates its tables, never the database.

See [02-durable-pipelines.yaml](./02-durable-pipelines.yaml) for the
manifest.
