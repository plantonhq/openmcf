# AWS external loadbalancer preset

The production three-broker shape plus a `loadbalancer` listener for
clients OUTSIDE the Kubernetes cluster: the AWS Load Balancer
Controller provisions one NLB per broker plus one for bootstrap
(Kafka is raw TCP — an L7 balancer cannot carry it), and external-dns
annotations point a DNS name at each. Mutual TLS authenticates the
external clients; the internal SCRAM listener keeps in-cluster
traffic on the plain path.

For teams whose producers and consumers live outside the cluster —
other VPCs, on-premises, partner systems. The trade-offs: four NLBs
of standing cost for this shape (bootstrap + three brokers, and the
count grows with the broker pool), and DNS/certificate coupling —
Kafka clients bootstrap once and then connect to EVERY broker
directly, so each broker needs its own resolvable name and a matching
`advertised_host`. `external_traffic_policy: Local` preserves client
IPs at the cost of NLB health checks passing only on nodes running a
broker pod.

Prerequisites beyond the operator: the AWS Load Balancer Controller
and external-dns (KubernetesExternalDns) running in the cluster, and
a Route 53 zone for `kafka.example.com`.

See [03-aws-external-loadbalancer.yaml](./03-aws-external-loadbalancer.yaml)
for the manifest.
