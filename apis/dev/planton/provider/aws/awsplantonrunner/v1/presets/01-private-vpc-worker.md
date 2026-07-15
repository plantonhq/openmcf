# Private VPC Worker

This preset deploys the standard runner appliance: a pull-based worker on
two private subnets that receives deploy operations through its queue and
executes them from inside the VPC. The 30-second decision for making a
private-endpoint Kubernetes cluster (or any in-network target)
deployable.

## When to Use

- A Kubernetes cluster with a private API endpoint needs to be deployed
  to and operated
- Any deployment target reachable only from inside the VPC (private
  databases, internal services)
- You want zero inbound network exposure: the runner only dials out

## Key Configuration Choices

- **`executionMode: temporal`** -- the pull-based worker mode; operations
  wait in the runner's queue until it polls, so nothing ever needs to
  reach the runner
- **Two private subnets in different AZs** -- the runner reschedules
  across an AZ event; private subnets need a NAT route for the runner's
  outbound traffic (image pulls, control plane)
- **Credentials as a managed-secret reference** -- the document from
  `planton runner generate-credentials`, stored once as a secret and
  resolved just-in-time; it never appears in plaintext anywhere
- **Default sizing (0.5 vCPU / 1 GiB)** -- comfortable for typical IaC
  operations; see the high-capacity preset when stacks grow large

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<runner-name>` | Name for the runner appliance | Match the runner registration's name |
| `<aws-region>` | AWS region code | The region hosting the private targets |
| `<private-subnet-a/b-resource-name>` | Names of the private AwsSubnet resources | Your subnet manifests' `metadata.name` |
| `<runner-credentials-secret-slug>` | The managed secret holding the credentials JSON | Created from `planton runner generate-credentials <runner-name>` |

## Related Presets

- `02-dual-mode` -- adds the real-time CloudOps channel (live resource browsing)
- `03-high-capacity` -- pinned version and larger sizing for heavy IaC workloads
