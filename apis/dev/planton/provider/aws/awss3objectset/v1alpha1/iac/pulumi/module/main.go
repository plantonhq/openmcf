package module

import (
	"fmt"

	"github.com/pkg/errors"
	awss3objectsetv1alpha1 "github.com/plantonhq/planton/apis/dev/planton/provider/aws/awss3objectset/v1alpha1"
	"github.com/plantonhq/planton/pkg/iac/pulumi/pulumimodule/provider/aws/pulumiawsprovider"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws/s3"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func Resources(ctx *pulumi.Context, stackInput *awss3objectsetv1alpha1.AwsS3ObjectSetStackInput) error {
	locals := initializeLocals(ctx, stackInput)
	spec := locals.AwsS3ObjectSet.Spec

	// Build the AWS provider from the stack input via the shared builder, which resolves
	// the right credential mechanism (static keys, keyless web identity, or ambient chain).
	provider, err := pulumiawsprovider.Get(ctx, stackInput.ProviderConfig, spec.Region)
	if err != nil {
		return errors.Wrap(err, "failed to create AWS provider")
	}

	// The bucket foreign key is resolved to a literal name before the module
	// runs (a referenced AwsS3Bucket resolves to status.outputs.bucket_id).
	bucketName := spec.Bucket.GetValue()
	if bucketName == "" {
		return errors.New("bucket name must be resolved before invoking IaC module")
	}

	arnMap := pulumi.StringMap{}
	etagMap := pulumi.StringMap{}
	versionIdMap := pulumi.StringMap{}

	for _, obj := range spec.Objects {
		if obj.Key == "" {
			return errors.New("every object must have a non-empty key")
		}

		// Merge order: identity tags < set-level tags < object-level tags, so
		// an object can specialize but never lose its resource-identity
		// attribution — the same precedence the Terraform module applies.
		objTags := pulumi.StringMap{}
		for k, v := range locals.Labels {
			objTags[k] = pulumi.String(v)
		}
		for k, v := range spec.Tags {
			objTags[k] = pulumi.String(v)
		}
		for k, v := range obj.Tags {
			objTags[k] = pulumi.String(v)
		}

		args := &s3.BucketObjectv2Args{
			Bucket: pulumi.String(bucketName),
			Key:    pulumi.String(obj.Key),
			Tags:   objTags,
		}

		// Content source: the spec guarantees exactly one arm of the oneof is
		// set (proto oneof + CEL).
		switch source := obj.Source.(type) {
		case *awss3objectsetv1alpha1.AwsS3Object_Content:
			args.Content = pulumi.StringPtr(source.Content)
		case *awss3objectsetv1alpha1.AwsS3Object_ContentBase64:
			args.ContentBase64 = pulumi.StringPtr(source.ContentBase64)
		default:
			return fmt.Errorf("object %q has no content source", obj.Key)
		}

		// Sent explicitly, resolving the spec default (application/octet-stream)
		// here rather than leaving it to the provider: an omitted Content-Type
		// would make S3 store its own default (binary/octet-stream), diverging
		// from what the manifest declares and from the Terraform engine.
		args.ContentType = pulumi.StringPtr(obj.GetContentType())

		// HTTP presentation headers — sent only when set so unset stays
		// indistinguishable from the AWS defaults (no phantom headers stored).
		if obj.CacheControl != "" {
			args.CacheControl = pulumi.StringPtr(obj.CacheControl)
		}
		if obj.ContentEncoding != "" {
			args.ContentEncoding = pulumi.StringPtr(obj.ContentEncoding)
		}
		if obj.ContentDisposition != "" {
			args.ContentDisposition = pulumi.StringPtr(obj.ContentDisposition)
		}
		if obj.ContentLanguage != "" {
			args.ContentLanguage = pulumi.StringPtr(obj.ContentLanguage)
		}

		// User metadata (x-amz-meta-*). Keys are spec-validated lowercase,
		// matching what S3 stores, so reads never drift from the manifest.
		if len(obj.Metadata) > 0 {
			metadataMap := pulumi.StringMap{}
			for k, v := range obj.Metadata {
				metadataMap[k] = pulumi.String(v)
			}
			args.Metadata = metadataMap
		}

		// Website redirect: inert unless the bucket has static website hosting.
		if obj.WebsiteRedirect != "" {
			args.WebsiteRedirect = pulumi.StringPtr(obj.WebsiteRedirect)
		}

		// Storage placement. Unset means STANDARD (AWS's default) — omit the
		// argument so the provider computes rather than pins a value.
		if obj.StorageClass != "" {
			args.StorageClass = pulumi.StringPtr(obj.StorageClass)
		}

		// Per-object encryption OVERRIDE. Unset inherits the bucket's default
		// encryption, which is where uniform posture belongs. A kms_key alone
		// implies aws:kms (the provider sets ServerSideEncryption when a KMS
		// key is sent); the spec CEL rejects the contradictory kms_key +
		// AES256 pair.
		if obj.ServerSideEncryption != "" {
			args.ServerSideEncryption = pulumi.StringPtr(obj.ServerSideEncryption)
		}
		if kmsKeyArn := obj.KmsKey.GetValue(); kmsKeyArn != "" {
			args.KmsKeyId = pulumi.StringPtr(kmsKeyArn)
		}
		if obj.BucketKeyEnabled != nil {
			args.BucketKeyEnabled = pulumi.BoolPtr(obj.GetBucketKeyEnabled())
		}

		// Upload-integrity checksum, stored alongside the object.
		if obj.ChecksumAlgorithm != "" {
			args.ChecksumAlgorithm = pulumi.StringPtr(obj.ChecksumAlgorithm)
		}

		// Object Lock (requires an Object Lock-enabled bucket; the spec CEL
		// guarantees mode and retain-until arrive as a pair).
		if obj.ObjectLockMode != "" {
			args.ObjectLockMode = pulumi.StringPtr(obj.ObjectLockMode)
		}
		if obj.ObjectLockRetainUntilDate != "" {
			args.ObjectLockRetainUntilDate = pulumi.StringPtr(obj.ObjectLockRetainUntilDate)
		}
		if obj.ObjectLockLegalHoldStatus != "" {
			args.ObjectLockLegalHoldStatus = pulumi.StringPtr(obj.ObjectLockLegalHoldStatus)
		}

		// Canned ACL — only valid on buckets whose ownership controls permit
		// ACLs; modern (BucketOwnerEnforced) buckets reject it at apply time.
		if obj.Acl != "" {
			args.Acl = pulumi.StringPtr(obj.Acl)
		}

		// The GOVERNANCE-retention delete bypass
		// (x-amz-bypass-governance-retention). Only valid on Object
		// Lock-enabled buckets — S3 rejects the flag on regular buckets,
		// failing the destroy. Versioned-bucket cleanup needs no flag:
		// deleting an object always removes all of its versions.
		if obj.ForceDestroy {
			args.ForceDestroy = pulumi.BoolPtr(true)
		}

		// The resource is named by the S3 object key, so adding, removing, or
		// REORDERING entries in the manifest never churns unrelated objects —
		// each object's Pulumi identity is its key, mirroring both its S3
		// identity and the Terraform module's for_each keying.
		s3Object, err := s3.NewBucketObjectv2(ctx, obj.Key, args, pulumi.Provider(provider))
		if err != nil {
			return errors.Wrapf(err, "failed to create S3 object %q", obj.Key)
		}

		arnMap[obj.Key] = s3Object.Arn
		etagMap[obj.Key] = s3Object.Etag
		versionIdMap[obj.Key] = s3Object.VersionId
	}

	ctx.Export(OpBucketID, pulumi.String(bucketName))
	ctx.Export(OpObjectArns, arnMap)
	ctx.Export(OpObjectEtags, etagMap)
	ctx.Export(OpObjectVersionIds, versionIdMap)

	return nil
}
