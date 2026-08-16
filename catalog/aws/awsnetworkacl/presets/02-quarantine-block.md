# Data Tier with a Quarantine Deny

What security groups cannot do: an explicit DENY for a known-bad range, numbered BELOW the allows so it always wins, in front of a Postgres-only allow from the VPC. Egress permits only ephemeral-port replies back into the VPC — the data tier initiates nothing.
