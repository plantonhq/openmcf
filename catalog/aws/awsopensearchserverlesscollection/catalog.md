# AWS OpenSearch Serverless Collection

Deploys an Amazon OpenSearch Serverless collection -- a fully managed, auto-scaling OpenSearch workspace (capacity billed in OpenSearch Compute Units, no domains or nodes to size) for search, time-series, and vector workloads -- together with the collection-scoped encryption, network, data-access, and retention policies that make a collection usable. VECTORSEARCH collections are the vector store Bedrock knowledge bases consume. It integrates with Planton's Provider Connections for AWS credential management and supports ValueFromRef wiring to KMS keys and IAM roles.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **Encryption Security Policy** -- always created, scoped to exactly this collection (AWS rejects CreateCollection without a matching encryption policy): AWS-owned key by default, or the referenced customer-managed KMS key
- **Collection** -- SEARCH, TIMESERIES (default), or VECTORSEARCH, with standby replicas (default ENABLED; DISABLED halves the OCU floor for dev/test) and optional collection-group membership
- **Network Security Policy** -- public reachability of the collection and Dashboards endpoints by default (SigV4-authenticated reachability only -- not public data), or restricted to OpenSearch Serverless VPC endpoints
- **Data Access Policy** -- rendered from `dataAccess` rules granting IAM principals collection- and index-level permissions; without at least one rule nothing can read or write data (IAM permissions alone grant nothing in OpenSearch Serverless)
- **Lifecycle Policy** -- rendered from `retentionRules` (index retention by pattern; indefinite without rules)
- **AWS Tags** -- resource metadata tags (organization, environment, resource kind, resource ID) applied automatically on the collection

## Before You Deploy

### Planton Setup

- **AWS Provider Connection** -- an active connection in the Connect module with credentials for the target AWS account. Map it as the default for your environment, or specify it explicitly when creating the Cloud Resource.

### AWS Account

- **IAM principals for data access** -- the roles or users your applications run as; grant them through `dataAccess` rules. Reference AwsIamRole Cloud Resources via ValueFromRef or provide ARNs directly.
- **A KMS key** (optional) for customer-managed encryption -- the key choice is fixed at create time. Provide the ARN or reference an AwsKmsKey Cloud Resource.
- **OpenSearch Serverless VPC endpoints** (optional) for private network access -- these are the service's OWN endpoint objects (created through the OpenSearch Serverless API, not ordinary Interface Endpoints); create them outside this component and list their IDs in `network.vpcEndpointIds`.

## Deploy

### Console

Open the deployment store, find **AWS OpenSearch Serverless Collection**, and click **Deploy**. The creation wizard walks you through preset selection, environment and connection configuration, and spec fields. Start from the **Dev Search Collection** preset in the [Presets](#presets) tab to pre-populate a working configuration.

### CLI

Create a manifest and apply it:

```yaml
apiVersion: aws.planton.dev/v1alpha1
kind: AwsOpenSearchServerlessCollection
metadata:
  name: app-search
spec:
  region: us-west-2
  type: SEARCH
  standbyReplicas: DISABLED
  dataAccess:
    - principals:
        - value: arn:aws:iam::123456789012:role/app-role
      indexPermissions:
        - aoss:ReadDocument
        - aoss:WriteDocument
```

```bash
planton apply -f collection.yaml
```

## Operational Notes

- **ForceNew surface**: the name (metadata.name), `type`, `standbyReplicas`, `collectionGroupName`, and the encryption key choice are all fixed at create time. Network, data-access, and retention rules update in place.
- **The OCU floor bills while the collection exists** -- standby DISABLED halves it; destroy dev collections when idle.
- **Bedrock knowledge bases** consume VECTORSEARCH collections through the `collection_arn` output.
- **Account-wide pattern policies** (one policy matching many collections) are a different tool and deliberately outside this component -- it owns exactly one collection's policies.
