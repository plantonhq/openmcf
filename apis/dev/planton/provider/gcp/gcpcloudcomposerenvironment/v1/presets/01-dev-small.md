# Dev Small

A minimal Cloud Composer environment for developing and testing Airflow
DAGs. Small allocations, public endpoint, no private networking — the
cheapest footprint that still runs real pipelines.

## When to Use

- Developing and testing Airflow DAGs before promoting them
- Learning Cloud Composer without production overhead
- Small, low-volume pipelines where cost matters more than isolation

## Key Configuration

- **ENVIRONMENT_SIZE_SMALL** — minimal managed-infrastructure footprint
- **Small workloads** — 0.5 CPU / 2 GB for scheduler, web server, and
  workers; 1-3 workers autoscale with queue depth
- **Public endpoint** — the Airflow UI is reachable without VPN; fine for
  dev, not for production
- **No CMEK, no private networking, no recovery snapshots** — the
  production presets add these

## What to Customize

- `region` — where the environment lives; immutable after creation
- `softwareConfig.imageVersion` — pin a current Composer/Airflow image
- `projectId` — add it if you deploy outside the provider's default
  project

## Important Notes

- Environment creation takes 25-45 minutes: Composer assembles a GKE
  cluster, Cloud SQL database, and web server behind the scenes.
- Not suitable for production — the UI is public and there is no
  disaster recovery.

## Related Presets

- **02-production-private** — private networking, high resilience, and
  control-plane allowlisting
- **03-enterprise-encrypted** — CMEK, data retention, and disaster
  recovery on top
