package awscloudmapnamespacev1alpha1

import (
	"testing"

	"buf.build/go/protovalidate"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	foreignkeyv1 "github.com/plantonhq/planton/shared/foreignkey/v1"
)

func TestAwsCloudMapNamespaceSpec(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "AwsCloudMapNamespaceSpec Validation Suite")
}

func literal(value string) *foreignkeyv1.StringValueOrRef {
	return &foreignkeyv1.StringValueOrRef{
		LiteralOrRef: &foreignkeyv1.StringValueOrRef_Value{Value: value},
	}
}

func privateNamespace() *AwsCloudMapNamespaceSpec {
	return &AwsCloudMapNamespaceSpec{
		Region: "us-west-2",
		Type:   "PRIVATE_DNS",
		VpcId:  literal("vpc-0123456789abcdef0"),
		Services: []*AwsCloudMapService{
			{
				Name: "api",
				DnsConfig: &AwsCloudMapServiceDnsConfig{
					Records: []*AwsCloudMapServiceDnsRecord{{Type: "A", Ttl: 10}},
				},
				HealthCheckCustomConfig: &AwsCloudMapServiceHealthCheckCustomConfig{},
				Instances: []*AwsCloudMapInstance{
					{InstanceId: "static-1", Ip: "10.0.0.10"},
				},
			},
		},
	}
}

var _ = ginkgo.Describe("AwsCloudMapNamespaceSpec validations", func() {

	ginkgo.Describe("When valid input is passed", func() {

		ginkgo.It("accepts a private DNS namespace with a service and instance", func() {
			gomega.Expect(protovalidate.Validate(privateNamespace())).To(gomega.BeNil())
		})

		ginkgo.It("accepts an HTTP namespace with an API-only service", func() {
			spec := &AwsCloudMapNamespaceSpec{
				Region: "us-west-2",
				Type:   "HTTP",
				Services: []*AwsCloudMapService{
					{Name: "workers", Instances: []*AwsCloudMapInstance{
						{InstanceId: "w-1", CustomAttributes: map[string]string{"stage": "prod"}},
					}},
				},
			}
			gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
		})

		ginkgo.It("accepts a public DNS namespace with Route 53 health checks", func() {
			spec := &AwsCloudMapNamespaceSpec{
				Region: "us-west-2",
				Type:   "PUBLIC_DNS",
				Services: []*AwsCloudMapService{
					{
						Name: "edge",
						DnsConfig: &AwsCloudMapServiceDnsConfig{
							Records:       []*AwsCloudMapServiceDnsRecord{{Type: "A", Ttl: 60}},
							RoutingPolicy: "WEIGHTED",
						},
						HealthCheckConfig: &AwsCloudMapServiceHealthCheckConfig{
							Type:             "HTTPS",
							ResourcePath:     "/health",
							FailureThreshold: 2,
						},
					},
				},
			}
			gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
		})

		ginkgo.It("accepts a CNAME instance pointing at an endpoint", func() {
			spec := privateNamespace()
			spec.Services[0].DnsConfig.Records = []*AwsCloudMapServiceDnsRecord{{Type: "CNAME", Ttl: 30}}
			spec.Services[0].Instances = []*AwsCloudMapInstance{
				{InstanceId: "db", Cname: "mydb.cluster-abc.us-west-2.rds.amazonaws.com"},
			}
			gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
		})

		ginkgo.It("accepts an alias instance standing alone", func() {
			spec := privateNamespace()
			spec.Services[0].Instances = []*AwsCloudMapInstance{
				{InstanceId: "lb", AliasDnsName: literal("internal-alb-123.us-west-2.elb.amazonaws.com")},
			}
			gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
		})

		ginkgo.It("accepts an EC2-derived instance without an ip", func() {
			spec := privateNamespace()
			spec.Services[0].Instances = []*AwsCloudMapInstance{
				{InstanceId: "vm", Ec2InstanceId: literal("i-0123456789abcdef0")},
			}
			gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
		})
	})

	ginkgo.Describe("When invalid input is passed", func() {

		ginkgo.It("rejects a private DNS namespace without a vpc", func() {
			spec := privateNamespace()
			spec.VpcId = nil
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects a vpc on a public DNS namespace", func() {
			spec := privateNamespace()
			spec.Type = "PUBLIC_DNS"
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects dns_config on an HTTP namespace service", func() {
			spec := privateNamespace()
			spec.Type = "HTTP"
			spec.VpcId = nil
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects Route 53 health checks outside a public namespace", func() {
			spec := privateNamespace()
			spec.Services[0].HealthCheckCustomConfig = nil
			spec.Services[0].HealthCheckConfig = &AwsCloudMapServiceHealthCheckConfig{Type: "HTTP"}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects both health check forms on one service", func() {
			spec := &AwsCloudMapNamespaceSpec{
				Region: "us-west-2",
				Type:   "PUBLIC_DNS",
				Services: []*AwsCloudMapService{
					{
						Name:                    "edge",
						DnsConfig:               &AwsCloudMapServiceDnsConfig{Records: []*AwsCloudMapServiceDnsRecord{{Type: "A", Ttl: 60}}},
						HealthCheckConfig:       &AwsCloudMapServiceHealthCheckConfig{Type: "HTTP"},
						HealthCheckCustomConfig: &AwsCloudMapServiceHealthCheckCustomConfig{},
					},
				},
			}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects an alias combined with an ip", func() {
			spec := privateNamespace()
			spec.Services[0].Instances = []*AwsCloudMapInstance{
				{InstanceId: "lb", AliasDnsName: literal("alb.example.com"), Ip: "10.0.0.10"},
			}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects an ec2 instance id combined with an ip", func() {
			spec := privateNamespace()
			spec.Services[0].Instances = []*AwsCloudMapInstance{
				{InstanceId: "vm", Ec2InstanceId: literal("i-0123456789abcdef0"), Ip: "10.0.0.10"},
			}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects an instance with both ip and cname", func() {
			spec := privateNamespace()
			spec.Services[0].Instances = []*AwsCloudMapInstance{
				{InstanceId: "x", Ip: "10.0.0.10", Cname: "db.example.com"},
			}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects duplicate service names", func() {
			spec := privateNamespace()
			spec.Services = append(spec.Services, &AwsCloudMapService{Name: "api"})
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects duplicate instance ids", func() {
			spec := privateNamespace()
			spec.Services[0].Instances = append(spec.Services[0].Instances,
				&AwsCloudMapInstance{InstanceId: "static-1", Ip: "10.0.0.11"})
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects a dns_config with no records", func() {
			spec := privateNamespace()
			spec.Services[0].DnsConfig.Records = nil
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects a zero record ttl", func() {
			spec := privateNamespace()
			spec.Services[0].DnsConfig.Records[0].Ttl = 0
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects a resource path on a TCP health check", func() {
			spec := &AwsCloudMapNamespaceSpec{
				Region: "us-west-2",
				Type:   "PUBLIC_DNS",
				Services: []*AwsCloudMapService{
					{
						Name:              "edge",
						DnsConfig:         &AwsCloudMapServiceDnsConfig{Records: []*AwsCloudMapServiceDnsRecord{{Type: "A", Ttl: 60}}},
						HealthCheckConfig: &AwsCloudMapServiceHealthCheckConfig{Type: "TCP", ResourcePath: "/health"},
					},
				},
			}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})
	})
})
