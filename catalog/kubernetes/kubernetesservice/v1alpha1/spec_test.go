package kubernetesservicev1alpha1

import (
	"testing"

	"buf.build/go/protovalidate"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	foreignkeyv1 "github.com/plantonhq/planton/shared/foreignkey/v1"
)

func TestKubernetesServiceSpec(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "KubernetesServiceSpec Validation Suite")
}

// enum shorthands keep the table-style specs readable.
func svcType(t KubernetesServiceSpec_KubernetesServiceType) *KubernetesServiceSpec_KubernetesServiceType {
	return &t
}

func affinity(a KubernetesServiceSpec_KubernetesServiceSessionAffinity) *KubernetesServiceSpec_KubernetesServiceSessionAffinity {
	return &a
}

func etp(p KubernetesServiceSpec_KubernetesServiceExternalTrafficPolicy) *KubernetesServiceSpec_KubernetesServiceExternalTrafficPolicy {
	return &p
}

func famPolicy(p KubernetesServiceSpec_KubernetesServiceIpFamilyPolicy) *KubernetesServiceSpec_KubernetesServiceIpFamilyPolicy {
	return &p
}

func i32(v int32) *int32 { return &v }

func boolPtr(v bool) *bool { return &v }

var _ = ginkgo.Describe("KubernetesServiceSpec validations", func() {

	ginkgo.Context("When valid specs are provided", func() {

		ginkgo.It("accepts a minimal spec: name plus one port", func() {
			spec := &KubernetesServiceSpec{
				Name:  "web",
				Ports: []*KubernetesServicePort{{Port: 80}},
			}
			gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
		})

		ginkgo.It("accepts a namespace provided as a literal value", func() {
			spec := &KubernetesServiceSpec{
				Name: "web",
				Namespace: &foreignkeyv1.StringValueOrRef{
					LiteralOrRef: &foreignkeyv1.StringValueOrRef_Value{Value: "prod"},
				},
				Ports: []*KubernetesServicePort{{Port: 80}},
			}
			gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
		})

		ginkgo.It("accepts a namespace provided as a resource reference", func() {
			spec := &KubernetesServiceSpec{
				Name: "web",
				Namespace: &foreignkeyv1.StringValueOrRef{
					LiteralOrRef: &foreignkeyv1.StringValueOrRef_ValueFrom{
						ValueFrom: &foreignkeyv1.ValueFromRef{Name: "team-namespace"},
					},
				},
				Ports: []*KubernetesServicePort{{Port: 80}},
			}
			gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
		})

		ginkgo.It("accepts a full LoadBalancer spec with all LB knobs", func() {
			spec := &KubernetesServiceSpec{
				Name:                          "public-lb",
				Type:                          svcType(KubernetesServiceSpec_load_balancer),
				Selector:                      map[string]string{"app": "web"},
				Ports:                         []*KubernetesServicePort{{Name: "http", Port: 80, NodePort: 30080}},
				ExternalTrafficPolicy:         etp(KubernetesServiceSpec_local),
				HealthCheckNodePort:           30081,
				LoadBalancerSourceRanges:      []string{"203.0.113.0/24", "2001:db8::/64"},
				LoadBalancerClass:             "example.com/internal-vip",
				AllocateLoadBalancerNodePorts: boolPtr(false),
			}
			gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
		})

		ginkgo.It("accepts a headless service", func() {
			spec := &KubernetesServiceSpec{
				Name:                     "peers",
				Headless:                 true,
				Selector:                 map[string]string{"app": "db"},
				Ports:                    []*KubernetesServicePort{{Port: 5432}},
				PublishNotReadyAddresses: true,
			}
			gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
		})

		ginkgo.It("accepts an ExternalName service with no ports", func() {
			spec := &KubernetesServiceSpec{
				Name:            "ext-db",
				Type:            svcType(KubernetesServiceSpec_external_name),
				ExternalDnsName: "db.prod.example.com",
			}
			gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
		})

		ginkgo.It("accepts a static cluster IP on a non-headless service", func() {
			spec := &KubernetesServiceSpec{
				Name:             "pinned-vip",
				ClusterIpAddress: "10.96.0.50",
				Ports:            []*KubernetesServicePort{{Port: 80}},
			}
			gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
		})

		ginkgo.It("accepts ClientIP session affinity with a timeout", func() {
			spec := &KubernetesServiceSpec{
				Name:                          "sticky",
				SessionAffinity:               affinity(KubernetesServiceSpec_client_ip),
				SessionAffinityTimeoutSeconds: i32(3600),
				Ports:                         []*KubernetesServicePort{{Port: 80}},
			}
			gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
		})

		ginkgo.It("accepts dual-stack with two distinct families", func() {
			spec := &KubernetesServiceSpec{
				Name:           "dual",
				IpFamilies:     []KubernetesServiceSpec_KubernetesServiceIpFamily{KubernetesServiceSpec_ipv6, KubernetesServiceSpec_ipv4},
				IpFamilyPolicy: famPolicy(KubernetesServiceSpec_require_dual_stack),
				Ports:          []*KubernetesServicePort{{Port: 80}},
			}
			gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
		})

		ginkgo.It("accepts named and numeric target ports", func() {
			spec := &KubernetesServiceSpec{
				Name: "multi",
				Ports: []*KubernetesServicePort{
					{Name: "http", Port: 80, TargetPort: "8080"},
					{Name: "metrics", Port: 9090, TargetPort: "metrics"},
				},
			}
			gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
		})

		ginkgo.It("accepts external IPs", func() {
			spec := &KubernetesServiceSpec{
				Name:        "ext-vip",
				ExternalIps: []string{"198.51.100.7"},
				Ports:       []*KubernetesServicePort{{Port: 80}},
			}
			gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
		})
	})

	ginkgo.Context("When invalid specs are provided", func() {

		ginkgo.It("rejects an empty name", func() {
			spec := &KubernetesServiceSpec{
				Ports: []*KubernetesServicePort{{Port: 80}},
			}
			gomega.Expect(protovalidate.Validate(spec)).ToNot(gomega.BeNil())
		})

		ginkgo.It("rejects a name starting with a digit (DNS-1035, verified against the live API)", func() {
			spec := &KubernetesServiceSpec{
				Name:  "0web",
				Ports: []*KubernetesServicePort{{Port: 80}},
			}
			gomega.Expect(protovalidate.Validate(spec)).ToNot(gomega.BeNil())
		})

		ginkgo.It("rejects an uppercase name", func() {
			spec := &KubernetesServiceSpec{
				Name:  "Web",
				Ports: []*KubernetesServicePort{{Port: 80}},
			}
			gomega.Expect(protovalidate.Validate(spec)).ToNot(gomega.BeNil())
		})

		ginkgo.It("rejects a proxying service with zero ports", func() {
			spec := &KubernetesServiceSpec{Name: "web"}
			gomega.Expect(protovalidate.Validate(spec)).ToNot(gomega.BeNil())
		})

		ginkgo.It("rejects ExternalName without external_dns_name", func() {
			spec := &KubernetesServiceSpec{
				Name: "ext",
				Type: svcType(KubernetesServiceSpec_external_name),
			}
			gomega.Expect(protovalidate.Validate(spec)).ToNot(gomega.BeNil())
		})

		ginkgo.It("rejects external_dns_name on a non-ExternalName type", func() {
			spec := &KubernetesServiceSpec{
				Name:            "web",
				ExternalDnsName: "db.example.com",
				Ports:           []*KubernetesServicePort{{Port: 80}},
			}
			gomega.Expect(protovalidate.Validate(spec)).ToNot(gomega.BeNil())
		})

		ginkgo.It("rejects external IPs on an ExternalName service", func() {
			spec := &KubernetesServiceSpec{
				Name:            "ext",
				Type:            svcType(KubernetesServiceSpec_external_name),
				ExternalDnsName: "db.example.com",
				ExternalIps:     []string{"198.51.100.7"},
			}
			gomega.Expect(protovalidate.Validate(spec)).ToNot(gomega.BeNil())
		})

		ginkgo.It("rejects a selector on an ExternalName service", func() {
			spec := &KubernetesServiceSpec{
				Name:            "ext",
				Type:            svcType(KubernetesServiceSpec_external_name),
				ExternalDnsName: "db.example.com",
				Selector:        map[string]string{"app": "db"},
			}
			gomega.Expect(protovalidate.Validate(spec)).ToNot(gomega.BeNil())
		})

		ginkgo.It("rejects headless combined with NodePort", func() {
			spec := &KubernetesServiceSpec{
				Name:     "bad",
				Type:     svcType(KubernetesServiceSpec_node_port),
				Headless: true,
				Ports:    []*KubernetesServicePort{{Port: 80}},
			}
			gomega.Expect(protovalidate.Validate(spec)).ToNot(gomega.BeNil())
		})

		ginkgo.It("rejects a static cluster IP on a headless service", func() {
			spec := &KubernetesServiceSpec{
				Name:             "bad",
				Headless:         true,
				ClusterIpAddress: "10.96.0.50",
				Ports:            []*KubernetesServicePort{{Port: 80}},
			}
			gomega.Expect(protovalidate.Validate(spec)).ToNot(gomega.BeNil())
		})

		ginkgo.It("rejects the literal \"None\" as a cluster IP", func() {
			spec := &KubernetesServiceSpec{
				Name:             "bad",
				ClusterIpAddress: "None",
				Ports:            []*KubernetesServicePort{{Port: 80}},
			}
			gomega.Expect(protovalidate.Validate(spec)).ToNot(gomega.BeNil())
		})

		ginkgo.It("rejects an invalid external IP", func() {
			spec := &KubernetesServiceSpec{
				Name:        "bad",
				ExternalIps: []string{"not-an-ip"},
				Ports:       []*KubernetesServicePort{{Port: 80}},
			}
			gomega.Expect(protovalidate.Validate(spec)).ToNot(gomega.BeNil())
		})

		ginkgo.It("rejects an invalid load balancer source range", func() {
			spec := &KubernetesServiceSpec{
				Name:                     "bad",
				Type:                     svcType(KubernetesServiceSpec_load_balancer),
				LoadBalancerSourceRanges: []string{"203.0.113.5"},
				Ports:                    []*KubernetesServicePort{{Port: 80}},
			}
			gomega.Expect(protovalidate.Validate(spec)).ToNot(gomega.BeNil())
		})

		ginkgo.It("rejects LoadBalancer-only knobs on a ClusterIP service", func() {
			spec := &KubernetesServiceSpec{
				Name:              "bad",
				LoadBalancerClass: "example.com/vip",
				Ports:             []*KubernetesServicePort{{Port: 80}},
			}
			gomega.Expect(protovalidate.Validate(spec)).ToNot(gomega.BeNil())
		})

		ginkgo.It("rejects health_check_node_port without Local-policy LoadBalancer", func() {
			spec := &KubernetesServiceSpec{
				Name:                "bad",
				Type:                svcType(KubernetesServiceSpec_load_balancer),
				HealthCheckNodePort: 30081,
				Ports:               []*KubernetesServicePort{{Port: 80}},
			}
			gomega.Expect(protovalidate.Validate(spec)).ToNot(gomega.BeNil())
		})

		ginkgo.It("rejects a session affinity timeout without ClientIP affinity", func() {
			spec := &KubernetesServiceSpec{
				Name:                          "bad",
				SessionAffinityTimeoutSeconds: i32(3600),
				Ports:                         []*KubernetesServicePort{{Port: 80}},
			}
			gomega.Expect(protovalidate.Validate(spec)).ToNot(gomega.BeNil())
		})

		ginkgo.It("rejects a session affinity timeout above one day", func() {
			spec := &KubernetesServiceSpec{
				Name:                          "bad",
				SessionAffinity:               affinity(KubernetesServiceSpec_client_ip),
				SessionAffinityTimeoutSeconds: i32(90000),
				Ports:                         []*KubernetesServicePort{{Port: 80}},
			}
			gomega.Expect(protovalidate.Validate(spec)).ToNot(gomega.BeNil())
		})

		ginkgo.It("rejects duplicate IP families", func() {
			spec := &KubernetesServiceSpec{
				Name:       "bad",
				IpFamilies: []KubernetesServiceSpec_KubernetesServiceIpFamily{KubernetesServiceSpec_ipv4, KubernetesServiceSpec_ipv4},
				Ports:      []*KubernetesServicePort{{Port: 80}},
			}
			gomega.Expect(protovalidate.Validate(spec)).ToNot(gomega.BeNil())
		})

		ginkgo.It("rejects two families under a single_stack policy", func() {
			spec := &KubernetesServiceSpec{
				Name:           "bad",
				IpFamilies:     []KubernetesServiceSpec_KubernetesServiceIpFamily{KubernetesServiceSpec_ipv4, KubernetesServiceSpec_ipv6},
				IpFamilyPolicy: famPolicy(KubernetesServiceSpec_single_stack),
				Ports:          []*KubernetesServicePort{{Port: 80}},
			}
			gomega.Expect(protovalidate.Validate(spec)).ToNot(gomega.BeNil())
		})
	})

	ginkgo.Context("Port-level validations", func() {

		ginkgo.It("rejects a port number out of range", func() {
			spec := &KubernetesServiceSpec{
				Name:  "bad",
				Ports: []*KubernetesServicePort{{Port: 70000}},
			}
			gomega.Expect(protovalidate.Validate(spec)).ToNot(gomega.BeNil())
		})

		ginkgo.It("rejects a node port outside the node-port range", func() {
			spec := &KubernetesServiceSpec{
				Name:  "bad",
				Type:  svcType(KubernetesServiceSpec_node_port),
				Ports: []*KubernetesServicePort{{Port: 80, NodePort: 8080}},
			}
			gomega.Expect(protovalidate.Validate(spec)).ToNot(gomega.BeNil())
		})

		ginkgo.It("rejects an uppercase port name", func() {
			spec := &KubernetesServiceSpec{
				Name:  "bad",
				Ports: []*KubernetesServicePort{{Name: "HTTP", Port: 80}},
			}
			gomega.Expect(protovalidate.Validate(spec)).ToNot(gomega.BeNil())
		})

		ginkgo.It("rejects an all-digit port name (must contain a letter)", func() {
			spec := &KubernetesServiceSpec{
				Name:  "bad",
				Ports: []*KubernetesServicePort{{Name: "8080", Port: 80}},
			}
			gomega.Expect(protovalidate.Validate(spec)).ToNot(gomega.BeNil())
		})

		ginkgo.It("rejects a numeric target port out of range", func() {
			spec := &KubernetesServiceSpec{
				Name:  "bad",
				Ports: []*KubernetesServicePort{{Port: 80, TargetPort: "99999"}},
			}
			gomega.Expect(protovalidate.Validate(spec)).ToNot(gomega.BeNil())
		})
	})
})
