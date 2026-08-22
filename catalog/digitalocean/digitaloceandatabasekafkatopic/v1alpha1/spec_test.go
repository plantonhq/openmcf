package digitaloceandatabasekafkatopicv1alpha1

import (
	"testing"

	"buf.build/go/protovalidate"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	fk "github.com/plantonhq/planton/shared/foreignkey/v1"
	"google.golang.org/protobuf/proto"
)

func TestDigitalOceanDatabaseKafkaTopicSpec(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "DigitalOceanDatabaseKafkaTopicSpec Validation Suite")
}

var _ = ginkgo.Describe("DigitalOceanDatabaseKafkaTopicSpec validations", func() {

	newClusterRef := func(clusterId string) *fk.StringValueOrRef {
		return &fk.StringValueOrRef{
			LiteralOrRef: &fk.StringValueOrRef_Value{Value: clusterId},
		}
	}

	makeValidSpec := func() *DigitalOceanDatabaseKafkaTopicSpec {
		return &DigitalOceanDatabaseKafkaTopicSpec{
			Cluster:   newClusterRef("aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee"),
			TopicName: "orders-events",
		}
	}

	ginkgo.Context("Required fields", func() {
		ginkgo.It("accepts a minimal valid spec", func() {
			err := protovalidate.Validate(makeValidSpec())
			gomega.Expect(err).To(gomega.BeNil())
		})

		ginkgo.It("accepts a cluster reference by name", func() {
			spec := makeValidSpec()
			spec.Cluster = &fk.StringValueOrRef{
				LiteralOrRef: &fk.StringValueOrRef_ValueFrom{
					ValueFrom: &fk.ValueFromRef{Name: "my-kafka-cluster"},
				},
			}
			err := protovalidate.Validate(spec)
			gomega.Expect(err).To(gomega.BeNil())
		})

		ginkgo.It("rejects spec with missing cluster", func() {
			spec := makeValidSpec()
			spec.Cluster = nil
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects spec with missing topic_name", func() {
			spec := makeValidSpec()
			spec.TopicName = ""
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
		})
	})

	ginkgo.Context("Partitions and replication", func() {
		ginkgo.It("accepts the partition floor and ceiling", func() {
			spec := makeValidSpec()
			spec.PartitionCount = proto.Int32(3)
			gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
			spec.PartitionCount = proto.Int32(2048)
			gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
		})

		ginkgo.It("rejects partition_count below 3", func() {
			spec := makeValidSpec()
			spec.PartitionCount = proto.Int32(2)
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects partition_count above 2048", func() {
			spec := makeValidSpec()
			spec.PartitionCount = proto.Int32(2049)
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
		})

		ginkgo.It("accepts replication_factor of 2", func() {
			spec := makeValidSpec()
			spec.ReplicationFactor = proto.Int32(2)
			err := protovalidate.Validate(spec)
			gomega.Expect(err).To(gomega.BeNil())
		})

		ginkgo.It("rejects replication_factor below 2", func() {
			spec := makeValidSpec()
			spec.ReplicationFactor = proto.Int32(1)
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
		})
	})

	ginkgo.Context("Config block", func() {
		ginkgo.It("accepts a rich valid config", func() {
			spec := makeValidSpec()
			spec.Config = &DigitalOceanDatabaseKafkaTopicConfig{
				CleanupPolicy:                   "compact",
				CompressionType:                 "zstd",
				DeleteRetentionMs:               proto.Uint64(86400000),
				MaxMessageBytes:                 proto.Uint64(1048588),
				MessageDownConversionEnable:     proto.Bool(true),
				MessageFormatVersion:            "3.6",
				MessageTimestampDifferenceMaxMs: proto.Int64(9223372036854775807),
				MessageTimestampType:            "log_append_time",
				MinCleanableDirtyRatio:          proto.Float64(0.5),
				MinInsyncReplicas:               proto.Int32(2),
				Preallocate:                     proto.Bool(false),
				RetentionBytes:                  proto.Int64(-1),
				RetentionMs:                     proto.Int64(-1),
				SegmentBytes:                    proto.Uint64(209715200),
				SegmentMs:                       proto.Uint64(604800000),
			}
			err := protovalidate.Validate(spec)
			gomega.Expect(err).To(gomega.BeNil())
		})

		ginkgo.It("accepts an empty config message", func() {
			spec := makeValidSpec()
			spec.Config = &DigitalOceanDatabaseKafkaTopicConfig{}
			err := protovalidate.Validate(spec)
			gomega.Expect(err).To(gomega.BeNil())
		})

		ginkgo.It("rejects an unknown cleanup_policy", func() {
			spec := makeValidSpec()
			spec.Config = &DigitalOceanDatabaseKafkaTopicConfig{CleanupPolicy: "compact-delete"}
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects a case-mismatched cleanup_policy (values are case-sensitive)", func() {
			spec := makeValidSpec()
			spec.Config = &DigitalOceanDatabaseKafkaTopicConfig{CleanupPolicy: "Delete"}
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects an unknown compression_type", func() {
			spec := makeValidSpec()
			spec.Config = &DigitalOceanDatabaseKafkaTopicConfig{CompressionType: "brotli"}
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
		})

		ginkgo.It("accepts every cleanup policy the provider allows", func() {
			for _, policy := range []string{"delete", "compact", "compact_delete"} {
				spec := makeValidSpec()
				spec.Config = &DigitalOceanDatabaseKafkaTopicConfig{CleanupPolicy: policy}
				gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
			}
		})

		ginkgo.It("rejects an unknown message_format_version", func() {
			spec := makeValidSpec()
			spec.Config = &DigitalOceanDatabaseKafkaTopicConfig{MessageFormatVersion: "9.9"}
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
		})

		ginkgo.It("accepts boundary message_format_version tokens", func() {
			for _, version := range []string{"0.8.0", "3.6-IV2", "0.10.1-IV2"} {
				spec := makeValidSpec()
				spec.Config = &DigitalOceanDatabaseKafkaTopicConfig{MessageFormatVersion: version}
				gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
			}
		})

		ginkgo.It("rejects an unknown message_timestamp_type", func() {
			spec := makeValidSpec()
			spec.Config = &DigitalOceanDatabaseKafkaTopicConfig{MessageTimestampType: "broker_time"}
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects min_cleanable_dirty_ratio above 1", func() {
			spec := makeValidSpec()
			spec.Config = &DigitalOceanDatabaseKafkaTopicConfig{MinCleanableDirtyRatio: proto.Float64(1.5)}
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects min_insync_replicas below 1", func() {
			spec := makeValidSpec()
			spec.Config = &DigitalOceanDatabaseKafkaTopicConfig{MinInsyncReplicas: proto.Int32(0)}
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
		})

		ginkgo.It("accepts -1 as unlimited retention", func() {
			spec := makeValidSpec()
			spec.Config = &DigitalOceanDatabaseKafkaTopicConfig{
				RetentionBytes: proto.Int64(-1),
				RetentionMs:    proto.Int64(-1),
			}
			err := protovalidate.Validate(spec)
			gomega.Expect(err).To(gomega.BeNil())
		})
	})
})
