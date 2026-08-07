# AtlasMongodb

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `atlas.planton.dev/v1alpha1`

atlas-mongodb spec

## Example

```yaml
apiVersion: atlas.planton.dev/v1alpha1
kind: AtlasMongodb
metadata:
  name: test-atlas-mongodb
  id: mdbatl-test-001
  org: engineering
  env: development
  labels:
    team: platform
    component: atlas-mongodb
    managed-by: planton
spec:
  clusterConfig:
    # Atlas MongoDB Project ID (required)
    # Create a project in Atlas UI and copy the Project ID from Project Settings
    projectId: "64f1a2b3c4d5e6f7g8h9i0j1"
    
    # Cluster type: REPLICASET (default), SHARDED, or GEOSHARDED
    clusterType: REPLICASET
    
    # Number of electable nodes (must be 3, 5, or 7 for consensus)
    # These nodes can become primary and facilitate local reads
    electableNodes: 3
    
    # Election priority (7 is highest, identifies preferred region)
    # In multi-region setups, primary region should be 7
    priority: 7
    
    # Number of read-only nodes (optional)
    # These nodes can never become primary but provide read scaling
    readOnlyNodes: 0
    
    # Enable cloud backups (recommended for production)
    cloudBackup: true
    
    # Enable automatic disk scaling
    autoScalingDiskGbEnabled: true
    
    # MongoDB major version
    # Supported versions: 4.4, 5.0, 6.0, 7.0
    mongoDbMajorVersion: "7.0"
    
    # Cloud provider: AWS, GCP, AZURE, or TENANT
    providerName: AWS
    
    # Instance size
    # Development: M10, M20
    # Production: M30+ (recommended)
    # Free tier: M0 (requires different resource type)
    providerInstanceSizeName: M10
```

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.clusterConfig` | `AtlasMongodbClusterConfig` |  |  |  |
| `spec.clusterConfig.projectId` | `string` |  |  |  |
| `spec.clusterConfig.clusterType` | `string` |  |  |  |
| `spec.clusterConfig.electableNodes` | `int32` |  |  |  |
| `spec.clusterConfig.priority` | `int32` |  |  |  |
| `spec.clusterConfig.readOnlyNodes` | `int32` |  |  |  |
| `spec.clusterConfig.cloudBackup` | `bool` |  |  |  |
| `spec.clusterConfig.autoScalingDiskGbEnabled` | `bool` |  |  |  |
| `spec.clusterConfig.mongoDbMajorVersion` | `string` |  |  |  |
| `spec.clusterConfig.providerName` | `string` |  |  |  |
| `spec.clusterConfig.providerInstanceSizeName` | `string` |  |  |  |

## Field Details

### spec.clusterConfig

`AtlasMongodbClusterConfig`

cluster-config

### spec.clusterConfig.projectId

`string`

The unique ID for the project to create the database user.
https://www.pulumi.com/registry/packages/atlasmongodb/api-docs/cluster/#projectid_yaml

### spec.clusterConfig.clusterType

`string`

Specifies the type of the cluster that you want to modify. You cannot convert a sharded cluster deployment to a replica set deployment.
Accepted values include:
REPLICASET Replica set
SHARDED Sharded cluster
GEOSHARDED Global Cluster
https://www.pulumi.com/registry/packages/atlasmongodb/api-docs/cluster/#clustertype_yaml

### spec.clusterConfig.electableNodes

`int32`

Number of electable nodes for Atlas to deploy to the region. Electable nodes can become the primary and can facilitate local reads.
The total number of electableNodes across all replication spec regions must total 3, 5, or 7.
Specify 0 if you do not want any electable nodes in the region.
You cannot create electable nodes in a region if priority is 0.
https://www.pulumi.com/registry/packages/atlasmongodb/api-docs/cluster/#electablenodes_yaml

### spec.clusterConfig.priority

`int32`

Election priority of the region. For regions with only read-only nodes, set this value to 0.
For regions where electable_nodes is at least 1, each region must have a priority of exactly one (1) less than the previous region. The first region must have a priority of 7. The lowest possible priority is 1.
The priority 7 region identifies the Preferred Region of the cluster. Atlas places the primary node in the Preferred Region. Priorities 1 through 7 are exclusive - no more than one region per cluster can be assigned a given priority.
Example: If you have three regions, their priorities would be 7, 6, and 5 respectively. If you added two more regions for supporting electable nodes, the priorities of those regions would be 4 and 3 respectively.
https://www.pulumi.com/registry/packages/atlasmongodb/api-docs/cluster/#priority_yaml

### spec.clusterConfig.readOnlyNodes

`int32`

Number of read-only nodes for Atlas to deploy to the region. Read-only nodes can never become the primary, but can facilitate local-reads. Specify 0 if you do not want any read-only nodes in the region.
https://www.pulumi.com/registry/packages/atlasmongodb/api-docs/cluster/#readonlynodes_yaml

### spec.clusterConfig.cloudBackup

`bool`

enable or disable cloud backup

### spec.clusterConfig.autoScalingDiskGbEnabled

`bool`

auto scaling disk db enabled

### spec.clusterConfig.mongoDbMajorVersion

`string`

Version of the cluster to deploy. Atlas supports the following MongoDB versions for M10+ clusters: 4.4, 5.0, 6.0 or 7.0.
If omitted, Atlas deploys a cluster that runs MongoDB 7.0.
If provider_instance_size_name: M0, M2 or M5, Atlas deploys MongoDB 5.0.
Atlas always deploys the cluster with the latest stable release of the specified version
https://www.pulumi.com/registry/packages/atlasmongodb/api-docs/cluster/#mongodbmajorversion_yaml

### spec.clusterConfig.providerName

`string`

Cloud service provider on which the servers are provisioned.

The possible values are:

AWS - Amazon AWS
GCP - Google Cloud Platform
AZURE - Microsoft Azure
TENANT - A multi-tenant deployment on one of the supported cloud service providers. Only valid when providerSettings.instanceSizeName is either M2 or M5.
https://www.pulumi.com/registry/packages/atlasmongodb/api-docs/cluster/#providername_yaml

### spec.clusterConfig.providerInstanceSizeName

`string`

https://www.pulumi.com/registry/packages/atlasmongodb/api-docs/cluster/#providerinstancesizename_yaml
Atlas provides different instance sizes, each with a default storage capacity and RAM size.
The instance size you select is used for all the data-bearing servers in your cluster.
https://www.pulumi.com/registry/packages/atlasmongodb/api-docs/cluster/#providerinstancesizename_yaml

## Outputs

Reference an output from another manifest as `valueFrom: {kind: AtlasMongodb, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.id` | `string` | The provider-assigned unique ID for the Atlas MongoDB cluster (cluster_id). https://registry.terraform.io/providers/mongodb/atlasmongodb/latest/docs/resources/advanced_cluster#cluster_id |
| `status.outputs.bootstrap_endpoint` | `string` | Atlas MongoDB standard connection string in SRV format (recommended for MongoDB drivers). Format: mongodb+srv://<username>:<password>@<cluster>.mongodb.net/<database> https://registry.terraform.io/providers/mongodb/atlasmongodb/latest/docs/resources/advanced_cluster#connection_strings |
| `status.outputs.crn` | `string` | The cluster identifier, same as id field. Used for resource identification and API operations. https://registry.terraform.io/providers/mongodb/atlasmongodb/latest/docs/resources/advanced_cluster#cluster_id |
| `status.outputs.rest_endpoint` | `string` | Atlas MongoDB standard connection string in legacy format. Format: mongodb://<host1>:<port1>,<host2>:<port2>/<database> https://registry.terraform.io/providers/mongodb/atlasmongodb/latest/docs/resources/advanced_cluster#connection_strings |

## See Also

- [Overview](./README.md)
- [Design notes](./docs/README.md)
