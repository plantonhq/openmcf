# Terraform Module: AWS OpenSearch Serverless Collection

Provisions an Amazon OpenSearch Serverless collection using Terraform,
together with the collection-scoped policies that make it usable.

## Resources Created

- `aws_opensearchserverless_security_policy` (type `encryption`) — Always
  created: AWS rejects CreateCollection without a matching encryption
  policy. AWS-owned key by default; the referenced customer-managed key
  otherwise. The collection depends on it (create after, destroy before).
- `aws_opensearchserverless_collection` — The collection: type (SEARCH /
  TIMESERIES / VECTORSEARCH), standby replicas, optional collection-group
  membership, optional vector acceleration.
- `aws_opensearchserverless_security_policy` (type `network`) — Always
  created: public reachability by default (SigV4-authenticated; not
  public data), or restricted to OpenSearch Serverless VPC endpoints.
- `aws_opensearchserverless_access_policy` (type `data`) — Only when
  `data_access` rules are declared; without rules nothing can read or
  write data (IAM permissions alone grant nothing in OpenSearch
  Serverless).
- `aws_opensearchserverless_lifecycle_policy` (type `retention`) — Only
  when `retention_rules` are declared.

Every policy's rules match exactly this collection
(`collection/<name>`, `index/<name>/<pattern>`) — one manifest owns one
collection and everything that makes it usable. All four policy objects
share the collection's name (types are separate namespaces at AWS).

## Usage

```hcl
module "collection" {
  source = "./path/to/module"

  metadata = {
    name = "app-search"
    org  = "my-org"
    env  = "production"
    id   = "awsoss-abc123"
  }

  spec = {
    region           = "us-west-2"
    type             = "SEARCH"
    standby_replicas = "DISABLED"

    data_access = [{
      principals        = ["arn:aws:iam::123456789012:role/app-role"]
      index_permissions = ["aoss:ReadDocument", "aoss:WriteDocument", "aoss:CreateIndex"]
      index_patterns    = ["*"]
    }]

    retention_rules = [{
      index_patterns      = ["logs-*"]
      min_index_retention = "30d"
    }]
  }
}
```

The collection's inline `encryption_config` argument is deliberately not
sent — the key choice lives on the module-rendered encryption security
policy (the universally supported arm at the pin); see the parity
manifest for the recorded judgment.
