package module

// Stack output keys — exactly the DigitalOceanLoadBalancerStackOutputs
// contract, identical across both provisioners.
const (
	// OpLoadBalancerId is the exported stack output containing the balancer UUID.
	OpLoadBalancerId = "load_balancer_id"
	// OpIp is the exported stack output containing the public IPv4 address.
	OpIp = "ip"
	// OpUrn is the exported stack output with the balancer's uniform resource
	// name ("do:loadbalancer:<id>").
	OpUrn = "urn"
	// OpIpv6 is the exported stack output with the IPv6 address (populated
	// when network_stack is DUALSTACK).
	OpIpv6 = "ipv6"
)
