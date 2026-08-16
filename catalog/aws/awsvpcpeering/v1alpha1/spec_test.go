package awsvpcpeeringv1alpha1

import (
	"testing"

	"buf.build/go/protovalidate"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	foreignkeyv1 "github.com/plantonhq/planton/shared/foreignkey/v1"
)

func TestAwsVpcPeeringSpec(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "AwsVpcPeeringSpec Validation Suite")
}

func svr(val string) *foreignkeyv1.StringValueOrRef {
	return &foreignkeyv1.StringValueOrRef{
		LiteralOrRef: &foreignkeyv1.StringValueOrRef_Value{Value: val},
	}
}

func requestSpec() *AwsVpcPeeringSpec {
	return &AwsVpcPeeringSpec{
		Region: "us-east-1",
		Request: &AwsVpcPeeringRequest{
			VpcId:      svr("vpc-0123456789abcdef0"),
			PeerVpcId:  svr("vpc-0fedcba9876543210"),
			AutoAccept: true,
		},
	}
}

func acceptSpec() *AwsVpcPeeringSpec {
	return &AwsVpcPeeringSpec{
		Region: "us-west-2",
		Accept: &AwsVpcPeeringAccept{
			VpcPeeringConnectionId: svr("pcx-0123456789abcdef0"),
			AutoAccept:             true,
		},
	}
}

var _ = ginkgo.Describe("AwsVpcPeeringSpec validations", func() {

	ginkgo.Describe("When valid input is passed", func() {

		ginkgo.It("accepts a same-account auto-accepted request", func() {
			gomega.Expect(protovalidate.Validate(requestSpec())).To(gomega.BeNil())
		})

		ginkgo.It("accepts an accept-arm instance", func() {
			gomega.Expect(protovalidate.Validate(acceptSpec())).To(gomega.BeNil())
		})

		ginkgo.It("accepts a cross-region request without auto_accept", func() {
			spec := requestSpec()
			spec.Request.AutoAccept = false
			spec.Request.PeerRegion = "eu-west-1"
			gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
		})

		ginkgo.It("accepts a cross-account request without auto_accept", func() {
			spec := requestSpec()
			spec.Request.AutoAccept = false
			spec.Request.PeerOwnerId = "210987654321"
			gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
		})

		ginkgo.It("accepts DNS-resolution options on an auto-accepted request", func() {
			spec := requestSpec()
			spec.Request.RequesterAllowRemoteVpcDnsResolution = true
			spec.Request.AccepterAllowRemoteVpcDnsResolution = true
			gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
		})
	})

	ginkgo.Describe("When invalid input is passed", func() {

		ginkgo.It("rejects both arms configured", func() {
			spec := requestSpec()
			spec.Accept = acceptSpec().Accept
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects neither arm configured", func() {
			spec := requestSpec()
			spec.Request = nil
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects auto_accept with peer_region", func() {
			spec := requestSpec()
			spec.Request.PeerRegion = "eu-west-1"
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects auto_accept with peer_owner_id", func() {
			spec := requestSpec()
			spec.Request.PeerOwnerId = "210987654321"
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects the accepter DNS option on a non-auto-accepted request", func() {
			spec := requestSpec()
			spec.Request.AutoAccept = false
			spec.Request.AccepterAllowRemoteVpcDnsResolution = true
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects a request without a requester VPC", func() {
			spec := requestSpec()
			spec.Request.VpcId = nil
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects an accept arm without a connection id", func() {
			spec := acceptSpec()
			spec.Accept.VpcPeeringConnectionId = nil
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})
	})
})
