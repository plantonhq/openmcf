package module

import (
	"github.com/pkg/errors"
	"github.com/pulumi/pulumi-cloudflare/sdk/v6/go/cloudflare"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// logpushJob creates the Logpush job: continuous delivery of one Cloudflare
// log dataset to a destination the account controls. Dual scope -- exactly
// one of account_id or zone_id is set (spec validation enforces it).
//
// Two API truths this module honors: `dataset` is immutable (the provider
// replaces the job when it changes), and `destination_conf` is not marked
// for replacement by the provider yet Cloudflare rejects changing it on an
// existing job -- repointing means delete and recreate.
func logpushJob(
	ctx *pulumi.Context,
	locals *Locals,
	cloudflareProvider *cloudflare.Provider,
) error {
	spec := locals.CloudflareLogpushJob.Spec

	args := &cloudflare.LogpushJobArgs{
		Dataset:         pulumi.String(spec.Dataset),
		DestinationConf: pulumi.String(spec.DestinationConf.GetValue()),
	}

	if spec.AccountId != "" {
		args.AccountId = pulumi.StringPtr(spec.AccountId)
	}
	if spec.ZoneId.GetValue() != "" {
		args.ZoneId = pulumi.StringPtr(spec.ZoneId.GetValue())
	}

	if spec.Name != "" {
		args.Name = pulumi.StringPtr(spec.Name)
	}

	// enabled defaults to TRUE here even though Cloudflare's own default is
	// FALSE: a declared log job is meant to ship logs. Set enabled = false
	// explicitly to pause delivery.
	if spec.Enabled != nil {
		args.Enabled = pulumi.BoolPtr(spec.GetEnabled())
	} else {
		args.Enabled = pulumi.BoolPtr(true)
	}

	if spec.Filter != "" {
		args.Filter = pulumi.StringPtr(spec.Filter)
	}
	if spec.Kind != "" {
		args.Kind = pulumi.StringPtr(spec.Kind)
	}
	if spec.MaxUploadBytes != nil {
		args.MaxUploadBytes = pulumi.IntPtr(int(spec.GetMaxUploadBytes()))
	}
	if spec.MaxUploadIntervalSeconds != nil {
		args.MaxUploadIntervalSeconds = pulumi.IntPtr(int(spec.GetMaxUploadIntervalSeconds()))
	}
	if spec.MaxUploadRecords != nil {
		args.MaxUploadRecords = pulumi.IntPtr(int(spec.GetMaxUploadRecords()))
	}
	if spec.OwnershipChallenge.GetValue() != "" {
		args.OwnershipChallenge = pulumi.StringPtr(spec.OwnershipChallenge.GetValue())
	}

	if spec.OutputOptions != nil {
		options := cloudflare.LogpushJobOutputOptionsArgs{}
		if spec.OutputOptions.OutputType != "" {
			options.OutputType = pulumi.StringPtr(spec.OutputOptions.OutputType)
		}
		if len(spec.OutputOptions.FieldNames) > 0 {
			options.FieldNames = pulumi.ToStringArray(spec.OutputOptions.FieldNames)
		}
		if spec.OutputOptions.TimestampFormat != "" {
			options.TimestampFormat = pulumi.StringPtr(spec.OutputOptions.TimestampFormat)
		}
		if spec.OutputOptions.SampleRate != nil {
			options.SampleRate = pulumi.Float64Ptr(spec.OutputOptions.GetSampleRate())
		}
		if spec.OutputOptions.BatchPrefix != "" {
			options.BatchPrefix = pulumi.StringPtr(spec.OutputOptions.BatchPrefix)
		}
		if spec.OutputOptions.BatchSuffix != "" {
			options.BatchSuffix = pulumi.StringPtr(spec.OutputOptions.BatchSuffix)
		}
		if spec.OutputOptions.RecordPrefix != "" {
			options.RecordPrefix = pulumi.StringPtr(spec.OutputOptions.RecordPrefix)
		}
		if spec.OutputOptions.RecordSuffix != "" {
			options.RecordSuffix = pulumi.StringPtr(spec.OutputOptions.RecordSuffix)
		}
		if spec.OutputOptions.RecordDelimiter != "" {
			options.RecordDelimiter = pulumi.StringPtr(spec.OutputOptions.RecordDelimiter)
		}
		if spec.OutputOptions.RecordTemplate != "" {
			options.RecordTemplate = pulumi.StringPtr(spec.OutputOptions.RecordTemplate)
		}
		if spec.OutputOptions.FieldDelimiter != "" {
			options.FieldDelimiter = pulumi.StringPtr(spec.OutputOptions.FieldDelimiter)
		}
		if spec.OutputOptions.MergeSubrequests != nil {
			options.MergeSubrequests = pulumi.BoolPtr(spec.OutputOptions.GetMergeSubrequests())
		}
		if spec.OutputOptions.Cve_2021_44228 != nil {
			options.Cve202144228 = pulumi.BoolPtr(spec.OutputOptions.GetCve_2021_44228())
		}
		args.OutputOptions = options
	}

	createdJob, err := cloudflare.NewLogpushJob(
		ctx,
		"logpush_job",
		args,
		pulumi.Provider(cloudflareProvider),
	)
	if err != nil {
		return errors.Wrap(err, "failed to create logpush job")
	}

	ctx.Export(OpJobId, createdJob.ID())
	ctx.Export(OpAccountId, pulumi.String(spec.AccountId))
	ctx.Export(OpZoneId, pulumi.String(spec.ZoneId.GetValue()))

	// The ownership-challenge issuing step, deployed only when the spec asks
	// for it. ONE-SHOT at Cloudflare: the POST drops a challenge file into
	// the destination and that is the whole lifecycle -- no read, no update,
	// no delete (destroying it only forgets it), no import. The token inside
	// the dropped file is fetched by the operator and fed back as
	// spec.ownership_challenge: Cloudflare deliberately routes the proof
	// through storage you control, so no API can read it for you.
	if !spec.GenerateOwnershipChallenge {
		ctx.Export(OpOwnershipChallengeFilename, pulumi.String(""))
		ctx.Export(OpOwnershipChallengeMessage, pulumi.String(""))
		ctx.Export(OpOwnershipChallengeValid, pulumi.Bool(false))
		return nil
	}

	challengeArgs := &cloudflare.LogpushOwnershipChallengeArgs{
		DestinationConf: pulumi.String(spec.DestinationConf.GetValue()),
	}
	if spec.AccountId != "" {
		challengeArgs.AccountId = pulumi.StringPtr(spec.AccountId)
	}
	if spec.ZoneId.GetValue() != "" {
		challengeArgs.ZoneId = pulumi.StringPtr(spec.ZoneId.GetValue())
	}

	createdChallenge, err := cloudflare.NewLogpushOwnershipChallenge(
		ctx,
		"logpush_ownership_challenge",
		challengeArgs,
		pulumi.Provider(cloudflareProvider),
	)
	if err != nil {
		return errors.Wrap(err, "failed to issue logpush ownership challenge")
	}

	ctx.Export(OpOwnershipChallengeFilename, createdChallenge.Filename)
	ctx.Export(OpOwnershipChallengeMessage, createdChallenge.Message)
	ctx.Export(OpOwnershipChallengeValid, createdChallenge.Valid)

	return nil
}
