# Inbound from On-Prem

The shape on-prem resolvers forward AWS-bound queries TO: two ENIs across AZs, guarded by a DNS security group. After deploy, point your on-prem conditional forwarders at the `ip_addresses` output.
