# Lambda-to-PostgreSQL Pool

The classic deployment: Lambda functions storming a PostgreSQL instance get pooled through the proxy with IAM-token auth over enforced TLS — no passwords in function config, no connection exhaustion under burst. Wire your functions' database host to the proxy's `endpoint` output and grant their execution role `rds-db:connect`.
