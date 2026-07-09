# EKS Environment InfraChart

This chart provisions a **complete, production-ready Kubernetes environment on AWS** in a single run — the cluster, the network it lives in, and the in-cluster add-ons that make it useful:

* Custom VPC across two AZs, composed from standalone networking primitives — Internet Gateway, one public and one private subnet per AZ (each with its own route table), Elastic IP(s) and NAT gateway(s)
* Selectable private-subnet egress via `nat_mode`: `single` (one shared NAT, cost-conscious), `per_az` (one NAT per AZ, highly available), or `none` (no outbound internet; nodes cannot pull images)
* IAM roles for control-plane and nodes
* Optional customer-managed KMS key for secrets encryption
* Private (default) or restricted-public API endpoint with CloudWatch control-plane logs
* Managed node group with autoscaling, Spot or On-Demand instances
* AWS-managed core add-ons (vpc-cni, kube-proxy, CoreDNS) as first-class nodes, adopted from the cluster's bootstrap copies
* The networking add-on trio: **ingress-nginx** (traffic entry), **cert-manager** (TLS certificates), and — when `dnsEnabled` is on — a **Route 53 public zone with external-dns** wired to it, so Ingress and Service hostnames get DNS records automatically (delegate the zone's NS records at your registrar)

## How the add-ons reach the cluster

The in-cluster add-ons carry a `planton.dev/connection` annotation naming the cluster's Kubernetes connection, which the platform publishes at the env-qualified slug `<env>-<cluster name>` when the cluster deploys. No manual connection setup is needed: the connection exists before the first add-on starts.

## Private API endpoints

With `disable_public_endpoint` on (the default), nothing outside the VPC can reach the Kubernetes API — including the platform's hosted runners. The chart therefore includes a standing **Planton runner** on serverless compute inside the VPC (`enable_planton_runner`, default on): an outbound-only worker that deploys and operates the in-cluster add-ons with zero inbound exposure. A dedicated security group admits the runner (and only the runner) to the cluster's API endpoint on port 443. The cluster resource names its runner via the `planton.dev/runner` annotation, and in-cluster work is routed through it automatically.

The runner needs outbound internet to pull its image and dial the control plane, so keep `nat_mode` at `single` or `per_az` when the endpoint is private.

Edit **values.yaml** to tailor the deployment; each boolean toggle cleanly removes its piece.

© Planton. Licensed under [Apache-2.0](../../../LICENSE).
