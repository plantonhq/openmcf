package awss3bucketv1alpha1

import (
	"testing"

	"buf.build/go/protovalidate"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	foreignkeyv1 "github.com/plantonhq/planton/apis/dev/planton/shared/foreignkey/v1"
	"google.golang.org/protobuf/types/known/structpb"
)

func TestAwsS3BucketSpec(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "AwsS3BucketSpec Validation Suite")
}

// helper to create a StringValueOrRef with a literal value.
func strRef(val string) *foreignkeyv1.StringValueOrRef {
	return &foreignkeyv1.StringValueOrRef{
		LiteralOrRef: &foreignkeyv1.StringValueOrRef_Value{Value: val},
	}
}

// helper to build a Struct from a map, failing the spec on error.
func mustStruct(m map[string]interface{}) *structpb.Struct {
	s, err := structpb.NewStruct(m)
	if err != nil {
		panic(err)
	}
	return s
}

var _ = ginkgo.Describe("AwsS3BucketSpec validations", func() {
	var spec *AwsS3BucketSpec

	ginkgo.BeforeEach(func() {
		// Minimal valid spec: a private bucket with all AWS defaults.
		spec = &AwsS3BucketSpec{
			Region: "us-west-2",
		}
	})

	// -------------------------------------------------------------------------
	// Happy path
	// -------------------------------------------------------------------------

	ginkgo.It("accepts a minimal private bucket (all defaults)", func() {
		gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
	})

	ginkgo.It("rejects a missing region", func() {
		spec.Region = ""
		gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
	})

	ginkgo.It("accepts a versioned encrypted bucket with force destroy", func() {
		spec.VersioningStatus = "Enabled"
		spec.ForceDestroy = true
		spec.Encryption = &AwsS3BucketEncryption{
			SseAlgorithm:     "aws:kms",
			KmsKeyId:         strRef("arn:aws:kms:us-west-2:123456789012:key/abc"),
			BucketKeyEnabled: true,
		}
		gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
	})

	// -------------------------------------------------------------------------
	// Versioning
	// -------------------------------------------------------------------------

	ginkgo.It("accepts Enabled and Suspended versioning states", func() {
		for _, s := range []string{"", "Enabled", "Suspended"} {
			spec.VersioningStatus = s
			gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
		}
	})

	ginkgo.It("rejects an invalid versioning state", func() {
		spec.VersioningStatus = "Disabled"
		gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
	})

	// -------------------------------------------------------------------------
	// Encryption
	// -------------------------------------------------------------------------

	ginkgo.It("rejects an invalid sse_algorithm", func() {
		spec.Encryption = &AwsS3BucketEncryption{SseAlgorithm: "kms"}
		gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
	})

	ginkgo.It("rejects a kms key without a KMS algorithm", func() {
		spec.Encryption = &AwsS3BucketEncryption{
			SseAlgorithm: "AES256",
			KmsKeyId:     strRef("arn:aws:kms:us-west-2:123456789012:key/abc"),
		}
		gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
	})

	ginkgo.It("accepts dual-layer SSE-KMS with a key", func() {
		spec.Encryption = &AwsS3BucketEncryption{
			SseAlgorithm: "aws:kms:dsse",
			KmsKeyId:     strRef("arn:aws:kms:us-west-2:123456789012:key/abc"),
		}
		gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
	})

	// -------------------------------------------------------------------------
	// Ownership + ACL coupling
	// -------------------------------------------------------------------------

	ginkgo.It("rejects an acl under (default) BucketOwnerEnforced ownership", func() {
		spec.Acl = "public-read"
		gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
	})

	ginkgo.It("accepts an acl when ownership re-enables ACLs", func() {
		spec.ObjectOwnership = "BucketOwnerPreferred"
		spec.Acl = "log-delivery-write"
		gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
	})

	ginkgo.It("rejects an invalid canned acl", func() {
		spec.ObjectOwnership = "ObjectWriter"
		spec.Acl = "bucket-owner-full-control"
		gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
	})

	ginkgo.It("rejects an invalid object_ownership", func() {
		spec.ObjectOwnership = "Enforced"
		gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
	})

	// -------------------------------------------------------------------------
	// Public access block + policy
	// -------------------------------------------------------------------------

	ginkgo.It("accepts a public website posture (guards relaxed + policy)", func() {
		spec.PublicAccessBlock = &AwsS3BucketPublicAccessBlock{
			BlockPublicAcls:  true,
			IgnorePublicAcls: true,
			// block_public_policy and restrict_public_buckets deliberately false
		}
		spec.Policy = mustStruct(map[string]interface{}{
			"Version": "2012-10-17",
			"Statement": []interface{}{map[string]interface{}{
				"Sid": "PublicReadGetObject", "Effect": "Allow", "Principal": "*",
				"Action": "s3:GetObject", "Resource": "arn:aws:s3:::my-bucket/*",
			}},
		})
		gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
	})

	// -------------------------------------------------------------------------
	// Lifecycle
	// -------------------------------------------------------------------------

	ginkgo.It("accepts a full lifecycle rule", func() {
		spec.VersioningStatus = "Enabled"
		spec.TransitionDefaultMinimumObjectSize = "all_storage_classes_128K"
		spec.LifecycleRules = []*AwsS3BucketLifecycleRule{{
			Id: "tier-logs",
			Filter: &AwsS3BucketLifecycleFilter{
				Prefix:                "logs/",
				Tags:                  map[string]string{"team": "data"},
				ObjectSizeGreaterThan: 1024,
			},
			Transitions: []*AwsS3BucketLifecycleTransition{
				{Days: 30, StorageClass: "STANDARD_IA"},
				{Days: 365, StorageClass: "DEEP_ARCHIVE"},
			},
			Expiration:                   &AwsS3BucketLifecycleExpiration{Days: 730},
			NoncurrentVersionTransitions: []*AwsS3BucketNoncurrentVersionTransition{{NoncurrentDays: 30, StorageClass: "GLACIER_IR"}},
			NoncurrentVersionExpiration:  &AwsS3BucketNoncurrentVersionExpiration{NoncurrentDays: 90, NewerNoncurrentVersions: 3},
		}}
		gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
	})

	ginkgo.It("accepts multipart-upload abort with a prefix-only filter", func() {
		spec.LifecycleRules = []*AwsS3BucketLifecycleRule{{
			Id:                                 "abort-prefixed",
			Filter:                             &AwsS3BucketLifecycleFilter{Prefix: "uploads/"},
			AbortIncompleteMultipartUploadDays: 7,
		}}
		gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
	})

	ginkgo.It("rejects multipart-upload abort combined with a tag filter", func() {
		spec.LifecycleRules = []*AwsS3BucketLifecycleRule{{
			Id:                                 "abort-tagged",
			Filter:                             &AwsS3BucketLifecycleFilter{Tags: map[string]string{"team": "data"}},
			AbortIncompleteMultipartUploadDays: 7,
		}}
		gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
	})

	ginkgo.It("rejects multipart-upload abort combined with an object-size filter", func() {
		spec.LifecycleRules = []*AwsS3BucketLifecycleRule{{
			Id:                                 "abort-sized",
			Filter:                             &AwsS3BucketLifecycleFilter{ObjectSizeGreaterThan: 1024},
			AbortIncompleteMultipartUploadDays: 7,
		}}
		gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
	})

	ginkgo.It("rejects a lifecycle rule with no action", func() {
		spec.LifecycleRules = []*AwsS3BucketLifecycleRule{{Id: "noop"}}
		gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
	})

	ginkgo.It("rejects a lifecycle rule without an id", func() {
		spec.LifecycleRules = []*AwsS3BucketLifecycleRule{{
			Expiration: &AwsS3BucketLifecycleExpiration{Days: 30},
		}}
		gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
	})

	ginkgo.It("rejects an invalid lifecycle rule status", func() {
		spec.LifecycleRules = []*AwsS3BucketLifecycleRule{{
			Id: "r", Status: "Paused",
			Expiration: &AwsS3BucketLifecycleExpiration{Days: 30},
		}}
		gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
	})

	ginkgo.It("rejects a transition with both days and date", func() {
		spec.LifecycleRules = []*AwsS3BucketLifecycleRule{{
			Id: "r",
			Transitions: []*AwsS3BucketLifecycleTransition{
				{Days: 30, Date: "2027-01-01T00:00:00Z", StorageClass: "GLACIER"},
			},
		}}
		gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
	})

	ginkgo.It("accepts a date-based transition", func() {
		spec.LifecycleRules = []*AwsS3BucketLifecycleRule{{
			Id: "r",
			Transitions: []*AwsS3BucketLifecycleTransition{
				{Date: "2027-01-01T00:00:00Z", StorageClass: "GLACIER"},
			},
		}}
		gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
	})

	ginkgo.It("rejects an invalid transition storage class", func() {
		spec.LifecycleRules = []*AwsS3BucketLifecycleRule{{
			Id: "r",
			Transitions: []*AwsS3BucketLifecycleTransition{
				{Days: 30, StorageClass: "STANDARD"},
			},
		}}
		gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
	})

	ginkgo.It("rejects an expiration mixing days with delete-marker cleanup", func() {
		spec.LifecycleRules = []*AwsS3BucketLifecycleRule{{
			Id: "r",
			Expiration: &AwsS3BucketLifecycleExpiration{
				Days: 30, ExpiredObjectDeleteMarker: true,
			},
		}}
		gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
	})

	ginkgo.It("accepts delete-marker-only expiration", func() {
		spec.LifecycleRules = []*AwsS3BucketLifecycleRule{{
			Id:         "r",
			Expiration: &AwsS3BucketLifecycleExpiration{ExpiredObjectDeleteMarker: true},
		}}
		gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
	})

	ginkgo.It("rejects noncurrent version expiration without days", func() {
		spec.LifecycleRules = []*AwsS3BucketLifecycleRule{{
			Id:                          "r",
			NoncurrentVersionExpiration: &AwsS3BucketNoncurrentVersionExpiration{},
		}}
		gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
	})

	ginkgo.It("rejects an invalid transition_default_minimum_object_size", func() {
		spec.TransitionDefaultMinimumObjectSize = "none"
		gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
	})

	// -------------------------------------------------------------------------
	// Replication
	// -------------------------------------------------------------------------

	replicationBase := func() *AwsS3BucketReplication {
		return &AwsS3BucketReplication{
			RoleArn: strRef("arn:aws:iam::123456789012:role/replication"),
			Rules: []*AwsS3BucketReplicationRule{{
				Id: "to-dr",
				Destination: &AwsS3BucketReplicationDestination{
					BucketArn: strRef("arn:aws:s3:::dr-bucket"),
				},
			}},
		}
	}

	ginkgo.It("accepts replication on a versioned bucket", func() {
		spec.VersioningStatus = "Enabled"
		spec.Replication = replicationBase()
		gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
	})

	ginkgo.It("rejects replication without versioning", func() {
		spec.Replication = replicationBase()
		gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
	})

	ginkgo.It("rejects replication without rules", func() {
		spec.VersioningStatus = "Enabled"
		spec.Replication = replicationBase()
		spec.Replication.Rules = nil
		gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
	})

	ginkgo.It("rejects a replication rule without a destination", func() {
		spec.VersioningStatus = "Enabled"
		spec.Replication = replicationBase()
		spec.Replication.Rules[0].Destination = nil
		gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
	})

	ginkgo.It("rejects SSE-KMS replication without a replica key", func() {
		spec.VersioningStatus = "Enabled"
		spec.Replication = replicationBase()
		spec.Replication.Rules[0].ReplicateSseKmsEncryptedObjects = true
		gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
	})

	ginkgo.It("accepts SSE-KMS replication with a replica key", func() {
		spec.VersioningStatus = "Enabled"
		spec.Replication = replicationBase()
		spec.Replication.Rules[0].ReplicateSseKmsEncryptedObjects = true
		spec.Replication.Rules[0].Destination.ReplicaKmsKeyId = strRef("arn:aws:kms:us-east-1:123456789012:key/dr")
		gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
	})

	ginkgo.It("rejects RTC without metrics", func() {
		spec.VersioningStatus = "Enabled"
		spec.Replication = replicationBase()
		spec.Replication.Rules[0].Destination.ReplicationTimeControlEnabled = true
		gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
	})

	ginkgo.It("accepts RTC with metrics", func() {
		spec.VersioningStatus = "Enabled"
		spec.Replication = replicationBase()
		spec.Replication.Rules[0].Destination.MetricsEnabled = true
		spec.Replication.Rules[0].Destination.ReplicationTimeControlEnabled = true
		gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
	})

	ginkgo.It("rejects ownership translation without a destination account", func() {
		spec.VersioningStatus = "Enabled"
		spec.Replication = replicationBase()
		spec.Replication.Rules[0].Destination.ChangeReplicaOwnershipToDestination = true
		gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
	})

	ginkgo.It("accepts cross-account replication with ownership translation", func() {
		spec.VersioningStatus = "Enabled"
		spec.Replication = replicationBase()
		spec.Replication.Rules[0].Destination.Account = "210987654321"
		spec.Replication.Rules[0].Destination.ChangeReplicaOwnershipToDestination = true
		gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
	})

	ginkgo.It("rejects an invalid replica storage class", func() {
		spec.VersioningStatus = "Enabled"
		spec.Replication = replicationBase()
		spec.Replication.Rules[0].Destination.StorageClass = "EXPRESS_ONEZONE"
		gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
	})

	// -------------------------------------------------------------------------
	// Website
	// -------------------------------------------------------------------------

	ginkgo.It("accepts a website with index, error, and routing rules", func() {
		spec.Website = &AwsS3BucketWebsite{
			IndexDocumentSuffix: "index.html",
			ErrorDocumentKey:    "error.html",
			RoutingRules: []*AwsS3BucketWebsiteRoutingRule{{
				Condition: &AwsS3BucketWebsiteRoutingRuleCondition{KeyPrefixEquals: "docs/"},
				Redirect:  &AwsS3BucketWebsiteRoutingRuleRedirect{ReplaceKeyPrefixWith: "documents/"},
			}},
		}
		gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
	})

	ginkgo.It("accepts a redirect-all website", func() {
		spec.Website = &AwsS3BucketWebsite{
			RedirectAllRequestsTo: &AwsS3BucketWebsiteRedirectAll{HostName: "www.example.com", Protocol: "https"},
		}
		gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
	})

	ginkgo.It("rejects a website mixing index and redirect-all", func() {
		spec.Website = &AwsS3BucketWebsite{
			IndexDocumentSuffix:   "index.html",
			RedirectAllRequestsTo: &AwsS3BucketWebsiteRedirectAll{HostName: "www.example.com"},
		}
		gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
	})

	ginkgo.It("rejects an empty website block", func() {
		spec.Website = &AwsS3BucketWebsite{}
		gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
	})

	ginkgo.It("rejects an error document without an index document", func() {
		spec.Website = &AwsS3BucketWebsite{
			ErrorDocumentKey:      "error.html",
			RedirectAllRequestsTo: &AwsS3BucketWebsiteRedirectAll{HostName: "www.example.com"},
		}
		gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
	})

	ginkgo.It("rejects a routing rule without a redirect", func() {
		spec.Website = &AwsS3BucketWebsite{
			IndexDocumentSuffix: "index.html",
			RoutingRules:        []*AwsS3BucketWebsiteRoutingRule{{}},
		}
		gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
	})

	ginkgo.It("rejects an invalid redirect protocol", func() {
		spec.Website = &AwsS3BucketWebsite{
			RedirectAllRequestsTo: &AwsS3BucketWebsiteRedirectAll{HostName: "www.example.com", Protocol: "ftp"},
		}
		gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
	})

	// -------------------------------------------------------------------------
	// Logging
	// -------------------------------------------------------------------------

	ginkgo.It("accepts logging with a partitioned prefix", func() {
		spec.Logging = &AwsS3BucketLogging{
			TargetBucket:                strRef("central-logs"),
			TargetPrefix:                "logs/my-bucket/",
			PartitionedPrefixDateSource: "EventTime",
		}
		gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
	})

	ginkgo.It("rejects logging without a target bucket", func() {
		spec.Logging = &AwsS3BucketLogging{TargetPrefix: "logs/"}
		gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
	})

	ginkgo.It("rejects an invalid partitioned prefix date source", func() {
		spec.Logging = &AwsS3BucketLogging{
			TargetBucket:                strRef("central-logs"),
			PartitionedPrefixDateSource: "RequestTime",
		}
		gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
	})

	// -------------------------------------------------------------------------
	// CORS
	// -------------------------------------------------------------------------

	ginkgo.It("accepts a CORS rule", func() {
		spec.CorsRules = []*AwsS3BucketCorsRule{{
			Id:             "web-app",
			AllowedMethods: []string{"GET", "HEAD"},
			AllowedOrigins: []string{"https://example.com"},
			AllowedHeaders: []string{"*"},
			ExposeHeaders:  []string{"ETag"},
			MaxAgeSeconds:  3600,
		}}
		gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
	})

	ginkgo.It("rejects a CORS rule with an invalid method", func() {
		spec.CorsRules = []*AwsS3BucketCorsRule{{
			AllowedMethods: []string{"PATCH"},
			AllowedOrigins: []string{"*"},
		}}
		gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
	})

	ginkgo.It("rejects a CORS rule without origins", func() {
		spec.CorsRules = []*AwsS3BucketCorsRule{{
			AllowedMethods: []string{"GET"},
		}}
		gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
	})

	// -------------------------------------------------------------------------
	// Notifications
	// -------------------------------------------------------------------------

	ginkgo.It("accepts all four notification arms", func() {
		spec.Notification = &AwsS3BucketNotification{
			Eventbridge: true,
			LambdaFunctions: []*AwsS3BucketLambdaNotification{{
				LambdaFunctionArn: strRef("arn:aws:lambda:us-west-2:123456789012:function:thumbs"),
				Events:            []string{"s3:ObjectCreated:*"},
				FilterSuffix:      ".jpg",
			}},
			Queues: []*AwsS3BucketQueueNotification{{
				QueueArn: strRef("arn:aws:sqs:us-west-2:123456789012:ingest"),
				Events:   []string{"s3:ObjectCreated:Put"},
			}},
			Topics: []*AwsS3BucketTopicNotification{{
				TopicArn: strRef("arn:aws:sns:us-west-2:123456789012:alerts"),
				Events:   []string{"s3:ObjectRemoved:*"},
			}},
		}
		gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
	})

	ginkgo.It("rejects a queue notification without events", func() {
		spec.Notification = &AwsS3BucketNotification{
			Queues: []*AwsS3BucketQueueNotification{{
				QueueArn: strRef("arn:aws:sqs:us-west-2:123456789012:ingest"),
			}},
		}
		gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
	})

	ginkgo.It("rejects a lambda notification without a function", func() {
		spec.Notification = &AwsS3BucketNotification{
			LambdaFunctions: []*AwsS3BucketLambdaNotification{{
				Events: []string{"s3:ObjectCreated:*"},
			}},
		}
		gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
	})

	// -------------------------------------------------------------------------
	// Object Lock
	// -------------------------------------------------------------------------

	ginkgo.It("accepts an object-lock bucket with default retention", func() {
		spec.VersioningStatus = "Enabled"
		spec.ObjectLockEnabled = true
		spec.ObjectLockDefaultRetention = &AwsS3BucketObjectLockDefaultRetention{
			Mode: "GOVERNANCE", Days: 30,
		}
		gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
	})

	ginkgo.It("rejects object lock without versioning", func() {
		spec.ObjectLockEnabled = true
		gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
	})

	ginkgo.It("rejects default retention without object lock", func() {
		spec.VersioningStatus = "Enabled"
		spec.ObjectLockDefaultRetention = &AwsS3BucketObjectLockDefaultRetention{
			Mode: "COMPLIANCE", Years: 1,
		}
		gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
	})

	ginkgo.It("rejects retention with both days and years", func() {
		spec.VersioningStatus = "Enabled"
		spec.ObjectLockEnabled = true
		spec.ObjectLockDefaultRetention = &AwsS3BucketObjectLockDefaultRetention{
			Mode: "GOVERNANCE", Days: 30, Years: 1,
		}
		gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
	})

	ginkgo.It("rejects retention with an invalid mode", func() {
		spec.VersioningStatus = "Enabled"
		spec.ObjectLockEnabled = true
		spec.ObjectLockDefaultRetention = &AwsS3BucketObjectLockDefaultRetention{
			Mode: "LEGAL_HOLD", Days: 30,
		}
		gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
	})

	// -------------------------------------------------------------------------
	// Acceleration + requester pays
	// -------------------------------------------------------------------------

	ginkgo.It("accepts acceleration and requester pays", func() {
		spec.AccelerationStatus = "Enabled"
		spec.RequestPayer = "Requester"
		gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
	})

	ginkgo.It("rejects an invalid acceleration status", func() {
		spec.AccelerationStatus = "On"
		gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
	})

	ginkgo.It("rejects an invalid request payer", func() {
		spec.RequestPayer = "Caller"
		gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
	})

	// -------------------------------------------------------------------------
	// Intelligent-Tiering
	// -------------------------------------------------------------------------

	ginkgo.It("accepts an intelligent-tiering configuration with both tiers", func() {
		spec.IntelligentTieringConfigurations = []*AwsS3BucketIntelligentTieringConfiguration{{
			Name:         "archive-old",
			FilterPrefix: "data/",
			Tiers: []*AwsS3BucketIntelligentTieringTier{
				{AccessTier: "ARCHIVE_ACCESS", Days: 90},
				{AccessTier: "DEEP_ARCHIVE_ACCESS", Days: 180},
			},
		}}
		gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
	})

	ginkgo.It("rejects an intelligent-tiering configuration without tiers", func() {
		spec.IntelligentTieringConfigurations = []*AwsS3BucketIntelligentTieringConfiguration{{
			Name: "empty",
		}}
		gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
	})

	ginkgo.It("rejects an archive tier below the 90-day minimum", func() {
		spec.IntelligentTieringConfigurations = []*AwsS3BucketIntelligentTieringConfiguration{{
			Name:  "too-soon",
			Tiers: []*AwsS3BucketIntelligentTieringTier{{AccessTier: "ARCHIVE_ACCESS", Days: 60}},
		}}
		gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
	})

	ginkgo.It("rejects a deep-archive tier below the 180-day minimum", func() {
		spec.IntelligentTieringConfigurations = []*AwsS3BucketIntelligentTieringConfiguration{{
			Name:  "too-soon",
			Tiers: []*AwsS3BucketIntelligentTieringTier{{AccessTier: "DEEP_ARCHIVE_ACCESS", Days: 120}},
		}}
		gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
	})

	ginkgo.It("rejects an invalid access tier", func() {
		spec.IntelligentTieringConfigurations = []*AwsS3BucketIntelligentTieringConfiguration{{
			Name:  "bad-tier",
			Tiers: []*AwsS3BucketIntelligentTieringTier{{AccessTier: "GLACIER", Days: 90}},
		}}
		gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
	})
})
