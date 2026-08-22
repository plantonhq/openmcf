package module

const (
	// OpPoolId is the autoscale pool's UUID (its API identity and import
	// id).
	OpPoolId = "pool_id"
	// OpStatus is the pool's health status at apply time ("active" once
	// the pool and every member droplet are provisioned).
	OpStatus = "status"
)
