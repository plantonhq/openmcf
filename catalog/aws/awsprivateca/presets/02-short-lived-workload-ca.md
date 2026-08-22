# Short-Lived Workload mTLS CA

The two-tier hierarchy's working tier: an EC subordinate signed by your root by reference (both deploy in one chart), issuing 7-day-max workload certificates in SHORT_LIVED_CERTIFICATE mode — USD 50/month instead of 400, because certificates that expire in days need no revocation infrastructure. Path length 0 keeps it a leaf-issuer; meshes and batch workloads rotate through it while the root stays cold.
