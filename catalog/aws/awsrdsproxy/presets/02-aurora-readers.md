# Aurora Cluster with a Reader Endpoint

An Aurora MySQL cluster fronted end to end: writes through the default endpoint, reads through the `readers` READ_ONLY endpoint (the proxy distributes across the cluster's replicas and rides failovers gracefully). Pool tuning relaxes MySQL's variable-set pinning so multiplexing stays effective, and the idle ceiling keeps warm connections from starving the cluster.
