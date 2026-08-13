package module

import (
	"github.com/pkg/errors"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws/s3"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// lifecycle creates the lifecycle configuration satellite when the spec
// defines rules.
func lifecycle(ctx *pulumi.Context, locals *Locals, provider *aws.Provider,
	createdBucket *s3.BucketV2, versioning *s3.BucketVersioningV2) error {
	spec := locals.Spec
	if len(spec.LifecycleRules) == 0 {
		return nil
	}

	rules := s3.BucketLifecycleConfigurationV2RuleArray{}
	for _, r := range spec.LifecycleRules {
		status := r.Status
		if status == "" {
			status = "Enabled"
		}

		rule := &s3.BucketLifecycleConfigurationV2RuleArgs{
			Id:     pulumi.String(r.Id),
			Status: pulumi.String(status),
		}

		// AWS expresses a single-predicate filter directly on the filter
		// block and multi-predicate filters inside an `and` wrapper. The spec
		// exposes one flat filter message and this shaping emits whichever
		// document shape the predicate count requires.
		if r.Filter != nil {
			predicates := 0
			if r.Filter.Prefix != "" {
				predicates++
			}
			predicates += len(r.Filter.Tags)
			if r.Filter.ObjectSizeGreaterThan > 0 {
				predicates++
			}
			if r.Filter.ObjectSizeLessThan > 0 {
				predicates++
			}

			filter := &s3.BucketLifecycleConfigurationV2RuleFilterArgs{}
			switch {
			case predicates > 1:
				and := &s3.BucketLifecycleConfigurationV2RuleFilterAndArgs{}
				if r.Filter.Prefix != "" {
					and.Prefix = pulumi.StringPtr(r.Filter.Prefix)
				}
				if len(r.Filter.Tags) > 0 {
					and.Tags = pulumi.ToStringMap(r.Filter.Tags)
				}
				if r.Filter.ObjectSizeGreaterThan > 0 {
					and.ObjectSizeGreaterThan = pulumi.IntPtr(int(r.Filter.ObjectSizeGreaterThan))
				}
				if r.Filter.ObjectSizeLessThan > 0 {
					and.ObjectSizeLessThan = pulumi.IntPtr(int(r.Filter.ObjectSizeLessThan))
				}
				filter.And = and
			case r.Filter.Prefix != "":
				filter.Prefix = pulumi.StringPtr(r.Filter.Prefix)
			case len(r.Filter.Tags) == 1:
				for k, v := range r.Filter.Tags {
					filter.Tag = &s3.BucketLifecycleConfigurationV2RuleFilterTagArgs{
						Key:   pulumi.String(k),
						Value: pulumi.String(v),
					}
				}
			case r.Filter.ObjectSizeGreaterThan > 0:
				filter.ObjectSizeGreaterThan = pulumi.IntPtr(int(r.Filter.ObjectSizeGreaterThan))
			case r.Filter.ObjectSizeLessThan > 0:
				filter.ObjectSizeLessThan = pulumi.IntPtr(int(r.Filter.ObjectSizeLessThan))
			}
			rule.Filter = filter
		}

		if len(r.Transitions) > 0 {
			transitions := s3.BucketLifecycleConfigurationV2RuleTransitionArray{}
			for _, t := range r.Transitions {
				// days is presence-typed in the contract so an explicit 0 —
				// AWS's "transition on the upload day" — passes through; CEL
				// enforces exactly-one of days/date.
				transition := &s3.BucketLifecycleConfigurationV2RuleTransitionArgs{
					StorageClass: pulumi.String(t.StorageClass),
				}
				if t.Days != nil {
					transition.Days = pulumi.IntPtr(int(*t.Days))
				}
				if t.Date != "" {
					transition.Date = pulumi.StringPtr(t.Date)
				}
				transitions = append(transitions, transition)
			}
			rule.Transitions = transitions
		}

		if r.Expiration != nil {
			expiration := &s3.BucketLifecycleConfigurationV2RuleExpirationArgs{}
			if r.Expiration.Days > 0 {
				expiration.Days = pulumi.IntPtr(int(r.Expiration.Days))
			}
			if r.Expiration.Date != "" {
				expiration.Date = pulumi.StringPtr(r.Expiration.Date)
			}
			if r.Expiration.ExpiredObjectDeleteMarker {
				expiration.ExpiredObjectDeleteMarker = pulumi.BoolPtr(true)
			}
			rule.Expiration = expiration
		}

		if len(r.NoncurrentVersionTransitions) > 0 {
			nvts := s3.BucketLifecycleConfigurationV2RuleNoncurrentVersionTransitionArray{}
			for _, t := range r.NoncurrentVersionTransitions {
				nvt := &s3.BucketLifecycleConfigurationV2RuleNoncurrentVersionTransitionArgs{
					NoncurrentDays: pulumi.Int(int(t.NoncurrentDays)),
					StorageClass:   pulumi.String(t.StorageClass),
				}
				if t.NewerNoncurrentVersions > 0 {
					nvt.NewerNoncurrentVersions = pulumi.IntPtr(int(t.NewerNoncurrentVersions))
				}
				nvts = append(nvts, nvt)
			}
			rule.NoncurrentVersionTransitions = nvts
		}

		if r.NoncurrentVersionExpiration != nil {
			nve := &s3.BucketLifecycleConfigurationV2RuleNoncurrentVersionExpirationArgs{
				NoncurrentDays: pulumi.Int(int(r.NoncurrentVersionExpiration.NoncurrentDays)),
			}
			if r.NoncurrentVersionExpiration.NewerNoncurrentVersions > 0 {
				nve.NewerNoncurrentVersions = pulumi.IntPtr(int(r.NoncurrentVersionExpiration.NewerNoncurrentVersions))
			}
			rule.NoncurrentVersionExpiration = nve
		}

		if r.AbortIncompleteMultipartUploadDays > 0 {
			rule.AbortIncompleteMultipartUpload = &s3.BucketLifecycleConfigurationV2RuleAbortIncompleteMultipartUploadArgs{
				DaysAfterInitiation: pulumi.IntPtr(int(r.AbortIncompleteMultipartUploadDays)),
			}
		}

		rules = append(rules, rule)
	}

	args := &s3.BucketLifecycleConfigurationV2Args{
		Bucket: createdBucket.ID(),
		Rules:  rules,
	}
	if spec.TransitionDefaultMinimumObjectSize != "" {
		args.TransitionDefaultMinimumObjectSize = pulumi.StringPtr(spec.TransitionDefaultMinimumObjectSize)
	}

	// Lifecycle actions on noncurrent versions only make sense once
	// versioning exists; the explicit edge avoids a transient AWS validation
	// window when both land in one deploy.
	opts := []pulumi.ResourceOption{pulumi.Provider(provider)}
	if versioning != nil {
		opts = append(opts, pulumi.DependsOn([]pulumi.Resource{versioning}))
	}
	if _, err := s3.NewBucketLifecycleConfigurationV2(ctx, "lifecycle", args, opts...); err != nil {
		return errors.Wrap(err, "failed to configure lifecycle rules")
	}
	return nil
}
