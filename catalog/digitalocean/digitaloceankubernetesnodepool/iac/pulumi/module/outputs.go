package module

const (
	// OpNodePoolId exports the UUID of the created node pool.
	OpNodePoolId = "node_pool_id"
	// OpNodeIds exports the DOKS node object UUIDs of the pool's members.
	OpNodeIds = "node_ids"
	// OpClusterId exports the UUID of the cluster that owns this pool.
	OpClusterId = "cluster_id"
	// OpDropletIds exports the ids of the Droplets backing the pool's nodes.
	OpDropletIds = "droplet_ids"
)
