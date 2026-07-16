# General Purpose Regional EFS

Regional, encrypted, bursting throughput, backup enabled. Simplest production-safe starting point.

## When to Use

- EKS pods needing shared persistent storage via the EFS CSI driver
- EC2 instances or ECS tasks that mount EFS directly
- Workloads with predictable or moderate I/O patterns (bursting scales with storage size)
- First EFS deployment when you want minimal configuration

## What It Configures

- **Regional** — No `availabilityZoneName`; file system spans multiple AZs for high availability
- **Encrypted** — AES-256 encryption at rest using the AWS-managed key
- **Bursting throughput** — Throughput scales with storage; 50 MiB/s per TiB with bursts up to 100 MiB/s
- **Backup enabled** — Daily backups via AWS Backup
- **Two mount targets** — One per AZ so clients mount locally and avoid cross-AZ data charges

## What to Customize

- Replace placeholders: `<aws-region>`, `<subnet-id-az-a>`, `<subnet-id-az-b>`, `<security-group-id>`
- Add more mount targets (one per AZ) for broader availability
- Switch to `throughputMode: elastic` for unpredictable or spiky workloads
- Create `AwsEfsAccessPoint` resources referencing this file system for per-application POSIX isolation (ECS tasks, Lambda)
