---
title: "Airflow"
description: "Airflow deployment documentation"
icon: "package"
order: 100
componentName: "kubernetesairflow"
---

# Airflow

The data-pipeline orchestrator. ETL, ML training schedules, analytics
jobs and their dependencies — written as Python DAGs, scheduled and
retried by the platform the data world standardized on, with a UI that
shows every run, task and log.

## Highlights

- **Airflow 3, the real component set** — API server, scheduler,
  standalone DAG processor, triggerer, and a Celery worker fleet when
  you want one; migrations and the admin user handled at install.
- **Bring your own database, composed by reference** — a
  `KubernetesPostgres` pairs naturally (host and credential resolve by
  reference); external MySQL 8 works the same way. The chart's
  non-production bundled database never ships.
- **Secret-native end to end** — connection URIs, the Fernet key,
  session and JWT keys, the admin password: module-generated or
  composed into Secrets, never rendered into values or pod arguments
  (the chart's own defaults would regenerate keys on every upgrade —
  these are stable).
- **Push-to-deploy DAGs** — git-sync on every component; a merged PR
  is a deployed pipeline. HTTPS-token and SSH-key arms for private
  repositories.
- **Scale on real queue depth** — KEDA drives the Celery fleet from
  queued-task counts, down to zero between runs; PgBouncer pools the
  connection storm Airflow is famous for.

---

© Planton. Licensed under [Apache-2.0](https://github.com/plantonhq/planton/blob/main/LICENSE).
