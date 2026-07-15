package awsbatchschedulingpolicyv1

import (
	"testing"

	"buf.build/go/protovalidate"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

func TestAwsBatchSchedulingPolicySpec(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "AwsBatchSchedulingPolicySpec Validation Suite")
}

func int32Ptr(i int32) *int32 {
	return &i
}

var _ = ginkgo.Describe("AwsBatchSchedulingPolicySpec validations", func() {

	ginkgo.Describe("When valid input is passed", func() {

		ginkgo.Context("with only a region (all fair-share dials at their defaults)", func() {
			ginkgo.It("should not return a validation error", func() {
				err := protovalidate.Validate(&AwsBatchSchedulingPolicySpec{Region: "us-west-2"})
				gomega.Expect(err).To(gomega.BeNil())
			})
		})

		ginkgo.Context("with the full fair-share surface", func() {
			ginkgo.It("should not return a validation error", func() {
				spec := &AwsBatchSchedulingPolicySpec{
					Region:             "us-west-2",
					ComputeReservation: int32Ptr(10),
					ShareDecaySeconds:  int32Ptr(3600),
					ShareDistributions: []*AwsBatchShareDistribution{
						{ShareIdentifier: "teamA", WeightFactor: 0.5},
						{ShareIdentifier: "teamB", WeightFactor: 1.0},
						{ShareIdentifier: "adhoc*", WeightFactor: 2.0},
					},
				}
				err := protovalidate.Validate(spec)
				gomega.Expect(err).To(gomega.BeNil())
			})
		})

		ginkgo.Context("with a share distribution leaving weight_factor unset", func() {
			ginkgo.It("should not return a validation error (AWS defaults the weight to 1.0)", func() {
				spec := &AwsBatchSchedulingPolicySpec{
					Region: "us-west-2",
					ShareDistributions: []*AwsBatchShareDistribution{
						{ShareIdentifier: "teamA"},
					},
				}
				err := protovalidate.Validate(spec)
				gomega.Expect(err).To(gomega.BeNil())
			})
		})

		ginkgo.Context("with boundary values", func() {
			ginkgo.It("should accept compute_reservation 0 and 99 and decay 0 and 604800", func() {
				for _, cr := range []int32{0, 99} {
					for _, decay := range []int32{0, 604800} {
						spec := &AwsBatchSchedulingPolicySpec{
							Region:             "us-west-2",
							ComputeReservation: int32Ptr(cr),
							ShareDecaySeconds:  int32Ptr(decay),
						}
						gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
					}
				}
			})
		})
	})

	ginkgo.Describe("When invalid input is passed", func() {

		ginkgo.Context("with no region", func() {
			ginkgo.It("should return a validation error", func() {
				err := protovalidate.Validate(&AwsBatchSchedulingPolicySpec{})
				gomega.Expect(err).ToNot(gomega.BeNil())
			})
		})

		ginkgo.Context("with compute_reservation above 99", func() {
			ginkgo.It("should return a validation error", func() {
				spec := &AwsBatchSchedulingPolicySpec{
					Region:             "us-west-2",
					ComputeReservation: int32Ptr(100),
				}
				err := protovalidate.Validate(spec)
				gomega.Expect(err).ToNot(gomega.BeNil())
			})
		})

		ginkgo.Context("with share_decay_seconds above 7 days", func() {
			ginkgo.It("should return a validation error", func() {
				spec := &AwsBatchSchedulingPolicySpec{
					Region:            "us-west-2",
					ShareDecaySeconds: int32Ptr(604801),
				}
				err := protovalidate.Validate(spec)
				gomega.Expect(err).ToNot(gomega.BeNil())
			})
		})

		ginkgo.Context("with an invalid share identifier", func() {
			ginkgo.It("should reject a wildcard in the middle", func() {
				spec := &AwsBatchSchedulingPolicySpec{
					Region: "us-west-2",
					ShareDistributions: []*AwsBatchShareDistribution{
						{ShareIdentifier: "team*a", WeightFactor: 1.0},
					},
				}
				err := protovalidate.Validate(spec)
				gomega.Expect(err).ToNot(gomega.BeNil())
			})
			ginkgo.It("should reject a missing share identifier", func() {
				spec := &AwsBatchSchedulingPolicySpec{
					Region: "us-west-2",
					ShareDistributions: []*AwsBatchShareDistribution{
						{WeightFactor: 1.0},
					},
				}
				err := protovalidate.Validate(spec)
				gomega.Expect(err).ToNot(gomega.BeNil())
			})
		})

		ginkgo.Context("with an out-of-range weight factor", func() {
			ginkgo.It("should reject a weight below the AWS minimum", func() {
				spec := &AwsBatchSchedulingPolicySpec{
					Region: "us-west-2",
					ShareDistributions: []*AwsBatchShareDistribution{
						{ShareIdentifier: "teamA", WeightFactor: 0.00001},
					},
				}
				err := protovalidate.Validate(spec)
				gomega.Expect(err).ToNot(gomega.BeNil())
			})
			ginkgo.It("should reject a weight above 999.9999", func() {
				spec := &AwsBatchSchedulingPolicySpec{
					Region: "us-west-2",
					ShareDistributions: []*AwsBatchShareDistribution{
						{ShareIdentifier: "teamA", WeightFactor: 1000},
					},
				}
				err := protovalidate.Validate(spec)
				gomega.Expect(err).ToNot(gomega.BeNil())
			})
			ginkgo.It("should reject a negative weight", func() {
				spec := &AwsBatchSchedulingPolicySpec{
					Region: "us-west-2",
					ShareDistributions: []*AwsBatchShareDistribution{
						{ShareIdentifier: "teamA", WeightFactor: -1},
					},
				}
				err := protovalidate.Validate(spec)
				gomega.Expect(err).ToNot(gomega.BeNil())
			})
		})
	})
})
