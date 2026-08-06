package awsclientvpnv1alpha1

import (
	"testing"

	"buf.build/go/protovalidate"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	foreignkeyv1 "github.com/plantonhq/planton/apis/dev/planton/shared/foreignkey/v1"
)

func TestAwsClientVpnSpec(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "AwsClientVpnSpec Validation Suite")
}

func int32Ptr(i int32) *int32 {
	return &i
}

func svr(val string) *foreignkeyv1.StringValueOrRef {
	return &foreignkeyv1.StringValueOrRef{
		LiteralOrRef: &foreignkeyv1.StringValueOrRef_Value{Value: val},
	}
}

const certArn = "arn:aws:acm:us-west-2:123456789012:certificate/abc-123"

// certificateVpn is the common mutual-TLS shape: certificate authentication,
// one associated subnet, one authorization rule.
func certificateVpn() *AwsClientVpnSpec {
	return &AwsClientVpnSpec{
		Region: "us-west-2",
		AuthenticationOptions: []*AwsClientVpnAuthenticationOption{
			{
				Type:                    "certificate-authentication",
				RootCertificateChainArn: svr(certArn),
			},
		},
		ServerCertificateArn: svr(certArn),
		ClientCidrBlock:      "10.100.0.0/22",
		VpcId:                svr("vpc-12345678"),
		SubnetIds: []*foreignkeyv1.StringValueOrRef{
			svr("subnet-11111111"),
		},
		AuthorizationRules: []*AwsClientVpnAuthorizationRule{
			{
				TargetNetworkCidr:  "10.0.0.0/16",
				AuthorizeAllGroups: true,
			},
		},
	}
}

// federatedVpn is the SSO shape: SAML federation with the self-service portal.
func federatedVpn() *AwsClientVpnSpec {
	spec := certificateVpn()
	spec.AuthenticationOptions = []*AwsClientVpnAuthenticationOption{
		{
			Type:            "federated-authentication",
			SamlProviderArn: "arn:aws:iam::123456789012:saml-provider/okta",
		},
	}
	spec.SelfServicePortalEnabled = true
	return spec
}

var _ = ginkgo.Describe("AwsClientVpnSpec validations", func() {

	// -----------------------------------------------------------------
	// Valid inputs
	// -----------------------------------------------------------------
	ginkgo.Describe("When valid input is passed", func() {

		ginkgo.Context("with certificate authentication", func() {
			ginkgo.It("should not return a validation error", func() {
				err := protovalidate.Validate(certificateVpn())
				gomega.Expect(err).To(gomega.BeNil())
			})
		})

		ginkgo.Context("with federated authentication and the self-service portal", func() {
			ginkgo.It("should not return a validation error", func() {
				err := protovalidate.Validate(federatedVpn())
				gomega.Expect(err).To(gomega.BeNil())
			})
		})

		ginkgo.Context("with directory-service authentication", func() {
			ginkgo.It("should not return a validation error", func() {
				spec := certificateVpn()
				spec.AuthenticationOptions = []*AwsClientVpnAuthenticationOption{
					{
						Type:              "directory-service-authentication",
						ActiveDirectoryId: "d-1234567890",
					},
				}
				err := protovalidate.Validate(spec)
				gomega.Expect(err).To(gomega.BeNil())
			})
		})

		ginkgo.Context("with two combined authentication options", func() {
			ginkgo.It("should not return a validation error", func() {
				spec := certificateVpn()
				spec.AuthenticationOptions = append(spec.AuthenticationOptions,
					&AwsClientVpnAuthenticationOption{
						Type:            "federated-authentication",
						SamlProviderArn: "arn:aws:iam::123456789012:saml-provider/okta",
					})
				err := protovalidate.Validate(spec)
				gomega.Expect(err).To(gomega.BeNil())
			})
		})

		ginkgo.Context("with a zero-association pre-staged endpoint", func() {
			ginkgo.It("should not return a validation error", func() {
				spec := certificateVpn()
				spec.SubnetIds = nil
				spec.AuthorizationRules = nil
				err := protovalidate.Validate(spec)
				gomega.Expect(err).To(gomega.BeNil())
			})
		})

		ginkgo.Context("with a full-tunnel route through a NAT-ed subnet", func() {
			ginkgo.It("should not return a validation error", func() {
				spec := certificateVpn()
				spec.SplitTunnel = false
				spec.Routes = []*AwsClientVpnRoute{
					{
						DestinationCidrBlock: "0.0.0.0/0",
						TargetSubnetId:       svr("subnet-11111111"),
						Description:          "internet egress",
					},
				}
				spec.AuthorizationRules = append(spec.AuthorizationRules,
					&AwsClientVpnAuthorizationRule{
						TargetNetworkCidr:  "0.0.0.0/0",
						AuthorizeAllGroups: true,
					})
				err := protovalidate.Validate(spec)
				gomega.Expect(err).To(gomega.BeNil())
			})
		})

		ginkgo.Context("with a group-scoped authorization rule", func() {
			ginkgo.It("should not return a validation error", func() {
				spec := federatedVpn()
				spec.AuthorizationRules = []*AwsClientVpnAuthorizationRule{
					{
						TargetNetworkCidr: "10.0.0.0/16",
						AccessGroupId:     "engineering",
					},
				}
				err := protovalidate.Validate(spec)
				gomega.Expect(err).To(gomega.BeNil())
			})
		})

		ginkgo.Context("with a transit gateway attachment", func() {
			ginkgo.It("should not return a validation error", func() {
				spec := certificateVpn()
				spec.VpcId = nil
				spec.SecurityGroupIds = nil
				spec.SubnetIds = nil
				spec.TransitGatewayConfiguration = &AwsClientVpnTransitGatewayConfiguration{
					TransitGatewayId: svr("tgw-12345678"),
				}
				err := protovalidate.Validate(spec)
				gomega.Expect(err).To(gomega.BeNil())
			})
		})

		ginkgo.Context("with pure-IPv6 tunnel traffic and no client CIDR", func() {
			ginkgo.It("should not return a validation error", func() {
				spec := certificateVpn()
				spec.TrafficIpAddressType = "ipv6"
				spec.ClientCidrBlock = ""
				err := protovalidate.Validate(spec)
				gomega.Expect(err).To(gomega.BeNil())
			})
		})

		ginkgo.Context("with session and client-experience dials", func() {
			ginkgo.It("should not return a validation error", func() {
				spec := certificateVpn()
				spec.SessionTimeoutHours = int32Ptr(12)
				spec.DisconnectOnSessionTimeout = true
				spec.ClientRouteEnforcementEnabled = true
				spec.SplitTunnel = true
				spec.TransportProtocol = "tcp"
				spec.VpnPort = int32Ptr(1194)
				spec.DnsServers = []string{"10.0.0.2"}
				spec.ClientLoginBanner = &AwsClientVpnLoginBanner{
					BannerText: "Authorized use only.",
				}
				spec.ConnectionLog = &AwsClientVpnConnectionLog{
					CloudwatchLogGroup: svr("/vpn/corp-access"),
				}
				spec.ClientConnectOptions = &AwsClientVpnClientConnectOptions{
					LambdaFunctionArn: svr("arn:aws:lambda:us-west-2:123456789012:function:AWSClientVPN-posture"),
				}
				err := protovalidate.Validate(spec)
				gomega.Expect(err).To(gomega.BeNil())
			})
		})
	})

	// -----------------------------------------------------------------
	// Invalid inputs
	// -----------------------------------------------------------------
	ginkgo.Describe("When invalid input is passed", func() {

		ginkgo.Context("authentication options", func() {

			ginkgo.It("should reject zero authentication options", func() {
				spec := certificateVpn()
				spec.AuthenticationOptions = nil
				err := protovalidate.Validate(spec)
				gomega.Expect(err).ToNot(gomega.BeNil())
			})

			ginkgo.It("should reject three authentication options", func() {
				spec := certificateVpn()
				spec.AuthenticationOptions = []*AwsClientVpnAuthenticationOption{
					{Type: "certificate-authentication", RootCertificateChainArn: svr(certArn)},
					{Type: "directory-service-authentication", ActiveDirectoryId: "d-1234567890"},
					{Type: "federated-authentication", SamlProviderArn: "arn:aws:iam::123456789012:saml-provider/okta"},
				}
				err := protovalidate.Validate(spec)
				gomega.Expect(err).ToNot(gomega.BeNil())
			})

			ginkgo.It("should reject an unknown authentication type", func() {
				spec := certificateVpn()
				spec.AuthenticationOptions[0].Type = "cognito"
				err := protovalidate.Validate(spec)
				gomega.Expect(err).ToNot(gomega.BeNil())
			})

			ginkgo.It("should reject certificate authentication without a root chain", func() {
				spec := certificateVpn()
				spec.AuthenticationOptions[0].RootCertificateChainArn = nil
				err := protovalidate.Validate(spec)
				gomega.Expect(err).ToNot(gomega.BeNil())
			})

			ginkgo.It("should reject directory authentication without a directory id", func() {
				spec := certificateVpn()
				spec.AuthenticationOptions = []*AwsClientVpnAuthenticationOption{
					{Type: "directory-service-authentication"},
				}
				err := protovalidate.Validate(spec)
				gomega.Expect(err).ToNot(gomega.BeNil())
			})

			ginkgo.It("should reject federated authentication without a SAML provider", func() {
				spec := certificateVpn()
				spec.AuthenticationOptions = []*AwsClientVpnAuthenticationOption{
					{Type: "federated-authentication"},
				}
				err := protovalidate.Validate(spec)
				gomega.Expect(err).ToNot(gomega.BeNil())
			})

			ginkgo.It("should reject a SAML provider on a certificate option", func() {
				spec := certificateVpn()
				spec.AuthenticationOptions[0].SamlProviderArn = "arn:aws:iam::123456789012:saml-provider/okta"
				err := protovalidate.Validate(spec)
				gomega.Expect(err).ToNot(gomega.BeNil())
			})

			ginkgo.It("should reject a directory id on a federated option", func() {
				spec := federatedVpn()
				spec.AuthenticationOptions[0].ActiveDirectoryId = "d-1234567890"
				err := protovalidate.Validate(spec)
				gomega.Expect(err).ToNot(gomega.BeNil())
			})
		})

		ginkgo.Context("endpoint shape", func() {

			ginkgo.It("should reject a missing server certificate", func() {
				spec := certificateVpn()
				spec.ServerCertificateArn = nil
				err := protovalidate.Validate(spec)
				gomega.Expect(err).ToNot(gomega.BeNil())
			})

			ginkgo.It("should reject a missing client_cidr_block on IPv4 traffic", func() {
				spec := certificateVpn()
				spec.ClientCidrBlock = ""
				err := protovalidate.Validate(spec)
				gomega.Expect(err).ToNot(gomega.BeNil())
			})

			ginkgo.It("should reject a client_cidr_block on pure-IPv6 traffic", func() {
				spec := certificateVpn()
				spec.TrafficIpAddressType = "ipv6"
				err := protovalidate.Validate(spec)
				gomega.Expect(err).ToNot(gomega.BeNil())
			})

			ginkgo.It("should reject a malformed client_cidr_block", func() {
				spec := certificateVpn()
				spec.ClientCidrBlock = "10.100.0.0"
				err := protovalidate.Validate(spec)
				gomega.Expect(err).ToNot(gomega.BeNil())
			})

			ginkgo.It("should reject an invalid vpn_port", func() {
				spec := certificateVpn()
				spec.VpnPort = int32Ptr(8443)
				err := protovalidate.Validate(spec)
				gomega.Expect(err).ToNot(gomega.BeNil())
			})

			ginkgo.It("should reject an invalid transport protocol", func() {
				spec := certificateVpn()
				spec.TransportProtocol = "sctp"
				err := protovalidate.Validate(spec)
				gomega.Expect(err).ToNot(gomega.BeNil())
			})

			ginkgo.It("should reject an invalid session timeout", func() {
				spec := certificateVpn()
				spec.SessionTimeoutHours = int32Ptr(9)
				err := protovalidate.Validate(spec)
				gomega.Expect(err).ToNot(gomega.BeNil())
			})

			ginkgo.It("should reject an invalid endpoint IP address type", func() {
				spec := certificateVpn()
				spec.EndpointIpAddressType = "dualstack"
				err := protovalidate.Validate(spec)
				gomega.Expect(err).ToNot(gomega.BeNil())
			})

			ginkgo.It("should reject more than two DNS servers", func() {
				spec := certificateVpn()
				spec.DnsServers = []string{"10.0.0.2", "10.0.0.3", "10.0.0.4"}
				err := protovalidate.Validate(spec)
				gomega.Expect(err).ToNot(gomega.BeNil())
			})

			ginkgo.It("should reject a non-IP DNS server", func() {
				spec := certificateVpn()
				spec.DnsServers = []string{"resolver.internal"}
				err := protovalidate.Validate(spec)
				gomega.Expect(err).ToNot(gomega.BeNil())
			})

			ginkgo.It("should reject more than five security groups", func() {
				spec := certificateVpn()
				spec.SecurityGroupIds = []*foreignkeyv1.StringValueOrRef{
					svr("sg-1"), svr("sg-2"), svr("sg-3"), svr("sg-4"), svr("sg-5"), svr("sg-6"),
				}
				err := protovalidate.Validate(spec)
				gomega.Expect(err).ToNot(gomega.BeNil())
			})

			ginkgo.It("should reject the self-service portal without federated auth", func() {
				spec := certificateVpn()
				spec.SelfServicePortalEnabled = true
				err := protovalidate.Validate(spec)
				gomega.Expect(err).ToNot(gomega.BeNil())
			})
		})

		ginkgo.Context("transit gateway exclusivity", func() {

			ginkgo.It("should reject a transit gateway together with a VPC", func() {
				spec := certificateVpn()
				spec.SubnetIds = nil
				spec.TransitGatewayConfiguration = &AwsClientVpnTransitGatewayConfiguration{
					TransitGatewayId: svr("tgw-12345678"),
				}
				err := protovalidate.Validate(spec)
				gomega.Expect(err).ToNot(gomega.BeNil())
			})

			ginkgo.It("should reject a transit gateway together with subnet associations", func() {
				spec := certificateVpn()
				spec.VpcId = nil
				spec.TransitGatewayConfiguration = &AwsClientVpnTransitGatewayConfiguration{
					TransitGatewayId: svr("tgw-12345678"),
				}
				err := protovalidate.Validate(spec)
				gomega.Expect(err).ToNot(gomega.BeNil())
			})

			ginkgo.It("should reject both AZ addressing forms", func() {
				spec := certificateVpn()
				spec.VpcId = nil
				spec.SubnetIds = nil
				spec.SecurityGroupIds = nil
				spec.TransitGatewayConfiguration = &AwsClientVpnTransitGatewayConfiguration{
					TransitGatewayId:    svr("tgw-12345678"),
					AvailabilityZones:   []string{"us-west-2a"},
					AvailabilityZoneIds: []string{"usw2-az1"},
				}
				err := protovalidate.Validate(spec)
				gomega.Expect(err).ToNot(gomega.BeNil())
			})
		})

		ginkgo.Context("authorization rules and routes", func() {

			ginkgo.It("should reject a rule with neither grantee arm", func() {
				spec := certificateVpn()
				spec.AuthorizationRules = []*AwsClientVpnAuthorizationRule{
					{TargetNetworkCidr: "10.0.0.0/16"},
				}
				err := protovalidate.Validate(spec)
				gomega.Expect(err).ToNot(gomega.BeNil())
			})

			ginkgo.It("should reject a rule with both grantee arms", func() {
				spec := certificateVpn()
				spec.AuthorizationRules = []*AwsClientVpnAuthorizationRule{
					{
						TargetNetworkCidr:  "10.0.0.0/16",
						AccessGroupId:      "engineering",
						AuthorizeAllGroups: true,
					},
				}
				err := protovalidate.Validate(spec)
				gomega.Expect(err).ToNot(gomega.BeNil())
			})

			ginkgo.It("should reject a rule with a malformed CIDR", func() {
				spec := certificateVpn()
				spec.AuthorizationRules[0].TargetNetworkCidr = "not-a-cidr"
				err := protovalidate.Validate(spec)
				gomega.Expect(err).ToNot(gomega.BeNil())
			})

			ginkgo.It("should reject a route without a target subnet", func() {
				spec := certificateVpn()
				spec.Routes = []*AwsClientVpnRoute{
					{DestinationCidrBlock: "0.0.0.0/0"},
				}
				err := protovalidate.Validate(spec)
				gomega.Expect(err).ToNot(gomega.BeNil())
			})

			ginkgo.It("should reject a banner over 1400 characters", func() {
				spec := certificateVpn()
				long := make([]byte, 1401)
				for i := range long {
					long[i] = 'a'
				}
				spec.ClientLoginBanner = &AwsClientVpnLoginBanner{BannerText: string(long)}
				err := protovalidate.Validate(spec)
				gomega.Expect(err).ToNot(gomega.BeNil())
			})
		})
	})
})
