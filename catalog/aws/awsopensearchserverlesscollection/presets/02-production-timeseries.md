# Production Time-Series (Log Analytics)

This preset creates a production TIMESERIES collection with standby replicas, a customer-managed KMS key, least-privilege data access split between an ingest role and an analyst role, and retention rules that expire logs while keeping audit indexes forever.

## When to Use

- Centralized application/service log analytics
- Observability pipelines (Firehose, ingestion Lambdas, OTel collectors writing OpenSearch)
- Any production time-series workload with retention requirements

## Key Configuration Choices

- **`standbyReplicas: ENABLED`** — warm standbys in a second AZ (the AWS production default; 2+2 OCU floor)
- **Customer-managed KMS** — fixed at create; gives rotation-independent audit and revocation
- **Split data access** — the ingest role writes and manages indexes, analysts read; neither can administer the collection
- **Retention rules** — `logs-*` expire after 30 days; `audit-*` carry an explicit unlimited rule so a broader future rule can never sweep them

## Operational Notes

Retention and data-access rules update in place. The name, type, standby posture, and key are fixed at create time.
