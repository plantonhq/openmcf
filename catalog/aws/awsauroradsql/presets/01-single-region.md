# Single-Region Production Cluster

The whole production posture in two lines: deletion protection on, AWS-owned encryption, nothing to size. Connect with any PostgreSQL driver at the `endpoint` output using IAM auth tokens — there is no password to manage, and an idle cluster bills ~nothing.
