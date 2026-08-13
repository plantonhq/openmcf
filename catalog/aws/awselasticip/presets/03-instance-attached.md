# Preset: Instance-Attached Elastic IP

**Use case:** Give an EC2 instance a public IP that survives stop/start cycles.

An instance's default public IP changes every time it stops; an Elastic IP
attached through this preset stays put. The attachment is declared on the
EIP itself (`spec.instance` — the same inline surface the AWS provider
offers), referencing the instance component, so the chart graph reads
"this address points at that instance".

## What You Get

- A VPC Elastic IP allocated from Amazon's pool and associated with the
  referenced instance's primary network interface
- Outputs: `allocation_id`, `public_ip`, `arn`, `public_dns`,
  `association_id` (proof of the attachment)

## Variations

- Attach to a specific network interface instead: use
  `spec.networkInterface` (at most one of the two targets may be set).
- Pin a specific private IP on a multi-IP target: add
  `spec.associateWithPrivateIp`.
- Add reverse DNS for mail workloads: set `spec.reverseDnsDomainName`
  AFTER pointing a forward A record at the address (AWS validates the
  forward record server-side before granting the PTR).

## When to Use

- A pet instance (bastion, license server, legacy appliance) needs a
  stable public address
- An external allowlist names one IP and the workload lives on EC2

## Cost

- **Free** while the instance runs; **$0.005/hour** (~$3.60/month) if the
  instance stops and the address sits idle
