---
title: "Private VPC egress — database-backed function"
description: "The composed private-networking pattern: the function routes egress through a [GcpServerlessVpcConnector](/docs/catalog/gcp/gcpserverlessvpcconnector) to reach a private-IP database (Cloud SQL..."
type: "preset"
rank: "03"
presetSlug: "03-private-vpc-egress"
componentSlug: "cloud-function"
componentTitle: "Cloud Function"
provider: "gcp"
icon: "package"
order: 3
---

# Private VPC egress — database-backed function

The composed private-networking pattern: the function routes egress through a [GcpServerlessVpcConnector](/docs/catalog/gcp/gcpserverlessvpcconnector) to reach a private-IP database (Cloud SQL private IP, Memorystore), runs as a dedicated least-privilege [GcpServiceAccount](/docs/catalog/gcp/gcpserviceaccount), and reads its credential from Secret Manager at instance start — the password never appears in the manifest.

`PRIVATE_RANGES_ONLY` keeps public egress on the normal path while RFC1918 traffic uses the connector. The runtime service account needs `roles/secretmanager.secretAccessor` on the secret ([GcpProjectIamMember](/docs/catalog/gcp/gcpprojectiammember) composes the grant).
