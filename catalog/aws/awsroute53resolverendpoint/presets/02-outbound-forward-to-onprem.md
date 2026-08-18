# Outbound Forward to On-Prem

The classic hybrid split: corp.example.com forwards to two on-prem name servers, while aws.corp.example.com stays on AWS's recursive resolution (the SYSTEM override). The security group must allow DNS egress to the targets, and the network path (VPN/DX/peering) must exist.
