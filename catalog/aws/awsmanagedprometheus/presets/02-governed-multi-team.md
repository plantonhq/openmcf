# Governed Multi-Team Workspace

A shared workspace with guardrails: per-team active-series caps (the experiments team cannot melt the ingest budget), expensive-query logging above a QSP threshold, and an anomaly detector watching the ingest rate itself — missing data marks anomalous, because silence from the fleet IS the incident.
