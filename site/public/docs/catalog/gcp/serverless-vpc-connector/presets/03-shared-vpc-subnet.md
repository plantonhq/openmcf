---
title: "Shared VPC — subnet placement"
description: "The required shape on Shared VPC: the connector cannot carve its own range out of a host-project network, so a network admin creates a dedicated `/28` subnetwork in the host project..."
type: "preset"
rank: "03"
presetSlug: "03-shared-vpc-subnet"
componentSlug: "serverless-vpc-connector"
componentTitle: "Serverless VPC Connector"
provider: "gcp"
icon: "package"
order: 3
---

# Shared VPC — subnet placement

The required shape on Shared VPC: the connector cannot carve its own range out of a host-project network, so a network admin creates a dedicated `/28` subnetwork in the host project ([GcpSubnetwork](/docs/catalog/gcp/gcpsubnetwork)) and the connector occupies it. The subnet must be exactly `/28` and used by nothing else.

`subnet.projectId` names the host project; omit it when the subnet lives in the connector's own project (useful when network admins simply prefer the range managed like any other subnet).
