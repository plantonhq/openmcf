package mongodbatlasv1

import (
	"testing"

	"buf.build/go/protovalidate"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	"github.com/plantonhq/openmcf/apis/org/openmcf/shared"
)

func TestMongodbAtlas(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "MongodbAtlas Suite")
}

var _ = ginkgo.Describe("MongodbAtlas Validation Tests", func() {

	ginkgo.Describe("When valid input is passed", func() {
		ginkgo.Context("mongodb_atlas with minimal valid fields", func() {
			var input *MongodbAtlas

			ginkgo.BeforeEach(func() {
				input = &MongodbAtlas{
					ApiVersion: "atlas.openmcf.org/v1",
					Kind:       "MongodbAtlas",
					Metadata: &shared.CloudResourceMetadata{
						Name: "test-atlas-resource",
					},
					Spec: &MongodbAtlasSpec{
						ClusterConfig: &MongodbAtlasClusterConfig{
							ProjectId:                "some-project-id",
							ClusterType:              "REPLICASET",
							ElectableNodes:           3,
							Priority:                 7,
							ReadOnlyNodes:            0,
							CloudBackup:              true,
							AutoScalingDiskGbEnabled: true,
							MongoDbMajorVersion:      "7.0",
							ProviderName:             "AWS",
							ProviderInstanceSizeName: "M10",
						},
					},
				}
			})

			ginkgo.It("should not return a validation error", func() {
				err := protovalidate.Validate(input)
				gomega.Expect(err).To(gomega.BeNil())
			})
		})

		ginkgo.Context("mongodb_atlas with GCP provider", func() {
			ginkgo.It("should not return a validation error", func() {
				input := &MongodbAtlas{
					ApiVersion: "atlas.openmcf.org/v1",
					Kind:       "MongodbAtlas",
					Metadata: &shared.CloudResourceMetadata{
						Name: "test-mongodb-gcp",
					},
					Spec: &MongodbAtlasSpec{
						ClusterConfig: &MongodbAtlasClusterConfig{
							ProjectId:                "test-project-id-gcp",
							ClusterType:              "REPLICASET",
							ElectableNodes:           5,
							Priority:                 7,
							ReadOnlyNodes:            2,
							CloudBackup:              true,
							AutoScalingDiskGbEnabled: false,
							MongoDbMajorVersion:      "6.0",
							ProviderName:             "GCP",
							ProviderInstanceSizeName: "M30",
						},
					},
				}
				err := protovalidate.Validate(input)
				gomega.Expect(err).To(gomega.BeNil())
			})
		})

		ginkgo.Context("mongodb_atlas with Azure provider", func() {
			ginkgo.It("should not return a validation error", func() {
				input := &MongodbAtlas{
					ApiVersion: "atlas.openmcf.org/v1",
					Kind:       "MongodbAtlas",
					Metadata: &shared.CloudResourceMetadata{
						Name: "test-mongodb-azure",
					},
					Spec: &MongodbAtlasSpec{
						ClusterConfig: &MongodbAtlasClusterConfig{
							ProjectId:                "test-project-id-azure",
							ClusterType:              "SHARDED",
							ElectableNodes:           7,
							Priority:                 7,
							ReadOnlyNodes:            0,
							CloudBackup:              true,
							AutoScalingDiskGbEnabled: true,
							MongoDbMajorVersion:      "5.0",
							ProviderName:             "AZURE",
							ProviderInstanceSizeName: "M50",
						},
					},
				}
				err := protovalidate.Validate(input)
				gomega.Expect(err).To(gomega.BeNil())
			})
		})

		ginkgo.Context("mongodb_atlas with GEOSHARDED cluster type", func() {
			ginkgo.It("should not return a validation error", func() {
				input := &MongodbAtlas{
					ApiVersion: "atlas.openmcf.org/v1",
					Kind:       "MongodbAtlas",
					Metadata: &shared.CloudResourceMetadata{
						Name: "test-mongodb-geosharded",
					},
					Spec: &MongodbAtlasSpec{
						ClusterConfig: &MongodbAtlasClusterConfig{
							ProjectId:                "test-project-id-geo",
							ClusterType:              "GEOSHARDED",
							ElectableNodes:           5,
							Priority:                 7,
							ReadOnlyNodes:            0,
							CloudBackup:              true,
							AutoScalingDiskGbEnabled: true,
							MongoDbMajorVersion:      "7.0",
							ProviderName:             "AWS",
							ProviderInstanceSizeName: "M40",
						},
					},
				}
				err := protovalidate.Validate(input)
				gomega.Expect(err).To(gomega.BeNil())
			})
		})
	})

	ginkgo.Describe("When invalid input is passed", func() {
		ginkgo.Context("missing required metadata", func() {
			ginkgo.It("should return a validation error", func() {
				input := &MongodbAtlas{
					ApiVersion: "atlas.openmcf.org/v1",
					Kind:       "MongodbAtlas",
					Metadata:   nil,
					Spec: &MongodbAtlasSpec{
						ClusterConfig: &MongodbAtlasClusterConfig{
							ProjectId:                "test-project-id",
							ClusterType:              "REPLICASET",
							ElectableNodes:           3,
							Priority:                 7,
							ReadOnlyNodes:            0,
							CloudBackup:              true,
							AutoScalingDiskGbEnabled: true,
							MongoDbMajorVersion:      "7.0",
							ProviderName:             "AWS",
							ProviderInstanceSizeName: "M10",
						},
					},
				}
				err := protovalidate.Validate(input)
				gomega.Expect(err).NotTo(gomega.BeNil())
			})
		})

		ginkgo.Context("missing required spec", func() {
			ginkgo.It("should return a validation error", func() {
				input := &MongodbAtlas{
					ApiVersion: "atlas.openmcf.org/v1",
					Kind:       "MongodbAtlas",
					Metadata: &shared.CloudResourceMetadata{
						Name: "test-mongodb-atlas",
					},
					Spec: nil,
				}
				err := protovalidate.Validate(input)
				gomega.Expect(err).NotTo(gomega.BeNil())
			})
		})

		ginkgo.Context("incorrect api_version", func() {
			ginkgo.It("should return a validation error", func() {
				input := &MongodbAtlas{
					ApiVersion: "wrong.api.version/v1",
					Kind:       "MongodbAtlas",
					Metadata: &shared.CloudResourceMetadata{
						Name: "test-mongodb-atlas",
					},
					Spec: &MongodbAtlasSpec{
						ClusterConfig: &MongodbAtlasClusterConfig{
							ProjectId:                "test-project-id",
							ClusterType:              "REPLICASET",
							ElectableNodes:           3,
							Priority:                 7,
							ReadOnlyNodes:            0,
							CloudBackup:              true,
							AutoScalingDiskGbEnabled: true,
							MongoDbMajorVersion:      "7.0",
							ProviderName:             "AWS",
							ProviderInstanceSizeName: "M10",
						},
					},
				}
				err := protovalidate.Validate(input)
				gomega.Expect(err).NotTo(gomega.BeNil())
			})
		})

		ginkgo.Context("incorrect kind", func() {
			ginkgo.It("should return a validation error", func() {
				input := &MongodbAtlas{
					ApiVersion: "atlas.openmcf.org/v1",
					Kind:       "WrongKind",
					Metadata: &shared.CloudResourceMetadata{
						Name: "test-mongodb-atlas",
					},
					Spec: &MongodbAtlasSpec{
						ClusterConfig: &MongodbAtlasClusterConfig{
							ProjectId:                "test-project-id",
							ClusterType:              "REPLICASET",
							ElectableNodes:           3,
							Priority:                 7,
							ReadOnlyNodes:            0,
							CloudBackup:              true,
							AutoScalingDiskGbEnabled: true,
							MongoDbMajorVersion:      "7.0",
							ProviderName:             "AWS",
							ProviderInstanceSizeName: "M10",
						},
					},
				}
				err := protovalidate.Validate(input)
				gomega.Expect(err).NotTo(gomega.BeNil())
			})
		})
	})
})
