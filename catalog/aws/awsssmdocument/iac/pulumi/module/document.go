package module

import (
	"strings"

	"github.com/pkg/errors"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws/ssm"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// document creates the SSM document and exports outputs.
//
// Lifecycle facts the render below depends on:
//   - the document name is metadata.name on both engines and changing
//     it forces replacement;
//   - updating the content creates a NEW document version and the
//     provider promotes it to the default version; schema-1.x documents
//     only update when the content itself changes (an AWS rule);
//   - Permissions is the provider's flat share map - the spec's
//     share_with_account_ids renders as {type: "Share", account_ids:
//     "<comma-joined>"} (AWS applies changes in batches of 20);
//   - attachment metadata is never read back by any SSM API - the
//     import map declares attachment_sources config-only.
func document(ctx *pulumi.Context, locals *Locals, provider *aws.Provider) error {
	spec := locals.Spec

	args := &ssm.DocumentArgs{
		Name:         pulumi.String(locals.DocumentName),
		Content:      pulumi.String(spec.Content),
		DocumentType: pulumi.String(spec.DocumentType),
		Tags:         pulumi.ToStringMap(locals.AwsTags),
	}

	if spec.DocumentFormat != "" {
		args.DocumentFormat = pulumi.String(spec.DocumentFormat)
	}
	if spec.TargetType != "" {
		args.TargetType = pulumi.String(spec.TargetType)
	}
	if spec.VersionName != "" {
		args.VersionName = pulumi.String(spec.VersionName)
	}

	var attachments ssm.DocumentAttachmentsSourceArray
	for _, a := range spec.AttachmentSources {
		attachmentArgs := &ssm.DocumentAttachmentsSourceArgs{
			Key:    pulumi.String(a.Key),
			Values: pulumi.ToStringArray(a.Values),
		}
		if a.Name != "" {
			attachmentArgs.Name = pulumi.String(a.Name)
		}
		attachments = append(attachments, attachmentArgs)
	}
	if len(attachments) > 0 {
		args.AttachmentsSources = attachments
	}

	if len(spec.ShareWithAccountIds) > 0 {
		args.Permissions = pulumi.StringMap{
			"type":        pulumi.String("Share"),
			"account_ids": pulumi.String(strings.Join(spec.ShareWithAccountIds, ",")),
		}
	}

	createdDocument, err := ssm.NewDocument(ctx, "document", args, pulumi.Provider(provider))
	if err != nil {
		return errors.Wrap(err, "create document")
	}

	ctx.Export(OpDocumentName, createdDocument.Name)
	ctx.Export(OpDocumentArn, createdDocument.Arn)
	ctx.Export(OpDefaultVersion, createdDocument.DefaultVersion)
	ctx.Export(OpLatestVersion, createdDocument.LatestVersion)
	ctx.Export(OpDocumentHash, createdDocument.Hash)
	ctx.Export(OpStatus, createdDocument.Status)
	return nil
}
