package awscloudwatchlogdeliveryv1alpha1

import (
	"testing"

	"buf.build/go/protovalidate"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	foreignkeyv1 "github.com/plantonhq/planton/shared/foreignkey/v1"
	"google.golang.org/protobuf/types/known/structpb"
)

func TestAwsCloudwatchLogDeliverySpec(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "AwsCloudwatchLogDeliverySpec Validation Suite")
}

func svr(val string) *foreignkeyv1.StringValueOrRef {
	return &foreignkeyv1.StringValueOrRef{
		LiteralOrRef: &foreignkeyv1.StringValueOrRef_Value{Value: val},
	}
}

func samplePolicy() *structpb.Struct {
	doc, err := structpb.NewStruct(map[string]any{
		"Version": "2012-10-17",
		"Statement": []any{
			map[string]any{
				"Effect":    "Allow",
				"Principal": map[string]any{"AWS": "arn:aws:iam::210987654321:root"},
				"Action":    "logs:PutSubscriptionFilter",
				"Resource":  "*",
			},
		},
	})
	if err != nil {
		panic(err)
	}
	return doc
}

func sampleSource() *AwsVendedLogSource {
	return &AwsVendedLogSource{
		Name:        "kb-application-logs",
		LogType:     "APPLICATION_LOGS",
		ResourceArn: svr("arn:aws:bedrock:us-east-1:123456789012:knowledge-base/KBID"),
	}
}

func sampleDestination() *AwsVendedLogDestination {
	return &AwsVendedLogDestination{
		Name:                   "central-s3",
		DestinationResourceArn: svr("arn:aws:s3:::central-log-archive"),
	}
}

func vendedSpec() *AwsCloudwatchLogDeliverySpec {
	return &AwsCloudwatchLogDeliverySpec{
		Region: "us-east-1",
		Vended: &AwsVendedLogDelivery{
			Source:       sampleSource(),
			Destinations: []*AwsVendedLogDestination{sampleDestination()},
			Deliveries: []*AwsVendedLogDeliveryEntry{
				{
					Name:            "to-s3",
					DestinationName: "central-s3",
				},
			},
		},
	}
}

func crossAccountSpec() *AwsCloudwatchLogDeliverySpec {
	return &AwsCloudwatchLogDeliverySpec{
		Region: "us-east-1",
		CrossAccountDestination: &AwsCrossAccountLogDestination{
			Name:         "org-log-sink",
			RoleArn:      svr("arn:aws:iam::123456789012:role/logs-to-kinesis"),
			TargetArn:    svr("arn:aws:kinesis:us-east-1:123456789012:stream/org-logs"),
			AccessPolicy: samplePolicy(),
		},
	}
}

var _ = ginkgo.Describe("AwsCloudwatchLogDeliverySpec validations", func() {

	ginkgo.Describe("When valid input is passed", func() {

		ginkgo.It("accepts a full vended pipeline", func() {
			gomega.Expect(protovalidate.Validate(vendedSpec())).To(gomega.BeNil())
		})

		ginkgo.It("accepts a source-only vended arm", func() {
			spec := vendedSpec()
			spec.Vended.Destinations = nil
			spec.Vended.Deliveries = nil
			gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
		})

		ginkgo.It("accepts a destinations-only vended arm (the central-team shape)", func() {
			spec := vendedSpec()
			spec.Vended.Source = nil
			spec.Vended.Deliveries = nil
			spec.Vended.Destinations[0].Policy = samplePolicy()
			gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
		})

		ginkgo.It("accepts a delivery to an external destination by ARN", func() {
			spec := vendedSpec()
			spec.Vended.Destinations = nil
			spec.Vended.Deliveries[0].DestinationName = ""
			spec.Vended.Deliveries[0].DestinationArn = svr("arn:aws:logs:us-east-1:210987654321:delivery-destination:central")
			gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
		})

		ginkgo.It("accepts an XRAY destination without a resource ARN", func() {
			spec := vendedSpec()
			spec.Vended.Destinations[0].DestinationResourceArn = nil
			spec.Vended.Destinations[0].DeliveryDestinationType = "XRAY"
			gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
		})

		ginkgo.It("accepts the cross-account arm", func() {
			gomega.Expect(protovalidate.Validate(crossAccountSpec())).To(gomega.BeNil())
		})
	})

	ginkgo.Describe("When invalid input is passed", func() {

		ginkgo.It("rejects an empty spec (no arm)", func() {
			spec := &AwsCloudwatchLogDeliverySpec{Region: "us-east-1"}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects an empty vended arm", func() {
			spec := &AwsCloudwatchLogDeliverySpec{
				Region: "us-east-1",
				Vended: &AwsVendedLogDelivery{},
			}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects deliveries without a source", func() {
			spec := vendedSpec()
			spec.Vended.Source = nil
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects a delivery naming an unknown owned destination", func() {
			spec := vendedSpec()
			spec.Vended.Deliveries[0].DestinationName = "not-declared"
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects a delivery with both destination shapes", func() {
			spec := vendedSpec()
			spec.Vended.Deliveries[0].DestinationArn = svr("arn:aws:logs:us-east-1:210987654321:delivery-destination:central")
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects a delivery with neither destination shape", func() {
			spec := vendedSpec()
			spec.Vended.Deliveries[0].DestinationName = ""
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects a non-XRAY destination without a resource ARN", func() {
			spec := vendedSpec()
			spec.Vended.Destinations[0].DestinationResourceArn = nil
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects an XRAY destination WITH a resource ARN", func() {
			spec := vendedSpec()
			spec.Vended.Destinations[0].DeliveryDestinationType = "XRAY"
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects duplicate delivery names", func() {
			spec := vendedSpec()
			spec.Vended.Deliveries = append(spec.Vended.Deliveries, &AwsVendedLogDeliveryEntry{
				Name:            "to-s3",
				DestinationName: "central-s3",
			})
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects a field delimiter above five characters", func() {
			spec := vendedSpec()
			spec.Vended.Deliveries[0].FieldDelimiter = "||||||"
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects a cross-account destination name with colons or asterisks", func() {
			spec := crossAccountSpec()
			spec.CrossAccountDestination.Name = "org:sink*"
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects a cross-account destination without its access policy", func() {
			spec := crossAccountSpec()
			spec.CrossAccountDestination.AccessPolicy = nil
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})
	})
})
